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

import io.netty.buffer.ByteBuf;
import java.util.HashMap;
import java.util.Hashtable;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.merlot.api.PlcDevice;
import org.apache.plc4x.merlot.api.PlcDeviceFactory;
import org.apache.plc4x.merlot.api.PlcFunction;
import org.apache.plc4x.merlot.api.PlcGroup;
import org.apache.plc4x.merlot.api.PlcItem;
import org.apache.plc4x.merlot.scheduler.api.ScheduleOptions;
import org.apache.plc4x.merlot.scheduler.api.Scheduler;
import org.apache.plc4x.merlot.scheduler.api.SchedulerError;
import org.osgi.framework.BundleContext;
import org.osgi.framework.Constants;
import org.osgi.framework.InvalidSyntaxException;
import org.osgi.framework.ServiceReference;
import org.osgi.service.dal.OperationMetadata;
import org.osgi.service.dal.PropertyMetadata;
import org.slf4j.LoggerFactory;
import org.apache.plc4x.merlot.api.PlcGeneralFunction;

/*
*
*/
public class PlcGeneralFunctionImpl implements PlcGeneralFunction  {
    
    private static final org.slf4j.Logger LOGGER = LoggerFactory.getLogger(PlcGeneralFunctionImpl.class); 
    
    private static String FILTER_DRIVER =  "(&(" + Constants.OBJECTCLASS + "=" + PlcDriver.class.getName() + ")" +
                        "(org.apache.plc4x.driver.code=*))";

    private static String FILTER_DEVICE_CATEGORY =  "(&(" + Constants.OBJECTCLASS + "=" + PlcDevice.class.getName() + ")" +
                        "(" + org.osgi.service.device.Constants.DEVICE_CATEGORY  + "=*))";    
    
    private static final String FILTER_DEVICE_UID =  "(&(" + 
            org.osgi.framework.Constants.OBJECTCLASS + "=" + 
            PlcDevice.class.getName() + ")" +
            "(" + PlcDevice.SERVICE_UID + "=*))";
    
    private static String FILTER_DEVICE_GROUP =  "(&(" + Constants.OBJECTCLASS + "=" + PlcGroup.class.getName() + ")" +
                        "(" + org.apache.plc4x.merlot.api.PlcGroup.GROUP_DEVICE_UID + "=*))";    
    
    private static String FILTER_GROUP_UID =  "(&(" + Constants.OBJECTCLASS + "=" + PlcGroup.class.getName() + ")" +
                        "(" + org.apache.plc4x.merlot.api.PlcGroup.GROUP_UID + "=*))";  
       
    private static String FILTER_ITEM_UID =  "(&(" + Constants.OBJECTCLASS + "=" + PlcItem.class.getName() + ")" +
                        "(" + PlcItem.ITEM_UID + "=*))";  
     
   private static String FILTER_FACTORY =  "(&(" + Constants.OBJECTCLASS + "=" + PlcDeviceFactory.class.getName() + ")" +
                        "(org.apache.plc4x.device.factory=*))";        
    
    
    
    private static final String[] operations = {"getPlcDrivers",
                                                "getPlcDevices",
                                                "getPlcDevice",
                                                "getPlcDeviceMeta",
                                                "setPlcDeviceName", 
                                                "setPlcDeviceDesc",      
                                                "getPlcDeviceGroups",
                                                "getPlcGroup",
                                                "setPlcGroupScanRate",
                                                "getPlcGroupItemsCount",
                                                "getPlcGroupItems",
                                                "getPlcItem",
                                                "getPlcItemType",
                                                "getPlcItemValue",
                                                "getPlcItemByteBuf",
                                                "enable",
                                                "disable", 
                                                "read", 
                                                "write"};
    
