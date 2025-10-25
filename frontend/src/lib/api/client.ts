import type { APIResponse } from "@/types/transaction";

/**
 * API 基礎 URL
 */
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

/**
 * API 錯誤類別
 */
export class APIError extends Error {
  constructor(public code: string, message: string, public status?: number) {
    super(message);
    this.name = "APIError";
  }
}

/**
 * Fetch 選項
 */
interface FetchOptions extends RequestInit {
  params?: Record<string, string | number | boolean | undefined | null>;
}

/**
 * 建立完整的 URL（包含查詢參數）
 */
function buildURL(path: string, params?: FetchOptions["params"]): string {
  const url = new URL(path, API_BASE_URL);

  if (params) {
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        url.searchParams.append(key, String(value));
      }
    });
  }

  return url.toString();
}

/**
 * 處理 API 回應
 */
async function handleResponse<T>(response: Response): Promise<T> {
  // 檢查 Content-Type
  const contentType = response.headers.get("content-type");
  const isJSON = contentType?.includes("application/json");

  // 如果不是 JSON，拋出錯誤
  if (!isJSON) {
    throw new APIError(
      "INVALID_RESPONSE",
      "伺服器回應格式錯誤",
      response.status
    );
  }

  // 解析 JSON
  const data: APIResponse<T> = await response.json();

  // 檢查是否有錯誤
  if (!response.ok || data.error) {
    throw new APIError(
      data.error?.code || "UNKNOWN_ERROR",
      data.error?.message || "未知錯誤",
      response.status
    );
  }

  // 回傳資料
  if (data.data === null) {
    throw new APIError("NO_DATA", "伺服器未回傳資料", response.status);
  }

  return data.data;
}

/**
 * 基礎 API 呼叫函式
 */
async function apiCall<T>(
  path: string,
  options: FetchOptions = {}
): Promise<T> {
  const { params, ...fetchOptions } = options;

  // 建立完整 URL
  const url = buildURL(path, params);

  // 設定預設 headers
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...fetchOptions.headers,
  };

  try {
    // 發送請求
    const response = await fetch(url, {
      ...fetchOptions,
      headers,
      credentials: fetchOptions.credentials || "include", // 預設包含 cookies
    });

    // 處理回應
    return await handleResponse<T>(response);
  } catch (error) {
    // 如果是 APIError，直接拋出
    if (error instanceof APIError) {
      throw error;
    }

    // 網路錯誤或其他錯誤
    if (error instanceof TypeError) {
      throw new APIError("NETWORK_ERROR", "網路連線失敗");
    }

    // 其他未知錯誤
    throw new APIError(
      "UNKNOWN_ERROR",
      error instanceof Error ? error.message : "未知錯誤"
    );
  }
}

/**
 * API Client
 */
export const apiClient = {
  /**
   * GET 請求
   */
  get: <T>(path: string, options?: FetchOptions) =>
    apiCall<T>(path, { ...options, method: "GET" }),

  /**
   * POST 請求
   */
  post: <T>(path: string, body?: unknown, options?: FetchOptions) =>
    apiCall<T>(path, {
      ...options,
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    }),

  /**
   * PUT 請求
   */
  put: <T>(path: string, body?: unknown, options?: FetchOptions) =>
    apiCall<T>(path, {
      ...options,
      method: "PUT",
      body: body ? JSON.stringify(body) : undefined,
    }),

  /**
   * DELETE 請求
   */
  delete: <T>(path: string, options?: FetchOptions) =>
    apiCall<T>(path, { ...options, method: "DELETE" }),
};

/**
 * 取得 API 基礎 URL（用於除錯）
 */
export function getAPIBaseURL(): string {
  return API_BASE_URL;
}
