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
