/**
 * 使用者管理頁面
 * 整合銀行帳戶和信用卡的管理功能
 */

"use client";

import { useState } from "react";
import { AppLayout } from "@/components/layout/AppLayout";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { BankAccountList } from "@/components/user-management/BankAccountList";
import { BankAccountForm } from "@/components/user-management/BankAccountForm";
import { CreditCardList } from "@/components/user-management/CreditCardList";
import { CreditCardForm } from "@/components/user-management/CreditCardForm";
import { CreditCardGroupForm } from "@/components/user-management/CreditCardGroupForm";
import { CategoryManagement } from "@/components/cash-flows";
import {
  useBankAccounts,
  useCreateBankAccount,
  useUpdateBankAccount,
  useDeleteBankAccount,
} from "@/hooks/useBankAccounts";
import {
  useCreditCards,
  useCreateCreditCard,
  useUpdateCreditCard,
  useDeleteCreditCard,
} from "@/hooks/useCreditCards";
import {
  useCreditCardGroups,
  useCreateCreditCardGroup,
  useUpdateCreditCardGroup,
  useDeleteCreditCardGroup,
  useRemoveCardsFromGroup,
} from "@/hooks/useCreditCardGroups";
import { PlusIcon, FolderPlusIcon } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import type {
  BankAccount,
  CreateBankAccountInput,
  CreditCard,
  CreateCreditCardInput,
  CreditCardGroupWithCards,
  CreateCreditCardGroupInput,
} from "@/types/user-management";

