/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.profinet.config;

import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.gsdml.ProfinetISO15745Profile;

import java.util.function.BiFunction;

public class ConfigurationProfinetDevice {
    private final String devicename;
    private final String deviceaccess;
    private final String submodules;
    private final BiFunction<String, String, ProfinetISO15745Profile> gsdHandler;
    private String ipaddress;

    public ConfigurationProfinetDevice(String devicename, String deviceaccess, String submodules, BiFunction<String, String, ProfinetISO15745Profile> gsdHandler) {
        this.devicename = devicename;
        this.deviceaccess = deviceaccess;
        this.submodules = submodules;
        this.gsdHandler = gsdHandler;
    }

    public void setIpAddress(String ipaddress) {
        this.ipaddress = ipaddress;
    }

    public String getDevicename() {
        return devicename;
    }

    public String getDeviceaccess() {
        return deviceaccess;
    }

    public String getSubmodules() {
        return submodules;
    }

    public BiFunction<String, String, ProfinetISO15745Profile> getGsdHandler() {
        return gsdHandler;
    }

    public String getIpaddress() {
        return ipaddress;
    }
}
