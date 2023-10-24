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

import java.util.Collection;
import java.util.Hashtable;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.atomic.AtomicReference;
import java.util.stream.Collectors;
import org.apache.commons.lang3.time.StopWatch;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.merlot.scheduler.api.Job;
import org.apache.plc4x.merlot.scheduler.api.JobContext;
import org.osgi.framework.BundleContext;
import org.apache.plc4x.merlot.api.PlcItem;
import org.apache.plc4x.merlot.api.PlcGroup;
import org.apache.plc4x.merlot.api.core.PlcItemClientService;
import org.osgi.framework.ServiceReference;
import org.slf4j.LoggerFactory;

/*
*
*/
public class PlcGroupImpl implements PlcGroup, Job {
    private static final org.slf4j.Logger LOGGER = LoggerFactory.getLogger(PlcGroupImpl.class);
    private static final String FILTER_ITEM =  "(&(" + org.osgi.framework.Constants.OBJECTCLASS + "=" + PlcItem.class.getName() + ")" +
                    "(" + PlcItem.ITEM_UID + "=*))";     
    
    protected final BundleContext bc;    
    
    private UUID uid;
        
    private boolean enable = false;
       
    private long grouptransmit = 0;
    private long groupreceives = 0;    
    private long grouperrors = 0;  
    
    private long groupitems = 0;
    private long groupupdatrate = -1;
     
    private final Hashtable<String, Object> myproperties;
    
    private final Map<UUID, PlcItem> group_items;
    
    private PlcItemClientService items_service = null;
    
    private AtomicReference<PlcConnection> plcConnection = null;
    private PlcReadRequest.Builder builder = null;     

    private StopWatch watch = new StopWatch();
    private long[] aux = new long[1];
    
    

    public PlcGroupImpl(PlcGroupBuilder builder) { 
        this.bc = builder.bc;
        this.items_service = builder.items_service;
        this.group_items = new Hashtable<>();
        this.myproperties = new Hashtable<>();

        myproperties.put(PlcGroup.GROUP_UID, builder.group_uid.toString());  
        myproperties.put(PlcGroup.GROUP_NAME, builder.group_name);
        if (null != builder.device_uid)
        myproperties.put(PlcGroup.GROUP_DEVICE_UID, builder.device_uid);        

        myproperties.put(PlcGroup.GROUP_DESCRIPTION, builder.group_description);              
        myproperties.put(PlcGroup.GROUP_CONCURRENT, false);        
        myproperties.put(PlcGroup.GROUP_IMMEDIATE, true);
        myproperties.put(PlcGroup.GROUP_PERIOD, builder.group_period);
        
    }
        
    public void start(int bc){
        //
        if (null != plcConnection) {
            if (plcConnection.get().isConnected()) {
                enable = true;
            } else {
                enable = false;
            }
        } else {
            LOGGER.info("The PlcConnection has not been assigned to the group.");
            enable = false;
        }
    }
    
    public void stop(int bc){
        enable = false;
    }

    @Override
    public void enable() {
        enable = true;
    }

    @Override
    public void disable() {
        enable = false;
    }
    
    
    @Override
    public UUID getGroupUid() {
        return  UUID.fromString((String) myproperties.get(PlcGroup.GROUP_UID));
    }    

    @Override
    public UUID getGroupDeviceUid() {
        return  (UUID) myproperties.get(PlcGroup.GROUP_DEVICE_UID);
    }

    @Override
    public void setGroupDeviceUid(UUID device_uid) {
        myproperties.put(PlcGroup.GROUP_DEVICE_UID, device_uid);  
    }

    @Override
    public String getGroupName() {
        return (String) myproperties.get(PlcGroup.GROUP_NAME);
    }

    @Override
    public void setGroupName(String groupname) {
        myproperties.put(PlcGroup.GROUP_NAME, groupname);
    }

    @Override
    public String getGroupDescription() {
        return (String) myproperties.get(PlcGroup.GROUP_DESCRIPTION); 
    }

    @Override
    public void setGroupDescription(String groupdescription) {
        myproperties.put(PlcGroup.GROUP_DESCRIPTION, groupdescription);
    }

    @Override
    public boolean isEnable() {
        return enable; 
    }

    @Override
    public long getPeriod() {
        return (long) myproperties.get(PlcGroup.GROUP_PERIOD); 
    }

    @Override
    public void setPeriod(long period) {
        myproperties.put(PlcGroup.GROUP_PERIOD, (period < 100)?100:period);        
    }

    @Override
    public long getGroupTransmit() {  
        aux[0] = 0;
        group_items.forEach((uid, item) -> {aux[0] = +item.getItemTransmits();});
        return aux[0];
    }

    @Override
    public long getGroupReceives() {
        aux[0] = 0;
        group_items.forEach((uid, item) -> {aux[0] = +item.getItemReceives();});
        return aux[0];
    }

