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

import org.apache.plc4x.java.spi.configuration.PlcConfiguration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.Description;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;

public class Limits implements PlcConfiguration {

    @ConfigurationParameter("receive-buffer-size")
    @IntDefaultValue(65535)
    @Description("Maximum size of received TCP transport message chunk value in bytes.")
    private int receiveBufferSize;
    @ConfigurationParameter("send-buffer-size")
    @IntDefaultValue(65535)
    @Description("Maximum size of sent transport message chunk.")
    private int sendBufferSize;
    @ConfigurationParameter("max-message-size")
    @IntDefaultValue(2097152)
    @Description("Maximum size of complete message.")
    private int maxMessageSize;
    @ConfigurationParameter("max-chunk-count")
    @IntDefaultValue(64)
    @Description("Maximum number of chunks for both sent and received messages.")
    private int maxChunkCount;

    public Limits() {
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
