/**
 * CSV 匯入 Dialog 元件
 * 提供 CSV 檔案上傳功能，包含樣板下載與檔案驗證
 */

"use client";

import { useState, useRef } from "react";
import { useTranslations } from "next-intl";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Upload, Download, FileText, AlertCircle } from "lucide-react";
import { toast } from "sonner";

interface CSVImportDialogProps {
  onSuccess: (transactions: any[]) => void;
}

export function CSVImportDialog({ onSuccess }: CSVImportDialogProps) {
  const t = useTranslations("transactions");
  const [open, setOpen] = useState(false);
  const [file, setFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [errors, setErrors] = useState<any[]>([]);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // 處理檔案選擇
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (selectedFile) {
      // 驗證檔案類型
      if (!selectedFile.name.endsWith(".csv")) {
        toast.error(t("csvErrorFileType"));
        return;
      }
      setFile(selectedFile);
      setErrors([]);
    }
  };

  // 下載 CSV 樣板
  const handleDownloadTemplate = async () => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/transactions/template`,
        {
          credentials: "include", // 重要：帶上 httpOnly cookie
        }
      );

      if (!response.ok) {
        if (response.status === 401) {
          toast.error(t("csvErrorExpired"));
        } else {
          toast.error(t("csvErrorDownload", { status: response.status }));
        }
        return;
      }

      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = "transaction_template.csv";
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);

      toast.success(t("csvSuccessDownload"));
    } catch (error) {
      console.error("下載樣板失敗:", error);
      toast.error(t("csvErrorDownloadRetry"));
    }
  };

  // 上傳並解析 CSV
  const handleUpload = async () => {
    if (!file) {
      toast.error(t("csvErrorNoFile"));
      return;
    }

    setIsUploading(true);
    setErrors([]);

    try {
      const formData = new FormData();
      formData.append("file", file);

      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/transactions/parse-csv`,
        {
          method: "POST",
          credentials: "include", // 重要：帶上 httpOnly cookie
          body: formData,
        }
      );

      const result = await response.json();

      if (!response.ok) {
        if (response.status === 401) {
          toast.error(t("csvErrorExpired"));
          return;
        }
        throw new Error(result.error?.message || t("csvErrorUpload"));
      }

      // 檢查解析結果
      if (!result.data.success) {
        setErrors(result.data.errors);
        toast.error(t("csvErrorValidation"));
        return;
      }

      // 解析成功，傳遞交易資料給父元件
      toast.success(t("csvSuccessUpload"));
      onSuccess(result.data.transactions);
      setOpen(false);
      setFile(null);
      setErrors([]);
    } catch (error: any) {
      console.error("上傳失敗:", error);
      toast.error(error.message || t("csvErrorUpload"));
    } finally {
      setIsUploading(false);
    }
  };

  // 重置狀態
  const handleReset = () => {
    setFile(null);
    setErrors([]);
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm">
          <Upload className="h-4 w-4 mr-2" />
          {t("csvImportButton")}
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>{t("csvImportTitle")}</DialogTitle>
          <DialogDescription>{t("csvImportDesc")}</DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-4">
          {/* 下載樣板按鈕 */}
          <div className="flex items-center justify-between p-4 border rounded-lg bg-muted/50">
            <div className="flex items-center gap-3">
              <FileText className="h-5 w-5 text-muted-foreground" />
              <div>
                <p className="text-sm font-medium">CSV 樣板檔案</p>
                <p className="text-xs text-muted-foreground">
                  下載樣板並填寫交易資料
                </p>
              </div>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={handleDownloadTemplate}
            >
              <Download className="h-4 w-4 mr-2" />
              {t("downloadTemplate")}
            </Button>
          </div>

          {/* 檔案上傳區域 */}
          <div className="space-y-2">
            <label className="text-sm font-medium">{t("selectFile")}</label>
            <div className="flex items-center gap-2">
              <input
                ref={fileInputRef}
                type="file"
                accept=".csv"
                onChange={handleFileChange}
                className="hidden"
                id="csv-file-input"
              />
              <label
                htmlFor="csv-file-input"
                className="flex-1 flex items-center justify-center h-32 border-2 border-dashed rounded-lg cursor-pointer hover:bg-muted/50 transition-colors"
              >
                <div className="text-center">
                  <Upload className="h-8 w-8 mx-auto mb-2 text-muted-foreground" />
                  <p className="text-sm font-medium">
                    {file ? file.name : t("selectFile")}
                  </p>
                  <p className="text-xs text-muted-foreground mt-1">
                    支援 .csv 格式
                  </p>
                </div>
              </label>
            </div>
            {file && (
              <Button
                variant="ghost"
                size="sm"
                onClick={handleReset}
                className="w-full"
              >
                {t("selectFile")}
              </Button>
            )}
          </div>

          {/* 錯誤訊息 */}
          {errors.length > 0 && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                <p className="font-medium mb-2">
                  發現 {errors.length} 個錯誤：
                </p>
                <ul className="list-disc list-inside space-y-1 text-sm">
                  {errors.slice(0, 5).map((error, index) => (
                    <li key={index}>
                      第 {error.row} 行 - {error.field}: {error.message}
                    </li>
                  ))}
                  {errors.length > 5 && (
                    <li className="text-muted-foreground">
                      還有 {errors.length - 5} 個錯誤...
                    </li>
                  )}
                </ul>
              </AlertDescription>
            </Alert>
          )}
        </div>

        {/* 操作按鈕 */}
        <div className="flex justify-end gap-2">
          <Button
            variant="outline"
            onClick={() => {
              setOpen(false);
              handleReset();
            }}
            disabled={isUploading}
          >
            {t("selectFile")}
          </Button>
          <Button onClick={handleUpload} disabled={!file || isUploading}>
            {isUploading ? t("selectFile") : t("uploadFile")}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
