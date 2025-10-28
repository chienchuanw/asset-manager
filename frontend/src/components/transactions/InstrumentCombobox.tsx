import { useState, useMemo } from "react";
import { Check, ChevronsUpDown, Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
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
  /** 選擇標的時的回調函式 */
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
  onSelect,
  disabled = false,
  placeholder = "選擇標的...",
  searchPlaceholder = "搜尋標的代碼或名稱...",
}: InstrumentComboboxProps) {
  const [open, setOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");

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

  // 找到當前選中的標的
  const selectedInstrument = instruments.find((i) => i.symbol === value);

  // 顯示文字：如果有選中的標的，顯示「代碼 - 名稱」，否則顯示 placeholder
  const displayText = selectedInstrument
    ? `${selectedInstrument.symbol} - ${selectedInstrument.name}`
    : placeholder;

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          disabled={disabled}
          className={cn(
            "w-full justify-between",
            !value && "text-muted-foreground"
          )}
        >
          <span className="truncate">{displayText}</span>
          {isLoading ? (
            <Loader2 className="ml-2 h-4 w-4 shrink-0 animate-spin opacity-50" />
          ) : (
            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-full p-0" align="start">
        <Command shouldFilter={false}>
          <CommandInput
            placeholder={searchPlaceholder}
            value={searchQuery}
            onValueChange={setSearchQuery}
          />
          <CommandList>
            <CommandEmpty>
              {isLoading ? "載入中..." : "查無結果"}
            </CommandEmpty>
            <CommandGroup>
              {filteredInstruments.map((instrument) => (
                <CommandItem
                  key={instrument.symbol}
                  value={instrument.symbol}
                  onSelect={() => {
                    onSelect(instrument);
                    setOpen(false);
                    setSearchQuery(""); // 清空搜尋
                  }}
                >
                  <Check
                    className={cn(
                      "mr-2 h-4 w-4",
                      value === instrument.symbol ? "opacity-100" : "opacity-0"
                    )}
                  />
                  <div className="flex flex-col">
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
  );
}

