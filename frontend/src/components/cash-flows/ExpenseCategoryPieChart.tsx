/**
 * 支出分類圓餅圖元件
 * 顯示指定期間內支出的各分類佔比
 */

"use client";

import { useMemo } from "react";
import { Pie, PieChart } from "recharts";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
} from "@/components/ui/chart";
import { CashFlow, CashFlowType, formatAmount } from "@/types/cash-flow";

interface ExpenseCategoryPieChartProps {
  cashFlows: CashFlow[];
  period: "week" | "month";
}

/**
 * 支出分類圓餅圖
 *
 * 功能說明：
 * 1. 篩選出所有支出類型的現金流記錄
 * 2. 按照分類名稱分組並計算總金額
 * 3. 只顯示前 8 名分類，其餘合併為「其他」
 * 4. 使用圓餅圖視覺化呈現各分類佔比
 */
export function ExpenseCategoryPieChart({
  cashFlows,
  period,
}: ExpenseCategoryPieChartProps) {
  // 處理圓餅圖資料
  const { chartData, chartConfig, totalExpense } = useMemo(() => {
    // 步驟 1: 篩選出支出類型的記錄
    const expenses = cashFlows.filter((cf) => cf.type === CashFlowType.EXPENSE);

    // 如果沒有支出記錄，直接返回空資料
    if (expenses.length === 0) {
      return {
        chartData: [],
        chartConfig: { amount: { label: "金額" } },
        totalExpense: 0,
      };
    }

    // 步驟 2: 按分類分組並計算總額
    const categoryMap = new Map<string, number>();
    let total = 0;

    expenses.forEach((cf) => {
      const categoryName = cf.category?.name || "未分類";
      const currentAmount = categoryMap.get(categoryName) || 0;
      categoryMap.set(categoryName, currentAmount + cf.amount);
      total += cf.amount;
    });

    // 步驟 3: 轉換為陣列並按金額排序（由大到小）
    const sortedCategories = Array.from(categoryMap.entries())
      .map(([name, amount]) => ({ name, amount }))
      .sort((a, b) => b.amount - a.amount);

    // 步驟 4: 取前 8 名，其餘合併為「其他」
    const topCategories = sortedCategories.slice(0, 8);
    const otherCategories = sortedCategories.slice(8);

    // 如果有其他分類，計算總額並加入
    if (otherCategories.length > 0) {
      const otherTotal = otherCategories.reduce(
        (sum, cat) => sum + cat.amount,
        0
      );
      topCategories.push({ name: "其他", amount: otherTotal });
    }

    // 步驟 5: 轉換成圓餅圖需要的格式，並同時建立配置
    const config: ChartConfig = {
      amount: {
        label: "金額",
      },
    };

    const data = topCategories.map((cat, index) => {
      const colorIndex = (index % 5) + 1; // 循環使用 chart-1 到 chart-5
      const categoryKey = cat.name; // 使用分類名稱作為 key
      const percentage = ((cat.amount / total) * 100).toFixed(1); // 計算百分比

      // 為每個分類建立配置，包含金額和百分比
      config[categoryKey] = {
        label: `${cat.name} $${formatAmount(cat.amount)} (${percentage}%)`,
        color: `var(--chart-${colorIndex})`,
      };

      return {
        category: cat.name,
        amount: cat.amount,
        fill: `var(--color-${categoryKey})`,
      };
    });

    return { chartData: data, chartConfig: config, totalExpense: total };
  }, [cashFlows]);

  // 取得期間標題
  const periodLabel = period === "week" ? "本週" : "本月";

  return (
    <Card className="flex flex-col">
      <CardHeader className="items-center pb-0">
        <CardTitle>支出分類分析</CardTitle>
        <CardDescription>{periodLabel}支出分類佔比</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 pb-0">
        {chartData.length === 0 ? (
          // 空資料狀態
          <div className="flex h-[300px] items-center justify-center text-muted-foreground">
            {periodLabel}暫無支出記錄
          </div>
        ) : (
          <ChartContainer
            config={chartConfig}
            className="mx-auto aspect-square max-h-[300px]"
          >
            <PieChart>
              <Pie data={chartData} dataKey="amount" nameKey="category" />
              <ChartLegend
                content={<ChartLegendContent nameKey="category" />}
                className="-translate-y-2 flex-wrap gap-2 *:basis-1/4 *:justify-center"
              />
            </PieChart>
          </ChartContainer>
        )}
      </CardContent>
      {chartData.length > 0 && (
        <div className="flex-col gap-2 text-sm px-6 pb-4">
          <div className="text-muted-foreground text-center">
            {periodLabel}總支出：${formatAmount(totalExpense)}
          </div>
        </div>
      )}
    </Card>
  );
}
