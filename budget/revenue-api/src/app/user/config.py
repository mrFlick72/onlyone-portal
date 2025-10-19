from app.user.domain.user_name_resolver import UserNameResolver
from app.user.adapter.thread.local_thread_user_name_resolver import LocalThreadUserNameResolver


def user_name_resolver() -> UserNameResolver:
    return LocalThreadUserNameResolver.get_instance()
