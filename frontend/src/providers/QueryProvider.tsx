"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { useState } from "react";

/**
 * React Query Provider
 * 
 * 提供 React Query 的功能給整個應用程式
 */
export function QueryProvider({ children }: { children: React.ReactNode }) {
  // 建立 QueryClient 實例（使用 useState 確保每個 client 只建立一次）
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            // 資料過期時間（5 分鐘）
            staleTime: 5 * 60 * 1000,
            // 快取時間（10 分鐘）
            gcTime: 10 * 60 * 1000,
            // 失敗時重試次數
            retry: 1,
            // 視窗重新獲得焦點時重新獲取資料
            refetchOnWindowFocus: false,
            // 重新連線時重新獲取資料
            refetchOnReconnect: true,
          },
          mutations: {
            // 失敗時重試次數
            retry: 0,
          },
        },
      })
  );

  return (
    <QueryClientProvider client={queryClient}>
      {children}
      {/* React Query Devtools（僅在開發環境顯示）*/}
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}

