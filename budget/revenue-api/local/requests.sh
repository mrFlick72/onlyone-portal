curl -v -X POST -H "Content-Type: application/json" -d '{
    "user_name": "johndoe",
    "date": "10/10/2025",
    "note": "It is a Note",
    "amount": 15000.00
}' http://localhost:3030/budget/revenue 
