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
    /// IEC 61131 <c>DATE_AND_LTIME</c> (S7 <c>DTL</c>) - a full instant with
    /// nanosecond-of-second precision. The whole-second part is a
    /// <see cref="DateTime"/>; the sub-second part is kept exactly as
    /// <see cref="GetNanosecondsOfSecond"/> so both the segmented (DTL) and the
    /// nanoseconds-since-epoch wire forms round-trip.
    /// </summary>
    public class PlcDATE_AND_LTIME : PlcSimpleValueAdapter
    {
        private const long NanosPerSecond = 1_000_000_000L;

        private readonly DateTime wholeSeconds;
        private readonly uint nanosecondsOfSecond;

        public PlcDATE_AND_LTIME(DateTime dateTime, uint nanosecondsOfSecond)
        {
            // Split the instant into whole seconds + a nanoseconds-of-second so
            // both wire forms round-trip; any sub-second part carried by the
            // DateTime is folded in rather than dropped. Kind is normalised to
            // Unspecified to match plc4j's zoneless LocalDateTime semantics -
            // otherwise the two factories below disagree (UnixEpoch is Utc,
            // `new DateTime(y,m,d,…)` is Unspecified).
            var subSecondTicks = dateTime.Ticks % TimeSpan.TicksPerSecond;
            wholeSeconds = DateTime.SpecifyKind(
                new DateTime(dateTime.Ticks - subSecondTicks), DateTimeKind.Unspecified);
            this.nanosecondsOfSecond = nanosecondsOfSecond + (uint) (subSecondTicks * 100);
        }

        public static PlcDATE_AND_LTIME OfNanosecondsSinceEpoch(ulong nanosecondsSinceEpoch)
        {
            var seconds = (long) (nanosecondsSinceEpoch / NanosPerSecond);
            var nanos = (uint) (nanosecondsSinceEpoch % NanosPerSecond);
            return new PlcDATE_AND_LTIME(DateTime.UnixEpoch.AddSeconds(seconds), nanos);
        }

        public static PlcDATE_AND_LTIME OfSegments(
            int year, int month, int day, int hour, int minutes, int seconds, long nanoseconds) =>
            new PlcDATE_AND_LTIME(
                new DateTime(year, month == 0 ? 1 : month, day == 0 ? 1 : day, hour, minutes, seconds),
                (uint) nanoseconds);

        public uint GetNanosecondsOfSecond()
        {
            return nanosecondsOfSecond;
        }

        public ulong GetNanosecondsSinceEpoch()
        {
            var seconds = (wholeSeconds.Ticks - DateTime.UnixEpoch.Ticks) / TimeSpan.TicksPerSecond;
            return (ulong) seconds * (ulong) NanosPerSecond + nanosecondsOfSecond;
        }

        public override bool IsDateTime()
        {
            return true;
        }

        public override DateTime GetDateTime()
        {
            return wholeSeconds.AddTicks(nanosecondsOfSecond / 100);
        }

        protected bool Equals(PlcDATE_AND_LTIME other)
        {
            return wholeSeconds.Equals(other.wholeSeconds)
                && nanosecondsOfSecond == other.nanosecondsOfSecond;
        }

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((PlcDATE_AND_LTIME) obj);
        }

        public override int GetHashCode()
        {
            return HashCode.Combine(wholeSeconds, nanosecondsOfSecond);
        }
    }
}
