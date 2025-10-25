import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryOptions,
  type UseMutationOptions,
} from "@tanstack/react-query";
import { subscriptionsAPI } from "@/lib/api/subscriptions";
import type {
  Subscription,
  CreateSubscriptionInput,
  UpdateSubscriptionInput,
  CancelSubscriptionInput,
  SubscriptionFilters,
} from "@/types/subscription";
import { APIError } from "@/lib/api/client";

/**
 * Query Keys
 * 用於識別和管理快取
 */
export const subscriptionKeys = {
  all: ["subscriptions"] as const,
  lists: () => [...subscriptionKeys.all, "list"] as const,
  list: (filters?: SubscriptionFilters) =>
    [...subscriptionKeys.lists(), filters] as const,
  details: () => [...subscriptionKeys.all, "detail"] as const,
  detail: (id: string) => [...subscriptionKeys.details(), id] as const,
};

// ==================== 訂閱 Hooks ====================

/**
 * 取得訂閱列表
 *
 * @param filters 篩選條件
 * @param options React Query 選項
 * @returns 訂閱列表查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useSubscriptions({
 *   status: 'active',
 *   limit: 10,
 * });
 * ```
 */
export function useSubscriptions(
  filters?: SubscriptionFilters,
  options?: Omit<
    UseQueryOptions<Subscription[], APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<Subscription[], APIError>({
    queryKey: subscriptionKeys.list(filters),
    queryFn: () => subscriptionsAPI.getAll(filters),
    ...options,
  });
}

/**
 * 取得單筆訂閱
 *
 * @param id 訂閱 ID
 * @param options React Query 選項
 * @returns 單筆訂閱查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useSubscription("subscription-id");
 * ```
 */
export function useSubscription(
  id: string,
  options?: Omit<
    UseQueryOptions<Subscription, APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<Subscription, APIError>({
    queryKey: subscriptionKeys.detail(id),
    queryFn: () => subscriptionsAPI.getById(id),
    enabled: !!id,
    ...options,
  });
}

/**
 * 建立訂閱
 *
 * @param options Mutation 選項
 * @returns 建立訂閱的 mutation
 *
 * @example
 * ```tsx
 * const createMutation = useCreateSubscription();
 *
 * const handleCreate = async (data: CreateSubscriptionInput) => {
 *   await createMutation.mutateAsync(data);
 * };
 * ```
 */
export function useCreateSubscription(
  options?: UseMutationOptions<
    Subscription,
    APIError,
    CreateSubscriptionInput,
    unknown
  >
) {
  const queryClient = useQueryClient();

  return useMutation<Subscription, APIError, CreateSubscriptionInput>({
    mutationFn: subscriptionsAPI.create,
    onSuccess: () => {
      // 使所有訂閱列表查詢失效，觸發重新獲取
      queryClient.invalidateQueries({ queryKey: subscriptionKeys.lists() });
    },
    ...options,
  });
}

/**
 * 更新訂閱
 *
 * @param options Mutation 選項
 * @returns 更新訂閱的 mutation
 *
 * @example
 * ```tsx
 * const updateMutation = useUpdateSubscription();
 *
 * const handleUpdate = async (id: string, data: UpdateSubscriptionInput) => {
 *   await updateMutation.mutateAsync({ id, data });
 * };
 * ```
 */
export function useUpdateSubscription(
  options?: UseMutationOptions<
    Subscription,
    APIError,
    { id: string; data: UpdateSubscriptionInput },
    unknown
  >
) {
  const queryClient = useQueryClient();

  return useMutation<
    Subscription,
    APIError,
    { id: string; data: UpdateSubscriptionInput }
  >({
    mutationFn: ({ id, data }) => subscriptionsAPI.update(id, data),
    onSuccess: (_, variables) => {
      // 使特定訂閱和所有列表查詢失效
      queryClient.invalidateQueries({
        queryKey: subscriptionKeys.detail(variables.id),
      });
      queryClient.invalidateQueries({ queryKey: subscriptionKeys.lists() });
    },
    ...options,
  });
}

/**
 * 刪除訂閱
 *
 * @param options Mutation 選項
 * @returns 刪除訂閱的 mutation
 *
 * @example
 * ```tsx
 * const deleteMutation = useDeleteSubscription();
 *
 * const handleDelete = async (id: string) => {
 *   await deleteMutation.mutateAsync(id);
 * };
 * ```
 */
export function useDeleteSubscription(
  options?: UseMutationOptions<void, APIError, string, unknown>
) {
  const queryClient = useQueryClient();

  return useMutation<void, APIError, string>({
    mutationFn: subscriptionsAPI.delete,
    onSuccess: (_, id) => {
      // 移除特定訂閱的快取並使列表失效
      queryClient.removeQueries({ queryKey: subscriptionKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: subscriptionKeys.lists() });
    },
    ...options,
  });
}

/**
 * 取消訂閱
 *
 * @param options Mutation 選項
 * @returns 取消訂閱的 mutation
 *
 * @example
 * ```tsx
 * const cancelMutation = useCancelSubscription();
 *
 * const handleCancel = async (id: string, endDate: string) => {
 *   await cancelMutation.mutateAsync({ id, data: { end_date: endDate } });
 * };
 * ```
 */
export function useCancelSubscription(
  options?: UseMutationOptions<
    Subscription,
    APIError,
    { id: string; data: CancelSubscriptionInput },
    unknown
  >
) {
  const queryClient = useQueryClient();

  return useMutation<
    Subscription,
    APIError,
    { id: string; data: CancelSubscriptionInput }
  >({
    mutationFn: ({ id, data }) => subscriptionsAPI.cancel(id, data),
    onSuccess: (_, variables) => {
      // 使特定訂閱和所有列表查詢失效
      queryClient.invalidateQueries({
        queryKey: subscriptionKeys.detail(variables.id),
      });
      queryClient.invalidateQueries({ queryKey: subscriptionKeys.lists() });
    },
    ...options,
  });
}

