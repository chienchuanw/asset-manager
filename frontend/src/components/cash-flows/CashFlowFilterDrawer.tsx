/**
 * 現金流記錄進階篩選 Drawer 元件
 * 用於手機友善的篩選介面
 */

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { CashFlowType, getCashFlowTypeLabel } from "@/types/cash-flow";
import { Filter, X } from "lucide-react";
import { DateRange } from "react-day-picker";

interface CashFlowFilterDrawerProps {
  filterType: CashFlowType | "all";
  dateRange?: DateRange;
  onFilterTypeChange: (value: CashFlowType | "all") => void;
  onDateRangeChange: (range: DateRange | undefined) => void;
  onReset: () => void;
}

export function CashFlowFilterDrawer({
  filterType,
  dateRange,
  onFilterTypeChange,
  onDateRangeChange,
  onReset,
}: CashFlowFilterDrawerProps) {
  const [open, setOpen] = useState(false);
  const [tempDateRange, setTempDateRange] = useState<DateRange | undefined>(
    dateRange
  );

  // 檢查是否有啟用篩選
  const hasActiveFilters = filterType !== "all" || dateRange;

  // 處理套用篩選
  const handleApply = () => {
    onDateRangeChange(tempDateRange);
    setOpen(false);
  };

  // 處理重置
  const handleReset = () => {
    setTempDateRange(undefined);
    onReset();
  };

  return (
    <Drawer open={open} onOpenChange={setOpen}>
      <DrawerTrigger asChild>
        <Button variant="outline" size="sm" className="relative">
          <Filter className="h-4 w-4 mr-2" />
          篩選
          {hasActiveFilters && (
            <span className="absolute -top-1 -right-1 h-2 w-2 bg-blue-600 rounded-full" />
          )}
        </Button>
      </DrawerTrigger>
      <DrawerContent className="max-h-[85vh]">
        <DrawerHeader>
          <DrawerTitle>進階篩選</DrawerTitle>
          <DrawerDescription>設定篩選條件以縮小搜尋範圍</DrawerDescription>
        </DrawerHeader>

        <div className="px-4 pb-4 overflow-y-auto">
          <div className="space-y-6">
            {/* 現金流類型篩選 */}
            <div className="space-y-2">
              <Label htmlFor="filter-type">現金流類型</Label>
              <Select value={filterType} onValueChange={onFilterTypeChange}>
                <SelectTrigger id="filter-type">
                  <SelectValue placeholder="選擇現金流類型" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">全部類型</SelectItem>
                  <SelectItem value={CashFlowType.INCOME}>
                    {getCashFlowTypeLabel(CashFlowType.INCOME)}
                  </SelectItem>
                  <SelectItem value={CashFlowType.EXPENSE}>
                    {getCashFlowTypeLabel(CashFlowType.EXPENSE)}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* 日期範圍篩選 */}
            <div className="space-y-2">
              <Label>自訂日期範圍</Label>
              <div className="flex flex-col items-center">
                <Calendar
                  mode="range"
                  selected={tempDateRange}
                  onSelect={setTempDateRange}
                  numberOfMonths={1}
                  className="rounded-md border"
                />
                {tempDateRange?.from && (
                  <div className="mt-2 text-sm text-muted-foreground">
                    {tempDateRange.from.toLocaleDateString("zh-TW")}
                    {tempDateRange.to &&
                      ` - ${tempDateRange.to.toLocaleDateString("zh-TW")}`}
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>

        <DrawerFooter className="flex-row gap-2">
          <Button
            variant="outline"
            onClick={handleReset}
            className="flex-1"
            disabled={!hasActiveFilters}
          >
            <X className="h-4 w-4 mr-2" />
            重置
          </Button>
          <DrawerClose asChild>
            <Button onClick={handleApply} className="flex-1">
              套用篩選
            </Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  );
}

