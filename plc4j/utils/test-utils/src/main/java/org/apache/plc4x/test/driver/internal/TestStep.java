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
package org.apache.plc4x.test.driver.internal;

import io.netty.channel.embedded.Plc4xEmbeddedChannel;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.test.dom4j.LocationAwareElement;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.test.driver.internal.handlers.*;
import org.apache.plc4x.test.driver.internal.utils.Delay;
import org.apache.plc4x.test.driver.internal.utils.Synchronizer;
import org.apache.plc4x.test.model.Location;
import org.apache.plc4x.test.model.LocationAware;
import org.dom4j.Element;
import org.dom4j.QName;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

public class TestStep implements LocationAware {

    private static final Logger LOGGER = LoggerFactory.getLogger(TestStep.class);

    private final StepType type;
    private final String name;
    private final Location location;

    private final ApiRequestHandler apiRequestHandler;
    private final ApiResponseHandler apiResponseHandler;
    private final IncomingPlcBytesHandler incomingPlcBytesHandler;
    private final IncomingPlcMessageHandler incomingPlcMessageHandler;
    private final OutgoingPlcBytesHandler outgoingPlcBytesHandler;
    private final OutgoingPlcMessageHandler outgoingPlcMessageHandler;

    private TestStep(StepType type, String name, Location location, List<String> parserArguments, Element payload, Synchronizer synchronizer, DriverTestsuiteConfiguration driverTestsuiteConfiguration) {
        this.type = type;
        this.name = name;
        this.location = location;
        apiRequestHandler = new ApiRequestHandler(payload, synchronizer);
        apiResponseHandler = new ApiResponseHandler(payload, synchronizer);
        incomingPlcBytesHandler = new IncomingPlcBytesHandler();
        incomingPlcMessageHandler = new IncomingPlcMessageHandler(driverTestsuiteConfiguration, payload, parserArguments);
        outgoingPlcBytesHandler = new OutgoingPlcBytesHandler(payload);
        outgoingPlcMessageHandler = new OutgoingPlcMessageHandler(driverTestsuiteConfiguration, payload, parserArguments);
    }

    public static TestStep parseTestStep(Element curElement, Synchronizer synchronizer, DriverTestsuiteConfiguration driverTestsuiteConfiguration) throws DriverTestsuiteException {
        final String elementName = curElement.getName();
        final StepType stepType = StepType.valueOf(elementName.toUpperCase().replace("-", "_"));
        final String stepName = curElement.attributeValue(new QName("name"));
        Element parserArgumentsNode = null;
        Element definitionNode = null;
        for (Element element : curElement.elements()) {
            if (element.getName().equals("parser-arguments")) {
                parserArgumentsNode = element;
            } else if (definitionNode == null) {
                definitionNode = element;
            } else {
                throw new DriverTestsuiteException("Error processing the xml. Only one content node allowed.");
            }
        }
        final List<String> parserArguments = new ArrayList<>();
        if (parserArgumentsNode != null) {
            for (Element parserArgumentNode : parserArgumentsNode.elements()) {
                parserArguments.add(parserArgumentNode.getTextTrim());
            }
        }
        Location location = null;
        if (curElement instanceof LocationAwareElement) {
            location = ((LocationAwareElement) curElement).getLocation();
        }
        return new TestStep(stepType, stepName, location, parserArguments, definitionNode, synchronizer, driverTestsuiteConfiguration);
    }

    @Override
    public Optional<Location> getLocation() {
        return Optional.ofNullable(location);
    }

    public void execute(PlcConnection plcConnection, Plc4xEmbeddedChannel embeddedChannel, ByteOrder byteOrder) throws DriverTestsuiteException {
        assert type != null;
        LOGGER.info(String.format("  - Running step: '%s' - %s", name, type));
        try {
            switch (type) {
                case OUTGOING_PLC_BYTES:
                    outgoingPlcBytesHandler.executeOutgoingPlcBytes(embeddedChannel, byteOrder);
                    break;
                case OUTGOING_PLC_MESSAGE:
                    outgoingPlcMessageHandler.executeOutgoingPlcMessage(embeddedChannel, byteOrder);
                    break;
                case INCOMING_PLC_BYTES:
                    incomingPlcBytesHandler.executeIncomingPlcBytes();
                    break;
                case INCOMING_PLC_MESSAGE:
                    incomingPlcMessageHandler.executeIncomingPlcMessage(embeddedChannel, byteOrder);
                    break;
                case API_REQUEST:
                    apiRequestHandler.executeApiRequest(plcConnection);
                    break;
                case API_RESPONSE:
                    apiResponseHandler.executeApiResponse();
                    break;
                case DELAY:
                    Delay.delay(1000);
                    break;
                case TERMINATE:
                    embeddedChannel.close();
                    break;
                default:
                    throw new DriverTestsuiteException("Unknown step type" + type);
            }
        } catch (Exception e) {
            LOGGER.error("    Failed: Error running step: {}: {}", name, e.getMessage());
            throw new DriverTestsuiteException("Error running the step " + name, e);
        }
        LOGGER.info("    Done");
    }

}
