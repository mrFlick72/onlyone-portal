export type BudgetExpense = {
    id?: string,
    date: string,
    amount: number,
    note: string,
    tagKey: string,
    tagValue: string,
}

export type SavedBudgetExpense = {
    id?: string,
    date: string,
    amount: number,
    note: string,
    searchTag: { value: string, label: string }
}

export type BudgetExpenseSearchCriteria = {
    month: string,
    year: string
}