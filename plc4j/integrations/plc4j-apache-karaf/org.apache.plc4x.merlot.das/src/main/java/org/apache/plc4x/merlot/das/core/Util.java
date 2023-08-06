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


import org.osgi.framework.Constants;
import org.osgi.framework.Filter;
import org.osgi.framework.FrameworkUtil;
import org.osgi.framework.InvalidSyntaxException;
import org.osgi.framework.ServiceReference;
import org.osgi.service.device.Device;


/**
 * TODO: add javadoc
 * 
 * @author <a href="mailto:dev@felix.apache.org">Felix Project Team</a>
 */
public class Util
{

    private Util()
    {
    }


    public static String showDriver( ServiceReference ref )
    {
        Object objectClass = ref.getProperty( Constants.OBJECTCLASS );
        Object driverId = ref.getProperty( org.osgi.service.device.Constants.DRIVER_ID );
        StringBuffer buffer = new StringBuffer();

        buffer.append( "Driver: " );
        buffer.append( Constants.OBJECTCLASS ).append( "=" );

        if ( String[].class.isInstance( objectClass ) )
        {
            buffer.append( enumerateStringArray( String[].class.cast( objectClass ) ) );
        }
        else
        {
            buffer.append( objectClass );
        }
        buffer.append( " " );

        buffer.append( org.osgi.service.device.Constants.DRIVER_ID ).append( "=" );
        buffer.append( driverId );
        return buffer.toString();
    }

    public static boolean isDevice( ServiceReference ref )
    {
        try
        {
            Filter device = createFilter( "(|(%s=%s)(%s=%s))", new Object[]
                { 
            		Constants.OBJECTCLASS, Device.class.getName(),
            		org.osgi.service.device.Constants.DEVICE_CATEGORY, "*" 
            		}
            );
            return device.match( ref );
        }
        catch ( Exception e )
        {
            e.printStackTrace();
        }
        return false;
    }
    
    public static boolean isDeviceInstance( ServiceReference ref )
    {
        try
        {
            Filter device = createFilter( "(%s=%s)", new Object[]
                { Constants.OBJECTCLASS, Device.class.getName() } );
            return device.match( ref );
        }
        catch ( Exception e )
        {
            e.printStackTrace();
        }
        return false;
    }


    public static String createFilterString( String input, Object[] data )
    {
        return String.format( input, data );
    }


    public static Filter createFilter( String input, Object[] data ) throws InvalidSyntaxException
    {
        return FrameworkUtil.createFilter( String.format( input, data ) );
    }


    public static String showDevice( ServiceReference ref )
    {
        Object objectClass = ref.getProperty( Constants.OBJECTCLASS );
        Object category = ref.getProperty( org.osgi.service.device.Constants.DEVICE_CATEGORY );
        StringBuffer buffer = new StringBuffer();

        buffer.append( "Device: " );
        buffer.append( Constants.OBJECTCLASS ).append( "=" );

        appendObject( buffer, objectClass );
        buffer.append( " " );

        buffer.append( org.osgi.service.device.Constants.DEVICE_CATEGORY ).append( "=" );
        appendObject( buffer, category );

        buffer.append( "\n{ " );
        String[] keys = ref.getPropertyKeys();
        
        for ( String key : keys )
        {
            if ( key.equals( Constants.OBJECTCLASS ) )
            {
                continue;
            }
            if ( key.equals( org.osgi.service.device.Constants.DEVICE_CATEGORY ) )
            {
                continue;
            }
            buffer.append( key ).append( "=" );
            appendObject( buffer, ref.getProperty( key ) );
            buffer.append( " " );
        }
        buffer.append( "}\n" );

        return buffer.toString();
    }


    private static void appendObject( StringBuffer buffer, Object obj )
    {
        if ( String[].class.isInstance( obj ) )
        {
            buffer.append( enumerateStringArray( String[].class.cast( obj ) ) );
        }
        else
        {
            buffer.append( obj );
        }
    }


    private static String enumerateStringArray( String[] strings )
    {
        StringBuffer buffer = new StringBuffer();

        buffer.append( "[" );
        for ( String str : strings )
        {
            buffer.append( str );
            buffer.append( " " );
        }
        buffer.deleteCharAt( buffer.length() - 1 );
        buffer.append( "]" );
        return buffer.toString();
    }
}
