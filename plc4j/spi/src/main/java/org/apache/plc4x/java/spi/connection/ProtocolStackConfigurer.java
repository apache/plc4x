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

import io.netty.channel.ChannelPipeline;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.listener.EventListener;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.generation.Message;

import java.util.Collections;
import java.util.List;

public interface ProtocolStackConfigurer<T extends Message> {

    default Plc4xProtocolBase<T> configurePipeline(Configuration configuration, ChannelPipeline pipeline, PlcAuthentication authentication, boolean passive) {
        return configurePipeline(configuration, pipeline, authentication, passive, Collections.emptyList());
    }

    Plc4xProtocolBase<T> configurePipeline(Configuration configuration, ChannelPipeline pipeline, PlcAuthentication authentication, boolean passive, List<EventListener> listeners);

}
