from abc import ABC, abstractmethod

from app.user.domain.user import UserName


class UserNameResolver(ABC):
    @abstractmethod
    def get_user_name(self) -> UserName: ...

    @abstractmethod
    def set_user_name(self, user_name: UserName) -> None: ...
