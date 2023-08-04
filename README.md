# ES Transaction Service

## How to Run

### Run Docker Compose
```shell
$ docker-compose up
```

### Run KSQL DB Script
```shell
$ docker exec -it ksqldb-cli ksql http://ksqldb-server:8088
> CREATE STREAM topUpEvents (id VARCHAR, userId INT, amount BIGINT)
    WITH (kafka_topic='top-up-events', value_format='json', partitions=1);
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

### Note
- Still can do validation on the broker side. The error is `Delivery failed: Broker: Broker failed to validate record`.
