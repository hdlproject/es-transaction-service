-- create a stream backed by a topic
-- a message that is published to the topic will be inserted to the stream
-- a message that is inserted to the stream will be published to the topic
CREATE STREAM TOP_UP_EVENTS (id VARCHAR, user_id INT, amount BIGINT)
    WITH (kafka_topic='top-up-events', value_format='json', partitions=1);

-- create a table from a stream
-- new messages arrive at the stream will be ingested to the table
CREATE TABLE TOP_UP_EVENTS_BY_USER_ID AS
    SELECT user_id,
           SUM(amount) AS total_amount
    FROM TOP_UP_EVENTS
    GROUP BY user_id;
