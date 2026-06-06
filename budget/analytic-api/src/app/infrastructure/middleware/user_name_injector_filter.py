import json
import logging
import os
from typing import Any

import jwt
import requests
from app.user.domain.user import UserName
from fastapi import Request
from starlette.middleware.base import BaseHTTPMiddleware, RequestResponseEndpoint
from starlette.responses import Response as StarletteResponse
from starlette.types import ASGIApp
from jwt import get_unverified_header, decode

from app.user.domain.user_name_resolver import UserNameResolver


class UserNameInjectorFilter(BaseHTTPMiddleware):

    __logger = logging.getLogger(__name__)

    def __init__(self, app: ASGIApp, user_name_claim: str, user_name_resolver: UserNameResolver) -> None:
        super().__init__(app)
        self.user_name_resolver = user_name_resolver
        self.user_name_claim = user_name_claim
        self.public_keys: dict[str, Any] = {}
        self.jwk_endpoint = f"{os.getenv('IDP_ISS')}/oauth2/jwks"
        self.load_jwks()
        self.__logger.debug(self.jwk_endpoint)

    async def dispatch(self, request: Request, call_next: RequestResponseEndpoint) -> StarletteResponse:
        if request.url.path not in ["/health"] and not request.method == "OPTIONS":
            auth_header = request.headers.get("authorization")
            if not auth_header or not auth_header.startswith("Bearer "):
                return StarletteResponse(status_code=401)
            token = auth_header[7:]
            current_kid = get_unverified_header(token).get("kid")
            if current_kid not in self.public_keys:
                return StarletteResponse(status_code=401)
            decoded_token = decode(
                jwt=token,
                key=self.public_keys[current_kid],
                algorithms=["RS256"],
                issuer=os.getenv("IDP_ISS"),
                options={"verify_aud": False},
            )
            self.user_name_resolver.set_user_name(UserName(decoded_token[self.user_name_claim]))

        response = await call_next(request)
        return response

    def load_jwks(self) -> None:
        response = requests.get(self.jwk_endpoint)
        jwks = response.json()
        for jwk in jwks["keys"]:
            kid = jwk["kid"]
            public_key = jwt.algorithms.RSAAlgorithm.from_jwk(json.dumps(jwk))
            self.public_keys[kid] = public_key
