type MfaDevice = {
    userName: string,
    mfaMethod: "EMAIL_MFA_METHOD" | "SMS_MFA_METHOD",
    mfaChannel: string,
    mfaDeviceId: string,
    default: boolean
}

export default MfaDevice
