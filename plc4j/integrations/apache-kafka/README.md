//
//  Licensed to the Apache Software Foundation (ASF) under one or more
//  contributor license agreements.  See the NOTICE file distributed with
//  this work for additional information regarding copyright ownership.
//  The ASF licenses this file to You under the Apache License, Version 2.0
//  (the "License"); you may not use this file except in compliance with
//  the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//
:imagesdir: ../../images/
:icons: font

== https://kafka.apache.org/[Apache Kafka]

Apache Kafka is an open-source distributed event streaming platform used by thousands of
companies for high-performance data pipelines, streaming analytics, data integration, and
mission-critical applications.

# PLC4X Kafka Connectors

The PLC4X connectors have the ability to pass data between Kafka and devices using industrial protocols.
They can be built from source from the future 0.8 https://plc4x.apache.org/users/download.html[release] of
PLC4X or from the latest snapshot from https://github.com/apache/plc4x[github].
//They can also be downloaded from the confluent https://www.confluent.io/hub/[hub].

## Introduction

A connect worker is basically a producer or consumer process with a standard api that Kafka can use to manage it. It is
able to be run in two modes:-

- Standalone
- Distributed

Standalone allows you to run the connector locally from the command line without having to install the jar file on your
Kafka brokers.
In distributed mode the connector runs on the Kafka brokers, which requires you to install the jar file on all of your
brokers. It allows the worker to be distributed across the Kafka brokers to provide redundancy and load balancing.

## Quickstart

In order to start a Kafka Connect system the following steps have to be performed:

1) Download the latest version of Apache Kafka binaries from here: https://kafka.apache.org/downloads.

2) Unpack the archive.

3) Copy the `target/plc4j-apache-kafka-0.8.0-uber-jar.jar` to the Kafka `libs` or plugin directory specified
in the config/connect-distributed.properties file.

4) Copy the files in the `config` to Kafka's `config` directory.

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

## Source Connector

The starting configuration for your connect worker is provided by a configuration file. However, once the worker has
started the configuration can be changed using the connect REST API which is generally available on
http://localhost:8083/connectors. When running in distributed mode all the configuration needs to be done via the REST API.

A sample configuration file is provided in the PLC4X Kafka integration repository in the `config/plc4x-source.properties` directory..
This includes comments as well as meaningful properties that can be used with the worker.

The configuration of the connectors via the REST interface expects the same properties as are specified within the
example `config/plc4x-source.properties` file. These will need to be in JSON format and included with a couple of headers.
An example below shows the format it expects:-

    curl -X POST -H "Content-Type: application/json" --data '{"name": "plc-source-test", "config": {"connector.class":"org.apache.plc4x.kafka.Plc4xSourceConnector",
    // TODO: Continue here ...
    "tasks.max":"1", "file":"test.sink.txt", "topics":"connect-test" }}' http://localhost:8083/connectors


### Start a Kafka Connect Source Worker (Standalone)

Ideal for testing.

1) Start Kafka connect:

        bin/connect-standalone.sh config/connect-standalone.properties config/plc4x-source.properties

Now watch the console window with the "kafka-console-consumer".

If you want to debug the connector, be sure to set some environment variables before starting Kafka-Connect:

        export KAFKA_DEBUG=y; export DEBUG_SUSPEND_FLAG=y;

In this case the startup will suspend till an IDE is connected via a remote-debugging session.

### Start Kafka Connect Source Worker (Distributed Mode)

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

## Sink Connector

See `config/sink.properties` for an example configuration.

### Start a Kafka Connect Sink Worker (Standalone)

Ideal for testing.

1) Start Kafka connect:

        bin/connect-standalone.sh config/connect-standalone.properties config/plc4x-sink.properties

Now open console window with "kafka-console-producer".

