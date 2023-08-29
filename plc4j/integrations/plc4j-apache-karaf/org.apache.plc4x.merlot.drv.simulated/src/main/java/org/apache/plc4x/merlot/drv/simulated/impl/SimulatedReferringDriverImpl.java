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
package org.apache.plc4x.merlot.drv.simulated.impl;

import org.apache.plc4x.java.api.PlcDriver;
import org.osgi.framework.BundleContext;
import org.osgi.framework.ServiceReference;
import org.osgi.service.device.Driver;
import org.osgi.service.event.EventAdmin;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class SimulatedReferringDriverImpl implements Driver {
    private static final Logger LOGGER = LoggerFactory.getLogger(SimulatedReferringDriverImpl.class);    
    private BundleContext bc;  
    private final EventAdmin eventAdmin;    
    private PlcDriver plcdriver = null;    

    public SimulatedReferringDriverImpl(BundleContext bc, EventAdmin eventAdmin) {
        System.out.println("Creo el objeto.....");       
        this.bc = bc;
        this.eventAdmin = eventAdmin;
    }

    @Override
    public int match(ServiceReference reference) throws Exception {
        System.out.println("Match not supported yet.");
        return 0;
    }

    @Override
    public String attach(ServiceReference reference) throws Exception {
        System.out.println("Attach not supported yet.");
        return null;
    }
    
   public void bind(ServiceReference reference) {
        System.out.println(">>> Bind ref to simulated driver");
   }    
    
    public void bind(PlcDriver driver){
        System.out.println(">>> Bind simulated driver");
        this.plcdriver = driver;
    }
       
    public void unbind(ServiceReference driver){  
        System.out.println(">>> Unbind ref to simulated driver");        
        this.plcdriver = null;
    }    
    
}
