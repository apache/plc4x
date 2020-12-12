/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.s7.connection.S7PlcConnection;
import org.apache.plc4x.java.spi.PlcDriver;
import org.osgi.service.component.annotations.Component;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Implementation of the S7 protocol, based on:
 * - S7 Protocol
 * - ISO Transport Protocol (Class 0) (https://tools.ietf.org/html/rfc905)
 * - ISO on TCP (https://tools.ietf.org/html/rfc1006)
 * - TCP
 * - Support for R and H systems.
 */
@Component(service = PlcDriver.class, immediate = true)
public class S7PlcDriver implements PlcDriver {

    private static final Pattern S7_URI_PATTERN = Pattern.compile("^s7://(?<host>.*)/(?<rack>\\d{1,4})/(?<slot>\\d{1,4})(?<params>\\?.*)?");
    private static final Pattern S7H_URI_PATTERN = Pattern.compile("^s7://(?<host0>.*):(?<rack0>\\d{1,4}):(?<slot0>\\d{1,4})/(?<host1>.*):(?<rack1>\\d{1,4}):(?<slot1>\\d{1,4})(?<params>\\?.*)?");    

    
    @Override
    public String getProtocolCode() {
        return "s7";
    }

    @Override
    public String getProtocolName() {
        return "Siemens S7 (Basic)";
    }

    @Override
    public PlcConnection connect(String url) throws PlcConnectionException {
        boolean isHSystem = false;
        String params = null;
        String[] hosts = new String[2];
        int[] racks = new int[2];
        int[] slots = new int[2];
        InetAddress[] plcInetAddress = new InetAddress[2];
        
        Matcher matcher = S7_URI_PATTERN.matcher(url);
        Matcher hmatcher = S7H_URI_PATTERN.matcher(url);
        
        isHSystem = hmatcher.matches();
        
        if (!matcher.matches()) {
            if  (!hmatcher.matches()) 
            throw new PlcConnectionException(
                "Connection url doesn't match the format 's7://{host|ip}/{rack}/{slot}' or s7://{host0|ip0}:{rack0}:{slot0}/{host1|ip1}:{rack1}:{slot1}'");
        }
        if (!isHSystem) {
            hosts[0] = matcher.group("host");
            racks[0] = Integer.parseInt(matcher.group("rack"));
            slots[0] = Integer.parseInt(matcher.group("slot"));
            params = matcher.group("params") != null ? matcher.group("params").substring(1) : null;
        } else {
            hosts[0] = hmatcher.group("host0");
            racks[0] = Integer.parseInt(hmatcher.group("rack0"));
            slots[0] = Integer.parseInt(hmatcher.group("slot0"));
            hosts[1] = hmatcher.group("host1");        
            racks[1] = Integer.parseInt(hmatcher.group("rack1"));
            slots[1] = Integer.parseInt(hmatcher.group("slot1"));
            System.out.println(hosts[0] + ":" + racks[0]+ ":" + slots[0]);
            System.out.println(hosts[1] + ":" + racks[1]+ ":" + slots[1]);
            params = hmatcher.group("params") != null ? hmatcher.group("params").substring(1) : null;  
        }
        
        try {
            if (isHSystem) {
                plcInetAddress[0] = InetAddress.getByName(hosts[0]);
                plcInetAddress[1] = InetAddress.getByName(hosts[1]);                
                return new S7PlcConnection(plcInetAddress, racks, slots, params);                              
            } else {
                plcInetAddress[0] = InetAddress.getByName(hosts[0]);
                return new S7PlcConnection(plcInetAddress[0], racks[0], slots[0], params);  
            }
        } catch (UnknownHostException e) {
            throw new PlcConnectionException("Error parsing address", e);
        } catch (Exception e) {
            throw new PlcConnectionException("Error connecting to host", e);
        }
    }

    @Override
    public PlcConnection connect(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("Basic S7 connections don't support authentication.");
    }

}
