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
package org.apache.plc4x.test.driver.internal.handlers;

import io.netty.channel.embedded.Plc4xEmbeddedChannel;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.test.driver.internal.utils.ChannelUtil;
import org.apache.plc4x.test.driver.internal.utils.Delay;
import org.dom4j.Element;

public class OutgoingPlcBytesHandler {

    private final Element payload;

    public OutgoingPlcBytesHandler(Element payload) {
        this.payload = payload;
    }

    public void executeOutgoingPlcBytes(Plc4xEmbeddedChannel embeddedChannel, boolean bigEndian) {
        // As we're in the asynchronous world, give the driver some time to respond.
        Delay.shortDelay();
        // Prepare a ByteBuf that contains the data which would have been sent to the PLC.
        final byte[] data = ChannelUtil.getOutboundBytes(embeddedChannel);
        // Validate the data actually matches the expected message.
        validateBytes(payload, data, bigEndian);
    }


    private void validateBytes(Element referenceXml, byte[] data, boolean bigEndian) throws DriverTestsuiteException {
        // TODO: Implement this ...
    }
}
