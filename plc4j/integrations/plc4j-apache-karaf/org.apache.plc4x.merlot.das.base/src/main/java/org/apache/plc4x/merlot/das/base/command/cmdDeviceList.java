/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.apache.plc4x.merlot.das.base.command;

import java.util.List;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.osgi.service.device.Device;

@Command(scope = "plc4x", name = "device-list", description = "List al device.")
@Service
public class cmdDeviceList implements Action {

    @Reference
    List<Device> devices;
    
    @Override
    public Object execute() throws Exception {
        if (null != devices) {
            devices.forEach(dev -> {
                if (dev instanceof org.apache.plc4x.merlot.api.PlcDevice) {
                    System.out.println("Encontro un dispositivo Local");
                }
            });
        } else {
            System.out.println("No device registered.");
        }
        return null;
    }
    
}
