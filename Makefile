
generate_proto:
	cd ./pkg/grpc/ && ./generate.sh
psql:
	docker-compose exec postgres psql -d usersServiceDB -U usersServiceDBuser