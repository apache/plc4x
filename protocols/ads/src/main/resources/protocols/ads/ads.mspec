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

[constants
    [const          uint 16     adsTcpDefaultPort 48898]
]

////////////////////////////////////////////////////////////////
// External types
////////////////////////////////////////////////////////////////

[enum PlcValueType external='true']

////////////////////////////////////////////////////////////////
// AMS/TCP Packet
////////////////////////////////////////////////////////////////

[type AmsTCPPacket byteOrder='"LITTLE_ENDIAN"' unsignedIntegerEncoding='"unsigned-binary"' signedIntegerEncoding='"twos-complement"' floatEncoding='"IEEE754"' stringEncoding='"UTF8"'
    // AMS/TCP Header	6 bytes	contains the tcpLength of the data packet.
    // These bytes must be set to 0.
    [reserved   uint       16       '0x0000'                            ]
    // This array contains the length of the data packet.
    // It consists of the AMS-Header and the enclosed ADS data. The unit is bytes.
    [implicit   uint       32       length  'userdata.lengthInBytes'  ]
    // The AMS packet to be sent.
    [simple AmsPacket    userdata                                     ]
]

////////////////////////////////////////////////////////////////
// AMS/Serial Packet
////////////////////////////////////////////////////////////////

// If an AMS serial frame has been received and the frame is OK (magic cookie OK, CRC OK, correct fragment number etc.),
// then the receiver has to send an acknowledge frame, to inform the transmitter that the frame has arrived.
//
// @see <a href="https://infosys.beckhoff.com/content/1033/tcadsamsserialspec/html/tcamssericalspec_amsframe.htm?id=8115637053270715044">TwinCAT AMS via RS232 Specification</a>
[type AmsSerialAcknowledgeFrame
    // Id for detecting an AMS serial frame.
    [simple     uint        16  magicCookie        ]
    // Address of the sending participant. This value can always be set to 0 for an RS232 communication,
    // since it is a 1 to 1 connection and hence the participants are unique.
    [simple     int          8  transmitterAddress ]
    // Receiver’s address. This value can always be set to 0 for an RS232 communication, since it is a 1 to 1
    // connection and hence the participants are unique.
    [simple     int          8  receiverAddress    ]
    // Number of the frame sent. Once the number 255 has been sent, it starts again from 0. The receiver checks this
    // number with an internal counter.
    [simple     int          8  fragmentNumber     ]
    // The max. length of the AMS packet to be sent is 255. If larger AMS packets are to be sent then they have to be
    // fragmented (not published at the moment).
    [simple     int          8  length             ]
    [simple     uint        16  crc                ]
]

// An AMS packet can be transferred via RS232 with the help of an AMS serial frame.
// The actual AMS packet is in the user data field of the frame.
// The max. length of the AMS packet is limited to 255 bytes.
// Therefore the max. size of an AMS serial frame is 263 bytes.
// The fragment number is compared with an internal counter by the receiver.
// The frame number is simply accepted and not checked when receiving the first AMS frame or in case a timeout is
// exceeded. The CRC16 algorithm is used for calculating the checksum.
// @see <a href="https://infosys.beckhoff.com/content/1033/tcadsamsserialspec/html/tcamssericalspec_amsframe.htm?id=8115637053270715044">TwinCAT AMS via RS232 Specification</a>
[type AmsSerialFrame
    // Id for detecting an AMS serial frame.
    [simple     uint        16  magicCookie        ]
    // Address of the sending participant. This value can always be set to 0 for an RS232 communication,
    // since it is a 1 to 1 connection and hence the participants are unique.
    [simple     int          8  transmitterAddress ]
    // Receiver’s address. This value can always be set to 0 for an RS232 communication, since it is a 1 to 1
    // connection and hence the participants are unique.
    [simple     int          8  receiverAddress    ]
    // Number of the frame sent. Once the number 255 has been sent, it starts again from 0. The receiver checks this
    // number with an internal counter.
    [simple     int          8  fragmentNumber     ]
    // The max. length of the AMS packet to be sent is 255. If larger AMS packets are to be sent then they have to be
    // fragmented (not published at the moment).
    [simple     int          8  length             ]
    // The AMS packet to be sent.
    [simple AmsPacket           userdata           ]
    [simple     uint        16  crc                ]
]

// In case the transmitter does not receive a valid acknowledgement after multiple transmission, then a reset frame is
// sent. In this way the receiver is informed that a new communication is running and the receiver then accepts the
// fragment number during the next AMS-Frame, without carrying out a check.
[type AmsSerialResetFrame
    // Id for detecting an AMS serial frame.
    [simple     uint        16  magicCookie        ]
    // Address of the sending participant. This value can always be set to 0 for an RS232 communication,
    // since it is a 1 to 1 connection and hence the participants are unique.
    [simple     int          8  transmitterAddress ]
    // Receiver’s address. This value can always be set to 0 for an RS232 communication, since it is a 1 to 1
    // connection and hence the participants are unique.
    [simple     int          8  receiverAddress    ]
    // Number of the frame sent. Once the number 255 has been sent, it starts again from 0. The receiver checks this
    // number with an internal counter.
    [simple     int          8  fragmentNumber     ]
    // The max. length of the AMS packet to be sent is 255. If larger AMS packets are to be sent then they have to be
    // fragmented (not published at the moment).
    [simple     int          8  length     ]
    [simple     uint        16  crc                ]
]

////////////////////////////////////////////////////////////////
// AMS Common
////////////////////////////////////////////////////////////////

