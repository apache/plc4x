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
package org.apache.plc4x.merlot.modbus.svr.command;

import org.apache.plc4x.merlot.modbus.svr.api.ModbusServer;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.Option;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.osgi.framework.BundleContext;
import org.osgi.framework.ServiceReference;


@Command(scope = "modbus", name = "server", description = "Start/Stop the Modbus server.")
@Service
public class ModbusServerCommand implements Action {

    @Reference
    BundleContext bundleContext;

    @Option(name = "-s", aliases = {"--start"}, description = "Start the modbus server.", required = false, multiValued = false)
    //@Argument(index = 3, name = "length", description = "The device unit identifier.", required = true, multiValued = false)     
    Boolean start = false;

    @Option(name = "-k", aliases = {"--kill"}, description = "Kill the modbus server.", required = false, multiValued = false)
    //@Argument(index = 3, name = "length", description = "The device unit identifier.", required = true, multiValued = false)     
    Boolean kill = false;

    @Option(name = "-h", aliases = {"--host"}, description = "IP address to bind.", required = false, multiValued = false)
    //@Argument(index = 3, name = "length", description = "The device unit identifier.", required = true, multiValued = false)     
    String host = "0.0.0.0";

    @Option(name = "-p", aliases = {"--port"}, description = "TCP/IP port to bind.", required = false, multiValued = false)
    //@Argument(index = 3, name = "length", description = "The device unit identifier.", required = true, multiValued = false)     
    int port = 502;

    @Option(name = "-i", aliases = {"--info"}, description = "General information about the server.", required = false, multiValued = false)
    //@Argument(index = 3, name = "length", description = "The device unit identifier.", required = true, multiValued = false)     
    Boolean info = false;

    public Object execute() throws Exception {
        // TODO Auto-generated method stub
        ServiceReference<?> reference = bundleContext.getServiceReference(ModbusServer.class.getName());
        ModbusServer mbserver = (ModbusServer) bundleContext.getService(reference);
        if (mbserver != null) {
            if (info) {
                System.out.println("Start date:  " + mbserver.getDate());
                System.out.println("Uptime: " + mbserver.getElapseTime());
                System.out.println("Host: " + mbserver.getHost());
                System.out.println("Port: " + mbserver.getPort());
            } else {
                if (kill) {
                    mbserver.stop();
                };
                if (start) {
                    mbserver.start();
                };
                if (port != 502) {
                    mbserver.setPort(port);
                }
                if (!host.equalsIgnoreCase("0.0.0.0")) {
                    //regex here
                    mbserver.setHost(host);
                }

            }
        }
        return null;
    }

}