export default function UserManagementPage() {
  const { toast } = useToast();

  // 銀行帳戶相關狀態
  const [bankAccountDialogOpen, setBankAccountDialogOpen] = useState(false);
  const [editingBankAccount, setEditingBankAccount] = useState<
    BankAccount | undefined
  >();
  const [deletingBankAccountId, setDeletingBankAccountId] = useState<
    string | undefined
  >();

  // 信用卡相關狀態
  const [creditCardDialogOpen, setCreditCardDialogOpen] = useState(false);
  const [editingCreditCard, setEditingCreditCard] = useState<
    CreditCard | undefined
  >();
  const [deletingCreditCardId, setDeletingCreditCardId] = useState<
    string | undefined
  >();

  // 信用卡群組相關狀態
  const [groupDialogOpen, setGroupDialogOpen] = useState(false);
  const [editingGroup, setEditingGroup] = useState<
    CreditCardGroupWithCards | undefined
  >();
  const [deletingGroupId, setDeletingGroupId] = useState<string | undefined>();

  // 資料查詢
  const { data: bankAccounts, isLoading: bankAccountsLoading } =
    useBankAccounts();
  const { data: creditCards, isLoading: creditCardsLoading } = useCreditCards();
  const { data: groups, isLoading: groupsLoading } = useCreditCardGroups();

  // Mutations
  const createBankAccountMutation = useCreateBankAccount();
  const updateBankAccountMutation = useUpdateBankAccount();
  const deleteBankAccountMutation = useDeleteBankAccount();
  const createCreditCardMutation = useCreateCreditCard();
  const updateCreditCardMutation = useUpdateCreditCard();
  const deleteCreditCardMutation = useDeleteCreditCard();
  const createGroupMutation = useCreateCreditCardGroup();
  const updateGroupMutation = useUpdateCreditCardGroup();
  const deleteGroupMutation = useDeleteCreditCardGroup();
  const removeCardsFromGroupMutation = useRemoveCardsFromGroup();

  // ==================== 銀行帳戶處理函式 ====================

  const handleCreateBankAccount = () => {
    setEditingBankAccount(undefined);
    setBankAccountDialogOpen(true);
  };

  const handleEditBankAccount = (bankAccount: BankAccount) => {
    setEditingBankAccount(bankAccount);
    setBankAccountDialogOpen(true);
  };

  const handleSubmitBankAccount = async (data: CreateBankAccountInput) => {
    try {
      if (editingBankAccount) {
        await updateBankAccountMutation.mutateAsync({
          id: editingBankAccount.id,
          data,
        });
        toast({
          title: "更新成功",
          description: "銀行帳戶已更新",
        });
      } else {
        await createBankAccountMutation.mutateAsync(data);
        toast({
          title: "新增成功",
          description: "銀行帳戶已新增",
        });
      }
      setBankAccountDialogOpen(false);
      setEditingBankAccount(undefined);
    } catch (error: any) {
      toast({
        title: "操作失敗",
        description: error.message || "請稍後再試",
        variant: "destructive",
      });
    }
  };

  const handleDeleteBankAccount = async () => {
    if (!deletingBankAccountId) return;

    try {
      await deleteBankAccountMutation.mutateAsync(deletingBankAccountId);
      toast({
        title: "刪除成功",
        description: "銀行帳戶已刪除",
      });
      setDeletingBankAccountId(undefined);
    } catch (error: any) {
      toast({
        title: "刪除失敗",
        description: error.message || "請稍後再試",
        variant: "destructive",
      });
    }
  };

  // ==================== 信用卡處理函式 ====================

  const handleCreateCreditCard = () => {
    setEditingCreditCard(undefined);
    setCreditCardDialogOpen(true);
  };

  const handleEditCreditCard = (creditCard: CreditCard) => {
    setEditingCreditCard(creditCard);
    setCreditCardDialogOpen(true);
  };

  const handleSubmitCreditCard = async (data: CreateCreditCardInput) => {
    try {
      if (editingCreditCard) {
        await updateCreditCardMutation.mutateAsync({
          id: editingCreditCard.id,
          data,
        });
        toast({
          title: "更新成功",
          description: "信用卡已更新",
        });
      } else {
        await createCreditCardMutation.mutateAsync(data);
        toast({
          title: "新增成功",
          description: "信用卡已新增",
        });
      }
      setCreditCardDialogOpen(false);
      setEditingCreditCard(undefined);
    } catch (error: any) {
      toast({
        title: "操作失敗",
        description: error.message || "請稍後再試",
        variant: "destructive",
      });
    }
  };

  const handleDeleteCreditCard = async () => {
    if (!deletingCreditCardId) return;

    try {
      await deleteCreditCardMutation.mutateAsync(deletingCreditCardId);
      toast({
        title: "刪除成功",
        description: "信用卡已刪除",
      });
      setDeletingCreditCardId(undefined);
    } catch (error: any) {
      toast({
        title: "刪除失敗",
        description: error.message || "請稍後再試",
        variant: "destructive",
      });
    }
  };

  // ==================== 信用卡群組處理函式 ====================

  const handleCreateGroup = () => {
    setEditingGroup(undefined);
    setGroupDialogOpen(true);
  };

  const handleEditGroup = (group: CreditCardGroupWithCards) => {
    setEditingGroup(group);
    setGroupDialogOpen(true);
  };

  const handleSubmitGroup = async (data: CreateCreditCardGroupInput) => {
    try {
      if (editingGroup) {
        await updateGroupMutation.mutateAsync({
          id: editingGroup.id,
          data: {
            name: data.name,
            shared_credit_limit: data.shared_credit_limit,
            note: data.note,
          },
        });
        toast({
          title: "更新成功",
          description: "信用卡群組已更新",
        });
      } else {
        await createGroupMutation.mutateAsync(data);
        toast({
          title: "新增成功",
          description: "信用卡群組已新增",
        });
      }
      setGroupDialogOpen(false);
      setEditingGroup(undefined);
    } catch (error: any) {
      toast({
        title: "操作失敗",
        description: error.message || "請稍後再試",
        variant: "destructive",
      });
    }
  };

  const handleDeleteGroup = async () => {
    if (!deletingGroupId) return;

    try {
      await deleteGroupMutation.mutateAsync(deletingGroupId);
      toast({
        title: "刪除成功",
        description: "信用卡群組已解散，卡片已恢復為獨立狀態",
      });
      setDeletingGroupId(undefined);
    } catch (error: any) {
      toast({
        title: "刪除失敗",
        description: error.message || "請稍後再試",
        variant: "destructive",
      });
    }
  };

  const handleRemoveCardFromGroup = async (groupId: string, cardId: string) => {
    try {
      await removeCardsFromGroupMutation.mutateAsync({
        id: groupId,
        data: { card_ids: [cardId] },
      });
      toast({
        title: "移除成功",
        description: "卡片已從群組移除",
      });
    } catch (error: any) {
      toast({
        title: "移除失敗",
        description: error.message || "請稍後再試",
        variant: "destructive",
      });
    }
  };

  // 取得可用於建立群組的卡片（不在任何群組中的卡片）
  const getAvailableCards = () => {
    if (!creditCards) return [];

    // 如果沒有群組資料，返回所有信用卡
    if (!groups || groups.length === 0) return creditCards;

    // 過濾掉已在群組中的卡片
    const cardsInGroups = new Set(
      groups.flatMap((group) => group.cards.map((card) => card.id))
    );
    return creditCards.filter((card) => !cardsInGroups.has(card.id));
  };

  return (
    <AppLayout title="使用者管理" description="管理您的銀行帳戶和信用卡資訊">
      {/* Main Content */}
      <div className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="flex flex-col gap-6">
          {/* 分類管理 */}
          <CategoryManagement />

          {/* 銀行帳戶區塊 */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>銀行帳戶</CardTitle>
                <Button onClick={handleCreateBankAccount}>
                  <PlusIcon className="mr-2 h-4 w-4" />
                  新增銀行帳戶
                </Button>
              </div>
            </CardHeader>
            <BankAccountList
              bankAccounts={bankAccounts}
              isLoading={bankAccountsLoading}
              onEdit={handleEditBankAccount}
              onDelete={setDeletingBankAccountId}
            />
          </Card>

          {/* 信用卡區塊 */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>信用卡</CardTitle>
                <div className="flex gap-2">
                  <Button variant="outline" onClick={handleCreateGroup}>
                    <FolderPlusIcon className="mr-2 h-4 w-4" />
                    建立群組
                  </Button>
                  <Button onClick={handleCreateCreditCard}>
                    <PlusIcon className="mr-2 h-4 w-4" />
                    新增信用卡
                  </Button>
                </div>
              </div>
            </CardHeader>
            <CreditCardList
              creditCards={creditCards}
              groups={groups}
              isLoading={creditCardsLoading || groupsLoading}
              onEdit={handleEditCreditCard}
              onDelete={setDeletingCreditCardId}
              onEditGroup={handleEditGroup}
              onDeleteGroup={setDeletingGroupId}
              onRemoveCardFromGroup={handleRemoveCardFromGroup}
            />
          </Card>
        </div>
      </div>

      {/* 銀行帳戶 Dialog */}
      <Dialog
        open={bankAccountDialogOpen}
        onOpenChange={setBankAccountDialogOpen}
      >
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              {editingBankAccount ? "編輯銀行帳戶" : "新增銀行帳戶"}
            </DialogTitle>
          </DialogHeader>
          <BankAccountForm
            bankAccount={editingBankAccount}
            onSubmit={handleSubmitBankAccount}
            onCancel={() => setBankAccountDialogOpen(false)}
            isSubmitting={
              createBankAccountMutation.isPending ||
              updateBankAccountMutation.isPending
            }
          />
        </DialogContent>
      </Dialog>

      {/* 信用卡 Dialog */}
      <Dialog
        open={creditCardDialogOpen}
        onOpenChange={setCreditCardDialogOpen}
      >
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              {editingCreditCard ? "編輯信用卡" : "新增信用卡"}
            </DialogTitle>
          </DialogHeader>
          <CreditCardForm
            creditCard={editingCreditCard}
            onSubmit={handleSubmitCreditCard}
            onCancel={() => setCreditCardDialogOpen(false)}
            isSubmitting={
              createCreditCardMutation.isPending ||
              updateCreditCardMutation.isPending
            }
          />
        </DialogContent>
      </Dialog>

      {/* 刪除銀行帳戶確認 Dialog */}
      <AlertDialog
        open={!!deletingBankAccountId}
        onOpenChange={(open) => !open && setDeletingBankAccountId(undefined)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>確認刪除</AlertDialogTitle>
            <AlertDialogDescription>
              確定要刪除此銀行帳戶嗎？此操作無法復原。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteBankAccount}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              刪除
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* 刪除信用卡確認 Dialog */}
      <AlertDialog
        open={!!deletingCreditCardId}
        onOpenChange={(open) => !open && setDeletingCreditCardId(undefined)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>確認刪除</AlertDialogTitle>
            <AlertDialogDescription>
              確定要刪除此信用卡嗎？此操作無法復原。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteCreditCard}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              刪除
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* 信用卡群組 Dialog */}
      <Dialog open={groupDialogOpen} onOpenChange={setGroupDialogOpen}>
        <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              {editingGroup ? "編輯信用卡群組" : "建立信用卡群組"}
            </DialogTitle>
          </DialogHeader>
          <CreditCardGroupForm
            group={editingGroup}
            availableCards={getAvailableCards()}
            onSubmit={handleSubmitGroup}
            onCancel={() => setGroupDialogOpen(false)}
            isSubmitting={
              createGroupMutation.isPending || updateGroupMutation.isPending
            }
          />
        </DialogContent>
      </Dialog>

      {/* 刪除群組確認 Dialog */}
      <AlertDialog
        open={!!deletingGroupId}
        onOpenChange={(open) => !open && setDeletingGroupId(undefined)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>確認解散群組</AlertDialogTitle>
            <AlertDialogDescription>
              確定要解散此信用卡群組嗎？群組內的卡片將恢復為獨立狀態。此操作無法復原。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteGroup}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              解散群組
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </AppLayout>
  );
}
