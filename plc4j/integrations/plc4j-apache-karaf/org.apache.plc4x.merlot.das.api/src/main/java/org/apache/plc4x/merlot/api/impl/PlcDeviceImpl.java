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
package org.apache.plc4x.merlot.api.impl;


import java.util.HashMap;
import java.util.Hashtable;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.atomic.AtomicReference;
import static java.util.stream.Collectors.toList;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.merlot.api.PlcDevice;
import org.osgi.framework.BundleContext;
import org.osgi.service.dal.Device;
import org.osgi.service.dal.DeviceException;
import org.osgi.service.device.Constants;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.apache.plc4x.merlot.api.PlcGroup;
import org.apache.plc4x.merlot.api.PlcItem;
import org.apache.plc4x.merlot.scheduler.api.Job;
import org.osgi.framework.ServiceReference;

/*
*
*/
public class PlcDeviceImpl implements PlcDevice {	
     
    private static final Logger LOGGER = LoggerFactory.getLogger(PlcDeviceImpl.class);
    
    private static final String FILTER_DEVICE =  "(&(" + 
            org.osgi.framework.Constants.OBJECTCLASS + 
            "=" + PlcDevice.class.getName() + ")" +
            "(" + PlcDevice.SERVICE_UID + "=*))";   
    
    protected final BundleContext bc;
    protected boolean enable  = false;      
    protected boolean autostart = false;  
    
    protected PlcDriver plcDriver = null;    
    AtomicReference<PlcConnection> refPlcConnection;
    PlcConnection plcConnection = null;
    
    protected Hashtable<String, Object> myProperties;
    
    private final Map<UUID, PlcGroup> device_groups;     
    
    public PlcDeviceImpl(PlcDeviceBuilder builder) {
        this.refPlcConnection = new AtomicReference<PlcConnection>();
        this.myProperties = new Hashtable<String, Object>();
        this.device_groups = new HashMap<UUID, PlcGroup>();
        this.bc = builder.bc;
        myProperties.put(PlcDevice.SERVICE_DRIVER, builder.service_driver);         
        myProperties.put(PlcDevice.SERVICE_NAME, builder.service_name);
        myProperties.put(PlcDevice.SERVICE_DESCRIPTION, builder.service_description); 
        if (null != builder.service_uid) {
            myProperties.put(PlcDevice.SERVICE_UID, builder.service_uid.toString());
        } else {
            myProperties.put(PlcDevice.SERVICE_UID, UUID.randomUUID().toString());            
        }        
        
        if (null != builder.device_category) myProperties.put(Constants.DEVICE_CATEGORY, builder.device_category);
        if (null != builder.service_firmware_vendor) myProperties.put(Device.SERVICE_FIRMWARE_VENDOR, builder.service_firmware_vendor);  
        if (null != builder.service_firmware_version) myProperties.put(Device.SERVICE_FIRMWARE_VERSION, builder.service_firmware_version);  
        if (null != builder.service_hardware_vendor) myProperties.put(Device.SERVICE_HARDWARE_VENDOR, builder.service_hardware_vendor); 
        if (null != builder.service_hardware_version) myProperties.put(Device.SERVICE_HARDWARE_VERSION, builder.service_hardware_version);            
        if (null != builder.service_model) myProperties.put(Device.SERVICE_MODEL, builder.service_model); 
        if (null != builder.service_reference_uids) myProperties.put(Device.SERVICE_REFERENCE_UIDS, builder.service_reference_uids);
        if (null != builder.service_serial_number) myProperties.put(Device.SERVICE_SERIAL_NUMBER, builder.service_serial_number);
        if (null != builder.service_status) myProperties.put(Device.SERVICE_STATUS, builder.service_status);
        if (null != builder.service_status_detail) myProperties.put(Device.SERVICE_STATUS_DETAIL, builder.service_status_detail);
        if (null != builder.service_types) myProperties.put(Device.SERVICE_TYPES, builder.service_types);      
    }
        
    @Override
    public void init() throws Exception {
        //
    }

    @Override
    public void destroy() throws Exception {
        //        
    }

    @Override
    public void enable() {
        if (null != plcDriver) {
            //Try to connect
            final String url = (String) myProperties.get(Device.SERVICE_DRIVER);
            try {
                plcConnection = plcDriver.getConnection(url);
                plcConnection.connect();
                refPlcConnection.set(plcConnection);
                if (plcConnection.isConnected()) {
                    enable = true;
                    LOGGER.info("Device [{}] was enbale.", myProperties.get(Device.SERVICE_NAME));
                } else {
                    LOGGER.info("The connection could not be established, check the url.");
                }               
            } catch (PlcConnectionException ex) {
                LOGGER.info(ex.getLocalizedMessage());
                enable = false;
            }
        } else {
            LOGGER.info("The PlcDriver has not been assigned to the device.");
        }
    }

    @Override
    public void disable() {
        enable = false;
        LOGGER.info("Device [{}] was disable.", myProperties.get(Device.SERVICE_NAME));        
        try {
            if (null != plcConnection) {
                //All groups are disabled for security, they are activated 
                //individually manually. Simple job to do when the IDE 
                //is available.                 
                device_groups.forEach((u, d) -> d.disable());                
                plcConnection.close();
                if (!plcConnection.isConnected()) {
                    enable = false;
                    LOGGER.info("Device [{}] connection was close.", myProperties.get(Device.SERVICE_NAME));
                }
            }
        } catch (Exception ex) {
            LOGGER.info(ex.getLocalizedMessage());
        }
    }    
            
    @Override
    public boolean isEnable() {
        return enable;
    }    
    
    @Override
    public Hashtable<String, ?> getProperties() {
        return myProperties;
    }
        
