"use client";

import { useState } from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

interface DailyDateNavigatorProps {
  date: Date;
  onDateChange: (date: Date) => void;
}

/**
 * 日期導航元件
 *
 * 用於「今日」分頁的日期切換，包含：
 * - 左箭頭：切換到前一天
 * - 日期顯示（可點擊）：彈出日曆選擇器
 * - 右箭頭：切換到後一天
 */
export function DailyDateNavigator({
  date,
  onDateChange,
}: DailyDateNavigatorProps) {
  const [isCalendarOpen, setIsCalendarOpen] = useState(false);

  // 格式化日期為 YYYY/MM/DD
  const formatDate = (date: Date): string => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const day = String(date.getDate()).padStart(2, "0");
    return `${year}/${month}/${day}`;
  };

  // 切換到前一天
  const handlePreviousDay = () => {
    const newDate = new Date(date);
    newDate.setDate(newDate.getDate() - 1);
    onDateChange(newDate);
  };

  // 切換到後一天
  const handleNextDay = () => {
    const newDate = new Date(date);
    newDate.setDate(newDate.getDate() + 1);
    onDateChange(newDate);
  };

  // 從日曆選擇日期
  const handleCalendarSelect = (selectedDate: Date | undefined) => {
    if (selectedDate) {
      onDateChange(selectedDate);
      setIsCalendarOpen(false);
    }
  };

  return (
    <div className="flex items-center justify-center gap-2">
      {/* 前一天按鈕 */}
      <Button
        variant="outline"
        size="icon"
        onClick={handlePreviousDay}
        aria-label="前一天"
      >
        <ChevronLeft className="h-4 w-4" />
      </Button>

      {/* 日期選擇器 */}
      <Popover open={isCalendarOpen} onOpenChange={setIsCalendarOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            className="min-w-[180px] justify-center font-semibold text-lg"
          >
            {formatDate(date)}
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-auto p-0" align="center">
          <Calendar
            mode="single"
            selected={date}
            onSelect={handleCalendarSelect}
          />
        </PopoverContent>
      </Popover>

      {/* 後一天按鈕 */}
      <Button
        variant="outline"
        size="icon"
        onClick={handleNextDay}
        aria-label="後一天"
      >
        <ChevronRight className="h-4 w-4" />
      </Button>
    </div>
  );
}
