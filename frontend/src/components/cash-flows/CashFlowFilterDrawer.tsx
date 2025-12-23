/**
 * 現金流記錄進階篩選 Drawer 元件
 * 用於手機友善的篩選介面
 */

import { useState } from "react";
import { useTranslations } from "next-intl";
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
  const t = useTranslations("cashFlows");

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
          {t("filter")}
          {hasActiveFilters && (
            <span className="absolute -top-1 -right-1 h-2 w-2 bg-blue-600 rounded-full" />
          )}
        </Button>
      </DrawerTrigger>
      <DrawerContent className="max-h-[85vh] sm:max-w-2xl sm:mx-auto">
        <DrawerHeader>
          <DrawerTitle>{t("advancedFilter")}</DrawerTitle>
          <DrawerDescription>{t("setFilterConditions")}</DrawerDescription>
        </DrawerHeader>

        <div className="px-4 pb-4 overflow-y-auto">
          {/* 桌面版:使用網格佈局 */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-6">
            {/* 現金流類型篩選 */}
            <div className="space-y-2">
              <Label htmlFor="filter-type">{t("selectCashFlowType")}</Label>
              <Select value={filterType} onValueChange={onFilterTypeChange}>
                <SelectTrigger id="filter-type">
                  <SelectValue placeholder={t("selectCashFlowType")} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{t("allTypes")}</SelectItem>
                  <SelectItem value={CashFlowType.INCOME}>
                    {getCashFlowTypeLabel(CashFlowType.INCOME)}
                  </SelectItem>
                  <SelectItem value={CashFlowType.EXPENSE}>
                    {getCashFlowTypeLabel(CashFlowType.EXPENSE)}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* 日期範圍篩選 - 跨兩欄 */}
            <div className="space-y-2 sm:col-span-2">
              <Label>{t("date")}</Label>
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
            {t("reset")}
          </Button>
          <DrawerClose asChild>
            <Button onClick={handleApply} className="flex-1">
              {t("applyFilter")}
            </Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  );
}
