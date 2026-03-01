import {getBudgetApiBaseUrl} from "../../config/ConfigLoader";

const BUDGET_EXPENSE_URI = (baseUrl, budgetExpenseId) => budgetExpenseId ?
    `${baseUrl}/api/budget/expense/${budgetExpenseId}` :
    `${baseUrl}/api/budget/expense`

export async function saveBudgetExpense(budgetExpense) {
    const baseUrl = await getBudgetApiBaseUrl();
    return fetch(BUDGET_EXPENSE_URI(baseUrl,budgetExpense.id), {
        method: budgetExpense.id ? "PUT" : "POST",
        credentials: 'include',
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(budgetExpense)
    })
}

export async function findBudgetExpense(searchCriteria) {
    const baseUrl = await getBudgetApiBaseUrl();
    let baseUri = BUDGET_EXPENSE_URI(baseUrl);
    return fetch(baseUri, {
        method: "PUT",
        credentials: 'include',
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        },
        body: JSON.stringify(searchCriteria)
    }).then((response) => {
        return response.json();
    })
}

export async function deleteBudgetExpense(budgetExpenseId) {
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