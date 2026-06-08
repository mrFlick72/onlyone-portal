from app.infrastructure.security.local_thread_security_context_resolver import (
    LocalThreadSecurityContextResolver,
)
from dependency_injector import containers, providers


class SecurityConfigContainer(containers.DeclarativeContainer):

    security_context_resolver = providers.Singleton(
        LocalThreadSecurityContextResolver,
    )
