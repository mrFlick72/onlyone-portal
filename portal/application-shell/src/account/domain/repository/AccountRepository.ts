import { getAccountApiBaseUrl } from "../../../config/ConfigLoader";
import Account from "../Account";

export async function getAccountData(): Promise<Account> {
    const baseUrl = await getAccountApiBaseUrl()
    const apiResult = await fetch(`${baseUrl}/api/account/user-account`, {
        method: "GET",
        mode: "cors",
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
        },
        credentials: 'include'
    })
    return apiResult.json()
}


export async function save(account: Account): Promise<void> {
    const baseUrl = await getAccountApiBaseUrl()
    await fetch(`${baseUrl}/api/account/user-account`, {
        method: "PUT",
        mode: "cors",
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
        },
        body: JSON.stringify(account),
        credentials: 'include'
    });

    return Promise.resolve()
}