/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.ctrlx.readwrite.connection;

import com.hrakaroo.glob.MatchingEngine;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.model.PlcQuery;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.ctrlx.readwrite.rest.datalayer.ApiClient;
import org.apache.plc4x.java.ctrlx.readwrite.rest.datalayer.ApiException;
import org.apache.plc4x.java.ctrlx.readwrite.rest.datalayer.api.DataLayerInformationAndSettingsApi;
import org.apache.plc4x.java.ctrlx.readwrite.rest.datalayer.api.NodesApi;
import org.apache.plc4x.java.ctrlx.readwrite.rest.datalayer.model.BrowseData;
import org.apache.plc4x.java.ctrlx.readwrite.rest.datalayer.model.ReadNode200Response;
import org.apache.plc4x.java.ctrlx.readwrite.tag.CtrlXQuery;
import org.apache.plc4x.java.ctrlx.readwrite.tag.CtrlXTag;
import org.apache.plc4x.java.ctrlx.readwrite.tag.CtrlXTagHandler;
import org.apache.plc4x.java.ctrlx.readwrite.utils.ApiClientFactory;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.stream.Collectors;

public class CtrlXConnection implements PlcConnection, PlcPinger, PlcBrowser {

    private static final Logger logger = LoggerFactory.getLogger(CtrlXConnection.class);

    private final String baseUrl;
    private final String username;
    private final String password;

    private final ExecutorService executorService;
    private PlcValueHandler valueHandler;

    private ApiClient apiClient;
    private NodesApi nodesApi;
    private DataLayerInformationAndSettingsApi dataLayerApi;

    private final CtrlXTagHandler controlXTagHandler = new CtrlXTagHandler();

    public CtrlXConnection(String baseUrl, String username, String password) {
        this.baseUrl = baseUrl;
        this.username = username;
        this.password = password;
        this.executorService = Executors.newFixedThreadPool(10);
        this.valueHandler = new DefaultPlcValueHandler();
    }

    @Override
    public Optional<PlcValue> parseTagValue(PlcTag tag, Object... values) {
        PlcValue plcValue;
        try {
            plcValue = valueHandler.newPlcValue(tag, values);
        } catch (Exception e) {
            throw new PlcRuntimeException("Error parsing tag value " + tag, e);
        }
        return Optional.of(plcValue);
    }

    @Override
    public void connect() throws PlcConnectionException {
        if (apiClient != null) {
            throw new PlcConnectionException("Already connected");
        }
        apiClient = ApiClientFactory.getApiClient(baseUrl, username, password);
        nodesApi = new NodesApi(apiClient);
        dataLayerApi = new DataLayerInformationAndSettingsApi(apiClient);
    }

    @Override
    public boolean isConnected() {
        return apiClient != null;
    }

    @Override
    public void close() throws Exception {
        apiClient.getHttpClient().close();
        apiClient = null;
        dataLayerApi = null;
        executorService.shutdown();
    }

    @Override
    public Optional<PlcTag> parseTagAddress(String tagAddress) {
        PlcTag plcTag;
        try {
            plcTag = controlXTagHandler.parseTag(tagAddress);
        } catch (Exception e) {
            logger.error("Error parsing tag address {}", tagAddress);
            return Optional.empty();
        }
        return Optional.ofNullable(plcTag);
    }

    @Override
    public PlcConnectionMetadata getMetadata() {
        return new PlcConnectionMetadata() {
            @Override
            public boolean isReadSupported() {
                return true;
            }

            @Override
            public boolean isWriteSupported() {
                return true;
            }

            @Override
            public boolean isSubscribeSupported() {
                return true;
            }

            @Override
            public boolean isBrowseSupported() {
                return true;
            }
        };
    }

    @Override
    public CompletableFuture<? extends PlcPingResponse> ping() {
        return new DefaultPlcPingRequest(this).execute();
    }

    @Override
    public PlcBrowseRequest.Builder browseRequestBuilder() {
        return new DefaultPlcBrowseRequest.Builder(this, controlXTagHandler);
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        return null;
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        return null;
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        return null;
    }

