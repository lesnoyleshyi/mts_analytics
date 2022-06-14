SERVICES = postgres_mts analytics

rmv:
		docker volume rm events_volume

stopc:
		docker stop $(SERVICES)

rmc:	stopc
		docker rm -rf $(SERVICES)

build:
		docker compose up --build

get_signed:
		curl -X GET localhost:8080/agreed

get_unsigned:
		curl -X GET localhost:8080/canceled

get_time:
		curl -X GET localhost:8080/total_time?id=a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11

enter_db:
		docker exec -it postgres_mts /bin/bash

zookeeper:
		docker compose up zookeeper -d

kafka:
		docker compose up kafka-ui kafka-1 kafka-2 kafka-3 -d

start:
		docker start mts_analytics-kafka-1-1 mts_analytics-kafka-2-1 mts_analytics-kafka-3-1 \
 					mts_analytics-kafka-ui-1


re: rmc build