[discriminatedType AmsPacket
    // AMS Header	32 bytes	The AMS/TCP-Header contains the addresses of the transmitter and receiver. In addition the AMS error code , the ADS command Id and some other information.
    // This is the AmsNetId of the station, for which the packet is intended. Remarks see below.
    [simple        AmsNetId        targetAmsNetId                            ]
    // This is the AmsPort of the station, for which the packet is intended.
    [simple        uint        16  targetAmsPort                             ]
    // This contains the AmsNetId of the station, from which the packet was sent.
    [simple        AmsNetId        sourceAmsNetId                            ]
    // This contains the AmsPort of the station, from which the packet was sent.
    [simple        uint        16  sourceAmsPort                             ]
    // 2 bytes.
    [discriminator CommandId       commandId                                 ]
    // 2 bytes. (I set these as constants in order to minimize the input needed for creating requests)
    [const         bit             initCommand            false              ]
    [const         bit             updCommand             false              ]
    [const         bit             timestampAdded         false              ]
    [const         bit             highPriorityCommand    false              ]
    [const         bit             systemCommand          false              ]
    [const         bit             adsCommand             true               ]
    [const         bit             noReturn               false              ]
    [discriminator bit             response                                  ]
    [const         bit             broadcast              false              ]
    [reserved      int 7           '0x0'                                     ]
    // 4 bytes	Size of the data range. The unit is byte.
    [implicit      uint        32  length   'lengthInBytes - 32'             ]
    // 4 bytes	AMS error number. See ADS Return Codes.
    [simple        ReturnCode      errorCode                                 ]
    // free usable field of 4 bytes
    // 4 bytes	Free usable 32 bit array. Usually this array serves to send an Id. This Id makes is possible to assign a received response to a request, which was sent before.
    [simple        uint        32  invokeId                                  ]
    // The payload
    // TODO: In case of an error code that is not 0, we might not have a payload at all
    [typeSwitch errorCode, commandId, response
        ['OK', 'INVALID', 'false' AdsInvalidRequest]
        ['OK', 'INVALID', 'true' AdsInvalidResponse]

        ['OK', 'ADS_READ_DEVICE_INFO', 'false' AdsReadDeviceInfoRequest]
        ['OK', 'ADS_READ_DEVICE_INFO', 'true' AdsReadDeviceInfoResponse
            // 4 bytes	ADS error number.
            [simple ReturnCode result]
            // Version	1 byte	Major version number
            [simple uint 8  majorVersion]
            // Version	1 byte	Minor version number
            [simple uint 8  minorVersion]
            // Build	2 bytes	Build number
            [simple uint 16  version]
            // Name	16 bytes	Name of ADS device
            [array byte  device count '16']
        ]

        ['OK', 'ADS_READ', 'false' AdsReadRequest
            // 4 bytes	Index Group of the data which should be read.
            [simple uint 32 indexGroup]
            // 4 bytes	Index Offset of the data which should be read.
            [simple uint 32 indexOffset]
            // 4 bytes	Length of the data (in bytes) which should be read.
            [simple uint 32 length]
        ]
        ['OK', 'ADS_READ', 'true' AdsReadResponse
            // 4 bytes	ADS error number
            [simple ReturnCode result]
            // 4 bytes	Length of data which are supplied back.
            [implicit uint 32 length 'COUNT(data)']
            // n bytes	Data which are supplied back.
            [array byte data count 'length']
        ]

        ['OK', 'ADS_WRITE', 'false' AdsWriteRequest
            // 4 bytes	Index Group of the data which should be written.
            [simple uint 32 indexGroup]
            // 4 bytes	Index Offset of the data which should be written.
            [simple uint 32 indexOffset]
            // 4 bytes	Length of the data (in bytes) which should be written.
            [implicit uint 32 length 'COUNT(data)']
            // n bytes	Data which are written in the ADS device.
            [array byte data count 'length']
        ]
        ['OK', 'ADS_WRITE', 'true' AdsWriteResponse
            // 4 bytes	ADS error number
            [simple ReturnCode result]
        ]

        ['OK', 'ADS_READ_STATE', 'false' AdsReadStateRequest]
        ['OK', 'ADS_READ_STATE', 'true' AdsReadStateResponse
            // 4 bytes	ADS error number
            [simple ReturnCode result]
            // 2 bytes	New ADS status (see data type ADSSTATE of the ADS-DLL).
            [simple uint 16 adsState]
            // 2 bytes	New device status.
            [simple uint 16 deviceState]
        ]

        ['OK', 'ADS_WRITE_CONTROL', 'false' AdsWriteControlRequest
            // 2 bytes	New ADS status (see data type ADSSTATE of the ADS-DLL).
            [simple uint 16 adsState]
            // 2 bytes	New device status.
            [simple uint 16 deviceState]
            // 4 bytes	Length of data in byte.
            [implicit uint 32 length 'COUNT(data)']
            // n bytes	Additional data which are sent to the ADS device
            [array byte data count 'length']
        ]
        ['OK', 'ADS_WRITE_CONTROL', 'true' AdsWriteControlResponse
            // 4 bytes	ADS error number
            [simple ReturnCode result]
        ]

        ['OK', 'ADS_ADD_DEVICE_NOTIFICATION', 'false' AdsAddDeviceNotificationRequest
            // 4 bytes	Index Group of the data, which should be sent per notification.
            [simple     uint 32      indexGroup      ]
            // 4 bytes	Index Offset of the data, which should be sent per notification.
            [simple     uint 32      indexOffset     ]
            // 4 bytes	Index Offset of the data, which should be sent per notification.
            // 4 bytes	Length of data in bytes, which should be sent per notification.
            [simple     uint 32      length          ]
            // 4 bytes	The type of subscription.
            [simple     AdsTransMode transmissionMode]
            // 4 bytes	At the latest after this time, the ADS Device Notification is called. The unit is 1ms.
            [simple     uint 32      maxDelayInMs    ]
            // 4 bytes	The ADS server checks if the value changes in this time slice. The unit is 1ms
            [simple     uint 32      cycleTimeInMs   ]
            // 16bytes	Must be set to 0
            [reserved   uint 64      '0x0000'        ]
            [reserved   uint 64      '0x0000'        ]
        ]
        ['OK', 'ADS_ADD_DEVICE_NOTIFICATION', 'true' AdsAddDeviceNotificationResponse
            // 4 bytes	ADS error number
            [simple ReturnCode result]
            // 4 bytes	Handle of notification
            [simple uint 32 notificationHandle]
        ]

        ['OK', 'ADS_DELETE_DEVICE_NOTIFICATION', 'false' AdsDeleteDeviceNotificationRequest
            // 4 bytes	Handle of notification
            [simple uint 32 notificationHandle]
        ]
        ['OK', 'ADS_DELETE_DEVICE_NOTIFICATION', 'true' AdsDeleteDeviceNotificationResponse
            // 4 bytes	ADS error number
            [simple ReturnCode result]
        ]

        ['OK', 'ADS_DEVICE_NOTIFICATION', 'false' AdsDeviceNotificationRequest
            // 4 bytes	Size of data in byte.
            [simple uint 32 length]
            // 4 bytes	Number of elements of type AdsStampHeader.
            [simple uint 32 stamps]
            // n bytes	Array with elements of type AdsStampHeader.
            [array AdsStampHeader adsStampHeaders count 'stamps']
        ]
        ['OK', 'ADS_DEVICE_NOTIFICATION', 'true' AdsDeviceNotificationResponse]

        ['OK', 'ADS_READ_WRITE', 'false' AdsReadWriteRequest
            // 4 bytes	Index Group of the data which should be written.
            [simple uint 32 indexGroup]
            // 4 bytes	Index Offset of the data which should be written.
            [simple uint 32 indexOffset]
            // 4 bytes	Length of data in bytes, which should be read.
            [simple uint 32 readLength]
            // 4 bytes	Length of the data (in bytes) which should be written. (if it's ADSIGRP_MULTIPLE_READ_WRITE, this is 16 otherwise 12)
            [implicit uint 32 writeLength '(COUNT(items) * ((indexGroup == 61570) ? 16 : 12)) + COUNT(data)']
            // Only if the indexGroup implies a sum-read response, will the indexOffset indicate the number of elements. (ADSIGRP_MULTIPLE_READ, ADSIGRP_MULTIPLE_WRITE, ADSIGRP_MULTIPLE_READ_WRITE)
            [array  AdsMultiRequestItem('indexGroup') items count '((indexGroup == 61568) || (indexGroup == 61569) || (indexGroup == 61570)) ? indexOffset : 0']
            // n bytes	Data which are written in the ADS device.
            [array byte data count 'writeLength - (COUNT(items) * 12)']
        ]
        ['OK', 'ADS_READ_WRITE', 'true' AdsReadWriteResponse
            // 4 bytes	ADS error number
            [simple ReturnCode result]
            // 4 bytes	Length of data in byte.
            [implicit uint 32 length  'COUNT(data)']
            // n bytes Additional data which are sent to the ADS device
            [array byte data count 'length']
        ]
        [AdsErrorResponse
        ]
    ]
]

[enum uint 16 CommandId
    ['0x0000' INVALID                       ]
    ['0x0001' ADS_READ_DEVICE_INFO          ]
    ['0x0002' ADS_READ                      ]
    ['0x0003' ADS_WRITE                     ]
    ['0x0004' ADS_READ_STATE                ]
    ['0x0005' ADS_WRITE_CONTROL             ]
    ['0x0006' ADS_ADD_DEVICE_NOTIFICATION   ]
    ['0x0007' ADS_DELETE_DEVICE_NOTIFICATION]
    ['0x0008' ADS_DEVICE_NOTIFICATION       ]
    ['0x0009' ADS_READ_WRITE                ]
]

// It is not only possible to exchange data between TwinCAT modules on one PC, it is even possible to do so by ADS
// methods between multiple TwinCAT PC's on the network.
// <p>
// Every PC on the network can be uniquely identified by a TCP/IP address, such as "172.1.2.16". The AdsAmsNetId is an
// extension of the TCP/IP address and identifies a TwinCAT message router, e.g. "172.1.2.16.1.1". TwinCAT message
// routers exist on every TwinCAT PC, and on every Beckhoff BCxxxx bus controller (e.g. BC3100, BC8100, BC9000, ...).
// <p>
// The AmsNetId consists of 6 bytes and addresses the transmitter or receiver. One possible AmsNetId would be e.g.
// "172.16.17.10.1.1". The storage arrangement in this example is as follows:
// <p>
// _____0     1     2     3     4     5
// __+-----------------------------------+
// 0 | 127 |  16 |  17 |  10 |   1 |   1 |
// __+-----------------------------------+
// <p>
// The AmsNetId is purely logical and has usually no relation to the IP address. The AmsNetId is configured at the
// target system. At the PC for this the TwinCAT System Control is used. If you use other hardware, see the considering
// documentation for notes about settings of the AMS NetId.
// @see <a href="https://infosys.beckhoff.com/content/1033/tcadscommon/html/tcadscommon_identadsdevice.htm?id=3991659524769593444">ADS device identification</a>
[type AmsNetId
    [simple     uint        8   octet1            ]
    [simple     uint        8   octet2            ]
    [simple     uint        8   octet3            ]
    [simple     uint        8   octet4            ]
    [simple     uint        8   octet5            ]
    [simple     uint        8   octet6            ]
]

