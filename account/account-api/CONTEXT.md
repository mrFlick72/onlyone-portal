# Account

Proxies the authenticated user's profile to and from the vauthenticator IDP; account-api owns no persistent state of its own.

## Language

**Account**:
The user's profile as held by vauthenticator (first name, last name, birth date, email, phone). Never persisted or cached locally — every read/update round-trips to the IDP.
_Avoid_: User, Profile