    private static final PlcOperationMetadata OP_GET_PLCDRIVERS_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("List all Plc4x drivers in the context.").
            setReturnValueMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("value").
                                        putLongDescription("All  PlcDriver in a Map<String, String>.").                    
                                        build()).
            build();    
    
    private static final PlcOperationMetadata OP_GET_PLCDEVICES_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("List all device in the context.").
            setReturnValueMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("value").
                                        putLongDescription("All  PlcDevice in a Map<UUID, String>.").                    
                                        build()).            
            build();
    
    private static final PlcOperationMetadata OP_GET_PLCDEVICE_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Return a PlcDevice .").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("UUID of the device.").                    
                                        build()).
            setReturnValueMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("value").
                                        putLongDescription("Return a PlcDevice with the uid equal to uuid.").                    
                                        build()).              
            build();  
    
    private static final PlcOperationMetadata OP_GET_PLCDEVICE_META_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Return metadata from asociated PlcDevice.").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("UUID of the device.").                    
                                        build()).
            setReturnValueMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("value").
                                        putLongDescription("All  PlcDevice in a Map<UUID, Object>.").                    
                                        build()).              
            build(); 
    
    private static final PlcOperationMetadata OP_SET_PLCDEVICE_NAME_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Assigns the device name indicated by the UUID.").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("UUID of the device.").                    
                                        build()).
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("device_name").
                                        putLongDescription("String with new name.").                    
                                        build()).            
            build(); 

    private static final PlcOperationMetadata OP_SET_PLCDEVICE_DESC_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Assigns the device description indicated by the UUID.").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("UUID of the device.").                    
                                        build()).
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("device_desc").
                                        putLongDescription("String with new device description.").                    
                                        build()).             
            build();     
    
    private static final PlcOperationMetadata OP_GET_PLCDEVICE_GROUPS_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Lists the all groups associated with a device.").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("UUID of the device.").                    
                                        build()).
            setReturnValueMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("value").
                                        putLongDescription("All  PlcGroup in a Map<UUID, String>.").                    
                                        build()).              
            build();   
    
    private static final PlcOperationMetadata OP_GET_PLCGROUP_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Return the PlcGroup indicate by uid.").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("UUID of the group.").                    
                                        build()).
            setReturnValueMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("value").
                                        putLongDescription("Return PlcGroup searched by uid.").                    
                                        build()).              
            build();  
    
    private static final PlcOperationMetadata OP_SET_PLCGROUP_SR_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Set the scan rate (period) of the group").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("UUID of the device.").                    
                                        build()).
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("scan_rate").
                                        putLongDescription("a long value for the scan rate in ms. Minimum 100 ms").                    
                                        build()).            
            build();     
    
    private static final PlcOperationMetadata OP_GET_PLCGROUP_ITEMS_COUNT_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Get the number of items in a group.").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("UUID of the group.").                    
                                        build()).
            setReturnValueMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("value").
                                        putLongDescription("Return a integer with the number of items in the group.").                    
                                        build()).            
            build(); 

    private static final PlcOperationMetadata OP_GET_PLCGROUP_ITEMS_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Lists the items associated with a group.").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("UUID of the group.").                    
                                        build()).
            setReturnValueMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("value").
                                        putLongDescription("All  PlcItem in a group like a Map<UUID, String>.").                    
                                        build()).             
            build();
    
    private static final PlcOperationMetadata OP_GET_PLCITEM_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Get a PlcItem from a group.").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("UUID of the item.").                    
                                        build()).
            setReturnValueMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("value").
                                        putLongDescription("A PlcItem from  the group").                    
                                        build()).             
            build();  
    
    private static final PlcOperationMetadata OP_GET_PLCITEM_TYPE_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Get the default type for a PlcItem.").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("item_uid").
                                        putLongDescription("UUID of the device.").                    
                                        build()).
            setReturnValueMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("value").
                                        putLongDescription("A PlcItem default Type.").                    
                                        build()).             
            build();  

    private static final PlcOperationMetadata OP_GET_PLCITEM_VALUE_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Lists the groups associated with a device.").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("UUID of the device.").                    
                                        build()).
            setReturnValueMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("value").
                                        putLongDescription("Return PlcValue of PlcItem default Type.").                    
                                        build()).             
            build();      
        
    private static final PlcOperationMetadata OP_GET_PLCITEM_BYTEBUF_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Lists the groups associated with a device.").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("UUID of the device.").                    
                                        build()).
            build();      
    
    private static final PlcOperationMetadata OP_ENABLE_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Enable a single PlcItem.").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("String representing the UUID of the device.").                    
                                        build()).
            build();
    
    private static final PlcOperationMetadata OP_DISABLE_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Disable a single PlcItem.").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("String representing the UUID of the device."). 
                                        build()).
            build();
    
    private static final PlcOperationMetadata OP_READ_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder(). 
            putDescription("Read a single PlcValue from PlcItem.").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("String representing the UUID of the device."). 
                                        build()).
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("tag").
                                        putLongDescription("PlcItem ID.").    
                                        build()).              
            build();    
    
    private static final PlcOperationMetadata OP_WRITE_META = new PlcOperationMetadata.
            PlcOperationMetaDataBuilder().
            putDescription("Write a single PlcValue from PlcItem.").
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("uuid").
                                        putLongDescription("String representing the UUID of the device.").
                                        build()).
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("tag").
                                        putLongDescription("PlcItem ID.").                    
                                        build()).     
            addOperationMetadata(new    PlcPropertyMetadata.PlcPropertyMetaDataBuilder().
                                        putDescription("value").
                                        putLongDescription("Value to write in String."). 
                                        build()).             
            build();        
    
    private final BundleContext bc;
    private final Scheduler scheduler;    
    protected Hashtable<String, Object> myProperties;    

    public PlcGeneralFunctionImpl(BundleContext bc,  Scheduler scheduler) {
        this.myProperties = new Hashtable<String, Object>();        
        this.bc = bc;
        this.scheduler = scheduler;     
        myProperties.put(PlcFunction.SERVICE_OPERATION_NAMES, operations);

    }    

    public void bind(PlcFunction fc, Map properties) {
        myProperties.clear();
        myProperties.putAll(properties);
    }

    public void unbind(PlcFunction fc, Map properties) {
        //
    }

    @Override
    public Hashtable<String, Object> getProperties() {
        return myProperties;
    }
            
    @Override
    public PropertyMetadata getPropertyMetadata(String propertyName) {
        return null;
    }

    @Override
    public OperationMetadata getOperationMetadata(String operationName) {
        switch (operationName) {
            case "getPlcDrivers":           return OP_GET_PLCDRIVERS_META;
            case "getPlcDevices":           return OP_GET_PLCDEVICES_META; 
            case "getPlcDevice":            return OP_GET_PLCDEVICE_META;
            case "getPlcDeviceMeta":        return OP_GET_PLCDEVICE_META_META;
            case "setPlcDeviceName":        return OP_SET_PLCDEVICE_NAME_META; 
            case "setPlcDeviceDesc":        return OP_SET_PLCDEVICE_DESC_META; 
            case "getPlcDeviceGroups":      return OP_GET_PLCDEVICE_GROUPS_META; 
            case "getPlcGroup":             return OP_GET_PLCGROUP_META; 
            case "setPlcGroupScanRate":     return OP_SET_PLCGROUP_SR_META; 
            case "getPlcGroupItemsCount":   return OP_GET_PLCGROUP_ITEMS_COUNT_META;
            case "getPlcGroupItems":        return OP_GET_PLCGROUP_ITEMS_META; 
            case "getPlcItem":              return OP_GET_PLCITEM_META; 
            case "getPlcItemType":          return OP_GET_PLCITEM_TYPE_META;  
            case "getPlcItemValue":         return OP_GET_PLCITEM_VALUE_META;   
            case "getPlcItemByteBuf":       return OP_GET_PLCITEM_BYTEBUF_META;  
            case "enable":                  return OP_ENABLE_META; 
            case "disable":                 return OP_DISABLE_META; 
            case "read":                    return OP_READ_META; 
            case "write":                   return OP_WRITE_META;         
            default: return null;    
        }
    }

    @Override
    public Object getServiceProperty(String propKey) {
        return null;
    }
    
    @Override
    public String[] getServicePropertyKeys() {
        return operations;
    }
          
    @Override
    public Map<String, String> getPlcDrivers() {
        Map<String, String> drivers = new HashMap<>();
        
        try {
            ServiceReference[] refs = bc.getServiceReferences((String) null, FILTER_DRIVER);
            
            if (null != refs) {
                for (ServiceReference ref:refs){
                    final PlcDriver driver = (PlcDriver) bc.getService(ref);
                    drivers.put(driver.getProtocolCode(),driver.getProtocolName());
                }
            }
        } catch (InvalidSyntaxException ex) {
            LOGGER.info(ex.getMessage());
        } finally {
            return drivers;
        }        
    }    
                
    @Override
    public Map<UUID, String> getPlcDevices(String driver_code) {
        Map<UUID, String> plcdevices = new HashMap<>();
        
        try {
            String filter = FILTER_DEVICE_CATEGORY.replace("*", driver_code);
            ServiceReference[] refs = bc.getServiceReferences((String) null, filter);
            if (null != refs) {
                for (ServiceReference ref:refs){
                    final PlcDevice device = (PlcDevice) bc.getService(ref);
                    plcdevices.put(device.getUid(),device.getDeviceName());
                }
            }
        } catch (InvalidSyntaxException ex) {
            LOGGER.info(ex.getMessage());
        } finally {
            return  plcdevices;
        }      
    }
        
    @Override
    public PlcDevice getPlcDevice(UUID device_uid) {
        String filter = FILTER_DEVICE_UID.replace("*", device_uid.toString());
        try {        
            ServiceReference[] refs = bc.getServiceReferences((String) null, filter);
            if (null != refs) {         
                final PlcDevice plcdevice = (PlcDevice) bc.getService(refs[0]);  
                return plcdevice;
            } else {
                LOGGER.info("No device detectec with filter: " + filter);
            }           
        } catch (InvalidSyntaxException ex) {
            LOGGER.info(ex.getMessage());
        }
      
        return null;
    }    
    
    @Override
    public Map<String, Object> getPlcDeviceMeta(UUID group_uid) {
        Map<String, Object> props = new HashMap<>();
        
        try {
            String filter = FILTER_DEVICE_UID.replace("*", group_uid.toString());
            ServiceReference ref = bc.getServiceReference(filter);
            if (null != ref) {
                final PlcDevice device = (PlcDevice) bc.getService(ref);
                props.putAll(device.getProperties());
            }
        } catch (Exception ex) {
            LOGGER.info(ex.getMessage());
        } finally {
            return  props;
        } 
    }

    @Override
    public void setPlcDeviceName(UUID device_uid, String group_name) {
        final PlcDevice plcdevice = getPlcDevice(device_uid);
        if (null != plcdevice){
            plcdevice.setDeviceName(group_name);
        }
    }

    @Override
    public void setPlcDeviceDesc(UUID device_uid, String device_desc) {
        final PlcDevice plcdevice = getPlcDevice(device_uid);
        if (null != plcdevice){
            plcdevice.setDeviceDescription(device_desc);
        } 
    }
        
    @Override
    public Map<UUID, String> getPlcDeviceGroups(UUID device_uid) {
        Map<UUID, String> plcgroups = new HashMap<>();
        
        try {
            String filter = FILTER_DEVICE_GROUP.replace("*", device_uid.toString());
            ServiceReference[] refs = bc.getServiceReferences((String) null, filter);
            if (null != refs) {
                for (ServiceReference ref:refs){
                    final PlcGroup group = (PlcGroup) bc.getService(ref);
                    plcgroups.put(group.getGroupUid(),group.getGroupName());
                }
            }
        } catch (InvalidSyntaxException ex) {
            LOGGER.info(ex.getMessage());
        } finally {
            return  plcgroups;
        }
    }    
   
    @Override
    public PlcGroup getPlcGroup(UUID group_uid) {
        String filter = FILTER_GROUP_UID.replace("*", group_uid.toString());
        try {        
            ServiceReference[] refs = bc.getServiceReferences((String) null, filter);
            if (null != refs) {
                final PlcGroup plcgroup = (PlcGroup) bc.getService(refs[0]);  
                return plcgroup;
            };            
        } catch (InvalidSyntaxException ex) {
            LOGGER.info(ex.getMessage());
        }

        
        return null;
    }    
    
    @Override
    public int setPlcGroupScanRate(UUID group_uid, long scan_rate) {
        final PlcGroup plcgroup = getPlcGroup(group_uid);
        if (null != plcgroup){
            try {
                scheduler.getJobs().entrySet().stream().
                        filter(e -> e.getKey().startsWith(plcgroup.getGroupName())).
                        forEach(e ->{
                            ScheduleOptions options = scheduler.NOW(-1, scan_rate);
                            try {
                                scheduler.reschedule(e.getKey(), options);
                            } catch (IllegalArgumentException ex) {
                                LOGGER.info(ex.getMessage());
                            } catch (SchedulerError ex) {
                                LOGGER.info(ex.getMessage());
                            }
                        });
            } catch (SchedulerError ex) {
                LOGGER.info(ex.getMessage());
            }
        } 
        return 0;
    }    
    
    @Override
    public int getPlcGroupItemsCount(UUID group_uid) {
        final PlcGroup plcgroup = getPlcGroup(group_uid);
        return plcgroup.getItems().size();
    }

    @Override
    public Map<UUID, String> getPlcGroupItems(UUID group_uid) {
        Map<UUID, String> plcitems = new HashMap<>();        
        final PlcGroup plcgroup = getPlcGroup(group_uid);
        List<PlcItem> items = plcgroup.getItems();
        items.forEach(i -> {plcitems.put(i.getItemUid(), i.getItemName());});
        return plcitems;
    }

    //Items    
    
    @Override
    public PlcItem getPlcItem(UUID item_uid) {
        String filter = FILTER_ITEM_UID.replace("*", item_uid.toString());
        ServiceReference ref = bc.getServiceReference(filter); 
        if (null != ref) {
            final PlcItem plcitem = (PlcItem) bc.getService(ref);  
            return plcitem;
        };
        
        return null;
    }    
    
    @Override
    public short getPlcItemType(UUID item_uid) {
        final PlcItem item = getPlcItem(item_uid);
        return item.getItemPlcValue().getPlcValueType().getValue();
    }

    @Override
    public PlcValue getPlcItemValue(UUID item_uid) {
        final PlcItem item = getPlcItem(item_uid);        
        return item.getItemPlcValue();
    }

    @Override
    public ByteBuf getPlcItemByteBuf(UUID item_uid) {
        final PlcItem item = getPlcItem(item_uid);        
        return item.getItemByteBuf();
    }

    @Override
    public byte[] getPlcItemBytes(UUID item_uid) {
        final PlcItem item = getPlcItem(item_uid);        
        return item.getItemBytes();
    }        
    
    @Override
    public void enable(UUID uuid) {
        final String strfilter = "(|" + 
                        FILTER_DEVICE_UID.replace("*", uuid.toString()) +
                        FILTER_GROUP_UID.replace("*", uuid.toString()) +
                        FILTER_ITEM_UID.replace("*", uuid.toString()) +
                        ")";
        final ServiceReference ref = bc.getServiceReference(strfilter);
        if (null != ref) {
            final Object plcasset =  bc.getService(ref);
            if (plcasset instanceof PlcDevice) {
                ((PlcDevice) plcasset).enable();
            } else if (plcasset instanceof PlcGroup){
                ((PlcGroup) plcasset).enable();
            } else if (plcasset instanceof PlcItem) {
                ((PlcItem) plcasset).enable();
            }
        } else {
            LOGGER.info("The Plc4x asset service does not exist uid: {}", uuid.toString());
        }        
    }
    
    @Override
    public void disable(UUID uuid) {
        final String strfilter = "(| " + 
                        FILTER_DEVICE_UID.replace("*", uuid.toString()) +
                        FILTER_GROUP_UID.replace("*", uuid.toString()) +
                        FILTER_ITEM_UID.replace("*", uuid.toString()) +
                        ")";
        final ServiceReference ref = bc.getServiceReference(strfilter);
        if (null != ref) {
            final Object plcasset =  bc.getService(ref);
            if (plcasset instanceof PlcDevice) {
                ((PlcDevice) plcasset).disable();
            } else if (plcasset instanceof PlcGroup){
                ((PlcGroup) plcasset).disable();                    
            } else if (plcasset instanceof PlcItem) {
                ((PlcItem) plcasset).disable();                       
            }
        } else {
            LOGGER.info("The Plc4x asset service does not exist uid: {}", uuid.toString());
        }   
    }
    
    @Override
    public PlcValue read(UUID uuid, String tag) {
        final String strfilter = FILTER_ITEM_UID.replace("*", uuid.toString());
        final ServiceReference ref = bc.getServiceReference(strfilter);
        if (null != ref) {
            final Object plcasset =  bc.getService(ref);            
            if (plcasset instanceof PlcItem) {
                final PlcValue plcvalue = ((PlcItem) plcasset).getItemPlcValue();
                return plcvalue;
            }
            LOGGER.info("The PlcItem asset service does not exist uid: {}", uuid.toString());            
        }        
        return null;
    }

    @Override
    public void write(UUID uuid, String tag, String value) {
        //
    }    
  
}
