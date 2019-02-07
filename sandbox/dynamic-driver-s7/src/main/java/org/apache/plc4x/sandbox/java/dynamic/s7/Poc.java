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

package org.apache.plc4x.sandbox.java.dynamic.s7;

import org.apache.commons.scxml2.SCXMLExecutor;
import org.apache.commons.scxml2.env.SimpleDispatcher;
import org.apache.commons.scxml2.env.SimpleErrorReporter;
import org.apache.commons.scxml2.invoke.SimpleSCXMLInvoker;
import org.apache.commons.scxml2.io.SCXMLReader;
import org.apache.commons.scxml2.model.CustomAction;
import org.apache.commons.scxml2.model.SCXML;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.sandbox.java.dynamic.actions.ConnectAction;
import org.apache.plc4x.sandbox.java.dynamic.actions.InitContextAction;
import org.apache.plc4x.sandbox.java.dynamic.actions.ReceiveAction;
import org.apache.plc4x.sandbox.java.dynamic.actions.SendAction;
import org.apache.plc4x.sandbox.java.dynamic.s7.actions.S7DecodeArticleNumber;

import java.util.HashMap;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;

public class Poc {

    private String dataFormatURI;

    private SCXMLExecutor executor;

    private Poc(String stateMachineURI, String dataFormatURI) throws Exception {
        this.dataFormatURI = dataFormatURI;

        // Initialize our PLC4X specific actions.
        List<CustomAction> customActions = new LinkedList<>();
        customActions.add(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "initContext", InitContextAction.class));
        customActions.add(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "connect", ConnectAction.class));
        customActions.add(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "send", SendAction.class));
        customActions.add(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "receive", ReceiveAction.class));
        customActions.add(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "S7DecodeArticleNumber", S7DecodeArticleNumber.class));

        // Initialize the state-machine with the definition from the protocol module.
        SCXML scxml = SCXMLReader.read(
            Poc.class.getClassLoader().getResource(stateMachineURI),
            new SCXMLReader.Configuration(null, null, customActions));

        // Create an executor for running the state-machine.
        executor = new SCXMLExecutor(null, new SimpleDispatcher(), new SimpleErrorReporter());
        executor.setStateMachine(scxml);
        executor.registerInvokerClass("scxml", SimpleSCXMLInvoker.class);
    }

    private void run() throws Exception {
        Map<String, Object> context = new HashMap<>();
        context.put("protocolDaffodilSchema", dataFormatURI);

        context.put("cotpLocalReference", "15");
        context.put("cotpCalledTsap", "512");
        context.put("cotpCallingTsap", "273");
        context.put("cotpTpduSize", "10");
        context.put("s7MaxAmqCaller", "10");
        context.put("s7MaxAmqCallee", "10");
        context.put("s7PduLength", "1024");

        //context.put("plcType", "HURZ");

        // Run the state-machine.
        executor.go(context);
    }

    public static void main(String[] args) throws Exception {
        /*Poc poc = new Poc(
            "org/apache/plc4x/protocols/s7/protocol.scxml.xml",
            "org/apache/plc4x/protocols/s7/protocol.dfdl.xsd");
        poc.run();*/

        PlcConnection connection = new DynamicS7Driver().connect("s7://10.10.64.20/1/1");
        connection.isConnected();
    }

}
