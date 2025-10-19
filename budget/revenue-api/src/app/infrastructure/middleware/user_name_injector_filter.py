import json
import os

import jwt
import requests
from flask import request
from jwt import get_unverified_header, decode

from app.user.domain.user_name_resolver import UserNameResolver


class UserNameInjectorFilter:

    def __init__(self, user_name_resolver : UserNameResolver):
        self.user_name_resolver = user_name_resolver
        self.public_keys = {}
        self.jwk_endpoint = f"{os.getenv('IDP_ISS')}/oauth2/jwks"
        self.load_jwks()
        # todo use  a logger instead of print
        # print(self.jwk_endpoint)

    def filter(self, user_name_claim="user_name"):
        if request.path not in ["/health"]:
            token = str(request.headers.get("authorization")[7:])  # remove "Bearer "
            current_kid = get_unverified_header(token)["kid"]
            # todo use  a logger instead of print
            # print("token")
            # print(token)
            # print("kid")
            # print(current_kid)
            decoded_token = decode(
                jwt=token,
                key=self.public_keys[current_kid],
                algorithms=["RS256"],
                options={"verify_aud": False},
            )
            self.user_name_resolver.set_user_name(decoded_token[user_name_claim])
        return None

    def load_jwks(self):
        response = requests.get(self.jwk_endpoint)
        jwks = response.json()
        for jwk in jwks["keys"]:
            kid = jwk["kid"]
            public_key = jwt.algorithms.RSAAlgorithm.from_jwk(json.dumps(jwk))
            self.public_keys[kid] = public_key
