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
package org.apache.plc4x.java.ads.api.generic.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.UnsignedShortLEByteValue;

import java.util.EnumSet;
import java.util.LinkedList;
import java.util.List;
import java.util.stream.Stream;

import static java.lang.Integer.toBinaryString;
import static org.apache.plc4x.java.ads.api.generic.types.State.StateMask.*;

/**
 * 2 bytes	see below.
 * <p>
 * State Flags
 * Flag     Description
 * 0x0001	0: Request / 1: Response
 * 0x0004	ADS command
 * The first bit marks, whether it´s a request or response. The third bit must be set to 1, to exchange data with ADS commands. The other bits are not defined or were used for other internal purposes.
 * <p>
 * Therefore the other bits must be set to 0!
 * <p>
 * Flag     Description
 * 0x000x	TCP Protocol
 * 0x004x	UDP Protocol
 * Bit number 7 marks, if it should be transfered with TCP or UDP.
 */
public class State extends UnsignedShortLEByteValue {

    public enum StateMask {
        RESPONSE(0b0000_0000_0000_0001),
        NO_RETURN(0b0000_0000_0000_0010),
        ADS_COMMAND(0b0000_0000_0000_0100),
        SYSTEM_COMMAND(0b0000_0000_0000_1000),
        HIGH_PRIORITY_COMMAND(0b0000_0000_0001_0000),
        TIMESTAMP_ADDED(0b0000_0000_0010_0000),
        UDP_COMMAND(0b0000_0000_0100_0000),
        INIT_COMMAND(0b0000_0000_1000_0000),
        BROADCAST(0b1000_0000_0000_0000);

        private final int mask;

        StateMask(int mask) {
            this.mask = mask;
        }

        public boolean applies(int value) {
            return (value & this.mask) != 0;
        }

        boolean applies(State state) {
            return (state.getAsInt() & this.mask) != 0;
        }

        int getMask() {
            return mask;
        }

        @Override
        public String toString() {
            return name() + "/mask=" + String.format("0b%16s", toBinaryString(mask)).replace(' ', '0');
        }
    }

    public static final State DEFAULT = State.of(ADS_COMMAND);

    public static final State DEFAULT_RESPONSE = State.of(ADS_COMMAND, RESPONSE);

    public static final int NUM_BYTES = UnsignedShortLEByteValue.UNSIGNED_SHORT_LE_NUM_BYTES;

    private State(byte... values) {
        super(values);
    }

    private State(int value) {
        super(value);
    }

    private State(String stateId) {
        super(stateId);
    }

    private State(ByteBuf byteBuf) {
        super(byteBuf);
    }

    public static State of(byte... values) {
        return new State(values);
    }

    public static State of(int value) {
        return new State(value);
    }

    public static State of(ByteBuf byteBuf) {
        return new State(byteBuf);
    }

    public static State of(String stateId) {
        return new State(stateId);
    }

    private static State of(StateMask... stateMasks) {
        return State.of(Stream.of(stateMasks).map(StateMask::getMask).reduce(0, (a, b) -> a | b));
    }

    private EnumSet<StateMask> getStateMaskEnumSet() {
        List<StateMask> result = new LinkedList<>();
        if (isResponse()) {
            result.add(StateMask.RESPONSE);
        }
        if (isNoReturn()) {
            result.add(NO_RETURN);
        }
        if (isADSCommand()) {
            result.add(StateMask.ADS_COMMAND);
        }
        if (isSystemCommand()) {
            result.add(StateMask.SYSTEM_COMMAND);
        }
        if (isHighPriorityCommand()) {
            result.add(StateMask.HIGH_PRIORITY_COMMAND);
        }
        if (isTimestampAdded()) {
            result.add(StateMask.TIMESTAMP_ADDED);
        }
        if (isUDPCommand()) {
            result.add(StateMask.UDP_COMMAND);
        }
        if (isInitCommand()) {
            result.add(StateMask.INIT_COMMAND);
        }
        if (isBroadcast()) {
            result.add(StateMask.BROADCAST);
        }
        return EnumSet.copyOf(result);
    }

    private boolean isResponse() {
        return RESPONSE.applies(this);
    }

    public boolean isRequest() {
        return !isResponse();
    }

    private boolean isNoReturn() {
        return NO_RETURN.applies(this);
    }

    private boolean isADSCommand() {
        return ADS_COMMAND.applies(this);
    }

    private boolean isSystemCommand() {
        return SYSTEM_COMMAND.applies(this);
    }

    private boolean isHighPriorityCommand() {
        return HIGH_PRIORITY_COMMAND.applies(this);
    }

    private boolean isTimestampAdded() {
        return TIMESTAMP_ADDED.applies(this);
    }

    private boolean isUDPCommand() {
        return UDP_COMMAND.applies(this);
    }

    private boolean isInitCommand() {
        return INIT_COMMAND.applies(this);
    }

    private boolean isBroadcast() {
        return BROADCAST.applies(this);
    }

    @Override
    public String toString() {
        return "State{" + getStateMaskEnumSet() + "} " + super.toString();
    }
}
