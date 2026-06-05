import pytest
from fastapi.testclient import TestClient

from app.server import app


@pytest.fixture
def client():
    return TestClient(app)


def test_hello_returns_200(client):
    response = client.get("/api/analytic/hello")
    assert response.status_code == 200


def test_hello_returns_message(client):
    response = client.get("/api/analytic/hello")
    assert response.json() == {"message": "hello world"}
