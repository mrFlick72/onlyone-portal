import threading

from user.domain.user_name_resolver import UserNameResolver

class LocalThreadUserNameResolver(UserNameResolver):
    __instance = None

    @staticmethod
    def get_instance():
        """ Static access method. """
        if LocalThreadUserNameResolver.__instance is None:
            LocalThreadUserNameResolver.__instance = LocalThreadUserNameResolver()
        return LocalThreadUserNameResolver.__instance

    def __init__(self):
        self.context = threading.local()

    def get_user_name(self) -> str:
        return self.context.user_name

    def set_user_name(self, user_name: str):
        self.context.user_name = user_name