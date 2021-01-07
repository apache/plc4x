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

type KnxDatapointSubtype uint32

type IKnxDatapointSubtype interface {
	Number() uint16
	DatapointType() KnxDatapointType
	Name() string
	Serialize(io utils.WriteBuffer) error
}

const (
	KnxDatapointSubtype_DPST_UNKNOWN                                  KnxDatapointSubtype = 0
	KnxDatapointSubtype_DPST_Switch                                   KnxDatapointSubtype = 1
	KnxDatapointSubtype_DPST_Bool                                     KnxDatapointSubtype = 2
	KnxDatapointSubtype_DPST_Enable                                   KnxDatapointSubtype = 3
	KnxDatapointSubtype_DPST_Ramp                                     KnxDatapointSubtype = 4
	KnxDatapointSubtype_DPST_Alarm                                    KnxDatapointSubtype = 5
	KnxDatapointSubtype_DPST_BinaryValue                              KnxDatapointSubtype = 6
	KnxDatapointSubtype_DPST_Step                                     KnxDatapointSubtype = 7
	KnxDatapointSubtype_DPST_UpDown                                   KnxDatapointSubtype = 8
	KnxDatapointSubtype_DPST_OpenClose                                KnxDatapointSubtype = 9
	KnxDatapointSubtype_DPST_Start                                    KnxDatapointSubtype = 10
	KnxDatapointSubtype_DPST_State                                    KnxDatapointSubtype = 11
	KnxDatapointSubtype_DPST_Invert                                   KnxDatapointSubtype = 12
	KnxDatapointSubtype_DPST_DimSendStyle                             KnxDatapointSubtype = 13
	KnxDatapointSubtype_DPST_InputSource                              KnxDatapointSubtype = 14
	KnxDatapointSubtype_DPST_Reset                                    KnxDatapointSubtype = 15
	KnxDatapointSubtype_DPST_Ack                                      KnxDatapointSubtype = 16
	KnxDatapointSubtype_DPST_Trigger                                  KnxDatapointSubtype = 17
	KnxDatapointSubtype_DPST_Occupancy                                KnxDatapointSubtype = 18
	KnxDatapointSubtype_DPST_Window_Door                              KnxDatapointSubtype = 19
	KnxDatapointSubtype_DPST_LogicalFunction                          KnxDatapointSubtype = 20
	KnxDatapointSubtype_DPST_Scene_AB                                 KnxDatapointSubtype = 21
	KnxDatapointSubtype_DPST_ShutterBlinds_Mode                       KnxDatapointSubtype = 22
	KnxDatapointSubtype_DPST_DayNight                                 KnxDatapointSubtype = 23
	KnxDatapointSubtype_DPST_Heat_Cool                                KnxDatapointSubtype = 24
	KnxDatapointSubtype_DPST_Switch_Control                           KnxDatapointSubtype = 25
	KnxDatapointSubtype_DPST_Bool_Control                             KnxDatapointSubtype = 26
	KnxDatapointSubtype_DPST_Enable_Control                           KnxDatapointSubtype = 27
	KnxDatapointSubtype_DPST_Ramp_Control                             KnxDatapointSubtype = 28
	KnxDatapointSubtype_DPST_Alarm_Control                            KnxDatapointSubtype = 29
	KnxDatapointSubtype_DPST_BinaryValue_Control                      KnxDatapointSubtype = 30
	KnxDatapointSubtype_DPST_Step_Control                             KnxDatapointSubtype = 31
	KnxDatapointSubtype_DPST_Direction1_Control                       KnxDatapointSubtype = 32
	KnxDatapointSubtype_DPST_Direction2_Control                       KnxDatapointSubtype = 33
	KnxDatapointSubtype_DPST_Start_Control                            KnxDatapointSubtype = 34
	KnxDatapointSubtype_DPST_State_Control                            KnxDatapointSubtype = 35
	KnxDatapointSubtype_DPST_Invert_Control                           KnxDatapointSubtype = 36
	KnxDatapointSubtype_DPST_Control_Dimming                          KnxDatapointSubtype = 37
	KnxDatapointSubtype_DPST_Control_Blinds                           KnxDatapointSubtype = 38
	KnxDatapointSubtype_DPST_Char_ASCII                               KnxDatapointSubtype = 39
	KnxDatapointSubtype_DPST_Char_8859_1                              KnxDatapointSubtype = 40
	KnxDatapointSubtype_DPST_Scaling                                  KnxDatapointSubtype = 41
	KnxDatapointSubtype_DPST_Angle                                    KnxDatapointSubtype = 42
	KnxDatapointSubtype_DPST_Percent_U8                               KnxDatapointSubtype = 43
	KnxDatapointSubtype_DPST_DecimalFactor                            KnxDatapointSubtype = 44
	KnxDatapointSubtype_DPST_Tariff                                   KnxDatapointSubtype = 45
	KnxDatapointSubtype_DPST_Value_1_Ucount                           KnxDatapointSubtype = 46
	KnxDatapointSubtype_DPST_FanStage                                 KnxDatapointSubtype = 47
	KnxDatapointSubtype_DPST_Percent_V8                               KnxDatapointSubtype = 48
	KnxDatapointSubtype_DPST_Value_1_Count                            KnxDatapointSubtype = 49
	KnxDatapointSubtype_DPST_Status_Mode3                             KnxDatapointSubtype = 50
	KnxDatapointSubtype_DPST_Value_2_Ucount                           KnxDatapointSubtype = 51
	KnxDatapointSubtype_DPST_TimePeriodMsec                           KnxDatapointSubtype = 52
	KnxDatapointSubtype_DPST_TimePeriod10Msec                         KnxDatapointSubtype = 53
	KnxDatapointSubtype_DPST_TimePeriod100Msec                        KnxDatapointSubtype = 54
	KnxDatapointSubtype_DPST_TimePeriodSec                            KnxDatapointSubtype = 55
	KnxDatapointSubtype_DPST_TimePeriodMin                            KnxDatapointSubtype = 56
	KnxDatapointSubtype_DPST_TimePeriodHrs                            KnxDatapointSubtype = 57
	KnxDatapointSubtype_DPST_PropDataType                             KnxDatapointSubtype = 58
	KnxDatapointSubtype_DPST_Length_mm                                KnxDatapointSubtype = 59
	KnxDatapointSubtype_DPST_UElCurrentmA                             KnxDatapointSubtype = 60
	KnxDatapointSubtype_DPST_Brightness                               KnxDatapointSubtype = 61
	KnxDatapointSubtype_DPST_Absolute_Colour_Temperature              KnxDatapointSubtype = 62
	KnxDatapointSubtype_DPST_Value_2_Count                            KnxDatapointSubtype = 63
	KnxDatapointSubtype_DPST_DeltaTimeMsec                            KnxDatapointSubtype = 64
	KnxDatapointSubtype_DPST_DeltaTime10Msec                          KnxDatapointSubtype = 65
	KnxDatapointSubtype_DPST_DeltaTime100Msec                         KnxDatapointSubtype = 66
	KnxDatapointSubtype_DPST_DeltaTimeSec                             KnxDatapointSubtype = 67
	KnxDatapointSubtype_DPST_DeltaTimeMin                             KnxDatapointSubtype = 68
	KnxDatapointSubtype_DPST_DeltaTimeHrs                             KnxDatapointSubtype = 69
	KnxDatapointSubtype_DPST_Percent_V16                              KnxDatapointSubtype = 70
	KnxDatapointSubtype_DPST_Rotation_Angle                           KnxDatapointSubtype = 71
	KnxDatapointSubtype_DPST_Length_m                                 KnxDatapointSubtype = 72
	KnxDatapointSubtype_DPST_Value_Temp                               KnxDatapointSubtype = 73
	KnxDatapointSubtype_DPST_Value_Tempd                              KnxDatapointSubtype = 74
	KnxDatapointSubtype_DPST_Value_Tempa                              KnxDatapointSubtype = 75
	KnxDatapointSubtype_DPST_Value_Lux                                KnxDatapointSubtype = 76
	KnxDatapointSubtype_DPST_Value_Wsp                                KnxDatapointSubtype = 77
	KnxDatapointSubtype_DPST_Value_Pres                               KnxDatapointSubtype = 78
	KnxDatapointSubtype_DPST_Value_Humidity                           KnxDatapointSubtype = 79
	KnxDatapointSubtype_DPST_Value_AirQuality                         KnxDatapointSubtype = 80
	KnxDatapointSubtype_DPST_Value_AirFlow                            KnxDatapointSubtype = 81
	KnxDatapointSubtype_DPST_Value_Time1                              KnxDatapointSubtype = 82
	KnxDatapointSubtype_DPST_Value_Time2                              KnxDatapointSubtype = 83
	KnxDatapointSubtype_DPST_Value_Volt                               KnxDatapointSubtype = 84
	KnxDatapointSubtype_DPST_Value_Curr                               KnxDatapointSubtype = 85
	KnxDatapointSubtype_DPST_PowerDensity                             KnxDatapointSubtype = 86
	KnxDatapointSubtype_DPST_KelvinPerPercent                         KnxDatapointSubtype = 87
	KnxDatapointSubtype_DPST_Power                                    KnxDatapointSubtype = 88
	KnxDatapointSubtype_DPST_Value_Volume_Flow                        KnxDatapointSubtype = 89
	KnxDatapointSubtype_DPST_Rain_Amount                              KnxDatapointSubtype = 90
	KnxDatapointSubtype_DPST_Value_Temp_F                             KnxDatapointSubtype = 91
	KnxDatapointSubtype_DPST_Value_Wsp_kmh                            KnxDatapointSubtype = 92
	KnxDatapointSubtype_DPST_Value_Absolute_Humidity                  KnxDatapointSubtype = 93
	KnxDatapointSubtype_DPST_Concentration_ygm3                       KnxDatapointSubtype = 94
	KnxDatapointSubtype_DPST_TimeOfDay                                KnxDatapointSubtype = 95
	KnxDatapointSubtype_DPST_Date                                     KnxDatapointSubtype = 96
	KnxDatapointSubtype_DPST_Value_4_Ucount                           KnxDatapointSubtype = 97
	KnxDatapointSubtype_DPST_LongTimePeriod_Sec                       KnxDatapointSubtype = 98
	KnxDatapointSubtype_DPST_LongTimePeriod_Min                       KnxDatapointSubtype = 99
	KnxDatapointSubtype_DPST_LongTimePeriod_Hrs                       KnxDatapointSubtype = 100
	KnxDatapointSubtype_DPST_VolumeLiquid_Litre                       KnxDatapointSubtype = 101
	KnxDatapointSubtype_DPST_Volume_m_3                               KnxDatapointSubtype = 102
	KnxDatapointSubtype_DPST_Value_4_Count                            KnxDatapointSubtype = 103
	KnxDatapointSubtype_DPST_FlowRate_m3h                             KnxDatapointSubtype = 104
	KnxDatapointSubtype_DPST_ActiveEnergy                             KnxDatapointSubtype = 105
	KnxDatapointSubtype_DPST_ApparantEnergy                           KnxDatapointSubtype = 106
	KnxDatapointSubtype_DPST_ReactiveEnergy                           KnxDatapointSubtype = 107
	KnxDatapointSubtype_DPST_ActiveEnergy_kWh                         KnxDatapointSubtype = 108
	KnxDatapointSubtype_DPST_ApparantEnergy_kVAh                      KnxDatapointSubtype = 109
	KnxDatapointSubtype_DPST_ReactiveEnergy_kVARh                     KnxDatapointSubtype = 110
	KnxDatapointSubtype_DPST_ActiveEnergy_MWh                         KnxDatapointSubtype = 111
	KnxDatapointSubtype_DPST_LongDeltaTimeSec                         KnxDatapointSubtype = 112
	KnxDatapointSubtype_DPST_DeltaVolumeLiquid_Litre                  KnxDatapointSubtype = 113
	KnxDatapointSubtype_DPST_DeltaVolume_m_3                          KnxDatapointSubtype = 114
	KnxDatapointSubtype_DPST_Value_Acceleration                       KnxDatapointSubtype = 115
	KnxDatapointSubtype_DPST_Value_Acceleration_Angular               KnxDatapointSubtype = 116
	KnxDatapointSubtype_DPST_Value_Activation_Energy                  KnxDatapointSubtype = 117
	KnxDatapointSubtype_DPST_Value_Activity                           KnxDatapointSubtype = 118
	KnxDatapointSubtype_DPST_Value_Mol                                KnxDatapointSubtype = 119
	KnxDatapointSubtype_DPST_Value_Amplitude                          KnxDatapointSubtype = 120
	KnxDatapointSubtype_DPST_Value_AngleRad                           KnxDatapointSubtype = 121
	KnxDatapointSubtype_DPST_Value_AngleDeg                           KnxDatapointSubtype = 122
	KnxDatapointSubtype_DPST_Value_Angular_Momentum                   KnxDatapointSubtype = 123
	KnxDatapointSubtype_DPST_Value_Angular_Velocity                   KnxDatapointSubtype = 124
	KnxDatapointSubtype_DPST_Value_Area                               KnxDatapointSubtype = 125
	KnxDatapointSubtype_DPST_Value_Capacitance                        KnxDatapointSubtype = 126
	KnxDatapointSubtype_DPST_Value_Charge_DensitySurface              KnxDatapointSubtype = 127
	KnxDatapointSubtype_DPST_Value_Charge_DensityVolume               KnxDatapointSubtype = 128
	KnxDatapointSubtype_DPST_Value_Compressibility                    KnxDatapointSubtype = 129
	KnxDatapointSubtype_DPST_Value_Conductance                        KnxDatapointSubtype = 130
	KnxDatapointSubtype_DPST_Value_Electrical_Conductivity            KnxDatapointSubtype = 131
	KnxDatapointSubtype_DPST_Value_Density                            KnxDatapointSubtype = 132
	KnxDatapointSubtype_DPST_Value_Electric_Charge                    KnxDatapointSubtype = 133
	KnxDatapointSubtype_DPST_Value_Electric_Current                   KnxDatapointSubtype = 134
	KnxDatapointSubtype_DPST_Value_Electric_CurrentDensity            KnxDatapointSubtype = 135
	KnxDatapointSubtype_DPST_Value_Electric_DipoleMoment              KnxDatapointSubtype = 136
	KnxDatapointSubtype_DPST_Value_Electric_Displacement              KnxDatapointSubtype = 137
	KnxDatapointSubtype_DPST_Value_Electric_FieldStrength             KnxDatapointSubtype = 138
	KnxDatapointSubtype_DPST_Value_Electric_Flux                      KnxDatapointSubtype = 139
	KnxDatapointSubtype_DPST_Value_Electric_FluxDensity               KnxDatapointSubtype = 140
	KnxDatapointSubtype_DPST_Value_Electric_Polarization              KnxDatapointSubtype = 141
	KnxDatapointSubtype_DPST_Value_Electric_Potential                 KnxDatapointSubtype = 142
	KnxDatapointSubtype_DPST_Value_Electric_PotentialDifference       KnxDatapointSubtype = 143
	KnxDatapointSubtype_DPST_Value_ElectromagneticMoment              KnxDatapointSubtype = 144
	KnxDatapointSubtype_DPST_Value_Electromotive_Force                KnxDatapointSubtype = 145
	KnxDatapointSubtype_DPST_Value_Energy                             KnxDatapointSubtype = 146
	KnxDatapointSubtype_DPST_Value_Force                              KnxDatapointSubtype = 147
	KnxDatapointSubtype_DPST_Value_Frequency                          KnxDatapointSubtype = 148
	KnxDatapointSubtype_DPST_Value_Angular_Frequency                  KnxDatapointSubtype = 149
	KnxDatapointSubtype_DPST_Value_Heat_Capacity                      KnxDatapointSubtype = 150
	KnxDatapointSubtype_DPST_Value_Heat_FlowRate                      KnxDatapointSubtype = 151
	KnxDatapointSubtype_DPST_Value_Heat_Quantity                      KnxDatapointSubtype = 152
	KnxDatapointSubtype_DPST_Value_Impedance                          KnxDatapointSubtype = 153
	KnxDatapointSubtype_DPST_Value_Length                             KnxDatapointSubtype = 154
	KnxDatapointSubtype_DPST_Value_Light_Quantity                     KnxDatapointSubtype = 155
	KnxDatapointSubtype_DPST_Value_Luminance                          KnxDatapointSubtype = 156
	KnxDatapointSubtype_DPST_Value_Luminous_Flux                      KnxDatapointSubtype = 157
	KnxDatapointSubtype_DPST_Value_Luminous_Intensity                 KnxDatapointSubtype = 158
	KnxDatapointSubtype_DPST_Value_Magnetic_FieldStrength             KnxDatapointSubtype = 159
	KnxDatapointSubtype_DPST_Value_Magnetic_Flux                      KnxDatapointSubtype = 160
	KnxDatapointSubtype_DPST_Value_Magnetic_FluxDensity               KnxDatapointSubtype = 161
	KnxDatapointSubtype_DPST_Value_Magnetic_Moment                    KnxDatapointSubtype = 162
	KnxDatapointSubtype_DPST_Value_Magnetic_Polarization              KnxDatapointSubtype = 163
	KnxDatapointSubtype_DPST_Value_Magnetization                      KnxDatapointSubtype = 164
	KnxDatapointSubtype_DPST_Value_MagnetomotiveForce                 KnxDatapointSubtype = 165
	KnxDatapointSubtype_DPST_Value_Mass                               KnxDatapointSubtype = 166
	KnxDatapointSubtype_DPST_Value_MassFlux                           KnxDatapointSubtype = 167
	KnxDatapointSubtype_DPST_Value_Momentum                           KnxDatapointSubtype = 168
	KnxDatapointSubtype_DPST_Value_Phase_AngleRad                     KnxDatapointSubtype = 169
	KnxDatapointSubtype_DPST_Value_Phase_AngleDeg                     KnxDatapointSubtype = 170
	KnxDatapointSubtype_DPST_Value_Power                              KnxDatapointSubtype = 171
	KnxDatapointSubtype_DPST_Value_Power_Factor                       KnxDatapointSubtype = 172
	KnxDatapointSubtype_DPST_Value_Pressure                           KnxDatapointSubtype = 173
	KnxDatapointSubtype_DPST_Value_Reactance                          KnxDatapointSubtype = 174
	KnxDatapointSubtype_DPST_Value_Resistance                         KnxDatapointSubtype = 175
	KnxDatapointSubtype_DPST_Value_Resistivity                        KnxDatapointSubtype = 176
	KnxDatapointSubtype_DPST_Value_SelfInductance                     KnxDatapointSubtype = 177
	KnxDatapointSubtype_DPST_Value_SolidAngle                         KnxDatapointSubtype = 178
	KnxDatapointSubtype_DPST_Value_Sound_Intensity                    KnxDatapointSubtype = 179
	KnxDatapointSubtype_DPST_Value_Speed                              KnxDatapointSubtype = 180
	KnxDatapointSubtype_DPST_Value_Stress                             KnxDatapointSubtype = 181
	KnxDatapointSubtype_DPST_Value_Surface_Tension                    KnxDatapointSubtype = 182
	KnxDatapointSubtype_DPST_Value_Common_Temperature                 KnxDatapointSubtype = 183
	KnxDatapointSubtype_DPST_Value_Absolute_Temperature               KnxDatapointSubtype = 184
	KnxDatapointSubtype_DPST_Value_TemperatureDifference              KnxDatapointSubtype = 185
	KnxDatapointSubtype_DPST_Value_Thermal_Capacity                   KnxDatapointSubtype = 186
	KnxDatapointSubtype_DPST_Value_Thermal_Conductivity               KnxDatapointSubtype = 187
	KnxDatapointSubtype_DPST_Value_ThermoelectricPower                KnxDatapointSubtype = 188
	KnxDatapointSubtype_DPST_Value_Time                               KnxDatapointSubtype = 189
	KnxDatapointSubtype_DPST_Value_Torque                             KnxDatapointSubtype = 190
	KnxDatapointSubtype_DPST_Value_Volume                             KnxDatapointSubtype = 191
	KnxDatapointSubtype_DPST_Value_Volume_Flux                        KnxDatapointSubtype = 192
	KnxDatapointSubtype_DPST_Value_Weight                             KnxDatapointSubtype = 193
	KnxDatapointSubtype_DPST_Value_Work                               KnxDatapointSubtype = 194
	KnxDatapointSubtype_DPST_Volume_Flux_Meter                        KnxDatapointSubtype = 195
	KnxDatapointSubtype_DPST_Volume_Flux_ls                           KnxDatapointSubtype = 196
	KnxDatapointSubtype_DPST_Access_Data                              KnxDatapointSubtype = 197
	KnxDatapointSubtype_DPST_String_ASCII                             KnxDatapointSubtype = 198
	KnxDatapointSubtype_DPST_String_8859_1                            KnxDatapointSubtype = 199
	KnxDatapointSubtype_DPST_SceneNumber                              KnxDatapointSubtype = 200
	KnxDatapointSubtype_DPST_SceneControl                             KnxDatapointSubtype = 201
	KnxDatapointSubtype_DPST_DateTime                                 KnxDatapointSubtype = 202
	KnxDatapointSubtype_DPST_SCLOMode                                 KnxDatapointSubtype = 203
	KnxDatapointSubtype_DPST_BuildingMode                             KnxDatapointSubtype = 204
	KnxDatapointSubtype_DPST_OccMode                                  KnxDatapointSubtype = 205
	KnxDatapointSubtype_DPST_Priority                                 KnxDatapointSubtype = 206
	KnxDatapointSubtype_DPST_LightApplicationMode                     KnxDatapointSubtype = 207
	KnxDatapointSubtype_DPST_ApplicationArea                          KnxDatapointSubtype = 208
	KnxDatapointSubtype_DPST_AlarmClassType                           KnxDatapointSubtype = 209
	KnxDatapointSubtype_DPST_PSUMode                                  KnxDatapointSubtype = 210
	KnxDatapointSubtype_DPST_ErrorClass_System                        KnxDatapointSubtype = 211
	KnxDatapointSubtype_DPST_ErrorClass_HVAC                          KnxDatapointSubtype = 212
	KnxDatapointSubtype_DPST_Time_Delay                               KnxDatapointSubtype = 213
	KnxDatapointSubtype_DPST_Beaufort_Wind_Force_Scale                KnxDatapointSubtype = 214
	KnxDatapointSubtype_DPST_SensorSelect                             KnxDatapointSubtype = 215
	KnxDatapointSubtype_DPST_ActuatorConnectType                      KnxDatapointSubtype = 216
	KnxDatapointSubtype_DPST_Cloud_Cover                              KnxDatapointSubtype = 217
	KnxDatapointSubtype_DPST_PowerReturnMode                          KnxDatapointSubtype = 218
	KnxDatapointSubtype_DPST_FuelType                                 KnxDatapointSubtype = 219
	KnxDatapointSubtype_DPST_BurnerType                               KnxDatapointSubtype = 220
	KnxDatapointSubtype_DPST_HVACMode                                 KnxDatapointSubtype = 221
	KnxDatapointSubtype_DPST_DHWMode                                  KnxDatapointSubtype = 222
	KnxDatapointSubtype_DPST_LoadPriority                             KnxDatapointSubtype = 223
	KnxDatapointSubtype_DPST_HVACContrMode                            KnxDatapointSubtype = 224
	KnxDatapointSubtype_DPST_HVACEmergMode                            KnxDatapointSubtype = 225
	KnxDatapointSubtype_DPST_ChangeoverMode                           KnxDatapointSubtype = 226
	KnxDatapointSubtype_DPST_ValveMode                                KnxDatapointSubtype = 227
	KnxDatapointSubtype_DPST_DamperMode                               KnxDatapointSubtype = 228
	KnxDatapointSubtype_DPST_HeaterMode                               KnxDatapointSubtype = 229
	KnxDatapointSubtype_DPST_FanMode                                  KnxDatapointSubtype = 230
	KnxDatapointSubtype_DPST_MasterSlaveMode                          KnxDatapointSubtype = 231
	KnxDatapointSubtype_DPST_StatusRoomSetp                           KnxDatapointSubtype = 232
	KnxDatapointSubtype_DPST_Metering_DeviceType                      KnxDatapointSubtype = 233
	KnxDatapointSubtype_DPST_HumDehumMode                             KnxDatapointSubtype = 234
	KnxDatapointSubtype_DPST_EnableHCStage                            KnxDatapointSubtype = 235
	KnxDatapointSubtype_DPST_ADAType                                  KnxDatapointSubtype = 236
	KnxDatapointSubtype_DPST_BackupMode                               KnxDatapointSubtype = 237
	KnxDatapointSubtype_DPST_StartSynchronization                     KnxDatapointSubtype = 238
	KnxDatapointSubtype_DPST_Behaviour_Lock_Unlock                    KnxDatapointSubtype = 239
	KnxDatapointSubtype_DPST_Behaviour_Bus_Power_Up_Down              KnxDatapointSubtype = 240
	KnxDatapointSubtype_DPST_DALI_Fade_Time                           KnxDatapointSubtype = 241
	KnxDatapointSubtype_DPST_BlinkingMode                             KnxDatapointSubtype = 242
	KnxDatapointSubtype_DPST_LightControlMode                         KnxDatapointSubtype = 243
	KnxDatapointSubtype_DPST_SwitchPBModel                            KnxDatapointSubtype = 244
	KnxDatapointSubtype_DPST_PBAction                                 KnxDatapointSubtype = 245
	KnxDatapointSubtype_DPST_DimmPBModel                              KnxDatapointSubtype = 246
	KnxDatapointSubtype_DPST_SwitchOnMode                             KnxDatapointSubtype = 247
	KnxDatapointSubtype_DPST_LoadTypeSet                              KnxDatapointSubtype = 248
	KnxDatapointSubtype_DPST_LoadTypeDetected                         KnxDatapointSubtype = 249
	KnxDatapointSubtype_DPST_Converter_Test_Control                   KnxDatapointSubtype = 250
	KnxDatapointSubtype_DPST_SABExcept_Behaviour                      KnxDatapointSubtype = 251
	KnxDatapointSubtype_DPST_SABBehaviour_Lock_Unlock                 KnxDatapointSubtype = 252
	KnxDatapointSubtype_DPST_SSSBMode                                 KnxDatapointSubtype = 253
	KnxDatapointSubtype_DPST_BlindsControlMode                        KnxDatapointSubtype = 254
	KnxDatapointSubtype_DPST_CommMode                                 KnxDatapointSubtype = 255
	KnxDatapointSubtype_DPST_AddInfoTypes                             KnxDatapointSubtype = 256
	KnxDatapointSubtype_DPST_RF_ModeSelect                            KnxDatapointSubtype = 257
	KnxDatapointSubtype_DPST_RF_FilterSelect                          KnxDatapointSubtype = 258
	KnxDatapointSubtype_DPST_StatusGen                                KnxDatapointSubtype = 259
	KnxDatapointSubtype_DPST_Device_Control                           KnxDatapointSubtype = 260
	KnxDatapointSubtype_DPST_ForceSign                                KnxDatapointSubtype = 261
	KnxDatapointSubtype_DPST_ForceSignCool                            KnxDatapointSubtype = 262
	KnxDatapointSubtype_DPST_StatusRHC                                KnxDatapointSubtype = 263
	KnxDatapointSubtype_DPST_StatusSDHWC                              KnxDatapointSubtype = 264
	KnxDatapointSubtype_DPST_FuelTypeSet                              KnxDatapointSubtype = 265
	KnxDatapointSubtype_DPST_StatusRCC                                KnxDatapointSubtype = 266
	KnxDatapointSubtype_DPST_StatusAHU                                KnxDatapointSubtype = 267
	KnxDatapointSubtype_DPST_CombinedStatus_RTSM                      KnxDatapointSubtype = 268
	KnxDatapointSubtype_DPST_LightActuatorErrorInfo                   KnxDatapointSubtype = 269
	KnxDatapointSubtype_DPST_RF_ModeInfo                              KnxDatapointSubtype = 270
	KnxDatapointSubtype_DPST_RF_FilterInfo                            KnxDatapointSubtype = 271
	KnxDatapointSubtype_DPST_Channel_Activation_8                     KnxDatapointSubtype = 272
	KnxDatapointSubtype_DPST_StatusDHWC                               KnxDatapointSubtype = 273
	KnxDatapointSubtype_DPST_StatusRHCC                               KnxDatapointSubtype = 274
	KnxDatapointSubtype_DPST_CombinedStatus_HVA                       KnxDatapointSubtype = 275
	KnxDatapointSubtype_DPST_CombinedStatus_RTC                       KnxDatapointSubtype = 276
	KnxDatapointSubtype_DPST_Media                                    KnxDatapointSubtype = 277
	KnxDatapointSubtype_DPST_Channel_Activation_16                    KnxDatapointSubtype = 278
	KnxDatapointSubtype_DPST_OnOffAction                              KnxDatapointSubtype = 279
	KnxDatapointSubtype_DPST_Alarm_Reaction                           KnxDatapointSubtype = 280
	KnxDatapointSubtype_DPST_UpDown_Action                            KnxDatapointSubtype = 281
	KnxDatapointSubtype_DPST_HVAC_PB_Action                           KnxDatapointSubtype = 282
	KnxDatapointSubtype_DPST_DoubleNibble                             KnxDatapointSubtype = 283
	KnxDatapointSubtype_DPST_SceneInfo                                KnxDatapointSubtype = 284
	KnxDatapointSubtype_DPST_CombinedInfoOnOff                        KnxDatapointSubtype = 285
	KnxDatapointSubtype_DPST_ActiveEnergy_V64                         KnxDatapointSubtype = 286
	KnxDatapointSubtype_DPST_ApparantEnergy_V64                       KnxDatapointSubtype = 287
	KnxDatapointSubtype_DPST_ReactiveEnergy_V64                       KnxDatapointSubtype = 288
	KnxDatapointSubtype_DPST_Channel_Activation_24                    KnxDatapointSubtype = 289
	KnxDatapointSubtype_DPST_HVACModeNext                             KnxDatapointSubtype = 290
	KnxDatapointSubtype_DPST_DHWModeNext                              KnxDatapointSubtype = 291
	KnxDatapointSubtype_DPST_OccModeNext                              KnxDatapointSubtype = 292
	KnxDatapointSubtype_DPST_BuildingModeNext                         KnxDatapointSubtype = 293
	KnxDatapointSubtype_DPST_Version                                  KnxDatapointSubtype = 294
	KnxDatapointSubtype_DPST_AlarmInfo                                KnxDatapointSubtype = 295
	KnxDatapointSubtype_DPST_TempRoomSetpSetF16_3                     KnxDatapointSubtype = 296
	KnxDatapointSubtype_DPST_TempRoomSetpSetShiftF16_3                KnxDatapointSubtype = 297
	KnxDatapointSubtype_DPST_Scaling_Speed                            KnxDatapointSubtype = 298
	KnxDatapointSubtype_DPST_Scaling_Step_Time                        KnxDatapointSubtype = 299
	KnxDatapointSubtype_DPST_MeteringValue                            KnxDatapointSubtype = 300
	KnxDatapointSubtype_DPST_MBus_Address                             KnxDatapointSubtype = 301
	KnxDatapointSubtype_DPST_Colour_RGB                               KnxDatapointSubtype = 302
	KnxDatapointSubtype_DPST_LanguageCodeAlpha2_ASCII                 KnxDatapointSubtype = 303
	KnxDatapointSubtype_DPST_Tariff_ActiveEnergy                      KnxDatapointSubtype = 304
	KnxDatapointSubtype_DPST_Prioritised_Mode_Control                 KnxDatapointSubtype = 305
	KnxDatapointSubtype_DPST_DALI_Control_Gear_Diagnostic             KnxDatapointSubtype = 306
	KnxDatapointSubtype_DPST_DALI_Diagnostics                         KnxDatapointSubtype = 307
	KnxDatapointSubtype_DPST_CombinedPosition                         KnxDatapointSubtype = 308
	KnxDatapointSubtype_DPST_StatusSAB                                KnxDatapointSubtype = 309
	KnxDatapointSubtype_DPST_Colour_xyY                               KnxDatapointSubtype = 310
	KnxDatapointSubtype_DPST_Converter_Status                         KnxDatapointSubtype = 311
	KnxDatapointSubtype_DPST_Converter_Test_Result                    KnxDatapointSubtype = 312
	KnxDatapointSubtype_DPST_Battery_Info                             KnxDatapointSubtype = 313
	KnxDatapointSubtype_DPST_Brightness_Colour_Temperature_Transition KnxDatapointSubtype = 314
	KnxDatapointSubtype_DPST_Brightness_Colour_Temperature_Control    KnxDatapointSubtype = 315
	KnxDatapointSubtype_DPST_Colour_RGBW                              KnxDatapointSubtype = 316
	KnxDatapointSubtype_DPST_Relative_Control_RGBW                    KnxDatapointSubtype = 317
	KnxDatapointSubtype_DPST_Relative_Control_RGB                     KnxDatapointSubtype = 318
	KnxDatapointSubtype_DPST_GeographicalLocation                     KnxDatapointSubtype = 319
	KnxDatapointSubtype_DPST_TempRoomSetpSetF16_4                     KnxDatapointSubtype = 320
	KnxDatapointSubtype_DPST_TempRoomSetpSetShiftF16_4                KnxDatapointSubtype = 321
)

