/**
 * 信用卡群組表單元件
 * 用於建立和編輯信用卡群組
 */

"use client";

import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import type {
  CreditCard,
  CreditCardGroup,
  CreateCreditCardGroupInput,
} from "@/types/user-management";

interface CreditCardGroupFormProps {
  group?: CreditCardGroup;
  availableCards: CreditCard[];
  onSubmit: (data: CreateCreditCardGroupInput) => void;
  onCancel?: () => void;
  isSubmitting?: boolean;
}

export function CreditCardGroupForm({
  group,
  availableCards,
  onSubmit,
  onCancel,
  isSubmitting = false,
}: CreditCardGroupFormProps) {
  const [selectedCardIds, setSelectedCardIds] = useState<string[]>([]);

  const form = useForm<CreateCreditCardGroupInput>({
    defaultValues: {
      name: "",
      issuing_bank: "",
      shared_credit_limit: 0,
      card_ids: [],
      note: "",
    },
  });

  // 如果是編輯模式，填入現有資料
  useEffect(() => {
    if (group) {
      form.reset({
        name: group.name,
        issuing_bank: group.issuing_bank,
        shared_credit_limit: group.shared_credit_limit,
        card_ids: [],
        note: group.note || "",
      });
    }
  }, [group, form]);

  // 處理卡片選擇
  const handleCardToggle = (cardId: string) => {
    setSelectedCardIds((prev) => {
      if (prev.includes(cardId)) {
        return prev.filter((id) => id !== cardId);
      } else {
        return [...prev, cardId];
      }
    });
  };

  // 根據選擇的卡片自動填入銀行和額度
  useEffect(() => {
    if (selectedCardIds.length > 0) {
      const firstCard = availableCards.find(
        (card) => card.id === selectedCardIds[0]
      );
      if (firstCard) {
        form.setValue("issuing_bank", firstCard.issuing_bank);
        form.setValue("shared_credit_limit", firstCard.credit_limit);
      }
    }
  }, [selectedCardIds, availableCards, form]);

  // 過濾可選擇的卡片（同銀行、同額度）
  const getFilteredCards = () => {
    if (selectedCardIds.length === 0) {
      return availableCards;
    }

    const firstCard = availableCards.find(
      (card) => card.id === selectedCardIds[0]
    );
    if (!firstCard) return availableCards;

    return availableCards.filter(
      (card) =>
        card.issuing_bank === firstCard.issuing_bank &&
        card.credit_limit === firstCard.credit_limit
    );
  };

  const handleSubmit = (data: CreateCreditCardGroupInput) => {
    onSubmit({
      ...data,
      card_ids: selectedCardIds,
    });
  };

  const filteredCards = getFilteredCards();

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-6">
        {/* 群組名稱 */}
        <FormField
          control={form.control}
          name="name"
          rules={{ required: "請輸入群組名稱" }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>群組名稱</FormLabel>
              <FormControl>
                <Input placeholder="例如：玉山銀行共同額度" {...field} />
              </FormControl>
              <FormDescription>為這個信用卡群組命名</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 選擇卡片 */}
        <div className="space-y-3">
          <Label>選擇卡片</Label>
          <FormDescription>
            選擇要加入群組的信用卡（必須為同一銀行且額度相同）
          </FormDescription>
          <div className="space-y-2 max-h-[300px] overflow-y-auto border rounded-md p-4">
            {filteredCards.length === 0 ? (
              <p className="text-sm text-muted-foreground text-center py-4">
                沒有可用的信用卡
              </p>
            ) : (
              filteredCards.map((card) => {
                const isSelected = selectedCardIds.includes(card.id);
                const isDisabled =
                  selectedCardIds.length > 0 &&
                  !isSelected &&
                  (card.issuing_bank !==
                    availableCards.find((c) => c.id === selectedCardIds[0])
                      ?.issuing_bank ||
                    card.credit_limit !==
                      availableCards.find((c) => c.id === selectedCardIds[0])
                        ?.credit_limit);

                return (
                  <div
                    key={card.id}
                    className={`flex items-center space-x-3 p-3 rounded-md border ${
                      isSelected
                        ? "bg-primary/5 border-primary"
                        : "hover:bg-muted/50"
                    } ${isDisabled ? "opacity-50 cursor-not-allowed" : "cursor-pointer"}`}
                    onClick={() => !isDisabled && handleCardToggle(card.id)}
                  >
                    <input
                      type="checkbox"
                      checked={isSelected}
                      disabled={isDisabled}
                      onChange={() => handleCardToggle(card.id)}
                      className="h-4 w-4 rounded border-gray-300"
                    />
                    <div className="flex-1">
                      <div className="flex items-center gap-2">
                        <span className="font-medium">{card.card_name}</span>
                        <Badge variant="outline" className="text-xs">
                          {card.issuing_bank}
                        </Badge>
                      </div>
                      <div className="text-sm text-muted-foreground">
                        額度:{" "}
                        {card.credit_limit.toLocaleString("zh-TW", {
                          minimumFractionDigits: 0,
                          maximumFractionDigits: 0,
                        })}{" "}
                        | 已用:{" "}
                        {card.used_credit.toLocaleString("zh-TW", {
                          minimumFractionDigits: 0,
                          maximumFractionDigits: 0,
                        })}
                      </div>
                    </div>
                  </div>
                );
              })
            )}
          </div>
          {selectedCardIds.length === 0 && (
            <p className="text-sm text-destructive">請至少選擇一張卡片</p>
          )}
        </div>

        {/* 發卡銀行（自動填入） */}
        <FormField
          control={form.control}
          name="issuing_bank"
          render={({ field }) => (
            <FormItem>
              <FormLabel>發卡銀行</FormLabel>
              <FormControl>
                <Input {...field} disabled />
              </FormControl>
              <FormDescription>根據選擇的卡片自動填入</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 共同額度（自動填入） */}
        <FormField
          control={form.control}
          name="shared_credit_limit"
          render={({ field }) => (
            <FormItem>
              <FormLabel>共同額度</FormLabel>
              <FormControl>
                <Input
                  type="number"
                  {...field}
                  onChange={(e) => field.onChange(parseFloat(e.target.value))}
                  disabled
                />
              </FormControl>
              <FormDescription>根據選擇的卡片自動填入</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 備註 */}
        <FormField
          control={form.control}
          name="note"
          render={({ field }) => (
            <FormItem>
              <FormLabel>備註</FormLabel>
              <FormControl>
                <Textarea
                  placeholder="輸入備註（選填）"
                  className="resize-none"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* 按鈕 */}
        <div className="flex justify-end gap-3">
          {onCancel && (
            <Button
              type="button"
              variant="outline"
              onClick={onCancel}
              disabled={isSubmitting}
            >
              取消
            </Button>
          )}
          <Button
            type="submit"
            disabled={isSubmitting || selectedCardIds.length === 0}
          >
            {isSubmitting ? "處理中..." : group ? "更新群組" : "建立群組"}
          </Button>
        </div>
      </form>
    </Form>
  );
}

