Feature: Holdings reconciliation
  As a user with investment holdings
  I want to reconcile holding quantities and average costs to the real-world state
  So that subsequent FIFO P&L is computed from a corrected basis

  Background:
    Given I am logged in
    And I am on the holdings page

  Scenario: Single-symbol reconciliation flow
    When I open the single-symbol reconcile dialog for the first holding
    And I enter a target quantity and target average cost
    And I click the preview button
    Then I see the preview diff with prev and target values
    When I click the confirm button
    Then I see a success toast
    And the dialog closes

  Scenario: Batch reconciliation with mixed actions
    When I open the batch reconcile dialog
    And I add a new holding row with symbol details
    And I edit an existing row to a different quantity
    And I edit another existing row to zero quantity to liquidate
    And I click the batch preview button
    Then I see preview rows for create, update, and liquidate actions
    When I click the batch confirm button
    Then I see a success toast
    And the batch dialog closes
