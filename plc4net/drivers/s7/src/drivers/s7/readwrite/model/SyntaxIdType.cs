//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

// Code generated from the mspec by plc4net-code-gen. DO NOT EDIT.

namespace org.apache.plc4net.drivers.s7.readwrite.model
{
    public enum SyntaxIdType : byte
    {
        S7ANY = 0x01,
        PBC_ID = 0x13,
        ALARM_LOCKFREESET = 0x15,
        ALARM_INDSET = 0x16,
        ALARM_ACKSET = 0x19,
        ALARM_QUERYREQSET = 0x1A,
        NOTIFY_INDSET = 0x1C,
        NCK = 0x82,
        NCK_METRIC = 0x83,
        NCK_INCH = 0x84,
        DRIVEESANY = 0xA2,
        SYM1200 = 0xB2,
        DBREAD = 0xB0,
    }
}
