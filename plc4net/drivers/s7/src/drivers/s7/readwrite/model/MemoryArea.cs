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
    public enum MemoryArea : byte
    {
        COUNTERS = 0x1C,
        TIMERS = 0x1D,
        DIRECT_PERIPHERAL_ACCESS = 0x80,
        INPUTS = 0x81,
        OUTPUTS = 0x82,
        FLAGS_MARKERS = 0x83,
        DATA_BLOCKS = 0x84,
        INSTANCE_DATA_BLOCKS = 0x85,
        LOCAL_DATA = 0x86,
    }

    public static class MemoryAreaExtensions
    {
        public static string GetShortName(this MemoryArea value) => value switch
        {
            MemoryArea.COUNTERS => "C",
            MemoryArea.TIMERS => "T",
            MemoryArea.DIRECT_PERIPHERAL_ACCESS => "D",
            MemoryArea.INPUTS => "I",
            MemoryArea.OUTPUTS => "Q",
            MemoryArea.FLAGS_MARKERS => "M",
            MemoryArea.DATA_BLOCKS => "DB",
            MemoryArea.INSTANCE_DATA_BLOCKS => "DBI",
            MemoryArea.LOCAL_DATA => "LD",
            _ => default,
        };
    }
}
