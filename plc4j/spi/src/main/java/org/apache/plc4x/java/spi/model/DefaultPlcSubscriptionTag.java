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
package org.apache.plc4x.java.spi.model;

import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.time.Duration;
import java.util.List;
import java.util.Optional;

/**
 * special {@link PlcTag} which adds a {@link PlcSubscriptionType}.
 */
public class DefaultPlcSubscriptionTag implements PlcSubscriptionTag {

    private final PlcSubscriptionType plcSubscriptionType;

    private final PlcTag plcTag;

    private final Duration duration;

    public DefaultPlcSubscriptionTag(PlcSubscriptionType plcSubscriptionType, PlcTag plcTag, Duration duration) {
        this.plcSubscriptionType = plcSubscriptionType;
        this.plcTag = plcTag;
        this.duration = duration;
    }

    @Override
    public String getAddressString() {
        return plcTag.getAddressString();
    }

    @Override
    public PlcValueType getPlcValueType() {
        return plcTag.getPlcValueType();
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return plcTag.getArrayInfo();
    }

    public PlcSubscriptionType getPlcSubscriptionType() {
        return plcSubscriptionType;
    }

    public PlcTag getTag() {
        return plcTag;
    }

    public Optional<Duration> getDuration() {
        return Optional.ofNullable(duration);
    }

    @Override
    public String toString() {
        return "DefaultPlcSubscriptionTag{" +
            "plcSubscriptionType=" + plcSubscriptionType +
            ", plcTag=" + plcTag +
            ", duration=" + duration +
            '}';
    }
}
