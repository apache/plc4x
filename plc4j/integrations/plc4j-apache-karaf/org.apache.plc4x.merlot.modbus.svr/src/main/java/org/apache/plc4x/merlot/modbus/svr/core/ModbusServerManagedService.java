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
package org.apache.plc4x.merlot.modbus.svr.core;

import org.apache.plc4x.merlot.modbus.svr.api.ModbusServer;
import org.apache.plc4x.merlot.modbus.svr.api.ModbusServerMBean;
import org.apache.plc4x.merlot.modbus.svr.impl.ModbusServerImpl;
import org.apache.plc4x.merlot.modbus.svr.impl.ModbusServerMBeanImpl;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.util.Dictionary;
import java.util.Enumeration;
import java.util.Hashtable;
import java.util.logging.Level;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import javax.management.NotCompliantMBeanException;
import org.osgi.framework.BundleContext;
import org.osgi.framework.Constants;
import org.osgi.service.cm.ConfigurationException;
import org.osgi.service.cm.ManagedService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 *
 * @author cgarcia
 */
public class ModbusServerManagedService implements ManagedService {
    private static final Logger LOGGER = LoggerFactory.getLogger(ModbusServerManagedService.class);
    private static final String MODBUS_SERVER_ID = "modbus.svr.id";
    private static final String MODBUS_SERVER_DESC = "modbus.svr.desc";

    private String filter_device = "(&(" + Constants.OBJECTCLASS + "=" + ModbusServer.class.getName() + ")"
            + "(" + MODBUS_SERVER_ID + "=*))";
    //https://mkyong.com/regular-expressions/how-to-validate-ip-address-with-regular-expression/
    private static final Pattern IPADDRESS_PATTERN = Pattern.compile(
		"^([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\." +
		"([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\." +
		"([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\." +
		"([01]?\\d\\d?|2[0-4]\\d|25[0-5])$");
    private static final Pattern IPADDRESS_PORT_PATTERN = Pattern.compile("^(\\d+)");
    private final BundleContext bundleContext;    
        
    public ModbusServerManagedService(BundleContext bundleContext) {
        this.bundleContext = bundleContext;
    }

    @Override
    public void updated(Dictionary props) throws ConfigurationException {
        if (null == props)  return;        
        Enumeration<String> keys = props.keys();
        int port = 502;
        while (keys.hasMoreElements()) {
            String key = keys.nextElement();
            String[] fields = ((String) props.get(key)).split(",");
            if (fields.length == 2){
                String short_description = fields[1];
                String[] sub_fields = fields[0].split("/");               
                if (sub_fields.length > 3) {
                    Matcher matchPortAddress = IPADDRESS_PORT_PATTERN.matcher(sub_fields[sub_fields.length - 1]); 
                    if (matchPortAddress.matches()){
                        port = Integer.parseInt(sub_fields[sub_fields.length - 1]);
                        int size = sub_fields.length-3;
                        LOGGER.info("Modbus server registering to TCP port: " + port);
                        SocketAddress[] sas = new SocketAddress[size];
                        for (int i=2; i<(sub_fields.length-1); i++){
                            LOGGER.info("Modbus server bind address: " + sub_fields[i]);
                            sas[i-2] = new InetSocketAddress(sub_fields[i],port);
                        }
                        
                        ModbusServer mbServer = new ModbusServerImpl();
                        //mbServer.setSocketAddress(sas);
                        mbServer.setBundleContext(bundleContext);
                        mbServer.setPort(port);
                        mbServer.start();

                        Hashtable properties = new Hashtable();
                        properties.put(MODBUS_SERVER_ID, key);                            
                        properties.put(MODBUS_SERVER_DESC, short_description);
                        bundleContext.registerService(ModbusServer.class.getName()
                                ,mbServer
                                ,properties);
                        LOGGER.info("Registered modbus server [{}] [{}]"
                                ,key
                                ,short_description);   

                        Hashtable mbean_props = new Hashtable();
                        ModbusServerMBean msmbean;
                        try {
                            msmbean = new ModbusServerMBeanImpl(mbServer);
                            String strProp  = "org.apache.plc4x.merlot:type=server,name=org.apache.plc4x.modbus.svr,id="+key;
                            mbean_props.put("jmx.objectname", strProp);                        
                            bundleContext.registerService(new String[]{ModbusServerMBean.class.getName()}, msmbean, mbean_props);                              
                        } catch (NotCompliantMBeanException ex) {
                            LOGGER.info(ex.getMessage());
                        }

                        
                        
                    }
                }                
            }
        }        
        
    }

}
