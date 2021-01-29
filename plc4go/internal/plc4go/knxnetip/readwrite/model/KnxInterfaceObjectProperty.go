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
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
)

type KnxInterfaceObjectProperty uint32

type IKnxInterfaceObjectProperty interface {
	PropertyDataType() KnxPropertyDataType
	Text() string
	PropertyId() uint16
	ObjectType() KnxInterfaceObjectType
	Serialize(io utils.WriteBuffer) error
}

const (
	KnxInterfaceObjectProperty_PID_GENERAL_OBJECT_TYPE                                              KnxInterfaceObjectProperty = 65535001
	KnxInterfaceObjectProperty_PID_GENERAL_OBJECT_NAME                                              KnxInterfaceObjectProperty = 65535002
	KnxInterfaceObjectProperty_PID_GENERAL_SEMAPHOR                                                 KnxInterfaceObjectProperty = 65535003
	KnxInterfaceObjectProperty_PID_GENERAL_GROUP_OBJECT_REFERENCE                                   KnxInterfaceObjectProperty = 65535004
	KnxInterfaceObjectProperty_PID_GENERAL_LOAD_STATE_CONTROL                                       KnxInterfaceObjectProperty = 65535005
	KnxInterfaceObjectProperty_PID_GENERAL_RUN_STATE_CONTROL                                        KnxInterfaceObjectProperty = 65535006
	KnxInterfaceObjectProperty_PID_GENERAL_TABLE_REFERENCE                                          KnxInterfaceObjectProperty = 65535007
	KnxInterfaceObjectProperty_PID_GENERAL_SERVICE_CONTROL                                          KnxInterfaceObjectProperty = 65535008
	KnxInterfaceObjectProperty_PID_GENERAL_FIRMWARE_REVISION                                        KnxInterfaceObjectProperty = 65535009
	KnxInterfaceObjectProperty_PID_GENERAL_SERVICES_SUPPORTED                                       KnxInterfaceObjectProperty = 65535010
	KnxInterfaceObjectProperty_PID_GENERAL_SERIAL_NUMBER                                            KnxInterfaceObjectProperty = 65535011
	KnxInterfaceObjectProperty_PID_GENERAL_MANUFACTURER_ID                                          KnxInterfaceObjectProperty = 65535012
	KnxInterfaceObjectProperty_PID_GENERAL_PROGRAM_VERSION                                          KnxInterfaceObjectProperty = 65535013
	KnxInterfaceObjectProperty_PID_GENERAL_DEVICE_CONTROL                                           KnxInterfaceObjectProperty = 65535014
	KnxInterfaceObjectProperty_PID_GENERAL_ORDER_INFO                                               KnxInterfaceObjectProperty = 65535015
	KnxInterfaceObjectProperty_PID_GENERAL_PEI_TYPE                                                 KnxInterfaceObjectProperty = 65535016
	KnxInterfaceObjectProperty_PID_GENERAL_PORT_CONFIGURATION                                       KnxInterfaceObjectProperty = 65535017
	KnxInterfaceObjectProperty_PID_GENERAL_POLL_GROUP_SETTINGS                                      KnxInterfaceObjectProperty = 65535018
	KnxInterfaceObjectProperty_PID_GENERAL_MANUFACTURER_DATA                                        KnxInterfaceObjectProperty = 65535019
	KnxInterfaceObjectProperty_PID_GENERAL_ENABLE                                                   KnxInterfaceObjectProperty = 65535020
	KnxInterfaceObjectProperty_PID_GENERAL_DESCRIPTION                                              KnxInterfaceObjectProperty = 65535021
	KnxInterfaceObjectProperty_PID_GENERAL_FILE                                                     KnxInterfaceObjectProperty = 65535022
	KnxInterfaceObjectProperty_PID_GENERAL_TABLE                                                    KnxInterfaceObjectProperty = 65535023
	KnxInterfaceObjectProperty_PID_GENERAL_ENROL                                                    KnxInterfaceObjectProperty = 65535024
	KnxInterfaceObjectProperty_PID_GENERAL_VERSION                                                  KnxInterfaceObjectProperty = 65535025
	KnxInterfaceObjectProperty_PID_GENERAL_GROUP_OBJECT_LINK                                        KnxInterfaceObjectProperty = 65535026
	KnxInterfaceObjectProperty_PID_GENERAL_MCB_TABLE                                                KnxInterfaceObjectProperty = 65535027
	KnxInterfaceObjectProperty_PID_GENERAL_ERROR_CODE                                               KnxInterfaceObjectProperty = 65535028
	KnxInterfaceObjectProperty_PID_GENERAL_OBJECT_INDEX                                             KnxInterfaceObjectProperty = 65535029
	KnxInterfaceObjectProperty_PID_GENERAL_DOWNLOAD_COUNTER                                         KnxInterfaceObjectProperty = 65535030
	KnxInterfaceObjectProperty_PID_DEVICE_ROUTING_COUNT                                             KnxInterfaceObjectProperty = 51
	KnxInterfaceObjectProperty_PID_DEVICE_MAX_RETRY_COUNT                                           KnxInterfaceObjectProperty = 52
	KnxInterfaceObjectProperty_PID_DEVICE_ERROR_FLAGS                                               KnxInterfaceObjectProperty = 53
	KnxInterfaceObjectProperty_PID_DEVICE_PROGMODE                                                  KnxInterfaceObjectProperty = 54
	KnxInterfaceObjectProperty_PID_DEVICE_PRODUCT_ID                                                KnxInterfaceObjectProperty = 55
	KnxInterfaceObjectProperty_PID_DEVICE_MAX_APDULENGTH                                            KnxInterfaceObjectProperty = 56
	KnxInterfaceObjectProperty_PID_DEVICE_SUBNET_ADDR                                               KnxInterfaceObjectProperty = 57
	KnxInterfaceObjectProperty_PID_DEVICE_DEVICE_ADDR                                               KnxInterfaceObjectProperty = 58
	KnxInterfaceObjectProperty_PID_DEVICE_PB_CONFIG                                                 KnxInterfaceObjectProperty = 59
	KnxInterfaceObjectProperty_PID_DEVICE_ADDR_REPORT                                               KnxInterfaceObjectProperty = 60
	KnxInterfaceObjectProperty_PID_DEVICE_ADDR_CHECK                                                KnxInterfaceObjectProperty = 61
	KnxInterfaceObjectProperty_PID_DEVICE_OBJECT_VALUE                                              KnxInterfaceObjectProperty = 62
	KnxInterfaceObjectProperty_PID_DEVICE_OBJECTLINK                                                KnxInterfaceObjectProperty = 63
	KnxInterfaceObjectProperty_PID_DEVICE_APPLICATION                                               KnxInterfaceObjectProperty = 64
	KnxInterfaceObjectProperty_PID_DEVICE_PARAMETER                                                 KnxInterfaceObjectProperty = 65
	KnxInterfaceObjectProperty_PID_DEVICE_OBJECTADDRESS                                             KnxInterfaceObjectProperty = 66
	KnxInterfaceObjectProperty_PID_DEVICE_PSU_TYPE                                                  KnxInterfaceObjectProperty = 67
	KnxInterfaceObjectProperty_PID_DEVICE_PSU_STATUS                                                KnxInterfaceObjectProperty = 68
	KnxInterfaceObjectProperty_PID_DEVICE_PSU_ENABLE                                                KnxInterfaceObjectProperty = 69
	KnxInterfaceObjectProperty_PID_DEVICE_DOMAIN_ADDRESS                                            KnxInterfaceObjectProperty = 70
	KnxInterfaceObjectProperty_PID_DEVICE_IO_LIST                                                   KnxInterfaceObjectProperty = 71
	KnxInterfaceObjectProperty_PID_DEVICE_MGT_DESCRIPTOR_01                                         KnxInterfaceObjectProperty = 72
	KnxInterfaceObjectProperty_PID_DEVICE_PL110_PARAM                                               KnxInterfaceObjectProperty = 73
	KnxInterfaceObjectProperty_PID_DEVICE_RF_REPEAT_COUNTER                                         KnxInterfaceObjectProperty = 74
	KnxInterfaceObjectProperty_PID_DEVICE_RECEIVE_BLOCK_TABLE                                       KnxInterfaceObjectProperty = 75
	KnxInterfaceObjectProperty_PID_DEVICE_RANDOM_PAUSE_TABLE                                        KnxInterfaceObjectProperty = 76
	KnxInterfaceObjectProperty_PID_DEVICE_RECEIVE_BLOCK_NR                                          KnxInterfaceObjectProperty = 77
	KnxInterfaceObjectProperty_PID_DEVICE_HARDWARE_TYPE                                             KnxInterfaceObjectProperty = 78
	KnxInterfaceObjectProperty_PID_DEVICE_RETRANSMITTER_NUMBER                                      KnxInterfaceObjectProperty = 79
	KnxInterfaceObjectProperty_PID_DEVICE_SERIAL_NR_TABLE                                           KnxInterfaceObjectProperty = 80
	KnxInterfaceObjectProperty_PID_DEVICE_BIBATMASTER_ADDRESS                                       KnxInterfaceObjectProperty = 81
	KnxInterfaceObjectProperty_PID_DEVICE_RF_DOMAIN_ADDRESS                                         KnxInterfaceObjectProperty = 82
	KnxInterfaceObjectProperty_PID_DEVICE_DEVICE_DESCRIPTOR                                         KnxInterfaceObjectProperty = 83
	KnxInterfaceObjectProperty_PID_DEVICE_METERING_FILTER_TABLE                                     KnxInterfaceObjectProperty = 84
	KnxInterfaceObjectProperty_PID_DEVICE_GROUP_TELEGR_RATE_LIMIT_TIME_BASE                         KnxInterfaceObjectProperty = 85
	KnxInterfaceObjectProperty_PID_DEVICE_GROUP_TELEGR_RATE_LIMIT_NO_OF_TELEGR                      KnxInterfaceObjectProperty = 86
	KnxInterfaceObjectProperty_PID_GROUP_OBJECT_TABLE_GRPOBJTABLE                                   KnxInterfaceObjectProperty = 9051
	KnxInterfaceObjectProperty_PID_GROUP_OBJECT_TABLE_EXT_GRPOBJREFERENCE                           KnxInterfaceObjectProperty = 9052
	KnxInterfaceObjectProperty_PID_ROUTER_LINE_STATUS                                               KnxInterfaceObjectProperty = 6051
	KnxInterfaceObjectProperty_PID_ROUTER_MAIN_LCCONFIG                                             KnxInterfaceObjectProperty = 6052
	KnxInterfaceObjectProperty_PID_ROUTER_SUB_LCCONFIG                                              KnxInterfaceObjectProperty = 6053
	KnxInterfaceObjectProperty_PID_ROUTER_MAIN_LCGRPCONFIG                                          KnxInterfaceObjectProperty = 6054
	KnxInterfaceObjectProperty_PID_ROUTER_SUB_LCGRPCONFIG                                           KnxInterfaceObjectProperty = 6055
	KnxInterfaceObjectProperty_PID_ROUTER_ROUTETABLE_CONTROL                                        KnxInterfaceObjectProperty = 6056
	KnxInterfaceObjectProperty_PID_ROUTER_COUPL_SERV_CONTROL                                        KnxInterfaceObjectProperty = 6057
	KnxInterfaceObjectProperty_PID_ROUTER_MAX_ROUTER_APDU_LENGTH                                    KnxInterfaceObjectProperty = 6058
	KnxInterfaceObjectProperty_PID_ROUTER_MEDIUM                                                    KnxInterfaceObjectProperty = 6063
	KnxInterfaceObjectProperty_PID_ROUTER_FILTER_TABLE_USE                                          KnxInterfaceObjectProperty = 6067
	KnxInterfaceObjectProperty_PID_ROUTER_RF_ENABLE_SBC                                             KnxInterfaceObjectProperty = 6112
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_PROJECT_INSTALLATION_ID                          KnxInterfaceObjectProperty = 11051
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNX_INDIVIDUAL_ADDRESS                           KnxInterfaceObjectProperty = 11052
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_ADDITIONAL_INDIVIDUAL_ADDRESSES                  KnxInterfaceObjectProperty = 11053
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_IP_ASSIGNMENT_METHOD                     KnxInterfaceObjectProperty = 11054
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_IP_ASSIGNMENT_METHOD                             KnxInterfaceObjectProperty = 11055
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_IP_CAPABILITIES                                  KnxInterfaceObjectProperty = 11056
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_IP_ADDRESS                               KnxInterfaceObjectProperty = 11057
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_SUBNET_MASK                              KnxInterfaceObjectProperty = 11058
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_DEFAULT_GATEWAY                          KnxInterfaceObjectProperty = 11059
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_IP_ADDRESS                                       KnxInterfaceObjectProperty = 11060
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SUBNET_MASK                                      KnxInterfaceObjectProperty = 11061
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_DEFAULT_GATEWAY                                  KnxInterfaceObjectProperty = 11062
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_DHCP_BOOTP_SERVER                                KnxInterfaceObjectProperty = 11063
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MAC_ADDRESS                                      KnxInterfaceObjectProperty = 11064
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SYSTEM_SETUP_MULTICAST_ADDRESS                   KnxInterfaceObjectProperty = 11065
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_ROUTING_MULTICAST_ADDRESS                        KnxInterfaceObjectProperty = 11066
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_TTL                                              KnxInterfaceObjectProperty = 11067
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNXNETIP_DEVICE_CAPABILITIES                     KnxInterfaceObjectProperty = 11068
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNXNETIP_DEVICE_STATE                            KnxInterfaceObjectProperty = 11069
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNXNETIP_ROUTING_CAPABILITIES                    KnxInterfaceObjectProperty = 11070
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_PRIORITY_FIFO_ENABLED                            KnxInterfaceObjectProperty = 11071
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_QUEUE_OVERFLOW_TO_IP                             KnxInterfaceObjectProperty = 11072
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_QUEUE_OVERFLOW_TO_KNX                            KnxInterfaceObjectProperty = 11073
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MSG_TRANSMIT_TO_IP                               KnxInterfaceObjectProperty = 11074
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MSG_TRANSMIT_TO_KNX                              KnxInterfaceObjectProperty = 11075
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_FRIENDLY_NAME                                    KnxInterfaceObjectProperty = 11076
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_BACKBONE_KEY                                     KnxInterfaceObjectProperty = 11091
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_DEVICE_AUTHENTICATION_CODE                       KnxInterfaceObjectProperty = 11092
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_PASSWORD_HASHES                                  KnxInterfaceObjectProperty = 11093
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SECURED_SERVICE_FAMILIES                         KnxInterfaceObjectProperty = 11094
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MULTICAST_LATENCY_TOLERANCE                      KnxInterfaceObjectProperty = 11095
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SYNC_LATENCY_FRACTION                            KnxInterfaceObjectProperty = 11096
	KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_TUNNELLING_USERS                                 KnxInterfaceObjectProperty = 11097
	KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_MODE                                           KnxInterfaceObjectProperty = 17051
	KnxInterfaceObjectProperty_PID_SECURITY_P2P_KEY_TABLE                                           KnxInterfaceObjectProperty = 17052
	KnxInterfaceObjectProperty_PID_SECURITY_GRP_KEY_TABLE                                           KnxInterfaceObjectProperty = 17053
	KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_INDIVIDUAL_ADDRESS_TABLE                       KnxInterfaceObjectProperty = 17054
	KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_FAILURES_LOG                                   KnxInterfaceObjectProperty = 17055
	KnxInterfaceObjectProperty_PID_SECURITY_SKI_TOOL                                                KnxInterfaceObjectProperty = 17056
	KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_REPORT                                         KnxInterfaceObjectProperty = 17057
	KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_REPORT_CONTROL                                 KnxInterfaceObjectProperty = 17058
	KnxInterfaceObjectProperty_PID_SECURITY_SEQUENCE_NUMBER_SENDING                                 KnxInterfaceObjectProperty = 17059
	KnxInterfaceObjectProperty_PID_SECURITY_ZONE_KEYS_TABLE                                         KnxInterfaceObjectProperty = 17060
	KnxInterfaceObjectProperty_PID_SECURITY_GO_SECURITY_FLAGS                                       KnxInterfaceObjectProperty = 17061
	KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_MULTI_TYPE                                          KnxInterfaceObjectProperty = 19051
	KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DOMAIN_ADDRESS                                      KnxInterfaceObjectProperty = 19056
	KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_RETRANSMITTER                                       KnxInterfaceObjectProperty = 19057
	KnxInterfaceObjectProperty_PID_RF_MEDIUM_SECURITY_REPORT_CONTROL                                KnxInterfaceObjectProperty = 19058
	KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_FILTERING_MODE_SELECT                               KnxInterfaceObjectProperty = 19059
	KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_BIDIR_TIMEOUT                                       KnxInterfaceObjectProperty = 19060
	KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DIAG_SA_FILTER_TABLE                                KnxInterfaceObjectProperty = 19061
	KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DIAG_QUALITY_TABLE                                  KnxInterfaceObjectProperty = 19062
	KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DIAG_PROBE                                          KnxInterfaceObjectProperty = 19063
	KnxInterfaceObjectProperty_PID_INDOOR_BRIGHTNESS_SENSOR_CHANGE_OF_VALUE                         KnxInterfaceObjectProperty = 409110
	KnxInterfaceObjectProperty_PID_INDOOR_BRIGHTNESS_SENSOR_REPETITION_TIME                         KnxInterfaceObjectProperty = 409111
	KnxInterfaceObjectProperty_PID_INDOOR_LUMINANCE_SENSOR_CHANGE_OF_VALUE                          KnxInterfaceObjectProperty = 410110
	KnxInterfaceObjectProperty_PID_INDOOR_LUMINANCE_SENSOR_REPETITION_TIME                          KnxInterfaceObjectProperty = 410111
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_ON_DELAY                          KnxInterfaceObjectProperty = 417101
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_OFF_DELAY                         KnxInterfaceObjectProperty = 417102
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TIMED_ON_DURATION                 KnxInterfaceObjectProperty = 417103
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_PREWARNING_DURATION               KnxInterfaceObjectProperty = 417104
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TRANSMISSION_CYCLE_TIME           KnxInterfaceObjectProperty = 417105
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BUS_POWER_UP_MESSAGE_DELAY        KnxInterfaceObjectProperty = 417106
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_AT_LOCKING              KnxInterfaceObjectProperty = 417107
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_AT_UNLOCKING            KnxInterfaceObjectProperty = 417108
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP            KnxInterfaceObjectProperty = 417109
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_DOWN          KnxInterfaceObjectProperty = 417110
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_INVERT_OUTPUT_STATE               KnxInterfaceObjectProperty = 417111
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TIMED_ON_RETRIGGER_FUNCTION       KnxInterfaceObjectProperty = 417112
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_MANUAL_OFF_ENABLE                 KnxInterfaceObjectProperty = 417113
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_INVERT_LOCK_DEVICE                KnxInterfaceObjectProperty = 417114
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_LOCK_STATE                        KnxInterfaceObjectProperty = 417115
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_UNLOCK_STATE                      KnxInterfaceObjectProperty = 417116
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_STATE_FOR_SCENE_NUMBER            KnxInterfaceObjectProperty = 417117
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_STORAGE_FUNCTION_FOR_SCENE        KnxInterfaceObjectProperty = 417118
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BUS_POWER_UP_STATE                KnxInterfaceObjectProperty = 417119
	KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP_2          KnxInterfaceObjectProperty = 417120
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_ON_DELAY                                  KnxInterfaceObjectProperty = 418101
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_OFF_DELAY                                 KnxInterfaceObjectProperty = 418102
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_SWITCH_OFF_BRIGHTNESS_DELAY_TIME          KnxInterfaceObjectProperty = 418103
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_TIMED_ON_DURATION                         KnxInterfaceObjectProperty = 418104
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_PREWARNING_DURATION                       KnxInterfaceObjectProperty = 418105
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_TRANSMISSION_CYCLE_TIME                   KnxInterfaceObjectProperty = 418106
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BUS_POWER_UP_MESSAGE_DELAY                KnxInterfaceObjectProperty = 418107
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED                             KnxInterfaceObjectProperty = 418108
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME                         KnxInterfaceObjectProperty = 418109
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED_FOR_SWITCH_ON_SET_VALUE     KnxInterfaceObjectProperty = 418110
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED_FOR_SWITCH_OFF              KnxInterfaceObjectProperty = 418111
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME_FOR_SWITCH_ON_SET_VALUE KnxInterfaceObjectProperty = 418112
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME_FOR_SWITCH_OFF          KnxInterfaceObjectProperty = 418113
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_SWITCFH_OFF_BRIGHTNESS                    KnxInterfaceObjectProperty = 418114
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MINIMUM_SET_VALUE                         KnxInterfaceObjectProperty = 418115
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MAXIMUM_SET_VALUE                         KnxInterfaceObjectProperty = 418116
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_SWITCH_ON_SET_VALUE                       KnxInterfaceObjectProperty = 418117
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMM_MODE_SELECTION                       KnxInterfaceObjectProperty = 418118
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_RELATIV_OFF_ENABLE                        KnxInterfaceObjectProperty = 418119
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MEMORY_FUNCTION                           KnxInterfaceObjectProperty = 418120
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_TIMED_ON_RETRIGGER_FUNCTION               KnxInterfaceObjectProperty = 418121
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MANUAL_OFF_ENABLE                         KnxInterfaceObjectProperty = 418122
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_INVERT_LOCK_DEVICE                        KnxInterfaceObjectProperty = 418123
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_AT_LOCKING                      KnxInterfaceObjectProperty = 418124
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_AT_UNLOCKING                    KnxInterfaceObjectProperty = 418125
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_LOCK_SETVALUE                             KnxInterfaceObjectProperty = 418126
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_UNLOCK_SETVALUE                           KnxInterfaceObjectProperty = 418127
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BIGHTNESS_FOR_SCENE                       KnxInterfaceObjectProperty = 418128
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_STORAGE_FUNCTION_FOR_SCENE                KnxInterfaceObjectProperty = 418129
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DELTA_DIMMING_VALUE                       KnxInterfaceObjectProperty = 418130
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP                    KnxInterfaceObjectProperty = 418131
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP_SET_VALUE          KnxInterfaceObjectProperty = 418132
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_DOWN                  KnxInterfaceObjectProperty = 418133
	KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BUS_POWER_DOWN_SET_VALUE                  KnxInterfaceObjectProperty = 418134
	KnxInterfaceObjectProperty_PID_DIMMING_SENSOR_BASIC_ON_OFF_ACTION                               KnxInterfaceObjectProperty = 420051
	KnxInterfaceObjectProperty_PID_DIMMING_SENSOR_BASIC_ENABLE_TOGGLE_MODE                          KnxInterfaceObjectProperty = 420052
	KnxInterfaceObjectProperty_PID_DIMMING_SENSOR_BASIC_ABSOLUTE_SETVALUE                           KnxInterfaceObjectProperty = 420053
	KnxInterfaceObjectProperty_PID_SWITCHING_SENSOR_BASIC_ON_OFF_ACTION                             KnxInterfaceObjectProperty = 421051
	KnxInterfaceObjectProperty_PID_SWITCHING_SENSOR_BASIC_ENABLE_TOGGLE_MODE                        KnxInterfaceObjectProperty = 421052
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REVERSION_PAUSE_TIME                     KnxInterfaceObjectProperty = 800051
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_UP_DOWN_TIME                        KnxInterfaceObjectProperty = 800052
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_SLAT_STEP_TIME                           KnxInterfaceObjectProperty = 800053
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_PRESET_POSITION_TIME                KnxInterfaceObjectProperty = 800054
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_TO_PRESET_POSITION_IN_PERCENT       KnxInterfaceObjectProperty = 800055
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_TO_PRESET_POSITION_LENGTH           KnxInterfaceObjectProperty = 800057
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_PRESET_SLAT_POSITION_PERCENT             KnxInterfaceObjectProperty = 800058
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_PRESET_SLAT_POSITION_ANGLE               KnxInterfaceObjectProperty = 800060
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REACTION_WIND_ALARM                      KnxInterfaceObjectProperty = 800061
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_WIND_ALARM                     KnxInterfaceObjectProperty = 800062
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REACTION_ON_RAIN_ALARM                   KnxInterfaceObjectProperty = 800063
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_RAIN_ALARM                     KnxInterfaceObjectProperty = 800064
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REACTION_FROST_ALARM                     KnxInterfaceObjectProperty = 800065
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_FROST_ALARM                    KnxInterfaceObjectProperty = 800066
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MAX_SLAT_MOVE_TIME                       KnxInterfaceObjectProperty = 800067
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_ENABLE_BLINDS_MODE                       KnxInterfaceObjectProperty = 800068
	KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_STORAGE_FUNCTIONS_FOR_SCENE              KnxInterfaceObjectProperty = 800069
	KnxInterfaceObjectProperty_PID_SUNBLIND_SENSOR_BASIC_ENABLE_BLINDS_MODE                         KnxInterfaceObjectProperty = 801051
	KnxInterfaceObjectProperty_PID_SUNBLIND_SENSOR_BASIC_UP_DOWN_ACTION                             KnxInterfaceObjectProperty = 801052
	KnxInterfaceObjectProperty_PID_SUNBLIND_SENSOR_BASIC_ENABLE_TOGGLE_MODE                         KnxInterfaceObjectProperty = 801053
)

