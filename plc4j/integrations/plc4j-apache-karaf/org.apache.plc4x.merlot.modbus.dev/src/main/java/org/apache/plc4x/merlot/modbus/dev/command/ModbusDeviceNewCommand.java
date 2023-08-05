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
package org.apache.plc4x.merlot.modbus.dev.command;


import org.apache.plc4x.merlot.modbus.dev.api.ModbusDevice;
import org.apache.plc4x.merlot.modbus.dev.api.ModbusDeviceArray;
import java.util.Hashtable;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Argument;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.plc4x.merlot.modbus.dev.impl.ModbusDeviceImpl;
import org.osgi.framework.BundleContext;
import org.osgi.framework.Constants;
import org.osgi.framework.ServiceReference;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Command(scope = "modbus", name = "new", description = "Create a new modbus device.")
@Service
public class ModbusDeviceNewCommand implements Action {
    private static final Logger LOGGER = LoggerFactory.getLogger(ModbusDeviceNewCommand.class);
    private static final String MODBUS_DEVICE_ID = "modbus.id";
    private static final String MODBUS_DEVICE_DESC = "modbus.desc";    
    private static final String MODBUS_DEVICE_UID = "modbus.uid";
    private String filter_device = "(&(" + Constants.OBJECTCLASS + "=" + ModbusDevice.class.getName() + ")"
            + "(|(" + MODBUS_DEVICE_ID + "=*)(" + MODBUS_DEVICE_UID + "=$)))";    
    
    @Reference
    BundleContext bundleContext;
        
    @Argument(index = 0, name = "uid", description = "The device modbus uid.", required = true, multiValued = false)            
    int uid;     
    
    @Argument(index = 1, name = "id", description = "A technological name for the instance.", required = true, multiValued = false)            
    String id;      
    
    @Argument(index = 2, name = "desc", description = "The device short description.", required = true, multiValued = false)            
    String short_description;      
    
    @Argument(index = 3, name = "di", description = "Number of discrete inputs. Read only bobines.", required = true, multiValued = false)            
    int discrete_inputs;

    @Argument(index = 4, name = "coils", description = "Number of coils. Read/Write bobines.", required = true, multiValued = false)            
    int coils;

    @Argument(index = 5, name = "ir", description = "Number of input registers. Read only registers.", required = true, multiValued = false)            
    int input_registers;

    @Argument(index = 6, name = "hr", description = "Number of holding registers. Read/Write registers.", required = true, multiValued = false)            
    int holding_registers;

    
    
    private ModbusDevice mbd = null;
    
    public Object execute() throws Exception {
        ModbusDevice mbDevice = null;        
        try {
            String deviceFilter = filter_device.replace("*", id);
            deviceFilter = deviceFilter.replace("$", String.valueOf(uid));            
            
                ServiceReference[] references = bundleContext.getServiceReferences((String) null, deviceFilter);
                if (references != null)System.out.println("ModbusDevice instance exists.: " + references.length);  
                if (references == null) {
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
                    } else {
                        System.out.println("ModbusDevice check the supplied settings.");
                        return null;
                    }
                    mbDevice.setUnitIdentifier((byte) uid);
                    mbDevice.setUnitDescription(id);  
                    
                    Hashtable properties = new Hashtable();
                    properties.put(MODBUS_DEVICE_ID, id);    
                    properties.put(MODBUS_DEVICE_DESC, short_description);                       
                    properties.put(MODBUS_DEVICE_UID, mbDevice.getUnitIdentifier());
                    bundleContext.registerService(ModbusDevice.class.getName()
                            ,mbDevice
                            ,properties);
                    LOGGER.info("Registered ModbusDevice [{}] [{}] [{}]"
                            ,uid
                            ,id
                            ,short_description);                    
                }
                         
            
        } catch (Exception ex){
            LOGGER.info(ex.toString());
        }

        return null;
    }

}