    @Override
    public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
        return null;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // PlcPinger
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public CompletableFuture<PlcPingResponse> ping(PlcPingRequest pingRequest) {
        CompletableFuture<PlcPingResponse> future = new CompletableFuture<>();
        executorService.execute(() -> {
            try {
                // Just execute some random request, that we can expect to be replied by the controller.
                nodesApi.readNode("datalayer/server/settings", "browse");
                future.complete(new DefaultPlcPingResponse(pingRequest, PlcResponseCode.OK));
            } catch (ApiException e) {
                future.completeExceptionally(new PlcProtocolException("Error pinging remote", e));
            }
        });
        return future;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // PlcBrowser
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public CompletableFuture<PlcBrowseResponse> browse(PlcBrowseRequest browseRequest) {
        return browseWithInterceptor(browseRequest, item -> true);
    }

    @Override
    public CompletableFuture<PlcBrowseResponse> browseWithInterceptor(PlcBrowseRequest browseRequest, PlcBrowseRequestInterceptor interceptor) {
        CompletableFuture<PlcBrowseResponse> future = new CompletableFuture<>();
        executorService.execute(() -> {
            int numQueries = browseRequest.getQueryNames().size();

            // Initialize the response structures.
            Map<String, PlcResponseCode> responseCodes = new HashMap<>(numQueries);
            Map<String, List<PlcBrowseItem>> responseItems = new HashMap<>(numQueries);
            Map<String, MatchingEngine> matchers = new HashMap<>(numQueries);
            for (String queryName : browseRequest.getQueryNames()) {
                responseCodes.put(queryName, PlcResponseCode.OK);
                responseItems.put(queryName, new ArrayList<>());
                PlcQuery query = browseRequest.getQuery(queryName);
                if (query instanceof CtrlXQuery) {
                    CtrlXQuery ctrlXQuery = (CtrlXQuery) query;
                    matchers.put(queryName, ctrlXQuery.getMatcher());
                } else {
                    future.completeExceptionally(
                        new PlcInvalidTagException("Invalid query type: " + query.getClass().getName()));
                }
            }

            // Now walk through the tree and for each leaf-node, check which queries it matches.
            // Start by initializing the list with the lists of all normal and real-time nodes.
            Queue<String> uncheckedNodeList = new LinkedList<>();
            try {
                // Initialize the list with all normal node names
                BrowseData nodeNames = dataLayerApi.getNodeNames();
                if (nodeNames.getValue() != null) {
                    uncheckedNodeList.addAll(nodeNames.getValue());
                }

                // Then add all real-time node names.
                BrowseData realtimeNodeNames = dataLayerApi.getRealtimeNodeNames();
                if (realtimeNodeNames.getValue() != null) {
                    uncheckedNodeList.addAll(realtimeNodeNames.getValue());
                }
            } catch (ApiException e) {
                throw new RuntimeException(e);
            }
            // Now keep on resolving paths till the list is empty.
            while (!uncheckedNodeList.isEmpty()) {
                String curNode = uncheckedNodeList.poll();
                try {
                    ReadNode200Response readNode200Response = nodesApi.readNode(curNode, "browse");
                    List<String> children = readNode200Response.getValue();

                    // If this node has no children, then this is a leaf-node,
                    // and it's a potential match for any of the queries.
                    if (children.isEmpty()) {
                        List<String> matchingQueryNames = matchers.entrySet().stream()
                            .filter(entry -> entry.getValue().matches(curNode)).map(Map.Entry::getKey)
                            .collect(Collectors.toList());
                        // If there's at least one matching query, read the "metadata", which contains information
                        // on if the property is readable or writable.
                        if (!matchingQueryNames.isEmpty()) {
                            // TODO: Implement the reading of "metadate" as this contains information on if the
                            // tag is readable or writable.
                            /*try {
                                ReadNode200Response metaDataResponse = nodesApi.readNode(curNode, "metadata");
                                System.out.println(metaDataResponse);
                            } catch (ApiException e) {
                                e.printStackTrace();
                            }*/
                            matchingQueryNames.forEach(queryName -> responseItems.get(queryName).add(
                                new DefaultPlcBrowseItem(
                                    new CtrlXTag(curNode, PlcValueType.BOOL, Collections.emptyList()),
                                    curNode, true, true, true, false,
                                    Collections.emptyList(), Collections.emptyMap(), Collections.emptyMap())));
                        }
                    }
                    // If this node has children, then it's branch, and we need to add its children to the queue.
                    else {
                        // Add all children to the list.
                        uncheckedNodeList.addAll(children.stream().map(child -> curNode + "/" + child).collect(Collectors.toList()));
                    }
                } catch (ApiException e) {
                    // Ignore ...
                }
            }
            future.complete(new DefaultPlcBrowseResponse(browseRequest, responseCodes, responseItems));
        });
        return future;
    }

}
