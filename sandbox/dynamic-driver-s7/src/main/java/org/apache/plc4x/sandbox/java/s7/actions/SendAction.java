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
import org.apache.commons.scxml2.model.ParsedValue;
import org.apache.commons.scxml2.model.ParsedValueContainer;
import org.apache.daffodil.japi.Compiler;
import org.apache.daffodil.japi.*;
import org.apache.daffodil.japi.infoset.InfosetInputter;
import org.apache.daffodil.japi.infoset.W3CDOMInfosetInputter;
import org.w3c.dom.Document;
import org.w3c.dom.Node;

import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.parsers.ParserConfigurationException;
import java.io.DataOutputStream;
import java.io.IOException;
import java.net.Socket;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.URL;
import java.nio.channels.Channels;
import java.nio.channels.WritableByteChannel;
import java.util.List;

public class SendAction extends Action implements ParsedValueContainer {

    private ParsedValue message;

    @Override
    public ParsedValue getParsedValue() {
        return message;
    }

    @Override
    public void setParsedValue(ParsedValue parsedValue) {
        message = parsedValue;
    }

    @Override
    public void execute(ActionExecutionContext ctx) {
        if(message != null) {
            if(message.getType() == ParsedValue.ValueType.NODE) {
                try {
                    Node messageTemplate = (Node) message.getValue();
                    DocumentBuilderFactory dbf = DocumentBuilderFactory.newInstance();
                    DocumentBuilder builder = dbf.newDocumentBuilder();
                    Document doc = builder.newDocument();
                    Node messageTemplateClone = doc.importNode(messageTemplate, true);
                    doc.appendChild(messageTemplateClone);
                    Compiler c = Daffodil.compiler();
                    c.setValidateDFDLSchemas(true);
                    URL shemaUrl = SendAction.class.getClassLoader().getResource("org/apache/plc4x/protocols/s7/protocol.dfdl.xsd");
                    if(shemaUrl != null) {
                        URI schemaUri = shemaUrl.toURI();
                        ProcessorFactory pf = c.compileSource(schemaUri);
                        if(pf.isError()) {
                            List<Diagnostic> diags = pf.getDiagnostics();
                            for (Diagnostic d : diags) {
                                System.err.println(d.getSomeMessage());
                            }
                            return;
                        }
                        DataProcessor dp = pf.onPath("/");
                        if(dp.isError()) {
                            List<Diagnostic> diags = dp.getDiagnostics();
                            for (Diagnostic d : diags) {
                                System.err.println(d.getSomeMessage());
                            }
                            return;
                        }
                        InfosetInputter inputter = new W3CDOMInfosetInputter(doc);

                        Socket connection = (Socket) ctx.getGlobalContext().get("connection");
                        DataOutputStream outputStream = new DataOutputStream(connection.getOutputStream());
                        WritableByteChannel wbc = Channels.newChannel(outputStream);

                        UnparseResult byteMessage = dp.unparse(inputter, wbc);
                        if(byteMessage.isError()) {
                            List<Diagnostic> diags = byteMessage.getDiagnostics();
                            for (Diagnostic d : diags) {
                                System.err.println(d.getSomeMessage());
                            }
                            return;
                        }

                        outputStream.flush();
                    }
                } catch(URISyntaxException e) {
                    e.printStackTrace();
                } catch (IOException e) {
                    e.printStackTrace();
                } catch (ParserConfigurationException e) {
                    e.printStackTrace();
                }
            } else if(message.getType() == ParsedValue.ValueType.JSON) {

            }
        }
        ctx.getAppLog().info("Sending.");
        TriggerEvent event = new EventBuilder("success", TriggerEvent.SIGNAL_EVENT).build();
        ctx.getInternalIOProcessor().addEvent(event);
    }

}
