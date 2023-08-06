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


import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.SortedMap;
import java.util.TreeMap;
import org.osgi.framework.Bundle;
import org.osgi.framework.ServiceReference;
import org.osgi.service.device.Constants;
import org.osgi.service.device.DriverSelector;
import org.osgi.service.device.Match;
import org.slf4j.Logger;


/**
 * TODO: add javadoc
 * 
 * @author <a href="mailto:dev@felix.apache.org">Felix Project Team</a>
 */
public class DriverMatcher
{

    private final Logger m_log;

    SortedMap<Integer, List<DriverAttributes>> m_map = new TreeMap<Integer, List<DriverAttributes>>();

    List<Match> m_matches = new ArrayList<Match>();


    public DriverMatcher( Logger log )
    {
        m_log = log;
    }


    // we keep track of the driver attributes in two
    // lists, one to aid us if there is no driver selector, one
    // if there is...
    public void add( Integer match, DriverAttributes value )
    {
        List<DriverAttributes> da = get( match );
        da.add( value );
        m_matches.add( new MatchImpl( value.getReference(), match ) );
    }


    private List<DriverAttributes> get( Integer key )
    {
        List<DriverAttributes> da = m_map.get( key );
        if ( da == null )
        {
            m_map.put( ( Integer ) key, new ArrayList<DriverAttributes>() );
        }
        return m_map.get( key );
    }


    public Match getBestMatch()
    {
        if ( m_map.isEmpty() )
        {
            return null;
        }

        int matchValue = m_map.lastKey();

        // these are the matches that
        // got the highest match value
        List<DriverAttributes> das = m_map.get( matchValue );
        if ( das.size() == 1 )
        {
            // a shortcut: there's only one with the highest match
            return new MatchImpl( das.get( 0 ).getReference(), matchValue );
        }

        // get the highest ranking driver
        final SortedMap<ServiceReference, Match> matches = new TreeMap<ServiceReference, Match>( new ServicePriority() );

        for ( DriverAttributes da : das )
        {
            matches.put( da.getReference(), new MatchImpl( da.getReference(), matchValue ) );
        }

        ServiceReference last = matches.lastKey();
        return matches.get( last );
    }


    public Match selectBestMatch( ServiceReference deviceRef, DriverSelector selector )
    {
//        Match[] matches = m_matches.toArray( new Match[0] );
        
        //(re)check bundle status
        List<Match> activeMatches = new ArrayList<Match>();
        for (Match match : m_matches) {
            if (match.getDriver().getBundle().getState() == Bundle.ACTIVE) {
                activeMatches.add(match);
            } else {
                m_log.debug("skipping: " + match + ", it's bundle is: " + match.getDriver().getBundle().getState());
            }
        }
        try
        {
          Match[] matches = activeMatches.toArray( new Match[0] );
            int index = selector.select( deviceRef, matches );
            if ( index != DriverSelector.SELECT_NONE && index >= 0 && index < matches.length )
            {
                return matches[index];
            }
        }
        catch ( Exception e )
        {
            m_log.error( "exception thrown in DriverSelector.select()", e );
        }
        return null;
    }

    private class MatchImpl implements Match
    {

        private final ServiceReference ref;
        private final int match;


        public MatchImpl( ServiceReference ref, int match )
        {
            this.ref = ref;
            this.match = match;
        }


        public ServiceReference getDriver()
        {
            return ref;
        }


        public int getMatchValue()
        {
            return match;
        }


        public String toString()
        {
            return "[MatchImpl: DRIVER_ID=" + ref.getProperty( Constants.DRIVER_ID ) + ", match=" + match + "]";
        }

    }

    private class ServicePriority implements Comparator<ServiceReference>
    {

        private int getValue( ServiceReference ref, String key, int defaultValue )
        {
            Object obj = ref.getProperty( key );
            if ( obj == null )
            {
                return defaultValue;
            }
            try
            {
                return Integer.class.cast( obj );
            }
            catch ( Exception e )
            {
                return defaultValue;
            }
        }


        public int compare( ServiceReference o1, ServiceReference o2 )
        {
            int serviceRanking1 = getValue( o1, org.osgi.framework.Constants.SERVICE_RANKING, 0 );
            int serviceRanking2 = getValue( o2, org.osgi.framework.Constants.SERVICE_RANKING, 0 );

            if ( serviceRanking1 != serviceRanking2 )
            {
                return ( serviceRanking1 - serviceRanking2 );
            }
            int serviceId1 = getValue( o1, org.osgi.framework.Constants.SERVICE_ID, Integer.MAX_VALUE );
            int serviceId2 = getValue( o2, org.osgi.framework.Constants.SERVICE_ID, Integer.MAX_VALUE );

            return ( serviceId2 - serviceId1 );
        }
    }
}
