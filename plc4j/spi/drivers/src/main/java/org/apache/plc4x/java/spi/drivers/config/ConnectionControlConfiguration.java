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
package org.apache.plc4x.java.spi.drivers.config;

import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.defaults.BooleanDefaultValue;

/**
 * SPI-level, cross-driver connection controls parsed from the connection-string parameters.
 *
 * <p>Unlike a driver's protocol {@link Configuration} or a transport configuration, the options here
 * are not protocol- or transport-specific — they govern how the SPI core
 * ({@code DriverBase.getConnection(...)}) establishes <em>any</em> connection. They are parsed from the
 * same connection-string parameter mechanism every other option uses (via
 * {@code ConfigurationFactory.createConfiguration(...)}).</p>
 */
public class ConnectionControlConfiguration implements Configuration {

    /**
     * When {@code true}, the SPI core skips its check that the selected transport is one of the
     * transports the driver declares it supports, allowing a driver to be used with a transport
     * outside its declared supported set (the behavior that existed before that check was added).
     *
     * <p>Defaults to {@code false}: by default the supported-transport check is enforced, so
     * accidentally pairing a driver with a transport it does not support (for example a {@code tcp}
     * transport with the S7 driver, which speaks COTP) fails fast at connect time. Set this to
     * {@code true} only when the non-standard pairing is intentional.</p>
     *
     * <p>This option only bypasses the <em>driver-supported</em> check. It does NOT relax the
     * pre-existing requirement that the transport be a registered/known transport.</p>
     */
    @ConfigurationParameter("allow-unsupported-transport")
    @BooleanDefaultValue(false)
    @Description("When true, allows using a transport that is not in the driver's set of supported "
        + "transports. Defaults to false, in which case using an unsupported transport fails the "
        + "connection attempt. Only bypasses the driver-supported check, not the check that the "
        + "transport is registered at all.")
    public boolean allowUnsupportedTransport;

    /**
     * @return whether the supported-transport check should be skipped for this connection.
     */
    public boolean isAllowUnsupportedTransport() {
        return allowUnsupportedTransport;
    }

}
