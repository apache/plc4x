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
import org.apache.commons.scxml2.model.ParsedValue;
import org.apache.daffodil.japi.DataProcessor;
import org.apache.daffodil.japi.UnparseResult;
import org.apache.daffodil.japi.infoset.InfosetInputter;
import org.apache.daffodil.japi.infoset.W3CDOMInfosetInputter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.w3c.dom.Document;
import org.w3c.dom.Node;

import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.parsers.ParserConfigurationException;
import java.io.DataOutputStream;
import java.io.IOException;
import java.net.Socket;
import java.nio.channels.Channels;
import java.nio.channels.WritableByteChannel;

public class SendAction extends BasePlc4xAction {

    @Override
    protected Logger getLogger() {
        return LoggerFactory.getLogger(SendAction.class);
    }

    @Override
    public void execute(ActionExecutionContext ctx) {
        ctx.getAppLog().info("Sending.");
        if(getParsedValue() != null) {
            if(getParsedValue().getType() == ParsedValue.ValueType.NODE) {
                try {
                    Node messageTemplate = (Node) getParsedValue().getValue();
                    DocumentBuilderFactory dbf = DocumentBuilderFactory.newInstance();
                    DocumentBuilder builder = dbf.newDocumentBuilder();
                    Document doc = builder.newDocument();
                    Node messageTemplateClone = doc.importNode(messageTemplate, true);
                    doc.appendChild(messageTemplateClone);

                    DataProcessor dp = getDaffodilDataProcessor(ctx);
                    if(dp == null) {
                        TriggerEvent event = new EventBuilder("failure", TriggerEvent.SIGNAL_EVENT).
                            data("Couldn't initialize daffodil data processor.").build();
                        ctx.getInternalIOProcessor().addEvent(event);
                        return;
                    }
                    InfosetInputter inputter = new W3CDOMInfosetInputter(doc);

                    Socket connection = getSocket(ctx);
                    DataOutputStream outputStream = new DataOutputStream(connection.getOutputStream());
                    WritableByteChannel wbc = Channels.newChannel(outputStream);
                    UnparseResult byteMessage = dp.unparse(inputter, wbc);
                    if(byteMessage.isError()) {
                        logDiagnosticInformation(byteMessage);
                        return;
                    }
                    outputStream.flush();
                    ctx.getAppLog().info("Successfully sent message.");
                } catch(IOException | ParserConfigurationException e) {
                    e.printStackTrace();
                }
            } else {
                TriggerEvent event = new EventBuilder("failure", TriggerEvent.SIGNAL_EVENT).
                    data("type '" + getParsedValue().getType() + "' not supported").build();
                ctx.getInternalIOProcessor().addEvent(event);
                return;
            }
        }
        TriggerEvent event = new EventBuilder("success", TriggerEvent.SIGNAL_EVENT).build();
        ctx.getInternalIOProcessor().addEvent(event);
    }

}
