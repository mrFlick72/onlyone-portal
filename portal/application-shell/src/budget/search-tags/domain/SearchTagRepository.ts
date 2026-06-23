import { getTagApiBaseUrl } from "../../../config/ConfigLoader";

// The tag catalog is partitioned by scope. The management UI maintains one
// catalog per scope (e.g. "expense", "revenue") through these two calls.
export type TagScope = "expense" | "revenue";

// UNKNOWN is a technical, non-persisted sentinel tag that tag-api synthesizes
// into every scoped read so consumers (e.g. expense tagging) have a catch-all.
// It is returned here unchanged; consumers that present the catalog for
// editing rather than selection (the management UI) hide it themselves.
export const UNKNOWN_SENTINEL_KEY = "UNKNOWN";

export async function getSearchTagRegistry(scope: TagScope): Promise<Array<SearchTag>> {
    const baseUrl = await getTagApiBaseUrl();
    const response = await fetch(`${baseUrl}/api/tags/scope/${scope}`, {
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

export async function saveSearchTag(searchTag: SearchTag, scope: TagScope): Promise<void> {
    const baseUrl = await getTagApiBaseUrl();
    await fetch(`${baseUrl}/api/tags`, {
        method: "PUT",
        credentials: "same-origin",
        body: JSON.stringify({ "key": searchTag.key, "value": searchTag.value, "scope": scope }),
        headers: {
            Authorization: `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            "Content-Type": "application/json",
        },
    });
    return Promise.resolve();
}

// updateSearchTagValue changes only the Value of an existing tag — Key and
// Scope are immutable once a tag exists (tag-api ADR 0008).
export async function updateSearchTagValue(key: string, value: string, scope: TagScope): Promise<void> {
    const baseUrl = await getTagApiBaseUrl();
    await fetch(`${baseUrl}/api/tags/scope/${scope}/${key}`, {
        method: "PATCH",
        credentials: "same-origin",
        body: JSON.stringify({ "value": value }),
        headers: {
            Authorization: `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
            "Content-Type": "application/json",
        },
    });
    return Promise.resolve();
}

export async function deleteSearchTag(key: string, scope: TagScope): Promise<void> {
    const baseUrl = await getTagApiBaseUrl();
    await fetch(`${baseUrl}/api/tags/scope/${scope}/${key}`, {
        method: "DELETE",
        credentials: "same-origin",
        headers: {
            Authorization: `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
        },
    });
    return Promise.resolve();
}
