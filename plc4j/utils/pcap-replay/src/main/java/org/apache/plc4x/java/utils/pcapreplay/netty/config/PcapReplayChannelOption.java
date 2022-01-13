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
package org.apache.plc4x.java.utils.pcapreplay.netty.config;

import io.netty.channel.ChannelOption;
import org.apache.plc4x.java.utils.pcap.netty.config.PcapChannelOption;

public class PcapReplayChannelOption extends PcapChannelOption {

    /**
     * Option to increase/decrease the replay speed of the recording.
     * 1.0 being real-time.
     */
    public static final ChannelOption<Float> SPEED_FACTOR =
        ChannelOption.valueOf(Float.class, "SPEED_FACTOR");

    /**
     * Option to tell the playback to automatically restart at the
     * beginning as soon as the end of the playback is reached.
     */
    public static final ChannelOption<Boolean> LOOP =
        ChannelOption.valueOf(Boolean.class, "LOOP");

    /**
     * set a BPF filter
     */
    public static final ChannelOption<String> FILTER =
        ChannelOption.valueOf(String.class, "FILTER");

}
