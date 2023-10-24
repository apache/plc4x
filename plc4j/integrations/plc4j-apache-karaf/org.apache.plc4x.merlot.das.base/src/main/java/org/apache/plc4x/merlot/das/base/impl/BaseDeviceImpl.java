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
package org.apache.plc4x.merlot.das.base.impl;


import java.util.HashMap;
import org.apache.plc4x.merlot.api.Driver;
import java.util.Hashtable;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.merlot.api.PlcDevice;
import org.apache.plc4x.merlot.das.base.api.BaseDevice;
import org.osgi.framework.BundleContext;
import org.osgi.service.dal.Device;
import org.osgi.service.dal.DeviceException;
import org.osgi.service.device.Constants;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.apache.plc4x.merlot.api.PlcItem;
import org.apache.plc4x.merlot.api.PlcGroup;

 public class BaseDeviceImpl implements PlcDevice {
	
    private static final Logger LOGGER = LoggerFactory.getLogger(BaseDeviceImpl.class);
    protected final BundleContext bc;
    protected boolean enable  = false;      
    protected boolean autostart = false;  
    protected UUID uid = null;
    protected PlcDriver plcdriver = null;
    protected String url = null;
    
    protected Hashtable<String, Object> myProperties = new Hashtable<String, Object>();
    
    private final Map<UUID, PlcGroup> groups = new HashMap();     
    
    public BaseDeviceImpl(BaseDeviceBuilder builder) {
        this.bc = builder.bc;
        myProperties.put(Device.SERVICE_DRIVER, builder.service_driver);         
        myProperties.put(Device.SERVICE_NAME, builder.service_name);
        myProperties.put(Device.SERVICE_DESCRIPTION, builder.service_description);   
        if (null != builder.device_category) myProperties.put(Constants.DEVICE_CATEGORY, builder.device_category);
        
        if (null != builder.service_uid) myProperties.put(Device.SERVICE_UID, builder.service_uid);     
        
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


    public void start() {
        if (null != plcdriver) {
            //Try to connect
            
            //if is ok then enable
            enable = true;
        }
    }


    public void stop() {
        enable = false;
        //Disconect from plc.
        
    }  

    @Override
    public void enable() {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public void disable() {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
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
        this.uid = uid;
    }

    @Override
    public UUID getUid() {
        return uid;
    }

            
    @Override
    public boolean isEnable() {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public void setUrl(String url) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public String getUrl() {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public void putGroup(PlcGroup group) {
        if (!enable) {
            Hashtable props = new Hashtable();            
            props.put(PlcGroup.GROUP_NAME, "PRUEBA_DE_SCHEDULER");
            props.put(PlcGroup.GROUP_UID, UUID.randomUUID().toString());
            props.put(PlcGroup.GROUP_PERIOD, 1000L);
            props.put(PlcGroup.GROUP_IMMEDIATE, true);
            props.put(PlcGroup.GROUP_CONCURRENT, false);             
            groups.put(group.getGroupUid(), group);
        }
    }

    @Override
    public PlcGroup getGroup(UUID uid) {
        return groups.get(uid);
    }

    @Override
    public void removeGroup(UUID uid) {
        groups.remove(uid);
    }

    @Override
    public List<PlcGroup> getGroups() {
        return null;
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
        this.plcdriver = driver;
    }


    public static class BaseDeviceBuilder {
        private final BundleContext bc;        
        private final String service_name;
        private final String service_description;
        private final String service_driver;         
        private String service_uid;          
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

        public BaseDeviceBuilder(BundleContext bc, String service_driver, String service_name, String service_description) {
            this.bc = bc;
            this.service_name = service_name;
            this.service_description = service_description;
            this.service_driver = service_driver;
            String[] drv = service_driver.split(":");
            this.device_category = drv[0];
        }

        public BaseDeviceBuilder setServiceUid(String serviceuid) {
            this.service_uid = serviceuid;
            return this;
        }

        public BaseDeviceBuilder setDeviceCategory(String device_category) {
            this.device_category = device_category;
            return this;            
        }

        public BaseDeviceBuilder setServiceFirmwareVendor(String service_firmware_vendor) {
            this.service_firmware_vendor = service_firmware_vendor;
            return this;            
        }

        public BaseDeviceBuilder setServiceFirmwareVersion(String service_firmware_version) {
            this.service_firmware_version = service_firmware_version;
            return this;            
        }

        public BaseDeviceBuilder setServiceHardwareVendor(String service_hardware_vendor) {
            this.service_hardware_vendor = service_hardware_vendor;
            return this;            
        }

        public BaseDeviceBuilder setServiceHardwareVersion(String service_hardware_version) {
            this.service_hardware_version = service_hardware_version;
            return this;            
        }

        public BaseDeviceBuilder setServiceModel(String service_model) {
            this.service_model = service_model;
            return this;            
        }

        public BaseDeviceBuilder setServiceReferenceUids(String[] service_reference_uids) {
            this.service_reference_uids = service_reference_uids;
            return this;            
        }

        public BaseDeviceBuilder setServiceSerialNumber(String service_serial_number) {
            this.service_serial_number = service_serial_number;
            return this;            
        }

        public BaseDeviceBuilder setServiceStatus(String service_status) {
            this.service_status = service_status;
            return this;            
        }

        public BaseDeviceBuilder setServiceStatusDetail(String service_status_detail) {
            this.service_status_detail = service_status_detail;
            return this;            
        }

        public BaseDeviceBuilder setServiceTypes(String[] service_types) {
            this.service_types = service_types;
            return this;            
        }
        

        public PlcDevice build() {
            PlcDevice plcdevice = new BaseDeviceImpl(this);
            validateBaseDeviceObject(plcdevice);
            return plcdevice;
        }
        
        private void validateBaseDeviceObject(PlcDevice plcdevice) {
            //
        }
    }
    

}
