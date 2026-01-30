import { getTagApiBaseUrl } from "../../../config/ConfigLoader";

export async function getSearchTagRegistry(): Promise<Array<SearchTag>> {
    const baseUrl = await getTagApiBaseUrl();
    const response = await fetch(`${baseUrl}/api/tags`, {
        method: "GET",
        mode: "cors",
        credentials: "include",
        headers: {
            Authorization: `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            Accept: "application/json",
        },
    });
    const jsonBody = await response.json();
    return Promise.resolve(
        jsonBody.map((item: any) => {
            return {
                key: item.key,
                value: item.value,
            } as SearchTag;
        })
    );
}

export async function saveSearchTag(searchTag: SearchTag): Promise<void> {
    const baseUrl = await getTagApiBaseUrl();
    await fetch(`${baseUrl}/api/tags`, {
        method: "PUT",
        credentials: "same-origin",
        body: JSON.stringify({ "key": searchTag.key, "value": searchTag.value }),
        headers: {
            Authorization: `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            "Content-Type": "application/json",
        },
    });
    return Promise.resolve();
}
