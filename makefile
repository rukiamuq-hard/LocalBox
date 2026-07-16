up:
	docker compose -f deployment/docker/docker-compose.yml up -d --build

down:
	docker compose -f deployment/docker/docker-compose.yml down

logs:
	docker compose -f deployment/docker/docker-compose.yml logs -f
