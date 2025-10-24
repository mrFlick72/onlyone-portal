from dependency_injector import containers, providers
from app.user.adapter.thread.local_thread_user_name_resolver import (
    LocalThreadUserNameResolver,
)


class UserConfigContainer(containers.DeclarativeContainer):

    get_user_name_resolver = providers.Singleton(
        LocalThreadUserNameResolver,
    )
