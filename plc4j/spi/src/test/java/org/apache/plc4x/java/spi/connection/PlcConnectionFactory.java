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

package org.apache.plc4x.java.spi.connection;

import org.apache.plc4x.java.spi.generation.Message;

public class PlcConnectionFactory {

    private boolean awaitDiscovery = false;
    private boolean fireDiscovery = false;
    private boolean awaitDisconnect;

    PlcConnectionFactory withDiscovery() {
        awaitDiscovery = true;
        fireDiscovery = true;
        return this;
    }

    PlcConnectionFactory doNotAwaitForDisconnect() {
        awaitDisconnect = false;
        return this;
    }

    <T extends Message> DefaultNettyPlcConnection create(ChannelFactory channelFactory, ProtocolStackConfigurer<T> stackConfigurer) {
        return new DefaultNettyPlcConnection(
            true, true, true, true, true,
            null, null, channelFactory,
            fireDiscovery, // force discovery
            true, // await setup
            awaitDisconnect, // await disconnect
            awaitDiscovery, // await discovery
            stackConfigurer,
            null,
            null
        );
    }

}
