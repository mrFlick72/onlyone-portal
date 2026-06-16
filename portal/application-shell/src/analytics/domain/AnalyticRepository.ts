import { getAnalyticApiBaseUrl } from "../../config/ConfigLoader";

export type TagTotal = {
    tag: string
    total: string
}

export type YearTotal = {
    year: number
    total: string
}

export type TotalByTagSearchCriteria = {
    year: number
    month?: number
    tags: string[]
}

export type TotalByYearSearchCriteria = {
    fromYear: number
    toYear: number
    tag?: string
}

export type ReindexCriteria = {
    fromYear: number
    toYear: number
}

export type ReindexResult = {
    imported: number
}

export async function findTotalByTag(searchCriteria: TotalByTagSearchCriteria): Promise<TagTotal[]> {
    const baseUrl = await getAnalyticApiBaseUrl();
    const response = await fetch(`${baseUrl}/api/analytic/budget/expense/total-by-tag`, {
        method: "PUT",
        credentials: 'include',
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        },
        body: JSON.stringify(searchCriteria)
    })
    if (!response.ok) {
        throw new Error(`total-by-tag request failed with status ${response.status}`)
    }
    return await response.json()
}

export async function findTotalByYear(searchCriteria: TotalByYearSearchCriteria): Promise<YearTotal[]> {
    const baseUrl = await getAnalyticApiBaseUrl();
    const response = await fetch(`${baseUrl}/api/analytic/budget/expense/total-by-year`, {
        method: "PUT",
        credentials: 'include',
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        },
        body: JSON.stringify(searchCriteria)
    })
    if (!response.ok) {
        throw new Error(`total-by-year request failed with status ${response.status}`)
    }
    return await response.json()
}

export async function reindexBudgetExpense(criteria: ReindexCriteria): Promise<ReindexResult> {
    const baseUrl = await getAnalyticApiBaseUrl();
    const response = await fetch(`${baseUrl}/api/analytic/budget/expense/reindex`, {
        method: "POST",
        credentials: 'include',
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        },
        body: JSON.stringify(criteria)
    })
    if (!response.ok) {
        throw new Error(`reindex request failed with status ${response.status}`)
    }
    return await response.json()
}
