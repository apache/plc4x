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
        IPlcValue[] GetList();

        // Struct Methods
        bool IsStruct();
        string[] GetKeys();
        bool HasKey(string key);
        IPlcValue GetValue(string key);
        Dictionary<string, IPlcValue> GetStruct();
    }
}