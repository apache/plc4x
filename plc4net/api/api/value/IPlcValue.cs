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

namespace org.apache.plc4net.api.value
{
    public interface IPlcValue
    {
        // Simple Types
        bool IsSimple();
        bool IsNullable();
        bool IsNull();

        // Boolean
        bool IsBool();
        int GetBoolLength();
        bool GetBool();
        bool GetBoolAt(int index);
        bool[] GetBoolArray();

        // Integer
        bool IsByte();
        byte GetByte();
        bool IsUshort();
        ushort GetUshort();
        bool IsUint();
        uint GetUint();
        bool IsUlong();
        ulong GetUlong();
        bool IsSbyte();
        sbyte GetSbyte();
        bool IsShort();
        short GetShort();
        bool IsInt();
        int GetInt();
        bool IsLong();
        long GetLong();

        // Floating Point
        bool IsFloat();
        float GetFloat();
        bool IsDouble();
        double GetDouble();

        // String
        bool IsString();
        string GetString();

        // Time
        bool IsDateTime();
        DateTime GetDateTime();

        // Raw Access
        byte[] GetRaw();

        // List Methods
        bool IsList();
        int GetLength();
        IPlcValue GetIndex(int index);
        List<IPlcValue> GetList();

        // Struct Methods
        bool IsStruct();
        string[] GetKeys();
        bool HasKey(string key);
        IPlcValue GetValue(string key);
        Dictionary<string, IPlcValue> GetStruct();
    }
}