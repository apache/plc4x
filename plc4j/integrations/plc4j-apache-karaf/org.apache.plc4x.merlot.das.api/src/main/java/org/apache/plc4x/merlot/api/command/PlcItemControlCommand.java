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
package org.apache.plc4x.merlot.api.command;

import java.util.UUID;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.Option;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.plc4x.merlot.api.PlcItem;
import org.apache.plc4x.merlot.api.PlcGeneralFunction;

@Command(scope = "plc4x", name = "item-ctrl", description = "List all channels")
@Service
public class PlcItemControlCommand implements Action {
    
    @Reference
    PlcGeneralFunction plcservice; 
    
    
    @Option(name = "-i", aliases = "--iid", description = "Item uid.", required = true, multiValued = false)
    String iid;
    
    @Option(name = "-e", aliases = "--enable", description = "Enable the device.", required = false, multiValued = false)
    Boolean enable = false;  

    @Option(name = "-x", aliases = "--disable", description = "Disable the device.", required = false, multiValued = false)
    Boolean disable = false;      
    
    @Override
    public Object execute() throws Exception {
        UUID item_uid = UUID.fromString(iid);
        final PlcItem item = plcservice.getPlcItem(item_uid);
        
        if (null != item){
            if (enable) {
                item.enable();
            } else if (disable) {
                item.disable();
            }            
        }
             
        return null;
    }
    
}
