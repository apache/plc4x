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

import java.util.List;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.karaf.shell.support.table.ShellTable;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.merlot.api.PlcDevice;
import org.apache.plc4x.merlot.api.impl.PlcDriverAdminServiceImpl;
import org.apache.plc4x.merlot.api.PlcGeneralFunction;

@Command(scope = "plc4x", name = "driver-list", description = "List all channels")
@Service
public class PlcDriverListCommand implements Action {
       
    @Reference
    PlcGeneralFunction plcservice;    
    
    @Override
    public Object execute() throws Exception {
        ShellTable table = new ShellTable();
        table.column("Code");
        table.column("Name");        
        
//        drivers.forEach(d -> {
//            table.addRow().addContent( d.getProtocolCode(),
//                    d.getProtocolName());
//        });

        plcservice.getPlcDrivers().forEach((c,d) -> {
            table.addRow().addContent( c,
                    d);            
        });
        
        table.print(System.out);        
        return null;
    }
    
}
