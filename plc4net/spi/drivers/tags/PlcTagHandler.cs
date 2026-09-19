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

namespace org.apache.plc4net.spi.drivers.tags
{
    /// <summary>
    /// A driver implements this to parse its protocol-specific tag address syntax.
    /// Mirrors Java SPI3 <c>PlcTagHandler</c>.
    /// </summary>
    /// <remarks>
    /// Each driver registers its own handler. The framework calls
    /// <see cref="ParseTag"/> when the user writes something like
    /// <c>builder.AddTagAddress("myTag", "%DB1.DBW10")</c>.
    /// </remarks>
    public interface PlcTagHandler
    {
        /// <summary>
        /// Parses a tag address string into the driver's concrete tag type.
        /// </summary>
        /// <param name="tagAddress">The address string, e.g. "%DB1.DBW10" for S7
        /// or "coil:0" for Modbus.</param>
        /// <returns>A protocol-specific tag.</returns>
        IPlcTag ParseTag(string tagAddress);
    }
}
