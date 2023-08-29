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
package org.apache.plc4x.merlot.das.base.impl;

import java.util.HashMap;
import java.util.Map;
import org.apache.plc4x.java.api.PlcDriver;
import org.osgi.framework.ServiceReference;
import org.osgi.service.dal.Function;
import org.osgi.service.dal.OperationMetadata;
import org.osgi.service.dal.PropertyMetadata;
import org.apache.plc4x.merlot.api.PlcDevice;
import org.osgi.framework.BundleContext;


public class BaseFunctionImpl implements Function  {

    protected final BundleContext bc;
    private  Map<String, PlcDevice> plcdevices = new HashMap<String, PlcDevice>();    

    public BaseFunctionImpl(BundleContext bc) {
        this.bc = bc;
    }
        
    public void bind(ServiceReference reference) {
       //
        System.out.println("Registro dispositivo reference...");
    }    
    
    public void bind(PlcDevice plcdevice){
        if (null != plcdevice) {
            System.out.println("Registro dispositivo PlcDevice...");            
            plcdevices.put(plcdevice.getUid().toString(), plcdevice);
        }      
    }
       
    public void unbind(ServiceReference deviceref){  
        System.out.println("DeRegistro dispositivo PlcDevice...");         
        final PlcDevice plcdevice = (PlcDevice) deviceref.getBundle();
        plcdevices.remove(plcdevice.getUid().toString());       
    }     
    
    
    @Override
    public PropertyMetadata getPropertyMetadata(String propertyName) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public OperationMetadata getOperationMetadata(String operationName) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Object getServiceProperty(String propKey) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public String[] getServicePropertyKeys() {
        throw new UnsupportedOperationException("Not supported yet.");
    }
    
}
