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
import org.apache.daffodil.japi.Compiler;
import org.apache.daffodil.japi.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.URL;
import java.util.List;

public class InitContextAction extends Action {

    private static final Logger logger = LoggerFactory.getLogger(InitContextAction.class);

    @Override
    public void execute(ActionExecutionContext ctx) {
        ctx.getAppLog().info("Initializing Context.");

        try {
            Compiler c = Daffodil.compiler();
            c.setValidateDFDLSchemas(true);
            URL shemaUrl = SendAction.class.getClassLoader().getResource("org/apache/plc4x/protocols/s7/protocol.dfdl.xsd");
            if (shemaUrl != null) {
                URI schemaUri = shemaUrl.toURI();
                ProcessorFactory pf = c.compileSource(schemaUri);
                if (pf.isError()) {
                    logDiagnosticInformation(pf);
                    TriggerEvent event = new EventBuilder("failure", TriggerEvent.SIGNAL_EVENT).build();
                    ctx.getInternalIOProcessor().addEvent(event);
                    return;
                }
                DataProcessor dp = pf.onPath("/");
                if (dp.isError()) {
                    logDiagnosticInformation(dp);
                    TriggerEvent event = new EventBuilder("failure", TriggerEvent.SIGNAL_EVENT).build();
                    ctx.getInternalIOProcessor().addEvent(event);
                    return;
                }
                ctx.getGlobalContext().set("dfdl", dp);
            }
        } catch (IOException | URISyntaxException e) {
            TriggerEvent event = new EventBuilder("failure", TriggerEvent.SIGNAL_EVENT).data(e).build();
            ctx.getInternalIOProcessor().addEvent(event);
            return;
        }

        TriggerEvent event = new EventBuilder("success", TriggerEvent.SIGNAL_EVENT).build();
        ctx.getInternalIOProcessor().addEvent(event);
    }

    private void logDiagnosticInformation(WithDiagnostics withDiagnostics) {
        List<Diagnostic> diags = withDiagnostics.getDiagnostics();
        for (Diagnostic d : diags) {
            logger.error(d.getSomeMessage());
        }
    }

}
