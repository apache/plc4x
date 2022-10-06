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
package org.apache.plc4x.java.spi.connection;

import io.netty.channel.Channel;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelPipeline;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;

public interface ChannelFactory {

    Channel createChannel(ChannelHandler channelHandler) throws PlcConnectionException;

    boolean isPassive();

    /** Possibility to add an initial Layer to the Pipeline */
    default void initializePipeline(ChannelPipeline pipeline) {
        // Intentionally do Nothing
    }

    default void closeEventLoopForChannel(Channel channel) {
        // By default do nothing for compatibility
        // Extending classes should implement their logic here
    }

}
