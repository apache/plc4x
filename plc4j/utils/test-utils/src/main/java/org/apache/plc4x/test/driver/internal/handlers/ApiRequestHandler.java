/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.test.driver.internal.handlers;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.test.driver.internal.utils.Synchronizer;
import org.dom4j.Element;

import java.util.ArrayList;
import java.util.List;

public class ApiRequestHandler {

    private final Element payload;

    private final Synchronizer synchronizer;

    public ApiRequestHandler(Element payload, Synchronizer synchronizer) {
        this.payload = payload;
        this.synchronizer = synchronizer;
    }

    public void executeApiRequest(PlcConnection plcConnection) {
        assert synchronizer != null;
        String typeName = payload.getName();
        switch (typeName) {
            case "TestReadRequest": {
                final PlcReadRequest.Builder builder = plcConnection.readRequestBuilder();
                if (payload.element("fields") != null) {
                    for (Element fieldElement : payload.element("fields").elements("field")) {
                        builder.addItem(fieldElement.elementText("name"), fieldElement.elementText("address"));
                    }
                }
                final PlcReadRequest plc4xRequest = builder.build();
                // Currently we can only process one response at at time, throw an error if more
                // are submitted.
                if (synchronizer.responseFuture != null) {
                    throw new DriverTestsuiteException("Previous response not handled.");
                }
                // Save the response for being used later on.
                synchronizer.responseFuture = plc4xRequest.execute();
                break;
            }
            case "TestWriteRequest": {
                final PlcWriteRequest.Builder builder = plcConnection.writeRequestBuilder();
                if (payload.element("fields") != null) {
                    for (Element fieldElement : payload.element("fields").elements("field")) {
                        List<Element> valueElements = fieldElement.elements("value");
                        List<String> valueStrings = new ArrayList<>(valueElements.size());
                        for (Element valueElement : valueElements) {
                            valueStrings.add(valueElement.getTextTrim());
                        }
                        builder.addItem(fieldElement.elementText("name"),
                            fieldElement.elementText("address"), valueStrings.toArray(new Object[0]));
                    }
                }
                final PlcWriteRequest plc4xRequest = builder.build();
                // Currently we can only process one response at at time, throw an error if more
                // are submitted.
                if (synchronizer.responseFuture != null) {
                    throw new DriverTestsuiteException("Previous response not handled.");
                }
                // Save the response for being used later on.
                synchronizer.responseFuture = plc4xRequest.execute();
                break;
            }
            default:
                throw new PlcRuntimeException("Unknown class name" + typeName);
        }
    }

}
