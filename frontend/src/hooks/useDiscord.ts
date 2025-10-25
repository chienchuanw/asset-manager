import { useMutation } from "@tanstack/react-query";
import { testDiscord, sendDailyReport, type TestDiscordInput } from "@/lib/api/discord";
import { toast } from "sonner";

/**
 * 測試 Discord 發送的 Hook
 * 提供測試 Discord Webhook 的功能
 */
export function useTestDiscord() {
  return useMutation({
    mutationFn: (input: TestDiscordInput) => testDiscord(input),
    onSuccess: () => {
      toast.success("測試訊息已成功發送到 Discord！");
    },
    onError: (error: Error) => {
      toast.error(`發送失敗：${error.message}`);
    },
  });
}

/**
 * 發送每日報告的 Hook
 * 提供手動發送每日報告到 Discord 的功能
 */
export function useSendDailyReport() {
  return useMutation({
    mutationFn: () => sendDailyReport(),
    onSuccess: () => {
      toast.success("每日報告已成功發送到 Discord！");
    },
    onError: (error: Error) => {
      toast.error(`發送失敗：${error.message}`);
    },
  });
}

