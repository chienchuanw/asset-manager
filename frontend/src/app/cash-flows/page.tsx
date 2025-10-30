/**
 * 現金流記錄頁面
 * 顯示所有現金流記錄，支援篩選、統計功能
 */

"use client";

import { useState, useMemo } from "react";
import { AppLayout } from "@/components/layout/AppLayout";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import {
  AddCashFlowDialog,
  CashFlowSummaryCard,
  CashFlowList,
} from "@/components/cash-flows";
import { useCashFlows } from "@/hooks";
import {
  CashFlowType,
  getCashFlowTypeLabel,
  type CashFlowFilters,
} from "@/types/cash-flow";
import { Download, Calendar } from "lucide-react";

export default function CashFlowsPage() {
  // 狀態管理
  const [filterType, setFilterType] = useState<CashFlowType | "all">("all");
  const [dateRange, setDateRange] = useState<"month" | "quarter" | "year">(
    "month"
  );

  // 計算日期範圍
  const { startDate, endDate } = useMemo(() => {
    const now = new Date();
    const year = now.getFullYear();
    const month = now.getMonth();

    let start: Date;
    let end: Date = now;

    switch (dateRange) {
      case "month":
        // 本月
        start = new Date(year, month, 1);
        break;
      case "quarter":
        // 本季
        const quarterStartMonth = Math.floor(month / 3) * 3;
        start = new Date(year, quarterStartMonth, 1);
        break;
      case "year":
        // 本年
        start = new Date(year, 0, 1);
        break;
      default:
        start = new Date(year, month, 1);
    }

    return {
      startDate: start.toISOString().split("T")[0],
      endDate: end.toISOString().split("T")[0],
    };
  }, [dateRange]);

  // 建立篩選條件
  const filters: CashFlowFilters = useMemo(() => {
    const baseFilters: CashFlowFilters = {
      start_date: startDate,
      end_date: endDate,
    };

    if (filterType !== "all") {
      baseFilters.type = filterType;
    }

    return baseFilters;
  }, [filterType, startDate, endDate]);

  // 取得現金流列表資料
  const { data: cashFlows, isLoading, error, refetch } = useCashFlows(filters);

  // 計算統計資料
  const stats = useMemo(() => {
    if (!cashFlows) {
      return {
        totalRecords: 0,
        incomeRecords: 0,
        expenseRecords: 0,
      };
    }

    return {
      totalRecords: cashFlows.length,
      incomeRecords: cashFlows.filter((cf) => cf.type === CashFlowType.INCOME)
        .length,
      expenseRecords: cashFlows.filter((cf) => cf.type === CashFlowType.EXPENSE)
        .length,
    };
  }, [cashFlows]);

  // 取得日期範圍標籤
  const getDateRangeLabel = (range: "month" | "quarter" | "year") => {
    switch (range) {
      case "month":
        return "本月";
      case "quarter":
        return "本季";
      case "year":
        return "本年";
    }
  };

  return (
    <AppLayout title="現金流記錄" description="追蹤和管理您的收入與支出">
      {/* Main Content */}
      <div className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="flex flex-col gap-6">
          {/* 統計卡片區域 - 六張卡片排列在同一排 */}
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6">
            {/* 摘要統計卡片 - 使用內嵌方式避免額外的 wrapper */}
            <CashFlowSummaryCard startDate={startDate} endDate={endDate} />

            {/* 記錄統計卡片 */}
            <Card>
              <CardHeader className="pb-2">
                <CardDescription>總記錄數</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.totalRecords}</div>
                <p className="text-xs text-muted-foreground mt-1">
                  {getDateRangeLabel(dateRange)}的所有記錄
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardDescription>收入記錄</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-green-600">
                  {stats.incomeRecords}
                </div>
                <p className="text-xs text-muted-foreground mt-1">
                  收入類型的記錄數
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardDescription>支出記錄</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-red-600">
                  {stats.expenseRecords}
                </div>
                <p className="text-xs text-muted-foreground mt-1">
                  支出類型的記錄數
                </p>
              </CardContent>
            </Card>
          </div>

          {/* 現金流記錄表格 */}
          <Card>
            <CardHeader>
              <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <CardTitle>現金流記錄</CardTitle>
                  <CardDescription>
                    {isLoading
                      ? "載入中..."
                      : `${getDateRangeLabel(dateRange)} - 共 ${
                          stats.totalRecords
                        } 筆記錄`}
                  </CardDescription>
                </div>
                <div className="flex gap-2">
                  <AddCashFlowDialog onSuccess={() => refetch()} />
                  <Button variant="outline" size="sm">
                    <Download className="h-4 w-4 mr-2" />
                    匯出
                  </Button>
                </div>
              </div>

              {/* 篩選工具列 */}
              <div className="flex flex-col gap-3 sm:flex-row sm:items-center mt-4">
                {/* 日期範圍篩選 */}
                <div className="flex items-center gap-2">
                  <Calendar className="h-4 w-4 text-muted-foreground" />
                  <Select
                    value={dateRange}
                    onValueChange={(value) =>
                      setDateRange(value as "month" | "quarter" | "year")
                    }
                  >
                    <SelectTrigger className="w-full sm:w-[150px]">
                      <SelectValue placeholder="日期範圍" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="month">本月</SelectItem>
                      <SelectItem value="quarter">本季</SelectItem>
                      <SelectItem value="year">本年</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                {/* 類型篩選 */}
                <Select
                  value={filterType}
                  onValueChange={(value) => setFilterType(value as any)}
                >
                  <SelectTrigger className="w-full sm:w-[150px]">
                    <SelectValue placeholder="類型" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">全部類型</SelectItem>
                    <SelectItem value={CashFlowType.INCOME}>
                      {getCashFlowTypeLabel(CashFlowType.INCOME)}
                    </SelectItem>
                    <SelectItem value={CashFlowType.EXPENSE}>
                      {getCashFlowTypeLabel(CashFlowType.EXPENSE)}
                    </SelectItem>
                  </SelectContent>
                </Select>
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
              <CashFlowList filters={filters} onRefresh={() => refetch()} />
            </CardContent>
          </Card>
        </div>
      </div>
    </AppLayout>
  );
}
