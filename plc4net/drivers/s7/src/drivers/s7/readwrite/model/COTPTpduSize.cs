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
    public enum COTPTpduSize : byte
    {
        SIZE_128 = 0x07,
        SIZE_256 = 0x08,
        SIZE_512 = 0x09,
        SIZE_1024 = 0x0a,
        SIZE_2048 = 0x0b,
        SIZE_4096 = 0x0c,
        SIZE_8192 = 0x0d,
    }

    public static class COTPTpduSizeExtensions
    {
        public static ushort GetSizeInBytes(this COTPTpduSize value) => value switch
        {
            COTPTpduSize.SIZE_128 => 128,
            COTPTpduSize.SIZE_256 => 256,
            COTPTpduSize.SIZE_512 => 512,
            COTPTpduSize.SIZE_1024 => 1024,
            COTPTpduSize.SIZE_2048 => 2048,
            COTPTpduSize.SIZE_4096 => 4096,
            COTPTpduSize.SIZE_8192 => 8192,
            _ => default,
        };
    }
}
