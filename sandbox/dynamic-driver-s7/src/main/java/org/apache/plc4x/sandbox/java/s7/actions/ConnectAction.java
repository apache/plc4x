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

package org.apache.plc4x.sandbox.java.s7.actions;

import org.apache.commons.scxml2.ActionExecutionContext;
import org.apache.commons.scxml2.EventBuilder;
import org.apache.commons.scxml2.TriggerEvent;
import org.apache.commons.scxml2.model.Action;

import java.io.IOException;
import java.net.Socket;

public class ConnectAction extends Action {

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
    public void execute(ActionExecutionContext ctx) {
        ctx.getAppLog().info("Connecting...");
        try {
            if ("TCP".equalsIgnoreCase(type)) {
                Socket socket = new Socket(host, Integer.parseInt(port));
                ctx.getGlobalContext().set("connection", socket);
                TriggerEvent event = new EventBuilder("success", TriggerEvent.SIGNAL_EVENT).build();
                ctx.getInternalIOProcessor().addEvent(event);
                ctx.getAppLog().info("Connected.");
            }
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

}
