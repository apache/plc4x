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

package org.apache.plc4x.java.tools.ui.controller;

import org.apache.plc4x.java.tools.ui.model.Driver;
import org.apache.plc4x.java.tools.ui.service.DriverService;
import org.springframework.web.bind.annotation.*;

import java.util.List;

// Allow from the default port 8080 as well as the one node usually uses for it's dev-mode 5173
//@CrossOrigin(origins = {"http://localhost:8080", "http://localhost:5173"})
@CrossOrigin(origins = {"*"})
@RestController
@RequestMapping("/api")
public class DriverController {

    private final DriverService driverService;

    public DriverController(DriverService driverService) {
        this.driverService = driverService;
    }

    @GetMapping("/drivers")
    public List<Driver> getAllDrivers() {
        return driverService.getDriverList();
    }

    @GetMapping("/discover/{protocolCode}")
    public void discover(@PathVariable("protocolCode") String protocolCode) {
        driverService.discover(protocolCode);
    }

}
