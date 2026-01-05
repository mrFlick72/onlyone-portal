import os
from turtle import pd
from fastapi import FastAPI, Request, Response
import matplotlib
import pandas as pd
import glob

from budget_expense_total_amont_by_tags_plot import get_budget_expense_total_amount_by_tags_plot
from budget_expense_yearly_total_amount_plot import (
    get_budget_expense_yearly_total_amount_plot,
)

matplotlib.use("Agg")
import matplotlib.pyplot as plt
import io
import uvicorn


from dotenv import load_dotenv


load_dotenv()

app = FastAPI()

df = pd.concat(
    [
        pd.read_csv(f)
        for f in glob.glob(f"{os.getenv('BUDGET_EXPENSE_CSV_FILE_PATH')}/*.csv")
    ],
    ignore_index=True,
)

@app.get("/budget/expense/yearly-total")
async def total_yearly_plot(request: Request):
    params = request.query_params
    tag = params.get("tag", None)
    buf = await get_budget_expense_yearly_total_amount_plot(df, tag=tag)
    return Response(buf.read(), media_type="image/png")

@app.get("/budget/expense")
async def total_amount_by_tags_plot(request: Request):
    params = request.query_params
    tags = params.get("tags", "").split(",")
    print(tags)
    buf = await get_budget_expense_total_amount_by_tags_plot(df, tags=tags)
    return Response(buf.read(), media_type="image/png")


if __name__ == "__main__":
    uvicorn.run("server:app", host="0.0.0.0", port=3035, reload=True)