[discriminatedType AdsMultiRequestItem(uint 32 indexGroup)
    [typeSwitch indexGroup
        // ReservedIndexGroups.ADSIGRP_MULTIPLE_READ
        ['61568' AdsMultiRequestItemRead
            // 4 bytes	Index Group of the data which should be written.
            [simple uint 32 itemIndexGroup]
            // 4 bytes	Index Offset of the data which should be written.
            [simple uint 32 itemIndexOffset]
            // 4 bytes	Length of data in bytes, which should be read.
            [simple uint 32 itemReadLength]
        ]
        // ReservedIndexGroups.ADSIGRP_MULTIPLE_WRITE
        ['61569' AdsMultiRequestItemWrite
            // 4 bytes	Index Group of the data which should be written.
            [simple uint 32 itemIndexGroup]
            // 4 bytes	Index Offset of the data which should be written.
            [simple uint 32 itemIndexOffset]
            // 4 bytes	Length of the data (in bytes) which should be written.
            [simple uint 32 itemWriteLength]
        ]
        // ReservedIndexGroups.ADSIGRP_MULTIPLE_READ_WRITE
        ['61570' AdsMultiRequestItemReadWrite
            // 4 bytes	Index Group of the data which should be written.
            [simple uint 32 itemIndexGroup]
            // 4 bytes	Index Offset of the data which should be written.
            [simple uint 32 itemIndexOffset]
            // 4 bytes	Length of data in bytes, which should be read.
            [simple uint 32 itemReadLength]
            // 4 bytes	Length of the data (in bytes) which should be written.
            [simple uint 32 itemWriteLength]
        ]
    ]
]

[type AdsStampHeader
    // 8 bytes	The timestamp is coded after the Windows FILETIME format. I.e. the value contains the number of the nano seconds, which passed since 1.1.1601. In addition, the local time change is not considered. Thus the time stamp is present as universal Coordinated time (UTC).
    [simple uint 64 timestamp]
    // 4 bytes	Number of elements of type AdsNotificationSample.
    [simple uint 32 samples]
    // n bytes	Array with elements of type AdsNotificationSample.
    [array AdsNotificationSample adsNotificationSamples count 'samples']
]

[type AdsNotificationSample
    // 4 bytes	Handle of notification
    [simple uint 32 notificationHandle]
    // 4 Bytes	Size of data range in bytes.
    [simple uint 32 sampleSize]
    // n Bytes	Data
    [array byte data count 'sampleSize']
]

[dataIo DataItem(PlcValueType plcValueType, int 32 stringLength) byteOrder='"LITTLE_ENDIAN"' unsignedIntegerEncoding='"unsigned-binary"' signedIntegerEncoding='"twos-complement"' floatEncoding='"IEEE754"' stringEncoding='"UTF8"'
    [typeSwitch plcValueType
        // -----------------------------------------
        // Bit
        // -----------------------------------------
        ['BOOL' BOOL
            [reserved uint 7  '0x00']
            [simple   bit     value]
        ]

        // -----------------------------------------
        // Bit-strings
        // -----------------------------------------
        // 1 byte
        ['BYTE' BYTE
            [simple uint 8 value]
        ]
        // 2 byte (16 bit)
        ['WORD' WORD
            [simple uint 16 value]
        ]
        // 4 byte (32 bit)
        ['DWORD' DWORD
            [simple uint 32 value]
        ]
        // 8 byte (64 bit)
        ['LWORD' LWORD
            [simple uint 64 value]
        ]

        // -----------------------------------------
        // Integers
        // -----------------------------------------
        // 8 bit:
        ['SINT' SINT
            [simple int 8 value]
        ]
        ['USINT' USINT
            [simple uint 8 value]
        ]
        // 16 bit:
        ['INT' INT
            [simple int 16 value]
        ]
        ['UINT' UINT
            [simple uint 16 value]
        ]
        // 32 bit:
        ['DINT' DINT
            [simple int 32 value]
        ]
        ['UDINT' UDINT
            [simple uint 32 value]
        ]
        // 64 bit:
        ['LINT' LINT
            [simple int 64 value]
        ]
        ['ULINT' ULINT
            [simple uint 64 value]
        ]

        // -----------------------------------------
        // Floating point values
        // -----------------------------------------
        ['REAL' REAL
            [simple float 32  value]
        ]
        ['LREAL' LREAL
            [simple float 64 value]
        ]

        // -----------------------------------------
        // Characters & Strings
        // -----------------------------------------
        ['CHAR' CHAR
            [simple   string 8                       value    stringEncoding='"WINDOWS1252"']
        ]
        ['WCHAR' WCHAR
            [simple   string 16                      value    stringEncoding='"UTF16LE"'    ]
        ]
        ['STRING' STRING
            [simple   vstring 'stringLength * 8'     value    stringEncoding='"WINDOWS1252"']
            [reserved uint 8                         '0x00'                            ]
        ]
        ['WSTRING' WSTRING
            [simple   vstring 'stringLength * 8 * 2' value    stringEncoding='"UTF16LE"'    ]
            [reserved uint 16                        '0x0000'                          ]
        ]

        // -----------------------------------------
        // Date & Times
        // -----------------------------------------
        // https://infosys.beckhoff.com/english.php?content=../content/1033/tc3_plc_intro/2529415819.html&id=
        // Interpreted as "milliseconds"
        ['TIME' TIME
            [simple uint 32 milliseconds]
        ]
        // Interpreted as "nanoseconds"
        ['LTIME' LTIME
            [simple uint 64 nanoseconds]
        ]
        // Interpreted as "seconds since epoch"
        ['DATE' DATE
            [simple uint 32 secondsSinceEpoch]
        ]
        ['LDATE' LDATE
            [simple uint 64 nanosecondsSinceEpoch]
        ]
        // Interpreted as "milliseconds since midnight"
        ['TIME_OF_DAY' TIME_OF_DAY
            [simple uint 32 millisecondsSinceMidnight]
        ]
        ['LTIME_OF_DAY' LTIME_OF_DAY
            [simple uint 64 nanosecondsSinceMidnight]
        ]
        // Interpreted as "seconds since epoch"
        ['DATE_AND_TIME' DATE_AND_TIME
            [simple uint 32 secondsSinceEpoch]
        ]
        ['LDATE_AND_TIME' LDATE_AND_TIME
            [simple uint 64 nanosecondsSinceEpoch]
        ]
    ]
]

[enum int 8 AdsDataType(uint 16 numBytes, PlcValueType plcValueType)
    ['0x01' BOOL           ['1',   'BOOL'         ]]
    ['0x02' BIT            ['1',   'BOOL'         ]]

    // -----------------------------------------
    // Bit-strings
    // -----------------------------------------
    // 1 byte
    ['0x03' BIT8           ['1',   'BYTE'         ]]
    ['0x04' BYTE           ['1',   'BYTE'         ]]
    ['0x05' BITARR8        ['1',   'BYTE'         ]]
    // 2 byte (16 bit)
    ['0x06' WORD           ['2',   'WORD'         ]]
    ['0x07' BITARR16       ['2',   'WORD'         ]]
    // 4 byte (32 bit)
    ['0x08' DWORD          ['4',   'DWORD'        ]]
    ['0x09' BITARR32       ['4',   'DWORD'        ]]
    // -----------------------------------------
    // Integers
    // -----------------------------------------
    // 8 bit:
    ['0x0A' SINT           ['1',   'SINT'         ]]
    ['0x0B' INT8           ['1',   'SINT'         ]]
    ['0x0C' USINT          ['1',   'USINT'        ]]
    ['0x0D' UINT8          ['1',   'USINT'        ]]
    // 16 bit:
    ['0x0E' INT            ['2',   'INT'          ]]
    ['0x0F' INT16          ['2',   'INT'          ]]
    ['0x10' UINT           ['2',   'UINT'         ]]
    ['0x11' UINT16         ['2',   'UINT'         ]]
    // 32 bit:
    ['0x12' DINT           ['4',   'DINT'         ]]
    ['0x13' INT32          ['4',   'DINT'         ]]
    ['0x14' UDINT          ['4',   'UDINT'        ]]
    ['0x15' UINT32         ['4',   'UDINT'        ]]
    // 64 bit:
    ['0x16' LINT           ['8',   'LINT'         ]]
    ['0x17' INT64          ['8',   'LINT'         ]]
    ['0x18' ULINT          ['8',   'ULINT'        ]]
    ['0x19' UINT64         ['8',   'ULINT'        ]]
    // -----------------------------------------
    // Floating point values
    // -----------------------------------------
    ['0x1A' REAL           ['4',   'REAL'         ]]
    ['0x1B' FLOAT          ['4',   'REAL'         ]]
    ['0x1C' LREAL          ['8',   'LREAL'        ]]
    ['0x1D' DOUBLE         ['8',   'LREAL'        ]]
    // -----------------------------------------
    // Characters & Strings
    // -----------------------------------------
    ['0x1E' CHAR           ['1',   'CHAR'         ]]
    ['0x1F' WCHAR          ['2',   'WCHAR'        ]]
    ['0x20' STRING         ['256', 'STRING'       ]]
    ['0x21' WSTRING        ['512', 'WSTRING'      ]]
    // -----------------------------------------
    // Dates & Times
    // -----------------------------------------
    ['0x22' TIME           ['4',   'TIME'         ]]
    ['0x23' LTIME          ['8',   'LTIME'        ]]
    ['0x24' DATE           ['4',   'DATE'         ]]
    ['0x25' TIME_OF_DAY    ['4',   'TIME_OF_DAY'  ]]
    ['0x26' TOD            ['4',   'TIME_OF_DAY'  ]]
    ['0x27' DATE_AND_TIME  ['4',   'DATE_AND_TIME']]
    ['0x28' DT             ['4',   'DATE_AND_TIME']]
]

