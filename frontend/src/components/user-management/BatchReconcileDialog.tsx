/**
 * 批次校準 Dialog
 * 表格列出所有銀行帳戶 + 信用卡，使用者可以一次更新多個項目；空白輸入代表跳過。
 */

"use client";

import { useMemo, useState } from "react";
import { useTranslations } from "next-intl";

import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useToast } from "@/hooks/use-toast";
import { useReconcileBatch } from "@/hooks/useReconcile";
import type {
  BankAccount,
  CreditCard,
} from "@/types/user-management";
import type { ReconcileBatchItem } from "@/types/reconciliation";

interface RowState {
  newBalance?: string;
  newUsedCredit?: string;
  newCreditLimit?: string;
}

interface Props {
  open: boolean;
  bankAccounts: BankAccount[];
  creditCards: CreditCard[];
  onClose: () => void;
}

export function BatchReconcileDialog({
  open,
  bankAccounts,
  creditCards,
  onClose,
}: Props) {
  const t = useTranslations("userManagement.reconcile");
  const tCommon = useTranslations("common");
  const { toast } = useToast();
  const mutation = useReconcileBatch();

  const [date, setDate] = useState(new Date().toISOString().slice(0, 10));
  const [note, setNote] = useState("");
  const [rows, setRows] = useState<Record<string, RowState>>({});

  const updateRow = (id: string, patch: Partial<RowState>) =>
    setRows((prev) => ({ ...prev, [id]: { ...prev[id], ...patch } }));

  const { items, expectedCashFlows } = useMemo(() => {
    const items: ReconcileBatchItem[] = [];
    let cashFlows = 0;

    for (const a of bankAccounts) {
      const v = rows[a.id]?.newBalance;
      if (v === undefined || v === "") continue;
      const next = Number(v);
      items.push({
        target_type: "bank_account",
        target_id: a.id,
        new_balance: next,
      });
      if (next - a.balance !== 0) cashFlows++;
    }

    for (const c of creditCards) {
      const r = rows[c.id] ?? {};
      const item: ReconcileBatchItem = {
        target_type: "credit_card",
        target_id: c.id,
      };
      let any = false;
      let usedDelta = 0;
      if (r.newUsedCredit !== undefined && r.newUsedCredit !== "") {
        const next = Number(r.newUsedCredit);
        item.new_used_credit = next;
        usedDelta = next - c.used_credit;
        any = true;
      }
      if (r.newCreditLimit !== undefined && r.newCreditLimit !== "") {
        item.new_credit_limit = Number(r.newCreditLimit);
        any = true;
      }
      if (any) {
        items.push(item);
        if (usedDelta !== 0) cashFlows++;
      }
    }
    return { items, expectedCashFlows: cashFlows };
  }, [bankAccounts, creditCards, rows]);

  const submit = () => {
    if (items.length === 0) return;
    mutation.mutate(
      {
        date: new Date(date).toISOString(),
        note: note || undefined,
        items,
      },
      {
        onSuccess: (results) => {
          toast({
            title: t("batchSuccess"),
            description: t("batchSummary", { items: results.length }),
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
      <DialogContent className="max-w-3xl">
        <DialogHeader>
          <DialogTitle>{t("batchTitle")}</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <Label htmlFor="batch-date">{t("date")}</Label>
              <Input
                id="batch-date"
                type="date"
                value={date}
                onChange={(e) => setDate(e.target.value)}
              />
            </div>
          </div>

          <div>
            <Label htmlFor="batch-note">{t("note")}</Label>
            <Textarea
              id="batch-note"
              value={note}
              onChange={(e) => setNote(e.target.value)}
            />
          </div>

          {bankAccounts.length > 0 && (
            <div>
              <h3 className="font-medium mb-2">{t("bankAccountsSection")}</h3>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>{t("nameColumn")}</TableHead>
                    <TableHead className="text-right">
                      {t("currentBalance")}
                    </TableHead>
                    <TableHead className="text-right">
                      {t("newBalance")}
                    </TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {bankAccounts.map((a) => (
                    <TableRow key={a.id}>
                      <TableCell>
                        {a.bank_name} ****{a.account_number_last4}
                      </TableCell>
                      <TableCell className="text-right">
                        {a.balance.toLocaleString()}
                      </TableCell>
                      <TableCell className="text-right">
                        <Input
                          type="number"
                          step="0.01"
                          value={rows[a.id]?.newBalance ?? ""}
                          onChange={(e) =>
                            updateRow(a.id, { newBalance: e.target.value })
                          }
                          className="w-32 ml-auto"
                        />
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          )}

          {creditCards.length > 0 && (
            <div>
              <h3 className="font-medium mb-2">{t("creditCardsSection")}</h3>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>{t("nameColumn")}</TableHead>
                    <TableHead className="text-right">
                      {t("currentUsedCredit")}
                    </TableHead>
                    <TableHead className="text-right">
                      {t("newUsedCredit")}
                    </TableHead>
                    <TableHead className="text-right">
                      {t("currentCreditLimit")}
                    </TableHead>
                    <TableHead className="text-right">
                      {t("newCreditLimit")}
                    </TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {creditCards.map((c) => (
                    <TableRow key={c.id}>
                      <TableCell>
                        {c.issuing_bank} {c.card_name} ****{c.card_number_last4}
                      </TableCell>
                      <TableCell className="text-right">
                        {c.used_credit.toLocaleString()}
                      </TableCell>
                      <TableCell className="text-right">
                        <Input
                          type="number"
                          step="0.01"
                          value={rows[c.id]?.newUsedCredit ?? ""}
                          onChange={(e) =>
                            updateRow(c.id, { newUsedCredit: e.target.value })
                          }
                          className="w-32 ml-auto"
                        />
                      </TableCell>
                      <TableCell className="text-right">
                        {c.credit_limit.toLocaleString()}
                      </TableCell>
                      <TableCell className="text-right">
                        <Input
                          type="number"
                          step="0.01"
                          value={rows[c.id]?.newCreditLimit ?? ""}
                          onChange={(e) =>
                            updateRow(c.id, { newCreditLimit: e.target.value })
                          }
                          className="w-32 ml-auto"
                        />
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          )}

          <div className="rounded border p-3 bg-muted/30 text-sm">
            {t("batchPreview", {
              items: items.length,
              cashFlows: expectedCashFlows,
            })}
          </div>
        </div>

        <DialogFooter>
          <Button type="button" variant="outline" onClick={onClose}>
            {tCommon("cancel")}
          </Button>
          <Button
            type="button"
            onClick={submit}
            disabled={mutation.isPending || items.length === 0}
          >
            {t("confirm")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
