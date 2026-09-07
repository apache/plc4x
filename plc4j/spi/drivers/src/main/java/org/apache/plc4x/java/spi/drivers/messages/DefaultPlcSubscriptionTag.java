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
package org.apache.plc4x.java.spi.drivers.messages;

import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;

import java.time.Duration;
import java.util.Optional;

public class DefaultPlcSubscriptionTag implements PlcSubscriptionTag {

    private final PlcSubscriptionType subscriptionType;
    private final PlcTag tag;
    private final Duration duration;

    public DefaultPlcSubscriptionTag(PlcSubscriptionType subscriptionType, PlcTag tag, Duration duration) {
        this.subscriptionType = subscriptionType;
        this.tag = tag;
        this.duration = duration;
    }

    @Override
    public PlcSubscriptionType getPlcSubscriptionType() {
        return subscriptionType;
    }

    @Override
    public PlcTag getTag() {
        return tag;
    }

    @Override
    public Optional<Duration> getDuration() {
        return Optional.ofNullable(duration);
    }

    @Override
    public String getAddressString() {
        return tag.getAddressString();
    }

}
