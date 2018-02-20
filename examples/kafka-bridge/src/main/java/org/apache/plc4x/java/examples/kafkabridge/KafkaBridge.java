/*
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
*/
package org.apache.plc4x.java.examples.kafkabridge;

import java.io.File;
import java.io.IOException;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.TimeUnit;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import org.apache.commons.lang3.StringUtils;
import org.apache.edgent.connectors.kafka.KafkaProducer;
import org.apache.edgent.function.Supplier;
import org.apache.edgent.providers.direct.DirectProvider;
import org.apache.edgent.topology.TStream;
import org.apache.edgent.topology.Topology;
import org.apache.plc4x.edgent.PlcConnectionAdapter;
import org.apache.plc4x.edgent.PlcFunctions;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.examples.kafkabridge.model.Address;
import org.apache.plc4x.java.examples.kafkabridge.model.Configuration;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class KafkaBridge {

    private static final Logger logger = LoggerFactory.getLogger(KafkaBridge.class);

    private Configuration config;
    private PlcConnectionAdapter plcAdapter;

    private KafkaBridge(String propsPath) {
        if(StringUtils.isEmpty(propsPath)) {
            logger.error("Empty configuration file parameter");
            throw new IllegalArgumentException("Empty configuration file parameter");
        }
        File propsFile = new File(propsPath);
        if(!(propsFile.exists() && propsFile.isFile())) {
            logger.error("Invalid configuration file {}", propsFile.getPath());
            throw new IllegalArgumentException("Invalid configuration file " + propsFile.getPath());
        }
        ObjectMapper mapper = new ObjectMapper(new YAMLFactory());
        try {
            config = mapper.readValue(propsFile, Configuration.class);
            plcAdapter = new PlcConnectionAdapter(config.getPlcConfig().getConnection());
        } catch (IOException e) {
            logger.error("Error parsing configuration", e);
        }
    }

    private void run() {
        DirectProvider dp = new DirectProvider();
        Topology top = dp.newTopology("kafka-bridge");

        // Create the Kafka Producer broker connector
        Map<String,Object> kafkaConfig = createKafkaConfig();
        KafkaProducer kafka = new KafkaProducer(top, () -> kafkaConfig);

        Map<String, ReadRequestItem> readRequestItems = new HashMap<>();

        for(Address address : config.getPlcConfig().getAddresses()) {
            try {
                org.apache.plc4x.java.api.model.Address plcAddress = plcAdapter.parseAddress(address.getAddress());
                ReadRequestItem readItem = new ReadRequestItem<>(address.getType(), plcAddress, address.getSize());
                readRequestItems.put(address.getName(), readItem);
            } catch (PlcException e) {
                logger.error("Error parsing address {}", address.getAddress(), e);
            }
        }
        // TODO: Here we somehow have to create an Edgent supplier, that can cope with batch reads ...
        // WARN: This example doesn't work at the moment ...
        Supplier<Byte> plcSupplier = PlcFunctions.byteSupplier(plcAdapter, "INPUTS/0");

        // Start polling our plc source in the given interval.
        TStream<Byte> source = top.poll(plcSupplier, config.getPollingInterval(), TimeUnit.MILLISECONDS);

        // Convert the byte into a string.
        TStream<String> stringSource = source.map(value -> Byte.toString(value));

        // Publish the stream to the topic.  The String tuple is the message value.
        kafka.publish(stringSource, config.getKafkaConfig().getTopicName());

        dp.submit(top);
    }

    private Map<String,Object> createKafkaConfig() {
        Map<String,Object> kafkaConfig = new HashMap<>();
        kafkaConfig.put("bootstrap.servers", config.getKafkaConfig().getBootstrapServers());
        return kafkaConfig;
    }

    public static void main(String[] args) {
        if(args.length != 1) {
            System.out.println("Usage: KafkaBridge {path-to-kafkabridge-yml}");
        }
        KafkaBridge kafkaBridge = new KafkaBridge(args[0]);
        kafkaBridge.run();
    }

}
