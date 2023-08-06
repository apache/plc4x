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

//import com.ceos.merlot.api.DriverEvent;
//import com.ceos.merlot.core.Merlot.FUNCTION;
import org.apache.plc4x.merlot.das.api.DeviceManager;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Argument;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.Option;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.osgi.framework.BundleContext;
import org.osgi.framework.ServiceReference;
import org.osgi.service.device.Device;


@Command(scope = "das", name = "read", description = "List registered Drivers & Device jobs")
@Service
public class ReadDevCommand  implements Action  {

    @Reference
    BundleContext bc;   
    
    @Reference
    DeviceManager manager; 
    
    @Argument( index = 0, name = "device", description = "The Device short name.", required = true,
               multiValued = false)
    String device = null;
 
    @Argument(index = 1, name = "item",  description = "Item to read from Device.",required = false)
    String item = null;
    
    @Option(name = "-n", description = "Number of times.", required = false, multiValued = false)
    protected int n;
    
    @Option(name = "-c", description = "Commans supported by Merlot.", required = false, multiValued = false)
    boolean commands;

    @Override
    public Object execute() throws Exception {

//        String strFilter = "(" + org.osgi.service.device.Constants.DEVICE_DESCRIPTION + "=" + device + ")";
//        //Filter filter = bc.createFilter(strFilter);
//        ServiceReference[] references = bc.getAllServiceReferences(Device.class.getName() , strFilter);
//        if (references != null){
//            com.ceos.merlot.api.Device bsd = (com.ceos.merlot.api.Device) bc.getService(references[0]);
//            for (int i=0;i<10;i++) {
//                final DriverEvent event = bsd.getEvent();
//                event.setFunctionCode(FUNCTION.FC_READ_DATA_BYTES);
//                event.setTransactionID((short) i);
//                event.setCallback((cb)->{System.out.println(cb.toString());});
//                bsd.putEvent(event);            
//            }
//        }            


        return null;
    }

}
