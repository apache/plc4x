// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
// https://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied. See the License for the
// specific language governing permissions and limitations
// under the License.

using System;
using org.apache.plc4net.api.value;
using org.apache.plc4net.messages;
using org.apache.plc4net.types;

namespace org.apache.plc4net.spi.drivers.subscriptions
{
    /// <summary>A sampled field value or the response code of a failed read.</summary>
    public sealed class SubscriptionValueEventArgs : IPlcSubscriptionEventArgs
    {
        public SubscriptionValueEventArgs(string name, PlcResponseCode responseCode, IPlcValue? value)
        {
            Name = name;
            ResponseCode = responseCode;
            Value = responseCode == PlcResponseCode.Ok ? value : null;
            Timestamp = DateTime.UtcNow;
        }

        public DateTime Timestamp { get; }
        public string Name { get; }
        public PlcResponseCode ResponseCode { get; }
        public IPlcValue? Value { get; }
    }
}
