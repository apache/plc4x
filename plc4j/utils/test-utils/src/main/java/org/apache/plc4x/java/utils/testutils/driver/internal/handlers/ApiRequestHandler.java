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
package org.apache.plc4x.java.utils.testutils.driver.internal.handlers;

import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.java.utils.testutils.driver.internal.TestContext;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.dom4j.Element;
import org.dom4j.QName;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

/**
 * Handler for executing API requests (read, write, etc.).
 * Requests are executed asynchronously and futures are stored in TestContext.
 */
public class ApiRequestHandler {

    private static final Logger LOGGER = LoggerFactory.getLogger(ApiRequestHandler.class);
    private static final long DEFAULT_TIMEOUT_MS = 5000;

    private final Element requestXml;

    public ApiRequestHandler(Element requestXml) {
        this.requestXml = requestXml;
    }

    /**
     * Executes the API request asynchronously.
     * The future is stored in the TestContext for later validation by api-response steps.
     *
     * @param connection the PLC connection
     * @param context the test context for storing the response future
     * @throws DriverTestsuiteException if execution fails
     */
    public void executeApiRequest(PlcConnection connection, TestContext context) {
        String requestType = requestXml.getName();

        switch (requestType) {
            case "PlcReadRequest":
            case "TestReadRequest":
                executeReadRequest(connection, context);
                break;
            case "PlcWriteRequest":
            case "TestWriteRequest":
                executeWriteRequest(connection, context);
                break;
            case "PlcBrowseRequest":
            case "TestBrowseRequest":
                executeBrowseRequest(connection, context);
                break;
            case "PlcSubscriptionRequest":
            case "TestSubscriptionRequest":
                executeSubscriptionRequest(connection, context);
                break;
            default:
                throw new DriverTestsuiteException("Unknown request type: " + requestType);
        }
    }

    private void executeReadRequest(PlcConnection connection, TestContext context) {
        try {
            PlcReadRequest.Builder builder = connection.readRequestBuilder();

            // Parse tags
            Element tagsElement = requestXml.element(new QName("tags"));
            if (tagsElement != null) {
                for (Element tag : tagsElement.elements()) {
                    String name = tag.element(new QName("name")).getTextTrim();
                    String address = tag.element(new QName("address")).getTextTrim();
                    builder.addTagAddress(name, address);
                }
            }

            PlcReadRequest request = builder.build();
            CompletableFuture<? extends PlcReadResponse> future = request.execute();

            // Store the future in context for api-response step to validate
            context.setPendingResponse(future);
            if (LOGGER.isInfoEnabled()) {
                StringBuilder sb = new StringBuilder("Read request initiated:\n");
                for (String tagName : request.getTagNames()) {
                    PlcResponseCode tagResponseCode = request.getTagResponseCode(tagName);
                    if(tagResponseCode != PlcResponseCode.OK) {
                        sb.append("  ").append(tagName).append(": ").append(tagResponseCode).append("\n");
                    } else {
                        sb.append("  ").append(tagName).append(": ").append(tagResponseCode).append(": ").append(request.getTag(tagName).getAddressString()).append("\n");
                    }
                }
                LOGGER.info(sb.toString());
            }

        } catch (Exception e) {
            throw new DriverTestsuiteException("Failed to execute read request", e);
        }
    }

