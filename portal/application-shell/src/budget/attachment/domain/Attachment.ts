export type BudgetType = "expense" | "revenue"

export type AttachmentTarget = {
    budgetId: string
    budgetType: BudgetType
    date: string
    attachmentId?: string
}