func (e KnxDatapointSubtype) Number() uint16 {
	switch e {
	case 0:
		{ /* '0' */
			return 0
		}
	case 1:
		{ /* '1' */
			return 1
		}
	case 10:
		{ /* '10' */
			return 10
		}
	case 100:
		{ /* '100' */
			return 102
		}
	case 101:
		{ /* '101' */
			return 1200
		}
	case 102:
		{ /* '102' */
			return 1201
		}
	case 103:
		{ /* '103' */
			return 1
		}
	case 104:
		{ /* '104' */
			return 2
		}
	case 105:
		{ /* '105' */
			return 10
		}
	case 106:
		{ /* '106' */
			return 11
		}
	case 107:
		{ /* '107' */
			return 12
		}
	case 108:
		{ /* '108' */
			return 13
		}
	case 109:
		{ /* '109' */
			return 14
		}
	case 11:
		{ /* '11' */
			return 11
		}
	case 110:
		{ /* '110' */
			return 15
		}
	case 111:
		{ /* '111' */
			return 16
		}
	case 112:
		{ /* '112' */
			return 100
		}
	case 113:
		{ /* '113' */
			return 1200
		}
	case 114:
		{ /* '114' */
			return 1201
		}
	case 115:
		{ /* '115' */
			return 0
		}
	case 116:
		{ /* '116' */
			return 1
		}
	case 117:
		{ /* '117' */
			return 2
		}
	case 118:
		{ /* '118' */
			return 3
		}
	case 119:
		{ /* '119' */
			return 4
		}
	case 12:
		{ /* '12' */
			return 12
		}
	case 120:
		{ /* '120' */
			return 5
		}
	case 121:
		{ /* '121' */
			return 6
		}
	case 122:
		{ /* '122' */
			return 7
		}
	case 123:
		{ /* '123' */
			return 8
		}
	case 124:
		{ /* '124' */
			return 9
		}
	case 125:
		{ /* '125' */
			return 10
		}
	case 126:
		{ /* '126' */
			return 11
		}
	case 127:
		{ /* '127' */
			return 12
		}
	case 128:
		{ /* '128' */
			return 13
		}
	case 129:
		{ /* '129' */
			return 14
		}
	case 13:
		{ /* '13' */
			return 13
		}
	case 130:
		{ /* '130' */
			return 15
		}
	case 131:
		{ /* '131' */
			return 16
		}
	case 132:
		{ /* '132' */
			return 17
		}
	case 133:
		{ /* '133' */
			return 18
		}
	case 134:
		{ /* '134' */
			return 19
		}
	case 135:
		{ /* '135' */
			return 20
		}
	case 136:
		{ /* '136' */
			return 21
		}
	case 137:
		{ /* '137' */
			return 22
		}
	case 138:
		{ /* '138' */
			return 23
		}
	case 139:
		{ /* '139' */
			return 24
		}
	case 14:
		{ /* '14' */
			return 14
		}
	case 140:
		{ /* '140' */
			return 25
		}
	case 141:
		{ /* '141' */
			return 26
		}
	case 142:
		{ /* '142' */
			return 27
		}
	case 143:
		{ /* '143' */
			return 28
		}
	case 144:
		{ /* '144' */
			return 29
		}
	case 145:
		{ /* '145' */
			return 30
		}
	case 146:
		{ /* '146' */
			return 31
		}
	case 147:
		{ /* '147' */
			return 32
		}
	case 148:
		{ /* '148' */
			return 33
		}
	case 149:
		{ /* '149' */
			return 34
		}
	case 15:
		{ /* '15' */
			return 15
		}
	case 150:
		{ /* '150' */
			return 35
		}
	case 151:
		{ /* '151' */
			return 36
		}
	case 152:
		{ /* '152' */
			return 37
		}
	case 153:
		{ /* '153' */
			return 38
		}
	case 154:
		{ /* '154' */
			return 39
		}
	case 155:
		{ /* '155' */
			return 40
		}
	case 156:
		{ /* '156' */
			return 41
		}
	case 157:
		{ /* '157' */
			return 42
		}
	case 158:
		{ /* '158' */
			return 43
		}
	case 159:
		{ /* '159' */
			return 44
		}
	case 16:
		{ /* '16' */
			return 16
		}
	case 160:
		{ /* '160' */
			return 45
		}
	case 161:
		{ /* '161' */
			return 46
		}
	case 162:
		{ /* '162' */
			return 47
		}
	case 163:
		{ /* '163' */
			return 48
		}
	case 164:
		{ /* '164' */
			return 49
		}
	case 165:
		{ /* '165' */
			return 50
		}
	case 166:
		{ /* '166' */
			return 51
		}
	case 167:
		{ /* '167' */
			return 52
		}
	case 168:
		{ /* '168' */
			return 53
		}
	case 169:
		{ /* '169' */
			return 54
		}
	case 17:
		{ /* '17' */
			return 17
		}
	case 170:
		{ /* '170' */
			return 55
		}
	case 171:
		{ /* '171' */
			return 56
		}
	case 172:
		{ /* '172' */
			return 57
		}
	case 173:
		{ /* '173' */
			return 58
		}
	case 174:
		{ /* '174' */
			return 59
		}
	case 175:
		{ /* '175' */
			return 60
		}
	case 176:
		{ /* '176' */
			return 61
		}
	case 177:
		{ /* '177' */
			return 62
		}
	case 178:
		{ /* '178' */
			return 63
		}
	case 179:
		{ /* '179' */
			return 64
		}
	case 18:
		{ /* '18' */
			return 18
		}
	case 180:
		{ /* '180' */
			return 65
		}
	case 181:
		{ /* '181' */
			return 66
		}
	case 182:
		{ /* '182' */
			return 67
		}
	case 183:
		{ /* '183' */
			return 68
		}
	case 184:
		{ /* '184' */
			return 69
		}
	case 185:
		{ /* '185' */
			return 70
		}
	case 186:
		{ /* '186' */
			return 71
		}
	case 187:
		{ /* '187' */
			return 72
		}
	case 188:
		{ /* '188' */
			return 73
		}
	case 189:
		{ /* '189' */
			return 74
		}
	case 19:
		{ /* '19' */
			return 19
		}
	case 190:
		{ /* '190' */
			return 75
		}
	case 191:
		{ /* '191' */
			return 76
		}
	case 192:
		{ /* '192' */
			return 77
		}
	case 193:
		{ /* '193' */
			return 78
		}
	case 194:
		{ /* '194' */
			return 79
		}
	case 195:
		{ /* '195' */
			return 1200
		}
	case 196:
		{ /* '196' */
			return 1201
		}
	case 197:
		{ /* '197' */
			return 0
		}
	case 198:
		{ /* '198' */
			return 0
		}
	case 199:
		{ /* '199' */
			return 1
		}
	case 2:
		{ /* '2' */
			return 2
		}
	case 20:
		{ /* '20' */
			return 21
		}
	case 200:
		{ /* '200' */
			return 1
		}
	case 201:
		{ /* '201' */
			return 1
		}
	case 202:
		{ /* '202' */
			return 1
		}
	case 203:
		{ /* '203' */
			return 1
		}
	case 204:
		{ /* '204' */
			return 2
		}
	case 205:
		{ /* '205' */
			return 3
		}
	case 206:
		{ /* '206' */
			return 4
		}
	case 207:
		{ /* '207' */
			return 5
		}
	case 208:
		{ /* '208' */
			return 6
		}
	case 209:
		{ /* '209' */
			return 7
		}
	case 21:
		{ /* '21' */
			return 22
		}
	case 210:
		{ /* '210' */
			return 8
		}
	case 211:
		{ /* '211' */
			return 11
		}
	case 212:
		{ /* '212' */
			return 12
		}
	case 213:
		{ /* '213' */
			return 13
		}
	case 214:
		{ /* '214' */
			return 14
		}
	case 215:
		{ /* '215' */
			return 17
		}
	case 216:
		{ /* '216' */
			return 20
		}
	case 217:
		{ /* '217' */
			return 21
		}
	case 218:
		{ /* '218' */
			return 22
		}
	case 219:
		{ /* '219' */
			return 100
		}
	case 22:
		{ /* '22' */
			return 23
		}
	case 220:
		{ /* '220' */
			return 101
		}
	case 221:
		{ /* '221' */
			return 102
		}
	case 222:
		{ /* '222' */
			return 103
		}
	case 223:
		{ /* '223' */
			return 104
		}
	case 224:
		{ /* '224' */
			return 105
		}
	case 225:
		{ /* '225' */
			return 106
		}
	case 226:
		{ /* '226' */
			return 107
		}
	case 227:
		{ /* '227' */
			return 108
		}
	case 228:
		{ /* '228' */
			return 109
		}
	case 229:
		{ /* '229' */
			return 110
		}
	case 23:
		{ /* '23' */
			return 24
		}
	case 230:
		{ /* '230' */
			return 111
		}
	case 231:
		{ /* '231' */
			return 112
		}
	case 232:
		{ /* '232' */
			return 113
		}
	case 233:
		{ /* '233' */
			return 114
		}
	case 234:
		{ /* '234' */
			return 115
		}
	case 235:
		{ /* '235' */
			return 116
		}
	case 236:
		{ /* '236' */
			return 120
		}
	case 237:
		{ /* '237' */
			return 121
		}
	case 238:
		{ /* '238' */
			return 122
		}
	case 239:
		{ /* '239' */
			return 600
		}
	case 24:
		{ /* '24' */
			return 100
		}
	case 240:
		{ /* '240' */
			return 601
		}
	case 241:
		{ /* '241' */
			return 602
		}
	case 242:
		{ /* '242' */
			return 603
		}
	case 243:
		{ /* '243' */
			return 604
		}
	case 244:
		{ /* '244' */
			return 605
		}
	case 245:
		{ /* '245' */
			return 606
		}
	case 246:
		{ /* '246' */
			return 607
		}
	case 247:
		{ /* '247' */
			return 608
		}
	case 248:
		{ /* '248' */
			return 609
		}
	case 249:
		{ /* '249' */
			return 610
		}
	case 25:
		{ /* '25' */
			return 1
		}
	case 250:
		{ /* '250' */
			return 611
		}
	case 251:
		{ /* '251' */
			return 801
		}
	case 252:
		{ /* '252' */
			return 802
		}
	case 253:
		{ /* '253' */
			return 803
		}
	case 254:
		{ /* '254' */
			return 804
		}
	case 255:
		{ /* '255' */
			return 1000
		}
	case 256:
		{ /* '256' */
			return 1001
		}
	case 257:
		{ /* '257' */
			return 1002
		}
	case 258:
		{ /* '258' */
			return 1003
		}
	case 259:
		{ /* '259' */
			return 1
		}
	case 26:
		{ /* '26' */
			return 2
		}
	case 260:
		{ /* '260' */
			return 2
		}
	case 261:
		{ /* '261' */
			return 100
		}
	case 262:
		{ /* '262' */
			return 101
		}
	case 263:
		{ /* '263' */
			return 102
		}
	case 264:
		{ /* '264' */
			return 103
		}
	case 265:
		{ /* '265' */
			return 104
		}
	case 266:
		{ /* '266' */
			return 105
		}
	case 267:
		{ /* '267' */
			return 106
		}
	case 268:
		{ /* '268' */
			return 107
		}
	case 269:
		{ /* '269' */
			return 601
		}
	case 27:
		{ /* '27' */
			return 3
		}
	case 270:
		{ /* '270' */
			return 1000
		}
	case 271:
		{ /* '271' */
			return 1001
		}
	case 272:
		{ /* '272' */
			return 1010
		}
	case 273:
		{ /* '273' */
			return 100
		}
	case 274:
		{ /* '274' */
			return 101
		}
	case 275:
		{ /* '275' */
			return 102
		}
	case 276:
		{ /* '276' */
			return 103
		}
	case 277:
		{ /* '277' */
			return 1000
		}
	case 278:
		{ /* '278' */
			return 1010
		}
	case 279:
		{ /* '279' */
			return 1
		}
	case 28:
		{ /* '28' */
			return 4
		}
	case 280:
		{ /* '280' */
			return 2
		}
	case 281:
		{ /* '281' */
			return 3
		}
	case 282:
		{ /* '282' */
			return 102
		}
	case 283:
		{ /* '283' */
			return 1000
		}
	case 284:
		{ /* '284' */
			return 1
		}
	case 285:
		{ /* '285' */
			return 1
		}
	case 286:
		{ /* '286' */
			return 10
		}
	case 287:
		{ /* '287' */
			return 11
		}
	case 288:
		{ /* '288' */
			return 12
		}
	case 289:
		{ /* '289' */
			return 1010
		}
	case 29:
		{ /* '29' */
			return 5
		}
	case 290:
		{ /* '290' */
			return 100
		}
	case 291:
		{ /* '291' */
			return 102
		}
	case 292:
		{ /* '292' */
			return 104
		}
	case 293:
		{ /* '293' */
			return 105
		}
	case 294:
		{ /* '294' */
			return 1
		}
	case 295:
		{ /* '295' */
			return 1
		}
	case 296:
		{ /* '296' */
			return 100
		}
	case 297:
		{ /* '297' */
			return 101
		}
	case 298:
		{ /* '298' */
			return 1
		}
	case 299:
		{ /* '299' */
			return 2
		}
	case 3:
		{ /* '3' */
			return 3
		}
	case 30:
		{ /* '30' */
			return 6
		}
	case 300:
		{ /* '300' */
			return 1
		}
	case 301:
		{ /* '301' */
			return 1000
		}
	case 302:
		{ /* '302' */
			return 600
		}
	case 303:
		{ /* '303' */
			return 1
		}
	case 304:
		{ /* '304' */
			return 1
		}
	case 305:
		{ /* '305' */
			return 1
		}
	case 306:
		{ /* '306' */
			return 600
		}
	case 307:
		{ /* '307' */
			return 600
		}
	case 308:
		{ /* '308' */
			return 800
		}
	case 309:
		{ /* '309' */
			return 800
		}
	case 31:
		{ /* '31' */
			return 7
		}
	case 310:
		{ /* '310' */
			return 600
		}
	case 311:
		{ /* '311' */
			return 600
		}
	case 312:
		{ /* '312' */
			return 600
		}
	case 313:
		{ /* '313' */
			return 600
		}
	case 314:
		{ /* '314' */
			return 600
		}
	case 315:
		{ /* '315' */
			return 600
		}
	case 316:
		{ /* '316' */
			return 600
		}
	case 317:
		{ /* '317' */
			return 600
		}
	case 318:
		{ /* '318' */
			return 600
		}
	case 319:
		{ /* '319' */
			return 1
		}
	case 32:
		{ /* '32' */
			return 8
		}
	case 320:
		{ /* '320' */
			return 100
		}
	case 321:
		{ /* '321' */
			return 101
		}
	case 33:
		{ /* '33' */
			return 9
		}
	case 34:
		{ /* '34' */
			return 10
		}
	case 35:
		{ /* '35' */
			return 11
		}
	case 36:
		{ /* '36' */
			return 12
		}
	case 37:
		{ /* '37' */
			return 7
		}
	case 38:
		{ /* '38' */
			return 8
		}
	case 39:
		{ /* '39' */
			return 1
		}
	case 4:
		{ /* '4' */
			return 4
		}
	case 40:
		{ /* '40' */
			return 2
		}
	case 41:
		{ /* '41' */
			return 1
		}
	case 42:
		{ /* '42' */
			return 3
		}
	case 43:
		{ /* '43' */
			return 4
		}
	case 44:
		{ /* '44' */
			return 5
		}
	case 45:
		{ /* '45' */
			return 6
		}
	case 46:
		{ /* '46' */
			return 10
		}
	case 47:
		{ /* '47' */
			return 100
		}
	case 48:
		{ /* '48' */
			return 1
		}
	case 49:
		{ /* '49' */
			return 10
		}
	case 5:
		{ /* '5' */
			return 5
		}
	case 50:
		{ /* '50' */
			return 20
		}
	case 51:
		{ /* '51' */
			return 1
		}
	case 52:
		{ /* '52' */
			return 2
		}
	case 53:
		{ /* '53' */
			return 3
		}
	case 54:
		{ /* '54' */
			return 4
		}
	case 55:
		{ /* '55' */
			return 5
		}
	case 56:
		{ /* '56' */
			return 6
		}
	case 57:
		{ /* '57' */
			return 7
		}
	case 58:
		{ /* '58' */
			return 10
		}
	case 59:
		{ /* '59' */
			return 11
		}
	case 6:
		{ /* '6' */
			return 6
		}
	case 60:
		{ /* '60' */
			return 12
		}
	case 61:
		{ /* '61' */
			return 13
		}
	case 62:
		{ /* '62' */
			return 600
		}
	case 63:
		{ /* '63' */
			return 1
		}
	case 64:
		{ /* '64' */
			return 2
		}
	case 65:
		{ /* '65' */
			return 3
		}
	case 66:
		{ /* '66' */
			return 4
		}
	case 67:
		{ /* '67' */
			return 5
		}
	case 68:
		{ /* '68' */
			return 6
		}
	case 69:
		{ /* '69' */
			return 7
		}
	case 7:
		{ /* '7' */
			return 7
		}
	case 70:
		{ /* '70' */
			return 10
		}
	case 71:
		{ /* '71' */
			return 11
		}
	case 72:
		{ /* '72' */
			return 12
		}
	case 73:
		{ /* '73' */
			return 1
		}
	case 74:
		{ /* '74' */
			return 2
		}
	case 75:
		{ /* '75' */
			return 3
		}
	case 76:
		{ /* '76' */
			return 4
		}
	case 77:
		{ /* '77' */
			return 5
		}
	case 78:
		{ /* '78' */
			return 6
		}
	case 79:
		{ /* '79' */
			return 7
		}
	case 8:
		{ /* '8' */
			return 8
		}
	case 80:
		{ /* '80' */
			return 8
		}
	case 81:
		{ /* '81' */
			return 9
		}
	case 82:
		{ /* '82' */
			return 10
		}
	case 83:
		{ /* '83' */
			return 11
		}
	case 84:
		{ /* '84' */
			return 20
		}
	case 85:
		{ /* '85' */
			return 21
		}
	case 86:
		{ /* '86' */
			return 22
		}
	case 87:
		{ /* '87' */
			return 23
		}
	case 88:
		{ /* '88' */
			return 24
		}
	case 89:
		{ /* '89' */
			return 25
		}
	case 9:
		{ /* '9' */
			return 9
		}
	case 90:
		{ /* '90' */
			return 26
		}
	case 91:
		{ /* '91' */
			return 27
		}
	case 92:
		{ /* '92' */
			return 28
		}
	case 93:
		{ /* '93' */
			return 29
		}
	case 94:
		{ /* '94' */
			return 30
		}
	case 95:
		{ /* '95' */
			return 1
		}
	case 96:
		{ /* '96' */
			return 1
		}
	case 97:
		{ /* '97' */
			return 1
		}
	case 98:
		{ /* '98' */
			return 100
		}
	case 99:
		{ /* '99' */
			return 101
		}
	default:
		{
			return 0
		}
	}
}

