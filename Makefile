DOCKER_HUB_USERNAME = carlschader
SERVICE_NAME = poker-go-api
CREATE_RANKS_SERVICE = poker-create-ranks

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
	docker run --privileged --rm tonistiigi/binfmt --install all
	docker buildx create --use --name ${SERVICE_NAME}

	docker buildx build \
	--push \
	--platform linux/amd64,linux/arm/v7,linux/arm64/v8,linux/ppc64le,linux/s390x \
	--tag ${DOCKER_HUB_USERNAME}/${SERVICE_NAME}:latest \
	-f services/server/Dockerfile .

	docker buildx stop ${SERVICE_NAME}
	docker buildx rm ${SERVICE_NAME}

publish-create-ranks:
	docker login
	docker run --privileged --rm tonistiigi/binfmt --install all
	docker buildx create --use --name ${CREATE_RANKS_SERVICE}

	docker buildx build \
	--push \
	--platform linux/amd64,linux/arm/v7,linux/arm64/v8,linux/ppc64le,linux/s390x \
	--tag ${DOCKER_HUB_USERNAME}/${CREATE_RANKS_SERVICE}:latest \
	-f services/create-ranks/Dockerfile .

	docker buildx stop ${CREATE_RANKS_SERVICE}
	docker buildx rm ${CREATE_RANKS_SERVICE}

publish-all:
	docker login
	docker run --privileged --rm tonistiigi/binfmt --install all
	docker buildx create --use --name ${SERVICE_NAME}

	docker buildx build \
	--push \
	--platform linux/amd64,linux/arm/v7,linux/arm64/v8,linux/ppc64le,linux/s390x \
	--tag ${DOCKER_HUB_USERNAME}/${SERVICE_NAME}:latest \
	-f services/server/Dockerfile .

	docker buildx build \
	--push \
	--platform linux/amd64,linux/arm/v7,linux/arm64/v8,linux/ppc64le,linux/s390x \
	--tag ${DOCKER_HUB_USERNAME}/${CREATE_RANKS_SERVICE}:latest \
	-f services/create-ranks/Dockerfile .

	docker buildx stop ${SERVICE_NAME}
	docker buildx rm ${SERVICE_NAME}
