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
package org.apache.plc4x.java.utils.pcapsockets.netty;

import io.netty.channel.ChannelOption;
import org.apache.plc4x.java.utils.pcapsockets.netty.handlers.PacketHandler;

public class PcapSocketChannelOption<T> extends ChannelOption<T> {

    /**
     * Option to restrict the captures based on packet port.
     */
    public static final ChannelOption<Integer> PORT =
        ChannelOption.valueOf(Integer.class, "PORT");

    /**
     * Option to restrict the captures based on TCP protocol ids.
     */
    public static final ChannelOption<Integer> PROTOCOL_ID =
        ChannelOption.valueOf(Integer.class, "PROTOCOL_ID");

    /**
     * Option to increase/decrease the replay speed of the recording.
     * 1.0 being real-time.
     */
    public static final ChannelOption<Float> SPEED_FACTOR =
        ChannelOption.valueOf(Float.class, "SPEED_FACTOR");

    /**
     * Option for providing a PacketHandler, that intercepts the captured packets
     * before passing the data into the channel.
     */
    public static final ChannelOption<PacketHandler> PACKET_HANDLER =
        ChannelOption.valueOf(PacketHandler.class, "PACKET_HANDLER");

    protected PcapSocketChannelOption() {
        super(null);
    }

}
