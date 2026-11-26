package shareholders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound indicates no shareholder data exists for the requested ticker.
var ErrNotFound = errors.New("shareholder ticker not found")

var slugRe = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// Service implements read-only shareholders and corporate ownership matrix queries.
type Service struct {
	DB *pgxpool.Pool
}

// ClassifyShareholder replicates shareholders._classify_shareholder.
func ClassifyShareholder(name string, pct float64, isControlling bool) string {
	clean := strings.ToLower(strings.TrimSpace(name))
	if isControlling {
		return "controlling"
	}
	if strings.Contains(clean, "treasury") {
		return "treasury"
	}
	if strings.HasPrefix(clean, "masyarakat") || strings.HasPrefix(clean, "publik") {
		return "public"
	}
	if pct >= 5.0 {
		return "substantial"
	}
	return "management"
}

func parseHoldingsJSON(raw []byte) []HoldingItem {
	if len(raw) == 0 {
		return []HoldingItem{}
	}
	var items []HoldingItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return []HoldingItem{}
	}
	for i := range items {
		items[i].Pct = roundEven(items[i].Pct, 3)
	}
	return items
}

// GetTopCrossHoldings returns multi-issuer entities ordered by ticker_count DESC, total_pct DESC, name ASC.
func (s *Service) GetTopCrossHoldings(ctx context.Context, limit, minHoldings int) ([]EntityHoldingRow, error) {
	rows, err := s.DB.Query(ctx, `
		SELECT name,
		       COUNT(DISTINCT ticker)::int AS ticker_count,
		       COALESCE(SUM(pct), 0)::float8 AS total_pct,
		       json_agg(
		           json_build_object(
		               'ticker', ticker,
		               'pct', pct,
		               'is_controlling', is_controlling
		           ) ORDER BY pct DESC
		       )::text AS holdings
		FROM company_shareholders
		WHERE name NOT ILIKE 'Masyarakat%'
		  AND name NOT ILIKE 'Saham Treasury%'
		  AND name NOT ILIKE 'Publik%'
		  AND name NOT IN ('-', 'Pengendali Saham', 'Non Pengendali Saham', 'Others', 'Lain-lain')
		GROUP BY name
		HAVING COUNT(DISTINCT ticker) >= $1
		ORDER BY ticker_count DESC, total_pct DESC, name ASC
		LIMIT $2
	`, minHoldings, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []EntityHoldingRow{}
	for rows.Next() {
		var name, holdingsJSON string
		var cnt int
		var totPct float64
		if scanErr := rows.Scan(&name, &cnt, &totPct, &holdingsJSON); scanErr != nil {
			return nil, scanErr
		}
		out = append(out, EntityHoldingRow{
			Name:        name,
			TickerCount: cnt,
			TotalPct:    roundEven(totPct, 2),
			Holdings:    parseHoldingsJSON([]byte(holdingsJSON)),
		})
	}
	return out, rows.Err()
}

// SearchShareholders searches shareholders by entity name substring.
func (s *Service) SearchShareholders(ctx context.Context, q string, limit int) (SearchResponse, error) {
	searchTerm := "%" + strings.TrimSpace(q) + "%"
	rows, err := s.DB.Query(ctx, `
		SELECT name,
		       COUNT(DISTINCT ticker)::int AS ticker_count,
		       COALESCE(SUM(pct), 0)::float8 AS total_pct,
		       json_agg(
		           json_build_object(
		               'ticker', ticker,
		               'pct', pct,
		               'is_controlling', is_controlling
		           ) ORDER BY pct DESC
		       )::text AS holdings
		FROM company_shareholders
		WHERE name ILIKE $1
		GROUP BY name
		ORDER BY ticker_count DESC, total_pct DESC, name ASC
		LIMIT $2
	`, searchTerm, limit)
	if err != nil {
		return SearchResponse{}, err
	}
	defer rows.Close()

	results := []EntityHoldingRow{}
	for rows.Next() {
		var name, holdingsJSON string
		var cnt int
		var totPct float64
		if scanErr := rows.Scan(&name, &cnt, &totPct, &holdingsJSON); scanErr != nil {
			return SearchResponse{}, scanErr
		}
		results = append(results, EntityHoldingRow{
			Name:        name,
			TickerCount: cnt,
			TotalPct:    roundEven(totPct, 2),
			Holdings:    parseHoldingsJSON([]byte(holdingsJSON)),
		})
	}
	if err := rows.Err(); err != nil {
		return SearchResponse{}, err
	}
	return SearchResponse{
		Query:   q,
		Count:   len(results),
		Results: results,
	}, nil
}

func (s *Service) orderedGroups(ctx context.Context) []shared.ConglomerateGroup {
	ordered, _ := shared.LoadConglomerateGroups(ctx, s.DB)
	if len(ordered) > 0 {
		return ordered
	}
	return staticConglomerateGroupsOrdered
}

// GetConglomerateGroups returns mapped conglomerate groups preserving SQL/static order.
func (s *Service) GetConglomerateGroups(ctx context.Context) GroupsResponse {
	src := s.orderedGroups(ctx)
	groupsList := make([]GroupListItem, 0, len(src))
	for _, g := range src {
		sample := g.KnownTickers
		if sample == nil {
			sample = []string{}
		}
		srcName := g.Source
		if srcName == "" {
			srcName = "conglomerate_static"
		}
		groupsList = append(groupsList, GroupListItem{
			ID:            g.ID,
			Name:          g.Name,
			Description:   g.Description,
			HubEntity:     g.HubEntity,
			LogoURL:       g.LogoURL,
			CatalogCode:   g.CatalogCode,
			SampleTickers: sample,
			TickerCount:   len(sample),
			Source:        srcName,
		})
	}
	return GroupsResponse{Groups: groupsList}
}

type graphHolderRow struct {
	ticker        string
	ownerName     string
	pct           float64
	isControlling bool
}

// GetOwnershipNetworkGraph builds the conglomerate/ticker cross-ownership node/link graph.
func (s *Service) GetOwnershipNetworkGraph(ctx context.Context, group, ticker string, minPct float64) (NetworkGraphResponse, error) {
	targetTicker := strings.ToUpper(strings.TrimSpace(ticker))
	groupKey := strings.ToLower(strings.TrimSpace(group))
	if groupKey == "" {
		groupKey = "barito"
	}

	var groupData shared.ConglomerateGroup
	if resolved, ok := shared.GetConglomerateGroup(ctx, s.DB, groupKey); ok && resolved != nil {
		groupData = *resolved
	} else {
		groupData = staticConglomerateGroupsOrdered[0]
		for _, sg := range staticConglomerateGroupsOrdered {
			if sg.ID == groupKey {
				groupData = sg
				break
			}
		}
	}

	var allRows []graphHolderRow
	var activeTickers []string
	var graphTitle, hubName string

	if targetTicker != "" {
		rows, err := s.DB.Query(ctx, `
			SELECT ticker, name, pct::float8, COALESCE(is_controlling, FALSE)
			FROM company_shareholders
			WHERE ticker = $1
			  AND name NOT ILIKE 'Masyarakat%'
			  AND name NOT ILIKE 'Publik%'
			  AND name NOT ILIKE '%Treasury%'
			  AND pct >= $2
			ORDER BY pct DESC
		`, targetTicker, minPct)
		if err != nil {
			return NetworkGraphResponse{}, err
		}
		var directRows []graphHolderRow
		var shNames []string
		for rows.Next() {
			var r graphHolderRow
			if scanErr := rows.Scan(&r.ticker, &r.ownerName, &r.pct, &r.isControlling); scanErr != nil {
				rows.Close()
				return NetworkGraphResponse{}, scanErr
			}
			directRows = append(directRows, r)
			shNames = append(shNames, r.ownerName)
		}
		rows.Close()

		var siblingRows []graphHolderRow
		if len(shNames) > 0 {
			sibRows, sibErr := s.DB.Query(ctx, `
				SELECT ticker, name, pct::float8, COALESCE(is_controlling, FALSE)
				FROM company_shareholders
				WHERE name = ANY($1)
				  AND ticker != $2
				  AND name NOT ILIKE 'Masyarakat%'
				  AND name NOT ILIKE 'Publik%'
				  AND name NOT ILIKE '%Treasury%'
				  AND pct >= $3
				ORDER BY pct DESC
				LIMIT 40
			`, shNames, targetTicker, minPct)
			if sibErr != nil {
				return NetworkGraphResponse{}, sibErr
			}
			for sibRows.Next() {
				var r graphHolderRow
				if scanErr := sibRows.Scan(&r.ticker, &r.ownerName, &r.pct, &r.isControlling); scanErr != nil {
					sibRows.Close()
					return NetworkGraphResponse{}, scanErr
				}
				siblingRows = append(siblingRows, r)
			}
			sibRows.Close()
		}

		allRows = append(allRows, directRows...)
		allRows = append(allRows, siblingRows...)

		seenTicker := map[string]bool{targetTicker: true}
		activeTickers = []string{targetTicker}
		for _, r := range allRows {
			if !seenTicker[r.ticker] {
				seenTicker[r.ticker] = true
				activeTickers = append(activeTickers, r.ticker)
			}
		}
		graphTitle = "Jejaring Kepemilikan " + targetTicker
		hubName = targetTicker
	} else {
		activeTickers = append([]string(nil), groupData.KnownTickers...)
		rows, err := s.DB.Query(ctx, `
			SELECT ticker, name, pct::float8, COALESCE(is_controlling, FALSE)
			FROM company_shareholders
			WHERE ticker = ANY($1)
			  AND name NOT ILIKE 'Masyarakat%'
			  AND name NOT ILIKE 'Publik%'
			  AND name NOT ILIKE '%Treasury%'
			  AND pct >= $2
			ORDER BY ticker, pct DESC
		`, activeTickers, minPct)
		if err != nil {
			return NetworkGraphResponse{}, err
		}
		for rows.Next() {
			var r graphHolderRow
			if scanErr := rows.Scan(&r.ticker, &r.ownerName, &r.pct, &r.isControlling); scanErr != nil {
				rows.Close()
				return NetworkGraphResponse{}, scanErr
			}
			allRows = append(allRows, r)
		}
		rows.Close()
		graphTitle = groupData.Name
		hubName = groupData.HubEntity
	}

	type compMeta struct{ name, sector string }
	compMap := make(map[string]compMeta)
	if len(activeTickers) > 0 {
		cRows, cErr := s.DB.Query(ctx, `
			SELECT t.ticker,
			       COALESCE(t.name, t.ticker) AS company_name,
			       COALESCE(cg.groups, t.sector, '-') AS group_name
			FROM tickers t
			LEFT JOIN (
			    SELECT ticker, STRING_AGG(DISTINCT group_name, ' / ' ORDER BY group_name) AS groups
			    FROM ownership_staging.v_conglomerate_members
			    WHERE ticker = ANY($1)
			    GROUP BY ticker
			) cg ON cg.ticker = t.ticker
			WHERE t.ticker = ANY($1)
		`, activeTickers)
		if cErr != nil {
			return NetworkGraphResponse{}, cErr
		}
		for cRows.Next() {
			var t, cname, sec string
			if scanErr := cRows.Scan(&t, &cname, &sec); scanErr == nil {
				compMap[t] = compMeta{name: cname, sector: sec}
			}
		}
		cRows.Close()
	}

	activeSet := make(map[string]bool, len(activeTickers))
	for _, t := range activeTickers {
		activeSet[t] = true
	}

	nodeOrder := make([]string, 0, len(activeTickers)+len(allRows))
	nodesMap := make(map[string]*GraphNode)
	links := []GraphLink{}

	for _, t := range activeTickers {
		tID := "emiten:" + t
		cInfo, ok := compMap[t]
		cName, cSec := t, "-"
		if ok {
			cName, cSec = cInfo.name, cInfo.sector
		}
		val := 20
		if t == targetTicker {
			val = 26
		}
		if _, exists := nodesMap[tID]; !exists {
			nodeOrder = append(nodeOrder, tID)
		}
		nodesMap[tID] = &GraphNode{
			ID:     tID,
			Label:  t,
			Name:   cName,
			Type:   "emiten",
			Sector: cSec,
			Val:    val,
		}
	}

	hubParts := strings.Split(strings.ToLower(hubName), "/")
	for _, r := range allRows {
		ownerClean := strings.ToLower(strings.TrimSpace(r.ownerName))
		var sourceID string
		if mappedTicker, ok := knownNameToTicker[ownerClean]; ok && activeSet[mappedTicker] {
			sourceID = "emiten:" + mappedTicker
		} else {
			slug := strings.Trim(slugRe.ReplaceAllString(ownerClean, "-"), "-")
			if len(slug) > 40 {
				slug = slug[:40]
			}
			sourceID = "owner:" + slug
			if _, exists := nodesMap[sourceID]; !exists {
				isHub := false
				for _, hPart := range hubParts {
					if strings.Contains(ownerClean, hPart) {
						isHub = true
						break
					}
				}
				shortLabel := r.ownerName
				runes := []rune(r.ownerName)
				if len(runes) > 22 {
					shortLabel = string(runes[:20]) + "..."
				}
				nodeType := "owner"
				val := 16
				if isHub {
					nodeType = "hub"
					val = 30
				}
				nodeOrder = append(nodeOrder, sourceID)
				nodesMap[sourceID] = &GraphNode{
					ID:            sourceID,
					Label:         shortLabel,
					Name:          r.ownerName,
					Type:          nodeType,
					Sector:        "Pengendali / Pemegang Saham",
					IsControlling: r.isControlling,
					Val:           val,
				}
			}
		}

		targetID := "emiten:" + r.ticker
		if n, ok := nodesMap[sourceID]; ok {
			n.Degree++
		}
		if n, ok := nodesMap[targetID]; ok {
			n.Degree++
		}
		links = append(links, GraphLink{
			Source:        sourceID,
			Target:        targetID,
			Pct:           roundEven(r.pct, 2),
			IsControlling: r.isControlling,
			Label:         fmt.Sprintf("%.1f%%", r.pct),
		})
	}

	nodes := make([]GraphNode, 0, len(nodeOrder))
	for _, id := range nodeOrder {
		if n, ok := nodesMap[id]; ok {
			nodes = append(nodes, *n)
		}
	}

	orderedGroups := s.orderedGroups(ctx)
	availGroups := make([]AvailGroup, 0, len(orderedGroups))
	for _, g := range orderedGroups {
		availGroups = append(availGroups, AvailGroup{
			ID:          g.ID,
			Name:        g.Name,
			TickerCount: len(g.KnownTickers),
		})
	}

	return NetworkGraphResponse{
		GroupID:   groupKey,
		GroupName: graphTitle,
		HubEntity: hubName,
		Nodes:     nodes,
		Links:     links,
		Stats: GraphStats{
			TotalNodes:   len(nodes),
			TotalEmitens: len(activeTickers),
			TotalLinks:   len(links),
			HubEntity:    hubName,
		},
		AvailableGroups: availGroups,
	}, nil
}

// GetStockShareholders returns shareholders breakdown and summary for one issuer.
func (s *Service) GetStockShareholders(ctx context.Context, ticker string) (StockShareholdersResponse, error) {
	normTicker := strings.TrimSpace(strings.ReplaceAll(strings.ToUpper(ticker), ".JK", ""))
	rows, err := s.DB.Query(ctx, `
		SELECT name, COALESCE(pct, 0)::float8, COALESCE(is_controlling, FALSE)
		FROM company_shareholders
		WHERE ticker = $1
		ORDER BY pct DESC, name ASC
	`, normTicker)
	if err != nil {
		return StockShareholdersResponse{}, err
	}
	defer rows.Close()

	items := []ShareholderItem{}
	controllingEntities := []string{}
	var controllingPct, otherSubstantialPct, publicPct, treasuryPct, managementPct float64

	for rows.Next() {
		var name string
		var rawPct float64
		var isCtrl bool
		if scanErr := rows.Scan(&name, &rawPct, &isCtrl); scanErr != nil {
			return StockShareholdersResponse{}, scanErr
		}
		pct := roundEven(rawPct, 4)
		cat := ClassifyShareholder(name, pct, isCtrl)
		switch cat {
		case "controlling":
			controllingPct += pct
			controllingEntities = append(controllingEntities, name)
		case "treasury":
			treasuryPct += pct
		case "public":
			publicPct += pct
		case "substantial":
			otherSubstantialPct += pct
		default:
			managementPct += pct
		}
		items = append(items, ShareholderItem{
			Name:          name,
			Pct:           pct,
			IsControlling: isCtrl,
			Category:      cat,
		})
	}
	if err := rows.Err(); err != nil {
		return StockShareholdersResponse{}, err
	}
	if len(items) == 0 {
		return StockShareholdersResponse{}, fmt.Errorf("%w: %s", ErrNotFound, normTicker)
	}

	return StockShareholdersResponse{
		Ticker:       normTicker,
		Count:        len(items),
		Shareholders: items,
		Summary: ShareholderSummary{
			ControllingPct:      roundEven(controllingPct, 2),
			OtherSubstantialPct: roundEven(otherSubstantialPct, 2),
			PublicPct:           roundEven(publicPct, 2),
			TreasuryPct:         roundEven(treasuryPct, 2),
			ManagementPct:       roundEven(managementPct, 2),
			FreeFloatEst:        roundEven(publicPct, 2),
			ControllingEntities: controllingEntities,
		},
	}, nil
}
