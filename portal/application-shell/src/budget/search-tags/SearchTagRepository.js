import {getBudgetApiBaseUrl} from "../../config/ConfigLoader";

export async function getSearchTagRegistry() {
    const baseUrl = await getBudgetApiBaseUrl();
    return fetch(`${baseUrl}/budget-expense/search-tag`, {
        method: "GET",
        mode: "cors",
        credentials: 'include',
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            'Accept': 'application/json'
        },
    }).then(response => response.json())
}

export async function saveSearchTag(searchTag) {
    const baseUrl = await getBudgetApiBaseUrl();
    return fetch(`${baseUrl}/budget-expense/search-tag`, {
        method: "PUT",
        credentials: 'same-origin',
        body: JSON.stringify(searchTag),
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            'Content-Type': 'application/json'
        },
    })
}