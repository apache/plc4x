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
import java.util.Map;
import java.util.UUID;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.Option;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.karaf.shell.support.table.ShellTable;
import org.apache.plc4x.merlot.api.PlcDevice;
import org.apache.plc4x.merlot.api.PlcGroup;
import org.osgi.service.device.Device;
import org.apache.plc4x.merlot.api.PlcGeneralFunction;

@Command(scope = "plc4x", name = "device-list", description = "List al device & groups")
@Service
public class PlcDeviceListCommand implements Action {
  
    @Reference
    PlcGeneralFunction plcservice;   
    
    @Option(name = "-d", aliases = "--did", description = "Device uid.", required = false, multiValued = false)
    String did;   
    
    @Option(name = "-g", aliases = "--gid", description = "Group uid.", required = false, multiValued = false)
    String gid;    
    
    @Override
    public Object execute() throws Exception {
        if ((null == did) && (null == gid))
            printAllDevices();
        if (null != did) {
            printAllGroups();
        } else if (null != gid) {
            printAllItems();            
        }
        return null;
    }
        
    private void printAllDevices() {    
        int[] i = new int[1];
        final Map<UUID, String> devices = plcservice.getPlcDevices("*");
        if (null != devices) {
            ShellTable table = new ShellTable();
            table.column("Uid");
            table.column("Name");  
            table.column("Description"); 
            table.column("Enable");
            table.column("Url");              
            table.column("Groups"); 
            table.column("Items"); 
            
            devices.forEach((u,s) -> {
                final PlcDevice dev = plcservice.getPlcDevice(u);
                i[0] = 0;
                dev.getGroups().forEach(g-> {
                    i[0] += g.getItems().size();
                });                
                table.addRow().addContent(dev.getUid(),
                        dev.getDeviceName(),
                        dev.getDeviceDescription(),
                        dev.isEnable(),
                        dev.getUrl(),
                        dev.getGroups().size(),
                        i[0]);                
            });
            table.print(System.out);  
        } else {
            System.out.println("No device registered.");
        }        
    }
    
    private void printAllGroups() {
        final PlcDevice dev = plcservice.getPlcDevice(UUID.fromString(did));        
        
        if (null != dev) {
            ShellTable table = new ShellTable();
            table.column("Uid");
            table.column("Name");   
            table.column("Description");  
            table.column("Enable");                   
            table.column("Items");                  
            table.column("Transmits"); 
            table.column("Receives");   
            table.column("Errors"); 

            dev.getGroups().forEach(g -> {
                table.addRow().addContent(g.getGroupUid(),
                        g.getGroupName(),
                        g.getGroupDescription(),
                        g.isEnable(),
                        g.getItems().size(),
                        g.getGroupTransmit(),
                        g.getGroupReceives(),
                        g.getGroupErrors());
            });

            table.print(System.out);                           
        }
    }    
        
    private void printAllItems() {
        final PlcGroup g = plcservice.getPlcGroup(UUID.fromString(gid));         
        if (null != g) {       
            ShellTable table = new ShellTable(); 
            table.column("Uid");
            table.column("Name");   
            table.column("Description");  
            table.column("Id");                      
            table.column("Enable");                      
            table.column("Transmits"); 
            table.column("Receives");   
            table.column("Errors"); 
            table.column("LastReadDate"); 
            table.column("LastWriteDate");   
            table.column("LastErrorDate");                     
            g.getItems().forEach(i -> {
                table.addRow().addContent( i.getItemUid(),
                        i.getItemName(),
                        i.getItemDescription(),
                        i.getItemId(),
                        i.isEnable(),
                        i.getItemTransmits(),
                        i.getItemReceives(),
                        i.getItemErrors(),
                        i.getLastReadDate(),
                        i.getLastWriteDate(),
                        i.getLastErrorDate());
            });
            table.print(System.out);                     
        };        
    }       
    
    
}
