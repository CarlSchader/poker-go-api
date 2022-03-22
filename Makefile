USERNAME = carlschader

run:
	docker compose -f docker/docker-compose.yaml up --build

kill:
	docker compose -f docker/docker-compose.yaml down

build:
	docker build -t poker-go-api:latest -f docker/Dockerfile .

publish:
	docker login

	docker build -t ${USERNAME}/poker-go-api:arm -f docker/Dockerfile-arm .
	docker build -t ${USERNAME}/poker-go-api:amd -f docker/Dockerfile-amd .

	docker push ${USERNAME}/poker-go-api:arm
	docker push ${USERNAME}/poker-go-api:amd
	
	docker manifest create \
	${USERNAME}/poker-go-api:latest \
	--amend ${USERNAME}/poker-go-api:arm \
	--amend ${USERNAME}/poker-go-api:amd \

	docker manifest push --purge ${USERNAME}/poker-go-api:latest
