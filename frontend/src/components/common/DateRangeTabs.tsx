/**
 * 日期範圍 Tabs 元件
 * 用於快速切換不同的日期範圍（今日、本週、本月）
 */

import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";

export type DateRangeType = "today" | "week" | "month";

interface DateRangeTabsProps {
  value: DateRangeType;
  onValueChange: (value: DateRangeType) => void;
  className?: string;
}

/**
 * 取得日期範圍的標籤
 */
export const getDateRangeLabel = (range: DateRangeType): string => {
  switch (range) {
    case "today":
      return "今日";
    case "week":
      return "本週";
    case "month":
      return "本月";
    default:
      return "今日";
  }
};

/**
 * 計算日期範圍
 * @param range 日期範圍類型
 * @returns { startDate, endDate } ISO 8601 格式的日期字串
 */
export const calculateDateRange = (
  range: DateRangeType
): { startDate: string; endDate: string } => {
  const now = new Date();
  const year = now.getFullYear();
  const month = now.getMonth();
  const day = now.getDate();
  const dayOfWeek = now.getDay(); // 0 = 星期日, 1 = 星期一, ...

  let start: Date;
  const end: Date = now;

  switch (range) {
    case "today":
      // 今日：從今天 00:00:00 到現在
      start = new Date(year, month, day);
      break;

    case "week":
      // 本週：從本週一到現在
      // 如果今天是星期日(0),則往前推 6 天到星期一
      // 如果今天是星期一(1),則往前推 0 天
      // 如果今天是星期二(2),則往前推 1 天
      const daysToMonday = dayOfWeek === 0 ? 6 : dayOfWeek - 1;
      start = new Date(year, month, day - daysToMonday);
      break;

    case "month":
      // 本月：從本月 1 日到現在
      start = new Date(year, month, 1);
      break;

    default:
      start = new Date(year, month, day);
  }

  return {
    startDate: start.toISOString().split("T")[0],
    endDate: end.toISOString().split("T")[0],
  };
};

export function DateRangeTabs({
  value,
  onValueChange,
  className,
}: DateRangeTabsProps) {
  return (
    <Tabs value={value} onValueChange={onValueChange} className={className}>
      <TabsList className="grid w-full grid-cols-3">
        <TabsTrigger value="today">今日</TabsTrigger>
        <TabsTrigger value="week">本週</TabsTrigger>
        <TabsTrigger value="month">本月</TabsTrigger>
      </TabsList>
    </Tabs>
  );
}

