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
package org.apache.plc4x.java.opm.issues.i1935;

import org.apache.plc4x.java.opm.PlcEntity;
import org.apache.plc4x.java.opm.PlcTag;

@PlcEntity
public class FooRxEntity {
    @PlcTag(value = "input-register:1:WORD", cacheDurationMillis = 1000)
    private short bp102;
    @PlcTag(value = "input-register:2:WORD", cacheDurationMillis = 1000)
    private short bp105;
    @PlcTag(value = "input-register:3:WORD", cacheDurationMillis = 1000)
    private short bp112;
    @PlcTag(value = "input-register:4:WORD", cacheDurationMillis = 1000)
    private short bp115;
    @PlcTag(value = "input-register:5:WORD", cacheDurationMillis = 1000)
    private short bp122;
    @PlcTag(value = "input-register:6:WORD", cacheDurationMillis = 1000)
    private short bp125;
    @PlcTag(value = "input-register:7:WORD", cacheDurationMillis = 1000)
    private short bp360;
    @PlcTag(value = "input-register:8:WORD", cacheDurationMillis = 1000)
    private short bp260;
    @PlcTag(value = "input-register:9:WORD", cacheDurationMillis = 1000)
    private short bf102;
    @PlcTag(value = "input-register:10:WORD", cacheDurationMillis = 1000)
    private short bf112;
    @PlcTag(value = "input-register:11:WORD", cacheDurationMillis = 1000)
    private short bf122;
    @PlcTag(value = "input-register:21:WORD", cacheDurationMillis = 1000)
    private short bt220;
    @PlcTag(value = "input-register:22:WORD", cacheDurationMillis = 1000)
    private short bt230;
    @PlcTag(value = "input-register:23:WORD", cacheDurationMillis = 1000)
    private short bt260;
    @PlcTag(value = "input-register:24:WORD", cacheDurationMillis = 1000)
    private short bt320;
    @PlcTag(value = "input-register:25:WORD", cacheDurationMillis = 1000)
    private short bt330;
    @PlcTag(value = "input-register:26:WORD", cacheDurationMillis = 1000)
    private short bt360;
    @PlcTag(value = "holding-register:2049:WORD", cacheDurationMillis = 1000)
    private int bf102Setpoint;
    @PlcTag(value = "holding-register:2050:WORD", cacheDurationMillis = 1000)
    private int bf112Setpoint;
    @PlcTag(value = "holding-register:2051:WORD", cacheDurationMillis = 1000)
    private int bf122Setpoint;
    @PlcTag(value = "holding-register:2057:WORD", cacheDurationMillis = 1000)
    private short qm;

    public short getBp102() {
        return bp102;
    }

    public short getBp105() {
        return bp105;
    }

    public short getBp112() {
        return bp112;
    }

    public short getBp115() {
        return bp115;
    }

    public short getBp122() {
        return bp122;
    }

    public short getBp125() {
        return bp125;
    }

    public short getBp360() {
        return bp360;
    }

    public short getBp260() {
        return bp260;
    }

    public short getBf102() {
        return bf102;
    }

    public short getBf112() {
        return bf112;
    }

    public short getBf122() {
        return bf122;
    }

    public short getBt220() {
        return bt220;
    }

    public short getBt230() {
        return bt230;
    }

    public short getBt260() {
        return bt260;
    }

    public short getBt320() {
        return bt320;
    }

    public short getBt330() {
        return bt330;
    }

    public short getBt360() {
        return bt360;
    }

    public int getBf102Setpoint() {
        return bf102Setpoint;
    }

    public int getBf112Setpoint() {
        return bf112Setpoint;
    }

    public int getBf122Setpoint() {
        return bf122Setpoint;
    }

    public short getQm() {
        return qm;
    }

    public void updateAllTheTags() {
        // Dummy ...
    }

}