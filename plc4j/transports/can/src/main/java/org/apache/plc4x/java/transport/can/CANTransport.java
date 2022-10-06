/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.transport.can;

import io.netty.buffer.ByteBuf;

import java.util.function.Function;
import java.util.function.ToIntFunction;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.generation.MessageInput;
import org.apache.plc4x.java.spi.transport.Transport;

public interface CANTransport<F extends Message> extends Transport {

    ToIntFunction<ByteBuf> getEstimator();

    Class<F> getMessageType();

    MessageInput<F> getMessageInput(Configuration configuration);

    CANFrameBuilder<F> getTransportFrameBuilder();

    Function<F, FrameData> adapter();

    interface FrameHandler<C, T> {
        T fromCAN(FrameData frame);
        C toCAN(T frame);
    }


}
