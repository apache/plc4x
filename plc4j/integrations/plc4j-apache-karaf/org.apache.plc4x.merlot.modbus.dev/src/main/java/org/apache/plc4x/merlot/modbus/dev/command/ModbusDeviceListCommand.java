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
import java.util.List;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.karaf.shell.support.table.ShellTable;
import org.osgi.framework.BundleContext;
import org.osgi.framework.ServiceReference;

@Command(scope = "modbus", name = "list", description = "List modbus device")
@Service
public class ModbusDeviceListCommand implements Action {
	
    @Reference
    BundleContext bundleContext;
    
    private ModbusDevice mbd = null;
    
	public Object execute() throws Exception {
        ShellTable table = new ShellTable();
        table.column("UID");
        table.column("Description");
        table.column("Enabled");        
        table.column("Discrete inputs");
        table.column("Coils");
        table.column("Input Registers");
        table.column("Holding Registers");        
        try {
            //Bundle bundle = bundleContext.getSe
        	ServiceReference<?> reference = bundleContext.getServiceReference(ModbusDeviceArray.class.getName());
        	ModbusDeviceArray mdbarray = (ModbusDeviceArray) bundleContext.getService(reference);
        	
        	List<?> mbdList = mdbarray.getModbusDeviceList();        	
        	
        	for (int i=0; i < mbdList.size(); i++) {
        		mbd = (ModbusDevice)  mbdList.get(i);
        		table.addRow().addContent(	
        				mbd.getUnitIdentifier(),
        				mbd.getUnitDescription(),
        				mbd.getEnabled(),
        				mbd.getDiscreteInputs().capacity() * 8,
        				mbd.getCoils().capacity() * 8,
        				mbd.getInputRegisters().capacity(),
        				mbd.getHoldingRegisters().capacity());        		
        	}        	
        	
        } catch (NumberFormatException ex) {
            // It was not a number, so ignore.
        }
        System.out.println();
        table.print(System.out);
        System.out.println();        
		return null;
	}

}
