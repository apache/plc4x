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
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
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
import org.apache.plc4x.java.ctrlx.readwrite.configuration.CtrlXConfiguration;
import org.apache.plc4x.java.ctrlx.readwrite.utils.ApiClientFactory;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcBrowseItem;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcBrowseRequest;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcBrowseResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcPingRequest;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcPingResponse;
import org.apache.plc4x.java.spi.drivers.functions.PlcBrowser;
import org.apache.plc4x.java.spi.drivers.functions.PlcPinger;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.stream.Collectors;

public class CtrlXConnection implements PlcConnection, PlcPinger, PlcBrowser {

    private static final Logger LOGGER = LoggerFactory.getLogger(CtrlXConnection.class);

    private static final Logger logger = LoggerFactory.getLogger(CtrlXConnection.class);

    private final String baseUrl;
    private final CtrlXConfiguration configuration;
    private final String username;
    private final String password;

    private final ExecutorService executorService;
    private PlcValueHandler valueHandler;

    private ApiClient apiClient;
    private NodesApi nodesApi;
    private DataLayerInformationAndSettingsApi dataLayerApi;

    private final CtrlXTagHandler controlXTagHandler = new CtrlXTagHandler();

    public CtrlXConnection(String baseUrl, String username, String password,
                           CtrlXConfiguration configuration) {
        this.baseUrl = baseUrl;
        this.username = username;
        this.password = password;
        this.configuration = configuration;
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
        apiClient = ApiClientFactory.getApiClient(baseUrl, username, password, configuration);
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

    /**
     * Reading, writing and subscribing are not implemented - the corresponding request builders
     * return {@code null}. Browsing has an implementation, but it does not work yet. Reporting any
     * of them as supported makes callers that check the metadata first (the connection cache and
     * the tooling among them) attempt an operation that cannot succeed, so all four are reported
     * as unsupported until they actually work.
     */
    @Override
    public PlcConnectionMetadata getMetadata() {
        return new PlcConnectionMetadata() {
            @Override
            public boolean isReadSupported() {
                return false;
            }

            @Override
            public boolean isWriteSupported() {
                return false;
            }

            @Override
            public boolean isSubscribeSupported() {
                return false;
            }

            @Override
            public boolean isBrowseSupported() {
                return false;
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
        return browseWithInterceptor(browseRequest, (queryName,query, item) -> true);
    }

    @Override
    public CompletableFuture<PlcBrowseResponse> browseWithInterceptor(PlcBrowseRequest browseRequest, PlcBrowseRequestInterceptor interceptor) {
        CompletableFuture<PlcBrowseResponse> future = new CompletableFuture<>();
        executorService.execute(() -> {
            // Whatever happens in here, the caller is waiting on this future and has no other way
            // to learn that it went wrong. Anything thrown out of a Runnable handed to an executor
            // goes to the thread's uncaught handler and leaves the future waiting for good, so a
            // failure has to be turned into a completion rather than allowed to escape.
            try {
                doBrowse(browseRequest, interceptor, future);
            } catch (Throwable e) {
                future.completeExceptionally(
                    new PlcProtocolException("Error browsing remote", e));
            }
        });
        return future;
    }

    private void doBrowse(PlcBrowseRequest browseRequest, PlcBrowseRequestInterceptor interceptor,
                          CompletableFuture<PlcBrowseResponse> future) {
            int numQueries = browseRequest.getQueryNames().size();

            // Initialize the response structures.
            Map<String, PlcResponseCode> responseCodes = new HashMap<>(numQueries);
            Map<String, List<PlcBrowseItem>> responseItems = new HashMap<>(numQueries);
            Map<String, MatchingEngine> matchers = new HashMap<>(numQueries);
            for (String queryName : browseRequest.getQueryNames()) {
                responseCodes.put(queryName, PlcResponseCode.OK);
                responseItems.put(queryName, new ArrayList<>());
                PlcQuery query = browseRequest.getQuery(queryName);
                if (query instanceof CtrlXQuery ctrlXQuery) {
                    matchers.put(queryName, ctrlXQuery.getMatcher());
                } else {
                    future.completeExceptionally(
                        new PlcInvalidTagException("Invalid query type: " + query.getClass().getName()));
                }
            }

            // Now walk through the tree and for each leaf-node, check which queries it matches.
            // Start by initializing the list with the lists of all normal and real-time nodes.
            List<String> roots = new ArrayList<>();
            try {
                // Initialize the list with all normal node names
                BrowseData nodeNames = dataLayerApi.getNodeNames();
                if (nodeNames.getValue() != null) {
                    roots.addAll(nodeNames.getValue());
                }

                // Then add all real-time node names.
                BrowseData realtimeNodeNames = dataLayerApi.getRealtimeNodeNames();
                if (realtimeNodeNames.getValue() != null) {
                    roots.addAll(realtimeNodeNames.getValue());
                }
            } catch (ApiException e) {
                future.completeExceptionally(new PlcProtocolException("Error listing root nodes", e));
                return;
            }
            // Now keep on resolving paths till the list is empty.
            walkNodes(roots,
                curNode -> nodesApi.readNode(curNode, "browse").getValue(),
                configuration.getBrowseMaxTotalNodes(),
                configuration.getBrowseMaxDepth(),
                curNode -> {
                    List<String> matchingQueryNames = matchers.entrySet().stream()
                        .filter(entry -> entry.getValue().matches(curNode)).map(Map.Entry::getKey)
                        .toList();
                    // If there's at least one matching query, read the "metadata", which contains information
                    // on if the property is readable or writable.
                    if (!matchingQueryNames.isEmpty()) {
                        // TODO: Implement the reading of "metadate" as this contains information on if the
                        // tag is readable or writable.
                        matchingQueryNames.forEach(queryName -> responseItems.get(queryName).add(
                            new DefaultPlcBrowseItem(
                                new CtrlXTag(curNode, PlcValueType.BOOL, Collections.emptyList()),
                                curNode, true, true,
                                EnumSet.of(PlcSubscriptionType.CYCLIC, PlcSubscriptionType.CHANGE_OF_STATE, PlcSubscriptionType.EVENT),
                                false,
                                Collections.emptyList(), Collections.emptyMap(), Collections.emptyMap())));
                    }
                });
            future.complete(new DefaultPlcBrowseResponse(browseRequest, responseCodes, responseItems));
    }

    /**
     * Reads the names of one node's children. The browse's only contact with the device, kept
     * separate so the walk below can be exercised without one.
     */
    @FunctionalInterface
    interface NodeBrowser {
        List<String> childrenOf(String node) throws ApiException;
    }

    /** A node still to be read, and how far below a root it sits. */
    private record PendingNode(String path, int depth) {
    }

    /**
     * Walks the tree the device describes, handing every leaf to {@code onLeaf}.
     *
     * <p>How much work this is belongs to the device, not to the caller: each node is read to find
     * its children and each child is then read the same way. Two things bound it. A child is
     * addressed by its parent's path with the child's name appended, so a device naming a child
     * that leads back to somewhere the walk has already been produces a longer path every time and
     * never repeats one - remembering visited paths would not notice, and only depth ends it. And
     * the total is capped, for a tree that is finite but larger than anyone wants to read.</p>
     *
     * <p>A node that cannot be read is skipped, as it always was; one node refusing to answer is
     * not a reason to abandon the rest of the tree.</p>
     *
     * @return how many nodes were read
     */
    static int walkNodes(List<String> roots, NodeBrowser browser, int maxTotalNodes, int maxDepth,
                         Consumer<String> onLeaf) {
        Queue<PendingNode> pending = new LinkedList<>();
        for (String root : roots) {
            pending.add(new PendingNode(root, 0));
        }
        int read = 0;
        while (!pending.isEmpty()) {
            if (maxTotalNodes > 0 && read >= maxTotalNodes) {
                LOGGER.warn("Stopped browsing after {} nodes, the configured browse-max-total-nodes; "
                    + "{} were still queued", maxTotalNodes, pending.size());
                break;
            }
            PendingNode current = pending.poll();
            read++;
            List<String> children;
            try {
                children = browser.childrenOf(current.path());
            } catch (ApiException e) {
                // Ignore ...
                continue;
            }
            // A node naming no children is a leaf, and a potential match for any of the queries.
            // The device may say so with an empty list or with none at all; both mean the same
            // thing here, and the second used to be read as a list and fail.
            if (children == null || children.isEmpty()) {
                onLeaf.accept(current.path());
                continue;
            }
            if (maxDepth > 0 && current.depth() >= maxDepth) {
                LOGGER.warn("Stopped browsing '{}' at depth {}, the configured browse-max-depth",
                    current.path(), maxDepth);
                continue;
            }
            for (String child : children) {
                pending.add(new PendingNode(current.path() + "/" + child, current.depth() + 1));
            }
        }
        return read;
    }

}
