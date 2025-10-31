import pytest

from starlette.requests import Request
from app.infrastructure.middleware.user_name_injector_filter import (
    UserNameInjectorFilter,
)


class FakeResolver:
    def __init__(self):
        self.value = None

    def set_user_name(self, v):
        self.value = v


@pytest.mark.asyncio
async def test_dispatch_sets_user_name(monkeypatch):
    resolver = FakeResolver()

    # Patch load_jwks to avoid network calls and populate public_keys
    monkeypatch.setattr(
        "app.infrastructure.middleware.user_name_injector_filter.UserNameInjectorFilter.load_jwks",
        lambda self: self.public_keys.update({"kid1": "pubkey_placeholder"}),
    )

    # Patch jwt helpers imported into the module
    monkeypatch.setattr(
        "app.infrastructure.middleware.user_name_injector_filter.get_unverified_header",
        lambda token: {"kid": "kid1"},
    )
    monkeypatch.setattr(
        "app.infrastructure.middleware.user_name_injector_filter.decode",
        lambda jwt, key, algorithms, options: {"user_name": "alice"},
    )

    # Instantiate middleware (load_jwks was patched so no network)
    uut = UserNameInjectorFilter(
        app=None, user_name_claim="user_name", user_name_resolver=resolver
    )

    # Build a minimal ASGI scope with Authorization header
    scope = {
        "type": "http",
        "method": "GET",
        "path": "/some/path",
        "headers": [(b"authorization", b"Bearer faketoken")],
        "query_string": b"",
        "client": ("127.0.0.1", 0),
        "server": ("127.0.0.1", 80),
    }

    async def receive():
        return {"type": "http.request"}

    request = Request(scope, receive)

    async def call_next(req):
        from starlette.responses import Response

        return Response(status_code=200)

    # Act
    response = await uut.dispatch(request, call_next)

    # Assert
    assert response.status_code == 200
    assert resolver.value == "alice"


@pytest.mark.asyncio
async def test_dispatch_bypasses_health(monkeypatch):
    resolver = FakeResolver()

    # Ensure load_jwks not called / no network
    monkeypatch.setattr(
        "app.infrastructure.middleware.user_name_injector_filter.UserNameInjectorFilter.load_jwks",
        lambda self: None,
    )

    uut = UserNameInjectorFilter(
        app=None, user_name_claim="user_name", user_name_resolver=resolver
    )

    scope = {
        "type": "http",
        "method": "GET",
        "path": "/health",
        "headers": [],
        "query_string": b"",
        "client": ("127.0.0.1", 0),
        "server": ("127.0.0.1", 80),
    }

    async def receive():
        return {"type": "http.request"}

    request = Request(scope, receive)

    async def call_next(req):
        from starlette.responses import Response

        return Response(status_code=200)

    # Act
    response = await uut.dispatch(request, call_next)

    # Assert
    assert response.status_code == 200
    assert resolver.value is None