    @Override
    public String getDeviceName() {
        return (String) myProperties.get(Device.SERVICE_NAME);
    }

    @Override
    public void setDeviceName(String devicename) {
        myProperties.put(Device.SERVICE_NAME, devicename);
    }

    @Override
    public String getDeviceDescription() {
        return (String) myProperties.get(Device.SERVICE_DESCRIPTION);
    }

    @Override
    public void setDeviceDescription(String devicedescription) {
        myProperties.put(Device.SERVICE_DESCRIPTION, devicedescription); 
    }

    @Override
    public void setUid(UUID uid) {
        myProperties.put(PlcDevice.SERVICE_UID, uid.toString());
    }

    @Override
    public UUID getUid() {
        return UUID.fromString((String) myProperties.get(PlcDevice.SERVICE_UID));
    }

    @Override
    public void setUrl(String url) {
        if (!enable) {
            myProperties.put(Device.SERVICE_DRIVER, url); 
        }
    }

    @Override
    public String getUrl() {
        return (String) myProperties.get(Device.SERVICE_DRIVER); 
    }

    @Override
    public void putGroup(PlcGroup group) {
        if ((!enable) && (!device_groups.containsKey(group.getGroupUid()))) {
                group.setGroupDeviceUid(UUID.fromString((String) myProperties.get(PlcDevice.SERVICE_UID)));
                group.setPlcConnection(refPlcConnection);
                device_groups.put(group.getGroupUid(), group);
                bc.registerService(new String[]{Job.class.getName(), 
                    PlcGroup.class.getName()}, 
                  group, 
               group.getProperties());
        } else {
            LOGGER.info("The device is enabled or the group identifier already exists.");
        }
    }

    @Override
    public PlcGroup getGroup(UUID uid) {
        return device_groups.get(uid);
    }

    @Override
    public void removeGroup(UUID uid) {
        String filter = FILTER_DEVICE.replace("*", uid.toString()); 
        ServiceReference<?> sr = bc.getServiceReference(filter);
        bc.ungetService(sr); 
        device_groups.remove(uid);
    }

    @Override
    public List<PlcGroup> getGroups() { 
        return device_groups.values().stream().
                collect(toList());
    }
            
    @Override
    public void noDriverFound() {
        LOGGER.info("The associated driver is not found. go to IDLE.");
    }

    @Override
    public Object getServiceProperty(String propKey) {
        return myProperties.get(propKey);
    }

    @Override
    public String[] getServicePropertyKeys() {
        return myProperties.keySet().toArray(new String[myProperties.size()]);
    }

    @Override
    public void remove() throws DeviceException {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public void attach(PlcDriver driver) {
        LOGGER.info("Device: {} attach to driver {} ", myProperties.get(Device.SERVICE_NAME),  driver.getProtocolCode());
        this.plcDriver = driver;
    }


    public static class PlcDeviceBuilder {
        private final BundleContext bc;        
        private final String service_name;
        private final String service_description;
        private final String service_driver;         
        private UUID service_uid;          
        private String device_category;
        private String service_firmware_vendor;  
        private String service_firmware_version;  
        private String service_hardware_vendor; 
        private String service_hardware_version;            
        private String service_model; 
        private String[] service_reference_uids;
        private String service_serial_number;
        private String service_status;
        private String service_status_detail;
        private String[] service_types;         

        public PlcDeviceBuilder(BundleContext bc, String service_driver, String service_name, String service_description) {
            this.bc = bc;
            this.service_name = service_name;
            this.service_description = service_description;
            this.service_driver = service_driver;
            String[] drv = service_driver.split(":");
            this.device_category = drv[0];
        }

        public PlcDeviceBuilder setServiceUid(UUID serviceuid) {
            this.service_uid = serviceuid;
            return this;
        }

        public PlcDeviceBuilder setDeviceCategory(String device_category) {
            this.device_category = device_category;
            return this;            
        }

        public PlcDeviceBuilder setServiceFirmwareVendor(String service_firmware_vendor) {
            this.service_firmware_vendor = service_firmware_vendor;
            return this;            
        }

        public PlcDeviceBuilder setServiceFirmwareVersion(String service_firmware_version) {
            this.service_firmware_version = service_firmware_version;
            return this;            
        }

        public PlcDeviceBuilder setServiceHardwareVendor(String service_hardware_vendor) {
            this.service_hardware_vendor = service_hardware_vendor;
            return this;            
        }

        public PlcDeviceBuilder setServiceHardwareVersion(String service_hardware_version) {
            this.service_hardware_version = service_hardware_version;
            return this;            
        }

        public PlcDeviceBuilder setServiceModel(String service_model) {
            this.service_model = service_model;
            return this;            
        }

        public PlcDeviceBuilder setServiceReferenceUids(String[] service_reference_uids) {
            this.service_reference_uids = service_reference_uids;
            return this;            
        }

        public PlcDeviceBuilder setServiceSerialNumber(String service_serial_number) {
            this.service_serial_number = service_serial_number;
            return this;            
        }

        public PlcDeviceBuilder setServiceStatus(String service_status) {
            this.service_status = service_status;
            return this;            
        }

        public PlcDeviceBuilder setServiceStatusDetail(String service_status_detail) {
            this.service_status_detail = service_status_detail;
            return this;            
        }

        public PlcDeviceBuilder setServiceTypes(String[] service_types) {
            this.service_types = service_types;
            return this;            
        }
        

        public PlcDevice build() {
            PlcDevice plcdevice = new PlcDeviceImpl(this);
            validateBaseDeviceObject(plcdevice);
            return plcdevice;
        }
        
        private void validateBaseDeviceObject(PlcDevice plcdevice) {
            //
        }
    }
    

}
