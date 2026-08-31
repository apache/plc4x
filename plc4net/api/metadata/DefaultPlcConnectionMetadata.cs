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

namespace org.apache.plc4net.api.metadata
{
    /// <summary>
    /// Default implementation of <see cref="IPlcConnectionMetadata"/>
    /// so drivers don't return null.
    /// </summary>
    public class DefaultPlcConnectionMetadata : IPlcConnectionMetadata
    {
        /// <summary>Whether the connection supports reading values.</summary>
        public bool CanRead { get; set; }

        /// <summary>Whether the connection supports writing values.</summary>
        public bool CanWrite { get; set; }

        /// <summary>Whether the connection supports subscriptions.</summary>
        public bool CanSubscribe { get; set; }
    }
}
