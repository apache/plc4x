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
package org.apache.plc4x.java.transport.tcp;

import org.apache.plc4x.java.spi.transport.TransportConfiguration;

/**
 * boolean keepalive = Boolean.parseBoolean(getPropertyOrDefault(SO_KEEPALIVE, "true"));
 * boolean nodelay = Boolean.parseBoolean(getPropertyOrDefault(TCP_NODELAY, "true"));
 * int connectTimeout = Integer.parseInt(getPropertyOrDefault(CONNECT_TIMEOUT_MILLIS, "1000"));
 */
public interface TcpTransportConfiguration extends TransportConfiguration {

    int NO_DEFAULT_PORT = -1;

    default int getDefaultPort() {
        return NO_DEFAULT_PORT;
    }

    default boolean isKeepAlive() {
        return false;
    }

    default boolean isNoDelay() {
        return true;
    }

    default int getConnectTimeout() {
        return 1000;
    }

}
