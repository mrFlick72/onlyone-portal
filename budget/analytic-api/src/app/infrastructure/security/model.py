from dataclasses import dataclass


@dataclass
class UserName:
    content: str


@dataclass
class SecurityContext:
    token: str
    user_name: UserName
