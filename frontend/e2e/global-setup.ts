import { chromium, request, type FullConfig } from "@playwright/test";
import { mkdir, writeFile } from "node:fs/promises";
import { dirname } from "node:path";

const STORAGE_STATE = "e2e/.auth/storage-state.json";

export default async function globalSetup(_config: FullConfig) {
  const apiURL = process.env.E2E_API_URL ?? "http://localhost:8080";
  const baseURL = process.env.E2E_BASE_URL ?? "http://localhost:3001";
  const username = process.env.E2E_USERNAME ?? "admin";
  const password = process.env.E2E_PASSWORD ?? "admin123";

  await mkdir(dirname(STORAGE_STATE), { recursive: true });

  const ctx = await request.newContext({ baseURL: apiURL });
  const res = await ctx.post("/api/auth/login", {
    data: { username, password },
  });
  if (!res.ok()) {
    const body = await res.text();
    throw new Error(
      `E2E login failed (${res.status()} from ${apiURL}/api/auth/login): ${body}\n` +
        `Set E2E_USERNAME / E2E_PASSWORD or ensure backend is running with AUTH_USERNAME/AUTH_PASSWORD configured.`,
    );
  }

  // Rewrite cookie domain to the frontend host so storageState applies in-browser.
  // Assumes API and frontend share a hostname (e.g. localhost on different ports).
  // For cross-host setups (staging), drive the real /login form instead.
  const cookies = (await ctx.storageState()).cookies.map((c) => ({
    ...c,
    domain: new URL(baseURL).hostname,
  }));
  await ctx.dispose();

  const browser = await chromium.launch();
  const browserCtx = await browser.newContext();
  await browserCtx.addCookies(cookies);
  const state = await browserCtx.storageState();
  await writeFile(STORAGE_STATE, JSON.stringify(state, null, 2));
  await browser.close();
}
