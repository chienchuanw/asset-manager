/**
 * 信用卡校準 Dialog
 * 可同時校準 used_credit 與 credit_limit；只有 used_credit 變動會產生對帳 cash flow。
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
import { useReconcileCreditCard } from "@/hooks/useReconcile";
import type { CreditCard } from "@/types/user-management";

interface FormValues {
  new_used_credit: string; // 空字串代表不更動
  new_credit_limit: string;
  date: string;
  note: string;
}

interface Props {
  open: boolean;
  card: CreditCard;
  onClose: () => void;
}

export function ReconcileCreditCardDialog({ open, card, onClose }: Props) {
  const t = useTranslations("userManagement.reconcile");
  const tCommon = useTranslations("common");
  const { toast } = useToast();
  const mutation = useReconcileCreditCard();

  const form = useForm<FormValues>({
    defaultValues: {
      new_used_credit: "",
      new_credit_limit: "",
      date: new Date().toISOString().slice(0, 10),
      note: "",
    },
  });

  const usedRaw = form.watch("new_used_credit");
  const limitRaw = form.watch("new_credit_limit");

  const usedDelta = useMemo(() => {
    if (usedRaw === "" || usedRaw === undefined) return null;
    return Number(usedRaw) - card.used_credit;
  }, [usedRaw, card.used_credit]);

  const limitChanged = useMemo(() => {
    if (limitRaw === "" || limitRaw === undefined) return false;
    return Number(limitRaw) !== card.credit_limit;
  }, [limitRaw, card.credit_limit]);

  const onSubmit = (values: FormValues) => {
    if (
      values.new_used_credit === "" &&
      values.new_credit_limit === ""
    ) {
      form.setError("new_used_credit", { message: t("atLeastOneFieldRequired") });
      return;
    }

    mutation.mutate(
      {
        id: card.id,
        input: {
          new_used_credit:
            values.new_used_credit !== ""
              ? Number(values.new_used_credit)
              : undefined,
          new_credit_limit:
            values.new_credit_limit !== ""
              ? Number(values.new_credit_limit)
              : undefined,
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
          <DialogTitle>{t("creditCardTitle")}</DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="space-y-4"
            aria-label="reconcile-credit-card-form"
          >
            <div className="grid grid-cols-2 gap-4">
              <div>
                <div className="text-sm text-muted-foreground">
                  {t("currentUsedCredit")}
                </div>
                <div className="text-lg font-medium">
                  {card.used_credit.toLocaleString()}
                </div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground">
                  {t("currentCreditLimit")}
                </div>
                <div className="text-lg font-medium">
                  {card.credit_limit.toLocaleString()}
                </div>
              </div>
            </div>

            <FormField
              control={form.control}
              name="new_used_credit"
              rules={{
                validate: (v) =>
                  v === "" ||
                  Number(v) >= 0 ||
                  t("newUsedCreditMustBeNonNegative"),
              }}
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("newUsedCredit")}</FormLabel>
                  <FormControl>
                    <Input type="number" step="0.01" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="new_credit_limit"
              rules={{
                validate: (v) =>
                  v === "" ||
                  Number(v) > 0 ||
                  t("newCreditLimitMustBePositive"),
              }}
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("newCreditLimit")}</FormLabel>
                  <FormControl>
                    <Input type="number" step="0.01" {...field} />
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

            <div className="rounded border p-3 bg-muted/30 text-sm space-y-1">
              {usedDelta === null ? (
                limitChanged ? (
                  <span>{t("creditLimitOverwriteNote")}</span>
                ) : (
                  <span>{t("noChange")}</span>
                )
              ) : usedDelta === 0 ? (
                limitChanged ? (
                  <span>{t("creditLimitOverwriteNote")}</span>
                ) : (
                  <span>{t("noChange")}</span>
                )
              ) : (
                <>
                  <div>
                    {t(usedDelta > 0 ? "deltaPositive" : "deltaNegative", {
                      amount: Math.abs(usedDelta).toLocaleString(),
                      category:
                        usedDelta > 0
                          ? t("categoryExpense")
                          : t("categoryIncome"),
                    })}
                  </div>
                  {limitChanged && <div>{t("creditLimitOverwriteNote")}</div>}
                </>
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