    private void executeWriteRequest(PlcConnection connection, TestContext context) {
        try {
            PlcWriteRequest.Builder builder = connection.writeRequestBuilder();

            // Parse tags and values
            Element tagsElement = requestXml.element(new QName("tags"));
            if (tagsElement != null) {
                for (Element tag : tagsElement.elements()) {
                    // Support both attribute and element formats for compatibility
                    String name = tag.attributeValue("name");
                    if (name == null) {
                        Element nameElement = tag.element(new QName("name"));
                        if (nameElement != null) {
                            name = nameElement.getTextTrim();
                        }
                    }

                    String address = tag.attributeValue("address");
                    if (address == null) {
                        Element addressElement = tag.element(new QName("address"));
                        if (addressElement != null) {
                            address = addressElement.getTextTrim();
                        }
                    }

                    // Get value(s) - handle both single and multiple values
                    List<Element> valueElements = tag.elements(new QName("value"));
                    PlcValue value;

                    if (valueElements.isEmpty()) {
                        // No value elements, try text content
                        String valueStr = tag.getTextTrim();
                        value = parseValue(valueStr);
                    } else if (valueElements.size() == 1) {
                        // Single value element - check for nested PLC type elements
                        Element valueElement = valueElements.get(0);
                        List<Element> typeElements = valueElement.elements();

                        if (typeElements.isEmpty()) {
                            // Simple value - just text content
                            String valueStr = valueElement.getTextTrim();
                            value = parseValue(valueStr);
                        } else {
                            // Nested PLC type element (e.g., <PlcUINT>1</PlcUINT>)
                            Element typeElement = typeElements.get(0);
                            value = parseTypedValue(typeElement);
                        }
                    } else {
                        // Multiple value elements - create array
                        List<PlcValue> values = new ArrayList<>(valueElements.size());
                        for (Element valueElement : valueElements) {
                            List<Element> typeElements = valueElement.elements();
                            if (typeElements.isEmpty()) {
                                // Simple value
                                String valueStr = valueElement.getTextTrim();
                                values.add(parseValue(valueStr));
                            } else {
                                // Nested PLC type element
                                Element typeElement = typeElements.get(0);
                                values.add(parseTypedValue(typeElement));
                            }
                        }
                        value = new PlcList(values);
                    }

                    builder.addTagAddress(name, address, value);
                }
            }

            PlcWriteRequest request = builder.build();
            CompletableFuture<? extends PlcWriteResponse> future = request.execute();

            if (future.isDone()) {
                PlcWriteResponse writeResponse = future.get(1, TimeUnit.MILLISECONDS);
                writeResponse.getTagNames().forEach(tagName -> {
                    if(writeResponse.getResponseCode(tagName) != PlcResponseCode.OK) {
                        LOGGER.warn("Failed to write tag {}: {}", tagName, writeResponse.getResponseCode(tagName));
                        throw new DriverTestsuiteException("Failed to write tag " + tagName + ": " + writeResponse.getResponseCode(tagName));
                    } else {
                        LOGGER.info("Successfully wrote tag {}: {}", tagName, writeResponse.getResponseCode(tagName));
                    }
                });
            } else {
                // Store the future in context for api-response step to validate
                context.setPendingResponse(future);
                if (LOGGER.isInfoEnabled()) {
                    StringBuilder sb = new StringBuilder("Write request initiated:\n");
                    for (String tagName : request.getTagNames()) {
                        PlcResponseCode tagResponseCode = request.getTagResponseCode(tagName);
                        if (tagResponseCode != PlcResponseCode.OK) {
                            sb.append("  ").append(tagName).append(": ").append(tagResponseCode).append("\n");
                        } else {
                            sb.append("  ").append(tagName).append(": ").append(tagResponseCode).append(": ").append(request.getTag(tagName).getAddressString()).append(": ").append(request.getPlcValue(tagName)).append("\n");
                        }
                    }
                    LOGGER.info(sb.toString());
                }
            }
            LOGGER.debug("Write request initiated asynchronously");

        } catch (Exception e) {
            throw new DriverTestsuiteException("Failed to execute write request", e);
        }
    }

    private void executeBrowseRequest(PlcConnection connection, TestContext context) {
        try {
            PlcBrowseRequest.Builder builder = connection.browseRequestBuilder();

            // Parse query
            Element queriesElement = requestXml.element(new QName("queries"));
            if (queriesElement != null) {
                for (Element tag : queriesElement.elements()) {
                    String name = tag.getName();
                    String address = tag.getTextTrim();
                    builder.addQuery(name, address);
                }
            }

            PlcBrowseRequest request = builder.build();
            CompletableFuture<? extends PlcBrowseResponse> future = request.execute();

            // Store the future in context for api-response step to validate
            context.setPendingResponse(future);
            LOGGER.debug("Browse request initiated asynchronously");

        } catch (Exception e) {
            throw new DriverTestsuiteException("Failed to execute browse request", e);
        }
    }

