# more arches: linux/amd64,linux/arm/v7,linux/arm64/v8,linux/ppc64le,linux/s390x

DOCKER_HUB_USERNAME = carlschader
SERVICE_NAME = poker-go-api
CREATE_RANKS_SERVICE = poker-create-ranks
ARCHES = linux/amd64,linux/arm64/v8
DOCKERFILE_PATH = services/server/Dockerfile
DOCKER_CONTEXT = .

CREATE_RANKS_PATH = services/create-ranks/Dockerfile

run:
	docker-compose up --build

kill:
	docker-compose down

run-create-ranks:
	go run services/create-ranks/generate/generate.go ranks.json

run-create-pockets:
	go run services/create-ranks/generate/generate.go ranks.json
	go run services/create-pockets/generate/generate.go ranks.json pockets.json

build:
	docker build -t poker-go-api:latest -f ${DOCKERFILE_PATH} ${DOCKER_CONTEXT}

build-create-ranks:
	docker build -t poker-create-ranks:latest -f ${CREATE_RANKS_PATH} ${DOCKER_CONTEXT}

publish:
	docker login
	docker run --privileged --rm tonistiigi/binfmt --install all
	docker buildx create --use --name ${SERVICE_NAME}

	docker buildx build \
	--push \
	--platform ${ARCHES} \
	--tag ${DOCKER_HUB_USERNAME}/${SERVICE_NAME}:latest \
	-f ${DOCKERFILE_PATH} ${DOCKER_CONTEXT}

	docker buildx stop ${SERVICE_NAME}
	docker buildx rm ${SERVICE_NAME}

publish-create-ranks:
	docker login
	docker run --privileged --rm tonistiigi/binfmt --install all
	docker buildx create --use --name ${CREATE_RANKS_SERVICE}

	docker buildx build \
	--push \
	--platform ${ARCHES} \
	--tag ${DOCKER_HUB_USERNAME}/${CREATE_RANKS_SERVICE}:latest \
	-f ${CREATE_RANKS_PATH} ${DOCKER_CONTEXT}

	docker buildx stop ${CREATE_RANKS_SERVICE}
	docker buildx rm ${CREATE_RANKS_SERVICE}

publish-all:
	docker login
	docker run --privileged --rm tonistiigi/binfmt --install all
	docker buildx create --use --name ${SERVICE_NAME}

	docker buildx build \
	--push \
	--platform ${ARCHES} \
	--tag ${DOCKER_HUB_USERNAME}/${SERVICE_NAME}:latest \
	-f ${DOCKERFILE_PATH} ${DOCKER_CONTEXT}

	docker buildx build \
	--push \
	--platform ${ARCHES} \
	--tag ${DOCKER_HUB_USERNAME}/${CREATE_RANKS_SERVICE}:latest \
	-f ${CREATE_RANKS_PATH} ${DOCKER_CONTEXT}

	docker buildx stop ${SERVICE_NAME}
	docker buildx rm ${SERVICE_NAME}
