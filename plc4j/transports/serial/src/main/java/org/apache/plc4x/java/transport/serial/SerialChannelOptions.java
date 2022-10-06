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
package org.apache.plc4x.java.transport.serial;

import io.netty.channel.ChannelOption;

public class SerialChannelOptions {

    /**
     * Option to configure the baud rate.
     */
    public static final ChannelOption<Integer> BAUD_RATE =
        ChannelOption.valueOf(Integer.class, "BAUD_RATE");

    /**
     * Option to configure the number of data bits.
     */
    public static final ChannelOption<Integer> DATA_BITS =
        ChannelOption.valueOf(Integer.class, "DATA_BITS");

    /**
     * Option to configure the number of stop bits.
     */
    public static final ChannelOption<Integer> STOP_BITS =
        ChannelOption.valueOf(Integer.class, "STOP_BITS");

    /**
     * Option to configure the number of parity bits.
     */
    public static final ChannelOption<Integer> PARITY_BITS =
        ChannelOption.valueOf(Integer.class, "PARITY_BITS");

}
