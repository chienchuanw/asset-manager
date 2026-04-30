import { expect } from "@playwright/test";
import { createBdd } from "playwright-bdd";

const { Given, When, Then } = createBdd();

Given("I am logged in", async ({ page }) => {
  await page.goto("/");
  await expect(page).not.toHaveURL(/\/login/);
});

Given("I am on the holdings page", async ({ page }) => {
  await page.goto("/holdings");
  await expect(page.getByTestId("reconcile-batch-open")).toBeVisible();
});

When("I open the batch reconcile dialog", async ({ page }) => {
  await page.getByTestId("reconcile-batch-open").click();
  await expect(page.getByTestId("reconcile-batch-preview-btn")).toBeVisible();
});

When(
  "I add a new holding row with symbol {string}, name {string}, currency {string}, quantity {string}, avg cost {string}",
  async ({ page }, symbol: string, name: string, currency: string, qty: string, avgCost: string) => {
    const existingCount = await page
      .locator('[data-testid^="reconcile-batch-row-"]')
      .count();
    await page.getByTestId("reconcile-batch-add-row").click();
    const i = existingCount;
    await page.getByTestId(`reconcile-batch-symbol-${i}`).fill(symbol);
    await page.getByTestId(`reconcile-batch-name-${i}`).fill(name);
    await page.getByTestId(`reconcile-batch-currency-${i}`).click();
    await page.getByRole("option", { name: currency }).click();
    await page.getByTestId(`reconcile-batch-quantity-${i}`).fill(qty);
    await page.getByTestId(`reconcile-batch-avg-cost-${i}`).fill(avgCost);
  },
);

When(
  "I edit the first existing row quantity to {string}",
  async ({ page }, qty: string) => {
    await page.getByTestId("reconcile-batch-quantity-0").fill(qty);
  },
);

When(
  "I edit the second existing row quantity to {string} to liquidate",
  async ({ page }, qty: string) => {
    await page.getByTestId("reconcile-batch-quantity-1").fill(qty);
  },
);

When("I click the batch preview button", async ({ page }) => {
  await page.getByTestId("reconcile-batch-preview-btn").click();
});

Then(
  "I see a preview row for the {string} action",
  async ({ page }, action: string) => {
    await expect(
      page.getByTestId(`reconcile-preview-row-${action}`).first(),
    ).toBeVisible();
  },
);
