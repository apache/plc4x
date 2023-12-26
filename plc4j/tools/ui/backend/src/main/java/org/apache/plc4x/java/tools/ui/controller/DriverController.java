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

package org.apache.plc4x.java.tools.ui.controller;

import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.metadata.PlcDriverMetadata;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.configuration.annotations.ComplexConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.Required;
import org.apache.plc4x.java.tools.ui.model.ConfigurationOption;
import org.apache.plc4x.java.tools.ui.model.Driver;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.lang.reflect.Field;
import java.util.*;
import java.util.function.Function;
import java.util.stream.Collectors;

// Allow from the default port 8080 as well as the one node usually uses for it's dev-mode 5173
@CrossOrigin(origins = {"http://localhost:8080", "http://localhost:5173"})
@RestController
@RequestMapping("/api")
public class DriverController {

    @GetMapping("/drivers")
    public ResponseEntity<List<Driver>> getDriverList() {
        List<Driver> drivers = new ArrayList<>();

        // Build a list of driver objects.
        PlcDriverManager driverManager = PlcDriverManager.getDefault();
        for (String protocolCode : driverManager.listDrivers()) {
            try {
                PlcDriver driver = driverManager.getDriver(protocolCode);
                PlcDriverMetadata metadata = driver.getMetadata();

                // Get a description of all supported configuration options of the given driver.
                Class<?> configurationType = driver.getConfigurationType();
                Map<String, ConfigurationOption> configurationOptions = Arrays.stream(FieldUtils.getAllFields(configurationType))
                    // - Filter out only the ones annotated with the ConfigurationParameter annotation.
                    .filter(field -> (field.getAnnotation(ConfigurationParameter.class) != null) || (field.getAnnotation(ComplexConfigurationParameter.class) != null))
                    .map(field -> new ConfigurationOption(field.getName(), field.getType().getTypeName(), field.isAnnotationPresent(Required.class), ConfigurationFactory.getDefaultValueFromAnnotation(field)))
                    // - Create a map with the field-name as key and the field itself as value.
                    .collect(Collectors.toMap(
                        ConfigurationOption::getName,
                        Function.identity()
                    ));

                // TODO: Get a list of all directly supported transports and for each a list of the configuration options.

                drivers.add(new Driver(protocolCode, driver.getProtocolName(), metadata.canDiscover(), configurationOptions, null));
            } catch (Exception e) {
                e.printStackTrace();
                // Ignore ...
            }
        }

        // Sort the list by the code of the driver elements.
        drivers.sort(Comparator.comparing(Driver::getCode));

        return ResponseEntity.ok(drivers);
    }

}
