import io
import pandas as pd
import matplotlib.pyplot as plt


async def get_budget_expense_yearly_total_amount_plot(
    df: pd.DataFrame, tag: str = None
) -> io.BytesIO:
    buf = io.BytesIO()
    if tag is not None:
        df = df[(df["tagValue"] == tag)]

    df["amount"] = pd.to_numeric(df["amount"], errors="coerce")
    df["iso-date"] = pd.to_datetime(
        df["date"], format="%d/%m/%Y", errors="coerce"
    ).dt.strftime("%Y-%m-%d")

    df["year"] = pd.to_datetime(df["date"], format="%d/%m/%Y", errors="coerce").dt.year
    df["mount"] = pd.to_datetime(
        df["date"], format="%d/%m/%Y", errors="coerce"
    ).dt.month

    # Plot
    plot = (
        df.groupby("year")["amount"]
        .sum()
        .plot(
            kind="bar",
            figsize=(20, 10),
            title=f"{tag} Amount by Years",
            xlabel="Years",
            ylabel="Total Amount",
        )
    )  # or "pie", "line", etc.

    for bar in plot.patches:
        plot.text(
            bar.get_x() + bar.get_width() / 2,  # X position
            bar.get_height(),  # Y position
            f"{bar.get_height():,.2f}",  # Format the label (2 decimals)
            ha="center",
            va="bottom",
            fontsize=15,
            rotation=60,
        )
    plt.xticks(rotation=45, ha="right")
    plt.tight_layout()
    plt.savefig(buf, format="png")
    plt.close()
    buf.seek(0)

    return buf
