# Scripts

- source venv/bin/activate
  
- PYTHONPATH=src pytest --cov=app --cov-report=term-missing --cov-report=html

- pip install --no-cache-dir  --upgrade build

docker build -t mrflick72/budget/revenue-api:1 . 