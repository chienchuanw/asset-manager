import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryOptions,
  type UseMutationOptions,
} from "@tanstack/react-query";
import { transactionsAPI } from "@/lib/api/transactions";
import type {
  Transaction,
  CreateTransactionInput,
  UpdateTransactionInput,
  TransactionFilters,
  BatchCreateTransactionsInput,
} from "@/types/transaction";
import { APIError } from "@/lib/api/client";

/**
 * Query Keys
 * 用於識別和管理快取
 */
export const transactionKeys = {
  all: ["transactions"] as const,
  lists: () => [...transactionKeys.all, "list"] as const,
  list: (filters?: TransactionFilters) =>
    [...transactionKeys.lists(), filters] as const,
  details: () => [...transactionKeys.all, "detail"] as const,
  detail: (id: string) => [...transactionKeys.details(), id] as const,
};

/**
 * 取得交易列表
 *
 * @param filters 篩選條件
 * @param options React Query 選項
 * @returns 交易列表查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useTransactions({
 *   asset_type: "tw-stock",
 *   limit: 10,
 * });
 * ```
 */
export function useTransactions(
  filters?: TransactionFilters,
  options?: Omit<
    UseQueryOptions<Transaction[], APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<Transaction[], APIError>({
    queryKey: transactionKeys.list(filters),
    queryFn: () => transactionsAPI.getAll(filters),
    ...options,
  });
}

/**
 * 取得單筆交易
 *
 * @param id 交易 ID
 * @param options React Query 選項
 * @returns 單筆交易查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useTransaction("transaction-id");
 * ```
 */
export function useTransaction(
  id: string,
  options?: Omit<UseQueryOptions<Transaction, APIError>, "queryKey" | "queryFn">
) {
  return useQuery<Transaction, APIError>({
    queryKey: transactionKeys.detail(id),
    queryFn: () => transactionsAPI.getById(id),
    enabled: !!id, // 只有當 id 存在時才執行查詢
    ...options,
  });
}

/**
 * 建立交易
 *
 * @param options React Query mutation 選項
 * @returns 建立交易的 mutation
 *
 * @example
 * ```tsx
 * const createMutation = useCreateTransaction({
 *   onSuccess: () => {
 *     toast.success("交易建立成功");
 *   },
 *   onError: (error) => {
 *     toast.error(error.message);
 *   },
 * });
 *
 * createMutation.mutate({
 *   date: "2025-10-23T00:00:00Z",
 *   asset_type: "tw-stock",
 *   symbol: "2330",
 *   name: "台積電",
 *   type: "buy",
 *   quantity: 10,
 *   price: 620,
 *   amount: 6200,
 * });
 * ```
 */
export function useCreateTransaction(
  options?: UseMutationOptions<Transaction, APIError, CreateTransactionInput>
) {
  const queryClient = useQueryClient();

  return useMutation<Transaction, APIError, CreateTransactionInput>({
    mutationFn: transactionsAPI.create,
    onSuccess: async () => {
      // 使所有交易列表的快取失效，強制重新獲取
      await queryClient.invalidateQueries({
        queryKey: transactionKeys.lists(),
      });
    },
    ...options,
  });
}

/**
 * 更新交易
 *
 * @param options React Query mutation 選項
 * @returns 更新交易的 mutation
 *
 * @example
 * ```tsx
 * const updateMutation = useUpdateTransaction({
 *   onSuccess: () => {
 *     toast.success("交易更新成功");
 *   },
 * });
 *
 * updateMutation.mutate({
 *   id: "transaction-id",
 *   data: { quantity: 20 },
 * });
 * ```
 */
export function useUpdateTransaction(
  options?: UseMutationOptions<
    Transaction,
    APIError,
    { id: string; data: UpdateTransactionInput }
  >
) {
  const queryClient = useQueryClient();

  return useMutation<
    Transaction,
    APIError,
    { id: string; data: UpdateTransactionInput }
  >({
    mutationFn: ({ id, data }) => transactionsAPI.update(id, data),
    onSuccess: async (_data, variables) => {
      // 使所有交易列表的快取失效
      await queryClient.invalidateQueries({
        queryKey: transactionKeys.lists(),
      });

      // 使該筆交易的快取失效
      await queryClient.invalidateQueries({
        queryKey: transactionKeys.detail(variables.id),
      });
    },
    ...options,
  });
}

