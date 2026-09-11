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
using System.Globalization;

namespace org.apache.plc4net.spi.model.values
{
    public abstract class SimpleNumericValueAdapter<T> : PlcSimpleValueAdapter where T : IComparable
    {
        private readonly IComparable value;

        protected SimpleNumericValueAdapter(IComparable value)
        {
            this.value = value;
        }

        /// <summary>
        /// The boxed value, widened to decimal for range checks.
        /// </summary>
        /// <remarks>
        /// Range checks cannot go through <see cref="IComparable.CompareTo"/>: comparing a
        /// boxed Int32 against a UInt64 bound throws ArgumentException, because CompareTo
        /// requires both operands to be the same runtime type. Widening to a common type
        /// first is what makes the comparison well defined.
        ///
        /// decimal covers every integral PLC type exactly, including ulong.MaxValue. It
        /// cannot hold the outer range of float/double, so an out-of-range REAL/LREAL
        /// reports false rather than throwing.
        /// </remarks>
        private bool TryAsDecimal(out decimal result)
        {
            try
            {
                result = Convert.ToDecimal(value, CultureInfo.InvariantCulture);
                return true;
            }
            catch (OverflowException)
            {
                result = default;
                return false;
            }
            catch (InvalidCastException)
            {
                result = default;
                return false;
            }
        }

        private bool InRange(decimal min, decimal max)
        {
            return TryAsDecimal(out var d) && d >= min && d <= max;
        }

        public override bool IsBool()
        {
            return true;
        }

        public override bool GetBool()
        {
            return TryAsDecimal(out var d) ? d != 0m : false;
        }

        public override bool IsByte()
        {
            return InRange(byte.MinValue, byte.MaxValue);
        }

        public override byte GetByte()
        {
            if (!IsByte())
            {
                throw new ArgumentOutOfRangeException(nameof(value));
            }
            return Convert.ToByte(value, CultureInfo.InvariantCulture);
        }

        public override bool IsUshort()
        {
            return InRange(ushort.MinValue, ushort.MaxValue);
        }

        public override ushort GetUshort()
        {
            if (!IsUshort())
            {
                throw new ArgumentOutOfRangeException(nameof(value));
            }
            return Convert.ToUInt16(value, CultureInfo.InvariantCulture);
        }

        public override bool IsUint()
        {
            return InRange(uint.MinValue, uint.MaxValue);
        }

        public override uint GetUint()
        {
            if (!IsUint())
            {
                throw new ArgumentOutOfRangeException(nameof(value));
            }
            return Convert.ToUInt32(value, CultureInfo.InvariantCulture);
        }

        public override bool IsUlong()
        {
            return InRange(ulong.MinValue, ulong.MaxValue);
        }

        public override ulong GetUlong()
        {
            if (!IsUlong())
            {
                throw new ArgumentOutOfRangeException(nameof(value));
            }
            return Convert.ToUInt64(value, CultureInfo.InvariantCulture);
        }

        public override bool IsSbyte()
        {
            return InRange(sbyte.MinValue, sbyte.MaxValue);
        }

        public override sbyte GetSbyte()
        {
            if (!IsSbyte())
            {
                throw new ArgumentOutOfRangeException(nameof(value));
            }
            return Convert.ToSByte(value, CultureInfo.InvariantCulture);
        }

        public override bool IsShort()
        {
            return InRange(short.MinValue, short.MaxValue);
        }

        public override short GetShort()
        {
            if (!IsShort())
            {
                throw new ArgumentOutOfRangeException(nameof(value));
            }
            return Convert.ToInt16(value, CultureInfo.InvariantCulture);
        }

        public override bool IsInt()
        {
            return InRange(int.MinValue, int.MaxValue);
        }

        public override int GetInt()
        {
            if (!IsInt())
            {
                throw new ArgumentOutOfRangeException(nameof(value));
            }
            return Convert.ToInt32(value, CultureInfo.InvariantCulture);
        }

        public override bool IsLong()
        {
            return InRange(long.MinValue, long.MaxValue);
        }

        public override long GetLong()
        {
            if (!IsLong())
            {
                throw new ArgumentOutOfRangeException(nameof(value));
            }
            return Convert.ToInt64(value, CultureInfo.InvariantCulture);
        }

        public override bool IsFloat()
        {
            var d = Convert.ToDouble(value, CultureInfo.InvariantCulture);
            // NaN / ±Infinity are representable as a float; only a finite value
            // outside the float range is not.
            return double.IsNaN(d) || double.IsInfinity(d)
                   || (d >= -float.MaxValue && d <= float.MaxValue);
        }

        public override float GetFloat()
        {
            if (!IsFloat())
            {
                throw new ArgumentOutOfRangeException(nameof(value));
            }
            return Convert.ToSingle(value, CultureInfo.InvariantCulture);
        }

        public override bool IsDouble()
        {
            return true;
        }

        public override double GetDouble()
        {
            return Convert.ToDouble(value, CultureInfo.InvariantCulture);
        }

        public override bool IsString()
        {
            return true;
        }

        public override string GetString()
        {
            return Convert.ToString(value, CultureInfo.InvariantCulture);
        }

        protected bool Equals(SimpleNumericValueAdapter<T> other)
        {
            return Equals(value, other.value);
        }

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((SimpleNumericValueAdapter<T>) obj);
        }

        public override int GetHashCode()
        {
            return (value != null ? value.GetHashCode() : 0);
        }

    }
}
