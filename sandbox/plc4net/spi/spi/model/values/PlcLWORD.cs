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
    public class PlcLWORD : PlcBitString
    {
        private ulong value;
        
        public PlcLWORD(ulong value)
        {
            this.value = value;
        }
        
        public new int GetBoolLength()
        {
            return 64;
        }
        
        public new bool GetBoolAt(int index)
        {
            if (index > 63)
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
                GetBoolAt(6), GetBoolAt(7),
                GetBoolAt(8), GetBoolAt(9),
                GetBoolAt(10), GetBoolAt(11),
                GetBoolAt(12), GetBoolAt(13),
                GetBoolAt(14), GetBoolAt(15),
                GetBoolAt(16), GetBoolAt(17),
                GetBoolAt(18), GetBoolAt(19),
                GetBoolAt(20), GetBoolAt(21),
                GetBoolAt(22), GetBoolAt(23),
                GetBoolAt(24), GetBoolAt(25),
                GetBoolAt(26), GetBoolAt(27),
                GetBoolAt(28), GetBoolAt(29),
                GetBoolAt(30), GetBoolAt(31),
                GetBoolAt(32), GetBoolAt(33),
                GetBoolAt(34), GetBoolAt(35),
                GetBoolAt(36), GetBoolAt(37),
                GetBoolAt(38), GetBoolAt(39),
                GetBoolAt(40), GetBoolAt(41),
                GetBoolAt(42), GetBoolAt(43),
                GetBoolAt(44), GetBoolAt(45),
                GetBoolAt(46), GetBoolAt(47),
                GetBoolAt(48), GetBoolAt(49),
                GetBoolAt(50), GetBoolAt(51),
                GetBoolAt(52), GetBoolAt(53),
                GetBoolAt(54), GetBoolAt(55),
                GetBoolAt(56), GetBoolAt(57),
                GetBoolAt(58), GetBoolAt(59),
                GetBoolAt(60), GetBoolAt(61),
                GetBoolAt(62), GetBoolAt(63)
            };
        }

    }
}