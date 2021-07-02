/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.spi.messages;

import com.fasterxml.jackson.annotation.*;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.spi.utils.XmlSerializable;
import org.w3c.dom.Document;
import org.w3c.dom.Element;

import java.util.concurrent.CompletableFuture;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcDiscoveryRequest implements PlcDiscoveryRequest, XmlSerializable {

    private final PlcDiscoverer discoverer;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcDiscoveryRequest(@JsonProperty("discoverer") PlcDiscoverer discoverer) {
        this.discoverer = discoverer;
    }

    @Override
    public CompletableFuture<? extends PlcDiscoveryResponse> execute() {
        return discoverer.discover(this);
    }

    @Override
    public void xmlSerialize(Element parent) {
        Document doc = parent.getOwnerDocument();
        Element messageElement = doc.createElement("PlcDiscoveryRequest");
        parent.appendChild(messageElement);
        // TODO: Implement
    }

    public static class Builder implements PlcDiscoveryRequest.Builder {

        private final PlcDiscoverer discoverer;

        public Builder(PlcDiscoverer discoverer) {
            this.discoverer = discoverer;
        }

        @Override
        public PlcDiscoveryRequest build() {
            return new DefaultPlcDiscoveryRequest(discoverer);
        }

    }

}
