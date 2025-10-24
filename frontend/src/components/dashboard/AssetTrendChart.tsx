/**
 * 資產趨勢圖表元件
 * 顯示總資產及各類資產的價值變化趨勢
 */

"use client";

import { useEffect, useState, useMemo } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { AlertCircle } from "lucide-react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { useAssetTrend } from "@/hooks";
import { AssetTrendData } from "@/types/asset-snapshot";

export function AssetTrendChart() {
  // 使用 state 來延遲渲染圖表,避免 SSR 問題
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  // 取得各類資產的趨勢資料
  const {
    data: totalData,
    isLoading: totalLoading,
    error: totalError,
  } = useAssetTrend("total", 30);
  const { data: twStockData, isLoading: twStockLoading } = useAssetTrend(
    "tw-stock",
    30
  );
  const { data: usStockData, isLoading: usStockLoading } = useAssetTrend(
    "us-stock",
    30
  );
  const { data: cryptoData, isLoading: cryptoLoading } = useAssetTrend(
    "crypto",
    30
  );

  // 合併所有資料為圖表格式
  const chartData = useMemo(() => {
    if (!totalData) return [];

    // 建立日期對應的資料 map
    const dataMap = new Map<string, AssetTrendData>();

    // 加入總資產資料
    totalData.forEach((snapshot) => {
      const date = new Date(snapshot.snapshot_date).toLocaleDateString(
        "zh-TW",
        {
          month: "numeric",
          day: "numeric",
        }
      );
      dataMap.set(snapshot.snapshot_date, {
        date,
        total: snapshot.value_twd,
      });
    });

    // 加入台股資料
    twStockData?.forEach((snapshot) => {
      const existing = dataMap.get(snapshot.snapshot_date);
      if (existing) {
        existing.twStock = snapshot.value_twd;
      }
    });

    // 加入美股資料
    usStockData?.forEach((snapshot) => {
      const existing = dataMap.get(snapshot.snapshot_date);
      if (existing) {
        existing.usStock = snapshot.value_twd;
      }
    });

    // 加入加密貨幣資料
    cryptoData?.forEach((snapshot) => {
      const existing = dataMap.get(snapshot.snapshot_date);
      if (existing) {
        existing.crypto = snapshot.value_twd;
      }
    });

    // 轉換為陣列並排序
    return Array.from(dataMap.values()).sort((a, b) => {
      const dateA = new Date(a.date);
      const dateB = new Date(b.date);
      return dateA.getTime() - dateB.getTime();
    });
  }, [totalData, twStockData, usStockData, cryptoData]);

  const isLoading =
    totalLoading || twStockLoading || usStockLoading || cryptoLoading;

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
        {/* Loading 狀態 */}
        {isLoading && (
          <div className="h-[300px] w-full">
            <Skeleton className="h-full w-full" />
          </div>
        )}

        {/* Error 狀態 */}
        {!isLoading && totalError && (
          <div className="h-[300px] w-full flex items-center justify-center">
            <div className="flex items-center gap-2 text-red-600">
              <AlertCircle className="h-5 w-5" />
              <p className="text-sm">載入資料失敗：{totalError.message}</p>
            </div>
          </div>
        )}

        {/* 空資料狀態 */}
        {!isLoading && !totalError && chartData.length === 0 && (
          <div className="h-[300px] w-full flex items-center justify-center">
            <p className="text-sm text-muted-foreground">暫無資料</p>
          </div>
        )}

        {/* 圖表 */}
        {!isLoading && !totalError && chartData.length > 0 && (
          <div className="h-[300px] w-full">
            {mounted ? (
              <ResponsiveContainer width="100%" height="100%">
                <LineChart
                  data={chartData}
                  margin={{
                    top: 5,
                    right: 10,
                    left: 10,
                    bottom: 0,
                  }}
                >
                  <CartesianGrid
                    strokeDasharray="3 3"
                    vertical={false}
                    className="stroke-muted"
                  />
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
        )}
      </CardContent>
    </Card>
  );
}
