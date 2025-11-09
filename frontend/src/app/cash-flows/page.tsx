/**
 * 現金流記錄頁面
 * 顯示所有現金流記錄，支援篩選、統計功能
 */

"use client";

import { useState, useMemo } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { AppLayout } from "@/components/layout/AppLayout";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  AddCashFlowDialog,
  CashFlowSummaryCard,
  CashFlowList,
  CashFlowFilterDrawer,
  DailyDateNavigator,
  DailySummaryCards,
} from "@/components/cash-flows";
import { calculateDateRange } from "@/components/common/DateRangeTabs";
import { WeekMonthTabs } from "@/components/common/WeekMonthTabs";
import { useCashFlows, cashFlowKeys } from "@/hooks";
import { CashFlowType, type CashFlowFilters } from "@/types/cash-flow";
import { Download, Search } from "lucide-react";
import { DateRange } from "react-day-picker";

export default function CashFlowsPage() {
  const queryClient = useQueryClient();

  // 狀態管理
  const [filterType, setFilterType] = useState<CashFlowType | "all">("all");
  const [customDateRange, setCustomDateRange] = useState<DateRange | undefined>(
    undefined
  );
  const [searchQuery, setSearchQuery] = useState("");

  // 左側「今日」專用的日期狀態
  const [selectedDate, setSelectedDate] = useState<Date>(new Date());

  // 右側「本週/本月」專用的分頁狀態
  const [rightPanelTab, setRightPanelTab] = useState<"week" | "month">("week");

  // 計算左側「今日」的日期範圍
  const { startDate: todayStartDate, endDate: todayEndDate } = useMemo(() => {
    const dateStr = selectedDate.toISOString().split("T")[0];
    return {
      startDate: dateStr,
      endDate: dateStr,
    };
  }, [selectedDate]);

  // 計算右側「本週/本月」的日期範圍
  const { startDate: rightStartDate, endDate: rightEndDate } = useMemo(() => {
    return calculateDateRange(rightPanelTab);
  }, [rightPanelTab]);

  // 建立左側「今日」的篩選條件
  const todayFilters: CashFlowFilters = useMemo(() => {
    const baseFilters: CashFlowFilters = {
      start_date: todayStartDate,
      end_date: todayEndDate,
    };

    if (filterType !== "all") {
      baseFilters.type = filterType;
    }

    return baseFilters;
  }, [filterType, todayStartDate, todayEndDate]);

  // 取得左側「今日」的現金流列表資料
  const {
    data: cashFlows,
    isLoading,
    error,
  } = useCashFlows(todayFilters, {
    // 確保資料總是最新的
    staleTime: 0,
  });

  // 取得右側「本週/本月」的現金流列表資料（用於統計）
  const rightFilters: CashFlowFilters = useMemo(() => {
    return {
      start_date: rightStartDate,
      end_date: rightEndDate,
    };
  }, [rightStartDate, rightEndDate]);

  const { data: rightCashFlows } = useCashFlows(rightFilters, {
    staleTime: 0,
  });

  // 計算右側統計資料
  const rightStats = useMemo(() => {
    if (!rightCashFlows) {
      return {
        totalRecords: 0,
        incomeRecords: 0,
        expenseRecords: 0,
      };
    }

    return {
      totalRecords: rightCashFlows.length,
      incomeRecords: rightCashFlows.filter(
        (cf) => cf.type === CashFlowType.INCOME
      ).length,
      expenseRecords: rightCashFlows.filter(
        (cf) => cf.type === CashFlowType.EXPENSE
      ).length,
    };
  }, [rightCashFlows]);

  // 重置篩選
  const handleResetFilters = () => {
    setFilterType("all");
    setCustomDateRange(undefined);
  };

  // 重新獲取所有相關資料
  const handleRefreshData = async () => {
    // 讓所有現金流相關查詢失效，強制重新獲取
    await queryClient.invalidateQueries({
      queryKey: cashFlowKeys.all,
    });
  };

  // 客戶端篩選(搜尋)
  const filteredCashFlows = useMemo(() => {
    if (!cashFlows) return [];
    if (!searchQuery) return cashFlows;

    const query = searchQuery.toLowerCase();
    return cashFlows.filter(
      (cf) =>
        cf.description.toLowerCase().includes(query) ||
        cf.note?.toLowerCase().includes(query) ||
        cf.category?.name.toLowerCase().includes(query)
    );
  }, [cashFlows, searchQuery]);

  return (
    <AppLayout title="現金流記錄" description="追蹤和管理您的收入與支出">
      {/* Main Content */}
      <div className="flex-1 p-4 md:p-6 bg-gray-50">
        {/* 左右分欄佈局 */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* 左側：今日內容 */}
          <div className="flex flex-col gap-6">
            {/* 日期導航 */}
            <DailyDateNavigator
              date={selectedDate}
              onDateChange={setSelectedDate}
            />

            {/* 今日摘要卡片 */}
            <DailySummaryCards date={todayStartDate} />

            {/* 今日現金流記錄列表 */}
            <Card>
              <CardHeader>
                <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                  <div>
                    <CardTitle>現金流記錄</CardTitle>
                    <CardDescription>
                      {isLoading
                        ? "載入中..."
                        : `共 ${filteredCashFlows.length} 筆記錄`}
                    </CardDescription>
                  </div>
                  <div className="flex gap-2">
                    <AddCashFlowDialog onSuccess={handleRefreshData} />
                    <Button variant="outline" size="sm">
                      <Download className="h-4 w-4 mr-2" />
                      匯出
                    </Button>
                  </div>
                </div>

                {/* 篩選工具列 */}
                <div className="flex flex-col gap-3 sm:flex-row sm:items-center mt-4">
                  {/* 搜尋框 */}
                  <div className="relative flex-1">
                    <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                    <Input
                      placeholder="搜尋描述、備註或分類..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="pl-9"
                    />
                  </div>

                  {/* 進階篩選 Drawer */}
                  <CashFlowFilterDrawer
                    filterType={filterType}
                    dateRange={customDateRange}
                    onFilterTypeChange={setFilterType}
                    onDateRangeChange={setCustomDateRange}
                    onReset={handleResetFilters}
                  />
                </div>
              </CardHeader>

              <CardContent>
                {/* 錯誤訊息 */}
                {error && (
                  <div className="p-4 mb-4 text-sm text-red-800 bg-red-100 rounded-lg">
                    <p className="font-medium">載入失敗</p>
                    <p>{error.message}</p>
                  </div>
                )}

                {/* 現金流列表 */}
                <CashFlowList filters={todayFilters} />
              </CardContent>
            </Card>
          </div>

          {/* 右側：本週/本月內容（手機版隱藏）*/}
          <div className="hidden lg:flex flex-col gap-6">
            {/* 本週/本月 Tabs */}
            <WeekMonthTabs
              value={rightPanelTab}
              onValueChange={setRightPanelTab}
            />

            {/* 卡片區域 - 垂直排列（1 欄 6 列）*/}
            <div className="flex flex-col gap-4">
              {/* 摘要統計卡片 - 佔三欄 */}
              <CashFlowSummaryCard
                startDate={rightStartDate}
                endDate={rightEndDate}
              />

              {/* 記錄統計卡片 - 垂直排列 */}
              <Card className="hover:shadow-lg transition-shadow">
                <CardHeader className="pb-3">
                  <CardDescription className="text-base">
                    總記錄數
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="text-4xl font-bold mb-2">
                    {rightStats.totalRecords}
                  </div>
                  <p className="text-sm text-muted-foreground">
                    期間內的所有交易記錄
                  </p>
                </CardContent>
              </Card>

              <Card className="hover:shadow-lg transition-shadow">
                <CardHeader className="pb-3">
                  <CardDescription className="text-base">
                    收入記錄
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="text-4xl font-bold text-green-600 mb-2">
                    {rightStats.incomeRecords}
                  </div>
                  <p className="text-sm text-muted-foreground">
                    期間內的收入交易筆數
                  </p>
                </CardContent>
              </Card>

              <Card className="hover:shadow-lg transition-shadow">
                <CardHeader className="pb-3">
                  <CardDescription className="text-base">
                    支出記錄
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="text-4xl font-bold text-red-600 mb-2">
                    {rightStats.expenseRecords}
                  </div>
                  <p className="text-sm text-muted-foreground">
                    期間內的支出交易筆數
                  </p>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </div>
    </AppLayout>
  );
}
