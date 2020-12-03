//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
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
    public class PlcBYTE : PlcBitString
    {
        private byte value;

        public PlcBYTE(byte value)
        {
            this.value = value;
        }

        public new int GetBoolLength()
        {
            return 8;
        }

        public new bool GetBoolAt(int index)
        {
            if (index > 7)
            {
                return false;
            }

            return ((value >> index) & 1) == 1;
        }

        public new bool[] GetBoolArray()
        {
            return new[]
            {
                GetBoolAt(0), GetBoolAt(1),
                GetBoolAt(2), GetBoolAt(3),
                GetBoolAt(4), GetBoolAt(5),
                GetBoolAt(6), GetBoolAt(7)
            };
        }

        protected bool Equals(PlcBYTE other)
        {
            return value == other.value;
        }

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((PlcBYTE) obj);
        }

        public override int GetHashCode()
        {
            return value.GetHashCode();
        }
        
    }
}