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
using System.Text;

namespace org.apache.plc4net.spi.generation
{

    public class ReadBuffer
    {
        private readonly byte[] _data;
        private readonly BitReader _reader;
        private readonly int _startBit;
        private readonly int _endBit;
        private readonly Stack<ByteOrder> _byteOrder = new Stack<ByteOrder>(new[] { ByteOrder.BIG_ENDIAN });

        public ReadBuffer(byte[] data) : this(data, 0, data?.Length ?? throw new ArgumentNullException(nameof(data)))
        {
        }

        /// <summary>
        /// Reads a byte-aligned window of <paramref name="data"/> without copying it.
        /// Positions reported by this buffer are relative to the window, and no read can
        /// cross its end. This lets stream codecs validate a frame inside a larger receive
        /// buffer without allocating a temporary array for that frame.
        /// </summary>
        public ReadBuffer(byte[] data, int offset, int length)
        {
            if (data == null)
            {
                throw new ArgumentNullException(nameof(data));
            }
            if (offset < 0 || length < 0 || offset > data.Length - length)
            {
                throw new ArgumentOutOfRangeException(
                    nameof(length), "The buffer window must be wholly contained in the data array.");
            }

            _data = data;
            _reader = new BitReader(data);
            _startBit = checked(offset * 8);
            _endBit = checked((offset + length) * 8);
            _reader.Position = _startBit;
        }

        public void Reset()
        {
            _reader.Position = _startBit;
        }

        /// <summary>
        /// Scopes <see cref="ByteOrder.LITTLE_ENDIAN"/> (or back to
        /// <see cref="ByteOrder.BIG_ENDIAN"/>) to the type whose generated
        /// <c>StaticParse</c> declares it. Must be paired with
        /// <see cref="PopByteOrder"/>, normally via <c>try</c> / <c>finally</c>
        /// so an early return (a type-switch case) still pops it.
        /// </summary>
        public void PushByteOrder(ByteOrder order) => _byteOrder.Push(order);

        public void PopByteOrder() => _byteOrder.Pop();

        public int GetPos()
        {
            return _reader.Position - _startBit;
        }

        public byte[] GetBytes()
        {
            return _data;
        }

        public int GetTotalBytes()
        {
            return (_endBit - _startBit) / 8;
        }

        public bool HasMore(int bitLength)
        {
            // A negative bit length is nonsensical and should not silently return true
            // just because the comparison below would hold for any non-negative remaining.
            if (bitLength < 0)
            {
                return false;
            }
            // Inclusive: a buffer with exactly bitLength bits left can still serve
            // a read of that size.
            return bitLength <= _endBit - _reader.Position;
        }

        public byte PeekByte(int offset)
        {
            EnsureReadable((offset + 1) * 8);
            return _data[(_reader.Position / 8) + offset];
        }

        public bool ReadBit(String logicalName)
        {
            EnsureReadable(1);
            return _reader.ReadBool();
        }

