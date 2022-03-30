USERNAME = carlschader

run:
	docker-compose up --build

kill:
	docker-compose down

build:
	docker build -t poker-go-api:latest -f services/server/Dockerfile .

build-create-ranks:
	docker build -t poker-create-ranks:latest -f services/create-ranks/Dockerfile .

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
