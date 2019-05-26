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

package org.apache.plc4x.java.examples.storage.elasticsearch;

import org.apache.edgent.function.Supplier;
import org.apache.edgent.providers.direct.DirectProvider;
import org.apache.edgent.topology.TStream;
import org.apache.edgent.topology.Topology;
import org.apache.plc4x.edgent.PlcConnectionAdapter;
import org.apache.plc4x.edgent.PlcFunctions;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.elasticsearch.action.admin.indices.create.CreateIndexRequest;
import org.elasticsearch.action.admin.indices.create.CreateIndexResponse;
import org.elasticsearch.action.admin.indices.exists.indices.IndicesExistsRequest;
import org.elasticsearch.action.index.IndexResponse;
import org.elasticsearch.client.Client;
import org.elasticsearch.client.IndicesAdminClient;
import org.elasticsearch.common.settings.Settings;
import org.elasticsearch.common.xcontent.XContentBuilder;
import org.elasticsearch.common.xcontent.XContentFactory;
import org.elasticsearch.common.xcontent.XContentType;
import org.elasticsearch.node.InternalSettingsPreparer;
import org.elasticsearch.node.Node;
import org.elasticsearch.node.NodeValidationException;
import org.elasticsearch.plugins.Plugin;
import org.elasticsearch.transport.Netty4Plugin;

import java.io.IOException;
import java.util.Calendar;
import java.util.Collection;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.TimeUnit;

public class ElasticsearchStorage {

    public enum ConveyorState {
        STOPPED,
        RUNNING_LEFT,
        RUNNING_RIGHT
    }

    private static class IotElasticsearchFactoryException extends RuntimeException {
        private IotElasticsearchFactoryException(String message) {
            super(message);
        }
        private IotElasticsearchFactoryException(String message, Throwable cause) {
            super(message, cause);
        }
    }

    private static class MyNode extends Node {
        private MyNode(Settings preparedSettings, Collection<Class<? extends Plugin>> classpathPlugins) {
            super(InternalSettingsPreparer.prepareEnvironment(
                preparedSettings, Collections.emptyMap(), null, () -> "hello-es"), classpathPlugins, true);
        }
    }

    private ConveyorState conveyorState = ConveyorState.STOPPED;

    private Node startElasticsearchNode() {
        try {
            Node node = new MyNode(Settings.builder()
                .put("transport.type", "netty4")
                .put("http.type", "netty4")
                .put("http.cors.enabled", "true")
                .put("path.home", "elasticsearch-data")
                .build(), Collections.singletonList(Netty4Plugin.class));
            node.start();
            return node;
        } catch (NodeValidationException e) {
            throw new IotElasticsearchFactoryException("Could not start Elasticsearch node.", e);
        }
    }

    private void prepareIndexes(Client esClient) {
        IndicesAdminClient indicesAdminClient = esClient.admin().indices();

        // Check if the factory-data index exists and create it, if it doesn't.
        IndicesExistsRequest factoryDataIndexExistsRequest =
            indicesAdminClient.prepareExists("iot-factory-data").request();
        if(!indicesAdminClient.exists(factoryDataIndexExistsRequest).actionGet().isExists()) {
            CreateIndexRequest createIndexRequest = new CreateIndexRequest("iot-factory-data");
            createIndexRequest.mapping("FactoryData",
                "{\n" +
                    "            \"properties\": {\n" +
                    "                \"time\": {\n" +
                    "                    \"type\": \"date\"\n" +
                    "                }\n" +
                    "            }\n" +
                    "        }", XContentType.JSON);
            CreateIndexResponse createIndexResponse = indicesAdminClient.create(createIndexRequest).actionGet();
            if(!createIndexResponse.isAcknowledged()) {
                throw new IotElasticsearchFactoryException("Could not create index 'iot-factory-data'");
            }
        }

        // Check if the product-data index exists and create it, if it doesn't.
        IndicesExistsRequest productDataIndexExistsRequest =
            indicesAdminClient.prepareExists("iot-product-data").request();
        if(!indicesAdminClient.exists(productDataIndexExistsRequest).actionGet().isExists()) {
            CreateIndexRequest createIndexRequest = new CreateIndexRequest("iot-product-data");
            createIndexRequest.mapping("ProductData",
                "{\n" +
                    "            \"properties\": {\n" +
                    "                \"time\": {\n" +
                    "                    \"type\": \"date\"\n" +
                    "                },\n" +
                    "                \"type\": {\n" +
                    "                    \"type\": \"keyword\"\n" +
                    "                }\n" +
                    "            }\n" +
                    "        }", XContentType.JSON);
            CreateIndexResponse createIndexResponse = indicesAdminClient.create(createIndexRequest).actionGet();
            if(!createIndexResponse.isAcknowledged()) {
                throw new IotElasticsearchFactoryException("Could not create index 'iot-product-data'");
            }
        }
    }

