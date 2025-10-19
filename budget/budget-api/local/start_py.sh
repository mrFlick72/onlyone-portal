export BUDGET_API_CONFIG_FILE_LOCATION="../../local/.env"

export AWS_ACCESS_KEY_ID="xxx" 
export AWS_SECRET_ACCESS_KEY="xxx" 
export AWS_REGION="eu-central-1"

cd ../
pip install .

python venv/lib/python3.12/site-packages/app/main.py 