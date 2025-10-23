/**
 * Dashboard 主頁面
 * 整合所有 Dashboard 元件,顯示資產概況
 */

'use client';

import { useState } from 'react';
import { Sidebar } from '@/components/dashboard/Sidebar';
import { MobileSidebar } from '@/components/dashboard/MobileSidebar';
import { Header } from '@/components/dashboard/Header';
import { StatCard } from '@/components/dashboard/StatCard';
import { AssetTrendChart } from '@/components/dashboard/AssetTrendChart';
import { HoldingsTable } from '@/components/dashboard/HoldingsTable';
import { AssetAllocationChart } from '@/components/dashboard/AssetAllocationChart';
import { RecentTransactions } from '@/components/dashboard/RecentTransactions';
import {
  mockStatCards,
  mockChartData,
  mockHoldings,
  mockAssetAllocation,
  mockTransactions,
} from '@/lib/mock-data';

export default function DashboardPage() {
  const [isMobileSidebarOpen, setIsMobileSidebarOpen] = useState(false);

  return (
    <div className="flex min-h-screen bg-gray-50">
      {/* 桌面版側邊欄 */}
      <Sidebar />

      {/* 手機版側邊欄 */}
      <MobileSidebar
        isOpen={isMobileSidebarOpen}
        onClose={() => setIsMobileSidebarOpen(false)}
      />

      {/* 主要內容區 */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Header */}
        <Header
          userName="Zac"
          onMenuClick={() => setIsMobileSidebarOpen(true)}
        />

        {/* 內容區域 */}
        <main className="flex-1 p-4 md:p-6">
          <div className="@container/main flex flex-1 flex-col gap-4 md:gap-6">
            {/* 統計卡片區 - 響應式網格 */}
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
              {mockStatCards.map((card, index) => (
                <StatCard
                  key={index}
                  title={card.title}
                  value={card.value}
                  change={card.change}
                  description={card.description}
                />
              ))}
            </div>

            {/* 主要內容區 - 響應式佈局 */}
            <div className="grid grid-cols-1 gap-4 lg:grid-cols-7 md:gap-6">
              {/* 左側：資產趨勢圖表 */}
              <div className="lg:col-span-4">
                <AssetTrendChart data={mockChartData} />
              </div>

              {/* 右側：資產配置圓餅圖 */}
              <div className="lg:col-span-3">
                <AssetAllocationChart data={mockAssetAllocation} />
              </div>
            </div>

            {/* 底部區域 - 響應式佈局 */}
            <div className="grid grid-cols-1 gap-4 lg:grid-cols-7 md:gap-6">
              {/* 左側：持倉明細表格 */}
              <div className="lg:col-span-4">
                <HoldingsTable holdings={mockHoldings} />
              </div>

              {/* 右側：近期交易 */}
              <div className="lg:col-span-3">
                <RecentTransactions transactions={mockTransactions} />
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}