// https://github.com/Beckhoff/ADS/blob/master/AdsLib/standalone/AdsDef.h
// https://gitlab.com/xilix-systems-llc/go-native-ads/-/blob/master/ads.go#L145
// https://gitlab.com/xilix-systems-llc/go-native-ads/-/blob/master/connection.go#L109
// https://gitlab.com/xilix-systems-llc/go-native-ads/-/blob/master/symbols.go#L222
// Especially interesting for the sum add/delete notification requests
// https://infosys.beckhoff.com/english.php?content=../content/1033/tc3_ads_intro/117463563.html&id=
[enum uint 32 ReservedIndexGroups                             (bit tc2, bit tc3)
    ['0x0000F000' ADSIGRP_SYMTAB                              ['true' , 'true' ]]
    ['0x0000F001' ADSIGRP_SYMNAME                             ['true' , 'true' ]]
    ['0x0000F002' ADSIGRP_SYMVAL                              ['true' , 'true' ]]
    ['0x0000F003' ADSIGRP_SYM_HNDBYNAME                       ['true' , 'true' ]]
    ['0x0000F004' ADSIGRP_SYM_VALBYNAME                       ['true' , 'true' ]]
    ['0x0000F005' ADSIGRP_SYM_VALBYHND                        ['true' , 'true' ]]
    ['0x0000F006' ADSIGRP_SYM_RELEASEHND                      ['true' , 'true' ]]
    ['0x0000F007' ADSIGRP_SYM_INFOBYNAME                      ['true' , 'true' ]]
    ['0x0000F008' ADSIGRP_SYM_VERSION                         ['true' , 'true' ]]
    // We can use this GID to read the type information of a given variable
    // in the operation mode in which we don't read the entire structures on
    // connection start.
    ['0x0000F009' ADSIGRP_SYM_INFOBYNAMEEX                    ['true' , 'true' ]]
    ['0x0000F00A' ADSIGRP_SYM_DOWNLOAD                        ['true' , 'true' ]]
    // Read the symbol-table (All variables defined in the PLC)
    ['0x0000F00B' ADSIGRP_SYM_UPLOAD                          ['true' , 'true' ]]
    ['0x0000F00C' ADSIGRP_SYM_UPLOADINFO                      ['true' , 'true' ]]
    ['0x0000F00D' ADSIGRP_SYM_DOWNLOAD2                       ['true' , 'true' ]]
    // Read the data-type-table (All data-types defined in the PLC)
    ['0x0000F00E' ADSIGRP_DATA_TYPE_TABLE_UPLOAD              ['false', 'true' ]]
    // Read the sizes of the symbol and data-type-tables
    ['0x0000F00F' ADSIGRP_SYMBOL_AND_DATA_TYPE_SIZES          ['false', 'true' ]]
    ['0x0000F010' ADSIGRP_SYMNOTE                             ['true' , 'true' ]]
    // We can use this GID to read the data-type information for a given
    // data type name in the operation mode in which we don't read the
    // entire structures on connection start.
    ['0x0000F011' ADSIGRP_DT_INFOBYNAMEEX                     ['false', 'true' ]]
    ['0x0000F012' ADSIGRP_SYM_ADDRBYHND                       ['false', 'true' ]]
    ['0x0000F013' ADSIGRP_SYM_POINTER_SUPPORT                 ['false', 'true' ]]
    ['0x0000F014' ADSIGRP_SYM_POINTER_ACCESS                  ['false', 'true' ]]
    ['0x0000F015' ADSIGRP_SYM_REFERENCE_SUPPORT               ['false', 'true' ]]
    ['0x0000F016' ADSIGRP_SYM_REFERENCE_ACCESS                ['false', 'true' ]]
    ['0x0000F018' ADSIGRP_SYM_VALBYHND_WITHMASK               ['false', 'true' ]]
    ['0x0000F019' ADSIGRP_SYM_NOACCESS_TO_SUBSYM              ['false', 'true' ]]
    ['0x0000F01A' ADSIGRP_SYM_POINTER_BITACCESS               ['false', 'true' ]]
    ['0x0000F01B' ADSIGRP_SYM_REFERENCE_BITACCESS             ['false', 'true' ]]
    ['0x0000F01C' ADSIGRP_SYM_DOWNLOAD3                       ['false', 'true' ]]
    ['0x0000F01D' ADSIGRP_SYM_FORWARD_ACCESS                  ['false', 'true' ]]
    ['0x0000F01E' ADSIGRP_SYM_FORWARD_BYHND                   ['false', 'true' ]]
    ['0x0000F01F' ADSIGRP_SYM_XAF_OBJECTID                    ['false', 'true' ]]

    // Access to the %I fields
    ['0x0000F020' ADSIGRP_IOIMAGE_RWIB                        ['true' , 'true' ]]
    ['0x0000F021' ADSIGRP_IOIMAGE_RWIX                        ['true' , 'true' ]]
    ['0x0000F024' ADSIGRP_SYM_INFOBYHNDEX                     ['true' , 'true' ]]
    ['0x0000F025' ADSIGRP_IOIMAGE_RISIZE                      ['true' , 'true' ]]
    ['0x0000F028' ADSIGRP_IOIMAGE_RWIX0                       ['true' , 'true' ]]
    ['0x0000F029' ADSIGRP_IOIMAGE_RWIX1                       ['true' , 'true' ]]
    ['0x0000F02A' ADSIGRP_IOIMAGE_RWIX2                       ['true' , 'true' ]]
    ['0x0000F02B' ADSIGRP_IOIMAGE_RWIX3                       ['true' , 'true' ]]
    ['0x0000F02C' ADSIGRP_IOIMAGE_RWIX4                       ['true' , 'true' ]]
    ['0x0000F02D' ADSIGRP_IOIMAGE_RWIX5                       ['true' , 'true' ]]
    ['0x0000F02E' ADSIGRP_IOIMAGE_RWIX6                       ['true' , 'true' ]]
    ['0x0000F02F' ADSIGRP_IOIMAGE_RWIX7                       ['true' , 'true' ]]

    // Access to the %Q fields
    ['0x0000F030' ADSIGRP_IOIMAGE_RWOB                        ['true' , 'true' ]]
    ['0x0000F031' ADSIGRP_IOIMAGE_RWOX                        ['true' , 'true' ]]
    ['0x0000F035' ADSIGRP_IOIMAGE_RWOSIZE                     ['true' , 'true' ]]
    ['0x0000F038' ADSIGRP_IOIMAGE_RWOX0                       ['true' , 'true' ]]
    ['0x0000F039' ADSIGRP_IOIMAGE_RWOX1                       ['true' , 'true' ]]
    ['0x0000F03A' ADSIGRP_IOIMAGE_RWOX2                       ['true' , 'true' ]]
    ['0x0000F03B' ADSIGRP_IOIMAGE_RWOX3                       ['true' , 'true' ]]
    ['0x0000F03C' ADSIGRP_IOIMAGE_RWOX4                       ['true' , 'true' ]]
    ['0x0000F03D' ADSIGRP_IOIMAGE_RWOX5                       ['true' , 'true' ]]
    ['0x0000F03E' ADSIGRP_IOIMAGE_RWOX6                       ['true' , 'true' ]]
    ['0x0000F03F' ADSIGRP_IOIMAGE_RWOX7                       ['true' , 'true' ]]

    ['0x0000F040' ADSIGRP_IOIMAGE_CLEARI                      ['true' , 'true' ]]

    ['0x0000F050' ADSIGRP_IOIMAGE_CLEARO                      ['true' , 'true' ]]

    ['0x0000F060' ADSIGRP_IOIMAGE_RWIOB                       ['true' , 'true' ]]
    ['0x0000F064' ADSIGRP_IOIMAGE_WATCHDOG                    ['true' , 'true' ]]
    ['0x0000F068' ADSIGRP_IOIMAGE_CREATE                      ['true' , 'true' ]]

    ['0x0000F070' ADSIGRP_JSON                                ['true' , 'true' ]]
    ['0x0000F071' ADSIGRP_JSON_DOM_ACCESS_GET                 ['true' , 'true' ]]
    ['0x0000F072' ADSIGRP_JSON_DOM_ACCESS_SET                 ['true' , 'true' ]]
    ['0x0000F073' ADSIGRP_JSON_DOM_ACCESS_DEL                 ['true' , 'true' ]]
    ['0x0000F074' ADSIGRP_JSON_DOM_ACCESS_LEN                 ['true' , 'true' ]]
    ['0x0000F078' ADSIGRP_ADSWATCH_GETOID                     ['true' , 'true' ]]

    // Sum Requests
    ['0x0000F080' ADSIGRP_MULTIPLE_READ                       ['false', 'true' ]]
    ['0x0000F081' ADSIGRP_MULTIPLE_WRITE                      ['false', 'true' ]]
    ['0x0000F082' ADSIGRP_MULTIPLE_READ_WRITE                 ['false', 'true' ]]
    ['0x0000F083' ADSIGRP_MULTIPLE_RELEASE_HANDLE             ['false', 'true' ]]
    ['0x0000F084' ADSIGRP_SUMUP_READEX2                       ['false', 'true' ]]
    ['0x0000F085' ADSIGRP_MULTIPLE_ADD_DEVICE_NOTIFICATIONS   ['false', 'true' ]]
    ['0x0000F086' ADSIGRP_MULTIPLE_DELETE_DEVICE_NOTIFICATIONS['false', 'true' ]]
    ['0x0000F088' ADSIGRP_EXTERNALTIME                        ['false', 'true' ]]

    ['0x0000F090' ADSIGRP_CHECK_NOTIFICATION                  ['false', 'true' ]]
    ['0x0000F091' ADSIGRP_DIAG_NOTIFICATION                   ['false', 'true' ]]
    ['0x0000F098' ADSIGRP_LOGGING_CONFIG                      ['false', 'true' ]]

    ['0x0000F0A0' ADSIGRP_SYM_CONTEXTCYCLE                    ['false', 'true' ]]

    ['0x0000F0B0' ADSIGRP_DEVICE_CONTEXT_DATA                 ['false', 'true' ]]

    ['0x0000F100' ADSIGRP_DEVICE_DATA                         ['true' , 'true' ]]

    ['0x0000F200' ADSIGRP_TASK_DATA                           ['true' , 'true' ]]

    ['0x0000F300' ADSIGRP_CANOPEN_BEGIN                       ['true' , 'true' ]]
    // This Group Index makes ADS access data via AoE (ADS over EtherCAT) and
    // can be used to access telemetry data on the DeviceManager or from attached
    // EtherCAT devices.
    // https://infosys.beckhoff.com/index.php?content=../content/1031/eap/1521731467.html
    //['0x0000F302' ADS_OVER_ETHERCAT                           ['true' , 'true' ]]
    ['0x0000F302' ADSIGRP_CANOPEN_SDO                         ['true' , 'true' ]]
    ['0x0000F303' ADSIGRP_CANOPEN_SDO_LASTERROR               ['true' , 'true' ]]
    ['0x0000F304' ADSIGRP_CANOPEN_SDO_SUMUP_READ              ['true' , 'true' ]]
    ['0x0000F305' ADSIGRP_CANOPEN_SDO_SUMUP_WRITE             ['true' , 'true' ]]

    ['0x0000F3F8' ADSIGRP_CANOPEN_TXPDO_ACCESS                ['true' , 'true' ]]
    ['0x0000F3F9' ADSIGRP_CANOPEN_RXPDO_ACCESS                ['true' , 'true' ]]
    ['0x0000F3FB' ADSIGRP_CANOPEN_SDO_ENI_CONTENT             ['true' , 'true' ]]
    ['0x0000F3FC' ADSIGRP_CANOPEN_SDO_INFO_LIST               ['true' , 'true' ]]
    ['0x0000F3FD' ADSIGRP_CANOPEN_SDO_INFO_OBJ                ['true' , 'true' ]]
    ['0x0000F3FE' ADSIGRP_CANOPEN_SDO_INFO_ENTRY              ['true' , 'true' ]]
    ['0x0000F3FF' ADSIGRP_CANOPEN_END                         ['true' , 'true' ]]

    ['0x8001F302' ADSIGRP_ECAT_EMCY_SERVER                    ['true' , 'true' ]]
    ['0xF8200101' ADSIOFFS_ECAT_EMCY_SERVER_CONNECT           ['true' , 'true' ]]

    ['0x0000F400' ADSIGRP_ECAT_FOE_BEGIN                      ['true' , 'true' ]]
    ['0x0000F401' ADSIGRP_ECAT_FOE_FOPENREAD                  ['true' , 'true' ]]
    ['0x0000F402' ADSIGRP_ECAT_FOE_FOPENWRITE                 ['true' , 'true' ]]
    ['0x0000F403' ADSIGRP_ECAT_FOE_FCLOSE                     ['true' , 'true' ]]
    ['0x0000F404' ADSIGRP_ECAT_FOE_FREAD                      ['true' , 'true' ]]
    ['0x0000F405' ADSIGRP_ECAT_FOE_FWRITE                     ['true' , 'true' ]]
    ['0x0000F406' ADSIGRP_ECAT_FOE_PROGRESSINFO               ['true' , 'true' ]]
    ['0x0000F407' ADSIGRP_ECAT_FOE_LASTERROR                  ['true' , 'true' ]]

    ['0x0000F41F' ADSIGRP_ECAT_FOE_END                        ['true' , 'true' ]]

    ['0x0000F420' ADSIGRP_ECAT_SOE                            ['true' , 'true' ]]
    ['0x0000F421' ADSIGRP_ECAT_SOE_LASTERROR                  ['true' , 'true' ]]

    ['0x0000F430' ADSIGRP_ECAT_VOE                            ['true' , 'true' ]]

    ['0x00000000' ADSIOFFS_DEVDATA_ADSSTATE                   ['true' , 'true' ]]
    ['0x00000002' ADSIOFFS_DEVDATA_DEVSTATE                   ['true' , 'true' ]]
]

