import { getAccountApiBaseUrl } from "../../../config/ConfigLoader";
import MfaDevice, { MfaAssociationRequest, MfaEnrollmentRequest, MfaEnrollmentTicket } from "../Mfa";

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

export async function startMfaEnrollment(request: MfaEnrollmentRequest): Promise<MfaEnrollmentTicket> {
    const baseUrl = await getAccountApiBaseUrl()
    const apiResult = await fetch(`${baseUrl}/api/account/mfa/enrollment`, {
        method: "POST",
        mode: "cors",
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
        },
        body: JSON.stringify(request),
        credentials: 'include'
    })

    if (!apiResult.ok) {
        throw new Error(`Failed to start MFA enrollment: ${apiResult.status}`)
    }

    return apiResult.json()
}

export async function associateMfaEnrollment(request: MfaAssociationRequest): Promise<void> {
    const baseUrl = await getAccountApiBaseUrl()
    const apiResult = await fetch(`${baseUrl}/api/account/mfa/associate`, {
        method: "POST",
        mode: "cors",
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
        },
        body: JSON.stringify(request),
        credentials: 'include'
    })

    if (!apiResult.ok) {
        throw new Error(`Failed to associate MFA enrollment: ${apiResult.status}`)
    }
}

export async function setDefaultMfaDevice(mfaDeviceId: string): Promise<void> {
    const baseUrl = await getAccountApiBaseUrl()
    const apiResult = await fetch(`${baseUrl}/api/account/mfa/device`, {
        method: "PUT",
        mode: "cors",
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
        },
        body: JSON.stringify({ mfaDeviceId }),
        credentials: 'include'
    })

    if (!apiResult.ok) {
        throw new Error(`Failed to set default MFA device: ${apiResult.status}`)
    }
}
