/**
 * 本週/本月 Tabs 元件
 * 用於快速切換本週和本月的日期範圍（不包含今日）
 */

import { useTranslations } from "next-intl";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";

export type WeekMonthType = "week" | "month";

interface WeekMonthTabsProps {
  value: WeekMonthType;
  onValueChange: (value: WeekMonthType) => void;
  className?: string;
}

/**
 * 本週/本月 Tabs 元件
 */
export function WeekMonthTabs({
  value,
  onValueChange,
  className,
}: WeekMonthTabsProps) {
  const t = useTranslations("common");

  // 處理 Tabs 元件的 onValueChange 類型轉換
  const handleValueChange = (newValue: string) => {
    // 確保值是有效的 WeekMonthType
    if (newValue === "week" || newValue === "month") {
      onValueChange(newValue as WeekMonthType);
    }
  };

  return (
    <Tabs value={value} onValueChange={handleValueChange} className={className}>
      <TabsList className="grid w-full grid-cols-2">
        <TabsTrigger value="week">{t("thisWeek")}</TabsTrigger>
        <TabsTrigger value="month">{t("thisMonth")}</TabsTrigger>
      </TabsList>
    </Tabs>
  );
}
