import { getBudgetApiBaseUrl } from "../../../config/ConfigLoader"
import { AttachmentTarget } from "./Attachment"

const ATTACHMENT_URI = (baseUrl: string) => `${baseUrl}/api/attachment`

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
