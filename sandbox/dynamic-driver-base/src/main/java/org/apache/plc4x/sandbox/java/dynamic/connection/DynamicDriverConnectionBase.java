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

package org.apache.plc4x.sandbox.java.dynamic.connection;

import org.apache.commons.scxml2.EventBuilder;
import org.apache.commons.scxml2.SCXMLExecutor;
import org.apache.commons.scxml2.TriggerEvent;
import org.apache.commons.scxml2.env.SimpleDispatcher;
import org.apache.commons.scxml2.env.SimpleErrorReporter;
import org.apache.commons.scxml2.invoke.SimpleSCXMLInvoker;
import org.apache.commons.scxml2.io.SCXMLReader;
import org.apache.commons.scxml2.model.CustomAction;
import org.apache.commons.scxml2.model.ModelException;
import org.apache.commons.scxml2.model.SCXML;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.base.connection.AbstractPlcConnection;
import org.apache.plc4x.sandbox.java.dynamic.actions.*;

import javax.xml.stream.XMLStreamException;
import java.io.IOException;
import java.util.*;

public abstract class DynamicDriverConnectionBase extends AbstractPlcConnection implements PlcConnection {

    private String stateMachineURI;
    private String dataFormatURI;
    private SCXMLExecutor executor;

    protected DynamicDriverConnectionBase(String stateMachineURI, String dataFormatURI) {
        this.stateMachineURI = stateMachineURI;
        this.dataFormatURI = dataFormatURI;
    }

    private void init() throws PlcConnectionException {
        // Initialize our PLC4X specific actions.
        List<CustomAction> customActions = new LinkedList<>();
        customActions.add(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "initContext", InitContextAction.class));
        customActions.add(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "connect", ConnectAction.class));
        customActions.add(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "send", SendAction.class));
        customActions.add(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "receiveExtractVerify", ReceiveExtractVerifyAction.class));
        customActions.add(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "sendRequest", SendRequestAction.class));
        customActions.addAll(getAdditionalCustomActions());

        try {
            // Initialize the state-machine with the definition from the protocol module.
            SCXML scxml = SCXMLReader.read(
                DynamicDriverConnectionBase.class.getClassLoader().getResource(stateMachineURI),
                new SCXMLReader.Configuration(null, null, customActions));

            // Create an executor for running the state-machine.
            executor = new SCXMLExecutor(null, new SimpleDispatcher(), new SimpleErrorReporter());
            executor.setStateMachine(scxml);
            executor.registerInvokerClass("scxml", SimpleSCXMLInvoker.class);
        } catch (XMLStreamException | IOException | ModelException e) {
            throw new PlcConnectionException("Error initializing driver state-machine", e);
        }
    }

    protected SCXMLExecutor getExecutor() {
        return executor;
    }

    protected Collection<CustomAction> getAdditionalCustomActions() {
        return Collections.emptyList();
    }

    protected Map<String, Object> getAdditionalContextDataItems() {
        return Collections.emptyMap();
    }

    protected abstract String getConnectedStateName();

    protected abstract String getDisconnectTransitionName();

    @Override
    public void connect() throws PlcConnectionException {
        // Setup the driver.
        init();

        // Initialize the drivers state.
        Map<String, Object> context = new HashMap<>();
        context.put("protocolDaffodilSchema", dataFormatURI);
        getAdditionalContextDataItems().forEach(context::put);

        try {
            // Run the state-machine.
            executor.go(context);
        } catch (ModelException e) {
            throw new PlcConnectionException("Error initializing driver state-machine", e);
        }
    }

    @Override
    public boolean isConnected() {
        if(executor == null) {
            return false;
        }
        return executor.getStatus().isInState(getConnectedStateName());
    }

    @Override
    public void close() throws Exception {
        if(executor == null) {
            return;
        }
        executor.triggerEvent(new EventBuilder(getDisconnectTransitionName(), TriggerEvent.CALL_EVENT).build());
    }

    @Override
    public PlcConnectionMetadata getMetadata() {
        return new PlcConnectionMetadata() {
            @Override
            public boolean canRead() {
                return true;
            }

            @Override
            public boolean canWrite() {
                return false;
            }

            @Override
            public boolean canSubscribe() {
                return false;
            }
        };
    }

}
