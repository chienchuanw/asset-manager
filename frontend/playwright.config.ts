// Playwright config for E2E BDD tests.
//
// Env vars:
//   E2E_BASE_URL       (default http://localhost:3001) — frontend URL.
//   E2E_API_URL        (default http://localhost:8080) — backend URL for global-setup login.
//   E2E_USERNAME       (default admin)                 — login username.
//   E2E_PASSWORD       (default admin123)              — login password.
//   E2E_SKIP_WEBSERVER (unset = start `pnpm dev`)      — set in CI when the
//                                                         frontend is started
//                                                         externally so the test
//                                                         run reuses it.
//
// Mutating scenarios are intentionally avoided; tests stop at preview to keep
// the dev DB clean. See e2e/features/holdings_reconcile.feature.

import { defineConfig, devices } from "@playwright/test";
import { defineBddConfig } from "playwright-bdd";

const testDir = defineBddConfig({
  features: "e2e/features/**/*.feature",
  steps: "e2e/steps/**/*.ts",
  outputDir: ".features-gen",
});

const baseURL = process.env.E2E_BASE_URL ?? "http://localhost:3001";

const webServer = process.env.E2E_SKIP_WEBSERVER
  ? undefined
  : {
      command: "pnpm dev",
      url: baseURL,
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
    };

export default defineConfig({
  testDir,
  fullyParallel: false,
  retries: process.env.CI ? 1 : 0,
  workers: 1,
  reporter: process.env.CI ? "list" : [["list"], ["html", { open: "never" }]],
  globalSetup: "./e2e/global-setup.ts",
  use: {
    baseURL,
    storageState: "e2e/.auth/storage-state.json",
    trace: "retain-on-failure",
  },
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
  ],
  webServer,
});
