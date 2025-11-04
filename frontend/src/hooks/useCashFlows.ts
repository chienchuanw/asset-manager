import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryOptions,
  type UseMutationOptions,
} from "@tanstack/react-query";
import { cashFlowsAPI, categoriesAPI } from "@/lib/api/cash-flows";
import type {
  CashFlow,
  CashFlowCategory,
  CashFlowSummary,
  CreateCashFlowInput,
  UpdateCashFlowInput,
  CreateCategoryInput,
  UpdateCategoryInput,
  CashFlowFilters,
  CashFlowType,
} from "@/types/cash-flow";
import { APIError } from "@/lib/api/client";
import { toast } from "sonner";

/**
 * Query Keys
 * 用於識別和管理快取
 */
export const cashFlowKeys = {
  all: ["cash-flows"] as const,
  lists: () => [...cashFlowKeys.all, "list"] as const,
  list: (filters?: CashFlowFilters) =>
    [...cashFlowKeys.lists(), filters] as const,
  details: () => [...cashFlowKeys.all, "detail"] as const,
  detail: (id: string) => [...cashFlowKeys.details(), id] as const,
  summaries: () => [...cashFlowKeys.all, "summary"] as const,
  summary: (startDate: string, endDate: string) =>
    [...cashFlowKeys.summaries(), { startDate, endDate }] as const,
};

export const categoryKeys = {
  all: ["categories"] as const,
  lists: () => [...categoryKeys.all, "list"] as const,
  list: (type?: CashFlowType) => [...categoryKeys.lists(), type] as const,
  details: () => [...categoryKeys.all, "detail"] as const,
  detail: (id: string) => [...categoryKeys.details(), id] as const,
};

// ==================== 現金流記錄 Hooks ====================

/**
 * 取得現金流記錄列表
 *
 * @param filters 篩選條件
 * @param options React Query 選項
 * @returns 現金流記錄列表查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useCashFlows({
 *   type: CashFlowType.INCOME,
 *   limit: 10,
 * });
 * ```
 */
export function useCashFlows(
  filters?: CashFlowFilters,
  options?: Omit<UseQueryOptions<CashFlow[], APIError>, "queryKey" | "queryFn">
) {
  const queryKey = cashFlowKeys.list(filters);

  return useQuery<CashFlow[], APIError>({
    queryKey,
    queryFn: () => {
      return cashFlowsAPI.getAll(filters);
    },
    ...options,
  });
}

/**
 * 取得單筆現金流記錄
 *
 * @param id 現金流記錄 ID
 * @param options React Query 選項
 * @returns 單筆現金流記錄查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useCashFlow("cash-flow-id");
 * ```
 */
export function useCashFlow(
  id: string,
  options?: Omit<UseQueryOptions<CashFlow, APIError>, "queryKey" | "queryFn">
) {
  return useQuery<CashFlow, APIError>({
    queryKey: cashFlowKeys.detail(id),
    queryFn: () => cashFlowsAPI.getById(id),
    enabled: !!id,
    ...options,
  });
}

/**
 * 取得現金流摘要
 *
 * @param startDate 開始日期 (YYYY-MM-DD)
 * @param endDate 結束日期 (YYYY-MM-DD)
 * @param options React Query 選項
 * @returns 現金流摘要查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useCashFlowSummary("2025-10-01", "2025-10-31");
 * ```
 */
export function useCashFlowSummary(
  startDate: string,
  endDate: string,
  options?: Omit<
    UseQueryOptions<CashFlowSummary, APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<CashFlowSummary, APIError>({
    queryKey: cashFlowKeys.summary(startDate, endDate),
    queryFn: () => cashFlowsAPI.getSummary(startDate, endDate),
    enabled: !!startDate && !!endDate,
    ...options,
  });
}