[enum uint 32 ReturnCode
    // Global Return Codes
    ['0x00000000' OK]
    ['0x00000001' INTERNAL_ERROR]
    ['0x00000002' NO_REALTIME]
    ['0x00000003' SAVE_ERROR]
    ['0x00000004' MAILBOX_FULL]
    ['0x00000005' WRONG_HMSG]
    ['0x00000006' TARGET_PORT_NOT_FOUND]
    ['0x00000007' TARGET_HOST_NOT_FOUND]
    ['0x00000008' UNKNOWN_COMMAND_ID]
    ['0x00000009' UNKNOWN_TASK_ID]
    ['0x0000000A' NO_IO]
    ['0x0000000B' UNKNOWN_ADS_COMMAND]
    ['0x0000000C' WIN32_ERROR]
    ['0x0000000D' PORT_NOT_CONNECTED]
    ['0x0000000E' INVALID_ADS_LENGTH]
    ['0x0000000F' INVALID_AMS_NET_ID]
    ['0x00000010' LOW_INSTALLATION_LEVEL]
    ['0x00000011' NO_DEBUGGING_AVAILABLE]
    ['0x00000012' PORT_DEACTIVATED]
    ['0x00000013' PORT_ALREADY_CONNECTED]
    ['0x00000014' ADS_SYNC_WIN32_ERROR]
    ['0x00000015' ADS_SYNC_TIMEOUT]
    ['0x00000016' ADS_SYNC_AMS_ERROR]
    ['0x00000017' NO_INDEX_MAP_FOR_ADS_AVAILABLE]
    ['0x00000018' INVALID_ADS_PORT]
    ['0x00000019' NO_MEMORY]
    ['0x0000001A' TCP_SENDING_ERROR]
    ['0x0000001B' HOST_NOT_REACHABLE]
    ['0x0000001C' INVALID_AMS_FRAGMENT]

    // Router Error-Codes
    ['0x00000500' ROUTERERR_NOLOCKEDMEMORY]
    ['0x00000501' ROUTERERR_RESIZEMEMORY]
    ['0x00000502' ROUTERERR_MAILBOXFULL]
    ['0x00000503' ROUTERERR_DEBUGBOXFULL]
    ['0x00000504' ROUTERERR_UNKNOWNPORTTYPE]
    ['0x00000505' ROUTERERR_NOTINITIALIZED]
    ['0x00000506' ROUTERERR_PORTALREADYINUSE]
    ['0x00000507' ROUTERERR_NOTREGISTERED]
    ['0x00000508' ROUTERERR_NOMOREQUEUES]
    ['0x00000509' ROUTERERR_INVALIDPORT]
    ['0x0000050A' ROUTERERR_NOTACTIVATED]

    // General ADS Error-Codes
    ['0x00000700' ADSERR_DEVICE_ERROR]
    ['0x00000701' ADSERR_DEVICE_SRVNOTSUPP]
    ['0x00000702' ADSERR_DEVICE_INVALIDGRP]
    ['0x00000703' ADSERR_DEVICE_INVALIDOFFSET]
    ['0x00000704' ADSERR_DEVICE_INVALIDACCESS]
    ['0x00000705' ADSERR_DEVICE_INVALIDSIZE]
    ['0x00000706' ADSERR_DEVICE_INVALIDDATA]
    ['0x00000707' ADSERR_DEVICE_NOTREADY]
    ['0x00000708' ADSERR_DEVICE_BUSY]
    ['0x00000709' ADSERR_DEVICE_INVALIDCONTEXT]
    ['0x0000070A' ADSERR_DEVICE_NOMEMORY]
    ['0x0000070B' ADSERR_DEVICE_INVALIDPARM]
    ['0x0000070C' ADSERR_DEVICE_NOTFOUND]
    ['0x0000070D' ADSERR_DEVICE_SYNTAX]
    ['0x0000070E' ADSERR_DEVICE_INCOMPATIBLE]
    ['0x0000070F' ADSERR_DEVICE_EXISTS]
    ['0x00000710' ADSERR_DEVICE_SYMBOLNOTFOUND]
    ['0x00000711' ADSERR_DEVICE_SYMBOLVERSIONINVALID]
    ['0x00000712' ADSERR_DEVICE_INVALIDSTATE]
    ['0x00000713' ADSERR_DEVICE_TRANSMODENOTSUPP]
    ['0x00000714' ADSERR_DEVICE_NOTIFYHNDINVALID]
    ['0x00000715' ADSERR_DEVICE_CLIENTUNKNOWN]
    ['0x00000716' ADSERR_DEVICE_NOMOREHDLS]
    ['0x00000717' ADSERR_DEVICE_INVALIDWATCHSIZE]
    ['0x00000718' ADSERR_DEVICE_NOTINIT]
    ['0x00000719' ADSERR_DEVICE_TIMEOUT]
    ['0x0000071A' ADSERR_DEVICE_NOINTERFACE]
    ['0x0000071B' ADSERR_DEVICE_INVALIDINTERFACE]
    ['0x0000071C' ADSERR_DEVICE_INVALIDCLSID]
    ['0x0000071D' ADSERR_DEVICE_INVALIDOBJID]
    ['0x0000071E' ADSERR_DEVICE_PENDING]
    ['0x0000071F' ADSERR_DEVICE_ABORTED]
    ['0x00000720' ADSERR_DEVICE_WARNING]
    ['0x00000721' ADSERR_DEVICE_INVALIDARRAYIDX]
    ['0x00000722' ADSERR_DEVICE_SYMBOLNOTACTIVE]
    ['0x00000723' ADSERR_DEVICE_ACCESSDENIED]
    ['0x00000724' ADSERR_DEVICE_LICENSENOTFOUND]
    ['0x00000725' ADSERR_DEVICE_LICENSEEXPIRED]
    ['0x00000726' ADSERR_DEVICE_LICENSEEXCEEDED]
    ['0x00000727' ADSERR_DEVICE_LICENSEINVALID]
    ['0x00000728' ADSERR_DEVICE_LICENSESYSTEMID]
    ['0x00000729' ADSERR_DEVICE_LICENSENOTIMELIMIT]
    ['0x0000072A' ADSERR_DEVICE_LICENSEFUTUREISSUE]
    ['0x0000072B' ADSERR_DEVICE_LICENSETIMETOLONG]
    ['0x0000072c' ADSERR_DEVICE_EXCEPTION]
    ['0x0000072D' ADSERR_DEVICE_LICENSEDUPLICATED]
    ['0x0000072E' ADSERR_DEVICE_SIGNATUREINVALID]
    ['0x0000072F' ADSERR_DEVICE_CERTIFICATEINVALID]
    ['0x00000740' ADSERR_CLIENT_ERROR]
    ['0x00000741' ADSERR_CLIENT_INVALIDPARM]
    ['0x00000742' ADSERR_CLIENT_LISTEMPTY]
    ['0x00000743' ADSERR_CLIENT_VARUSED]
    ['0x00000744' ADSERR_CLIENT_DUPLINVOKEID]
    ['0x00000745' ADSERR_CLIENT_SYNCTIMEOUT]
    ['0x00000746' ADSERR_CLIENT_W32ERROR]
    ['0x00000747' ADSERR_CLIENT_TIMEOUTINVALID]
    ['0x00000748' ADSERR_CLIENT_PORTNOTOPEN]
    ['0x00000750' ADSERR_CLIENT_NOAMSADDR]
    ['0x00000751' ADSERR_CLIENT_SYNCINTERNAL]
    ['0x00000752' ADSERR_CLIENT_ADDHASH]
    ['0x00000753' ADSERR_CLIENT_REMOVEHASH]
    ['0x00000754' ADSERR_CLIENT_NOMORESYM]
    ['0x00000755' ADSERR_CLIENT_SYNCRESINVALID]

    // RTime Error-Codes
    ['0x00001000' RTERR_INTERNAL]
    ['0x00001001' RTERR_BADTIMERPERIODS]
    ['0x00001002' RTERR_INVALIDTASKPTR]
    ['0x00001003' RTERR_INVALIDSTACKPTR]
    ['0x00001004' RTERR_PRIOEXISTS]
    ['0x00001005' RTERR_NOMORETCB]
    ['0x00001006' RTERR_NOMORESEMAS]
    ['0x00001007' RTERR_NOMOREQUEUES]
    ['0x0000100D' RTERR_EXTIRQALREADYDEF]
    ['0x0000100E' RTERR_EXTIRQNOTDEF]
    ['0x0000100F' RTERR_EXTIRQINSTALLFAILED]
    ['0x00001010' RTERR_IRQLNOTLESSOREQUAL]
    ['0x00001017' RTERR_VMXNOTSUPPORTED]
    ['0x00001018' RTERR_VMXDISABLED]
    ['0x00001019' RTERR_VMXCONTROLSMISSING]
    ['0x0000101A' RTERR_VMXENABLEFAILS]

    // TCP Windsock Error-Codes
    ['0x0000274C' WSAETIMEDOUT]
    ['0x0000274D' WSAECONNREFUSED]
    ['0x00002751' WSAEHOSTUNREACH]
]

