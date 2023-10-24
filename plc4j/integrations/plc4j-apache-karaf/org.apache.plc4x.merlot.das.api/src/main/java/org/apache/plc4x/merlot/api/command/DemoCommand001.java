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
import java.util.UUID;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.plc4x.merlot.scheduler.api.Job;
import org.osgi.framework.BundleContext;
import org.apache.plc4x.merlot.api.PlcGroup;
import org.apache.plc4x.merlot.api.impl.PlcGroupImpl;


@Command(scope = "plc4x", name = "demo_001", description = "Command for test.")
@Service
public class DemoCommand001  implements Action  {

//    private String filter_device =  "(&(" + Constants.OBJECTCLASS + "=" + Device.class.getName() + ")" +
//                        "(" + org.osgi.service.device.Constants.DEVICE_SERIAL + "=dummy))";    
    
    @Reference
    BundleContext bc;    
    
    @Override
    public Object execute() throws Exception {
//        Hashtable properties = new Hashtable();
//        properties.putIfAbsent(Constants.DEVICE_CATEGORY, "s7");
//        properties.putIfAbsent(Constants.DEVICE_DESCRIPTION, this);
//        properties.putIfAbsent(Constants.DEVICE_SERIAL, this);
//        
//        PlcDevice plcdevice = new BaseDeviceImpl.BaseDeviceBuilder(bc, "s7://192.168.1.1","NAME","LA DESCRIPCION").build();
//        bc.registerService(PlcDevice.class.getName(), plcdevice, plcdevice.getProperties());
        

        Hashtable props = new Hashtable();
                      
        props.put(PlcGroup.GROUP_NAME, "PRUEBA_DE_SCHEDULER_1");
        props.put(PlcGroup.GROUP_UID, UUID.randomUUID().toString());
        props.put(PlcGroup.GROUP_PERIOD, 5000L);
        props.put(PlcGroup.GROUP_IMMEDIATE, true);
        props.put(PlcGroup.GROUP_CONCURRENT, false); 
        
        PlcGroup grupo1 = new PlcGroupImpl.PlcGroupBuilder(bc, "PRUEBA_DE_SCHEDULER_1", UUID.fromString((String) props.get(PlcGroup.GROUP_UID))).
                                    setGroupPeriod((long) props.get(PlcGroup.GROUP_PERIOD)).                                    
                                    build();
        
        bc.registerService(new String[]{Job.class.getName(), PlcGroup.class.getName()}, grupo1, props);         
        
        props.put(PlcGroup.GROUP_NAME, "PRUEBA_DE_SCHEDULER_2");        
        props.put(PlcGroup.GROUP_UID, UUID.randomUUID().toString());        
        PlcGroup grupo2 = new PlcGroupImpl.PlcGroupBuilder(bc, "PRUEBA_DE_SCHEDULER_2", UUID.fromString((String) props.get(PlcGroup.GROUP_UID))).
                                    setGroupPeriod((long) props.get(PlcGroup.GROUP_PERIOD)).                                    
                                    build();
        
        bc.registerService(new String[]{Job.class.getName(), PlcGroup.class.getName()}, grupo2, props);         
        
        props.put(PlcGroup.GROUP_NAME, "PRUEBA_DE_SCHEDULER_3");
        props.put(PlcGroup.GROUP_UID, UUID.randomUUID().toString());        
        PlcGroup grupo3 = new PlcGroupImpl.PlcGroupBuilder(bc, "PRUEBA_DE_SCHEDULER_3", UUID.fromString((String) props.get(PlcGroup.GROUP_UID))).
                                    setGroupPeriod((long) props.get(PlcGroup.GROUP_PERIOD)).                                    
                                    build();        
        
        bc.registerService(new String[]{Job.class.getName(), PlcGroup.class.getName()}, grupo3, props);         
        
        return null;
    }
    
}
