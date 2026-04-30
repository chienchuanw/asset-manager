Feature: Holdings reconciliation
  As a user with investment holdings
  I want to reconcile holding quantities and average costs to the real-world state
  So that subsequent FIFO P&L is computed from a corrected basis

  Note: scenarios stop at the preview step to avoid mutating the dev database.
  Mutating end-to-end coverage (confirm + toast + table refetch) is intentionally
  left to manual verification. See frontend/e2e/README guidance in PR for details.

  Background:
    Given I am logged in
    And I am on the holdings page

  Scenario: Batch reconcile preview shows mixed actions
    When I open the batch reconcile dialog
    And I add a new holding row with symbol "BDDX", name "BDD Test Symbol", currency "USD", quantity "1", avg cost "10"
    And I edit the first existing row quantity to "5"
    And I edit the second existing row quantity to "0" to liquidate
    And I click the batch preview button
    Then I see a preview row for the "create" action
    And I see a preview row for the "update" action
    And I see a preview row for the "liquidate" action
