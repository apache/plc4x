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
package org.apache.plc4x.merlot.das.ref.impl;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.apache.plc4x.java.api.PlcDriver;
import org.osgi.framework.BundleContext;
import org.osgi.framework.ServiceReference;
import org.osgi.service.device.Constants;
import org.apache.plc4x.merlot.api.PlcDevice;
import org.osgi.service.device.Driver;
import org.osgi.service.event.EventAdmin;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class Plc4xReferringDriverImpl implements Driver {
    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4xReferringDriverImpl.class);    
    private BundleContext bc;  
    private final EventAdmin eventAdmin;    
    private  Map<String, PlcDriver> drivers = new HashMap<String, PlcDriver>();

    public Plc4xReferringDriverImpl(BundleContext bc, EventAdmin eventAdmin) {     
        this.bc = bc;
        this.eventAdmin = eventAdmin;
    }

    /*
    * It verifies that there is a PlcDriver registered in the context and
    * compares it with the one required by the device.
    */
    @Override
    public int match(ServiceReference<?> reference) throws Exception {
        if (!drivers.isEmpty()) {
            final String drv_code = (String) reference.getProperty(Constants.DEVICE_CATEGORY);           
            final PlcDriver plcdriver = drivers.get(drv_code);
            if (null == plcdriver) {
                final PlcDevice plcdev =  (PlcDevice) bc.getService(reference);
                plcdev.noDriverFound();
            } else {
                return 1;
            }
        }
        return org.osgi.service.device.Device.MATCH_NONE; // =0
    }

    /*
    * Attach the PlcDriver to PlcDevice.
    */
    @Override
    public String attach(ServiceReference<?> reference) throws Exception {
        final String drv_code = (String) reference.getProperty(Constants.DEVICE_CATEGORY);
        final PlcDriver plcdriver = drivers.get(drv_code);
        final PlcDevice plcdev =  (PlcDevice) bc.getService(reference);
        plcdev.attach(plcdriver);
        
        return null;
    }
    
   public void bind(ServiceReference reference) {
       //
   }    
    
    public void bind(PlcDriver driver){
        if (null != driver) {
            drivers.put(driver.getProtocolCode(), driver);
        }      
    }
       
    public void unbind(ServiceReference driver){        
        PlcDriver drv = (PlcDriver) bc.getService(driver);
        drivers.remove(drv.getProtocolCode());       
    }    
    
    public void setPlcDriverList(List<PlcDriver> drivers) {
        //
    }
        
}
