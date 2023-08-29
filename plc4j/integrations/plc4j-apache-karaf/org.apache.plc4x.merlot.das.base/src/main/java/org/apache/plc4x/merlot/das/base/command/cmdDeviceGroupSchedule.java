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

import java.util.Hashtable;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.Option;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.plc4x.merlot.scheduler.api.ScheduleOptions;
import org.apache.plc4x.merlot.scheduler.api.Scheduler;
import org.osgi.framework.BundleContext;
import org.osgi.framework.ServiceReference;
import org.apache.plc4x.merlot.api.PlcChannel;
import org.apache.plc4x.merlot.api.PlcGroup;

@Command(scope = "plc4x", name = "group-rt", description = "Delete a device.")
@Service
public class cmdDeviceGroupSchedule implements Action {

    @Reference
    BundleContext bc;
    
    @Reference
    Scheduler scheduler;
    
    @Option(name = "-u", aliases = "--uid", description = "Channel uid.", required = true, multiValued = false)
    String uid; 
    
    @Option(name = "-s", aliases = "--scantime", description = "Channel scan time in ms.", required = true, multiValued = false)
    Long scantime = 100L;    
    
    @Override
    public Object execute() throws Exception {
        Hashtable props = new Hashtable();
        String filter = "(" + PlcGroup.GROUP_UID + "=" + uid + ")";
        props.put(PlcGroup.GROUP_UID, uid);
        
        ServiceReference[] services = bc.getServiceReferences(PlcGroup.class.getName(), filter);
        
        if (services.length > 0) {
            String[] keys = services[0].getPropertyKeys();
            
            String jobkey = services[0].getProperty(PlcGroup.GROUP_NAME) + "." +
                            services[0].getProperty("service.id");
            ScheduleOptions options = scheduler.NOW(-1, scantime);
            
            scheduler.reschedule(jobkey, options);
            
        }
        
        return null;
    }
    
}
