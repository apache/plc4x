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
package org.apache.plc4x.java.examples.connectivity.mqtt.model;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

public class PlcConfig {

    private String connection;
    @JsonProperty("memory-blocks")
    private List<PlcMemoryBlock> plcMemoryBlocks;
    @JsonProperty("addresses")
    private List<PlcFieldConfig> plcFields;

    public String getConnection() {
        return connection;
    }

    public void setConnection(String connection) {
        this.connection = connection;
    }

    public List<PlcMemoryBlock> getPlcMemoryBlocks() {
        return plcMemoryBlocks;
    }

    public void setPlcMemoryBlocks(List<PlcMemoryBlock> plcMemoryBlocks) {
        this.plcMemoryBlocks = plcMemoryBlocks;
    }

    public List<PlcFieldConfig> getPlcFields() {
        return plcFields;
    }

    public void setPlcFields(List<PlcFieldConfig> plcFields) {
        this.plcFields = plcFields;
    }

}
