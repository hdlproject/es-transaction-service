# ES Transaction Service

## How to Run

### Run Docker Compose
```shell
$ docker-compose up
```

### Run KSQL DB script
```shell
$ docker exec -it ksqldb-cli ksql http://ksqldb-server:8088
> CREATE STREAM topUpEvents (id VARCHAR, userId INT, amount BIGINT)
    WITH (kafka_topic='top-up-events', value_format='json', partitions=1);
> 
```
