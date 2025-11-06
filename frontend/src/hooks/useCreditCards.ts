import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryOptions,
  type UseMutationOptions,
} from "@tanstack/react-query";
import { creditCardsAPI } from "@/lib/api/user-management";
import type {
  CreditCard,
  CreateCreditCardInput,
  UpdateCreditCardInput,
  CreditCardQueryParams,
} from "@/types/user-management";
import { APIError } from "@/lib/api/client";

/**
 * Query Keys
 * 用於識別和管理快取
 */
export const creditCardKeys = {
  all: ["creditCards"] as const,
  lists: () => [...creditCardKeys.all, "list"] as const,
  list: () => [...creditCardKeys.lists()] as const,
  upcomingBilling: (params?: CreditCardQueryParams) =>
    [...creditCardKeys.all, "upcomingBilling", params] as const,
  upcomingPayment: (params?: CreditCardQueryParams) =>
    [...creditCardKeys.all, "upcomingPayment", params] as const,
  details: () => [...creditCardKeys.all, "detail"] as const,
  detail: (id: string) => [...creditCardKeys.details(), id] as const,
};

// ==================== 信用卡 Hooks ====================

/**
 * 取得信用卡列表
 *
 * @param options React Query 選項
 * @returns 信用卡列表查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useCreditCards();
 * ```
 */
export function useCreditCards(
  options?: Omit<
    UseQueryOptions<CreditCard[], APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<CreditCard[], APIError>({
    queryKey: creditCardKeys.list(),
    queryFn: () => creditCardsAPI.getAll(),
    staleTime: 0, // 強制每次都重新驗證資料
    ...options,
  });
}

/**
 * 取得單筆信用卡
 *
 * @param id 信用卡 ID
 * @param options React Query 選項
 * @returns 單筆信用卡查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useCreditCard("card-id");
 * ```
 */
export function useCreditCard(
  id: string,
  options?: Omit<UseQueryOptions<CreditCard, APIError>, "queryKey" | "queryFn">
) {
  return useQuery<CreditCard, APIError>({
    queryKey: creditCardKeys.detail(id),
    queryFn: () => creditCardsAPI.getById(id),
    enabled: !!id,
    ...options,
  });
}

/**
 * 取得即將到來的帳單日信用卡
 *
 * @param params 查詢參數
 * @param options React Query 選項
 * @returns 信用卡列表查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useUpcomingBillingCreditCards({
 *   days_ahead: 7,
 * });
 * ```
 */
export function useUpcomingBillingCreditCards(
  params?: CreditCardQueryParams,
  options?: Omit<
    UseQueryOptions<CreditCard[], APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<CreditCard[], APIError>({
    queryKey: creditCardKeys.upcomingBilling(params),
    queryFn: () => creditCardsAPI.getUpcomingBilling(params),
    ...options,
  });
}

/**
 * 取得即將到來的繳款截止日信用卡
 *
 * @param params 查詢參數
 * @param options React Query 選項
 * @returns 信用卡列表查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useUpcomingPaymentCreditCards({
 *   days_ahead: 7,
 * });
 * ```
 */
export function useUpcomingPaymentCreditCards(
  params?: CreditCardQueryParams,
  options?: Omit<
    UseQueryOptions<CreditCard[], APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<CreditCard[], APIError>({
    queryKey: creditCardKeys.upcomingPayment(params),
    queryFn: () => creditCardsAPI.getUpcomingPayment(params),
    ...options,
  });
}

/**
 * 建立信用卡
 *
 * @param options Mutation 選項
 * @returns 建立信用卡的 mutation
 *
 * @example
 * ```tsx
 * const createMutation = useCreateCreditCard();
 *
 * const handleCreate = async (data: CreateCreditCardInput) => {
 *   await createMutation.mutateAsync(data);
 * };
 * ```
 */
export function useCreateCreditCard(
  options?: UseMutationOptions<
    CreditCard,
    APIError,
    CreateCreditCardInput,
    unknown
  >
) {
  const queryClient = useQueryClient();

  return useMutation<CreditCard, APIError, CreateCreditCardInput>({
    mutationFn: creditCardsAPI.create,
    onSuccess: () => {
      // 使所有信用卡列表查詢失效，觸發重新獲取
      queryClient.invalidateQueries({ queryKey: creditCardKeys.lists() });
      queryClient.invalidateQueries({ queryKey: creditCardKeys.all });
    },
    ...options,
  });
}

/**
 * 更新信用卡
 *
 * @param options Mutation 選項
 * @returns 更新信用卡的 mutation
 *
 * @example
 * ```tsx
 * const updateMutation = useUpdateCreditCard();
 *
 * const handleUpdate = async (id: string, data: UpdateCreditCardInput) => {
 *   await updateMutation.mutateAsync({ id, data });
 * };
 * ```
 */
export function useUpdateCreditCard(
  options?: UseMutationOptions<
    CreditCard,
    APIError,
    { id: string; data: UpdateCreditCardInput },
    unknown
  >
) {
  const queryClient = useQueryClient();

  return useMutation<
    CreditCard,
    APIError,
    { id: string; data: UpdateCreditCardInput }
  >({
    mutationFn: ({ id, data }) => creditCardsAPI.update(id, data),
    onSuccess: (_, variables) => {
      // 使所有信用卡列表查詢失效
      queryClient.invalidateQueries({ queryKey: creditCardKeys.lists() });
      queryClient.invalidateQueries({ queryKey: creditCardKeys.all });
      // 使特定信用卡詳情查詢失效
      queryClient.invalidateQueries({
        queryKey: creditCardKeys.detail(variables.id),
      });
    },
    ...options,
  });
}

/**
 * 刪除信用卡
 *
 * @param options Mutation 選項
 * @returns 刪除信用卡的 mutation
 *
 * @example
 * ```tsx
 * const deleteMutation = useDeleteCreditCard();
 *
 * const handleDelete = async (id: string) => {
 *   await deleteMutation.mutateAsync(id);
 * };
 * ```
 */
export function useDeleteCreditCard(
  options?: UseMutationOptions<void, APIError, string, unknown>
) {
  const queryClient = useQueryClient();

  return useMutation<void, APIError, string>({
    mutationFn: creditCardsAPI.delete,
    onSuccess: (_, id) => {
      // 使所有信用卡列表查詢失效
      queryClient.invalidateQueries({ queryKey: creditCardKeys.lists() });
      queryClient.invalidateQueries({ queryKey: creditCardKeys.all });
      // 移除特定信用卡詳情查詢
      queryClient.removeQueries({ queryKey: creditCardKeys.detail(id) });
    },
    ...options,
  });
}
