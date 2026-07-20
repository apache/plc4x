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

// SLMP (Seamless Message Protocol) / MELSEC Communication protocol.
//
// Specification: Mitsubishi "MELSEC Communication Protocol Reference Manual"
//   (SH(NA)-080008, public). All field layouts below are taken from that manual:
//     - 3E frame message format ....... section 5.2
//     - Subheader (50 00 / D0 00) ...... section 5.3
//     - Access route (3E fixed value) .. chapter 6
//     - Request/response data length ... section 5.3 (2-byte, little-endian)
//     - Commands (Batch/Random/block Read) chapter 7 / 8.1 / 8.3 / 8.4
//     - Device code list ............... section 8.1 (MELSEC-Q/L, 1-byte binary)
//     - Batch Read data layout ......... section 8.1 (binary, word units)
//     - Random Read data layout ........ section 8.3 (binary, word units)
//     - Multi-block Read data layout ... section 8.4 (binary, word units)
//
// Scope of this initial version: 3E binary frame, read-only, Batch Read
// (command 0x0401), Random Read (command 0x0403) and Batch Read Multiple Blocks
// (command 0x0406) in word units (subcommand 0x0000), plus Batch Write
// (command 0x1401) in word units (subcommand 0x0000). This is the wire layer
// only; typed value decoding
// (INT/WORD/DINT/REAL) and the device-addressing tag layer are intentionally
// NOT modelled here yet and will follow once the driver logic is built.
// Validated hardware-free via the ParserSerializer test suite.
//
// All multi-byte numeric fields are little-endian (transmitted least-significant byte first).
// The 2-byte subheader (0x50/0xD0 then 0x00) is modelled as two single bytes so
// its on-wire order is preserved regardless of the frame byte order.

[constants
    // Typical SLMP TCP port; the actual port is configured on the device and is
    // overridable by the driver, so this is only a default.
    [const uint 16 slmpDefaultPort 5007]
]

// Device codes for MELSEC-Q/L series commands (subcommand 0x0000 / 0x0001),
// 1-byte binary, from the device code list in SH-080008 section 8.1.
// Only the word/bit devices needed by the read-only road-map are listed.
//
// These 1-byte binary device codes are confirmed identical in the Mitsubishi
// MELSEC iQ-F FX5 SLMP manual (JY997D56001) -- D=A8 W=B4 R=AF M=90 X=9C Y=9D
// B=A0 -- so the same frame works on FX5 hardware (e.g. FX5S) as well as
// MELSEC-Q/L. (FX5 also offers a 2-byte/extended code form via subcommand
// 0x0002/0x0003 for ZR and large device numbers, which is out of scope here.)
[enum uint 8 SlmpDeviceCode
    ['0xA8' D]  // Data register   (word, decimal addressing)
    ['0xB4' W]  // Link register   (word, hex addressing)
    ['0xAF' R]  // File register    (word, decimal addressing)
    ['0xC2' TN] // Timer, current value (word, decimal addressing) -- used by the
                // Random Read worked example in SH-080008 section 8.3
    ['0x90' M]  // Internal relay  (bit,  decimal addressing)
    ['0x9C' X]  // Input            (bit,  hex addressing)
    ['0x9D' Y]  // Output          (bit,  hex addressing)
    ['0xA0' B]  // Link relay      (bit,  hex addressing)
]

// 3E frame. The Ethernet/TCP header is added by the transport and is not part
// of this message. Request and response are distinguished by the subheader byte
// (0x50 request / 0xD0 response).
[discriminatedType SlmpMessage byteOrder='"LITTLE_ENDIAN"' unsignedIntegerEncoding='"unsigned-binary"' signedIntegerEncoding='"twos-complement"' floatEncoding='"IEEE754"' stringEncoding='"UTF8"'
    [discriminator uint 8 subHeader]                  // 0x50 request, 0xD0 response
    [const         uint 8 subHeaderReserved 0x00]     // second subheader byte is always 0x00
    [typeSwitch subHeader
        ['0x50' SlmpRequestFrame3E
            // Access route - fixed value for 3E frame (chapter 6): 00 FF FF 03 00
            [const    uint 8  network                  0x00]
            [const    uint 8  pcNumber                 0xFF]
            [const    uint 16 requestDestModuleIoNo    0x03FF]
            [const    uint 8  requestDestModuleStation 0x00]
            // Number of bytes from the monitoring timer up to the end of the request data.
            // = monitoringTimer(2) + command(2) + subCommand(2) + requestData
            [implicit uint 16 requestDataLength        '6 + requestData.lengthInBytes']
            [simple   uint 16 monitoringTimer]                 // 0x0000 = wait infinitely
            [simple   uint 16 command]                         // 0x0401 = Batch Read
            [simple   uint 16 subCommand]                      // 0x0000 = word units
            [simple   SlmpRequestData('command') requestData]
        ]
        ['0xD0' SlmpResponseFrame3E
            [const    uint 8  network                  0x00]
            [const    uint 8  pcNumber                 0xFF]
            [const    uint 16 requestDestModuleIoNo    0x03FF]
            [const    uint 8  requestDestModuleStation 0x00]
            // Number of bytes from the end code up to the end of the response data.
            // = endCode(2) + responseData
            [implicit uint 16 responseDataLength       '2 + COUNT(responseData)']
            [simple   uint 16 endCode]                         // 0x0000 = normal completion
            // Raw payload. On normal completion this is the read words (little-endian,
            // 2 bytes per point); on abnormal completion it carries error information.
            // Typed decoding is left to the driver layer.
            [array    byte    responseData             count 'responseDataLength - 2']
        ]
    ]
]

