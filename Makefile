
generate_proto:
	cd ./pkg/grpc/ && ./generate.sh
psql:
	docker-compose exec postgres psql -d usersServiceDB -U usersServiceDBuser
docker_stop:
	docker-compose down -v
docker_build:
	docker-compose up -d --build