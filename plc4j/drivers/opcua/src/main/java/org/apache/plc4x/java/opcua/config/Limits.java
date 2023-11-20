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

package org.apache.plc4x.java.opcua.config;

import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;

public class Limits implements Configuration {

    private static final int DEFAULT_RECEIVE_BUFFER_SIZE = 65535;
    private static final int DEFAULT_SEND_BUFFER_SIZE = 65535;
    private static final int DEFAULT_MAX_MESSAGE_SIZE = 2097152;
    private static final int DEFAULT_MAX_CHUNK_COUNT = 64;

    @ConfigurationParameter("receiveBufferSize")
    private int receiveBufferSize;
    @ConfigurationParameter("sendBufferSize")
    private int sendBufferSize;
    @ConfigurationParameter("maxMessageSize")
    private int maxMessageSize;
    @ConfigurationParameter("maxChunkCount")
    private int maxChunkCount;

    public Limits() {
        this(DEFAULT_RECEIVE_BUFFER_SIZE, DEFAULT_SEND_BUFFER_SIZE, DEFAULT_MAX_MESSAGE_SIZE, DEFAULT_MAX_CHUNK_COUNT);
    }

    public Limits(int receiveBufferSize, int sendBufferSize, int maxMessageSize, int maxChunkCount) {
        this.receiveBufferSize = receiveBufferSize;
        this.sendBufferSize = sendBufferSize;
        this.maxMessageSize = maxMessageSize;
        this.maxChunkCount = maxChunkCount;
    }

    public int getReceiveBufferSize() {
        return receiveBufferSize;
    }

    public int getSendBufferSize() {
        return sendBufferSize;
    }

    public int getMaxMessageSize() {
        return maxMessageSize;
    }

    public int getMaxChunkCount() {
        return maxChunkCount;
    }

    @Override
    public String toString() {
        return "Limits{" +
            " receiveBufferSize=" + receiveBufferSize +
            ", sendBufferSize=" + sendBufferSize +
            ", maxMessageSize=" + maxMessageSize +
            ", maxChunkCount=" + maxChunkCount +
            '}';
    }
}
