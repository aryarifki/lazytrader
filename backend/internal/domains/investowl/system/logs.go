package system

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func logDir(smtrackerDir string) string {
	return filepath.Join(smtrackerDir, "debug")
}

func matchingLogFiles(dir, prefix string) []string {
	matches, _ := filepath.Glob(filepath.Join(dir, prefix+"_*.log"))
	standalone := filepath.Join(dir, prefix+".log")
	if info, err := os.Stat(standalone); err == nil && !info.IsDir() {
		found := false
		for _, m := range matches {
			if m == standalone {
				found = true
				break
			}
		}
		if !found {
			matches = append(matches, standalone)
		}
	}
	return matches
}

func readLastLines(path string, maxLines int) (string, int, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", 0, err
	}
	content := strings.ToValidUTF8(string(raw), "\uFFFD")
	lines := strings.SplitAfter(content, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}
	return strings.Join(lines, ""), len(lines), nil
}

// GetLatestLog returns the last 1000 lines of the newest log matching prefix.
func GetLatestLog(smtrackerDir, prefix string) string {
	dir := logDir(smtrackerDir)
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return fmt.Sprintf("Error: Folder log tidak ditemukan di %s. Pastikan path sudah benar.", dir)
	}
	files := matchingLogFiles(dir, prefix)
	if len(files) == 0 {
		return fmt.Sprintf("Belum ada log untuk '%s'. Siap dijalankan kapan saja.", prefix)
	}
	sort.Slice(files, func(i, j int) bool {
		si, _ := os.Stat(files[i])
		sj, _ := os.Stat(files[j])
		if si == nil || sj == nil {
			return files[i] > files[j]
		}
		return si.ModTime().After(sj.ModTime())
	})
	text, _, err := readLastLines(files[0], 1000)
	if err != nil {
		return fmt.Sprintf("Error membaca file %s: %v", files[0], err)
	}
	return text
}

// ListHistoryFiles returns up to 15 newest log files for a job prefix.
func ListHistoryFiles(smtrackerDir, job string) []LogHistoryFile {
	dir := logDir(smtrackerDir)
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return []LogHistoryFile{}
	}
	prefix := job
	if cfg, ok := JobConfigs[job]; ok && cfg.LogPrefix != "" {
		prefix = cfg.LogPrefix
	}
	files := matchingLogFiles(dir, prefix)
	type fileMeta struct {
		path string
		info os.FileInfo
	}
	items := make([]fileMeta, 0, len(files))
	for _, fp := range files {
		if st, err := os.Stat(fp); err == nil {
			items = append(items, fileMeta{path: fp, info: st})
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].info.ModTime().After(items[j].info.ModTime())
	})
	if len(items) > 15 {
		items = items[:15]
	}
	out := make([]LogHistoryFile, 0, len(items))
	for _, it := range items {
		out = append(out, LogHistoryFile{
			Filename:   filepath.Base(it.path),
			SizeBytes:  it.info.Size(),
			SizeKB:     math.Round((float64(it.info.Size())/1024.0)*10.0) / 10.0,
			ModifiedAt: it.info.ModTime().Format(time.RFC3339),
		})
	}
	return out
}

// SanitizeLogFilename strictly guards against path traversal.
func SanitizeLogFilename(filename string) (string, bool) {
	trimmed := strings.TrimSpace(filename)
	if trimmed == "" || strings.Contains(trimmed, "..") || strings.ContainsAny(trimmed, "/\\") {
		return "", false
	}
	clean := filepath.Base(trimmed)
	if clean == "." || clean == ".." {
		return "", false
	}
	return clean, true
}