    @Override
    public long getGroupErrors() {
        aux[0] = 0;
        group_items.forEach((uid, item) -> {aux[0] = +item.getItemErrors();});
        return aux[0];
    }

    @Override
    public void setPlcConnection(AtomicReference<PlcConnection> plcConnection) {
        LOGGER.info("Grupo [{}] Volatile: Se asigno la conexión.", myproperties.get(PlcGroup.GROUP_NAME));
        this.plcConnection = plcConnection;
    }
    
    @Override
    public Map<UUID, PlcItem> getGroupItems() {
        return group_items; 
    }

    @Override
    public void setGroupItems(long groupitems) {
        throw new UnsupportedOperationException("Not supported yet."); 
    }

    @Override
    public Hashtable<String, Object> getProperties() {
        return myproperties;
    }
           
    @Override
    public void putItem(PlcItem item) {
        if (!group_items.containsKey(item.getItemUid())) {
            group_items.put(item.getItemUid(), item);
            bc.registerService(PlcItem.class.getName(), item, item.getProperties());            
        }                   
    }

    @Override
    public PlcItem getItem(UUID itemuid) {
        return group_items.get(itemuid);
    }

    //TODO: remove from context
    @Override
    public void removeItem(UUID itemuid) {
        String filter = FILTER_ITEM.replace("*", itemuid.toString());
        ServiceReference<?> sr = bc.getServiceReference(filter);
        bc.ungetService(sr);        
        group_items.remove(itemuid);
    }

    @Override
    public List<PlcItem> getItems() {
        return group_items.values().stream().
                collect(Collectors.toList());
    }

    @Override
    public void execute(JobContext context) {
        if (enable) {
            if (null != plcConnection) {
                if (plcConnection.get().isConnected()) {
                    watch.start();
                    executeReadAllItems();
                    watch.stop();
                    System.out.println("Elapse time: " + watch.getTime());
                     watch.reset();
                } else {
                    LOGGER.info("The driver is disconnected.");
                }
            } else {
                LOGGER.info("Unassigned connection.");
            }
        } else {
            LOGGER.info("The group {}:{} is disable.", 
                    ((String) myproperties.get(PlcGroup.GROUP_NAME)),
                    ((String) myproperties.get(PlcGroup.GROUP_UID)));
        }
    }
    
    /*
    * Execute the read function for all items
    * 
    */
    private void executeReadAllItems() {
        //1. The item was 
        builder = plcConnection.get().readRequestBuilder();
        group_items.forEach((u,i) ->{
            if (i.isEnable()) {
                builder.addTagAddress(u.toString(), i.getItemId());
            }
        }); 
        final PlcReadRequest readRequest = builder.build();        
        try {        
            PlcReadResponse syncResponse = readRequest.execute().get();
            group_items.forEach((u,i) -> {
                 System.out.println(i.getItemName()); 
                final PlcValue plcvalue = syncResponse.getPlcValue(u.toString());
                if (null == plcvalue) System.out.println("Valor nulo");
                System.out.println(i.getItemName() + " : " + plcvalue.getInt());
                items_service.putItemEvent(u, plcvalue);                
            });

        } catch (Exception ex) {
            LOGGER.info(ex.getMessage());
        }
    }

    public static class PlcGroupBuilder {
        protected final BundleContext bc;        
        private final String group_name;
        private String group_description; 
        private UUID group_uid;
        private UUID device_uid = null;        
        private boolean group_enable = false;
        private long group_period = 100;
        private PlcItemClientService items_service = null;
        
        public PlcGroupBuilder(BundleContext bc, String group_name, UUID group_uid) {
            this.bc = bc;
            this.group_name = group_name;
            this.group_uid = group_uid; 
            this.group_description = "";
        }
        
        public PlcGroupBuilder(BundleContext bc, String group_name) {
            this.bc = bc;
            this.group_name = group_name;
            this.group_uid = UUID.randomUUID();
            this.group_description = "";
        }        

        public PlcGroupBuilder  setGroupDescription(String group_description) {
            this.group_description = group_description;            
            return this;
        }

        public PlcGroupBuilder  setGroupUid(UUID group_uid) {
            this.group_uid = group_uid;
            return this;
        }
        
        public PlcGroupBuilder  setGroupDeviceUid(UUID device_uid) {
            this.device_uid = device_uid;
            return this;
        }        

        public PlcGroupBuilder  setGroupEnable(boolean group_enable) {
            this.group_enable = group_enable;
            return this;
        }

        public PlcGroupBuilder  setGroupPeriod(long group_period) {
            this.group_period = group_period;
            return this;
        }
        
        public PlcGroupBuilder  setItemService(PlcItemClientService items_service) {
            this.items_service = items_service;
            return this;
        }        
        
        public PlcGroup build() {
            PlcGroup plcgroup = new PlcGroupImpl(this);
            validatePlcGroupObject(plcgroup);
            return plcgroup;
        }
        
        private void validatePlcGroupObject(PlcGroup plcgroup) {
            //
        }        
        
        
        
    }    
    
}
