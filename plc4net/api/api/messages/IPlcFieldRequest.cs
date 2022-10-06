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

using System.Collections.Generic;
using org.apache.plc4net.model;

namespace org.apache.plc4net.messages
{
    /// <summary>
    /// Request for certain fields inside the PLC
    /// </summary>
    public interface IPlcFieldRequest: IPlcRequest
    {
        /// <summary>
        /// Number of fields in the request
        /// </summary>
        int FieldCount { get; }

        /// <summary>
        /// Enumeration of field names
        /// </summary>
        IEnumerable<string> FieldNames { get; }

        /// <summary>
        /// Get a field inside the request by its name
        /// </summary>
        /// <param name="name">Name of the PLC field to return</param>
        /// <returns></returns>
        IPlcField GetFieldByName(string name);

        /// <summary>
        /// Returns all fields inside the request
        /// </summary>
        IEnumerable<IPlcField> Fields { get; }
    }
}