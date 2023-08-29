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
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.merlot.scheduler.api.Job;
import org.apache.plc4x.merlot.scheduler.api.JobContext;
import org.osgi.framework.BundleContext;
import org.apache.plc4x.merlot.api.PlcItem;
import org.apache.plc4x.merlot.api.PlcGroup;
import org.osgi.framework.ServiceReference;
import org.slf4j.LoggerFactory;


public class PlcGroupImpl implements PlcGroup, Job {
    private static final org.slf4j.Logger LOGGER = LoggerFactory.getLogger(PlcGroupImpl.class);
    private static final String FILTER_ITEM =  "(&(" + org.osgi.framework.Constants.OBJECTCLASS + "=" + PlcItem.class.getName() + ")" +
                    "(" + PlcItem.ITEM_UID + "=*))";     
    
    protected final BundleContext bc;    
    
    private UUID uid;
    
    private String  groupname;
    private String  groupdescription;
    
    private boolean enable = false;
    
    private long period;    
    private long deadband;
    
    private long grouptransmit = 0;
    private long groupreceives = 0;    
    private long grouperrors = 0;  
    
    private long groupitems = 0;
    private long groupupdatrate = -1;
     
    private final Hashtable<String, Object> myproperties;
    
    private final Map<UUID, PlcItem> plcitems;
    
    private PlcConnection plcConnection = null;
    private PlcReadRequest.Builder builder = null;     

    private long[] aux = new long[1];

    public PlcGroupImpl(PlcGroupBuilder builder) { 
        this.bc = builder.bc;
        this.plcitems = new Hashtable<>();
        this.myproperties = new Hashtable<>();
        myproperties.put(PlcGroup.GROUP_NAME, builder.group_name);
        myproperties.put(PlcGroup.GROUP_CONCURRENT, false);        
        myproperties.put(PlcGroup.GROUP_IMMEDIATE, true);        
        if (null != builder.group_name) myproperties.put(PlcGroup.GROUP_DESCRIPTION, builder.group_name); 
        if (null != builder.group_uid) myproperties.put(PlcGroup.GROUP_UID, builder.group_uid); 
    }
    
    
    public void start(int bc){
        //
        if (null != plcConnection) {
            if (plcConnection.isConnected()) {
                enable = true;
            } else {
                enable = false;
            }
        } else {
            enable = false;
        }
    }
    
    public void stop(int bc){
        enable = false;
    }
    
    @Override
    public UUID getGroupUid() {
        return  (UUID) myproperties.get(PlcGroup.GROUP_UID);
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
    public void setEnable(boolean enable) {
        this.enable = enable;
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
        plcitems.forEach((uid, item) -> {aux[0] = +item.getItemTransmits();});
        return aux[0];
    }

    @Override
    public long getGroupReceives() {
        aux[0] = 0;
        plcitems.forEach((uid, item) -> {aux[0] = +item.getItemReceives();});
        return aux[0];
    }

    @Override
    public long getGroupErrors() {
        aux[0] = 0;
        plcitems.forEach((uid, item) -> {aux[0] = +item.getItemErrors();});
        return aux[0];
    }

    @Override
    public Map<UUID, PlcItem> getGroupItems() {
        return plcitems; 
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
        if (!plcitems.containsKey(item)) {
            plcitems.put(item.getItemUid(), item);
            Hashtable properties = new Hashtable();
            properties.putIfAbsent(PlcGroup.GROUP_UID, uid.toString());             
            properties.putIfAbsent(PlcItem.ITEM_UID, item.getItemUid().toString()); 
            bc.registerService(PlcItem.class.getName(), item, properties);            
        }                   
    }

    @Override
    public PlcItem getItem(UUID itemuid) {
        return plcitems.get(itemuid);
    }

    //TODO: remove from context
    @Override
    public void removeItem(UUID itemuid) {
        String filter = FILTER_ITEM.replace("*", itemuid.toString());
        ServiceReference<?> sr = bc.getServiceReference(filter);
        bc.ungetService(sr);        
        plcitems.remove(itemuid);
    }

    @Override
    public List<PlcItem> getItems() {
        return null;
    }

    @Override
    public void execute(JobContext context) {
        if (enable) {
            System.out.println("Tamaño de la coleccion: " + plcitems.size());
        }
    }
    
    /*
    * Execute the read function for all items
    * 
    */
    private void executeReadAllItems() {
        //1. The item was 
        builder = plcConnection.readRequestBuilder();
        plcitems.forEach((u,i) ->{
            if (i.getIsEnable()) {
                builder.addTagAddress(u.toString(), i.getItemId());
            }
        }); 
        PlcReadRequest readRequest = builder.build();        
        try {        
            PlcReadResponse syncResponse = readRequest.execute().get();
            syncResponse.getTagNames().forEach(s ->{                
                Collection<Object> values = syncResponse.getAllObjects(s);
                plcitems.get(s).setDataQuality(syncResponse.getResponseCode(s));
                if (syncResponse.getResponseCode(s) == PlcResponseCode.OK) {
                    plcitems.get(s).setPlcValues(syncResponse);
                }
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
        private boolean group_enable = false;
        private long group_period = 100;    
        
        public PlcGroupBuilder(BundleContext bc, String group_name) {
            this.bc = bc;
            this.group_name = group_name;
        }

        public PlcGroupBuilder  setGroupDescription(String group_description) {
            this.group_description = group_description;            
            return this;
        }

        public PlcGroupBuilder  setGroupUid(UUID group_uid) {
            this.group_uid = group_uid;
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
