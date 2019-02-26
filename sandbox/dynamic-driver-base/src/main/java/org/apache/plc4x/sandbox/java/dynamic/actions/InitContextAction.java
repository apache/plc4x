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
import org.apache.commons.scxml2.EventBuilder;
import org.apache.commons.scxml2.TriggerEvent;
import org.apache.daffodil.japi.Compiler;
import org.apache.daffodil.japi.*;
import org.apache.plc4x.sandbox.java.dynamic.utils.RequestRegistry;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.URI;
import java.net.URL;
import java.util.List;

public class InitContextAction extends BasePlc4xAction {

    private long maxRequestId;
    private String protocolDaffodilSchemaName;

    public String getMaxRequestId() {
        return Long.toString(maxRequestId);
    }

    public void setMaxRequestId(String maxRequestId) {
        this.maxRequestId = Long.valueOf(maxRequestId);
    }

    public String getProtocolDaffodilSchemaName() {
        return protocolDaffodilSchemaName;
    }

    public void setProtocolDaffodilSchemaName(String protocolDaffodilSchemaName) {
        this.protocolDaffodilSchemaName = protocolDaffodilSchemaName;
    }

    @Override
    protected Logger getLogger() {
        return LoggerFactory.getLogger(InitContextAction.class);
    }

    @Override
    public void execute(ActionExecutionContext ctx) {
        getLogger().info(getStateName() + ": Initializing Context...");

        // Initialize the Daffodil system for parsing and serializing the
        // protocol messages.
        try {
            Compiler c = Daffodil.compiler();
            c.setValidateDFDLSchemas(true);
            String schemaUrlString = (String) ctx.getGlobalContext().get(protocolDaffodilSchemaName);
            URL schemaUrl = SendAction.class.getClassLoader().getResource(schemaUrlString);
            if (schemaUrl != null) {
                URI schemaUri = schemaUrl.toURI();
                ProcessorFactory pf = c.compileSource(schemaUri);
                logDiagnosticInformation(pf);
                DataProcessor dp = pf.onPath("/");
                logDiagnosticInformation(dp);
                ctx.getGlobalContext().set("dfdl", dp);
            }
        } catch (Exception e) {
            fireFailureEvent(ctx, "Error initializing daffodil schema");
            TriggerEvent event = new EventBuilder("failure", TriggerEvent.SIGNAL_EVENT).data(e).build();
            ctx.getInternalIOProcessor().addEvent(event);
            return;
        }

        // Create a new request-registry that will be used for matching
        // requests and responses.
        ctx.getGlobalContext().set("requestRegistry", new RequestRegistry(maxRequestId));

        getLogger().info("Context initialized.");

        fireSuccessEvent(ctx);
    }

    private void logDiagnosticInformation(WithDiagnostics withDiagnostics) throws Exception {
        if(withDiagnostics.isError()) {
            List<Diagnostic> diags = withDiagnostics.getDiagnostics();
            for (Diagnostic d : diags) {
                getLogger().error(d.getSomeMessage());
            }
            throw new Exception();
        }
    }

}
