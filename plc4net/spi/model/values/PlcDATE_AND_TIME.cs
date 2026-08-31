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

using System;

namespace org.apache.plc4net.spi.model.values
{
    /// <summary>IEC 61131 <c>DATE_AND_TIME</c> (S7 <c>DT</c>) - a full instant,
    /// millisecond resolution.</summary>
    public class PlcDATE_AND_TIME : PlcSimpleValueAdapter
    {
        private readonly DateTime value;

        public PlcDATE_AND_TIME(DateTime value)
        {
            this.value = value;
        }

        public static PlcDATE_AND_TIME OfSegments(
            int year, int month, int day, int hour, int minutes, int seconds, int nanoseconds) =>
            new PlcDATE_AND_TIME(new DateTime(
                year, month == 0 ? 1 : month, day == 0 ? 1 : day, hour, minutes, seconds)
                .AddTicks(nanoseconds / 100));

        public override bool IsDateTime()
        {
            return true;
        }

        public override DateTime GetDateTime()
        {
            return value;
        }

        protected bool Equals(PlcDATE_AND_TIME other)
        {
            return value.Equals(other.value);
        }

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((PlcDATE_AND_TIME) obj);
        }

        public override int GetHashCode()
        {
            return value.GetHashCode();
        }
    }
}