func (e KnxInterfaceObjectProperty) PropertyDataType() KnxPropertyDataType {
	switch e {
	case 11051:
		{ /* '11051' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11052:
		{ /* '11052' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11053:
		{ /* '11053' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11054:
		{ /* '11054' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11055:
		{ /* '11055' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11056:
		{ /* '11056' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11057:
		{ /* '11057' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11058:
		{ /* '11058' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11059:
		{ /* '11059' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11060:
		{ /* '11060' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11061:
		{ /* '11061' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11062:
		{ /* '11062' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11063:
		{ /* '11063' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11064:
		{ /* '11064' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11065:
		{ /* '11065' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11066:
		{ /* '11066' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11067:
		{ /* '11067' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11068:
		{ /* '11068' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11069:
		{ /* '11069' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11070:
		{ /* '11070' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11071:
		{ /* '11071' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11072:
		{ /* '11072' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11073:
		{ /* '11073' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11074:
		{ /* '11074' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11075:
		{ /* '11075' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11076:
		{ /* '11076' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11091:
		{ /* '11091' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11092:
		{ /* '11092' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11093:
		{ /* '11093' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11094:
		{ /* '11094' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11095:
		{ /* '11095' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11096:
		{ /* '11096' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 11097:
		{ /* '11097' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 17051:
		{ /* '17051' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 17052:
		{ /* '17052' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 17053:
		{ /* '17053' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 17054:
		{ /* '17054' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 17055:
		{ /* '17055' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 17056:
		{ /* '17056' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 17057:
		{ /* '17057' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 17058:
		{ /* '17058' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 17059:
		{ /* '17059' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 17060:
		{ /* '17060' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 17061:
		{ /* '17061' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 19051:
		{ /* '19051' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 19056:
		{ /* '19056' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 19057:
		{ /* '19057' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 19058:
		{ /* '19058' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 19059:
		{ /* '19059' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 19060:
		{ /* '19060' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 19061:
		{ /* '19061' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 19062:
		{ /* '19062' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 19063:
		{ /* '19063' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 409110:
		{ /* '409110' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 409111:
		{ /* '409111' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 410110:
		{ /* '410110' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 410111:
		{ /* '410111' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417101:
		{ /* '417101' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417102:
		{ /* '417102' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417103:
		{ /* '417103' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417104:
		{ /* '417104' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417105:
		{ /* '417105' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417106:
		{ /* '417106' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417107:
		{ /* '417107' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417108:
		{ /* '417108' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417109:
		{ /* '417109' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417110:
		{ /* '417110' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417111:
		{ /* '417111' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417112:
		{ /* '417112' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417113:
		{ /* '417113' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417114:
		{ /* '417114' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417115:
		{ /* '417115' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417116:
		{ /* '417116' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417117:
		{ /* '417117' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417118:
		{ /* '417118' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417119:
		{ /* '417119' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 417120:
		{ /* '417120' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418101:
		{ /* '418101' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418102:
		{ /* '418102' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418103:
		{ /* '418103' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418104:
		{ /* '418104' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418105:
		{ /* '418105' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418106:
		{ /* '418106' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418107:
		{ /* '418107' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418108:
		{ /* '418108' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418109:
		{ /* '418109' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418110:
		{ /* '418110' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418111:
		{ /* '418111' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418112:
		{ /* '418112' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418113:
		{ /* '418113' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418114:
		{ /* '418114' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418115:
		{ /* '418115' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418116:
		{ /* '418116' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418117:
		{ /* '418117' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418118:
		{ /* '418118' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418119:
		{ /* '418119' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418120:
		{ /* '418120' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418121:
		{ /* '418121' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418122:
		{ /* '418122' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418123:
		{ /* '418123' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418124:
		{ /* '418124' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418125:
		{ /* '418125' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418126:
		{ /* '418126' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418127:
		{ /* '418127' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418128:
		{ /* '418128' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418129:
		{ /* '418129' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418130:
		{ /* '418130' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418131:
		{ /* '418131' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418132:
		{ /* '418132' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418133:
		{ /* '418133' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 418134:
		{ /* '418134' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 420051:
		{ /* '420051' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 420052:
		{ /* '420052' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 420053:
		{ /* '420053' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 421051:
		{ /* '421051' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 421052:
		{ /* '421052' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 51:
		{ /* '51' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 52:
		{ /* '52' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 53:
		{ /* '53' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 54:
		{ /* '54' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 55:
		{ /* '55' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 56:
		{ /* '56' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 57:
		{ /* '57' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 58:
		{ /* '58' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 59:
		{ /* '59' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 60:
		{ /* '60' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 6051:
		{ /* '6051' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 6052:
		{ /* '6052' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 6053:
		{ /* '6053' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 6054:
		{ /* '6054' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 6055:
		{ /* '6055' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 6056:
		{ /* '6056' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 6057:
		{ /* '6057' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 6058:
		{ /* '6058' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 6063:
		{ /* '6063' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 6067:
		{ /* '6067' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 61:
		{ /* '61' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 6112:
		{ /* '6112' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 62:
		{ /* '62' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 63:
		{ /* '63' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 64:
		{ /* '64' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65:
		{ /* '65' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535001:
		{ /* '65535001' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535002:
		{ /* '65535002' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535003:
		{ /* '65535003' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535004:
		{ /* '65535004' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535005:
		{ /* '65535005' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535006:
		{ /* '65535006' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535007:
		{ /* '65535007' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535008:
		{ /* '65535008' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535009:
		{ /* '65535009' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535010:
		{ /* '65535010' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535011:
		{ /* '65535011' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535012:
		{ /* '65535012' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535013:
		{ /* '65535013' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535014:
		{ /* '65535014' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535015:
		{ /* '65535015' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535016:
		{ /* '65535016' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535017:
		{ /* '65535017' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535018:
		{ /* '65535018' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535019:
		{ /* '65535019' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535020:
		{ /* '65535020' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535021:
		{ /* '65535021' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535022:
		{ /* '65535022' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535023:
		{ /* '65535023' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535024:
		{ /* '65535024' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535025:
		{ /* '65535025' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535026:
		{ /* '65535026' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535027:
		{ /* '65535027' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535028:
		{ /* '65535028' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535029:
		{ /* '65535029' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 65535030:
		{ /* '65535030' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 66:
		{ /* '66' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 67:
		{ /* '67' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 68:
		{ /* '68' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 69:
		{ /* '69' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 70:
		{ /* '70' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 71:
		{ /* '71' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 72:
		{ /* '72' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 73:
		{ /* '73' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 74:
		{ /* '74' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 75:
		{ /* '75' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 76:
		{ /* '76' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 77:
		{ /* '77' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 78:
		{ /* '78' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 79:
		{ /* '79' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 80:
		{ /* '80' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800051:
		{ /* '800051' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800052:
		{ /* '800052' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800053:
		{ /* '800053' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800054:
		{ /* '800054' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800055:
		{ /* '800055' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800057:
		{ /* '800057' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800058:
		{ /* '800058' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800060:
		{ /* '800060' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800061:
		{ /* '800061' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800062:
		{ /* '800062' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800063:
		{ /* '800063' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800064:
		{ /* '800064' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800065:
		{ /* '800065' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800066:
		{ /* '800066' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800067:
		{ /* '800067' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800068:
		{ /* '800068' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 800069:
		{ /* '800069' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 801051:
		{ /* '801051' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 801052:
		{ /* '801052' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 801053:
		{ /* '801053' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 81:
		{ /* '81' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 82:
		{ /* '82' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 83:
		{ /* '83' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 84:
		{ /* '84' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 85:
		{ /* '85' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 86:
		{ /* '86' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 9051:
		{ /* '9051' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	case 9052:
		{ /* '9052' */
			return KnxPropertyDataType_PDT_CONTROL
		}
	default:
		{
			return 0
		}
	}
}

