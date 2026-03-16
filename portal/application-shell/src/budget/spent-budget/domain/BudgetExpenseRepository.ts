import { getBudgetApiBaseUrl } from "../../../config/ConfigLoader";
import { BudgetExpense, BudgetExpenseSearchCriteria, SpentBudget } from "./BudgetExpense";

const BUDGET_EXPENSE_URI = (baseUrl: string, budgetExpenseId?: string) => budgetExpenseId ?
    `${baseUrl}/api/budget/expense/${budgetExpenseId}` :
    `${baseUrl}/api/budget/expense`

export async function saveBudgetExpense(budgetExpense: BudgetExpense) {
    const baseUrl = await getBudgetApiBaseUrl();
    return fetch(BUDGET_EXPENSE_URI(baseUrl, budgetExpense.id), {
        method: budgetExpense.id ? "PUT" : "POST",
        credentials: 'include',
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(budgetExpense)
    })
}

export async function findBudgetExpense(searchCriteria: BudgetExpenseSearchCriteria): Promise<SpentBudget> {
    const baseUrl = await getBudgetApiBaseUrl();
    let baseUri = BUDGET_EXPENSE_URI(baseUrl);
    const response = await fetch(baseUri, {
        method: "PUT",
        credentials: 'include',
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        },
        body: JSON.stringify(searchCriteria)
    })
    return await response.json()
}

export async function deleteBudgetExpense(budgetExpenseId: string) {
    const baseUrl = await getBudgetApiBaseUrl();
    return fetch([BUDGET_EXPENSE_URI(baseUrl), budgetExpenseId].join("/"), {
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        },
        method: "delete",
        credentials: 'include',
    })
}