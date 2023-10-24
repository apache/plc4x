/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.apache.plc4x.merlot.api.core;

import org.apache.plc4x.merlot.scheduler.api.Job;
import org.apache.plc4x.merlot.scheduler.api.JobContext;
import java.util.Collections;
import java.util.Dictionary;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.Hashtable;
import java.util.Map;
import org.apache.plc4x.java.api.PlcDriver;
import org.osgi.framework.BundleContext;
import org.osgi.framework.Constants;
import org.osgi.framework.InvalidSyntaxException;
import org.osgi.framework.ServiceReference;
import org.osgi.service.cm.ConfigurationEvent;
import org.osgi.service.cm.ConfigurationException;
import org.osgi.service.cm.ConfigurationListener;
import org.osgi.service.cm.ManagedService;
import org.osgi.service.device.Device;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.apache.plc4x.merlot.api.PlcDeviceFactory;
import org.apache.plc4x.merlot.api.PlcDevice;
/**
 *
 * @author cgarcia
 */
public class PlcDeviceManagedService implements ManagedService, ConfigurationListener, Job {
    
    private static final Logger LOGGER = LoggerFactory.getLogger(PlcDeviceManagedService.class);     
    private final BundleContext bc;
    private static Map<String, Dictionary<String, ?>> waitingConfigs = null;  
    private String filter_driver =  "(&(" + Constants.OBJECTCLASS + "=" + PlcDriver.class.getName() + ")" +
                        "(org.apache.plc4x.driver.code=*))";    
    private String filter_device =  "(&(" + Constants.OBJECTCLASS + "=" + Device.class.getName() + ")" +
                        "(" + org.apache.plc4x.merlot.api.PlcDevice.SERVICE_NAME + "=*))";
    private String filter_factory =  "(&(" + Constants.OBJECTCLASS + "=" + PlcDeviceFactory.class.getName() + ")" +
                        "(org.apache.plc4x.device.factory=*))";
    
       
    public PlcDeviceManagedService(BundleContext bc) {
        this.bc = bc;
        waitingConfigs = Collections.synchronizedMap(new HashMap<String, Dictionary<String, ?>>());
    }

    public String getName() {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    /*
    * Example:
    * AS01 = s7://192.168.1.23/0/2,simatic,compresor system
    * \__/   \___________________/ \_____/ \______________/
    * key            url            name      description
    *
    */
    @Override
    public void updated(Dictionary<String, ?> props) throws ConfigurationException {
        if (null == props) return;

//        if (props.size() < 3){
//            waitingConfigs.put(pid, props);  
//            return;
//        }

        if (props!=null) {
            Enumeration<String> keys = props.keys();
            for (Enumeration e = props.keys(); e.hasMoreElements();) {
                Object key = e.nextElement();
                String driver_information = props.get(key).toString();
                String[] drv_data = driver_information.split(":",2);
                if (drv_data.length == 2) {
                    try {  
                        String deviceFilter = filter_device.replace("*",key.toString());
                        ServiceReference[] references = bc.getServiceReferences((String) null, deviceFilter);
                        if (references != null){
                            for (ServiceReference reference:references){
                                LOGGER.info("The device already exists: " + reference.getProperty(org.apache.plc4x.merlot.api.PlcDevice.SERVICE_NAME));                                
                            }
                        } else {
                            String factoryFilter = filter_factory.replace("*",drv_data[0]);                            
                            references = bc.getServiceReferences((String) null, factoryFilter);

                            if (null != references){
                                //1. Si existe el factory, existe el driver.
                                drv_data = driver_information.split(",",3);
                              
                                ServiceReference reference = references[0];
                                PlcDeviceFactory bdf = (PlcDeviceFactory) bc.getService(reference);
                                if (bdf != null){
                                    PlcDevice device =  bdf.create(key.toString(), drv_data[0], drv_data[1], drv_data[2]);
                                    if (device != null) {
                                        
                                        bc.registerService(org.apache.plc4x.merlot.api.PlcDevice.class.getName(), device, device.getProperties());
                                        
                                    } else {
                                       LOGGER.info("Failed to register driver." + factoryFilter);
//                                       waitingConfigs.put(pid, props);  
                                    }
                                } else {
                                    LOGGER.info("The factory is not available:" + drv_data[0]);
//                                    waitingConfigs.put(pid, props);
                                };
                            } else {
                                LOGGER.info("There is no factory specific for the driver [{}], using base device. to register device driver.", factoryFilter);
                                //1. Si existe un driver se puede registrar el dispositivo.                                                              
                                //2. Registra el dispositivo con un BaseDevice
                                String factoryBaseFilter = filter_factory.replace("*","base");    
                            
                                references = bc.getServiceReferences((String) null, factoryBaseFilter);
                                
                                //Filters changes to other configuration files.
                                if ((null != references) && (references.length > 0) && (!drv_data[0].equalsIgnoreCase("file"))) {
                                    ServiceReference reference = references[0];                                         
                                    PlcDeviceFactory bdf = (PlcDeviceFactory) bc.getService(reference);
                                    drv_data = driver_information.split(",",3);
                                    PlcDevice device  =  bdf.create(key.toString() , drv_data[0], drv_data[1], drv_data[2]);

                                    bc.registerService(org.apache.plc4x.merlot.api.PlcDevice.class.getName(), device, device.getProperties());
                                                                       
                                } else {
                                    LOGGER.info("The base factory is not available.");  
                                }
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
    public void execute(JobContext arg0) {      
//        String pid = null;
//        Dictionary<String, ?> props = null;
//        Set<String> keys = new HashSet<>();
//        keys.addAll(waitingConfigs.keySet());
//        for (String key:keys){
//            pid = key;
//            props = waitingConfigs.remove(key);            
//            try {
//                updated(pid,props);
//            } catch (ConfigurationException ex) {
//                LOGGER.debug("Problem updating [" + key +"] from waiting list." );
//            }
//        }
    }    

    @Override
    public void configurationEvent(ConfigurationEvent event) {
        LOGGER.info(">>> Cambio de configuración: " + event.getPid());
    }
    
}
