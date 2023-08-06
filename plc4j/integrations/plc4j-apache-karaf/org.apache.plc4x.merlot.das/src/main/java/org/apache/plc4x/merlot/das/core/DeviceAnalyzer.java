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
package org.apache.plc4x.merlot.das.core;


import org.osgi.framework.BundleContext;
import org.osgi.framework.Filter;
import org.osgi.framework.InvalidSyntaxException;
import org.osgi.framework.ServiceReference;
import org.osgi.service.device.Constants;
import org.osgi.service.device.Device;
import org.osgi.service.log.LogService;


/**
 * TODO: add javadoc
 * 
 * @author <a href="mailto:dev@felix.apache.org">Felix Project Team</a>
 */
public class DeviceAnalyzer
{

    private LogService m_log;

    private Filter deviceImpl;

    private Filter validCategory;

    private final BundleContext m_context;


    public DeviceAnalyzer( BundleContext context )
    {
        m_context = context;
    }


    @SuppressWarnings("unused")
    private void start() throws InvalidSyntaxException
    {
        String deviceString = Util.createFilterString( "(%s=%s)", new Object[]
            { org.osgi.framework.Constants.OBJECTCLASS, Device.class.getName() } );

        deviceImpl = m_context.createFilter( deviceString );

        String categoryString = Util.createFilterString( "(%s=%s)", new Object[]
            { Constants.DEVICE_CATEGORY, "*" } );

        validCategory = m_context.createFilter( categoryString );
    }


    /**
     * used to analyze invalid devices
     * 
     * @param ref
     */
    public void deviceAdded( ServiceReference ref )
    {

        if ( deviceImpl.match( ref ) )
        {
            return;
        }
        if ( validCategory.match( ref ) )
        {
            Object cat = ref.getProperty( Constants.DEVICE_CATEGORY );
            if ( !String[].class.isInstance( cat ) )
            {
                m_log.log( LogService.LOG_ERROR, "invalid device: invalid device category: " + Util.showDevice( ref ) );
                return;
            }
            if ( String[].class.cast( cat ).length == 0 )
            {
                m_log.log( LogService.LOG_ERROR, "invalid device: empty device category: " + Util.showDevice( ref ) );
            }
        }
        else
        {
            m_log.log( LogService.LOG_ERROR, "invalid device: no device category: " + Util.showDevice( ref ) );
        }
    }
}
