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
import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Argument;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.osgi.framework.BundleContext;
import org.osgi.framework.ServiceReference;

@Command(scope = "modbus", name = "dump", description = "List modbus device")
@Service
public class ModbusDeviceDumpCommand implements Action {

    @Reference
    BundleContext bundleContext;
    
    @Argument(index = 0, name = "uid", description = "The device unit identifier.", required = true, multiValued = false)
    int uid;
    
    //@Option(name = "-r", aliases = { "--reg" }, description = "The data bank to dump.", required = true, multiValued = false)
    @Argument(index = 1, name = "type", description = "1:Digital register, 2:Coil, 3:Input register, 4:Holding register.", required = true, multiValued = false)    
    int type;
    
    //@Option(name = "-s", aliases = { "--start" }, description = "The initial register.", required = true, multiValued = false)
    @Argument(index = 2, name = "start", description = "The start register address.", required = true, multiValued = false)     
    int start;

    //@Option(name = "-l", aliases = { "--length" }, description = "Number of registers to dump.", required = true, multiValued = false)
    @Argument(index = 3, name = "length", description = "Number of registers to dump.", required = true, multiValued = false)     
    int length;    
    
    String strHeader = null;
    		
    ByteBuf registers = null;

	@SuppressWarnings("unused")
	private ModbusDevice mdb;
    
	public Object execute() throws Exception {
        try {
            //Bundle bundle = bundleContext.getSe
        	ServiceReference<?> reference = bundleContext.getServiceReference(ModbusDeviceArray.class.getName());
        	ModbusDeviceArray mdbarray = (ModbusDeviceArray) bundleContext.getService(reference);
        	mdb = (ModbusDevice) mdbarray.getModbusDeviceList().get(0);
        	
        	switch(type) {
        		case 1: registers = mdbarray.getModbusDevicesArray()[uid].getDiscreteInputs();
        				strHeader = "Discrete Inputs";
        				break;
        		case 2: registers = mdbarray.getModbusDevicesArray()[uid].getCoils();
						strHeader = "Coils"; 					
        				break;
        		case 3: registers = mdbarray.getModbusDevicesArray()[uid].getInputRegisters();
						strHeader = "Input Registers";           		
        				break;
        		case 4: registers = mdbarray.getModbusDevicesArray()[uid].getHoldingRegisters();
						strHeader = "Holding Registers";         		
        				break;
        	}
        	
        	if (registers != null){  

        		System.out.println("\r\nDevice UID: "+uid+", Dump of \"" + strHeader + "\", from: "+start+", length: "+length+ " registers");
        		if ((type==1) || (type==2)) {
    				length = length / 8;
    				if (length==0) {length = 1;}
        		}
        		String str = ByteBufUtil.prettyHexDump (registers,start,length);

        		System.out.println(str+"\r\n");
        	}
        } catch (NumberFormatException ex) {
            // It was not a number, so ignore.
        }

		return null;
	}

}