func (e KnxDatapointSubtype) DatapointType() KnxDatapointType {
	switch e {
	case 0:
		{ /* '0' */
			return KnxDatapointType_DPT_UNKNOWN
		}
	case 1:
		{ /* '1' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10:
		{ /* '10' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 100:
		{ /* '100' */
			return KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE
		}
	case 101:
		{ /* '101' */
			return KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE
		}
	case 102:
		{ /* '102' */
			return KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE
		}
	case 103:
		{ /* '103' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 104:
		{ /* '104' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 105:
		{ /* '105' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 106:
		{ /* '106' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 107:
		{ /* '107' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 108:
		{ /* '108' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 109:
		{ /* '109' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 11:
		{ /* '11' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 110:
		{ /* '110' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 111:
		{ /* '111' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 112:
		{ /* '112' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 113:
		{ /* '113' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 114:
		{ /* '114' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 115:
		{ /* '115' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 116:
		{ /* '116' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 117:
		{ /* '117' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 118:
		{ /* '118' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 119:
		{ /* '119' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 12:
		{ /* '12' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 120:
		{ /* '120' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 121:
		{ /* '121' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 122:
		{ /* '122' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 123:
		{ /* '123' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 124:
		{ /* '124' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 125:
		{ /* '125' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 126:
		{ /* '126' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 127:
		{ /* '127' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 128:
		{ /* '128' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 129:
		{ /* '129' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 13:
		{ /* '13' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 130:
		{ /* '130' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 131:
		{ /* '131' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 132:
		{ /* '132' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 133:
		{ /* '133' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 134:
		{ /* '134' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 135:
		{ /* '135' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 136:
		{ /* '136' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 137:
		{ /* '137' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 138:
		{ /* '138' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 139:
		{ /* '139' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 14:
		{ /* '14' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 140:
		{ /* '140' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 141:
		{ /* '141' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 142:
		{ /* '142' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 143:
		{ /* '143' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 144:
		{ /* '144' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 145:
		{ /* '145' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 146:
		{ /* '146' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 147:
		{ /* '147' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 148:
		{ /* '148' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 149:
		{ /* '149' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 15:
		{ /* '15' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 150:
		{ /* '150' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 151:
		{ /* '151' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 152:
		{ /* '152' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 153:
		{ /* '153' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 154:
		{ /* '154' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 155:
		{ /* '155' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 156:
		{ /* '156' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 157:
		{ /* '157' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 158:
		{ /* '158' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 159:
		{ /* '159' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 16:
		{ /* '16' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 160:
		{ /* '160' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 161:
		{ /* '161' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 162:
		{ /* '162' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 163:
		{ /* '163' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 164:
		{ /* '164' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 165:
		{ /* '165' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 166:
		{ /* '166' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 167:
		{ /* '167' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 168:
		{ /* '168' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 169:
		{ /* '169' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 17:
		{ /* '17' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 170:
		{ /* '170' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 171:
		{ /* '171' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 172:
		{ /* '172' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 173:
		{ /* '173' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 174:
		{ /* '174' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 175:
		{ /* '175' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 176:
		{ /* '176' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 177:
		{ /* '177' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 178:
		{ /* '178' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 179:
		{ /* '179' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 18:
		{ /* '18' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 180:
		{ /* '180' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 181:
		{ /* '181' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 182:
		{ /* '182' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 183:
		{ /* '183' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 184:
		{ /* '184' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 185:
		{ /* '185' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 186:
		{ /* '186' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 187:
		{ /* '187' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 188:
		{ /* '188' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 189:
		{ /* '189' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 19:
		{ /* '19' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 190:
		{ /* '190' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 191:
		{ /* '191' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 192:
		{ /* '192' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 193:
		{ /* '193' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 194:
		{ /* '194' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 195:
		{ /* '195' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 196:
		{ /* '196' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 197:
		{ /* '197' */
			return KnxDatapointType_DPT_ENTRANCE_ACCESS
		}
	case 198:
		{ /* '198' */
			return KnxDatapointType_DPT_CHARACTER_STRING
		}
	case 199:
		{ /* '199' */
			return KnxDatapointType_DPT_CHARACTER_STRING
		}
	case 2:
		{ /* '2' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 20:
		{ /* '20' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 200:
		{ /* '200' */
			return KnxDatapointType_DPT_SCENE_NUMBER
		}
	case 201:
		{ /* '201' */
			return KnxDatapointType_DPT_SCENE_CONTROL
		}
	case 202:
		{ /* '202' */
			return KnxDatapointType_DPT_DATE_TIME
		}
	case 203:
		{ /* '203' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 204:
		{ /* '204' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 205:
		{ /* '205' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 206:
		{ /* '206' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 207:
		{ /* '207' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 208:
		{ /* '208' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 209:
		{ /* '209' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 21:
		{ /* '21' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 210:
		{ /* '210' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 211:
		{ /* '211' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 212:
		{ /* '212' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 213:
		{ /* '213' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 214:
		{ /* '214' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 215:
		{ /* '215' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 216:
		{ /* '216' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 217:
		{ /* '217' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 218:
		{ /* '218' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 219:
		{ /* '219' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 22:
		{ /* '22' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 220:
		{ /* '220' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 221:
		{ /* '221' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 222:
		{ /* '222' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 223:
		{ /* '223' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 224:
		{ /* '224' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 225:
		{ /* '225' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 226:
		{ /* '226' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 227:
		{ /* '227' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 228:
		{ /* '228' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 229:
		{ /* '229' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 23:
		{ /* '23' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 230:
		{ /* '230' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 231:
		{ /* '231' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 232:
		{ /* '232' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 233:
		{ /* '233' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 234:
		{ /* '234' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 235:
		{ /* '235' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 236:
		{ /* '236' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 237:
		{ /* '237' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 238:
		{ /* '238' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 239:
		{ /* '239' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 24:
		{ /* '24' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 240:
		{ /* '240' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 241:
		{ /* '241' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 242:
		{ /* '242' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 243:
		{ /* '243' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 244:
		{ /* '244' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 245:
		{ /* '245' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 246:
		{ /* '246' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 247:
		{ /* '247' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 248:
		{ /* '248' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 249:
		{ /* '249' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 25:
		{ /* '25' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 250:
		{ /* '250' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 251:
		{ /* '251' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 252:
		{ /* '252' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 253:
		{ /* '253' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 254:
		{ /* '254' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 255:
		{ /* '255' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 256:
		{ /* '256' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 257:
		{ /* '257' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 258:
		{ /* '258' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 259:
		{ /* '259' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 26:
		{ /* '26' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 260:
		{ /* '260' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 261:
		{ /* '261' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 262:
		{ /* '262' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 263:
		{ /* '263' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 264:
		{ /* '264' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 265:
		{ /* '265' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 266:
		{ /* '266' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 267:
		{ /* '267' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 268:
		{ /* '268' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 269:
		{ /* '269' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 27:
		{ /* '27' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 270:
		{ /* '270' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 271:
		{ /* '271' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 272:
		{ /* '272' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 273:
		{ /* '273' */
			return KnxDatapointType_DPT_16_BIT_SET
		}
	case 274:
		{ /* '274' */
			return KnxDatapointType_DPT_16_BIT_SET
		}
	case 275:
		{ /* '275' */
			return KnxDatapointType_DPT_16_BIT_SET
		}
	case 276:
		{ /* '276' */
			return KnxDatapointType_DPT_16_BIT_SET
		}
	case 277:
		{ /* '277' */
			return KnxDatapointType_DPT_16_BIT_SET
		}
	case 278:
		{ /* '278' */
			return KnxDatapointType_DPT_16_BIT_SET
		}
	case 279:
		{ /* '279' */
			return KnxDatapointType_DPT_2_BIT_SET
		}
	case 28:
		{ /* '28' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 280:
		{ /* '280' */
			return KnxDatapointType_DPT_2_BIT_SET
		}
	case 281:
		{ /* '281' */
			return KnxDatapointType_DPT_2_BIT_SET
		}
	case 282:
		{ /* '282' */
			return KnxDatapointType_DPT_2_BIT_SET
		}
	case 283:
		{ /* '283' */
			return KnxDatapointType_DPT_2_NIBBLE_SET
		}
	case 284:
		{ /* '284' */
			return KnxDatapointType_DPT_8_BIT_SET_2
		}
	case 285:
		{ /* '285' */
			return KnxDatapointType_DPT_32_BIT_SET
		}
	case 286:
		{ /* '286' */
			return KnxDatapointType_DPT_ELECTRICAL_ENERGY
		}
	case 287:
		{ /* '287' */
			return KnxDatapointType_DPT_ELECTRICAL_ENERGY
		}
	case 288:
		{ /* '288' */
			return KnxDatapointType_DPT_ELECTRICAL_ENERGY
		}
	case 289:
		{ /* '289' */
			return KnxDatapointType_DPT_24_TIMES_CHANNEL_ACTIVATION
		}
	case 29:
		{ /* '29' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 290:
		{ /* '290' */
			return KnxDatapointType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM
		}
	case 291:
		{ /* '291' */
			return KnxDatapointType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM
		}
	case 292:
		{ /* '292' */
			return KnxDatapointType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM
		}
	case 293:
		{ /* '293' */
			return KnxDatapointType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM
		}
	case 294:
		{ /* '294' */
			return KnxDatapointType_DPT_DATAPOINT_TYPE_VERSION
		}
	case 295:
		{ /* '295' */
			return KnxDatapointType_DPT_ALARM_INFO
		}
	case 296:
		{ /* '296' */
			return KnxDatapointType_DPT_3X_2_BYTE_FLOAT_VALUE
		}
	case 297:
		{ /* '297' */
			return KnxDatapointType_DPT_3X_2_BYTE_FLOAT_VALUE
		}
	case 298:
		{ /* '298' */
			return KnxDatapointType_DPT_SCALING_SPEED
		}
	case 299:
		{ /* '299' */
			return KnxDatapointType_DPT_SCALING_SPEED
		}
	case 3:
		{ /* '3' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 30:
		{ /* '30' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 300:
		{ /* '300' */
			return KnxDatapointType_DPT_4_1_1_BYTE_COMBINED_INFORMATION
		}
	case 301:
		{ /* '301' */
			return KnxDatapointType_DPT_MBUS_ADDRESS
		}
	case 302:
		{ /* '302' */
			return KnxDatapointType_DPT_3_BYTE_COLOUR_RGB
		}
	case 303:
		{ /* '303' */
			return KnxDatapointType_DPT_LANGUAGE_CODE_ISO_639_1
		}
	case 304:
		{ /* '304' */
			return KnxDatapointType_DPT_SIGNED_VALUE_WITH_CLASSIFICATION_AND_VALIDITY
		}
	case 305:
		{ /* '305' */
			return KnxDatapointType_DPT_PRIORITISED_MODE_CONTROL
		}
	case 306:
		{ /* '306' */
			return KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_16_BIT
		}
	case 307:
		{ /* '307' */
			return KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_8_BIT
		}
	case 308:
		{ /* '308' */
			return KnxDatapointType_DPT_POSITIONS
		}
	case 309:
		{ /* '309' */
			return KnxDatapointType_DPT_STATUS_32_BIT
		}
	case 31:
		{ /* '31' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 310:
		{ /* '310' */
			return KnxDatapointType_DPT_STATUS_48_BIT
		}
	case 311:
		{ /* '311' */
			return KnxDatapointType_DPT_CONVERTER_STATUS
		}
	case 312:
		{ /* '312' */
			return KnxDatapointType_DPT_CONVERTER_TEST_RESULT
		}
	case 313:
		{ /* '313' */
			return KnxDatapointType_DPT_BATTERY_INFORMATION
		}
	case 314:
		{ /* '314' */
			return KnxDatapointType_DPT_BRIGHTNESS_COLOUR_TEMPERATURE_TRANSITION
		}
	case 315:
		{ /* '315' */
			return KnxDatapointType_DPT_STATUS_24_BIT
		}
	case 316:
		{ /* '316' */
			return KnxDatapointType_DPT_COLOUR_RGBW
		}
	case 317:
		{ /* '317' */
			return KnxDatapointType_DPT_RELATIVE_CONTROL_RGBW
		}
	case 318:
		{ /* '318' */
			return KnxDatapointType_DPT_RELATIVE_CONTROL_RGB
		}
	case 319:
		{ /* '319' */
			return KnxDatapointType_DPT_F32F32
		}
	case 32:
		{ /* '32' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 320:
		{ /* '320' */
			return KnxDatapointType_DPT_F16F16F16F16
		}
	case 321:
		{ /* '321' */
			return KnxDatapointType_DPT_F16F16F16F16
		}
	case 33:
		{ /* '33' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 34:
		{ /* '34' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 35:
		{ /* '35' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 36:
		{ /* '36' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 37:
		{ /* '37' */
			return KnxDatapointType_DPT_3_BIT_CONTROLLED
		}
	case 38:
		{ /* '38' */
			return KnxDatapointType_DPT_3_BIT_CONTROLLED
		}
	case 39:
		{ /* '39' */
			return KnxDatapointType_DPT_CHARACTER
		}
	case 4:
		{ /* '4' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 40:
		{ /* '40' */
			return KnxDatapointType_DPT_CHARACTER
		}
	case 41:
		{ /* '41' */
			return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
		}
	case 42:
		{ /* '42' */
			return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
		}
	case 43:
		{ /* '43' */
			return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
		}
	case 44:
		{ /* '44' */
			return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
		}
	case 45:
		{ /* '45' */
			return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
		}
	case 46:
		{ /* '46' */
			return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
		}
	case 47:
		{ /* '47' */
			return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
		}
	case 48:
		{ /* '48' */
			return KnxDatapointType_DPT_8_BIT_SIGNED_VALUE
		}
	case 49:
		{ /* '49' */
			return KnxDatapointType_DPT_8_BIT_SIGNED_VALUE
		}
	case 5:
		{ /* '5' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 50:
		{ /* '50' */
			return KnxDatapointType_DPT_8_BIT_SIGNED_VALUE
		}
	case 51:
		{ /* '51' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 52:
		{ /* '52' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 53:
		{ /* '53' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 54:
		{ /* '54' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 55:
		{ /* '55' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 56:
		{ /* '56' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 57:
		{ /* '57' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 58:
		{ /* '58' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 59:
		{ /* '59' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 6:
		{ /* '6' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 60:
		{ /* '60' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 61:
		{ /* '61' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 62:
		{ /* '62' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 63:
		{ /* '63' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 64:
		{ /* '64' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 65:
		{ /* '65' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 66:
		{ /* '66' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 67:
		{ /* '67' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 68:
		{ /* '68' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 69:
		{ /* '69' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 7:
		{ /* '7' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 70:
		{ /* '70' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 71:
		{ /* '71' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 72:
		{ /* '72' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 73:
		{ /* '73' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 74:
		{ /* '74' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 75:
		{ /* '75' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 76:
		{ /* '76' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 77:
		{ /* '77' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 78:
		{ /* '78' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 79:
		{ /* '79' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 8:
		{ /* '8' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 80:
		{ /* '80' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 81:
		{ /* '81' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 82:
		{ /* '82' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 83:
		{ /* '83' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 84:
		{ /* '84' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 85:
		{ /* '85' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 86:
		{ /* '86' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 87:
		{ /* '87' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 88:
		{ /* '88' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 89:
		{ /* '89' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 9:
		{ /* '9' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 90:
		{ /* '90' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 91:
		{ /* '91' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 92:
		{ /* '92' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 93:
		{ /* '93' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 94:
		{ /* '94' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 95:
		{ /* '95' */
			return KnxDatapointType_DPT_TIME
		}
	case 96:
		{ /* '96' */
			return KnxDatapointType_DPT_DATE
		}
	case 97:
		{ /* '97' */
			return KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE
		}
	case 98:
		{ /* '98' */
			return KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE
		}
	case 99:
		{ /* '99' */
			return KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE
		}
	default:
		{
			return 0
		}
	}
}

