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
package org.apache.plc4x.java.cbus.protocol;

import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.cbus.readwrite.CBusCommand;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;

public class CBusProtocolLogic extends Plc4xProtocolBase<CBusCommand> {

    private final Logger logger = LoggerFactory.getLogger(CBusProtocolLogic.class);

    private RequestTransactionManager tm;

    @Override
    public void setDriverContext(DriverContext driverContext) {
        super.setDriverContext(driverContext);

        // Initialize Transaction Manager.
        // Until the number of concurrent requests is successfully negotiated we set it to a
        // maximum of only one request being able to be sent at a time. During the login process
        // No concurrent requests can be sent anyway. It will be updated when receiving the
        // S7ParameterSetupCommunication response.
        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public void onConnect(ConversationContext<CBusCommand> context) {

    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        return future;
    }

    /**
     * This method is only called when there is no Response Handler.
     */
    @Override
    protected void decode(ConversationContext<CBusCommand> context, CBusCommand msg) throws Exception {
    }

    @Override
    public void close(ConversationContext<CBusCommand> context) {
    }

}
