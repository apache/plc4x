//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package model

import (
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
)

type KnxObjectProperties uint8

type IKnxObjectProperties interface {
    Name() string
    DataTypeId() string
    Text() string
    Serialize(io utils.WriteBuffer) error
}

const(
    KnxObjectProperties_PID_G_1 KnxObjectProperties = 1
    KnxObjectProperties_PID_G_2 KnxObjectProperties = 2
    KnxObjectProperties_PID_G_3 KnxObjectProperties = 3
    KnxObjectProperties_PID_G_4 KnxObjectProperties = 4
    KnxObjectProperties_PID_G_5 KnxObjectProperties = 5
    KnxObjectProperties_PID_G_6 KnxObjectProperties = 6
    KnxObjectProperties_PID_G_7 KnxObjectProperties = 7
    KnxObjectProperties_PID_G_8 KnxObjectProperties = 8
    KnxObjectProperties_PID_G_9 KnxObjectProperties = 9
    KnxObjectProperties_PID_G_10 KnxObjectProperties = 10
    KnxObjectProperties_PID_G_11 KnxObjectProperties = 11
    KnxObjectProperties_PID_G_12 KnxObjectProperties = 12
    KnxObjectProperties_PID_G_13 KnxObjectProperties = 13
    KnxObjectProperties_PID_G_14 KnxObjectProperties = 14
    KnxObjectProperties_PID_G_15 KnxObjectProperties = 15
    KnxObjectProperties_PID_G_16 KnxObjectProperties = 16
    KnxObjectProperties_PID_G_17 KnxObjectProperties = 17
    KnxObjectProperties_PID_G_18 KnxObjectProperties = 18
    KnxObjectProperties_PID_G_19 KnxObjectProperties = 19
    KnxObjectProperties_PID_G_20 KnxObjectProperties = 20
    KnxObjectProperties_PID_G_21 KnxObjectProperties = 21
    KnxObjectProperties_PID_G_22 KnxObjectProperties = 22
    KnxObjectProperties_PID_G_23 KnxObjectProperties = 23
    KnxObjectProperties_PID_G_24 KnxObjectProperties = 24
    KnxObjectProperties_PID_G_25 KnxObjectProperties = 25
    KnxObjectProperties_PID_G_26 KnxObjectProperties = 26
    KnxObjectProperties_PID_G_27 KnxObjectProperties = 27
    KnxObjectProperties_PID_G_28 KnxObjectProperties = 28
    KnxObjectProperties_PID_G_29 KnxObjectProperties = 29
    KnxObjectProperties_PID_G_30 KnxObjectProperties = 30
    KnxObjectProperties_PID_0_51 KnxObjectProperties = 51
    KnxObjectProperties_PID_0_52 KnxObjectProperties = 52
    KnxObjectProperties_PID_0_53 KnxObjectProperties = 53
    KnxObjectProperties_PID_0_54 KnxObjectProperties = 54
    KnxObjectProperties_PID_0_55 KnxObjectProperties = 55
    KnxObjectProperties_PID_0_56 KnxObjectProperties = 56
    KnxObjectProperties_PID_0_57 KnxObjectProperties = 57
    KnxObjectProperties_PID_0_58 KnxObjectProperties = 58
    KnxObjectProperties_PID_0_59 KnxObjectProperties = 59
    KnxObjectProperties_PID_0_60 KnxObjectProperties = 60
    KnxObjectProperties_PID_0_61 KnxObjectProperties = 61
    KnxObjectProperties_PID_0_62 KnxObjectProperties = 62
    KnxObjectProperties_PID_0_63 KnxObjectProperties = 63
    KnxObjectProperties_PID_0_64 KnxObjectProperties = 64
    KnxObjectProperties_PID_0_65 KnxObjectProperties = 65
    KnxObjectProperties_PID_0_66 KnxObjectProperties = 66
    KnxObjectProperties_PID_0_67 KnxObjectProperties = 67
    KnxObjectProperties_PID_0_68 KnxObjectProperties = 68
    KnxObjectProperties_PID_0_69 KnxObjectProperties = 69
    KnxObjectProperties_PID_0_70 KnxObjectProperties = 70
    KnxObjectProperties_PID_0_71 KnxObjectProperties = 71
    KnxObjectProperties_PID_0_72 KnxObjectProperties = 72
    KnxObjectProperties_PID_0_73 KnxObjectProperties = 73
    KnxObjectProperties_PID_0_74 KnxObjectProperties = 74
    KnxObjectProperties_PID_0_75 KnxObjectProperties = 75
    KnxObjectProperties_PID_0_76 KnxObjectProperties = 76
    KnxObjectProperties_PID_0_77 KnxObjectProperties = 77
    KnxObjectProperties_PID_0_78 KnxObjectProperties = 78
    KnxObjectProperties_PID_0_79 KnxObjectProperties = 79
    KnxObjectProperties_PID_0_80 KnxObjectProperties = 80
    KnxObjectProperties_PID_0_81 KnxObjectProperties = 81
    KnxObjectProperties_PID_0_82 KnxObjectProperties = 82
    KnxObjectProperties_PID_0_83 KnxObjectProperties = 83
    KnxObjectProperties_PID_0_84 KnxObjectProperties = 84
    KnxObjectProperties_PID_0_85 KnxObjectProperties = 85
    KnxObjectProperties_PID_0_86 KnxObjectProperties = 86
    KnxObjectProperties_PID_9_51 KnxObjectProperties = 51
    KnxObjectProperties_PID_9_52 KnxObjectProperties = 52
    KnxObjectProperties_PID_6_51 KnxObjectProperties = 51
    KnxObjectProperties_PID_6_52 KnxObjectProperties = 52
    KnxObjectProperties_PID_6_53 KnxObjectProperties = 53
    KnxObjectProperties_PID_6_54 KnxObjectProperties = 54
    KnxObjectProperties_PID_6_55 KnxObjectProperties = 55
    KnxObjectProperties_PID_6_56 KnxObjectProperties = 56
    KnxObjectProperties_PID_6_57 KnxObjectProperties = 57
    KnxObjectProperties_PID_6_58 KnxObjectProperties = 58
    KnxObjectProperties_PID_6_63 KnxObjectProperties = 63
    KnxObjectProperties_PID_6_67 KnxObjectProperties = 67
    KnxObjectProperties_PID_6_112 KnxObjectProperties = 112
    KnxObjectProperties_PID_11_51 KnxObjectProperties = 51
    KnxObjectProperties_PID_11_52 KnxObjectProperties = 52
    KnxObjectProperties_PID_11_53 KnxObjectProperties = 53
    KnxObjectProperties_PID_11_54 KnxObjectProperties = 54
    KnxObjectProperties_PID_11_55 KnxObjectProperties = 55
    KnxObjectProperties_PID_11_56 KnxObjectProperties = 56
    KnxObjectProperties_PID_11_57 KnxObjectProperties = 57
    KnxObjectProperties_PID_11_58 KnxObjectProperties = 58
    KnxObjectProperties_PID_11_59 KnxObjectProperties = 59
    KnxObjectProperties_PID_11_60 KnxObjectProperties = 60
    KnxObjectProperties_PID_11_61 KnxObjectProperties = 61
    KnxObjectProperties_PID_11_62 KnxObjectProperties = 62
    KnxObjectProperties_PID_11_63 KnxObjectProperties = 63
    KnxObjectProperties_PID_11_64 KnxObjectProperties = 64
    KnxObjectProperties_PID_11_65 KnxObjectProperties = 65
    KnxObjectProperties_PID_11_66 KnxObjectProperties = 66
    KnxObjectProperties_PID_11_67 KnxObjectProperties = 67
    KnxObjectProperties_PID_11_68 KnxObjectProperties = 68
    KnxObjectProperties_PID_11_69 KnxObjectProperties = 69
    KnxObjectProperties_PID_11_70 KnxObjectProperties = 70
    KnxObjectProperties_PID_11_71 KnxObjectProperties = 71
    KnxObjectProperties_PID_11_72 KnxObjectProperties = 72
    KnxObjectProperties_PID_11_73 KnxObjectProperties = 73
    KnxObjectProperties_PID_11_74 KnxObjectProperties = 74
    KnxObjectProperties_PID_11_75 KnxObjectProperties = 75
    KnxObjectProperties_PID_11_76 KnxObjectProperties = 76
    KnxObjectProperties_PID_11_91 KnxObjectProperties = 91
    KnxObjectProperties_PID_11_92 KnxObjectProperties = 92
    KnxObjectProperties_PID_11_93 KnxObjectProperties = 93
    KnxObjectProperties_PID_11_94 KnxObjectProperties = 94
    KnxObjectProperties_PID_11_95 KnxObjectProperties = 95
    KnxObjectProperties_PID_11_96 KnxObjectProperties = 96
    KnxObjectProperties_PID_11_97 KnxObjectProperties = 97
    KnxObjectProperties_PID_17_51 KnxObjectProperties = 51
    KnxObjectProperties_PID_17_52 KnxObjectProperties = 52
    KnxObjectProperties_PID_17_53 KnxObjectProperties = 53
    KnxObjectProperties_PID_17_54 KnxObjectProperties = 54
    KnxObjectProperties_PID_17_55 KnxObjectProperties = 55
    KnxObjectProperties_PID_17_56 KnxObjectProperties = 56
    KnxObjectProperties_PID_17_57 KnxObjectProperties = 57
    KnxObjectProperties_PID_17_58 KnxObjectProperties = 58
    KnxObjectProperties_PID_17_59 KnxObjectProperties = 59
    KnxObjectProperties_PID_17_60 KnxObjectProperties = 60
    KnxObjectProperties_PID_17_61 KnxObjectProperties = 61
    KnxObjectProperties_PID_19_51 KnxObjectProperties = 51
    KnxObjectProperties_PID_19_56 KnxObjectProperties = 56
    KnxObjectProperties_PID_19_57 KnxObjectProperties = 57
    KnxObjectProperties_PID_19_58 KnxObjectProperties = 58
    KnxObjectProperties_PID_19_59 KnxObjectProperties = 59
    KnxObjectProperties_PID_19_60 KnxObjectProperties = 60
    KnxObjectProperties_PID_19_61 KnxObjectProperties = 61
    KnxObjectProperties_PID_19_62 KnxObjectProperties = 62
    KnxObjectProperties_PID_19_63 KnxObjectProperties = 63
    KnxObjectProperties_PID_409_110 KnxObjectProperties = 110
    KnxObjectProperties_PID_409_111 KnxObjectProperties = 111
    KnxObjectProperties_PID_410_110 KnxObjectProperties = 110
    KnxObjectProperties_PID_410_111 KnxObjectProperties = 111
    KnxObjectProperties_PID_417_101 KnxObjectProperties = 101
    KnxObjectProperties_PID_417_102 KnxObjectProperties = 102
    KnxObjectProperties_PID_417_103 KnxObjectProperties = 103
    KnxObjectProperties_PID_417_104 KnxObjectProperties = 104
    KnxObjectProperties_PID_417_105 KnxObjectProperties = 105
    KnxObjectProperties_PID_417_106 KnxObjectProperties = 106
    KnxObjectProperties_PID_417_107 KnxObjectProperties = 107
    KnxObjectProperties_PID_417_108 KnxObjectProperties = 108
    KnxObjectProperties_PID_417_109 KnxObjectProperties = 109
    KnxObjectProperties_PID_417_110 KnxObjectProperties = 110
    KnxObjectProperties_PID_417_111 KnxObjectProperties = 111
    KnxObjectProperties_PID_417_112 KnxObjectProperties = 112
    KnxObjectProperties_PID_417_113 KnxObjectProperties = 113
    KnxObjectProperties_PID_417_114 KnxObjectProperties = 114
    KnxObjectProperties_PID_417_115 KnxObjectProperties = 115
    KnxObjectProperties_PID_417_116 KnxObjectProperties = 116
    KnxObjectProperties_PID_417_117 KnxObjectProperties = 117
    KnxObjectProperties_PID_417_118 KnxObjectProperties = 118
    KnxObjectProperties_PID_417_119 KnxObjectProperties = 119
    KnxObjectProperties_PID_417_120 KnxObjectProperties = 120
    KnxObjectProperties_PID_418_101 KnxObjectProperties = 101
    KnxObjectProperties_PID_418_102 KnxObjectProperties = 102
    KnxObjectProperties_PID_418_103 KnxObjectProperties = 103
    KnxObjectProperties_PID_418_104 KnxObjectProperties = 104
    KnxObjectProperties_PID_418_105 KnxObjectProperties = 105
    KnxObjectProperties_PID_418_106 KnxObjectProperties = 106
    KnxObjectProperties_PID_418_107 KnxObjectProperties = 107
    KnxObjectProperties_PID_418_108 KnxObjectProperties = 108
    KnxObjectProperties_PID_418_109 KnxObjectProperties = 109
    KnxObjectProperties_PID_418_110 KnxObjectProperties = 110
    KnxObjectProperties_PID_418_111 KnxObjectProperties = 111
    KnxObjectProperties_PID_418_112 KnxObjectProperties = 112
    KnxObjectProperties_PID_418_113 KnxObjectProperties = 113
    KnxObjectProperties_PID_418_114 KnxObjectProperties = 114
    KnxObjectProperties_PID_418_115 KnxObjectProperties = 115
    KnxObjectProperties_PID_418_116 KnxObjectProperties = 116
    KnxObjectProperties_PID_418_117 KnxObjectProperties = 117
    KnxObjectProperties_PID_418_118 KnxObjectProperties = 118
    KnxObjectProperties_PID_418_119 KnxObjectProperties = 119
    KnxObjectProperties_PID_418_120 KnxObjectProperties = 120
    KnxObjectProperties_PID_418_121 KnxObjectProperties = 121
    KnxObjectProperties_PID_418_122 KnxObjectProperties = 122
    KnxObjectProperties_PID_418_123 KnxObjectProperties = 123
    KnxObjectProperties_PID_418_124 KnxObjectProperties = 124
    KnxObjectProperties_PID_418_125 KnxObjectProperties = 125
    KnxObjectProperties_PID_418_126 KnxObjectProperties = 126
    KnxObjectProperties_PID_418_127 KnxObjectProperties = 127
    KnxObjectProperties_PID_418_128 KnxObjectProperties = 128
    KnxObjectProperties_PID_418_129 KnxObjectProperties = 129
    KnxObjectProperties_PID_418_130 KnxObjectProperties = 130
    KnxObjectProperties_PID_418_131 KnxObjectProperties = 131
    KnxObjectProperties_PID_418_132 KnxObjectProperties = 132
    KnxObjectProperties_PID_418_133 KnxObjectProperties = 133
    KnxObjectProperties_PID_418_134 KnxObjectProperties = 134
    KnxObjectProperties_PID_420_51 KnxObjectProperties = 51
    KnxObjectProperties_PID_420_52 KnxObjectProperties = 52
    KnxObjectProperties_PID_420_53 KnxObjectProperties = 53
    KnxObjectProperties_PID_421_51 KnxObjectProperties = 51
    KnxObjectProperties_PID_421_52 KnxObjectProperties = 52
    KnxObjectProperties_PID_800_51 KnxObjectProperties = 51
    KnxObjectProperties_PID_800_52 KnxObjectProperties = 52
    KnxObjectProperties_PID_800_53 KnxObjectProperties = 53
    KnxObjectProperties_PID_800_54 KnxObjectProperties = 54
    KnxObjectProperties_PID_800_55 KnxObjectProperties = 55
    KnxObjectProperties_PID_800_57 KnxObjectProperties = 57
    KnxObjectProperties_PID_800_58 KnxObjectProperties = 58
    KnxObjectProperties_PID_800_60 KnxObjectProperties = 60
    KnxObjectProperties_PID_800_61 KnxObjectProperties = 61
    KnxObjectProperties_PID_800_62 KnxObjectProperties = 62
    KnxObjectProperties_PID_800_63 KnxObjectProperties = 63
    KnxObjectProperties_PID_800_64 KnxObjectProperties = 64
    KnxObjectProperties_PID_800_65 KnxObjectProperties = 65
    KnxObjectProperties_PID_800_66 KnxObjectProperties = 66
    KnxObjectProperties_PID_800_67 KnxObjectProperties = 67
    KnxObjectProperties_PID_800_68 KnxObjectProperties = 68
    KnxObjectProperties_PID_800_69 KnxObjectProperties = 69
    KnxObjectProperties_PID_801_51 KnxObjectProperties = 51
    KnxObjectProperties_PID_801_52 KnxObjectProperties = 52
    KnxObjectProperties_PID_801_53 KnxObjectProperties = 53
)


