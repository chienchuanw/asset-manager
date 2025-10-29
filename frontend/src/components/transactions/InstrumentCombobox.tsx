import { useState, useMemo, useRef, useEffect } from "react";
import { Check, ChevronsUpDown, Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { Input } from "@/components/ui/input";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useInstruments } from "@/hooks";
import { AssetType } from "@/types/transaction";
import { Instrument } from "@/types/instrument";

interface InstrumentComboboxProps {
  /** 當前選中的標的代碼 */
  value: string;
  /** 資產類型（用於篩選標的清單） */
  assetType: AssetType;
  /** 當代碼改變時的回調函式 */
  onChange: (symbol: string) => void;
  /** 選擇標的時的回調函式（從清單選擇時會觸發） */
  onSelect: (instrument: Instrument) => void;
  /** 是否禁用 */
  disabled?: boolean;
  /** 自訂 placeholder */
  placeholder?: string;
  /** 自訂搜尋 placeholder */
  searchPlaceholder?: string;
}

/**
 * 標的選擇 Combobox 元件
 *
 * 功能：
 * 1. 根據資產類型顯示對應的標的清單
 * 2. 支援模糊搜尋（symbol 或 name）
 * 3. 歷史交易記錄的標的會優先顯示
 * 4. 支援鍵盤導航
 * 5. 允許使用者輸入不在清單中的標的（透過父元件處理）
 *
 * @example
 * ```tsx
 * <InstrumentCombobox
 *   value={form.watch("symbol")}
 *   assetType={form.watch("asset_type")}
 *   onSelect={(instrument) => {
 *     form.setValue("symbol", instrument.symbol);
 *     form.setValue("name", instrument.name);
 *   }}
 * />
 * ```
 */
export function InstrumentCombobox({
  value,
  assetType,
  onChange,
  onSelect,
  disabled = false,
  placeholder = "輸入或選擇代碼...",
  searchPlaceholder = "搜尋標的代碼或名稱...",
}: InstrumentComboboxProps) {
  const [open, setOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);

  // 取得標的清單
  const { data: instruments = [], isLoading } = useInstruments(assetType);

  // 根據搜尋關鍵字篩選標的
  const filteredInstruments = useMemo(() => {
    if (!searchQuery || searchQuery.trim() === "") {
      return instruments;
    }

    const normalizedQuery = searchQuery.toLowerCase().trim();
    return instruments.filter(
      (instrument) =>
        instrument.symbol.toLowerCase().includes(normalizedQuery) ||
        instrument.name.toLowerCase().includes(normalizedQuery)
    );
  }, [instruments, searchQuery]);

  // 當 Popover 打開時，同步搜尋框的值為當前 input 的值
  useEffect(() => {
    if (open && value) {
      setSearchQuery(value);
    }
  }, [open, value]);

  return (
    <div className="relative flex items-center gap-2">
      {/* 手動輸入框 */}
      <Input
        ref={inputRef}
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value.toUpperCase())}
        placeholder={placeholder}
        disabled={disabled}
        className="flex-1"
      />

      {/* 下拉選單按鈕 */}
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <button
            type="button"
            disabled={disabled}
            className={cn(
              "absolute right-0 top-0 h-full px-3 flex items-center",
              "hover:bg-accent rounded-r-md transition-colors",
              disabled && "opacity-50 cursor-not-allowed"
            )}
            aria-label="開啟標的清單"
          >
            {isLoading ? (
              <Loader2 className="h-4 w-4 animate-spin opacity-50" />
            ) : (
              <ChevronsUpDown className="h-4 w-4 opacity-50" />
            )}
          </button>
        </PopoverTrigger>
        <PopoverContent className="w-[400px] p-0" align="start">
          <Command shouldFilter={false}>
            <CommandInput
              placeholder={searchPlaceholder}
              value={searchQuery}
              onValueChange={setSearchQuery}
            />
            <CommandList>
              <CommandEmpty>
                {isLoading ? (
                  "載入中..."
                ) : instruments.length === 0 ? (
                  <div className="py-6 text-center text-sm">
                    <p className="text-muted-foreground">
                      尚無交易記錄，請手動輸入代碼
                    </p>
                  </div>
                ) : (
                  "查無結果"
                )}
              </CommandEmpty>
              <CommandGroup>
                {filteredInstruments.map((instrument) => (
                  <CommandItem
                    key={instrument.symbol}
                    value={instrument.symbol}
                    onSelect={() => {
                      onSelect(instrument);
                      setOpen(false);
                      setSearchQuery("");
                    }}
                  >
                    <Check
                      className={cn(
                        "mr-2 h-4 w-4",
                        value === instrument.symbol
                          ? "opacity-100"
                          : "opacity-0"
                      )}
                    />
                    <div className="flex flex-col flex-1">
                      <span className="font-medium">{instrument.symbol}</span>
                      <span className="text-xs text-muted-foreground">
                        {instrument.name}
                      </span>
                    </div>
                    {instrument.from_history && (
                      <span className="ml-auto text-xs text-muted-foreground">
                        常用
                      </span>
                    )}
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
    </div>
  );
}
