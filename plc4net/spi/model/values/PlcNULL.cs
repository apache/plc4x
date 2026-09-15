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
    /// <summary>
    /// The absence of a value — the fall-through result of a <c>[dataIo]</c>
    /// parse and of an unresolved lookup. All instances are equal.
    /// </summary>
    public class PlcNULL : PlcValueAdapter
    {
        public PlcNULL()
        {
        }

        public override bool IsNull()
        {
            return true;
        }

        public override bool IsNullable()
        {
            return true;
        }

        public override string GetString()
        {
            return "null";
        }

        public override bool Equals(object obj)
        {
            return obj is PlcNULL;
        }

        public override int GetHashCode()
        {
            return 0;
        }

        public override string ToString()
        {
            return "PlcNULL";
        }
    }
}