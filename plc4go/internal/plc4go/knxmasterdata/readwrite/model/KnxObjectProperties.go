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
    DataTypeId() string
    Text() string
    Serialize(io utils.WriteBuffer) error
}

const(
    KnxObjectProperties_PID_OBJECT_TYPE KnxObjectProperties = 1
    KnxObjectProperties_PID_OBJECT_NAME KnxObjectProperties = 2
    KnxObjectProperties_PID_SEMAPHOR KnxObjectProperties = 3
    KnxObjectProperties_PID_GROUP_OBJECT_REFERENCE KnxObjectProperties = 4
    KnxObjectProperties_PID_LOAD_STATE_CONTROL KnxObjectProperties = 5
    KnxObjectProperties_PID_RUN_STATE_CONTROL KnxObjectProperties = 6
    KnxObjectProperties_PID_TABLE_REFERENCE KnxObjectProperties = 7
    KnxObjectProperties_PID_SERVICE_CONTROL KnxObjectProperties = 8
    KnxObjectProperties_PID_FIRMWARE_REVISION KnxObjectProperties = 9
    KnxObjectProperties_PID_SERVICES_SUPPORTED KnxObjectProperties = 10
    KnxObjectProperties_PID_SERIAL_NUMBER KnxObjectProperties = 11
    KnxObjectProperties_PID_MANUFACTURER_ID KnxObjectProperties = 12
    KnxObjectProperties_PID_PROGRAM_VERSION KnxObjectProperties = 13
    KnxObjectProperties_PID_DEVICE_CONTROL KnxObjectProperties = 14
    KnxObjectProperties_PID_ORDER_INFO KnxObjectProperties = 15
    KnxObjectProperties_PID_PEI_TYPE KnxObjectProperties = 16
    KnxObjectProperties_PID_PORT_CONFIGURATION KnxObjectProperties = 17
    KnxObjectProperties_PID_POLL_GROUP_SETTINGS KnxObjectProperties = 18
    KnxObjectProperties_PID_MANUFACTURER_DATA KnxObjectProperties = 19
    KnxObjectProperties_PID_ENABLE KnxObjectProperties = 20
    KnxObjectProperties_PID_DESCRIPTION KnxObjectProperties = 21
    KnxObjectProperties_PID_FILE KnxObjectProperties = 22
    KnxObjectProperties_PID_TABLE KnxObjectProperties = 23
    KnxObjectProperties_PID_ENROL KnxObjectProperties = 24
    KnxObjectProperties_PID_VERSION KnxObjectProperties = 25
    KnxObjectProperties_PID_GROUP_OBJECT_LINK KnxObjectProperties = 26
    KnxObjectProperties_PID_MCB_TABLE KnxObjectProperties = 27
    KnxObjectProperties_PID_ERROR_CODE KnxObjectProperties = 28
    KnxObjectProperties_PID_OBJECT_INDEX KnxObjectProperties = 29
    KnxObjectProperties_PID_DOWNLOAD_COUNTER KnxObjectProperties = 30
    KnxObjectProperties_PID_ROUTING_COUNT KnxObjectProperties = 51
    KnxObjectProperties_PID_MAX_RETRY_COUNT KnxObjectProperties = 52
    KnxObjectProperties_PID_ERROR_FLAGS KnxObjectProperties = 53
    KnxObjectProperties_PID_PROGMODE KnxObjectProperties = 54
    KnxObjectProperties_PID_PRODUCT_ID KnxObjectProperties = 55
    KnxObjectProperties_PID_MAX_APDULENGTH KnxObjectProperties = 56
    KnxObjectProperties_PID_SUBNET_ADDR KnxObjectProperties = 57
    KnxObjectProperties_PID_DEVICE_ADDR KnxObjectProperties = 58
    KnxObjectProperties_PID_PB_CONFIG KnxObjectProperties = 59
    KnxObjectProperties_PID_ADDR_REPORT KnxObjectProperties = 60
    KnxObjectProperties_PID_ADDR_CHECK KnxObjectProperties = 61
    KnxObjectProperties_PID_OBJECT_VALUE KnxObjectProperties = 62
    KnxObjectProperties_PID_OBJECTLINK KnxObjectProperties = 63
    KnxObjectProperties_PID_APPLICATION KnxObjectProperties = 64
    KnxObjectProperties_PID_PARAMETER KnxObjectProperties = 65
    KnxObjectProperties_PID_OBJECTADDRESS KnxObjectProperties = 66
    KnxObjectProperties_PID_PSU_TYPE KnxObjectProperties = 67
    KnxObjectProperties_PID_PSU_STATUS KnxObjectProperties = 68
    KnxObjectProperties_PID_PSU_ENABLE KnxObjectProperties = 69
    KnxObjectProperties_PID_DOMAIN_ADDRESS KnxObjectProperties = 70
    KnxObjectProperties_PID_IO_LIST KnxObjectProperties = 71
    KnxObjectProperties_PID_MGT_DESCRIPTOR_01 KnxObjectProperties = 72
    KnxObjectProperties_PID_PL110_PARAM KnxObjectProperties = 73
    KnxObjectProperties_PID_RF_REPEAT_COUNTER KnxObjectProperties = 74
    KnxObjectProperties_PID_RECEIVE_BLOCK_TABLE KnxObjectProperties = 75
    KnxObjectProperties_PID_RANDOM_PAUSE_TABLE KnxObjectProperties = 76
    KnxObjectProperties_PID_RECEIVE_BLOCK_NR KnxObjectProperties = 77
    KnxObjectProperties_PID_HARDWARE_TYPE KnxObjectProperties = 78
    KnxObjectProperties_PID_RETRANSMITTER_NUMBER KnxObjectProperties = 79
    KnxObjectProperties_PID_SERIAL_NR_TABLE KnxObjectProperties = 80
    KnxObjectProperties_PID_BIBATMASTER_ADDRESS KnxObjectProperties = 81
    KnxObjectProperties_PID_RF_DOMAIN_ADDRESS KnxObjectProperties = 82
    KnxObjectProperties_PID_DEVICE_DESCRIPTOR KnxObjectProperties = 83
    KnxObjectProperties_PID_METERING_FILTER_TABLE KnxObjectProperties = 84
    KnxObjectProperties_PID_GROUP_TELEGR_RATE_LIMIT_TIME_BASE KnxObjectProperties = 85
    KnxObjectProperties_PID_GROUP_TELEGR_RATE_LIMIT_NO_OF_TELEGR KnxObjectProperties = 86
    KnxObjectProperties_PID_GRPOBJTABLE KnxObjectProperties = 51
    KnxObjectProperties_PID_EXT_GRPOBJREFERENCE KnxObjectProperties = 52
    KnxObjectProperties_PID_LINE_STATUS KnxObjectProperties = 51
    KnxObjectProperties_PID_MAIN_LCCONFIG KnxObjectProperties = 52
    KnxObjectProperties_PID_SUB_LCCONFIG KnxObjectProperties = 53
    KnxObjectProperties_PID_MAIN_LCGRPCONFIG KnxObjectProperties = 54
    KnxObjectProperties_PID_SUB_LCGRPCONFIG KnxObjectProperties = 55
    KnxObjectProperties_PID_ROUTETABLE_CONTROL KnxObjectProperties = 56
    KnxObjectProperties_PID_COUPL_SERV_CONTROL KnxObjectProperties = 57
    KnxObjectProperties_PID_MAX_ROUTER_APDU_LENGTH KnxObjectProperties = 58
    KnxObjectProperties_PID_MEDIUM KnxObjectProperties = 63
    KnxObjectProperties_PID_FILTER_TABLE_USE KnxObjectProperties = 67
    KnxObjectProperties_PID_RF_ENABLE_SBC KnxObjectProperties = 112
    KnxObjectProperties_PID_PROJECT_INSTALLATION_ID KnxObjectProperties = 51
    KnxObjectProperties_PID_KNX_INDIVIDUAL_ADDRESS KnxObjectProperties = 52
    KnxObjectProperties_PID_ADDITIONAL_INDIVIDUAL_ADDRESSES KnxObjectProperties = 53
    KnxObjectProperties_PID_CURRENT_IP_ASSIGNMENT_METHOD KnxObjectProperties = 54
    KnxObjectProperties_PID_IP_ASSIGNMENT_METHOD KnxObjectProperties = 55
    KnxObjectProperties_PID_IP_CAPABILITIES KnxObjectProperties = 56
    KnxObjectProperties_PID_CURRENT_IP_ADDRESS KnxObjectProperties = 57
    KnxObjectProperties_PID_CURRENT_SUBNET_MASK KnxObjectProperties = 58
    KnxObjectProperties_PID_CURRENT_DEFAULT_GATEWAY KnxObjectProperties = 59
    KnxObjectProperties_PID_IP_ADDRESS KnxObjectProperties = 60
    KnxObjectProperties_PID_SUBNET_MASK KnxObjectProperties = 61
    KnxObjectProperties_PID_DEFAULT_GATEWAY KnxObjectProperties = 62
    KnxObjectProperties_PID_DHCP_BOOTP_SERVER KnxObjectProperties = 63
    KnxObjectProperties_PID_MAC_ADDRESS KnxObjectProperties = 64
    KnxObjectProperties_PID_SYSTEM_SETUP_MULTICAST_ADDRESS KnxObjectProperties = 65
    KnxObjectProperties_PID_ROUTING_MULTICAST_ADDRESS KnxObjectProperties = 66
    KnxObjectProperties_PID_TTL KnxObjectProperties = 67
    KnxObjectProperties_PID_KNXNETIP_DEVICE_CAPABILITIES KnxObjectProperties = 68
    KnxObjectProperties_PID_KNXNETIP_DEVICE_STATE KnxObjectProperties = 69
    KnxObjectProperties_PID_KNXNETIP_ROUTING_CAPABILITIES KnxObjectProperties = 70
    KnxObjectProperties_PID_PRIORITY_FIFO_ENABLED KnxObjectProperties = 71
    KnxObjectProperties_PID_QUEUE_OVERFLOW_TO_IP KnxObjectProperties = 72
    KnxObjectProperties_PID_QUEUE_OVERFLOW_TO_KNX KnxObjectProperties = 73
    KnxObjectProperties_PID_MSG_TRANSMIT_TO_IP KnxObjectProperties = 74
    KnxObjectProperties_PID_MSG_TRANSMIT_TO_KNX KnxObjectProperties = 75
    KnxObjectProperties_PID_FRIENDLY_NAME KnxObjectProperties = 76
    KnxObjectProperties_PID_BACKBONE_KEY KnxObjectProperties = 91
    KnxObjectProperties_PID_DEVICE_AUTHENTICATION_CODE KnxObjectProperties = 92
    KnxObjectProperties_PID_PASSWORD_HASHES KnxObjectProperties = 93
    KnxObjectProperties_PID_SECURED_SERVICE_FAMILIES KnxObjectProperties = 94
    KnxObjectProperties_PID_MULTICAST_LATENCY_TOLERANCE KnxObjectProperties = 95
    KnxObjectProperties_PID_SYNC_LATENCY_FRACTION KnxObjectProperties = 96
    KnxObjectProperties_PID_TUNNELLING_USERS KnxObjectProperties = 97
    KnxObjectProperties_PID_SECURITY_MODE KnxObjectProperties = 51
    KnxObjectProperties_PID_P2P_KEY_TABLE KnxObjectProperties = 52
    KnxObjectProperties_PID_GRP_KEY_TABLE KnxObjectProperties = 53
    KnxObjectProperties_PID_SECURITY_INDIVIDUAL_ADDRESS_TABLE KnxObjectProperties = 54
    KnxObjectProperties_PID_SECURITY_FAILURES_LOG KnxObjectProperties = 55
    KnxObjectProperties_PID_SKI_TOOL KnxObjectProperties = 56
    KnxObjectProperties_PID_SECURITY_REPORT KnxObjectProperties = 57
    KnxObjectProperties_PID_SECURITY_REPORT_CONTROL KnxObjectProperties = 58
    KnxObjectProperties_PID_SEQUENCE_NUMBER_SENDING KnxObjectProperties = 59
    KnxObjectProperties_PID_ZONE_KEYS_TABLE KnxObjectProperties = 60
    KnxObjectProperties_PID_GO_SECURITY_FLAGS KnxObjectProperties = 61
    KnxObjectProperties_PID_RF_MULTI_TYPE KnxObjectProperties = 51
    KnxObjectProperties_PID_RF_DOMAIN_ADDRESS KnxObjectProperties = 56
    KnxObjectProperties_PID_RF_RETRANSMITTER KnxObjectProperties = 57
    KnxObjectProperties_PID_SECURITY_REPORT_CONTROL KnxObjectProperties = 58
    KnxObjectProperties_PID_RF_FILTERING_MODE_SELECT KnxObjectProperties = 59
    KnxObjectProperties_PID_RF_BIDIR_TIMEOUT KnxObjectProperties = 60
    KnxObjectProperties_PID_RF_DIAG_SA_FILTER_TABLE KnxObjectProperties = 61
    KnxObjectProperties_PID_RF_DIAG_QUALITY_TABLE KnxObjectProperties = 62
    KnxObjectProperties_PID_RF_DIAG_PROBE KnxObjectProperties = 63
    KnxObjectProperties_PID_CHANGE_OF_VALUE KnxObjectProperties = 110
    KnxObjectProperties_PID_REPETITION_TIME KnxObjectProperties = 111
    KnxObjectProperties_PID_CHANGE_OF_VALUE KnxObjectProperties = 110
    KnxObjectProperties_PID_REPETITION_TIME KnxObjectProperties = 111
    KnxObjectProperties_PID_ON_DELAY KnxObjectProperties = 101
    KnxObjectProperties_PID_OFF_DELAY KnxObjectProperties = 102
    KnxObjectProperties_PID_TIMED_ON_DURATION KnxObjectProperties = 103
    KnxObjectProperties_PID_PREWARNING_DURATION KnxObjectProperties = 104
    KnxObjectProperties_PID_TRANSMISSION_CYCLE_TIME KnxObjectProperties = 105
    KnxObjectProperties_PID_BUS_POWER_UP_MESSAGE_DELAY KnxObjectProperties = 106
    KnxObjectProperties_PID_BEHAVIOUR_AT_LOCKING KnxObjectProperties = 107
    KnxObjectProperties_PID_BEHAVIOUR_AT_UNLOCKING KnxObjectProperties = 108
    KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_UP KnxObjectProperties = 109
    KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_DOWN KnxObjectProperties = 110
    KnxObjectProperties_PID_INVERT_OUTPUT_STATE KnxObjectProperties = 111
    KnxObjectProperties_PID_TIMED_ON_RETRIGGER_FUNCTION KnxObjectProperties = 112
    KnxObjectProperties_PID_MANUAL_OFF_ENABLE KnxObjectProperties = 113
    KnxObjectProperties_PID_INVERT_LOCK_DEVICE KnxObjectProperties = 114
    KnxObjectProperties_PID_LOCK_STATE KnxObjectProperties = 115
    KnxObjectProperties_PID_UNLOCK_STATE KnxObjectProperties = 116
    KnxObjectProperties_PID_STATE_FOR_SCENE_NUMBER KnxObjectProperties = 117
    KnxObjectProperties_PID_STORAGE_FUNCTION_FOR_SCENE KnxObjectProperties = 118
    KnxObjectProperties_PID_BUS_POWER_UP_STATE KnxObjectProperties = 119
    KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_UP KnxObjectProperties = 120
    KnxObjectProperties_PID_ON_DELAY KnxObjectProperties = 101
    KnxObjectProperties_PID_OFF_DELAY KnxObjectProperties = 102
    KnxObjectProperties_PID_SWITCH_OFF_BRIGHTNESS_DELAY_TIME KnxObjectProperties = 103
    KnxObjectProperties_PID_TIMED_ON_DURATION KnxObjectProperties = 104
    KnxObjectProperties_PID_PREWARNING_DURATION KnxObjectProperties = 105
    KnxObjectProperties_PID_TRANSMISSION_CYCLE_TIME KnxObjectProperties = 106
    KnxObjectProperties_PID_BUS_POWER_UP_MESSAGE_DELAY KnxObjectProperties = 107
    KnxObjectProperties_PID_DIMMING_SPEED KnxObjectProperties = 108
    KnxObjectProperties_PID_DIMMING_STEP_TIME KnxObjectProperties = 109
    KnxObjectProperties_PID_DIMMING_SPEED_FOR_SWITCH_ON_SET_VALUE KnxObjectProperties = 110
    KnxObjectProperties_PID_DIMMING_SPEED_FOR_SWITCH_OFF KnxObjectProperties = 111
    KnxObjectProperties_PID_DIMMING_STEP_TIME_FOR_SWITCH_ON_SET_VALUE KnxObjectProperties = 112
    KnxObjectProperties_PID_DIMMING_STEP_TIME_FOR_SWITCH_OFF KnxObjectProperties = 113
    KnxObjectProperties_PID_SWITCFH_OFF_BRIGHTNESS KnxObjectProperties = 114
    KnxObjectProperties_PID_MINIMUM_SET_VALUE KnxObjectProperties = 115
    KnxObjectProperties_PID_MAXIMUM_SET_VALUE KnxObjectProperties = 116
    KnxObjectProperties_PID_SWITCH_ON_SET_VALUE KnxObjectProperties = 117
    KnxObjectProperties_PID_DIMM_MODE_SELECTION KnxObjectProperties = 118
    KnxObjectProperties_PID_RELATIV_OFF_ENABLE KnxObjectProperties = 119
    KnxObjectProperties_PID_MEMORY_FUNCTION KnxObjectProperties = 120
    KnxObjectProperties_PID_TIMED_ON_RETRIGGER_FUNCTION KnxObjectProperties = 121
    KnxObjectProperties_PID_MANUAL_OFF_ENABLE KnxObjectProperties = 122
    KnxObjectProperties_PID_INVERT_LOCK_DEVICE KnxObjectProperties = 123
    KnxObjectProperties_PID_BEHAVIOUR_AT_LOCKING KnxObjectProperties = 124
    KnxObjectProperties_PID_BEHAVIOUR_AT_UNLOCKING KnxObjectProperties = 125
    KnxObjectProperties_PID_LOCK_SETVALUE KnxObjectProperties = 126
    KnxObjectProperties_PID_UNLOCK_SETVALUE KnxObjectProperties = 127
    KnxObjectProperties_PID_BIGHTNESS_FOR_SCENE KnxObjectProperties = 128
    KnxObjectProperties_PID_STORAGE_FUNCTION_FOR_SCENE KnxObjectProperties = 129
    KnxObjectProperties_PID_DELTA_DIMMING_VALUE KnxObjectProperties = 130
    KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_UP KnxObjectProperties = 131
    KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_UP_SET_VALUE KnxObjectProperties = 132
    KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_DOWN KnxObjectProperties = 133
    KnxObjectProperties_PID_BUS_POWER_DOWN_SET_VALUE KnxObjectProperties = 134
    KnxObjectProperties_PID_ON_OFF_ACTION KnxObjectProperties = 51
    KnxObjectProperties_PID_ENABLE_TOGGLE_MODE KnxObjectProperties = 52
    KnxObjectProperties_PID_ABSOLUTE_SETVALUE KnxObjectProperties = 53
    KnxObjectProperties_PID_ON_OFF_ACTION KnxObjectProperties = 51
    KnxObjectProperties_PID_ENABLE_TOGGLE_MODE KnxObjectProperties = 52
    KnxObjectProperties_PID_REVERSION_PAUSE_TIME KnxObjectProperties = 51
    KnxObjectProperties_PID_MOVE_UP_DOWN_TIME KnxObjectProperties = 52
    KnxObjectProperties_PID_SLAT_STEP_TIME KnxObjectProperties = 53
    KnxObjectProperties_PID_MOVE_PRESET_POSITION_TIME KnxObjectProperties = 54
    KnxObjectProperties_PID_MOVE_TO_PRESET_POSITION_IN_PERCENT KnxObjectProperties = 55
    KnxObjectProperties_PID_MOVE_TO_PRESET_POSITION_LENGTH KnxObjectProperties = 57
    KnxObjectProperties_PID_PRESET_SLAT_POSITION_PERCENT KnxObjectProperties = 58
    KnxObjectProperties_PID_PRESET_SLAT_POSITION_ANGLE KnxObjectProperties = 60
    KnxObjectProperties_PID_REACTION_WIND_ALARM KnxObjectProperties = 61
    KnxObjectProperties_PID_HEARTBEAT_WIND_ALARM KnxObjectProperties = 62
    KnxObjectProperties_PID_REACTION_ON_RAIN_ALARM KnxObjectProperties = 63
    KnxObjectProperties_PID_HEARTBEAT_RAIN_ALARM KnxObjectProperties = 64
    KnxObjectProperties_PID_REACTION_FROST_ALARM KnxObjectProperties = 65
    KnxObjectProperties_PID_HEARTBEAT_FROST_ALARM KnxObjectProperties = 66
    KnxObjectProperties_PID_MAX_SLAT_MOVE_TIME KnxObjectProperties = 67
    KnxObjectProperties_PID_ENABLE_BLINDS_MODE KnxObjectProperties = 68
    KnxObjectProperties_PID_STORAGE_FUNCTIONS_FOR_SCENE KnxObjectProperties = 69
    KnxObjectProperties_PID_ENABLE_BLINDS_MODE KnxObjectProperties = 51
    KnxObjectProperties_PID_UP_DOWN_ACTION KnxObjectProperties = 52
    KnxObjectProperties_PID_ENABLE_TOGGLE_MODE KnxObjectProperties = 53
)


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
            return """"
        }
        case 102: { /* '102' */
            return """"
        }
        case 103: { /* '103' */
            return """"
        }
        case 104: { /* '104' */
            return """"
        }
        case 105: { /* '105' */
            return """"
        }
        case 106: { /* '106' */
            return """"
        }
        case 107: { /* '107' */
            return """"
        }
        case 108: { /* '108' */
            return """"
        }
        case 109: { /* '109' */
            return """"
        }
        case 11: { /* '11' */
            return "KNX Serial Number"
        }
        case 110: { /* '110' */
            return """"
        }
        case 111: { /* '111' */
            return """"
        }
        case 112: { /* '112' */
            return """"
        }
        case 113: { /* '113' */
            return """"
        }
        case 114: { /* '114' */
            return """"
        }
        case 115: { /* '115' */
            return """"
        }
        case 116: { /* '116' */
            return """"
        }
        case 117: { /* '117' */
            return """"
        }
        case 118: { /* '118' */
            return """"
        }
        case 119: { /* '119' */
            return """"
        }
        case 12: { /* '12' */
            return "Manufacturer Identifier"
        }
        case 120: { /* '120' */
            return """"
        }
        case 121: { /* '121' */
            return """"
        }
        case 122: { /* '122' */
            return """"
        }
        case 123: { /* '123' */
            return """"
        }
        case 124: { /* '124' */
            return """"
        }
        case 125: { /* '125' */
            return """"
        }
        case 126: { /* '126' */
            return """"
        }
        case 127: { /* '127' */
            return """"
        }
        case 128: { /* '128' */
            return """"
        }
        case 129: { /* '129' */
            return """"
        }
        case 13: { /* '13' */
            return "Application Version"
        }
        case 130: { /* '130' */
            return """"
        }
        case 131: { /* '131' */
            return """"
        }
        case 132: { /* '132' */
            return """"
        }
        case 133: { /* '133' */
            return """"
        }
        case 134: { /* '134' */
            return """"
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
            return """"
        }
        case 21: { /* '21' */
            return "Description"
        }
        case 22: { /* '22' */
            return """"
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
            return """"
        }
        case 61: { /* '61' */
            return """"
        }
        case 62: { /* '62' */
            return """"
        }
        case 63: { /* '63' */
            return """"
        }
        case 64: { /* '64' */
            return """"
        }
        case 65: { /* '65' */
            return """"
        }
        case 66: { /* '66' */
            return """"
        }
        case 67: { /* '67' */
            return """"
        }
        case 68: { /* '68' */
            return """"
        }
        case 69: { /* '69' */
            return """"
        }
        case 7: { /* '7' */
            return "Table Reference"
        }
        case 70: { /* '70' */
            return "Domain Address"
        }
        case 71: { /* '71' */
            return """"
        }
        case 72: { /* '72' */
            return "Management Descriptor 1"
        }
        case 73: { /* '73' */
            return "PL110 Parameters"
        }
        case 74: { /* '74' */
            return """"
        }
        case 75: { /* '75' */
            return """"
        }
        case 76: { /* '76' */
            return """"
        }
        case 77: { /* '77' */
            return """"
        }
        case 78: { /* '78' */
            return "Hardware Type"
        }
        case 79: { /* '79' */
            return """"
        }
        case 8: { /* '8' */
            return "Service Control"
        }
        case 80: { /* '80' */
            return """"
        }
        case 81: { /* '81' */
            return """"
        }
        case 82: { /* '82' */
            return "RF Domain Address"
        }
        case 83: { /* '83' */
            return """"
        }
        case 84: { /* '84' */
            return """"
        }
        case 85: { /* '85' */
            return """"
        }
        case 86: { /* '86' */
            return """"
        }
        case 9: { /* '9' */
            return "Firmware Revision"
        }
        case 91: { /* '91' */
            return """"
        }
        case 92: { /* '92' */
            return """"
        }
        case 93: { /* '93' */
            return """"
        }
        case 94: { /* '94' */
            return """"
        }
        case 95: { /* '95' */
            return """"
        }
        case 96: { /* '96' */
            return """"
        }
        case 97: { /* '97' */
            return """"
        }
        default: {
            return ""
        }
    }
}
func KnxObjectPropertiesValueOf(value uint8) KnxObjectProperties {
    switch value {
        case 1:
            return KnxObjectProperties_PID_OBJECT_TYPE
        case 10:
            return KnxObjectProperties_PID_SERVICES_SUPPORTED
        case 101:
            return KnxObjectProperties_PID_ON_DELAY
        case 102:
            return KnxObjectProperties_PID_OFF_DELAY
        case 103:
            return KnxObjectProperties_PID_TIMED_ON_DURATION
        case 104:
            return KnxObjectProperties_PID_PREWARNING_DURATION
        case 105:
            return KnxObjectProperties_PID_TRANSMISSION_CYCLE_TIME
        case 106:
            return KnxObjectProperties_PID_BUS_POWER_UP_MESSAGE_DELAY
        case 107:
            return KnxObjectProperties_PID_BEHAVIOUR_AT_LOCKING
        case 108:
            return KnxObjectProperties_PID_BEHAVIOUR_AT_UNLOCKING
        case 109:
            return KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_UP
        case 11:
            return KnxObjectProperties_PID_SERIAL_NUMBER
        case 110:
            return KnxObjectProperties_PID_CHANGE_OF_VALUE
        case 111:
            return KnxObjectProperties_PID_REPETITION_TIME
        case 112:
            return KnxObjectProperties_PID_RF_ENABLE_SBC
        case 113:
            return KnxObjectProperties_PID_MANUAL_OFF_ENABLE
        case 114:
            return KnxObjectProperties_PID_INVERT_LOCK_DEVICE
        case 115:
            return KnxObjectProperties_PID_LOCK_STATE
        case 116:
            return KnxObjectProperties_PID_UNLOCK_STATE
        case 117:
            return KnxObjectProperties_PID_STATE_FOR_SCENE_NUMBER
        case 118:
            return KnxObjectProperties_PID_STORAGE_FUNCTION_FOR_SCENE
        case 119:
            return KnxObjectProperties_PID_BUS_POWER_UP_STATE
        case 12:
            return KnxObjectProperties_PID_MANUFACTURER_ID
        case 120:
            return KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_UP
        case 121:
            return KnxObjectProperties_PID_TIMED_ON_RETRIGGER_FUNCTION
        case 122:
            return KnxObjectProperties_PID_MANUAL_OFF_ENABLE
        case 123:
            return KnxObjectProperties_PID_INVERT_LOCK_DEVICE
        case 124:
            return KnxObjectProperties_PID_BEHAVIOUR_AT_LOCKING
        case 125:
            return KnxObjectProperties_PID_BEHAVIOUR_AT_UNLOCKING
        case 126:
            return KnxObjectProperties_PID_LOCK_SETVALUE
        case 127:
            return KnxObjectProperties_PID_UNLOCK_SETVALUE
        case 128:
            return KnxObjectProperties_PID_BIGHTNESS_FOR_SCENE
        case 129:
            return KnxObjectProperties_PID_STORAGE_FUNCTION_FOR_SCENE
        case 13:
            return KnxObjectProperties_PID_PROGRAM_VERSION
        case 130:
            return KnxObjectProperties_PID_DELTA_DIMMING_VALUE
        case 131:
            return KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_UP
        case 132:
            return KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_UP_SET_VALUE
        case 133:
            return KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_DOWN
        case 134:
            return KnxObjectProperties_PID_BUS_POWER_DOWN_SET_VALUE
        case 14:
            return KnxObjectProperties_PID_DEVICE_CONTROL
        case 15:
            return KnxObjectProperties_PID_ORDER_INFO
        case 16:
            return KnxObjectProperties_PID_PEI_TYPE
        case 17:
            return KnxObjectProperties_PID_PORT_CONFIGURATION
        case 18:
            return KnxObjectProperties_PID_POLL_GROUP_SETTINGS
        case 19:
            return KnxObjectProperties_PID_MANUFACTURER_DATA
        case 2:
            return KnxObjectProperties_PID_OBJECT_NAME
        case 20:
            return KnxObjectProperties_PID_ENABLE
        case 21:
            return KnxObjectProperties_PID_DESCRIPTION
        case 22:
            return KnxObjectProperties_PID_FILE
        case 23:
            return KnxObjectProperties_PID_TABLE
        case 24:
            return KnxObjectProperties_PID_ENROL
        case 25:
            return KnxObjectProperties_PID_VERSION
        case 26:
            return KnxObjectProperties_PID_GROUP_OBJECT_LINK
        case 27:
            return KnxObjectProperties_PID_MCB_TABLE
        case 28:
            return KnxObjectProperties_PID_ERROR_CODE
        case 29:
            return KnxObjectProperties_PID_OBJECT_INDEX
        case 3:
            return KnxObjectProperties_PID_SEMAPHOR
        case 30:
            return KnxObjectProperties_PID_DOWNLOAD_COUNTER
        case 4:
            return KnxObjectProperties_PID_GROUP_OBJECT_REFERENCE
        case 5:
            return KnxObjectProperties_PID_LOAD_STATE_CONTROL
        case 51:
            return KnxObjectProperties_PID_ROUTING_COUNT
        case 52:
            return KnxObjectProperties_PID_MAX_RETRY_COUNT
        case 53:
            return KnxObjectProperties_PID_ERROR_FLAGS
        case 54:
            return KnxObjectProperties_PID_PROGMODE
        case 55:
            return KnxObjectProperties_PID_PRODUCT_ID
        case 56:
            return KnxObjectProperties_PID_MAX_APDULENGTH
        case 57:
            return KnxObjectProperties_PID_SUBNET_ADDR
        case 58:
            return KnxObjectProperties_PID_DEVICE_ADDR
        case 59:
            return KnxObjectProperties_PID_PB_CONFIG
        case 6:
            return KnxObjectProperties_PID_RUN_STATE_CONTROL
        case 60:
            return KnxObjectProperties_PID_ADDR_REPORT
        case 61:
            return KnxObjectProperties_PID_ADDR_CHECK
        case 62:
            return KnxObjectProperties_PID_OBJECT_VALUE
        case 63:
            return KnxObjectProperties_PID_OBJECTLINK
        case 64:
            return KnxObjectProperties_PID_APPLICATION
        case 65:
            return KnxObjectProperties_PID_PARAMETER
        case 66:
            return KnxObjectProperties_PID_OBJECTADDRESS
        case 67:
            return KnxObjectProperties_PID_PSU_TYPE
        case 68:
            return KnxObjectProperties_PID_PSU_STATUS
        case 69:
            return KnxObjectProperties_PID_PSU_ENABLE
        case 7:
            return KnxObjectProperties_PID_TABLE_REFERENCE
        case 70:
            return KnxObjectProperties_PID_DOMAIN_ADDRESS
        case 71:
            return KnxObjectProperties_PID_IO_LIST
        case 72:
            return KnxObjectProperties_PID_MGT_DESCRIPTOR_01
        case 73:
            return KnxObjectProperties_PID_PL110_PARAM
        case 74:
            return KnxObjectProperties_PID_RF_REPEAT_COUNTER
        case 75:
            return KnxObjectProperties_PID_RECEIVE_BLOCK_TABLE
        case 76:
            return KnxObjectProperties_PID_RANDOM_PAUSE_TABLE
        case 77:
            return KnxObjectProperties_PID_RECEIVE_BLOCK_NR
        case 78:
            return KnxObjectProperties_PID_HARDWARE_TYPE
        case 79:
            return KnxObjectProperties_PID_RETRANSMITTER_NUMBER
        case 8:
            return KnxObjectProperties_PID_SERVICE_CONTROL
        case 80:
            return KnxObjectProperties_PID_SERIAL_NR_TABLE
        case 81:
            return KnxObjectProperties_PID_BIBATMASTER_ADDRESS
        case 82:
            return KnxObjectProperties_PID_RF_DOMAIN_ADDRESS
        case 83:
            return KnxObjectProperties_PID_DEVICE_DESCRIPTOR
        case 84:
            return KnxObjectProperties_PID_METERING_FILTER_TABLE
        case 85:
            return KnxObjectProperties_PID_GROUP_TELEGR_RATE_LIMIT_TIME_BASE
        case 86:
            return KnxObjectProperties_PID_GROUP_TELEGR_RATE_LIMIT_NO_OF_TELEGR
        case 9:
            return KnxObjectProperties_PID_FIRMWARE_REVISION
        case 91:
            return KnxObjectProperties_PID_BACKBONE_KEY
        case 92:
            return KnxObjectProperties_PID_DEVICE_AUTHENTICATION_CODE
        case 93:
            return KnxObjectProperties_PID_PASSWORD_HASHES
        case 94:
            return KnxObjectProperties_PID_SECURED_SERVICE_FAMILIES
        case 95:
            return KnxObjectProperties_PID_MULTICAST_LATENCY_TOLERANCE
        case 96:
            return KnxObjectProperties_PID_SYNC_LATENCY_FRACTION
        case 97:
            return KnxObjectProperties_PID_TUNNELLING_USERS
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
    case KnxObjectProperties_PID_OBJECT_TYPE:
        return "PID_OBJECT_TYPE"
    case KnxObjectProperties_PID_SERVICES_SUPPORTED:
        return "PID_SERVICES_SUPPORTED"
    case KnxObjectProperties_PID_ON_DELAY:
        return "PID_ON_DELAY"
    case KnxObjectProperties_PID_OFF_DELAY:
        return "PID_OFF_DELAY"
    case KnxObjectProperties_PID_TIMED_ON_DURATION:
        return "PID_TIMED_ON_DURATION"
    case KnxObjectProperties_PID_PREWARNING_DURATION:
        return "PID_PREWARNING_DURATION"
    case KnxObjectProperties_PID_TRANSMISSION_CYCLE_TIME:
        return "PID_TRANSMISSION_CYCLE_TIME"
    case KnxObjectProperties_PID_BUS_POWER_UP_MESSAGE_DELAY:
        return "PID_BUS_POWER_UP_MESSAGE_DELAY"
    case KnxObjectProperties_PID_BEHAVIOUR_AT_LOCKING:
        return "PID_BEHAVIOUR_AT_LOCKING"
    case KnxObjectProperties_PID_BEHAVIOUR_AT_UNLOCKING:
        return "PID_BEHAVIOUR_AT_UNLOCKING"
    case KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_UP:
        return "PID_BEHAVIOUR_BUS_POWER_UP"
    case KnxObjectProperties_PID_SERIAL_NUMBER:
        return "PID_SERIAL_NUMBER"
    case KnxObjectProperties_PID_CHANGE_OF_VALUE:
        return "PID_CHANGE_OF_VALUE"
    case KnxObjectProperties_PID_REPETITION_TIME:
        return "PID_REPETITION_TIME"
    case KnxObjectProperties_PID_RF_ENABLE_SBC:
        return "PID_RF_ENABLE_SBC"
    case KnxObjectProperties_PID_MANUAL_OFF_ENABLE:
        return "PID_MANUAL_OFF_ENABLE"
    case KnxObjectProperties_PID_INVERT_LOCK_DEVICE:
        return "PID_INVERT_LOCK_DEVICE"
    case KnxObjectProperties_PID_LOCK_STATE:
        return "PID_LOCK_STATE"
    case KnxObjectProperties_PID_UNLOCK_STATE:
        return "PID_UNLOCK_STATE"
    case KnxObjectProperties_PID_STATE_FOR_SCENE_NUMBER:
        return "PID_STATE_FOR_SCENE_NUMBER"
    case KnxObjectProperties_PID_STORAGE_FUNCTION_FOR_SCENE:
        return "PID_STORAGE_FUNCTION_FOR_SCENE"
    case KnxObjectProperties_PID_BUS_POWER_UP_STATE:
        return "PID_BUS_POWER_UP_STATE"
    case KnxObjectProperties_PID_MANUFACTURER_ID:
        return "PID_MANUFACTURER_ID"
    case KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_UP:
        return "PID_BEHAVIOUR_BUS_POWER_UP"
    case KnxObjectProperties_PID_TIMED_ON_RETRIGGER_FUNCTION:
        return "PID_TIMED_ON_RETRIGGER_FUNCTION"
    case KnxObjectProperties_PID_MANUAL_OFF_ENABLE:
        return "PID_MANUAL_OFF_ENABLE"
    case KnxObjectProperties_PID_INVERT_LOCK_DEVICE:
        return "PID_INVERT_LOCK_DEVICE"
    case KnxObjectProperties_PID_BEHAVIOUR_AT_LOCKING:
        return "PID_BEHAVIOUR_AT_LOCKING"
    case KnxObjectProperties_PID_BEHAVIOUR_AT_UNLOCKING:
        return "PID_BEHAVIOUR_AT_UNLOCKING"
    case KnxObjectProperties_PID_LOCK_SETVALUE:
        return "PID_LOCK_SETVALUE"
    case KnxObjectProperties_PID_UNLOCK_SETVALUE:
        return "PID_UNLOCK_SETVALUE"
    case KnxObjectProperties_PID_BIGHTNESS_FOR_SCENE:
        return "PID_BIGHTNESS_FOR_SCENE"
    case KnxObjectProperties_PID_STORAGE_FUNCTION_FOR_SCENE:
        return "PID_STORAGE_FUNCTION_FOR_SCENE"
    case KnxObjectProperties_PID_PROGRAM_VERSION:
        return "PID_PROGRAM_VERSION"
    case KnxObjectProperties_PID_DELTA_DIMMING_VALUE:
        return "PID_DELTA_DIMMING_VALUE"
    case KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_UP:
        return "PID_BEHAVIOUR_BUS_POWER_UP"
    case KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_UP_SET_VALUE:
        return "PID_BEHAVIOUR_BUS_POWER_UP_SET_VALUE"
    case KnxObjectProperties_PID_BEHAVIOUR_BUS_POWER_DOWN:
        return "PID_BEHAVIOUR_BUS_POWER_DOWN"
    case KnxObjectProperties_PID_BUS_POWER_DOWN_SET_VALUE:
        return "PID_BUS_POWER_DOWN_SET_VALUE"
    case KnxObjectProperties_PID_DEVICE_CONTROL:
        return "PID_DEVICE_CONTROL"
    case KnxObjectProperties_PID_ORDER_INFO:
        return "PID_ORDER_INFO"
    case KnxObjectProperties_PID_PEI_TYPE:
        return "PID_PEI_TYPE"
    case KnxObjectProperties_PID_PORT_CONFIGURATION:
        return "PID_PORT_CONFIGURATION"
    case KnxObjectProperties_PID_POLL_GROUP_SETTINGS:
        return "PID_POLL_GROUP_SETTINGS"
    case KnxObjectProperties_PID_MANUFACTURER_DATA:
        return "PID_MANUFACTURER_DATA"
    case KnxObjectProperties_PID_OBJECT_NAME:
        return "PID_OBJECT_NAME"
    case KnxObjectProperties_PID_ENABLE:
        return "PID_ENABLE"
    case KnxObjectProperties_PID_DESCRIPTION:
        return "PID_DESCRIPTION"
    case KnxObjectProperties_PID_FILE:
        return "PID_FILE"
    case KnxObjectProperties_PID_TABLE:
        return "PID_TABLE"
    case KnxObjectProperties_PID_ENROL:
        return "PID_ENROL"
    case KnxObjectProperties_PID_VERSION:
        return "PID_VERSION"
    case KnxObjectProperties_PID_GROUP_OBJECT_LINK:
        return "PID_GROUP_OBJECT_LINK"
    case KnxObjectProperties_PID_MCB_TABLE:
        return "PID_MCB_TABLE"
    case KnxObjectProperties_PID_ERROR_CODE:
        return "PID_ERROR_CODE"
    case KnxObjectProperties_PID_OBJECT_INDEX:
        return "PID_OBJECT_INDEX"
    case KnxObjectProperties_PID_SEMAPHOR:
        return "PID_SEMAPHOR"
    case KnxObjectProperties_PID_DOWNLOAD_COUNTER:
        return "PID_DOWNLOAD_COUNTER"
    case KnxObjectProperties_PID_GROUP_OBJECT_REFERENCE:
        return "PID_GROUP_OBJECT_REFERENCE"
    case KnxObjectProperties_PID_LOAD_STATE_CONTROL:
        return "PID_LOAD_STATE_CONTROL"
    case KnxObjectProperties_PID_ROUTING_COUNT:
        return "PID_ROUTING_COUNT"
    case KnxObjectProperties_PID_MAX_RETRY_COUNT:
        return "PID_MAX_RETRY_COUNT"
    case KnxObjectProperties_PID_ERROR_FLAGS:
        return "PID_ERROR_FLAGS"
    case KnxObjectProperties_PID_PROGMODE:
        return "PID_PROGMODE"
    case KnxObjectProperties_PID_PRODUCT_ID:
        return "PID_PRODUCT_ID"
    case KnxObjectProperties_PID_MAX_APDULENGTH:
        return "PID_MAX_APDULENGTH"
    case KnxObjectProperties_PID_SUBNET_ADDR:
        return "PID_SUBNET_ADDR"
    case KnxObjectProperties_PID_DEVICE_ADDR:
        return "PID_DEVICE_ADDR"
    case KnxObjectProperties_PID_PB_CONFIG:
        return "PID_PB_CONFIG"
    case KnxObjectProperties_PID_RUN_STATE_CONTROL:
        return "PID_RUN_STATE_CONTROL"
    case KnxObjectProperties_PID_ADDR_REPORT:
        return "PID_ADDR_REPORT"
    case KnxObjectProperties_PID_ADDR_CHECK:
        return "PID_ADDR_CHECK"
    case KnxObjectProperties_PID_OBJECT_VALUE:
        return "PID_OBJECT_VALUE"
    case KnxObjectProperties_PID_OBJECTLINK:
        return "PID_OBJECTLINK"
    case KnxObjectProperties_PID_APPLICATION:
        return "PID_APPLICATION"
    case KnxObjectProperties_PID_PARAMETER:
        return "PID_PARAMETER"
    case KnxObjectProperties_PID_OBJECTADDRESS:
        return "PID_OBJECTADDRESS"
    case KnxObjectProperties_PID_PSU_TYPE:
        return "PID_PSU_TYPE"
    case KnxObjectProperties_PID_PSU_STATUS:
        return "PID_PSU_STATUS"
    case KnxObjectProperties_PID_PSU_ENABLE:
        return "PID_PSU_ENABLE"
    case KnxObjectProperties_PID_TABLE_REFERENCE:
        return "PID_TABLE_REFERENCE"
    case KnxObjectProperties_PID_DOMAIN_ADDRESS:
        return "PID_DOMAIN_ADDRESS"
    case KnxObjectProperties_PID_IO_LIST:
        return "PID_IO_LIST"
    case KnxObjectProperties_PID_MGT_DESCRIPTOR_01:
        return "PID_MGT_DESCRIPTOR_01"
    case KnxObjectProperties_PID_PL110_PARAM:
        return "PID_PL110_PARAM"
    case KnxObjectProperties_PID_RF_REPEAT_COUNTER:
        return "PID_RF_REPEAT_COUNTER"
    case KnxObjectProperties_PID_RECEIVE_BLOCK_TABLE:
        return "PID_RECEIVE_BLOCK_TABLE"
    case KnxObjectProperties_PID_RANDOM_PAUSE_TABLE:
        return "PID_RANDOM_PAUSE_TABLE"
    case KnxObjectProperties_PID_RECEIVE_BLOCK_NR:
        return "PID_RECEIVE_BLOCK_NR"
    case KnxObjectProperties_PID_HARDWARE_TYPE:
        return "PID_HARDWARE_TYPE"
    case KnxObjectProperties_PID_RETRANSMITTER_NUMBER:
        return "PID_RETRANSMITTER_NUMBER"
    case KnxObjectProperties_PID_SERVICE_CONTROL:
        return "PID_SERVICE_CONTROL"
    case KnxObjectProperties_PID_SERIAL_NR_TABLE:
        return "PID_SERIAL_NR_TABLE"
    case KnxObjectProperties_PID_BIBATMASTER_ADDRESS:
        return "PID_BIBATMASTER_ADDRESS"
    case KnxObjectProperties_PID_RF_DOMAIN_ADDRESS:
        return "PID_RF_DOMAIN_ADDRESS"
    case KnxObjectProperties_PID_DEVICE_DESCRIPTOR:
        return "PID_DEVICE_DESCRIPTOR"
    case KnxObjectProperties_PID_METERING_FILTER_TABLE:
        return "PID_METERING_FILTER_TABLE"
    case KnxObjectProperties_PID_GROUP_TELEGR_RATE_LIMIT_TIME_BASE:
        return "PID_GROUP_TELEGR_RATE_LIMIT_TIME_BASE"
    case KnxObjectProperties_PID_GROUP_TELEGR_RATE_LIMIT_NO_OF_TELEGR:
        return "PID_GROUP_TELEGR_RATE_LIMIT_NO_OF_TELEGR"
    case KnxObjectProperties_PID_FIRMWARE_REVISION:
        return "PID_FIRMWARE_REVISION"
    case KnxObjectProperties_PID_BACKBONE_KEY:
        return "PID_BACKBONE_KEY"
    case KnxObjectProperties_PID_DEVICE_AUTHENTICATION_CODE:
        return "PID_DEVICE_AUTHENTICATION_CODE"
    case KnxObjectProperties_PID_PASSWORD_HASHES:
        return "PID_PASSWORD_HASHES"
    case KnxObjectProperties_PID_SECURED_SERVICE_FAMILIES:
        return "PID_SECURED_SERVICE_FAMILIES"
    case KnxObjectProperties_PID_MULTICAST_LATENCY_TOLERANCE:
        return "PID_MULTICAST_LATENCY_TOLERANCE"
    case KnxObjectProperties_PID_SYNC_LATENCY_FRACTION:
        return "PID_SYNC_LATENCY_FRACTION"
    case KnxObjectProperties_PID_TUNNELLING_USERS:
        return "PID_TUNNELLING_USERS"
    }
    return ""
}
