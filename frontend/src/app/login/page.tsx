import { LoginForm } from "@/components/login-form";
import { LanguageSwitcher } from "@/components/common/LanguageSwitcher";

export default function Page() {
  return (
    <div className="relative flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      {/* 語言切換按鈕 - 右上角 */}
      <div className="absolute right-4 top-4">
        <LanguageSwitcher />
      </div>
      <div className="w-full max-w-sm">
        <LoginForm />
      </div>
    </div>
  );
}
