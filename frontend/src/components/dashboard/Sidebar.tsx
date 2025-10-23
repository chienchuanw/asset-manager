/**
 * 側邊欄導航元件
 * 提供主要導航功能
 */

'use client';

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
} from 'lucide-react';
import { useState } from 'react';

export function Sidebar() {
  const [activeItem, setActiveItem] = useState('dashboard');

  // 主要導航項目
  const mainNavItems = [
    { id: 'dashboard', label: '首頁', icon: HomeIcon },
    { id: 'holdings', label: '持倉明細', icon: WalletIcon },
    { id: 'transactions', label: '交易記錄', icon: ArrowLeftRightIcon },
    { id: 'analytics', label: '分析報表', icon: BarChart3Icon },
  ];

  // 工具區項目
  const toolItems = [
    { id: 'settings', label: '設定', icon: SettingsIcon },
    { id: 'help', label: '幫助', icon: HelpCircleIcon },
    { id: 'user', label: '使用者管理', icon: UserIcon },
  ];

  return (
    <aside className="hidden lg:flex lg:flex-col lg:w-64 bg-white border-r border-gray-200 h-screen sticky top-0">
      {/* Logo 區域 */}
      <div className="p-6">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 bg-gradient-to-br from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
            <WalletIcon className="h-5 w-5 text-white" />
          </div>
          <span className="text-xl font-bold text-gray-900">Asset Manager</span>
        </div>
      </div>

      {/* 主要導航 */}
      <nav className="flex-1 px-4 space-y-1">
        {mainNavItems.map((item) => {
          const Icon = item.icon;
          const isActive = activeItem === item.id;
          return (
            <Button
              key={item.id}
              variant={isActive ? 'default' : 'ghost'}
              className={`w-full justify-start gap-3 ${
                isActive
                  ? 'bg-gray-900 text-white hover:bg-gray-800'
                  : 'text-gray-700 hover:bg-gray-100'
              }`}
              onClick={() => setActiveItem(item.id)}
            >
              <Icon className="h-5 w-5" />
              {item.label}
            </Button>
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
            const isActive = activeItem === item.id;
            return (
              <Button
                key={item.id}
                variant="ghost"
                className={`w-full justify-start gap-3 ${
                  isActive
                    ? 'bg-gray-100 text-gray-900'
                    : 'text-gray-700 hover:bg-gray-100'
                }`}
                onClick={() => setActiveItem(item.id)}
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
  );
}

