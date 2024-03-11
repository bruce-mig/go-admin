mysql:
	docker run --name admin-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8.3

createdb:
	docker exec -it admin-mysql mysql -u root -psecret -e "CREATE DATABASE go_admin;"

server:
	 go run main.go