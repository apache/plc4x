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
package org.apache.plc4x.test.driver.internal;

import io.netty.channel.Channel;
import io.netty.channel.embedded.Plc4xEmbeddedChannel;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.connection.ChannelExposingConnection;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;

import java.util.Map;

public class ConnectionManager {

    PlcDriverManager plcDriverManager;

    public ConnectionManager() {
        this.plcDriverManager = new PlcDriverManager();
    }

    public PlcConnection getConnection(String driverName, Map<String, String> driverParameters) throws DriverTestsuiteException {
        try {
            StringBuilder sb = new StringBuilder();
            if (driverParameters != null) {
                for (Map.Entry<String, String> parameter : driverParameters.entrySet()) {
                    sb.append("&").append(parameter.getKey()).append("=").append(parameter.getValue());
                }
            }
            if (sb.length() > 0) {
                sb.replace(0, 1, "?");
            }
            return plcDriverManager.getConnection(driverName + ":test://hurz" + sb);
        } catch (PlcConnectionException e) {
            throw new DriverTestsuiteException("Error loading driver", e);
        }
    }

    public Plc4xEmbeddedChannel getEmbeddedChannel(PlcConnection plcConnection) {
        if (!(plcConnection instanceof ChannelExposingConnection)) {
            throw new PlcRuntimeException("Expecting ChannelExposingConnection");
        }
        ChannelExposingConnection connection = (ChannelExposingConnection) plcConnection;
        Channel channel = connection.getChannel();
        if (!(channel instanceof Plc4xEmbeddedChannel)) {
            throw new PlcRuntimeException("Expecting EmbeddedChannel");
        }
        return (Plc4xEmbeddedChannel) channel;
    }
}
