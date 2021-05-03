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
package org.apache.plc4x.test.driver.internal.utils;

import io.netty.buffer.ByteBuf;
import io.netty.channel.embedded.Plc4xEmbeddedChannel;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;

public class ChannelUtil {

    public static final int MAX_TRIES = 500;

    public static byte[] getOutboundBytes(Plc4xEmbeddedChannel embeddedChannel) throws DriverTestsuiteException {
        ByteBuf byteBuf = null;
        for (int i = 0; i < MAX_TRIES; i++) {
            byteBuf = embeddedChannel.readOutbound();
            if (byteBuf != null) {
                break;
            }
            Delay.delay(10);
        }
        if (byteBuf == null) {
            throw new DriverTestsuiteException(String.format("No outbound message available within %dms", 10 * MAX_TRIES));
        }
        final byte[] data = new byte[byteBuf.readableBytes()];
        byteBuf.readBytes(data);
        return data;
    }
}
