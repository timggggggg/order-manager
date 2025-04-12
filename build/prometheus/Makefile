DOCKER_DIR = ${CURDIR}
DOCKER_YML = "${DOCKER_DIR}/docker-compose.yml"
ENV_NAME = "order-manager"

all: compose-up

compose-up:
	@docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} up -d

compose-down:
	@docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} stop

compose-rm:
	@docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} rm -fvs

compose-ps:
	@docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} ps

compose-rs:
	make compose-rm && make compose-up

.PHONY: compose-up compose-down compose-rm compose-ps compose-rs