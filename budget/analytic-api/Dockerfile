FROM python:3.14.3-alpine3.23

COPY pyproject.toml /opt/application/
COPY src/ /opt/application/src/

RUN pip install /opt/application

WORKDIR /opt/application

CMD ["analytic-api"]
