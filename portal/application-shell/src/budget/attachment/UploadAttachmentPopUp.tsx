import React, { useEffect, useState } from "react"
import { Box, Button, Dialog, DialogActions, DialogContent, DialogTitle, Typography } from "@mui/material"
import { CloudUpload, UploadFile } from "@mui/icons-material"
import YesAndNoButtonGroup from "../../components/layout/YesAndNoButtonGroup"
import { saveAttachment } from "./domain/AttachmentRepository"
import { AttachmentTarget } from "./domain/Attachment"

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

    useEffect(() => {
        if (!open) {
            setFile(null)
            setSubmitting(false)
        }
    }, [open])

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
            handleClose()
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
