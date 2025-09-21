import {getBudgetApiBaseUrl} from "../../config/ConfigLoader";

const BUDGET_REVENUE_URI = (baseUrl, budgetRevenueId) => budgetRevenueId ?
    `${baseUrl}/budget/revenue/${budgetRevenueId}` :
    `${baseUrl}/budget/revenue`

const budgetRevenueWith = (baseUrl, year) => `${baseUrl}/budget/revenue?q=year=${year}`


export async function deleteBudgetRevenue(budgetRevenueId) {
    const baseUrl = await getBudgetApiBaseUrl();
    return fetch(BUDGET_REVENUE_URI(baseUrl, budgetRevenueId), {
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            'Accept': 'application/json'
        },
        credentials: 'include',
        method: "delete"
    })
}

export async function findBudgetRevenue(year) {
    const baseUrl = await getBudgetApiBaseUrl();
    let responsePromise = await fetch(budgetRevenueWith(baseUrl, year), {
        method: "GET",
        credentials: 'include',
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            'Accept': 'application/json'
        }
    });
    return responsePromise.json()
}

export async function saveBudgetRevenue(budgetRevenue) {
    const baseUrl = await getBudgetApiBaseUrl();
    return fetch(BUDGET_REVENUE_URI(baseUrl, budgetRevenue.id), {
        method: budgetRevenue.id ? "PUT" : "POST",
        credentials: 'include',
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(budgetRevenue)
    })
}