    private void executeSubscriptionRequest(PlcConnection connection, TestContext context) {
        try {
            PlcSubscriptionRequest.Builder builder = connection.subscriptionRequestBuilder();

            // Parse tags
            Element tagsElement = requestXml.element(new QName("tags"));
            if (tagsElement != null) {
                for (Element tag : tagsElement.elements()) {
                    String name = tag.element(new QName("name")).getTextTrim();
                    String address = tag.element(new QName("address")).getTextTrim();
                    builder.addChangeOfStateTagAddress(name, address);
                }
            }
            // Add incoming events to the context.
            builder.setConsumer(context::addSubscriptionEvent);

            PlcSubscriptionRequest request = builder.build();
            CompletableFuture<? extends PlcSubscriptionResponse> future = request.execute();

            // Store the future in context for api-response step to validate
            context.setPendingResponse(future);
            LOGGER.debug("Subscription request initiated asynchronously");

        } catch (Exception e) {
            throw new DriverTestsuiteException("Failed to execute subscription request", e);
        }
    }

    /**
     * Parses a typed PLC value element (e.g., <PlcUINT>1</PlcUINT>).
     * Extracts the value and converts it to the appropriate Java type.
     *
     * @param typeElement the PLC type element
     * @return the parsed value as a Java object
     */
    private PlcValue parseTypedValue(Element typeElement) {
        String typeName = typeElement.getName();
        String valueStr = typeElement.getTextTrim();

        // Handle different PLC types and convert to appropriate Java types
        switch (typeName) {
            // Boolean types
            case "PlcBOOL":
            case "PlcBool":
                return new PlcBOOL(Boolean.parseBoolean(valueStr));

            // Integer types (8-bit)
            case "PlcSINT":
                return new PlcSINT(Byte.parseByte(valueStr));

            // Unsigned 8-bit integer - use Short to handle full range (0-255)
            case "PlcUSINT":
                return new PlcUSINT(Short.parseShort(valueStr));

            case "PlcBYTE":
                return new PlcBYTE(Short.parseShort(valueStr));

            // Integer types (16-bit)
            case "PlcINT":
                return new PlcINT(Short.parseShort(valueStr));

            // Unsigned 16-bit integer - use Integer to handle full range (0-65535)
            case "PlcUINT":
                return new PlcUINT(Integer.parseInt(valueStr));

            case "PlcWORD":
                return new PlcWORD(Integer.parseInt(valueStr));

            // Integer types (32-bit)
            case "PlcDINT":
                return new PlcDINT(Integer.parseInt(valueStr));

            // Unsigned 32-bit integer - use Long to handle full range (0-4294967295)
            case "PlcUDINT":
                return new PlcUDINT(Long.parseLong(valueStr));

            case "PlcDWORD":
                return new PlcDWORD(Long.parseLong(valueStr));

            // Integer types (64-bit)
            case "PlcLINT":
                return new PlcLINT(Long.parseLong(valueStr));

            // Unsigned 64-bit integer - use BigInteger to handle full range (0-???)
            case "PlcULINT":
                return new PlcULINT(new BigInteger(valueStr));

            case "PlcLWORD":
                return new PlcLWORD(new BigInteger(valueStr));

            // Floating point types
            case "PlcREAL":
                return new PlcREAL(Float.parseFloat(valueStr));

            case "PlcLREAL":
                return new PlcLREAL(Double.parseDouble(valueStr));

            // String types
            case "PlcSTRING":
                return new PlcSTRING(valueStr);

            case "PlcWSTRING":
                return new PlcWSTRING(valueStr);

            case "PlcCHAR":
                return new PlcCHAR(valueStr);

            case "PlcWCHAR":
                return new PlcWCHAR(valueStr);

            // Temporal types
            case "PlcTIME":
                // Parse ISO 8601 duration format (e.g., PT2.5S, PT3.20005S)
                return new PlcTIME(Duration.parse(valueStr));

            case "PlcLTIME":
                // Parse ISO 8601 duration format (e.g., PT2.5S, PT3.20005S)
                return new PlcLTIME(Duration.parse(valueStr));

            case "PlcDATE":
                return new PlcDATE(LocalDate.parse(valueStr));

            case "PlcLDATE":
                return new PlcLDATE(LocalDate.parse(valueStr));

            case "PlcTIME_OF_DAY":
                return new PlcTIME_OF_DAY(LocalTime.parse(valueStr));

            case "PlcLTIME_OF_DAY":
                return new PlcLTIME_OF_DAY(LocalTime.parse(valueStr));

            case "PlcDATE_AND_TIME":
                return new PlcDATE_AND_TIME(LocalDateTime.parse(valueStr));

            case "PlcDATE_AND_LTIME":
                return new PlcDATE_AND_LTIME(LocalDateTime.parse(valueStr));

            case "PlcLDATE_AND_TIME":
                return new PlcLDATE_AND_TIME(LocalDateTime.parse(valueStr));

            // Raw byte array: <PlcRAW_BYTE_ARRAY>0001020304</PlcRAW_BYTE_ARRAY>
            case "PlcRAW_BYTE_ARRAY":
                try {
                    return new PlcRawByteArray(org.apache.commons.codec.binary.Hex.decodeHex(valueStr));
                } catch (org.apache.commons.codec.DecoderException e) {
                    throw new DriverTestsuiteException("Invalid hex in PlcRAW_BYTE_ARRAY: " + valueStr, e);
                }

            // List type
            case "PlcList":
                List<PlcValue> listValues = new ArrayList<>(typeElement.elements().size());
                for (int i = 0; i < typeElement.elements().size(); i++) {
                    Element itemElement = typeElement.elements().get(i);
                    listValues.add(parseTypedValue(itemElement));
                }
                return new PlcList(listValues);

            // Struct type
            case "PlcStruct":
                // PlcStruct contains named fields, each field is an element whose name is the field name
                // and whose content is a PlcValue type element
                Map<String, PlcValue> structFields = new LinkedHashMap<>();
                for (Element fieldElement : typeElement.elements()) {
                    String fieldName = fieldElement.getName();
                    // Each field element should contain exactly one child element which is the PlcValue type
                    List<Element> fieldValueElements = fieldElement.elements();
                    if (fieldValueElements.isEmpty()) {
                        // No nested element, use the text content directly
                        String fieldValue = fieldElement.getTextTrim();
                        structFields.put(fieldName, parseValue(fieldValue));
                    } else {
                        // Parse the nested PlcValue type element
                        Element valueTypeElement = fieldValueElements.get(0);
                        structFields.put(fieldName, parseTypedValue(valueTypeElement));
                    }
                }
                return new PlcStruct(structFields);

            // Default: try to infer the type from the value
            default:
                LOGGER.warn("Unknown PLC type: {}, falling back to value parsing", typeName);
                return parseValue(valueStr);
        }
    }

    private PlcValue parseValue(String valueStr) {
        // Try to parse as different types, from narrowest to widest
        try {
            return new PlcDINT(Integer.parseInt(valueStr));
        } catch (NumberFormatException e) {
            // Not an integer
        }

        try {
            return new PlcLINT(Long.parseLong(valueStr));
        } catch (NumberFormatException e) {
            // Not a long
        }

        try {
            BigInteger bigInt = new BigInteger(valueStr);
            // If it fits in unsigned 64-bit range, use ULINT; otherwise fall through
            if (bigInt.signum() >= 0 && bigInt.bitLength() <= 64) {
                return new PlcULINT(bigInt);
            }
        } catch (NumberFormatException e) {
            // Not a BigInteger
        }

        try {
            return new PlcLREAL(Double.parseDouble(valueStr));
        } catch (NumberFormatException e) {
            // Not a double
        }

        if ("true".equalsIgnoreCase(valueStr) || "false".equalsIgnoreCase(valueStr)) {
            return new PlcBOOL(Boolean.parseBoolean(valueStr));
        }

        // Default to string
        return new PlcSTRING(valueStr);
    }
}
