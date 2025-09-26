from flask import Flask 
from infrastructure.management.HealthEndPoint import HealthEndPoint
from infrastructure.middleware.UserNameInjectorFilter import UserNameInjectorFilter

app = Flask(__name__)

@app.route("/")
def hello():
    return "Hello, World!"

HealthEndPoint(app)
user_name_injector_filter = UserNameInjectorFilter()
app.before_request(user_name_injector_filter.filter)

if __name__ == '__main__':
  app.run(host="0.0.0.0", port=3030, debug=True)