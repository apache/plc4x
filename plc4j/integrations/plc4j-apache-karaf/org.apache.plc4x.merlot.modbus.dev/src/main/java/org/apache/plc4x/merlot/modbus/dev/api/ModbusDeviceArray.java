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
package org.apache.plc4x.merlot.modbus.dev.api;

import java.util.List;
import org.osgi.framework.ServiceReference;

public interface ModbusDeviceArray {
	
	public void init();
	  
	public void destroy();	

	public void bind(ServiceReference<?> reference);
	
	public void bind(ModbusDevice reference);	
	
	public void unbind(ServiceReference<?> reference);
	
	public ModbusDevice getModbusDevice(int device);
	
	public void setModbusDevicesArray(ModbusDevice[] arrModbusDevices);
	
	public ModbusDevice[] getModbusDevicesArray();
	
	public void setModbusDeviceList(List<?> list);
	
	public List<?> getModbusDeviceList();	
	
	
}
