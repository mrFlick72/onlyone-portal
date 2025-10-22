import json
import os

import jwt
import requests
from fastapi import Request
from starlette.middleware.base import BaseHTTPMiddleware
from jwt import get_unverified_header, decode

from app.user.domain.user_name_resolver import UserNameResolver


class UserNameInjectorFilter(BaseHTTPMiddleware):

    def __init__(self, app, user_name_claim: str, user_name_resolver: UserNameResolver):
        super().__init__(app)
        self.user_name_resolver = user_name_resolver
        self.user_name_claim = user_name_claim
        self.public_keys = {}
        self.jwk_endpoint = f"{os.getenv('IDP_ISS')}/oauth2/jwks"
        self.load_jwks()
        # todo use  a logger instead of print
        # print(self.jwk_endpoint)

    async def dispatch(self, request: Request, call_next):
        if request.url.path not in ["/health"]:
            token = str(request.headers.get("authorization")[7:])  # remove "Bearer "
            current_kid = get_unverified_header(token)["kid"]
            # todo use  a logger instead of print
            print("token")
            print(token)
            print("kid")
            # print(current_kid)
            decoded_token = decode(
                jwt=token,
                key=self.public_keys[current_kid],
                algorithms=["RS256"],
                options={"verify_aud": False},
            )
            self.user_name_resolver.set_user_name(decoded_token[self.user_name_claim])

        response = await call_next(request)
        return response

    def load_jwks(self):
        response = requests.get(self.jwk_endpoint)
        jwks = response.json()
        for jwk in jwks["keys"]:
            kid = jwk["kid"]
            public_key = jwt.algorithms.RSAAlgorithm.from_jwk(json.dumps(jwk))
            self.public_keys[kid] = public_key
