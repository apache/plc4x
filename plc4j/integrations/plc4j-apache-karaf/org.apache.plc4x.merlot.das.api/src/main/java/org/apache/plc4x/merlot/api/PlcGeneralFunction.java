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
package org.apache.plc4x.merlot.api;

import io.netty.buffer.ByteBuf;
import java.util.Map;
import java.util.UUID;
import org.apache.plc4x.java.api.value.PlcValue;


public interface PlcGeneralFunction extends PlcFunction {
        
    
    /*
    * Returns a list of the Plc Drivers registered in the context.
    *
    * @return Returns a map with key corresponding to the "code" of the driver 
    *         and value corresponding to the "name" of the driver, 
    *         where "name" is the long description
    */
    public Map<String, String> getPlcDrivers();
           
    /*
    * Returns a list of the PlcDevice registered in the context.
    *
    * @param driver_code A String with the "code" of the driver.
    * @return Returns a map with key corresponding to the "uid" of the device 
    *         and value corresponding to the "name" of the device. 
    */    
    public Map<UUID, String> getPlcDevices(String driver_code); 
    
    /*
    * Returns a list of the PlcDevice registered in the context.
    *
    * @param driver_code A String with the "code" of the driver.
    * @return Returns a map with key corresponding to the "uid" of the device 
    *         and value corresponding to the "name" of the device. 
    */    
    public PlcDevice getPlcDevice(UUID device_uid); 
    
    /*
    * Returns the metadata associated with the device, according to 
    * chapter 141 Device Abstraction Layer Specification
    * 
    * To maintain compatibility between the Device access specification and 
    * Device abstraction layer specifications, the value of the 
    * DEVICE_CATEGORY meta corresponds to that defined in chapter 103.
    *
    * From Device chapter 103
    * DEVICE_CATEGORY           = "DEVICE_CATEGORY"
    * DEVICE_DESCRIPTION        = "DEVICE_DESCRIPTION"
    * DEVICE_SERIAL             = "DEVICE_SERIAL"
    * DRIVER_ID                 = "DRIVER_ID"
    * 
    * From Device chapter 141
    * SERVICE_UID               = "dal.device.UID"
    * SERVICE_REFERENCE_UIDS    = "dal.device.reference.UIDs"
    * SERVICE_DRIVER            = "dal.device.driver"
    * SERVICE_NAME              = "dal.device.name"
    * SERVICE_STATUS            = "dal.device.status"
    *                   STATUS_REMOVED, STATUS_OFFLINE, STATUS_ONLINE
    *                   STATUS_PROCESSING, STATUS_NOT_INITIALIZED, 
    *                   STATUS_NOT_CONFIGURED
    * SERVICE_STATUS_DETAIL     = "dal.device.status.detail"
    *                   STATUS_DETAIL_CONNECTING, STATUS_DETAIL_INITIALIZING 
    *                   STATUS_DETAIL_REMOVING,STATUS_DETAIL_FIRMWARE_UPDATING
    *                   STATUS_DETAIL_CONFIGURATION_UNAPPLIED,
    *                   STATUS_DETAIL_BROKEN,
    *                   STATUS_DETAIL_COMMUNICATION_ERROR,
    *                   STATUS_DETAIL_DATA_INSUFFICIENT,
    *                   STATUS_DETAIL_INACCESSIBLE,
    *                   STATUS_DETAIL_CONFIGURATION_ERROR,
    *                   STATUS_DETAIL_DUTY_CYCLE
    *
    * @param group_uid String representation of the group's UUID.
    * @return Returns a map with key corresponding to the property of the device 
    *         and value corresponding to this property of the device.     
    */
    public Map<String, Object> getPlcDeviceMeta(UUID group_uid);
    
    /*
    * @param group_uid String representation of the group's UUID.
    * @return Returns a map with key corresponding to the property of the device 
    *         and value corresponding to this property of the device.     
    */
    public void setPlcDeviceName(UUID device_uid, String device_name);     
    
    
    /*
    * @param group_uid String representation of the group's UUID.
    * @return Returns a map with key corresponding to the property of the device 
    *         and value corresponding to this property of the device.     
    */
    public void setPlcDeviceDesc(UUID device_uid, String device_desc); 
    
