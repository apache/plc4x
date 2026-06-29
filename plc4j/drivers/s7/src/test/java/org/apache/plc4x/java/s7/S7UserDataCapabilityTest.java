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
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.s7.readwrite.ControllerType;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.EnumSource;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Covers the controller-type → UserData-service capability mapping used when the user pins the
 * controller type and the connect-time SZL probe is therefore skipped (see
 * {@link S7CotpConnection#probeUserDataCapability()}, restored to the pre-SPI3 behaviour).
 */
class S7UserDataCapabilityTest {

    @ParameterizedTest
    @EnumSource(value = ControllerType.class, names = {"S7_300", "S7_400", "S7_1200", "S7_1500"})
    void realCpusSpeakUserDataServices(ControllerType type) {
        assertTrue(S7CotpConnection.supportsUserDataServices(type));
    }

    @ParameterizedTest
    @EnumSource(value = ControllerType.class, names = {"S7_200", "LOGO", "ANY"})
    void lowEndAndUnknownDevicesDoNot(ControllerType type) {
        assertFalse(S7CotpConnection.supportsUserDataServices(type));
    }
}
