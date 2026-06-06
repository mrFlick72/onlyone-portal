import threading
from typing import cast

from app.user.domain.user_name_resolver import UserNameResolver
from app.user.domain.user import UserName


class LocalThreadUserNameResolver(UserNameResolver):
    __instance = None

    @staticmethod
    def get_instance() -> UserNameResolver:
        if LocalThreadUserNameResolver.__instance is None:
            LocalThreadUserNameResolver.__instance = LocalThreadUserNameResolver()
        return LocalThreadUserNameResolver.__instance

    def __init__(self) -> None:
        self.context = threading.local()

    def get_user_name(self) -> UserName:
        return cast(UserName, self.context.user_name)

    def set_user_name(self, user_name: UserName) -> None:
        self.context.user_name = user_name
