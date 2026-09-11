/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

using org.apache.plc4net.model;

namespace org.apache.plc4net.spi.drivers.messages.items
{
    /// <summary>
    /// A tag item in a request — a named pair of a tag address.
    /// Mirrors Java SPI3 <c>PlcTagItem&lt;T extends PlcTag&gt;</c>.
    /// </summary>
    public interface PlcTagItem<T> where T : IPlcTag
    {
        /// <summary>Name the caller gave this tag.</summary>
        string Name { get; }

        /// <summary>The tag with its protocol-specific address.</summary>
        T Tag { get; }
    }
}