    /*
    * Returns a map of groups associated with the device defined by device_uid, 
    * where key corresponds to the uid of the group, and the value to its name.
    * 
    * @param device_uid String representation of the device's UUID.
    * @return Returns a map with key corresponding to the uid of the group 
    *         and value corresponding to the name of the group.     
    */
    public Map<UUID, String> getPlcDeviceGroups(UUID device_uid);     
    
    
    /*
    * Returns a PlcGroup registered in the context.
    *
    * @param group_uid A UUID with the uid of the group.
    * @return Returns a PlcGroup regsitered like service. 
    *
    */       
    public PlcGroup getPlcGroup(UUID group_uid);      
    
    /*
    * Allows you to change the group scanning period. The minimum 
    * value is 100 msec. This value can be changed even 
    * if the device is enabled.
    * 
    * @param group_uid String representation of the group's UUID.
    * @return Returns a map with key corresponding to the property of the device 
    *         and value corresponding to this property of the device.     
    */    
    public int setPlcGroupScanRate(UUID group_uid, long scan_rate); 
    
    /*
    * Returns the number of items associated with the 
    * indicated group of group_id.
    * 
    * @param group_uid String representation of the group's UUID.
    * @return Returns the number of items in the group.
    *
    */
    public int getPlcGroupItemsCount(UUID group_uid);     
    
    /*
    * Returns a map of the items associated with the group defined by group_is, 
    * where key corresponds to the uid of the item, and the value to its name.
    * 
    * @param group_uid String representation of the group's UUID.
    * @return Returns a map with key corresponding to the uid of the item 
    *         and value corresponding to the name of the item.     
    */
    public Map<UUID, String> getPlcGroupItems(UUID group_uid); 
    
    /*
    * Returns the data type associated with the item as defined in 
    * org.apache.plc4x.java.api.types.PlcValueType.
    *
    * @param item_uid String representation of the item's UUID.
    * @return Return a short value, representing the value type.     
    */    
    public PlcItem getPlcItem(UUID item_uid);       
        
    /*
    * Returns the data type associated with the item as defined in 
    * org.apache.plc4x.java.api.types.PlcValueType.
    *
    * @param item_uid String representation of the item's UUID.
    * @return Return a short value, representing the value type.     
    */    
    public short getPlcItemType(UUID item_uid);   
    
    /*
    * Returns the value associated with the Item.
    *
    * @param item_uid String representation of the item's UUID.
    * @return Return a short value, representing the value type. 
    *         It can be null.
    */    
    public PlcValue getPlcItemValue(UUID item_uid);  
    
    /*
    * Returns the value associated with the Item as a ByteBuf.
    *
    * @param item_uid String representation of the item's UUID.
    * @return Return a short value, representing the value type.     
    */    
    public ByteBuf getPlcItemByteBuf(UUID item_uid); 

    /*
    * Returns the value associated with the Item as a byte array.
    *

    * @return Return a short value, representing the value type.     
    */    
    public byte[] getPlcItemBytes(UUID item_uid);  
    
    /*
    * Enables the element corresponding to the UUID.
    *
    * @param item_uid String representation of the item's UUID.    
    */    
    public void enable(UUID uuid);
    
    /*
    * Disables the element corresponding to the UUID.
    *
    * @param item_uid String representation of the item's UUID.    
    */
    public void disable(UUID uuid);
        
    /*
    * Returns the value associated with the Item.
    *
    * @param uuid UUID representation of the device.
    * @param tag String representation of the item's UUID.
    * @return Return a short value, representing the value type. 
    *         It can be null.    
    */
    public PlcValue read(UUID uuid, String tag);

    /*
    * Writes a value to the item defined by UUID and tag. 
    * Value corresponds to the String of the value to be written.
    *
    * @param uuid String representation of the item's UUID.
    * @param tag String representation of the item's UUID. 
    * @param value String representation of the item's UUID.     
    */
    public void write(UUID uuid, String tag, String value);    
    
    
}