func (e KnxInterfaceObjectProperty) Text() string {
	switch e {
	case 11051:
		{ /* '11051' */
			return "Project Installation Identification"
		}
	case 11052:
		{ /* '11052' */
			return "KNX Individual Address"
		}
	case 11053:
		{ /* '11053' */
			return "Additional Individual Addresses"
		}
	case 11054:
		{ /* '11054' */
			return ""
		}
	case 11055:
		{ /* '11055' */
			return ""
		}
	case 11056:
		{ /* '11056' */
			return ""
		}
	case 11057:
		{ /* '11057' */
			return ""
		}
	case 11058:
		{ /* '11058' */
			return ""
		}
	case 11059:
		{ /* '11059' */
			return ""
		}
	case 11060:
		{ /* '11060' */
			return ""
		}
	case 11061:
		{ /* '11061' */
			return ""
		}
	case 11062:
		{ /* '11062' */
			return ""
		}
	case 11063:
		{ /* '11063' */
			return ""
		}
	case 11064:
		{ /* '11064' */
			return ""
		}
	case 11065:
		{ /* '11065' */
			return ""
		}
	case 11066:
		{ /* '11066' */
			return ""
		}
	case 11067:
		{ /* '11067' */
			return ""
		}
	case 11068:
		{ /* '11068' */
			return ""
		}
	case 11069:
		{ /* '11069' */
			return ""
		}
	case 11070:
		{ /* '11070' */
			return ""
		}
	case 11071:
		{ /* '11071' */
			return ""
		}
	case 11072:
		{ /* '11072' */
			return ""
		}
	case 11073:
		{ /* '11073' */
			return ""
		}
	case 11074:
		{ /* '11074' */
			return ""
		}
	case 11075:
		{ /* '11075' */
			return ""
		}
	case 11076:
		{ /* '11076' */
			return ""
		}
	case 11091:
		{ /* '11091' */
			return ""
		}
	case 11092:
		{ /* '11092' */
			return ""
		}
	case 11093:
		{ /* '11093' */
			return ""
		}
	case 11094:
		{ /* '11094' */
			return ""
		}
	case 11095:
		{ /* '11095' */
			return ""
		}
	case 11096:
		{ /* '11096' */
			return ""
		}
	case 11097:
		{ /* '11097' */
			return ""
		}
	case 17051:
		{ /* '17051' */
			return ""
		}
	case 17052:
		{ /* '17052' */
			return ""
		}
	case 17053:
		{ /* '17053' */
			return ""
		}
	case 17054:
		{ /* '17054' */
			return ""
		}
	case 17055:
		{ /* '17055' */
			return ""
		}
	case 17056:
		{ /* '17056' */
			return ""
		}
	case 17057:
		{ /* '17057' */
			return ""
		}
	case 17058:
		{ /* '17058' */
			return ""
		}
	case 17059:
		{ /* '17059' */
			return ""
		}
	case 17060:
		{ /* '17060' */
			return ""
		}
	case 17061:
		{ /* '17061' */
			return ""
		}
	case 19051:
		{ /* '19051' */
			return ""
		}
	case 19056:
		{ /* '19056' */
			return ""
		}
	case 19057:
		{ /* '19057' */
			return ""
		}
	case 19058:
		{ /* '19058' */
			return ""
		}
	case 19059:
		{ /* '19059' */
			return ""
		}
	case 19060:
		{ /* '19060' */
			return ""
		}
	case 19061:
		{ /* '19061' */
			return ""
		}
	case 19062:
		{ /* '19062' */
			return ""
		}
	case 19063:
		{ /* '19063' */
			return ""
		}
	case 409110:
		{ /* '409110' */
			return ""
		}
	case 409111:
		{ /* '409111' */
			return ""
		}
	case 410110:
		{ /* '410110' */
			return ""
		}
	case 410111:
		{ /* '410111' */
			return ""
		}
	case 417101:
		{ /* '417101' */
			return ""
		}
	case 417102:
		{ /* '417102' */
			return ""
		}
	case 417103:
		{ /* '417103' */
			return ""
		}
	case 417104:
		{ /* '417104' */
			return ""
		}
	case 417105:
		{ /* '417105' */
			return ""
		}
	case 417106:
		{ /* '417106' */
			return ""
		}
	case 417107:
		{ /* '417107' */
			return ""
		}
	case 417108:
		{ /* '417108' */
			return ""
		}
	case 417109:
		{ /* '417109' */
			return ""
		}
	case 417110:
		{ /* '417110' */
			return ""
		}
	case 417111:
		{ /* '417111' */
			return ""
		}
	case 417112:
		{ /* '417112' */
			return ""
		}
	case 417113:
		{ /* '417113' */
			return ""
		}
	case 417114:
		{ /* '417114' */
			return ""
		}
	case 417115:
		{ /* '417115' */
			return ""
		}
	case 417116:
		{ /* '417116' */
			return ""
		}
	case 417117:
		{ /* '417117' */
			return ""
		}
	case 417118:
		{ /* '417118' */
			return ""
		}
	case 417119:
		{ /* '417119' */
			return ""
		}
	case 417120:
		{ /* '417120' */
			return ""
		}
	case 418101:
		{ /* '418101' */
			return ""
		}
	case 418102:
		{ /* '418102' */
			return ""
		}
	case 418103:
		{ /* '418103' */
			return ""
		}
	case 418104:
		{ /* '418104' */
			return ""
		}
	case 418105:
		{ /* '418105' */
			return ""
		}
	case 418106:
		{ /* '418106' */
			return ""
		}
	case 418107:
		{ /* '418107' */
			return ""
		}
	case 418108:
		{ /* '418108' */
			return ""
		}
	case 418109:
		{ /* '418109' */
			return ""
		}
	case 418110:
		{ /* '418110' */
			return ""
		}
	case 418111:
		{ /* '418111' */
			return ""
		}
	case 418112:
		{ /* '418112' */
			return ""
		}
	case 418113:
		{ /* '418113' */
			return ""
		}
	case 418114:
		{ /* '418114' */
			return ""
		}
	case 418115:
		{ /* '418115' */
			return ""
		}
	case 418116:
		{ /* '418116' */
			return ""
		}
	case 418117:
		{ /* '418117' */
			return ""
		}
	case 418118:
		{ /* '418118' */
			return ""
		}
	case 418119:
		{ /* '418119' */
			return ""
		}
	case 418120:
		{ /* '418120' */
			return ""
		}
	case 418121:
		{ /* '418121' */
			return ""
		}
	case 418122:
		{ /* '418122' */
			return ""
		}
	case 418123:
		{ /* '418123' */
			return ""
		}
	case 418124:
		{ /* '418124' */
			return ""
		}
	case 418125:
		{ /* '418125' */
			return ""
		}
	case 418126:
		{ /* '418126' */
			return ""
		}
	case 418127:
		{ /* '418127' */
			return ""
		}
	case 418128:
		{ /* '418128' */
			return ""
		}
	case 418129:
		{ /* '418129' */
			return ""
		}
	case 418130:
		{ /* '418130' */
			return ""
		}
	case 418131:
		{ /* '418131' */
			return ""
		}
	case 418132:
		{ /* '418132' */
			return ""
		}
	case 418133:
		{ /* '418133' */
			return ""
		}
	case 418134:
		{ /* '418134' */
			return ""
		}
	case 420051:
		{ /* '420051' */
			return ""
		}
	case 420052:
		{ /* '420052' */
			return ""
		}
	case 420053:
		{ /* '420053' */
			return ""
		}
	case 421051:
		{ /* '421051' */
			return ""
		}
	case 421052:
		{ /* '421052' */
			return ""
		}
	case 51:
		{ /* '51' */
			return "Routing Count"
		}
	case 52:
		{ /* '52' */
			return "Maximum Retry Count"
		}
	case 53:
		{ /* '53' */
			return "Error Flags"
		}
	case 54:
		{ /* '54' */
			return "Programming Mode"
		}
	case 55:
		{ /* '55' */
			return "Product Identification"
		}
	case 56:
		{ /* '56' */
			return "Max. APDU-Length"
		}
	case 57:
		{ /* '57' */
			return "Subnetwork Address"
		}
	case 58:
		{ /* '58' */
			return "Device Address"
		}
	case 59:
		{ /* '59' */
			return "Config Link"
		}
	case 60:
		{ /* '60' */
			return ""
		}
	case 6051:
		{ /* '6051' */
			return ""
		}
	case 6052:
		{ /* '6052' */
			return ""
		}
	case 6053:
		{ /* '6053' */
			return ""
		}
	case 6054:
		{ /* '6054' */
			return ""
		}
	case 6055:
		{ /* '6055' */
			return ""
		}
	case 6056:
		{ /* '6056' */
			return ""
		}
	case 6057:
		{ /* '6057' */
			return ""
		}
	case 6058:
		{ /* '6058' */
			return ""
		}
	case 6063:
		{ /* '6063' */
			return ""
		}
	case 6067:
		{ /* '6067' */
			return ""
		}
	case 61:
		{ /* '61' */
			return ""
		}
	case 6112:
		{ /* '6112' */
			return ""
		}
	case 62:
		{ /* '62' */
			return ""
		}
	case 63:
		{ /* '63' */
			return ""
		}
	case 64:
		{ /* '64' */
			return ""
		}
	case 65:
		{ /* '65' */
			return ""
		}
	case 65535001:
		{ /* '65535001' */
			return "Interface Object Type"
		}
	case 65535002:
		{ /* '65535002' */
			return "Interface Object Name"
		}
	case 65535003:
		{ /* '65535003' */
			return "Semaphor"
		}
	case 65535004:
		{ /* '65535004' */
			return "Group Object Reference"
		}
	case 65535005:
		{ /* '65535005' */
			return "Load Control"
		}
	case 65535006:
		{ /* '65535006' */
			return "Run Control"
		}
	case 65535007:
		{ /* '65535007' */
			return "Table Reference"
		}
	case 65535008:
		{ /* '65535008' */
			return "Service Control"
		}
	case 65535009:
		{ /* '65535009' */
			return "Firmware Revision"
		}
	case 65535010:
		{ /* '65535010' */
			return "Services Supported"
		}
	case 65535011:
		{ /* '65535011' */
			return "KNX Serial Number"
		}
	case 65535012:
		{ /* '65535012' */
			return "Manufacturer Identifier"
		}
	case 65535013:
		{ /* '65535013' */
			return "Application Version"
		}
	case 65535014:
		{ /* '65535014' */
			return "Device Control"
		}
	case 65535015:
		{ /* '65535015' */
			return "Order Info"
		}
	case 65535016:
		{ /* '65535016' */
			return "PEI Type"
		}
	case 65535017:
		{ /* '65535017' */
			return "PortADDR"
		}
	case 65535018:
		{ /* '65535018' */
			return "Polling Group Settings"
		}
	case 65535019:
		{ /* '65535019' */
			return "Manufacturer Data"
		}
	case 65535020:
		{ /* '65535020' */
			return ""
		}
	case 65535021:
		{ /* '65535021' */
			return "Description"
		}
	case 65535022:
		{ /* '65535022' */
			return ""
		}
	case 65535023:
		{ /* '65535023' */
			return "Table"
		}
	case 65535024:
		{ /* '65535024' */
			return "Interface Object Link"
		}
	case 65535025:
		{ /* '65535025' */
			return "Version"
		}
	case 65535026:
		{ /* '65535026' */
			return "Group Address Assignment"
		}
	case 65535027:
		{ /* '65535027' */
			return "Memory Control Table"
		}
	case 65535028:
		{ /* '65535028' */
			return "Error Code"
		}
	case 65535029:
		{ /* '65535029' */
			return "Object Index"
		}
	case 65535030:
		{ /* '65535030' */
			return "Download Counter"
		}
	case 66:
		{ /* '66' */
			return ""
		}
	case 67:
		{ /* '67' */
			return ""
		}
	case 68:
		{ /* '68' */
			return ""
		}
	case 69:
		{ /* '69' */
			return ""
		}
	case 70:
		{ /* '70' */
			return "Domain Address"
		}
	case 71:
		{ /* '71' */
			return ""
		}
	case 72:
		{ /* '72' */
			return "Management Descriptor 1"
		}
	case 73:
		{ /* '73' */
			return "PL110 Parameters"
		}
	case 74:
		{ /* '74' */
			return ""
		}
	case 75:
		{ /* '75' */
			return ""
		}
	case 76:
		{ /* '76' */
			return ""
		}
	case 77:
		{ /* '77' */
			return ""
		}
	case 78:
		{ /* '78' */
			return "Hardware Type"
		}
	case 79:
		{ /* '79' */
			return ""
		}
	case 80:
		{ /* '80' */
			return ""
		}
	case 800051:
		{ /* '800051' */
			return ""
		}
	case 800052:
		{ /* '800052' */
			return ""
		}
	case 800053:
		{ /* '800053' */
			return ""
		}
	case 800054:
		{ /* '800054' */
			return ""
		}
	case 800055:
		{ /* '800055' */
			return ""
		}
	case 800057:
		{ /* '800057' */
			return ""
		}
	case 800058:
		{ /* '800058' */
			return ""
		}
	case 800060:
		{ /* '800060' */
			return ""
		}
	case 800061:
		{ /* '800061' */
			return ""
		}
	case 800062:
		{ /* '800062' */
			return ""
		}
	case 800063:
		{ /* '800063' */
			return ""
		}
	case 800064:
		{ /* '800064' */
			return ""
		}
	case 800065:
		{ /* '800065' */
			return ""
		}
	case 800066:
		{ /* '800066' */
			return ""
		}
	case 800067:
		{ /* '800067' */
			return ""
		}
	case 800068:
		{ /* '800068' */
			return ""
		}
	case 800069:
		{ /* '800069' */
			return ""
		}
	case 801051:
		{ /* '801051' */
			return ""
		}
	case 801052:
		{ /* '801052' */
			return ""
		}
	case 801053:
		{ /* '801053' */
			return ""
		}
	case 81:
		{ /* '81' */
			return ""
		}
	case 82:
		{ /* '82' */
			return "RF Domain Address"
		}
	case 83:
		{ /* '83' */
			return ""
		}
	case 84:
		{ /* '84' */
			return ""
		}
	case 85:
		{ /* '85' */
			return ""
		}
	case 86:
		{ /* '86' */
			return ""
		}
	case 9051:
		{ /* '9051' */
			return ""
		}
	case 9052:
		{ /* '9052' */
			return ""
		}
	default:
		{
			return ""
		}
	}
}

func (e KnxInterfaceObjectProperty) PropertyId() uint16 {
	switch e {
	case 11051:
		{ /* '11051' */
			return 51
		}
	case 11052:
		{ /* '11052' */
			return 52
		}
	case 11053:
		{ /* '11053' */
			return 53
		}
	case 11054:
		{ /* '11054' */
			return 54
		}
	case 11055:
		{ /* '11055' */
			return 55
		}
	case 11056:
		{ /* '11056' */
			return 56
		}
	case 11057:
		{ /* '11057' */
			return 57
		}
	case 11058:
		{ /* '11058' */
			return 58
		}
	case 11059:
		{ /* '11059' */
			return 59
		}
	case 11060:
		{ /* '11060' */
			return 60
		}
	case 11061:
		{ /* '11061' */
			return 61
		}
	case 11062:
		{ /* '11062' */
			return 62
		}
	case 11063:
		{ /* '11063' */
			return 63
		}
	case 11064:
		{ /* '11064' */
			return 64
		}
	case 11065:
		{ /* '11065' */
			return 65
		}
	case 11066:
		{ /* '11066' */
			return 66
		}
	case 11067:
		{ /* '11067' */
			return 67
		}
	case 11068:
		{ /* '11068' */
			return 68
		}
	case 11069:
		{ /* '11069' */
			return 69
		}
	case 11070:
		{ /* '11070' */
			return 70
		}
	case 11071:
		{ /* '11071' */
			return 71
		}
	case 11072:
		{ /* '11072' */
			return 72
		}
	case 11073:
		{ /* '11073' */
			return 73
		}
	case 11074:
		{ /* '11074' */
			return 74
		}
	case 11075:
		{ /* '11075' */
			return 75
		}
	case 11076:
		{ /* '11076' */
			return 76
		}
	case 11091:
		{ /* '11091' */
			return 91
		}
	case 11092:
		{ /* '11092' */
			return 92
		}
	case 11093:
		{ /* '11093' */
			return 93
		}
	case 11094:
		{ /* '11094' */
			return 94
		}
	case 11095:
		{ /* '11095' */
			return 95
		}
	case 11096:
		{ /* '11096' */
			return 96
		}
	case 11097:
		{ /* '11097' */
			return 97
		}
	case 17051:
		{ /* '17051' */
			return 51
		}
	case 17052:
		{ /* '17052' */
			return 52
		}
	case 17053:
		{ /* '17053' */
			return 53
		}
	case 17054:
		{ /* '17054' */
			return 54
		}
	case 17055:
		{ /* '17055' */
			return 55
		}
	case 17056:
		{ /* '17056' */
			return 56
		}
	case 17057:
		{ /* '17057' */
			return 57
		}
	case 17058:
		{ /* '17058' */
			return 58
		}
	case 17059:
		{ /* '17059' */
			return 59
		}
	case 17060:
		{ /* '17060' */
			return 60
		}
	case 17061:
		{ /* '17061' */
			return 61
		}
	case 19051:
		{ /* '19051' */
			return 51
		}
	case 19056:
		{ /* '19056' */
			return 56
		}
	case 19057:
		{ /* '19057' */
			return 57
		}
	case 19058:
		{ /* '19058' */
			return 58
		}
	case 19059:
		{ /* '19059' */
			return 59
		}
	case 19060:
		{ /* '19060' */
			return 60
		}
	case 19061:
		{ /* '19061' */
			return 61
		}
	case 19062:
		{ /* '19062' */
			return 62
		}
	case 19063:
		{ /* '19063' */
			return 63
		}
	case 409110:
		{ /* '409110' */
			return 110
		}
	case 409111:
		{ /* '409111' */
			return 111
		}
	case 410110:
		{ /* '410110' */
			return 110
		}
	case 410111:
		{ /* '410111' */
			return 111
		}
	case 417101:
		{ /* '417101' */
			return 101
		}
	case 417102:
		{ /* '417102' */
			return 102
		}
	case 417103:
		{ /* '417103' */
			return 103
		}
	case 417104:
		{ /* '417104' */
			return 104
		}
	case 417105:
		{ /* '417105' */
			return 105
		}
	case 417106:
		{ /* '417106' */
			return 106
		}
	case 417107:
		{ /* '417107' */
			return 107
		}
	case 417108:
		{ /* '417108' */
			return 108
		}
	case 417109:
		{ /* '417109' */
			return 109
		}
	case 417110:
		{ /* '417110' */
			return 110
		}
	case 417111:
		{ /* '417111' */
			return 111
		}
	case 417112:
		{ /* '417112' */
			return 112
		}
	case 417113:
		{ /* '417113' */
			return 113
		}
	case 417114:
		{ /* '417114' */
			return 114
		}
	case 417115:
		{ /* '417115' */
			return 115
		}
	case 417116:
		{ /* '417116' */
			return 116
		}
	case 417117:
		{ /* '417117' */
			return 117
		}
	case 417118:
		{ /* '417118' */
			return 118
		}
	case 417119:
		{ /* '417119' */
			return 119
		}
	case 417120:
		{ /* '417120' */
			return 120
		}
	case 418101:
		{ /* '418101' */
			return 101
		}
	case 418102:
		{ /* '418102' */
			return 102
		}
	case 418103:
		{ /* '418103' */
			return 103
		}
	case 418104:
		{ /* '418104' */
			return 104
		}
	case 418105:
		{ /* '418105' */
			return 105
		}
	case 418106:
		{ /* '418106' */
			return 106
		}
	case 418107:
		{ /* '418107' */
			return 107
		}
	case 418108:
		{ /* '418108' */
			return 108
		}
	case 418109:
		{ /* '418109' */
			return 109
		}
	case 418110:
		{ /* '418110' */
			return 110
		}
	case 418111:
		{ /* '418111' */
			return 111
		}
	case 418112:
		{ /* '418112' */
			return 112
		}
	case 418113:
		{ /* '418113' */
			return 113
		}
	case 418114:
		{ /* '418114' */
			return 114
		}
	case 418115:
		{ /* '418115' */
			return 115
		}
	case 418116:
		{ /* '418116' */
			return 116
		}
	case 418117:
		{ /* '418117' */
			return 117
		}
	case 418118:
		{ /* '418118' */
			return 118
		}
	case 418119:
		{ /* '418119' */
			return 119
		}
	case 418120:
		{ /* '418120' */
			return 120
		}
	case 418121:
		{ /* '418121' */
			return 121
		}
	case 418122:
		{ /* '418122' */
			return 122
		}
	case 418123:
		{ /* '418123' */
			return 123
		}
	case 418124:
		{ /* '418124' */
			return 124
		}
	case 418125:
		{ /* '418125' */
			return 125
		}
	case 418126:
		{ /* '418126' */
			return 126
		}
	case 418127:
		{ /* '418127' */
			return 127
		}
	case 418128:
		{ /* '418128' */
			return 128
		}
	case 418129:
		{ /* '418129' */
			return 129
		}
	case 418130:
		{ /* '418130' */
			return 130
		}
	case 418131:
		{ /* '418131' */
			return 131
		}
	case 418132:
		{ /* '418132' */
			return 132
		}
	case 418133:
		{ /* '418133' */
			return 133
		}
	case 418134:
		{ /* '418134' */
			return 134
		}
	case 420051:
		{ /* '420051' */
			return 51
		}
	case 420052:
		{ /* '420052' */
			return 52
		}
	case 420053:
		{ /* '420053' */
			return 53
		}
	case 421051:
		{ /* '421051' */
			return 51
		}
	case 421052:
		{ /* '421052' */
			return 52
		}
	case 51:
		{ /* '51' */
			return 51
		}
	case 52:
		{ /* '52' */
			return 52
		}
	case 53:
		{ /* '53' */
			return 53
		}
	case 54:
		{ /* '54' */
			return 54
		}
	case 55:
		{ /* '55' */
			return 55
		}
	case 56:
		{ /* '56' */
			return 56
		}
	case 57:
		{ /* '57' */
			return 57
		}
	case 58:
		{ /* '58' */
			return 58
		}
	case 59:
		{ /* '59' */
			return 59
		}
	case 60:
		{ /* '60' */
			return 60
		}
	case 6051:
		{ /* '6051' */
			return 51
		}
	case 6052:
		{ /* '6052' */
			return 52
		}
	case 6053:
		{ /* '6053' */
			return 53
		}
	case 6054:
		{ /* '6054' */
			return 54
		}
	case 6055:
		{ /* '6055' */
			return 55
		}
	case 6056:
		{ /* '6056' */
			return 56
		}
	case 6057:
		{ /* '6057' */
			return 57
		}
	case 6058:
		{ /* '6058' */
			return 58
		}
	case 6063:
		{ /* '6063' */
			return 63
		}
	case 6067:
		{ /* '6067' */
			return 67
		}
	case 61:
		{ /* '61' */
			return 61
		}
	case 6112:
		{ /* '6112' */
			return 112
		}
	case 62:
		{ /* '62' */
			return 62
		}
	case 63:
		{ /* '63' */
			return 63
		}
	case 64:
		{ /* '64' */
			return 64
		}
	case 65:
		{ /* '65' */
			return 65
		}
	case 65535001:
		{ /* '65535001' */
			return 1
		}
	case 65535002:
		{ /* '65535002' */
			return 2
		}
	case 65535003:
		{ /* '65535003' */
			return 3
		}
	case 65535004:
		{ /* '65535004' */
			return 4
		}
	case 65535005:
		{ /* '65535005' */
			return 5
		}
	case 65535006:
		{ /* '65535006' */
			return 6
		}
	case 65535007:
		{ /* '65535007' */
			return 7
		}
	case 65535008:
		{ /* '65535008' */
			return 8
		}
	case 65535009:
		{ /* '65535009' */
			return 9
		}
	case 65535010:
		{ /* '65535010' */
			return 10
		}
	case 65535011:
		{ /* '65535011' */
			return 11
		}
	case 65535012:
		{ /* '65535012' */
			return 12
		}
	case 65535013:
		{ /* '65535013' */
			return 13
		}
	case 65535014:
		{ /* '65535014' */
			return 14
		}
	case 65535015:
		{ /* '65535015' */
			return 15
		}
	case 65535016:
		{ /* '65535016' */
			return 16
		}
	case 65535017:
		{ /* '65535017' */
			return 17
		}
	case 65535018:
		{ /* '65535018' */
			return 18
		}
	case 65535019:
		{ /* '65535019' */
			return 19
		}
	case 65535020:
		{ /* '65535020' */
			return 20
		}
	case 65535021:
		{ /* '65535021' */
			return 21
		}
	case 65535022:
		{ /* '65535022' */
			return 22
		}
	case 65535023:
		{ /* '65535023' */
			return 23
		}
	case 65535024:
		{ /* '65535024' */
			return 24
		}
	case 65535025:
		{ /* '65535025' */
			return 25
		}
	case 65535026:
		{ /* '65535026' */
			return 26
		}
	case 65535027:
		{ /* '65535027' */
			return 27
		}
	case 65535028:
		{ /* '65535028' */
			return 28
		}
	case 65535029:
		{ /* '65535029' */
			return 29
		}
	case 65535030:
		{ /* '65535030' */
			return 30
		}
	case 66:
		{ /* '66' */
			return 66
		}
	case 67:
		{ /* '67' */
			return 67
		}
	case 68:
		{ /* '68' */
			return 68
		}
	case 69:
		{ /* '69' */
			return 69
		}
	case 70:
		{ /* '70' */
			return 70
		}
	case 71:
		{ /* '71' */
			return 71
		}
	case 72:
		{ /* '72' */
			return 72
		}
	case 73:
		{ /* '73' */
			return 73
		}
	case 74:
		{ /* '74' */
			return 74
		}
	case 75:
		{ /* '75' */
			return 75
		}
	case 76:
		{ /* '76' */
			return 76
		}
	case 77:
		{ /* '77' */
			return 77
		}
	case 78:
		{ /* '78' */
			return 78
		}
	case 79:
		{ /* '79' */
			return 79
		}
	case 80:
		{ /* '80' */
			return 80
		}
	case 800051:
		{ /* '800051' */
			return 51
		}
	case 800052:
		{ /* '800052' */
			return 52
		}
	case 800053:
		{ /* '800053' */
			return 53
		}
	case 800054:
		{ /* '800054' */
			return 54
		}
	case 800055:
		{ /* '800055' */
			return 55
		}
	case 800057:
		{ /* '800057' */
			return 57
		}
	case 800058:
		{ /* '800058' */
			return 58
		}
	case 800060:
		{ /* '800060' */
			return 60
		}
	case 800061:
		{ /* '800061' */
			return 61
		}
	case 800062:
		{ /* '800062' */
			return 62
		}
	case 800063:
		{ /* '800063' */
			return 63
		}
	case 800064:
		{ /* '800064' */
			return 64
		}
	case 800065:
		{ /* '800065' */
			return 65
		}
	case 800066:
		{ /* '800066' */
			return 66
		}
	case 800067:
		{ /* '800067' */
			return 67
		}
	case 800068:
		{ /* '800068' */
			return 68
		}
	case 800069:
		{ /* '800069' */
			return 69
		}
	case 801051:
		{ /* '801051' */
			return 51
		}
	case 801052:
		{ /* '801052' */
			return 52
		}
	case 801053:
		{ /* '801053' */
			return 53
		}
	case 81:
		{ /* '81' */
			return 81
		}
	case 82:
		{ /* '82' */
			return 82
		}
	case 83:
		{ /* '83' */
			return 83
		}
	case 84:
		{ /* '84' */
			return 84
		}
	case 85:
		{ /* '85' */
			return 85
		}
	case 86:
		{ /* '86' */
			return 86
		}
	case 9051:
		{ /* '9051' */
			return 51
		}
	case 9052:
		{ /* '9052' */
			return 52
		}
	default:
		{
			return 0
		}
	}
}

