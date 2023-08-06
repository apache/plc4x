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



import org.osgi.framework.Bundle;
import org.osgi.framework.BundleException;
import org.osgi.framework.ServiceReference;
import org.osgi.service.device.Constants;
import org.osgi.service.device.Driver;


/**
 * TODO: add javadoc
 * 
 * @author <a href="mailto:dev@felix.apache.org">Felix Project Team</a>
 */
public class DriverAttributes
{

    private final Bundle m_bundle;

    private final ServiceReference m_ref;

    private final Driver m_driver;

    private final boolean m_dynamic;

    public DriverAttributes( ServiceReference ref, Driver driver )
    {
        m_ref = ref;
        m_driver = driver;
        m_bundle = ref.getBundle();
        m_dynamic = m_bundle.getLocation().startsWith( DriverLoader.DRIVER_LOCATION_PREFIX );
    }


    public ServiceReference getReference()
    {
        return m_ref;
    }


    public String getDriverId()
    {
        return m_ref.getProperty( Constants.DRIVER_ID ).toString();
    }


    public int match( ServiceReference ref ) throws Exception
    {
        return m_driver.match( ref );
    }


    /**
     * a driver bundle is idle if it isn't connected to a device bundle.
     * 
     * @return
     */
    private boolean isInUse()
    {
        ServiceReference[] used = m_bundle.getServicesInUse();
        if ( used == null || used.length == 0 )
        {
            return false;
        }

        for ( ServiceReference ref : used )
        {
            if ( Util.isDevice( ref ) )
            {
                return true;
            }
        }
        return false;
    }


    public String attach( ServiceReference ref ) throws Exception
    {
        return m_driver.attach( ref );
    }


    public void tryUninstall() throws BundleException
    {
    	
        // only install if _we_ loaded the driver
        if ( !isInUse() && m_dynamic )
        {
        	m_bundle.uninstall();
        }
    }

}
