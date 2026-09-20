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
    /// IEC 61131 <c>LTIME</c> - a duration in nanoseconds. The raw
    /// <see cref="GetNanoseconds"/> is exact; <see cref="GetDuration"/> rounds
    /// to the 100 ns <see cref="TimeSpan"/> tick.
    /// </summary>
    public class PlcLTIME : PlcSimpleValueAdapter
    {
        private readonly ulong nanoseconds;

        public PlcLTIME(ulong nanoseconds)
        {
            this.nanoseconds = nanoseconds;
        }

        public static PlcLTIME OfNanoseconds(ulong nanoseconds) => new PlcLTIME(nanoseconds);

        public ulong GetNanoseconds()
        {
            return nanoseconds;
        }

        public override bool IsDuration()
        {
            return true;
        }

        public override TimeSpan GetDuration()
        {
            return TimeSpan.FromTicks((long) (nanoseconds / 100));
        }

        protected bool Equals(PlcLTIME other)
        {
            return nanoseconds == other.nanoseconds;
        }

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((PlcLTIME) obj);
        }

        public override int GetHashCode()
        {
            return nanoseconds.GetHashCode();
        }
    }
}
