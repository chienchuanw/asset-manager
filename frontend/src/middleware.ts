import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

/**
 * 不需要驗證的路徑
 */
const PUBLIC_PATHS = ["/login"];

/**
 * Middleware 函式
 * 保護所有路由，未登入時重導向到登入頁面
 */
export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // 取得 token cookie
  const token = request.cookies.get("token");

  // 檢查是否為公開路徑
  const isPublicPath = PUBLIC_PATHS.some((path) => pathname.startsWith(path));

  // 如果是公開路徑
  if (isPublicPath) {
    // 如果已登入，重導向到 dashboard
    if (token) {
      return NextResponse.redirect(new URL("/dashboard", request.url));
    }
    // 未登入，允許訪問
    return NextResponse.next();
  }

  // 如果不是公開路徑，檢查是否已登入
  if (!token) {
    // 未登入，重導向到登入頁面
    return NextResponse.redirect(new URL("/login", request.url));
  }

  // 已登入，允許訪問
  return NextResponse.next();
}

/**
 * Middleware 配置
 * 指定哪些路徑需要執行 middleware
 */
export const config = {
  matcher: [
    /*
     * 明確指定需要保護的路徑
     */
    "/",
    "/dashboard/:path*",
    "/holdings/:path*",
    "/transactions/:path*",
    "/analytics/:path*",
    "/settings/:path*",
    "/login",
  ],
};
