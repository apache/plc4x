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

The PLC4X Connector streams data from and to any device accessible through the PLC4X interface.

## Source Connector

See `config/source.properties` for example configuration.

## Sink Connector

See `config/sink.properties` for example configuration.

## Quickstart

1) Download the latest version of Apache Kafka binaries from here: https://kafka.apache.org/downloads
2) Unpack the archive.
3) Copy the target/apache-kafka-0.4.0-SNAPSHOT.jar to the Kafka "libs" directory.
4) Copy the files in the "config" to Kafka's "configs" directory (maybe inside a "plc4x" subdirectory)
5) Open 4 console windows and change directory into that directory
6) Start Zookeeper: 
        
        bin/zookeeper-server-start.sh config/zookeeper.properties 
7) Start Kafka:
        
        bin/kafka-server-start.sh config/server.properties
8) Create the "test" topic:
        
        bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic test
9) Start the consumer:
        
        bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test --from-beginning

**Note:** Not quite sure here ... have to continue working on this ...

10) Start Kafka connect:
        
        bin/connect-standalone.sh config/connect-standalone.properties config/plc4x/source.properties
Now watch the console window with the "kafka-console-consumer". 

If you want to debug the connector, be sure to set some environment variables before starting Kafka-Connect:

        export KAFKA_DEBUG=y; export DEBUG_SUSPEND_FLAG=y;

In this case the startup will suspend till an IDE is connected via a remote-debugging session.