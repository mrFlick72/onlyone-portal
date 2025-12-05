ACCESS_TOKEN="eyJraWQiOiI1YTZlNjQ5NC0yZmViLTQ2ZTQtOTI0ZS0wYWU3Yzk4Zjk4NDkiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJ2YWxlcmlvLnZhdWRpQGdtYWlsLmNvbSIsImF1ZCI6ImxvY2FsLm9uZS1wb3J0YWwiLCJuYmYiOjE3NjQ5NzI4MTksInVzZXJfbmFtZSI6InZhbGVyaW8udmF1ZGlAZ21haWwuY29tIiwic2NvcGUiOlsib3BlbmlkIl0sImlzcyI6Imh0dHA6Ly9sb2NhbC5hcGkudmF1dGhlbnRpY2F0b3IuY29tOjkwOTAiLCJncm91cHMiOltdLCJleHAiOjE3Njg1NzI4MTksImlhdCI6MTc2NDk3MjgxOSwianRpIjoiZjNkMzM5NjEtMTRmYy00ZTg2LTgxNmMtM2Q1NWUwYTVkNDU2IiwiYXV0aG9yaXRpZXMiOlsiUk9MRV9CVURHRVQiLCJST0xFX1VTRVIiLCJGQUNUT1JfUEFTU1dPUkQiXX0.guGwwhL81UmzgdmO8nByaL2xn92P011-1PhTQ-bVYPeW5RUiNv07WMPQSiuQUTOCaxeZKq6bK1h5cvar4PIfa0YtJe3hPclhEWEs-RlD3KHH7un56CAKNq3Yh0xje5U1PyA0y1Ofew9MJ0UiT8nqBJEKknhHwoO3AOIWVymxWbSzkWQU4NRAa1ziudEJ2W52LdTBi6YonSS5tIzhwPhl4UburqdEqadKJFbfPb4BTYo3V5IPhIXEq1b37WZW2Nyg8DFG3CVexzo15mHeYA6dcNSaw9dQfOqUjDVH5Cv1BAchKQqHW7DbNBIKPDQ5578uC5hyn7eFOm1egrEMi8ZmNA"

curl -v -X GET -H "Content-Type: application/json" -H "Authorization: Bearer ${ACCESS_TOKEN}" http://localhost:3030/budget/revenue?q=year=2018 





curl -v -X POST -H "Content-Type: application/json" -H "Authorization: Bearer ${ACCESS_TOKEN}" -d '{
    "date": "10/10/2025",
    "note": "It is a Note",
    "amount": "15000.00"
}' http://localhost:3030/budget/revenue 