func (e KnxInterfaceObjectProperty) ObjectType() KnxInterfaceObjectType {
	switch e {
	case 11051:
		{ /* '11051' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11052:
		{ /* '11052' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11053:
		{ /* '11053' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11054:
		{ /* '11054' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11055:
		{ /* '11055' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11056:
		{ /* '11056' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11057:
		{ /* '11057' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11058:
		{ /* '11058' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11059:
		{ /* '11059' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11060:
		{ /* '11060' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11061:
		{ /* '11061' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11062:
		{ /* '11062' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11063:
		{ /* '11063' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11064:
		{ /* '11064' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11065:
		{ /* '11065' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11066:
		{ /* '11066' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11067:
		{ /* '11067' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11068:
		{ /* '11068' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11069:
		{ /* '11069' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11070:
		{ /* '11070' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11071:
		{ /* '11071' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11072:
		{ /* '11072' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11073:
		{ /* '11073' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11074:
		{ /* '11074' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11075:
		{ /* '11075' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11076:
		{ /* '11076' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11091:
		{ /* '11091' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11092:
		{ /* '11092' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11093:
		{ /* '11093' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11094:
		{ /* '11094' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11095:
		{ /* '11095' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11096:
		{ /* '11096' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 11097:
		{ /* '11097' */
			return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
		}
	case 17051:
		{ /* '17051' */
			return KnxInterfaceObjectType_OT_SECURITY
		}
	case 17052:
		{ /* '17052' */
			return KnxInterfaceObjectType_OT_SECURITY
		}
	case 17053:
		{ /* '17053' */
			return KnxInterfaceObjectType_OT_SECURITY
		}
	case 17054:
		{ /* '17054' */
			return KnxInterfaceObjectType_OT_SECURITY
		}
	case 17055:
		{ /* '17055' */
			return KnxInterfaceObjectType_OT_SECURITY
		}
	case 17056:
		{ /* '17056' */
			return KnxInterfaceObjectType_OT_SECURITY
		}
	case 17057:
		{ /* '17057' */
			return KnxInterfaceObjectType_OT_SECURITY
		}
	case 17058:
		{ /* '17058' */
			return KnxInterfaceObjectType_OT_SECURITY
		}
	case 17059:
		{ /* '17059' */
			return KnxInterfaceObjectType_OT_SECURITY
		}
	case 17060:
		{ /* '17060' */
			return KnxInterfaceObjectType_OT_SECURITY
		}
	case 17061:
		{ /* '17061' */
			return KnxInterfaceObjectType_OT_SECURITY
		}
	case 19051:
		{ /* '19051' */
			return KnxInterfaceObjectType_OT_RF_MEDIUM
		}
	case 19056:
		{ /* '19056' */
			return KnxInterfaceObjectType_OT_RF_MEDIUM
		}
	case 19057:
		{ /* '19057' */
			return KnxInterfaceObjectType_OT_RF_MEDIUM
		}
	case 19058:
		{ /* '19058' */
			return KnxInterfaceObjectType_OT_RF_MEDIUM
		}
	case 19059:
		{ /* '19059' */
			return KnxInterfaceObjectType_OT_RF_MEDIUM
		}
	case 19060:
		{ /* '19060' */
			return KnxInterfaceObjectType_OT_RF_MEDIUM
		}
	case 19061:
		{ /* '19061' */
			return KnxInterfaceObjectType_OT_RF_MEDIUM
		}
	case 19062:
		{ /* '19062' */
			return KnxInterfaceObjectType_OT_RF_MEDIUM
		}
	case 19063:
		{ /* '19063' */
			return KnxInterfaceObjectType_OT_RF_MEDIUM
		}
	case 409110:
		{ /* '409110' */
			return KnxInterfaceObjectType_OT_INDOOR_BRIGHTNESS_SENSOR
		}
	case 409111:
		{ /* '409111' */
			return KnxInterfaceObjectType_OT_INDOOR_BRIGHTNESS_SENSOR
		}
	case 410110:
		{ /* '410110' */
			return KnxInterfaceObjectType_OT_INDOOR_LUMINANCE_SENSOR
		}
	case 410111:
		{ /* '410111' */
			return KnxInterfaceObjectType_OT_INDOOR_LUMINANCE_SENSOR
		}
	case 417101:
		{ /* '417101' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417102:
		{ /* '417102' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417103:
		{ /* '417103' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417104:
		{ /* '417104' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417105:
		{ /* '417105' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417106:
		{ /* '417106' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417107:
		{ /* '417107' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417108:
		{ /* '417108' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417109:
		{ /* '417109' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417110:
		{ /* '417110' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417111:
		{ /* '417111' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417112:
		{ /* '417112' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417113:
		{ /* '417113' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417114:
		{ /* '417114' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417115:
		{ /* '417115' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417116:
		{ /* '417116' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417117:
		{ /* '417117' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417118:
		{ /* '417118' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417119:
		{ /* '417119' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 417120:
		{ /* '417120' */
			return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
		}
	case 418101:
		{ /* '418101' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418102:
		{ /* '418102' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418103:
		{ /* '418103' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418104:
		{ /* '418104' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418105:
		{ /* '418105' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418106:
		{ /* '418106' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418107:
		{ /* '418107' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418108:
		{ /* '418108' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418109:
		{ /* '418109' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418110:
		{ /* '418110' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418111:
		{ /* '418111' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418112:
		{ /* '418112' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418113:
		{ /* '418113' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418114:
		{ /* '418114' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418115:
		{ /* '418115' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418116:
		{ /* '418116' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418117:
		{ /* '418117' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418118:
		{ /* '418118' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418119:
		{ /* '418119' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418120:
		{ /* '418120' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418121:
		{ /* '418121' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418122:
		{ /* '418122' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418123:
		{ /* '418123' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418124:
		{ /* '418124' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418125:
		{ /* '418125' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418126:
		{ /* '418126' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418127:
		{ /* '418127' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418128:
		{ /* '418128' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418129:
		{ /* '418129' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418130:
		{ /* '418130' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418131:
		{ /* '418131' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418132:
		{ /* '418132' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418133:
		{ /* '418133' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 418134:
		{ /* '418134' */
			return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
		}
	case 420051:
		{ /* '420051' */
			return KnxInterfaceObjectType_OT_DIMMING_SENSOR_BASIC
		}
	case 420052:
		{ /* '420052' */
			return KnxInterfaceObjectType_OT_DIMMING_SENSOR_BASIC
		}
	case 420053:
		{ /* '420053' */
			return KnxInterfaceObjectType_OT_DIMMING_SENSOR_BASIC
		}
	case 421051:
		{ /* '421051' */
			return KnxInterfaceObjectType_OT_SWITCHING_SENSOR_BASIC
		}
	case 421052:
		{ /* '421052' */
			return KnxInterfaceObjectType_OT_SWITCHING_SENSOR_BASIC
		}
	case 51:
		{ /* '51' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 52:
		{ /* '52' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 53:
		{ /* '53' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 54:
		{ /* '54' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 55:
		{ /* '55' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 56:
		{ /* '56' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 57:
		{ /* '57' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 58:
		{ /* '58' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 59:
		{ /* '59' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 60:
		{ /* '60' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 6051:
		{ /* '6051' */
			return KnxInterfaceObjectType_OT_ROUTER
		}
	case 6052:
		{ /* '6052' */
			return KnxInterfaceObjectType_OT_ROUTER
		}
	case 6053:
		{ /* '6053' */
			return KnxInterfaceObjectType_OT_ROUTER
		}
	case 6054:
		{ /* '6054' */
			return KnxInterfaceObjectType_OT_ROUTER
		}
	case 6055:
		{ /* '6055' */
			return KnxInterfaceObjectType_OT_ROUTER
		}
	case 6056:
		{ /* '6056' */
			return KnxInterfaceObjectType_OT_ROUTER
		}
	case 6057:
		{ /* '6057' */
			return KnxInterfaceObjectType_OT_ROUTER
		}
	case 6058:
		{ /* '6058' */
			return KnxInterfaceObjectType_OT_ROUTER
		}
	case 6063:
		{ /* '6063' */
			return KnxInterfaceObjectType_OT_ROUTER
		}
	case 6067:
		{ /* '6067' */
			return KnxInterfaceObjectType_OT_ROUTER
		}
	case 61:
		{ /* '61' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 6112:
		{ /* '6112' */
			return KnxInterfaceObjectType_OT_ROUTER
		}
	case 62:
		{ /* '62' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 63:
		{ /* '63' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 64:
		{ /* '64' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 65:
		{ /* '65' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 65535001:
		{ /* '65535001' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535002:
		{ /* '65535002' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535003:
		{ /* '65535003' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535004:
		{ /* '65535004' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535005:
		{ /* '65535005' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535006:
		{ /* '65535006' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535007:
		{ /* '65535007' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535008:
		{ /* '65535008' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535009:
		{ /* '65535009' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535010:
		{ /* '65535010' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535011:
		{ /* '65535011' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535012:
		{ /* '65535012' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535013:
		{ /* '65535013' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535014:
		{ /* '65535014' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535015:
		{ /* '65535015' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535016:
		{ /* '65535016' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535017:
		{ /* '65535017' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535018:
		{ /* '65535018' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535019:
		{ /* '65535019' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535020:
		{ /* '65535020' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535021:
		{ /* '65535021' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535022:
		{ /* '65535022' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535023:
		{ /* '65535023' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535024:
		{ /* '65535024' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535025:
		{ /* '65535025' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535026:
		{ /* '65535026' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535027:
		{ /* '65535027' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535028:
		{ /* '65535028' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535029:
		{ /* '65535029' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 65535030:
		{ /* '65535030' */
			return KnxInterfaceObjectType_OT_GENERAL
		}
	case 66:
		{ /* '66' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 67:
		{ /* '67' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 68:
		{ /* '68' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 69:
		{ /* '69' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 70:
		{ /* '70' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 71:
		{ /* '71' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 72:
		{ /* '72' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 73:
		{ /* '73' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 74:
		{ /* '74' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 75:
		{ /* '75' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 76:
		{ /* '76' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 77:
		{ /* '77' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 78:
		{ /* '78' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 79:
		{ /* '79' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 80:
		{ /* '80' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 800051:
		{ /* '800051' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800052:
		{ /* '800052' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800053:
		{ /* '800053' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800054:
		{ /* '800054' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800055:
		{ /* '800055' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800057:
		{ /* '800057' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800058:
		{ /* '800058' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800060:
		{ /* '800060' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800061:
		{ /* '800061' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800062:
		{ /* '800062' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800063:
		{ /* '800063' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800064:
		{ /* '800064' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800065:
		{ /* '800065' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800066:
		{ /* '800066' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800067:
		{ /* '800067' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800068:
		{ /* '800068' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 800069:
		{ /* '800069' */
			return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
		}
	case 801051:
		{ /* '801051' */
			return KnxInterfaceObjectType_OT_SUNBLIND_SENSOR_BASIC
		}
	case 801052:
		{ /* '801052' */
			return KnxInterfaceObjectType_OT_SUNBLIND_SENSOR_BASIC
		}
	case 801053:
		{ /* '801053' */
			return KnxInterfaceObjectType_OT_SUNBLIND_SENSOR_BASIC
		}
	case 81:
		{ /* '81' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 82:
		{ /* '82' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 83:
		{ /* '83' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 84:
		{ /* '84' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 85:
		{ /* '85' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 86:
		{ /* '86' */
			return KnxInterfaceObjectType_OT_DEVICE
		}
	case 9051:
		{ /* '9051' */
			return KnxInterfaceObjectType_OT_GROUP_OBJECT_TABLE
		}
	case 9052:
		{ /* '9052' */
			return KnxInterfaceObjectType_OT_GROUP_OBJECT_TABLE
		}
	default:
		{
			return 0
		}
	}
}
func KnxInterfaceObjectPropertyByValue(value uint32) KnxInterfaceObjectProperty {
	switch value {
	case 11051:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_PROJECT_INSTALLATION_ID
	case 11052:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNX_INDIVIDUAL_ADDRESS
	case 11053:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_ADDITIONAL_INDIVIDUAL_ADDRESSES
	case 11054:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_IP_ASSIGNMENT_METHOD
	case 11055:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_IP_ASSIGNMENT_METHOD
	case 11056:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_IP_CAPABILITIES
	case 11057:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_IP_ADDRESS
	case 11058:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_SUBNET_MASK
	case 11059:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_DEFAULT_GATEWAY
	case 11060:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_IP_ADDRESS
	case 11061:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SUBNET_MASK
	case 11062:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_DEFAULT_GATEWAY
	case 11063:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_DHCP_BOOTP_SERVER
	case 11064:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MAC_ADDRESS
	case 11065:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SYSTEM_SETUP_MULTICAST_ADDRESS
	case 11066:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_ROUTING_MULTICAST_ADDRESS
	case 11067:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_TTL
	case 11068:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNXNETIP_DEVICE_CAPABILITIES
	case 11069:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNXNETIP_DEVICE_STATE
	case 11070:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNXNETIP_ROUTING_CAPABILITIES
	case 11071:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_PRIORITY_FIFO_ENABLED
	case 11072:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_QUEUE_OVERFLOW_TO_IP
	case 11073:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_QUEUE_OVERFLOW_TO_KNX
	case 11074:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MSG_TRANSMIT_TO_IP
	case 11075:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MSG_TRANSMIT_TO_KNX
	case 11076:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_FRIENDLY_NAME
	case 11091:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_BACKBONE_KEY
	case 11092:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_DEVICE_AUTHENTICATION_CODE
	case 11093:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_PASSWORD_HASHES
	case 11094:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SECURED_SERVICE_FAMILIES
	case 11095:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MULTICAST_LATENCY_TOLERANCE
	case 11096:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SYNC_LATENCY_FRACTION
	case 11097:
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_TUNNELLING_USERS
	case 17051:
		return KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_MODE
	case 17052:
		return KnxInterfaceObjectProperty_PID_SECURITY_P2P_KEY_TABLE
	case 17053:
		return KnxInterfaceObjectProperty_PID_SECURITY_GRP_KEY_TABLE
	case 17054:
		return KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_INDIVIDUAL_ADDRESS_TABLE
	case 17055:
		return KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_FAILURES_LOG
	case 17056:
		return KnxInterfaceObjectProperty_PID_SECURITY_SKI_TOOL
	case 17057:
		return KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_REPORT
	case 17058:
		return KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_REPORT_CONTROL
	case 17059:
		return KnxInterfaceObjectProperty_PID_SECURITY_SEQUENCE_NUMBER_SENDING
	case 17060:
		return KnxInterfaceObjectProperty_PID_SECURITY_ZONE_KEYS_TABLE
	case 17061:
		return KnxInterfaceObjectProperty_PID_SECURITY_GO_SECURITY_FLAGS
	case 19051:
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_MULTI_TYPE
	case 19056:
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DOMAIN_ADDRESS
	case 19057:
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_RETRANSMITTER
	case 19058:
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_SECURITY_REPORT_CONTROL
	case 19059:
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_FILTERING_MODE_SELECT
	case 19060:
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_BIDIR_TIMEOUT
	case 19061:
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DIAG_SA_FILTER_TABLE
	case 19062:
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DIAG_QUALITY_TABLE
	case 19063:
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DIAG_PROBE
	case 409110:
		return KnxInterfaceObjectProperty_PID_INDOOR_BRIGHTNESS_SENSOR_CHANGE_OF_VALUE
	case 409111:
		return KnxInterfaceObjectProperty_PID_INDOOR_BRIGHTNESS_SENSOR_REPETITION_TIME
	case 410110:
		return KnxInterfaceObjectProperty_PID_INDOOR_LUMINANCE_SENSOR_CHANGE_OF_VALUE
	case 410111:
		return KnxInterfaceObjectProperty_PID_INDOOR_LUMINANCE_SENSOR_REPETITION_TIME
	case 417101:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_ON_DELAY
	case 417102:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_OFF_DELAY
	case 417103:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TIMED_ON_DURATION
	case 417104:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_PREWARNING_DURATION
	case 417105:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TRANSMISSION_CYCLE_TIME
	case 417106:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BUS_POWER_UP_MESSAGE_DELAY
	case 417107:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_AT_LOCKING
	case 417108:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_AT_UNLOCKING
	case 417109:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP
	case 417110:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_DOWN
	case 417111:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_INVERT_OUTPUT_STATE
	case 417112:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TIMED_ON_RETRIGGER_FUNCTION
	case 417113:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_MANUAL_OFF_ENABLE
	case 417114:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_INVERT_LOCK_DEVICE
	case 417115:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_LOCK_STATE
	case 417116:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_UNLOCK_STATE
	case 417117:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_STATE_FOR_SCENE_NUMBER
	case 417118:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_STORAGE_FUNCTION_FOR_SCENE
	case 417119:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BUS_POWER_UP_STATE
	case 417120:
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP_2
	case 418101:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_ON_DELAY
	case 418102:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_OFF_DELAY
	case 418103:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_SWITCH_OFF_BRIGHTNESS_DELAY_TIME
	case 418104:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_TIMED_ON_DURATION
	case 418105:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_PREWARNING_DURATION
	case 418106:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_TRANSMISSION_CYCLE_TIME
	case 418107:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BUS_POWER_UP_MESSAGE_DELAY
	case 418108:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED
	case 418109:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME
	case 418110:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED_FOR_SWITCH_ON_SET_VALUE
	case 418111:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED_FOR_SWITCH_OFF
	case 418112:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME_FOR_SWITCH_ON_SET_VALUE
	case 418113:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME_FOR_SWITCH_OFF
	case 418114:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_SWITCFH_OFF_BRIGHTNESS
	case 418115:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MINIMUM_SET_VALUE
	case 418116:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MAXIMUM_SET_VALUE
	case 418117:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_SWITCH_ON_SET_VALUE
	case 418118:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMM_MODE_SELECTION
	case 418119:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_RELATIV_OFF_ENABLE
	case 418120:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MEMORY_FUNCTION
	case 418121:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_TIMED_ON_RETRIGGER_FUNCTION
	case 418122:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MANUAL_OFF_ENABLE
	case 418123:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_INVERT_LOCK_DEVICE
	case 418124:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_AT_LOCKING
	case 418125:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_AT_UNLOCKING
	case 418126:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_LOCK_SETVALUE
	case 418127:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_UNLOCK_SETVALUE
	case 418128:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BIGHTNESS_FOR_SCENE
	case 418129:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_STORAGE_FUNCTION_FOR_SCENE
	case 418130:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DELTA_DIMMING_VALUE
	case 418131:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP
	case 418132:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP_SET_VALUE
	case 418133:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_DOWN
	case 418134:
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BUS_POWER_DOWN_SET_VALUE
	case 420051:
		return KnxInterfaceObjectProperty_PID_DIMMING_SENSOR_BASIC_ON_OFF_ACTION
	case 420052:
		return KnxInterfaceObjectProperty_PID_DIMMING_SENSOR_BASIC_ENABLE_TOGGLE_MODE
	case 420053:
		return KnxInterfaceObjectProperty_PID_DIMMING_SENSOR_BASIC_ABSOLUTE_SETVALUE
	case 421051:
		return KnxInterfaceObjectProperty_PID_SWITCHING_SENSOR_BASIC_ON_OFF_ACTION
	case 421052:
		return KnxInterfaceObjectProperty_PID_SWITCHING_SENSOR_BASIC_ENABLE_TOGGLE_MODE
	case 51:
		return KnxInterfaceObjectProperty_PID_DEVICE_ROUTING_COUNT
	case 52:
		return KnxInterfaceObjectProperty_PID_DEVICE_MAX_RETRY_COUNT
	case 53:
		return KnxInterfaceObjectProperty_PID_DEVICE_ERROR_FLAGS
	case 54:
		return KnxInterfaceObjectProperty_PID_DEVICE_PROGMODE
	case 55:
		return KnxInterfaceObjectProperty_PID_DEVICE_PRODUCT_ID
	case 56:
		return KnxInterfaceObjectProperty_PID_DEVICE_MAX_APDULENGTH
	case 57:
		return KnxInterfaceObjectProperty_PID_DEVICE_SUBNET_ADDR
	case 58:
		return KnxInterfaceObjectProperty_PID_DEVICE_DEVICE_ADDR
	case 59:
		return KnxInterfaceObjectProperty_PID_DEVICE_PB_CONFIG
	case 60:
		return KnxInterfaceObjectProperty_PID_DEVICE_ADDR_REPORT
	case 6051:
		return KnxInterfaceObjectProperty_PID_ROUTER_LINE_STATUS
	case 6052:
		return KnxInterfaceObjectProperty_PID_ROUTER_MAIN_LCCONFIG
	case 6053:
		return KnxInterfaceObjectProperty_PID_ROUTER_SUB_LCCONFIG
	case 6054:
		return KnxInterfaceObjectProperty_PID_ROUTER_MAIN_LCGRPCONFIG
	case 6055:
		return KnxInterfaceObjectProperty_PID_ROUTER_SUB_LCGRPCONFIG
	case 6056:
		return KnxInterfaceObjectProperty_PID_ROUTER_ROUTETABLE_CONTROL
	case 6057:
		return KnxInterfaceObjectProperty_PID_ROUTER_COUPL_SERV_CONTROL
	case 6058:
		return KnxInterfaceObjectProperty_PID_ROUTER_MAX_ROUTER_APDU_LENGTH
	case 6063:
		return KnxInterfaceObjectProperty_PID_ROUTER_MEDIUM
	case 6067:
		return KnxInterfaceObjectProperty_PID_ROUTER_FILTER_TABLE_USE
	case 61:
		return KnxInterfaceObjectProperty_PID_DEVICE_ADDR_CHECK
	case 6112:
		return KnxInterfaceObjectProperty_PID_ROUTER_RF_ENABLE_SBC
	case 62:
		return KnxInterfaceObjectProperty_PID_DEVICE_OBJECT_VALUE
	case 63:
		return KnxInterfaceObjectProperty_PID_DEVICE_OBJECTLINK
	case 64:
		return KnxInterfaceObjectProperty_PID_DEVICE_APPLICATION
	case 65:
		return KnxInterfaceObjectProperty_PID_DEVICE_PARAMETER
	case 65535001:
		return KnxInterfaceObjectProperty_PID_GENERAL_OBJECT_TYPE
	case 65535002:
		return KnxInterfaceObjectProperty_PID_GENERAL_OBJECT_NAME
	case 65535003:
		return KnxInterfaceObjectProperty_PID_GENERAL_SEMAPHOR
	case 65535004:
		return KnxInterfaceObjectProperty_PID_GENERAL_GROUP_OBJECT_REFERENCE
	case 65535005:
		return KnxInterfaceObjectProperty_PID_GENERAL_LOAD_STATE_CONTROL
	case 65535006:
		return KnxInterfaceObjectProperty_PID_GENERAL_RUN_STATE_CONTROL
	case 65535007:
		return KnxInterfaceObjectProperty_PID_GENERAL_TABLE_REFERENCE
	case 65535008:
		return KnxInterfaceObjectProperty_PID_GENERAL_SERVICE_CONTROL
	case 65535009:
		return KnxInterfaceObjectProperty_PID_GENERAL_FIRMWARE_REVISION
	case 65535010:
		return KnxInterfaceObjectProperty_PID_GENERAL_SERVICES_SUPPORTED
	case 65535011:
		return KnxInterfaceObjectProperty_PID_GENERAL_SERIAL_NUMBER
	case 65535012:
		return KnxInterfaceObjectProperty_PID_GENERAL_MANUFACTURER_ID
	case 65535013:
		return KnxInterfaceObjectProperty_PID_GENERAL_PROGRAM_VERSION
	case 65535014:
		return KnxInterfaceObjectProperty_PID_GENERAL_DEVICE_CONTROL
	case 65535015:
		return KnxInterfaceObjectProperty_PID_GENERAL_ORDER_INFO
	case 65535016:
		return KnxInterfaceObjectProperty_PID_GENERAL_PEI_TYPE
	case 65535017:
		return KnxInterfaceObjectProperty_PID_GENERAL_PORT_CONFIGURATION
	case 65535018:
		return KnxInterfaceObjectProperty_PID_GENERAL_POLL_GROUP_SETTINGS
	case 65535019:
		return KnxInterfaceObjectProperty_PID_GENERAL_MANUFACTURER_DATA
	case 65535020:
		return KnxInterfaceObjectProperty_PID_GENERAL_ENABLE
	case 65535021:
		return KnxInterfaceObjectProperty_PID_GENERAL_DESCRIPTION
	case 65535022:
		return KnxInterfaceObjectProperty_PID_GENERAL_FILE
	case 65535023:
		return KnxInterfaceObjectProperty_PID_GENERAL_TABLE
	case 65535024:
		return KnxInterfaceObjectProperty_PID_GENERAL_ENROL
	case 65535025:
		return KnxInterfaceObjectProperty_PID_GENERAL_VERSION
	case 65535026:
		return KnxInterfaceObjectProperty_PID_GENERAL_GROUP_OBJECT_LINK
	case 65535027:
		return KnxInterfaceObjectProperty_PID_GENERAL_MCB_TABLE
	case 65535028:
		return KnxInterfaceObjectProperty_PID_GENERAL_ERROR_CODE
	case 65535029:
		return KnxInterfaceObjectProperty_PID_GENERAL_OBJECT_INDEX
	case 65535030:
		return KnxInterfaceObjectProperty_PID_GENERAL_DOWNLOAD_COUNTER
	case 66:
		return KnxInterfaceObjectProperty_PID_DEVICE_OBJECTADDRESS
	case 67:
		return KnxInterfaceObjectProperty_PID_DEVICE_PSU_TYPE
	case 68:
		return KnxInterfaceObjectProperty_PID_DEVICE_PSU_STATUS
	case 69:
		return KnxInterfaceObjectProperty_PID_DEVICE_PSU_ENABLE
	case 70:
		return KnxInterfaceObjectProperty_PID_DEVICE_DOMAIN_ADDRESS
	case 71:
		return KnxInterfaceObjectProperty_PID_DEVICE_IO_LIST
	case 72:
		return KnxInterfaceObjectProperty_PID_DEVICE_MGT_DESCRIPTOR_01
	case 73:
		return KnxInterfaceObjectProperty_PID_DEVICE_PL110_PARAM
	case 74:
		return KnxInterfaceObjectProperty_PID_DEVICE_RF_REPEAT_COUNTER
	case 75:
		return KnxInterfaceObjectProperty_PID_DEVICE_RECEIVE_BLOCK_TABLE
	case 76:
		return KnxInterfaceObjectProperty_PID_DEVICE_RANDOM_PAUSE_TABLE
	case 77:
		return KnxInterfaceObjectProperty_PID_DEVICE_RECEIVE_BLOCK_NR
	case 78:
		return KnxInterfaceObjectProperty_PID_DEVICE_HARDWARE_TYPE
	case 79:
		return KnxInterfaceObjectProperty_PID_DEVICE_RETRANSMITTER_NUMBER
	case 80:
		return KnxInterfaceObjectProperty_PID_DEVICE_SERIAL_NR_TABLE
	case 800051:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REVERSION_PAUSE_TIME
	case 800052:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_UP_DOWN_TIME
	case 800053:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_SLAT_STEP_TIME
	case 800054:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_PRESET_POSITION_TIME
	case 800055:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_TO_PRESET_POSITION_IN_PERCENT
	case 800057:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_TO_PRESET_POSITION_LENGTH
	case 800058:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_PRESET_SLAT_POSITION_PERCENT
	case 800060:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_PRESET_SLAT_POSITION_ANGLE
	case 800061:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REACTION_WIND_ALARM
	case 800062:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_WIND_ALARM
	case 800063:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REACTION_ON_RAIN_ALARM
	case 800064:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_RAIN_ALARM
	case 800065:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REACTION_FROST_ALARM
	case 800066:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_FROST_ALARM
	case 800067:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MAX_SLAT_MOVE_TIME
	case 800068:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_ENABLE_BLINDS_MODE
	case 800069:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_STORAGE_FUNCTIONS_FOR_SCENE
	case 801051:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_SENSOR_BASIC_ENABLE_BLINDS_MODE
	case 801052:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_SENSOR_BASIC_UP_DOWN_ACTION
	case 801053:
		return KnxInterfaceObjectProperty_PID_SUNBLIND_SENSOR_BASIC_ENABLE_TOGGLE_MODE
	case 81:
		return KnxInterfaceObjectProperty_PID_DEVICE_BIBATMASTER_ADDRESS
	case 82:
		return KnxInterfaceObjectProperty_PID_DEVICE_RF_DOMAIN_ADDRESS
	case 83:
		return KnxInterfaceObjectProperty_PID_DEVICE_DEVICE_DESCRIPTOR
	case 84:
		return KnxInterfaceObjectProperty_PID_DEVICE_METERING_FILTER_TABLE
	case 85:
		return KnxInterfaceObjectProperty_PID_DEVICE_GROUP_TELEGR_RATE_LIMIT_TIME_BASE
	case 86:
		return KnxInterfaceObjectProperty_PID_DEVICE_GROUP_TELEGR_RATE_LIMIT_NO_OF_TELEGR
	case 9051:
		return KnxInterfaceObjectProperty_PID_GROUP_OBJECT_TABLE_GRPOBJTABLE
	case 9052:
		return KnxInterfaceObjectProperty_PID_GROUP_OBJECT_TABLE_EXT_GRPOBJREFERENCE
	}
	return 0
}

