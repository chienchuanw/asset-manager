import {
  useMutation,
  useQueryClient,
  type UseMutationOptions,
} from "@tanstack/react-query";
import { billingAPI } from "@/lib/api/billing";
import type {
  BillingResult,
  DailyBillingResult,
  ProcessBillingInput,
} from "@/types/billing";
import { APIError } from "@/lib/api/client";
import { subscriptionKeys } from "./useSubscriptions";
import { installmentKeys } from "./useInstallments";
import { cashFlowKeys } from "./useCashFlows";

// ==================== 扣款 Hooks ====================

/**
 * 處理每日扣款
 *
 * @param options Mutation 選項
 * @returns 處理每日扣款的 mutation
 *
 * @example
 * ```tsx
 * const processDailyMutation = useProcessDailyBilling();
 *
 * const handleProcessDaily = async () => {
 *   const result = await processDailyMutation.mutateAsync();
 *   console.log(`處理了 ${result.subscription_count} 筆訂閱和 ${result.installment_count} 筆分期`);
 * };
 * ```
 */
export function useProcessDailyBilling(
  options?: UseMutationOptions<
    DailyBillingResult,
    APIError,
    ProcessBillingInput | undefined,
    unknown
  >
) {
  const queryClient = useQueryClient();

  return useMutation<DailyBillingResult, APIError, ProcessBillingInput | undefined>({
    mutationFn: (data) => billingAPI.processDaily(data),
    onSuccess: () => {
      // 使所有相關查詢失效
      queryClient.invalidateQueries({ queryKey: subscriptionKeys.all });
      queryClient.invalidateQueries({ queryKey: installmentKeys.all });
      queryClient.invalidateQueries({ queryKey: cashFlowKeys.all });
    },
    ...options,
  });
}

/**
 * 處理訂閱扣款
 *
 * @param options Mutation 選項
 * @returns 處理訂閱扣款的 mutation
 *
 * @example
 * ```tsx
 * const processSubscriptionsMutation = useProcessSubscriptionBilling();
 *
 * const handleProcessSubscriptions = async () => {
 *   const result = await processSubscriptionsMutation.mutateAsync();
 *   console.log(`處理了 ${result.count} 筆訂閱扣款`);
 * };
 * ```
 */
export function useProcessSubscriptionBilling(
  options?: UseMutationOptions<
    BillingResult,
    APIError,
    ProcessBillingInput | undefined,
    unknown
  >
) {
  const queryClient = useQueryClient();

  return useMutation<BillingResult, APIError, ProcessBillingInput | undefined>({
    mutationFn: (data) => billingAPI.processSubscriptions(data),
    onSuccess: () => {
      // 使訂閱和現金流查詢失效
      queryClient.invalidateQueries({ queryKey: subscriptionKeys.all });
      queryClient.invalidateQueries({ queryKey: cashFlowKeys.all });
    },
    ...options,
  });
}

/**
 * 處理分期扣款
 *
 * @param options Mutation 選項
 * @returns 處理分期扣款的 mutation
 *
 * @example
 * ```tsx
 * const processInstallmentsMutation = useProcessInstallmentBilling();
 *
 * const handleProcessInstallments = async () => {
 *   const result = await processInstallmentsMutation.mutateAsync();
 *   console.log(`處理了 ${result.count} 筆分期扣款`);
 * };
 * ```
 */
export function useProcessInstallmentBilling(
  options?: UseMutationOptions<
    BillingResult,
    APIError,
    ProcessBillingInput | undefined,
    unknown
  >
) {
  const queryClient = useQueryClient();

  return useMutation<BillingResult, APIError, ProcessBillingInput | undefined>({
    mutationFn: (data) => billingAPI.processInstallments(data),
    onSuccess: () => {
      // 使分期和現金流查詢失效
      queryClient.invalidateQueries({ queryKey: installmentKeys.all });
      queryClient.invalidateQueries({ queryKey: cashFlowKeys.all });
    },
    ...options,
  });
}

