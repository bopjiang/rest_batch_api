set -x

curl -i -s -X POST http://localhost:8080/users -H "Content-Type:application/json" -d '{"name":"papa","age": 26}'
curl -i -s -X POST http://localhost:8080/users -H "Content-Type:application/json" -d '{"name":"mama","age": 25}'
curl -i -s -X POST http://localhost:8080/users -H "Content-Type:application/json" -d '{"name":"son","age": 1}'
curl -i -s -X GET http://localhost:8080/users
curl -i -s -X DELETE http://localhost:8080/users/1
curl -i -s -X GET http://localhost:8080/users
curl -i -s -X POST http://localhost:8080/batch -H 'Content-Type:multipart/mixed; boundary=END_OF_PART' --data-binary @batch_delete.txt
curl -i -s -X GET http://localhost:8080/users

