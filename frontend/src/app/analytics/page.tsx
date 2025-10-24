/**
 * 分析報表頁面
 * 顯示已實現損益分析、績效分析等報表
 */

"use client";

import { useState } from "react";
import { AppLayout } from "@/components/layout/AppLayout";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from "recharts";
import { TrendingUp, TrendingDown, Loader2, AlertCircle } from "lucide-react";
import { useAnalytics } from "@/hooks/useAnalytics";
import { TimeRange } from "@/types/analytics";
import {
  formatCurrency,
  formatPercentage,
  isPositive,
} from "@/types/analytics";
import { getAssetTypeLabel } from "@/types/transaction";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

export default function AnalyticsPage() {
  const [timeRange, setTimeRange] = useState<TimeRange>("month");

  // 使用 Analytics Hook 取得資料
  const { summary, performance, topAssets, isLoading, isError, error } =
    useAnalytics(timeRange, 10);

  return (
    <AppLayout>
      {/* Main Content */}
      <main className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="flex flex-col gap-6">
          {/* 時間範圍選擇 */}
          <Tabs
            value={timeRange}
            onValueChange={(value) => setTimeRange(value as TimeRange)}
          >
            <TabsList>
              <TabsTrigger value="week">本週</TabsTrigger>
              <TabsTrigger value="month">本月</TabsTrigger>
              <TabsTrigger value="quarter">本季</TabsTrigger>
              <TabsTrigger value="year">本年</TabsTrigger>
              <TabsTrigger value="all">全部</TabsTrigger>
            </TabsList>
          </Tabs>

          {/* Loading 狀態 */}
          {isLoading && (
            <div className="flex items-center justify-center py-12">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
              <span className="ml-2 text-muted-foreground">載入中...</span>
            </div>
          )}

          {/* Error 狀態 */}
          {isError && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>載入失敗</AlertTitle>
              <AlertDescription>
                {error?.message || "無法載入分析資料，請稍後再試"}
              </AlertDescription>
            </Alert>
          )}

          {/* 資料顯示 */}
          {!isLoading && !isError && summary.data && (
            <>
              {/* 績效摘要卡片 */}
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
                <Card>
                  <CardHeader className="pb-2">
                    <CardDescription>總成本基礎</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold tabular-nums">
                      {formatCurrency(
                        summary.data.total_cost_basis,
                        summary.data.currency
                      )}
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">
                      賣出交易成本
                    </p>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader className="pb-2">
                    <CardDescription>總賣出金額</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold tabular-nums">
                      {formatCurrency(
                        summary.data.total_sell_amount,
                        summary.data.currency
                      )}
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">
                      實際賣出收入
                    </p>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader className="pb-2">
                    <CardDescription>已實現損益</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div
                      className={`text-2xl font-bold tabular-nums ${
                        isPositive(summary.data.total_realized_pl)
                          ? "text-green-600"
                          : "text-red-600"
                      }`}
                    >
                      {formatCurrency(
                        summary.data.total_realized_pl,
                        summary.data.currency
                      )}
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">
                      實際獲利/虧損
                    </p>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader className="pb-2">
                    <CardDescription>已實現報酬率</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div
                      className={`text-2xl font-bold tabular-nums ${
                        isPositive(summary.data.total_realized_pl_pct)
                          ? "text-green-600"
                          : "text-red-600"
                      }`}
                    >
                      {formatPercentage(summary.data.total_realized_pl_pct)}
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">
                      <Badge
                        variant="outline"
                        className={
                          isPositive(summary.data.total_realized_pl_pct)
                            ? "bg-green-100 text-green-800"
                            : "bg-red-100 text-red-800"
                        }
                      >
                        {isPositive(summary.data.total_realized_pl_pct) ? (
                          <TrendingUp className="h-3 w-3 mr-1" />
                        ) : (
                          <TrendingDown className="h-3 w-3 mr-1" />
                        )}
                        {summary.data.transaction_count} 筆交易
                      </Badge>
                    </p>
                  </CardContent>
                </Card>
              </div>

              {/* 圖表區域 */}
              <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
                {/* 各資產報酬率長條圖 */}
                <Card>
                  <CardHeader>
                    <CardTitle>各資產類型績效</CardTitle>
                    <CardDescription>
                      不同資產類別的已實現損益比較
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    {performance.data && performance.data.length > 0 ? (
                      <ResponsiveContainer width="100%" height={300}>
                        <BarChart data={performance.data}>
                          <CartesianGrid
                            strokeDasharray="3 3"
                            className="stroke-muted"
                          />
                          <XAxis
                            dataKey="name"
                            className="text-xs"
                            tick={{ fill: "hsl(var(--muted-foreground))" }}
                          />
                          <YAxis
                            className="text-xs"
                            tick={{ fill: "hsl(var(--muted-foreground))" }}
                            label={{
                              value: "報酬率 (%)",
                              angle: -90,
                              position: "insideLeft",
                            }}
                          />
                          <Tooltip
                            contentStyle={{
                              backgroundColor: "hsl(var(--background))",
                              border: "1px solid hsl(var(--border))",
                              borderRadius: "6px",
                            }}
                            formatter={(value: number) => [
                              `${value.toFixed(2)}%`,
                              "報酬率",
                            ]}
                          />
                          <Bar dataKey="realized_pl_pct" radius={[4, 4, 0, 0]}>
                            {performance.data.map((entry, index) => (
                              <Cell
                                key={`cell-${index}`}
                                fill={
                                  entry.realized_pl_pct >= 0
                                    ? "#10b981"
                                    : "#ef4444"
                                }
                              />
                            ))}
                          </Bar>
                        </BarChart>
                      </ResponsiveContainer>
                    ) : (
                      <div className="flex items-center justify-center h-[300px] text-muted-foreground">
                        目前沒有績效資料
                      </div>
                    )}
                  </CardContent>
                </Card>

                {/* 各資產類型損益統計 */}
                <Card>
                  <CardHeader>
                    <CardTitle>各資產類型損益</CardTitle>
                    <CardDescription>各資產類別的已實現損益統計</CardDescription>
                  </CardHeader>
                  <CardContent>
                    {performance.data && performance.data.length > 0 ? (
                      <div className="space-y-4">
                        {performance.data.map((item) => (
                          <div
                            key={item.asset_type}
                            className="flex items-center justify-between"
                          >
                            <div className="flex items-center gap-3">
                              <Badge variant="outline">{item.name}</Badge>
                              <span className="text-sm text-muted-foreground">
                                {item.transaction_count} 筆
                              </span>
                            </div>
                            <div className="flex items-center gap-4">
                              <span
                                className={`text-sm font-medium tabular-nums ${
                                  isPositive(item.realized_pl)
                                    ? "text-green-600"
                                    : "text-red-600"
                                }`}
                              >
                                {formatPercentage(item.realized_pl_pct)}
                              </span>
                              <span
                                className={`text-sm font-medium tabular-nums w-32 text-right ${
                                  isPositive(item.realized_pl)
                                    ? "text-green-600"
                                    : "text-red-600"
                                }`}
                              >
                                {formatCurrency(item.realized_pl, "TWD")}
                              </span>
                            </div>
                          </div>
                        ))}
                      </div>
                    ) : (
                      <div className="flex items-center justify-center h-[300px] text-muted-foreground">
                        目前沒有績效資料
                      </div>
                    )}
                  </CardContent>
                </Card>
              </div>

              {/* Top 資產表格 */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <TrendingUp className="h-5 w-5 text-green-600" />
                    Top 表現資產
                  </CardTitle>
                  <CardDescription>已實現損益最佳的投資標的</CardDescription>
                </CardHeader>
                <CardContent>
                  {topAssets.data && topAssets.data.length > 0 ? (
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead>代碼/名稱</TableHead>
                          <TableHead>類別</TableHead>
                          <TableHead className="text-right">成本</TableHead>
                          <TableHead className="text-right">賣出金額</TableHead>
                          <TableHead className="text-right">損益</TableHead>
                          <TableHead className="text-right">報酬率</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {topAssets.data.map((asset) => (
                          <TableRow key={asset.symbol}>
                            <TableCell>
                              <div>
                                <div className="font-medium">{asset.symbol}</div>
                                <div className="text-sm text-muted-foreground">
                                  {asset.name}
                                </div>
                              </div>
                            </TableCell>
                            <TableCell>
                              <Badge variant="outline">
                                {getAssetTypeLabel(asset.asset_type)}
                              </Badge>
                            </TableCell>
                            <TableCell className="text-right tabular-nums">
                              {formatCurrency(asset.cost_basis, "TWD")}
                            </TableCell>
                            <TableCell className="text-right tabular-nums">
                              {formatCurrency(asset.sell_amount, "TWD")}
                            </TableCell>
                            <TableCell
                              className={`text-right font-medium tabular-nums ${
                                isPositive(asset.realized_pl)
                                  ? "text-green-600"
                                  : "text-red-600"
                              }`}
                            >
                              {formatCurrency(asset.realized_pl, "TWD")}
                            </TableCell>
                            <TableCell
                              className={`text-right font-medium tabular-nums ${
                                isPositive(asset.realized_pl_pct)
                                  ? "text-green-600"
                                  : "text-red-600"
                              }`}
                            >
                              {formatPercentage(asset.realized_pl_pct)}
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  ) : (
                    <div className="flex items-center justify-center h-32 text-muted-foreground">
                      目前沒有資產資料
                    </div>
                  )}
                </CardContent>
              </Card>
            </>
          )}
        </div>
      </main>
    </AppLayout>
  );
}
