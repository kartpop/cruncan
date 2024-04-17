network-up:
	@docker network create cruncan-network

network-down:
	@docker network rm cruncan-network

kafka-up:
	@docker-compose -f ./pkg/kafka/docker-compose.yml up -d

kafka-down:
	@docker-compose -f ./pkg/kafka/docker-compose.yml down

one-up:
	@docker-compose -f ./one/docker-compose.yml up -d

one-down:
	@docker-compose -f ./one/docker-compose.yml down


