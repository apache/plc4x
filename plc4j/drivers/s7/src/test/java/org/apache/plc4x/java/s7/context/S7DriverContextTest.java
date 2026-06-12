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
package org.apache.plc4x.java.s7.context;

import org.apache.plc4x.java.s7.readwrite.ControllerType;
import org.apache.plc4x.java.s7.configuration.S7Configuration;
import org.apache.plc4x.java.spi.config.ConfigurationFactory;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class S7DriverContextTest {

    @Test
    void defaultsAreSane() {
        S7DriverContext ctx = new S7DriverContext();
        assertEquals(1024, ctx.getPduSize());
        assertEquals(8, ctx.getMaxAmqCaller());
        assertEquals(8, ctx.getMaxAmqCallee());
        assertEquals(ControllerType.ANY, ctx.getControllerType());
        assertEquals(10000, ctx.getReadTimeout());
    }

    @Test
    void applyConfigurationCopiesValues() {
        S7DriverContext ctx = new S7DriverContext();
        S7Configuration cfg = new S7Configuration();
        cfg.setPduSize(960);
        cfg.setMaxAmqCaller(4);
        cfg.setMaxAmqCallee(4);
        cfg.setReadTimeout(5000);
        cfg.setControllerType(ControllerType.S7_1500);
        ctx.applyConfiguration(cfg);
        assertEquals(960, ctx.getPduSize());
        assertEquals(4, ctx.getMaxAmqCaller());
        assertEquals(4, ctx.getMaxAmqCallee());
        assertEquals(5000, ctx.getReadTimeout());
        assertEquals(ControllerType.S7_1500, ctx.getControllerType());
    }

    @Test
    void logoForcesSmallerPdu() {
        S7DriverContext ctx = new S7DriverContext();
        S7Configuration cfg = new S7Configuration();
        cfg.setControllerType(ControllerType.LOGO);
        // Default pdu-size is 1024.
        ctx.applyConfiguration(cfg);
        assertEquals(480, ctx.getPduSize());
        assertEquals(ControllerType.LOGO, ctx.getControllerType());
    }

    @Test
    void unknownControllerTypeRejectedAtConfigParseTime() {
        // Field is enum-typed: ConfigurationFactory rejects unknown values up front
        // rather than silently falling back to ANY (the old String-based behaviour).
        assertThrows(IllegalArgumentException.class, () -> new ConfigurationFactory()
            .createConfiguration(S7Configuration.class, "controller-type=not-a-real-type"));
    }

    @Test
    void settersWork() {
        S7DriverContext ctx = new S7DriverContext();
        ctx.setPduSize(123);
        ctx.setMaxAmqCaller(2);
        ctx.setMaxAmqCallee(3);
        ctx.setReadTimeout(456);
        ctx.setControllerType(ControllerType.S7_300);
        assertEquals(123, ctx.getPduSize());
        assertEquals(2, ctx.getMaxAmqCaller());
        assertEquals(3, ctx.getMaxAmqCallee());
        assertEquals(456, ctx.getReadTimeout());
        assertEquals(ControllerType.S7_300, ctx.getControllerType());
    }
}
