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

/*
 * TODO: Decide if distinguishment between date, time and datetime are needed (see plc4j)
 * TODO: Decide if implementation of byte[] and Byte[] as in plc4j are needed
 * TODO: Check type compatibility between Java BigDecimal and Decimal
 */

using System;

namespace org.apache.plc4net.messages
{
    /// <summary>
    /// Interface for a builder that constructs write requests
    /// </summary>
    public interface IPlcWriteRequestBuilder: IPlcRequestBuilder<IPlcWriteRequest>
    {
        /// <summary>
        /// Add item with bool value
        /// </summary>
        /// <param name="name">Identifying name of the field</param>
        /// <param name="fieldQuery">Query string for the field</param>
        /// <param name="values">Values to set</param>
        /// <returns>Own instance to allow Fluid API calls</returns>
        IPlcWriteRequestBuilder AddItem(string name, string fieldQuery, params bool[] values);
        
        /// <summary>
        /// Add item with byte value
        /// </summary>
        /// <param name="name">Identifying name of the field</param>
        /// <param name="fieldQuery">Query string for the field</param>
        /// <param name="values">Values to set</param>
        /// <returns>Own instance to allow Fluid API calls</returns>
        IPlcWriteRequestBuilder AddItem(string name, string fieldQuery, params byte[] values);

        /// <summary>
        /// Add item with short value
        /// </summary>
        /// <param name="name">Identifying name of the field</param>
        /// <param name="fieldQuery">Query string for the field</param>
        /// <param name="values">Values to set</param>
        /// <returns>Own instance to allow Fluid API calls</returns>
        IPlcWriteRequestBuilder AddItem(string name, string fieldQuery, params short[] values);
        
        /// <summary>
        /// Add item with int value
        /// </summary>
        /// <param name="name">Identifying name of the field</param>
        /// <param name="fieldQuery">Query string for the field</param>
        /// <param name="values">Values to set</param>
        /// <returns>Own instance to allow Fluid API calls</returns>
        IPlcWriteRequestBuilder AddItem(string name, string fieldQuery, params int[] values);

        /// <summary>
        /// Add item with long value
        /// </summary>
        /// <param name="name">Identifying name of the field</param>
        /// <param name="fieldQuery">Query string for the field</param>
        /// <param name="values">Values to set</param>
        /// <returns>Own instance to allow Fluid API calls</returns>
        IPlcWriteRequestBuilder AddItem(string name, string fieldQuery, params long[] values);

        /// <summary>
        /// Add item with float value
        /// </summary>
        /// <param name="name">Identifying name of the field</param>
        /// <param name="fieldQuery">Query string for the field</param>
        /// <param name="values">Values to set</param>
        /// <returns>Own instance to allow Fluid API calls</returns>
        IPlcWriteRequestBuilder AddItem(string name, string fieldQuery, params float[] values);

        /// <summary>
        /// Add item with double value
        /// </summary>
        /// <param name="name">Identifying name of the field</param>
        /// <param name="fieldQuery">Query string for the field</param>
        /// <param name="values">Values to set</param>
        /// <returns>Own instance to allow Fluid API calls</returns>
        IPlcWriteRequestBuilder AddItem(string name, string fieldQuery, params double[] values);

        /// <summary>
        /// Add item with Decimal value
        /// </summary>
        /// <param name="name">Identifying name of the field</param>
        /// <param name="fieldQuery">Query string for the field</param>
        /// <param name="values">Values to set</param>
        /// <returns>Own instance to allow Fluid API calls</returns>
        IPlcWriteRequestBuilder AddItem(string name, string fieldQuery, params Decimal[] values);

        /// <summary>
        /// Add item with string value
        /// </summary>
        /// <param name="name">Identifying name of the field</param>
        /// <param name="fieldQuery">Query string for the field</param>
        /// <param name="values">Values to set</param>
        /// <returns>Own instance to allow Fluid API calls</returns>
        IPlcWriteRequestBuilder AddItem(string name, string fieldQuery, params string[] values);

        /// <summary>
        /// Add item with DateTime value
        /// </summary>
        /// <param name="name">Identifying name of the field</param>
        /// <param name="fieldQuery">Query string for the field</param>
        /// <param name="values">Values to set</param>
        /// <returns>Own instance to allow Fluid API calls</returns>
        IPlcWriteRequestBuilder AddItem(string name, string fieldQuery, params DateTime[] values);

        /// <summary>
        /// Add item with value of generic type
        /// </summary>
        /// <param name="name">Identifying name of the field</param>
        /// <param name="fieldQuery">Query string for the field</param>
        /// <param name="values">Values to set</param>        
        /// <returns>Own instance to allow Fluid API calls</returns>
        IPlcWriteRequestBuilder AddItem<T>(string name, string fieldQuery, params T[] values);
    }
}