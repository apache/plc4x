/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.isoontcp.mina;

import org.apache.mina.core.buffer.IoBuffer;
import org.apache.mina.core.filterchain.IoFilterAdapter;
import org.apache.mina.core.session.IoSession;
import org.apache.mina.core.write.WriteRequest;
import org.apache.mina.core.write.WriteRequestWrapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Filter for embedding data in the ISO on TCP protocol.
 */
public class IsoOnTcpFilterAdapter extends IoFilterAdapter {

    private final static byte ISO_ON_TCP_MAGIC_NUMBER = 0x03;

    private final static Logger logger = LoggerFactory.getLogger(IsoOnTcpFilterAdapter.class);

    @Override
    public void sessionOpened(NextFilter nextFilter, IoSession session) throws Exception {
        nextFilter.sessionOpened(session);
    }

    @Override
    public void filterWrite(NextFilter nextFilter, IoSession session, WriteRequest writeRequest) throws Exception {
        if (writeRequest.getMessage() instanceof IoBuffer) {
            logger.debug("ISO on TCP Message sent: {}", new Object[]{session});

            // At this point of processing all higher levels have already serialized their payload.
            // This data is passed to the lower levels in form of an IoBuffer.
            final IoBuffer userData = (IoBuffer) writeRequest.getMessage();

            int packetSize = userData.remaining() + 4;

            // Output the ISO on TCP header
            final IoBuffer buffer = IoBuffer.allocate(packetSize);
            // Version (is always constant 0x03)
            buffer.put(ISO_ON_TCP_MAGIC_NUMBER);
            // Reserved (is always constant 0x00)
            buffer.put((byte) 0x00);
            // Packet length (including ISOonTCP header)
            // ("remaining" returns the number of bytes left to read in this buffer.
            // It is usually set to a read position of 0 and a limit at the end.
            // So in general remaining is equivalent to a non-existing
            // "userData.size()" method.)
            buffer.putShort((short) packetSize);

            // Output the payload.
            buffer.put(userData);

            // Prepare the buffer for sending.
            buffer.flip();

            // Write the buffer to the output.
            //session.write(buffer);

            // The lower protocol filters only care about sending data, so we replace
            // the payload with the data output by this level.
            writeRequest = new WriteRequestWrapper(writeRequest) {
                @Override
                public Object getMessage() {
                    return buffer;
                }
            };
        }

        // Let the lower level protocol filters output their part.
        nextFilter.filterWrite(session, writeRequest);
    }

    @Override
    public void messageSent(NextFilter nextFilter, IoSession session, WriteRequest writeRequest) throws Exception {
        super.messageSent(nextFilter, session, writeRequest);
    }

    @Override
    public void messageReceived(NextFilter nextFilter, IoSession session, Object message) throws Exception {
        if (message instanceof IoBuffer) {
            logger.debug("ISO on TCP Message received: {}", new Object[]{session});

            final IoBuffer buffer = (IoBuffer) message;

            // The ISO on TCP protocol is really simple and in this case the buffer length
            // will take care of the higher levels not reading more than is in the packet.
            // So we just gobble up the header and continue reading in higher levels.
            if (buffer.get() != ISO_ON_TCP_MAGIC_NUMBER) {
                logger.warn("Expecting ISO on TCP magic number: {}", ISO_ON_TCP_MAGIC_NUMBER);
                return;
            }
            buffer.get();
            // We don't really care about the payload length.
            buffer.getShort();
        }

        // Let the higher level protocol filters parse their part.
        nextFilter.messageReceived(session, message);
    }

}
