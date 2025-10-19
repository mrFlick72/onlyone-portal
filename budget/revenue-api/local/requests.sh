curl -v -X POST -H "Content-Type: application/json" -d '{
    "year": 2024,
    "month": 6,
    "department": "Sales",
    "amount": 15000
}' http://localhost:3030/budget/revenue 