/**
 * Asset Snapshots API Client
 * Handles API calls for asset snapshots
 */

import { apiClient } from "./client";
import {
  AssetSnapshot,
  AssetTrendResponse,
  SnapshotAssetType,
} from "@/types/asset-snapshot";

/**
 * Get asset trend data for a specific asset type
 * @param assetType - Type of asset (total, tw-stock, us-stock, crypto)
 * @param days - Number of days to retrieve (default: 30)
 * @returns Array of asset trend responses
 */
export async function getAssetTrend(
  assetType: SnapshotAssetType,
  days: number = 30
): Promise<AssetTrendResponse[]> {
  return await apiClient.get<AssetTrendResponse[]>("/api/snapshots/trend", {
    params: {
      asset_type: assetType,
      days,
    },
  });
}

/**
 * Get the latest snapshot for a specific asset type
 * @param assetType - Type of asset (total, tw-stock, us-stock, crypto)
 * @returns Latest asset snapshot or null if not found
 */
export async function getLatestSnapshot(
  assetType: SnapshotAssetType
): Promise<AssetSnapshot | null> {
  return await apiClient.get<AssetSnapshot | null>("/api/snapshots/latest", {
    params: {
      asset_type: assetType,
    },
  });
}

/**
 * Manually trigger daily snapshot creation
 * @returns Success message
 */
export async function triggerDailySnapshots(): Promise<{ message: string }> {
  return await apiClient.post<{ message: string }>("/api/snapshots/trigger");
}
