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
package org.apache.plc4x.java.transport.can.virtualcan.config;

import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.defaults.StringDefaultValue;
import org.apache.plc4x.java.transport.can.config.CanTransportConfiguration;

/**
 * Configuration for the Virtual CAN transport.
 * <p>
 * Extends {@link CanTransportConfiguration} with a bus name that identifies which
 * virtual CAN bus this transport instance should connect to. Multiple instances
 * connected to the same bus name can exchange frames in-memory without any
 * underlying hardware or OS-level CAN support.
 */
public class VirtualCanTransportConfiguration extends CanTransportConfiguration {

    /**
     * Name of the virtual CAN bus to connect to.
     * Instances sharing the same bus name will see each other's frames.
     */
    @ConfigurationParameter("bus-name")
    @Description("Name of the virtual CAN bus. Instances on the same bus exchange frames in-memory.")
    @StringDefaultValue("default")
    public String busName = "default";

}
