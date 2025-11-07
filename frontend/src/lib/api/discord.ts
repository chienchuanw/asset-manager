import { apiClient } from "./client";

/**
 * Discord API 客戶端
 * 提供 Discord 相關的 API 呼叫功能
 */

// 測試 Discord 輸入
export interface TestDiscordInput {
  message: string;
}

/**
 * 測試 Discord 發送
 * @param input 測試訊息
 * @returns 成功訊息
 */
export async function testDiscord(input: TestDiscordInput): Promise<string> {
  return await apiClient.post<string>("/api/discord/test", input);
}

/**
 * 發送每日報告到 Discord
 * @returns 成功訊息
 */
export async function sendDailyReport(): Promise<string> {
  return await apiClient.post<string>("/api/discord/daily-report", {});
}

// 發送月度報告輸入
export interface SendMonthlyReportInput {
  year: number;
  month: number;
  webhook_url: string;
}

/**
 * 發送月度現金流報告到 Discord
 * @param input 年份、月份和 Webhook URL
 * @returns 成功訊息
 */
export async function sendMonthlyReport(
  input: SendMonthlyReportInput
): Promise<string> {
  const params = new URLSearchParams({
    year: input.year.toString(),
    month: input.month.toString(),
    webhook_url: input.webhook_url,
  });
  return await apiClient.post<string>(
    `/api/cash-flows/send-monthly-report?${params}`,
    {}
  );
}

// 發送年度報告輸入
export interface SendYearlyReportInput {
  year: number;
  webhook_url: string;
}

/**
 * 發送年度現金流報告到 Discord
 * @param input 年份和 Webhook URL
 * @returns 成功訊息
 */
export async function sendYearlyReport(
  input: SendYearlyReportInput
): Promise<string> {
  const params = new URLSearchParams({
    year: input.year.toString(),
    webhook_url: input.webhook_url,
  });
  return await apiClient.post<string>(
    `/api/cash-flows/send-yearly-report?${params}`,
    {}
  );
}