func (e KnxDatapointSubtype) Name() string {
	switch e {
	case 0:
		{ /* '0' */
			return "Unknown Datapoint Subtype"
		}
	case 1:
		{ /* '1' */
			return "switch"
		}
	case 10:
		{ /* '10' */
			return "start/stop"
		}
	case 100:
		{ /* '100' */
			return "counter timehrs (h)"
		}
	case 101:
		{ /* '101' */
			return "volume liquid (l)"
		}
	case 102:
		{ /* '102' */
			return "volume (m³)"
		}
	case 103:
		{ /* '103' */
			return "counter pulses (signed)"
		}
	case 104:
		{ /* '104' */
			return "flow rate (m³/h)"
		}
	case 105:
		{ /* '105' */
			return "active energy (Wh)"
		}
	case 106:
		{ /* '106' */
			return "apparant energy (VAh)"
		}
	case 107:
		{ /* '107' */
			return "reactive energy (VARh)"
		}
	case 108:
		{ /* '108' */
			return "active energy (kWh)"
		}
	case 109:
		{ /* '109' */
			return "apparant energy (kVAh)"
		}
	case 11:
		{ /* '11' */
			return "state"
		}
	case 110:
		{ /* '110' */
			return "reactive energy (kVARh)"
		}
	case 111:
		{ /* '111' */
			return "active energy (MWh)"
		}
	case 112:
		{ /* '112' */
			return "time lag (s)"
		}
	case 113:
		{ /* '113' */
			return "delta volume liquid (l)"
		}
	case 114:
		{ /* '114' */
			return "delta volume (m³)"
		}
	case 115:
		{ /* '115' */
			return "acceleration (m/s²)"
		}
	case 116:
		{ /* '116' */
			return "angular acceleration (rad/s²)"
		}
	case 117:
		{ /* '117' */
			return "activation energy (J/mol)"
		}
	case 118:
		{ /* '118' */
			return "radioactive activity (1/s)"
		}
	case 119:
		{ /* '119' */
			return "amount of substance (mol)"
		}
	case 12:
		{ /* '12' */
			return "invert"
		}
	case 120:
		{ /* '120' */
			return "amplitude"
		}
	case 121:
		{ /* '121' */
			return "angle (radiant)"
		}
	case 122:
		{ /* '122' */
			return "angle (degree)"
		}
	case 123:
		{ /* '123' */
			return "angular momentum (Js)"
		}
	case 124:
		{ /* '124' */
			return "angular velocity (rad/s)"
		}
	case 125:
		{ /* '125' */
			return "area (m*m)"
		}
	case 126:
		{ /* '126' */
			return "capacitance (F)"
		}
	case 127:
		{ /* '127' */
			return "flux density (C/m²)"
		}
	case 128:
		{ /* '128' */
			return "charge density (C/m³)"
		}
	case 129:
		{ /* '129' */
			return "compressibility (m²/N)"
		}
	case 13:
		{ /* '13' */
			return "dim send style"
		}
	case 130:
		{ /* '130' */
			return "conductance (S)"
		}
	case 131:
		{ /* '131' */
			return "conductivity  (S/m)"
		}
	case 132:
		{ /* '132' */
			return "density (kg/m³)"
		}
	case 133:
		{ /* '133' */
			return "electric charge (C)"
		}
	case 134:
		{ /* '134' */
			return "electric current (A)"
		}
	case 135:
		{ /* '135' */
			return "electric current density (A/m²)"
		}
	case 136:
		{ /* '136' */
			return "electric dipole moment (Cm)"
		}
	case 137:
		{ /* '137' */
			return "electric displacement (C/m²)"
		}
	case 138:
		{ /* '138' */
			return "electric field strength (V/m)"
		}
	case 139:
		{ /* '139' */
			return "electric flux (C)"
		}
	case 14:
		{ /* '14' */
			return "input source"
		}
	case 140:
		{ /* '140' */
			return "electric flux density (C/m²)"
		}
	case 141:
		{ /* '141' */
			return "electric polarization (C/m²)"
		}
	case 142:
		{ /* '142' */
			return "electric potential (V)"
		}
	case 143:
		{ /* '143' */
			return "electric potential difference (V)"
		}
	case 144:
		{ /* '144' */
			return "electromagnetic moment (Am²)"
		}
	case 145:
		{ /* '145' */
			return "electromotive force (V)"
		}
	case 146:
		{ /* '146' */
			return "energy (J)"
		}
	case 147:
		{ /* '147' */
			return "force (N)"
		}
	case 148:
		{ /* '148' */
			return "frequency (Hz)"
		}
	case 149:
		{ /* '149' */
			return "angular frequency (rad/s)"
		}
	case 15:
		{ /* '15' */
			return "reset"
		}
	case 150:
		{ /* '150' */
			return "heat capacity (J/K)"
		}
	case 151:
		{ /* '151' */
			return "heat flow rate (W)"
		}
	case 152:
		{ /* '152' */
			return "heat quantity"
		}
	case 153:
		{ /* '153' */
			return "impedance (Ω)"
		}
	case 154:
		{ /* '154' */
			return "length (m)"
		}
	case 155:
		{ /* '155' */
			return "light quantity (J)"
		}
	case 156:
		{ /* '156' */
			return "luminance (cd/m²)"
		}
	case 157:
		{ /* '157' */
			return "luminous flux (lm)"
		}
	case 158:
		{ /* '158' */
			return "luminous intensity (cd)"
		}
	case 159:
		{ /* '159' */
			return "magnetic field strength (A/m)"
		}
	case 16:
		{ /* '16' */
			return "acknowledge"
		}
	case 160:
		{ /* '160' */
			return "magnetic flux (Wb)"
		}
	case 161:
		{ /* '161' */
			return "magnetic flux density (T)"
		}
	case 162:
		{ /* '162' */
			return "magnetic moment (Am²)"
		}
	case 163:
		{ /* '163' */
			return "magnetic polarization (T)"
		}
	case 164:
		{ /* '164' */
			return "magnetization (A/m)"
		}
	case 165:
		{ /* '165' */
			return "magnetomotive force (A)"
		}
	case 166:
		{ /* '166' */
			return "mass (kg)"
		}
	case 167:
		{ /* '167' */
			return "mass flux (kg/s)"
		}
	case 168:
		{ /* '168' */
			return "momentum (N/s)"
		}
	case 169:
		{ /* '169' */
			return "phase angle (rad)"
		}
	case 17:
		{ /* '17' */
			return "trigger"
		}
	case 170:
		{ /* '170' */
			return "phase angle (°)"
		}
	case 171:
		{ /* '171' */
			return "power (W)"
		}
	case 172:
		{ /* '172' */
			return "power factor (cos Φ)"
		}
	case 173:
		{ /* '173' */
			return "pressure (Pa)"
		}
	case 174:
		{ /* '174' */
			return "reactance (Ω)"
		}
	case 175:
		{ /* '175' */
			return "resistance (Ω)"
		}
	case 176:
		{ /* '176' */
			return "resistivity (Ωm)"
		}
	case 177:
		{ /* '177' */
			return "self inductance (H)"
		}
	case 178:
		{ /* '178' */
			return "solid angle (sr)"
		}
	case 179:
		{ /* '179' */
			return "sound intensity (W/m²)"
		}
	case 18:
		{ /* '18' */
			return "occupancy"
		}
	case 180:
		{ /* '180' */
			return "speed (m/s)"
		}
	case 181:
		{ /* '181' */
			return "stress (Pa)"
		}
	case 182:
		{ /* '182' */
			return "surface tension (N/m)"
		}
	case 183:
		{ /* '183' */
			return "temperature (°C)"
		}
	case 184:
		{ /* '184' */
			return "temperature absolute (K)"
		}
	case 185:
		{ /* '185' */
			return "temperature difference (K)"
		}
	case 186:
		{ /* '186' */
			return "thermal capacity (J/K)"
		}
	case 187:
		{ /* '187' */
			return "thermal conductivity (W/mK)"
		}
	case 188:
		{ /* '188' */
			return "thermoelectric power (V/K)"
		}
	case 189:
		{ /* '189' */
			return "time (s)"
		}
	case 19:
		{ /* '19' */
			return "window/door"
		}
	case 190:
		{ /* '190' */
			return "torque (Nm)"
		}
	case 191:
		{ /* '191' */
			return "volume (m³)"
		}
	case 192:
		{ /* '192' */
			return "volume flux (m³/s)"
		}
	case 193:
		{ /* '193' */
			return "weight (N)"
		}
	case 194:
		{ /* '194' */
			return "work (J)"
		}
	case 195:
		{ /* '195' */
			return "volume flux for meters (m³/h)"
		}
	case 196:
		{ /* '196' */
			return "volume flux for meters (1/ls)"
		}
	case 197:
		{ /* '197' */
			return "entrance access"
		}
	case 198:
		{ /* '198' */
			return "Character String (ASCII)"
		}
	case 199:
		{ /* '199' */
			return "Character String (ISO 8859-1)"
		}
	case 2:
		{ /* '2' */
			return "boolean"
		}
	case 20:
		{ /* '20' */
			return "logical function"
		}
	case 200:
		{ /* '200' */
			return "scene number"
		}
	case 201:
		{ /* '201' */
			return "scene control"
		}
	case 202:
		{ /* '202' */
			return "date time"
		}
	case 203:
		{ /* '203' */
			return "SCLO mode"
		}
	case 204:
		{ /* '204' */
			return "building mode"
		}
	case 205:
		{ /* '205' */
			return "occupied"
		}
	case 206:
		{ /* '206' */
			return "priority"
		}
	case 207:
		{ /* '207' */
			return "light application mode"
		}
	case 208:
		{ /* '208' */
			return "light application area"
		}
	case 209:
		{ /* '209' */
			return "alarm class"
		}
	case 21:
		{ /* '21' */
			return "scene"
		}
	case 210:
		{ /* '210' */
			return "PSU mode"
		}
	case 211:
		{ /* '211' */
			return "system error class"
		}
	case 212:
		{ /* '212' */
			return "HVAC error class"
		}
	case 213:
		{ /* '213' */
			return "time delay"
		}
	case 214:
		{ /* '214' */
			return "wind force scale (0..12)"
		}
	case 215:
		{ /* '215' */
			return "sensor mode"
		}
	case 216:
		{ /* '216' */
			return "actuator connect type"
		}
	case 217:
		{ /* '217' */
			return "cloud cover"
		}
	case 218:
		{ /* '218' */
			return "power return mode"
		}
	case 219:
		{ /* '219' */
			return "fuel type"
		}
	case 22:
		{ /* '22' */
			return "shutter/blinds mode"
		}
	case 220:
		{ /* '220' */
			return "burner type"
		}
	case 221:
		{ /* '221' */
			return "HVAC mode"
		}
	case 222:
		{ /* '222' */
			return "DHW mode"
		}
	case 223:
		{ /* '223' */
			return "load priority"
		}
	case 224:
		{ /* '224' */
			return "HVAC control mode"
		}
	case 225:
		{ /* '225' */
			return "HVAC emergency mode"
		}
	case 226:
		{ /* '226' */
			return "changeover mode"
		}
	case 227:
		{ /* '227' */
			return "valve mode"
		}
	case 228:
		{ /* '228' */
			return "damper mode"
		}
	case 229:
		{ /* '229' */
			return "heater mode"
		}
	case 23:
		{ /* '23' */
			return "day/night"
		}
	case 230:
		{ /* '230' */
			return "fan mode"
		}
	case 231:
		{ /* '231' */
			return "master/slave mode"
		}
	case 232:
		{ /* '232' */
			return "status room setpoint"
		}
	case 233:
		{ /* '233' */
			return "metering device type"
		}
	case 234:
		{ /* '234' */
			return "hum dehum mode"
		}
	case 235:
		{ /* '235' */
			return "enable H/C stage"
		}
	case 236:
		{ /* '236' */
			return "ADA type"
		}
	case 237:
		{ /* '237' */
			return "backup mode"
		}
	case 238:
		{ /* '238' */
			return "start syncronization type"
		}
	case 239:
		{ /* '239' */
			return "behavior lock/unlock"
		}
	case 24:
		{ /* '24' */
			return "cooling/heating"
		}
	case 240:
		{ /* '240' */
			return "behavior bus power up/down"
		}
	case 241:
		{ /* '241' */
			return "dali fade time"
		}
	case 242:
		{ /* '242' */
			return "blink mode"
		}
	case 243:
		{ /* '243' */
			return "light control mode"
		}
	case 244:
		{ /* '244' */
			return "PB switch mode"
		}
	case 245:
		{ /* '245' */
			return "PB action mode"
		}
	case 246:
		{ /* '246' */
			return "PB dimm mode"
		}
	case 247:
		{ /* '247' */
			return "switch on mode"
		}
	case 248:
		{ /* '248' */
			return "load type"
		}
	case 249:
		{ /* '249' */
			return "load type detection"
		}
	case 25:
		{ /* '25' */
			return "switch control"
		}
	case 250:
		{ /* '250' */
			return "converter test control"
		}
	case 251:
		{ /* '251' */
			return "SAB except behavior"
		}
	case 252:
		{ /* '252' */
			return "SAB behavior on lock/unlock"
		}
	case 253:
		{ /* '253' */
			return "SSSB mode"
		}
	case 254:
		{ /* '254' */
			return "blinds control mode"
		}
	case 255:
		{ /* '255' */
			return "communication mode"
		}
	case 256:
		{ /* '256' */
			return "additional information type"
		}
	case 257:
		{ /* '257' */
			return "RF mode selection"
		}
	case 258:
		{ /* '258' */
			return "RF filter mode selection"
		}
	case 259:
		{ /* '259' */
			return "general status"
		}
	case 26:
		{ /* '26' */
			return "boolean control"
		}
	case 260:
		{ /* '260' */
			return "device control"
		}
	case 261:
		{ /* '261' */
			return "forcing signal"
		}
	case 262:
		{ /* '262' */
			return "forcing signal cool"
		}
	case 263:
		{ /* '263' */
			return "room heating controller status"
		}
	case 264:
		{ /* '264' */
			return "solar DHW controller status"
		}
	case 265:
		{ /* '265' */
			return "fuel type set"
		}
	case 266:
		{ /* '266' */
			return "room cooling controller status"
		}
	case 267:
		{ /* '267' */
			return "ventilation controller status"
		}
	case 268:
		{ /* '268' */
			return "combined status RTSM"
		}
	case 269:
		{ /* '269' */
			return "lighting actuator error information"
		}
	case 27:
		{ /* '27' */
			return "enable control"
		}
	case 270:
		{ /* '270' */
			return "RF communication mode info"
		}
	case 271:
		{ /* '271' */
			return "cEMI server supported RF filtering modes"
		}
	case 272:
		{ /* '272' */
			return "channel activation for 8 channels"
		}
	case 273:
		{ /* '273' */
			return "DHW controller status"
		}
	case 274:
		{ /* '274' */
			return "RHCC status"
		}
	case 275:
		{ /* '275' */
			return "combined status HVA"
		}
	case 276:
		{ /* '276' */
			return "combined status RTC"
		}
	case 277:
		{ /* '277' */
			return "media"
		}
	case 278:
		{ /* '278' */
			return "channel activation for 16 channels"
		}
	case 279:
		{ /* '279' */
			return "on/off action"
		}
	case 28:
		{ /* '28' */
			return "ramp control"
		}
	case 280:
		{ /* '280' */
			return "alarm reaction"
		}
	case 281:
		{ /* '281' */
			return "up/down action"
		}
	case 282:
		{ /* '282' */
			return "HVAC push button action"
		}
	case 283:
		{ /* '283' */
			return "busy/nak repetitions"
		}
	case 284:
		{ /* '284' */
			return "scene information"
		}
	case 285:
		{ /* '285' */
			return "bit-combined info on/off"
		}
	case 286:
		{ /* '286' */
			return "active energy (Wh)"
		}
	case 287:
		{ /* '287' */
			return "apparant energy (VAh)"
		}
	case 288:
		{ /* '288' */
			return "reactive energy (VARh)"
		}
	case 289:
		{ /* '289' */
			return "activation state 0..23"
		}
	case 29:
		{ /* '29' */
			return "alarm control"
		}
	case 290:
		{ /* '290' */
			return "time delay & HVAC mode"
		}
	case 291:
		{ /* '291' */
			return "time delay & DHW mode"
		}
	case 292:
		{ /* '292' */
			return "time delay & occupancy mode"
		}
	case 293:
		{ /* '293' */
			return "time delay & building mode"
		}
	case 294:
		{ /* '294' */
			return "DPT version"
		}
	case 295:
		{ /* '295' */
			return "alarm info"
		}
	case 296:
		{ /* '296' */
			return "room temperature setpoint"
		}
	case 297:
		{ /* '297' */
			return "room temperature setpoint shift"
		}
	case 298:
		{ /* '298' */
			return "scaling speed"
		}
	case 299:
		{ /* '299' */
			return "scaling step time"
		}
	case 3:
		{ /* '3' */
			return "enable"
		}
	case 30:
		{ /* '30' */
			return "binary value control"
		}
	case 300:
		{ /* '300' */
			return "metering value (value,encoding,cmd)"
		}
	case 301:
		{ /* '301' */
			return "MBus address"
		}
	case 302:
		{ /* '302' */
			return "RGB value 3x(0..255)"
		}
	case 303:
		{ /* '303' */
			return "language code (ASCII)"
		}
	case 304:
		{ /* '304' */
			return "electrical energy with tariff"
		}
	case 305:
		{ /* '305' */
			return "priority control"
		}
	case 306:
		{ /* '306' */
			return "diagnostic value"
		}
	case 307:
		{ /* '307' */
			return "diagnostic value"
		}
	case 308:
		{ /* '308' */
			return "combined position"
		}
	case 309:
		{ /* '309' */
			return "status sunblind & shutter actuator"
		}
	case 31:
		{ /* '31' */
			return "step control"
		}
	case 310:
		{ /* '310' */
			return "colour xyY"
		}
	case 311:
		{ /* '311' */
			return "DALI converter status"
		}
	case 312:
		{ /* '312' */
			return "DALI converter test result"
		}
	case 313:
		{ /* '313' */
			return "Battery Information"
		}
	case 314:
		{ /* '314' */
			return "brightness colour temperature transition"
		}
	case 315:
		{ /* '315' */
			return "brightness colour temperature control"
		}
	case 316:
		{ /* '316' */
			return "RGBW value 4x(0..100%)"
		}
	case 317:
		{ /* '317' */
			return "RGBW relative control"
		}
	case 318:
		{ /* '318' */
			return "RGB relative control"
		}
	case 319:
		{ /* '319' */
			return "geographical location (longitude and latitude) expressed in degrees"
		}
	case 32:
		{ /* '32' */
			return "direction control 1"
		}
	case 320:
		{ /* '320' */
			return "Temperature setpoint setting for 4 HVAC Modes"
		}
	case 321:
		{ /* '321' */
			return "Temperature setpoint shift setting for 4 HVAC Modes"
		}
	case 33:
		{ /* '33' */
			return "direction control 2"
		}
	case 34:
		{ /* '34' */
			return "start control"
		}
	case 35:
		{ /* '35' */
			return "state control"
		}
	case 36:
		{ /* '36' */
			return "invert control"
		}
	case 37:
		{ /* '37' */
			return "dimming control"
		}
	case 38:
		{ /* '38' */
			return "blind control"
		}
	case 39:
		{ /* '39' */
			return "character (ASCII)"
		}
	case 4:
		{ /* '4' */
			return "ramp"
		}
	case 40:
		{ /* '40' */
			return "character (ISO 8859-1)"
		}
	case 41:
		{ /* '41' */
			return "percentage (0..100%)"
		}
	case 42:
		{ /* '42' */
			return "angle (degrees)"
		}
	case 43:
		{ /* '43' */
			return "percentage (0..255%)"
		}
	case 44:
		{ /* '44' */
			return "ratio (0..255)"
		}
	case 45:
		{ /* '45' */
			return "tariff (0..255)"
		}
	case 46:
		{ /* '46' */
			return "counter pulses (0..255)"
		}
	case 47:
		{ /* '47' */
			return "fan stage (0..255)"
		}
	case 48:
		{ /* '48' */
			return "percentage (-128..127%)"
		}
	case 49:
		{ /* '49' */
			return "counter pulses (-128..127)"
		}
	case 5:
		{ /* '5' */
			return "alarm"
		}
	case 50:
		{ /* '50' */
			return "status with mode"
		}
	case 51:
		{ /* '51' */
			return "pulses"
		}
	case 52:
		{ /* '52' */
			return "time (ms)"
		}
	case 53:
		{ /* '53' */
			return "time (10 ms)"
		}
	case 54:
		{ /* '54' */
			return "time (100 ms)"
		}
	case 55:
		{ /* '55' */
			return "time (s)"
		}
	case 56:
		{ /* '56' */
			return "time (min)"
		}
	case 57:
		{ /* '57' */
			return "time (h)"
		}
	case 58:
		{ /* '58' */
			return "property data type"
		}
	case 59:
		{ /* '59' */
			return "length (mm)"
		}
	case 6:
		{ /* '6' */
			return "binary value"
		}
	case 60:
		{ /* '60' */
			return "current (mA)"
		}
	case 61:
		{ /* '61' */
			return "brightness (lux)"
		}
	case 62:
		{ /* '62' */
			return "absolute colour temperature (K)"
		}
	case 63:
		{ /* '63' */
			return "pulses difference"
		}
	case 64:
		{ /* '64' */
			return "time lag (ms)"
		}
	case 65:
		{ /* '65' */
			return "time lag(10 ms)"
		}
	case 66:
		{ /* '66' */
			return "time lag (100 ms)"
		}
	case 67:
		{ /* '67' */
			return "time lag (s)"
		}
	case 68:
		{ /* '68' */
			return "time lag (min)"
		}
	case 69:
		{ /* '69' */
			return "time lag (h)"
		}
	case 7:
		{ /* '7' */
			return "step"
		}
	case 70:
		{ /* '70' */
			return "percentage difference (%)"
		}
	case 71:
		{ /* '71' */
			return "rotation angle (°)"
		}
	case 72:
		{ /* '72' */
			return "length (m)"
		}
	case 73:
		{ /* '73' */
			return "temperature (°C)"
		}
	case 74:
		{ /* '74' */
			return "temperature difference (K)"
		}
	case 75:
		{ /* '75' */
			return "kelvin/hour (K/h)"
		}
	case 76:
		{ /* '76' */
			return "lux (Lux)"
		}
	case 77:
		{ /* '77' */
			return "speed (m/s)"
		}
	case 78:
		{ /* '78' */
			return "pressure (Pa)"
		}
	case 79:
		{ /* '79' */
			return "humidity (%)"
		}
	case 8:
		{ /* '8' */
			return "up/down"
		}
	case 80:
		{ /* '80' */
			return "parts/million (ppm)"
		}
	case 81:
		{ /* '81' */
			return "air flow (m³/h)"
		}
	case 82:
		{ /* '82' */
			return "time (s)"
		}
	case 83:
		{ /* '83' */
			return "time (ms)"
		}
	case 84:
		{ /* '84' */
			return "voltage (mV)"
		}
	case 85:
		{ /* '85' */
			return "current (mA)"
		}
	case 86:
		{ /* '86' */
			return "power denisity (W/m²)"
		}
	case 87:
		{ /* '87' */
			return "kelvin/percent (K/%)"
		}
	case 88:
		{ /* '88' */
			return "power (kW)"
		}
	case 89:
		{ /* '89' */
			return "volume flow (l/h)"
		}
	case 9:
		{ /* '9' */
			return "open/close"
		}
	case 90:
		{ /* '90' */
			return "rain amount (l/m²)"
		}
	case 91:
		{ /* '91' */
			return "temperature (°F)"
		}
	case 92:
		{ /* '92' */
			return "wind speed (km/h)"
		}
	case 93:
		{ /* '93' */
			return "absolute humidity (g/m³)"
		}
	case 94:
		{ /* '94' */
			return "concentration (µg/m³)"
		}
	case 95:
		{ /* '95' */
			return "time of day"
		}
	case 96:
		{ /* '96' */
			return "date"
		}
	case 97:
		{ /* '97' */
			return "counter pulses (unsigned)"
		}
	case 98:
		{ /* '98' */
			return "counter timesec (s)"
		}
	case 99:
		{ /* '99' */
			return "counter timemin (min)"
		}
	default:
		{
			return ""
		}
	}
}
func KnxDatapointSubtypeByValue(value uint32) KnxDatapointSubtype {
	switch value {
	case 0:
		return KnxDatapointSubtype_DPST_UNKNOWN
	case 1:
		return KnxDatapointSubtype_DPST_Switch
	case 10:
		return KnxDatapointSubtype_DPST_Start
	case 100:
		return KnxDatapointSubtype_DPST_LongTimePeriod_Hrs
	case 101:
		return KnxDatapointSubtype_DPST_VolumeLiquid_Litre
	case 102:
		return KnxDatapointSubtype_DPST_Volume_m_3
	case 103:
		return KnxDatapointSubtype_DPST_Value_4_Count
	case 104:
		return KnxDatapointSubtype_DPST_FlowRate_m3h
	case 105:
		return KnxDatapointSubtype_DPST_ActiveEnergy
	case 106:
		return KnxDatapointSubtype_DPST_ApparantEnergy
	case 107:
		return KnxDatapointSubtype_DPST_ReactiveEnergy
	case 108:
		return KnxDatapointSubtype_DPST_ActiveEnergy_kWh
	case 109:
		return KnxDatapointSubtype_DPST_ApparantEnergy_kVAh
	case 11:
		return KnxDatapointSubtype_DPST_State
	case 110:
		return KnxDatapointSubtype_DPST_ReactiveEnergy_kVARh
	case 111:
		return KnxDatapointSubtype_DPST_ActiveEnergy_MWh
	case 112:
		return KnxDatapointSubtype_DPST_LongDeltaTimeSec
	case 113:
		return KnxDatapointSubtype_DPST_DeltaVolumeLiquid_Litre
	case 114:
		return KnxDatapointSubtype_DPST_DeltaVolume_m_3
	case 115:
		return KnxDatapointSubtype_DPST_Value_Acceleration
	case 116:
		return KnxDatapointSubtype_DPST_Value_Acceleration_Angular
	case 117:
		return KnxDatapointSubtype_DPST_Value_Activation_Energy
	case 118:
		return KnxDatapointSubtype_DPST_Value_Activity
	case 119:
		return KnxDatapointSubtype_DPST_Value_Mol
	case 12:
		return KnxDatapointSubtype_DPST_Invert
	case 120:
		return KnxDatapointSubtype_DPST_Value_Amplitude
	case 121:
		return KnxDatapointSubtype_DPST_Value_AngleRad
	case 122:
		return KnxDatapointSubtype_DPST_Value_AngleDeg
	case 123:
		return KnxDatapointSubtype_DPST_Value_Angular_Momentum
	case 124:
		return KnxDatapointSubtype_DPST_Value_Angular_Velocity
	case 125:
		return KnxDatapointSubtype_DPST_Value_Area
	case 126:
		return KnxDatapointSubtype_DPST_Value_Capacitance
	case 127:
		return KnxDatapointSubtype_DPST_Value_Charge_DensitySurface
	case 128:
		return KnxDatapointSubtype_DPST_Value_Charge_DensityVolume
	case 129:
		return KnxDatapointSubtype_DPST_Value_Compressibility
	case 13:
		return KnxDatapointSubtype_DPST_DimSendStyle
	case 130:
		return KnxDatapointSubtype_DPST_Value_Conductance
	case 131:
		return KnxDatapointSubtype_DPST_Value_Electrical_Conductivity
	case 132:
		return KnxDatapointSubtype_DPST_Value_Density
	case 133:
		return KnxDatapointSubtype_DPST_Value_Electric_Charge
	case 134:
		return KnxDatapointSubtype_DPST_Value_Electric_Current
	case 135:
		return KnxDatapointSubtype_DPST_Value_Electric_CurrentDensity
	case 136:
		return KnxDatapointSubtype_DPST_Value_Electric_DipoleMoment
	case 137:
		return KnxDatapointSubtype_DPST_Value_Electric_Displacement
	case 138:
		return KnxDatapointSubtype_DPST_Value_Electric_FieldStrength
	case 139:
		return KnxDatapointSubtype_DPST_Value_Electric_Flux
	case 14:
		return KnxDatapointSubtype_DPST_InputSource
	case 140:
		return KnxDatapointSubtype_DPST_Value_Electric_FluxDensity
	case 141:
		return KnxDatapointSubtype_DPST_Value_Electric_Polarization
	case 142:
		return KnxDatapointSubtype_DPST_Value_Electric_Potential
	case 143:
		return KnxDatapointSubtype_DPST_Value_Electric_PotentialDifference
	case 144:
		return KnxDatapointSubtype_DPST_Value_ElectromagneticMoment
	case 145:
		return KnxDatapointSubtype_DPST_Value_Electromotive_Force
	case 146:
		return KnxDatapointSubtype_DPST_Value_Energy
	case 147:
		return KnxDatapointSubtype_DPST_Value_Force
	case 148:
		return KnxDatapointSubtype_DPST_Value_Frequency
	case 149:
		return KnxDatapointSubtype_DPST_Value_Angular_Frequency
	case 15:
		return KnxDatapointSubtype_DPST_Reset
	case 150:
		return KnxDatapointSubtype_DPST_Value_Heat_Capacity
	case 151:
		return KnxDatapointSubtype_DPST_Value_Heat_FlowRate
	case 152:
		return KnxDatapointSubtype_DPST_Value_Heat_Quantity
	case 153:
		return KnxDatapointSubtype_DPST_Value_Impedance
	case 154:
		return KnxDatapointSubtype_DPST_Value_Length
	case 155:
		return KnxDatapointSubtype_DPST_Value_Light_Quantity
	case 156:
		return KnxDatapointSubtype_DPST_Value_Luminance
	case 157:
		return KnxDatapointSubtype_DPST_Value_Luminous_Flux
	case 158:
		return KnxDatapointSubtype_DPST_Value_Luminous_Intensity
	case 159:
		return KnxDatapointSubtype_DPST_Value_Magnetic_FieldStrength
	case 16:
		return KnxDatapointSubtype_DPST_Ack
	case 160:
		return KnxDatapointSubtype_DPST_Value_Magnetic_Flux
	case 161:
		return KnxDatapointSubtype_DPST_Value_Magnetic_FluxDensity
	case 162:
		return KnxDatapointSubtype_DPST_Value_Magnetic_Moment
	case 163:
		return KnxDatapointSubtype_DPST_Value_Magnetic_Polarization
	case 164:
		return KnxDatapointSubtype_DPST_Value_Magnetization
	case 165:
		return KnxDatapointSubtype_DPST_Value_MagnetomotiveForce
	case 166:
		return KnxDatapointSubtype_DPST_Value_Mass
	case 167:
		return KnxDatapointSubtype_DPST_Value_MassFlux
	case 168:
		return KnxDatapointSubtype_DPST_Value_Momentum
	case 169:
		return KnxDatapointSubtype_DPST_Value_Phase_AngleRad
	case 17:
		return KnxDatapointSubtype_DPST_Trigger
	case 170:
		return KnxDatapointSubtype_DPST_Value_Phase_AngleDeg
	case 171:
		return KnxDatapointSubtype_DPST_Value_Power
	case 172:
		return KnxDatapointSubtype_DPST_Value_Power_Factor
	case 173:
		return KnxDatapointSubtype_DPST_Value_Pressure
	case 174:
		return KnxDatapointSubtype_DPST_Value_Reactance
	case 175:
		return KnxDatapointSubtype_DPST_Value_Resistance
	case 176:
		return KnxDatapointSubtype_DPST_Value_Resistivity
	case 177:
		return KnxDatapointSubtype_DPST_Value_SelfInductance
	case 178:
		return KnxDatapointSubtype_DPST_Value_SolidAngle
	case 179:
		return KnxDatapointSubtype_DPST_Value_Sound_Intensity
	case 18:
		return KnxDatapointSubtype_DPST_Occupancy
	case 180:
		return KnxDatapointSubtype_DPST_Value_Speed
	case 181:
		return KnxDatapointSubtype_DPST_Value_Stress
	case 182:
		return KnxDatapointSubtype_DPST_Value_Surface_Tension
	case 183:
		return KnxDatapointSubtype_DPST_Value_Common_Temperature
	case 184:
		return KnxDatapointSubtype_DPST_Value_Absolute_Temperature
	case 185:
		return KnxDatapointSubtype_DPST_Value_TemperatureDifference
	case 186:
		return KnxDatapointSubtype_DPST_Value_Thermal_Capacity
	case 187:
		return KnxDatapointSubtype_DPST_Value_Thermal_Conductivity
	case 188:
		return KnxDatapointSubtype_DPST_Value_ThermoelectricPower
	case 189:
		return KnxDatapointSubtype_DPST_Value_Time
	case 19:
		return KnxDatapointSubtype_DPST_Window_Door
	case 190:
		return KnxDatapointSubtype_DPST_Value_Torque
	case 191:
		return KnxDatapointSubtype_DPST_Value_Volume
	case 192:
		return KnxDatapointSubtype_DPST_Value_Volume_Flux
	case 193:
		return KnxDatapointSubtype_DPST_Value_Weight
	case 194:
		return KnxDatapointSubtype_DPST_Value_Work
	case 195:
		return KnxDatapointSubtype_DPST_Volume_Flux_Meter
	case 196:
		return KnxDatapointSubtype_DPST_Volume_Flux_ls
	case 197:
		return KnxDatapointSubtype_DPST_Access_Data
	case 198:
		return KnxDatapointSubtype_DPST_String_ASCII
	case 199:
		return KnxDatapointSubtype_DPST_String_8859_1
	case 2:
		return KnxDatapointSubtype_DPST_Bool
	case 20:
		return KnxDatapointSubtype_DPST_LogicalFunction
	case 200:
		return KnxDatapointSubtype_DPST_SceneNumber
	case 201:
		return KnxDatapointSubtype_DPST_SceneControl
	case 202:
		return KnxDatapointSubtype_DPST_DateTime
	case 203:
		return KnxDatapointSubtype_DPST_SCLOMode
	case 204:
		return KnxDatapointSubtype_DPST_BuildingMode
	case 205:
		return KnxDatapointSubtype_DPST_OccMode
	case 206:
		return KnxDatapointSubtype_DPST_Priority
	case 207:
		return KnxDatapointSubtype_DPST_LightApplicationMode
	case 208:
		return KnxDatapointSubtype_DPST_ApplicationArea
	case 209:
		return KnxDatapointSubtype_DPST_AlarmClassType
	case 21:
		return KnxDatapointSubtype_DPST_Scene_AB
	case 210:
		return KnxDatapointSubtype_DPST_PSUMode
	case 211:
		return KnxDatapointSubtype_DPST_ErrorClass_System
	case 212:
		return KnxDatapointSubtype_DPST_ErrorClass_HVAC
	case 213:
		return KnxDatapointSubtype_DPST_Time_Delay
	case 214:
		return KnxDatapointSubtype_DPST_Beaufort_Wind_Force_Scale
	case 215:
		return KnxDatapointSubtype_DPST_SensorSelect
	case 216:
		return KnxDatapointSubtype_DPST_ActuatorConnectType
	case 217:
		return KnxDatapointSubtype_DPST_Cloud_Cover
	case 218:
		return KnxDatapointSubtype_DPST_PowerReturnMode
	case 219:
		return KnxDatapointSubtype_DPST_FuelType
	case 22:
		return KnxDatapointSubtype_DPST_ShutterBlinds_Mode
	case 220:
		return KnxDatapointSubtype_DPST_BurnerType
	case 221:
		return KnxDatapointSubtype_DPST_HVACMode
	case 222:
		return KnxDatapointSubtype_DPST_DHWMode
	case 223:
		return KnxDatapointSubtype_DPST_LoadPriority
	case 224:
		return KnxDatapointSubtype_DPST_HVACContrMode
	case 225:
		return KnxDatapointSubtype_DPST_HVACEmergMode
	case 226:
		return KnxDatapointSubtype_DPST_ChangeoverMode
	case 227:
		return KnxDatapointSubtype_DPST_ValveMode
	case 228:
		return KnxDatapointSubtype_DPST_DamperMode
	case 229:
		return KnxDatapointSubtype_DPST_HeaterMode
	case 23:
		return KnxDatapointSubtype_DPST_DayNight
	case 230:
		return KnxDatapointSubtype_DPST_FanMode
	case 231:
		return KnxDatapointSubtype_DPST_MasterSlaveMode
	case 232:
		return KnxDatapointSubtype_DPST_StatusRoomSetp
	case 233:
		return KnxDatapointSubtype_DPST_Metering_DeviceType
	case 234:
		return KnxDatapointSubtype_DPST_HumDehumMode
	case 235:
		return KnxDatapointSubtype_DPST_EnableHCStage
	case 236:
		return KnxDatapointSubtype_DPST_ADAType
	case 237:
		return KnxDatapointSubtype_DPST_BackupMode
	case 238:
		return KnxDatapointSubtype_DPST_StartSynchronization
	case 239:
		return KnxDatapointSubtype_DPST_Behaviour_Lock_Unlock
	case 24:
		return KnxDatapointSubtype_DPST_Heat_Cool
	case 240:
		return KnxDatapointSubtype_DPST_Behaviour_Bus_Power_Up_Down
	case 241:
		return KnxDatapointSubtype_DPST_DALI_Fade_Time
	case 242:
		return KnxDatapointSubtype_DPST_BlinkingMode
	case 243:
		return KnxDatapointSubtype_DPST_LightControlMode
	case 244:
		return KnxDatapointSubtype_DPST_SwitchPBModel
	case 245:
		return KnxDatapointSubtype_DPST_PBAction
	case 246:
		return KnxDatapointSubtype_DPST_DimmPBModel
	case 247:
		return KnxDatapointSubtype_DPST_SwitchOnMode
	case 248:
		return KnxDatapointSubtype_DPST_LoadTypeSet
	case 249:
		return KnxDatapointSubtype_DPST_LoadTypeDetected
	case 25:
		return KnxDatapointSubtype_DPST_Switch_Control
	case 250:
		return KnxDatapointSubtype_DPST_Converter_Test_Control
	case 251:
		return KnxDatapointSubtype_DPST_SABExcept_Behaviour
	case 252:
		return KnxDatapointSubtype_DPST_SABBehaviour_Lock_Unlock
	case 253:
		return KnxDatapointSubtype_DPST_SSSBMode
	case 254:
		return KnxDatapointSubtype_DPST_BlindsControlMode
	case 255:
		return KnxDatapointSubtype_DPST_CommMode
	case 256:
		return KnxDatapointSubtype_DPST_AddInfoTypes
	case 257:
		return KnxDatapointSubtype_DPST_RF_ModeSelect
	case 258:
		return KnxDatapointSubtype_DPST_RF_FilterSelect
	case 259:
		return KnxDatapointSubtype_DPST_StatusGen
	case 26:
		return KnxDatapointSubtype_DPST_Bool_Control
	case 260:
		return KnxDatapointSubtype_DPST_Device_Control
	case 261:
		return KnxDatapointSubtype_DPST_ForceSign
	case 262:
		return KnxDatapointSubtype_DPST_ForceSignCool
	case 263:
		return KnxDatapointSubtype_DPST_StatusRHC
	case 264:
		return KnxDatapointSubtype_DPST_StatusSDHWC
	case 265:
		return KnxDatapointSubtype_DPST_FuelTypeSet
	case 266:
		return KnxDatapointSubtype_DPST_StatusRCC
	case 267:
		return KnxDatapointSubtype_DPST_StatusAHU
	case 268:
		return KnxDatapointSubtype_DPST_CombinedStatus_RTSM
	case 269:
		return KnxDatapointSubtype_DPST_LightActuatorErrorInfo
	case 27:
		return KnxDatapointSubtype_DPST_Enable_Control
	case 270:
		return KnxDatapointSubtype_DPST_RF_ModeInfo
	case 271:
		return KnxDatapointSubtype_DPST_RF_FilterInfo
	case 272:
		return KnxDatapointSubtype_DPST_Channel_Activation_8
	case 273:
		return KnxDatapointSubtype_DPST_StatusDHWC
	case 274:
		return KnxDatapointSubtype_DPST_StatusRHCC
	case 275:
		return KnxDatapointSubtype_DPST_CombinedStatus_HVA
	case 276:
		return KnxDatapointSubtype_DPST_CombinedStatus_RTC
	case 277:
		return KnxDatapointSubtype_DPST_Media
	case 278:
		return KnxDatapointSubtype_DPST_Channel_Activation_16
	case 279:
		return KnxDatapointSubtype_DPST_OnOffAction
	case 28:
		return KnxDatapointSubtype_DPST_Ramp_Control
	case 280:
		return KnxDatapointSubtype_DPST_Alarm_Reaction
	case 281:
		return KnxDatapointSubtype_DPST_UpDown_Action
	case 282:
		return KnxDatapointSubtype_DPST_HVAC_PB_Action
	case 283:
		return KnxDatapointSubtype_DPST_DoubleNibble
	case 284:
		return KnxDatapointSubtype_DPST_SceneInfo
	case 285:
		return KnxDatapointSubtype_DPST_CombinedInfoOnOff
	case 286:
		return KnxDatapointSubtype_DPST_ActiveEnergy_V64
	case 287:
		return KnxDatapointSubtype_DPST_ApparantEnergy_V64
	case 288:
		return KnxDatapointSubtype_DPST_ReactiveEnergy_V64
	case 289:
		return KnxDatapointSubtype_DPST_Channel_Activation_24
	case 29:
		return KnxDatapointSubtype_DPST_Alarm_Control
	case 290:
		return KnxDatapointSubtype_DPST_HVACModeNext
	case 291:
		return KnxDatapointSubtype_DPST_DHWModeNext
	case 292:
		return KnxDatapointSubtype_DPST_OccModeNext
	case 293:
		return KnxDatapointSubtype_DPST_BuildingModeNext
	case 294:
		return KnxDatapointSubtype_DPST_Version
	case 295:
		return KnxDatapointSubtype_DPST_AlarmInfo
	case 296:
		return KnxDatapointSubtype_DPST_TempRoomSetpSetF16_3
	case 297:
		return KnxDatapointSubtype_DPST_TempRoomSetpSetShiftF16_3
	case 298:
		return KnxDatapointSubtype_DPST_Scaling_Speed
	case 299:
		return KnxDatapointSubtype_DPST_Scaling_Step_Time
	case 3:
		return KnxDatapointSubtype_DPST_Enable
	case 30:
		return KnxDatapointSubtype_DPST_BinaryValue_Control
	case 300:
		return KnxDatapointSubtype_DPST_MeteringValue
	case 301:
		return KnxDatapointSubtype_DPST_MBus_Address
	case 302:
		return KnxDatapointSubtype_DPST_Colour_RGB
	case 303:
		return KnxDatapointSubtype_DPST_LanguageCodeAlpha2_ASCII
	case 304:
		return KnxDatapointSubtype_DPST_Tariff_ActiveEnergy
	case 305:
		return KnxDatapointSubtype_DPST_Prioritised_Mode_Control
	case 306:
		return KnxDatapointSubtype_DPST_DALI_Control_Gear_Diagnostic
	case 307:
		return KnxDatapointSubtype_DPST_DALI_Diagnostics
	case 308:
		return KnxDatapointSubtype_DPST_CombinedPosition
	case 309:
		return KnxDatapointSubtype_DPST_StatusSAB
	case 31:
		return KnxDatapointSubtype_DPST_Step_Control
	case 310:
		return KnxDatapointSubtype_DPST_Colour_xyY
	case 311:
		return KnxDatapointSubtype_DPST_Converter_Status
	case 312:
		return KnxDatapointSubtype_DPST_Converter_Test_Result
	case 313:
		return KnxDatapointSubtype_DPST_Battery_Info
	case 314:
		return KnxDatapointSubtype_DPST_Brightness_Colour_Temperature_Transition
	case 315:
		return KnxDatapointSubtype_DPST_Brightness_Colour_Temperature_Control
	case 316:
		return KnxDatapointSubtype_DPST_Colour_RGBW
	case 317:
		return KnxDatapointSubtype_DPST_Relative_Control_RGBW
	case 318:
		return KnxDatapointSubtype_DPST_Relative_Control_RGB
	case 319:
		return KnxDatapointSubtype_DPST_GeographicalLocation
	case 32:
		return KnxDatapointSubtype_DPST_Direction1_Control
	case 320:
		return KnxDatapointSubtype_DPST_TempRoomSetpSetF16_4
	case 321:
		return KnxDatapointSubtype_DPST_TempRoomSetpSetShiftF16_4
	case 33:
		return KnxDatapointSubtype_DPST_Direction2_Control
	case 34:
		return KnxDatapointSubtype_DPST_Start_Control
	case 35:
		return KnxDatapointSubtype_DPST_State_Control
	case 36:
		return KnxDatapointSubtype_DPST_Invert_Control
	case 37:
		return KnxDatapointSubtype_DPST_Control_Dimming
	case 38:
		return KnxDatapointSubtype_DPST_Control_Blinds
	case 39:
		return KnxDatapointSubtype_DPST_Char_ASCII
	case 4:
		return KnxDatapointSubtype_DPST_Ramp
	case 40:
		return KnxDatapointSubtype_DPST_Char_8859_1
	case 41:
		return KnxDatapointSubtype_DPST_Scaling
	case 42:
		return KnxDatapointSubtype_DPST_Angle
	case 43:
		return KnxDatapointSubtype_DPST_Percent_U8
	case 44:
		return KnxDatapointSubtype_DPST_DecimalFactor
	case 45:
		return KnxDatapointSubtype_DPST_Tariff
	case 46:
		return KnxDatapointSubtype_DPST_Value_1_Ucount
	case 47:
		return KnxDatapointSubtype_DPST_FanStage
	case 48:
		return KnxDatapointSubtype_DPST_Percent_V8
	case 49:
		return KnxDatapointSubtype_DPST_Value_1_Count
	case 5:
		return KnxDatapointSubtype_DPST_Alarm
	case 50:
		return KnxDatapointSubtype_DPST_Status_Mode3
	case 51:
		return KnxDatapointSubtype_DPST_Value_2_Ucount
	case 52:
		return KnxDatapointSubtype_DPST_TimePeriodMsec
	case 53:
		return KnxDatapointSubtype_DPST_TimePeriod10Msec
	case 54:
		return KnxDatapointSubtype_DPST_TimePeriod100Msec
	case 55:
		return KnxDatapointSubtype_DPST_TimePeriodSec
	case 56:
		return KnxDatapointSubtype_DPST_TimePeriodMin
	case 57:
		return KnxDatapointSubtype_DPST_TimePeriodHrs
	case 58:
		return KnxDatapointSubtype_DPST_PropDataType
	case 59:
		return KnxDatapointSubtype_DPST_Length_mm
	case 6:
		return KnxDatapointSubtype_DPST_BinaryValue
	case 60:
		return KnxDatapointSubtype_DPST_UElCurrentmA
	case 61:
		return KnxDatapointSubtype_DPST_Brightness
	case 62:
		return KnxDatapointSubtype_DPST_Absolute_Colour_Temperature
	case 63:
		return KnxDatapointSubtype_DPST_Value_2_Count
	case 64:
		return KnxDatapointSubtype_DPST_DeltaTimeMsec
	case 65:
		return KnxDatapointSubtype_DPST_DeltaTime10Msec
	case 66:
		return KnxDatapointSubtype_DPST_DeltaTime100Msec
	case 67:
		return KnxDatapointSubtype_DPST_DeltaTimeSec
	case 68:
		return KnxDatapointSubtype_DPST_DeltaTimeMin
	case 69:
		return KnxDatapointSubtype_DPST_DeltaTimeHrs
	case 7:
		return KnxDatapointSubtype_DPST_Step
	case 70:
		return KnxDatapointSubtype_DPST_Percent_V16
	case 71:
		return KnxDatapointSubtype_DPST_Rotation_Angle
	case 72:
		return KnxDatapointSubtype_DPST_Length_m
	case 73:
		return KnxDatapointSubtype_DPST_Value_Temp
	case 74:
		return KnxDatapointSubtype_DPST_Value_Tempd
	case 75:
		return KnxDatapointSubtype_DPST_Value_Tempa
	case 76:
		return KnxDatapointSubtype_DPST_Value_Lux
	case 77:
		return KnxDatapointSubtype_DPST_Value_Wsp
	case 78:
		return KnxDatapointSubtype_DPST_Value_Pres
	case 79:
		return KnxDatapointSubtype_DPST_Value_Humidity
	case 8:
		return KnxDatapointSubtype_DPST_UpDown
	case 80:
		return KnxDatapointSubtype_DPST_Value_AirQuality
	case 81:
		return KnxDatapointSubtype_DPST_Value_AirFlow
	case 82:
		return KnxDatapointSubtype_DPST_Value_Time1
	case 83:
		return KnxDatapointSubtype_DPST_Value_Time2
	case 84:
		return KnxDatapointSubtype_DPST_Value_Volt
	case 85:
		return KnxDatapointSubtype_DPST_Value_Curr
	case 86:
		return KnxDatapointSubtype_DPST_PowerDensity
	case 87:
		return KnxDatapointSubtype_DPST_KelvinPerPercent
	case 88:
		return KnxDatapointSubtype_DPST_Power
	case 89:
		return KnxDatapointSubtype_DPST_Value_Volume_Flow
	case 9:
		return KnxDatapointSubtype_DPST_OpenClose
	case 90:
		return KnxDatapointSubtype_DPST_Rain_Amount
	case 91:
		return KnxDatapointSubtype_DPST_Value_Temp_F
	case 92:
		return KnxDatapointSubtype_DPST_Value_Wsp_kmh
	case 93:
		return KnxDatapointSubtype_DPST_Value_Absolute_Humidity
	case 94:
		return KnxDatapointSubtype_DPST_Concentration_ygm3
	case 95:
		return KnxDatapointSubtype_DPST_TimeOfDay
	case 96:
		return KnxDatapointSubtype_DPST_Date
	case 97:
		return KnxDatapointSubtype_DPST_Value_4_Ucount
	case 98:
		return KnxDatapointSubtype_DPST_LongTimePeriod_Sec
	case 99:
		return KnxDatapointSubtype_DPST_LongTimePeriod_Min
	}
	return 0
}

