/**
 * 資產趨勢圖表元件
 * 顯示總資產及各類資產的價值變化趨勢
 */

'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts';
import { ChartDataPoint } from '@/lib/mock-data';

interface AssetTrendChartProps {
  data: ChartDataPoint[];
}

export function AssetTrendChart({ data }: AssetTrendChartProps) {
  // 使用 state 來延遲渲染圖表,避免 SSR 問題
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  // 自訂 Tooltip 內容
  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      return (
        <div className="rounded-lg border bg-background p-2 shadow-sm">
          <div className="grid gap-2">
            <div className="flex flex-col">
              <span className="text-[0.70rem] uppercase text-muted-foreground">
                {label}
              </span>
            </div>
            {payload.map((entry: any, index: number) => (
              <div key={index} className="flex items-center gap-2">
                <div
                  className="h-2 w-2 rounded-full"
                  style={{ backgroundColor: entry.color }}
                />
                <span className="text-sm font-medium">
                  {entry.name}: NT$ {entry.value.toLocaleString()}
                </span>
              </div>
            ))}
          </div>
        </div>
      );
    }
    return null;
  };

  return (
    <Card className="@container/card">
      <CardHeader>
        <CardTitle>資產價值趨勢</CardTitle>
        <CardDescription>最近 30 天的資產變化</CardDescription>
      </CardHeader>
      <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6">
        <div className="h-[250px] w-full">
          {mounted ? (
            <ResponsiveContainer width="100%" height="100%">
              <LineChart
                data={data}
                margin={{
                  top: 5,
                  right: 10,
                  left: 10,
                  bottom: 0,
                }}
              >
                <CartesianGrid strokeDasharray="3 3" vertical={false} className="stroke-muted" />
                <XAxis
                  dataKey="date"
                  tickLine={false}
                  axisLine={false}
                  tickMargin={8}
                  className="text-xs"
                />
                <YAxis
                  tickLine={false}
                  axisLine={false}
                  tickMargin={8}
                  className="text-xs"
                  tickFormatter={(value) => `${(value / 1000).toFixed(0)}K`}
                />
                <Tooltip content={<CustomTooltip />} />
                <Line
                  type="monotone"
                  dataKey="total"
                  name="總資產"
                  stroke="#111827"
                  strokeWidth={2}
                  dot={false}
                />
                <Line
                  type="monotone"
                  dataKey="twStock"
                  name="台股"
                  stroke="#3b82f6"
                  strokeWidth={1.5}
                  dot={false}
                  strokeOpacity={0.7}
                />
                <Line
                  type="monotone"
                  dataKey="usStock"
                  name="美股"
                  stroke="#10b981"
                  strokeWidth={1.5}
                  dot={false}
                  strokeOpacity={0.7}
                />
                <Line
                  type="monotone"
                  dataKey="crypto"
                  name="加密貨幣"
                  stroke="#f59e0b"
                  strokeWidth={1.5}
                  dot={false}
                  strokeOpacity={0.7}
                />
              </LineChart>
            </ResponsiveContainer>
          ) : (
            <div className="flex h-full items-center justify-center text-sm text-muted-foreground">
              載入中...
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