    private void runFactory() {
        // Start an Elasticsearch node.
        Node esNode = startElasticsearchNode();
        Client esClient = esNode.client();
        // Register a shutdown hook.
        final Thread mainThread = Thread.currentThread();
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            try {
                esNode.close();
                mainThread.join();
            } catch (IOException | InterruptedException e) {
                throw new PlcRuntimeException("Error closing ES Node", e);
            }
        }));
        System.out.println("Started Elasticsearch node on port 9200");

        // Make sure the indexes exist prior to writing to them.
        prepareIndexes(esClient);

        // Get a plc connection.
        try (PlcConnectionAdapter plcAdapter = new PlcConnectionAdapter("s7://10.10.64.20/1/1")) {
            // Initialize the Edgent core.
            DirectProvider dp = new DirectProvider();
            Topology top = dp.newTopology();

            // Define the event stream.
            // 1) PLC4X source generating a stream of bytes.
            Supplier<List<Boolean>> plcSupplier = PlcFunctions.booleanListSupplier(plcAdapter, "%Q0:BYTE");
            // 2) Use polling to get an item from the byte-stream in regular intervals.
            TStream<List<Boolean>> plcOutputStates = top.poll(plcSupplier, 100, TimeUnit.MILLISECONDS);

            // 3a) Create a stream that pumps all data into a 'factory-data' index.
            TStream<XContentBuilder> factoryData = plcOutputStates.map(this::translatePlcInput);
            TStream<IndexResponse> factoryDataResponses = factoryData.map(
                value -> esClient.prepareIndex("iot-factory-data", "FactoryData").setSource(value).get());
            factoryDataResponses.print();

            // 3b) Create a stream that does some local analysis to detect big and small boxes and to only output
            //     something to the 'product-data' index, if a new item is detected.
            TStream<XContentBuilder> productData = plcOutputStates.map(this::handlePlcInput);
            TStream<IndexResponse> productDataResponses = productData.map(
                value -> esClient.prepareIndex("iot-product-data", "ProductData").setSource(value).get());
            productDataResponses.print();

            // Submit the topology and hereby start the event streams.
            dp.submit(top);
        } catch (Exception e) {
            throw new IotElasticsearchFactoryException("Error while connecting or disconnecting from the PLC.", e);
        }
    }

    private XContentBuilder translatePlcInput(List<Boolean> input) {
        boolean conveyorEntry = input.get(0);
        boolean load = input.get(1);
        boolean unload = input.get(2);
        boolean transferLeft = input.get(3);
        boolean transferRight = input.get(4);
        boolean conveyorLeft = input.get(5);
        boolean conveyorRight = input.get(6);

        try(XContentBuilder builder = XContentFactory.jsonBuilder()
            .startObject()
            .field("time", Calendar.getInstance().getTimeInMillis())
            .field("conveyorEntry", conveyorEntry)
            .field("load", load)
            .field( "unload", unload)
            .field( "transferLeft", transferLeft)
            .field( "transferRight", transferRight)
            .field( "conveyorLeft", conveyorLeft)
            .field( "conveyorRight", conveyorRight)
            .endObject()) {
            return builder;
        } catch (IOException e) {
            throw new IotElasticsearchFactoryException("Error building JSON message.", e);
        }
    }

    private XContentBuilder handlePlcInput(List<Boolean> input) {
        boolean transferLeft = input.get(3);
        boolean transferRight = input.get(4);

        if (conveyorState == ConveyorState.STOPPED) {
            if (transferLeft) {
                conveyorState = ConveyorState.RUNNING_LEFT;
                try (XContentBuilder builder = XContentFactory.jsonBuilder()
                    .startObject()
                    .field("time", Calendar.getInstance().getTimeInMillis())
                    .field("type", "small")
                    .endObject()) {
                    return builder;
                } catch (IOException e) {
                    throw new IotElasticsearchFactoryException("Error building JSON message.", e);
                }
            } else if (transferRight){
                conveyorState = ConveyorState.RUNNING_RIGHT;
                try (XContentBuilder builder = XContentFactory.jsonBuilder()
                    .startObject()
                    .field("time", Calendar.getInstance().getTimeInMillis())
                    .field("type", "large")
                    .endObject()) {
                    return builder;
                } catch (IOException e) {
                    throw new IotElasticsearchFactoryException("Error building JSON message.", e);
                }
            }
        } else if (!(transferLeft || transferRight)) {
            conveyorState = ConveyorState.STOPPED;
        }
        return null;
    }

    public static void main(String[] args) {
        ElasticsearchStorage factory = new ElasticsearchStorage();
        factory.runFactory();
    }

}
