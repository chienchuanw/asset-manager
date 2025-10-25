import { apiClient } from "./client";

/**
 * 登入請求參數
 */
export interface LoginRequest {
  username: string;
  password: string;
}

/**
 * 登入回應
 */
export interface LoginResponse {
  message: string;
}

/**
 * 使用者資訊
 */
export interface User {
  username: string;
}

/**
 * 登入
 * 成功後會自動設定 httpOnly cookie
 */
export async function login(
  credentials: LoginRequest
): Promise<LoginResponse> {
  return apiClient.post<LoginResponse>("/api/auth/login", credentials, {
    credentials: "include", // 重要：允許發送和接收 cookies
  });
}

/**
 * 登出
 * 會清除 httpOnly cookie
 */
export async function logout(): Promise<LoginResponse> {
  return apiClient.post<LoginResponse>("/api/auth/logout", undefined, {
    credentials: "include",
  });
}

/**
 * 取得當前使用者資訊
 * 需要有效的 JWT token (在 cookie 中)
 */
export async function getCurrentUser(): Promise<User> {
  return apiClient.get<User>("/api/auth/me", {
    credentials: "include",
  });
}

