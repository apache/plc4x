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
using org.apache.plc4net.types;

namespace org.apache.plc4net.messages
{
    /// <summary>
    /// Interface for responses to requests realted to
    /// a specific PLC field
    /// </summary>
    public interface IPlcFieldResponse: IPlcResponse
    {
        /// <summary>
        /// Enumeration of fields in the response
        /// </summary>
        IEnumerable<string> FieldNames { get; }

        /// <summary>
        /// Get a field by name
        /// </summary>
        /// <param name="name">Name of the field to retrieve</param>
        /// <returns>Field with the given name</returns>
        IPlcField GetFieldByName(string name);

        /// <summary>
        /// Get the response code from the PLC
        /// </summary>
        PlcResponseCode ResponseCode { get; }
    }
}