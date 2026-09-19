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
using org.apache.plc4net.spi.generation;

namespace org.apache.plc4net.spi.model.values
{
    /// <summary>IEC 61131 <c>TIME_OF_DAY</c> - a wall-clock time, no date.</summary>
    public class PlcTIME_OF_DAY : PlcSimpleValueAdapter
    {
        private const long MillisPerDay = 86_400_000L;

        private readonly TimeOnly value;

        public PlcTIME_OF_DAY(TimeOnly value)
        {
            this.value = value;
        }

        public static PlcTIME_OF_DAY OfMillisecondsSinceMidnight(long millisecondsSinceMidnight)
        {
            if (millisecondsSinceMidnight < 0 || millisecondsSinceMidnight >= MillisPerDay)
            {
                // A corrupt frame must surface as a ParseException the codec can
                // catch, not a framework ArgumentOutOfRangeException that escapes
                // the receive loop.
                throw new ParseException(
                    $"TIME_OF_DAY {millisecondsSinceMidnight} ms is not within a 24-hour day.");
            }
            return new PlcTIME_OF_DAY(TimeOnly.FromTimeSpan(TimeSpan.FromMilliseconds(millisecondsSinceMidnight)));
        }

        public override bool IsTime()
        {
            return true;
        }

        public override TimeOnly GetTime()
        {
            return value;
        }

        protected bool Equals(PlcTIME_OF_DAY other)
        {
            return value.Equals(other.value);
        }

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((PlcTIME_OF_DAY) obj);
        }

        public override int GetHashCode()
        {
            return value.GetHashCode();
        }
    }
}
