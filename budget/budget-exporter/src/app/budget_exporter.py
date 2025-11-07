"""Small helper to perform an HTTP GET with an OAuth2 Bearer token.

This module provides a minimal, well-documented helper that other code
in the repo can import. It intentionally keeps dependencies small and
uses the popular `requests` library.

Usage patterns:
- Pass an access token directly: `budget_expense_export(url, token)`
- Read token from environment (useful when you call `idp/get_access_token.sh` and export the token)

Note: make sure `requests` is available in the environment (install with
`pip install requests` or add it to the service's `requirements.txt`).
"""

from typing import Any, Dict
import os
import requests
from dotenv import load_dotenv
from functools import reduce
import operator
import csv

load_dotenv()


def get_budget_expense_for(year: str, month: str) -> Dict[str, Any]:
    """Perform an HTTP GET to `api_url` using an OAuth2 Bearer token.

    Args:
        api_url: Full URL to issue the GET request to.
        access_token: OAuth2 access token (just the token string, not the JSON).
        params: Optional query parameters.
        timeout: Request timeout in seconds.

    Returns:
        Parsed JSON response as a dict.

    Raises:
        requests.RequestException: for network-related errors
        requests.HTTPError: if the response status is 4xx/5xx (raised by raise_for_status)
    """
    access_token = os.environ.get("OAUTH2_ACCESS_TOKEN")
    api_url = f"{os.environ.get("API_BASE_URL")}/budget/expense"
    params = f"q=year={year};month={month}"
    headers = {
        "Authorization": f"Bearer {access_token}",
        "Accept": "application/json",
    }

    resp = requests.get(api_url, headers=headers, params=params)
    # Raise for non-2xx responses with a helpful exception
    resp.raise_for_status()

    # If the endpoint returns non-JSON, this will raise a ValueError — caller can catch it.
    daily_items = [
        item["budgetExpenseRepresentationList"]
        for item in resp.json()["dailyBudgetExpenseRepresentationList"]
    ]
    return reduce(operator.concat, daily_items)


def store(budget_expenses, year):
    out_dir = os.environ.get("BUDGET_EXPENSE_CSV_FILE_PATH")
    os.makedirs(out_dir, exist_ok=True)
    with open(
        f"{out_dir}/budget_expense_{year}.csv",
        "a",
    ) as file:
        writer = csv.DictWriter(file, fieldnames=budget_expenses[0].keys())
        writer.writeheader()
        writer.writerows(budget_expenses)


if __name__ == "__main__":
    # Simple CLI example for local testing
    import argparse

    p = argparse.ArgumentParser(description="Call an API endpoint with an OAuth2 token")
    p.add_argument("year", help="this is the year to use for the export job")
    args = p.parse_args()

    years = str(args.year).split("=")[1].split("-")
    print(years)
    start_year = int(years[0])
    end_year = int(years[1]) + 1

    for year in range(start_year, end_year):
        print("I'm working on the year: ", year)
        for month in range(1, 13):
            yearly_budget_expense = get_budget_expense_for(year=year, month=month)
            store(yearly_budget_expense, year)
