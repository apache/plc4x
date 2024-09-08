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

package org.apache.plc4x.java.transport.pcapreplay;

import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.Description;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.FloatDefaultValue;
import org.apache.plc4x.java.transport.pcap.DefaultPcapTransportConfiguration;

public abstract class DefaultPcapReplayTransportConfiguration extends DefaultPcapTransportConfiguration implements PcapReplayTransportConfiguration {

    @ConfigurationParameter("replay-speed-factor")
    @FloatDefaultValue(1.0f)
    @Description("Numeric value for changing the replay speed:\n" +
        "- 1 = Normal speed (default)\n" +
        "- 0 = Maximum speed\n" +
        "- 0.5 = Half speed\n" +
        "- 2 = Double speed")
    private float replaySpeedFactor;

    @ConfigurationParameter("loop")
    @BooleanDefaultValue(false)
    @Description("Should the replay start at the beginning of the file as soon as the end is reached?\n" +
        "- true = Automatically start again\n" +
        "- false = Stop at the end (default)")
    private boolean loop;

    @ConfigurationParameter("filter")
    @Description("Filter expression used to filter out unwanted packets from the replay.")
    private String filter;

    @Override
    public float getReplaySpeedFactor() {
        return replaySpeedFactor;
    }

    public void setReplaySpeedFactor(float replaySpeedFactor) {
        this.replaySpeedFactor = replaySpeedFactor;
    }

    @Override
    public boolean isLoop() {
        return loop;
    }

    public void setLoop(boolean loop) {
        this.loop = loop;
    }

    @Override
    public String getFilter() {
        return filter;
    }

    public void setFilter(String filter) {
        this.filter = filter;
    }

}
