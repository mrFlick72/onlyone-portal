export type MfaMethod = "EMAIL_MFA_METHOD" | "SMS_MFA_METHOD"

type MfaDevice = {
    userName: string,
    mfaMethod: MfaMethod,
    mfaChannel: string,
    mfaDeviceId: string,
    default: boolean
}

export type MfaEnrollmentRequest = {
    mfaMethod: MfaMethod,
    mfaChannel: string
}

export type MfaEnrollmentTicket = {
    ticket: string
}

export type MfaAssociationRequest = {
    ticket: string,
    code: string
}

export default MfaDevice
