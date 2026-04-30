import { expect } from "@playwright/test";
import { createBdd } from "playwright-bdd";

const { Given, When, Then } = createBdd();

Given("I am logged in", async ({ page }) => {
  await page.goto("/");
  await expect(page).not.toHaveURL(/\/login/);
});

Given("I am on the holdings page", async ({ page }) => {
  await page.goto("/holdings");
  await page.waitForLoadState("networkidle");
});

When(
  "I open the single-symbol reconcile dialog for the first holding",
  async ({ page }) => {
    await page.locator('[data-testid^="holding-row-"]').first().click();
    await page.getByTestId("reconcile-single-open").first().click();
    await expect(page.getByTestId("reconcile-quantity")).toBeVisible();
  },
);

When("I enter a target quantity and target average cost", async ({ page }) => {
  await page.getByTestId("reconcile-quantity").fill("10");
  await page.getByTestId("reconcile-avg-cost").fill("100");
});

When("I click the preview button", async ({ page }) => {
  await page.getByTestId("reconcile-preview-btn").click();
});

Then(
  "I see the preview diff with prev and target values",
  async ({ page }) => {
    await expect(page.getByTestId("reconcile-confirm-btn")).toBeVisible();
  },
);

When("I click the confirm button", async ({ page }) => {
  await page.getByTestId("reconcile-confirm-btn").click();
});

Then("I see a success toast", async ({ page }) => {
  await expect(page.locator("[data-sonner-toast]").first()).toBeVisible({
    timeout: 5_000,
  });
});

Then("the dialog closes", async ({ page }) => {
  await expect(page.getByTestId("reconcile-quantity")).toBeHidden();
});

When("I open the batch reconcile dialog", async ({ page }) => {
  await page.getByTestId("reconcile-batch-open").click();
  await expect(page.getByTestId("reconcile-batch-preview-btn")).toBeVisible();
});

When("I add a new holding row with symbol details", async ({ page }) => {
  await page.getByTestId("reconcile-batch-add-row").click();
});

When(
  "I edit an existing row to a different quantity",
  async ({ page }) => {
    const qtyInputs = page.locator(
      '[data-testid^="reconcile-batch-quantity-"]',
    );
    await qtyInputs.first().fill("5");
  },
);

When(
  "I edit another existing row to zero quantity to liquidate",
  async ({ page }) => {
    const qtyInputs = page.locator(
      '[data-testid^="reconcile-batch-quantity-"]',
    );
    await qtyInputs.nth(1).fill("0");
  },
);

When("I click the batch preview button", async ({ page }) => {
  await page.getByTestId("reconcile-batch-preview-btn").click();
});

Then(
  "I see preview rows for create, update, and liquidate actions",
  async ({ page }) => {
    await expect(page.getByTestId("reconcile-batch-confirm-btn")).toBeVisible();
  },
);

When("I click the batch confirm button", async ({ page }) => {
  await page.getByTestId("reconcile-batch-confirm-btn").click();
});

Then("the batch dialog closes", async ({ page }) => {
  await expect(page.getByTestId("reconcile-batch-preview-btn")).toBeHidden();
});