// Request data, dispatched by command. Batch Read (0x0401) and Random Read
// (0x0403) are modelled here; both are read-only, word units (subcommand 0x0000).
[discriminatedType SlmpRequestData(uint 16 command)
    [typeSwitch command
        ['0x0401' SlmpReadRequest
            // Binary, MELSEC-Q/L series: device number (3 bytes, little-endian)
            // then device code (1 byte). See SH-080008 section 8.1.
            [simple uint 24        headDeviceNumber]
            [simple SlmpDeviceCode deviceCode]
            // Number of points to read. In word units this is the number of 16-bit words.
            [simple uint 16        numberOfPoints]
        ]
        ['0x0403' SlmpRandomReadRequest
            // Random Read in word units (SH-080008 section 8.3): read discontinuous
            // devices, given as a list of word-access points followed by a list of
            // double-word-access points. The point counts are single bytes; each
            // device is a 3-byte device number + 1-byte device code (SlmpDeviceSpec).
            [simple uint 8         numberOfWordAccessPoints]
            [simple uint 8         numberOfDoubleWordAccessPoints]
            [array  SlmpDeviceSpec wordAccessDevices       count 'numberOfWordAccessPoints']
            [array  SlmpDeviceSpec doubleWordAccessDevices count 'numberOfDoubleWordAccessPoints']
        ]
        ['0x0406' SlmpMultiBlockReadRequest
            // Batch Read multiple blocks in word units (SH-080008 section 8.4): read a
            // number of word-device blocks followed by a number of bit-device blocks,
            // where each block is a run of consecutive devices. The two block counts are
            // single bytes; each block (SlmpDeviceBlock) is a 3-byte head device number,
            // a 1-byte device code and a 2-byte number of points. Bit-device blocks are
            // read in 16-bit word units (one point = 16 bits), so the wire layout of a
            // bit-device block is identical to that of a word-device block.
            [simple uint 8          numberOfWordDeviceBlocks]
            [simple uint 8          numberOfBitDeviceBlocks]
            [array  SlmpDeviceBlock wordDeviceBlocks count 'numberOfWordDeviceBlocks']
            [array  SlmpDeviceBlock bitDeviceBlocks  count 'numberOfBitDeviceBlocks']
        ]
        ['0x1401' SlmpWriteRequest
            // Batch Write in word units (SH-080008 section 8.2 "Batch Write"): the same
            // device addressing as SlmpReadRequest (0x0401), followed by the data to write.
            // numberOfPoints is the number of 16-bit words; writeData is exactly that many
            // little-endian words carried as raw bytes (2 * numberOfPoints). Typed encoding
            // (INT/WORD/DINT/REAL) is owned by the driver layer, symmetric with how the read
            // response leaves responseData as raw bytes for the driver to decode.
            [simple uint 24        headDeviceNumber]
            [simple SlmpDeviceCode deviceCode]
            [simple uint 16        numberOfPoints]
            [array  byte           writeData count 'numberOfPoints * 2']
        ]
    ]
]

// A single device specification used by commands that address discontinuous
// devices (e.g. Random Read 0x0403): a device number (3 bytes, little-endian)
// followed by the 1-byte device code. See SH-080008 section 8.1 / 8.3.
[type SlmpDeviceSpec
    [simple uint 24        deviceNumber]
    [simple SlmpDeviceCode deviceCode]
]

// A single block of consecutive devices used by commands that read whole blocks
// (e.g. Batch Read Multiple Blocks 0x0406): a head device number (3 bytes,
// little-endian) and 1-byte device code identifying the start of the run,
// followed by the 2-byte number of points (consecutive devices) in the block.
// This is the same field layout as a single-block Batch Read request
// (SlmpReadRequest). See SH-080008 section 8.4.
[type SlmpDeviceBlock
    [simple uint 24        headDeviceNumber]
    [simple SlmpDeviceCode deviceCode]
    [simple uint 16        numberOfPoints]
]
