# analytic-api

Analytics API for the OnlyOne Portal budget service.

## Setup

```bash
python3 -m venv venv && source venv/bin/activate
pip install -e ".[dev]"
```

## Run

```bash
ANALYTIC_API_CONFIG_FILE_LOCATION=local/.env analytic-api
```

## Test

```bash
pytest
```

## Type check

```bash
mypy
```
