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
using System.Collections.Generic;
using org.apache.plc4net.api.value;

namespace org.apache.plc4net.spi.model.values
{
    public abstract class PlcValueAdapter : IPlcValue
    {
        public virtual bool IsSimple()
        {
            return false;
        }

        public virtual bool IsNullable()
        {
            return false;
        }

        public virtual bool IsNull()
        {
            return false;
        }

        public virtual bool IsBool()
        {
            return false;
        }

        public virtual int GetBoolLength()
        {
            return 1;
        }

        public virtual bool GetBool()
        {
            return default;
        }

        public virtual bool GetBoolAt(int index)
        {
            if (index == 0)
            {
                return GetBool();
            }
            return default;
        }

        public virtual bool[] GetBoolArray()
        {
            return default;
        }

        public virtual bool IsByte()
        {
            return false;
        }

        public virtual byte GetByte()
        {
            return default;
        }

        public virtual bool IsUshort()
        {
            return false;
        }

        public virtual ushort GetUshort()
        {
            return default;
        }

        public virtual bool IsUint()
        {
            return false;
        }

        public virtual uint GetUint()
        {
            return default;
        }

        public virtual bool IsUlong()
        {
            return false;
        }

        public virtual ulong GetUlong()
        {
            return default;
        }

        public virtual bool IsSbyte()
        {
            return false;
        }

        public virtual sbyte GetSbyte()
        {
            return default;
        }

        public virtual bool IsShort()
        {
            return false;
        }

        public virtual short GetShort()
        {
            return default;
        }

        public virtual bool IsInt()
        {
            return false;
        }

        public virtual int GetInt()
        {
            return default;
        }

        public virtual bool IsLong()
        {
            return false;
        }

        public virtual long GetLong()
        {
            return default;
        }

        public virtual bool IsFloat()
        {
            return false;
        }

        public virtual float GetFloat()
        {
            return default;
        }

        public virtual bool IsDouble()
        {
            return false;
        }

        public virtual double GetDouble()
        {
            return default;
        }

        public virtual bool IsString()
        {
            return false;
        }

        public virtual string GetString()
        {
            return default;
        }

        public virtual bool IsDuration()
        {
            return false;
        }

        public virtual TimeSpan GetDuration()
        {
            return default;
        }

        public virtual bool IsDate()
        {
            return false;
        }

        public virtual DateOnly GetDate()
        {
            return default;
        }

        public virtual bool IsTime()
        {
            return false;
        }

        public virtual TimeOnly GetTime()
        {
            return default;
        }

        public virtual bool IsDateTime()
        {
            return false;
        }

        public virtual DateTime GetDateTime()
        {
            return default;
        }

        public virtual byte[] GetRaw()
        {
            return default;
        }

        public virtual bool IsList()
        {
            return false;
        }

        public virtual int GetLength()
        {
            return 1;
        }

        public virtual IPlcValue GetIndex(int index)
        {
            return default;
        }

        public virtual List<IPlcValue> GetList()
        {
            return default;
        }

        public virtual bool IsStruct()
        {
            return false;
        }

        public virtual string[] GetKeys()
        {
            return default;
        }

        public virtual bool HasKey(string key)
        {
            return false;
        }

        public virtual IPlcValue GetValue(string key)
        {
            return default;
        }

        public virtual Dictionary<string, IPlcValue> GetStruct()
        {
            return default;
        }
    }
}