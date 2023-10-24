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

import java.util.Hashtable;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.atomic.AtomicReference;
import org.apache.plc4x.java.api.PlcConnection;


public interface PlcGroup {
    
    public static final String	GROUP_DEVICE_UID    = "DEVICE_UID";      
    
    public static final String	GROUP_UID           = "plc4x.group.uid"; 
    
    public static final String	GROUP_NAME          = "scheduler.name";    
    
    public static final String	GROUP_DESCRIPTION   = "plc4x.group.description";  
    
    public static final String	GROUP_PERIOD        = "scheduler.period";  
    
    public static final String	GROUP_IMMEDIATE     = "scheduler.immediate";  
    
    public static final String	GROUP_CONCURRENT    = "scheduler.concurrent";  
    
    /*
    *
    */
    public void enable();
    
    /*
    *
    */
    public void disable(); 
    
    /*
    *
    */
    public boolean isEnable();
      
    /*
    *
    */
    public UUID getGroupUid();    
    
    /*
    *
    */
    public UUID getGroupDeviceUid(); 
    
    /*
    *
    */
    public void setGroupDeviceUid(UUID device_uid);     
    
    /*
    *
    */
    public String getGroupName();
    
    /*
    *
    */
    public void setGroupName(String groupname);
    
    /*
    *
    */
    public String getGroupDescription();
    
    /*
    *
    */
    public void setGroupDescription(String groupdescription);    
    
    /*
    *
    */
    public long getPeriod();
    
    /*
    *
    */
    public void setPeriod(long scantime);

    /*
    *
    */
    public long getGroupTransmit();
    
    /*
    *
    */
    public long getGroupReceives();

    /*
    *
    */
    public long getGroupErrors();
    
    /*
    * 
    */
    public void setPlcConnection(AtomicReference<PlcConnection> plcConnection);
 
    /*
    *
    */
    public Map<UUID, PlcItem> getGroupItems();
    
    /*
    *
    */
    public void setGroupItems(long groupitems);
     
    /*
    *
    */
    public Hashtable<String, Object> getProperties();
     
    /*
    *
    */
    public void putItem(PlcItem item);
    
    /*
    *
    */
    public PlcItem getItem(UUID itemkey);
    
    /*
    *
    */
    public void removeItem(UUID itemkey);
    
    /*
    *
    */
    public List<PlcItem> getItems();
     
}
