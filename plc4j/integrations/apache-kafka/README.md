<!--

  Licensed to the Apache Software Foundation (ASF) under one or more
  contributor license agreements.  See the NOTICE file distributed with
  this work for additional information regarding copyright ownership.
  The ASF licenses this file to You under the Apache License, Version 2.0
  (the "License"); you may not use this file except in compliance with
  the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

-->

# Kafka Connect PLC4X Connector

The PLC4X Connector streams data from any device accessible through the PLC4X interface.

## Source Connector

See `config/source.properties` for example configuration.

## Quickstart

A Kafka Connect worker can be run in two modes: 
- Standalone
- Distributed

Both modes require a Kafka Broker instance to be available.
Kafka Connect is part of the Kafka distribution. 

In order to start a Kafka Connect system the following steps have to be performed:

1) Download the latest version of Apache Kafka binaries from here: https://kafka.apache.org/downloads
2) Unpack the archive.
3) Copy the `target/plc4j-apache-kafka-0.6.0-SNAPSHOT-uber-jar.jar` to the Kafka `libs` directory.
4) Copy the files in the `config` to Kafka's `configs` directory.

### Start a Kafka Broker

1) Open 4 console windows and change directory into that directory
2) Start Zookeeper: 
        
        bin/zookeeper-server-start.sh config/zookeeper.properties 

3) Start Kafka:
        
        bin/kafka-server-start.sh config/server.properties

4) Create the "test" topic:
        
        bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic test

5) Start the consumer:
        
        bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test --from-beginning

### Start a Kafka Connect Worker (Standalone)

Ideal for testing. 

1) Start Kafka connect:
        
        bin/connect-standalone.sh config/connect-standalone.properties config/plc4x-source.properties

Now watch the console window with the "kafka-console-consumer". 

If you want to debug the connector, be sure to set some environment variables before starting Kafka-Connect:

        export KAFKA_DEBUG=y; export DEBUG_SUSPEND_FLAG=y;

In this case the startup will suspend till an IDE is connected via a remote-debugging session.

### Start Kafka Connect Worker (Distributed Mode)

Ideal for production.

In this case the state of the node is handled by Zookeeper and the configuration of the connectors are distributed via Kafka topics.

    bin/kafka-topics --create --zookeeper localhost:2181 --topic connect-configs --replication-factor 3 --partitions 1 --config cleanup.policy=compact
    bin/kafka-topics --create --zookeeper localhost:2181 --topic connect-offsets --replication-factor 3 --partitions 50 --config cleanup.policy=compact
    bin/kafka-topics --create --zookeeper localhost:2181 --topic connect-status --replication-factor 3 --partitions 10 --config cleanup.policy=compact

Starting the worker is then as simple as this:

    bin /connect-distributed.sh config/connect-distributed.properties
    
The configuration of the Connectors is then provided via REST interface:

    curl -X POST -H "Content-Type: application/json" --data '{"name": "plc-source-test", "config": {"connector.class":"org.apache.plc4x.kafka.Plc4xSourceConnector", 
    // TODO: Continue here ...
    "tasks.max":"1", "file":"test.sink.txt", "topics":"connect-test" }}' http://localhost:8083/connectors
