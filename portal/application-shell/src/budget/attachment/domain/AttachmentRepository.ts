import { getBudgetApiBaseUrl } from "../../../config/ConfigLoader"
import { AttachmentMetadata, AttachmentTarget, BudgetType } from "./Attachment"

const ATTACHMENT_URI = (baseUrl: string) => `${baseUrl}/api/attachment`
const ATTACHMENT_METADATA_URI = (baseUrl: string, budgetType: BudgetType, budgetId: string) =>
    `${baseUrl}/api/attachment/metadata/${budgetType}/${budgetId}`
const ATTACHMENT_CONTENT_URI = (baseUrl: string, attachmentId: string) =>
    `${baseUrl}/api/attachment/${attachmentId}/content`
const ATTACHMENT_RESOURCE_URI = (baseUrl: string, attachmentId: string) =>
    `${baseUrl}/api/attachment/${attachmentId}`

export async function saveAttachment(target: AttachmentTarget, file: File) {
    const baseUrl = await getBudgetApiBaseUrl()

    const body = new FormData()
    body.append("file", file)
    body.append("budgetId", target.budgetId)
    body.append("budgetType", target.budgetType)
    body.append("date", target.date)
    if (target.attachmentId) {
        body.append("attachmentId", target.attachmentId)
    }

    return fetch(ATTACHMENT_URI(baseUrl), {
        method: "POST",
        credentials: "include",
        headers: {
            Authorization: `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
        },
        body,
    })
}

export async function getAttachmentsFor(target: AttachmentTarget): Promise<AttachmentMetadata[]> {
    const baseUrl = await getBudgetApiBaseUrl()

    const response = await fetch(ATTACHMENT_METADATA_URI(baseUrl, target.budgetType, target.budgetId), {
        method: "GET",
        credentials: "include",
        headers: {
            Authorization: `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
        },
    })

    if (!response.ok) {
        return []
    }

    const data = (await response.json()) as AttachmentMetadata[] | null
    return data ?? []
}

export async function downloadAttachment(attachmentId: string, fileName: string): Promise<void> {
    const baseUrl = await getBudgetApiBaseUrl()

    const response = await fetch(ATTACHMENT_CONTENT_URI(baseUrl, attachmentId), {
        method: "GET",
        credentials: "include",
        headers: {
            Authorization: `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
        },
    })

    if (!response.ok) {
        return
    }

    const blob = await response.blob()
    const objectUrl = window.URL.createObjectURL(blob)
    const anchor = document.createElement("a")
    anchor.href = objectUrl
    anchor.download = fileName
    document.body.appendChild(anchor)
    anchor.click()
    anchor.remove()
    window.URL.revokeObjectURL(objectUrl)
}

export async function deleteAttachment(attachmentId: string) {
    const baseUrl = await getBudgetApiBaseUrl()

    return fetch(ATTACHMENT_RESOURCE_URI(baseUrl, attachmentId), {
        method: "DELETE",
        credentials: "include",
        headers: {
            Authorization: `Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`,
        },
    })
}
