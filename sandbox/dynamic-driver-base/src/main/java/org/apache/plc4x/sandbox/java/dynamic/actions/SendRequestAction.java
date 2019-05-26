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
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.sandbox.java.dynamic.utils.RequestRegistry;
import org.jdom2.Document;

public class SendRequestAction extends SendAction {

    private String idExpression = null;

    public String getIdExpression() {
        return idExpression;
    }

    public void setIdExpression(String idExpression) {
        this.idExpression = idExpression;
    }

    @Override
    protected void processMessage(Document message, ActionExecutionContext ctx) {
        // Generate a new request and add that to the context.
        RequestRegistry requestRegistry = (RequestRegistry) ctx.getGlobalContext().get("requestRegistry");
        String requestId = requestRegistry.generateRequestId();
        ctx.getGlobalContext().set("requestId", requestId);

        // If a container is present, we want to respond to a call from outside the driver.
        // Register the container with the registry using the requestId as key.
        PlcRequestContainer container = (PlcRequestContainer) ctx.getGlobalContext().get("container");
        if(container != null) {
            requestRegistry.addContainer(requestId, container);
        }

        // Do the normal processing.
        super.processMessage(message, ctx);
    }

}
