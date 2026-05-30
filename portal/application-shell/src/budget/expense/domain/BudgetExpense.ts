export type SearchTag = {
    tagKey: string,
    tagValue: string,
}

export type BudgetExpense = {
    id?: string,
    date: string,
    amount: number,
    note: string,
    tag: SearchTag,
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
    year: string,
    searchTagList: string[]
}

export type SpentBudget = {
    dailyBudgetExpenseRepresentationList: DailyBudgetExpenseRepresentationList[],
    totalDetailList: TotalDetail[]
    total: string
}

export type DailyBudgetExpenseRepresentationList = {
    budgetExpenseRepresentationList: BudgetExpense[]
    date: string,
    total: string
}


export type TotalDetail = {
    searchTagKey: string,
    searchTagValue: string
    total: string
}