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

using System;
using System.Collections.Generic;
using org.apache.plc4net.api.value;

namespace org.apache.plc4net.spi.model.values
{
    public abstract class PlcValueAdapter : IPlcValue
    {
        public bool IsSimple()
        {
            return false;
        }

        public bool IsNullable()
        {
            return false;
        }

        public bool IsNull()
        {
            return false;
        }

        public bool IsBool()
        {
            return false;
        }

        public int GetBoolLength()
        {
            return 1;
        }

        public bool GetBool()
        {
            return default;
        }

        public bool GetBoolAt(int index)
        {
            if (index == 0)
            {
                return GetBool();
            }
            return default;
        }

        public bool[] GetBoolArray()
        {
            return default;
        }

        public bool IsByte()
        {
            return false;
        }

        public byte GetByte()
        {
            return default;
        }

        public bool IsUshort()
        {
            return false;
        }

        public ushort GetUshort()
        {
            return default;
        }

        public bool IsUint()
        {
            return false;
        }

        public uint GetUint()
        {
            return default;
        }

        public bool IsUlong()
        {
            return false;
        }

        public ulong GetUlong()
        {
            return default;
        }

        public bool IsSbyte()
        {
            return false;
        }

        public sbyte GetSbyte()
        {
            return default;
        }

        public bool IsShort()
        {
            return false;
        }

        public short GetShort()
        {
            return default;
        }

        public bool IsInt()
        {
            return false;
        }

        public int GetInt()
        {
            return default;
        }

        public bool IsLong()
        {
            return false;
        }

        public long GetLong()
        {
            return default;
        }

        public bool IsFloat()
        {
            return false;
        }

        public float GetFloat()
        {
            return default;
        }

        public bool IsDouble()
        {
            return false;
        }

        public double GetDouble()
        {
            return default;
        }

        public bool IsString()
        {
            return false;
        }

        public string GetString()
        {
            return default;
        }

        public bool IsDateTime()
        {
            return false;
        }

        public DateTime GetDateTime()
        {
            return default;
        }

        public byte[] GetRaw()
        {
            return default;
        }

        public bool IsList()
        {
            return false;
        }

        public int GetLength()
        {
            return 1;
        }

        public IPlcValue GetIndex(int index)
        {
            return default;
        }

        public List<IPlcValue> GetList()
        {
            return default;
        }

        public bool IsStruct()
        {
            return false;
        }

        public string[] GetKeys()
        {
            return default;
        }

        public bool HasKey(string key)
        {
            return false;
        }

        public IPlcValue GetValue(string key)
        {
            return default;
        }

        public Dictionary<string, IPlcValue> GetStruct()
        {
            return default;
        }
    }
}