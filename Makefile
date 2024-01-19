dev:
	docker compose down && docker compose up --build -d

live:
	docker compose down && docker compose up --build -d

down:
	docker compose down
