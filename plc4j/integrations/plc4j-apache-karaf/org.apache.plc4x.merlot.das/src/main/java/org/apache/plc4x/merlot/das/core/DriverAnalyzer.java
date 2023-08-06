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


import org.osgi.framework.ServiceReference;
import org.osgi.service.device.Constants;
import org.osgi.service.log.LogService;


/**
 * TODO: add javadoc
 * 
 * @author <a href="mailto:dev@felix.apache.org">Felix Project Team</a>
 */
public class DriverAnalyzer
{

    private LogService m_log;


    /**
     * used to analyze invalid drivers
     * 
     * @param ref
     */
    public void driverAdded( ServiceReference ref )
    {
        Object driverId = ref.getProperty( Constants.DRIVER_ID );
        if ( driverId == null || !String.class.isInstance( driverId ) )
        {
            m_log.log( LogService.LOG_ERROR, "invalid driver: no driver id: " + Util.showDriver( ref ) );
            return;
        }
        if ( String.class.isInstance( driverId ) )
        {
            String value = (String)( driverId );
            if ( value.length() == 0 )
            {
                m_log.log( LogService.LOG_ERROR, "invalid driver: empty driver id: " + Util.showDriver( ref ) );
            }
        }
    }
}
