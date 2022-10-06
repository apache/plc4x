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
package org.apache.plc4x.java.osgi;

import java.util.ArrayList;
import java.util.List;
import org.apache.plc4x.java.api.PlcDriver;
import org.osgi.framework.BundleActivator;
import org.osgi.framework.BundleContext;
import org.osgi.framework.ServiceRegistration;
import org.osgi.framework.wiring.BundleWiring;

import java.util.Hashtable;
import java.util.ServiceLoader;

public class DriverActivator implements BundleActivator {

    private List<ServiceRegistration<PlcDriver>> registrations = new ArrayList<>();
    public static final String PROTOCOL_NAME = "org.apache.plc4x.driver.name";
    public static final String PROTOCOL_CODE = "org.apache.plc4x.driver.code";


    @Override
    public void start(BundleContext context) throws Exception {
        ServiceLoader<PlcDriver> drivers = ServiceLoader.load(PlcDriver.class, context.getBundle().adapt(BundleWiring.class).getClassLoader());
        for (PlcDriver driver : drivers) {
            Hashtable<String, String> props = new Hashtable<>();
            props.put(PROTOCOL_CODE, driver.getProtocolCode());
            props.put(PROTOCOL_NAME, driver.getProtocolName());
            registrations.add(context.registerService(PlcDriver.class, driver, props));
        }
    }

    @Override
    public void stop(BundleContext context) {
        registrations.forEach(ServiceRegistration::unregister);
        registrations.clear();
    }
}

