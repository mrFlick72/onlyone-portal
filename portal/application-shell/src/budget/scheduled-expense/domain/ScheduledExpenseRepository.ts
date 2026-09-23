import { getBudgetApiBaseUrl } from "../../../config/ConfigLoader";
import ScheduledExpense, { ScheduledExpenseList } from "./ScheduledExpense";

// Scheduled Expense lives in budget-api (not revenue-api, despite revenue's
// own frontend repository still pointing at REVENUE_API_BASE_URL) — see
// budget-api/CLAUDE.md's Scheduled Expense API Contract.
const SCHEDULED_EXPENSE_URI = (baseUrl: string) => `${baseUrl}/api/budget/scheduled-expense`;

const authHeaders = (extra: Record<string, string> = {}) => ({
    "Authorization": `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
    "Accept": "application/json",
    ...extra,
});

export async function getAllScheduledExpenses(): Promise<ScheduledExpense[]> {
    const baseUrl = await getBudgetApiBaseUrl();
    const response = await fetch(SCHEDULED_EXPENSE_URI(baseUrl), {
        method: "GET",
        credentials: "include",
        headers: authHeaders(),
    });
    if (!response.ok) {
        return [];
    }
    const body = await response.json() as ScheduledExpenseList;
    return body.scheduledExpenses;
}

export async function createScheduledExpense(scheduledExpense: ScheduledExpense) {
    const baseUrl = await getBudgetApiBaseUrl();
    return fetch(SCHEDULED_EXPENSE_URI(baseUrl), {
        method: "POST",
        credentials: "include",
        headers: authHeaders({ "Content-Type": "application/json" }),
        body: JSON.stringify(scheduledExpense),
    });
}
