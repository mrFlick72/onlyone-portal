import { getRevenueApiBaseUrl } from "../../../config/ConfigLoader";
import BudgetRevenue from "./BudgetRevenue";

const BUDGET_REVENUE_URI = (baseUrl: string, budgetRevenueId?: string) => budgetRevenueId ?
    `${baseUrl}/budget/revenue/${budgetRevenueId}` :
    `${baseUrl}/budget/revenue`

const budgetRevenueWith = (baseUrl: string, year: string) => `${baseUrl}/budget/revenue?q=year=${year}`


export async function deleteBudgetRevenue(budgetRevenueId: string) {
    const baseUrl = await getRevenueApiBaseUrl();
    return fetch(BUDGET_REVENUE_URI(baseUrl, budgetRevenueId), {
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            'Accept': 'application/json'
        },
        credentials: 'include',
        method: "delete"
    })
}

export async function findBudgetRevenue(year: string) {
    const baseUrl = await getRevenueApiBaseUrl();
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

export async function saveBudgetRevenue(budgetRevenue: BudgetRevenue) {
    const baseUrl = await getRevenueApiBaseUrl();
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
