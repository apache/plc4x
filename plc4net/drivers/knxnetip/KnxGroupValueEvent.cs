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

using org.apache.plc4net.api.value;

namespace org.apache.plc4net.drivers.knxnetip
{
    /// <summary>
    /// One group telegram observed on the KNX bus.
    /// </summary>
    /// <remarks>
    /// The plc4net SPI has no subscription framework yet, so a KNX bus monitor
    /// registers a plain callback with <see cref="KnxNetIpConnection.RegisterGroupValueListener"/>
    /// instead of building a <c>PlcSubscriptionRequest</c>. Every
    /// <c>GroupValueWrite</c> and <c>GroupValueResponse</c> the gateway forwards
    /// becomes one of these.
    /// </remarks>
    public sealed class KnxGroupValueEvent
    {
        public KnxGroupValueEvent(
            string sourceAddress, string groupAddress, byte[] rawPayload, IPlcValue value)
        {
            SourceAddress = sourceAddress;
            GroupAddress = groupAddress;
            RawPayload = rawPayload;
            Value = value;
        }

        /// <summary>The individual address of the device that sent it, e.g. <c>1.1.5</c>.</summary>
        public string SourceAddress { get; }

        /// <summary>The destination group address, e.g. <c>1/2/3</c>.</summary>
        public string GroupAddress { get; }

        /// <summary>The APDU payload (first byte included), before datapoint decoding.</summary>
        public byte[] RawPayload { get; }

        /// <summary>
        /// The decoded value when a datapoint type was known (a <c>:DPT…</c> tag was
        /// registered for the address), otherwise a <c>PlcRawByteArray</c> of the payload.
        /// </summary>
        public IPlcValue Value { get; }
    }
}