func (e KnxObjectProperties) Name() string {
    switch e  {
        case 1: { /* '1' */
            return "PID_OBJECT_TYPE"
        }
        case 10: { /* '10' */
            return "PID_SERVICES_SUPPORTED"
        }
        case 101: { /* '101' */
            return "PID_ON_DELAY"
        }
        case 102: { /* '102' */
            return "PID_OFF_DELAY"
        }
        case 103: { /* '103' */
            return "PID_TIMED_ON_DURATION"
        }
        case 104: { /* '104' */
            return "PID_PREWARNING_DURATION"
        }
        case 105: { /* '105' */
            return "PID_TRANSMISSION_CYCLE_TIME"
        }
        case 106: { /* '106' */
            return "PID_BUS_POWER_UP_MESSAGE_DELAY"
        }
        case 107: { /* '107' */
            return "PID_BEHAVIOUR_AT_LOCKING"
        }
        case 108: { /* '108' */
            return "PID_BEHAVIOUR_AT_UNLOCKING"
        }
        case 109: { /* '109' */
            return "PID_BEHAVIOUR_BUS_POWER_UP"
        }
        case 11: { /* '11' */
            return "PID_SERIAL_NUMBER"
        }
        case 110: { /* '110' */
            return "PID_CHANGE_OF_VALUE"
        }
        case 111: { /* '111' */
            return "PID_REPETITION_TIME"
        }
        case 112: { /* '112' */
            return "PID_RF_ENABLE_SBC"
        }
        case 113: { /* '113' */
            return "PID_MANUAL_OFF_ENABLE"
        }
        case 114: { /* '114' */
            return "PID_INVERT_LOCK_DEVICE"
        }
        case 115: { /* '115' */
            return "PID_LOCK_STATE"
        }
        case 116: { /* '116' */
            return "PID_UNLOCK_STATE"
        }
        case 117: { /* '117' */
            return "PID_STATE_FOR_SCENE_NUMBER"
        }
        case 118: { /* '118' */
            return "PID_STORAGE_FUNCTION_FOR_SCENE"
        }
        case 119: { /* '119' */
            return "PID_BUS_POWER_UP_STATE"
        }
        case 12: { /* '12' */
            return "PID_MANUFACTURER_ID"
        }
        case 120: { /* '120' */
            return "PID_BEHAVIOUR_BUS_POWER_UP"
        }
        case 121: { /* '121' */
            return "PID_TIMED_ON_RETRIGGER_FUNCTION"
        }
        case 122: { /* '122' */
            return "PID_MANUAL_OFF_ENABLE"
        }
        case 123: { /* '123' */
            return "PID_INVERT_LOCK_DEVICE"
        }
        case 124: { /* '124' */
            return "PID_BEHAVIOUR_AT_LOCKING"
        }
        case 125: { /* '125' */
            return "PID_BEHAVIOUR_AT_UNLOCKING"
        }
        case 126: { /* '126' */
            return "PID_LOCK_SETVALUE"
        }
        case 127: { /* '127' */
            return "PID_UNLOCK_SETVALUE"
        }
        case 128: { /* '128' */
            return "PID_BIGHTNESS_FOR_SCENE"
        }
        case 129: { /* '129' */
            return "PID_STORAGE_FUNCTION_FOR_SCENE"
        }
        case 13: { /* '13' */
            return "PID_PROGRAM_VERSION"
        }
        case 130: { /* '130' */
            return "PID_DELTA_DIMMING_VALUE"
        }
        case 131: { /* '131' */
            return "PID_BEHAVIOUR_BUS_POWER_UP"
        }
        case 132: { /* '132' */
            return "PID_BEHAVIOUR_BUS_POWER_UP_SET_VALUE"
        }
        case 133: { /* '133' */
            return "PID_BEHAVIOUR_BUS_POWER_DOWN"
        }
        case 134: { /* '134' */
            return "PID_BUS_POWER_DOWN_SET_VALUE"
        }
        case 14: { /* '14' */
            return "PID_DEVICE_CONTROL"
        }
        case 15: { /* '15' */
            return "PID_ORDER_INFO"
        }
        case 16: { /* '16' */
            return "PID_PEI_TYPE"
        }
        case 17: { /* '17' */
            return "PID_PORT_CONFIGURATION"
        }
        case 18: { /* '18' */
            return "PID_POLL_GROUP_SETTINGS"
        }
        case 19: { /* '19' */
            return "PID_MANUFACTURER_DATA"
        }
        case 2: { /* '2' */
            return "PID_OBJECT_NAME"
        }
        case 20: { /* '20' */
            return "PID_ENABLE"
        }
        case 21: { /* '21' */
            return "PID_DESCRIPTION"
        }
        case 22: { /* '22' */
            return "PID_FILE"
        }
        case 23: { /* '23' */
            return "PID_TABLE"
        }
        case 24: { /* '24' */
            return "PID_ENROL"
        }
        case 25: { /* '25' */
            return "PID_VERSION"
        }
        case 26: { /* '26' */
            return "PID_GROUP_OBJECT_LINK"
        }
        case 27: { /* '27' */
            return "PID_MCB_TABLE"
        }
        case 28: { /* '28' */
            return "PID_ERROR_CODE"
        }
        case 29: { /* '29' */
            return "PID_OBJECT_INDEX"
        }
        case 3: { /* '3' */
            return "PID_SEMAPHOR"
        }
        case 30: { /* '30' */
            return "PID_DOWNLOAD_COUNTER"
        }
        case 4: { /* '4' */
            return "PID_GROUP_OBJECT_REFERENCE"
        }
        case 5: { /* '5' */
            return "PID_LOAD_STATE_CONTROL"
        }
        case 51: { /* '51' */
            return "PID_ROUTING_COUNT"
        }
        case 52: { /* '52' */
            return "PID_MAX_RETRY_COUNT"
        }
        case 53: { /* '53' */
            return "PID_ERROR_FLAGS"
        }
        case 54: { /* '54' */
            return "PID_PROGMODE"
        }
        case 55: { /* '55' */
            return "PID_PRODUCT_ID"
        }
        case 56: { /* '56' */
            return "PID_MAX_APDULENGTH"
        }
        case 57: { /* '57' */
            return "PID_SUBNET_ADDR"
        }
        case 58: { /* '58' */
            return "PID_DEVICE_ADDR"
        }
        case 59: { /* '59' */
            return "PID_PB_CONFIG"
        }
        case 6: { /* '6' */
            return "PID_RUN_STATE_CONTROL"
        }
        case 60: { /* '60' */
            return "PID_ADDR_REPORT"
        }
        case 61: { /* '61' */
            return "PID_ADDR_CHECK"
        }
        case 62: { /* '62' */
            return "PID_OBJECT_VALUE"
        }
        case 63: { /* '63' */
            return "PID_OBJECTLINK"
        }
        case 64: { /* '64' */
            return "PID_APPLICATION"
        }
        case 65: { /* '65' */
            return "PID_PARAMETER"
        }
        case 66: { /* '66' */
            return "PID_OBJECTADDRESS"
        }
        case 67: { /* '67' */
            return "PID_PSU_TYPE"
        }
        case 68: { /* '68' */
            return "PID_PSU_STATUS"
        }
        case 69: { /* '69' */
            return "PID_PSU_ENABLE"
        }
        case 7: { /* '7' */
            return "PID_TABLE_REFERENCE"
        }
        case 70: { /* '70' */
            return "PID_DOMAIN_ADDRESS"
        }
        case 71: { /* '71' */
            return "PID_IO_LIST"
        }
        case 72: { /* '72' */
            return "PID_MGT_DESCRIPTOR_01"
        }
        case 73: { /* '73' */
            return "PID_PL110_PARAM"
        }
        case 74: { /* '74' */
            return "PID_RF_REPEAT_COUNTER"
        }
        case 75: { /* '75' */
            return "PID_RECEIVE_BLOCK_TABLE"
        }
        case 76: { /* '76' */
            return "PID_RANDOM_PAUSE_TABLE"
        }
        case 77: { /* '77' */
            return "PID_RECEIVE_BLOCK_NR"
        }
        case 78: { /* '78' */
            return "PID_HARDWARE_TYPE"
        }
        case 79: { /* '79' */
            return "PID_RETRANSMITTER_NUMBER"
        }
        case 8: { /* '8' */
            return "PID_SERVICE_CONTROL"
        }
        case 80: { /* '80' */
            return "PID_SERIAL_NR_TABLE"
        }
        case 81: { /* '81' */
            return "PID_BIBATMASTER_ADDRESS"
        }
        case 82: { /* '82' */
            return "PID_RF_DOMAIN_ADDRESS"
        }
        case 83: { /* '83' */
            return "PID_DEVICE_DESCRIPTOR"
        }
        case 84: { /* '84' */
            return "PID_METERING_FILTER_TABLE"
        }
        case 85: { /* '85' */
            return "PID_GROUP_TELEGR_RATE_LIMIT_TIME_BASE"
        }
        case 86: { /* '86' */
            return "PID_GROUP_TELEGR_RATE_LIMIT_NO_OF_TELEGR"
        }
        case 9: { /* '9' */
            return "PID_FIRMWARE_REVISION"
        }
        case 91: { /* '91' */
            return "PID_BACKBONE_KEY"
        }
        case 92: { /* '92' */
            return "PID_DEVICE_AUTHENTICATION_CODE"
        }
        case 93: { /* '93' */
            return "PID_PASSWORD_HASHES"
        }
        case 94: { /* '94' */
            return "PID_SECURED_SERVICE_FAMILIES"
        }
        case 95: { /* '95' */
            return "PID_MULTICAST_LATENCY_TOLERANCE"
        }
        case 96: { /* '96' */
            return "PID_SYNC_LATENCY_FRACTION"
        }
        case 97: { /* '97' */
            return "PID_TUNNELLING_USERS"
        }
        default: {
            return ""
        }
    }
}

