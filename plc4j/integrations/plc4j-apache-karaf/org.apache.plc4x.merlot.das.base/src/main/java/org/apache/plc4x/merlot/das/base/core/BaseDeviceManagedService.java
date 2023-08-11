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
package org.apache.plc4x.merlot.das.base.core;

import org.apache.plc4x.merlot.das.base.api.BaseDevice;
import org.apache.plc4x.merlot.das.base.api.BaseDeviceFactory;
import org.apache.plc4x.merlot.scheduler.api.Job;
import org.apache.plc4x.merlot.scheduler.api.JobContext;
import java.util.Collections;
import java.util.Dictionary;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import org.osgi.framework.BundleContext;
import org.osgi.framework.Constants;
import org.osgi.framework.InvalidSyntaxException;
import org.osgi.framework.ServiceReference;
import org.osgi.service.cm.ConfigurationException;
import org.osgi.service.cm.ManagedServiceFactory;
import org.osgi.service.device.Device;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class BaseDeviceManagedService implements ManagedServiceFactory, Job {
    
    private static final Logger LOGGER = LoggerFactory.getLogger(BaseDeviceManagedService.class);     
    private final BundleContext bundleContext;
    private static Map<String, Dictionary<String, ?>> waitingConfigs = null;     
    private String filter_device =  "(&(" + Constants.OBJECTCLASS + "=" + Device.class.getName() + ")" +
                        "(" + org.osgi.service.device.Constants.DEVICE_SERIAL + "=*))";
    private String filter_factory =  "(&(" + Constants.OBJECTCLASS + "=" + BaseDeviceFactory.class.getName() + ")" +
                        "(device.factory=*))";
    
       
    public BaseDeviceManagedService(BundleContext bundleContext) {
        this.bundleContext = bundleContext;
        waitingConfigs = Collections.synchronizedMap(new HashMap<String, Dictionary<String, ?>>());
    }

    @Override
    public String getName() {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public void updated(String pid, Dictionary<String, ?> props) throws ConfigurationException {
        if (props.size() < 3){
            waitingConfigs.put(pid, props);  
            return;
        }
        if (props!=null) {
            Enumeration<String> keys = props.keys();
            for (Enumeration e = props.keys(); e.hasMoreElements();) {
                Object key = e.nextElement();
                String driver_information = props.get(key).toString();
                String[] drv_data = driver_information.split(":",2);
                if (drv_data.length == 2) {
                    try {  
                        String deviceFilter = filter_device.replace("*",key.toString());
                        ServiceReference[] references = bundleContext.getServiceReferences((String) null, deviceFilter);
                        if (references != null){
                            for (ServiceReference reference:references){
                                LOGGER.info("The device already exists: " + reference.getProperty("DEVICE_DESCRIPTION"));                                
                            }
                        } else {
                            String factoryFilter = filter_factory.replace("*",drv_data[0]);                            
                            references = bundleContext.getServiceReferences((String) null, factoryFilter);

                            if (references != null){
                                drv_data = driver_information.split(",",3);
                              
                                ServiceReference reference = references[0];
                                BaseDeviceFactory bdf = (BaseDeviceFactory) bundleContext.getService(reference);
                                if (bdf != null){
                                    BaseDevice device = bdf.create(key.toString(), drv_data[0], drv_data[1]);
                                    if (device != null) {
                                        device.start();
                                        LOGGER.info("Starting device["+ key.toString() + "] @ url[ " + drv_data[0]+ "]");  
                                    } else {
                                       LOGGER.info("Failed to register driver." + factoryFilter);
                                       waitingConfigs.put(pid, props);  
                                    }
                                } else {
                                    LOGGER.info("The factory is not available:" + drv_data[0]);
                                    waitingConfigs.put(pid, props);
                                };
                            } else {
                                LOGGER.info("Failed to register device driver." + factoryFilter);
                                waitingConfigs.put(pid, props);  
                            }
                        }
                    } catch (InvalidSyntaxException ex) {
                        LOGGER.info("Exception: " + ex);
                    }
                } else {
                    LOGGER.info("Key not valid: " + key.toString() + " : " + drv_data[0]);
                }
             }
        }
    }

    @Override
    public void deleted(String arg0) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }
    
    @Override
    public void execute(JobContext arg0) {       
        String pid = null;
        Dictionary<String, ?> props = null;
        Set<String> keys = new HashSet<>();
        keys.addAll(waitingConfigs.keySet());
        for (String key:keys){
            pid = key;
            props = waitingConfigs.remove(key);            
            try {
                updated(pid,props);
            } catch (ConfigurationException ex) {
                LOGGER.debug("Problem updating [" + key +"] from waiting list." );
            }
        }
    }    
    
}
