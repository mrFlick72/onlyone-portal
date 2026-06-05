FROM python:3.14.3-alpine3.23

COPY build/lib /opt/application
COPY pyproject.toml /opt/application/pyproject.toml
COPY requirements.txt /opt/application/requirements.txt
RUN  pip install -r /opt/application/requirements.txt

WORKDIR /opt/application
RUN pip install .

CMD ["python", "app/main.py"]
