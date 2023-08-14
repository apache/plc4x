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
package org.apache.plc4x.merlot.modbus.dev.core;

import org.apache.plc4x.merlot.modbus.dev.api.ModbusDevice;
import org.apache.plc4x.merlot.modbus.dev.impl.ModbusDeviceImpl;
import java.util.Dictionary;
import java.util.Enumeration;
import java.util.Hashtable;
import org.osgi.framework.BundleContext;
import org.osgi.framework.Constants;
import org.osgi.framework.InvalidSyntaxException;
import org.osgi.framework.ServiceReference;
import org.osgi.service.cm.ConfigurationException;
import org.osgi.service.cm.ManagedService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class ModbusDeviceManagedService implements ManagedService {

    private static final Logger LOGGER = LoggerFactory.getLogger(ModbusDeviceManagedService.class);
    private static final String MODBUS_DEVICE_ID = "modbus.id";
    private static final String MODBUS_DEVICE_UID = "modbus.uid";

    private final BundleContext bundleContext;
    private String filter_device = "(&(" + Constants.OBJECTCLASS + "=" + ModbusDevice.class.getName() + ")"
            + "(|(" + MODBUS_DEVICE_ID + "=*)(" + MODBUS_DEVICE_UID + "=$)))";
        
    public ModbusDeviceManagedService(BundleContext bundleContext) {
        this.bundleContext = bundleContext;
    }

    @Override
    public void updated(Dictionary props) throws ConfigurationException {
        if (null == props)  return;
        
        Enumeration<String> keys = props.keys();
        ModbusDevice mbDevice = null;
        while (keys.hasMoreElements()) {
            String key = keys.nextElement();
            String[] fields = ((String) props.get(key)).split(",");
            mbDevice = null;
            try {
                String deviceFilter = filter_device.replace("*", key);
                deviceFilter = deviceFilter.replace("$", fields[0]);

                ServiceReference[] references = bundleContext.getServiceReferences((String) null, deviceFilter);
                if (references != null) System.out.println("Referencias: " + references.length);

                if ((references == null) 
                        && (fields.length == 6) 
                        && !key.equalsIgnoreCase("felix.fileinstall.filename")
                        && !key.equalsIgnoreCase("service.pid")) {                                                                                                    
                    int uid = Integer.parseInt(fields[0]);
                    int discrete_inputs = Integer.parseInt(fields[1]);
                    int coils = Integer.parseInt(fields[2]);
                    int input_registers = Integer.parseInt(fields[3]);
                    int holding_registers = Integer.parseInt(fields[4]);
                    String short_description = fields[5];
                    if ((discrete_inputs == 0) 
                            && (coils == 0) 
                            && (input_registers == 0) 
                            && (holding_registers != 0)) {
                        mbDevice = (ModbusDevice) new ModbusDeviceImpl(holding_registers);
                    } else if ((discrete_inputs > 0) 
                            && (coils > 0) 
                            && (input_registers > 0) 
                            && (holding_registers > 0)) {
                        mbDevice = (ModbusDevice) new ModbusDeviceImpl(discrete_inputs
                                ,coils
                                ,input_registers
                                ,holding_registers);                        
                    }
                    mbDevice.setUnitIdentifier((byte) uid);
                    mbDevice.setUnitDescription(short_description);
                    
                    Hashtable properties = new Hashtable();
                    properties.put(MODBUS_DEVICE_ID, key);                            
                    properties.put(MODBUS_DEVICE_UID, mbDevice.getUnitIdentifier());
                    bundleContext.registerService(ModbusDevice.class.getName()
                            ,mbDevice
                            ,properties);
                    LOGGER.info("Registered ModbusDevice [{}] [{}] [{}]"
                            ,key
                            ,uid
                            ,short_description);
                    
                } else {
                    LOGGER.info("Fail registering device [" + key + "].");
                }
            } catch (InvalidSyntaxException ex) {
                LOGGER.info(ex.toString());
            }
        }

    }

}
