/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.test.driver.internal.handlers;

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.test.driver.internal.utils.Synchronizer;
import org.dom4j.Element;

import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;

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
                if (payload.element("tags") != null) {
                    for (Element tagElement : payload.element("tags").elements("tag")) {
                        builder.addTagAddress(tagElement.elementText("name"), tagElement.elementText("address"));
                    }
                }
                final PlcReadRequest plc4xRequest = builder.build();
                // Currently we can only process one response at a time, throw an error if more
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
                if (payload.element("tags") != null) {
                    for (Element tagElement : payload.element("tags").elements("tag")) {
                        List<Element> valueElements = tagElement.elements("value");
                        List<Object> valueObjects = new ArrayList<>(valueElements.size());
                        for (Element valueElement : valueElements) {
                            // If the element is a simple text element, use this driectly.
                            if(valueElement.isTextOnly()) {
                                valueObjects.add(valueElement.getTextTrim());
                            }
                            // If the value is a "PlcValue" element, try to parse this.
                            else if (valueElement.elements() != null) {
                                for (Element element : valueElement.elements()) {
                                    valueObjects.add(parsePlcValue(element));
                                }
                            }
                        }
                        builder.addTagAddress(tagElement.elementText("name"),
                            tagElement.elementText("address"), valueObjects.toArray(new Object[0]));
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
            case "TestSubscriptionRequest":{
                // TODO: chris add your stuff here...
                throw new NotImplementedException();
            }
            default:
                throw new PlcRuntimeException("Unknown class name " + typeName);
        }
    }

    protected PlcValue parsePlcValue(Element element) {
        switch (element.getName()) {
            case "PlcStruct": {
                Map<String, PlcValue> structContent = new LinkedHashMap<>();
                for (Element field : element.elements()) {
                    String fieldName = field.getName();
                    if((field.elements() == null) || (field.elements().size() != 1)) {
                        throw new RuntimeException("Expected exactly one child element");
                    }
                    Element valueElement = field.elements().get(0);
                    PlcValue fieldValue = parsePlcValue(valueElement);
                    structContent.put(fieldName, fieldValue);
                }
                return new PlcStruct(structContent);
            }
            case "PlcList": {
                // TODO: Implement this ...
            }
            case "PlcNULL": {
                return new PlcNull();
            }

            case "PlcBOOL": {
                return new PlcBOOL(element.getTextTrim());
            }
            case "PlcBYTE": {
                return new PlcBYTE(element.getTextTrim());
            }
            case "PlcWORD": {
                return new PlcWORD(element.getTextTrim());
            }
            case "PlcDWORD": {
                return new PlcDWORD(element.getTextTrim());
            }
            case "PlcLWORD": {
                return new PlcLWORD(element.getTextTrim());
            }

            case "PlcUSINT": {
                return new PlcUSINT(element.getTextTrim());
            }
            case "PlcUINT": {
                return new PlcUINT(element.getTextTrim());
            }
            case "PlcUDINT": {
                return new PlcUDINT(element.getTextTrim());
            }
            case "PlcULINT": {
                return new PlcULINT(element.getTextTrim());
            }

            case "PlcSINT": {
                return new PlcSINT(element.getTextTrim());
            }
            case "PlcINT": {
                return new PlcINT(element.getTextTrim());
            }
            case "PlcDINT": {
                return new PlcDINT(element.getTextTrim());
            }
            case "PlcLINT": {
                return new PlcLINT(element.getTextTrim());
            }

            case "PlcREAL": {
                return new PlcREAL(element.getTextTrim());
            }
            case "PlcLREAL": {
                return new PlcLREAL(element.getTextTrim());
            }

            case "PlcCHAR": {
                return new PlcCHAR(element.getTextTrim());
            }
            case "PlcWCHAR": {
                return new PlcWCHAR(element.getTextTrim());
            }
            case "PlcSTRING": {
                return new PlcSTRING(element.getTextTrim());
            }
            case "PlcWSTRING": {
                return new PlcWSTRING(element.getTextTrim());
            }

            case "PlcTIME": {
                long milliseconds = Long.parseLong(element.getTextTrim());
                return new PlcTIME(milliseconds);
            }
            case "PlcLTIME": {
                return new PlcLTIME(Duration.parse(element.getTextTrim()));
            }
            case "PlcDATE": {
                return new PlcDATE(LocalDate.parse(element.getTextTrim()));
            }
            case "PlcLDATE": {
                return new PlcLDATE(LocalDate.parse(element.getTextTrim()));
            }
            case "PlcTIME_OF_DAY": {
                return new PlcTIME_OF_DAY(LocalTime.parse(element.getTextTrim()));
            }
            case "PlcLTIME_OF_DAY": {
                return new PlcLTIME_OF_DAY(LocalTime.parse(element.getTextTrim()));
            }
            case "PlcDATE_AND_TIME": {
                return new PlcDATE_AND_TIME(LocalDateTime.parse(element.getTextTrim()));
            }
            case "PlcDATE_AND_LTIME": {
                return new PlcDATE_AND_LTIME(LocalDateTime.parse(element.getTextTrim()));
            }
            case "PlcLDATE_AND_TIME": {
                return new PlcLDATE_AND_TIME(LocalDateTime.parse(element.getTextTrim()));
            }

            /*case "PlcRAW_BYTE_ARRAY": {
                return new PlcRAW_BYTE_ARRAY(element.getTextTrim());
            }*/
            default: throw new RuntimeException("Unsupported type of PlcValue " + element.getName());
        }
    }

}
