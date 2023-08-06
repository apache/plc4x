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
package org.apache.plc4x.merlot.modbus.sim.command;

import org.apache.plc4x.merlot.modbus.sim.api.ModbusSim;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.Option;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.osgi.framework.BundleContext;
import org.osgi.framework.ServiceReference;


@Command(scope = "modbus", name = "sim", description = "Admin commands for Modbus sim values.")
@Service
public class ModbusSimAdminCommand implements Action {

    @Reference
    BundleContext bundleContext;    
    
    @Option(name = "-s", aliases = "--start", description = "Start all ModbusSim services installed.", required = false, multiValued = false)
    boolean start;  
    
    @Option(name = "-k", aliases = "--kill", description = "Stop but not uninstall all ModbusSim services.", required = false, multiValued = false)
    boolean kill;    
    
    /*
    @Argument(index = 0, name = "channel", description = "Channel <PV name>", required = true, multiValued = true)
    List<String> strChannels = null;     
    */
    
    @Override
    public Object execute() throws Exception {
        if (kill){
            stopAll();
        } else if (start){
            startAll();
        }
        return null;
    }
    
    private void stopAll(){ 
        try {            
            ServiceReference[] services = bundleContext.getServiceReferences(ModbusSim.class.getName(), null);
            if (services !=null){
                for(ServiceReference service:services){
                    ModbusSim mbsim = (ModbusSim) bundleContext.getService(service);
                    mbsim.stop();
                }
            }
        } catch (Exception ex) {
            ex.printStackTrace();
        }        
    }   
    
    private void startAll(){ 
        try {            
            ServiceReference[] services = bundleContext.getServiceReferences(ModbusSim.class.getName(), null);
            if (services !=null){
                for(ServiceReference service:services){
                    ModbusSim mbsim = (ModbusSim) bundleContext.getService(service);
                    mbsim.start();
                }
            }
        } catch (Exception ex) {
            ex.printStackTrace();
        }        
    }                   
    
}
