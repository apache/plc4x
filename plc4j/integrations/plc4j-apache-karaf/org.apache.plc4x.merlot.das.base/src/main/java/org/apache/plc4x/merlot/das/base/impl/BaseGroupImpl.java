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

import java.util.ArrayList;
import java.util.HashMap;
import java.util.Hashtable;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.atomic.AtomicReference;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.merlot.scheduler.api.Job;
import org.apache.plc4x.merlot.scheduler.api.JobContext;
import org.osgi.framework.BundleContext;
import org.apache.plc4x.merlot.api.PlcItem;
import org.apache.plc4x.merlot.api.PlcGroup;


public class BaseGroupImpl implements PlcGroup, Job {

    private UUID uid;
    
    private String  groupname;
    private String  groupdescription;
    
    private boolean enable = false;
    
    private long scantime;    
    private long deadband;
    
    private long grouptransmit = 0;
    private long groupreceives = 0;    
    private long grouperrors = 0;  
    
    private long groupitems = 0;
    private long groupupdatrate = -1;
     
    private final Hashtable<String, Object> myproperties = new Hashtable();     
    private final Map<UUID, PlcItem> items = new Hashtable();

    private long[] aux = new long[1];

    public BaseGroupImpl(String name) { 
        myproperties.put(PlcGroup.GROUP_NAME, name);
    }

    @Override
    public UUID getGroupDeviceUid() {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public void setGroupDeviceUid(UUID device_uid) {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }
    
    
    public void start(int bc){
        //
    }
    
    public void stop(int bc){
        //
    }

    @Override
    public void enable() {
        //
    }

    @Override
    public void disable() {
        //
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
    public long getGroupTransmit() {  
        aux[0] = 0;
        items.forEach((uid, item) -> {aux[0] = +item.getItemTransmits();});
        return aux[0];
    }

    @Override
    public long getGroupReceives() {
        aux[0] = 0;
        items.forEach((uid, item) -> {aux[0] = +item.getItemReceives();});
        return aux[0];
    }

    @Override
    public long getGroupErrors() {
        aux[0] = 0;
        items.forEach((uid, item) -> {aux[0] = +item.getItemErrors();});
        return aux[0];
    }

    @Override
    public void setGroupItems(long groupitems) {
        throw new UnsupportedOperationException("Not supported yet."); 
    }


           
    @Override
    public void putItem(PlcItem item) {
        items.put(item.getItemUid(), item);
    }

    @Override
    public PlcItem getItem(UUID itemuid) {
        return items.get(itemuid);
    }

    //TODO: remove from context
    @Override
    public void removeItem(UUID itemuid) {
        items.remove(itemuid);
    }

    @Override
    public List<PlcItem> getItems() {
        return null;
    }

    @Override
    public void execute(JobContext context) {
        if (enable) {
            System.out.println("Tamaño de la coleccion: " + items.size());
        }
    }

    @Override
    public long getPeriod() {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public void setPeriod(long scantime) {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public Map<UUID, PlcItem> getGroupItems() {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public Hashtable<String, Object> getProperties() {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }

    @Override
    public void setPlcConnection(AtomicReference<PlcConnection> plcConnection) {
        throw new UnsupportedOperationException("Not supported yet."); // Generated from nbfs://nbhost/SystemFileSystem/Templates/Classes/Code/GeneratedMethodBody
    }


    
}
