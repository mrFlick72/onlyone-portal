// ScheduledExpenseTag mirrors budget-api's web/tags.SearchTagRepresentation
// wire shape ({tagKey, tagValue}) — the same shape BudgetRevenue.ts's
// RevenueTag uses, not the {key, value} shape of the internal domain type.
export type ScheduledExpenseTag = {
    tagKey: string;
    tagValue: string;
};

export type ScheduledExpenseStatus = "ACTIVE" | "PAUSED";

// month and endDate are omitted (not null) by the backend when unset — see
// budget-api/CLAUDE.md's Scheduled Expense API Contract. A definition with no
// month recurs monthly; one with no endDate recurs forever.
type ScheduledExpense = {
    id?: string;
    description: string;
    amount: string;
    notes: string;
    tags: ScheduledExpenseTag[];
    day: number;
    month?: number;
    endDate?: string;
    status?: ScheduledExpenseStatus;
};

export type ScheduledExpenseList = {
    scheduledExpenses: ScheduledExpense[];
};

export default ScheduledExpense;
