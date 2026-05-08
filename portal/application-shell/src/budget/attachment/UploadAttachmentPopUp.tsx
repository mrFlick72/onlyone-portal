import React, { useCallback, useEffect, useState } from "react"
import {
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Divider,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    Typography,
} from "@mui/material"
import { AttachFile, CloudUpload, UploadFile } from "@mui/icons-material"
import YesAndNoButtonGroup from "../../components/layout/YesAndNoButtonGroup"
import { getAttachmentsFor, saveAttachment } from "./domain/AttachmentRepository"
import { AttachmentMetadata, AttachmentTarget } from "./domain/Attachment"

type UploadAttachmentPopUpProps = {
    open: boolean
    handleClose: () => void
    target: AttachmentTarget | null
    modal: {
        title: string
        uploadButtonLabel: string
        closeButtonLabel: string
        chooseFileLabel: string
        noFileSelectedLabel: string
        existingAttachmentsLabel: string
        noAttachmentsLabel: string
    }
    onUploaded?: () => void
}

const UploadAttachmentPopUp: React.FC<UploadAttachmentPopUpProps> = ({
    open,
    handleClose,
    target,
    modal,
    onUploaded,
}) => {
    const [file, setFile] = useState<File | null>(null)
    const [submitting, setSubmitting] = useState(false)
    const [attachments, setAttachments] = useState<AttachmentMetadata[]>([])

    const loadAttachments = useCallback(async (current: AttachmentTarget) => {
        const list = await getAttachmentsFor(current)
        setAttachments(list)
    }, [])

    useEffect(() => {
        if (!open) {
            setFile(null)
            setSubmitting(false)
            setAttachments([])
            return
        }
        if (target) {
            loadAttachments(target)
        }
    }, [open, target, loadAttachments])

    const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const selected = event.target.files?.[0] ?? null
        setFile(selected)
    }

    const handleUpload = async () => {
        if (!target || !file || submitting) {
            return
        }
        setSubmitting(true)
        const response = await saveAttachment(target, file)
        setSubmitting(false)
        if (response.status === 201 || response.status === 204) {
            onUploaded?.()
            await loadAttachments(target)
            setFile(null)
        }
    }

    return (
        <Dialog onClose={handleClose} open={open} fullWidth scroll="paper">
            <DialogTitle>{modal.title}</DialogTitle>

            <DialogContent>
                <Box sx={{ display: "flex", flexDirection: "column", gap: 2, py: 1 }}>
                    <Button
                        variant="outlined"
                        component="label"
                        startIcon={<UploadFile />}
                    >
                        {modal.chooseFileLabel}
                        <input
                            type="file"
                            hidden
                            onChange={handleFileChange}
                        />
                    </Button>

                    <Typography variant="body2" color="text.secondary">
                        {file ? file.name : modal.noFileSelectedLabel}
                    </Typography>
                </Box>

                <Divider sx={{ my: 2 }} />

                <Typography variant="subtitle2" sx={{ mb: 1 }}>
                    {modal.existingAttachmentsLabel}
                </Typography>

                {attachments.length === 0 ? (
                    <Typography variant="body2" color="text.secondary">
                        {modal.noAttachmentsLabel}
                    </Typography>
                ) : (
                    <List dense disablePadding>
                        {attachments.map((attachment) => (
                            <ListItem key={attachment.attachmentId} disableGutters>
                                <ListItemIcon sx={{ minWidth: 36 }}>
                                    <AttachFile fontSize="small" />
                                </ListItemIcon>
                                <ListItemText
                                    primary={attachment.fileName}
                                    secondary={attachment.owner}
                                />
                            </ListItem>
                        ))}
                    </List>
                )}
            </DialogContent>
            <DialogActions>
                <YesAndNoButtonGroup
                    yesIcon={<CloudUpload />}
                    yesFun={handleUpload}
                    noFun={handleClose}
                    buttonMessages={{
                        noLabel: modal.closeButtonLabel,
                        yesLabel: modal.uploadButtonLabel,
                    }}
                />
            </DialogActions>
        </Dialog>
    )
}

export default UploadAttachmentPopUp
