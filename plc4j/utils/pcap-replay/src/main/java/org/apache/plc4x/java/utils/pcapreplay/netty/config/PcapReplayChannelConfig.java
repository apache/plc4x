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
package org.apache.plc4x.java.utils.pcapreplay.netty.config;

import io.netty.channel.Channel;
import io.netty.channel.ChannelOption;
import org.apache.plc4x.java.utils.pcap.netty.config.PcapChannelConfig;

import java.util.Map;

public class PcapReplayChannelConfig extends PcapChannelConfig {

    public static final float SPEED_SLOW_HALF = 0.5f;
    public static final float SPEED_REALTIME = 1f;
    public static final float SPEED_FAST_DOUBLE = 2f;
    public static final float SPEED_FAST_FULL = 0f;

    private float speedFactor = SPEED_REALTIME;
    private boolean loop = false;

    public PcapReplayChannelConfig(Channel channel) {
        super(channel);
    }

    @Override
    public Map<ChannelOption<?>, Object> getOptions() {
        return getOptions(super.getOptions(), PcapReplayChannelOption.SPEED_FACTOR);
    }

    @Override
    public <T> boolean setOption(ChannelOption<T> option, T value) {
        if(option == PcapReplayChannelOption.SPEED_FACTOR) {
            if (value instanceof Float) {
                speedFactor = (Float) value;
                return speedFactor >= 0;
            }
            return false;
        } else if(option == PcapReplayChannelOption.LOOP) {
            if (value instanceof Boolean) {
                loop = (Boolean) value;
                return true;
            }
            return false;
        } else {
            return super.setOption(option, value);
        }
    }

    public float getSpeedFactor() {
        return speedFactor;
    }

    public void setSpeedFactor(float speedFactor) {
        this.speedFactor = speedFactor;
    }

    public boolean isLoop() {
        return loop;
    }

    public void setLoop(boolean loop) {
        this.loop = loop;
    }

}
