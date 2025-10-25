/**
 * 側邊欄導航元件
 * 提供主要導航功能
 */

"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import {
  HomeIcon,
  BarChart3Icon,
  WalletIcon,
  ArrowLeftRightIcon,
  SettingsIcon,
  HelpCircleIcon,
  UserIcon,
  LogOutIcon,
} from "lucide-react";

export function Sidebar() {
  const pathname = usePathname();

  // 主要導航項目
  const mainNavItems = [
    { id: "dashboard", label: "首頁", icon: HomeIcon, href: "/dashboard" },
    { id: "holdings", label: "持倉明細", icon: WalletIcon, href: "/holdings" },
    {
      id: "transactions",
      label: "交易記錄",
      icon: ArrowLeftRightIcon,
      href: "/transactions",
    },
    {
      id: "analytics",
      label: "分析報表",
      icon: BarChart3Icon,
      href: "/analytics",
    },
  ];

  // 工具區項目
  const toolItems = [
    { id: "settings", label: "設定", icon: SettingsIcon },
    { id: "help", label: "幫助", icon: HelpCircleIcon },
    { id: "user", label: "使用者管理", icon: UserIcon },
  ];

  return (
    <aside className="hidden lg:flex lg:flex-col lg:w-64 bg-white border-r border-gray-200 h-screen sticky top-0">
      {/* Logo 區域 */}
      <div className="p-6">
        <div className="flex items-center gap-2">
          <span className="text-xl font-bold text-gray-900">Asset Manager</span>
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
                variant={isActive ? "default" : "ghost"}
                className={`w-full justify-start gap-3 ${
                  isActive
                    ? "bg-gray-900 text-white hover:bg-gray-800"
                    : "text-gray-700 hover:bg-gray-100"
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
        <Button
          variant="ghost"
          className="w-full justify-start gap-3 text-gray-700 hover:bg-gray-100"
        >
          <LogOutIcon className="h-5 w-5" />
          登出
        </Button>
      </div>
    </aside>
  );
}
