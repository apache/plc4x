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
import org.apache.commons.scxml2.model.NodeListValue;
import org.apache.commons.scxml2.model.NodeValue;
import org.apache.commons.scxml2.model.ParsedValue;
import org.apache.daffodil.japi.DataProcessor;
import org.apache.daffodil.japi.ParseResult;
import org.apache.daffodil.japi.infoset.W3CDOMInfosetOutputter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.w3c.dom.Document;
import org.w3c.dom.Element;
import org.w3c.dom.Node;

import java.io.BufferedReader;
import java.io.DataInputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.Socket;
import java.nio.channels.Channels;
import java.nio.channels.ReadableByteChannel;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class ReceiveAction extends BasePlc4xAction {

    private String timeout;

    private final Map<String, String> verificationRules;
    private final Map<String, String> extractionRules;

    public ReceiveAction() {
        verificationRules = new HashMap<>();
        extractionRules = new HashMap<>();
    }

    @Override
    protected Logger getLogger() {
        return LoggerFactory.getLogger(ReceiveAction.class);
    }

    public String getTimeout() {
        return timeout;
    }

    public void setTimeout(String timeout) {
        this.timeout = timeout;
    }

    @Override
    @SuppressWarnings("unchecked")
    public void setParsedValue(ParsedValue parsedValue) {
        super.setParsedValue(parsedValue);

        if(parsedValue != null) {
            if(parsedValue instanceof NodeListValue) {
                List<Node> ruleList = (List<Node>) parsedValue.getValue();
                for (Node node : ruleList) {
                    if(node instanceof Element) {
                        parseElement((Element) node);
                    }
                }
            } else if(parsedValue instanceof NodeValue) {
                parseElement((Element) parsedValue.getValue());
            }
        }
    }

    private void parseElement(Element ruleElement) {
        String name = ruleElement.getAttribute("name");
        String xpath = ruleElement.getAttribute("xpath");
        if("verification".equals(ruleElement.getTagName())) {
            verificationRules.put(name, xpath);
        } else if("extraction".equals(ruleElement.getTagName())) {
            extractionRules.put(name, xpath);
        } else {
            getLogger().error("unsupported rule type: " + ruleElement.getTagName());
        }
    }

    @Override
    public void execute(ActionExecutionContext ctx) {
        ctx.getAppLog().info("Receiving.");
        try {
            DataProcessor dp = getDaffodilDataProcessor(ctx);
            if(dp == null) {
                TriggerEvent event = new EventBuilder("failure", TriggerEvent.SIGNAL_EVENT).
                    data("Couldn't initialize daffodil data processor.").build();
                ctx.getInternalIOProcessor().addEvent(event);
                return;
            }

            try {
                Thread.sleep(1000);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
            Socket connection = getSocket(ctx);
            BufferedReader in = new BufferedReader(new InputStreamReader(connection.getInputStream()));
            System.out.println(in.readLine());

            DataInputStream inputStream = new DataInputStream(connection.getInputStream());
            ReadableByteChannel rbc = Channels.newChannel(inputStream);
            W3CDOMInfosetOutputter outputter = new W3CDOMInfosetOutputter();
            ParseResult byteMessage = dp.parse(rbc, outputter);
            if (byteMessage.isError()) {
                logDiagnosticInformation(byteMessage);
                return;
            }

            Document message = outputter.getResult();
            System.out.println(message);
            ctx.getAppLog().info("Successfully sent message.");
        } catch (IOException e) {
            e.printStackTrace();
        }

        TriggerEvent event = new EventBuilder("success", TriggerEvent.SIGNAL_EVENT).build();
        ctx.getInternalIOProcessor().addEvent(event);
    }

}
