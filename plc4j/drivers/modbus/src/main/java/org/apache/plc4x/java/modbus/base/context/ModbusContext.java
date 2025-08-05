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

package org.apache.plc4x.java.modbus.base.context;

import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
import org.apache.plc4x.java.spi.context.DriverContext;

public abstract class ModbusContext implements DriverContext {

    private ModbusByteOrder byteOrder;
    private int maxCoilsPerRequest;
    private int maxRegistersPerRequest;

    public ModbusByteOrder getByteOrder() {
        return byteOrder;
    }

    public void setByteOrder(ModbusByteOrder byteOrder) {
        this.byteOrder = byteOrder;
    }

    public int getMaxCoilsPerRequest() {
        return maxCoilsPerRequest;
    }

    public void setMaxCoilsPerRequest(int maxCoilsPerRequest) {
        this.maxCoilsPerRequest = maxCoilsPerRequest;
    }

    public int getMaxRegistersPerRequest() {
        return maxRegistersPerRequest;
    }

    public void setMaxRegistersPerRequest(int maxRegistersPerRequest) {
        this.maxRegistersPerRequest = maxRegistersPerRequest;
    }

}
