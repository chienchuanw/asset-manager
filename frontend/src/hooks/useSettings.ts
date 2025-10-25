import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getSettings, updateSettings } from '@/lib/api/settings';
import type { UpdateSettingsGroupInput } from '@/types/analytics';

/**
 * 取得設定的 Hook
 */
export function useSettings() {
  return useQuery({
    queryKey: ['settings'],
    queryFn: getSettings,
  });
}

/**
 * 更新設定的 Hook
 */
export function useUpdateSettings() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: UpdateSettingsGroupInput) => updateSettings(input),
    onSuccess: () => {
      // 更新成功後重新取得設定
      queryClient.invalidateQueries({ queryKey: ['settings'] });
    },
  });
}

