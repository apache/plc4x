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
package org.apache.plc4x.merlot.das.api;

import org.apache.plc4x.merlot.das.core.DriverAttributes;
import java.util.Map;
import org.osgi.framework.ServiceReference;


public interface DeviceManager {
    
    public void init() throws Exception;
	
    public void destroy() throws Exception;
	
    public void start();
	
    public void stop();
    
    public void driverAdded( ServiceReference ref, Object obj );
    
    public void driverModified( ServiceReference ref, Object obj ); 
    
    public void driverRemoved( ServiceReference ref );

    public void deviceAdded( ServiceReference ref, Object obj );
    
    public void deviceModified( ServiceReference ref, Object obj ); 
    
    public void deviceRemoved( ServiceReference ref );    
    
    public Map<ServiceReference, DriverAttributes> getDrivers();
    
    public Map<ServiceReference, Object> getDevices();
    
}
