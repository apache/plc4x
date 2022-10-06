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

import io.netty.channel.Channel;
import io.netty.channel.ChannelConfig;
import io.netty.channel.ChannelOption;
import io.netty.channel.DefaultChannelConfig;

import java.util.Map;

public class SerialChannelConfig extends DefaultChannelConfig implements ChannelConfig {

    private int baudRate;
    private int dataBits;
    private int stopBits;
    private int parityBits;

    public SerialChannelConfig(Channel channel) {
        super(channel);
    }

    @Override
    public Map<ChannelOption<?>, Object> getOptions() {
        return getOptions(super.getOptions(),
            SerialChannelOptions.BAUD_RATE, SerialChannelOptions.DATA_BITS,
            SerialChannelOptions.STOP_BITS, SerialChannelOptions.PARITY_BITS);
    }

    @Override
    public <T> boolean setOption(ChannelOption<T> option, T value) {
        if(option == SerialChannelOptions.BAUD_RATE) {
            if(value instanceof Integer) {
                baudRate = (Integer) value;
                return true;
            }
            return false;
        } else if(option == SerialChannelOptions.DATA_BITS) {
            if(value instanceof Integer) {
                dataBits = (Integer) value;
                return true;
            }
            return false;
        } else if(option == SerialChannelOptions.STOP_BITS) {
            if(value instanceof Integer) {
                stopBits = (Integer) value;
                return true;
            }
            return false;
        } else if(option == SerialChannelOptions.PARITY_BITS) {
            if(value instanceof Integer) {
                parityBits = (Integer) value;
                return true;
            }
            return false;
        } else {
            return super.setOption(option, value);
        }
    }

    public int getBaudRate() {
        return baudRate;
    }

    public void setBaudRate(int baudRate) {
        this.baudRate = baudRate;
    }

    public int getDataBits() {
        return dataBits;
    }

    public void setDataBits(int dataBits) {
        this.dataBits = dataBits;
    }

    public int getStopBits() {
        return stopBits;
    }

    public void setStopBits(int stopBits) {
        this.stopBits = stopBits;
    }

    public int getParityBits() {
        return parityBits;
    }

    public void setParityBits(int parityBits) {
        this.parityBits = parityBits;
    }

}
