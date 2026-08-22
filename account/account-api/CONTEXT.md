# Account

Proxies the authenticated user's profile to and from the vauthenticator IDP; account-api owns no persistent state of its own.

## Language

**Account**:
The user's profile as held by vauthenticator (first name, last name, birth date, email, phone). Never persisted or cached locally — every read/update round-trips to the IDP.
_Avoid_: User, Profile

### MFA

**MFA Device**:
An enrolled multi-factor method that has completed Association, identified by an opaque device id. One device per Account may be marked default. Never persisted locally — every read/write round-trips to vauthenticator.
_Avoid_: Authenticator, MFA Method (that's the delivery mechanism, not the enrolled result)

**MFA Method**:
The delivery mechanism for an MFA code: `EMAIL_MFA_METHOD` or `SMS_MFA_METHOD`. Distinct from the Channel (destination) and the Device (the associated result).

**MFA Channel**:
The destination an MFA code is sent to for a given Method — an email address or phone number. This UI restricts Channel input to the Account's own known email/phone rather than accepting an arbitrary destination.

**MFA Enrollment**:
The two-step process of registering a new MFA Device: starting enrollment (returns an Enrollment Ticket) and confirming it (Association) with the code received on the Channel. Does not become an MFA Device until Association succeeds.

**Association**:
Confirming an MFA Enrollment with the code received on the Channel, turning it into an active MFA Device.
_Avoid_: Verification, Confirmation

**Enrollment Ticket**:
Opaque, short-lived token returned by starting an MFA Enrollment; exchanged together with the code during Association.
_Avoid_: Token (ambiguous with the OAuth2 access token)