func KnxInterfaceObjectPropertyByName(value string) KnxInterfaceObjectProperty {
	switch value {
	case "PID_KNXIP_PARAMETER_PROJECT_INSTALLATION_ID":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_PROJECT_INSTALLATION_ID
	case "PID_KNXIP_PARAMETER_KNX_INDIVIDUAL_ADDRESS":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNX_INDIVIDUAL_ADDRESS
	case "PID_KNXIP_PARAMETER_ADDITIONAL_INDIVIDUAL_ADDRESSES":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_ADDITIONAL_INDIVIDUAL_ADDRESSES
	case "PID_KNXIP_PARAMETER_CURRENT_IP_ASSIGNMENT_METHOD":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_IP_ASSIGNMENT_METHOD
	case "PID_KNXIP_PARAMETER_IP_ASSIGNMENT_METHOD":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_IP_ASSIGNMENT_METHOD
	case "PID_KNXIP_PARAMETER_IP_CAPABILITIES":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_IP_CAPABILITIES
	case "PID_KNXIP_PARAMETER_CURRENT_IP_ADDRESS":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_IP_ADDRESS
	case "PID_KNXIP_PARAMETER_CURRENT_SUBNET_MASK":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_SUBNET_MASK
	case "PID_KNXIP_PARAMETER_CURRENT_DEFAULT_GATEWAY":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_DEFAULT_GATEWAY
	case "PID_KNXIP_PARAMETER_IP_ADDRESS":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_IP_ADDRESS
	case "PID_KNXIP_PARAMETER_SUBNET_MASK":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SUBNET_MASK
	case "PID_KNXIP_PARAMETER_DEFAULT_GATEWAY":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_DEFAULT_GATEWAY
	case "PID_KNXIP_PARAMETER_DHCP_BOOTP_SERVER":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_DHCP_BOOTP_SERVER
	case "PID_KNXIP_PARAMETER_MAC_ADDRESS":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MAC_ADDRESS
	case "PID_KNXIP_PARAMETER_SYSTEM_SETUP_MULTICAST_ADDRESS":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SYSTEM_SETUP_MULTICAST_ADDRESS
	case "PID_KNXIP_PARAMETER_ROUTING_MULTICAST_ADDRESS":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_ROUTING_MULTICAST_ADDRESS
	case "PID_KNXIP_PARAMETER_TTL":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_TTL
	case "PID_KNXIP_PARAMETER_KNXNETIP_DEVICE_CAPABILITIES":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNXNETIP_DEVICE_CAPABILITIES
	case "PID_KNXIP_PARAMETER_KNXNETIP_DEVICE_STATE":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNXNETIP_DEVICE_STATE
	case "PID_KNXIP_PARAMETER_KNXNETIP_ROUTING_CAPABILITIES":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNXNETIP_ROUTING_CAPABILITIES
	case "PID_KNXIP_PARAMETER_PRIORITY_FIFO_ENABLED":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_PRIORITY_FIFO_ENABLED
	case "PID_KNXIP_PARAMETER_QUEUE_OVERFLOW_TO_IP":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_QUEUE_OVERFLOW_TO_IP
	case "PID_KNXIP_PARAMETER_QUEUE_OVERFLOW_TO_KNX":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_QUEUE_OVERFLOW_TO_KNX
	case "PID_KNXIP_PARAMETER_MSG_TRANSMIT_TO_IP":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MSG_TRANSMIT_TO_IP
	case "PID_KNXIP_PARAMETER_MSG_TRANSMIT_TO_KNX":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MSG_TRANSMIT_TO_KNX
	case "PID_KNXIP_PARAMETER_FRIENDLY_NAME":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_FRIENDLY_NAME
	case "PID_KNXIP_PARAMETER_BACKBONE_KEY":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_BACKBONE_KEY
	case "PID_KNXIP_PARAMETER_DEVICE_AUTHENTICATION_CODE":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_DEVICE_AUTHENTICATION_CODE
	case "PID_KNXIP_PARAMETER_PASSWORD_HASHES":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_PASSWORD_HASHES
	case "PID_KNXIP_PARAMETER_SECURED_SERVICE_FAMILIES":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SECURED_SERVICE_FAMILIES
	case "PID_KNXIP_PARAMETER_MULTICAST_LATENCY_TOLERANCE":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MULTICAST_LATENCY_TOLERANCE
	case "PID_KNXIP_PARAMETER_SYNC_LATENCY_FRACTION":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SYNC_LATENCY_FRACTION
	case "PID_KNXIP_PARAMETER_TUNNELLING_USERS":
		return KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_TUNNELLING_USERS
	case "PID_SECURITY_SECURITY_MODE":
		return KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_MODE
	case "PID_SECURITY_P2P_KEY_TABLE":
		return KnxInterfaceObjectProperty_PID_SECURITY_P2P_KEY_TABLE
	case "PID_SECURITY_GRP_KEY_TABLE":
		return KnxInterfaceObjectProperty_PID_SECURITY_GRP_KEY_TABLE
	case "PID_SECURITY_SECURITY_INDIVIDUAL_ADDRESS_TABLE":
		return KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_INDIVIDUAL_ADDRESS_TABLE
	case "PID_SECURITY_SECURITY_FAILURES_LOG":
		return KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_FAILURES_LOG
	case "PID_SECURITY_SKI_TOOL":
		return KnxInterfaceObjectProperty_PID_SECURITY_SKI_TOOL
	case "PID_SECURITY_SECURITY_REPORT":
		return KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_REPORT
	case "PID_SECURITY_SECURITY_REPORT_CONTROL":
		return KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_REPORT_CONTROL
	case "PID_SECURITY_SEQUENCE_NUMBER_SENDING":
		return KnxInterfaceObjectProperty_PID_SECURITY_SEQUENCE_NUMBER_SENDING
	case "PID_SECURITY_ZONE_KEYS_TABLE":
		return KnxInterfaceObjectProperty_PID_SECURITY_ZONE_KEYS_TABLE
	case "PID_SECURITY_GO_SECURITY_FLAGS":
		return KnxInterfaceObjectProperty_PID_SECURITY_GO_SECURITY_FLAGS
	case "PID_RF_MEDIUM_RF_MULTI_TYPE":
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_MULTI_TYPE
	case "PID_RF_MEDIUM_RF_DOMAIN_ADDRESS":
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DOMAIN_ADDRESS
	case "PID_RF_MEDIUM_RF_RETRANSMITTER":
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_RETRANSMITTER
	case "PID_RF_MEDIUM_SECURITY_REPORT_CONTROL":
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_SECURITY_REPORT_CONTROL
	case "PID_RF_MEDIUM_RF_FILTERING_MODE_SELECT":
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_FILTERING_MODE_SELECT
	case "PID_RF_MEDIUM_RF_BIDIR_TIMEOUT":
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_BIDIR_TIMEOUT
	case "PID_RF_MEDIUM_RF_DIAG_SA_FILTER_TABLE":
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DIAG_SA_FILTER_TABLE
	case "PID_RF_MEDIUM_RF_DIAG_QUALITY_TABLE":
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DIAG_QUALITY_TABLE
	case "PID_RF_MEDIUM_RF_DIAG_PROBE":
		return KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DIAG_PROBE
	case "PID_INDOOR_BRIGHTNESS_SENSOR_CHANGE_OF_VALUE":
		return KnxInterfaceObjectProperty_PID_INDOOR_BRIGHTNESS_SENSOR_CHANGE_OF_VALUE
	case "PID_INDOOR_BRIGHTNESS_SENSOR_REPETITION_TIME":
		return KnxInterfaceObjectProperty_PID_INDOOR_BRIGHTNESS_SENSOR_REPETITION_TIME
	case "PID_INDOOR_LUMINANCE_SENSOR_CHANGE_OF_VALUE":
		return KnxInterfaceObjectProperty_PID_INDOOR_LUMINANCE_SENSOR_CHANGE_OF_VALUE
	case "PID_INDOOR_LUMINANCE_SENSOR_REPETITION_TIME":
		return KnxInterfaceObjectProperty_PID_INDOOR_LUMINANCE_SENSOR_REPETITION_TIME
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_ON_DELAY":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_ON_DELAY
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_OFF_DELAY":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_OFF_DELAY
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TIMED_ON_DURATION":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TIMED_ON_DURATION
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_PREWARNING_DURATION":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_PREWARNING_DURATION
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TRANSMISSION_CYCLE_TIME":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TRANSMISSION_CYCLE_TIME
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BUS_POWER_UP_MESSAGE_DELAY":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BUS_POWER_UP_MESSAGE_DELAY
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_AT_LOCKING":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_AT_LOCKING
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_AT_UNLOCKING":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_AT_UNLOCKING
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_DOWN":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_DOWN
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_INVERT_OUTPUT_STATE":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_INVERT_OUTPUT_STATE
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TIMED_ON_RETRIGGER_FUNCTION":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TIMED_ON_RETRIGGER_FUNCTION
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_MANUAL_OFF_ENABLE":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_MANUAL_OFF_ENABLE
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_INVERT_LOCK_DEVICE":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_INVERT_LOCK_DEVICE
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_LOCK_STATE":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_LOCK_STATE
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_UNLOCK_STATE":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_UNLOCK_STATE
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_STATE_FOR_SCENE_NUMBER":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_STATE_FOR_SCENE_NUMBER
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_STORAGE_FUNCTION_FOR_SCENE":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_STORAGE_FUNCTION_FOR_SCENE
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BUS_POWER_UP_STATE":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BUS_POWER_UP_STATE
	case "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP_2":
		return KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP_2
	case "PID_DIMMING_ACTUATOR_BASIC_ON_DELAY":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_ON_DELAY
	case "PID_DIMMING_ACTUATOR_BASIC_OFF_DELAY":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_OFF_DELAY
	case "PID_DIMMING_ACTUATOR_BASIC_SWITCH_OFF_BRIGHTNESS_DELAY_TIME":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_SWITCH_OFF_BRIGHTNESS_DELAY_TIME
	case "PID_DIMMING_ACTUATOR_BASIC_TIMED_ON_DURATION":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_TIMED_ON_DURATION
	case "PID_DIMMING_ACTUATOR_BASIC_PREWARNING_DURATION":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_PREWARNING_DURATION
	case "PID_DIMMING_ACTUATOR_BASIC_TRANSMISSION_CYCLE_TIME":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_TRANSMISSION_CYCLE_TIME
	case "PID_DIMMING_ACTUATOR_BASIC_BUS_POWER_UP_MESSAGE_DELAY":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BUS_POWER_UP_MESSAGE_DELAY
	case "PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED
	case "PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME
	case "PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED_FOR_SWITCH_ON_SET_VALUE":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED_FOR_SWITCH_ON_SET_VALUE
	case "PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED_FOR_SWITCH_OFF":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED_FOR_SWITCH_OFF
	case "PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME_FOR_SWITCH_ON_SET_VALUE":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME_FOR_SWITCH_ON_SET_VALUE
	case "PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME_FOR_SWITCH_OFF":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME_FOR_SWITCH_OFF
	case "PID_DIMMING_ACTUATOR_BASIC_SWITCFH_OFF_BRIGHTNESS":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_SWITCFH_OFF_BRIGHTNESS
	case "PID_DIMMING_ACTUATOR_BASIC_MINIMUM_SET_VALUE":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MINIMUM_SET_VALUE
	case "PID_DIMMING_ACTUATOR_BASIC_MAXIMUM_SET_VALUE":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MAXIMUM_SET_VALUE
	case "PID_DIMMING_ACTUATOR_BASIC_SWITCH_ON_SET_VALUE":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_SWITCH_ON_SET_VALUE
	case "PID_DIMMING_ACTUATOR_BASIC_DIMM_MODE_SELECTION":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMM_MODE_SELECTION
	case "PID_DIMMING_ACTUATOR_BASIC_RELATIV_OFF_ENABLE":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_RELATIV_OFF_ENABLE
	case "PID_DIMMING_ACTUATOR_BASIC_MEMORY_FUNCTION":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MEMORY_FUNCTION
	case "PID_DIMMING_ACTUATOR_BASIC_TIMED_ON_RETRIGGER_FUNCTION":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_TIMED_ON_RETRIGGER_FUNCTION
	case "PID_DIMMING_ACTUATOR_BASIC_MANUAL_OFF_ENABLE":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MANUAL_OFF_ENABLE
	case "PID_DIMMING_ACTUATOR_BASIC_INVERT_LOCK_DEVICE":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_INVERT_LOCK_DEVICE
	case "PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_AT_LOCKING":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_AT_LOCKING
	case "PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_AT_UNLOCKING":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_AT_UNLOCKING
	case "PID_DIMMING_ACTUATOR_BASIC_LOCK_SETVALUE":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_LOCK_SETVALUE
	case "PID_DIMMING_ACTUATOR_BASIC_UNLOCK_SETVALUE":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_UNLOCK_SETVALUE
	case "PID_DIMMING_ACTUATOR_BASIC_BIGHTNESS_FOR_SCENE":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BIGHTNESS_FOR_SCENE
	case "PID_DIMMING_ACTUATOR_BASIC_STORAGE_FUNCTION_FOR_SCENE":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_STORAGE_FUNCTION_FOR_SCENE
	case "PID_DIMMING_ACTUATOR_BASIC_DELTA_DIMMING_VALUE":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DELTA_DIMMING_VALUE
	case "PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP
	case "PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP_SET_VALUE":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP_SET_VALUE
	case "PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_DOWN":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_DOWN
	case "PID_DIMMING_ACTUATOR_BASIC_BUS_POWER_DOWN_SET_VALUE":
		return KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BUS_POWER_DOWN_SET_VALUE
	case "PID_DIMMING_SENSOR_BASIC_ON_OFF_ACTION":
		return KnxInterfaceObjectProperty_PID_DIMMING_SENSOR_BASIC_ON_OFF_ACTION
	case "PID_DIMMING_SENSOR_BASIC_ENABLE_TOGGLE_MODE":
		return KnxInterfaceObjectProperty_PID_DIMMING_SENSOR_BASIC_ENABLE_TOGGLE_MODE
	case "PID_DIMMING_SENSOR_BASIC_ABSOLUTE_SETVALUE":
		return KnxInterfaceObjectProperty_PID_DIMMING_SENSOR_BASIC_ABSOLUTE_SETVALUE
	case "PID_SWITCHING_SENSOR_BASIC_ON_OFF_ACTION":
		return KnxInterfaceObjectProperty_PID_SWITCHING_SENSOR_BASIC_ON_OFF_ACTION
	case "PID_SWITCHING_SENSOR_BASIC_ENABLE_TOGGLE_MODE":
		return KnxInterfaceObjectProperty_PID_SWITCHING_SENSOR_BASIC_ENABLE_TOGGLE_MODE
	case "PID_DEVICE_ROUTING_COUNT":
		return KnxInterfaceObjectProperty_PID_DEVICE_ROUTING_COUNT
	case "PID_DEVICE_MAX_RETRY_COUNT":
		return KnxInterfaceObjectProperty_PID_DEVICE_MAX_RETRY_COUNT
	case "PID_DEVICE_ERROR_FLAGS":
		return KnxInterfaceObjectProperty_PID_DEVICE_ERROR_FLAGS
	case "PID_DEVICE_PROGMODE":
		return KnxInterfaceObjectProperty_PID_DEVICE_PROGMODE
	case "PID_DEVICE_PRODUCT_ID":
		return KnxInterfaceObjectProperty_PID_DEVICE_PRODUCT_ID
	case "PID_DEVICE_MAX_APDULENGTH":
		return KnxInterfaceObjectProperty_PID_DEVICE_MAX_APDULENGTH
	case "PID_DEVICE_SUBNET_ADDR":
		return KnxInterfaceObjectProperty_PID_DEVICE_SUBNET_ADDR
	case "PID_DEVICE_DEVICE_ADDR":
		return KnxInterfaceObjectProperty_PID_DEVICE_DEVICE_ADDR
	case "PID_DEVICE_PB_CONFIG":
		return KnxInterfaceObjectProperty_PID_DEVICE_PB_CONFIG
	case "PID_DEVICE_ADDR_REPORT":
		return KnxInterfaceObjectProperty_PID_DEVICE_ADDR_REPORT
	case "PID_ROUTER_LINE_STATUS":
		return KnxInterfaceObjectProperty_PID_ROUTER_LINE_STATUS
	case "PID_ROUTER_MAIN_LCCONFIG":
		return KnxInterfaceObjectProperty_PID_ROUTER_MAIN_LCCONFIG
	case "PID_ROUTER_SUB_LCCONFIG":
		return KnxInterfaceObjectProperty_PID_ROUTER_SUB_LCCONFIG
	case "PID_ROUTER_MAIN_LCGRPCONFIG":
		return KnxInterfaceObjectProperty_PID_ROUTER_MAIN_LCGRPCONFIG
	case "PID_ROUTER_SUB_LCGRPCONFIG":
		return KnxInterfaceObjectProperty_PID_ROUTER_SUB_LCGRPCONFIG
	case "PID_ROUTER_ROUTETABLE_CONTROL":
		return KnxInterfaceObjectProperty_PID_ROUTER_ROUTETABLE_CONTROL
	case "PID_ROUTER_COUPL_SERV_CONTROL":
		return KnxInterfaceObjectProperty_PID_ROUTER_COUPL_SERV_CONTROL
	case "PID_ROUTER_MAX_ROUTER_APDU_LENGTH":
		return KnxInterfaceObjectProperty_PID_ROUTER_MAX_ROUTER_APDU_LENGTH
	case "PID_ROUTER_MEDIUM":
		return KnxInterfaceObjectProperty_PID_ROUTER_MEDIUM
	case "PID_ROUTER_FILTER_TABLE_USE":
		return KnxInterfaceObjectProperty_PID_ROUTER_FILTER_TABLE_USE
	case "PID_DEVICE_ADDR_CHECK":
		return KnxInterfaceObjectProperty_PID_DEVICE_ADDR_CHECK
	case "PID_ROUTER_RF_ENABLE_SBC":
		return KnxInterfaceObjectProperty_PID_ROUTER_RF_ENABLE_SBC
	case "PID_DEVICE_OBJECT_VALUE":
		return KnxInterfaceObjectProperty_PID_DEVICE_OBJECT_VALUE
	case "PID_DEVICE_OBJECTLINK":
		return KnxInterfaceObjectProperty_PID_DEVICE_OBJECTLINK
	case "PID_DEVICE_APPLICATION":
		return KnxInterfaceObjectProperty_PID_DEVICE_APPLICATION
	case "PID_DEVICE_PARAMETER":
		return KnxInterfaceObjectProperty_PID_DEVICE_PARAMETER
	case "PID_GENERAL_OBJECT_TYPE":
		return KnxInterfaceObjectProperty_PID_GENERAL_OBJECT_TYPE
	case "PID_GENERAL_OBJECT_NAME":
		return KnxInterfaceObjectProperty_PID_GENERAL_OBJECT_NAME
	case "PID_GENERAL_SEMAPHOR":
		return KnxInterfaceObjectProperty_PID_GENERAL_SEMAPHOR
	case "PID_GENERAL_GROUP_OBJECT_REFERENCE":
		return KnxInterfaceObjectProperty_PID_GENERAL_GROUP_OBJECT_REFERENCE
	case "PID_GENERAL_LOAD_STATE_CONTROL":
		return KnxInterfaceObjectProperty_PID_GENERAL_LOAD_STATE_CONTROL
	case "PID_GENERAL_RUN_STATE_CONTROL":
		return KnxInterfaceObjectProperty_PID_GENERAL_RUN_STATE_CONTROL
	case "PID_GENERAL_TABLE_REFERENCE":
		return KnxInterfaceObjectProperty_PID_GENERAL_TABLE_REFERENCE
	case "PID_GENERAL_SERVICE_CONTROL":
		return KnxInterfaceObjectProperty_PID_GENERAL_SERVICE_CONTROL
	case "PID_GENERAL_FIRMWARE_REVISION":
		return KnxInterfaceObjectProperty_PID_GENERAL_FIRMWARE_REVISION
	case "PID_GENERAL_SERVICES_SUPPORTED":
		return KnxInterfaceObjectProperty_PID_GENERAL_SERVICES_SUPPORTED
	case "PID_GENERAL_SERIAL_NUMBER":
		return KnxInterfaceObjectProperty_PID_GENERAL_SERIAL_NUMBER
	case "PID_GENERAL_MANUFACTURER_ID":
		return KnxInterfaceObjectProperty_PID_GENERAL_MANUFACTURER_ID
	case "PID_GENERAL_PROGRAM_VERSION":
		return KnxInterfaceObjectProperty_PID_GENERAL_PROGRAM_VERSION
	case "PID_GENERAL_DEVICE_CONTROL":
		return KnxInterfaceObjectProperty_PID_GENERAL_DEVICE_CONTROL
	case "PID_GENERAL_ORDER_INFO":
		return KnxInterfaceObjectProperty_PID_GENERAL_ORDER_INFO
	case "PID_GENERAL_PEI_TYPE":
		return KnxInterfaceObjectProperty_PID_GENERAL_PEI_TYPE
	case "PID_GENERAL_PORT_CONFIGURATION":
		return KnxInterfaceObjectProperty_PID_GENERAL_PORT_CONFIGURATION
	case "PID_GENERAL_POLL_GROUP_SETTINGS":
		return KnxInterfaceObjectProperty_PID_GENERAL_POLL_GROUP_SETTINGS
	case "PID_GENERAL_MANUFACTURER_DATA":
		return KnxInterfaceObjectProperty_PID_GENERAL_MANUFACTURER_DATA
	case "PID_GENERAL_ENABLE":
		return KnxInterfaceObjectProperty_PID_GENERAL_ENABLE
	case "PID_GENERAL_DESCRIPTION":
		return KnxInterfaceObjectProperty_PID_GENERAL_DESCRIPTION
	case "PID_GENERAL_FILE":
		return KnxInterfaceObjectProperty_PID_GENERAL_FILE
	case "PID_GENERAL_TABLE":
		return KnxInterfaceObjectProperty_PID_GENERAL_TABLE
	case "PID_GENERAL_ENROL":
		return KnxInterfaceObjectProperty_PID_GENERAL_ENROL
	case "PID_GENERAL_VERSION":
		return KnxInterfaceObjectProperty_PID_GENERAL_VERSION
	case "PID_GENERAL_GROUP_OBJECT_LINK":
		return KnxInterfaceObjectProperty_PID_GENERAL_GROUP_OBJECT_LINK
	case "PID_GENERAL_MCB_TABLE":
		return KnxInterfaceObjectProperty_PID_GENERAL_MCB_TABLE
	case "PID_GENERAL_ERROR_CODE":
		return KnxInterfaceObjectProperty_PID_GENERAL_ERROR_CODE
	case "PID_GENERAL_OBJECT_INDEX":
		return KnxInterfaceObjectProperty_PID_GENERAL_OBJECT_INDEX
	case "PID_GENERAL_DOWNLOAD_COUNTER":
		return KnxInterfaceObjectProperty_PID_GENERAL_DOWNLOAD_COUNTER
	case "PID_DEVICE_OBJECTADDRESS":
		return KnxInterfaceObjectProperty_PID_DEVICE_OBJECTADDRESS
	case "PID_DEVICE_PSU_TYPE":
		return KnxInterfaceObjectProperty_PID_DEVICE_PSU_TYPE
	case "PID_DEVICE_PSU_STATUS":
		return KnxInterfaceObjectProperty_PID_DEVICE_PSU_STATUS
	case "PID_DEVICE_PSU_ENABLE":
		return KnxInterfaceObjectProperty_PID_DEVICE_PSU_ENABLE
	case "PID_DEVICE_DOMAIN_ADDRESS":
		return KnxInterfaceObjectProperty_PID_DEVICE_DOMAIN_ADDRESS
	case "PID_DEVICE_IO_LIST":
		return KnxInterfaceObjectProperty_PID_DEVICE_IO_LIST
	case "PID_DEVICE_MGT_DESCRIPTOR_01":
		return KnxInterfaceObjectProperty_PID_DEVICE_MGT_DESCRIPTOR_01
	case "PID_DEVICE_PL110_PARAM":
		return KnxInterfaceObjectProperty_PID_DEVICE_PL110_PARAM
	case "PID_DEVICE_RF_REPEAT_COUNTER":
		return KnxInterfaceObjectProperty_PID_DEVICE_RF_REPEAT_COUNTER
	case "PID_DEVICE_RECEIVE_BLOCK_TABLE":
		return KnxInterfaceObjectProperty_PID_DEVICE_RECEIVE_BLOCK_TABLE
	case "PID_DEVICE_RANDOM_PAUSE_TABLE":
		return KnxInterfaceObjectProperty_PID_DEVICE_RANDOM_PAUSE_TABLE
	case "PID_DEVICE_RECEIVE_BLOCK_NR":
		return KnxInterfaceObjectProperty_PID_DEVICE_RECEIVE_BLOCK_NR
	case "PID_DEVICE_HARDWARE_TYPE":
		return KnxInterfaceObjectProperty_PID_DEVICE_HARDWARE_TYPE
	case "PID_DEVICE_RETRANSMITTER_NUMBER":
		return KnxInterfaceObjectProperty_PID_DEVICE_RETRANSMITTER_NUMBER
	case "PID_DEVICE_SERIAL_NR_TABLE":
		return KnxInterfaceObjectProperty_PID_DEVICE_SERIAL_NR_TABLE
	case "PID_SUNBLIND_ACTUATOR_BASIC_REVERSION_PAUSE_TIME":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REVERSION_PAUSE_TIME
	case "PID_SUNBLIND_ACTUATOR_BASIC_MOVE_UP_DOWN_TIME":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_UP_DOWN_TIME
	case "PID_SUNBLIND_ACTUATOR_BASIC_SLAT_STEP_TIME":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_SLAT_STEP_TIME
	case "PID_SUNBLIND_ACTUATOR_BASIC_MOVE_PRESET_POSITION_TIME":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_PRESET_POSITION_TIME
	case "PID_SUNBLIND_ACTUATOR_BASIC_MOVE_TO_PRESET_POSITION_IN_PERCENT":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_TO_PRESET_POSITION_IN_PERCENT
	case "PID_SUNBLIND_ACTUATOR_BASIC_MOVE_TO_PRESET_POSITION_LENGTH":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_TO_PRESET_POSITION_LENGTH
	case "PID_SUNBLIND_ACTUATOR_BASIC_PRESET_SLAT_POSITION_PERCENT":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_PRESET_SLAT_POSITION_PERCENT
	case "PID_SUNBLIND_ACTUATOR_BASIC_PRESET_SLAT_POSITION_ANGLE":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_PRESET_SLAT_POSITION_ANGLE
	case "PID_SUNBLIND_ACTUATOR_BASIC_REACTION_WIND_ALARM":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REACTION_WIND_ALARM
	case "PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_WIND_ALARM":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_WIND_ALARM
	case "PID_SUNBLIND_ACTUATOR_BASIC_REACTION_ON_RAIN_ALARM":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REACTION_ON_RAIN_ALARM
	case "PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_RAIN_ALARM":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_RAIN_ALARM
	case "PID_SUNBLIND_ACTUATOR_BASIC_REACTION_FROST_ALARM":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REACTION_FROST_ALARM
	case "PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_FROST_ALARM":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_FROST_ALARM
	case "PID_SUNBLIND_ACTUATOR_BASIC_MAX_SLAT_MOVE_TIME":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MAX_SLAT_MOVE_TIME
	case "PID_SUNBLIND_ACTUATOR_BASIC_ENABLE_BLINDS_MODE":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_ENABLE_BLINDS_MODE
	case "PID_SUNBLIND_ACTUATOR_BASIC_STORAGE_FUNCTIONS_FOR_SCENE":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_STORAGE_FUNCTIONS_FOR_SCENE
	case "PID_SUNBLIND_SENSOR_BASIC_ENABLE_BLINDS_MODE":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_SENSOR_BASIC_ENABLE_BLINDS_MODE
	case "PID_SUNBLIND_SENSOR_BASIC_UP_DOWN_ACTION":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_SENSOR_BASIC_UP_DOWN_ACTION
	case "PID_SUNBLIND_SENSOR_BASIC_ENABLE_TOGGLE_MODE":
		return KnxInterfaceObjectProperty_PID_SUNBLIND_SENSOR_BASIC_ENABLE_TOGGLE_MODE
	case "PID_DEVICE_BIBATMASTER_ADDRESS":
		return KnxInterfaceObjectProperty_PID_DEVICE_BIBATMASTER_ADDRESS
	case "PID_DEVICE_RF_DOMAIN_ADDRESS":
		return KnxInterfaceObjectProperty_PID_DEVICE_RF_DOMAIN_ADDRESS
	case "PID_DEVICE_DEVICE_DESCRIPTOR":
		return KnxInterfaceObjectProperty_PID_DEVICE_DEVICE_DESCRIPTOR
	case "PID_DEVICE_METERING_FILTER_TABLE":
		return KnxInterfaceObjectProperty_PID_DEVICE_METERING_FILTER_TABLE
	case "PID_DEVICE_GROUP_TELEGR_RATE_LIMIT_TIME_BASE":
		return KnxInterfaceObjectProperty_PID_DEVICE_GROUP_TELEGR_RATE_LIMIT_TIME_BASE
	case "PID_DEVICE_GROUP_TELEGR_RATE_LIMIT_NO_OF_TELEGR":
		return KnxInterfaceObjectProperty_PID_DEVICE_GROUP_TELEGR_RATE_LIMIT_NO_OF_TELEGR
	case "PID_GROUP_OBJECT_TABLE_GRPOBJTABLE":
		return KnxInterfaceObjectProperty_PID_GROUP_OBJECT_TABLE_GRPOBJTABLE
	case "PID_GROUP_OBJECT_TABLE_EXT_GRPOBJREFERENCE":
		return KnxInterfaceObjectProperty_PID_GROUP_OBJECT_TABLE_EXT_GRPOBJREFERENCE
	}
	return 0
}

