/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

////////////////////////////////////////////////////////////////
// AMS/TCP Packet
////////////////////////////////////////////////////////////////

[type 'AmsTCPPacket'
    // AMS/TCP Header	6 bytes	contains the tcpLength of the data packet.
    // These bytes must be set to 0.
    [reserved   uint       16       '0x0000'                            ]
    // This array contains the length of the data packet.
    // It consists of the AMS-Header and the enclosed ADS data. The unit is bytes.
    [implicit   uint       32       'length'  'userdata.lengthInBytes'  ]
    // The AMS packet to be sent.
    [simple AmsPacket    'userdata'                                     ]
]

////////////////////////////////////////////////////////////////
// AMS/Serial Packet
////////////////////////////////////////////////////////////////

// If an AMS serial frame has been received and the frame is OK (magic cookie OK, CRC OK, correct fragment number etc.),
// then the receiver has to send an acknowledge frame, to inform the transmitter that the frame has arrived.
//
// @see <a href="https://infosys.beckhoff.com/content/1033/tcadsamsserialspec/html/tcamssericalspec_amsframe.htm?id=8115637053270715044">TwinCAT AMS via RS232 Specification</a>
[type 'AmsSerialAcknowledgeFrame'
    // Id for detecting an AMS serial frame.
    [simple     uint        16  'magicCookie'        ]
    // Address of the sending participant. This value can always be set to 0 for an RS232 communication,
    // since it is a 1 to 1 connection and hence the participants are unique.
    [simple     int          8  'transmitterAddress' ]
    // Receiver’s address. This value can always be set to 0 for an RS232 communication, since it is a 1 to 1
    // connection and hence the participants are unique.
    [simple     int          8  'receiverAddress'    ]
    // Number of the frame sent. Once the number 255 has been sent, it starts again from 0. The receiver checks this
    // number with an internal counter.
    [simple     int          8  'fragmentNumber'     ]
    // The max. length of the AMS packet to be sent is 255. If larger AMS packets are to be sent then they have to be
    // fragmented (not published at the moment).
    [simple     int          8  'length'             ]
    [simple     uint        16  'crc'                ]
]

// An AMS packet can be transferred via RS232 with the help of an AMS serial frame.
// The actual AMS packet is in the user data field of the frame.
// The max. length of the AMS packet is limited to 255 bytes.
// Therefore the max. size of an AMS serial frame is 263 bytes.
// The fragment number is compared with an internal counter by the receiver.
// The frame number is simply accepted and not checked when receiving the first AMS frame or in case a timeout is
// exceeded. The CRC16 algorithm is used for calculating the checksum.
// @see <a href="https://infosys.beckhoff.com/content/1033/tcadsamsserialspec/html/tcamssericalspec_amsframe.htm?id=8115637053270715044">TwinCAT AMS via RS232 Specification</a>
[type 'AmsSerialFrame'
    // Id for detecting an AMS serial frame.
    [simple     uint        16  'magicCookie'        ]
    // Address of the sending participant. This value can always be set to 0 for an RS232 communication,
    // since it is a 1 to 1 connection and hence the participants are unique.
    [simple     int          8  'transmitterAddress' ]
    // Receiver’s address. This value can always be set to 0 for an RS232 communication, since it is a 1 to 1
    // connection and hence the participants are unique.
    [simple     int          8  'receiverAddress'    ]
    // Number of the frame sent. Once the number 255 has been sent, it starts again from 0. The receiver checks this
    // number with an internal counter.
    [simple     int          8  'fragmentNumber'     ]
    // The max. length of the AMS packet to be sent is 255. If larger AMS packets are to be sent then they have to be
    // fragmented (not published at the moment).
    [simple     int          8  'length'             ]
    // The AMS packet to be sent.
    [simple AmsPacket           'userdata'           ]
    [simple     uint        16  'crc'                ]
]

// In case the transmitter does not receive a valid acknowledgement after multiple transmission, then a reset frame is
// sent. In this way the receiver is informed that a new communication is running and the receiver then accepts the
// fragment number during the next AMS-Frame, without carrying out a check.
[type 'AmsSerialResetFrame'
    // Id for detecting an AMS serial frame.
    [simple     uint        16  'magicCookie'        ]
    // Address of the sending participant. This value can always be set to 0 for an RS232 communication,
    // since it is a 1 to 1 connection and hence the participants are unique.
    [simple     int          8  'transmitterAddress' ]
    // Receiver’s address. This value can always be set to 0 for an RS232 communication, since it is a 1 to 1
    // connection and hence the participants are unique.
    [simple     int          8  'receiverAddress'    ]
    // Number of the frame sent. Once the number 255 has been sent, it starts again from 0. The receiver checks this
    // number with an internal counter.
    [simple     int          8  'fragmentNumber'     ]
    // The max. length of the AMS packet to be sent is 255. If larger AMS packets are to be sent then they have to be
    // fragmented (not published at the moment).
    [simple     int          8  'length'     ]
    [simple     uint        16  'crc'                ]
]

