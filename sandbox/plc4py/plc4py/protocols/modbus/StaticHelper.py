#  Licensed to the Apache Software Foundation (ASF) under one
#  or more contributor license agreements.  See the NOTICE file
#  distributed with this work for additional information
#  regarding copyright ownership.  The ASF licenses this file
#  to you under the Apache License, Version 2.0 (the
#  "License"); you may not use this file except in compliance
#  with the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing,
#  software distributed under the License is distributed on an
#  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
#  KIND, either express or implied.  See the License for the
#  specific language governing permissions and limitations
#  under the License.

from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.spi.generation.WriteBuffer import WriteBufferByteBased
from plc4py.utils.GenericTypes import ByteOrder

auch_crc_hi: bytearray = bytearray(
    b"\x00\xC1\x81\x40\x01\xC0\x80\x41\x01\xC0\x80\x41\x00\xC1\x81 \
        \x40\x01\xC0\x80\x41\x00\xC1\x81\x40\x00\xC1\x81\x40\x01\xC0 \
        \x80\x41\x01\xC0\x80\x41\x00\xC1\x81\x40\x00\xC1\x81\x40\x01 \
        \xC0\x80\x41\x00\xC1\x81\x40\x01\xC0\x80\x41\x01\xC0\x80\x41 \
        \x00\xC1\x81\x40\x01\xC0\x80\x41\x00\xC1\x81\x40\x00\xC1\x81 \
        \x40\x01\xC0\x80\x41\x00\xC1\x81\x40\x01\xC0\x80\x41\x01\xC0 \
        \x80\x41\x00\xC1\x81\x40\x00\xC1\x81\x40\x01\xC0\x80\x41\x01 \
        \xC0\x80\x41\x00\xC1\x81\x40\x01\xC0\x80\x41\x00\xC1\x81\x40 \
        \x00\xC1\x81\x40\x01\xC0\x80\x41\x01\xC0\x80\x41\x00\xC1\x81 \
        \x40\x00\xC1\x81\x40\x01\xC0\x80\x41\x00\xC1\x81\x40\x01\xC0 \
        \x80\x41\x01\xC0\x80\x41\x00\xC1\x81\x40\x00\xC1\x81\x40\x01 \
        \xC0\x80\x41\x01\xC0\x80\x41\x00\xC1\x81\x40\x01\xC0\x80\x41 \
        \x00\xC1\x81\x40\x00\xC1\x81\x40\x01\xC0\x80\x41\x00\xC1\x81 \
        \x40\x01\xC0\x80\x41\x01\xC0\x80\x41\x00\xC1\x81\x40\x01\xC0 \
        \x80\x41\x00\xC1\x81\x40\x00\xC1\x81\x40\x01\xC0\x80\x41\x01 \
        \xC0\x80\x41\x00\xC1\x81\x40\x00\xC1\x81\x40\x01\xC0\x80\x41 \
        \x00\xC1\x81\x40\x01\xC0\x80\x41\x01\xC0\x80\x41\x00\xC1\x81 \
        \x40"
)

auch_crc_lo: bytearray = bytearray(
    b"\x00\xC0\xC1\x01\xC3\x03\x02\xC2\xC6\x06\x07\xC7\x05\xC5\xC4 \
        \x04\xCC\x0C\x0D\xCD\x0F\xCF\xCE\x0E\x0A\xCA\xCB\x0B\xC9\x09 \
        \x08\xC8\xD8\x18\x19\xD9\x1B\xDB\xDA\x1A\x1E\xDE\xDF\x1F\xDD \
        \x1D\x1C\xDC\x14\xD4\xD5\x15\xD7\x17\x16\xD6\xD2\x12\x13\xD3 \
        \x11\xD1\xD0\x10\xF0\x30\x31\xF1\x33\xF3\xF2\x32\x36\xF6\xF7 \
        \x37\xF5\x35\x34\xF4\x3C\xFC\xFD\x3D\xFF\x3F\x3E\xFE\xFA\x3A \
        \x3B\xFB\x39\xF9\xF8\x38\x28\xE8\xE9\x29\xEB\x2B\x2A\xEA\xEE \
        \x2E\x2F\xEF\x2D\xED\xEC\x2C\xE4\x24\x25\xE5\x27\xE7\xE6\x26 \
        \x22\xE2\xE3\x23\xE1\x21\x20\xE0\xA0\x60\x61\xA1\x63\xA3\xA2 \
        \x62\x66\xA6\xA7\x67\xA5\x65\x64\xA4\x6C\xAC\xAD\x6D\xAF\x6F \
        \x6E\xAE\xAA\x6A\x6B\xAB\x69\xA9\xA8\x68\x78\xB8\xB9\x79\xBB \
        \x7B\x7A\xBA\xBE\x7E\x7F\xBF\x7D\xBD\xBC\x7C\xB4\x74\x75\xB5 \
        \x77\xB7\xB6\x76\x72\xB2\xB3\x73\xB1\x71\x70\xB0\x50\x90\x91 \
        \x51\x93\x53\x52\x92\x96\x56\x57\x97\x55\x95\x94\x54\x9C\x5C \
        \x5D\x9D\x5F\x9F\x9E\x5E\x5A\x9A\x9B\x5B\x99\x59\x58\x98\x88 \
        \x48\x49\x89\x4B\x8B\x8A\x4A\x4E\x8E\x8F\x4F\x8D\x4D\x4C\x8C \
        \x44\x84\x85\x45\x87\x47\x46\x86\x82\x42\x43\x83\x41\x81\x80 \
        \x40"
)


def rtu_crc_check(address: int, pdu: ModbusPDU) -> int:
    # Using the algorithm from PI_MBUS_300.pdf page 121
    write_buffer: WriteBufferByteBased = WriteBufferByteBased(
        pdu.length_in_bytes() + 2, byte_order=ByteOrder.LITTLE_ENDIAN
    )
    write_buffer.write_unsigned_short(address, 8)
    pdu.serialize(write_buffer)
    m_view: memoryview = write_buffer.get_bytes()
    uch_crc_hi: int = 0xFF
    uch_crc_lo: int = 0xFF
    u_index: int
    for b in m_view:
        u_index = (uch_crc_hi ^ b) & 0xFF
        uch_crc_hi = uch_crc_lo ^ auch_crc_hi[u_index]
        uch_crc_lo = auch_crc_lo[u_index]
    return ((uch_crc_hi << 8) & 0xFFFF) | (uch_crc_lo & 0x00FF)


# 8 Bit checksum, (byte) transported as 2 characters
def ascii_lrc_check(address: int, pdu: ModbusPDU) -> int:
    write_buffer: WriteBufferByteBased = WriteBufferByteBased(
        pdu.length_in_bytes() + 2, byte_order=ByteOrder.LITTLE_ENDIAN
    )
    write_buffer.write_unsigned_byte(address, 8)
    pdu.serialize(write_buffer)
    m_view: memoryview = write_buffer.get_bytes()

    # A procedure for generating an LRC is:
    # 1. Add all bytes in the message, (byte) excluding the starting ‘colon’ and ending
    #    CRLF. Add them into an 8–bit tag, (byte) so that carries will be discarded.
    # 2. Subtract the final tag value from FF hex (all 1’s), (byte) to produce the
    #    ones–complement.
    # 3. Add 1 to produce the twos–complement
    # 4. Convert the 8 bit checksum into it's 16 bit (2 char) hex representation.
    #    (Handled in the transport layer)
    lrc: int = 0
    for b in m_view:
        lrc = lrc + b

    lrc = -lrc
    return lrc & 0xFF
