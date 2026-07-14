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
package org.apache.plc4x.java.slmp.config;

import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;

public class SlmpConfiguration implements Configuration {

    @ConfigurationParameter("monitoring-timer")
    @IntDefaultValue(0x0000)
    @Description("SLMP monitoring timer written into each 3E request frame (0 = wait infinitely).")
    private int monitoringTimer;

    @ConfigurationParameter("request-timeout")
    @IntDefaultValue(5_000)
    @Description("Client-side timeout in milliseconds awaiting a response.")
    private int requestTimeout;

    public int getMonitoringTimer() {
        return monitoringTimer;
    }

    public void setMonitoringTimer(int monitoringTimer) {
        this.monitoringTimer = monitoringTimer;
    }

    public int getRequestTimeout() {
        return requestTimeout;
    }

    public void setRequestTimeout(int requestTimeout) {
        this.requestTimeout = requestTimeout;
    }

    @Override
    public String toString() {
        return "SlmpConfiguration{monitoringTimer=" + monitoringTimer
            + ", requestTimeout=" + requestTimeout + '}';
    }
}
