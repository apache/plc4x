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

package org.apache.plc4x.sandbox.java.dynamic.actions;

import org.apache.commons.scxml2.ActionExecutionContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.Socket;

public class ConnectAction extends BasePlc4xAction {

    private String type;
    private String host;
    private String port;

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getHost() {
        return host;
    }

    public void setHost(String host) {
        this.host = host;
    }

    public String getPort() {
        return port;
    }

    public void setPort(String port) {
        this.port = port;
    }

    @Override
    protected Logger getLogger() {
        return LoggerFactory.getLogger(ConnectAction.class);
    }

    @Override
    public void execute(ActionExecutionContext ctx) {
        getLogger().info(getStateName() + ": Connecting...");
        try {
            if ("TCP".equalsIgnoreCase(type)) {
                Socket socket = new Socket(host, Integer.parseInt(port));
                ctx.getGlobalContext().set(BaseConnectedAction.SOCKET_PARAMETER_NAME, socket);

                getLogger().info("Connected.");

                fireSuccessEvent(ctx);
            }
        } catch (IOException e) {
            getLogger().error("Error connecting to remote.", e);
        }
    }

}