/**
 * 建立現金流記錄
 *
 * @param options React Query mutation 選項
 * @returns 建立現金流記錄的 mutation
 *
 * @example
 * ```tsx
 * const createMutation = useCreateCashFlow({
 *   onSuccess: () => {
 *     toast.success("記錄建立成功");
 *   },
 *   onError: (error) => {
 *     toast.error(error.message);
 *   },
 * });
 *
 * createMutation.mutate({
 *   date: "2025-10-25",
 *   type: CashFlowType.INCOME,
 *   category_id: "category-id",
 *   amount: 50000,
 *   description: "十月薪資",
 * });
 * ```
 */
export function useCreateCashFlow(
  options?: UseMutationOptions<CashFlow, APIError, CreateCashFlowInput>
) {
  const queryClient = useQueryClient();

  return useMutation<CashFlow, APIError, CreateCashFlowInput>({
    mutationFn: cashFlowsAPI.create,
    onSuccess: async () => {
      // 使所有現金流相關查詢失效
      await queryClient.invalidateQueries({
        queryKey: cashFlowKeys.all,
      });
      // 使銀行帳戶列表失效（餘額可能已更新）
      await queryClient.invalidateQueries({
        queryKey: ["bankAccounts"],
      });
      // 使信用卡列表失效（餘額可能已更新）
      await queryClient.invalidateQueries({
        queryKey: ["creditCards"],
      });
    },
    ...options,
  });
}

/**
 * 更新現金流記錄
 *
 * @param options React Query mutation 選項
 * @returns 更新現金流記錄的 mutation
 *
 * @example
 * ```tsx
 * const updateMutation = useUpdateCashFlow({
 *   onSuccess: () => {
 *     toast.success("記錄更新成功");
 *   },
 * });
 *
 * updateMutation.mutate({
 *   id: "cash-flow-id",
 *   data: { amount: 55000 },
 * });
 * ```
 */
export function useUpdateCashFlow(
  options?: UseMutationOptions<
    CashFlow,
    APIError,
    { id: string; data: UpdateCashFlowInput }
  >
) {
  const queryClient = useQueryClient();

  return useMutation<
    CashFlow,
    APIError,
    { id: string; data: UpdateCashFlowInput }
  >({
    mutationFn: ({ id, data }) => cashFlowsAPI.update(id, data),
    onSuccess: async () => {
      // 顯示成功訊息
      toast.success("記錄更新成功");

      // 使所有現金流相關查詢失效
      await queryClient.invalidateQueries({
        queryKey: cashFlowKeys.all,
      });
      // 使銀行帳戶列表失效（餘額可能已更新）
      await queryClient.invalidateQueries({
        queryKey: ["bankAccounts"],
      });
      // 使信用卡列表失效（餘額可能已更新）
      await queryClient.invalidateQueries({
        queryKey: ["creditCards"],
      });
    },
    ...options,
  });
}

/**
 * 刪除現金流記錄
 *
 * @param options React Query mutation 選項
 * @returns 刪除現金流記錄的 mutation
 *
 * @example
 * ```tsx
 * const deleteMutation = useDeleteCashFlow({
 *   onSuccess: () => {
 *     toast.success("記錄刪除成功");
 *   },
 * });
 *
 * deleteMutation.mutate("cash-flow-id");
 * ```
 */
export function useDeleteCashFlow(
  options?: UseMutationOptions<void, APIError, string>
) {
  const queryClient = useQueryClient();

  return useMutation<void, APIError, string>({
    mutationFn: cashFlowsAPI.delete,
    onSuccess: async () => {
      // 顯示成功訊息
      toast.success("記錄刪除成功");

      // 使所有現金流相關查詢失效
      await queryClient.invalidateQueries({
        queryKey: cashFlowKeys.all,
      });
      // 使銀行帳戶列表失效（餘額可能已更新）
      await queryClient.invalidateQueries({
        queryKey: ["bankAccounts"],
      });
      // 使信用卡列表失效（餘額可能已更新）
      await queryClient.invalidateQueries({
        queryKey: ["creditCards"],
      });
    },
    ...options,
  });
}

// ==================== 分類 Hooks ====================

/**
 * 取得分類列表
 *
 * @param type 現金流類型（可選）
 * @param options React Query 選項
 * @returns 分類列表查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useCategories(CashFlowType.INCOME);
 * ```
 */
