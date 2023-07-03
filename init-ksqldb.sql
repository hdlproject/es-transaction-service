CREATE STREAM topUpEvents (id VARCHAR, userId INT, amount BIGINT)
    WITH (kafka_topic='top-up-events', value_format='json', partitions=1);

CREATE TABLE topUpEventsByUserId AS
    SELECT userId,
           SUM(amount) AS totalAmount
    FROM topUpEvents
    GROUP BY userId;

SELECT * FROM topUpEvents EMIT CHANGES;
