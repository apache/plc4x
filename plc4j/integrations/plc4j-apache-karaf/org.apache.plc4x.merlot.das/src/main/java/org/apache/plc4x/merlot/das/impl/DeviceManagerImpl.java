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
import org.apache.plc4x.merlot.das.core.DriverAttributes;
import org.apache.plc4x.merlot.das.core.DriverLoader;
import org.apache.plc4x.merlot.das.core.DriverMatcher;
import org.apache.plc4x.merlot.das.core.NamedThreadFactory;
import org.apache.plc4x.merlot.das.core.Util;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Dictionary;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Properties;
import java.util.Set;
import java.util.concurrent.Callable;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import org.osgi.framework.Bundle;
import org.osgi.framework.BundleContext;
import org.osgi.framework.Filter;
import org.osgi.framework.FrameworkEvent;
import org.osgi.framework.FrameworkListener;
import org.osgi.framework.ServiceReference;
import org.osgi.service.device.Constants;
import org.osgi.service.device.Device;
import org.osgi.service.device.Driver;
import org.osgi.service.device.DriverLocator;
import org.osgi.service.device.DriverSelector;
import org.osgi.service.device.Match;
import org.osgi.service.log.LogService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


/**
 * This class represents the Apache Felix implementation of the device access specification.
 * It is based on version 1.1 of the spec.
 * 
 * 
 * @author <a href="mailto:dev@felix.apache.org">Felix Project Team</a>
 */
public class DeviceManagerImpl implements DeviceManager {

    private static final Logger LOGGER = LoggerFactory.getLogger(DeviceManagerImpl.class);  
    
    private final long DEFAULT_TIMEOUT_SEC = 1;

    // the logger
    private volatile LogService m_log;

    // the bundle context
    private final BundleContext m_context;

    // the driver selector
    private volatile DriverSelector m_selector;

    // the driver locators
    private List<DriverLocator> m_locators;

    // the devices
    private Map<ServiceReference, Object> m_devices;

    // the drivers
    private Map<ServiceReference, DriverAttributes> m_drivers;

    // performs all the background actions
    private ExecutorService m_worker;

    // used to add delayed actions
    private ScheduledExecutorService m_delayed;

    //the devices filter
    private Filter m_deviceImplFilter;

    //the drivers filter
    private Filter m_driverImplFilter;


    /**
     * Public constructor. Used by the Activator in this <code>Bundle</code>
     * to instantiate one instance.
     * 
     * @param context the <code>BundleContext</code>
     */
    public DeviceManagerImpl( BundleContext context )
    {
        m_context = context;
    }

    // dependency manager methods

    /**
     *
     * @throws Exception
     */
    @Override
    public void init() throws Exception
    {
        LOGGER.info( "Init device manager" );
        m_locators = Collections.synchronizedList( new ArrayList<DriverLocator>() );
        m_worker = Executors.newSingleThreadExecutor( new NamedThreadFactory( "Merlot Device Manager" ) );
        m_delayed = Executors.newScheduledThreadPool( 1, new NamedThreadFactory(
            "Merlot Device Manager - delayed" ) );
        m_deviceImplFilter = Util.createFilter( "(%s=%s)", new Object[]
            { org.osgi.framework.Constants.OBJECTCLASS, Device.class.getName() } );
        m_driverImplFilter = Util.createFilter( "(%s=%s)", new Object[]
            { org.osgi.framework.Constants.OBJECTCLASS, Driver.class.getName() } );
        
        m_drivers = new HashMap<ServiceReference, DriverAttributes>();
        m_devices = new HashMap<ServiceReference, Object>();
        submit( new WaitForStartFramework() );
    }

    public void destroy()
    {
        LOGGER.info( "Destroying device manager" );
        m_worker.shutdownNow();
        m_delayed.shutdownNow();
    }

    public void start()
    {
        LOGGER.info( "Starting device manager" );
        // nothing to do ?
    }

    public void stop()
    {
        LOGGER.info( "Stopping device manager" );
        // nothing to do ?
    }


    // callback methods

    public void selectorAdded( DriverSelector selector )
    {
        m_selector = selector;
        LOGGER.debug( "driver selector appeared" );
    }


    public void selectorRemoved( DriverSelector selector )
    {
        m_selector = null;
        LOGGER.debug( "driver selector lost" );
    }
    
    public void locatorAdded( DriverLocator locator )
    {
        m_locators.add( locator );
        LOGGER.debug( "driver locator appeared" );
    }


