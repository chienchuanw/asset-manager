import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryOptions,
  type UseMutationOptions,
} from "@tanstack/react-query";
import { bankAccountsAPI } from "@/lib/api/user-management";
import type {
  BankAccount,
  CreateBankAccountInput,
  UpdateBankAccountInput,
  BankAccountFilters,
} from "@/types/user-management";
import { APIError } from "@/lib/api/client";

/**
 * Query Keys
 * 用於識別和管理快取
 */
export const bankAccountKeys = {
  all: ["bankAccounts"] as const,
  lists: () => [...bankAccountKeys.all, "list"] as const,
  list: (filters?: BankAccountFilters) =>
    [...bankAccountKeys.lists(), filters] as const,
  details: () => [...bankAccountKeys.all, "detail"] as const,
  detail: (id: string) => [...bankAccountKeys.details(), id] as const,
};

// ==================== 銀行帳戶 Hooks ====================

/**
 * 取得銀行帳戶列表
 *
 * @param filters 篩選條件
 * @param options React Query 選項
 * @returns 銀行帳戶列表查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useBankAccounts({
 *   currency: 'TWD',
 * });
 * ```
 */
export function useBankAccounts(
  filters?: BankAccountFilters,
  options?: Omit<
    UseQueryOptions<BankAccount[], APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<BankAccount[], APIError>({
    queryKey: bankAccountKeys.list(filters),
    queryFn: () => bankAccountsAPI.getAll(filters),
    ...options,
  });
}

/**
 * 取得單筆銀行帳戶
 *
 * @param id 銀行帳戶 ID
 * @param options React Query 選項
 * @returns 單筆銀行帳戶查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useBankAccount("account-id");
 * ```
 */
export function useBankAccount(
  id: string,
  options?: Omit<
    UseQueryOptions<BankAccount, APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<BankAccount, APIError>({
    queryKey: bankAccountKeys.detail(id),
    queryFn: () => bankAccountsAPI.getById(id),
    enabled: !!id,
    ...options,
  });
}

/**
 * 建立銀行帳戶
 *
 * @param options Mutation 選項
 * @returns 建立銀行帳戶的 mutation
 *
 * @example
 * ```tsx
 * const createMutation = useCreateBankAccount();
 *
 * const handleCreate = async (data: CreateBankAccountInput) => {
 *   await createMutation.mutateAsync(data);
 * };
 * ```
 */
export function useCreateBankAccount(
  options?: UseMutationOptions<
    BankAccount,
    APIError,
    CreateBankAccountInput,
    unknown
  >
) {
  const queryClient = useQueryClient();

  return useMutation<BankAccount, APIError, CreateBankAccountInput>({
    mutationFn: bankAccountsAPI.create,
    onSuccess: () => {
      // 使所有銀行帳戶列表查詢失效，觸發重新獲取
      queryClient.invalidateQueries({ queryKey: bankAccountKeys.lists() });
    },
    ...options,
  });
}

/**
 * 更新銀行帳戶
 *
 * @param options Mutation 選項
 * @returns 更新銀行帳戶的 mutation
 *
 * @example
 * ```tsx
 * const updateMutation = useUpdateBankAccount();
 *
 * const handleUpdate = async (id: string, data: UpdateBankAccountInput) => {
 *   await updateMutation.mutateAsync({ id, data });
 * };
 * ```
 */
export function useUpdateBankAccount(
  options?: UseMutationOptions<
    BankAccount,
    APIError,
    { id: string; data: UpdateBankAccountInput },
    unknown
  >
) {
  const queryClient = useQueryClient();

  return useMutation<
    BankAccount,
    APIError,
    { id: string; data: UpdateBankAccountInput }
  >({
    mutationFn: ({ id, data }) => bankAccountsAPI.update(id, data),
    onSuccess: (_, variables) => {
      // 使所有銀行帳戶列表查詢失效
      queryClient.invalidateQueries({ queryKey: bankAccountKeys.lists() });
      // 使特定銀行帳戶詳情查詢失效
      queryClient.invalidateQueries({
        queryKey: bankAccountKeys.detail(variables.id),
      });
    },
    ...options,
  });
}

/**
 * 刪除銀行帳戶
 *
 * @param options Mutation 選項
 * @returns 刪除銀行帳戶的 mutation
 *
 * @example
 * ```tsx
 * const deleteMutation = useDeleteBankAccount();
 *
 * const handleDelete = async (id: string) => {
 *   await deleteMutation.mutateAsync(id);
 * };
 * ```
 */
export function useDeleteBankAccount(
  options?: UseMutationOptions<void, APIError, string, unknown>
) {
  const queryClient = useQueryClient();

  return useMutation<void, APIError, string>({
    mutationFn: bankAccountsAPI.delete,
    onSuccess: (_, id) => {
      // 使所有銀行帳戶列表查詢失效
      queryClient.invalidateQueries({ queryKey: bankAccountKeys.lists() });
      // 移除特定銀行帳戶詳情查詢
      queryClient.removeQueries({ queryKey: bankAccountKeys.detail(id) });
    },
    ...options,
  });
}

