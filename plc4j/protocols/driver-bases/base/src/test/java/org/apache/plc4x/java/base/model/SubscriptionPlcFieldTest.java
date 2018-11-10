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

package org.apache.plc4x.java.base.model;

import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.time.Duration;
import java.time.temporal.ChronoUnit;
import java.util.Optional;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;

@ExtendWith(MockitoExtension.class)
class SubscriptionPlcFieldTest {

    @Mock
    private PlcField field;
    private SubscriptionPlcField SUT;

    @BeforeEach
    void setUp() {
        SUT = new SubscriptionPlcField(PlcSubscriptionType.CYCLIC, field, Duration.of(1, ChronoUnit.SECONDS));
    }

    @Test
    void getPlcSubscriptionType() {
        PlcSubscriptionType plcSubscriptionType = SUT.getPlcSubscriptionType();
        assertThat(plcSubscriptionType, equalTo(PlcSubscriptionType.CYCLIC));
    }

    @Test
    void getPlcField() {
        PlcField plcField = SUT.getPlcField();
        assertThat(plcField, equalTo(field));
    }

    @Test
    void getDuration() {
        Optional<Duration> durationOptional = SUT.getDuration();
        assertThat(durationOptional.isPresent(), equalTo(true));
        assertThat(durationOptional.get(), equalTo(Duration.of(1, ChronoUnit.SECONDS)));
    }

}