func KnxDatapointSubtypeByName(value string) KnxDatapointSubtype {
	switch value {
	case "DPST_UNKNOWN":
		return KnxDatapointSubtype_DPST_UNKNOWN
	case "DPST_Switch":
		return KnxDatapointSubtype_DPST_Switch
	case "DPST_Start":
		return KnxDatapointSubtype_DPST_Start
	case "DPST_LongTimePeriod_Hrs":
		return KnxDatapointSubtype_DPST_LongTimePeriod_Hrs
	case "DPST_VolumeLiquid_Litre":
		return KnxDatapointSubtype_DPST_VolumeLiquid_Litre
	case "DPST_Volume_m_3":
		return KnxDatapointSubtype_DPST_Volume_m_3
	case "DPST_Value_4_Count":
		return KnxDatapointSubtype_DPST_Value_4_Count
	case "DPST_FlowRate_m3h":
		return KnxDatapointSubtype_DPST_FlowRate_m3h
	case "DPST_ActiveEnergy":
		return KnxDatapointSubtype_DPST_ActiveEnergy
	case "DPST_ApparantEnergy":
		return KnxDatapointSubtype_DPST_ApparantEnergy
	case "DPST_ReactiveEnergy":
		return KnxDatapointSubtype_DPST_ReactiveEnergy
	case "DPST_ActiveEnergy_kWh":
		return KnxDatapointSubtype_DPST_ActiveEnergy_kWh
	case "DPST_ApparantEnergy_kVAh":
		return KnxDatapointSubtype_DPST_ApparantEnergy_kVAh
	case "DPST_State":
		return KnxDatapointSubtype_DPST_State
	case "DPST_ReactiveEnergy_kVARh":
		return KnxDatapointSubtype_DPST_ReactiveEnergy_kVARh
	case "DPST_ActiveEnergy_MWh":
		return KnxDatapointSubtype_DPST_ActiveEnergy_MWh
	case "DPST_LongDeltaTimeSec":
		return KnxDatapointSubtype_DPST_LongDeltaTimeSec
	case "DPST_DeltaVolumeLiquid_Litre":
		return KnxDatapointSubtype_DPST_DeltaVolumeLiquid_Litre
	case "DPST_DeltaVolume_m_3":
		return KnxDatapointSubtype_DPST_DeltaVolume_m_3
	case "DPST_Value_Acceleration":
		return KnxDatapointSubtype_DPST_Value_Acceleration
	case "DPST_Value_Acceleration_Angular":
		return KnxDatapointSubtype_DPST_Value_Acceleration_Angular
	case "DPST_Value_Activation_Energy":
		return KnxDatapointSubtype_DPST_Value_Activation_Energy
	case "DPST_Value_Activity":
		return KnxDatapointSubtype_DPST_Value_Activity
	case "DPST_Value_Mol":
		return KnxDatapointSubtype_DPST_Value_Mol
	case "DPST_Invert":
		return KnxDatapointSubtype_DPST_Invert
	case "DPST_Value_Amplitude":
		return KnxDatapointSubtype_DPST_Value_Amplitude
	case "DPST_Value_AngleRad":
		return KnxDatapointSubtype_DPST_Value_AngleRad
	case "DPST_Value_AngleDeg":
		return KnxDatapointSubtype_DPST_Value_AngleDeg
	case "DPST_Value_Angular_Momentum":
		return KnxDatapointSubtype_DPST_Value_Angular_Momentum
	case "DPST_Value_Angular_Velocity":
		return KnxDatapointSubtype_DPST_Value_Angular_Velocity
	case "DPST_Value_Area":
		return KnxDatapointSubtype_DPST_Value_Area
	case "DPST_Value_Capacitance":
		return KnxDatapointSubtype_DPST_Value_Capacitance
	case "DPST_Value_Charge_DensitySurface":
		return KnxDatapointSubtype_DPST_Value_Charge_DensitySurface
	case "DPST_Value_Charge_DensityVolume":
		return KnxDatapointSubtype_DPST_Value_Charge_DensityVolume
	case "DPST_Value_Compressibility":
		return KnxDatapointSubtype_DPST_Value_Compressibility
	case "DPST_DimSendStyle":
		return KnxDatapointSubtype_DPST_DimSendStyle
	case "DPST_Value_Conductance":
		return KnxDatapointSubtype_DPST_Value_Conductance
	case "DPST_Value_Electrical_Conductivity":
		return KnxDatapointSubtype_DPST_Value_Electrical_Conductivity
	case "DPST_Value_Density":
		return KnxDatapointSubtype_DPST_Value_Density
	case "DPST_Value_Electric_Charge":
		return KnxDatapointSubtype_DPST_Value_Electric_Charge
	case "DPST_Value_Electric_Current":
		return KnxDatapointSubtype_DPST_Value_Electric_Current
	case "DPST_Value_Electric_CurrentDensity":
		return KnxDatapointSubtype_DPST_Value_Electric_CurrentDensity
	case "DPST_Value_Electric_DipoleMoment":
		return KnxDatapointSubtype_DPST_Value_Electric_DipoleMoment
	case "DPST_Value_Electric_Displacement":
		return KnxDatapointSubtype_DPST_Value_Electric_Displacement
	case "DPST_Value_Electric_FieldStrength":
		return KnxDatapointSubtype_DPST_Value_Electric_FieldStrength
	case "DPST_Value_Electric_Flux":
		return KnxDatapointSubtype_DPST_Value_Electric_Flux
	case "DPST_InputSource":
		return KnxDatapointSubtype_DPST_InputSource
	case "DPST_Value_Electric_FluxDensity":
		return KnxDatapointSubtype_DPST_Value_Electric_FluxDensity
	case "DPST_Value_Electric_Polarization":
		return KnxDatapointSubtype_DPST_Value_Electric_Polarization
	case "DPST_Value_Electric_Potential":
		return KnxDatapointSubtype_DPST_Value_Electric_Potential
	case "DPST_Value_Electric_PotentialDifference":
		return KnxDatapointSubtype_DPST_Value_Electric_PotentialDifference
	case "DPST_Value_ElectromagneticMoment":
		return KnxDatapointSubtype_DPST_Value_ElectromagneticMoment
	case "DPST_Value_Electromotive_Force":
		return KnxDatapointSubtype_DPST_Value_Electromotive_Force
	case "DPST_Value_Energy":
		return KnxDatapointSubtype_DPST_Value_Energy
	case "DPST_Value_Force":
		return KnxDatapointSubtype_DPST_Value_Force
	case "DPST_Value_Frequency":
		return KnxDatapointSubtype_DPST_Value_Frequency
	case "DPST_Value_Angular_Frequency":
		return KnxDatapointSubtype_DPST_Value_Angular_Frequency
	case "DPST_Reset":
		return KnxDatapointSubtype_DPST_Reset
	case "DPST_Value_Heat_Capacity":
		return KnxDatapointSubtype_DPST_Value_Heat_Capacity
	case "DPST_Value_Heat_FlowRate":
		return KnxDatapointSubtype_DPST_Value_Heat_FlowRate
	case "DPST_Value_Heat_Quantity":
		return KnxDatapointSubtype_DPST_Value_Heat_Quantity
	case "DPST_Value_Impedance":
		return KnxDatapointSubtype_DPST_Value_Impedance
	case "DPST_Value_Length":
		return KnxDatapointSubtype_DPST_Value_Length
	case "DPST_Value_Light_Quantity":
		return KnxDatapointSubtype_DPST_Value_Light_Quantity
	case "DPST_Value_Luminance":
		return KnxDatapointSubtype_DPST_Value_Luminance
	case "DPST_Value_Luminous_Flux":
		return KnxDatapointSubtype_DPST_Value_Luminous_Flux
	case "DPST_Value_Luminous_Intensity":
		return KnxDatapointSubtype_DPST_Value_Luminous_Intensity
	case "DPST_Value_Magnetic_FieldStrength":
		return KnxDatapointSubtype_DPST_Value_Magnetic_FieldStrength
	case "DPST_Ack":
		return KnxDatapointSubtype_DPST_Ack
	case "DPST_Value_Magnetic_Flux":
		return KnxDatapointSubtype_DPST_Value_Magnetic_Flux
	case "DPST_Value_Magnetic_FluxDensity":
		return KnxDatapointSubtype_DPST_Value_Magnetic_FluxDensity
	case "DPST_Value_Magnetic_Moment":
		return KnxDatapointSubtype_DPST_Value_Magnetic_Moment
	case "DPST_Value_Magnetic_Polarization":
		return KnxDatapointSubtype_DPST_Value_Magnetic_Polarization
	case "DPST_Value_Magnetization":
		return KnxDatapointSubtype_DPST_Value_Magnetization
	case "DPST_Value_MagnetomotiveForce":
		return KnxDatapointSubtype_DPST_Value_MagnetomotiveForce
	case "DPST_Value_Mass":
		return KnxDatapointSubtype_DPST_Value_Mass
	case "DPST_Value_MassFlux":
		return KnxDatapointSubtype_DPST_Value_MassFlux
	case "DPST_Value_Momentum":
		return KnxDatapointSubtype_DPST_Value_Momentum
	case "DPST_Value_Phase_AngleRad":
		return KnxDatapointSubtype_DPST_Value_Phase_AngleRad
	case "DPST_Trigger":
		return KnxDatapointSubtype_DPST_Trigger
	case "DPST_Value_Phase_AngleDeg":
		return KnxDatapointSubtype_DPST_Value_Phase_AngleDeg
	case "DPST_Value_Power":
		return KnxDatapointSubtype_DPST_Value_Power
	case "DPST_Value_Power_Factor":
		return KnxDatapointSubtype_DPST_Value_Power_Factor
	case "DPST_Value_Pressure":
		return KnxDatapointSubtype_DPST_Value_Pressure
	case "DPST_Value_Reactance":
		return KnxDatapointSubtype_DPST_Value_Reactance
	case "DPST_Value_Resistance":
		return KnxDatapointSubtype_DPST_Value_Resistance
	case "DPST_Value_Resistivity":
		return KnxDatapointSubtype_DPST_Value_Resistivity
	case "DPST_Value_SelfInductance":
		return KnxDatapointSubtype_DPST_Value_SelfInductance
	case "DPST_Value_SolidAngle":
		return KnxDatapointSubtype_DPST_Value_SolidAngle
	case "DPST_Value_Sound_Intensity":
		return KnxDatapointSubtype_DPST_Value_Sound_Intensity
	case "DPST_Occupancy":
		return KnxDatapointSubtype_DPST_Occupancy
	case "DPST_Value_Speed":
		return KnxDatapointSubtype_DPST_Value_Speed
	case "DPST_Value_Stress":
		return KnxDatapointSubtype_DPST_Value_Stress
	case "DPST_Value_Surface_Tension":
		return KnxDatapointSubtype_DPST_Value_Surface_Tension
	case "DPST_Value_Common_Temperature":
		return KnxDatapointSubtype_DPST_Value_Common_Temperature
	case "DPST_Value_Absolute_Temperature":
		return KnxDatapointSubtype_DPST_Value_Absolute_Temperature
	case "DPST_Value_TemperatureDifference":
		return KnxDatapointSubtype_DPST_Value_TemperatureDifference
	case "DPST_Value_Thermal_Capacity":
		return KnxDatapointSubtype_DPST_Value_Thermal_Capacity
	case "DPST_Value_Thermal_Conductivity":
		return KnxDatapointSubtype_DPST_Value_Thermal_Conductivity
	case "DPST_Value_ThermoelectricPower":
		return KnxDatapointSubtype_DPST_Value_ThermoelectricPower
	case "DPST_Value_Time":
		return KnxDatapointSubtype_DPST_Value_Time
	case "DPST_Window_Door":
		return KnxDatapointSubtype_DPST_Window_Door
	case "DPST_Value_Torque":
		return KnxDatapointSubtype_DPST_Value_Torque
	case "DPST_Value_Volume":
		return KnxDatapointSubtype_DPST_Value_Volume
	case "DPST_Value_Volume_Flux":
		return KnxDatapointSubtype_DPST_Value_Volume_Flux
	case "DPST_Value_Weight":
		return KnxDatapointSubtype_DPST_Value_Weight
	case "DPST_Value_Work":
		return KnxDatapointSubtype_DPST_Value_Work
	case "DPST_Volume_Flux_Meter":
		return KnxDatapointSubtype_DPST_Volume_Flux_Meter
	case "DPST_Volume_Flux_ls":
		return KnxDatapointSubtype_DPST_Volume_Flux_ls
	case "DPST_Access_Data":
		return KnxDatapointSubtype_DPST_Access_Data
	case "DPST_String_ASCII":
		return KnxDatapointSubtype_DPST_String_ASCII
	case "DPST_String_8859_1":
		return KnxDatapointSubtype_DPST_String_8859_1
	case "DPST_Bool":
		return KnxDatapointSubtype_DPST_Bool
	case "DPST_LogicalFunction":
		return KnxDatapointSubtype_DPST_LogicalFunction
	case "DPST_SceneNumber":
		return KnxDatapointSubtype_DPST_SceneNumber
	case "DPST_SceneControl":
		return KnxDatapointSubtype_DPST_SceneControl
	case "DPST_DateTime":
		return KnxDatapointSubtype_DPST_DateTime
	case "DPST_SCLOMode":
		return KnxDatapointSubtype_DPST_SCLOMode
	case "DPST_BuildingMode":
		return KnxDatapointSubtype_DPST_BuildingMode
	case "DPST_OccMode":
		return KnxDatapointSubtype_DPST_OccMode
	case "DPST_Priority":
		return KnxDatapointSubtype_DPST_Priority
	case "DPST_LightApplicationMode":
		return KnxDatapointSubtype_DPST_LightApplicationMode
	case "DPST_ApplicationArea":
		return KnxDatapointSubtype_DPST_ApplicationArea
	case "DPST_AlarmClassType":
		return KnxDatapointSubtype_DPST_AlarmClassType
	case "DPST_Scene_AB":
		return KnxDatapointSubtype_DPST_Scene_AB
	case "DPST_PSUMode":
		return KnxDatapointSubtype_DPST_PSUMode
	case "DPST_ErrorClass_System":
		return KnxDatapointSubtype_DPST_ErrorClass_System
	case "DPST_ErrorClass_HVAC":
		return KnxDatapointSubtype_DPST_ErrorClass_HVAC
	case "DPST_Time_Delay":
		return KnxDatapointSubtype_DPST_Time_Delay
	case "DPST_Beaufort_Wind_Force_Scale":
		return KnxDatapointSubtype_DPST_Beaufort_Wind_Force_Scale
	case "DPST_SensorSelect":
		return KnxDatapointSubtype_DPST_SensorSelect
	case "DPST_ActuatorConnectType":
		return KnxDatapointSubtype_DPST_ActuatorConnectType
	case "DPST_Cloud_Cover":
		return KnxDatapointSubtype_DPST_Cloud_Cover
	case "DPST_PowerReturnMode":
		return KnxDatapointSubtype_DPST_PowerReturnMode
	case "DPST_FuelType":
		return KnxDatapointSubtype_DPST_FuelType
	case "DPST_ShutterBlinds_Mode":
		return KnxDatapointSubtype_DPST_ShutterBlinds_Mode
	case "DPST_BurnerType":
		return KnxDatapointSubtype_DPST_BurnerType
	case "DPST_HVACMode":
		return KnxDatapointSubtype_DPST_HVACMode
	case "DPST_DHWMode":
		return KnxDatapointSubtype_DPST_DHWMode
	case "DPST_LoadPriority":
		return KnxDatapointSubtype_DPST_LoadPriority
	case "DPST_HVACContrMode":
		return KnxDatapointSubtype_DPST_HVACContrMode
	case "DPST_HVACEmergMode":
		return KnxDatapointSubtype_DPST_HVACEmergMode
	case "DPST_ChangeoverMode":
		return KnxDatapointSubtype_DPST_ChangeoverMode
	case "DPST_ValveMode":
		return KnxDatapointSubtype_DPST_ValveMode
	case "DPST_DamperMode":
		return KnxDatapointSubtype_DPST_DamperMode
	case "DPST_HeaterMode":
		return KnxDatapointSubtype_DPST_HeaterMode
	case "DPST_DayNight":
		return KnxDatapointSubtype_DPST_DayNight
	case "DPST_FanMode":
		return KnxDatapointSubtype_DPST_FanMode
	case "DPST_MasterSlaveMode":
		return KnxDatapointSubtype_DPST_MasterSlaveMode
	case "DPST_StatusRoomSetp":
		return KnxDatapointSubtype_DPST_StatusRoomSetp
	case "DPST_Metering_DeviceType":
		return KnxDatapointSubtype_DPST_Metering_DeviceType
	case "DPST_HumDehumMode":
		return KnxDatapointSubtype_DPST_HumDehumMode
	case "DPST_EnableHCStage":
		return KnxDatapointSubtype_DPST_EnableHCStage
	case "DPST_ADAType":
		return KnxDatapointSubtype_DPST_ADAType
	case "DPST_BackupMode":
		return KnxDatapointSubtype_DPST_BackupMode
	case "DPST_StartSynchronization":
		return KnxDatapointSubtype_DPST_StartSynchronization
	case "DPST_Behaviour_Lock_Unlock":
		return KnxDatapointSubtype_DPST_Behaviour_Lock_Unlock
	case "DPST_Heat_Cool":
		return KnxDatapointSubtype_DPST_Heat_Cool
	case "DPST_Behaviour_Bus_Power_Up_Down":
		return KnxDatapointSubtype_DPST_Behaviour_Bus_Power_Up_Down
	case "DPST_DALI_Fade_Time":
		return KnxDatapointSubtype_DPST_DALI_Fade_Time
	case "DPST_BlinkingMode":
		return KnxDatapointSubtype_DPST_BlinkingMode
	case "DPST_LightControlMode":
		return KnxDatapointSubtype_DPST_LightControlMode
	case "DPST_SwitchPBModel":
		return KnxDatapointSubtype_DPST_SwitchPBModel
	case "DPST_PBAction":
		return KnxDatapointSubtype_DPST_PBAction
	case "DPST_DimmPBModel":
		return KnxDatapointSubtype_DPST_DimmPBModel
	case "DPST_SwitchOnMode":
		return KnxDatapointSubtype_DPST_SwitchOnMode
	case "DPST_LoadTypeSet":
		return KnxDatapointSubtype_DPST_LoadTypeSet
	case "DPST_LoadTypeDetected":
		return KnxDatapointSubtype_DPST_LoadTypeDetected
	case "DPST_Switch_Control":
		return KnxDatapointSubtype_DPST_Switch_Control
	case "DPST_Converter_Test_Control":
		return KnxDatapointSubtype_DPST_Converter_Test_Control
	case "DPST_SABExcept_Behaviour":
		return KnxDatapointSubtype_DPST_SABExcept_Behaviour
	case "DPST_SABBehaviour_Lock_Unlock":
		return KnxDatapointSubtype_DPST_SABBehaviour_Lock_Unlock
	case "DPST_SSSBMode":
		return KnxDatapointSubtype_DPST_SSSBMode
	case "DPST_BlindsControlMode":
		return KnxDatapointSubtype_DPST_BlindsControlMode
	case "DPST_CommMode":
		return KnxDatapointSubtype_DPST_CommMode
	case "DPST_AddInfoTypes":
		return KnxDatapointSubtype_DPST_AddInfoTypes
	case "DPST_RF_ModeSelect":
		return KnxDatapointSubtype_DPST_RF_ModeSelect
	case "DPST_RF_FilterSelect":
		return KnxDatapointSubtype_DPST_RF_FilterSelect
	case "DPST_StatusGen":
		return KnxDatapointSubtype_DPST_StatusGen
	case "DPST_Bool_Control":
		return KnxDatapointSubtype_DPST_Bool_Control
	case "DPST_Device_Control":
		return KnxDatapointSubtype_DPST_Device_Control
	case "DPST_ForceSign":
		return KnxDatapointSubtype_DPST_ForceSign
	case "DPST_ForceSignCool":
		return KnxDatapointSubtype_DPST_ForceSignCool
	case "DPST_StatusRHC":
		return KnxDatapointSubtype_DPST_StatusRHC
	case "DPST_StatusSDHWC":
		return KnxDatapointSubtype_DPST_StatusSDHWC
	case "DPST_FuelTypeSet":
		return KnxDatapointSubtype_DPST_FuelTypeSet
	case "DPST_StatusRCC":
		return KnxDatapointSubtype_DPST_StatusRCC
	case "DPST_StatusAHU":
		return KnxDatapointSubtype_DPST_StatusAHU
	case "DPST_CombinedStatus_RTSM":
		return KnxDatapointSubtype_DPST_CombinedStatus_RTSM
	case "DPST_LightActuatorErrorInfo":
		return KnxDatapointSubtype_DPST_LightActuatorErrorInfo
	case "DPST_Enable_Control":
		return KnxDatapointSubtype_DPST_Enable_Control
	case "DPST_RF_ModeInfo":
		return KnxDatapointSubtype_DPST_RF_ModeInfo
	case "DPST_RF_FilterInfo":
		return KnxDatapointSubtype_DPST_RF_FilterInfo
	case "DPST_Channel_Activation_8":
		return KnxDatapointSubtype_DPST_Channel_Activation_8
	case "DPST_StatusDHWC":
		return KnxDatapointSubtype_DPST_StatusDHWC
	case "DPST_StatusRHCC":
		return KnxDatapointSubtype_DPST_StatusRHCC
	case "DPST_CombinedStatus_HVA":
		return KnxDatapointSubtype_DPST_CombinedStatus_HVA
	case "DPST_CombinedStatus_RTC":
		return KnxDatapointSubtype_DPST_CombinedStatus_RTC
	case "DPST_Media":
		return KnxDatapointSubtype_DPST_Media
	case "DPST_Channel_Activation_16":
		return KnxDatapointSubtype_DPST_Channel_Activation_16
	case "DPST_OnOffAction":
		return KnxDatapointSubtype_DPST_OnOffAction
	case "DPST_Ramp_Control":
		return KnxDatapointSubtype_DPST_Ramp_Control
	case "DPST_Alarm_Reaction":
		return KnxDatapointSubtype_DPST_Alarm_Reaction
	case "DPST_UpDown_Action":
		return KnxDatapointSubtype_DPST_UpDown_Action
	case "DPST_HVAC_PB_Action":
		return KnxDatapointSubtype_DPST_HVAC_PB_Action
	case "DPST_DoubleNibble":
		return KnxDatapointSubtype_DPST_DoubleNibble
	case "DPST_SceneInfo":
		return KnxDatapointSubtype_DPST_SceneInfo
	case "DPST_CombinedInfoOnOff":
		return KnxDatapointSubtype_DPST_CombinedInfoOnOff
	case "DPST_ActiveEnergy_V64":
		return KnxDatapointSubtype_DPST_ActiveEnergy_V64
	case "DPST_ApparantEnergy_V64":
		return KnxDatapointSubtype_DPST_ApparantEnergy_V64
	case "DPST_ReactiveEnergy_V64":
		return KnxDatapointSubtype_DPST_ReactiveEnergy_V64
	case "DPST_Channel_Activation_24":
		return KnxDatapointSubtype_DPST_Channel_Activation_24
	case "DPST_Alarm_Control":
		return KnxDatapointSubtype_DPST_Alarm_Control
	case "DPST_HVACModeNext":
		return KnxDatapointSubtype_DPST_HVACModeNext
	case "DPST_DHWModeNext":
		return KnxDatapointSubtype_DPST_DHWModeNext
	case "DPST_OccModeNext":
		return KnxDatapointSubtype_DPST_OccModeNext
	case "DPST_BuildingModeNext":
		return KnxDatapointSubtype_DPST_BuildingModeNext
	case "DPST_Version":
		return KnxDatapointSubtype_DPST_Version
	case "DPST_AlarmInfo":
		return KnxDatapointSubtype_DPST_AlarmInfo
	case "DPST_TempRoomSetpSetF16_3":
		return KnxDatapointSubtype_DPST_TempRoomSetpSetF16_3
	case "DPST_TempRoomSetpSetShiftF16_3":
		return KnxDatapointSubtype_DPST_TempRoomSetpSetShiftF16_3
	case "DPST_Scaling_Speed":
		return KnxDatapointSubtype_DPST_Scaling_Speed
	case "DPST_Scaling_Step_Time":
		return KnxDatapointSubtype_DPST_Scaling_Step_Time
	case "DPST_Enable":
		return KnxDatapointSubtype_DPST_Enable
	case "DPST_BinaryValue_Control":
		return KnxDatapointSubtype_DPST_BinaryValue_Control
	case "DPST_MeteringValue":
		return KnxDatapointSubtype_DPST_MeteringValue
	case "DPST_MBus_Address":
		return KnxDatapointSubtype_DPST_MBus_Address
	case "DPST_Colour_RGB":
		return KnxDatapointSubtype_DPST_Colour_RGB
	case "DPST_LanguageCodeAlpha2_ASCII":
		return KnxDatapointSubtype_DPST_LanguageCodeAlpha2_ASCII
	case "DPST_Tariff_ActiveEnergy":
		return KnxDatapointSubtype_DPST_Tariff_ActiveEnergy
	case "DPST_Prioritised_Mode_Control":
		return KnxDatapointSubtype_DPST_Prioritised_Mode_Control
	case "DPST_DALI_Control_Gear_Diagnostic":
		return KnxDatapointSubtype_DPST_DALI_Control_Gear_Diagnostic
	case "DPST_DALI_Diagnostics":
		return KnxDatapointSubtype_DPST_DALI_Diagnostics
	case "DPST_CombinedPosition":
		return KnxDatapointSubtype_DPST_CombinedPosition
	case "DPST_StatusSAB":
		return KnxDatapointSubtype_DPST_StatusSAB
	case "DPST_Step_Control":
		return KnxDatapointSubtype_DPST_Step_Control
	case "DPST_Colour_xyY":
		return KnxDatapointSubtype_DPST_Colour_xyY
	case "DPST_Converter_Status":
		return KnxDatapointSubtype_DPST_Converter_Status
	case "DPST_Converter_Test_Result":
		return KnxDatapointSubtype_DPST_Converter_Test_Result
	case "DPST_Battery_Info":
		return KnxDatapointSubtype_DPST_Battery_Info
	case "DPST_Brightness_Colour_Temperature_Transition":
		return KnxDatapointSubtype_DPST_Brightness_Colour_Temperature_Transition
	case "DPST_Brightness_Colour_Temperature_Control":
		return KnxDatapointSubtype_DPST_Brightness_Colour_Temperature_Control
	case "DPST_Colour_RGBW":
		return KnxDatapointSubtype_DPST_Colour_RGBW
	case "DPST_Relative_Control_RGBW":
		return KnxDatapointSubtype_DPST_Relative_Control_RGBW
	case "DPST_Relative_Control_RGB":
		return KnxDatapointSubtype_DPST_Relative_Control_RGB
	case "DPST_GeographicalLocation":
		return KnxDatapointSubtype_DPST_GeographicalLocation
	case "DPST_Direction1_Control":
		return KnxDatapointSubtype_DPST_Direction1_Control
	case "DPST_TempRoomSetpSetF16_4":
		return KnxDatapointSubtype_DPST_TempRoomSetpSetF16_4
	case "DPST_TempRoomSetpSetShiftF16_4":
		return KnxDatapointSubtype_DPST_TempRoomSetpSetShiftF16_4
	case "DPST_Direction2_Control":
		return KnxDatapointSubtype_DPST_Direction2_Control
	case "DPST_Start_Control":
		return KnxDatapointSubtype_DPST_Start_Control
	case "DPST_State_Control":
		return KnxDatapointSubtype_DPST_State_Control
	case "DPST_Invert_Control":
		return KnxDatapointSubtype_DPST_Invert_Control
	case "DPST_Control_Dimming":
		return KnxDatapointSubtype_DPST_Control_Dimming
	case "DPST_Control_Blinds":
		return KnxDatapointSubtype_DPST_Control_Blinds
	case "DPST_Char_ASCII":
		return KnxDatapointSubtype_DPST_Char_ASCII
	case "DPST_Ramp":
		return KnxDatapointSubtype_DPST_Ramp
	case "DPST_Char_8859_1":
		return KnxDatapointSubtype_DPST_Char_8859_1
	case "DPST_Scaling":
		return KnxDatapointSubtype_DPST_Scaling
	case "DPST_Angle":
		return KnxDatapointSubtype_DPST_Angle
	case "DPST_Percent_U8":
		return KnxDatapointSubtype_DPST_Percent_U8
	case "DPST_DecimalFactor":
		return KnxDatapointSubtype_DPST_DecimalFactor
	case "DPST_Tariff":
		return KnxDatapointSubtype_DPST_Tariff
	case "DPST_Value_1_Ucount":
		return KnxDatapointSubtype_DPST_Value_1_Ucount
	case "DPST_FanStage":
		return KnxDatapointSubtype_DPST_FanStage
	case "DPST_Percent_V8":
		return KnxDatapointSubtype_DPST_Percent_V8
	case "DPST_Value_1_Count":
		return KnxDatapointSubtype_DPST_Value_1_Count
	case "DPST_Alarm":
		return KnxDatapointSubtype_DPST_Alarm
	case "DPST_Status_Mode3":
		return KnxDatapointSubtype_DPST_Status_Mode3
	case "DPST_Value_2_Ucount":
		return KnxDatapointSubtype_DPST_Value_2_Ucount
	case "DPST_TimePeriodMsec":
		return KnxDatapointSubtype_DPST_TimePeriodMsec
	case "DPST_TimePeriod10Msec":
		return KnxDatapointSubtype_DPST_TimePeriod10Msec
	case "DPST_TimePeriod100Msec":
		return KnxDatapointSubtype_DPST_TimePeriod100Msec
	case "DPST_TimePeriodSec":
		return KnxDatapointSubtype_DPST_TimePeriodSec
	case "DPST_TimePeriodMin":
		return KnxDatapointSubtype_DPST_TimePeriodMin
	case "DPST_TimePeriodHrs":
		return KnxDatapointSubtype_DPST_TimePeriodHrs
	case "DPST_PropDataType":
		return KnxDatapointSubtype_DPST_PropDataType
	case "DPST_Length_mm":
		return KnxDatapointSubtype_DPST_Length_mm
	case "DPST_BinaryValue":
		return KnxDatapointSubtype_DPST_BinaryValue
	case "DPST_UElCurrentmA":
		return KnxDatapointSubtype_DPST_UElCurrentmA
	case "DPST_Brightness":
		return KnxDatapointSubtype_DPST_Brightness
	case "DPST_Absolute_Colour_Temperature":
		return KnxDatapointSubtype_DPST_Absolute_Colour_Temperature
	case "DPST_Value_2_Count":
		return KnxDatapointSubtype_DPST_Value_2_Count
	case "DPST_DeltaTimeMsec":
		return KnxDatapointSubtype_DPST_DeltaTimeMsec
	case "DPST_DeltaTime10Msec":
		return KnxDatapointSubtype_DPST_DeltaTime10Msec
	case "DPST_DeltaTime100Msec":
		return KnxDatapointSubtype_DPST_DeltaTime100Msec
	case "DPST_DeltaTimeSec":
		return KnxDatapointSubtype_DPST_DeltaTimeSec
	case "DPST_DeltaTimeMin":
		return KnxDatapointSubtype_DPST_DeltaTimeMin
	case "DPST_DeltaTimeHrs":
		return KnxDatapointSubtype_DPST_DeltaTimeHrs
	case "DPST_Step":
		return KnxDatapointSubtype_DPST_Step
	case "DPST_Percent_V16":
		return KnxDatapointSubtype_DPST_Percent_V16
	case "DPST_Rotation_Angle":
		return KnxDatapointSubtype_DPST_Rotation_Angle
	case "DPST_Length_m":
		return KnxDatapointSubtype_DPST_Length_m
	case "DPST_Value_Temp":
		return KnxDatapointSubtype_DPST_Value_Temp
	case "DPST_Value_Tempd":
		return KnxDatapointSubtype_DPST_Value_Tempd
	case "DPST_Value_Tempa":
		return KnxDatapointSubtype_DPST_Value_Tempa
	case "DPST_Value_Lux":
		return KnxDatapointSubtype_DPST_Value_Lux
	case "DPST_Value_Wsp":
		return KnxDatapointSubtype_DPST_Value_Wsp
	case "DPST_Value_Pres":
		return KnxDatapointSubtype_DPST_Value_Pres
	case "DPST_Value_Humidity":
		return KnxDatapointSubtype_DPST_Value_Humidity
	case "DPST_UpDown":
		return KnxDatapointSubtype_DPST_UpDown
	case "DPST_Value_AirQuality":
		return KnxDatapointSubtype_DPST_Value_AirQuality
	case "DPST_Value_AirFlow":
		return KnxDatapointSubtype_DPST_Value_AirFlow
	case "DPST_Value_Time1":
		return KnxDatapointSubtype_DPST_Value_Time1
	case "DPST_Value_Time2":
		return KnxDatapointSubtype_DPST_Value_Time2
	case "DPST_Value_Volt":
		return KnxDatapointSubtype_DPST_Value_Volt
	case "DPST_Value_Curr":
		return KnxDatapointSubtype_DPST_Value_Curr
	case "DPST_PowerDensity":
		return KnxDatapointSubtype_DPST_PowerDensity
	case "DPST_KelvinPerPercent":
		return KnxDatapointSubtype_DPST_KelvinPerPercent
	case "DPST_Power":
		return KnxDatapointSubtype_DPST_Power
	case "DPST_Value_Volume_Flow":
		return KnxDatapointSubtype_DPST_Value_Volume_Flow
	case "DPST_OpenClose":
		return KnxDatapointSubtype_DPST_OpenClose
	case "DPST_Rain_Amount":
		return KnxDatapointSubtype_DPST_Rain_Amount
	case "DPST_Value_Temp_F":
		return KnxDatapointSubtype_DPST_Value_Temp_F
	case "DPST_Value_Wsp_kmh":
		return KnxDatapointSubtype_DPST_Value_Wsp_kmh
	case "DPST_Value_Absolute_Humidity":
		return KnxDatapointSubtype_DPST_Value_Absolute_Humidity
	case "DPST_Concentration_ygm3":
		return KnxDatapointSubtype_DPST_Concentration_ygm3
	case "DPST_TimeOfDay":
		return KnxDatapointSubtype_DPST_TimeOfDay
	case "DPST_Date":
		return KnxDatapointSubtype_DPST_Date
	case "DPST_Value_4_Ucount":
		return KnxDatapointSubtype_DPST_Value_4_Ucount
	case "DPST_LongTimePeriod_Sec":
		return KnxDatapointSubtype_DPST_LongTimePeriod_Sec
	case "DPST_LongTimePeriod_Min":
		return KnxDatapointSubtype_DPST_LongTimePeriod_Min
	}
	return 0
}

