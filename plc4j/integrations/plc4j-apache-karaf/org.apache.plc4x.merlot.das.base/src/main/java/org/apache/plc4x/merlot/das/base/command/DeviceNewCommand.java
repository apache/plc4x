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
import java.util.UUID;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import static java.util.stream.Collectors.toList;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Argument;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.plc4x.java.api.PlcDriver;
import org.osgi.framework.BundleContext;
import org.apache.plc4x.merlot.api.PlcDevice;
import org.apache.plc4x.merlot.api.impl.PlcDeviceImpl;

@Command(scope = "plc4x", name = "device-new", description = "Create a new device.")
@Service
public class DeviceNewCommand implements Action {
    private static final Pattern p = Pattern.compile("^([a-zA-Z0-9]+)(:)(.*)"); 
    
    @Reference
    volatile List<PlcDriver> drivers;
    
    @Reference
    volatile List<PlcDevice> devices;    

    @Reference
    BundleContext bc;    
    
    @Argument(index = 0, name = "Url", description = "URL of device", required = true, multiValued = false)
    String url;     
    
    @Argument(index = 1, name = "Name", description = "Name of device.", required = true, multiValued = false)
    String name;

    @Argument(index = 2, name = "Description", description = "Description of device", required = true, multiValued = false)
    String description;  
             
    @Override
    public Object execute() throws Exception {
        Matcher m = p.matcher(url);
        UUID uuid = UUID.randomUUID();

        m.matches();
        String drivercode = m.group(1);

        List<PlcDriver> devdrivers = drivers.stream().
                                        filter(drv -> drv.getProtocolCode().equals(drivercode)).
                                        collect(toList());
        
        List<PlcDevice> ctxdevices = devices.stream().
                                    filter(dev -> dev.getDeviceName().equals(name)).
                                    collect(toList());
        
        if (!ctxdevices.isEmpty()) {
            System.out.println("A device with the same name exists.");
            return null;
        }
        
        if (!devdrivers.isEmpty()) {
             
            PlcDevice device =  new PlcDeviceImpl.PlcDeviceBuilder(bc, 
                    url,
                    name,
                    description).
                    build();
            
           
            bc.registerService(PlcDevice.class.getName(), device, device.getProperties());
            
            System.out.println("UID: " + uuid.toString());
        } else {
            System.out.println("Driver not found.");            
        }
        
        
        return null;
    }
    
}
