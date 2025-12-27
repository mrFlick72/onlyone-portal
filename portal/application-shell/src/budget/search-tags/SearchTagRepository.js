import {getTagApiBaseUrl} from "../../config/ConfigLoader";

export async function getSearchTagRegistry() {
    const baseUrl = await getTagApiBaseUrl();
    return fetch(`${baseUrl}/api/tags`, {
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
    const baseUrl = await getTagApiBaseUrl();
    return fetch(`${baseUrl}/api/tags`, {
        method: "PUT",
        credentials: 'same-origin',
        body: JSON.stringify(searchTag),
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            'Content-Type': 'application/json'
        },
    })
}