import React, { useEffect, useState } from 'react'
import {
    Alert,
    Button,
    CircularProgress,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    FormControl,
    FormControlLabel,
    Radio,
    RadioGroup,
    Typography
} from "@mui/material";
import { Close } from "@mui/icons-material";
import FormInputTextField from "../components/form/FormInputTextField";
import { associateMfaEnrollment, startMfaEnrollment } from "./domain/repository/MfaRepository";
import MfaDevice, { MfaMethod } from "./domain/Mfa";
import { EnrollMfaDeviceDialogMessageBundle } from "../messages/MessageBundles";

type EnrollMfaDeviceDialogProps = {
    open: boolean;
    onClose: () => void;
    onEnrolled: () => void;
    email: string;
    phone: string;
    existingDevices: MfaDevice[];
    messages: EnrollMfaDeviceDialogMessageBundle;
}

type Step = "method" | "code"

const channelFor = (method: MfaMethod, email: string, phone: string) => method === "EMAIL_MFA_METHOD" ? email : phone

const EnrollMfaDeviceDialog: React.FC<EnrollMfaDeviceDialogProps> = ({ open, onClose, onEnrolled, email, phone, existingDevices, messages }) => {
    const [step, setStep] = useState<Step>("method")
    const [selectedMethod, setSelectedMethod] = useState<MfaMethod | "">("")
    const [ticket, setTicket] = useState("")
    const [code, setCode] = useState("")
    const [submitting, setSubmitting] = useState(false)
    const [errorMessage, setErrorMessage] = useState<string | null>(null)

    useEffect(() => {
        if (open) {
            setStep("method")
            setSelectedMethod("")
            setTicket("")
            setCode("")
            setErrorMessage(null)
        }
    }, [open])

    const alreadyEnrolled = (method: MfaMethod) => existingDevices.some(device => device.mfaMethod === method)

    const emailAvailable = !!email && !alreadyEnrolled("EMAIL_MFA_METHOD")
    const smsAvailable = !!phone && !alreadyEnrolled("SMS_MFA_METHOD")

    const handleClose = () => {
        setSubmitting(false)
        onClose()
    }

    const handleStartEnrollment = () => {
        if (!selectedMethod) {
            return
        }

        setSubmitting(true)
        setErrorMessage(null)
        startMfaEnrollment({ mfaMethod: selectedMethod, mfaChannel: channelFor(selectedMethod, email, phone) })
            .then(response => {
                setTicket(response.ticket)
                setStep("code")
                setSubmitting(false)
            })
            .catch(() => {
                setErrorMessage(messages.genericErrorMessage)
                setSubmitting(false)
            })
    }

    const handleAssociate = () => {
        setSubmitting(true)
        setErrorMessage(null)
        associateMfaEnrollment({ ticket, code })
            .then(() => {
                setSubmitting(false)
                onEnrolled()
            })
            .catch(() => {
                setErrorMessage(messages.genericErrorMessage)
                setSubmitting(false)
            })
    }

    return <Dialog onClose={handleClose} open={open} fullWidth scroll="paper">
        <DialogTitle>{messages.title}</DialogTitle>
        <DialogContent>
            {errorMessage && <Alert severity="error" sx={{ marginBottom: 2 }}>{errorMessage}</Alert>}

            {step === "method" && <>
                <Typography variant="subtitle1">{messages.methodStepTitle}</Typography>
                {!emailAvailable && !smsAvailable
                    ? <Typography>{messages.noMethodsAvailableLabel}</Typography>
                    : <FormControl>
                        <RadioGroup
                            value={selectedMethod}
                            onChange={event => setSelectedMethod(event.target.value as MfaMethod)}>
                            {emailAvailable && <FormControlLabel value="EMAIL_MFA_METHOD" control={<Radio />}
                                label={`${messages.emailMethodLabel} (${email})`} />}
                            {smsAvailable && <FormControlLabel value="SMS_MFA_METHOD" control={<Radio />}
                                label={`${messages.smsMethodLabel} (${phone})`} />}
                        </RadioGroup>
                    </FormControl>
                }
            </>}

            {step === "code" && <>
                <Typography variant="subtitle1">{messages.codeStepTitle}</Typography>
                <Typography>{messages.codeStepDescription}</Typography>
                <FormInputTextField id="mfaEnrollmentCode"
                    label={messages.codeInputLabel}
                    autoFocus={true}
                    required={true}
                    handler={(value) => setCode(value.target.value)}
                    value={code} />
            </>}
        </DialogContent>
        <DialogActions>
            <Button variant="contained" onClick={handleClose} color="primary">
                <Close /> {messages.cancelButtonLabel}
            </Button>

            {step === "method" && (emailAvailable || smsAvailable) &&
                <Button variant="contained" color="success" disabled={!selectedMethod || submitting} onClick={handleStartEnrollment}>
                    {submitting && <CircularProgress size={16} color="inherit" sx={{ marginRight: 1 }} />}
                    {messages.nextButtonLabel}
                </Button>
            }

            {step === "code" &&
                <Button variant="contained" color="success" disabled={!code || submitting} onClick={handleAssociate}>
                    {submitting && <CircularProgress size={16} color="inherit" sx={{ marginRight: 1 }} />}
                    {messages.confirmButtonLabel}
                </Button>
            }
        </DialogActions>
    </Dialog>
}

export default EnrollMfaDeviceDialog
