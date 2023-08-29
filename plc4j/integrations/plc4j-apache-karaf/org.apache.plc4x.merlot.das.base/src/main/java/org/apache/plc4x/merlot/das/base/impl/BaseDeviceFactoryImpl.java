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

import org.apache.plc4x.merlot.das.base.api.BaseDeviceFactory;
import org.apache.plc4x.merlot.das.base.core.BaseDeviceManagedService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.apache.plc4x.merlot.api.PlcDevice;
import org.apache.plc4x.merlot.api.PlcDeviceFactory;
import org.apache.plc4x.merlot.api.impl.PlcDeviceImpl;
import org.osgi.framework.BundleContext;

public class BaseDeviceFactoryImpl implements PlcDeviceFactory {
    private static final Logger LOGGER = LoggerFactory.getLogger(BaseDeviceFactory.class); 
    
    private final BundleContext bc;

    public BaseDeviceFactoryImpl(BundleContext bc) {
        this.bc = bc;
    }
    
    @Override
    public PlcDevice create(String device, String url, String name, String description) {
        return new PlcDeviceImpl.PlcDeviceBuilder(bc, url, name, description).build();        

    }
    
}
