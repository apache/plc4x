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
package org.apache.plc4x.java.s7.configuration;

import org.apache.plc4x.java.s7.readwrite.ControllerType;
import org.apache.plc4x.java.s7.readwrite.DeviceGroup;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class S7ConfigurationTest {

    @Test
    void defaultsMatchAnnotations() {
        S7Configuration cfg = new S7Configuration();
        assertEquals(1, cfg.getLocalRack());
        assertEquals(1, cfg.getLocalSlot());
        assertEquals(DeviceGroup.PG_OR_PC, cfg.getLocalDeviceGroup());
        assertEquals(0, cfg.getLocalTsap());
        assertEquals(0, cfg.getRemoteRack());
        assertEquals(0, cfg.getRemoteSlot());
        assertEquals(DeviceGroup.PG_OR_PC, cfg.getRemoteDeviceGroup());
        assertEquals(0, cfg.getRemoteTsap());
        assertEquals(1024, cfg.getPduSize());
        assertEquals(8, cfg.getMaxAmqCaller());
        assertEquals(8, cfg.getMaxAmqCallee());
        assertEquals(ControllerType.ANY, cfg.getControllerType());
        assertEquals(10000, cfg.getReadTimeout());
        assertEquals(4000, cfg.getHaHeartbeatInterval());
        assertEquals(2000, cfg.getHaFailoverTimeout());
    }

    @Test
    void allSettersRoundTrip() {
        S7Configuration cfg = new S7Configuration();
        cfg.setLocalRack(2);
        cfg.setLocalSlot(3);
        cfg.setLocalDeviceGroup(DeviceGroup.OS);
        cfg.setLocalTsap(0x4711);
        cfg.setRemoteRack(4);
        cfg.setRemoteSlot(5);
        cfg.setRemoteDeviceGroup(DeviceGroup.OTHERS);
        cfg.setRemoteTsap(0x1234);
        cfg.setPduSize(480);
        cfg.setMaxAmqCaller(2);
        cfg.setMaxAmqCallee(2);
        cfg.setControllerType(ControllerType.S7_1200);
        cfg.setReadTimeout(5000);
        cfg.setHaHeartbeatInterval(8000);
        cfg.setHaFailoverTimeout(3000);

        assertEquals(2, cfg.getLocalRack());
        assertEquals(3, cfg.getLocalSlot());
        assertEquals(DeviceGroup.OS, cfg.getLocalDeviceGroup());
        assertEquals(0x4711, cfg.getLocalTsap());
        assertEquals(4, cfg.getRemoteRack());
        assertEquals(5, cfg.getRemoteSlot());
        assertEquals(DeviceGroup.OTHERS, cfg.getRemoteDeviceGroup());
        assertEquals(0x1234, cfg.getRemoteTsap());
        assertEquals(480, cfg.getPduSize());
        assertEquals(2, cfg.getMaxAmqCaller());
        assertEquals(2, cfg.getMaxAmqCallee());
        assertEquals(ControllerType.S7_1200, cfg.getControllerType());
        assertEquals(5000, cfg.getReadTimeout());
        assertEquals(8000, cfg.getHaHeartbeatInterval());
        assertEquals(3000, cfg.getHaFailoverTimeout());
    }
}
