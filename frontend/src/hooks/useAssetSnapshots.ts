/**
 * Asset Snapshots Hooks
 * React Query hooks for asset snapshots data
 */

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  getAssetTrend,
  getLatestSnapshot,
  triggerDailySnapshots,
} from "@/lib/api/asset-snapshots";
import {
  SnapshotAssetType,
  AssetSnapshot,
  AssetTrendResponse,
} from "@/types/asset-snapshot";

/**
 * Hook to fetch asset trend data
 * @param assetType - Type of asset to fetch
 * @param days - Number of days to retrieve
 * @returns React Query result with asset trend responses array
 */
export function useAssetTrend(assetType: SnapshotAssetType, days: number = 30) {
  return useQuery<AssetTrendResponse[], Error>({
    queryKey: ["asset-trend", assetType, days],
    queryFn: () => getAssetTrend(assetType, days),
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
  });
}

/**
 * Hook to fetch the latest snapshot for an asset type
 * @param assetType - Type of asset to fetch
 * @returns React Query result with latest snapshot
 */
export function useLatestSnapshot(assetType: SnapshotAssetType) {
  return useQuery<AssetSnapshot | null, Error>({
    queryKey: ["latest-snapshot", assetType],
    queryFn: () => getLatestSnapshot(assetType),
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
  });
}

/**
 * Hook to manually trigger daily snapshot creation
 * @returns Mutation function and state
 */
export function useTriggerDailySnapshots() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: triggerDailySnapshots,
    onSuccess: () => {
      // Invalidate all snapshot queries to refetch data
      queryClient.invalidateQueries({ queryKey: ["asset-trend"] });
      queryClient.invalidateQueries({ queryKey: ["latest-snapshot"] });
    },
  });
}
