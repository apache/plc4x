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
import java.util.UUID;
import java.util.function.Consumer;
import java.util.logging.Level;
import java.util.logging.Logger;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.Option;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.karaf.shell.support.table.ShellTable;
import org.apache.plc4x.merlot.api.PlcGroup;
import org.apache.plc4x.merlot.api.impl.PlcGroupImpl;
import org.apache.plc4x.merlot.scheduler.api.ScheduleOptions;
import org.apache.plc4x.merlot.scheduler.api.Scheduler;
import org.apache.plc4x.merlot.scheduler.api.SchedulerError;
import org.slf4j.LoggerFactory;
import org.apache.plc4x.merlot.api.PlcGeneralFunction;

@Command(scope = "plc4x", name = "group-ctrl", description = "Control state of ")
@Service
public class PlcGroupControlCommand implements Action {
    
    private static final org.slf4j.Logger LOGGER = LoggerFactory.getLogger(PlcGroupControlCommand.class);  
    
    @Reference
    PlcGeneralFunction plcservice; 
    
    @Reference
    Scheduler scheduler;
            
    @Option(name = "-g", aliases = "--gid", description = "Group uid.", required = true, multiValued = false)
    String gid;
    
    @Option(name = "-e", aliases = "--enable", description = "Enable the group.", required = false, multiValued = false)
    Boolean enable = false;  

    @Option(name = "-x", aliases = "--disable", description = "Disable the group.", required = false, multiValued = false)
    Boolean disable = false;
    
    @Option(name = "-s", aliases = "--schedule", description = "Change scan time for the group in ms.", required = false, multiValued = false)
    long period = 0L;      
    
    @Override
    public Object execute() throws Exception {
        UUID grp_uid = UUID.fromString(gid);
        final PlcGroup grp = plcservice.getPlcGroup(grp_uid);
        
        if (null != grp) {
            if (enable) {
                grp.enable();
            } else if (disable){
                grp.disable();
            }
            
            if ((period > 100) && (period != grp.getPeriod())){
                plcservice.setPlcGroupScanRate(grp_uid, period);                  
            }                
        }
            
        return null;
    }

    
}
