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

package org.apache.plc4x.java.tools.ui.service;

import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.metadata.PlcDriverMetadata;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.configuration.annotations.ComplexConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.Required;
import org.apache.plc4x.java.tools.ui.model.ConfigurationOption;
import org.apache.plc4x.java.tools.ui.model.Device;
import org.apache.plc4x.java.tools.ui.model.Driver;
import org.springframework.stereotype.Component;

import java.util.*;
import java.util.function.Function;
import java.util.stream.Collectors;

@Component
public class DriverService {

    private static final String ALL_DRIVERS = "all";

    private final PlcDriverManager driverManager;
    private final DeviceService deviceService;

    public DriverService(PlcDriverManager driverManager, DeviceService deviceService) {
        this.driverManager = driverManager;
        this.deviceService = deviceService;
    }

    public List<Driver> getDriverList() {
        List<Driver> drivers = new ArrayList<>();
        for (String protocolCode : driverManager.listProtocolCodes()) {
            try {
                PlcDriver driver = driverManager.getDriver(protocolCode);
                PlcDriverMetadata metadata = driver.getMetadata();

                // Get a description of all supported configuration options of the given driver.
                Class<?> configurationType = driver.getConfigurationType();
                Map<String, ConfigurationOption> configurationOptions = Arrays.stream(FieldUtils.getAllFields(configurationType))
                    .filter(field -> (field.getAnnotation(ConfigurationParameter.class) != null) || (field.getAnnotation(ComplexConfigurationParameter.class) != null))
                    .map(field -> new ConfigurationOption(field.getName(), field.getType().getTypeName(), field.isAnnotationPresent(Required.class), ConfigurationFactory.getDefaultValueFromAnnotation(field)))
                    .collect(Collectors.toMap(
                        ConfigurationOption::getName,
                        Function.identity()
                    ));

                // TODO: Get a list of all directly supported transports and for each a list of the configuration options.

                drivers.add(new Driver(protocolCode, driver.getProtocolName(), metadata.canDiscover(), configurationOptions, null));
            } catch (Exception e) {
                throw new RuntimeException("Error retrieving driver list", e);
            }
        }
        return drivers;
    }

    public void discover(String protocolCode) {
        if(ALL_DRIVERS.equals(protocolCode)) {
            for (String curProtocolCode : driverManager.listProtocolCodes()) {
                try {
                    if("modbus-tcp".equals(curProtocolCode)) {
                        continue;
                    }
                    PlcDriver driver = driverManager.getDriver(curProtocolCode);
                    if (driver.getMetadata().canDiscover()) {
                        discoverProtocol(curProtocolCode);
                    }
                } catch (PlcConnectionException e) {
                    throw new RuntimeException(e);
                }
            }
        } else {
            discoverProtocol(protocolCode);
        }
    }

    private void discoverProtocol(String protocolCode) {
        try {
            PlcDriver driver = driverManager.getDriver(protocolCode);
            if (!driver.getMetadata().canDiscover()) {
                throw new RuntimeException("Driver doesn't support discovery");
            } else {
                PlcDiscoveryRequest request = driver.discoveryRequestBuilder().addQuery("all", "*").build();
                // Execute the discovery request and have all connections found be added as connections.
                request.executeWithHandler(discoveryItem -> {
                    // Create the new device.
                    Device device = new Device();
                    device.setName(discoveryItem.getName());
                    device.setProtocolCode(discoveryItem.getProtocolCode());
                    device.setTransportCode(discoveryItem.getTransportCode());
                    device.setTransportUrl(discoveryItem.getTransportUrl());
                    device.setOptions(discoveryItem.getOptions());
                    Map<String, String> attributes = new HashMap<>();
                    for (String attributeName : discoveryItem.getAttributes().keySet()) {
                        String attributeValue = discoveryItem.getAttributes().get(attributeName).getString();
                        attributes.put(attributeName, attributeValue);
                    }
                    device.setAttributes(attributes);

                    // Save the found device in the database, if this is a new device,
                    // that is not stored in our system before.
                    if(deviceService.isNewDevice(device)) {
                        deviceService.createDevice(device);
                    }
                }).whenComplete((plcDiscoveryResponse, throwable) -> {
                    if(throwable != null) {
                        throw new RuntimeException("Error executing discovery", throwable);
                    }
                });
            }
        } catch (PlcConnectionException e) {
            throw new RuntimeException("Error getting driver", e);
        }
    }

}
