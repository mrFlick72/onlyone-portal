import json
import os

import jwt
import requests
from flask import request
from jwt import get_unverified_header, decode

from user.adapter.thread.LocalThreadUserNameResolver import LocalThreadUserNameResolver

class UserNameInjectorFilter:

    def __init__(self):
        self.user_name_resolver = LocalThreadUserNameResolver.get_instance()
        self.public_keys = {}
        self.jwk_endpoint = f"{os.getenv('ACCOUNT_SERVICE_ISS')}/jwks"
        jwks = requests.get(self.jwk_endpoint).json()
        for jwk in jwks['keys']:
            kid = jwk['kid']
            self.public_keys[kid] = jwt.algorithms.RSAAlgorithm.from_jwk(json.dumps(jwk))

    def filter(self, user_name_claim="user_name"):
        if request.path not in ["/health"]:
            token = str(request.headers.get("authorization"))
            print(token)
            print(os.getenv('PHONE_BOOK_AUD_CLAIM'))
            decoded_token = decode(jwt=token, key=self.public_keys['123'], algorithms=['RS256'],
                                   audience=os.getenv('PHONE_BOOK_AUD_CLAIM'))
            self.user_name_resolver.set_user_name(decoded_token["user_name"])
        return None

    def decode(self, token):
        kid = get_unverified_header(token)['kid']
        key = self.public_keys[kid]

        return decode(token, key=key, algorithms=['RS256'])