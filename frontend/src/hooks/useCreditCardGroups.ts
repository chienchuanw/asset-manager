import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { creditCardGroupsAPI } from "@/lib/api/user-management";
import type {
  CreateCreditCardGroupInput,
  UpdateCreditCardGroupInput,
  AddCardsToGroupInput,
  RemoveCardsFromGroupInput,
} from "@/types/user-management";
import { toast } from "sonner";

/**
 * Query keys for credit card groups
 */
export const creditCardGroupKeys = {
  all: ["creditCardGroups"] as const,
  lists: () => [...creditCardGroupKeys.all, "list"] as const,
  list: () => [...creditCardGroupKeys.lists()] as const,
  details: () => [...creditCardGroupKeys.all, "detail"] as const,
  detail: (id: string) => [...creditCardGroupKeys.details(), id] as const,
};

/**
 * Hook to fetch all credit card groups
 */
export function useCreditCardGroups() {
  return useQuery({
    queryKey: creditCardGroupKeys.list(),
    queryFn: () => creditCardGroupsAPI.getAll(),
  });
}

/**
 * Hook to fetch a single credit card group by ID
 */
export function useCreditCardGroup(id: string) {
  return useQuery({
    queryKey: creditCardGroupKeys.detail(id),
    queryFn: () => creditCardGroupsAPI.getById(id),
    enabled: !!id,
  });
}

/**
 * Hook to create a new credit card group
 */
export function useCreateCreditCardGroup() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateCreditCardGroupInput) =>
      creditCardGroupsAPI.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: creditCardGroupKeys.lists() });
      // Also invalidate credit cards list since cards now have group_id
      queryClient.invalidateQueries({ queryKey: ["creditCards"] });
      toast.success("信用卡群組建立成功");
    },
    onError: (error: Error) => {
      toast.error(`建立信用卡群組失敗: ${error.message}`);
    },
  });
}

/**
 * Hook to update a credit card group
 */
export function useUpdateCreditCardGroup() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      id,
      data,
    }: {
      id: string;
      data: UpdateCreditCardGroupInput;
    }) => creditCardGroupsAPI.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: creditCardGroupKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: creditCardGroupKeys.detail(variables.id),
      });
      toast.success("信用卡群組更新成功");
    },
    onError: (error: Error) => {
      toast.error(`更新信用卡群組失敗: ${error.message}`);
    },
  });
}

/**
 * Hook to delete a credit card group
 */
export function useDeleteCreditCardGroup() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => creditCardGroupsAPI.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: creditCardGroupKeys.lists() });
      // Also invalidate credit cards list since cards are now independent
      queryClient.invalidateQueries({ queryKey: ["creditCards"] });
      toast.success("信用卡群組刪除成功");
    },
    onError: (error: Error) => {
      toast.error(`刪除信用卡群組失敗: ${error.message}`);
    },
  });
}

/**
 * Hook to add cards to a group
 */
export function useAddCardsToGroup() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: AddCardsToGroupInput }) =>
      creditCardGroupsAPI.addCards(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: creditCardGroupKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: creditCardGroupKeys.detail(variables.id),
      });
      // Also invalidate credit cards list since cards now have group_id
      queryClient.invalidateQueries({ queryKey: ["creditCards"] });
      toast.success("卡片已加入群組");
    },
    onError: (error: Error) => {
      toast.error(`加入卡片失敗: ${error.message}`);
    },
  });
}

/**
 * Hook to remove cards from a group
 */
export function useRemoveCardsFromGroup() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      id,
      data,
    }: {
      id: string;
      data: RemoveCardsFromGroupInput;
    }) => creditCardGroupsAPI.removeCards(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: creditCardGroupKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: creditCardGroupKeys.detail(variables.id),
      });
      // Also invalidate credit cards list since cards are now independent
      queryClient.invalidateQueries({ queryKey: ["creditCards"] });
      toast.success("卡片已從群組移除");
    },
    onError: (error: Error) => {
      toast.error(`移除卡片失敗: ${error.message}`);
    },
  });
}

