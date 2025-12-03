ACCESS_TOKEN="eyJraWQiOiI2MGVhYTc0Ni1hM2VkLTQ1ZjEtYWM0Ny05ZDAxM2IzNWQwYmIiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJ2YWxlcmlvLnZhdWRpQGdtYWlsLmNvbSIsImF1ZCI6ImxvY2FsLm9uZS1wb3J0YWwiLCJuYmYiOjE3NjQ4MDI3MDcsInVzZXJfbmFtZSI6InZhbGVyaW8udmF1ZGlAZ21haWwuY29tIiwic2NvcGUiOlsib3BlbmlkIl0sImlzcyI6Imh0dHA6Ly9sb2NhbC5hcGkudmF1dGhlbnRpY2F0b3IuY29tOjkwOTAiLCJncm91cHMiOltdLCJleHAiOjE3Njg0MDI3MDcsImlhdCI6MTc2NDgwMjcwNywianRpIjoiMTgzYTA4ZGYtZWZmNi00MWFhLWI2ZjktNDgwM2FmODA3Y2VjIiwiYXV0aG9yaXRpZXMiOlsiUk9MRV9CVURHRVQiLCJST0xFX1VTRVIiLCJGQUNUT1JfUEFTU1dPUkQiXX0.rQfd1RcwDXFn_zfMIegON-WfEOMwtGHh5NkatMqYAyFuHxpWD3g0ZAsD-zAmduSDfO-c0Z3M9pwOZRhPWQMfbm-SnEBn3kPk8fmKpDQbXeLg-qIbmU4XV5PepgWwNRqObV8mecpOejapbDocEFXXQd6qkSrIGB_NwS_P1-RKkxUaoqbhBCTEY0dYgfLtPneJXg28ByidolQUBxTvI-hRhGVhvJZArSX9uWlpnA33WiVKU5L6H5eVx7Ef2qk3aP12nP5lZlra914qqBZwLrFgAiJnz_8a6dH_nkxHLrMdsGm1sKN-YJIwRX3RzG-eFCn0AuFzayMxrjjPyHRF3330SA"
curl -v -X POST -H "Content-Type: application/json" -H "Authorization: Bearer ${ACCESS_TOKEN}" -d '{
    "date": "10/10/2025",
    "note": "It is a Note",
    "amount": "15000.00"
}' http://localhost:3030/budget/revenue 

curl -v -X GET -H "Content-Type: application/json" -H "Authorization: Bearer ${ACCESS_TOKEN}" http://localhost:3030/budget/revenue?q=year=2018 








