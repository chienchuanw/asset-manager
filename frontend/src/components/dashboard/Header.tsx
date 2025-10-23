/**
 * 頂部 Header 元件
 * 顯示歡迎訊息、搜尋、通知和使用者資訊
 */

'use client';

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { SearchIcon, BellIcon, ChevronDownIcon, MenuIcon } from 'lucide-react';

interface HeaderProps {
  userName?: string;
  onMenuClick?: () => void;
}

export function Header({ userName = '使用者', onMenuClick }: HeaderProps) {
  return (
    <header className="bg-white border-b border-gray-200 sticky top-0 z-10">
      <div className="flex items-center justify-between px-6 py-4">
        {/* 左側：歡迎訊息 + 手機版選單按鈕 */}
        <div className="flex items-center gap-4">
          {/* 手機版選單按鈕 */}
          <Button
            variant="ghost"
            size="icon"
            className="lg:hidden"
            onClick={onMenuClick}
          >
            <MenuIcon className="h-6 w-6" />
          </Button>

          {/* 歡迎訊息 */}
          <div>
            <h1 className="text-xl font-bold text-gray-900">
              歡迎回來, {userName}!
            </h1>
            <p className="text-sm text-gray-600 hidden sm:block">
              這是你今天的資產概況
            </p>
          </div>
        </div>

        {/* 右側：搜尋、通知、使用者選單 */}
        <div className="flex items-center gap-3">
          {/* 搜尋按鈕 */}
          <Button variant="ghost" size="icon" className="hidden sm:flex">
            <SearchIcon className="h-5 w-5 text-gray-600" />
          </Button>

          {/* 通知按鈕 */}
          <Button variant="ghost" size="icon" className="relative">
            <BellIcon className="h-5 w-5 text-gray-600" />
            {/* 通知小紅點 */}
            <span className="absolute top-2 right-2 w-2 h-2 bg-red-500 rounded-full"></span>
          </Button>

          {/* 使用者選單 */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="flex items-center gap-2 pl-2">
                <Avatar className="h-8 w-8">
                  <AvatarImage src="/avatar-placeholder.png" alt={userName} />
                  <AvatarFallback className="bg-gradient-to-br from-blue-600 to-purple-600 text-white">
                    {userName.charAt(0).toUpperCase()}
                  </AvatarFallback>
                </Avatar>
                <span className="hidden sm:inline text-sm font-medium text-gray-700">
                  {userName}
                </span>
                <ChevronDownIcon className="h-4 w-4 text-gray-600" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuLabel>我的帳戶</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem>個人資料</DropdownMenuItem>
              <DropdownMenuItem>帳戶設定</DropdownMenuItem>
              <DropdownMenuItem>偏好設定</DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem className="text-red-600">登出</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header>
  );
}

