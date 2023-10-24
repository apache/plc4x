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

import java.util.List;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Argument;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.Option;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.plc4x.java.api.value.PlcValue;
import org.osgi.framework.BundleContext;
import org.apache.plc4x.merlot.api.PlcGroup;
import org.apache.plc4x.merlot.api.PlcItem;
import org.apache.plc4x.merlot.api.PlcItemClient;
import org.apache.plc4x.merlot.api.core.PlcItemClientService;
import org.apache.plc4x.merlot.api.impl.PlcItemImpl;


@Command(scope = "plc4x", name = "demo_004", description = "Prueba Registro glo.")
@Service
public class DemoCommand004  implements Action, PlcItemClient  {
  
    @Reference
    BundleContext bc;  
    
    @Reference
    PlcItemClientService clients;
    
    @Reference
    volatile List<PlcGroup> groups;
    
    @Option(name = "-g", aliases = "--gid", description = "Group uid.", required = true, multiValued = false)
    String gid; 
    
    @Argument(index = 0, name = "name", description = "Item name", required = true, multiValued = false)
    String item_name;    
    
    @Argument(index = 1, name = "desc", description = "Item description", required = true, multiValued = false)
    String item_desc;   
    
    @Argument(index = 2, name = "id", description = "Item id", required = true, multiValued = false)
    String item_id;  
    
    @Argument(index = 3, name = "enbale", description = "Enable the item", required = true, multiValued = false)
    Boolean item_enable = false;     
    
    @Override
    public Object execute() throws Exception {
        
        groups.stream().
                filter(g -> g.getGroupUid().toString().equalsIgnoreCase(gid)).
                forEach(g -> {
                    PlcItem item = new PlcItemImpl.PlcItemBuilder(item_name).
                            setItemDescription(item_desc).
                            setItemId(item_id).
                            setItemEnable(item_enable).
                            build();
                    g.putItem(item);
                    
                    clients.subcription(this, item.getItemUid());
                });
                        
        return null;
    }

    @Override
    public void execute(PlcValue plcvalue) {
        System.out.println("Value: " + plcvalue.getString());
    }
    
}
