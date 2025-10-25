import { apiClient } from "./client";
import type {
  SettingsGroup,
  UpdateSettingsGroupInput,
} from "@/types/analytics";

/**
 * 取得所有設定
 */
export async function getSettings(): Promise<SettingsGroup> {
  return apiClient.get<SettingsGroup>("/api/settings");
}

/**
 * 更新設定
 */
export async function updateSettings(
  input: UpdateSettingsGroupInput
): Promise<SettingsGroup> {
  return apiClient.put<SettingsGroup>("/api/settings", input);
}
