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
package org.apache.plc4x.merlot.das.base.command;

import java.util.List;
import java.util.regex.Pattern;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.karaf.shell.support.table.ShellTable;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.merlot.api.impl.PlcDriverAdminServiceImpl;

//@Command(scope = "plc4x", name = "driver-list", description = "List all Plc4x drivers.")
//@Service
public class cmdDriverList implements Action {
     
    @Reference
    PlcDriverAdminServiceImpl channels;
    
    @Override
    public Object execute() throws Exception {
        ShellTable table = new ShellTable();
        table.column("Code");
        table.column("Name");        
        List<PlcDriver> drivers = channels.getPlcDriverList();
        
        drivers.forEach(drv -> {table.addRow().addContent(drv.getProtocolCode(), drv.getProtocolName());});
        table.print(System.out);        
        return null;
    }
    
}