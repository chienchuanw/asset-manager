import { useMutation, useQueryClient } from "@tanstack/react-query";
import { refreshExchangeRate } from "@/lib/api/exchange-rates";
import { toast } from "sonner";

/**
 * 更新匯率的 Hook
 */
export function useRefreshExchangeRate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: refreshExchangeRate,
    onSuccess: (data) => {
      // 更新成功後，invalidate 相關的 queries
      queryClient.invalidateQueries({ queryKey: ["holdings"] });
      queryClient.invalidateQueries({ queryKey: ["analytics"] });
      queryClient.invalidateQueries({ queryKey: ["snapshots"] });
      queryClient.invalidateQueries({ queryKey: ["unrealized-analytics"] });

      toast.success("匯率更新成功", {
        description: `USD/TWD: ${data.rate.toFixed(4)} (${new Date(
          data.updated_at
        ).toLocaleString("zh-TW")})`,
      });
    },
    onError: (error) => {
      toast.error("匯率更新失敗", {
        description: error instanceof Error ? error.message : "未知錯誤",
      });
    },
  });
}

