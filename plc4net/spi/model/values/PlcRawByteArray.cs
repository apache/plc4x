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
using System.Linq;

namespace org.apache.plc4net.spi.model.values
{
    /// <summary>An opaque run of bytes - a KNX property's raw group address, a
    /// datapoint field with no finer structure.</summary>
    public class PlcRawByteArray : PlcValueAdapter
    {
        private readonly byte[] value;

        public PlcRawByteArray(byte[] value)
        {
            this.value = value ?? Array.Empty<byte>();
        }

        public override byte[] GetRaw()
        {
            return value;
        }

        public override int GetLength()
        {
            return value.Length;
        }

        protected bool Equals(PlcRawByteArray other)
        {
            return value.SequenceEqual(other.value);
        }

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((PlcRawByteArray) obj);
        }

        public override int GetHashCode()
        {
            var hash = 17;
            foreach (var b in value)
            {
                hash = hash * 31 + b;
            }

            return hash;
        }
    }
}
