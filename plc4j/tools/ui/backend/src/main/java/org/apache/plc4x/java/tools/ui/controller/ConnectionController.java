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

import org.apache.plc4x.java.tools.ui.model.Device;
import org.apache.plc4x.java.tools.ui.service.ConnectionService;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/connections")
public class ConnectionController {

    private final ConnectionService connectionService;

    public ConnectionController(ConnectionService connectionService) {
        this.connectionService = connectionService;
    }

    @GetMapping("/{id}")
    public Device getConnectionById(@PathVariable Integer id) {
        return connectionService.readConnection(id).orElseThrow(() -> new RuntimeException("Error finding connection with id: " + id));
    }

    @GetMapping
    public List<Device> getAllConnections() {
        return connectionService.getAllConnections();
    }

    @PostMapping
    public Device saveConnection(@RequestBody Device device) {
        if (device.getId() == null) {
            return connectionService.createConnection(device);
        } else {
            return connectionService.updateConnection(device);
        }
    }

    @DeleteMapping
    public void deleteConnection(@RequestBody Device device) {
        connectionService.deleteConnection(device);
    }

}
