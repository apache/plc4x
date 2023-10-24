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
package org.apache.plc4x.merlot.api.command;

import java.util.Hashtable;
import java.util.List;
import java.util.UUID;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.Option;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.plc4x.merlot.api.PlcDevice;
import org.apache.plc4x.merlot.scheduler.api.Job;
import org.osgi.framework.BundleContext;
import org.apache.plc4x.merlot.api.PlcGroup;
import org.apache.plc4x.merlot.api.PlcItem;
import org.apache.plc4x.merlot.api.core.PlcItemClientService;
import org.apache.plc4x.merlot.api.impl.PlcGroupImpl;
import org.apache.plc4x.merlot.api.impl.PlcItemImpl;


@Command(scope = "plc4x", name = "demo_002", description = "Command for test.")
@Service
public class DemoCommand002  implements Action  {
  
    @Reference
    BundleContext bc; 
    
    @Reference
    PlcItemClientService items_service;
    
    @Reference
    volatile List<PlcDevice> devices;
    
    @Option(name = "-d", aliases = "--did", description = "Device uid.", required = true, multiValued = false)
    String gid; 
    
    @Option(name = "-n", aliases = "--name", description = "Technological name of the group.", required = true, multiValued = false)
    String group_name;     
    
    @Override
    public Object execute() throws Exception {
        

        
        devices.stream().
                filter(p -> p.getUid().toString().equalsIgnoreCase(gid)).
                forEach(d -> {
                    PlcItem item;
                    PlcGroup grupo1 = new PlcGroupImpl.PlcGroupBuilder(bc, group_name, UUID.randomUUID()).
                                                setGroupPeriod(5000L).
                                                setItemService(items_service).
                                                setGroupDeviceUid(d.getUid()).
                                                build(); 
                    d.putGroup(grupo1);

                    for (int i=0; i<1; i++) {
                    
                        item = new PlcItemImpl.PlcItemBuilder("ITEM_00" + i).
                                    setItemDescription("THE ITEM 00" + i).
                                    setItemId("RANDOM/foo_" +i + ":INT").
                                    setItemEnable(true).
                                    build(); 
                        
                        grupo1.putItem(item);
                    
                    };
                    
                    
                });
                        
        return null;
    }
    
}