////////////////////////////////////////////////////////////////
// AMS Common
////////////////////////////////////////////////////////////////

[type 'AmsPacket'
    // AMS Header	32 bytes	The AMS/TCP-Header contains the addresses of the transmitter and receiver. In addition the AMS error code , the ADS command Id and some other information.
    // This is the AmsNetId of the station, for which the packet is intended. Remarks see below.
    [simple     AmsNetId        'targetAmsNetId'                            ]
    // This is the AmsPort of the station, for which the packet is intended.
    [simple     uint        16  'targetAmsPort'                             ]
    // This contains the AmsNetId of the station, from which the packet was sent.
    [simple     AmsNetId        'sourceAmsNetId'                            ]
    // This contains the AmsPort of the station, from which the packet was sent.
    [simple     uint        16  'sourceAmsPort'                             ]
    // 2 bytes.
    [simple     CommandId       'commandId'                                 ]
    // 2 bytes.
    [simple     State           'state'                                     ]
    // 4 bytes	Size of the data range. The unit is byte.
    [implicit   uint        32  'length'   'data.lengthInBytes'             ]
    // 4 bytes	AMS error number. See ADS Return Codes.
    [simple     uint        32  'errorCode'                                 ]
    // free usable field of 4 bytes
    // 4 bytes	Free usable 32 bit array. Usually this array serves to send an Id. This Id makes is possible to assign a received response to a request, which was sent before.
    [simple      uint        32  'invokeId'                                 ]
    // The payload
    [simple     AdsData    'data'   ['commandId', 'state.response']         ]
]

[enum uint 16 'CommandId'
    ['0x0000' INVALID]
    ['0x0001' ADS_READ_DEVICE_INFO]
    ['0x0002' ADS_READ]
    ['0x0003' ADS_WRITE]
    ['0x0004' ADS_READ_STATE]
    ['0x0005' ADS_WRITE_CONTROL]
    ['0x0006' ADS_ADD_DEVICE_NOTIFICATION]
    ['0x0007' ADS_DELETE_DEVICE_NOTIFICATION]
    ['0x0008' ADS_DEVICE_NOTIFICATION]
    ['0x0009' ADS_READ_WRITE]
]

