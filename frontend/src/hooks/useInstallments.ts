import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryOptions,
  type UseMutationOptions,
} from "@tanstack/react-query";
import { installmentsAPI } from "@/lib/api/installments";
import type {
  Installment,
  CreateInstallmentInput,
  UpdateInstallmentInput,
  InstallmentFilters,
  CompletingSoonParams,
} from "@/types/installment";
import { APIError } from "@/lib/api/client";

/**
 * Query Keys
 * 用於識別和管理快取
 */
export const installmentKeys = {
  all: ["installments"] as const,
  lists: () => [...installmentKeys.all, "list"] as const,
  list: (filters?: InstallmentFilters) =>
    [...installmentKeys.lists(), filters] as const,
  details: () => [...installmentKeys.all, "detail"] as const,
  detail: (id: string) => [...installmentKeys.details(), id] as const,
  completingSoon: (params?: CompletingSoonParams) =>
    [...installmentKeys.all, "completing-soon", params] as const,
};

// ==================== 分期 Hooks ====================

/**
 * 取得分期列表
 *
 * @param filters 篩選條件
 * @param options React Query 選項
 * @returns 分期列表查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useInstallments({
 *   status: 'active',
 *   limit: 10,
 * });
 * ```
 */
export function useInstallments(
  filters?: InstallmentFilters,
  options?: Omit<
    UseQueryOptions<Installment[], APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<Installment[], APIError>({
    queryKey: installmentKeys.list(filters),
    queryFn: () => installmentsAPI.getAll(filters),
    ...options,
  });
}

/**
 * 取得單筆分期
 *
 * @param id 分期 ID
 * @param options React Query 選項
 * @returns 單筆分期查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useInstallment("installment-id");
 * ```
 */
export function useInstallment(
  id: string,
  options?: Omit<UseQueryOptions<Installment, APIError>, "queryKey" | "queryFn">
) {
  return useQuery<Installment, APIError>({
    queryKey: installmentKeys.detail(id),
    queryFn: () => installmentsAPI.getById(id),
    enabled: !!id,
    ...options,
  });
}

/**
 * 取得即將完成的分期
 *
 * @param params 查詢參數
 * @param options React Query 選項
 * @returns 即將完成的分期列表查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useCompletingSoonInstallments({ months: 3 });
 * ```
 */
export function useCompletingSoonInstallments(
  params?: CompletingSoonParams,
  options?: Omit<
    UseQueryOptions<Installment[], APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<Installment[], APIError>({
    queryKey: installmentKeys.completingSoon(params),
    queryFn: () => installmentsAPI.getCompletingSoon(params),
    ...options,
  });
}

/**
 * 建立分期
 *
 * @param options Mutation 選項
 * @returns 建立分期的 mutation
 *
 * @example
 * ```tsx
 * const createMutation = useCreateInstallment();
 *
 * const handleCreate = async (data: CreateInstallmentInput) => {
 *   await createMutation.mutateAsync(data);
 * };
 * ```
 */
export function useCreateInstallment(
  options?: UseMutationOptions<
    Installment,
    APIError,
    CreateInstallmentInput,
    unknown
  >
) {
  const queryClient = useQueryClient();

  return useMutation<Installment, APIError, CreateInstallmentInput>({
    mutationFn: installmentsAPI.create,
    onSuccess: () => {
      // 使所有分期列表查詢失效，觸發重新獲取
      queryClient.invalidateQueries({ queryKey: installmentKeys.lists() });
      queryClient.invalidateQueries({ queryKey: installmentKeys.all });
    },
    ...options,
  });
}

/**
 * 更新分期
 *
 * @param options Mutation 選項
 * @returns 更新分期的 mutation
 *
 * @example
 * ```tsx
 * const updateMutation = useUpdateInstallment();
 *
 * const handleUpdate = async (id: string, data: UpdateInstallmentInput) => {
 *   await updateMutation.mutateAsync({ id, data });
 * };
 * ```
 */
export function useUpdateInstallment(
  options?: UseMutationOptions<
    Installment,
    APIError,
    { id: string; data: UpdateInstallmentInput },
    unknown
  >
) {
  const queryClient = useQueryClient();

  return useMutation<
    Installment,
    APIError,
    { id: string; data: UpdateInstallmentInput }
  >({
    mutationFn: ({ id, data }) => installmentsAPI.update(id, data),
    onSuccess: (_, variables) => {
      // 使特定分期和所有列表查詢失效
      queryClient.invalidateQueries({
        queryKey: installmentKeys.detail(variables.id),
      });
      queryClient.invalidateQueries({ queryKey: installmentKeys.lists() });
      queryClient.invalidateQueries({ queryKey: installmentKeys.all });
    },
    ...options,
  });
}

/**
 * 刪除分期
 *
 * @param options Mutation 選項
 * @returns 刪除分期的 mutation
 *
 * @example
 * ```tsx
 * const deleteMutation = useDeleteInstallment();
 *
 * const handleDelete = async (id: string) => {
 *   await deleteMutation.mutateAsync(id);
 * };
 * ```
 */
export function useDeleteInstallment(
  options?: UseMutationOptions<void, APIError, string, unknown>
) {
  const queryClient = useQueryClient();

  return useMutation<void, APIError, string>({
    mutationFn: installmentsAPI.delete,
    onSuccess: (_, id) => {
      // 移除特定分期的快取並使列表失效
      queryClient.removeQueries({ queryKey: installmentKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: installmentKeys.lists() });
      queryClient.invalidateQueries({ queryKey: installmentKeys.all });
    },
    ...options,
  });
}

