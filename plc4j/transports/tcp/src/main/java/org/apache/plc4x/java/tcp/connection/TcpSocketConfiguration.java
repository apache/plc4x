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

package org.apache.plc4x.java.tcp.connection;

import org.apache.plc4x.java.spi.parser.BooleanDefaultValue;
import org.apache.plc4x.java.spi.parser.ConfigurationParameter;
import org.apache.plc4x.java.spi.parser.IntDefaultValue;

/**
 * boolean keepalive = Boolean.parseBoolean(getPropertyOrDefault(SO_KEEPALIVE, "true"));
 * boolean nodelay = Boolean.parseBoolean(getPropertyOrDefault(TCP_NODELAY, "true"));
 * int connectTimeout = Integer.parseInt(getPropertyOrDefault(CONNECT_TIMEOUT_MILLIS, "1000"));
 */
public class TcpSocketConfiguration {

    @ConfigurationParameter("SO_KEEPALIVE")
    @BooleanDefaultValue(true)
    private boolean keepAlive = true;

    @ConfigurationParameter("TCP_NODELAY")
    @BooleanDefaultValue(true)
    private boolean noDelay = true;

    @ConfigurationParameter("CONNECT_TIMEOUT_MILLIS")
    @IntDefaultValue(1000)
    private int connectTimeout = 1000;

    public boolean isKeepAlive() {
        return keepAlive;
    }

    public void setKeepAlive(boolean keepAlive) {
        this.keepAlive = keepAlive;
    }

    public boolean isNoDelay() {
        return noDelay;
    }

    public void setNoDelay(boolean noDelay) {
        this.noDelay = noDelay;
    }

    public int getConnectTimeout() {
        return connectTimeout;
    }

    public void setConnectTimeout(int connectTimeout) {
        this.connectTimeout = connectTimeout;
    }

    @Override
    public String toString() {
        return "TcpSocketConfiguration{" +
            "keepAlive=" + keepAlive +
            ", noDelay=" + noDelay +
            ", connectTimeout=" + connectTimeout +
            '}';
    }

}
