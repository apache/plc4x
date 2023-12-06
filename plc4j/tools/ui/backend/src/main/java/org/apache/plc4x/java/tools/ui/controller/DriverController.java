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

import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.metadata.PlcDriverMetadata;
import org.apache.plc4x.java.tools.ui.model.Driver;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.ArrayList;
import java.util.Collections;
import java.util.Comparator;
import java.util.List;

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
                drivers.add(new Driver(protocolCode, driver.getProtocolName(), metadata.canDiscover()));
            } catch (PlcConnectionException e) {
                // Ignore ...
            }
        }

        // Sort the list by the code of the driver elements.
        drivers.sort(Comparator.comparing(Driver::getCode));

        return ResponseEntity.ok(drivers);
    }

}
