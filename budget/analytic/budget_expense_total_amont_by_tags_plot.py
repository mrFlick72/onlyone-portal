import io
import pandas as pd
import numpy as np
import glob
import matplotlib.pyplot as plt


async def get_budget_expense_total_amount_by_tags_plot(
    df: pd.DataFrame, tags: list[str] = list()
) -> io.BytesIO:
    buf = io.BytesIO()
    if len(tags) != 0 and tags != [""]:
        df = df[(df["tagValue"].isin(tags))]

    df["amount"] = pd.to_numeric(df["amount"], errors="coerce")

    grouped = df.groupby("tagValue")["amount"].sum()

    # Plot
    plot = grouped.plot(
        kind="bar",
        figsize=(20, 10),
        title="Amount by Tag",
        xlabel="Tags",
        ylabel="Total Amount",
    )  # or "pie", "line", etc.

    for bar in plot.patches:
        plot.text(
            bar.get_x() + bar.get_width() / 2,  # X position
            bar.get_height(),  # Y position
            f"{bar.get_height():,.2f}",  # Format the label (2 decimals)
            ha="center",
            va="bottom",
            fontsize=8,
            rotation=60,
        )
    plt.xticks(rotation=45, ha="right")
    plt.tight_layout()
    plt.savefig(buf, format="png")
    plt.close()
    buf.seek(0)

    return buf