    public void locatorRemoved( DriverLocator locator )
    {
        m_locators.remove( locator );
        LOGGER.debug( "driver locator lost" );
    }


    public void driverAdded( ServiceReference ref, Object obj )
    {
        final Driver driver = Driver.class.cast( obj );
        m_drivers.put( ref, new DriverAttributes( ref, driver ) );
        LOGGER.debug( "driver appeared: " + Util.showDriver( ref ) );
        
        //immediately check for idle devices
        submit( new CheckForIdleDevices() );
    }

    public void driverModified( ServiceReference ref, Object obj )
    {
        final Driver driver = Driver.class.cast( obj );
        
        LOGGER.debug( "driver modified: " + Util.showDriver( ref ) );
        m_drivers.remove( ref );
        m_drivers.put( ref , new DriverAttributes( ref, driver ) );

        // check if devices have become idle
        // after some time
        schedule( new CheckForIdleDevices() );
    }

    public void driverRemoved( ServiceReference ref )
    {
        LOGGER.debug( "driver lost: " + Util.showDriver( ref ) );
        m_drivers.remove( ref );

        // check if devices have become idle
        // after some time
        schedule( new CheckForIdleDevices() );

    }


    public void deviceAdded( ServiceReference ref, Object device )
    {
        m_devices.put( ref, device );
        LOGGER.debug( "device appeared: " + Util.showDevice( ref ) );
        submit( new DriverAttachAlgorithm( ref, device ) );
    }


    public void deviceModified( ServiceReference ref, Object device )
    {
        LOGGER.debug( "device modified: " + Util.showDevice( ref ) );
        // nothing further to do ?
        // DeviceAttributes da = m_devices.get(ref);
        // submit(new DriverAttachAlgorithm(da));
    }


    public void deviceRemoved( ServiceReference ref )
    {
        LOGGER.debug( "device removed: " + Util.showDevice( ref ) );
        m_devices.remove( ref );
        // nothing further to do ?
        // the services that use this
        // device should track it.
    }


    /**
     * perform this task as soon as possible.
     * 
     * @param task
     *            the task
     */
    private void submit( Callable<Object> task )
    {
        m_worker.submit( new LoggedCall( task ) );
    }


    /**
     * perform this task after the default delay.
     * 
     * @param task
     *            the task
     */
    private void schedule( Callable<Object> task )
    {
        m_delayed.schedule( new DelayedCall( task ), DEFAULT_TIMEOUT_SEC, TimeUnit.SECONDS );
    }

    @Override
    public Map<ServiceReference, DriverAttributes> getDrivers() {
        return this.m_drivers;
    }

    @Override
    public Map<ServiceReference, Object> getDevices() {
        return this.m_devices;
    }
    
    // worker callables

    /**
     * Callable used to start the DeviceManager. It either waits (blocking the
     * worker thread) for the framework to start, or if it has already started,
     * returns immediately, freeing up the worker thread.
     * 
     * @author <a href="mailto:dev@felix.apache.org">Felix Project Team</a>
     */
    private class WaitForStartFramework implements Callable<Object>, FrameworkListener
    {

        private final CountDownLatch m_latch = new CountDownLatch( 1 );


        public Object call() throws Exception
        {
            boolean addedAsListener = false;
            if ( m_context.getBundle( 0 ).getState() == Bundle.ACTIVE )
            {
                m_latch.countDown();
                LOGGER.debug( "Starting Device Manager immediately" );
            }
            else
            {
                m_context.addFrameworkListener( this );
                addedAsListener = true;
                LOGGER.debug( "Waiting for framework to start" );
            }

            m_latch.await();
            for ( Map.Entry<ServiceReference, Object> entry : m_devices.entrySet() )
            {
                submit( new DriverAttachAlgorithm( entry.getKey(), entry.getValue() ) );
            }
            // cleanup
            if ( addedAsListener )
            {
                m_context.removeFrameworkListener( this );
            }
            return null;
        }


        // FrameworkListener method
        public void frameworkEvent( FrameworkEvent event )
        {
            switch ( event.getType() )
            {
                case FrameworkEvent.STARTED:
                    LOGGER.debug( "Framework has started" );
                    m_latch.countDown();
                    break;
            }
        }


        @Override
        public String toString()
        {
            return getClass().getSimpleName();
        }
    }

    private class LoggedCall implements Callable<Object>
    {

