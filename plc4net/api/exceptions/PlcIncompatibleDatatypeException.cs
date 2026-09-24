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

namespace org.apache.plc4net.exceptions
{
    /// <summary>
    /// Thrown when a <c>PlcValue</c> accessor is used for a type the value
    /// cannot represent, e.g. calling <c>GetString()</c> on a numeric value.
    /// Mirrors the Java SPI3 <c>PlcIncompatibleDatatypeException</c>: check the
    /// <c>IsXxx()</c> guards instead of catching this.
    /// </summary>
    public class PlcIncompatibleDatatypeException : PlcException
    {
        public PlcIncompatibleDatatypeException()
        {
        }

        public PlcIncompatibleDatatypeException(string message) : base(message)
        {
        }

        public PlcIncompatibleDatatypeException(string message, Exception inner) : base(message, inner)
        {
        }
    }
}