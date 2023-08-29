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
package org.apache.plc4x.merlot.das.impl;

import org.apache.plc4x.merlot.das.api.DeviceManager;
import org.apache.plc4x.merlot.das.api.WhiteBoardHandler;

import org.osgi.framework.BundleContext;
import org.osgi.framework.Constants;
import org.osgi.framework.InvalidSyntaxException;
import org.osgi.framework.ServiceReference;
import org.osgi.service.device.Driver;
import org.osgi.service.device.Device;
import org.osgi.util.tracker.ServiceTracker;
import org.osgi.util.tracker.ServiceTrackerCustomizer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * The quartz based implementation of the scheduler.
 *
 */
public class WhiteBoardHandlerImpl implements WhiteBoardHandler {

    /** Default logger. */
    private final Logger LOGGER = LoggerFactory.getLogger(this.getClass());

    private final DeviceManager manager;

    private ServiceTracker<?,?> DriverTracker;
    private ServiceTracker<?,?> DeviceTracker;

    public WhiteBoardHandlerImpl(final BundleContext context, final DeviceManager manager) throws InvalidSyntaxException {
        
        this.manager = manager;
        //(" + Constants.OBJECTCLASS + "=" + org.apache.plc4x.java.api.PlcDriver.class.getName() + ")" +
        this.DriverTracker = new ServiceTracker<>(context,
                context.createFilter("(|(" + Constants.OBJECTCLASS + "=" + Driver.class.getName() + ")" +
                        "("+ org.osgi.service.device.Constants.DRIVER_ID + "=*))"),
                new ServiceTrackerCustomizer<Object,Object>() {

                    public synchronized void  removedService(final ServiceReference reference, final Object service) {
                        context.ungetService(reference);
                        driverRemoved(reference, service);
                    }

                    public synchronized void modifiedService(final ServiceReference reference, final Object service) {
                        driverModified(reference, service);
                    }

                    public synchronized Object addingService(final ServiceReference reference) {
                        final Object obj = context.getService(reference);
                        if ( obj != null ) {
                            driverAdd(reference, obj);
                        }
                        return obj;
                    }
                });
   
        //|((" + Constants.OBJECTCLASS + "=" + org.apache.plc4x.merlot.api.PlcDevice.class.getName() + "))
        this.DeviceTracker = new ServiceTracker<>(context,
                context.createFilter("(|(" + Constants.OBJECTCLASS + "=" + org.apache.plc4x.merlot.api.PlcDevice.class.getName() + ")" +
                        "(" + Constants.OBJECTCLASS + "=" + Device.class.getName() + ")" +
                        "(" + org.osgi.service.device.Constants.DEVICE_CATEGORY + "=*))"),
                new ServiceTrackerCustomizer<Object,Object>() {

                    public synchronized void  removedService(final ServiceReference reference, final Object service) {
                        context.ungetService(reference);
                        deviceRemoved(reference, service);
                    }

                    public synchronized void modifiedService(final ServiceReference reference, final Object service) {
                        deviceModified(reference, service);
                    }

                    public synchronized Object addingService(final ServiceReference reference) {
                        final Object obj = context.getService(reference);
                        if ( obj != null ) {
                            deviceAdd(reference, obj);
                        }
                        return obj;
                    }
                });
                    
    }

      @Override
    public void init() throws Exception {
        this.DriverTracker.open();
        this.DeviceTracker.open();
    }
    
       @Override
    public void destroy() throws Exception {      
        this.DriverTracker.close();
        this.DeviceTracker.close();
    }
    
        @Override
    public void start() {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public void stop() {
        throw new UnsupportedOperationException("Not supported yet.");
    }
    
    /**
     * Register a job or task
     */
    private void driverAdd(final ServiceReference reference, final Object job) {
        this.manager.driverAdded(reference, job);
    }
    
        /**
     * Unregister a service.
     */
    private void driverModified(final ServiceReference reference, final Object service) {
        this.manager.driverModified(reference, service);
    }

    /**
     * Unregister a service.
     */
    private void driverRemoved(final ServiceReference reference, final Object service) {
        this.manager.deviceRemoved(reference);
    }
    

    private void deviceAdd(final ServiceReference reference, final Object job) {
        this.manager.deviceAdded(reference, job);
    }
    
        /**
     * Unregister a service.
     */
    private void deviceModified(final ServiceReference reference, final Object service) {
        this.manager.deviceModified(reference, service);
    }

    /**
     * Unregister a service.
     */
    private void deviceRemoved(final ServiceReference reference, final Object service) {
        this.manager.deviceRemoved(reference);
    }


}
