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
package org.apache.plc4x.merlot.das.command;

import org.apache.plc4x.merlot.das.api.DeviceManager;
import java.util.logging.Level;
import java.util.logging.Logger;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Argument;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.Option;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.karaf.shell.support.table.ShellTable;
import org.osgi.framework.BundleContext;
import org.osgi.framework.Constants;
import org.osgi.framework.InvalidSyntaxException;
import org.osgi.framework.ServiceReference;
import org.osgi.service.device.Device;
import org.osgi.service.device.Driver;


@Command(scope = "das", name = "list", description = "List registered Drivers & Device jobs")
@Service
public class ListDevCommand implements Action  {

    @Reference
    BundleContext bc;   
    
    @Reference
    DeviceManager manager;
    
    @Option(name = "--driver", description = "List all driver registered in the device manager.")
    boolean driver;
    
    @Option(name = "--device", description = "List all device registered in the device manager.")
    boolean device;    
    
    @Argument(name = "name", required = false, description = "The Driver or Device short name.")
    String name = null;
    
    @Override
    public Object execute() throws Exception {

        if (manager != null) {
            if (driver) {
                listAllDrivers();
            } else if (device) {
                listAllDevices();                
            }

        }
        
        return null;
    }
    
    public void listAllDrivers() {
        ShellTable table = new ShellTable();
        try {

            table.column("DRIVER_ID");
            table.column("service.bundleid"); 
            table.column("service.id");
            ServiceReference[] references = bc.getServiceReferences(Driver.class.getName(), null);
            if (references != null) {
                if (references.length > 0){
                    for(ServiceReference ref:references){
                        table.addRow().addContent(ref.getProperty("DRIVER_ID"), 
                                                  ref.getProperty("service.bundleid"),
                                                  ref.getProperty("service.id"));
                    }
                }
            }
        } catch (InvalidSyntaxException ex) {
            Logger.getLogger(ListDevCommand.class.getName()).log(Level.SEVERE, null, ex);
        } finally {
            table.print(System.out);
        }
        
    }
    
    public void listAllDevices() {
        ShellTable table = new ShellTable();
        try {
            String filter =  "(&(" + Constants.OBJECTCLASS + "=" + Device.class.getName() + ")" +
                        "(" + org.osgi.service.device.Constants.DEVICE_CATEGORY + "=*))";
                        
            table.column("DEVICE_CATEGORY");
            table.column("DEVICE_DESCRIPTION");
            table.column("DEVICE_FRIENDLY_NAME");
            table.column("DEVICE_SERIAL");
            table.column("service.bundleid");            
            table.column("service.id");
            ServiceReference[] references = bc.getServiceReferences((String) null, filter);

            if (references != null) {
                if (references.length > 0){
                    for(ServiceReference ref:references){
                        table.addRow().addContent(ref.getProperty("DEVICE_CATEGORY"),
                                                  ref.getProperty("DEVICE_DESCRIPTION"),
                                                  ref.getProperty("DEVICE_FRIENDLY_NAME"),
                                                  ref.getProperty("DEVICE_SERIAL"),
                                                  ref.getProperty("service.bundleid"),
                                                  ref.getProperty("service.id"));
                    }
                }
            }
        } catch (InvalidSyntaxException ex) {
            Logger.getLogger(ListDevCommand.class.getName()).log(Level.SEVERE, null, ex);
        } finally {
            table.print(System.out);
        }
        
    }    
    
}
