curl -X POST -H "Content-Type: application/json" -d @todo.json http://localhost:5050/todo-service/todo 

curl -X GET http://localhost:5050/todo-service/todo/1  | jq