func (e KnxObjectProperties) DataTypeId() string {
    switch e  {
        case 1: { /* '1' */
            return "PDT-4"
        }
        case 10: { /* '10' */
            return "null"
        }
        case 101: { /* '101' */
            return "null"
        }
        case 102: { /* '102' */
            return "null"
        }
        case 103: { /* '103' */
            return "null"
        }
        case 104: { /* '104' */
            return "null"
        }
        case 105: { /* '105' */
            return "null"
        }
        case 106: { /* '106' */
            return "null"
        }
        case 107: { /* '107' */
            return "null"
        }
        case 108: { /* '108' */
            return "null"
        }
        case 109: { /* '109' */
            return "null"
        }
        case 11: { /* '11' */
            return "PDT-22"
        }
        case 110: { /* '110' */
            return "null"
        }
        case 111: { /* '111' */
            return "null"
        }
        case 112: { /* '112' */
            return "PDT-62"
        }
        case 113: { /* '113' */
            return "null"
        }
        case 114: { /* '114' */
            return "null"
        }
        case 115: { /* '115' */
            return "null"
        }
        case 116: { /* '116' */
            return "null"
        }
        case 117: { /* '117' */
            return "null"
        }
        case 118: { /* '118' */
            return "null"
        }
        case 119: { /* '119' */
            return "null"
        }
        case 12: { /* '12' */
            return "PDT-4"
        }
        case 120: { /* '120' */
            return "null"
        }
        case 121: { /* '121' */
            return "null"
        }
        case 122: { /* '122' */
            return "null"
        }
        case 123: { /* '123' */
            return "null"
        }
        case 124: { /* '124' */
            return "null"
        }
        case 125: { /* '125' */
            return "null"
        }
        case 126: { /* '126' */
            return "null"
        }
        case 127: { /* '127' */
            return "null"
        }
        case 128: { /* '128' */
            return "null"
        }
        case 129: { /* '129' */
            return "null"
        }
        case 13: { /* '13' */
            return "PDT-21"
        }
        case 130: { /* '130' */
            return "null"
        }
        case 131: { /* '131' */
            return "null"
        }
        case 132: { /* '132' */
            return "null"
        }
        case 133: { /* '133' */
            return "null"
        }
        case 134: { /* '134' */
            return "null"
        }
        case 14: { /* '14' */
            return "PDT-51"
        }
        case 15: { /* '15' */
            return "PDT-26"
        }
        case 16: { /* '16' */
            return "PDT-2"
        }
        case 17: { /* '17' */
            return "PDT-2"
        }
        case 18: { /* '18' */
            return "PDT-13"
        }
        case 19: { /* '19' */
            return "PDT-20"
        }
        case 2: { /* '2' */
            return "PDT-2"
        }
        case 20: { /* '20' */
            return "null"
        }
        case 21: { /* '21' */
            return "PDT-2"
        }
        case 22: { /* '22' */
            return "null"
        }
        case 23: { /* '23' */
            return "PDT-4"
        }
        case 24: { /* '24' */
            return "PDT-62"
        }
        case 25: { /* '25' */
            return "PDT-48"
        }
        case 26: { /* '26' */
            return "PDT-62"
        }
        case 27: { /* '27' */
            return "PDT-24"
        }
        case 28: { /* '28' */
            return "PDT-17"
        }
        case 29: { /* '29' */
            return "PDT-2"
        }
        case 3: { /* '3' */
            return "null"
        }
        case 30: { /* '30' */
            return "PDT-2"
        }
        case 4: { /* '4' */
            return "null"
        }
        case 5: { /* '5' */
            return "PDT-0"
        }
        case 51: { /* '51' */
            return "PDT-2"
        }
        case 52: { /* '52' */
            return "PDT-17"
        }
        case 53: { /* '53' */
            return "PDT-2"
        }
        case 54: { /* '54' */
            return "PDT-51"
        }
        case 55: { /* '55' */
            return "PDT-26"
        }
        case 56: { /* '56' */
            return "PDT-4"
        }
        case 57: { /* '57' */
            return "PDT-2"
        }
        case 58: { /* '58' */
            return "PDT-2"
        }
        case 59: { /* '59' */
            return "PDT-20"
        }
        case 6: { /* '6' */
            return "PDT-0"
        }
        case 60: { /* '60' */
            return "PDT-22"
        }
        case 61: { /* '61' */
            return "PDT-17"
        }
        case 62: { /* '62' */
            return "PDT-62"
        }
        case 63: { /* '63' */
            return "PDT-62"
        }
        case 64: { /* '64' */
            return "PDT-62"
        }
        case 65: { /* '65' */
            return "PDT-62"
        }
        case 66: { /* '66' */
            return "PDT-62"
        }
        case 67: { /* '67' */
            return "PDT-4"
        }
        case 68: { /* '68' */
            return "PDT-50"
        }
        case 69: { /* '69' */
            return "PDT-53"
        }
        case 7: { /* '7' */
            return "PDT-4"
        }
        case 70: { /* '70' */
            return "PDT-4"
        }
        case 71: { /* '71' */
            return "PDT-4"
        }
        case 72: { /* '72' */
            return "PDT-26"
        }
        case 73: { /* '73' */
            return "PDT-17"
        }
        case 74: { /* '74' */
            return "PDT-2"
        }
        case 75: { /* '75' */
            return "PDT-2"
        }
        case 76: { /* '76' */
            return "PDT-2"
        }
        case 77: { /* '77' */
            return "PDT-2"
        }
        case 78: { /* '78' */
            return "PDT-22"
        }
        case 79: { /* '79' */
            return "PDT-2"
        }
        case 8: { /* '8' */
            return "PDT-4"
        }
        case 80: { /* '80' */
            return "PDT-22"
        }
        case 81: { /* '81' */
            return "PDT-4"
        }
        case 82: { /* '82' */
            return "PDT-22"
        }
        case 83: { /* '83' */
            return "PDT-18"
        }
        case 84: { /* '84' */
            return "PDT-24"
        }
        case 85: { /* '85' */
            return "PDT-4"
        }
        case 86: { /* '86' */
            return "PDT-4"
        }
        case 9: { /* '9' */
            return "PDT-2"
        }
        case 91: { /* '91' */
            return "PDT-32"
        }
        case 92: { /* '92' */
            return "PDT-32"
        }
        case 93: { /* '93' */
            return "PDT-32"
        }
        case 94: { /* '94' */
            return "PDT-62"
        }
        case 95: { /* '95' */
            return "PDT-4"
        }
        case 96: { /* '96' */
            return "PDT-54"
        }
        case 97: { /* '97' */
            return "PDT-18"
        }
        default: {
            return ""
        }
    }
}

