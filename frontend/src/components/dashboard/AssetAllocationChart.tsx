/**
 * 資產配置圓餅圖元件
 * 顯示各類資產的佔比分布
 */

'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts';
import { AssetAllocation } from '@/lib/mock-data';

interface AssetAllocationChartProps {
  data: AssetAllocation[];
}

export function AssetAllocationChart({ data }: AssetAllocationChartProps) {
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
              NT$ {data.value.toLocaleString()}
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
        <div className="h-[200px]">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={data}
                cx="50%"
                cy="50%"
                labelLine={false}
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
              >
                {data.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip content={<CustomTooltip />} />
            </PieChart>
          </ResponsiveContainer>
        </div>

        {/* 圖例列表 */}
        <div className="mt-6 space-y-3">
          {data.map((item) => (
            <div key={item.assetType} className="flex items-center justify-between text-sm">
              <div className="flex items-center gap-2">
                <div
                  className="h-2 w-2 rounded-full"
                  style={{ backgroundColor: item.color }}
                ></div>
                <span className="text-muted-foreground">{item.name}</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="font-medium tabular-nums">
                  {item.value.toLocaleString()}
                </span>
                <span className="text-muted-foreground w-12 text-right tabular-nums">
                  {item.percentage.toFixed(1)}%
                </span>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

