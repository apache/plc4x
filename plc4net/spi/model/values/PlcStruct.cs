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
    /// <summary>A named group of <see cref="IPlcValue"/> - the shape a KNX
    /// composite datapoint (<c>DPT_Switch_Control</c>, <c>DPT_DateTime</c>)
    /// parses to.</summary>
    public class PlcStruct : PlcValueAdapter
    {
        private readonly Dictionary<string, IPlcValue> values;

        public PlcStruct(Dictionary<string, IPlcValue> values)
        {
            this.values = values ?? new Dictionary<string, IPlcValue>();
        }

        public override bool IsStruct()
        {
            return true;
        }

        public override string[] GetKeys()
        {
            return values.Keys.ToArray();
        }

        public override bool HasKey(string key)
        {
            return values.ContainsKey(key);
        }

        public override IPlcValue GetValue(string key)
        {
            return values.TryGetValue(key, out var value) ? value : null;
        }

        public override Dictionary<string, IPlcValue> GetStruct()
        {
            return values;
        }

        protected bool Equals(PlcStruct other)
        {
            return values.Count == other.values.Count
                   && values.All(kv => other.values.TryGetValue(kv.Key, out var v)
                                       && Equals(kv.Value, v));
        }

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((PlcStruct) obj);
        }

        public override int GetHashCode()
        {
            var hash = 17;
            foreach (var kv in values.OrderBy(kv => kv.Key, System.StringComparer.Ordinal))
            {
                hash = hash * 31 + kv.Key.GetHashCode();
                hash = hash * 31 + (kv.Value != null ? kv.Value.GetHashCode() : 0);
            }

            return hash;
        }
    }

}
