/**
 * 銀行帳戶餘額校準 Dialog
 * 顯示目前餘額、收新值與日期，並即時預覽差額。
 */

"use client";

import { useMemo } from "react";
import { useForm } from "react-hook-form";
import { useTranslations } from "next-intl";

import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/hooks/use-toast";
import { useReconcileBankAccount } from "@/hooks/useReconcile";
import type { BankAccount } from "@/types/user-management";

interface FormValues {
  new_balance: number;
  date: string;
  note: string;
}

interface Props {
  open: boolean;
  account: BankAccount;
  onClose: () => void;
}

export function ReconcileBankAccountDialog({ open, account, onClose }: Props) {
  const t = useTranslations("userManagement.reconcile");
  const tCommon = useTranslations("common");
  const { toast } = useToast();
  const mutation = useReconcileBankAccount();

  const form = useForm<FormValues>({
    defaultValues: {
      new_balance: account.balance,
      date: new Date().toISOString().slice(0, 10),
      note: "",
    },
  });

  const newBalance = form.watch("new_balance");
  const delta = useMemo(
    () => Number(newBalance ?? 0) - account.balance,
    [newBalance, account.balance]
  );

  const onSubmit = (values: FormValues) => {
    mutation.mutate(
      {
        id: account.id,
        input: {
          new_balance: Number(values.new_balance),
          date: new Date(values.date).toISOString(),
          note: values.note || undefined,
        },
      },
      {
        onSuccess: (res) => {
          toast({
            title: t("success"),
            description: t("createdAdjustment", {
              delta: res.delta.toLocaleString(),
            }),
          });
          onClose();
        },
        onError: (err) => {
          toast({
            title: t("error"),
            description: err?.message ?? t("error"),
            variant: "destructive",
          });
        },
      }
    );
  };

  return (
    <Dialog open={open} onOpenChange={(o) => !o && onClose()}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("bankAccountTitle")}</DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="space-y-4"
            aria-label="reconcile-bank-account-form"
          >
            <div>
              <div className="text-sm text-muted-foreground">
                {t("currentBalance")}
              </div>
              <div className="text-lg font-medium">
                {account.balance.toLocaleString()} {account.currency}
              </div>
            </div>

            <FormField
              control={form.control}
              name="new_balance"
              rules={{
                required: t("newBalanceRequired"),
                validate: (v) =>
                  Number(v) >= 0 || t("newBalanceMustBeNonNegative"),
              }}
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("newBalance")}</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      step="0.01"
                      {...field}
                      onChange={(e) => field.onChange(Number(e.target.value))}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="date"
              rules={{ required: t("dateRequired") }}
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("date")}</FormLabel>
                  <FormControl>
                    <Input type="date" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="note"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("note")}</FormLabel>
                  <FormControl>
                    <Textarea {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="rounded border p-3 bg-muted/30 text-sm">
              {delta === 0 ? (
                <span>{t("noChange")}</span>
              ) : (
                <span>
                  {t(delta > 0 ? "deltaPositive" : "deltaNegative", {
                    amount: Math.abs(delta).toLocaleString(),
                    category:
                      delta > 0 ? t("categoryIncome") : t("categoryExpense"),
                  })}
                </span>
              )}
            </div>

            <DialogFooter>
              <Button type="button" variant="outline" onClick={onClose}>
                {tCommon("cancel")}
              </Button>
              <Button type="submit" disabled={mutation.isPending}>
                {t("confirm")}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
