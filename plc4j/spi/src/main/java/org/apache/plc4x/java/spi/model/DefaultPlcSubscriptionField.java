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
package org.apache.plc4x.java.spi.model;

import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionField;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;

import java.time.Duration;
import java.util.Optional;

/**
 * special {@link PlcField} which adds a {@link PlcSubscriptionType}.
 */
public class DefaultPlcSubscriptionField implements PlcSubscriptionField {

    private final PlcSubscriptionType plcSubscriptionType;

    private final PlcField plcField;

    private final Duration duration;

    public DefaultPlcSubscriptionField(PlcSubscriptionType plcSubscriptionType, PlcField plcField, Duration duration) {
        this.plcSubscriptionType = plcSubscriptionType;
        this.plcField = plcField;
        this.duration = duration;
    }

    public PlcSubscriptionType getPlcSubscriptionType() {
        return plcSubscriptionType;
    }

    public PlcField getPlcField() {
        return plcField;
    }

    public Optional<Duration> getDuration() {
        return Optional.ofNullable(duration);
    }

    @Override
    public String toString() {
        return "DefaultPlcSubscriptionField{" +
            "plcSubscriptionType=" + plcSubscriptionType +
            ", plcField=" + plcField +
            ", duration=" + duration +
            '}';
    }
}