func (e KnxObjectProperties) Text() string {
    switch e  {
        case 1: { /* '1' */
            return "Interface Object Type"
        }
        case 10: { /* '10' */
            return "Services Supported"
        }
        case 101: { /* '101' */
            return ""
        }
        case 102: { /* '102' */
            return ""
        }
        case 103: { /* '103' */
            return ""
        }
        case 104: { /* '104' */
            return ""
        }
        case 105: { /* '105' */
            return ""
        }
        case 106: { /* '106' */
            return ""
        }
        case 107: { /* '107' */
            return ""
        }
        case 108: { /* '108' */
            return ""
        }
        case 109: { /* '109' */
            return ""
        }
        case 11: { /* '11' */
            return "KNX Serial Number"
        }
        case 110: { /* '110' */
            return ""
        }
        case 111: { /* '111' */
            return ""
        }
        case 112: { /* '112' */
            return ""
        }
        case 113: { /* '113' */
            return ""
        }
        case 114: { /* '114' */
            return ""
        }
        case 115: { /* '115' */
            return ""
        }
        case 116: { /* '116' */
            return ""
        }
        case 117: { /* '117' */
            return ""
        }
        case 118: { /* '118' */
            return ""
        }
        case 119: { /* '119' */
            return ""
        }
        case 12: { /* '12' */
            return "Manufacturer Identifier"
        }
        case 120: { /* '120' */
            return ""
        }
        case 121: { /* '121' */
            return ""
        }
        case 122: { /* '122' */
            return ""
        }
        case 123: { /* '123' */
            return ""
        }
        case 124: { /* '124' */
            return ""
        }
        case 125: { /* '125' */
            return ""
        }
        case 126: { /* '126' */
            return ""
        }
        case 127: { /* '127' */
            return ""
        }
        case 128: { /* '128' */
            return ""
        }
        case 129: { /* '129' */
            return ""
        }
        case 13: { /* '13' */
            return "Application Version"
        }
        case 130: { /* '130' */
            return ""
        }
        case 131: { /* '131' */
            return ""
        }
        case 132: { /* '132' */
            return ""
        }
        case 133: { /* '133' */
            return ""
        }
        case 134: { /* '134' */
            return ""
        }
        case 14: { /* '14' */
            return "Device Control"
        }
        case 15: { /* '15' */
            return "Order Info"
        }
        case 16: { /* '16' */
            return "PEI Type"
        }
        case 17: { /* '17' */
            return "PortADDR"
        }
        case 18: { /* '18' */
            return "Polling Group Settings"
        }
        case 19: { /* '19' */
            return "Manufacturer Data"
        }
        case 2: { /* '2' */
            return "Interface Object Name"
        }
        case 20: { /* '20' */
            return ""
        }
        case 21: { /* '21' */
            return "Description"
        }
        case 22: { /* '22' */
            return ""
        }
        case 23: { /* '23' */
            return "Table"
        }
        case 24: { /* '24' */
            return "Interface Object Link"
        }
        case 25: { /* '25' */
            return "Version"
        }
        case 26: { /* '26' */
            return "Group Address Assignment"
        }
        case 27: { /* '27' */
            return "Memory Control Table"
        }
        case 28: { /* '28' */
            return "Error Code"
        }
        case 29: { /* '29' */
            return "Object Index"
        }
        case 3: { /* '3' */
            return "Semaphor"
        }
        case 30: { /* '30' */
            return "Download Counter"
        }
        case 4: { /* '4' */
            return "Group Object Reference"
        }
        case 5: { /* '5' */
            return "Load Control"
        }
        case 51: { /* '51' */
            return "Routing Count"
        }
        case 52: { /* '52' */
            return "Maximum Retry Count"
        }
        case 53: { /* '53' */
            return "Error Flags"
        }
        case 54: { /* '54' */
            return "Programming Mode"
        }
        case 55: { /* '55' */
            return "Product Identification"
        }
        case 56: { /* '56' */
            return "Max. APDU-Length"
        }
        case 57: { /* '57' */
            return "Subnetwork Address"
        }
        case 58: { /* '58' */
            return "Device Address"
        }
        case 59: { /* '59' */
            return "Config Link"
        }
        case 6: { /* '6' */
            return "Run Control"
        }
        case 60: { /* '60' */
            return ""
        }
        case 61: { /* '61' */
            return ""
        }
        case 62: { /* '62' */
            return ""
        }
        case 63: { /* '63' */
            return ""
        }
        case 64: { /* '64' */
            return ""
        }
        case 65: { /* '65' */
            return ""
        }
        case 66: { /* '66' */
            return ""
        }
        case 67: { /* '67' */
            return ""
        }
        case 68: { /* '68' */
            return ""
        }
        case 69: { /* '69' */
            return ""
        }
        case 7: { /* '7' */
            return "Table Reference"
        }
        case 70: { /* '70' */
            return "Domain Address"
        }
        case 71: { /* '71' */
            return ""
        }
        case 72: { /* '72' */
            return "Management Descriptor 1"
        }
        case 73: { /* '73' */
            return "PL110 Parameters"
        }
        case 74: { /* '74' */
            return ""
        }
        case 75: { /* '75' */
            return ""
        }
        case 76: { /* '76' */
            return ""
        }
        case 77: { /* '77' */
            return ""
        }
        case 78: { /* '78' */
            return "Hardware Type"
        }
        case 79: { /* '79' */
            return ""
        }
        case 8: { /* '8' */
            return "Service Control"
        }
        case 80: { /* '80' */
            return ""
        }
        case 81: { /* '81' */
            return ""
        }
        case 82: { /* '82' */
            return "RF Domain Address"
        }
        case 83: { /* '83' */
            return ""
        }
        case 84: { /* '84' */
            return ""
        }
        case 85: { /* '85' */
            return ""
        }
        case 86: { /* '86' */
            return ""
        }
        case 9: { /* '9' */
            return "Firmware Revision"
        }
        case 91: { /* '91' */
            return ""
        }
        case 92: { /* '92' */
            return ""
        }
        case 93: { /* '93' */
            return ""
        }
        case 94: { /* '94' */
            return ""
        }
        case 95: { /* '95' */
            return ""
        }
        case 96: { /* '96' */
            return ""
        }
        case 97: { /* '97' */
            return ""
        }
        default: {
            return ""
        }
    }
}
func KnxObjectPropertiesValueOf(value uint8) KnxObjectProperties {
    switch value {
        case 1:
            return KnxObjectProperties_PID_G_1
        case 10:
            return KnxObjectProperties_PID_G_10
        case 101:
            return KnxObjectProperties_PID_417_101
        case 102:
            return KnxObjectProperties_PID_417_102
        case 103:
            return KnxObjectProperties_PID_417_103
        case 104:
            return KnxObjectProperties_PID_417_104
        case 105:
            return KnxObjectProperties_PID_417_105
        case 106:
            return KnxObjectProperties_PID_417_106
        case 107:
            return KnxObjectProperties_PID_417_107
        case 108:
            return KnxObjectProperties_PID_417_108
        case 109:
            return KnxObjectProperties_PID_417_109
        case 11:
            return KnxObjectProperties_PID_G_11
        case 110:
            return KnxObjectProperties_PID_409_110
        case 111:
            return KnxObjectProperties_PID_409_111
        case 112:
            return KnxObjectProperties_PID_6_112
        case 113:
            return KnxObjectProperties_PID_417_113
        case 114:
            return KnxObjectProperties_PID_417_114
        case 115:
            return KnxObjectProperties_PID_417_115
        case 116:
            return KnxObjectProperties_PID_417_116
        case 117:
            return KnxObjectProperties_PID_417_117
        case 118:
            return KnxObjectProperties_PID_417_118
        case 119:
            return KnxObjectProperties_PID_417_119
        case 12:
            return KnxObjectProperties_PID_G_12
        case 120:
            return KnxObjectProperties_PID_417_120
        case 121:
            return KnxObjectProperties_PID_418_121
        case 122:
            return KnxObjectProperties_PID_418_122
        case 123:
            return KnxObjectProperties_PID_418_123
        case 124:
            return KnxObjectProperties_PID_418_124
        case 125:
            return KnxObjectProperties_PID_418_125
        case 126:
            return KnxObjectProperties_PID_418_126
        case 127:
            return KnxObjectProperties_PID_418_127
        case 128:
            return KnxObjectProperties_PID_418_128
        case 129:
            return KnxObjectProperties_PID_418_129
        case 13:
            return KnxObjectProperties_PID_G_13
        case 130:
            return KnxObjectProperties_PID_418_130
        case 131:
            return KnxObjectProperties_PID_418_131
        case 132:
            return KnxObjectProperties_PID_418_132
        case 133:
            return KnxObjectProperties_PID_418_133
        case 134:
            return KnxObjectProperties_PID_418_134
        case 14:
            return KnxObjectProperties_PID_G_14
        case 15:
            return KnxObjectProperties_PID_G_15
        case 16:
            return KnxObjectProperties_PID_G_16
        case 17:
            return KnxObjectProperties_PID_G_17
        case 18:
            return KnxObjectProperties_PID_G_18
        case 19:
            return KnxObjectProperties_PID_G_19
        case 2:
            return KnxObjectProperties_PID_G_2
        case 20:
            return KnxObjectProperties_PID_G_20
        case 21:
            return KnxObjectProperties_PID_G_21
        case 22:
            return KnxObjectProperties_PID_G_22
        case 23:
            return KnxObjectProperties_PID_G_23
        case 24:
            return KnxObjectProperties_PID_G_24
        case 25:
            return KnxObjectProperties_PID_G_25
        case 26:
            return KnxObjectProperties_PID_G_26
        case 27:
            return KnxObjectProperties_PID_G_27
        case 28:
            return KnxObjectProperties_PID_G_28
        case 29:
            return KnxObjectProperties_PID_G_29
        case 3:
            return KnxObjectProperties_PID_G_3
        case 30:
            return KnxObjectProperties_PID_G_30
        case 4:
            return KnxObjectProperties_PID_G_4
        case 5:
            return KnxObjectProperties_PID_G_5
        case 51:
            return KnxObjectProperties_PID_0_51
        case 52:
            return KnxObjectProperties_PID_0_52
        case 53:
            return KnxObjectProperties_PID_0_53
        case 54:
            return KnxObjectProperties_PID_0_54
        case 55:
            return KnxObjectProperties_PID_0_55
        case 56:
            return KnxObjectProperties_PID_0_56
        case 57:
            return KnxObjectProperties_PID_0_57
        case 58:
            return KnxObjectProperties_PID_0_58
        case 59:
            return KnxObjectProperties_PID_0_59
        case 6:
            return KnxObjectProperties_PID_G_6
        case 60:
            return KnxObjectProperties_PID_0_60
        case 61:
            return KnxObjectProperties_PID_0_61
        case 62:
            return KnxObjectProperties_PID_0_62
        case 63:
            return KnxObjectProperties_PID_0_63
        case 64:
            return KnxObjectProperties_PID_0_64
        case 65:
            return KnxObjectProperties_PID_0_65
        case 66:
            return KnxObjectProperties_PID_0_66
        case 67:
            return KnxObjectProperties_PID_0_67
        case 68:
            return KnxObjectProperties_PID_0_68
        case 69:
            return KnxObjectProperties_PID_0_69
        case 7:
            return KnxObjectProperties_PID_G_7
        case 70:
            return KnxObjectProperties_PID_0_70
        case 71:
            return KnxObjectProperties_PID_0_71
        case 72:
            return KnxObjectProperties_PID_0_72
        case 73:
            return KnxObjectProperties_PID_0_73
        case 74:
            return KnxObjectProperties_PID_0_74
        case 75:
            return KnxObjectProperties_PID_0_75
        case 76:
            return KnxObjectProperties_PID_0_76
        case 77:
            return KnxObjectProperties_PID_0_77
        case 78:
            return KnxObjectProperties_PID_0_78
        case 79:
            return KnxObjectProperties_PID_0_79
        case 8:
            return KnxObjectProperties_PID_G_8
        case 80:
            return KnxObjectProperties_PID_0_80
        case 81:
            return KnxObjectProperties_PID_0_81
        case 82:
            return KnxObjectProperties_PID_0_82
        case 83:
            return KnxObjectProperties_PID_0_83
        case 84:
            return KnxObjectProperties_PID_0_84
        case 85:
            return KnxObjectProperties_PID_0_85
        case 86:
            return KnxObjectProperties_PID_0_86
        case 9:
            return KnxObjectProperties_PID_G_9
        case 91:
            return KnxObjectProperties_PID_11_91
        case 92:
            return KnxObjectProperties_PID_11_92
        case 93:
            return KnxObjectProperties_PID_11_93
        case 94:
            return KnxObjectProperties_PID_11_94
        case 95:
            return KnxObjectProperties_PID_11_95
        case 96:
            return KnxObjectProperties_PID_11_96
        case 97:
            return KnxObjectProperties_PID_11_97
    }
    return 0
}

