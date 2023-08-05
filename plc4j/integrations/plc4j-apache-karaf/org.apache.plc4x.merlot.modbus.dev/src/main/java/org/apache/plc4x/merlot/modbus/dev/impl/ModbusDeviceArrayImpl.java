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
package org.apache.plc4x.merlot.modbus.dev.impl;

import org.apache.plc4x.merlot.modbus.dev.api.ModbusDevice;
import org.apache.plc4x.merlot.modbus.dev.api.ModbusDeviceArray;
import java.util.List;
import org.osgi.framework.BundleContext;
import org.osgi.framework.ServiceReference;

public class ModbusDeviceArrayImpl implements ModbusDeviceArray {

    private BundleContext bundleContext;
    private List<?> ModbusDeviceList;

    private ModbusDevice[] arrModbusDevices = new ModbusDevice[255];

    public ModbusDeviceArrayImpl() {
        for (int i = 0; i < 255; i++) {
            arrModbusDevices[i] = null;
        }
    }

    public BundleContext getBundleContext() {
        return bundleContext;
    }

    public void setBundleContext(BundleContext bundleContext) {
        this.bundleContext = bundleContext;
    }

    public void init() {

    }

    public void destroy() {

    }

    public void bind(ServiceReference<?> reference) {
        //ModbusDevice mbd = (ModbusDevice) bundleContext.getService(reference);
        //System.out.println("ServiceReference  reference: " + mbd.getUnitIdentifier());		
    }

    public void bind(ModbusDevice reference) {
        arrModbusDevices[reference.getUnitIdentifier()] = reference;
    }

    public void unbind(ServiceReference<?> reference) {
    }

    public void unbind(ModbusDevice reference) {
        if (reference != null) {
            arrModbusDevices[reference.getUnitIdentifier()] = null;
        };
    }

    public ModbusDevice getModbusDevice(int device) {
        if (device >= 0 && device <= 255) {
            return arrModbusDevices[device];
        } else {
            return null;
        }
    };

    public void setModbusDevicesArray(ModbusDevice[] arrModbusDevices) {
        this.arrModbusDevices = arrModbusDevices;
    }

    public ModbusDevice[] getModbusDevicesArray() {
        return this.arrModbusDevices;
    }

    public void setModbusDeviceList(List<?> list) {
        this.ModbusDeviceList = list;
    }

    public List<?> getModbusDeviceList() {
        return this.ModbusDeviceList;
    }

}