        public byte ReadByte(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 8))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            EnsureReadable(bitLength);
            return (byte)_reader.ReadBits(bitLength);
        }

        public ushort ReadUshort(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 16))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            EnsureReadable(bitLength);
            return (ushort)ByteOrderSupport.Swap(_reader.ReadBits(bitLength), bitLength, _byteOrder.Peek());
        }

        public uint ReadUint(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 32))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            EnsureReadable(bitLength);
            return (uint)ByteOrderSupport.Swap(_reader.ReadBits(bitLength), bitLength, _byteOrder.Peek());
        }

        public ulong ReadUlong(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 64))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            // The previous implementation read (bitLength - 32) bits and then a
            // further bitLength bits, consuming 2*bitLength-32 bits instead of
            // bitLength, so it never round-tripped against WriteBuffer.WriteUlong.
            EnsureReadable(bitLength);
            return ByteOrderSupport.Swap(_reader.ReadBits(bitLength), bitLength, _byteOrder.Peek());
        }

        public sbyte ReadSbyte(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 8))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            EnsureReadable(bitLength);
            return (sbyte)_reader.ReadSignedBits(bitLength);
        }

        public short ReadShort(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 16))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            EnsureReadable(bitLength);
            return (short)ReadSignedRaw(bitLength);
        }

        public int ReadInt(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 32))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            EnsureReadable(bitLength);
            return (int)ReadSignedRaw(bitLength);
        }

        public long ReadLong(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 64))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            // Same over-read as ReadUlong had.
            EnsureReadable(bitLength);
            return ReadSignedRaw(bitLength);
        }

        /// <summary>
        /// The big-endian path is exactly the previous implementation
        /// (<c>BitReader.ReadSignedBits</c>), untouched. The little-endian path
        /// cannot use it directly - it sign-extends before the byte order is
        /// known to be wrong - so it reads the raw bit pattern, byte-swaps it,
        /// and sign-extends the swapped pattern itself.
        /// </summary>
        private long ReadSignedRaw(int bitLength)
        {
            var order = _byteOrder.Peek();
            if (order == ByteOrder.LITTLE_ENDIAN && bitLength > 8 && bitLength % 8 == 0)
            {
                var swapped = ByteOrderSupport.Swap(_reader.ReadBits(bitLength), bitLength, order);
                return ByteOrderSupport.SignExtend(swapped, bitLength);
            }
            return _reader.ReadSignedBits(bitLength);
        }

        public float ReadFloat(String logicalName, int bitLength)
        {
            if (_byteOrder.Peek() == ByteOrder.LITTLE_ENDIAN)
            {
                // No mspec here declares a little-endian float/double field yet;
                // fail loud rather than silently reading it big-endian.
                throw new NotImplementedException("Little-endian floating point is not supported yet");
            }
            if (bitLength == 32)
            {
                // Reinterpret the 32 raw bits as IEEE-754. Going via
                // BitConverter.GetBytes would reorder them on a little-endian host.
                EnsureReadable(32);
                return BitConverter.Int32BitsToSingle((int)_reader.ReadBits(32));
            }
            // KNX DPT 9.x: a 16 bit fixed-point form, not IEEE-754 half precision.
            if (bitLength == 16)
            {
                EnsureReadable(16);
                var sign = _reader.ReadBool();
                var exp = (int)_reader.ReadBits(4);
                var mantissa = (int)_reader.ReadBits(11);
                if (sign)
                {
                    // Negative values are stored as a two's complement mantissa.
                    mantissa -= 2048;
                }
                return (float)(0.01 * mantissa * Math.Pow(2, exp));
            }
            throw new NotImplementedException("This encoding is currently not supported");
        }

        public double ReadDouble(String logicalName, int bitLength)
        {
            if (_byteOrder.Peek() == ByteOrder.LITTLE_ENDIAN)
            {
                throw new NotImplementedException("Little-endian floating point is not supported yet");
            }
            if (bitLength == 32)
            {
                // Was BitConverter.ToDouble over a 4 byte array, which always threw.
                return ReadFloat(logicalName, 32);
            }
            if (bitLength == 64)
            {
                EnsureReadable(64);
                return BitConverter.Int64BitsToDouble((long)_reader.ReadBits(64));
            }
            throw new NotImplementedException("This encoding is currently not supported");
        }

        public string ReadString(String logicalName, int bitLength, Encoding encoding)
        {
            if (encoding == null)
            {
                throw new ArgumentNullException(nameof(encoding));
            }
            if (bitLength < 0 || bitLength % 8 != 0)
            {
                throw new ArgumentOutOfRangeException(
                    nameof(bitLength), bitLength, "String lengths must be a whole number of bytes.");
            }

            EnsureReadable(bitLength);
            var raw = _reader.ReadBytes(bitLength / 8);
            // PLC strings are commonly padded to a fixed field width with NULs.
            return encoding.GetString(raw).TrimEnd('\0');
        }

        public byte[] ReadByteArray(String logicalName, int bitLength)
        {
            if (bitLength < 0 || bitLength % 8 != 0)
            {
                throw new ArgumentOutOfRangeException(
                    nameof(bitLength), bitLength, "Byte arrays must be a whole number of bytes.");
            }
            EnsureReadable(bitLength);
            return _reader.ReadBytes(bitLength / 8);
        }

        private void EnsureReadable(int bitLength)
        {
            if (bitLength < 0 || bitLength > _endBit - _reader.Position)
            {
                throw new ParseException("Not enough data in the read buffer.");
            }
        }

    }
}