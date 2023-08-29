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
import static java.util.stream.Collectors.toList;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Argument;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.Option;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.plc4x.java.api.PlcDriver;
import org.osgi.framework.BundleContext;
import org.apache.plc4x.merlot.api.PlcDevice;

@Command(scope = "plc4x", name = "device-del", description = "Delete a device.")
@Service
public class cmdDeviceDel implements Action {

    @Reference
    volatile List<PlcDriver> drivers;
    
    @Reference
    volatile List<PlcDevice> devices;        
    
    @Reference
    BundleContext bc;      
    
    @Option(name = "-n", aliases = "--name", description = "Device name.", required = false, multiValued = false)
    String name;

    @Option(name = "-i", aliases = "--uid", description = "Device identifier (UID). ", required = false, multiValued = false)
    String uuid;    
    
    
    @Override
    public Object execute() throws Exception {
        List<PlcDevice> ctxdevices = null;
        
        if (null != name) {
            ctxdevices = devices.stream().
                            filter(dev -> dev.getDeviceName().equals(name)).
                            collect(toList());
        } else if (null != uuid) {
            ctxdevices = devices.stream().
                            filter(dev -> dev.getUid().equals(UUID.fromString(uuid))).
                            collect(toList());            
        }
        
        if ((null == ctxdevices) || (ctxdevices.isEmpty())) {
            System.out.println("A device not exists.");
            return null;
        }        
        
        System.out.println(">>> " + ctxdevices.get(0).getDeviceName());
        
        return null;
    }
    
}
