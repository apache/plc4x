<!--
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
  -->

# Using Kakfa Connect to pump PLC data from PLC to Elasticsearch using PLC4X

## Setup Kafka (Connect)

1) Download the latest version of Apache Kafka binaries from here: https://kafka.apache.org/downloads
2) Unpack the archive.
3) Copy the `target/plc4j-apache-kafka-0.7.0-SNAPSHOT-uber-jar.jar` to the Kafka `libs` directory.
4) Download the Kafka Connect Elasticsearch connector: https://www.confluent.io/hub/confluentinc/kafka-connect-elasticsearch
4) Copy the files in the `kafka-connect-config` to Kafka's `configs` directory.

## Start a Kafka

1) Open 3 console windows and change directory into that directory
2) Start Zookeeper:

        bin/zookeeper-server-start.sh config/zookeeper.properties

3) Start Kafka:

        bin/kafka-server-start.sh config/server.properties

4) Create the Kafka topics:

        bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic heartbeat

        bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic machineData

5) Start ElasticSearch

        ../../../Elastic/elasticsearch-7.5.2/bin/elasticsearch

6) Start Kibana

        ../../../Elastic/kibana-7.5.2-darwin-x86_64/bin/elasticsearch

5) Start the Kafka Connect PLC4X Source:

        bin/connect-standalone.sh config/plc4x-worker.properties config/plc4x-source.properties

5) Start the Kafka Connect ElasticSearch Sink:

        bin/connect-standalone.sh config/connect-standalone.properties config/elasticsearch-sink.properties
