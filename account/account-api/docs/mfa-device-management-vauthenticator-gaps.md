# vauthenticator gaps found while designing MFA device management

Findings from reading vauthenticator's actual source (`auth-server/src/main/kotlin/com/vauthenticator/server/mfa/**`, `oauth2/clientapp/**`), not just `auth-server/docs/mfa.md`, while designing account-api's MFA onboard/list/delete feature. Ordered by how much they block this feature.

## 1. No delete capability, at any layer ([vauthenticator#338](https://github.com/mrFlick72/vauthenticator/issues/338))

There is no way to delete/revoke an already-associated MFA device — not at the REST layer (`MfaEnrolmentAssociationEndPoint`, `MfaChallengeEndPoint`), not in the `MfaAccountMethodsRepository` domain interface (`findBy`/`findAll`/`save`/`setAsDefault`/`getDefaultDevice` only), not in the DynamoDB adapter (`DynamoMfaAccountMethodsRepository` has no `deleteItem` call anywhere). This blocks shipping delete as part of this feature; see ADR `0001-mfa-device-management-ships-against-vauthenticator-gaps.md`.

## 2. Scope provisioning is two steps, not one

`mfa:enrollment` scope must be **both**:
- requested by the frontend's OAuth2 authorize call, **and**
- already present in the portal's client-app registration inside vauthenticator (`ClientApplication.scopes`).

The second part is an admin/ops action against vauthenticator's client-app registry — the portal's client app isn't seeded anywhere in vauthenticator's own code (`ClientApplicationSetUpJob` only seeds vauthenticator's own built-in `vauthenticator-management-ui` and `admin` apps). Until the portal's client app is reprovisioned with `mfa:enrollment`, `POST /api/mfa/enrollment`, `POST /api/mfa/associate`, and `PUT /api/mfa/device` will all 403 via `PermissionValidator.principalScopesValidation`, regardless of what scope the frontend requests. `GET /api/mfa/enrollment` (list) is unaffected — see #3.

## 3. List has no scope/permission check at all ([vauthenticator#340](https://github.com/mrFlick72/vauthenticator/issues/340))

`MfaMethodsEnrolmentEndPoint.findAllAssociatedEnrolledMfaMethods` (`GET /api/mfa/enrollment`) calls no `permissionValidator.validate(...)`, unlike its three sibling endpoints in the same controller (enroll, associate, set-default), which all require `mfa:enrollment`. Today, any authenticated caller can list a user's enrolled MFA devices regardless of granted scope. This looks like an oversight relative to the pattern set by the rest of the controller, not a deliberate choice — worth flagging to the vauthenticator maintainer directly.

## 4. Association failures are undifferentiated over HTTP ([vauthenticator#339](https://github.com/mrFlick72/vauthenticator/issues/339))

`POST /api/mfa/associate` can fail for distinct reasons — wrong code, expired ticket, already-associated device — via `InvalidTicketException` (carrying an `InvalidTicketCause` enum: `TICKET_EXPIRED`, `ALREADY_ASSOCIATED_MFA`, etc.). `ExceptionAdviceController` only maps `InsufficientClientApplicationScopeException`; `InvalidTicketException` is unhandled and falls through to Spring's default, unstructured `500`. The cause never crosses the HTTP boundary, so callers can't distinguish "wrong code, try again" from "this ticket expired, restart enrollment."

## 5. No associate-as-default in one call ([vauthenticator#342](https://github.com/mrFlick72/vauthenticator/issues/342))

`MfaMethodsEnrollmentAssociation.associate(ticketId, code, asDefaultMethod)` supports setting the newly-associated device as default in the same call, but the REST DTO (`MfaEnrollmentAssociationRequest`) only exposes `ticket` and `code` — no `asDefaultMethod` field — so `MfaEnrolmentAssociationEndPoint.associateMfaEnrollment` always calls it with the default (`false`). Consequence for us: enrolling a user's very first MFA device is three round trips against vauthenticator (enroll → associate → set-default), not two.

## 6. `docs/mfa.md` is stale/incomplete ([vauthenticator#341](https://github.com/mrFlick72/vauthenticator/issues/341))

- Doesn't document `GET /api/mfa/enrollment` (list), `PUT /api/mfa/device` (set default), or `PUT /api/mfa/challenge` (login-time challenge send) — all three exist and work.
- Describes `mfaChannel` as "your mfa channel: mail sms and so on", implying a channel-type label. The code proves it's actually the destination address itself (an email address or phone number) that the code is sent to — confirmed via `MfaFixture` and `MfaMethodsEnrollmentTest` (e.g. `mfaChannel = email`).
- Says SMS support is "soon"; it's fully wired today (`SnsSmsSenderService` via SNS, configured in `MfaConfig`).

## Design notes for our own code (not vauthenticator's problem)

- Since `GET /api/mfa/enrollment` returns `userName`/`mfaChannel` **masked** (`getEnrollmentsFor(userName, withMaskedSensibleInformation = true)`), any "is this method already enrolled" check in account-api or the frontend must key off `mfaMethod` (and/or `mfaDeviceId`) — comparing the masked channel string against `Account.email`/`Account.phone` will not match.
