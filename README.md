# ES Transaction Service
This project contains a transaction service. 
Through this project, I want to illustrate how I leverage some technologies to build a transaction service.
This may not be practical and inefficient since the purpose is solely for showing-off.
This project implements Clean Architecture.
This project is supposed to be a part of a bigger project which establish a complete microservice system.

## About the Project
### Config
This project needs some configurations. 
The configurations can be supplied through several methods.
Each method will have a priority and the higher one will take precedence.
Here is the list of the methods ordered from the highest priority:
- environment variable
- config file
- default value
This feature is possible by the help of the [viper package](https://github.com/spf13/viper).
### External Dependencies
This project works with some storages to function properly. Here is the list:
- MongoDB to store the transaction events (e.g. top-up balance).
- RabbitMQ to broadcast the events to the other system members.
Beside the main external dependencies above, this project also contains alternative dependencies such as:
- PostgreSQL to store the transaction events, but it differs from MongoDB cause it requires a predefined schema.
- Kafka to broadcast the events, but it differs from RabbitMQ cause it has higher performance and more complex features.
You may sense the high performance from the async publishing and async consuming.
    - Kafka Schema Registry complements Kafka by validating the message based on the registered schema.
    - ksqlDB provides a new interface to work with Kafka in the form of a database query.
      The database is built on top of a Kafka Streams which stream/interact with the data stored directly in the Kafka topic.

## How to Run
### Run Docker Compose
```shell
$ docker-compose up
```
### Create Topic
Execute below command to activate schema validation.
```shell
$ docker exec -it kafka \
kafka-configs --bootstrap-server localhost:9092 \
    --alter --entity-type topics \
    --entity-name custom-events \
    --add-config confluent.value.schema.validation=true
```
### Run KSQL DB Script
```shell
$ docker exec -it ksqldb-cli ksql http://ksqldb-server:8088
> CREATE STREAM topUpEvents (id VARCHAR, userId INT, amount BIGINT)
    WITH (kafka_topic='top-up-events', value_format='json', partitions=1);
```

### TODO
- Still can't do validation on the broker side. The error is `Delivery failed: Broker: Broker failed to validate record`.