export function useCategories(
  type?: CashFlowType,
  options?: Omit<
    UseQueryOptions<CashFlowCategory[], APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<CashFlowCategory[], APIError>({
    queryKey: categoryKeys.list(type),
    queryFn: () => categoriesAPI.getAll(type),
    ...options,
  });
}

/**
 * 取得單筆分類
 *
 * @param id 分類 ID
 * @param options React Query 選項
 * @returns 單筆分類查詢結果
 *
 * @example
 * ```tsx
 * const { data, isLoading, error } = useCategory("category-id");
 * ```
 */
export function useCategory(
  id: string,
  options?: Omit<
    UseQueryOptions<CashFlowCategory, APIError>,
    "queryKey" | "queryFn"
  >
) {
  return useQuery<CashFlowCategory, APIError>({
    queryKey: categoryKeys.detail(id),
    queryFn: () => categoriesAPI.getById(id),
    enabled: !!id,
    ...options,
  });
}

/**
 * 建立分類
 *
 * @param options React Query mutation 選項
 * @returns 建立分類的 mutation
 *
 * @example
 * ```tsx
 * const createMutation = useCreateCategory({
 *   onSuccess: () => {
 *     toast.success("分類建立成功");
 *   },
 * });
 *
 * createMutation.mutate({
 *   name: "投資收入",
 *   type: CashFlowType.INCOME,
 * });
 * ```
 */
export function useCreateCategory(
  options?: UseMutationOptions<CashFlowCategory, APIError, CreateCategoryInput>
) {
  const queryClient = useQueryClient();

  return useMutation<CashFlowCategory, APIError, CreateCategoryInput>({
    mutationFn: categoriesAPI.create,
    onSuccess: async () => {
      // 使所有分類列表的快取失效
      await queryClient.invalidateQueries({
        queryKey: categoryKeys.lists(),
      });
    },
    ...options,
  });
}

/**
 * 更新分類
 *
 * @param options React Query mutation 選項
 * @returns 更新分類的 mutation
 *
 * @example
 * ```tsx
 * const updateMutation = useUpdateCategory({
 *   onSuccess: () => {
 *     toast.success("分類更新成功");
 *   },
 * });
 *
 * updateMutation.mutate({
 *   id: "category-id",
 *   data: { name: "新的分類名稱" },
 * });
 * ```
 */
export function useUpdateCategory(
  options?: UseMutationOptions<
    CashFlowCategory,
    APIError,
    { id: string; data: UpdateCategoryInput }
  >
) {
  const queryClient = useQueryClient();

  return useMutation<
    CashFlowCategory,
    APIError,
    { id: string; data: UpdateCategoryInput }
  >({
    mutationFn: ({ id, data }) => categoriesAPI.update(id, data),
    onSuccess: async (_data, variables) => {
      // 使所有分類列表的快取失效
      await queryClient.invalidateQueries({
        queryKey: categoryKeys.lists(),
      });
      // 使該筆分類的快取失效
      await queryClient.invalidateQueries({
        queryKey: categoryKeys.detail(variables.id),
      });
    },
    ...options,
  });
}

/**
 * 刪除分類
 *
 * @param options React Query mutation 選項
 * @returns 刪除分類的 mutation
 *
 * @example
 * ```tsx
 * const deleteMutation = useDeleteCategory({
 *   onSuccess: () => {
 *     toast.success("分類刪除成功");
 *   },
 * });
 *
 * deleteMutation.mutate("category-id");
 * ```
 */
export function useDeleteCategory(
  options?: UseMutationOptions<void, APIError, string>
) {
  const queryClient = useQueryClient();

  return useMutation<void, APIError, string>({
    mutationFn: categoriesAPI.delete,
    onSuccess: async (_data, variables) => {
      // 使所有分類列表的快取失效
      await queryClient.invalidateQueries({
        queryKey: categoryKeys.lists(),
      });
      // 移除該筆分類的快取
      queryClient.removeQueries({
        queryKey: categoryKeys.detail(variables),
      });
    },
    ...options,
  });
}
