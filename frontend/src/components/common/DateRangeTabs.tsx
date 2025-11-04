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

  // 使用本地日期字串避免時區問題
  const today = formatDateToLocal(now); // 今天的日期字串
  const year = now.getFullYear();
  const month = now.getMonth();
  const day = now.getDate();
  const dayOfWeek = now.getDay(); // 0 = 星期日, 1 = 星期一, ...

  let startDate: string;
  const endDate: string = today;

  switch (range) {
    case "today":
      // 今日：開始和結束都是今天
      startDate = today;
      break;

    case "week":
      // 本週：從本週一到今天
      // 如果今天是星期日(0),則往前推 6 天到星期一
      // 如果今天是星期一(1),則往前推 0 天
      // 如果今天是星期二(2),則往前推 1 天
      const daysToMonday = dayOfWeek === 0 ? 6 : dayOfWeek - 1;
      const mondayDate = new Date(year, month, day - daysToMonday);
      startDate = formatDateToLocal(mondayDate);
      break;

    case "month":
      // 本月：從本月 1 日到今天
      const firstDayOfMonth = new Date(year, month, 1);
      startDate = formatDateToLocal(firstDayOfMonth);
      break;

    default:
      startDate = today;
  }

  return { startDate, endDate };
};

/**
 * 將 Date 物件格式化為本地日期字串 (YYYY-MM-DD)
 * 避免時區轉換問題
 */
const formatDateToLocal = (date: Date): string => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
};

export function DateRangeTabs({
  value,
  onValueChange,
  className,
}: DateRangeTabsProps) {
  // 處理 Tabs 元件的 onValueChange 類型轉換
  const handleValueChange = (newValue: string) => {
    // 確保值是有效的 DateRangeType
    if (newValue === "today" || newValue === "week" || newValue === "month") {
      onValueChange(newValue as DateRangeType);
    }
  };

  return (
    <Tabs value={value} onValueChange={handleValueChange} className={className}>
      <TabsList className="grid w-full grid-cols-3">
        <TabsTrigger value="today">今日</TabsTrigger>
        <TabsTrigger value="week">本週</TabsTrigger>
        <TabsTrigger value="month">本月</TabsTrigger>
      </TabsList>
    </Tabs>
  );
}
