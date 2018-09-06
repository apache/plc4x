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
package org.apache.plc4x.java.ads.model;

public enum AdsDataType {
    // https://infosys.beckhoff.com/english.php?content=../content/1033/tcsystemmanager/basics/TcSysMgr_DatatypeComparison.htm&id=
    BIT(1),
    BIT8(1),
    BITARR8(1),
    BITARR16(2),
    BITARR32(4),
    INT8(1),
    INT16(2),
    INT32(4),
    INT64(8),
    UINT8(1),
    UINT16(2),
    UINT32(4),
    UINT64(8),
    FLOAT(4),
    DOUBLE(8),
    // https://infosys.beckhoff.com/english.php?content=../content/1033/tcplccontrol/html/tcplcctrl_plc_data_types_overview.htm&id
    BOOL(0),
    BYTE(0),
    WORD(0),
    DWORD(0),
    SINT(0),
    USINT(0),
    INT(0),
    UINT(0),
    DINT(0),
    UDINT(0),
    LINT(0),
    ULINT(0),
    REAL(0),
    LREAL(0),
    STRING(0),
    TIME(0),
    TIME_OF_DAY(0),
    DATE(0),
    DATE_AND_TIME(0);

    private final int tagetByteSize;

    AdsDataType(int tagetByteSize) {
        this.tagetByteSize = tagetByteSize;
    }

    public int getTagetByteSize() {
        return tagetByteSize;
    }
}
