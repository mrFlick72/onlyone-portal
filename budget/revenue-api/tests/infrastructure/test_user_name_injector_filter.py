# import pytest

# from flask import Flask

# from app.infrastructure.middleware.user_name_injector_filter import (
#     UserNameInjectorFilter,
# )
# from app.user.adapter.thread.local_thread_user_name_resolver import (
#     LocalThreadUserNameResolver,
# )


# class DummyResponse:
#     def __init__(self, json_data):
#         self._json = json_data

#     def json(self):
#         return self._json


# def test_filter_sets_user_name(monkeypatch):
#     # The `app` variable in the provided code is an instance of the Flask application. It is created
#     # using `Flask(__name__)`, which initializes a Flask object with the name of the current module.
#     # This `app` instance is then used to simulate request contexts for testing the middleware filter
#     # logic in the unit tests.
#     # The `app` variable in the provided code is an instance of the Flask application. It is created
#     # using `Flask(__name__)`, which initializes a Flask object with the name of the current module.
#     # This `app` instance is then used to simulate request contexts for testing the middleware filter
#     # logic in the unit tests.
#     app = Flask(__name__)

#     # Prepare filter and stub methods
#     uut = UserNameInjectorFilter(LocalThreadUserNameResolver().get_instance())

#     # Stub requests.get used in load_jwks to avoid network call
#     monkeypatch.setattr(
#         "app.infrastructure.middleware.user_name_injector_filter.requests.get",
#         lambda url: DummyResponse({"keys": []}),
#     )

#     # Replace public_keys directly
#     uut.public_keys = {"kid1": "fake-public-key"}

#     # Stub get_unverified_header to return desired kid
#     monkeypatch.setattr(
#         "app.infrastructure.middleware.user_name_injector_filter.get_unverified_header",
#         lambda token: {"kid": "kid1"},
#     )

#     # Stub decode to return payload with user_name
#     monkeypatch.setattr(
#         "app.infrastructure.middleware.user_name_injector_filter.decode",
#         lambda jwt, key, algorithms, options: {"user_name": "alice"},
#     )

#     uut.load_jwks()
#     resolver = LocalThreadUserNameResolver.get_instance()

#     with app.test_request_context(
#         "/api/resource", headers={"Authorization": "Bearer TOKEN"}
#     ):
#         # call filter; should set resolver user name
#         uut.filter(user_name_claim="user_name")

#     assert resolver.get_user_name().content == "alice"


# def test_filter_skips_health_path(monkeypatch):
#     app = Flask(__name__)
#     uut = UserNameInjectorFilter(LocalThreadUserNameResolver().get_instance())

#     # Ensure nothing is set if path is /health
#     resolver = LocalThreadUserNameResolver.get_instance()

#     with app.test_request_context("/health", headers={"Authorization": "Bearer TOKEN"}):
#         uut.filter()

#     # For health, filter returns None and resolver should not have user_name attribute set here
#     # If get_user_name raises, treat as not set
#     try:
#         _ = resolver.get_user_name()
#         # If it exists, fail the test
#         pytest.fail("user_name should not be set for /health path")
#     except Exception:
#         pass
