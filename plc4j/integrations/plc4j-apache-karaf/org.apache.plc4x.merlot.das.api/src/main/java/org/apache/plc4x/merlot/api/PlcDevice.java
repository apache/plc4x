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
import org.apache.plc4x.java.api.PlcDriver;


public interface PlcDevice extends org.osgi.service.device.Device, org.osgi.service.dal.Device {
	
    /*
    *
    */    
    public void init() throws Exception;
    
    /*
    *
    */    
    public void destroy() throws Exception;
      
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
    public Hashtable<String, ?> getProperties();
    
    /*
    *
    */    
    public String getDeviceName();
    
    /*
    *
    */    
    public void setDeviceName(String groupname);
    
    /*
    *
    */    
    public String getDeviceDescription();
    
    /*
    *
    */    
    public void setDeviceDescription(String groupdescription);   
    
    /*
    *
    */    
    public void setUid(UUID uid);
    
    /*
    *
    */    
    public UUID getUid();    
      
    /*
    *
    */    
    public void setUrl(String url);
    
    /*
    *
    */    
    public String getUrl();      
      
    /*
    *
    */    
    public void attach(PlcDriver driver);
      
    /*
    *
    */    
    public void putGroup(PlcGroup group);
    
    /*
    *
    */    
    public PlcGroup getGroup(UUID uid);
    
    /*
    *
    */    
    public void removeGroup(UUID uid);          

    /*
    *
    */    
    public List<PlcGroup> getGroups();
        	
}
