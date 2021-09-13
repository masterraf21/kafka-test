docker-up:
	docker-compose up -d
docker-stop:
	docker-compose stop
docker-down:
	docker-compose down
build:
	cd consumer1; go build consumer1.go
	cd consumer2; go build consumer2.go
	cd producer; go build producer.go
	mkdir -p bin
	mv producer/producer bin
	mv consumer1/consumer1 bin
	mv consumer2/consumer2 bin
publish-single:
	./bin/producer publish-single
publish-multiple:
	./bin/producer publish-multiple
run-consumer1:
	./bin/consumer1
run-consumer2:
	./bin/consumer2