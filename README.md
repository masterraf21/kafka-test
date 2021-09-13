# Kafka Demo
This repo is a demonstration of Apache Kafka functionalism using Docker to create Kafka Container and Golang for simulating Kafka Producer and Consumer

## Prerequisites
* Docker
* Go 1.16+

## How to Run
* Build the binaries using 
  ```bash 
  make build
  ```
* Run Docker Containers using
  ```bash
  make docker-up
  ```
* Run the consumers using
  ```bash
  make run-consumer1 
  ```
  and
  ```bash
  make run-consumer2
* Publish a message to a single consumer (consumer1)
  ```bash
  make publish-single
  ```
* Publish a message to multiple consumers (consumer1 & consumer2)
  ```bash
  make publish-multiple
  ```

For full information of commands, please refer [Makefile](Makefile)