[type 'State'
    [simple     bit 'initCommand'           ]
    [simple     bit 'updCommand'            ]
    [simple     bit 'timestampAdded'        ]
    [simple     bit 'highPriorityCommand'   ]
    [simple     bit 'systemCommand'         ]
    [simple     bit 'adsCommand'            ]
    [simple     bit 'noReturn'              ]
    [simple     bit 'response'              ]
    [simple     bit 'broadcast'             ]
    [reserved   int 7 '0x0'                 ]
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
[type 'AmsNetId'
    [simple     uint        8   'octet1'            ]
    [simple     uint        8   'octet2'            ]
    [simple     uint        8   'octet3'            ]
    [simple     uint        8   'octet4'            ]
    [simple     uint        8   'octet5'            ]
    [simple     uint        8   'octet6'            ]
]

[discriminatedType 'AdsData' [CommandId 'commandId', bit 'response']
    [typeSwitch 'commandId', 'response'
        ['INVALID', 'false' AdsInvalidRequest]
        ['INVALID', 'true' AdsInvalidResponse]

        ['ADS_READ_DEVICE_INFO', 'false' AdsReadDeviceInfoRequest]
        ['ADS_READ_DEVICE_INFO', 'true' AdsReadDeviceInfoResponse
            // 4 bytes	ADS error number.
            [simple ReturnCode 'result']
            // Version	1 byte	Major version number
            [simple uint 8  'majorVersion']
            // Version	1 byte	Minor version number
            [simple uint 8  'minorVersion']
            // Build	2 bytes	Build number
            [simple uint 16  'version']
            // Name	16 bytes	Name of ADS device
            [array int 8  'device' count '16']
        ]

        ['ADS_READ', 'false' AdsReadRequest
            // 4 bytes	Index Group of the data which should be read.
            [simple uint 32 'indexGroup']
            // 4 bytes	Index Offset of the data which should be read.
            [simple uint 32 'indexOffset']
            // 4 bytes	Length of the data (in bytes) which should be read.
            [simple uint 32 'length']
        ]
        ['ADS_READ', 'true' AdsReadResponse
            // 4 bytes	ADS error number
            [simple ReturnCode 'result']
            // 4 bytes	Length of data which are supplied back.
            [implicit uint 32 'length' 'COUNT(data)']
            // n bytes	Data which are supplied back.
            [array int 8 'data' count 'length']
        ]

        ['ADS_WRITE', 'false' AdsWriteRequest
            // 4 bytes	Index Group of the data which should be written.
            [simple uint 32 'indexGroup']
            // 4 bytes	Index Offset of the data which should be written.
            [simple uint 32 'indexOffset']
            // 4 bytes	Length of the data (in bytes) which should be written.
            [implicit uint 32 'length' 'COUNT(data)']
            // n bytes	Data which are written in the ADS device.
            [array int 8 'data' count 'length']
        ]
        ['ADS_WRITE', 'true' AdsWriteResponse
            // 4 bytes	ADS error number
            [simple ReturnCode 'result']
        ]

        ['ADS_READ_STATE', 'false' AdsReadStateRequest]
        ['ADS_READ_STATE', 'true' AdsReadStateResponse
            // 4 bytes	ADS error number
            [simple ReturnCode 'result']
            // 2 bytes	New ADS status (see data type ADSSTATE of the ADS-DLL).
            [simple uint 16 'adsState']
            // 2 bytes	New device status.
            [simple uint 16 'deviceState']
        ]

        ['ADS_WRITE_CONTROL', 'false' AdsWriteControlRequest
            // 2 bytes	New ADS status (see data type ADSSTATE of the ADS-DLL).
            [simple uint 16 'adsState']
            // 2 bytes	New device status.
            [simple uint 16 'deviceState']
            // 4 bytes	Length of data in byte.
            [implicit uint 32 'length' 'COUNT(data)']
            // n bytes	Additional data which are sent to the ADS device
            [array int 8 'data' count 'length']
        ]
        ['ADS_WRITE_CONTROL', 'true' AdsWriteControlResponse
            // 4 bytes	ADS error number
            [simple ReturnCode 'result']
        ]

        ['ADS_ADD_DEVICE_NOTIFICATION', 'false' AdsAddDeviceNotificationRequest
            // 4 bytes	Index Group of the data, which should be sent per notification.
            [simple uint 32 'indexGroup']
            // 4 bytes	Index Offset of the data, which should be sent per notification.
            [simple uint 32 'indexOffset']
            // 4 bytes	Index Offset of the data, which should be sent per notification.
            // 4 bytes	Length of data in bytes, which should be sent per notification.
            [simple uint 32 'length']
            // 4 bytes	See description of the structure ADSTRANSMODE at the ADS-DLL.
            [simple uint 32 'transmissionMode']
            // 4 bytes	At the latest after this time, the ADS Device Notification is called. The unit is 1ms.
            [simple uint 32 'maxDelay']
            // 4 bytes	The ADS server checks if the value changes in this time slice. The unit is 1ms
            [simple uint 32 'cycleTime']
            // 16bytes	Must be set to 0
            [reserved   uint       128       '0x0000' ]
        ]
        ['ADS_ADD_DEVICE_NOTIFICATION', 'true' AdsAddDeviceNotificationResponse
            // 4 bytes	ADS error number
            [simple ReturnCode 'result']
            // 4 bytes	Handle of notification
            [simple uint 32 'notificationHandle']
        ]

        ['ADS_DELETE_DEVICE_NOTIFICATION', 'false' AdsDeleteDeviceNotificationRequest
            // 4 bytes	Handle of notification
            [simple uint 32 'notificationHandle']
        ]
        ['ADS_DELETE_DEVICE_NOTIFICATION', 'true' AdsDeleteDeviceNotificationResponse
            // 4 bytes	ADS error number
            [simple ReturnCode 'result']
        ]

        ['ADS_DEVICE_NOTIFICATION', 'false' AdsDeviceNotificationRequest
            // 4 bytes	Size of data in byte.
            [simple uint 32 'length']
            // 4 bytes	Number of elements of type AdsStampHeader.
            [simple uint 32 'stamps']
            // n bytes	Array with elements of type AdsStampHeader.
            [array AdsStampHeader 'adsStampHeaders' count 'stamps']
        ]
        ['ADS_DEVICE_NOTIFICATION', 'true' AdsDeviceNotificationResponse]

        ['ADS_READ_WRITE', 'false' AdsReadWriteRequest
            // 4 bytes	Index Group of the data which should be written.
            [simple uint 32 'indexGroup']
            // 4 bytes	Index Offset of the data which should be written.
            [simple uint 32 'indexOffset']
            // 4 bytes	Length of data in bytes, which should be read.
            [simple uint 32 'readLength']
            // 4 bytes	Length of the data (in bytes) which should be written. (if it's ADSIGRP_MULTIPLE_READ_WRITE, this is 16 otherwise 12)
            [implicit uint 32 'writeLength' '(COUNT(items) * ((indexGroup == 61570) ? 16 : 12)) + COUNT(data)']
            // Only if the indexGroup implies a sum-read response, will the indexOffset indicate the number of elements. (ADSIGRP_MULTIPLE_READ, ADSIGRP_MULTIPLE_WRITE, ADSIGRP_MULTIPLE_READ_WRITE)
            [array  AdsMultiRequestItem 'items' count '((indexGroup == 61568) || (indexGroup == 61569) || (indexGroup == 61570)) ? indexOffset : 0' ['indexGroup']]
            // n bytes	Data which are written in the ADS device.
            [array int 8 'data' count 'writeLength - (COUNT(items) * 12)']
        ]
        ['ADS_READ_WRITE', 'true' AdsReadWriteResponse
            // 4 bytes	ADS error number
            [simple ReturnCode 'result']
            // 4 bytes	Length of data in byte.
            [implicit uint 32 'length'  'COUNT(data)']
            // n bytes Additional data which are sent to the ADS device
            [array int 8 'data' count 'length']
        ]
    ]
]

[discriminatedType 'AdsMultiRequestItem' [uint 32 'indexGroup']
    [typeSwitch 'indexGroup'
        // ReservedIndexGroups.ADSIGRP_MULTIPLE_READ
        ['61568' AdsMultiRequestItemRead
            // 4 bytes	Index Group of the data which should be written.
            [simple uint 32 'itemIndexGroup']
            // 4 bytes	Index Offset of the data which should be written.
            [simple uint 32 'itemIndexOffset']
            // 4 bytes	Length of data in bytes, which should be read.
            [simple uint 32 'itemReadLength']
        ]
        // ReservedIndexGroups.ADSIGRP_MULTIPLE_WRITE
        ['61569' AdsMultiRequestItemWrite
            // 4 bytes	Index Group of the data which should be written.
            [simple uint 32 'itemIndexGroup']
            // 4 bytes	Index Offset of the data which should be written.
            [simple uint 32 'itemIndexOffset']
            // 4 bytes	Length of the data (in bytes) which should be written.
            [simple uint 32 'itemWriteLength']
        ]
        // ReservedIndexGroups.ADSIGRP_MULTIPLE_READ_WRITE
        ['61570' AdsMultiRequestItemReadWrite
            // 4 bytes	Index Group of the data which should be written.
            [simple uint 32 'itemIndexGroup']
            // 4 bytes	Index Offset of the data which should be written.
            [simple uint 32 'itemIndexOffset']
            // 4 bytes	Length of data in bytes, which should be read.
            [simple uint 32 'itemReadLength']
            // 4 bytes	Length of the data (in bytes) which should be written.
            [simple uint 32 'itemWriteLength']
        ]
    ]
]

[type 'AdsStampHeader'
    // 8 bytes	The timestamp is coded after the Windows FILETIME format. I.e. the value contains the number of the nano seconds, which passed since 1.1.1601. In addition, the local time change is not considered. Thus the time stamp is present as universal Coordinated time (UTC).
    [simple uint 64 'timestamp']
    // 4 bytes	Number of elements of type AdsNotificationSample.
    [simple uint 32 'samples']
    // n bytes	Array with elements of type AdsNotificationSample.
    [array AdsNotificationSample 'adsNotificationSamples' count 'samples']
]

[type 'AdsNotificationSample'
    // 4 bytes	Handle of notification
    [simple uint 32 'notificationHandle']
    // 4 Bytes	Size of data range in bytes.
    [simple uint 32 'sampleSize']
    // n Bytes	Data
    [array int 8 'data' count 'sampleSize']
]

[dataIo 'DataItem' [string '-1' 'dataFormatName', int 32 'stringLength']
    [typeSwitch 'dataFormatName'
        // -----------------------------------------
        // Bit
        // -----------------------------------------
        ['IEC61131_BOOL' BOOL
            [reserved uint 7 '0x00']
            [simple   bit    'value']
        ]

        // -----------------------------------------
        // Bit-strings
        // -----------------------------------------
        // 1 byte
        ['IEC61131_BYTE' BitString
            [simple uint 8 'value']
        ]
        // 2 byte (16 bit)
        ['IEC61131_WORD' BitString
            [simple uint 16 'value']
        ]
        // 4 byte (32 bit)
        ['IEC61131_DWORD' BitString
            [simple uint 32 'value']
        ]

        // -----------------------------------------
        // Integers
        // -----------------------------------------
        // 8 bit:
        ['IEC61131_SINT' SINT
            [simple int 8 'value']
        ]
        ['IEC61131_USINT' USINT
            [simple uint 8 'value']
        ]
        // 16 bit:
        ['IEC61131_INT' INT
            [simple int 16 'value']
        ]
        ['IEC61131_UINT' UINT
            [simple uint 16 'value']
        ]
        // 32 bit:
        ['IEC61131_DINT' DINT
            [simple int 32 'value']
        ]
        ['IEC61131_UDINT' UDINT
            [simple uint 32 'value']
        ]
        // 64 bit:
        ['IEC61131_LINT' LINT
            [simple int 64 'value']
        ]
        ['IEC61131_ULINT' ULINT
            [simple uint 64 'value']
        ]

        // -----------------------------------------
        // Floating point values
        // -----------------------------------------
        ['IEC61131_REAL' REAL
            [simple float 8.23  'value']
        ]
        ['IEC61131_LREAL' LREAL
            [simple float 11.52 'value']
        ]

        // -----------------------------------------
        // Characters & Strings
        // -----------------------------------------
        ['IEC61131_CHAR' STRING
//            [simple string '8' 'UTF-8' 'value']
        ]
        ['IEC61131_WCHAR' STRING
//            [simple string '16' 'UTF-16' 'value']
        ]
        ['IEC61131_STRING' STRING
            [manual   string '-1' 'UTF-8' 'value' 'STATIC_CALL("org.apache.plc4x.java.ads.utils.StaticHelper.parseAmsString", readBuffer, stringLength, _type.encoding)' 'STATIC_CALL("org.apache.plc4x.java.ads.utils.StaticHelper.serializeAmsString", writeBuffer, _value, stringLength, _type.encoding)' 'stringLength + 1']
        ]
        ['IEC61131_WSTRING' STRING
            [manual string '-1' 'UTF-16' 'value' 'STATIC_CALL("org.apache.plc4x.java.ads.utils.StaticHelper.parseAmsString", readBuffer, stringLength, _type.encoding)' 'STATIC_CALL("org.apache.plc4x.java.ads.utils.StaticHelper.serializeAmsString", writeBuffer, _value, stringLength, _type.encoding)' '(stringLength * 2) + 2']
        ]

        // -----------------------------------------
        // Date & Times
        // -----------------------------------------
        // Interpreted as "milliseconds"
        ['IEC61131_TIME' TIME
            [simple uint 32 'value']
        ]
        // Interpreted as "nanoseconds"
        ['IEC61131_LTIME' LTIME
            [simple uint 64 'value']
        ]
        // Interpreted as "seconds since epoch"
        ['IEC61131_DATE' DATE
            [simple uint 32 'value']
        ]
        // Interpreted as "milliseconds since midnight"
        ['IEC61131_TIME_OF_DAY' TIME_OF_DAY
            [simple uint 32 'value']
        ]
        // Interpreted as "seconds since epoch"
        ['IEC61131_DATE_AND_TIME' DATE_AND_TIME
            [simple uint 32 'secondsSinceEpoch']
        ]
    ]
]

[enum int 8 'AdsDataType' [uint 16 'numBytes', string '-1' 'dataFormatName']
    ['0x01' BOOL       ['1', 'IEC61131_BOOL']]
    ['0x02' BIT        ['1', 'IEC61131_BOOL']]
    ['0x03' BIT8       ['1', 'IEC61131_BOOL']]

    // -----------------------------------------
    // Bit-strings
    // -----------------------------------------
    // 1 byte
    ['0x04' BYTE       ['1', 'IEC61131_BYTE']]
    ['0x05' BITARR8    ['1', 'IEC61131_BYTE']]
    // 2 byte (16 bit)
    ['0x06' WORD       ['2', 'IEC61131_WORD']]
    ['0x07' BITARR16   ['2', 'IEC61131_WORD']]
    // 4 byte (32 bit)
    ['0x08' DWORD      ['4', 'IEC61131_DWORD']]
    ['0x09' BITARR32   ['4', 'IEC61131_DWORD']]
    // -----------------------------------------
    // Integers
    // -----------------------------------------
    // 8 bit:
    ['0x0A' SINT       ['1', 'IEC61131_SINT']]
    ['0x0B' INT8       ['1', 'IEC61131_SINT']]
    ['0x0C' USINT      ['1', 'IEC61131_USINT']]
    ['0x0D' UINT8      ['1', 'IEC61131_USINT']]
    // 16 bit:
    ['0x0E' INT        ['2', 'IEC61131_INT']]
    ['0x0F' INT16      ['2', 'IEC61131_INT']]
    ['0x10' UINT       ['2', 'IEC61131_UINT']]
    ['0x11' UINT16     ['2', 'IEC61131_UINT']]
    // 32 bit:
    ['0x12' DINT       ['4', 'IEC61131_DINT']]
    ['0x13' INT32      ['4', 'IEC61131_DINT']]
    ['0x14' UDINT      ['4', 'IEC61131_UDINT']]
    ['0x15' UINT32     ['4', 'IEC61131_UDINT']]
    // 64 bit:
    ['0x16' LINT       ['8', 'IEC61131_LINT']]
    ['0x17' INT64      ['8', 'IEC61131_LINT']]
    ['0x18' ULINT      ['8', 'IEC61131_ULINT']]
    ['0x19' UINT64     ['8', 'IEC61131_ULINT']]
    // -----------------------------------------
    // Floating point values
    // -----------------------------------------
    ['0x1A' REAL       ['4', 'IEC61131_REAL']]
    ['0x1B' FLOAT      ['4', 'IEC61131_REAL']]
    ['0x1C' LREAL      ['8', 'IEC61131_LREAL']]
    ['0x1D' DOUBLE     ['8', 'IEC61131_LREAL']]
    // -----------------------------------------
    // Characters & Strings
    // -----------------------------------------
    ['0x1E' CHAR       ['1',   'IEC61131_CHAR']]
    ['0x1F' WCHAR      ['2',   'IEC61131_WCHAR']]
    ['0x20' STRING     ['256', 'IEC61131_STRING']]
    ['0x21' WSTRING    ['512', 'IEC61131_WSTRING']]
    // -----------------------------------------
    // Dates & Times
    // -----------------------------------------
    ['0x22' TIME           ['4', 'IEC61131_TIME']]
    ['0x23' LTIME          ['8', 'IEC61131_LTIME']]
    ['0x24' DATE           ['4', 'IEC61131_DATE']]
    ['0x25' TIME_OF_DAY    ['4', 'IEC61131_TIME_OF_DAY']]
    ['0x26' TOD            ['4', 'IEC61131_TIME_OF_DAY']]
    ['0x27' DATE_AND_TIME  ['4', 'IEC61131_DATE_AND_TIME']]
    ['0x28' DT             ['4', 'IEC61131_DATE_AND_TIME']]
]

[enum uint 32 'ReservedIndexGroups'
    ['0x0000F000' ADSIGRP_SYMTAB]
    ['0x0000F001' ADSIGRP_SYMNAME]
    ['0x0000F002' ADSIGRP_SYMVAL]
    ['0x0000F003' ADSIGRP_SYM_HNDBYNAME]
    ['0x0000F004' ADSIGRP_SYM_VALBYNAME]
    ['0x0000F005' ADSIGRP_SYM_VALBYHND]
    ['0x0000F006' ADSIGRP_SYM_RELEASEHND]
    ['0x0000F007' ADSIGRP_SYM_INFOBYNAME]
    ['0x0000F008' ADSIGRP_SYM_VERSION]
    ['0x0000F009' ADSIGRP_SYM_INFOBYNAMEEX]
    ['0x0000F00A' ADSIGRP_SYM_DOWNLOAD]
    ['0x0000F00B' ADSIGRP_SYM_UPLOAD]
    ['0x0000F00C' ADSIGRP_SYM_UPLOADINFO]
    ['0x0000F010' ADSIGRP_SYMNOTE]
    ['0x0000F020' ADSIGRP_IOIMAGE_RWIB]
    ['0x0000F021' ADSIGRP_IOIMAGE_RWIX]
    ['0x0000F025' ADSIGRP_IOIMAGE_RISIZE]
    ['0x0000F030' ADSIGRP_IOIMAGE_RWOB]
    ['0x0000F031' ADSIGRP_IOIMAGE_RWOX]
    ['0x0000F035' ADSIGRP_IOIMAGE_RWOSIZE]
    ['0x0000F040' ADSIGRP_IOIMAGE_CLEARI]
    ['0x0000F050' ADSIGRP_IOIMAGE_CLEARO]
    ['0x0000F060' ADSIGRP_IOIMAGE_RWIOB]
    ['0x0000F080' ADSIGRP_MULTIPLE_READ]
    ['0x0000F081' ADSIGRP_MULTIPLE_WRITE]
    ['0x0000F082' ADSIGRP_MULTIPLE_READ_WRITE]
    ['0x0000F083' ADSIGRP_MULTIPLE_RELEASE_HANDLE]
    ['0x0000F100' ADSIGRP_DEVICE_DATA]
    ['0x00000000' ADSIOFFS_DEVDATA_ADSSTATE]
    ['0x00000002' ADSIOFFS_DEVDATA_DEVSTATE]
]

[enum uint 32 'ReturnCode'
    // Global Return Codes
    ['0x00' OK]
    ['0x01' INTERNAL_ERROR]
    ['0x02' NO_REALTIME]
    ['0x03' SAVE_ERROR]
    ['0x04' MAILBOX_FULL]
    ['0x05' WRONG_HMSG]
    ['0x06' TARGET_PORT_NOT_FOUND]
    ['0x07' TARGET_HOST_NOT_FOUND]
    ['0x08' UNKNOWN_COMMAND_ID]
    ['0x09' UNKNOWN_TASK_ID]
    ['0x0A' NO_IO]
    ['0x0B' UNKNOWN_ADS_COMMAND]
    ['0x0C' WIN32_ERROR]
    ['0x0D' PORT_NOT_CONNECTED]
    ['0x0E' INVALID_ADS_LENGTH]
    ['0x0F' INVALID_AMS_NET_ID]
    ['0x10' LOW_INSTALLATION_LEVEL]
    ['0x11' NO_DEBUGGING_AVAILABLE]
    ['0x12' PORT_DEACTIVATED]
    ['0x13' PORT_ALREADY_CONNECTED]
    ['0x14' ADS_SYNC_WIN32_ERROR]
    ['0x15' ADS_SYNC_TIMEOUT]
    ['0x16' ADS_SYNC_AMS_ERROR]
    ['0x17' NO_INDEX_MAP_FOR_ADS_AVAILABLE]
    ['0x18' INVALID_ADS_PORT]
    ['0x19' NO_MEMORY]
    ['0x1A' TCP_SENDING_ERROR]
    ['0x1B' HOST_NOT_REACHABLE]
    ['0x1C' INVALID_AMS_FRAGMENT]

    // Router Error-Codes
    ['0x500' ROUTERERR_NOLOCKEDMEMORY]
    ['0x501' ROUTERERR_RESIZEMEMORY]
    ['0x502' ROUTERERR_MAILBOXFULL]
    ['0x503' ROUTERERR_DEBUGBOXFULL]
    ['0x504' ROUTERERR_UNKNOWNPORTTYPE]
    ['0x505' ROUTERERR_NOTINITIALIZED]
    ['0x506' ROUTERERR_PORTALREADYINUSE]
    ['0x507' ROUTERERR_NOTREGISTERED]
    ['0x508' ROUTERERR_NOMOREQUEUES]
    ['0x509' ROUTERERR_INVALIDPORT]
    ['0x50A' ROUTERERR_NOTACTIVATED]

    // General ADS Error-Codes
    ['0x700' ADSERR_DEVICE_ERROR]
    ['0x701' ADSERR_DEVICE_SRVNOTSUPP]
    ['0x702' ADSERR_DEVICE_INVALIDGRP]
    ['0x703' ADSERR_DEVICE_INVALIDOFFSET]
    ['0x704' ADSERR_DEVICE_INVALIDACCESS]
    ['0x705' ADSERR_DEVICE_INVALIDSIZE]
    ['0x706' ADSERR_DEVICE_INVALIDDATA]
    ['0x707' ADSERR_DEVICE_NOTREADY]
    ['0x708' ADSERR_DEVICE_BUSY]
    ['0x709' ADSERR_DEVICE_INVALIDCONTEXT]
    ['0x70A' ADSERR_DEVICE_NOMEMORY]
    ['0x70B' ADSERR_DEVICE_INVALIDPARM]
    ['0x70C' ADSERR_DEVICE_NOTFOUND]
    ['0x70D' ADSERR_DEVICE_SYNTAX]
    ['0x70E' ADSERR_DEVICE_INCOMPATIBLE]
    ['0x70F' ADSERR_DEVICE_EXISTS]
    ['0x710' ADSERR_DEVICE_SYMBOLNOTFOUND]
    ['0x711' ADSERR_DEVICE_SYMBOLVERSIONINVALID]
    ['0x712' ADSERR_DEVICE_INVALIDSTATE]
    ['0x713' ADSERR_DEVICE_TRANSMODENOTSUPP]
    ['0x714' ADSERR_DEVICE_NOTIFYHNDINVALID]
    ['0x715' ADSERR_DEVICE_CLIENTUNKNOWN]
    ['0x716' ADSERR_DEVICE_NOMOREHDLS]
    ['0x717' ADSERR_DEVICE_INVALIDWATCHSIZE]
    ['0x718' ADSERR_DEVICE_NOTINIT]
    ['0x719' ADSERR_DEVICE_TIMEOUT]
    ['0x71A' ADSERR_DEVICE_NOINTERFACE]
    ['0x71B' ADSERR_DEVICE_INVALIDINTERFACE]
    ['0x71C' ADSERR_DEVICE_INVALIDCLSID]
    ['0x71D' ADSERR_DEVICE_INVALIDOBJID]
    ['0x71E' ADSERR_DEVICE_PENDING]
    ['0x71F' ADSERR_DEVICE_ABORTED]
    ['0x720' ADSERR_DEVICE_WARNING]
    ['0x721' ADSERR_DEVICE_INVALIDARRAYIDX]
    ['0x722' ADSERR_DEVICE_SYMBOLNOTACTIVE]
    ['0x723' ADSERR_DEVICE_ACCESSDENIED]
    ['0x724' ADSERR_DEVICE_LICENSENOTFOUND]
    ['0x725' ADSERR_DEVICE_LICENSEEXPIRED]
    ['0x726' ADSERR_DEVICE_LICENSEEXCEEDED]
    ['0x727' ADSERR_DEVICE_LICENSEINVALID]
    ['0x728' ADSERR_DEVICE_LICENSESYSTEMID]
    ['0x729' ADSERR_DEVICE_LICENSENOTIMELIMIT]
    ['0x72A' ADSERR_DEVICE_LICENSEFUTUREISSUE]
    ['0x72B' ADSERR_DEVICE_LICENSETIMETOLONG]
    ['0x72c' ADSERR_DEVICE_EXCEPTION]
    ['0x72D' ADSERR_DEVICE_LICENSEDUPLICATED]
    ['0x72E' ADSERR_DEVICE_SIGNATUREINVALID]
    ['0x72F' ADSERR_DEVICE_CERTIFICATEINVALID]
    ['0x740' ADSERR_CLIENT_ERROR]
    ['0x741' ADSERR_CLIENT_INVALIDPARM]
    ['0x742' ADSERR_CLIENT_LISTEMPTY]
    ['0x743' ADSERR_CLIENT_VARUSED]
    ['0x744' ADSERR_CLIENT_DUPLINVOKEID]
    ['0x745' ADSERR_CLIENT_SYNCTIMEOUT]
    ['0x746' ADSERR_CLIENT_W32ERROR]
    ['0x747' ADSERR_CLIENT_TIMEOUTINVALID]
    ['0x748' ADSERR_CLIENT_PORTNOTOPEN]
    ['0x750' ADSERR_CLIENT_NOAMSADDR]
    ['0x751' ADSERR_CLIENT_SYNCINTERNAL]
    ['0x752' ADSERR_CLIENT_ADDHASH]
    ['0x753' ADSERR_CLIENT_REMOVEHASH]
    ['0x754' ADSERR_CLIENT_NOMORESYM]
    ['0x755' ADSERR_CLIENT_SYNCRESINVALID]

    // RTime Error-Codes
    ['0x1000' RTERR_INTERNAL]
    ['0x1001' RTERR_BADTIMERPERIODS]
    ['0x1002' RTERR_INVALIDTASKPTR]
    ['0x1003' RTERR_INVALIDSTACKPTR]
    ['0x1004' RTERR_PRIOEXISTS]
    ['0x1005' RTERR_NOMORETCB]
    ['0x1006' RTERR_NOMORESEMAS]
    ['0x1007' RTERR_NOMOREQUEUES]
    ['0x100D' RTERR_EXTIRQALREADYDEF]
    ['0x100E' RTERR_EXTIRQNOTDEF]
    ['0x100F' RTERR_EXTIRQINSTALLFAILED]
    ['0x1010' RTERR_IRQLNOTLESSOREQUAL]
    ['0x1017' RTERR_VMXNOTSUPPORTED]
    ['0x1018' RTERR_VMXDISABLED]
    ['0x1019' RTERR_VMXCONTROLSMISSING]
    ['0x101A' RTERR_VMXENABLEFAILS]

    // TCP Windsock Error-Codes
    ['0x274C' WSAETIMEDOUT]
    ['0x274D' WSAECONNREFUSED]
    ['0x2751' WSAEHOSTUNREACH]
]
