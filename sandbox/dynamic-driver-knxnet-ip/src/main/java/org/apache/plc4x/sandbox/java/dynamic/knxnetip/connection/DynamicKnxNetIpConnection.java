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

package org.apache.plc4x.sandbox.java.dynamic.knxnetip.connection;

import org.apache.commons.lang3.StringUtils;
import org.apache.commons.scxml2.EventBuilder;
import org.apache.commons.scxml2.TriggerEvent;
import org.apache.commons.scxml2.model.ModelException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.base.messages.InternalPlcReadRequest;
import org.apache.plc4x.java.base.messages.InternalPlcReadResponse;
import org.apache.plc4x.java.base.messages.PlcReader;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.sandbox.java.dynamic.connection.DynamicDriverConnectionBase;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;

public class DynamicKnxNetIpConnection extends DynamicDriverConnectionBase implements PlcReader {

    private static final Logger logger = LoggerFactory.getLogger(DynamicKnxNetIpConnection.class);

    private final InetAddress address;

    public DynamicKnxNetIpConnection(InetAddress address, String params) {
        super("org/apache/plc4x/protocols/knxnetip/protocol.scxml.xml",
            "org/apache/plc4x/protocols/knxnetip/protocol.dfdl.xsd");

        this.address = address;

        if (!StringUtils.isEmpty(params)) {
            for (String param : params.split("&")) {
                String[] paramElements = param.split("=");
                String paramName = paramElements[0];
                if (paramElements.length == 2) {
                    String paramValue = paramElements[1];
                    switch (paramName) {
                        default:
                            logger.debug("Unknown parameter {} with value {}", paramName, paramValue);
                    }
                } else {
                    logger.debug("Unknown no-value parameter {}", paramName);
                }
            }
        }
    }

    @Override
    protected String getConnectedStateName() {
        return "connected";
    }

    @Override
    protected String getDisconnectTransitionName() {
        return "disconnect";
    }

    /*@Override
    protected Collection<CustomAction> getAdditionalCustomActions() {
        return Arrays.asList(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "S7DecodeArticleNumber",
                S7DecodeArticleNumber.class),
            new CustomAction("https://plc4x.apache.org/scxml-extension", "S7DecodeReadResponse",
                S7DecodeReadResponseAction.class),
            new CustomAction("https://plc4x.apache.org/scxml-extension", "S7DecodeWriteResponse",
                S7DecodeWriteResponseAction.class)
            );
    }*/

    @Override
    protected Map<String, Object> getAdditionalContextDataItems() {
        Map<String, Object> dataItems = new HashMap<>();

    /*<sc:data id="clientIpAddress"/>
    <sc:data id="clientUdpPort"/>*/

        dataItems.put("serverIpAddress", address.getHostAddress());
        dataItems.put("serverPort", "102");

        return dataItems;
    }

    @Override
    public PlcConnectionMetadata getMetadata() {
        return super.getMetadata();
    }

    /*@Override
    public PlcReadRequest.Builder readRequestBuilder() {
        return new DefaultPlcReadRequest.Builder(this, new S7PlcFieldHandler());
    }*/

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        InternalPlcReadRequest internalReadRequest = checkInternal(readRequest, InternalPlcReadRequest.class);
        CompletableFuture<InternalPlcReadResponse> future = new CompletableFuture<>();
        PlcRequestContainer<InternalPlcReadRequest, InternalPlcReadResponse> container =
            new PlcRequestContainer<>(internalReadRequest, future);

        try {
            getExecutor().triggerEvent(
                new EventBuilder("read", TriggerEvent.CALL_EVENT).data(container).build());
        } catch (ModelException e) {
            throw new PlcRuntimeException("Error reading.", e);
        }

        return future.thenApply(PlcReadResponse.class::cast);
    }

}
