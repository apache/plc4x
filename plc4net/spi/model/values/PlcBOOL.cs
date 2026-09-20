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

using org.apache.plc4net.api.value;

namespace org.apache.plc4net.spi.model.values
{
    public class PlcBOOL : PlcSimpleValueAdapter
    {
        private bool value;
        
        public PlcBOOL(bool value)
        {
            this.value = value;
        }

        public override bool IsBool()
        {
            return true;
        }

        public override bool GetBool()
        {
            return value;
        }

        public override bool[] GetBoolArray()
        {
            return new[] { value };
        }

        // A BOOL is also 1/0 as a number and "true"/"false" as text (matching plc4j).

        public override bool IsByte() => true;
        public override byte GetByte() => (byte) (value ? 1 : 0);

        public override bool IsUshort() => true;
        public override ushort GetUshort() => (ushort) (value ? 1 : 0);

        public override bool IsInt() => true;
        public override int GetInt() => value ? 1 : 0;

        public override bool IsString()
        {
            return true;
        }

        public override string GetString()
        {
            return value ? "true" : "false";
        }

        protected bool Equals(PlcBOOL other)
        {
            return value == other.value;
        }

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((PlcBOOL) obj);
        }

        public override int GetHashCode()
        {
            return value.GetHashCode();
        }
        
    }
}