func CastKnxDatapointSubtype(structType interface{}) KnxDatapointSubtype {
	castFunc := func(typ interface{}) KnxDatapointSubtype {
		if sKnxDatapointSubtype, ok := typ.(KnxDatapointSubtype); ok {
			return sKnxDatapointSubtype
		}
		return 0
	}
	return castFunc(structType)
}

func (m KnxDatapointSubtype) LengthInBits() uint16 {
	return 32
}

func (m KnxDatapointSubtype) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func KnxDatapointSubtypeParse(io *utils.ReadBuffer) (KnxDatapointSubtype, error) {
	val, err := io.ReadUint32(32)
	if err != nil {
		return 0, nil
	}
	return KnxDatapointSubtypeByValue(val), nil
}

func (e KnxDatapointSubtype) Serialize(io utils.WriteBuffer) error {
	err := io.WriteUint32(32, uint32(e))
	return err
}

func (e KnxDatapointSubtype) String() string {
	switch e {
	case KnxDatapointSubtype_DPST_UNKNOWN:
		return "DPST_UNKNOWN"
	case KnxDatapointSubtype_DPST_Switch:
		return "DPST_Switch"
	case KnxDatapointSubtype_DPST_Start:
		return "DPST_Start"
	case KnxDatapointSubtype_DPST_LongTimePeriod_Hrs:
		return "DPST_LongTimePeriod_Hrs"
	case KnxDatapointSubtype_DPST_VolumeLiquid_Litre:
		return "DPST_VolumeLiquid_Litre"
	case KnxDatapointSubtype_DPST_Volume_m_3:
		return "DPST_Volume_m_3"
	case KnxDatapointSubtype_DPST_Value_4_Count:
		return "DPST_Value_4_Count"
	case KnxDatapointSubtype_DPST_FlowRate_m3h:
		return "DPST_FlowRate_m3h"
	case KnxDatapointSubtype_DPST_ActiveEnergy:
		return "DPST_ActiveEnergy"
	case KnxDatapointSubtype_DPST_ApparantEnergy:
		return "DPST_ApparantEnergy"
	case KnxDatapointSubtype_DPST_ReactiveEnergy:
		return "DPST_ReactiveEnergy"
	case KnxDatapointSubtype_DPST_ActiveEnergy_kWh:
		return "DPST_ActiveEnergy_kWh"
	case KnxDatapointSubtype_DPST_ApparantEnergy_kVAh:
		return "DPST_ApparantEnergy_kVAh"
	case KnxDatapointSubtype_DPST_State:
		return "DPST_State"
	case KnxDatapointSubtype_DPST_ReactiveEnergy_kVARh:
		return "DPST_ReactiveEnergy_kVARh"
	case KnxDatapointSubtype_DPST_ActiveEnergy_MWh:
		return "DPST_ActiveEnergy_MWh"
	case KnxDatapointSubtype_DPST_LongDeltaTimeSec:
		return "DPST_LongDeltaTimeSec"
	case KnxDatapointSubtype_DPST_DeltaVolumeLiquid_Litre:
		return "DPST_DeltaVolumeLiquid_Litre"
	case KnxDatapointSubtype_DPST_DeltaVolume_m_3:
		return "DPST_DeltaVolume_m_3"
	case KnxDatapointSubtype_DPST_Value_Acceleration:
		return "DPST_Value_Acceleration"
	case KnxDatapointSubtype_DPST_Value_Acceleration_Angular:
		return "DPST_Value_Acceleration_Angular"
	case KnxDatapointSubtype_DPST_Value_Activation_Energy:
		return "DPST_Value_Activation_Energy"
	case KnxDatapointSubtype_DPST_Value_Activity:
		return "DPST_Value_Activity"
	case KnxDatapointSubtype_DPST_Value_Mol:
		return "DPST_Value_Mol"
	case KnxDatapointSubtype_DPST_Invert:
		return "DPST_Invert"
	case KnxDatapointSubtype_DPST_Value_Amplitude:
		return "DPST_Value_Amplitude"
	case KnxDatapointSubtype_DPST_Value_AngleRad:
		return "DPST_Value_AngleRad"
	case KnxDatapointSubtype_DPST_Value_AngleDeg:
		return "DPST_Value_AngleDeg"
	case KnxDatapointSubtype_DPST_Value_Angular_Momentum:
		return "DPST_Value_Angular_Momentum"
	case KnxDatapointSubtype_DPST_Value_Angular_Velocity:
		return "DPST_Value_Angular_Velocity"
	case KnxDatapointSubtype_DPST_Value_Area:
		return "DPST_Value_Area"
	case KnxDatapointSubtype_DPST_Value_Capacitance:
		return "DPST_Value_Capacitance"
	case KnxDatapointSubtype_DPST_Value_Charge_DensitySurface:
		return "DPST_Value_Charge_DensitySurface"
	case KnxDatapointSubtype_DPST_Value_Charge_DensityVolume:
		return "DPST_Value_Charge_DensityVolume"
	case KnxDatapointSubtype_DPST_Value_Compressibility:
		return "DPST_Value_Compressibility"
	case KnxDatapointSubtype_DPST_DimSendStyle:
		return "DPST_DimSendStyle"
	case KnxDatapointSubtype_DPST_Value_Conductance:
		return "DPST_Value_Conductance"
	case KnxDatapointSubtype_DPST_Value_Electrical_Conductivity:
		return "DPST_Value_Electrical_Conductivity"
	case KnxDatapointSubtype_DPST_Value_Density:
		return "DPST_Value_Density"
	case KnxDatapointSubtype_DPST_Value_Electric_Charge:
		return "DPST_Value_Electric_Charge"
	case KnxDatapointSubtype_DPST_Value_Electric_Current:
		return "DPST_Value_Electric_Current"
	case KnxDatapointSubtype_DPST_Value_Electric_CurrentDensity:
		return "DPST_Value_Electric_CurrentDensity"
	case KnxDatapointSubtype_DPST_Value_Electric_DipoleMoment:
		return "DPST_Value_Electric_DipoleMoment"
	case KnxDatapointSubtype_DPST_Value_Electric_Displacement:
		return "DPST_Value_Electric_Displacement"
	case KnxDatapointSubtype_DPST_Value_Electric_FieldStrength:
		return "DPST_Value_Electric_FieldStrength"
	case KnxDatapointSubtype_DPST_Value_Electric_Flux:
		return "DPST_Value_Electric_Flux"
	case KnxDatapointSubtype_DPST_InputSource:
		return "DPST_InputSource"
	case KnxDatapointSubtype_DPST_Value_Electric_FluxDensity:
		return "DPST_Value_Electric_FluxDensity"
	case KnxDatapointSubtype_DPST_Value_Electric_Polarization:
		return "DPST_Value_Electric_Polarization"
	case KnxDatapointSubtype_DPST_Value_Electric_Potential:
		return "DPST_Value_Electric_Potential"
	case KnxDatapointSubtype_DPST_Value_Electric_PotentialDifference:
		return "DPST_Value_Electric_PotentialDifference"
	case KnxDatapointSubtype_DPST_Value_ElectromagneticMoment:
		return "DPST_Value_ElectromagneticMoment"
	case KnxDatapointSubtype_DPST_Value_Electromotive_Force:
		return "DPST_Value_Electromotive_Force"
	case KnxDatapointSubtype_DPST_Value_Energy:
		return "DPST_Value_Energy"
	case KnxDatapointSubtype_DPST_Value_Force:
		return "DPST_Value_Force"
	case KnxDatapointSubtype_DPST_Value_Frequency:
		return "DPST_Value_Frequency"
	case KnxDatapointSubtype_DPST_Value_Angular_Frequency:
		return "DPST_Value_Angular_Frequency"
	case KnxDatapointSubtype_DPST_Reset:
		return "DPST_Reset"
	case KnxDatapointSubtype_DPST_Value_Heat_Capacity:
		return "DPST_Value_Heat_Capacity"
	case KnxDatapointSubtype_DPST_Value_Heat_FlowRate:
		return "DPST_Value_Heat_FlowRate"
	case KnxDatapointSubtype_DPST_Value_Heat_Quantity:
		return "DPST_Value_Heat_Quantity"
	case KnxDatapointSubtype_DPST_Value_Impedance:
		return "DPST_Value_Impedance"
	case KnxDatapointSubtype_DPST_Value_Length:
		return "DPST_Value_Length"
	case KnxDatapointSubtype_DPST_Value_Light_Quantity:
		return "DPST_Value_Light_Quantity"
	case KnxDatapointSubtype_DPST_Value_Luminance:
		return "DPST_Value_Luminance"
	case KnxDatapointSubtype_DPST_Value_Luminous_Flux:
		return "DPST_Value_Luminous_Flux"
	case KnxDatapointSubtype_DPST_Value_Luminous_Intensity:
		return "DPST_Value_Luminous_Intensity"
	case KnxDatapointSubtype_DPST_Value_Magnetic_FieldStrength:
		return "DPST_Value_Magnetic_FieldStrength"
	case KnxDatapointSubtype_DPST_Ack:
		return "DPST_Ack"
	case KnxDatapointSubtype_DPST_Value_Magnetic_Flux:
		return "DPST_Value_Magnetic_Flux"
	case KnxDatapointSubtype_DPST_Value_Magnetic_FluxDensity:
		return "DPST_Value_Magnetic_FluxDensity"
	case KnxDatapointSubtype_DPST_Value_Magnetic_Moment:
		return "DPST_Value_Magnetic_Moment"
	case KnxDatapointSubtype_DPST_Value_Magnetic_Polarization:
		return "DPST_Value_Magnetic_Polarization"
	case KnxDatapointSubtype_DPST_Value_Magnetization:
		return "DPST_Value_Magnetization"
	case KnxDatapointSubtype_DPST_Value_MagnetomotiveForce:
		return "DPST_Value_MagnetomotiveForce"
	case KnxDatapointSubtype_DPST_Value_Mass:
		return "DPST_Value_Mass"
	case KnxDatapointSubtype_DPST_Value_MassFlux:
		return "DPST_Value_MassFlux"
	case KnxDatapointSubtype_DPST_Value_Momentum:
		return "DPST_Value_Momentum"
	case KnxDatapointSubtype_DPST_Value_Phase_AngleRad:
		return "DPST_Value_Phase_AngleRad"
	case KnxDatapointSubtype_DPST_Trigger:
		return "DPST_Trigger"
	case KnxDatapointSubtype_DPST_Value_Phase_AngleDeg:
		return "DPST_Value_Phase_AngleDeg"
	case KnxDatapointSubtype_DPST_Value_Power:
		return "DPST_Value_Power"
	case KnxDatapointSubtype_DPST_Value_Power_Factor:
		return "DPST_Value_Power_Factor"
	case KnxDatapointSubtype_DPST_Value_Pressure:
		return "DPST_Value_Pressure"
	case KnxDatapointSubtype_DPST_Value_Reactance:
		return "DPST_Value_Reactance"
	case KnxDatapointSubtype_DPST_Value_Resistance:
		return "DPST_Value_Resistance"
	case KnxDatapointSubtype_DPST_Value_Resistivity:
		return "DPST_Value_Resistivity"
	case KnxDatapointSubtype_DPST_Value_SelfInductance:
		return "DPST_Value_SelfInductance"
	case KnxDatapointSubtype_DPST_Value_SolidAngle:
		return "DPST_Value_SolidAngle"
	case KnxDatapointSubtype_DPST_Value_Sound_Intensity:
		return "DPST_Value_Sound_Intensity"
	case KnxDatapointSubtype_DPST_Occupancy:
		return "DPST_Occupancy"
	case KnxDatapointSubtype_DPST_Value_Speed:
		return "DPST_Value_Speed"
	case KnxDatapointSubtype_DPST_Value_Stress:
		return "DPST_Value_Stress"
	case KnxDatapointSubtype_DPST_Value_Surface_Tension:
		return "DPST_Value_Surface_Tension"
	case KnxDatapointSubtype_DPST_Value_Common_Temperature:
		return "DPST_Value_Common_Temperature"
	case KnxDatapointSubtype_DPST_Value_Absolute_Temperature:
		return "DPST_Value_Absolute_Temperature"
	case KnxDatapointSubtype_DPST_Value_TemperatureDifference:
		return "DPST_Value_TemperatureDifference"
	case KnxDatapointSubtype_DPST_Value_Thermal_Capacity:
		return "DPST_Value_Thermal_Capacity"
	case KnxDatapointSubtype_DPST_Value_Thermal_Conductivity:
		return "DPST_Value_Thermal_Conductivity"
	case KnxDatapointSubtype_DPST_Value_ThermoelectricPower:
		return "DPST_Value_ThermoelectricPower"
	case KnxDatapointSubtype_DPST_Value_Time:
		return "DPST_Value_Time"
	case KnxDatapointSubtype_DPST_Window_Door:
		return "DPST_Window_Door"
	case KnxDatapointSubtype_DPST_Value_Torque:
		return "DPST_Value_Torque"
	case KnxDatapointSubtype_DPST_Value_Volume:
		return "DPST_Value_Volume"
	case KnxDatapointSubtype_DPST_Value_Volume_Flux:
		return "DPST_Value_Volume_Flux"
	case KnxDatapointSubtype_DPST_Value_Weight:
		return "DPST_Value_Weight"
	case KnxDatapointSubtype_DPST_Value_Work:
		return "DPST_Value_Work"
	case KnxDatapointSubtype_DPST_Volume_Flux_Meter:
		return "DPST_Volume_Flux_Meter"
	case KnxDatapointSubtype_DPST_Volume_Flux_ls:
		return "DPST_Volume_Flux_ls"
	case KnxDatapointSubtype_DPST_Access_Data:
		return "DPST_Access_Data"
	case KnxDatapointSubtype_DPST_String_ASCII:
		return "DPST_String_ASCII"
	case KnxDatapointSubtype_DPST_String_8859_1:
		return "DPST_String_8859_1"
	case KnxDatapointSubtype_DPST_Bool:
		return "DPST_Bool"
	case KnxDatapointSubtype_DPST_LogicalFunction:
		return "DPST_LogicalFunction"
	case KnxDatapointSubtype_DPST_SceneNumber:
		return "DPST_SceneNumber"
	case KnxDatapointSubtype_DPST_SceneControl:
		return "DPST_SceneControl"
	case KnxDatapointSubtype_DPST_DateTime:
		return "DPST_DateTime"
	case KnxDatapointSubtype_DPST_SCLOMode:
		return "DPST_SCLOMode"
	case KnxDatapointSubtype_DPST_BuildingMode:
		return "DPST_BuildingMode"
	case KnxDatapointSubtype_DPST_OccMode:
		return "DPST_OccMode"
	case KnxDatapointSubtype_DPST_Priority:
		return "DPST_Priority"
	case KnxDatapointSubtype_DPST_LightApplicationMode:
		return "DPST_LightApplicationMode"
	case KnxDatapointSubtype_DPST_ApplicationArea:
		return "DPST_ApplicationArea"
	case KnxDatapointSubtype_DPST_AlarmClassType:
		return "DPST_AlarmClassType"
	case KnxDatapointSubtype_DPST_Scene_AB:
		return "DPST_Scene_AB"
	case KnxDatapointSubtype_DPST_PSUMode:
		return "DPST_PSUMode"
	case KnxDatapointSubtype_DPST_ErrorClass_System:
		return "DPST_ErrorClass_System"
	case KnxDatapointSubtype_DPST_ErrorClass_HVAC:
		return "DPST_ErrorClass_HVAC"
	case KnxDatapointSubtype_DPST_Time_Delay:
		return "DPST_Time_Delay"
	case KnxDatapointSubtype_DPST_Beaufort_Wind_Force_Scale:
		return "DPST_Beaufort_Wind_Force_Scale"
	case KnxDatapointSubtype_DPST_SensorSelect:
		return "DPST_SensorSelect"
	case KnxDatapointSubtype_DPST_ActuatorConnectType:
		return "DPST_ActuatorConnectType"
	case KnxDatapointSubtype_DPST_Cloud_Cover:
		return "DPST_Cloud_Cover"
	case KnxDatapointSubtype_DPST_PowerReturnMode:
		return "DPST_PowerReturnMode"
	case KnxDatapointSubtype_DPST_FuelType:
		return "DPST_FuelType"
	case KnxDatapointSubtype_DPST_ShutterBlinds_Mode:
		return "DPST_ShutterBlinds_Mode"
	case KnxDatapointSubtype_DPST_BurnerType:
		return "DPST_BurnerType"
	case KnxDatapointSubtype_DPST_HVACMode:
		return "DPST_HVACMode"
	case KnxDatapointSubtype_DPST_DHWMode:
		return "DPST_DHWMode"
	case KnxDatapointSubtype_DPST_LoadPriority:
		return "DPST_LoadPriority"
	case KnxDatapointSubtype_DPST_HVACContrMode:
		return "DPST_HVACContrMode"
	case KnxDatapointSubtype_DPST_HVACEmergMode:
		return "DPST_HVACEmergMode"
	case KnxDatapointSubtype_DPST_ChangeoverMode:
		return "DPST_ChangeoverMode"
	case KnxDatapointSubtype_DPST_ValveMode:
		return "DPST_ValveMode"
	case KnxDatapointSubtype_DPST_DamperMode:
		return "DPST_DamperMode"
	case KnxDatapointSubtype_DPST_HeaterMode:
		return "DPST_HeaterMode"
	case KnxDatapointSubtype_DPST_DayNight:
		return "DPST_DayNight"
	case KnxDatapointSubtype_DPST_FanMode:
		return "DPST_FanMode"
	case KnxDatapointSubtype_DPST_MasterSlaveMode:
		return "DPST_MasterSlaveMode"
	case KnxDatapointSubtype_DPST_StatusRoomSetp:
		return "DPST_StatusRoomSetp"
	case KnxDatapointSubtype_DPST_Metering_DeviceType:
		return "DPST_Metering_DeviceType"
	case KnxDatapointSubtype_DPST_HumDehumMode:
		return "DPST_HumDehumMode"
	case KnxDatapointSubtype_DPST_EnableHCStage:
		return "DPST_EnableHCStage"
	case KnxDatapointSubtype_DPST_ADAType:
		return "DPST_ADAType"
	case KnxDatapointSubtype_DPST_BackupMode:
		return "DPST_BackupMode"
	case KnxDatapointSubtype_DPST_StartSynchronization:
		return "DPST_StartSynchronization"
	case KnxDatapointSubtype_DPST_Behaviour_Lock_Unlock:
		return "DPST_Behaviour_Lock_Unlock"
	case KnxDatapointSubtype_DPST_Heat_Cool:
		return "DPST_Heat_Cool"
	case KnxDatapointSubtype_DPST_Behaviour_Bus_Power_Up_Down:
		return "DPST_Behaviour_Bus_Power_Up_Down"
	case KnxDatapointSubtype_DPST_DALI_Fade_Time:
		return "DPST_DALI_Fade_Time"
	case KnxDatapointSubtype_DPST_BlinkingMode:
		return "DPST_BlinkingMode"
	case KnxDatapointSubtype_DPST_LightControlMode:
		return "DPST_LightControlMode"
	case KnxDatapointSubtype_DPST_SwitchPBModel:
		return "DPST_SwitchPBModel"
	case KnxDatapointSubtype_DPST_PBAction:
		return "DPST_PBAction"
	case KnxDatapointSubtype_DPST_DimmPBModel:
		return "DPST_DimmPBModel"
	case KnxDatapointSubtype_DPST_SwitchOnMode:
		return "DPST_SwitchOnMode"
	case KnxDatapointSubtype_DPST_LoadTypeSet:
		return "DPST_LoadTypeSet"
	case KnxDatapointSubtype_DPST_LoadTypeDetected:
		return "DPST_LoadTypeDetected"
	case KnxDatapointSubtype_DPST_Switch_Control:
		return "DPST_Switch_Control"
	case KnxDatapointSubtype_DPST_Converter_Test_Control:
		return "DPST_Converter_Test_Control"
	case KnxDatapointSubtype_DPST_SABExcept_Behaviour:
		return "DPST_SABExcept_Behaviour"
	case KnxDatapointSubtype_DPST_SABBehaviour_Lock_Unlock:
		return "DPST_SABBehaviour_Lock_Unlock"
	case KnxDatapointSubtype_DPST_SSSBMode:
		return "DPST_SSSBMode"
	case KnxDatapointSubtype_DPST_BlindsControlMode:
		return "DPST_BlindsControlMode"
	case KnxDatapointSubtype_DPST_CommMode:
		return "DPST_CommMode"
	case KnxDatapointSubtype_DPST_AddInfoTypes:
		return "DPST_AddInfoTypes"
	case KnxDatapointSubtype_DPST_RF_ModeSelect:
		return "DPST_RF_ModeSelect"
	case KnxDatapointSubtype_DPST_RF_FilterSelect:
		return "DPST_RF_FilterSelect"
	case KnxDatapointSubtype_DPST_StatusGen:
		return "DPST_StatusGen"
	case KnxDatapointSubtype_DPST_Bool_Control:
		return "DPST_Bool_Control"
	case KnxDatapointSubtype_DPST_Device_Control:
		return "DPST_Device_Control"
	case KnxDatapointSubtype_DPST_ForceSign:
		return "DPST_ForceSign"
	case KnxDatapointSubtype_DPST_ForceSignCool:
		return "DPST_ForceSignCool"
	case KnxDatapointSubtype_DPST_StatusRHC:
		return "DPST_StatusRHC"
	case KnxDatapointSubtype_DPST_StatusSDHWC:
		return "DPST_StatusSDHWC"
	case KnxDatapointSubtype_DPST_FuelTypeSet:
		return "DPST_FuelTypeSet"
	case KnxDatapointSubtype_DPST_StatusRCC:
		return "DPST_StatusRCC"
	case KnxDatapointSubtype_DPST_StatusAHU:
		return "DPST_StatusAHU"
	case KnxDatapointSubtype_DPST_CombinedStatus_RTSM:
		return "DPST_CombinedStatus_RTSM"
	case KnxDatapointSubtype_DPST_LightActuatorErrorInfo:
		return "DPST_LightActuatorErrorInfo"
	case KnxDatapointSubtype_DPST_Enable_Control:
		return "DPST_Enable_Control"
	case KnxDatapointSubtype_DPST_RF_ModeInfo:
		return "DPST_RF_ModeInfo"
	case KnxDatapointSubtype_DPST_RF_FilterInfo:
		return "DPST_RF_FilterInfo"
	case KnxDatapointSubtype_DPST_Channel_Activation_8:
		return "DPST_Channel_Activation_8"
	case KnxDatapointSubtype_DPST_StatusDHWC:
		return "DPST_StatusDHWC"
	case KnxDatapointSubtype_DPST_StatusRHCC:
		return "DPST_StatusRHCC"
	case KnxDatapointSubtype_DPST_CombinedStatus_HVA:
		return "DPST_CombinedStatus_HVA"
	case KnxDatapointSubtype_DPST_CombinedStatus_RTC:
		return "DPST_CombinedStatus_RTC"
	case KnxDatapointSubtype_DPST_Media:
		return "DPST_Media"
	case KnxDatapointSubtype_DPST_Channel_Activation_16:
		return "DPST_Channel_Activation_16"
	case KnxDatapointSubtype_DPST_OnOffAction:
		return "DPST_OnOffAction"
	case KnxDatapointSubtype_DPST_Ramp_Control:
		return "DPST_Ramp_Control"
	case KnxDatapointSubtype_DPST_Alarm_Reaction:
		return "DPST_Alarm_Reaction"
	case KnxDatapointSubtype_DPST_UpDown_Action:
		return "DPST_UpDown_Action"
	case KnxDatapointSubtype_DPST_HVAC_PB_Action:
		return "DPST_HVAC_PB_Action"
	case KnxDatapointSubtype_DPST_DoubleNibble:
		return "DPST_DoubleNibble"
	case KnxDatapointSubtype_DPST_SceneInfo:
		return "DPST_SceneInfo"
	case KnxDatapointSubtype_DPST_CombinedInfoOnOff:
		return "DPST_CombinedInfoOnOff"
	case KnxDatapointSubtype_DPST_ActiveEnergy_V64:
		return "DPST_ActiveEnergy_V64"
	case KnxDatapointSubtype_DPST_ApparantEnergy_V64:
		return "DPST_ApparantEnergy_V64"
	case KnxDatapointSubtype_DPST_ReactiveEnergy_V64:
		return "DPST_ReactiveEnergy_V64"
	case KnxDatapointSubtype_DPST_Channel_Activation_24:
		return "DPST_Channel_Activation_24"
	case KnxDatapointSubtype_DPST_Alarm_Control:
		return "DPST_Alarm_Control"
	case KnxDatapointSubtype_DPST_HVACModeNext:
		return "DPST_HVACModeNext"
	case KnxDatapointSubtype_DPST_DHWModeNext:
		return "DPST_DHWModeNext"
	case KnxDatapointSubtype_DPST_OccModeNext:
		return "DPST_OccModeNext"
	case KnxDatapointSubtype_DPST_BuildingModeNext:
		return "DPST_BuildingModeNext"
	case KnxDatapointSubtype_DPST_Version:
		return "DPST_Version"
	case KnxDatapointSubtype_DPST_AlarmInfo:
		return "DPST_AlarmInfo"
	case KnxDatapointSubtype_DPST_TempRoomSetpSetF16_3:
		return "DPST_TempRoomSetpSetF16_3"
	case KnxDatapointSubtype_DPST_TempRoomSetpSetShiftF16_3:
		return "DPST_TempRoomSetpSetShiftF16_3"
	case KnxDatapointSubtype_DPST_Scaling_Speed:
		return "DPST_Scaling_Speed"
	case KnxDatapointSubtype_DPST_Scaling_Step_Time:
		return "DPST_Scaling_Step_Time"
	case KnxDatapointSubtype_DPST_Enable:
		return "DPST_Enable"
	case KnxDatapointSubtype_DPST_BinaryValue_Control:
		return "DPST_BinaryValue_Control"
	case KnxDatapointSubtype_DPST_MeteringValue:
		return "DPST_MeteringValue"
	case KnxDatapointSubtype_DPST_MBus_Address:
		return "DPST_MBus_Address"
	case KnxDatapointSubtype_DPST_Colour_RGB:
		return "DPST_Colour_RGB"
	case KnxDatapointSubtype_DPST_LanguageCodeAlpha2_ASCII:
		return "DPST_LanguageCodeAlpha2_ASCII"
	case KnxDatapointSubtype_DPST_Tariff_ActiveEnergy:
		return "DPST_Tariff_ActiveEnergy"
	case KnxDatapointSubtype_DPST_Prioritised_Mode_Control:
		return "DPST_Prioritised_Mode_Control"
	case KnxDatapointSubtype_DPST_DALI_Control_Gear_Diagnostic:
		return "DPST_DALI_Control_Gear_Diagnostic"
	case KnxDatapointSubtype_DPST_DALI_Diagnostics:
		return "DPST_DALI_Diagnostics"
	case KnxDatapointSubtype_DPST_CombinedPosition:
		return "DPST_CombinedPosition"
	case KnxDatapointSubtype_DPST_StatusSAB:
		return "DPST_StatusSAB"
	case KnxDatapointSubtype_DPST_Step_Control:
		return "DPST_Step_Control"
	case KnxDatapointSubtype_DPST_Colour_xyY:
		return "DPST_Colour_xyY"
	case KnxDatapointSubtype_DPST_Converter_Status:
		return "DPST_Converter_Status"
	case KnxDatapointSubtype_DPST_Converter_Test_Result:
		return "DPST_Converter_Test_Result"
	case KnxDatapointSubtype_DPST_Battery_Info:
		return "DPST_Battery_Info"
	case KnxDatapointSubtype_DPST_Brightness_Colour_Temperature_Transition:
		return "DPST_Brightness_Colour_Temperature_Transition"
	case KnxDatapointSubtype_DPST_Brightness_Colour_Temperature_Control:
		return "DPST_Brightness_Colour_Temperature_Control"
	case KnxDatapointSubtype_DPST_Colour_RGBW:
		return "DPST_Colour_RGBW"
	case KnxDatapointSubtype_DPST_Relative_Control_RGBW:
		return "DPST_Relative_Control_RGBW"
	case KnxDatapointSubtype_DPST_Relative_Control_RGB:
		return "DPST_Relative_Control_RGB"
	case KnxDatapointSubtype_DPST_GeographicalLocation:
		return "DPST_GeographicalLocation"
	case KnxDatapointSubtype_DPST_Direction1_Control:
		return "DPST_Direction1_Control"
	case KnxDatapointSubtype_DPST_TempRoomSetpSetF16_4:
		return "DPST_TempRoomSetpSetF16_4"
	case KnxDatapointSubtype_DPST_TempRoomSetpSetShiftF16_4:
		return "DPST_TempRoomSetpSetShiftF16_4"
	case KnxDatapointSubtype_DPST_Direction2_Control:
		return "DPST_Direction2_Control"
	case KnxDatapointSubtype_DPST_Start_Control:
		return "DPST_Start_Control"
	case KnxDatapointSubtype_DPST_State_Control:
		return "DPST_State_Control"
	case KnxDatapointSubtype_DPST_Invert_Control:
		return "DPST_Invert_Control"
	case KnxDatapointSubtype_DPST_Control_Dimming:
		return "DPST_Control_Dimming"
	case KnxDatapointSubtype_DPST_Control_Blinds:
		return "DPST_Control_Blinds"
	case KnxDatapointSubtype_DPST_Char_ASCII:
		return "DPST_Char_ASCII"
	case KnxDatapointSubtype_DPST_Ramp:
		return "DPST_Ramp"
	case KnxDatapointSubtype_DPST_Char_8859_1:
		return "DPST_Char_8859_1"
	case KnxDatapointSubtype_DPST_Scaling:
		return "DPST_Scaling"
	case KnxDatapointSubtype_DPST_Angle:
		return "DPST_Angle"
	case KnxDatapointSubtype_DPST_Percent_U8:
		return "DPST_Percent_U8"
	case KnxDatapointSubtype_DPST_DecimalFactor:
		return "DPST_DecimalFactor"
	case KnxDatapointSubtype_DPST_Tariff:
		return "DPST_Tariff"
	case KnxDatapointSubtype_DPST_Value_1_Ucount:
		return "DPST_Value_1_Ucount"
	case KnxDatapointSubtype_DPST_FanStage:
		return "DPST_FanStage"
	case KnxDatapointSubtype_DPST_Percent_V8:
		return "DPST_Percent_V8"
	case KnxDatapointSubtype_DPST_Value_1_Count:
		return "DPST_Value_1_Count"
	case KnxDatapointSubtype_DPST_Alarm:
		return "DPST_Alarm"
	case KnxDatapointSubtype_DPST_Status_Mode3:
		return "DPST_Status_Mode3"
	case KnxDatapointSubtype_DPST_Value_2_Ucount:
		return "DPST_Value_2_Ucount"
	case KnxDatapointSubtype_DPST_TimePeriodMsec:
		return "DPST_TimePeriodMsec"
	case KnxDatapointSubtype_DPST_TimePeriod10Msec:
		return "DPST_TimePeriod10Msec"
	case KnxDatapointSubtype_DPST_TimePeriod100Msec:
		return "DPST_TimePeriod100Msec"
	case KnxDatapointSubtype_DPST_TimePeriodSec:
		return "DPST_TimePeriodSec"
	case KnxDatapointSubtype_DPST_TimePeriodMin:
		return "DPST_TimePeriodMin"
	case KnxDatapointSubtype_DPST_TimePeriodHrs:
		return "DPST_TimePeriodHrs"
	case KnxDatapointSubtype_DPST_PropDataType:
		return "DPST_PropDataType"
	case KnxDatapointSubtype_DPST_Length_mm:
		return "DPST_Length_mm"
	case KnxDatapointSubtype_DPST_BinaryValue:
		return "DPST_BinaryValue"
	case KnxDatapointSubtype_DPST_UElCurrentmA:
		return "DPST_UElCurrentmA"
	case KnxDatapointSubtype_DPST_Brightness:
		return "DPST_Brightness"
	case KnxDatapointSubtype_DPST_Absolute_Colour_Temperature:
		return "DPST_Absolute_Colour_Temperature"
	case KnxDatapointSubtype_DPST_Value_2_Count:
		return "DPST_Value_2_Count"
	case KnxDatapointSubtype_DPST_DeltaTimeMsec:
		return "DPST_DeltaTimeMsec"
	case KnxDatapointSubtype_DPST_DeltaTime10Msec:
		return "DPST_DeltaTime10Msec"
	case KnxDatapointSubtype_DPST_DeltaTime100Msec:
		return "DPST_DeltaTime100Msec"
	case KnxDatapointSubtype_DPST_DeltaTimeSec:
		return "DPST_DeltaTimeSec"
	case KnxDatapointSubtype_DPST_DeltaTimeMin:
		return "DPST_DeltaTimeMin"
	case KnxDatapointSubtype_DPST_DeltaTimeHrs:
		return "DPST_DeltaTimeHrs"
	case KnxDatapointSubtype_DPST_Step:
		return "DPST_Step"
	case KnxDatapointSubtype_DPST_Percent_V16:
		return "DPST_Percent_V16"
	case KnxDatapointSubtype_DPST_Rotation_Angle:
		return "DPST_Rotation_Angle"
	case KnxDatapointSubtype_DPST_Length_m:
		return "DPST_Length_m"
	case KnxDatapointSubtype_DPST_Value_Temp:
		return "DPST_Value_Temp"
	case KnxDatapointSubtype_DPST_Value_Tempd:
		return "DPST_Value_Tempd"
	case KnxDatapointSubtype_DPST_Value_Tempa:
		return "DPST_Value_Tempa"
	case KnxDatapointSubtype_DPST_Value_Lux:
		return "DPST_Value_Lux"
	case KnxDatapointSubtype_DPST_Value_Wsp:
		return "DPST_Value_Wsp"
	case KnxDatapointSubtype_DPST_Value_Pres:
		return "DPST_Value_Pres"
	case KnxDatapointSubtype_DPST_Value_Humidity:
		return "DPST_Value_Humidity"
	case KnxDatapointSubtype_DPST_UpDown:
		return "DPST_UpDown"
	case KnxDatapointSubtype_DPST_Value_AirQuality:
		return "DPST_Value_AirQuality"
	case KnxDatapointSubtype_DPST_Value_AirFlow:
		return "DPST_Value_AirFlow"
	case KnxDatapointSubtype_DPST_Value_Time1:
		return "DPST_Value_Time1"
	case KnxDatapointSubtype_DPST_Value_Time2:
		return "DPST_Value_Time2"
	case KnxDatapointSubtype_DPST_Value_Volt:
		return "DPST_Value_Volt"
	case KnxDatapointSubtype_DPST_Value_Curr:
		return "DPST_Value_Curr"
	case KnxDatapointSubtype_DPST_PowerDensity:
		return "DPST_PowerDensity"
	case KnxDatapointSubtype_DPST_KelvinPerPercent:
		return "DPST_KelvinPerPercent"
	case KnxDatapointSubtype_DPST_Power:
		return "DPST_Power"
	case KnxDatapointSubtype_DPST_Value_Volume_Flow:
		return "DPST_Value_Volume_Flow"
	case KnxDatapointSubtype_DPST_OpenClose:
		return "DPST_OpenClose"
	case KnxDatapointSubtype_DPST_Rain_Amount:
		return "DPST_Rain_Amount"
	case KnxDatapointSubtype_DPST_Value_Temp_F:
		return "DPST_Value_Temp_F"
	case KnxDatapointSubtype_DPST_Value_Wsp_kmh:
		return "DPST_Value_Wsp_kmh"
	case KnxDatapointSubtype_DPST_Value_Absolute_Humidity:
		return "DPST_Value_Absolute_Humidity"
	case KnxDatapointSubtype_DPST_Concentration_ygm3:
		return "DPST_Concentration_ygm3"
	case KnxDatapointSubtype_DPST_TimeOfDay:
		return "DPST_TimeOfDay"
	case KnxDatapointSubtype_DPST_Date:
		return "DPST_Date"
	case KnxDatapointSubtype_DPST_Value_4_Ucount:
		return "DPST_Value_4_Ucount"
	case KnxDatapointSubtype_DPST_LongTimePeriod_Sec:
		return "DPST_LongTimePeriod_Sec"
	case KnxDatapointSubtype_DPST_LongTimePeriod_Min:
		return "DPST_LongTimePeriod_Min"
	}
	return ""
}