func CastKnxInterfaceObjectProperty(structType interface{}) KnxInterfaceObjectProperty {
	castFunc := func(typ interface{}) KnxInterfaceObjectProperty {
		if sKnxInterfaceObjectProperty, ok := typ.(KnxInterfaceObjectProperty); ok {
			return sKnxInterfaceObjectProperty
		}
		return 0
	}
	return castFunc(structType)
}

func (m KnxInterfaceObjectProperty) LengthInBits() uint16 {
	return 32
}

func (m KnxInterfaceObjectProperty) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func KnxInterfaceObjectPropertyParse(io *utils.ReadBuffer) (KnxInterfaceObjectProperty, error) {
	val, err := io.ReadUint32(32)
	if err != nil {
		return 0, nil
	}
	return KnxInterfaceObjectPropertyByValue(val), nil
}

func (e KnxInterfaceObjectProperty) Serialize(io utils.WriteBuffer) error {
	err := io.WriteUint32(32, uint32(e))
	return err
}

func (e KnxInterfaceObjectProperty) String() string {
	switch e {
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_PROJECT_INSTALLATION_ID:
		return "PID_KNXIP_PARAMETER_PROJECT_INSTALLATION_ID"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNX_INDIVIDUAL_ADDRESS:
		return "PID_KNXIP_PARAMETER_KNX_INDIVIDUAL_ADDRESS"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_ADDITIONAL_INDIVIDUAL_ADDRESSES:
		return "PID_KNXIP_PARAMETER_ADDITIONAL_INDIVIDUAL_ADDRESSES"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_IP_ASSIGNMENT_METHOD:
		return "PID_KNXIP_PARAMETER_CURRENT_IP_ASSIGNMENT_METHOD"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_IP_ASSIGNMENT_METHOD:
		return "PID_KNXIP_PARAMETER_IP_ASSIGNMENT_METHOD"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_IP_CAPABILITIES:
		return "PID_KNXIP_PARAMETER_IP_CAPABILITIES"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_IP_ADDRESS:
		return "PID_KNXIP_PARAMETER_CURRENT_IP_ADDRESS"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_SUBNET_MASK:
		return "PID_KNXIP_PARAMETER_CURRENT_SUBNET_MASK"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_CURRENT_DEFAULT_GATEWAY:
		return "PID_KNXIP_PARAMETER_CURRENT_DEFAULT_GATEWAY"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_IP_ADDRESS:
		return "PID_KNXIP_PARAMETER_IP_ADDRESS"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SUBNET_MASK:
		return "PID_KNXIP_PARAMETER_SUBNET_MASK"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_DEFAULT_GATEWAY:
		return "PID_KNXIP_PARAMETER_DEFAULT_GATEWAY"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_DHCP_BOOTP_SERVER:
		return "PID_KNXIP_PARAMETER_DHCP_BOOTP_SERVER"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MAC_ADDRESS:
		return "PID_KNXIP_PARAMETER_MAC_ADDRESS"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SYSTEM_SETUP_MULTICAST_ADDRESS:
		return "PID_KNXIP_PARAMETER_SYSTEM_SETUP_MULTICAST_ADDRESS"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_ROUTING_MULTICAST_ADDRESS:
		return "PID_KNXIP_PARAMETER_ROUTING_MULTICAST_ADDRESS"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_TTL:
		return "PID_KNXIP_PARAMETER_TTL"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNXNETIP_DEVICE_CAPABILITIES:
		return "PID_KNXIP_PARAMETER_KNXNETIP_DEVICE_CAPABILITIES"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNXNETIP_DEVICE_STATE:
		return "PID_KNXIP_PARAMETER_KNXNETIP_DEVICE_STATE"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_KNXNETIP_ROUTING_CAPABILITIES:
		return "PID_KNXIP_PARAMETER_KNXNETIP_ROUTING_CAPABILITIES"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_PRIORITY_FIFO_ENABLED:
		return "PID_KNXIP_PARAMETER_PRIORITY_FIFO_ENABLED"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_QUEUE_OVERFLOW_TO_IP:
		return "PID_KNXIP_PARAMETER_QUEUE_OVERFLOW_TO_IP"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_QUEUE_OVERFLOW_TO_KNX:
		return "PID_KNXIP_PARAMETER_QUEUE_OVERFLOW_TO_KNX"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MSG_TRANSMIT_TO_IP:
		return "PID_KNXIP_PARAMETER_MSG_TRANSMIT_TO_IP"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MSG_TRANSMIT_TO_KNX:
		return "PID_KNXIP_PARAMETER_MSG_TRANSMIT_TO_KNX"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_FRIENDLY_NAME:
		return "PID_KNXIP_PARAMETER_FRIENDLY_NAME"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_BACKBONE_KEY:
		return "PID_KNXIP_PARAMETER_BACKBONE_KEY"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_DEVICE_AUTHENTICATION_CODE:
		return "PID_KNXIP_PARAMETER_DEVICE_AUTHENTICATION_CODE"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_PASSWORD_HASHES:
		return "PID_KNXIP_PARAMETER_PASSWORD_HASHES"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SECURED_SERVICE_FAMILIES:
		return "PID_KNXIP_PARAMETER_SECURED_SERVICE_FAMILIES"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_MULTICAST_LATENCY_TOLERANCE:
		return "PID_KNXIP_PARAMETER_MULTICAST_LATENCY_TOLERANCE"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_SYNC_LATENCY_FRACTION:
		return "PID_KNXIP_PARAMETER_SYNC_LATENCY_FRACTION"
	case KnxInterfaceObjectProperty_PID_KNXIP_PARAMETER_TUNNELLING_USERS:
		return "PID_KNXIP_PARAMETER_TUNNELLING_USERS"
	case KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_MODE:
		return "PID_SECURITY_SECURITY_MODE"
	case KnxInterfaceObjectProperty_PID_SECURITY_P2P_KEY_TABLE:
		return "PID_SECURITY_P2P_KEY_TABLE"
	case KnxInterfaceObjectProperty_PID_SECURITY_GRP_KEY_TABLE:
		return "PID_SECURITY_GRP_KEY_TABLE"
	case KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_INDIVIDUAL_ADDRESS_TABLE:
		return "PID_SECURITY_SECURITY_INDIVIDUAL_ADDRESS_TABLE"
	case KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_FAILURES_LOG:
		return "PID_SECURITY_SECURITY_FAILURES_LOG"
	case KnxInterfaceObjectProperty_PID_SECURITY_SKI_TOOL:
		return "PID_SECURITY_SKI_TOOL"
	case KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_REPORT:
		return "PID_SECURITY_SECURITY_REPORT"
	case KnxInterfaceObjectProperty_PID_SECURITY_SECURITY_REPORT_CONTROL:
		return "PID_SECURITY_SECURITY_REPORT_CONTROL"
	case KnxInterfaceObjectProperty_PID_SECURITY_SEQUENCE_NUMBER_SENDING:
		return "PID_SECURITY_SEQUENCE_NUMBER_SENDING"
	case KnxInterfaceObjectProperty_PID_SECURITY_ZONE_KEYS_TABLE:
		return "PID_SECURITY_ZONE_KEYS_TABLE"
	case KnxInterfaceObjectProperty_PID_SECURITY_GO_SECURITY_FLAGS:
		return "PID_SECURITY_GO_SECURITY_FLAGS"
	case KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_MULTI_TYPE:
		return "PID_RF_MEDIUM_RF_MULTI_TYPE"
	case KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DOMAIN_ADDRESS:
		return "PID_RF_MEDIUM_RF_DOMAIN_ADDRESS"
	case KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_RETRANSMITTER:
		return "PID_RF_MEDIUM_RF_RETRANSMITTER"
	case KnxInterfaceObjectProperty_PID_RF_MEDIUM_SECURITY_REPORT_CONTROL:
		return "PID_RF_MEDIUM_SECURITY_REPORT_CONTROL"
	case KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_FILTERING_MODE_SELECT:
		return "PID_RF_MEDIUM_RF_FILTERING_MODE_SELECT"
	case KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_BIDIR_TIMEOUT:
		return "PID_RF_MEDIUM_RF_BIDIR_TIMEOUT"
	case KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DIAG_SA_FILTER_TABLE:
		return "PID_RF_MEDIUM_RF_DIAG_SA_FILTER_TABLE"
	case KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DIAG_QUALITY_TABLE:
		return "PID_RF_MEDIUM_RF_DIAG_QUALITY_TABLE"
	case KnxInterfaceObjectProperty_PID_RF_MEDIUM_RF_DIAG_PROBE:
		return "PID_RF_MEDIUM_RF_DIAG_PROBE"
	case KnxInterfaceObjectProperty_PID_INDOOR_BRIGHTNESS_SENSOR_CHANGE_OF_VALUE:
		return "PID_INDOOR_BRIGHTNESS_SENSOR_CHANGE_OF_VALUE"
	case KnxInterfaceObjectProperty_PID_INDOOR_BRIGHTNESS_SENSOR_REPETITION_TIME:
		return "PID_INDOOR_BRIGHTNESS_SENSOR_REPETITION_TIME"
	case KnxInterfaceObjectProperty_PID_INDOOR_LUMINANCE_SENSOR_CHANGE_OF_VALUE:
		return "PID_INDOOR_LUMINANCE_SENSOR_CHANGE_OF_VALUE"
	case KnxInterfaceObjectProperty_PID_INDOOR_LUMINANCE_SENSOR_REPETITION_TIME:
		return "PID_INDOOR_LUMINANCE_SENSOR_REPETITION_TIME"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_ON_DELAY:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_ON_DELAY"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_OFF_DELAY:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_OFF_DELAY"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TIMED_ON_DURATION:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TIMED_ON_DURATION"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_PREWARNING_DURATION:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_PREWARNING_DURATION"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TRANSMISSION_CYCLE_TIME:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TRANSMISSION_CYCLE_TIME"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BUS_POWER_UP_MESSAGE_DELAY:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BUS_POWER_UP_MESSAGE_DELAY"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_AT_LOCKING:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_AT_LOCKING"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_AT_UNLOCKING:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_AT_UNLOCKING"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_DOWN:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_DOWN"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_INVERT_OUTPUT_STATE:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_INVERT_OUTPUT_STATE"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TIMED_ON_RETRIGGER_FUNCTION:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_TIMED_ON_RETRIGGER_FUNCTION"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_MANUAL_OFF_ENABLE:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_MANUAL_OFF_ENABLE"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_INVERT_LOCK_DEVICE:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_INVERT_LOCK_DEVICE"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_LOCK_STATE:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_LOCK_STATE"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_UNLOCK_STATE:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_UNLOCK_STATE"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_STATE_FOR_SCENE_NUMBER:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_STATE_FOR_SCENE_NUMBER"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_STORAGE_FUNCTION_FOR_SCENE:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_STORAGE_FUNCTION_FOR_SCENE"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BUS_POWER_UP_STATE:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BUS_POWER_UP_STATE"
	case KnxInterfaceObjectProperty_PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP_2:
		return "PID_LIGHT_SWITCHING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP_2"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_ON_DELAY:
		return "PID_DIMMING_ACTUATOR_BASIC_ON_DELAY"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_OFF_DELAY:
		return "PID_DIMMING_ACTUATOR_BASIC_OFF_DELAY"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_SWITCH_OFF_BRIGHTNESS_DELAY_TIME:
		return "PID_DIMMING_ACTUATOR_BASIC_SWITCH_OFF_BRIGHTNESS_DELAY_TIME"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_TIMED_ON_DURATION:
		return "PID_DIMMING_ACTUATOR_BASIC_TIMED_ON_DURATION"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_PREWARNING_DURATION:
		return "PID_DIMMING_ACTUATOR_BASIC_PREWARNING_DURATION"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_TRANSMISSION_CYCLE_TIME:
		return "PID_DIMMING_ACTUATOR_BASIC_TRANSMISSION_CYCLE_TIME"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BUS_POWER_UP_MESSAGE_DELAY:
		return "PID_DIMMING_ACTUATOR_BASIC_BUS_POWER_UP_MESSAGE_DELAY"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED:
		return "PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME:
		return "PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED_FOR_SWITCH_ON_SET_VALUE:
		return "PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED_FOR_SWITCH_ON_SET_VALUE"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED_FOR_SWITCH_OFF:
		return "PID_DIMMING_ACTUATOR_BASIC_DIMMING_SPEED_FOR_SWITCH_OFF"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME_FOR_SWITCH_ON_SET_VALUE:
		return "PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME_FOR_SWITCH_ON_SET_VALUE"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME_FOR_SWITCH_OFF:
		return "PID_DIMMING_ACTUATOR_BASIC_DIMMING_STEP_TIME_FOR_SWITCH_OFF"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_SWITCFH_OFF_BRIGHTNESS:
		return "PID_DIMMING_ACTUATOR_BASIC_SWITCFH_OFF_BRIGHTNESS"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MINIMUM_SET_VALUE:
		return "PID_DIMMING_ACTUATOR_BASIC_MINIMUM_SET_VALUE"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MAXIMUM_SET_VALUE:
		return "PID_DIMMING_ACTUATOR_BASIC_MAXIMUM_SET_VALUE"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_SWITCH_ON_SET_VALUE:
		return "PID_DIMMING_ACTUATOR_BASIC_SWITCH_ON_SET_VALUE"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DIMM_MODE_SELECTION:
		return "PID_DIMMING_ACTUATOR_BASIC_DIMM_MODE_SELECTION"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_RELATIV_OFF_ENABLE:
		return "PID_DIMMING_ACTUATOR_BASIC_RELATIV_OFF_ENABLE"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MEMORY_FUNCTION:
		return "PID_DIMMING_ACTUATOR_BASIC_MEMORY_FUNCTION"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_TIMED_ON_RETRIGGER_FUNCTION:
		return "PID_DIMMING_ACTUATOR_BASIC_TIMED_ON_RETRIGGER_FUNCTION"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_MANUAL_OFF_ENABLE:
		return "PID_DIMMING_ACTUATOR_BASIC_MANUAL_OFF_ENABLE"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_INVERT_LOCK_DEVICE:
		return "PID_DIMMING_ACTUATOR_BASIC_INVERT_LOCK_DEVICE"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_AT_LOCKING:
		return "PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_AT_LOCKING"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_AT_UNLOCKING:
		return "PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_AT_UNLOCKING"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_LOCK_SETVALUE:
		return "PID_DIMMING_ACTUATOR_BASIC_LOCK_SETVALUE"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_UNLOCK_SETVALUE:
		return "PID_DIMMING_ACTUATOR_BASIC_UNLOCK_SETVALUE"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BIGHTNESS_FOR_SCENE:
		return "PID_DIMMING_ACTUATOR_BASIC_BIGHTNESS_FOR_SCENE"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_STORAGE_FUNCTION_FOR_SCENE:
		return "PID_DIMMING_ACTUATOR_BASIC_STORAGE_FUNCTION_FOR_SCENE"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_DELTA_DIMMING_VALUE:
		return "PID_DIMMING_ACTUATOR_BASIC_DELTA_DIMMING_VALUE"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP:
		return "PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP_SET_VALUE:
		return "PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_UP_SET_VALUE"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_DOWN:
		return "PID_DIMMING_ACTUATOR_BASIC_BEHAVIOUR_BUS_POWER_DOWN"
	case KnxInterfaceObjectProperty_PID_DIMMING_ACTUATOR_BASIC_BUS_POWER_DOWN_SET_VALUE:
		return "PID_DIMMING_ACTUATOR_BASIC_BUS_POWER_DOWN_SET_VALUE"
	case KnxInterfaceObjectProperty_PID_DIMMING_SENSOR_BASIC_ON_OFF_ACTION:
		return "PID_DIMMING_SENSOR_BASIC_ON_OFF_ACTION"
	case KnxInterfaceObjectProperty_PID_DIMMING_SENSOR_BASIC_ENABLE_TOGGLE_MODE:
		return "PID_DIMMING_SENSOR_BASIC_ENABLE_TOGGLE_MODE"
	case KnxInterfaceObjectProperty_PID_DIMMING_SENSOR_BASIC_ABSOLUTE_SETVALUE:
		return "PID_DIMMING_SENSOR_BASIC_ABSOLUTE_SETVALUE"
	case KnxInterfaceObjectProperty_PID_SWITCHING_SENSOR_BASIC_ON_OFF_ACTION:
		return "PID_SWITCHING_SENSOR_BASIC_ON_OFF_ACTION"
	case KnxInterfaceObjectProperty_PID_SWITCHING_SENSOR_BASIC_ENABLE_TOGGLE_MODE:
		return "PID_SWITCHING_SENSOR_BASIC_ENABLE_TOGGLE_MODE"
	case KnxInterfaceObjectProperty_PID_DEVICE_ROUTING_COUNT:
		return "PID_DEVICE_ROUTING_COUNT"
	case KnxInterfaceObjectProperty_PID_DEVICE_MAX_RETRY_COUNT:
		return "PID_DEVICE_MAX_RETRY_COUNT"
	case KnxInterfaceObjectProperty_PID_DEVICE_ERROR_FLAGS:
		return "PID_DEVICE_ERROR_FLAGS"
	case KnxInterfaceObjectProperty_PID_DEVICE_PROGMODE:
		return "PID_DEVICE_PROGMODE"
	case KnxInterfaceObjectProperty_PID_DEVICE_PRODUCT_ID:
		return "PID_DEVICE_PRODUCT_ID"
	case KnxInterfaceObjectProperty_PID_DEVICE_MAX_APDULENGTH:
		return "PID_DEVICE_MAX_APDULENGTH"
	case KnxInterfaceObjectProperty_PID_DEVICE_SUBNET_ADDR:
		return "PID_DEVICE_SUBNET_ADDR"
	case KnxInterfaceObjectProperty_PID_DEVICE_DEVICE_ADDR:
		return "PID_DEVICE_DEVICE_ADDR"
	case KnxInterfaceObjectProperty_PID_DEVICE_PB_CONFIG:
		return "PID_DEVICE_PB_CONFIG"
	case KnxInterfaceObjectProperty_PID_DEVICE_ADDR_REPORT:
		return "PID_DEVICE_ADDR_REPORT"
	case KnxInterfaceObjectProperty_PID_ROUTER_LINE_STATUS:
		return "PID_ROUTER_LINE_STATUS"
	case KnxInterfaceObjectProperty_PID_ROUTER_MAIN_LCCONFIG:
		return "PID_ROUTER_MAIN_LCCONFIG"
	case KnxInterfaceObjectProperty_PID_ROUTER_SUB_LCCONFIG:
		return "PID_ROUTER_SUB_LCCONFIG"
	case KnxInterfaceObjectProperty_PID_ROUTER_MAIN_LCGRPCONFIG:
		return "PID_ROUTER_MAIN_LCGRPCONFIG"
	case KnxInterfaceObjectProperty_PID_ROUTER_SUB_LCGRPCONFIG:
		return "PID_ROUTER_SUB_LCGRPCONFIG"
	case KnxInterfaceObjectProperty_PID_ROUTER_ROUTETABLE_CONTROL:
		return "PID_ROUTER_ROUTETABLE_CONTROL"
	case KnxInterfaceObjectProperty_PID_ROUTER_COUPL_SERV_CONTROL:
		return "PID_ROUTER_COUPL_SERV_CONTROL"
	case KnxInterfaceObjectProperty_PID_ROUTER_MAX_ROUTER_APDU_LENGTH:
		return "PID_ROUTER_MAX_ROUTER_APDU_LENGTH"
	case KnxInterfaceObjectProperty_PID_ROUTER_MEDIUM:
		return "PID_ROUTER_MEDIUM"
	case KnxInterfaceObjectProperty_PID_ROUTER_FILTER_TABLE_USE:
		return "PID_ROUTER_FILTER_TABLE_USE"
	case KnxInterfaceObjectProperty_PID_DEVICE_ADDR_CHECK:
		return "PID_DEVICE_ADDR_CHECK"
	case KnxInterfaceObjectProperty_PID_ROUTER_RF_ENABLE_SBC:
		return "PID_ROUTER_RF_ENABLE_SBC"
	case KnxInterfaceObjectProperty_PID_DEVICE_OBJECT_VALUE:
		return "PID_DEVICE_OBJECT_VALUE"
	case KnxInterfaceObjectProperty_PID_DEVICE_OBJECTLINK:
		return "PID_DEVICE_OBJECTLINK"
	case KnxInterfaceObjectProperty_PID_DEVICE_APPLICATION:
		return "PID_DEVICE_APPLICATION"
	case KnxInterfaceObjectProperty_PID_DEVICE_PARAMETER:
		return "PID_DEVICE_PARAMETER"
	case KnxInterfaceObjectProperty_PID_GENERAL_OBJECT_TYPE:
		return "PID_GENERAL_OBJECT_TYPE"
	case KnxInterfaceObjectProperty_PID_GENERAL_OBJECT_NAME:
		return "PID_GENERAL_OBJECT_NAME"
	case KnxInterfaceObjectProperty_PID_GENERAL_SEMAPHOR:
		return "PID_GENERAL_SEMAPHOR"
	case KnxInterfaceObjectProperty_PID_GENERAL_GROUP_OBJECT_REFERENCE:
		return "PID_GENERAL_GROUP_OBJECT_REFERENCE"
	case KnxInterfaceObjectProperty_PID_GENERAL_LOAD_STATE_CONTROL:
		return "PID_GENERAL_LOAD_STATE_CONTROL"
	case KnxInterfaceObjectProperty_PID_GENERAL_RUN_STATE_CONTROL:
		return "PID_GENERAL_RUN_STATE_CONTROL"
	case KnxInterfaceObjectProperty_PID_GENERAL_TABLE_REFERENCE:
		return "PID_GENERAL_TABLE_REFERENCE"
	case KnxInterfaceObjectProperty_PID_GENERAL_SERVICE_CONTROL:
		return "PID_GENERAL_SERVICE_CONTROL"
	case KnxInterfaceObjectProperty_PID_GENERAL_FIRMWARE_REVISION:
		return "PID_GENERAL_FIRMWARE_REVISION"
	case KnxInterfaceObjectProperty_PID_GENERAL_SERVICES_SUPPORTED:
		return "PID_GENERAL_SERVICES_SUPPORTED"
	case KnxInterfaceObjectProperty_PID_GENERAL_SERIAL_NUMBER:
		return "PID_GENERAL_SERIAL_NUMBER"
	case KnxInterfaceObjectProperty_PID_GENERAL_MANUFACTURER_ID:
		return "PID_GENERAL_MANUFACTURER_ID"
	case KnxInterfaceObjectProperty_PID_GENERAL_PROGRAM_VERSION:
		return "PID_GENERAL_PROGRAM_VERSION"
	case KnxInterfaceObjectProperty_PID_GENERAL_DEVICE_CONTROL:
		return "PID_GENERAL_DEVICE_CONTROL"
	case KnxInterfaceObjectProperty_PID_GENERAL_ORDER_INFO:
		return "PID_GENERAL_ORDER_INFO"
	case KnxInterfaceObjectProperty_PID_GENERAL_PEI_TYPE:
		return "PID_GENERAL_PEI_TYPE"
	case KnxInterfaceObjectProperty_PID_GENERAL_PORT_CONFIGURATION:
		return "PID_GENERAL_PORT_CONFIGURATION"
	case KnxInterfaceObjectProperty_PID_GENERAL_POLL_GROUP_SETTINGS:
		return "PID_GENERAL_POLL_GROUP_SETTINGS"
	case KnxInterfaceObjectProperty_PID_GENERAL_MANUFACTURER_DATA:
		return "PID_GENERAL_MANUFACTURER_DATA"
	case KnxInterfaceObjectProperty_PID_GENERAL_ENABLE:
		return "PID_GENERAL_ENABLE"
	case KnxInterfaceObjectProperty_PID_GENERAL_DESCRIPTION:
		return "PID_GENERAL_DESCRIPTION"
	case KnxInterfaceObjectProperty_PID_GENERAL_FILE:
		return "PID_GENERAL_FILE"
	case KnxInterfaceObjectProperty_PID_GENERAL_TABLE:
		return "PID_GENERAL_TABLE"
	case KnxInterfaceObjectProperty_PID_GENERAL_ENROL:
		return "PID_GENERAL_ENROL"
	case KnxInterfaceObjectProperty_PID_GENERAL_VERSION:
		return "PID_GENERAL_VERSION"
	case KnxInterfaceObjectProperty_PID_GENERAL_GROUP_OBJECT_LINK:
		return "PID_GENERAL_GROUP_OBJECT_LINK"
	case KnxInterfaceObjectProperty_PID_GENERAL_MCB_TABLE:
		return "PID_GENERAL_MCB_TABLE"
	case KnxInterfaceObjectProperty_PID_GENERAL_ERROR_CODE:
		return "PID_GENERAL_ERROR_CODE"
	case KnxInterfaceObjectProperty_PID_GENERAL_OBJECT_INDEX:
		return "PID_GENERAL_OBJECT_INDEX"
	case KnxInterfaceObjectProperty_PID_GENERAL_DOWNLOAD_COUNTER:
		return "PID_GENERAL_DOWNLOAD_COUNTER"
	case KnxInterfaceObjectProperty_PID_DEVICE_OBJECTADDRESS:
		return "PID_DEVICE_OBJECTADDRESS"
	case KnxInterfaceObjectProperty_PID_DEVICE_PSU_TYPE:
		return "PID_DEVICE_PSU_TYPE"
	case KnxInterfaceObjectProperty_PID_DEVICE_PSU_STATUS:
		return "PID_DEVICE_PSU_STATUS"
	case KnxInterfaceObjectProperty_PID_DEVICE_PSU_ENABLE:
		return "PID_DEVICE_PSU_ENABLE"
	case KnxInterfaceObjectProperty_PID_DEVICE_DOMAIN_ADDRESS:
		return "PID_DEVICE_DOMAIN_ADDRESS"
	case KnxInterfaceObjectProperty_PID_DEVICE_IO_LIST:
		return "PID_DEVICE_IO_LIST"
	case KnxInterfaceObjectProperty_PID_DEVICE_MGT_DESCRIPTOR_01:
		return "PID_DEVICE_MGT_DESCRIPTOR_01"
	case KnxInterfaceObjectProperty_PID_DEVICE_PL110_PARAM:
		return "PID_DEVICE_PL110_PARAM"
	case KnxInterfaceObjectProperty_PID_DEVICE_RF_REPEAT_COUNTER:
		return "PID_DEVICE_RF_REPEAT_COUNTER"
	case KnxInterfaceObjectProperty_PID_DEVICE_RECEIVE_BLOCK_TABLE:
		return "PID_DEVICE_RECEIVE_BLOCK_TABLE"
	case KnxInterfaceObjectProperty_PID_DEVICE_RANDOM_PAUSE_TABLE:
		return "PID_DEVICE_RANDOM_PAUSE_TABLE"
	case KnxInterfaceObjectProperty_PID_DEVICE_RECEIVE_BLOCK_NR:
		return "PID_DEVICE_RECEIVE_BLOCK_NR"
	case KnxInterfaceObjectProperty_PID_DEVICE_HARDWARE_TYPE:
		return "PID_DEVICE_HARDWARE_TYPE"
	case KnxInterfaceObjectProperty_PID_DEVICE_RETRANSMITTER_NUMBER:
		return "PID_DEVICE_RETRANSMITTER_NUMBER"
	case KnxInterfaceObjectProperty_PID_DEVICE_SERIAL_NR_TABLE:
		return "PID_DEVICE_SERIAL_NR_TABLE"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REVERSION_PAUSE_TIME:
		return "PID_SUNBLIND_ACTUATOR_BASIC_REVERSION_PAUSE_TIME"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_UP_DOWN_TIME:
		return "PID_SUNBLIND_ACTUATOR_BASIC_MOVE_UP_DOWN_TIME"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_SLAT_STEP_TIME:
		return "PID_SUNBLIND_ACTUATOR_BASIC_SLAT_STEP_TIME"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_PRESET_POSITION_TIME:
		return "PID_SUNBLIND_ACTUATOR_BASIC_MOVE_PRESET_POSITION_TIME"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_TO_PRESET_POSITION_IN_PERCENT:
		return "PID_SUNBLIND_ACTUATOR_BASIC_MOVE_TO_PRESET_POSITION_IN_PERCENT"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MOVE_TO_PRESET_POSITION_LENGTH:
		return "PID_SUNBLIND_ACTUATOR_BASIC_MOVE_TO_PRESET_POSITION_LENGTH"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_PRESET_SLAT_POSITION_PERCENT:
		return "PID_SUNBLIND_ACTUATOR_BASIC_PRESET_SLAT_POSITION_PERCENT"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_PRESET_SLAT_POSITION_ANGLE:
		return "PID_SUNBLIND_ACTUATOR_BASIC_PRESET_SLAT_POSITION_ANGLE"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REACTION_WIND_ALARM:
		return "PID_SUNBLIND_ACTUATOR_BASIC_REACTION_WIND_ALARM"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_WIND_ALARM:
		return "PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_WIND_ALARM"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REACTION_ON_RAIN_ALARM:
		return "PID_SUNBLIND_ACTUATOR_BASIC_REACTION_ON_RAIN_ALARM"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_RAIN_ALARM:
		return "PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_RAIN_ALARM"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_REACTION_FROST_ALARM:
		return "PID_SUNBLIND_ACTUATOR_BASIC_REACTION_FROST_ALARM"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_FROST_ALARM:
		return "PID_SUNBLIND_ACTUATOR_BASIC_HEARTBEAT_FROST_ALARM"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_MAX_SLAT_MOVE_TIME:
		return "PID_SUNBLIND_ACTUATOR_BASIC_MAX_SLAT_MOVE_TIME"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_ENABLE_BLINDS_MODE:
		return "PID_SUNBLIND_ACTUATOR_BASIC_ENABLE_BLINDS_MODE"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_ACTUATOR_BASIC_STORAGE_FUNCTIONS_FOR_SCENE:
		return "PID_SUNBLIND_ACTUATOR_BASIC_STORAGE_FUNCTIONS_FOR_SCENE"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_SENSOR_BASIC_ENABLE_BLINDS_MODE:
		return "PID_SUNBLIND_SENSOR_BASIC_ENABLE_BLINDS_MODE"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_SENSOR_BASIC_UP_DOWN_ACTION:
		return "PID_SUNBLIND_SENSOR_BASIC_UP_DOWN_ACTION"
	case KnxInterfaceObjectProperty_PID_SUNBLIND_SENSOR_BASIC_ENABLE_TOGGLE_MODE:
		return "PID_SUNBLIND_SENSOR_BASIC_ENABLE_TOGGLE_MODE"
	case KnxInterfaceObjectProperty_PID_DEVICE_BIBATMASTER_ADDRESS:
		return "PID_DEVICE_BIBATMASTER_ADDRESS"
	case KnxInterfaceObjectProperty_PID_DEVICE_RF_DOMAIN_ADDRESS:
		return "PID_DEVICE_RF_DOMAIN_ADDRESS"
	case KnxInterfaceObjectProperty_PID_DEVICE_DEVICE_DESCRIPTOR:
		return "PID_DEVICE_DEVICE_DESCRIPTOR"
	case KnxInterfaceObjectProperty_PID_DEVICE_METERING_FILTER_TABLE:
		return "PID_DEVICE_METERING_FILTER_TABLE"
	case KnxInterfaceObjectProperty_PID_DEVICE_GROUP_TELEGR_RATE_LIMIT_TIME_BASE:
		return "PID_DEVICE_GROUP_TELEGR_RATE_LIMIT_TIME_BASE"
	case KnxInterfaceObjectProperty_PID_DEVICE_GROUP_TELEGR_RATE_LIMIT_NO_OF_TELEGR:
		return "PID_DEVICE_GROUP_TELEGR_RATE_LIMIT_NO_OF_TELEGR"
	case KnxInterfaceObjectProperty_PID_GROUP_OBJECT_TABLE_GRPOBJTABLE:
		return "PID_GROUP_OBJECT_TABLE_GRPOBJTABLE"
	case KnxInterfaceObjectProperty_PID_GROUP_OBJECT_TABLE_EXT_GRPOBJREFERENCE:
		return "PID_GROUP_OBJECT_TABLE_EXT_GRPOBJREFERENCE"
	}
	return ""
}
