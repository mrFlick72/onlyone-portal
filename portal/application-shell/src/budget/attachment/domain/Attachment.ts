export type BudgetType = "expense" | "revenue"

export type AttachmentTarget = {
    budgetId: string
    budgetType: BudgetType
    date: string
    attachmentId?: string
}

export type AttachmentMetadata = {
    attachmentId: string
    fileName: string
    owner: string
    budgetId: string
    budgetType: string
}
