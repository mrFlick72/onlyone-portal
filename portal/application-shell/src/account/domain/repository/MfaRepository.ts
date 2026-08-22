import { getAccountApiBaseUrl } from "../../../config/ConfigLoader";
import MfaDevice from "../Mfa";

export async function getMfaDevices(): Promise<MfaDevice[]> {
    const baseUrl = await getAccountApiBaseUrl()
    const apiResult = await fetch(`${baseUrl}/api/account/mfa`, {
        method: "GET",
        mode: "cors",
        headers: {
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
        },
        credentials: 'include'
    })

    if (!apiResult.ok) {
        return []
    }

    return apiResult.json()
}
