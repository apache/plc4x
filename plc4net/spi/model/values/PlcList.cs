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

using System.Collections.Generic;
using System.Linq;
using org.apache.plc4net.api.value;

namespace org.apache.plc4net.spi.model.values
{
    /// <summary>An ordered list of <see cref="IPlcValue"/> - a KNX group of
    /// datapoints, a Modbus multi-register read.</summary>
    public class PlcList : PlcValueAdapter
    {
        private readonly List<IPlcValue> values;

        public PlcList(List<IPlcValue> values)
        {
            this.values = values ?? new List<IPlcValue>();
        }

        public override bool IsList()
        {
            return true;
        }

        public override int GetLength()
        {
            return values.Count;
        }

        public override IPlcValue GetIndex(int index)
        {
            return values[index];
        }

        public override List<IPlcValue> GetList()
        {
            return values;
        }

        protected bool Equals(PlcList other)
        {
            return values.SequenceEqual(other.values);
        }

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((PlcList) obj);
        }

        public override int GetHashCode()
        {
            var hash = 17;
            foreach (var value in values)
            {
                hash = hash * 31 + (value != null ? value.GetHashCode() : 0);
            }

            return hash;
        }
    }
}
