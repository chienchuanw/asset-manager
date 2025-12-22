/**
 * 分析報表頁面
 * 顯示已實現損益分析、績效分析等報表
 */

"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { AppLayout } from "@/components/layout/AppLayout";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
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
import { TrendingUp, TrendingDown, AlertCircle } from "lucide-react";
import { Loading } from "@/components/ui/loading";
import { useAnalytics } from "@/hooks/useAnalytics";
import { useUnrealizedAnalytics } from "@/hooks/useUnrealizedAnalytics";
import { TimeRange } from "@/types/analytics";
import {
  formatCurrency,
  formatPercentage,
  isPositive,
} from "@/types/analytics";
import { getAssetTypeLabel } from "@/types/transaction";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

export default function AnalyticsPage() {
  // i18n hooks
  const t = useTranslations("analytics");
  const tCommon = useTranslations("common");
  const tErrors = useTranslations("errors");

  const [timeRange, setTimeRange] = useState<TimeRange>("month");
  const [analysisType, setAnalysisType] = useState<"realized" | "unrealized">(
    "unrealized"
  );

  // 使用 Analytics Hook 取得已實現損益資料
  const realizedData = useAnalytics(timeRange, 10);

  // 使用 Unrealized Analytics Hook 取得未實現損益資料
  const unrealizedData = useUnrealizedAnalytics(10);

  return (
    <AppLayout title={t("title")} description={t("description")}>
      {/* Main Content */}
      <div className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="flex flex-col gap-6">
          {/* 分析類型切換 */}
          <Tabs
            value={analysisType}
            onValueChange={(value) =>
              setAnalysisType(value as "realized" | "unrealized")
            }
          >
            <TabsList>
              <TabsTrigger value="unrealized">{t("unrealizedPL")}</TabsTrigger>
              <TabsTrigger value="realized">{t("realizedPL")}</TabsTrigger>
            </TabsList>

            {/* 已實現損益 Tab */}
            <TabsContent value="realized" className="space-y-6 mt-6">
              {/* 時間範圍選擇 */}
              <Tabs
                value={timeRange}
                onValueChange={(value) => setTimeRange(value as TimeRange)}
              >
                <TabsList>
                  <TabsTrigger value="week">{t("thisWeek")}</TabsTrigger>
                  <TabsTrigger value="month">{t("thisMonth")}</TabsTrigger>
                  <TabsTrigger value="quarter">{t("thisQuarter")}</TabsTrigger>
                  <TabsTrigger value="year">{t("thisYear")}</TabsTrigger>
                  <TabsTrigger value="all">{tCommon("all")}</TabsTrigger>
                </TabsList>
              </Tabs>

              {/* Loading 狀態 */}
              {realizedData.isLoading && (
                <Loading variant="inline" size="md" text={tCommon("loading")} />
              )}

              {/* Error 狀態 */}
              {realizedData.isError && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertTitle>{tErrors("loadFailed")}</AlertTitle>
                  <AlertDescription>
                    {realizedData.error?.message || t("loadAnalyticsError")}
                  </AlertDescription>
                </Alert>
              )}

              {/* 資料顯示 */}
              {!realizedData.isLoading &&
                !realizedData.isError &&
                realizedData.summary.data && (
                  <>
                    {/* 績效摘要卡片 */}
                    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
                      <Card>
                        <CardHeader className="pb-2">
                          <CardDescription>
                            {t("totalCostBasis")}
                          </CardDescription>
                        </CardHeader>
                        <CardContent>
                          <div className="text-2xl font-bold tabular-nums">
                            {formatCurrency(
                              realizedData.summary.data.total_cost_basis,
                              realizedData.summary.data.currency
                            )}
                          </div>
                          <p className="text-xs text-muted-foreground mt-1">
                            {t("sellTransactionCost")}
                          </p>
                        </CardContent>
                      </Card>

                      <Card>
                        <CardHeader className="pb-2">
                          <CardDescription>
                            {t("totalSellAmount")}
                          </CardDescription>
                        </CardHeader>
                        <CardContent>
                          <div className="text-2xl font-bold tabular-nums">
                            {formatCurrency(
                              realizedData.summary.data.total_sell_amount,
                              realizedData.summary.data.currency
                            )}
                          </div>
                          <p className="text-xs text-muted-foreground mt-1">
                            {t("actualSellIncome")}
                          </p>
                        </CardContent>
                      </Card>

                      <Card>
                        <CardHeader className="pb-2">
                          <CardDescription>{t("realizedPL")}</CardDescription>
                        </CardHeader>
                        <CardContent>
                          <div
                            className={`text-2xl font-bold tabular-nums ${
                              isPositive(
                                realizedData.summary.data.total_realized_pl
                              )
                                ? "text-green-600"
                                : "text-red-600"
                            }`}
                          >
                            {formatCurrency(
                              realizedData.summary.data.total_realized_pl,
                              realizedData.summary.data.currency
                            )}
                          </div>
                          <p className="text-xs text-muted-foreground mt-1">
                            {t("actualProfitLoss")}
                          </p>
                        </CardContent>
                      </Card>

                      <Card>
                        <CardHeader className="pb-2">
                          <CardDescription>
                            {t("realizedReturnRate")}
                          </CardDescription>
                        </CardHeader>
                        <CardContent>
                          <div
                            className={`text-2xl font-bold tabular-nums ${
                              isPositive(
                                realizedData.summary.data.total_realized_pl_pct
                              )
                                ? "text-green-600"
                                : "text-red-600"
                            }`}
                          >
                            {formatPercentage(
                              realizedData.summary.data.total_realized_pl_pct
                            )}
                          </div>
                          <p className="text-xs text-muted-foreground mt-1">
                            <Badge
                              variant="outline"
                              className={
                                isPositive(
                                  realizedData.summary.data
                                    .total_realized_pl_pct
                                )
                                  ? "bg-red-100 text-red-800"
                                  : "bg-green-100 text-green-800"
                              }
                            >
                              {isPositive(
                                realizedData.summary.data.total_realized_pl_pct
                              ) ? (
                                <TrendingUp className="h-3 w-3 mr-1" />
                              ) : (
                                <TrendingDown className="h-3 w-3 mr-1" />
                              )}
                              {t("transactionCount", {
                                count:
                                  realizedData.summary.data.transaction_count,
                              })}
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
                          <CardTitle>{t("assetTypePerformance")}</CardTitle>
                          <CardDescription>
                            {t("realizedPLComparison")}
                          </CardDescription>
                        </CardHeader>
                        <CardContent>
                          {realizedData.performance.data &&
                          realizedData.performance.data.length > 0 ? (
                            <ResponsiveContainer width="100%" height={300}>
                              <BarChart data={realizedData.performance.data}>
                                <CartesianGrid
                                  strokeDasharray="3 3"
                                  className="stroke-muted"
                                />
                                <XAxis
                                  dataKey="name"
                                  className="text-xs"
                                  tick={{
                                    fill: "hsl(var(--muted-foreground))",
                                  }}
                                />
                                <YAxis
                                  className="text-xs"
                                  tick={{
                                    fill: "hsl(var(--muted-foreground))",
                                  }}
                                  label={{
                                    value: t("returnRatePercent"),
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
                                    t("returnRate"),
                                  ]}
                                />
                                <Bar
                                  dataKey="realized_pl_pct"
                                  radius={[4, 4, 0, 0]}
                                >
                                  {realizedData.performance.data.map(
                                    (entry, index) => (
                                      <Cell
                                        key={`cell-${index}`}
                                        fill={
                                          entry.realized_pl_pct >= 0
                                            ? "#ef4444"
                                            : "#10b981"
                                        }
                                      />
                                    )
                                  )}
                                </Bar>
                              </BarChart>
                            </ResponsiveContainer>
                          ) : (
                            <div className="flex items-center justify-center h-[300px] text-muted-foreground">
                              {t("noPerformanceData")}
                            </div>
                          )}
                        </CardContent>
                      </Card>

                      {/* 各資產類型損益統計 */}
                      <Card>
                        <CardHeader>
                          <CardTitle>{t("assetTypePL")}</CardTitle>
                          <CardDescription>
                            {t("realizedPLStats")}
                          </CardDescription>
                        </CardHeader>
                        <CardContent>
                          {realizedData.performance.data &&
                          realizedData.performance.data.length > 0 ? (
                            <div className="space-y-4">
                              {realizedData.performance.data.map((item) => (
                                <div
                                  key={item.asset_type}
                                  className="flex items-center justify-between"
                                >
                                  <div className="flex items-center gap-3">
                                    <Badge variant="outline">{item.name}</Badge>
                                    <span className="text-sm text-muted-foreground">
                                      {t("transactionCountShort", {
                                        count: item.transaction_count,
                                      })}
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
                              {t("noPerformanceData")}
                            </div>
                          )}
                        </CardContent>
                      </Card>
                    </div>

                    {/* Top 資產表格 */}
                    <Card>
                      <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                          <TrendingUp className="h-5 w-5 text-red-600" />
                          {t("topPerformingAssets")}
                        </CardTitle>
                        <CardDescription>
                          {t("bestRealizedPLAssets")}
                        </CardDescription>
                      </CardHeader>
                      <CardContent>
                        {realizedData.topAssets.data &&
                        realizedData.topAssets.data.length > 0 ? (
                          <Table>
                            <TableHeader>
                              <TableRow>
                                <TableHead>{t("symbolName")}</TableHead>
                                <TableHead>{t("category")}</TableHead>
                                <TableHead className="text-right">
                                  {t("cost")}
                                </TableHead>
                                <TableHead className="text-right">
                                  {t("sellAmount")}
                                </TableHead>
                                <TableHead className="text-right">
                                  {t("profitLoss")}
                                </TableHead>
                                <TableHead className="text-right">
                                  {t("returnRate")}
                                </TableHead>
                              </TableRow>
                            </TableHeader>
                            <TableBody>
                              {realizedData.topAssets.data.map((asset) => (
                                <TableRow key={asset.symbol}>
                                  <TableCell>
                                    <div>
                                      <div className="font-medium">
                                        {asset.symbol}
                                      </div>
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
                            {t("noAssetData")}
                          </div>
                        )}
                      </CardContent>
                    </Card>
                  </>
                )}
            </TabsContent>

            {/* 未實現損益 Tab */}
            <TabsContent value="unrealized" className="space-y-6 mt-6">
              {/* Loading 狀態 */}
              {unrealizedData.isLoading && (
                <Loading variant="inline" size="md" text={tCommon("loading")} />
              )}

              {/* Error 狀態 */}
              {unrealizedData.isError && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertTitle>{tErrors("loadFailed")}</AlertTitle>
                  <AlertDescription>
                    {unrealizedData.error?.message || t("loadUnrealizedError")}
                  </AlertDescription>
                </Alert>
              )}

              {/* 資料顯示 */}
              {!unrealizedData.isLoading &&
                !unrealizedData.isError &&
                unrealizedData.summary.data && (
                  <>
                    {/* 績效摘要卡片 */}
                    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
                      <Card>
                        <CardHeader className="pb-2">
                          <CardDescription>{t("totalCost")}</CardDescription>
                        </CardHeader>
                        <CardContent>
                          <div className="text-2xl font-bold tabular-nums">
                            {formatCurrency(
                              unrealizedData.summary.data.total_cost,
                              unrealizedData.summary.data.currency
                            )}
                          </div>
                          <p className="text-xs text-muted-foreground mt-1">
                            {t("holdingCostBasis")}
                          </p>
                        </CardContent>
                      </Card>

                      <Card>
                        <CardHeader className="pb-2">
                          <CardDescription>
                            {t("totalMarketValue")}
                          </CardDescription>
                        </CardHeader>
                        <CardContent>
                          <div className="text-2xl font-bold tabular-nums">
                            {formatCurrency(
                              unrealizedData.summary.data.total_market_value,
                              unrealizedData.summary.data.currency
                            )}
                          </div>
                          <p className="text-xs text-muted-foreground mt-1">
                            {t("currentMarketValue")}
                          </p>
                        </CardContent>
                      </Card>

                      <Card>
                        <CardHeader className="pb-2">
                          <CardDescription>{t("unrealizedPL")}</CardDescription>
                        </CardHeader>
                        <CardContent>
                          <div
                            className={`text-2xl font-bold tabular-nums ${
                              isPositive(
                                unrealizedData.summary.data.total_unrealized_pl
                              )
                                ? "text-green-600"
                                : "text-red-600"
                            }`}
                          >
                            {formatCurrency(
                              unrealizedData.summary.data.total_unrealized_pl,
                              unrealizedData.summary.data.currency
                            )}
                          </div>
                          <p className="text-xs text-muted-foreground mt-1">
                            {t("floatingPL")}
                          </p>
                        </CardContent>
                      </Card>

                      <Card>
                        <CardHeader className="pb-2">
                          <CardDescription>
                            {t("unrealizedReturnRate")}
                          </CardDescription>
                        </CardHeader>
                        <CardContent>
                          <div
                            className={`text-2xl font-bold tabular-nums ${
                              isPositive(
                                unrealizedData.summary.data.total_unrealized_pct
                              )
                                ? "text-green-600"
                                : "text-red-600"
                            }`}
                          >
                            {formatPercentage(
                              unrealizedData.summary.data.total_unrealized_pct
                            )}
                          </div>
                          <p className="text-xs text-muted-foreground mt-1">
                            <Badge variant="outline">
                              {t("holdingCount", {
                                count:
                                  unrealizedData.summary.data.holding_count,
                              })}
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
                          <CardTitle>{t("assetTypePerformance")}</CardTitle>
                          <CardDescription>
                            {t("unrealizedPLComparison")}
                          </CardDescription>
                        </CardHeader>
                        <CardContent>
                          {unrealizedData.performance.data &&
                          unrealizedData.performance.data.length > 0 ? (
                            <ResponsiveContainer width="100%" height={300}>
                              <BarChart data={unrealizedData.performance.data}>
                                <CartesianGrid
                                  strokeDasharray="3 3"
                                  className="stroke-muted"
                                />
                                <XAxis
                                  dataKey="name"
                                  className="text-xs"
                                  tick={{
                                    fill: "hsl(var(--muted-foreground))",
                                  }}
                                />
                                <YAxis
                                  className="text-xs"
                                  tick={{
                                    fill: "hsl(var(--muted-foreground))",
                                  }}
                                  label={{
                                    value: t("returnRatePercent"),
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
                                    t("returnRate"),
                                  ]}
                                />
                                <Bar
                                  dataKey="unrealized_pct"
                                  radius={[4, 4, 0, 0]}
                                >
                                  {unrealizedData.performance.data.map(
                                    (entry, index) => (
                                      <Cell
                                        key={`cell-${index}`}
                                        fill={
                                          entry.unrealized_pct >= 0
                                            ? "#ef4444"
                                            : "#10b981"
                                        }
                                      />
                                    )
                                  )}
                                </Bar>
                              </BarChart>
                            </ResponsiveContainer>
                          ) : (
                            <div className="flex items-center justify-center h-[300px] text-muted-foreground">
                              {t("noPerformanceData")}
                            </div>
                          )}
                        </CardContent>
                      </Card>

                      {/* 各資產類型損益統計 */}
                      <Card>
                        <CardHeader>
                          <CardTitle>{t("assetTypePL")}</CardTitle>
                          <CardDescription>
                            {t("unrealizedPLStats")}
                          </CardDescription>
                        </CardHeader>
                        <CardContent>
                          {unrealizedData.performance.data &&
                          unrealizedData.performance.data.length > 0 ? (
                            <div className="space-y-4">
                              {unrealizedData.performance.data.map((item) => (
                                <div
                                  key={item.asset_type}
                                  className="flex items-center justify-between"
                                >
                                  <div className="flex items-center gap-3">
                                    <Badge variant="outline">{item.name}</Badge>
                                    <span className="text-sm text-muted-foreground">
                                      {t("holdingCount", {
                                        count: item.holding_count,
                                      })}
                                    </span>
                                  </div>
                                  <div className="flex items-center gap-4">
                                    <span
                                      className={`text-sm font-medium tabular-nums ${
                                        isPositive(item.unrealized_pl)
                                          ? "text-green-600"
                                          : "text-red-600"
                                      }`}
                                    >
                                      {formatPercentage(item.unrealized_pct)}
                                    </span>
                                    <span
                                      className={`text-sm font-medium tabular-nums w-32 text-right ${
                                        isPositive(item.unrealized_pl)
                                          ? "text-green-600"
                                          : "text-red-600"
                                      }`}
                                    >
                                      {formatCurrency(
                                        item.unrealized_pl,
                                        "TWD"
                                      )}
                                    </span>
                                  </div>
                                </div>
                              ))}
                            </div>
                          ) : (
                            <div className="flex items-center justify-center h-[300px] text-muted-foreground">
                              {t("noPerformanceData")}
                            </div>
                          )}
                        </CardContent>
                      </Card>
                    </div>

                    {/* Top 資產表格 */}
                    <Card>
                      <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                          <TrendingUp className="h-5 w-5 text-red-600" />
                          {t("topPerformingAssets")}
                        </CardTitle>
                        <CardDescription>
                          {t("bestUnrealizedPLAssets")}
                        </CardDescription>
                      </CardHeader>
                      <CardContent>
                        {unrealizedData.topAssets.data &&
                        unrealizedData.topAssets.data.length > 0 ? (
                          <Table>
                            <TableHeader>
                              <TableRow>
                                <TableHead>{t("symbolName")}</TableHead>
                                <TableHead>{t("category")}</TableHead>
                                <TableHead className="text-right">
                                  {t("quantity")}
                                </TableHead>
                                <TableHead className="text-right">
                                  {t("avgCost")}
                                </TableHead>
                                <TableHead className="text-right">
                                  {t("currentPrice")}
                                </TableHead>
                                <TableHead className="text-right">
                                  {t("profitLoss")}
                                </TableHead>
                                <TableHead className="text-right">
                                  {t("returnRate")}
                                </TableHead>
                              </TableRow>
                            </TableHeader>
                            <TableBody>
                              {unrealizedData.topAssets.data.map((asset) => (
                                <TableRow key={asset.symbol}>
                                  <TableCell>
                                    <div>
                                      <div className="font-medium">
                                        {asset.symbol}
                                      </div>
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
                                    {asset.quantity.toLocaleString("zh-TW")}
                                  </TableCell>
                                  <TableCell className="text-right tabular-nums">
                                    {formatCurrency(asset.avg_cost, "TWD")}
                                  </TableCell>
                                  <TableCell className="text-right tabular-nums">
                                    {formatCurrency(asset.current_price, "TWD")}
                                  </TableCell>
                                  <TableCell
                                    className={`text-right font-medium tabular-nums ${
                                      isPositive(asset.unrealized_pl)
                                        ? "text-green-600"
                                        : "text-red-600"
                                    }`}
                                  >
                                    {formatCurrency(asset.unrealized_pl, "TWD")}
                                  </TableCell>
                                  <TableCell
                                    className={`text-right font-medium tabular-nums ${
                                      isPositive(asset.unrealized_pct)
                                        ? "text-green-600"
                                        : "text-red-600"
                                    }`}
                                  >
                                    {formatPercentage(asset.unrealized_pct)}
                                  </TableCell>
                                </TableRow>
                              ))}
                            </TableBody>
                          </Table>
                        ) : (
                          <div className="flex items-center justify-center h-32 text-muted-foreground">
                            {t("noAssetData")}
                          </div>
                        )}
                      </CardContent>
                    </Card>
                  </>
                )}
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </AppLayout>
  );
}
