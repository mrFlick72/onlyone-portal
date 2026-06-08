import threading
from typing import cast

from app.infrastructure.security.security_context_resolver import (
    SecurityContextResolver,
)
from app.infrastructure.security.model import SecurityContext


class LocalThreadSecurityContextResolver(SecurityContextResolver):
    __instance = None

    @staticmethod
    def get_instance() -> SecurityContextResolver:
        if LocalThreadSecurityContextResolver.__instance is None:
            LocalThreadSecurityContextResolver.__instance = (
                LocalThreadSecurityContextResolver()
            )
        return LocalThreadSecurityContextResolver.__instance

    def __init__(self) -> None:
        self.context = threading.local()

    def get_security_context(self) -> SecurityContext:
        return cast(SecurityContext, self.context.security_context)

    def set_security_context(self, security_context: SecurityContext) -> None:
        self.context.security_context = security_context
