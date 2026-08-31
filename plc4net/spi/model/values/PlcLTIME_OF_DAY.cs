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
    /// <summary>
    /// IEC 61131 <c>LTIME_OF_DAY</c> - a wall-clock time in nanoseconds since
    /// midnight. The raw <see cref="GetNanosecondsSinceMidnight"/> is exact;
    /// <see cref="GetTime"/> rounds to the 100 ns <see cref="TimeOnly"/> tick.
    /// </summary>
    public class PlcLTIME_OF_DAY : PlcSimpleValueAdapter
    {
        private readonly ulong nanosecondsSinceMidnight;

        public PlcLTIME_OF_DAY(ulong nanosecondsSinceMidnight)
        {
            this.nanosecondsSinceMidnight = nanosecondsSinceMidnight;
        }

        public static PlcLTIME_OF_DAY OfNanosecondsSinceMidnight(ulong nanosecondsSinceMidnight) =>
            new PlcLTIME_OF_DAY(nanosecondsSinceMidnight);

        public ulong GetNanosecondsSinceMidnight()
        {
            return nanosecondsSinceMidnight;
        }

        public override bool IsTime()
        {
            return true;
        }

        public override TimeOnly GetTime()
        {
            return TimeOnly.FromTimeSpan(TimeSpan.FromTicks((long) (nanosecondsSinceMidnight / 100)));
        }

        protected bool Equals(PlcLTIME_OF_DAY other)
        {
            return nanosecondsSinceMidnight == other.nanosecondsSinceMidnight;
        }

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((PlcLTIME_OF_DAY) obj);
        }

        public override int GetHashCode()
        {
            return nanosecondsSinceMidnight.GetHashCode();
        }
    }
}
