Feature: 持倉快速對帳
  使用者可以將投資持倉一次同步至目前現實狀態，不必逐筆補登交易。
  此操作不結算實現損益，也不影響現金流。

  Background:
    Given 一個乾淨的測試資料庫

  Scenario: 更新既有持倉的數量與平均成本
    Given 目前持有 "AAPL" 100 股，平均成本 150 USD
    When 將 "AAPL" 對帳為 150 股，平均成本 180 USD
    Then 系統會新增一筆 adjustment 交易
    And 該 adjustment 的稽核欄位記錄 prev_quantity 為 100，prev_avg_cost 為 150
    And 沒有任何 realized profit 紀錄
    And 沒有任何 cash flow 紀錄

  Scenario: 對帳清倉
    Given 目前持有 "AAPL" 100 股，平均成本 150 USD
    When 將 "AAPL" 對帳為 0 股
    Then "AAPL" 的持倉數量為 0
    And 沒有任何 realized profit 紀錄

  Scenario: 透過對帳首次建倉
    Given 我目前未持有 "TSLA"
    When 將新標的 "TSLA"（名稱 "Tesla"）對帳為 50 股，平均成本 250 USD
    Then "TSLA" 的持倉存在，數量為 50，平均成本為 250

  Scenario: 批次驗證失敗時整批 rollback
    Given 目前持有 "AAPL" 100 股，平均成本 150 USD
    When 提交一個批次，包含一筆合法更新與一筆 quantity 為負的非法項目
    Then 整批請求被拒絕
    And 沒有任何 adjustment 交易被寫入

  Scenario: 對帳後的賣出以對帳後均價計算成本
    Given 目前持有 "AAPL" 100 股，平均成本 150 USD
    When 將 "AAPL" 對帳為 100 股，平均成本 200 USD
    And 賣出 "AAPL" 50 股，每股 220 USD
    Then 該賣出的成本基礎為 10000

  Scenario: 對帳不影響現金流
    Given 目前持有 "AAPL" 100 股，平均成本 150 USD
    When 將 "AAPL" 對帳為 200 股，平均成本 180 USD
    Then 沒有任何 cash flow 紀錄
