"use client";

import { useTransactions, useCreateTransaction } from "@/hooks";
import type { CreateTransactionInput } from "@/types/transaction";

/**
 * 交易範例元件
 * 
 * 這是一個簡單的範例，展示如何使用 React Query Hooks
 */
export function TransactionExample() {
  // 取得交易列表
  const { data, isLoading, error } = useTransactions();

  // 建立交易 mutation
  const createMutation = useCreateTransaction({
    onSuccess: () => {
      console.log("交易建立成功");
    },
    onError: (error) => {
      console.error("交易建立失敗:", error.message);
    },
  });

  // 處理建立測試交易
  const handleCreateTestTransaction = () => {
    const testTransaction: CreateTransactionInput = {
      date: new Date().toISOString(),
      asset_type: "tw-stock",
      symbol: "2330",
      name: "台積電",
      type: "buy",
      quantity: 10,
      price: 620,
      amount: 6200,
      fee: 28,
      currency: "TWD",
      note: "測試交易",
    };

    createMutation.mutate(testTransaction);
  };

  // 載入中
  if (isLoading) {
    return (
      <div className="p-4">
        <p>載入中...</p>
      </div>
    );
  }

  // 錯誤
  if (error) {
    return (
      <div className="p-4 text-red-500">
        <p>錯誤: {error.message}</p>
        <p className="text-sm">錯誤代碼: {error.code}</p>
      </div>
    );
  }

  return (
    <div className="p-4 space-y-4">
      <div>
        <h2 className="text-xl font-bold mb-2">交易列表範例</h2>
        <button
          onClick={handleCreateTestTransaction}
          disabled={createMutation.isPending}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-400"
        >
          {createMutation.isPending ? "建立中..." : "建立測試交易"}
        </button>
      </div>

      <div>
        <h3 className="font-semibold mb-2">
          交易數量: {data?.length || 0}
        </h3>
        <div className="space-y-2">
          {data?.map((transaction) => (
            <div
              key={transaction.id}
              className="p-3 border rounded bg-gray-50"
            >
              <div className="flex justify-between">
                <span className="font-medium">{transaction.name}</span>
                <span className="text-sm text-gray-500">
                  {transaction.symbol}
                </span>
              </div>
              <div className="text-sm text-gray-600">
                {transaction.type} | 數量: {transaction.quantity} | 價格:{" "}
                {transaction.price}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

