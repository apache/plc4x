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
    public enum VirtualLanPriority : byte
    {
        BEST_EFFORT = 0x0,
        BACKGROUND = 0x1,
        EXCELLENT_EFFORT = 0x2,
        CRITICAL_APPLICATIONS = 0x3,
        VIDEO = 0x4,
        VOICE = 0x5,
        INTERNETWORK_CONTROL = 0x6,
        NETWORK_CONTROL = 0x7,
    }

    public static class VirtualLanPriorityExtensions
    {
        public static string GetAcronym(this VirtualLanPriority value) => value switch
        {
            VirtualLanPriority.BEST_EFFORT => "BE",
            VirtualLanPriority.BACKGROUND => "BK",
            VirtualLanPriority.EXCELLENT_EFFORT => "EE",
            VirtualLanPriority.CRITICAL_APPLICATIONS => "CA",
            VirtualLanPriority.VIDEO => "VI",
            VirtualLanPriority.VOICE => "VO",
            VirtualLanPriority.INTERNETWORK_CONTROL => "IC",
            VirtualLanPriority.NETWORK_CONTROL => "NC",
            _ => default,
        };
    }
}
