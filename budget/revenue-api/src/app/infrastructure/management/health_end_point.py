from flask import Flask


class HealthEndPoint:

    def __init__(self,
                 app: Flask,
                 ):
        self.app = app
        app.add_url_rule("/health", "health", self.health, methods=['GET'])

    def health(self):
        return self.app.response_class(status=200)