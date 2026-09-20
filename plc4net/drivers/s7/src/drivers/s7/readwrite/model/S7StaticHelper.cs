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

using System;
using System.Collections.Generic;
using System.Linq;
using org.apache.plc4net.spi.drivers;

namespace org.apache.plc4net.drivers.s7.readwrite.model
{
    /// <summary>
    /// The array-size helper, and the <c>STATIC_CALL</c> targets in
    /// s7.mspec. The generated bodies throw; supply the real
    /// implementation in a sibling non-generated <c>partial</c> file.
    /// </summary>
    public static partial class S7StaticHelper
    {
        /// <summary>Serialized byte length of a sequence of messages.</summary>
        public static int ArraySizeInBytes<T>(IEnumerable<T> items) where T : IMessage
            => items?.Sum(i => i.GetLengthInBytes()) ?? 0;

        public static byte ArrayLength(params object[] args)
            => throw new NotImplementedException("s7.mspec STATIC_CALL 'ArrayLength' has no implementation yet");

        public static byte EventItemLength(params object[] args)
            => throw new NotImplementedException("s7.mspec STATIC_CALL 'EventItemLength' has no implementation yet");

        public static ushort LeftShift3(params object[] args)
            => throw new NotImplementedException("s7.mspec STATIC_CALL 'LeftShift3' has no implementation yet");

        public static ushort RightShift3(params object[] args)
            => throw new NotImplementedException("s7.mspec STATIC_CALL 'RightShift3' has no implementation yet");
    }
}
