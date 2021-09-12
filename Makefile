docker-up:
	docker-compose up -d
docker-stop:
	docker-compose stop
docker-down:
	docker-compose down
build-producer:
	cd producer
	go build producer.go