[type AdsTableSizes byteOrder='"LITTLE_ENDIAN"' unsignedIntegerEncoding='"unsigned-binary"'
	[simple   uint 32 symbolCount   ]
	[simple   uint 32 symbolLength  ]
	[simple   uint 32 dataTypeCount ]
	[simple   uint 32 dataTypeLength]
	[simple   uint 32 aliasCount    ]
	[simple   uint 32 aliasLength   ]
]

[type AdsSymbolTableEntry byteOrder='"LITTLE_ENDIAN"' unsignedIntegerEncoding='"unsigned-binary"' signedIntegerEncoding='"twos-complement"' stringEncoding='"UTF8"'
  	[simple   uint 32                          entryLength                                                             ]
    [simple   uint 32                          group                                                                   ]
    [simple   uint 32                          offset                                                                  ]
    [simple   uint 32                          size                                                                    ]
    [simple   uint 32                          dataType                                                                ]
    // Start: Flags
    // https://github.com/jisotalo/ads-client/blob/master/src/ads-commons.ts#L619
    // Order of the bits if read Little-Endian and then accessing the bit flags
    // 7 6 5 4 3 2 1 0  |  15 14 13 12 11 10 9 8
    [simple   bit                              flagMethodDeref                                                         ]
    [simple   bit                              flagItfMethodAccess                                                     ]
    [simple   bit                              flagReadOnly                                                            ]
    [simple   bit                              flagTComInterfacePointer                                                ]
    [implicit bit                              flagTypeGuid               'COUNT(guid) > 0'                            ]
    [simple   bit                              flagReferenceTo                                                         ]
    [simple   bit                              flagBitValue                                                            ]
    [simple   bit                              flagPersistent                                                          ]
    [simple   bit                              flagCompilerGenerated                                                   ]
    [reserved uint 1                           '0x0'                                                                   ]
    [simple   bit                              flagSystemServiceSymbol                                                 ]
    [simple   bit                              flagExtendedFlags                                                       ]
    [simple   bit                              flagInitOnReset                                                         ]
    [simple   bit                              flagStatic                                                              ]
    [implicit bit                              flagAttributes             'attributes != null'                         ]
    [implicit bit                              flagContextMask            'contextMask != null'                        ]
    // https://github.com/jisotalo/ads-client/blob/master/src/ads-commons.ts#L679
    // Order of the bits if read Little-Endian and then accessing the bit flags
    // 7 6 5 4 3 2 1 0  |  15 14 13 12 11 10 9 8
    [reserved uint 3                          '0x0'                                                                    ]
    [simple   bit                              flagVariantType                                                         ]
    [simple   bit                              flagOnlineChangePtrRefType                                              ]
    [simple   bit                              flagRefactorInfo                                                        ]
    [simple   bit                              flagRedundancyIgnore                                                    ]
    [simple   bit                              flagPlcPointerType                                                      ]
    [reserved uint 8                           '0x00'                                                                  ]
    // End: Flags
    [implicit uint 16                          nameLength                 'STR_LEN(name)'                              ]
    [implicit uint 16                          dataTypeNameLength         'STR_LEN(dataTypeName)'                      ]
    [implicit uint 16                          commentLength              'STR_LEN(comment)'                           ]
	[simple   vstring 'nameLength * 8'         name                                                                    ]
	[const    uint 8                           nameTerminator             0x00                                         ]
	[simple   vstring 'dataTypeNameLength * 8' dataTypeName                                                            ]
	[const    uint 8                           dataTypeNameTerminator     0x00                                         ]
	[simple   vstring 'commentLength * 8'      comment                                                                 ]
	[const    uint 8                           commentTerminator          0x00                                         ]
	[optional uint 32                          contextMask                'flagContextMask'                            ]
    [array    byte                             guid                       count         'flagTypeGuid == true ? 16 : 0']
    [optional AdsDataTypeAttributes            attributes                 'flagAttributes'                             ]
    // Gobbling up the rest, but it seems there is only empty padding bytes in it.
	[array    byte                             rest                       count            'entryLength - curPos']
]

