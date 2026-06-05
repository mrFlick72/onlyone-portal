#!/usr/bin/env bash
WITH_FRESH_BUILD=$1

if [ "$WITH_FRESH_BUILD" = "true" ]; then
    rm -rf ../venv
    python3 -m venv ../venv
    source ../venv/bin/activate
    pip install -r ../requirements.txt
fi

cd ../
pip install .

export ANALYTIC_API_CONFIG_FILE_LOCATION="local/.env"
python src/app/main.py
