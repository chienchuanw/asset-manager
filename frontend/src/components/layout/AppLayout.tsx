/**
 * 應用程式主要佈局元件
 * 使用 shadcn/ui Sidebar 元件
 * 支援側邊欄收合/展開功能和鍵盤快捷鍵
 */

"use client";

import Link from "next/link";
import Image from "next/image";
import { usePathname } from "next/navigation";
import {
  HomeIcon,
  BarChart3Icon,
  WalletIcon,
  ArrowLeftRightIcon,
  SettingsIcon,
  UserIcon,
  LogOutIcon,
  TrendingUpIcon,
  RepeatIcon,
} from "lucide-react";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { Separator } from "@/components/ui/separator";
import { useAuth } from "@/providers/AuthProvider";

interface AppLayoutProps {
  children: React.ReactNode;
  title?: string;
  description?: string;
}

export function AppLayout({ children, title, description }: AppLayoutProps) {
  const pathname = usePathname();
  const { logout } = useAuth();

  // 從 cookie 讀取 sidebar 初始狀態
  const getInitialSidebarState = () => {
    if (typeof document === "undefined") return true; // SSR 時預設為 true
    const cookies = document.cookie.split("; ");
    const sidebarCookie = cookies.find((c) => c.startsWith("sidebar_state="));
    if (sidebarCookie) {
      return sidebarCookie.split("=")[1] === "true";
    }
    return true; // 預設為展開
  };

  const handleLogout = async () => {
    await logout();
  };

  // 主要導航項目
  const mainNavItems = [
    { id: "dashboard", label: "首頁", icon: HomeIcon, href: "/dashboard" },
    {
      id: "cash-flows",
      label: "現金流記錄",
      icon: TrendingUpIcon,
      href: "/cash-flows",
    },
    { id: "holdings", label: "持倉明細", icon: WalletIcon, href: "/holdings" },
    {
      id: "transactions",
      label: "交易記錄",
      icon: ArrowLeftRightIcon,
      href: "/transactions",
    },
    {
      id: "recurring",
      label: "訂閱分期",
      icon: RepeatIcon,
      href: "/recurring",
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
    { id: "settings", label: "設定", icon: SettingsIcon, href: "/settings" },
    {
      id: "user",
      label: "使用者管理",
      icon: UserIcon,
      href: "/user-management",
    },
    // { id: "help", label: "幫助", icon: HelpCircleIcon, href: "#" },
  ];

  return (
    <SidebarProvider defaultOpen={getInitialSidebarState()}>
      <Sidebar collapsible="icon">
        {/* Header: Logo */}
        <SidebarHeader>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton size="lg" asChild>
                <Link href="/dashboard">
                  <div className="flex aspect-square size-8 items-center justify-center rounded-lg">
                    <Image
                      src="/logotype_01.png"
                      alt="Asset Manager Logo"
                      width={32}
                      height={32}
                      className="object-contain"
                    />
                  </div>
                  <div className="grid flex-1 text-left text-sm leading-tight">
                    <span className="truncate font-semibold">
                      Asset Manager
                    </span>
                    <span className="truncate text-xs text-muted-foreground">
                      資產管理系統
                    </span>
                  </div>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarHeader>

        {/* Content: 主要導航 */}
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>主要功能</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {mainNavItems.map((item) => {
                  const Icon = item.icon;
                  const isActive = pathname === item.href;
                  return (
                    <SidebarMenuItem key={item.id}>
                      <SidebarMenuButton
                        asChild
                        isActive={isActive}
                        tooltip={item.label}
                      >
                        <Link href={item.href}>
                          <Icon />
                          <span>{item.label}</span>
                        </Link>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  );
                })}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>

        {/* Footer: 工具區 + 登出 */}
        <SidebarFooter>
          <SidebarMenu>
            {/* 工具區項目 */}
            {toolItems.map((item) => {
              const Icon = item.icon;
              const isActive = pathname === item.href;
              return (
                <SidebarMenuItem key={item.id}>
                  <SidebarMenuButton
                    asChild
                    isActive={isActive}
                    tooltip={item.label}
                  >
                    <Link href={item.href}>
                      <Icon />
                      <span>{item.label}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              );
            })}

            {/* 分隔線 */}
            <Separator className="my-2" />

            {/* 登出按鈕 */}
            <SidebarMenuItem>
              <SidebarMenuButton tooltip="登出" onClick={handleLogout}>
                <LogOutIcon />
                <span>登出</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarFooter>
      </Sidebar>

      {/* 主要內容區域 */}
      <main className="flex-1 flex flex-col min-w-0">
        {/* Header with Sidebar Trigger and Page Title */}
        <header className="sticky top-0 z-10 flex h-14 items-center gap-4 border-b bg-background/95 backdrop-blur px-4 lg:px-6">
          <SidebarTrigger />
          <Separator orientation="vertical" className="h-6" />

          {/* 動態頁面標題 */}
          {title ? (
            <div className="flex flex-col">
              <span className="font-semibold text-base">{title}</span>
              {description && (
                <span className="text-xs text-muted-foreground hidden sm:block">
                  {description}
                </span>
              )}
            </div>
          ) : (
            <div className="flex items-center gap-2">
              <span className="font-semibold">Asset Manager</span>
            </div>
          )}
        </header>

        {/* 頁面內容 - 加入適當的 padding 和 overflow 處理 */}
        <div className="flex-1 overflow-auto">{children}</div>
      </main>
    </SidebarProvider>
  );
}
