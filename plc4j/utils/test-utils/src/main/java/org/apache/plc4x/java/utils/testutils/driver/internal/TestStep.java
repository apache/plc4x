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
package org.apache.plc4x.java.utils.testutils.driver.internal;

import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.java.utils.testutils.driver.internal.handlers.*;
import org.apache.plc4x.java.utils.testutils.driver.internal.utils.Delay;
import org.apache.plc4x.java.utils.testutils.utils.dom4j.LocationAwareElement;
import org.apache.plc4x.java.utils.testutils.utils.model.Location;
import org.apache.plc4x.java.utils.testutils.utils.model.LocationAware;
import org.apache.plc4x.java.api.PlcConnection;
import org.dom4j.Element;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Optional;

/**
 * Represents a single test step in a test case.
 */
public class TestStep implements LocationAware {

    private static final Logger LOGGER = LoggerFactory.getLogger(TestStep.class);

    private final String basePackage;
    private final StepType type;
    private final String name;
    private final Element payload;
    private final String parserArguments;
    private final Location location;

    private TestStep(String basePackage, StepType type, String name, Element payload, String parserArguments, Location location) {
        this.basePackage = basePackage;
        this.type = type;
        this.name = name;
        this.payload = payload;
        this.parserArguments = parserArguments;
        this.location = location;
    }

    /**
     * Parses a test step from XML.
     *
     * @param stepElement the step XML element
     * @param byteOrder   the byte order name
     * @return the parsed test step
     * @throws DriverTestsuiteException if parsing fails
     */
    public static TestStep parseTestStep(String basePackage, Element stepElement, String byteOrder) {
        String stepTypeName = stepElement.getName();
        StepType stepType;

        try {
            stepType = StepType.valueOf(stepTypeName.toUpperCase().replace("-", "_"));
        } catch (IllegalArgumentException e) {
            throw new DriverTestsuiteException("Unknown step type: " + stepTypeName);
        }

        String name = stepElement.attributeValue("name");

        String parserArguments = stepElement.attributeValue("parserArguments");

        // Extract location if available
        Location location = null;
        if (stepElement instanceof LocationAwareElement) {
            LocationAwareElement locAware = (LocationAwareElement) stepElement;
            location = locAware.getLocation();
        }

        return new TestStep(basePackage, stepType, name, stepElement, parserArguments, location);
    }

    /**
     * Executes this test step.
     *
     * @param connection        the PLC connection
     * @param transportInstance the transport instance
     * @param byteOrder         the byte order name
     * @param context           the test context for sharing state between steps
     * @throws DriverTestsuiteException if execution fails
     */
    public void execute(PlcConnection connection, TransportInstance<?> transportInstance, String byteOrder, TestContext context) {
        LOGGER.info("Executing step: {} '{}'", type, name);

        switch (type) {
            case OUTGOING_PLC_BYTES:
                executeOutgoingPlcBytes(transportInstance, byteOrder);
                break;
            case INCOMING_PLC_BYTES:
                executeIncomingPlcBytes(transportInstance, byteOrder);
                break;
            case OUTGOING_PLC_MESSAGE:
                executeOutgoingPlcMessage(transportInstance, byteOrder);
                break;
            case INCOMING_PLC_MESSAGE:
                executeIncomingPlcMessage(transportInstance, byteOrder);
                break;
            case API_REQUEST:
                executeApiRequest(connection, context);
                break;
            case API_RESPONSE:
                executeApiResponse(context);
                break;
            case API_EVENT:
                executeApiEvent(context);
                break;
            case DELAY:
                executeDelay();
                break;
            case TERMINATE:
                executeTerminate();
                break;
            default:
                throw new DriverTestsuiteException("Unsupported step type: " + type);
        }
    }

    private void executeOutgoingPlcBytes(TransportInstance<?> transportInstance, String byteOrder) {
        OutgoingPlcBytesHandler handler = new OutgoingPlcBytesHandler(basePackage, payload);
        handler.executeOutgoingPlcBytes(transportInstance, byteOrder);
    }

    private void executeIncomingPlcBytes(TransportInstance<?> transportInstance, String byteOrder) {
        IncomingPlcBytesHandler handler = new IncomingPlcBytesHandler(payload);
        handler.executeIncomingPlcBytes(transportInstance, byteOrder);
    }

    private void executeOutgoingPlcMessage(TransportInstance<?> transportInstance, String byteOrder) {
        OutgoingPlcMessageHandler handler = new OutgoingPlcMessageHandler(basePackage, payload);
        handler.executeOutgoingPlcMessage(transportInstance, byteOrder);
    }

    private void executeIncomingPlcMessage(TransportInstance<?> transportInstance, String byteOrder) {
        IncomingPlcMessageHandler handler = new IncomingPlcMessageHandler(basePackage, payload);
        handler.executeIncomingPlcMessage(transportInstance, byteOrder);
    }

    private void executeApiRequest(PlcConnection connection, TestContext context) {
        // Get the actual request element (first child)
        Element requestElement = payload.elements().isEmpty() ? payload : (Element) payload.elements().get(0);
        ApiRequestHandler handler = new ApiRequestHandler(requestElement);
        handler.executeApiRequest(connection, context);
    }

    private void executeApiResponse(TestContext context) {
        Element responseElement = payload.elements().isEmpty() ? payload : (Element) payload.elements().get(0);
        ApiResponseHandler handler = new ApiResponseHandler(responseElement);
        handler.executeApiResponse(context);
    }

    private void executeApiEvent(TestContext context) {
        Element eventElement = payload.elements().isEmpty() ? payload : (Element) payload.elements().get(0);
        ApiEventHandler handler = new ApiEventHandler(eventElement);
        handler.executeApiEvent(context);
    }

    private void executeDelay() {
        String delayStr = payload.getTextTrim();
        if (delayStr.isEmpty()) {
            Delay.mediumDelay();
        } else {
            try {
                long delayMs = Long.parseLong(delayStr);
                Delay.delay(delayMs);
            } catch (NumberFormatException e) {
                throw new DriverTestsuiteException("Invalid delay value: " + delayStr);
            }
        }
    }

    private void executeTerminate() {
        LOGGER.info("Terminating test execution");
        // This step signals the test to stop
    }

    @Override
    public Optional<Location> getLocation() {
        return Optional.ofNullable(location);
    }

    public StepType getType() {
        return type;
    }

    public String getName() {
        return name;
    }

    public Element getPayload() {
        return payload;
    }

    public String getParserArguments() {
        return parserArguments;
    }
}
