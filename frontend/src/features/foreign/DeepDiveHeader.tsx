import { useState, useMemo, useRef, useEffect } from "react";
import { Link } from "react-router-dom";
import { Icon } from "@iconify/react";
import { POPULAR_TICKERS, WINDOW_OPTIONS } from "./deepdiveHelpers";

interface DeepDiveHeaderProps {
  activeTicker: string;
  setActiveTicker: (ticker: string) => void;
  windowDays: number;
  setWindowDays: (days: number) => void;
  companyName: string;
  groupName: string;
  latestDate?: string;
  isFetching: boolean;
  isLoading: boolean;
  allTickers: string[];
  isLoadingUniverse: boolean;
}

export default function DeepDiveHeader({
  activeTicker,
  setActiveTicker,
  windowDays,
  setWindowDays,
  companyName,
  groupName,
  latestDate,
  isFetching,
  isLoading,
  allTickers,
  isLoadingUniverse,
}: DeepDiveHeaderProps) {
  const [searchTerm, setSearchTerm] = useState("");
  const [isOpen, setIsOpen] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const filteredTickers = useMemo(() => {
    const term = searchTerm.trim().toUpperCase().replace(".JK", "");
    if (!term) return [];
    return allTickers.filter((t: string) => t.includes(term)).slice(0, 10);
  }, [allTickers, searchTerm]);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "/" && document.activeElement !== inputRef.current) {
        e.preventDefault();
        inputRef.current?.focus();
      } else if (e.key === "Escape") {
        setIsOpen(false);
        inputRef.current?.blur();
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  const handleSelectTicker = (t: string) => {
    const clean = t.trim().toUpperCase().replace(".JK", "");
    if (!clean) return;
    setActiveTicker(clean);
    setSearchTerm("");
    setIsOpen(false);
    inputRef.current?.blur();
  };

  const handleEnterSubmit = () => {
    const term = searchTerm.trim().toUpperCase().replace(".JK", "");
    if (!term) return;

    const exact = allTickers.find((t: string) => t === term);
    if (exact) {
      handleSelectTicker(exact);
      return;
    }
    if (filteredTickers.length > 0) {
      handleSelectTicker(filteredTickers[0]);
      return;
    }
    handleSelectTicker(term);
  };

  return (
    <div className="bg-[var(--md-sys-color-surface-container)] border border-[var(--md-sys-color-outline-variant)] rounded-2xl p-4 sm:p-5 shadow-xs space-y-4">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-3.5 pb-3.5 border-b border-[var(--md-sys-color-outline-variant)]/60">
        <div className="flex items-center gap-3">
          <div className="w-11 h-11 sm:w-12 sm:h-12 rounded-xl bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-primary)] flex items-center justify-center shrink-0 border border-[var(--md-sys-color-outline-variant)] shadow-2xs font-extrabold text-base sm:text-lg">
            {activeTicker}
          </div>
          <div>
            <div className="flex flex-wrap items-center gap-2">
              <h2 className="text-lg sm:text-xl font-extrabold text-[var(--md-sys-color-on-surface)] tracking-tight">
                {activeTicker}
              </h2>
              <span className="text-[10px] font-semibold px-2 py-0.5 rounded-full bg-[var(--md-sys-color-secondary-container)] text-[var(--md-sys-color-on-secondary-container)]">
                {groupName}
              </span>
              {latestDate && (
                <span className="text-[10px] font-mono px-2 py-0.5 rounded-full bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface-variant)] border border-[var(--md-sys-color-outline-variant)]/60">
                  Data s/d {latestDate}
                </span>
              )}
              <Link
                to={`/${activeTicker}`}
                className="text-xs text-[var(--md-sys-color-primary)] hover:underline flex items-center gap-0.5 font-medium ml-1"
                title="Lihat halaman analisis lengkap emiten ini"
              >
                <span>Detail Saham</span>
                <Icon icon="ph:arrow-square-out-duotone" width="13" />
              </Link>
            </div>
            <p className="text-xs text-[var(--md-sys-color-on-surface-variant)] mt-0.5 font-medium">
              {companyName}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2 self-start md:self-auto">
          {isFetching && !isLoading && (
            <span className="text-[10px] font-bold uppercase text-blue-500 bg-blue-500/10 px-2.5 py-1 rounded-lg border border-blue-500/20 shadow-2xs flex items-center gap-1.5">
              <Icon icon="ph:spinner-gap-bold" className="animate-spin" width="12" />
              <span>Update...</span>
            </span>
          )}

          <div className="flex items-center p-1 rounded-xl bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] text-xs font-semibold">
            {WINDOW_OPTIONS.map((opt) => (
              <button
                key={opt.value}
                onClick={() => setWindowDays(opt.value)}
                className={`px-2.5 sm:px-3 py-1 rounded-lg transition-all text-[11px] ${
                  windowDays === opt.value
                    ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] shadow-2xs"
                    : "text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)]"
                }`}
              >
                {opt.label}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Baris 2: Search Input & Quick Emiten Chips */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
        <div className="relative sm:w-80 z-30">
          <div className="flex items-center bg-[var(--md-sys-color-surface)] border border-[var(--md-sys-color-outline-variant)] hover:border-[var(--md-sys-color-primary)] focus-within:border-[var(--md-sys-color-primary)] rounded-xl px-3 py-1.5 shadow-2xs transition-all duration-200">
            <Icon
              icon="ph:magnifying-glass-duotone"
              className="text-[var(--md-sys-color-primary)] mr-2 flex-shrink-0"
              width="16"
              height="16"
            />
            <input
              ref={inputRef}
              type="text"
              placeholder="Cari kode saham (misal: BBCA, ASII)..."
              value={searchTerm}
              onChange={(e) => {
                setSearchTerm(e.target.value);
                setIsOpen(true);
              }}
              onFocus={() => setIsOpen(true)}
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  e.preventDefault();
                  handleEnterSubmit();
                } else if (e.key === "Escape") {
                  setIsOpen(false);
                  inputRef.current?.blur();
                }
              }}
              className="w-full bg-transparent border-none outline-none text-xs text-[var(--md-sys-color-on-surface)] placeholder-[var(--md-sys-color-on-surface-variant)] font-mono uppercase tracking-wider"
            />
            {searchTerm ? (
              <button
                type="button"
                onClick={() => {
                  setSearchTerm("");
                  setIsOpen(false);
                }}
                className="text-[var(--md-sys-color-on-surface-variant)] hover:text-[var(--md-sys-color-on-surface)] p-0.5 cursor-pointer"
              >
                <Icon icon="ph:x-circle-fill" width="14" />
              </button>
            ) : (
              <div className="hidden sm:flex items-center gap-1 text-[9px] font-mono font-bold text-[var(--md-sys-color-on-surface-variant)] bg-[var(--md-sys-color-surface-container-high)] border border-[var(--md-sys-color-outline-variant)] px-1.5 py-0.5 rounded shrink-0">
                <kbd className="text-[var(--md-sys-color-primary)]">/</kbd>
              </div>
            )}
          </div>

          {isOpen && searchTerm.trim() && (
            <>
              <div className="fixed inset-0 z-40" onClick={() => setIsOpen(false)} />
              <div className="absolute top-full left-0 right-0 mt-1.5 bg-[var(--md-sys-color-surface-container-highest)] border border-[var(--md-sys-color-outline-variant)] rounded-xl overflow-hidden shadow-xl z-50 max-h-72 overflow-y-auto">
                {isLoadingUniverse && allTickers.length === 0 ? (
                  <div className="px-4 py-3 text-xs text-[var(--md-sys-color-on-surface-variant)] font-mono flex items-center gap-2">
                    <Icon icon="ph:spinner-gap-bold" className="animate-spin text-[var(--md-sys-color-primary)]" width="14" />
                    <span>Memuat daftar saham bursa...</span>
                  </div>
                ) : filteredTickers.length > 0 ? (
                  filteredTickers.map((t: string) => (
                    <button
                      key={t}
                      type="button"
                      className="w-full flex items-center justify-between px-3.5 py-2.5 text-xs font-mono font-bold text-[var(--md-sys-color-on-surface)] hover:bg-[var(--md-sys-color-surface-container-low)] hover:text-[var(--md-sys-color-primary)] transition-colors border-b border-[var(--md-sys-color-outline-variant)] last:border-0 cursor-pointer text-left"
                      onClick={() => handleSelectTicker(t)}
                    >
                      <div className="flex items-center gap-2.5">
                        <span className="w-2 h-2 rounded-full bg-emerald-500"></span>
                        <span className="text-[var(--md-sys-color-primary)] font-extrabold">{t}</span>
                        <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] font-normal">Saham IDX</span>
                      </div>
                      <Icon icon="ph:arrow-up-right-bold" className="text-[var(--md-sys-color-on-surface-variant)]" width="12" />
                    </button>
                  ))
                ) : (
                  <div
                    className="px-4 py-3 text-xs text-[var(--md-sys-color-on-surface-variant)] font-mono flex items-center justify-between hover:bg-[var(--md-sys-color-surface-container-low)] cursor-pointer"
                    onClick={handleEnterSubmit}
                  >
                    <span>Cari langsung saham &quot;{searchTerm.toUpperCase().trim()}&quot;</span>
                    <kbd className="text-[var(--md-sys-color-primary)] font-bold text-[10px] bg-[var(--md-sys-color-surface-container-high)] px-2 py-0.5 rounded border border-[var(--md-sys-color-outline-variant)]">
                      ↵ Enter
                    </kbd>
                  </div>
                )}
              </div>
            </>
          )}
        </div>

        <div className="flex items-center gap-1.5 overflow-x-auto pb-0.5 no-scrollbar text-xs font-medium">
          <span className="text-[10px] text-[var(--md-sys-color-on-surface-variant)] uppercase mr-0.5 shrink-0 font-semibold">
            Populer:
          </span>
          {POPULAR_TICKERS.map((t) => (
            <button
              key={t}
              onClick={() => setActiveTicker(t)}
              className={`px-2 py-0.5 rounded-md font-mono text-[11px] shrink-0 transition-all ${
                activeTicker === t
                  ? "bg-[var(--md-sys-color-primary)] text-[var(--md-sys-color-on-primary)] font-bold shadow-2xs"
                  : "bg-[var(--md-sys-color-surface-container-high)] text-[var(--md-sys-color-on-surface)] hover:bg-[var(--md-sys-color-surface-container-highest)] border border-[var(--md-sys-color-outline-variant)]"
              }`}
            >
              {t}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}
