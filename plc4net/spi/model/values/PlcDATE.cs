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
    /// <summary>IEC 61131 <c>DATE</c> - a calendar date, no time of day.</summary>
    public class PlcDATE : PlcSimpleValueAdapter
    {
        private static readonly DateOnly Epoch = new DateOnly(1970, 1, 1);

        private readonly DateOnly value;

        public PlcDATE(DateOnly value)
        {
            this.value = value;
        }

        public static PlcDATE OfDaysSinceEpoch(int daysSinceEpoch) =>
            new PlcDATE(Epoch.AddDays(daysSinceEpoch));

        public int GetDaysSinceEpoch()
        {
            return value.DayNumber - Epoch.DayNumber;
        }

        public override bool IsDate()
        {
            return true;
        }

        public override DateOnly GetDate()
        {
            return value;
        }

        public override bool IsDateTime()
        {
            return true;
        }

        public override DateTime GetDateTime()
        {
            return value.ToDateTime(TimeOnly.MinValue);
        }

        protected bool Equals(PlcDATE other)
        {
            return value.Equals(other.value);
        }

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((PlcDATE) obj);
        }

        public override int GetHashCode()
        {
            return value.GetHashCode();
        }
    }
}
