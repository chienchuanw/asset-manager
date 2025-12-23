"use client";

import { useMemo } from "react";
import { useTranslations } from "next-intl";
import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { useCashFlows } from "@/hooks";
import { useIsMobile } from "@/hooks/use-mobile";
import { CashFlowType, formatAmount } from "@/types/cash-flow";
import { Skeleton } from "@/components/ui/skeleton";

interface MonthlyExpenseChartProps {
  selectedDate: string; // ISO 8601 格式 (YYYY-MM-DD)
}

/**
 * 當月/當週每日收入/支出圖表元件
 *
 * 根據使用者選擇的日期和裝置類型：
 * - 桌面版：顯示該月份每日的收入和支出統計
 * - 手機版：顯示該週每日的收入和支出統計
 * 使用 shadcn Bar Chart 渲染雙柱狀圖（綠色=收入，紅色=支出）
 */
export function MonthlyExpenseChart({
  selectedDate,
}: MonthlyExpenseChartProps) {
  const t = useTranslations("cashFlows");
  const tCommon = useTranslations("common");
  const isMobile = useIsMobile();

  // 計算日期範圍（桌面=當月，手機=當週）
  const { startDate, endDate, days, periodType } = useMemo(() => {
    const date = new Date(selectedDate);
    const year = date.getFullYear();
    const month = date.getMonth();

    // 格式化為 YYYY-MM-DD
    const formatDate = (d: Date) => {
      const y = d.getFullYear();
      const m = String(d.getMonth() + 1).padStart(2, "0");
      const day = String(d.getDate()).padStart(2, "0");
      return `${y}-${m}-${day}`;
    };

    if (isMobile) {
      // 手機版：計算當週範圍（週一到週日）
      const dayOfWeek = date.getDay(); // 0 = 週日, 1 = 週一, ..., 6 = 週六
      const diff = dayOfWeek === 0 ? -6 : 1 - dayOfWeek; // 調整到週一

      const weekStart = new Date(date);
      weekStart.setDate(date.getDate() + diff);

      const weekEnd = new Date(weekStart);
      weekEnd.setDate(weekStart.getDate() + 6); // 週日

      return {
        startDate: formatDate(weekStart),
        endDate: formatDate(weekEnd),
        days: 7,
        periodType: "week" as const,
        weekDays: [
          tCommon("monday"),
          tCommon("tuesday"),
          tCommon("wednesday"),
          tCommon("thursday"),
          tCommon("friday"),
          tCommon("saturday"),
          tCommon("sunday"),
        ],
      };
    } else {
      // 桌面版：計算當月範圍
      const monthStart = new Date(year, month, 1);
      const monthEnd = new Date(year, month + 1, 0);

      return {
        startDate: formatDate(monthStart),
        endDate: formatDate(monthEnd),
        days: monthEnd.getDate(),
        periodType: "month" as const,
      };
    }
  }, [selectedDate, isMobile]);

  // 取得當月所有現金流資料
  const { data: cashFlows, isLoading } = useCashFlows(
    {
      start_date: startDate,
      end_date: endDate,
    },
    {
      staleTime: 0,
    }
  );

  // 按日期分組，計算每日收入和支出
  const chartData = useMemo(() => {
    if (!cashFlows) return [];

    if (periodType === "week") {
      // 週模式：按星期幾分組（週一到週日）
      const weekData: Array<{
        day: string;
        income: number;
        expense: number;
      }> = [
        { day: "週一", income: 0, expense: 0 },
        { day: "週二", income: 0, expense: 0 },
        { day: "週三", income: 0, expense: 0 },
        { day: "週四", income: 0, expense: 0 },
        { day: "週五", income: 0, expense: 0 },
        { day: "週六", income: 0, expense: 0 },
        { day: "週日", income: 0, expense: 0 },
      ];

      cashFlows.forEach((cf) => {
        const date = new Date(cf.date);
        let dayOfWeek = date.getDay(); // 0 = 週日, 1 = 週一, ..., 6 = 週六
        // 調整為週一 = 0, 週日 = 6
        dayOfWeek = dayOfWeek === 0 ? 6 : dayOfWeek - 1;

        if (cf.type === CashFlowType.INCOME) {
          weekData[dayOfWeek].income += cf.amount;
        } else if (cf.type === CashFlowType.EXPENSE) {
          weekData[dayOfWeek].expense += cf.amount;
        }
      });

      return weekData;
    } else {
      // 月模式：按日期分組（1 到當月天數）
      const dailyData: Record<
        number,
        { day: number; income: number; expense: number }
      > = {};

      for (let day = 1; day <= days; day++) {
        dailyData[day] = { day, income: 0, expense: 0 };
      }

      cashFlows.forEach((cf) => {
        const date = new Date(cf.date);
        const day = date.getDate();

        if (cf.type === CashFlowType.INCOME) {
          dailyData[day].income += cf.amount;
        } else if (cf.type === CashFlowType.EXPENSE) {
          dailyData[day].expense += cf.amount;
        }
      });

      return Object.values(dailyData);
    }
  }, [cashFlows, days, periodType]);

  // 圖表配置
  const chartConfig = {
    income: {
      label: t("income"),
      color: "#22c55e", // 綠色
    },
    expense: {
      label: t("expense"),
      color: "#ef4444", // 紅色
    },
  } satisfies ChartConfig;

  // Loading 狀態
  if (isLoading) {
    return (
      <div className="w-full">
        <Skeleton className="h-[200px] w-full" />
      </div>
    );
  }

  // 空資料狀態
  if (!cashFlows || cashFlows.length === 0) {
    return (
      <div className="w-full h-[200px] flex items-center justify-center text-muted-foreground">
        <p>
          {periodType === "week" ? t("weekly") : t("monthly")}
          {t("noRecords")}
        </p>
      </div>
    );
  }

  return (
    <ChartContainer config={chartConfig} className="h-[200px] w-full">
      <BarChart accessibilityLayer data={chartData}>
        <CartesianGrid vertical={false} />
        <XAxis
          dataKey="day"
          tickLine={false}
          tickMargin={10}
          axisLine={false}
          tickFormatter={(value) =>
            periodType === "week" ? value : `${value}日`
          }
        />
        <YAxis
          tickLine={false}
          axisLine={false}
          tickFormatter={(value) => `$${formatAmount(value)}`}
        />
        <ChartTooltip
          content={<ChartTooltipContent />}
          labelFormatter={(value) =>
            periodType === "week" ? value : `${value} 日`
          }
        />
        <Bar dataKey="income" fill="var(--color-income)" radius={4} />
        <Bar dataKey="expense" fill="var(--color-expense)" radius={4} />
      </BarChart>
    </ChartContainer>
  );
}