Producing to the kafka topic using the sample packet shown below should result all the values included in the payload
being sent to the PLC using the mapping defined in the sink properties.

    {"schema":
        {"type":"struct","fields":
            [{"type":"struct","fields":
                [{"type":"boolean","optional":true,"field":"running"},
                 {"type":"boolean","optional":true,"field":"conveyorLeft"},
                 {"type":"boolean","optional":true,"field":"conveyorRight"},
                 {"type":"boolean","optional":true,"field":"load"},
                 {"type":"int32","optional":true,"field":"numLargeBoxes"},
                 {"type":"boolean","optional":true,"field":"unload"},
                 {"type":"boolean","optional":true,"field":"transferRight"},
                 {"type":"boolean","optional":true,"field":"transferLeft"},
                 {"type":"boolean","optional":true,"field":"conveyorEntry"},
                 {"type":"int32","optional":true,"field":"numSmallBoxes"}],
            "optional":false,"name":"org.apache.plc4x.kafka.schema.Field","field":"fields"},
        {"type":"int64","optional":false,"field":"timestamp"},
        {"type":"int64","optional":true,"field":"expires"}],
         "optional":false,"name":"org.apache.plc4x.kafka.schema.JobResult",
         "doc":"PLC Job result. This contains all of the received PLCValues as well as a received timestamp"},
    "payload":
        {"fields":
            {"running":false,"conveyorLeft":true,
             "conveyorRight":true,"load":false,
             "numLargeBoxes":1630806456,
             "unload":true,
             "transferRight":false,
             "transferLeft":true,
             "conveyorEntry":false,
             "numSmallBoxes":-1135309911},
         "timestamp":1606047842350,
         "expires":null}}


If you want to debug the connector, be sure to set some environment variables before starting Kafka-Connect:

        export KAFKA_DEBUG=y; export DEBUG_SUSPEND_FLAG=y;

In this case the startup will suspend till an IDE is connected via a remote-debugging session.

### Start Kafka Connect Sink Worker (Distributed Mode)

Ideal for production.

In this case the state of the node is handled by Zookeeper and the configuration of the connectors are distributed via Kafka topics.

    bin/kafka-topics --create --zookeeper localhost:2181 --topic connect-configs --replication-factor 3 --partitions 1 --config cleanup.policy=compact
    bin/kafka-topics --create --zookeeper localhost:2181 --topic connect-offsets --replication-factor 3 --partitions 50 --config cleanup.policy=compact
    bin/kafka-topics --create --zookeeper localhost:2181 --topic connect-status --replication-factor 3 --partitions 10 --config cleanup.policy=compact

Starting the worker is then as simple as this:

    bin /connect-distributed.sh config/connect-distributed.properties

The configuration of the Connectors is then provided via REST interface:

    curl -X POST -H "Content-Type: application/json" --data '{"name": "plc-sink-test", "config": {"connector.class":"org.apache.plc4x.kafka.Plc4xSinkConnector",
    // TODO: Continue here ...
    "tasks.max":"1", "file":"test.sink.txt", "topics":"connect-test" }}' http://localhost:8083/connectors

## Graceful Backoff

If an error occurs when reading or writing PLC addresses a graceful backoff has been implemented so that the PLC isn't
bombarded with requests. However as the number of connectors for each PLC should be limited to reduce the load on the PLC,
the graceful backoff shouldn't have a major impact.

For the source connector the PLC4X scraper logic is able to handle randomized polling rates on failures, this is buffered within the
connector, the poll rate of the connector has no affect on the PLC poll rate.

For the sink connector, if a write fails it is retried a configurable number of times with a timeout between each time.
A Retriable Exception is raised which provides jitter for the timing of the retries.

## Schema Compatibility

PLC4X specifies a very basic schema and leaves the majority of the implementation to the user. It contains the
following fields:-

-   "tags": - This is a customized structure that is formed by the tags defined in the connector configuration.
This allows the user to defined arbitrary tags within here all based on the PLC4X data types.

- "timestamp": - This is the timestamp at which the PLC4X connector processed the PLC request.

- "expires": - This tag is used by the sink connector. It allows it to discard the record if it is too old. A value
of 0 or null indicates that the record some never be discarded no matter how old it is.

As the majority of the schema is left to the user to define we expect to be able to provide backward compatibility
between the base schemas.

The schemas for the sink and source connectors are the same. This allows us to producer from one PLC and send the
data to a sink.
