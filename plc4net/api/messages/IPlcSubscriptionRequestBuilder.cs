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

using System;

namespace org.apache.plc4net.messages
{
    /// <summary>
    /// Interface for classes building requests for subscriptions
    /// </summary>
    public interface IPlcSubscriptionRequestBuilder: IPlcRequestBuilder<IPlcSubscriptionRequest>
    {
        /// <summary>
        /// Add a field that gets polled cyclically
        /// </summary>
        /// <param name="name">Alias for the field</param>
        /// <param name="fieldQuery">Query string for the field</param>
        /// <param name="pollingInterval">Time Interval at which the field should be polled</param>
        /// <returns>Builder instance for Fluid API requests</returns>
        IPlcSubscriptionRequestBuilder AddCyclicField(string name, string fieldQuery, TimeSpan pollingInterval);

        /// <summary>
        /// Add a field that sends an update when its value is changed
        /// </summary>
        /// <param name="name">Alias for the field</param>
        /// <param name="fieldQuery">Query string for the field</param>        
        /// <returns>Builder instance for Fluid API requests</returns>
        IPlcSubscriptionRequestBuilder AddChangeOfStateField(string name, string fieldQuery);

        /// <summary>
        /// Add a field that sends an update when an event occurs
        /// </summary>
        /// <param name="name">Alias for the field</param>
        /// <param name="fieldQuery">Query string for the field</param>        
        /// <returns>Builder instance for Fluid API requests</returns>
        IPlcSubscriptionRequestBuilder AddEventField(string name, string fieldQuery);
    }
}