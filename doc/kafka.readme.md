============================================================
KAFKA EVENT STREAMING GUIDE FOR GO GATEWAY
============================================================

1. ARCHITECTURE OVERVIEW
------------------------------------------------------------
Kafka is the central nervous system for your microservices.
- Producer (Gateway): Sends audit logs, user activity, or ride requests.
- Topic: A category or feed name where records are stored.
- Partition: Scalability unit; allows multiple consumers to read in parallel.
- Broker: The Kafka server handling the data.



2. QUICK START (DOCKER)
------------------------------------------------------------
# Start Kafka
docker-compose up -d kafka

# Check if it's running
docker logs -f kafka-dev

3. CLI CHEAT SHEET (Inside Docker)
------------------------------------------------------------
# Enter the container:
docker exec -it kafka-dev /bin/bash

# --- Topic Management ---
# Create a topic
kafka-topics.sh --create --topic gateway_logs --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1

# List all topics
kafka-topics.sh --list --bootstrap-server localhost:9092

# Describe a topic (check partitions/offsets)
kafka-topics.sh --describe --topic gateway_logs --bootstrap-server localhost:9092

# --- Producing & Consuming ---
# Produce messages manually (type message and press Enter)
kafka-console-producer.sh --topic gateway_logs --bootstrap-server localhost:9092

# Consume messages (Read from the beginning)
kafka-console-consumer.sh --topic gateway_logs --from-beginning --bootstrap-server localhost:9092

4. GO INTEGRATION (segmentio/kafka-go)
------------------------------------------------------------
- Use the Singleton Writer for non-blocking event publishing.
- Always handle Close() on application shutdown to flush buffered messages.
- For high performance, use 'Async: true' in your Writer configuration.

5. ENVIRONMENT VARIABLES (.env)
------------------------------------------------------------
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=gateway_logs
KAFKA_GROUP_ID=gateway_service_group

============================================================