rm -rf ../venv
python3 -m venv ../venv
source ../venv/bin/activate
pip install -r ../requirements.txt

export BUDGET_API_CONFIG_FILE_LOCATION="../../local/.env"

export AWS_ACCESS_KEY_ID="xxx" 
export AWS_SECRET_ACCESS_KEY="xxx" 
export AWS_REGION="eu-central-1"

cd ../
pip install .

python build/lib/app/main.py 