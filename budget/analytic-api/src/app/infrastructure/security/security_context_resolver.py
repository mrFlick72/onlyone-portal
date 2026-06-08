from abc import ABC, abstractmethod

from app.infrastructure.security.model import SecurityContext



class SecurityContextResolver(ABC):
    @abstractmethod
    def get_security_context(self) -> SecurityContext: ...

    @abstractmethod
    def set_security_context(self, ctx: SecurityContext) -> None: ...
