'use client';

import { useEffect, useState } from 'react';

export default function Home() {
  const [status, setStatus] = useState<string>('載入中...');
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchHealth = async () => {
      try {
        const response = await fetch('http://localhost:8080/health');
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        setStatus(data.message);
      } catch (err) {
        setError(err instanceof Error ? err.message : '未知錯誤');
      }
    };

    fetchHealth();
  }, []);

  return (
    <div className="flex min-h-screen items-center justify-center bg-linear-to-br from-blue-50 to-indigo-100">
      <main className="flex flex-col items-center justify-center gap-8 rounded-lg bg-white p-8 shadow-lg">
        <h1 className="text-4xl font-bold text-gray-800">
          Asset Manager
        </h1>
        
        <div className="text-center">
          <p className="mb-4 text-lg text-gray-600">後端狀態：</p>
          {error ? (
            <p className="text-red-600 font-semibold">❌ {error}</p>
          ) : (
            <p className="text-green-600 font-semibold">✅ {status}</p>
          )}
        </div>

        <p className="text-sm text-gray-500">
          前端運行在 http://localhost:3000
        </p>
        <p className="text-sm text-gray-500">
          後端運行在 http://localhost:8080
        </p>
      </main>
    </div>
  );
}
