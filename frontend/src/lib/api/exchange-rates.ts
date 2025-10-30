import { apiClient } from "./client";

/**
 * 匯率回應型別
 */
export interface ExchangeRateResponse {
  from_currency: string;
  to_currency: string;
  rate: number;
  date: string;
  updated_at: string;
  source: string;
}

/**
 * 更新今日匯率
 */
export async function refreshExchangeRate(): Promise<ExchangeRateResponse> {
  return apiClient.post<ExchangeRateResponse>("/api/exchange-rates/refresh");
}

