// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
// https://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied. See the License for the
// specific language governing permissions and limitations
// under the License.

using System;
using System.Linq;
using org.apache.plc4net.api.value;
using org.apache.plc4net.spi.model.values;

namespace org.apache.plc4net.spi.drivers.subscriptions
{
    internal static class ValueSnapshot
    {
        // Copy compound values before dispatch: callers may mutate lists and raw buffers.
        public static object? Create(IPlcValue? value)
        {
            if (value == null || value.IsNull()) return null;
            object contents;
            if (value.IsList()) contents = value.GetList().Select(Create).ToArray();
            else if (value.IsStruct()) contents = value.GetKeys().OrderBy(k => k, StringComparer.Ordinal)
                .Select(k => new object?[] { k, Create(value.GetValue(k)) }).ToArray();
            else if (value is SimpleNumericValueAdapter<float>) contents = value.GetFloat();
            else if (value is SimpleNumericValueAdapter<double>) contents = value.GetDouble();
            else if (value is SimpleNumericValueAdapter<ulong>) contents = value.GetUlong();
            else if (value is SimpleNumericValueAdapter<byte> || value is SimpleNumericValueAdapter<sbyte> ||
                value is SimpleNumericValueAdapter<short> || value is SimpleNumericValueAdapter<ushort> ||
                value is SimpleNumericValueAdapter<int> || value is SimpleNumericValueAdapter<uint> ||
                value is SimpleNumericValueAdapter<long>) contents = value.GetLong();
            else if (value is PlcBitString) contents = value.GetUlong();
            else if (value.IsBool()) contents = value.GetBoolArray().ToArray();
            else if (value.IsByte()) contents = value.GetByte();
            else if (value.IsUshort()) contents = value.GetUshort();
            else if (value.IsUint()) contents = value.GetUint();
            else if (value.IsUlong()) contents = value.GetUlong();
            else if (value.IsSbyte()) contents = value.GetSbyte();
            else if (value.IsShort()) contents = value.GetShort();
            else if (value.IsInt()) contents = value.GetInt();
            else if (value.IsLong()) contents = value.GetLong();
            else if (value.IsFloat()) contents = value.GetFloat();
            else if (value.IsDouble()) contents = value.GetDouble();
            else if (value.IsString()) contents = value.GetString();
            else if (value.IsDuration()) contents = value.GetDuration();
            else if (value.IsDate()) contents = value.GetDate();
            else if (value.IsTime()) contents = value.GetTime();
            else if (value.IsDateTime()) contents = value.GetDateTime();
            else contents = value.GetRaw().ToArray();
            return new object[] { value.GetType(), contents };
        }
    }
}