/**
 * 批次建立交易
 *
 * @param options React Query mutation 選項
 * @returns 批次建立交易的 mutation
 *
 * @example
 * ```tsx
 * const createBatchMutation = useCreateTransactionsBatch({
 *   onSuccess: () => {
 *     toast.success("批次交易建立成功");
 *   },
 *   onError: (error) => {
 *     toast.error(error.message);
 *   },
 * });
 *
 * createBatchMutation.mutate({
 *   transactions: [
 *     {
 *       date: "2025-10-23T00:00:00Z",
 *       asset_type: "tw-stock",
 *       symbol: "2330",
 *       name: "台積電",
 *       type: "buy",
 *       quantity: 10,
 *       price: 620,
 *       amount: 6200,
 *       currency: "TWD",
 *     },
 *     // ... 更多交易
 *   ],
 * });
 * ```
 */
export function useCreateTransactionsBatch(
  options?: UseMutationOptions<
    Transaction[],
    APIError,
    BatchCreateTransactionsInput
  >
) {
  const queryClient = useQueryClient();

  return useMutation<Transaction[], APIError, BatchCreateTransactionsInput>({
    mutationFn: transactionsAPI.createBatch,
    onSuccess: async () => {
      // 使所有交易列表的快取失效，強制重新獲取
      await queryClient.invalidateQueries({
        queryKey: transactionKeys.lists(),
      });
    },
    ...options,
  });
}

/**
 * 刪除交易
 *
 * @param options React Query mutation 選項
 * @returns 刪除交易的 mutation
 *
 * @example
 * ```tsx
 * const deleteMutation = useDeleteTransaction({
 *   onSuccess: () => {
 *     toast.success("交易刪除成功");
 *   },
 * });
 *
 * deleteMutation.mutate("transaction-id");
 * ```
 */
export function useDeleteTransaction(
  options?: UseMutationOptions<void, APIError, string>
) {
  const queryClient = useQueryClient();

  return useMutation<void, APIError, string>({
    mutationFn: transactionsAPI.delete,
    onSuccess: async (_data, variables) => {
      // 使所有交易列表的快取失效
      await queryClient.invalidateQueries({
        queryKey: transactionKeys.lists(),
      });

      // 移除該筆交易的快取
      queryClient.removeQueries({
        queryKey: transactionKeys.detail(variables),
      });
    },
    ...options,
  });
}

/**
 * 樂觀更新：建立交易
 *
 * 在伺服器回應之前先更新 UI，提供更好的使用者體驗
 *
 * @param options React Query mutation 選項
 * @returns 建立交易的 mutation（含樂觀更新）
 */
export function useCreateTransactionOptimistic(
  options?: UseMutationOptions<
    Transaction,
    APIError,
    CreateTransactionInput,
    { previousTransactions?: Transaction[] }
  >
) {
  const queryClient = useQueryClient();

  return useMutation<
    Transaction,
    APIError,
    CreateTransactionInput,
    { previousTransactions?: Transaction[] }
  >({
    mutationFn: transactionsAPI.create,
    onMutate: async (newTransaction) => {
      // 取消所有進行中的查詢（避免覆蓋樂觀更新）
      await queryClient.cancelQueries({
        queryKey: transactionKeys.lists(),
      });

      // 儲存之前的資料（用於錯誤回滾）
      const previousTransactions = queryClient.getQueryData<Transaction[]>(
        transactionKeys.lists()
      );

      // 樂觀更新：立即加入新交易到列表
      queryClient.setQueryData<Transaction[]>(
        transactionKeys.lists(),
        (old) => {
          const optimisticTransaction: Transaction = {
            id: `temp-${Date.now()}`, // 暫時 ID
            ...newTransaction,
            fee: newTransaction.fee ?? null, // 確保 fee 是 number | null
            tax: newTransaction.tax ?? null, // 確保 tax 是 number | null
            exchange_rate_id: null, // 樂觀更新時先設為 null，實際值由後端決定
            note: newTransaction.note ?? null, // 確保 note 是 string | null
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          };
          return old
            ? [optimisticTransaction, ...old]
            : [optimisticTransaction];
        }
      );

      // 回傳 context（用於錯誤回滾）
      return { previousTransactions };
    },
    onError: (_error, _variables, context) => {
      // 發生錯誤時，回滾到之前的資料
      if (context?.previousTransactions) {
        queryClient.setQueryData(
          transactionKeys.lists(),
          context.previousTransactions
        );
      }
    },
    onSettled: async () => {
      // 無論成功或失敗，都重新獲取資料以確保同步
      await queryClient.invalidateQueries({
        queryKey: transactionKeys.lists(),
      });
    },
    ...options,
  });
}
