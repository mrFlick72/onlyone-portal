# MFA device management ships against known vauthenticator gaps

vauthenticator's MFA support has gaps that surfaced while designing this feature — full detail in `docs/mfa-device-management-vauthenticator-gaps.md`. The two that shape this ADR:

1. **No delete** ([vauthenticator#338](https://github.com/mrFlick72/vauthenticator/issues/338)). There is no way to delete/revoke an already-associated MFA device — not at the REST layer, not in the `MfaAccountMethodsRepository` domain interface, not in the DynamoDB adapter.
2. **Undifferentiated association errors** ([vauthenticator#339](https://github.com/mrFlick72/vauthenticator/issues/339)). `POST /api/mfa/associate` can fail for distinct reasons (wrong code, expired ticket, already-associated device) but vauthenticator's `ExceptionAdviceController` only maps `InsufficientClientApplicationScopeException`; `InvalidTicketException` falls through to a generic, unstructured `500`. The `InvalidTicketCause` distinction exists internally but never crosses the HTTP boundary.
3. **List can't distinguish real devices from abandoned enrollment attempts** ([vauthenticator#343](https://github.com/mrFlick72/vauthenticator/issues/343)). Starting enrollment persists a row before any code is entered; if it's never associated, that row is permanent (same root cause as #338) and indistinguishable from a real device once returned by the list endpoint — its response drops the `associated` field entirely, so account-api can't filter it out either.

Rather than block the whole feature (onboard, list, set-default, delete) on changes to a different service, account-api and the application-shell UI ship enrollment, listing, and set-default now. Delete stays visible per device as a disabled icon with an explanatory tooltip, rather than being omitted. Association failures show one generic message ("verification failed, check the code or restart enrollment") rather than distinguishing the failure cause. The enrollment dialog can still be closed before a code is entered, but that action carries no guarantee of cleanup — abandoning enrollment at any point, cancel button or not, may leave a ghost entry in the list. All three gaps, along with the rest found in the linked doc, are filed as GitHub issues against `mrFlick72/vauthenticator` (#338–#343) rather than treated as prerequisites of this work.

A third gap does block a slice of this work directly: `mfa:enrollment` must be provisioned on the portal's client-app record inside vauthenticator (an admin action there, not a code change here) before enroll/associate/set-default will work — list is unaffected, since it has no scope check at all. Code for all four operations ships now regardless; provisioning proceeds in parallel as an explicit prerequisite task, not a blocker on writing the code.

## Considered Options

- Block this feature entirely until vauthenticator closes both gaps.
- Treat closing the gaps in vauthenticator as a prerequisite step of this same effort.
- Ship against the gaps now, with reduced UX fidelity where they bite (chosen).