// https://gitlab.com/xilix-systems-llc/go-native-ads/-/blob/master/symbols.go#L15
[type AdsDataTypeTableEntry(uint 16 maxDepth) byteOrder='"LITTLE_ENDIAN"' unsignedIntegerEncoding='"unsigned-binary"' signedIntegerEncoding='"twos-complement"' stringEncoding='"UTF8"'
	// An entry may contain further entries, so the nesting depth is dictated by the device. Each
	// level costs only ~45 bytes of wire input, so the budget is spent on the way down and a table
	// that nests past it is rejected as a malformed table rather than exhausting the parser stack.
	[validation                                  'maxDepth > 0'           "maximum data type table nesting depth exceeded"  ]
	[simple   uint 32                            entryLength                                                           ]
	[simple   uint 32                            version                                                               ]
	[simple   uint 32                            hashValue                                                             ]
	[simple   uint 32                            typeHashValue                                                         ]
	[simple   uint 32                            size                                                                  ]
	[simple   uint 32                            offset                                                                ]
	[simple   AdsDatatypeId                      dataType                                                              ]
	// Begin: Data Type Flags
	// Source (https://github.com/jisotalo/ads-client/blob/master/src/ads-commons.ts#L724)
    // 7 6 5 4 3 2 1 0  |  15 14 13 12 11 10 9 8  |  23 22 21 20 19 18 17 16 | 31 30 29 28 27 26 25 24
	// Byte 1
	[implicit bit                                flagTypeGuid               'COUNT(guid) > 0'                          ]
	[simple   bit                                flagPropItem                                                          ]
	[simple   bit                                flagBitValues                                                         ]
	[simple   bit                                flagOversample                                                        ]
	[simple   bit                                flagMethodRef                                                         ]
	[simple   bit                                flagReferenceTo                                                       ]
	[implicit bit                                flagDataType               'numChildren > 0'                         ]
	[implicit bit                                flagDataItem               'numChildren == 0'                          ]
	// Byte 2
    [reserved uint 2                             '0x0'                                                                 ]
  	[implicit bit                                flagExtendedInfos          'extendedInfos != null'                    ]
  	[implicit bit                                flagAttributes             'attributes != null'                       ]
  	[implicit bit                                flagMethodInfos            'methodInfos != null'                      ]
  	[simple   bit                                flagTComInterfacePtr                                                  ]
  	[simple   bit                                flagCopyMask                                                          ]
  	[simple   bit                                flagPersistent                                                        ]
	// Byte 3
    [simple   bit                                flagPlcPointerType                                                    ]
    [simple   bit                                flagInitOnReset                                                       ]
    [simple   bit                                flagPersistentDataType                                                ]
    [simple   bit                                flagAnySizeArray                                                      ]
    [simple   bit                                flagIgnorePersist                                                     ]
    [simple   bit                                flagSoftwareProtectionLevels                                          ]
    [simple   bit                                flagStatic                                                            ]
    [simple   bit                                flagAligned                                                           ]
	// Byte 4
    [simple   bit                                ExtendedFlags                                                         ]
    [reserved uint 1                             '0x0'                                                                 ]
    [simple   bit                                flagExtendedEnumInfos                                                 ]
    [simple   bit                                flagDeRefTypeItem                                                     ]
    [simple   bit                                flagContainsOnlineChangePtrRef                                        ]
    [simple   bit                                flagIncomplete                                                        ]
    [simple   bit                                flagHideSubItems                                                      ]
    [simple   bit                                flagRefactorInfo                                                      ]
    // End: Data Type Flags
	[implicit uint 16                            mainNameLength           'STR_LEN(mainName)'                          ]
	[implicit uint 16                            secondaryNameLength      'STR_LEN(secondaryName)'                     ]
	[implicit uint 16                            commentLength            'STR_LEN(comment)'                           ]
	[simple   uint 16                            arrayDimensions                                                       ]
	[simple   uint 16                            numChildren                                                           ]
	[simple   vstring 'mainNameLength * 8'       mainName                                                              ]
	[const    uint 8                             mainNameTerminator       0x00                                         ]
	[simple   vstring 'secondaryNameLength * 8'  secondaryName                                                         ]
	[const    uint 8                             secondaryNameTerminator  0x00                                         ]
	[simple   vstring 'commentLength * 8'        comment                                                               ]
	[const    uint 8                             commentTerminator        0x00                                         ]
    [array    AdsDataTypeArrayInfo               arrayInfo                count                     'arrayDimensions'  ]
   	[array    AdsDataTypeTableEntry('maxDepth - 1') children               count                     'numChildren'      ]
    [array    byte                               guid                     count         'flagTypeGuid == true ? 16 : 0']
	[optional AdsMethodInfos                     methodInfos              'flagMethodInfos'                            ]
    [optional AdsDataTypeAttributes              attributes               'flagAttributes'                             ]
    [optional AdsExtendedInfos('dataType')       extendedInfos            'flagExtendedInfos'                          ]
	// This only consumes the rest in cased of both methodInfos and extendedInfos not being set.
	[array    byte                               rest                     count            'entryLength - curPos']
]

[type AdsDataTypeArrayInfo byteOrder='"LITTLE_ENDIAN"'
    [simple  uint 32 lowerBound                                                                                        ]
    [simple  uint 32 numElements                                                                                       ]
    [virtual uint 32 upperBound  'lowerBound + (numElements - 1)'                                                      ]
]

[type AdsMethodInfos
    [implicit uint 16       numMethodInfos 'COUNT(methodInfos)'                                                        ]
    [array    AdsMethodInfo methodInfos    count                'numMethodInfos'                                       ]
]

[type AdsMethodInfo
    [implicit uint 32                      methodInfoLength   'lengthInBytes'                                          ]
    [simple   uint 32                      header1                                                                     ]
    [simple   uint 32                      header2                                                                     ]
    [simple   uint 32                      header3                                                                     ]
    [simple   uint 32                      header4                                                                     ]
    [simple   uint 32                      header5                                                                     ]
    [array    byte                         guid               count               '16'                                 ]
    // Potentially optional
    [simple   uint 32                      methodId                                                                    ]
    // Potentially optional
    [simple   uint 32                      methodFlags                                                                 ]
    // If the two above are optional, we need a reserved 0x0000 instead.
    [implicit uint 16                      nameLength         'STR_LEN(name)'                                          ]
    [implicit uint 16                      typeNameLength     'STR_LEN(typeName)'                                      ]
    [implicit uint 16                      commentLength      'STR_LEN(comment)'                                       ]
    [implicit uint 16                      parameterCount     'COUNT(parameters)'                                      ]
	[simple   vstring 'nameLength * 8'     name                                                                        ]
	[const    uint 8                       nameTerminator     0x00                                                     ]
	[simple   vstring 'typeNameLength * 8' typeName                                                                    ]
	[const    uint 8                       typeNameTerminator 0x00                                                     ]
	[simple   vstring 'commentLength * 8'  comment                                                                     ]
	[const    uint 8                       commentTerminator  0x00                                                     ]
    [array    AdsMethodParam               parameters         count         'parameterCount'                           ]
	// This only consumes the rest in cases, where we haven't quite figgured out the full structure of this type.
	[array    byte                         rest               count         'methodInfoLength - curPos'          ]
]

