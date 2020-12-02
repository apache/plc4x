using System;

namespace org.apache.plc4net.spi.model.values
{
    public abstract class SimpleNumericValueAdapter<T> : PlcSimpleValueAdapter where T : IComparable
    {
        private IComparable value;

        public SimpleNumericValueAdapter(IComparable value)
        {
            this.value = value;
        }
        
        public new bool IsBool()
        {
            return true;
        }
        
        public new bool GetBool()
        {
            return value.CompareTo(0) != 0;
        }
        
        public new bool IsByte()
        {
            return (value.CompareTo(byte.MinValue) >= 0) && (value.CompareTo(byte.MaxValue) <= 0);
        }

        public new byte GetByte()
        {
            if (!IsByte())
            {
                throw new ArgumentOutOfRangeException();
            }
            return (byte) value;
        }

        public new bool IsUshort()
        {
            return (value.CompareTo(ushort.MinValue) >= 0) && (value.CompareTo(ushort.MaxValue) <= 0);
        }

        public new ushort GetUshort()
        {
            if (!IsUshort())
            {
                throw new ArgumentOutOfRangeException();
            }
            return (ushort) value;
        }

        public new bool IsUint()
        {
            return (value.CompareTo(uint.MinValue) >= 0) && (value.CompareTo(uint.MaxValue) <= 0);
        }

        public new uint GetUint()
        {
            if (!IsUint())
            {
                throw new ArgumentOutOfRangeException();
            }
            return (uint) value;
        }

        public new bool IsUlong()
        {
            return (value.CompareTo(ulong.MinValue) >= 0) && (value.CompareTo(ulong.MaxValue) <= 0);
        }

        public new ulong GetUlong()
        {
            if (!IsUlong())
            {
                throw new ArgumentOutOfRangeException();
            }
            return (ulong) value;
        }

        public new bool IsSbyte()
        {
            return (value.CompareTo(sbyte.MinValue) >= 0) && (value.CompareTo(sbyte.MaxValue) <= 0);
        }

        public new sbyte GetSbyte()
        {
            if (!IsSbyte())
            {
                throw new ArgumentOutOfRangeException();
            }
            return (sbyte) value;
        }

        public new bool IsShort()
        {
            return (value.CompareTo(short.MinValue) >= 0) && (value.CompareTo(short.MaxValue) <= 0);
        }

        public new short GetShort()
        {
            if (!IsShort())
            {
                throw new ArgumentOutOfRangeException();
            }
            return (short) value;
        }

        public new bool IsInt()
        {
            return (value.CompareTo(int.MinValue) >= 0) && (value.CompareTo(int.MaxValue) <= 0);
        }

        public new int GetInt()
        {
            if (!IsInt())
            {
                throw new ArgumentOutOfRangeException();
            }
            return (int) value;
        }

        public new bool IsLong()
        {
            return (value.CompareTo(long.MinValue) >= 0) && (value.CompareTo(long.MaxValue) <= 0);
        }

        public new long GetLong()
        {
            if (!IsLong())
            {
                throw new ArgumentOutOfRangeException();
            }
            return (long) value;
        }

        public new bool IsFloat()
        {
            return (value.CompareTo(-float.MaxValue) >= 0) && (value.CompareTo(float.MaxValue) <= 0);
        }

        public new float GetFloat()
        {
            if (!IsFloat())
            {
                throw new ArgumentOutOfRangeException();
            }
            return (float) value;
        }

        public new bool IsDouble()
        {
            return (value.CompareTo(-double.MaxValue) >= 0) && (value.CompareTo(double.MaxValue) <= 0);
        }

        public new double GetDouble()
        {
            if (!IsDouble())
            {
                throw new ArgumentOutOfRangeException();
            }
            return (double) value;
        }

        public new bool IsString()
        {
            return true;
        }

        public new String GetString()
        {
            return value.ToString();
        }

    }
}