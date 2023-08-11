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
package org.apache.plc4x.merlot.das.base.impl;

import org.apache.plc4x.merlot.das.base.api.BaseDevice;
import org.apache.plc4x.merlot.das.base.api.BaseDriver;
import static org.apache.plc4x.merlot.das.base.impl.BaseDriverImpl.BASIC_DEVICE_CATEGORY;
import java.util.Dictionary;
import java.util.Hashtable;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.osgi.framework.BundleContext;
import org.osgi.framework.ServiceReference;
import org.osgi.service.device.Constants;
import org.osgi.service.device.Device;
import org.osgi.service.device.Driver;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class BaseDasDriverImpl implements Driver {

    private static final Logger LOGGER = LoggerFactory.getLogger(BaseDasDriverImpl.class);    
    
    @Reference
    private BundleContext bc;
    
    private BaseDevice bdev;
    
    /**
    *   Try to locate a driver for the device.
    *   If there is a driver associated with the device, return with 1.
    *   If the driver does not exist, create an instance and register it. 
    *   If everything goes well, return with 1.
    *   If the driver does not exist, return with Device.MATCH_NONE.
    * 
    * @param reference The reference to the Device
    *  @return Device.MATCH_NONE, 1
    */
    @Override
    public int match(ServiceReference reference) throws Exception {
        String driver_url;
        if (reference != null) {
            String deviceCategory = (String) reference.getProperty(Constants.DEVICE_CATEGORY);
            if (deviceCategory.equals(BASIC_DEVICE_CATEGORY)) {             
                {
                    ServiceReference[] services = bc.getServiceReferences(
                            BaseDriver.class.getName(), null);
                    if (services.length > 0){
                        for(ServiceReference service:services){
                            driver_url = (String) service.getProperty(Constants.DRIVER_ID);
                        }
                        
                    } else {
                        
                        final Dictionary<String, Object> props = new Hashtable<>();
                        props.put ( Constants.DRIVER_ID, "com.ceos.basic" );
                        props.put ( "DRIVER_URL", "Eclipse SCADA Project" );
                        
                        BaseDriver driver = new BaseDriverImpl(bc);
                        
                        bc.registerService(BaseDriver.class.getName(), driver, props);
                        
                        return 1;
                    }

                }
                //Busco el driver
                       	
                    return 1;
            } 
          }
          return Device.MATCH_NONE;
    }
 

    @Override
    public String attach(ServiceReference reference) throws Exception {
        System.out.println("Paso por aqui... ATTACH");  
        if (reference != null) {
            System.out.println("Paso por aqui... PASO1");  
            bdev = (BaseDevice) bc.getService(reference);
            System.out.println("Paso por aqui... PASO2");
            LOGGER.info("attach: " + bdev.toString());
        }
        return null;
    }
    
}
