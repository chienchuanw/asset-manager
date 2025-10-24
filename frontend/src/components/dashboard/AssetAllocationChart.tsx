/**
 * 資產配置圓餅圖元件
 * 顯示各類資產的佔比分布
 */

"use client";

import { useMemo } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from "recharts";
import { AssetType, getAssetTypeLabel } from "@/types/transaction";

interface AssetAllocationData {
  name: string;
  value: number;
}

interface AssetAllocationChartProps {
  data: AssetAllocationData[];
}

// 資產類型顏色對應
const ASSET_COLORS: Record<string, string> = {
  "tw-stock": "#3b82f6", // 藍色
  "us-stock": "#10b981", // 綠色
  crypto: "#f59e0b", // 橘色
};

export function AssetAllocationChart({ data }: AssetAllocationChartProps) {
  // 計算總值和百分比
  const chartData = useMemo(() => {
    const total = data.reduce((sum, item) => sum + item.value, 0);
    return data.map((item) => ({
      name: getAssetTypeLabel(item.name as AssetType),
      assetType: item.name,
      value: item.value,
      percentage: total > 0 ? (item.value / total) * 100 : 0,
      color: ASSET_COLORS[item.name] || "#6b7280",
    }));
  }, [data]);

  // 自訂 Tooltip
  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div className="rounded-lg border bg-background p-2 shadow-sm">
          <div className="grid gap-2">
            <div className="flex items-center gap-2">
              <div
                className="h-2 w-2 rounded-full"
                style={{ backgroundColor: data.color }}
              />
              <span className="text-sm font-medium">{data.name}</span>
            </div>
            <div className="text-sm text-muted-foreground">
              TWD{" "}
              {data.value.toLocaleString("zh-TW", { maximumFractionDigits: 0 })}
            </div>
            <div className="text-sm font-medium">
              {data.percentage.toFixed(1)}%
            </div>
          </div>
        </div>
      );
    }
    return null;
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>資產配置</CardTitle>
        <CardDescription>各類資產佔比分布</CardDescription>
      </CardHeader>
      <CardContent>
        {chartData.length === 0 ? (
          <div className="h-[300px] flex items-center justify-center text-muted-foreground">
            暫無資料
          </div>
        ) : (
          <>
            <div className="h-[200px]">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={chartData}
                    cx="50%"
                    cy="50%"
                    labelLine={false}
                    outerRadius={90}
                    fill="#8884d8"
                    dataKey="value"
                  >
                    {chartData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.color} />
                    ))}
                  </Pie>
                  <Tooltip content={<CustomTooltip />} />
                </PieChart>
              </ResponsiveContainer>
            </div>

            {/* 圖例列表 */}
            <div className="mt-6 space-y-3">
              {chartData.map((item) => (
                <div
                  key={item.assetType}
                  className="flex items-center justify-between text-sm"
                >
                  <div className="flex items-center gap-2">
                    <div
                      className="h-2 w-2 rounded-full"
                      style={{ backgroundColor: item.color }}
                    ></div>
                    <span className="text-muted-foreground">{item.name}</span>
                  </div>
                  <div className="flex items-center gap-3">
                    <span className="font-medium tabular-nums">
                      {item.value.toLocaleString("zh-TW", {
                        maximumFractionDigits: 0,
                      })}
                    </span>
                    <span className="text-muted-foreground w-12 text-right tabular-nums">
                      {item.percentage.toFixed(1)}%
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}
