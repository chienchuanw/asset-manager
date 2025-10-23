/**
 * 應用程式主要佈局元件
 * 提供側邊欄、Header 和主要內容區域
 * 支援側邊欄收合/展開功能
 */

'use client';

import { useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import {
  HomeIcon,
  BarChart3Icon,
  WalletIcon,
  ArrowLeftRightIcon,
  SettingsIcon,
  HelpCircleIcon,
  UserIcon,
  LogOutIcon,
  MenuIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from 'lucide-react';

interface AppLayoutProps {
  children: React.ReactNode;
}

export function AppLayout({ children }: AppLayoutProps) {
  const pathname = usePathname();
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  // 主要導航項目
  const mainNavItems = [
    { id: 'dashboard', label: '首頁', icon: HomeIcon, href: '/dashboard' },
    { id: 'holdings', label: '持倉明細', icon: WalletIcon, href: '/holdings' },
    { id: 'transactions', label: '交易記錄', icon: ArrowLeftRightIcon, href: '/transactions' },
    { id: 'analytics', label: '分析報表', icon: BarChart3Icon, href: '/analytics' },
  ];

  // 工具區項目
  const toolItems = [
    { id: 'settings', label: '設定', icon: SettingsIcon },
    { id: 'help', label: '幫助', icon: HelpCircleIcon },
    { id: 'user', label: '使用者管理', icon: UserIcon },
  ];

  return (
    <div className="flex min-h-screen">
      {/* 桌面版側邊欄 */}
      <aside
        className={`hidden lg:flex lg:flex-col bg-white border-r border-gray-200 h-screen sticky top-0 transition-all duration-300 ${
          sidebarOpen ? 'lg:w-64' : 'lg:w-20'
        }`}
      >
        {/* Logo 區域 */}
        <div className="p-6">
          <div className="flex items-center gap-2">
            {sidebarOpen && (
              <span className="text-xl font-bold text-gray-900 whitespace-nowrap">
                Asset Manager
              </span>
            )}
          </div>
        </div>

        {/* 主要導航 */}
        <nav className="flex-1 px-4 space-y-1">
          {mainNavItems.map((item) => {
            const Icon = item.icon;
            const isActive = pathname === item.href;
            return (
              <Link key={item.id} href={item.href}>
                <Button
                  variant={isActive ? 'default' : 'ghost'}
                  className={`w-full gap-3 ${
                    sidebarOpen ? 'justify-start' : 'justify-center'
                  } ${
                    isActive
                      ? 'bg-gray-900 text-white hover:bg-gray-800'
                      : 'text-gray-700 hover:bg-gray-100'
                  }`}
                  title={!sidebarOpen ? item.label : undefined}
                >
                  <Icon className="h-5 w-5 shrink-0" />
                  {sidebarOpen && <span>{item.label}</span>}
                </Button>
              </Link>
            );
          })}

          <Separator className="my-4" />

          {/* 工具區 */}
          <div className="pt-2">
            {sidebarOpen && (
              <p className="px-3 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
                工具
              </p>
            )}
            {toolItems.map((item) => {
              const Icon = item.icon;
              return (
                <Button
                  key={item.id}
                  variant="ghost"
                  className={`w-full gap-3 text-gray-700 hover:bg-gray-100 ${
                    sidebarOpen ? 'justify-start' : 'justify-center'
                  }`}
                  title={!sidebarOpen ? item.label : undefined}
                >
                  <Icon className="h-5 w-5 shrink-0" />
                  {sidebarOpen && <span>{item.label}</span>}
                </Button>
              );
            })}
          </div>
        </nav>

        {/* 收合按鈕 */}
        <div className="p-4 border-t border-gray-200">
          <Button
            variant="ghost"
            className={`w-full gap-3 text-gray-700 hover:bg-gray-100 ${
              sidebarOpen ? 'justify-start' : 'justify-center'
            }`}
            onClick={() => setSidebarOpen(!sidebarOpen)}
          >
            {sidebarOpen ? (
              <>
                <ChevronLeftIcon className="h-5 w-5 shrink-0" />
                <span>收合</span>
              </>
            ) : (
              <ChevronRightIcon className="h-5 w-5 shrink-0" />
            )}
          </Button>
        </div>

        {/* 登出按鈕 */}
        <div className="p-4 border-t border-gray-200">
          <Button
            variant="ghost"
            className={`w-full gap-3 text-gray-700 hover:bg-gray-100 ${
              sidebarOpen ? 'justify-start' : 'justify-center'
            }`}
            title={!sidebarOpen ? '登出' : undefined}
          >
            <LogOutIcon className="h-5 w-5 shrink-0" />
            {sidebarOpen && <span>登出</span>}
          </Button>
        </div>
      </aside>

      {/* 手機版背景遮罩 */}
      {mobileMenuOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-40 lg:hidden"
          onClick={() => setMobileMenuOpen(false)}
        />
      )}

      {/* 手機版側邊欄 */}
      <aside
        className={`fixed inset-y-0 left-0 z-50 w-64 bg-white border-r border-gray-200 transform transition-transform duration-300 lg:hidden ${
          mobileMenuOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        {/* Logo 區域 + 關閉按鈕 */}
        <div className="p-6 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-linear-to-br from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
              <WalletIcon className="h-5 w-5 text-white" />
            </div>
            <span className="text-xl font-bold text-gray-900">Asset Manager</span>
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setMobileMenuOpen(false)}
            className="lg:hidden"
          >
            <ChevronLeftIcon className="h-5 w-5" />
          </Button>
        </div>

        {/* 主要導航 */}
        <nav className="flex-1 px-4 space-y-1 overflow-y-auto">
          {mainNavItems.map((item) => {
            const Icon = item.icon;
            const isActive = pathname === item.href;
            return (
              <Link key={item.id} href={item.href} onClick={() => setMobileMenuOpen(false)}>
                <Button
                  variant={isActive ? 'default' : 'ghost'}
                  className={`w-full justify-start gap-3 ${
                    isActive
                      ? 'bg-gray-900 text-white hover:bg-gray-800'
                      : 'text-gray-700 hover:bg-gray-100'
                  }`}
                >
                  <Icon className="h-5 w-5" />
                  {item.label}
                </Button>
              </Link>
            );
          })}

          <Separator className="my-4" />

          {/* 工具區 */}
          <div className="pt-2">
            <p className="px-3 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
              工具
            </p>
            {toolItems.map((item) => {
              const Icon = item.icon;
              return (
                <Button
                  key={item.id}
                  variant="ghost"
                  className="w-full justify-start gap-3 text-gray-700 hover:bg-gray-100"
                  onClick={() => setMobileMenuOpen(false)}
                >
                  <Icon className="h-5 w-5" />
                  {item.label}
                </Button>
              );
            })}
          </div>
        </nav>

        {/* 登出按鈕 */}
        <div className="p-4 border-t border-gray-200">
          <Button variant="ghost" className="w-full justify-start gap-3 text-gray-700 hover:bg-gray-100">
            <LogOutIcon className="h-5 w-5" />
            登出
          </Button>
        </div>
      </aside>

      {/* 主要內容區域 */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* 手機版 Header */}
        <header className="sticky top-0 z-10 lg:hidden border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
          <div className="flex h-14 items-center px-4 gap-4">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setMobileMenuOpen(true)}
            >
              <MenuIcon className="h-5 w-5" />
            </Button>
            <div className="flex items-center gap-2">
              <span className="font-semibold">Asset Manager</span>
            </div>
          </div>
        </header>

        {/* 頁面內容 */}
        {children}
      </div>
    </div>
  );
}