// A single parameter descriptor. Matches your “tail”: length + 5 headers + GUID +
// optional ids/flags + lengths + strings + terminators.
[type AdsMethodParam
    [implicit uint 32  paramLength          'lengthInBytes'                                                            ]
    [simple   uint 32  header1                                                                                         ]
    [simple   uint 32  header2                                                                                         ]
    [simple   uint 32  header3                                                                                         ]
    [simple   uint 32  header4                                                                                         ]
    [simple   uint 32  header5                                                                                         ]
    [array    byte     guid                 count               '16'                                                   ]
    [reserved uint 16  '0x0000'                                                                                        ]
    [implicit uint 16  nameLength           'STR_LEN(name)'                                                            ]
    [implicit uint 16  typeNameLength       'STR_LEN(typeName)'                                                        ]
    [implicit uint 16  commentLength        'STR_LEN(comment)'                                                         ]
    [simple   vstring  'nameLength * 8'     name                                                                       ]
    [const    uint 8   nameTerminator       0x00                                                                       ]
    [simple   vstring  'typeNameLength * 8' typeName                                                                   ]
    [const    uint 8   typeNameTerminator   0x00                                                                       ]
    [simple   vstring  'commentLength * 8'  comment                                                                    ]
    [const    uint 8   commentTerminator    0x00                                                                       ]
]

[type AdsString(uint 16 stringLength)
	[manual vstring stringValue 'STATIC_CALL("parseZeroTerminatedString", readBuffer, stringLength)' 'STATIC_CALL("serializeZeroTerminatedString", writeBuffer, stringValue)' 'STATIC_CALL("lengthZeroTerminatedString", stringValue)'                                                                   ]
]

[type AdsDataTypeAttributes
    [implicit uint 16                        numAttributes   'COUNT(attributes)'                                       ]
    [array    AdsAttributeEntry              attributes      count               'numAttributes'                       ]
]

[type AdsAttributeEntry
    [implicit      uint 8                    nameLength      'STR_LEN(name)'                                           ]
    [implicit      uint 8                    valueLength     'STR_LEN(value)'                                          ]
    [simple        vstring 'nameLength * 8'  name                                                                      ]
    [const         uint 8                    nameTerminator  0x00                                                      ]
    [simple        vstring 'valueLength * 8' value                                                                     ]
    // This is an evil hack allowing us to parse dataypte tables containg TcLinkTo elements.
    [implicit      uint 8                    valueTerminator '0x00'                                                    ]
    [validation    'valueTerminator == 0x00' 'unexpected char instead of terminator' shouldFail=false                  ]
]

[type AdsExtendedInfos(AdsDatatypeId dataType)
    [implicit uint 16                          count           'COUNT(entries)']
    [array    AdsExtendedInfoEntry('dataType') entries         count 'count'    ]
]

[discriminatedType AdsExtendedInfoEntry(AdsDatatypeId dataType)
    [implicit uint 8                   nameLength      'STR_LEN(name)'                                                 ]
    [simple   vstring 'nameLength * 8' name                                                                            ]
    [const    uint 8                   nameTerminator  0x00                                                            ]
    [typeSwitch dataType
        ['ADST_VOID' *Void
        ]
        ['ADST_INT8' *Int8
            [simple   int 8                  value]
        ]
        ['ADST_UINT8' *Uint8
            [simple   uint 8                 value]
        ]
        ['ADST_INT16' *Int16
            [simple   int 16                 value]
        ]
        ['ADST_UINT16' *Uint16
            [simple   uint 16                value]
        ]
        ['ADST_INT32' *Int32
            [simple   int 32                  value]
        ]
        ['ADST_UINT32' *Uint32
            [simple   uint 32                 value]
        ]
        ['ADST_INT64' *Int64
            [simple   int 64                  value]
        ]
        ['ADST_UINT64' *Uint64
            [simple   uint 64                 value]
        ]
        ['ADST_REAL32' *Real32
            [simple   float 32                value]
        ]
        ['ADST_REAL64' *Real64
            [simple   float 64                value]
        ]
        ['ADST_BIGTYPE' *BigType
            // TODO: Find out
        ]
        ['ADST_STRING' *String
            // TODO: Find out
        ]
        ['ADST_WSTRING' *Wstring
            // TODO: Find out
        ]
        ['ADST_REAL80' *Real80
            // TODO: Find out
        ]
        ['ADST_BIT' *Bit
            // TODO: Find out
        ]
        ['ADST_MAXTYPES' *MaxTypes
            // TODO: Find out
        ]
    ]
]

[enum uint 32 AdsDatatypeId
    ['0'    ADST_VOID       ]
    ['16'   ADST_INT8       ]
    ['17'   ADST_UINT8      ]
    ['2'    ADST_INT16      ]
    ['18'   ADST_UINT16     ]
    ['3'    ADST_INT32      ]
    ['19'   ADST_UINT32     ]
    ['20'   ADST_INT64      ]
    ['21'   ADST_UINT64     ]
    ['4'    ADST_REAL32     ]
    ['5'    ADST_REAL64     ]
    ['65'   ADST_BIGTYPE    ]
    ['30'   ADST_STRING     ]
    ['31'   ADST_WSTRING    ]
    ['32'   ADST_REAL80     ]
    ['33'   ADST_BIT        ]
    ['34'   ADST_MAXTYPES   ]
]

// From: https://infosys.beckhoff.com/english.php?content=../content/1033/tcplclib_tc2_system/31064331.html&id=
[enum uint 16 DefaultAmsPorts
    ['900' CAM_CONTROLLER     ]
    ['851' RUNTIME_SYSTEM_01  ]
    ['852' RUNTIME_SYSTEM_02  ]
    ['853' RUNTIME_SYSTEM_03  ]
    ['854' RUNTIME_SYSTEM_04  ]
    ['855' RUNTIME_SYSTEM_05  ]
    ['856' RUNTIME_SYSTEM_06  ]
    ['857' RUNTIME_SYSTEM_07  ]
    ['858' RUNTIME_SYSTEM_08  ]
    ['859' RUNTIME_SYSTEM_09  ]
    ['860' RUNTIME_SYSTEM_10  ]
    ['861' RUNTIME_SYSTEM_11  ]
    ['862' RUNTIME_SYSTEM_12  ]
    ['863' RUNTIME_SYSTEM_13  ]
    ['864' RUNTIME_SYSTEM_14  ]
    ['865' RUNTIME_SYSTEM_15  ]
    ['866' RUNTIME_SYSTEM_16  ]
    ['867' RUNTIME_SYSTEM_17  ]
    ['868' RUNTIME_SYSTEM_18  ]
    ['869' RUNTIME_SYSTEM_19  ]
    ['870' RUNTIME_SYSTEM_20  ]
    ['871' RUNTIME_SYSTEM_21  ]
    ['872' RUNTIME_SYSTEM_22  ]
    ['873' RUNTIME_SYSTEM_23  ]
    ['874' RUNTIME_SYSTEM_24  ]
    ['875' RUNTIME_SYSTEM_25  ]
    ['876' RUNTIME_SYSTEM_26  ]
    ['877' RUNTIME_SYSTEM_27  ]
    ['878' RUNTIME_SYSTEM_28  ]
    ['879' RUNTIME_SYSTEM_29  ]
    ['880' RUNTIME_SYSTEM_30  ]
    ['881' RUNTIME_SYSTEM_31  ]
    ['882' RUNTIME_SYSTEM_32  ]
    ['883' RUNTIME_SYSTEM_33  ]
    ['884' RUNTIME_SYSTEM_34  ]
    ['885' RUNTIME_SYSTEM_35  ]
    ['886' RUNTIME_SYSTEM_36  ]
    ['887' RUNTIME_SYSTEM_37  ]
    ['888' RUNTIME_SYSTEM_38  ]
    ['889' RUNTIME_SYSTEM_39  ]
    ['890' RUNTIME_SYSTEM_40  ]
    ['891' RUNTIME_SYSTEM_41  ]
    ['892' RUNTIME_SYSTEM_42  ]
    ['893' RUNTIME_SYSTEM_43  ]
    ['894' RUNTIME_SYSTEM_44  ]
    ['895' RUNTIME_SYSTEM_45  ]
    ['896' RUNTIME_SYSTEM_46  ]
    ['897' RUNTIME_SYSTEM_47  ]
    ['898' RUNTIME_SYSTEM_48  ]
    ['899' RUNTIME_SYSTEM_49  ]
    ['500' NC                 ]
    ['400' RESERVED           ]
    ['300' IO                 ]
    ['200' REAL_TIME_CORE     ]
    ['100' EVENT_SYSTEM_LOGGER]
]

// See here: https://infosys.beckhoff.com/english.php?content=../content/1033/tc3_adsnetref/7313078411.html&id=
[enum uint 32 AdsTransMode
    ['0' NONE                ]
    ['1' CLIENT_CYCLE        ]
    ['2' CLIENT_ON_CHANGE    ]
    ['3' CYCLIC              ]
    ['4' ON_CHANGE           ]
    ['5' CYCLIC_IN_CONTEXT   ]
    ['6' ON_CHANGE_IN_CONTEXT]
]
