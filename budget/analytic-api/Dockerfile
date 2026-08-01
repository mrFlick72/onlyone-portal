# Debian/glibc base (not Alpine/musl): psycopg[binary] and confluent-kafka only
# publish manylinux (glibc) wheels, so an Alpine build would fail to install them.
FROM python:3.12-slim

COPY pyproject.toml README.md /opt/application/
COPY src/ /opt/application/src/

RUN pip install /opt/application

WORKDIR /opt/application

CMD ["app"]
