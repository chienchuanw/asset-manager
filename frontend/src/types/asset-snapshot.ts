/**
 * Asset Snapshot types
 * Defines data structures for asset snapshots
 */

export type SnapshotAssetType = 'total' | 'tw-stock' | 'us-stock' | 'crypto';

export interface AssetSnapshot {
  id: string;
  snapshot_date: string; // ISO date string (YYYY-MM-DD)
  asset_type: SnapshotAssetType;
  value_twd: number;
  created_at: string;
  updated_at: string;
}

export interface AssetTrendData {
  date: string;
  total?: number;
  twStock?: number;
  usStock?: number;
  crypto?: number;
}