        private final Callable<Object> m_call;


        public LoggedCall( Callable<Object> call )
        {
            m_call = call;
        }


        private String getName()
        {
            return m_call.getClass().getSimpleName();
        }


        public Object call() throws Exception
        {

            try
            {
                return m_call.call();
            }
            catch ( Exception e )
            {
                LOGGER.error( "call failed: " + getName(), e );
                throw e;
            }
            catch ( Throwable e )
            {
                LOGGER.error( "call failed: " + getName(), e );
                throw new RuntimeException( e );
            }
        }

    }

    private class DelayedCall implements Callable<Object>
    {

        private final Callable<Object> m_call;


        public DelayedCall( Callable<Object> call )
        {
            m_call = call;
        }


        private String getName()
        {
            return m_call.getClass().getSimpleName();
        }


        public Object call() throws Exception
        {
            LOGGER.info( "Delayed call: " + getName() );
            return m_worker.submit( m_call );
        }
    }

    /**
     * Checks for Idle devices, and attaches them
     * 
     * @author <a href="mailto:dev@felix.apache.org">Felix Project Team</a>
     */
    private class CheckForIdleDevices implements Callable<Object>
    {

        public Object call() throws Exception
        {
            LOGGER.debug( "START - check for idle devices" );
            for ( ServiceReference ref : getIdleDevices() )
            {
                LOGGER.info( "IDLE: " + ref.getBundle().getSymbolicName() );
                submit( new DriverAttachAlgorithm( ref, m_devices.get( ref ) ) );
            }

            submit( new IdleDriverUninstallAlgorithm() );
            LOGGER.debug( "STOP - check for idle devices" );
            return null;
        }


        /**
         * get a list of all idle devices.
         * 
         * @return
         */
        private List<ServiceReference> getIdleDevices()
        {
            List<ServiceReference> list = new ArrayList<ServiceReference>();

            for ( ServiceReference ref : m_devices.keySet() )
            {
                LOGGER.info( "checking if idle: " + ref.getBundle().getSymbolicName() );

                final Bundle[] usingBundles = ref.getUsingBundles();
                for ( Bundle bundle : usingBundles )
                {
                    if ( isDriverBundle( bundle ) )
                    {
                        LOGGER.info( "used by driver: " + bundle.getSymbolicName() );
                        LOGGER.debug( "not idle: " + ref.getBundle().getSymbolicName() );
                        break;
                    }
                    
                    list.add( ref );

                }
            }
            return list;
        }
    }


    private boolean isDriverBundle( Bundle bundle )
    {
        ServiceReference[] refs = bundle.getRegisteredServices();
        
        if (refs == null) {
            return false;
        }
        
        for ( ServiceReference ref : refs )
        {
            if ( m_driverImplFilter.match( ref ) )
            {
                return true;
            }
        }
        return false;
    }

    /**
     * 
     * Used to uninstall unused drivers
     * 
     * @author <a href="mailto:dev@felix.apache.org">Felix Project Team</a>
     */
    private class IdleDriverUninstallAlgorithm implements Callable<Object>
    {

        public Object call() throws Exception
        {

            LOGGER.info( "cleaning driver cache" );
            for ( DriverAttributes da : m_drivers.values() )
            {
                // just call the tryUninstall; the da itself
                // will know if it should really uninstall the driver.
                try 
                { 
                	da.tryUninstall(); 
                }
                catch (Exception e) 
                {
                	LOGGER.debug(da.getDriverId() + " uninstall failed");
                }
            }

            return null;
        }
    }

    private class DriverAttachAlgorithm implements Callable<Object>
    {

        private final ServiceReference m_ref;

        private final Device m_device;

        private List<DriverAttributes> m_included;

        private List<DriverAttributes> m_excluded;

        private final DriverLoader m_driverLoader;

        private DriverAttributes m_finalDriver;


        public DriverAttachAlgorithm( ServiceReference ref, Object obj )
        {
            m_ref = ref;
            if ( m_deviceImplFilter.match( ref ) )
            {
                m_device = Device.class.cast( obj );
            }
            else
            {
                m_device = null;
            }

            m_driverLoader = new DriverLoader( LOGGER, m_context );
        }


        @SuppressWarnings("all")
        private Dictionary createDictionary( ServiceReference ref )
        {
            final Properties p = new Properties();

            for ( String key : ref.getPropertyKeys() )
            {
                p.put( key, ref.getProperty( key ) );
            }
            return p;
        }


