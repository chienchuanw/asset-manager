"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import * as authAPI from "@/lib/api/auth";
import type { User, LoginRequest } from "@/lib/api/auth";
import { toast } from "sonner";

/**
 * Auth Context 的型別定義
 */
interface AuthContextType {
  // 狀態
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;

  // 方法
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
}

/**
 * Auth Context
 */
const AuthContext = createContext<AuthContextType | undefined>(undefined);

/**
 * Auth Provider Props
 */
interface AuthProviderProps {
  children: React.ReactNode;
}

/**
 * Auth Provider
 * 管理使用者的登入狀態
 */
export function AuthProvider({ children }: AuthProviderProps) {
  const router = useRouter();
  const queryClient = useQueryClient();
  const t = useTranslations("auth");
  const [isInitialized, setIsInitialized] = useState(false);

  // 使用 React Query 查詢當前使用者
  const {
    data: user,
    isLoading,
    refetch,
  } = useQuery({
    queryKey: ["currentUser"],
    queryFn: authAPI.getCurrentUser,
    retry: false, // 不自動重試（避免多次 401 請求）
    enabled: false, // 預設不自動執行，由 checkAuth 手動觸發
  });

  // 登入 mutation
  const loginMutation = useMutation({
    mutationFn: authAPI.login,
    onSuccess: async () => {
      // 登入成功後，重新取得使用者資訊
      await refetch();
      toast.success(t("loginSuccess"));
      router.push("/dashboard");
    },
    onError: (error: Error) => {
      toast.error(t("loginError", { error: error.message }));
    },
  });

  // 登出 mutation
  const logoutMutation = useMutation({
    mutationFn: authAPI.logout,
    onSuccess: () => {
      // 清除所有 query cache
      queryClient.clear();
      toast.success(t("logoutSuccess"));
      router.push("/login");
    },
    onError: (error: Error) => {
      toast.error(t("logoutError", { error: error.message }));
    },
  });

  // 檢查登入狀態
  const checkAuth = async () => {
    try {
      await refetch();
    } catch (error) {
      // 如果檢查失敗（例如 token 過期），不做任何處理
      // 讓 middleware 處理重導向
    } finally {
      setIsInitialized(true);
    }
  };

  // 登入方法
  const login = async (credentials: LoginRequest) => {
    await loginMutation.mutateAsync(credentials);
  };

  // 登出方法
  const logout = async () => {
    await logoutMutation.mutateAsync();
  };

  // 初始化時檢查登入狀態
  useEffect(() => {
    checkAuth();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // 計算是否已登入
  const isAuthenticated = !!user;

  return (
    <AuthContext.Provider
      value={{
        user: user || null,
        isAuthenticated,
        isLoading,
        login,
        logout,
        checkAuth,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

/**
 * useAuth Hook
 * 用於在元件中存取 Auth Context
 */
export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
