#  Licensed to the Apache Software Foundation (ASF) under one
#  or more contributor license agreements.  See the NOTICE file
#  distributed with this work for additional information
#  regarding copyright ownership.  The ASF licenses this file
#  to you under the Apache License, Version 2.0 (the
#  "License"); you may not use this file except in compliance
#  with the License.  You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing,
#  software distributed under the License is distributed on an
#  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
#  KIND, either express or implied.  See the License for the
#  specific language governing permissions and limitations
#  under the License.
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBufferByteBased
from plc4py.utils.GenericTypes import ByteOrder

auch_crc_hi: bytearray = bytearray(
    b"\x00\xc1\x81\x40\x01\xc0\x80\x41\x01\xc0\x80\x41\x00\xc1\x81 \
        \x40\x01\xc0\x80\x41\x00\xc1\x81\x40\x00\xc1\x81\x40\x01\xc0 \
        \x80\x41\x01\xc0\x80\x41\x00\xc1\x81\x40\x00\xc1\x81\x40\x01 \
        \xc0\x80\x41\x00\xc1\x81\x40\x01\xc0\x80\x41\x01\xc0\x80\x41 \
        \x00\xc1\x81\x40\x01\xc0\x80\x41\x00\xc1\x81\x40\x00\xc1\x81 \
        \x40\x01\xc0\x80\x41\x00\xc1\x81\x40\x01\xc0\x80\x41\x01\xc0 \
        \x80\x41\x00\xc1\x81\x40\x00\xc1\x81\x40\x01\xc0\x80\x41\x01 \
        \xc0\x80\x41\x00\xc1\x81\x40\x01\xc0\x80\x41\x00\xc1\x81\x40 \
        \x00\xc1\x81\x40\x01\xc0\x80\x41\x01\xc0\x80\x41\x00\xc1\x81 \
        \x40\x00\xc1\x81\x40\x01\xc0\x80\x41\x00\xc1\x81\x40\x01\xc0 \
        \x80\x41\x01\xc0\x80\x41\x00\xc1\x81\x40\x00\xc1\x81\x40\x01 \
        \xc0\x80\x41\x01\xc0\x80\x41\x00\xc1\x81\x40\x01\xc0\x80\x41 \
        \x00\xc1\x81\x40\x00\xc1\x81\x40\x01\xc0\x80\x41\x00\xc1\x81 \
        \x40\x01\xc0\x80\x41\x01\xc0\x80\x41\x00\xc1\x81\x40\x01\xc0 \
        \x80\x41\x00\xc1\x81\x40\x00\xc1\x81\x40\x01\xc0\x80\x41\x01 \
        \xc0\x80\x41\x00\xc1\x81\x40\x00\xc1\x81\x40\x01\xc0\x80\x41 \
        \x00\xc1\x81\x40\x01\xc0\x80\x41\x01\xc0\x80\x41\x00\xc1\x81 \
        \x40"
)

auch_crc_lo: bytearray = bytearray(
    b"\x00\xc0\xc1\x01\xc3\x03\x02\xc2\xc6\x06\x07\xc7\x05\xc5\xc4 \
        \x04\xcc\x0c\x0d\xcd\x0f\xcf\xce\x0e\x0a\xca\xcb\x0b\xc9\x09 \
        \x08\xc8\xd8\x18\x19\xd9\x1b\xdb\xda\x1a\x1e\xde\xdf\x1f\xdd \
        \x1d\x1c\xdc\x14\xd4\xd5\x15\xd7\x17\x16\xd6\xd2\x12\x13\xd3 \
        \x11\xd1\xd0\x10\xf0\x30\x31\xf1\x33\xf3\xf2\x32\x36\xf6\xf7 \
        \x37\xf5\x35\x34\xf4\x3c\xfc\xfd\x3d\xff\x3f\x3e\xfe\xfa\x3a \
        \x3b\xfb\x39\xf9\xf8\x38\x28\xe8\xe9\x29\xeb\x2b\x2a\xea\xee \
        \x2e\x2f\xef\x2d\xed\xec\x2c\xe4\x24\x25\xe5\x27\xe7\xe6\x26 \
        \x22\xe2\xe3\x23\xe1\x21\x20\xe0\xa0\x60\x61\xa1\x63\xa3\xa2 \
        \x62\x66\xa6\xa7\x67\xa5\x65\x64\xa4\x6c\xac\xad\x6d\xaf\x6f \
        \x6e\xae\xaa\x6a\x6b\xab\x69\xa9\xa8\x68\x78\xb8\xb9\x79\xbb \
        \x7b\x7a\xba\xbe\x7e\x7f\xbf\x7d\xbd\xbc\x7c\xb4\x74\x75\xb5 \
        \x77\xb7\xb6\x76\x72\xb2\xb3\x73\xb1\x71\x70\xb0\x50\x90\x91 \
        \x51\x93\x53\x52\x92\x96\x56\x57\x97\x55\x95\x94\x54\x9c\x5c \
        \x5d\x9d\x5f\x9f\x9e\x5e\x5a\x9a\x9b\x5b\x99\x59\x58\x98\x88 \
        \x48\x49\x89\x4b\x8b\x8a\x4a\x4e\x8e\x8f\x4f\x8d\x4d\x4c\x8c \
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


def parse_terminated_string(read_buffer: ReadBuffer, string_length) -> str:
    terminate: bool = False
    byte_list: bytearray = bytearray()
    while not terminate:
        next_byte: int = read_buffer.read_byte()
        if next_byte == 0x00:
            terminate = True
        else:
            byte_list.append(next_byte)
    result: str = byte_list.decode("UTF-8")
    return result


def parse_terminated_string_bytes(read_buffer: ReadBuffer, string_length) -> str:
    terminate: bool = False
    byte_list: bytearray = bytearray()
    while not terminate:
        if len(byte_list) == 16:
            break
        next_byte: int = read_buffer.read_byte()
        if next_byte == 0x00 or len(byte_list) == 16:
            terminate = True
        else:
            byte_list.append(next_byte)
    length = len(byte_list)
    for _ in range(1, 16 - length):
        read_buffer.read_byte()
    result: str = byte_list.decode("UTF-8")
    return result


def serialize_terminated_string(write_buffer, value, string_length):
    terminated_char_position = len(value.value)
    for i in range(len(value.value)):
        if value.value[i] == "":
            terminated_char_position = i
            break
    write_buffer.write_str(value.value, (terminated_char_position + 1) * 8)