        @SuppressWarnings("all")
        public Object call() throws Exception
        {
            LOGGER.info( "finding suitable driver for: " + Util.showDevice( m_ref ) );

            final Dictionary dict = createDictionary( m_ref );

            // first create a copy of all the drivers that are already there.
            // during the process, drivers will be added, but also excluded.
            m_included = new ArrayList<DriverAttributes>( m_drivers.values() );
            m_excluded = new ArrayList<DriverAttributes>();

            // first find matching driver bundles
            // if there are no driver locators
            // we'll have to do with the drivers that were
            // added 'manually'
            Set<String> driverIds = m_driverLoader.findDrivers( m_locators, dict );

            // remove the driverIds that are already available
            for ( DriverAttributes da : m_drivers.values() )
            {
                driverIds.remove( da.getDriverId() );
            }
            driverIds.removeAll( m_drivers.keySet() );
            try
            {
            	LOGGER.debug("entering attach phase for " + Util.showDevice( m_ref ) );
                return driverAttachment( dict, driverIds.toArray( new String[0] ) );
            }
            finally
            {
                // unload loaded drivers
                // that were unnecessarily loaded
                m_driverLoader.unload( m_finalDriver );
            }
        }


        @SuppressWarnings("all")
        private Object driverAttachment( Dictionary dict, String[] driverIds ) throws Exception
        {
            m_finalDriver = null;

            // remove the excluded drivers
            m_included.removeAll( m_excluded );

            // now load the drivers
            List<ServiceReference> driverRefs = m_driverLoader.loadDrivers( m_locators, driverIds );
            // these are the possible driver references that have been added
            // add them to the list of included drivers
            for ( ServiceReference serviceReference : driverRefs )
            {
                DriverAttributes da = m_drivers.get( serviceReference );
                if ( da != null )
                {
                    m_included.add( da );
                }
            }

            // now start matching all drivers
            final DriverMatcher mi = new DriverMatcher( LOGGER );

            for ( DriverAttributes driver : m_included )
            {
                try
                {
                    int match = driver.match( m_ref );
                    if ( match <= Device.MATCH_NONE ) 
                    {
                        continue;
                    }
                    mi.add( match, driver );
                }
                catch ( Throwable t )
                {
                    LOGGER.error( "match threw an exception", new Exception( t ) );
                }
            }

            // get the best match
            Match bestMatch = null;

            // local copy
            final DriverSelector selector = m_selector;
            
            if ( selector != null )
            {
                bestMatch = mi.selectBestMatch( m_ref, selector );
                if (bestMatch != null) {
                	LOGGER.debug(String.format("DriverSelector (%s) found best match: %s", selector.getClass().getName(), Util.showDriver(bestMatch.getDriver())));
                }
            }
            
            if (bestMatch == null) 
            {
                bestMatch = mi.getBestMatch();
            }

            if ( bestMatch == null )
            {
                noDriverFound();
                // really return
                return null;
            }

            String driverId = String.class.cast( bestMatch.getDriver().getProperty( Constants.DRIVER_ID ) );

            LOGGER.debug( "best match: " + driverId );
            m_finalDriver = m_drivers.get( bestMatch.getDriver() );

            if ( m_finalDriver == null )
            {
                LOGGER.error( "we found a driverId, but not the corresponding driver: " + driverId);
                noDriverFound();
                return null;
            }

            // here we get serious...
            try
            {
                LOGGER.debug( "attaching to: " + driverId );
                String newDriverId = m_finalDriver.attach( m_ref );
                if ( newDriverId == null )
                {
                    // successful attach
                    return null;
                }
                // its a referral
                LOGGER.info( "attach led to a referral to: " + newDriverId );
                m_excluded.add( m_finalDriver );
                return driverAttachment( dict, new String[]{ newDriverId } );
            }
            catch ( Throwable t )
            {
                LOGGER.error( "attach failed due to an exception", t );
            }
            m_excluded.add( m_finalDriver );
            return driverAttachment( dict, driverIds );
        }


        private void noDriverFound()
        {
            LOGGER.debug( "no suitable driver found for: " + Util.showDevice( m_ref ) );
            if ( m_device != null )
            {
                m_device.noDriverFound();
            }
        }


        @Override
        public String toString()
        {
            return getClass().getSimpleName();// + ": " +
            // Util.showDevice(m_ref);
        }

    }
}
