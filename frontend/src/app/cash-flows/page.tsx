/**
 * 現金流記錄頁面
 * 顯示所有現金流記錄，支援篩選、統計功能
 */

"use client";

import { useState, useMemo, useEffect } from "react";
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
  MonthlyExpenseChart,
  ExpenseCategoryPieChart,
} from "@/components/cash-flows";
import { calculateDateRange } from "@/components/common/DateRangeTabs";
import { WeekMonthTabs } from "@/components/common/WeekMonthTabs";
import { useCashFlows, cashFlowKeys } from "@/hooks";
import { useIsMobile } from "@/hooks/use-mobile";
import { CashFlowType, type CashFlowFilters } from "@/types/cash-flow";
import { Download, Search } from "lucide-react";
import { DateRange } from "react-day-picker";

// 日期記憶相關的 localStorage key
const SELECTED_DATE_KEY = "cash-flows-selected-date";

/**
 * 從 localStorage 讀取上次選擇的日期
 * 如果沒有記錄或日期無效，則返回今天
 */
function getInitialSelectedDate(): Date {
  if (typeof window === "undefined") return new Date();

  const savedDate = localStorage.getItem(SELECTED_DATE_KEY);
  if (!savedDate) return new Date();

  const parsedDate = new Date(savedDate);
  // 檢查日期是否有效
  if (isNaN(parsedDate.getTime())) return new Date();

  return parsedDate;
}

/**
 * 將選擇的日期儲存到 localStorage
 */
function saveSelectedDate(date: Date): void {
  if (typeof window === "undefined") return;
  localStorage.setItem(SELECTED_DATE_KEY, date.toISOString());
}

export default function CashFlowsPage() {
  const queryClient = useQueryClient();
  const isMobile = useIsMobile();

  // 狀態管理
  const [filterType, setFilterType] = useState<CashFlowType | "all">("all");
  const [customDateRange, setCustomDateRange] = useState<DateRange | undefined>(
    undefined
  );
  const [searchQuery, setSearchQuery] = useState("");

  // 左側「今日」專用的日期狀態 - 使用記憶的日期作為初始值
  const [selectedDate, setSelectedDate] = useState<Date>(
    getInitialSelectedDate
  );

  // 當使用者切換日期時，儲存到 localStorage
  useEffect(() => {
    saveSelectedDate(selectedDate);
  }, [selectedDate]);

  // 右側「本週/本月」專用的分頁狀態
  const [rightPanelTab, setRightPanelTab] = useState<"week" | "month">("week");

  // 計算左側「今日」的日期範圍
  const { startDate: todayStartDate, endDate: todayEndDate } = useMemo(() => {
    // 使用本地時間格式化日期，避免時區轉換問題
    const year = selectedDate.getFullYear();
    const month = String(selectedDate.getMonth() + 1).padStart(2, "0");
    const day = String(selectedDate.getDate()).padStart(2, "0");
    const dateStr = `${year}-${month}-${day}`;
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
      <div className="flex-1 p-4 md:p-6 bg-gray-50 space-y-6">
        {/* 當月/當週每日收入/支出圖表 - 置於頁面最上方 */}
        <Card>
          <CardHeader>
            <CardTitle>{isMobile ? "當週" : "當月"}每日收入/支出統計</CardTitle>
            <CardDescription>
              顯示 {new Date(selectedDate).getFullYear()} 年{" "}
              {new Date(selectedDate).getMonth() + 1} 月{isMobile ? "當週" : ""}
              的每日現金流動
            </CardDescription>
          </CardHeader>
          <CardContent>
            <MonthlyExpenseChart selectedDate={todayStartDate} />
          </CardContent>
        </Card>

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

            {/* 摘要統計卡片區域 - 2 欄 3 列佈局 */}
            <CashFlowSummaryCard
              startDate={rightStartDate}
              endDate={rightEndDate}
              totalRecords={rightStats.totalRecords}
              incomeRecords={rightStats.incomeRecords}
              expenseRecords={rightStats.expenseRecords}
            />

            {/* 支出分類圓餅圖 */}
            <ExpenseCategoryPieChart
              cashFlows={rightCashFlows || []}
              period={rightPanelTab}
            />
          </div>
        </div>
      </div>
    </AppLayout>
  );
}
