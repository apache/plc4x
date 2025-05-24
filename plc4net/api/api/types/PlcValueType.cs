/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

namespace org.apache.plc4net.types
{
    using System;
    using System.Collections.Generic;
    
    /// <summary>
    /// The supported types of PlcValues.
    /// </summary>
    public class PlcValueType
{
    public static readonly PlcValueType NULL = new PlcValueType(0x00, null);
    public static readonly PlcValueType BOOL = new PlcValueType(0x01, typeof(bool));
    public static readonly PlcValueType BYTE = new PlcValueType(0x02, typeof(byte));
    public static readonly PlcValueType WORD = new PlcValueType(0x03, typeof(short));
    public static readonly PlcValueType DWORD = new PlcValueType(0x04, typeof(int));
    public static readonly PlcValueType LWORD = new PlcValueType(0x05, typeof(long));
    public static readonly PlcValueType USINT = new PlcValueType(0x11, typeof(ushort));
    public static readonly PlcValueType UINT = new PlcValueType(0x12, typeof(uint));
    public static readonly PlcValueType UDINT = new PlcValueType(0x13, typeof(ulong));
    public static readonly PlcValueType ULINT = new PlcValueType(0x14, typeof(System.Numerics.BigInteger));
    public static readonly PlcValueType SINT = new PlcValueType(0x21, typeof(sbyte));
    public static readonly PlcValueType INT = new PlcValueType(0x22, typeof(short));
    public static readonly PlcValueType DINT = new PlcValueType(0x23, typeof(int));
    public static readonly PlcValueType LINT = new PlcValueType(0x24, typeof(long));
    public static readonly PlcValueType REAL = new PlcValueType(0x31, typeof(float));
    public static readonly PlcValueType LREAL = new PlcValueType(0x32, typeof(double));
    public static readonly PlcValueType CHAR = new PlcValueType(0x41, typeof(char));
    public static readonly PlcValueType WCHAR = new PlcValueType(0x42, typeof(short));
    public static readonly PlcValueType STRING = new PlcValueType(0x43, typeof(string));
    public static readonly PlcValueType WSTRING = new PlcValueType(0x44, typeof(string));
    public static readonly PlcValueType TIME = new PlcValueType(0x51, typeof(TimeSpan));
    public static readonly PlcValueType LTIME = new PlcValueType(0x52, typeof(TimeSpan));
    public static readonly PlcValueType DATE = new PlcValueType(0x53, typeof(DateTime));
    public static readonly PlcValueType LDATE = new PlcValueType(0x54, typeof(DateTime));
    public static readonly PlcValueType TIME_OF_DAY = new PlcValueType(0x55, typeof(TimeSpan));
    public static readonly PlcValueType LTIME_OF_DAY = new PlcValueType(0x56, typeof(TimeSpan));
    public static readonly PlcValueType DATE_AND_TIME = new PlcValueType(0x57, typeof(DateTime));
    public static readonly PlcValueType DATE_AND_LTIME = new PlcValueType(0x58, typeof(DateTime));
    public static readonly PlcValueType LDATE_AND_TIME = new PlcValueType(0x59, typeof(DateTime));
    public static readonly PlcValueType Struct = new PlcValueType(0x61, typeof(Dictionary<string, object>));
    public static readonly PlcValueType List = new PlcValueType(0x62, typeof(List<object>));
    public static readonly PlcValueType RAW_BYTE_ARRAY = new PlcValueType(0x71, typeof(byte[]));

    private static readonly Dictionary<short, PlcValueType> map = new Dictionary<short, PlcValueType>();

    static PlcValueType()
    {
        foreach (var field in typeof(PlcValueType).GetFields(System.Reflection.BindingFlags.Static | System.Reflection.BindingFlags.Public))
        {
            if (field.GetValue(null) is PlcValueType value)
            {
                map[value.Value] = value;
            }
        }
    }

    public short Value { get; }
    public Type DefaultDotNetType { get; }

    private PlcValueType(short value, Type defaultDotNetType)
    {
        Value = value;
        DefaultDotNetType = defaultDotNetType;
    }

    public static PlcValueType EnumForValue(short value)
    {
        map.TryGetValue(value, out var result);
        return result;
    }

    public static bool IsDefined(short value)
    {
        return map.ContainsKey(value);
    }

    public override string ToString() => $"PlcValueType(0x{Value:X}, {DefaultDotNetType?.Name ?? "null"})";
}
}