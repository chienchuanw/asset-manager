// 健康檢查端點
// 用於 Docker 容器健康檢查和負載平衡器

export async function GET() {
  return Response.json(
    {
      status: 'healthy',
      timestamp: new Date().toISOString(),
      service: 'asset-manager-frontend',
    },
    { status: 200 }
  );
}

