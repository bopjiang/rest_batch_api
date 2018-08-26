set -x

curl -s -X POST http://localhost:8080/users -H "Content-Type:application/json" -d '{"name":"papa","age": 26}' | jq
curl -s -X POST http://localhost:8080/users -H "Content-Type:application/json" -d '{"name":"mama","age": 25}' | jq
curl -s -X POST http://localhost:8080/users -H "Content-Type:application/json" -d '{"name":"son","age": 1}'   | jq
curl -s -X GET http://localhost:8080/users | jq
curl -X DELETE http://localhost:8080/users/1 
curl -s -X GET http://localhost:8080/users | jq
curl -X POST http://localhost:8080/batch -H 'Content-Type:multipart/mixed; boundary=END_OF_PART' --data-binary @batch_delete.txt
curl -s -X GET http://localhost:8080/users | jq

