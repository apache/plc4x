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

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{
    public enum AccessLevel : byte
    {
        Level0 = 0x0,
        Level1 = 0x1,
        Level2 = 0x2,
        Level3 = 0x3,
        Level15 = 0xF,
    }

    public static class AccessLevelExtensions
    {
        public static string GetPurpose(this AccessLevel value) => value switch
        {
            AccessLevel.Level0 => "system manufacturer",
            AccessLevel.Level1 => "product manufacturer",
            AccessLevel.Level2 => "configuration",
            AccessLevel.Level3 => "end-user",
            AccessLevel.Level15 => "read access",
            _ => default,
        };
        public static bool GetNeedsAuthentication(this AccessLevel value) => value switch
        {
            AccessLevel.Level0 => true,
            AccessLevel.Level1 => true,
            AccessLevel.Level2 => true,
            AccessLevel.Level3 => false,
            AccessLevel.Level15 => false,
            _ => default,
        };
    }
}
