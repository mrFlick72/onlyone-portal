CREATE TABLE IF NOT EXISTS budget_expense_projection (
    id           TEXT PRIMARY KEY,
    user_name    TEXT NOT NULL,
    expense_date DATE NOT NULL,
    amount       NUMERIC(14, 2) NOT NULL,
    note         TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS budget_expense_tag (
    expense_id TEXT NOT NULL REFERENCES budget_expense_projection (id) ON DELETE CASCADE,
    tag_key    TEXT NOT NULL,
    tag_value  TEXT NOT NULL,
    PRIMARY KEY (expense_id, tag_key)
);

CREATE INDEX IF NOT EXISTS idx_expense_user_date
    ON budget_expense_projection (user_name, expense_date);

CREATE INDEX IF NOT EXISTS idx_expense_tag_key
    ON budget_expense_tag (tag_key);