func CastKnxObjectProperties(structType interface{}) KnxObjectProperties {
    castFunc := func(typ interface{}) KnxObjectProperties {
        if sKnxObjectProperties, ok := typ.(KnxObjectProperties); ok {
            return sKnxObjectProperties
        }
        return 0
    }
    return castFunc(structType)
}

func (m KnxObjectProperties) LengthInBits() uint16 {
    return 8
}

func (m KnxObjectProperties) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KnxObjectPropertiesParse(io *utils.ReadBuffer) (KnxObjectProperties, error) {
    val, err := io.ReadUint8(8)
    if err != nil {
        return 0, nil
    }
    return KnxObjectPropertiesValueOf(val), nil
}

func (e KnxObjectProperties) Serialize(io utils.WriteBuffer) error {
    err := io.WriteUint8(8, uint8(e))
    return err
}

func (e KnxObjectProperties) String() string {
    switch e {
    case KnxObjectProperties_PID_G_1:
        return "PID_G_1"
    case KnxObjectProperties_PID_G_10:
        return "PID_G_10"
    case KnxObjectProperties_PID_417_101:
        return "PID_417_101"
    case KnxObjectProperties_PID_417_102:
        return "PID_417_102"
    case KnxObjectProperties_PID_417_103:
        return "PID_417_103"
    case KnxObjectProperties_PID_417_104:
        return "PID_417_104"
    case KnxObjectProperties_PID_417_105:
        return "PID_417_105"
    case KnxObjectProperties_PID_417_106:
        return "PID_417_106"
    case KnxObjectProperties_PID_417_107:
        return "PID_417_107"
    case KnxObjectProperties_PID_417_108:
        return "PID_417_108"
    case KnxObjectProperties_PID_417_109:
        return "PID_417_109"
    case KnxObjectProperties_PID_G_11:
        return "PID_G_11"
    case KnxObjectProperties_PID_409_110:
        return "PID_409_110"
    case KnxObjectProperties_PID_409_111:
        return "PID_409_111"
    case KnxObjectProperties_PID_6_112:
        return "PID_6_112"
    case KnxObjectProperties_PID_417_113:
        return "PID_417_113"
    case KnxObjectProperties_PID_417_114:
        return "PID_417_114"
    case KnxObjectProperties_PID_417_115:
        return "PID_417_115"
    case KnxObjectProperties_PID_417_116:
        return "PID_417_116"
    case KnxObjectProperties_PID_417_117:
        return "PID_417_117"
    case KnxObjectProperties_PID_417_118:
        return "PID_417_118"
    case KnxObjectProperties_PID_417_119:
        return "PID_417_119"
    case KnxObjectProperties_PID_G_12:
        return "PID_G_12"
    case KnxObjectProperties_PID_417_120:
        return "PID_417_120"
    case KnxObjectProperties_PID_418_121:
        return "PID_418_121"
    case KnxObjectProperties_PID_418_122:
        return "PID_418_122"
    case KnxObjectProperties_PID_418_123:
        return "PID_418_123"
    case KnxObjectProperties_PID_418_124:
        return "PID_418_124"
    case KnxObjectProperties_PID_418_125:
        return "PID_418_125"
    case KnxObjectProperties_PID_418_126:
        return "PID_418_126"
    case KnxObjectProperties_PID_418_127:
        return "PID_418_127"
    case KnxObjectProperties_PID_418_128:
        return "PID_418_128"
    case KnxObjectProperties_PID_418_129:
        return "PID_418_129"
    case KnxObjectProperties_PID_G_13:
        return "PID_G_13"
    case KnxObjectProperties_PID_418_130:
        return "PID_418_130"
    case KnxObjectProperties_PID_418_131:
        return "PID_418_131"
    case KnxObjectProperties_PID_418_132:
        return "PID_418_132"
    case KnxObjectProperties_PID_418_133:
        return "PID_418_133"
    case KnxObjectProperties_PID_418_134:
        return "PID_418_134"
    case KnxObjectProperties_PID_G_14:
        return "PID_G_14"
    case KnxObjectProperties_PID_G_15:
        return "PID_G_15"
    case KnxObjectProperties_PID_G_16:
        return "PID_G_16"
    case KnxObjectProperties_PID_G_17:
        return "PID_G_17"
    case KnxObjectProperties_PID_G_18:
        return "PID_G_18"
    case KnxObjectProperties_PID_G_19:
        return "PID_G_19"
    case KnxObjectProperties_PID_G_2:
        return "PID_G_2"
    case KnxObjectProperties_PID_G_20:
        return "PID_G_20"
    case KnxObjectProperties_PID_G_21:
        return "PID_G_21"
    case KnxObjectProperties_PID_G_22:
        return "PID_G_22"
    case KnxObjectProperties_PID_G_23:
        return "PID_G_23"
    case KnxObjectProperties_PID_G_24:
        return "PID_G_24"
    case KnxObjectProperties_PID_G_25:
        return "PID_G_25"
    case KnxObjectProperties_PID_G_26:
        return "PID_G_26"
    case KnxObjectProperties_PID_G_27:
        return "PID_G_27"
    case KnxObjectProperties_PID_G_28:
        return "PID_G_28"
    case KnxObjectProperties_PID_G_29:
        return "PID_G_29"
    case KnxObjectProperties_PID_G_3:
        return "PID_G_3"
    case KnxObjectProperties_PID_G_30:
        return "PID_G_30"
    case KnxObjectProperties_PID_G_4:
        return "PID_G_4"
    case KnxObjectProperties_PID_G_5:
        return "PID_G_5"
    case KnxObjectProperties_PID_0_51:
        return "PID_0_51"
    case KnxObjectProperties_PID_0_52:
        return "PID_0_52"
    case KnxObjectProperties_PID_0_53:
        return "PID_0_53"
    case KnxObjectProperties_PID_0_54:
        return "PID_0_54"
    case KnxObjectProperties_PID_0_55:
        return "PID_0_55"
    case KnxObjectProperties_PID_0_56:
        return "PID_0_56"
    case KnxObjectProperties_PID_0_57:
        return "PID_0_57"
    case KnxObjectProperties_PID_0_58:
        return "PID_0_58"
    case KnxObjectProperties_PID_0_59:
        return "PID_0_59"
    case KnxObjectProperties_PID_G_6:
        return "PID_G_6"
    case KnxObjectProperties_PID_0_60:
        return "PID_0_60"
    case KnxObjectProperties_PID_0_61:
        return "PID_0_61"
    case KnxObjectProperties_PID_0_62:
        return "PID_0_62"
    case KnxObjectProperties_PID_0_63:
        return "PID_0_63"
    case KnxObjectProperties_PID_0_64:
        return "PID_0_64"
    case KnxObjectProperties_PID_0_65:
        return "PID_0_65"
    case KnxObjectProperties_PID_0_66:
        return "PID_0_66"
    case KnxObjectProperties_PID_0_67:
        return "PID_0_67"
    case KnxObjectProperties_PID_0_68:
        return "PID_0_68"
    case KnxObjectProperties_PID_0_69:
        return "PID_0_69"
    case KnxObjectProperties_PID_G_7:
        return "PID_G_7"
    case KnxObjectProperties_PID_0_70:
        return "PID_0_70"
    case KnxObjectProperties_PID_0_71:
        return "PID_0_71"
    case KnxObjectProperties_PID_0_72:
        return "PID_0_72"
    case KnxObjectProperties_PID_0_73:
        return "PID_0_73"
    case KnxObjectProperties_PID_0_74:
        return "PID_0_74"
    case KnxObjectProperties_PID_0_75:
        return "PID_0_75"
    case KnxObjectProperties_PID_0_76:
        return "PID_0_76"
    case KnxObjectProperties_PID_0_77:
        return "PID_0_77"
    case KnxObjectProperties_PID_0_78:
        return "PID_0_78"
    case KnxObjectProperties_PID_0_79:
        return "PID_0_79"
    case KnxObjectProperties_PID_G_8:
        return "PID_G_8"
    case KnxObjectProperties_PID_0_80:
        return "PID_0_80"
    case KnxObjectProperties_PID_0_81:
        return "PID_0_81"
    case KnxObjectProperties_PID_0_82:
        return "PID_0_82"
    case KnxObjectProperties_PID_0_83:
        return "PID_0_83"
    case KnxObjectProperties_PID_0_84:
        return "PID_0_84"
    case KnxObjectProperties_PID_0_85:
        return "PID_0_85"
    case KnxObjectProperties_PID_0_86:
        return "PID_0_86"
    case KnxObjectProperties_PID_G_9:
        return "PID_G_9"
    case KnxObjectProperties_PID_11_91:
        return "PID_11_91"
    case KnxObjectProperties_PID_11_92:
        return "PID_11_92"
    case KnxObjectProperties_PID_11_93:
        return "PID_11_93"
    case KnxObjectProperties_PID_11_94:
        return "PID_11_94"
    case KnxObjectProperties_PID_11_95:
        return "PID_11_95"
    case KnxObjectProperties_PID_11_96:
        return "PID_11_96"
    case KnxObjectProperties_PID_11_97:
        return "PID_11_97"
    }
    return ""
}
