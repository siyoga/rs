codegen:
	go tool oapi-codegen -config oapi-codegen.yaml api/schema.yaml

run:
	docker-compose up --build

run-deamon:
	docker-compose up -d --build