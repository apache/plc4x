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

package org.apache.plc4x.sandbox.java.s7;

import org.apache.commons.scxml2.SCXMLExecutor;
import org.apache.commons.scxml2.env.SimpleDispatcher;
import org.apache.commons.scxml2.env.SimpleErrorReporter;
import org.apache.commons.scxml2.invoke.SimpleSCXMLInvoker;
import org.apache.commons.scxml2.io.SCXMLReader;
import org.apache.commons.scxml2.model.CustomAction;
import org.apache.commons.scxml2.model.SCXML;
import org.apache.plc4x.sandbox.java.s7.actions.ConnectAction;
import org.apache.plc4x.sandbox.java.s7.actions.InitContextAction;
import org.apache.plc4x.sandbox.java.s7.actions.ReceiveAction;
import org.apache.plc4x.sandbox.java.s7.actions.SendAction;

import java.util.LinkedList;
import java.util.List;

public class Poc {

    private SCXMLExecutor executor;

    public Poc() throws Exception {

        List<CustomAction> customActions = new LinkedList<>();
        customActions.add(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "initContext", InitContextAction.class));
        customActions.add(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "connect", ConnectAction.class));
        customActions.add(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "send", SendAction.class));
        customActions.add(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "receive", ReceiveAction.class));

        SCXML scxml = SCXMLReader.read(
            Poc.class.getClassLoader().getResource("org/apache/plc4x/protocols/s7/protocol.scxml.xml"),
            new SCXMLReader.Configuration(null, null, customActions));
        executor = new SCXMLExecutor(null, new SimpleDispatcher(), new SimpleErrorReporter());
        executor.setStateMachine(scxml);
        executor.registerInvokerClass("scxml", SimpleSCXMLInvoker.class);
    }

    protected void run() throws Exception {
        executor.go();
    }

    public static void main(String[] args) throws Exception {
        Poc poc = new Poc();
        poc.run();
    }

}
