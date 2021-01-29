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

type KnxDatapointType uint32

type IKnxDatapointType interface {
    Number() uint16
    Name() string
    DatapointMainType() KnxDatapointMainType
    Serialize(io utils.WriteBuffer) error
}

const(
    KnxDatapointType_DPT_UNKNOWN KnxDatapointType = 0
    KnxDatapointType_DPT_Switch KnxDatapointType = 1
    KnxDatapointType_DPT_Bool KnxDatapointType = 2
    KnxDatapointType_DPT_Enable KnxDatapointType = 3
    KnxDatapointType_DPT_Ramp KnxDatapointType = 4
    KnxDatapointType_DPT_Alarm KnxDatapointType = 5
    KnxDatapointType_DPT_BinaryValue KnxDatapointType = 6
    KnxDatapointType_DPT_Step KnxDatapointType = 7
    KnxDatapointType_DPT_UpDown KnxDatapointType = 8
    KnxDatapointType_DPT_OpenClose KnxDatapointType = 9
    KnxDatapointType_DPT_Start KnxDatapointType = 10
    KnxDatapointType_DPT_State KnxDatapointType = 11
    KnxDatapointType_DPT_Invert KnxDatapointType = 12
    KnxDatapointType_DPT_DimSendStyle KnxDatapointType = 13
    KnxDatapointType_DPT_InputSource KnxDatapointType = 14
    KnxDatapointType_DPT_Reset KnxDatapointType = 15
    KnxDatapointType_DPT_Ack KnxDatapointType = 16
    KnxDatapointType_DPT_Trigger KnxDatapointType = 17
    KnxDatapointType_DPT_Occupancy KnxDatapointType = 18
    KnxDatapointType_DPT_Window_Door KnxDatapointType = 19
    KnxDatapointType_DPT_LogicalFunction KnxDatapointType = 20
    KnxDatapointType_DPT_Scene_AB KnxDatapointType = 21
    KnxDatapointType_DPT_ShutterBlinds_Mode KnxDatapointType = 22
    KnxDatapointType_DPT_DayNight KnxDatapointType = 23
    KnxDatapointType_DPT_Heat_Cool KnxDatapointType = 24
    KnxDatapointType_DPT_Switch_Control KnxDatapointType = 25
    KnxDatapointType_DPT_Bool_Control KnxDatapointType = 26
    KnxDatapointType_DPT_Enable_Control KnxDatapointType = 27
    KnxDatapointType_DPT_Ramp_Control KnxDatapointType = 28
    KnxDatapointType_DPT_Alarm_Control KnxDatapointType = 29
    KnxDatapointType_DPT_BinaryValue_Control KnxDatapointType = 30
    KnxDatapointType_DPT_Step_Control KnxDatapointType = 31
    KnxDatapointType_DPT_Direction1_Control KnxDatapointType = 32
    KnxDatapointType_DPT_Direction2_Control KnxDatapointType = 33
    KnxDatapointType_DPT_Start_Control KnxDatapointType = 34
    KnxDatapointType_DPT_State_Control KnxDatapointType = 35
    KnxDatapointType_DPT_Invert_Control KnxDatapointType = 36
    KnxDatapointType_DPT_Control_Dimming KnxDatapointType = 37
    KnxDatapointType_DPT_Control_Blinds KnxDatapointType = 38
    KnxDatapointType_DPT_Char_ASCII KnxDatapointType = 39
    KnxDatapointType_DPT_Char_8859_1 KnxDatapointType = 40
    KnxDatapointType_DPT_Scaling KnxDatapointType = 41
    KnxDatapointType_DPT_Angle KnxDatapointType = 42
    KnxDatapointType_DPT_Percent_U8 KnxDatapointType = 43
    KnxDatapointType_DPT_DecimalFactor KnxDatapointType = 44
    KnxDatapointType_DPT_Tariff KnxDatapointType = 45
    KnxDatapointType_DPT_Value_1_Ucount KnxDatapointType = 46
    KnxDatapointType_DPT_FanStage KnxDatapointType = 47
    KnxDatapointType_DPT_Percent_V8 KnxDatapointType = 48
    KnxDatapointType_DPT_Value_1_Count KnxDatapointType = 49
    KnxDatapointType_DPT_Status_Mode3 KnxDatapointType = 50
    KnxDatapointType_DPT_Value_2_Ucount KnxDatapointType = 51
    KnxDatapointType_DPT_TimePeriodMsec KnxDatapointType = 52
    KnxDatapointType_DPT_TimePeriod10Msec KnxDatapointType = 53
    KnxDatapointType_DPT_TimePeriod100Msec KnxDatapointType = 54
    KnxDatapointType_DPT_TimePeriodSec KnxDatapointType = 55
    KnxDatapointType_DPT_TimePeriodMin KnxDatapointType = 56
    KnxDatapointType_DPT_TimePeriodHrs KnxDatapointType = 57
    KnxDatapointType_DPT_PropDataType KnxDatapointType = 58
    KnxDatapointType_DPT_Length_mm KnxDatapointType = 59
    KnxDatapointType_DPT_UElCurrentmA KnxDatapointType = 60
    KnxDatapointType_DPT_Brightness KnxDatapointType = 61
    KnxDatapointType_DPT_Absolute_Colour_Temperature KnxDatapointType = 62
    KnxDatapointType_DPT_Value_2_Count KnxDatapointType = 63
    KnxDatapointType_DPT_DeltaTimeMsec KnxDatapointType = 64
    KnxDatapointType_DPT_DeltaTime10Msec KnxDatapointType = 65
    KnxDatapointType_DPT_DeltaTime100Msec KnxDatapointType = 66
    KnxDatapointType_DPT_DeltaTimeSec KnxDatapointType = 67
    KnxDatapointType_DPT_DeltaTimeMin KnxDatapointType = 68
    KnxDatapointType_DPT_DeltaTimeHrs KnxDatapointType = 69
    KnxDatapointType_DPT_Percent_V16 KnxDatapointType = 70
    KnxDatapointType_DPT_Rotation_Angle KnxDatapointType = 71
    KnxDatapointType_DPT_Length_m KnxDatapointType = 72
    KnxDatapointType_DPT_Value_Temp KnxDatapointType = 73
    KnxDatapointType_DPT_Value_Tempd KnxDatapointType = 74
    KnxDatapointType_DPT_Value_Tempa KnxDatapointType = 75
    KnxDatapointType_DPT_Value_Lux KnxDatapointType = 76
    KnxDatapointType_DPT_Value_Wsp KnxDatapointType = 77
    KnxDatapointType_DPT_Value_Pres KnxDatapointType = 78
    KnxDatapointType_DPT_Value_Humidity KnxDatapointType = 79
    KnxDatapointType_DPT_Value_AirQuality KnxDatapointType = 80
    KnxDatapointType_DPT_Value_AirFlow KnxDatapointType = 81
    KnxDatapointType_DPT_Value_Time1 KnxDatapointType = 82
    KnxDatapointType_DPT_Value_Time2 KnxDatapointType = 83
    KnxDatapointType_DPT_Value_Volt KnxDatapointType = 84
    KnxDatapointType_DPT_Value_Curr KnxDatapointType = 85
    KnxDatapointType_DPT_PowerDensity KnxDatapointType = 86
    KnxDatapointType_DPT_KelvinPerPercent KnxDatapointType = 87
    KnxDatapointType_DPT_Power KnxDatapointType = 88
    KnxDatapointType_DPT_Value_Volume_Flow KnxDatapointType = 89
    KnxDatapointType_DPT_Rain_Amount KnxDatapointType = 90
    KnxDatapointType_DPT_Value_Temp_F KnxDatapointType = 91
    KnxDatapointType_DPT_Value_Wsp_kmh KnxDatapointType = 92
    KnxDatapointType_DPT_Value_Absolute_Humidity KnxDatapointType = 93
    KnxDatapointType_DPT_Concentration_ygm3 KnxDatapointType = 94
    KnxDatapointType_DPT_TimeOfDay KnxDatapointType = 95
    KnxDatapointType_DPT_Date KnxDatapointType = 96
    KnxDatapointType_DPT_Value_4_Ucount KnxDatapointType = 97
    KnxDatapointType_DPT_LongTimePeriod_Sec KnxDatapointType = 98
    KnxDatapointType_DPT_LongTimePeriod_Min KnxDatapointType = 99
    KnxDatapointType_DPT_LongTimePeriod_Hrs KnxDatapointType = 100
    KnxDatapointType_DPT_VolumeLiquid_Litre KnxDatapointType = 101
    KnxDatapointType_DPT_Volume_m_3 KnxDatapointType = 102
    KnxDatapointType_DPT_Value_4_Count KnxDatapointType = 103
    KnxDatapointType_DPT_FlowRate_m3h KnxDatapointType = 104
    KnxDatapointType_DPT_ActiveEnergy KnxDatapointType = 105
    KnxDatapointType_DPT_ApparantEnergy KnxDatapointType = 106
    KnxDatapointType_DPT_ReactiveEnergy KnxDatapointType = 107
    KnxDatapointType_DPT_ActiveEnergy_kWh KnxDatapointType = 108
    KnxDatapointType_DPT_ApparantEnergy_kVAh KnxDatapointType = 109
    KnxDatapointType_DPT_ReactiveEnergy_kVARh KnxDatapointType = 110
    KnxDatapointType_DPT_ActiveEnergy_MWh KnxDatapointType = 111
    KnxDatapointType_DPT_LongDeltaTimeSec KnxDatapointType = 112
    KnxDatapointType_DPT_DeltaVolumeLiquid_Litre KnxDatapointType = 113
    KnxDatapointType_DPT_DeltaVolume_m_3 KnxDatapointType = 114
    KnxDatapointType_DPT_Value_Acceleration KnxDatapointType = 115
    KnxDatapointType_DPT_Value_Acceleration_Angular KnxDatapointType = 116
    KnxDatapointType_DPT_Value_Activation_Energy KnxDatapointType = 117
    KnxDatapointType_DPT_Value_Activity KnxDatapointType = 118
    KnxDatapointType_DPT_Value_Mol KnxDatapointType = 119
    KnxDatapointType_DPT_Value_Amplitude KnxDatapointType = 120
    KnxDatapointType_DPT_Value_AngleRad KnxDatapointType = 121
    KnxDatapointType_DPT_Value_AngleDeg KnxDatapointType = 122
    KnxDatapointType_DPT_Value_Angular_Momentum KnxDatapointType = 123
    KnxDatapointType_DPT_Value_Angular_Velocity KnxDatapointType = 124
    KnxDatapointType_DPT_Value_Area KnxDatapointType = 125
    KnxDatapointType_DPT_Value_Capacitance KnxDatapointType = 126
    KnxDatapointType_DPT_Value_Charge_DensitySurface KnxDatapointType = 127
    KnxDatapointType_DPT_Value_Charge_DensityVolume KnxDatapointType = 128
    KnxDatapointType_DPT_Value_Compressibility KnxDatapointType = 129
    KnxDatapointType_DPT_Value_Conductance KnxDatapointType = 130
    KnxDatapointType_DPT_Value_Electrical_Conductivity KnxDatapointType = 131
    KnxDatapointType_DPT_Value_Density KnxDatapointType = 132
    KnxDatapointType_DPT_Value_Electric_Charge KnxDatapointType = 133
    KnxDatapointType_DPT_Value_Electric_Current KnxDatapointType = 134
    KnxDatapointType_DPT_Value_Electric_CurrentDensity KnxDatapointType = 135
    KnxDatapointType_DPT_Value_Electric_DipoleMoment KnxDatapointType = 136
    KnxDatapointType_DPT_Value_Electric_Displacement KnxDatapointType = 137
    KnxDatapointType_DPT_Value_Electric_FieldStrength KnxDatapointType = 138
    KnxDatapointType_DPT_Value_Electric_Flux KnxDatapointType = 139
    KnxDatapointType_DPT_Value_Electric_FluxDensity KnxDatapointType = 140
    KnxDatapointType_DPT_Value_Electric_Polarization KnxDatapointType = 141
    KnxDatapointType_DPT_Value_Electric_Potential KnxDatapointType = 142
    KnxDatapointType_DPT_Value_Electric_PotentialDifference KnxDatapointType = 143
    KnxDatapointType_DPT_Value_ElectromagneticMoment KnxDatapointType = 144
    KnxDatapointType_DPT_Value_Electromotive_Force KnxDatapointType = 145
    KnxDatapointType_DPT_Value_Energy KnxDatapointType = 146
    KnxDatapointType_DPT_Value_Force KnxDatapointType = 147
    KnxDatapointType_DPT_Value_Frequency KnxDatapointType = 148
    KnxDatapointType_DPT_Value_Angular_Frequency KnxDatapointType = 149
    KnxDatapointType_DPT_Value_Heat_Capacity KnxDatapointType = 150
    KnxDatapointType_DPT_Value_Heat_FlowRate KnxDatapointType = 151
    KnxDatapointType_DPT_Value_Heat_Quantity KnxDatapointType = 152
    KnxDatapointType_DPT_Value_Impedance KnxDatapointType = 153
    KnxDatapointType_DPT_Value_Length KnxDatapointType = 154
    KnxDatapointType_DPT_Value_Light_Quantity KnxDatapointType = 155
    KnxDatapointType_DPT_Value_Luminance KnxDatapointType = 156
    KnxDatapointType_DPT_Value_Luminous_Flux KnxDatapointType = 157
    KnxDatapointType_DPT_Value_Luminous_Intensity KnxDatapointType = 158
    KnxDatapointType_DPT_Value_Magnetic_FieldStrength KnxDatapointType = 159
    KnxDatapointType_DPT_Value_Magnetic_Flux KnxDatapointType = 160
    KnxDatapointType_DPT_Value_Magnetic_FluxDensity KnxDatapointType = 161
    KnxDatapointType_DPT_Value_Magnetic_Moment KnxDatapointType = 162
    KnxDatapointType_DPT_Value_Magnetic_Polarization KnxDatapointType = 163
    KnxDatapointType_DPT_Value_Magnetization KnxDatapointType = 164
    KnxDatapointType_DPT_Value_MagnetomotiveForce KnxDatapointType = 165
    KnxDatapointType_DPT_Value_Mass KnxDatapointType = 166
    KnxDatapointType_DPT_Value_MassFlux KnxDatapointType = 167
    KnxDatapointType_DPT_Value_Momentum KnxDatapointType = 168
    KnxDatapointType_DPT_Value_Phase_AngleRad KnxDatapointType = 169
    KnxDatapointType_DPT_Value_Phase_AngleDeg KnxDatapointType = 170
    KnxDatapointType_DPT_Value_Power KnxDatapointType = 171
    KnxDatapointType_DPT_Value_Power_Factor KnxDatapointType = 172
    KnxDatapointType_DPT_Value_Pressure KnxDatapointType = 173
    KnxDatapointType_DPT_Value_Reactance KnxDatapointType = 174
    KnxDatapointType_DPT_Value_Resistance KnxDatapointType = 175
    KnxDatapointType_DPT_Value_Resistivity KnxDatapointType = 176
    KnxDatapointType_DPT_Value_SelfInductance KnxDatapointType = 177
    KnxDatapointType_DPT_Value_SolidAngle KnxDatapointType = 178
    KnxDatapointType_DPT_Value_Sound_Intensity KnxDatapointType = 179
    KnxDatapointType_DPT_Value_Speed KnxDatapointType = 180
    KnxDatapointType_DPT_Value_Stress KnxDatapointType = 181
    KnxDatapointType_DPT_Value_Surface_Tension KnxDatapointType = 182
    KnxDatapointType_DPT_Value_Common_Temperature KnxDatapointType = 183
    KnxDatapointType_DPT_Value_Absolute_Temperature KnxDatapointType = 184
    KnxDatapointType_DPT_Value_TemperatureDifference KnxDatapointType = 185
    KnxDatapointType_DPT_Value_Thermal_Capacity KnxDatapointType = 186
    KnxDatapointType_DPT_Value_Thermal_Conductivity KnxDatapointType = 187
    KnxDatapointType_DPT_Value_ThermoelectricPower KnxDatapointType = 188
    KnxDatapointType_DPT_Value_Time KnxDatapointType = 189
    KnxDatapointType_DPT_Value_Torque KnxDatapointType = 190
    KnxDatapointType_DPT_Value_Volume KnxDatapointType = 191
    KnxDatapointType_DPT_Value_Volume_Flux KnxDatapointType = 192
    KnxDatapointType_DPT_Value_Weight KnxDatapointType = 193
    KnxDatapointType_DPT_Value_Work KnxDatapointType = 194
    KnxDatapointType_DPT_Volume_Flux_Meter KnxDatapointType = 195
    KnxDatapointType_DPT_Volume_Flux_ls KnxDatapointType = 196
    KnxDatapointType_DPT_Access_Data KnxDatapointType = 197
    KnxDatapointType_DPT_String_ASCII KnxDatapointType = 198
    KnxDatapointType_DPT_String_8859_1 KnxDatapointType = 199
    KnxDatapointType_DPT_SceneNumber KnxDatapointType = 200
    KnxDatapointType_DPT_SceneControl KnxDatapointType = 201
    KnxDatapointType_DPT_DateTime KnxDatapointType = 202
    KnxDatapointType_DPT_SCLOMode KnxDatapointType = 203
    KnxDatapointType_DPT_BuildingMode KnxDatapointType = 204
    KnxDatapointType_DPT_OccMode KnxDatapointType = 205
    KnxDatapointType_DPT_Priority KnxDatapointType = 206
    KnxDatapointType_DPT_LightApplicationMode KnxDatapointType = 207
    KnxDatapointType_DPT_ApplicationArea KnxDatapointType = 208
    KnxDatapointType_DPT_AlarmClassType KnxDatapointType = 209
    KnxDatapointType_DPT_PSUMode KnxDatapointType = 210
    KnxDatapointType_DPT_ErrorClass_System KnxDatapointType = 211
    KnxDatapointType_DPT_ErrorClass_HVAC KnxDatapointType = 212
    KnxDatapointType_DPT_Time_Delay KnxDatapointType = 213
    KnxDatapointType_DPT_Beaufort_Wind_Force_Scale KnxDatapointType = 214
    KnxDatapointType_DPT_SensorSelect KnxDatapointType = 215
    KnxDatapointType_DPT_ActuatorConnectType KnxDatapointType = 216
    KnxDatapointType_DPT_Cloud_Cover KnxDatapointType = 217
    KnxDatapointType_DPT_PowerReturnMode KnxDatapointType = 218
    KnxDatapointType_DPT_FuelType KnxDatapointType = 219
    KnxDatapointType_DPT_BurnerType KnxDatapointType = 220
    KnxDatapointType_DPT_HVACMode KnxDatapointType = 221
    KnxDatapointType_DPT_DHWMode KnxDatapointType = 222
    KnxDatapointType_DPT_LoadPriority KnxDatapointType = 223
    KnxDatapointType_DPT_HVACContrMode KnxDatapointType = 224
    KnxDatapointType_DPT_HVACEmergMode KnxDatapointType = 225
    KnxDatapointType_DPT_ChangeoverMode KnxDatapointType = 226
    KnxDatapointType_DPT_ValveMode KnxDatapointType = 227
    KnxDatapointType_DPT_DamperMode KnxDatapointType = 228
    KnxDatapointType_DPT_HeaterMode KnxDatapointType = 229
    KnxDatapointType_DPT_FanMode KnxDatapointType = 230
    KnxDatapointType_DPT_MasterSlaveMode KnxDatapointType = 231
    KnxDatapointType_DPT_StatusRoomSetp KnxDatapointType = 232
    KnxDatapointType_DPT_Metering_DeviceType KnxDatapointType = 233
    KnxDatapointType_DPT_HumDehumMode KnxDatapointType = 234
    KnxDatapointType_DPT_EnableHCStage KnxDatapointType = 235
    KnxDatapointType_DPT_ADAType KnxDatapointType = 236
    KnxDatapointType_DPT_BackupMode KnxDatapointType = 237
    KnxDatapointType_DPT_StartSynchronization KnxDatapointType = 238
    KnxDatapointType_DPT_Behaviour_Lock_Unlock KnxDatapointType = 239
    KnxDatapointType_DPT_Behaviour_Bus_Power_Up_Down KnxDatapointType = 240
    KnxDatapointType_DPT_DALI_Fade_Time KnxDatapointType = 241
    KnxDatapointType_DPT_BlinkingMode KnxDatapointType = 242
    KnxDatapointType_DPT_LightControlMode KnxDatapointType = 243
    KnxDatapointType_DPT_SwitchPBModel KnxDatapointType = 244
    KnxDatapointType_DPT_PBAction KnxDatapointType = 245
    KnxDatapointType_DPT_DimmPBModel KnxDatapointType = 246
    KnxDatapointType_DPT_SwitchOnMode KnxDatapointType = 247
    KnxDatapointType_DPT_LoadTypeSet KnxDatapointType = 248
    KnxDatapointType_DPT_LoadTypeDetected KnxDatapointType = 249
    KnxDatapointType_DPT_Converter_Test_Control KnxDatapointType = 250
    KnxDatapointType_DPT_SABExcept_Behaviour KnxDatapointType = 251
    KnxDatapointType_DPT_SABBehaviour_Lock_Unlock KnxDatapointType = 252
    KnxDatapointType_DPT_SSSBMode KnxDatapointType = 253
    KnxDatapointType_DPT_BlindsControlMode KnxDatapointType = 254
    KnxDatapointType_DPT_CommMode KnxDatapointType = 255
    KnxDatapointType_DPT_AddInfoTypes KnxDatapointType = 256
    KnxDatapointType_DPT_RF_ModeSelect KnxDatapointType = 257
    KnxDatapointType_DPT_RF_FilterSelect KnxDatapointType = 258
    KnxDatapointType_DPT_StatusGen KnxDatapointType = 259
    KnxDatapointType_DPT_Device_Control KnxDatapointType = 260
    KnxDatapointType_DPT_ForceSign KnxDatapointType = 261
    KnxDatapointType_DPT_ForceSignCool KnxDatapointType = 262
    KnxDatapointType_DPT_StatusRHC KnxDatapointType = 263
    KnxDatapointType_DPT_StatusSDHWC KnxDatapointType = 264
    KnxDatapointType_DPT_FuelTypeSet KnxDatapointType = 265
    KnxDatapointType_DPT_StatusRCC KnxDatapointType = 266
    KnxDatapointType_DPT_StatusAHU KnxDatapointType = 267
    KnxDatapointType_DPT_CombinedStatus_RTSM KnxDatapointType = 268
    KnxDatapointType_DPT_LightActuatorErrorInfo KnxDatapointType = 269
    KnxDatapointType_DPT_RF_ModeInfo KnxDatapointType = 270
    KnxDatapointType_DPT_RF_FilterInfo KnxDatapointType = 271
    KnxDatapointType_DPT_Channel_Activation_8 KnxDatapointType = 272
    KnxDatapointType_DPT_StatusDHWC KnxDatapointType = 273
    KnxDatapointType_DPT_StatusRHCC KnxDatapointType = 274
    KnxDatapointType_DPT_CombinedStatus_HVA KnxDatapointType = 275
    KnxDatapointType_DPT_CombinedStatus_RTC KnxDatapointType = 276
    KnxDatapointType_DPT_Media KnxDatapointType = 277
    KnxDatapointType_DPT_Channel_Activation_16 KnxDatapointType = 278
    KnxDatapointType_DPT_OnOffAction KnxDatapointType = 279
    KnxDatapointType_DPT_Alarm_Reaction KnxDatapointType = 280
    KnxDatapointType_DPT_UpDown_Action KnxDatapointType = 281
    KnxDatapointType_DPT_HVAC_PB_Action KnxDatapointType = 282
    KnxDatapointType_DPT_DoubleNibble KnxDatapointType = 283
    KnxDatapointType_DPT_SceneInfo KnxDatapointType = 284
    KnxDatapointType_DPT_CombinedInfoOnOff KnxDatapointType = 285
    KnxDatapointType_DPT_ActiveEnergy_V64 KnxDatapointType = 286
    KnxDatapointType_DPT_ApparantEnergy_V64 KnxDatapointType = 287
    KnxDatapointType_DPT_ReactiveEnergy_V64 KnxDatapointType = 288
    KnxDatapointType_DPT_Channel_Activation_24 KnxDatapointType = 289
    KnxDatapointType_DPT_HVACModeNext KnxDatapointType = 290
    KnxDatapointType_DPT_DHWModeNext KnxDatapointType = 291
    KnxDatapointType_DPT_OccModeNext KnxDatapointType = 292
    KnxDatapointType_DPT_BuildingModeNext KnxDatapointType = 293
    KnxDatapointType_DPT_Version KnxDatapointType = 294
    KnxDatapointType_DPT_AlarmInfo KnxDatapointType = 295
    KnxDatapointType_DPT_TempRoomSetpSetF16_3 KnxDatapointType = 296
    KnxDatapointType_DPT_TempRoomSetpSetShiftF16_3 KnxDatapointType = 297
    KnxDatapointType_DPT_Scaling_Speed KnxDatapointType = 298
    KnxDatapointType_DPT_Scaling_Step_Time KnxDatapointType = 299
    KnxDatapointType_DPT_MeteringValue KnxDatapointType = 300
    KnxDatapointType_DPT_MBus_Address KnxDatapointType = 301
    KnxDatapointType_DPT_Colour_RGB KnxDatapointType = 302
    KnxDatapointType_DPT_LanguageCodeAlpha2_ASCII KnxDatapointType = 303
    KnxDatapointType_DPT_Tariff_ActiveEnergy KnxDatapointType = 304
    KnxDatapointType_DPT_Prioritised_Mode_Control KnxDatapointType = 305
    KnxDatapointType_DPT_DALI_Control_Gear_Diagnostic KnxDatapointType = 306
    KnxDatapointType_DPT_DALI_Diagnostics KnxDatapointType = 307
    KnxDatapointType_DPT_CombinedPosition KnxDatapointType = 308
    KnxDatapointType_DPT_StatusSAB KnxDatapointType = 309
    KnxDatapointType_DPT_Colour_xyY KnxDatapointType = 310
    KnxDatapointType_DPT_Converter_Status KnxDatapointType = 311
    KnxDatapointType_DPT_Converter_Test_Result KnxDatapointType = 312
    KnxDatapointType_DPT_Battery_Info KnxDatapointType = 313
    KnxDatapointType_DPT_Brightness_Colour_Temperature_Transition KnxDatapointType = 314
    KnxDatapointType_DPT_Brightness_Colour_Temperature_Control KnxDatapointType = 315
    KnxDatapointType_DPT_Colour_RGBW KnxDatapointType = 316
    KnxDatapointType_DPT_Relative_Control_RGBW KnxDatapointType = 317
    KnxDatapointType_DPT_Relative_Control_RGB KnxDatapointType = 318
    KnxDatapointType_DPT_GeographicalLocation KnxDatapointType = 319
    KnxDatapointType_DPT_TempRoomSetpSetF16_4 KnxDatapointType = 320
    KnxDatapointType_DPT_TempRoomSetpSetShiftF16_4 KnxDatapointType = 321
)


func (e KnxDatapointType) Number() uint16 {
    switch e  {
        case 0: { /* '0' */
            return 0
        }
        case 1: { /* '1' */
            return 1
        }
        case 10: { /* '10' */
            return 10
        }
        case 100: { /* '100' */
            return 102
        }
        case 101: { /* '101' */
            return 1200
        }
        case 102: { /* '102' */
            return 1201
        }
        case 103: { /* '103' */
            return 1
        }
        case 104: { /* '104' */
            return 2
        }
        case 105: { /* '105' */
            return 10
        }
        case 106: { /* '106' */
            return 11
        }
        case 107: { /* '107' */
            return 12
        }
        case 108: { /* '108' */
            return 13
        }
        case 109: { /* '109' */
            return 14
        }
        case 11: { /* '11' */
            return 11
        }
        case 110: { /* '110' */
            return 15
        }
        case 111: { /* '111' */
            return 16
        }
        case 112: { /* '112' */
            return 100
        }
        case 113: { /* '113' */
            return 1200
        }
        case 114: { /* '114' */
            return 1201
        }
        case 115: { /* '115' */
            return 0
        }
        case 116: { /* '116' */
            return 1
        }
        case 117: { /* '117' */
            return 2
        }
        case 118: { /* '118' */
            return 3
        }
        case 119: { /* '119' */
            return 4
        }
        case 12: { /* '12' */
            return 12
        }
        case 120: { /* '120' */
            return 5
        }
        case 121: { /* '121' */
            return 6
        }
        case 122: { /* '122' */
            return 7
        }
        case 123: { /* '123' */
            return 8
        }
        case 124: { /* '124' */
            return 9
        }
        case 125: { /* '125' */
            return 10
        }
        case 126: { /* '126' */
            return 11
        }
        case 127: { /* '127' */
            return 12
        }
        case 128: { /* '128' */
            return 13
        }
        case 129: { /* '129' */
            return 14
        }
        case 13: { /* '13' */
            return 13
        }
        case 130: { /* '130' */
            return 15
        }
        case 131: { /* '131' */
            return 16
        }
        case 132: { /* '132' */
            return 17
        }
        case 133: { /* '133' */
            return 18
        }
        case 134: { /* '134' */
            return 19
        }
        case 135: { /* '135' */
            return 20
        }
        case 136: { /* '136' */
            return 21
        }
        case 137: { /* '137' */
            return 22
        }
        case 138: { /* '138' */
            return 23
        }
        case 139: { /* '139' */
            return 24
        }
        case 14: { /* '14' */
            return 14
        }
        case 140: { /* '140' */
            return 25
        }
        case 141: { /* '141' */
            return 26
        }
        case 142: { /* '142' */
            return 27
        }
        case 143: { /* '143' */
            return 28
        }
        case 144: { /* '144' */
            return 29
        }
        case 145: { /* '145' */
            return 30
        }
        case 146: { /* '146' */
            return 31
        }
        case 147: { /* '147' */
            return 32
        }
        case 148: { /* '148' */
            return 33
        }
        case 149: { /* '149' */
            return 34
        }
        case 15: { /* '15' */
            return 15
        }
        case 150: { /* '150' */
            return 35
        }
        case 151: { /* '151' */
            return 36
        }
        case 152: { /* '152' */
            return 37
        }
        case 153: { /* '153' */
            return 38
        }
        case 154: { /* '154' */
            return 39
        }
        case 155: { /* '155' */
            return 40
        }
        case 156: { /* '156' */
            return 41
        }
        case 157: { /* '157' */
            return 42
        }
        case 158: { /* '158' */
            return 43
        }
        case 159: { /* '159' */
            return 44
        }
        case 16: { /* '16' */
            return 16
        }
        case 160: { /* '160' */
            return 45
        }
        case 161: { /* '161' */
            return 46
        }
        case 162: { /* '162' */
            return 47
        }
        case 163: { /* '163' */
            return 48
        }
        case 164: { /* '164' */
            return 49
        }
        case 165: { /* '165' */
            return 50
        }
        case 166: { /* '166' */
            return 51
        }
        case 167: { /* '167' */
            return 52
        }
        case 168: { /* '168' */
            return 53
        }
        case 169: { /* '169' */
            return 54
        }
        case 17: { /* '17' */
            return 17
        }
        case 170: { /* '170' */
            return 55
        }
        case 171: { /* '171' */
            return 56
        }
        case 172: { /* '172' */
            return 57
        }
        case 173: { /* '173' */
            return 58
        }
        case 174: { /* '174' */
            return 59
        }
        case 175: { /* '175' */
            return 60
        }
        case 176: { /* '176' */
            return 61
        }
        case 177: { /* '177' */
            return 62
        }
        case 178: { /* '178' */
            return 63
        }
        case 179: { /* '179' */
            return 64
        }
        case 18: { /* '18' */
            return 18
        }
        case 180: { /* '180' */
            return 65
        }
        case 181: { /* '181' */
            return 66
        }
        case 182: { /* '182' */
            return 67
        }
        case 183: { /* '183' */
            return 68
        }
        case 184: { /* '184' */
            return 69
        }
        case 185: { /* '185' */
            return 70
        }
        case 186: { /* '186' */
            return 71
        }
        case 187: { /* '187' */
            return 72
        }
        case 188: { /* '188' */
            return 73
        }
        case 189: { /* '189' */
            return 74
        }
        case 19: { /* '19' */
            return 19
        }
        case 190: { /* '190' */
            return 75
        }
        case 191: { /* '191' */
            return 76
        }
        case 192: { /* '192' */
            return 77
        }
        case 193: { /* '193' */
            return 78
        }
        case 194: { /* '194' */
            return 79
        }
        case 195: { /* '195' */
            return 1200
        }
        case 196: { /* '196' */
            return 1201
        }
        case 197: { /* '197' */
            return 0
        }
        case 198: { /* '198' */
            return 0
        }
        case 199: { /* '199' */
            return 1
        }
        case 2: { /* '2' */
            return 2
        }
        case 20: { /* '20' */
            return 21
        }
        case 200: { /* '200' */
            return 1
        }
        case 201: { /* '201' */
            return 1
        }
        case 202: { /* '202' */
            return 1
        }
        case 203: { /* '203' */
            return 1
        }
        case 204: { /* '204' */
            return 2
        }
        case 205: { /* '205' */
            return 3
        }
        case 206: { /* '206' */
            return 4
        }
        case 207: { /* '207' */
            return 5
        }
        case 208: { /* '208' */
            return 6
        }
        case 209: { /* '209' */
            return 7
        }
        case 21: { /* '21' */
            return 22
        }
        case 210: { /* '210' */
            return 8
        }
        case 211: { /* '211' */
            return 11
        }
        case 212: { /* '212' */
            return 12
        }
        case 213: { /* '213' */
            return 13
        }
        case 214: { /* '214' */
            return 14
        }
        case 215: { /* '215' */
            return 17
        }
        case 216: { /* '216' */
            return 20
        }
        case 217: { /* '217' */
            return 21
        }
        case 218: { /* '218' */
            return 22
        }
        case 219: { /* '219' */
            return 100
        }
        case 22: { /* '22' */
            return 23
        }
        case 220: { /* '220' */
            return 101
        }
        case 221: { /* '221' */
            return 102
        }
        case 222: { /* '222' */
            return 103
        }
        case 223: { /* '223' */
            return 104
        }
        case 224: { /* '224' */
            return 105
        }
        case 225: { /* '225' */
            return 106
        }
        case 226: { /* '226' */
            return 107
        }
        case 227: { /* '227' */
            return 108
        }
        case 228: { /* '228' */
            return 109
        }
        case 229: { /* '229' */
            return 110
        }
        case 23: { /* '23' */
            return 24
        }
        case 230: { /* '230' */
            return 111
        }
        case 231: { /* '231' */
            return 112
        }
        case 232: { /* '232' */
            return 113
        }
        case 233: { /* '233' */
            return 114
        }
        case 234: { /* '234' */
            return 115
        }
        case 235: { /* '235' */
            return 116
        }
        case 236: { /* '236' */
            return 120
        }
        case 237: { /* '237' */
            return 121
        }
        case 238: { /* '238' */
            return 122
        }
        case 239: { /* '239' */
            return 600
        }
        case 24: { /* '24' */
            return 100
        }
        case 240: { /* '240' */
            return 601
        }
        case 241: { /* '241' */
            return 602
        }
        case 242: { /* '242' */
            return 603
        }
        case 243: { /* '243' */
            return 604
        }
        case 244: { /* '244' */
            return 605
        }
        case 245: { /* '245' */
            return 606
        }
        case 246: { /* '246' */
            return 607
        }
        case 247: { /* '247' */
            return 608
        }
        case 248: { /* '248' */
            return 609
        }
        case 249: { /* '249' */
            return 610
        }
        case 25: { /* '25' */
            return 1
        }
        case 250: { /* '250' */
            return 611
        }
        case 251: { /* '251' */
            return 801
        }
        case 252: { /* '252' */
            return 802
        }
        case 253: { /* '253' */
            return 803
        }
        case 254: { /* '254' */
            return 804
        }
        case 255: { /* '255' */
            return 1000
        }
        case 256: { /* '256' */
            return 1001
        }
        case 257: { /* '257' */
            return 1002
        }
        case 258: { /* '258' */
            return 1003
        }
        case 259: { /* '259' */
            return 1
        }
        case 26: { /* '26' */
            return 2
        }
        case 260: { /* '260' */
            return 2
        }
        case 261: { /* '261' */
            return 100
        }
        case 262: { /* '262' */
            return 101
        }
        case 263: { /* '263' */
            return 102
        }
        case 264: { /* '264' */
            return 103
        }
        case 265: { /* '265' */
            return 104
        }
        case 266: { /* '266' */
            return 105
        }
        case 267: { /* '267' */
            return 106
        }
        case 268: { /* '268' */
            return 107
        }
        case 269: { /* '269' */
            return 601
        }
        case 27: { /* '27' */
            return 3
        }
        case 270: { /* '270' */
            return 1000
        }
        case 271: { /* '271' */
            return 1001
        }
        case 272: { /* '272' */
            return 1010
        }
        case 273: { /* '273' */
            return 100
        }
        case 274: { /* '274' */
            return 101
        }
        case 275: { /* '275' */
            return 102
        }
        case 276: { /* '276' */
            return 103
        }
        case 277: { /* '277' */
            return 1000
        }
        case 278: { /* '278' */
            return 1010
        }
        case 279: { /* '279' */
            return 1
        }
        case 28: { /* '28' */
            return 4
        }
        case 280: { /* '280' */
            return 2
        }
        case 281: { /* '281' */
            return 3
        }
        case 282: { /* '282' */
            return 102
        }
        case 283: { /* '283' */
            return 1000
        }
        case 284: { /* '284' */
            return 1
        }
        case 285: { /* '285' */
            return 1
        }
        case 286: { /* '286' */
            return 10
        }
        case 287: { /* '287' */
            return 11
        }
        case 288: { /* '288' */
            return 12
        }
        case 289: { /* '289' */
            return 1010
        }
        case 29: { /* '29' */
            return 5
        }
        case 290: { /* '290' */
            return 100
        }
        case 291: { /* '291' */
            return 102
        }
        case 292: { /* '292' */
            return 104
        }
        case 293: { /* '293' */
            return 105
        }
        case 294: { /* '294' */
            return 1
        }
        case 295: { /* '295' */
            return 1
        }
        case 296: { /* '296' */
            return 100
        }
        case 297: { /* '297' */
            return 101
        }
        case 298: { /* '298' */
            return 1
        }
        case 299: { /* '299' */
            return 2
        }
        case 3: { /* '3' */
            return 3
        }
        case 30: { /* '30' */
            return 6
        }
        case 300: { /* '300' */
            return 1
        }
        case 301: { /* '301' */
            return 1000
        }
        case 302: { /* '302' */
            return 600
        }
        case 303: { /* '303' */
            return 1
        }
        case 304: { /* '304' */
            return 1
        }
        case 305: { /* '305' */
            return 1
        }
        case 306: { /* '306' */
            return 600
        }
        case 307: { /* '307' */
            return 600
        }
        case 308: { /* '308' */
            return 800
        }
        case 309: { /* '309' */
            return 800
        }
        case 31: { /* '31' */
            return 7
        }
        case 310: { /* '310' */
            return 600
        }
        case 311: { /* '311' */
            return 600
        }
        case 312: { /* '312' */
            return 600
        }
        case 313: { /* '313' */
            return 600
        }
        case 314: { /* '314' */
            return 600
        }
        case 315: { /* '315' */
            return 600
        }
        case 316: { /* '316' */
            return 600
        }
        case 317: { /* '317' */
            return 600
        }
        case 318: { /* '318' */
            return 600
        }
        case 319: { /* '319' */
            return 1
        }
        case 32: { /* '32' */
            return 8
        }
        case 320: { /* '320' */
            return 100
        }
        case 321: { /* '321' */
            return 101
        }
        case 33: { /* '33' */
            return 9
        }
        case 34: { /* '34' */
            return 10
        }
        case 35: { /* '35' */
            return 11
        }
        case 36: { /* '36' */
            return 12
        }
        case 37: { /* '37' */
            return 7
        }
        case 38: { /* '38' */
            return 8
        }
        case 39: { /* '39' */
            return 1
        }
        case 4: { /* '4' */
            return 4
        }
        case 40: { /* '40' */
            return 2
        }
        case 41: { /* '41' */
            return 1
        }
        case 42: { /* '42' */
            return 3
        }
        case 43: { /* '43' */
            return 4
        }
        case 44: { /* '44' */
            return 5
        }
        case 45: { /* '45' */
            return 6
        }
        case 46: { /* '46' */
            return 10
        }
        case 47: { /* '47' */
            return 100
        }
        case 48: { /* '48' */
            return 1
        }
        case 49: { /* '49' */
            return 10
        }
        case 5: { /* '5' */
            return 5
        }
        case 50: { /* '50' */
            return 20
        }
        case 51: { /* '51' */
            return 1
        }
        case 52: { /* '52' */
            return 2
        }
        case 53: { /* '53' */
            return 3
        }
        case 54: { /* '54' */
            return 4
        }
        case 55: { /* '55' */
            return 5
        }
        case 56: { /* '56' */
            return 6
        }
        case 57: { /* '57' */
            return 7
        }
        case 58: { /* '58' */
            return 10
        }
        case 59: { /* '59' */
            return 11
        }
        case 6: { /* '6' */
            return 6
        }
        case 60: { /* '60' */
            return 12
        }
        case 61: { /* '61' */
            return 13
        }
        case 62: { /* '62' */
            return 600
        }
        case 63: { /* '63' */
            return 1
        }
        case 64: { /* '64' */
            return 2
        }
        case 65: { /* '65' */
            return 3
        }
        case 66: { /* '66' */
            return 4
        }
        case 67: { /* '67' */
            return 5
        }
        case 68: { /* '68' */
            return 6
        }
        case 69: { /* '69' */
            return 7
        }
        case 7: { /* '7' */
            return 7
        }
        case 70: { /* '70' */
            return 10
        }
        case 71: { /* '71' */
            return 11
        }
        case 72: { /* '72' */
            return 12
        }
        case 73: { /* '73' */
            return 1
        }
        case 74: { /* '74' */
            return 2
        }
        case 75: { /* '75' */
            return 3
        }
        case 76: { /* '76' */
            return 4
        }
        case 77: { /* '77' */
            return 5
        }
        case 78: { /* '78' */
            return 6
        }
        case 79: { /* '79' */
            return 7
        }
        case 8: { /* '8' */
            return 8
        }
        case 80: { /* '80' */
            return 8
        }
        case 81: { /* '81' */
            return 9
        }
        case 82: { /* '82' */
            return 10
        }
        case 83: { /* '83' */
            return 11
        }
        case 84: { /* '84' */
            return 20
        }
        case 85: { /* '85' */
            return 21
        }
        case 86: { /* '86' */
            return 22
        }
        case 87: { /* '87' */
            return 23
        }
        case 88: { /* '88' */
            return 24
        }
        case 89: { /* '89' */
            return 25
        }
        case 9: { /* '9' */
            return 9
        }
        case 90: { /* '90' */
            return 26
        }
        case 91: { /* '91' */
            return 27
        }
        case 92: { /* '92' */
            return 28
        }
        case 93: { /* '93' */
            return 29
        }
        case 94: { /* '94' */
            return 30
        }
        case 95: { /* '95' */
            return 1
        }
        case 96: { /* '96' */
            return 1
        }
        case 97: { /* '97' */
            return 1
        }
        case 98: { /* '98' */
            return 100
        }
        case 99: { /* '99' */
            return 101
        }
        default: {
            return 0
        }
    }
}

func (e KnxDatapointType) Name() string {
    switch e  {
        case 0: { /* '0' */
            return "Unknown Datapoint Subtype"
        }
        case 1: { /* '1' */
            return "switch"
        }
        case 10: { /* '10' */
            return "start/stop"
        }
        case 100: { /* '100' */
            return "counter timehrs (h)"
        }
        case 101: { /* '101' */
            return "volume liquid (l)"
        }
        case 102: { /* '102' */
            return "volume (m³)"
        }
        case 103: { /* '103' */
            return "counter pulses (signed)"
        }
        case 104: { /* '104' */
            return "flow rate (m³/h)"
        }
        case 105: { /* '105' */
            return "active energy (Wh)"
        }
        case 106: { /* '106' */
            return "apparant energy (VAh)"
        }
        case 107: { /* '107' */
            return "reactive energy (VARh)"
        }
        case 108: { /* '108' */
            return "active energy (kWh)"
        }
        case 109: { /* '109' */
            return "apparant energy (kVAh)"
        }
        case 11: { /* '11' */
            return "state"
        }
        case 110: { /* '110' */
            return "reactive energy (kVARh)"
        }
        case 111: { /* '111' */
            return "active energy (MWh)"
        }
        case 112: { /* '112' */
            return "time lag (s)"
        }
        case 113: { /* '113' */
            return "delta volume liquid (l)"
        }
        case 114: { /* '114' */
            return "delta volume (m³)"
        }
        case 115: { /* '115' */
            return "acceleration (m/s²)"
        }
        case 116: { /* '116' */
            return "angular acceleration (rad/s²)"
        }
        case 117: { /* '117' */
            return "activation energy (J/mol)"
        }
        case 118: { /* '118' */
            return "radioactive activity (1/s)"
        }
        case 119: { /* '119' */
            return "amount of substance (mol)"
        }
        case 12: { /* '12' */
            return "invert"
        }
        case 120: { /* '120' */
            return "amplitude"
        }
        case 121: { /* '121' */
            return "angle (radiant)"
        }
        case 122: { /* '122' */
            return "angle (degree)"
        }
        case 123: { /* '123' */
            return "angular momentum (Js)"
        }
        case 124: { /* '124' */
            return "angular velocity (rad/s)"
        }
        case 125: { /* '125' */
            return "area (m*m)"
        }
        case 126: { /* '126' */
            return "capacitance (F)"
        }
        case 127: { /* '127' */
            return "flux density (C/m²)"
        }
        case 128: { /* '128' */
            return "charge density (C/m³)"
        }
        case 129: { /* '129' */
            return "compressibility (m²/N)"
        }
        case 13: { /* '13' */
            return "dim send style"
        }
        case 130: { /* '130' */
            return "conductance (S)"
        }
        case 131: { /* '131' */
            return "conductivity  (S/m)"
        }
        case 132: { /* '132' */
            return "density (kg/m³)"
        }
        case 133: { /* '133' */
            return "electric charge (C)"
        }
        case 134: { /* '134' */
            return "electric current (A)"
        }
        case 135: { /* '135' */
            return "electric current density (A/m²)"
        }
        case 136: { /* '136' */
            return "electric dipole moment (Cm)"
        }
        case 137: { /* '137' */
            return "electric displacement (C/m²)"
        }
        case 138: { /* '138' */
            return "electric field strength (V/m)"
        }
        case 139: { /* '139' */
            return "electric flux (C)"
        }
        case 14: { /* '14' */
            return "input source"
        }
        case 140: { /* '140' */
            return "electric flux density (C/m²)"
        }
        case 141: { /* '141' */
            return "electric polarization (C/m²)"
        }
        case 142: { /* '142' */
            return "electric potential (V)"
        }
        case 143: { /* '143' */
            return "electric potential difference (V)"
        }
        case 144: { /* '144' */
            return "electromagnetic moment (Am²)"
        }
        case 145: { /* '145' */
            return "electromotive force (V)"
        }
        case 146: { /* '146' */
            return "energy (J)"
        }
        case 147: { /* '147' */
            return "force (N)"
        }
        case 148: { /* '148' */
            return "frequency (Hz)"
        }
        case 149: { /* '149' */
            return "angular frequency (rad/s)"
        }
        case 15: { /* '15' */
            return "reset"
        }
        case 150: { /* '150' */
            return "heat capacity (J/K)"
        }
        case 151: { /* '151' */
            return "heat flow rate (W)"
        }
        case 152: { /* '152' */
            return "heat quantity"
        }
        case 153: { /* '153' */
            return "impedance (Ω)"
        }
        case 154: { /* '154' */
            return "length (m)"
        }
        case 155: { /* '155' */
            return "light quantity (J)"
        }
        case 156: { /* '156' */
            return "luminance (cd/m²)"
        }
        case 157: { /* '157' */
            return "luminous flux (lm)"
        }
        case 158: { /* '158' */
            return "luminous intensity (cd)"
        }
        case 159: { /* '159' */
            return "magnetic field strength (A/m)"
        }
        case 16: { /* '16' */
            return "acknowledge"
        }
        case 160: { /* '160' */
            return "magnetic flux (Wb)"
        }
        case 161: { /* '161' */
            return "magnetic flux density (T)"
        }
        case 162: { /* '162' */
            return "magnetic moment (Am²)"
        }
        case 163: { /* '163' */
            return "magnetic polarization (T)"
        }
        case 164: { /* '164' */
            return "magnetization (A/m)"
        }
        case 165: { /* '165' */
            return "magnetomotive force (A)"
        }
        case 166: { /* '166' */
            return "mass (kg)"
        }
        case 167: { /* '167' */
            return "mass flux (kg/s)"
        }
        case 168: { /* '168' */
            return "momentum (N/s)"
        }
        case 169: { /* '169' */
            return "phase angle (rad)"
        }
        case 17: { /* '17' */
            return "trigger"
        }
        case 170: { /* '170' */
            return "phase angle (°)"
        }
        case 171: { /* '171' */
            return "power (W)"
        }
        case 172: { /* '172' */
            return "power factor (cos Φ)"
        }
        case 173: { /* '173' */
            return "pressure (Pa)"
        }
        case 174: { /* '174' */
            return "reactance (Ω)"
        }
        case 175: { /* '175' */
            return "resistance (Ω)"
        }
        case 176: { /* '176' */
            return "resistivity (Ωm)"
        }
        case 177: { /* '177' */
            return "self inductance (H)"
        }
        case 178: { /* '178' */
            return "solid angle (sr)"
        }
        case 179: { /* '179' */
            return "sound intensity (W/m²)"
        }
        case 18: { /* '18' */
            return "occupancy"
        }
        case 180: { /* '180' */
            return "speed (m/s)"
        }
        case 181: { /* '181' */
            return "stress (Pa)"
        }
        case 182: { /* '182' */
            return "surface tension (N/m)"
        }
        case 183: { /* '183' */
            return "temperature (°C)"
        }
        case 184: { /* '184' */
            return "temperature absolute (K)"
        }
        case 185: { /* '185' */
            return "temperature difference (K)"
        }
        case 186: { /* '186' */
            return "thermal capacity (J/K)"
        }
        case 187: { /* '187' */
            return "thermal conductivity (W/mK)"
        }
        case 188: { /* '188' */
            return "thermoelectric power (V/K)"
        }
        case 189: { /* '189' */
            return "time (s)"
        }
        case 19: { /* '19' */
            return "window/door"
        }
        case 190: { /* '190' */
            return "torque (Nm)"
        }
        case 191: { /* '191' */
            return "volume (m³)"
        }
        case 192: { /* '192' */
            return "volume flux (m³/s)"
        }
        case 193: { /* '193' */
            return "weight (N)"
        }
        case 194: { /* '194' */
            return "work (J)"
        }
        case 195: { /* '195' */
            return "volume flux for meters (m³/h)"
        }
        case 196: { /* '196' */
            return "volume flux for meters (1/ls)"
        }
        case 197: { /* '197' */
            return "entrance access"
        }
        case 198: { /* '198' */
            return "Character String (ASCII)"
        }
        case 199: { /* '199' */
            return "Character String (ISO 8859-1)"
        }
        case 2: { /* '2' */
            return "boolean"
        }
        case 20: { /* '20' */
            return "logical function"
        }
        case 200: { /* '200' */
            return "scene number"
        }
        case 201: { /* '201' */
            return "scene control"
        }
        case 202: { /* '202' */
            return "date time"
        }
        case 203: { /* '203' */
            return "SCLO mode"
        }
        case 204: { /* '204' */
            return "building mode"
        }
        case 205: { /* '205' */
            return "occupied"
        }
        case 206: { /* '206' */
            return "priority"
        }
        case 207: { /* '207' */
            return "light application mode"
        }
        case 208: { /* '208' */
            return "light application area"
        }
        case 209: { /* '209' */
            return "alarm class"
        }
        case 21: { /* '21' */
            return "scene"
        }
        case 210: { /* '210' */
            return "PSU mode"
        }
        case 211: { /* '211' */
            return "system error class"
        }
        case 212: { /* '212' */
            return "HVAC error class"
        }
        case 213: { /* '213' */
            return "time delay"
        }
        case 214: { /* '214' */
            return "wind force scale (0..12)"
        }
        case 215: { /* '215' */
            return "sensor mode"
        }
        case 216: { /* '216' */
            return "actuator connect type"
        }
        case 217: { /* '217' */
            return "cloud cover"
        }
        case 218: { /* '218' */
            return "power return mode"
        }
        case 219: { /* '219' */
            return "fuel type"
        }
        case 22: { /* '22' */
            return "shutter/blinds mode"
        }
        case 220: { /* '220' */
            return "burner type"
        }
        case 221: { /* '221' */
            return "HVAC mode"
        }
        case 222: { /* '222' */
            return "DHW mode"
        }
        case 223: { /* '223' */
            return "load priority"
        }
        case 224: { /* '224' */
            return "HVAC control mode"
        }
        case 225: { /* '225' */
            return "HVAC emergency mode"
        }
        case 226: { /* '226' */
            return "changeover mode"
        }
        case 227: { /* '227' */
            return "valve mode"
        }
        case 228: { /* '228' */
            return "damper mode"
        }
        case 229: { /* '229' */
            return "heater mode"
        }
        case 23: { /* '23' */
            return "day/night"
        }
        case 230: { /* '230' */
            return "fan mode"
        }
        case 231: { /* '231' */
            return "master/slave mode"
        }
        case 232: { /* '232' */
            return "status room setpoint"
        }
        case 233: { /* '233' */
            return "metering device type"
        }
        case 234: { /* '234' */
            return "hum dehum mode"
        }
        case 235: { /* '235' */
            return "enable H/C stage"
        }
        case 236: { /* '236' */
            return "ADA type"
        }
        case 237: { /* '237' */
            return "backup mode"
        }
        case 238: { /* '238' */
            return "start syncronization type"
        }
        case 239: { /* '239' */
            return "behavior lock/unlock"
        }
        case 24: { /* '24' */
            return "cooling/heating"
        }
        case 240: { /* '240' */
            return "behavior bus power up/down"
        }
        case 241: { /* '241' */
            return "dali fade time"
        }
        case 242: { /* '242' */
            return "blink mode"
        }
        case 243: { /* '243' */
            return "light control mode"
        }
        case 244: { /* '244' */
            return "PB switch mode"
        }
        case 245: { /* '245' */
            return "PB action mode"
        }
        case 246: { /* '246' */
            return "PB dimm mode"
        }
        case 247: { /* '247' */
            return "switch on mode"
        }
        case 248: { /* '248' */
            return "load type"
        }
        case 249: { /* '249' */
            return "load type detection"
        }
        case 25: { /* '25' */
            return "switch control"
        }
        case 250: { /* '250' */
            return "converter test control"
        }
        case 251: { /* '251' */
            return "SAB except behavior"
        }
        case 252: { /* '252' */
            return "SAB behavior on lock/unlock"
        }
        case 253: { /* '253' */
            return "SSSB mode"
        }
        case 254: { /* '254' */
            return "blinds control mode"
        }
        case 255: { /* '255' */
            return "communication mode"
        }
        case 256: { /* '256' */
            return "additional information type"
        }
        case 257: { /* '257' */
            return "RF mode selection"
        }
        case 258: { /* '258' */
            return "RF filter mode selection"
        }
        case 259: { /* '259' */
            return "general status"
        }
        case 26: { /* '26' */
            return "boolean control"
        }
        case 260: { /* '260' */
            return "device control"
        }
        case 261: { /* '261' */
            return "forcing signal"
        }
        case 262: { /* '262' */
            return "forcing signal cool"
        }
        case 263: { /* '263' */
            return "room heating controller status"
        }
        case 264: { /* '264' */
            return "solar DHW controller status"
        }
        case 265: { /* '265' */
            return "fuel type set"
        }
        case 266: { /* '266' */
            return "room cooling controller status"
        }
        case 267: { /* '267' */
            return "ventilation controller status"
        }
        case 268: { /* '268' */
            return "combined status RTSM"
        }
        case 269: { /* '269' */
            return "lighting actuator error information"
        }
        case 27: { /* '27' */
            return "enable control"
        }
        case 270: { /* '270' */
            return "RF communication mode info"
        }
        case 271: { /* '271' */
            return "cEMI server supported RF filtering modes"
        }
        case 272: { /* '272' */
            return "channel activation for 8 channels"
        }
        case 273: { /* '273' */
            return "DHW controller status"
        }
        case 274: { /* '274' */
            return "RHCC status"
        }
        case 275: { /* '275' */
            return "combined status HVA"
        }
        case 276: { /* '276' */
            return "combined status RTC"
        }
        case 277: { /* '277' */
            return "media"
        }
        case 278: { /* '278' */
            return "channel activation for 16 channels"
        }
        case 279: { /* '279' */
            return "on/off action"
        }
        case 28: { /* '28' */
            return "ramp control"
        }
        case 280: { /* '280' */
            return "alarm reaction"
        }
        case 281: { /* '281' */
            return "up/down action"
        }
        case 282: { /* '282' */
            return "HVAC push button action"
        }
        case 283: { /* '283' */
            return "busy/nak repetitions"
        }
        case 284: { /* '284' */
            return "scene information"
        }
        case 285: { /* '285' */
            return "bit-combined info on/off"
        }
        case 286: { /* '286' */
            return "active energy (Wh)"
        }
        case 287: { /* '287' */
            return "apparant energy (VAh)"
        }
        case 288: { /* '288' */
            return "reactive energy (VARh)"
        }
        case 289: { /* '289' */
            return "activation state 0..23"
        }
        case 29: { /* '29' */
            return "alarm control"
        }
        case 290: { /* '290' */
            return "time delay & HVAC mode"
        }
        case 291: { /* '291' */
            return "time delay & DHW mode"
        }
        case 292: { /* '292' */
            return "time delay & occupancy mode"
        }
        case 293: { /* '293' */
            return "time delay & building mode"
        }
        case 294: { /* '294' */
            return "DPT version"
        }
        case 295: { /* '295' */
            return "alarm info"
        }
        case 296: { /* '296' */
            return "room temperature setpoint"
        }
        case 297: { /* '297' */
            return "room temperature setpoint shift"
        }
        case 298: { /* '298' */
            return "scaling speed"
        }
        case 299: { /* '299' */
            return "scaling step time"
        }
        case 3: { /* '3' */
            return "enable"
        }
        case 30: { /* '30' */
            return "binary value control"
        }
        case 300: { /* '300' */
            return "metering value (value,encoding,cmd)"
        }
        case 301: { /* '301' */
            return "MBus address"
        }
        case 302: { /* '302' */
            return "RGB value 3x(0..255)"
        }
        case 303: { /* '303' */
            return "language code (ASCII)"
        }
        case 304: { /* '304' */
            return "electrical energy with tariff"
        }
        case 305: { /* '305' */
            return "priority control"
        }
        case 306: { /* '306' */
            return "diagnostic value"
        }
        case 307: { /* '307' */
            return "diagnostic value"
        }
        case 308: { /* '308' */
            return "combined position"
        }
        case 309: { /* '309' */
            return "status sunblind & shutter actuator"
        }
        case 31: { /* '31' */
            return "step control"
        }
        case 310: { /* '310' */
            return "colour xyY"
        }
        case 311: { /* '311' */
            return "DALI converter status"
        }
        case 312: { /* '312' */
            return "DALI converter test result"
        }
        case 313: { /* '313' */
            return "Battery Information"
        }
        case 314: { /* '314' */
            return "brightness colour temperature transition"
        }
        case 315: { /* '315' */
            return "brightness colour temperature control"
        }
        case 316: { /* '316' */
            return "RGBW value 4x(0..100%)"
        }
        case 317: { /* '317' */
            return "RGBW relative control"
        }
        case 318: { /* '318' */
            return "RGB relative control"
        }
        case 319: { /* '319' */
            return "geographical location (longitude and latitude) expressed in degrees"
        }
        case 32: { /* '32' */
            return "direction control 1"
        }
        case 320: { /* '320' */
            return "Temperature setpoint setting for 4 HVAC Modes"
        }
        case 321: { /* '321' */
            return "Temperature setpoint shift setting for 4 HVAC Modes"
        }
        case 33: { /* '33' */
            return "direction control 2"
        }
        case 34: { /* '34' */
            return "start control"
        }
        case 35: { /* '35' */
            return "state control"
        }
        case 36: { /* '36' */
            return "invert control"
        }
        case 37: { /* '37' */
            return "dimming control"
        }
        case 38: { /* '38' */
            return "blind control"
        }
        case 39: { /* '39' */
            return "character (ASCII)"
        }
        case 4: { /* '4' */
            return "ramp"
        }
        case 40: { /* '40' */
            return "character (ISO 8859-1)"
        }
        case 41: { /* '41' */
            return "percentage (0..100%)"
        }
        case 42: { /* '42' */
            return "angle (degrees)"
        }
        case 43: { /* '43' */
            return "percentage (0..255%)"
        }
        case 44: { /* '44' */
            return "ratio (0..255)"
        }
        case 45: { /* '45' */
            return "tariff (0..255)"
        }
        case 46: { /* '46' */
            return "counter pulses (0..255)"
        }
        case 47: { /* '47' */
            return "fan stage (0..255)"
        }
        case 48: { /* '48' */
            return "percentage (-128..127%)"
        }
        case 49: { /* '49' */
            return "counter pulses (-128..127)"
        }
        case 5: { /* '5' */
            return "alarm"
        }
        case 50: { /* '50' */
            return "status with mode"
        }
        case 51: { /* '51' */
            return "pulses"
        }
        case 52: { /* '52' */
            return "time (ms)"
        }
        case 53: { /* '53' */
            return "time (10 ms)"
        }
        case 54: { /* '54' */
            return "time (100 ms)"
        }
        case 55: { /* '55' */
            return "time (s)"
        }
        case 56: { /* '56' */
            return "time (min)"
        }
        case 57: { /* '57' */
            return "time (h)"
        }
        case 58: { /* '58' */
            return "property data type"
        }
        case 59: { /* '59' */
            return "length (mm)"
        }
        case 6: { /* '6' */
            return "binary value"
        }
        case 60: { /* '60' */
            return "current (mA)"
        }
        case 61: { /* '61' */
            return "brightness (lux)"
        }
        case 62: { /* '62' */
            return "absolute colour temperature (K)"
        }
        case 63: { /* '63' */
            return "pulses difference"
        }
        case 64: { /* '64' */
            return "time lag (ms)"
        }
        case 65: { /* '65' */
            return "time lag(10 ms)"
        }
        case 66: { /* '66' */
            return "time lag (100 ms)"
        }
        case 67: { /* '67' */
            return "time lag (s)"
        }
        case 68: { /* '68' */
            return "time lag (min)"
        }
        case 69: { /* '69' */
            return "time lag (h)"
        }
        case 7: { /* '7' */
            return "step"
        }
        case 70: { /* '70' */
            return "percentage difference (%)"
        }
        case 71: { /* '71' */
            return "rotation angle (°)"
        }
        case 72: { /* '72' */
            return "length (m)"
        }
        case 73: { /* '73' */
            return "temperature (°C)"
        }
        case 74: { /* '74' */
            return "temperature difference (K)"
        }
        case 75: { /* '75' */
            return "kelvin/hour (K/h)"
        }
        case 76: { /* '76' */
            return "lux (Lux)"
        }
        case 77: { /* '77' */
            return "speed (m/s)"
        }
        case 78: { /* '78' */
            return "pressure (Pa)"
        }
        case 79: { /* '79' */
            return "humidity (%)"
        }
        case 8: { /* '8' */
            return "up/down"
        }
        case 80: { /* '80' */
            return "parts/million (ppm)"
        }
        case 81: { /* '81' */
            return "air flow (m³/h)"
        }
        case 82: { /* '82' */
            return "time (s)"
        }
        case 83: { /* '83' */
            return "time (ms)"
        }
        case 84: { /* '84' */
            return "voltage (mV)"
        }
        case 85: { /* '85' */
            return "current (mA)"
        }
        case 86: { /* '86' */
            return "power denisity (W/m²)"
        }
        case 87: { /* '87' */
            return "kelvin/percent (K/%)"
        }
        case 88: { /* '88' */
            return "power (kW)"
        }
        case 89: { /* '89' */
            return "volume flow (l/h)"
        }
        case 9: { /* '9' */
            return "open/close"
        }
        case 90: { /* '90' */
            return "rain amount (l/m²)"
        }
        case 91: { /* '91' */
            return "temperature (°F)"
        }
        case 92: { /* '92' */
            return "wind speed (km/h)"
        }
        case 93: { /* '93' */
            return "absolute humidity (g/m³)"
        }
        case 94: { /* '94' */
            return "concentration (µg/m³)"
        }
        case 95: { /* '95' */
            return "time of day"
        }
        case 96: { /* '96' */
            return "date"
        }
        case 97: { /* '97' */
            return "counter pulses (unsigned)"
        }
        case 98: { /* '98' */
            return "counter timesec (s)"
        }
        case 99: { /* '99' */
            return "counter timemin (min)"
        }
        default: {
            return ""
        }
    }
}

func (e KnxDatapointType) DatapointMainType() KnxDatapointMainType {
    switch e  {
        case 0: { /* '0' */
            return KnxDatapointMainType_DPT_UNKNOWN
        }
        case 1: { /* '1' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 10: { /* '10' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 100: { /* '100' */
            return KnxDatapointMainType_DPT_4_BYTE_UNSIGNED_VALUE
        }
        case 101: { /* '101' */
            return KnxDatapointMainType_DPT_4_BYTE_UNSIGNED_VALUE
        }
        case 102: { /* '102' */
            return KnxDatapointMainType_DPT_4_BYTE_UNSIGNED_VALUE
        }
        case 103: { /* '103' */
            return KnxDatapointMainType_DPT_4_BYTE_SIGNED_VALUE
        }
        case 104: { /* '104' */
            return KnxDatapointMainType_DPT_4_BYTE_SIGNED_VALUE
        }
        case 105: { /* '105' */
            return KnxDatapointMainType_DPT_4_BYTE_SIGNED_VALUE
        }
        case 106: { /* '106' */
            return KnxDatapointMainType_DPT_4_BYTE_SIGNED_VALUE
        }
        case 107: { /* '107' */
            return KnxDatapointMainType_DPT_4_BYTE_SIGNED_VALUE
        }
        case 108: { /* '108' */
            return KnxDatapointMainType_DPT_4_BYTE_SIGNED_VALUE
        }
        case 109: { /* '109' */
            return KnxDatapointMainType_DPT_4_BYTE_SIGNED_VALUE
        }
        case 11: { /* '11' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 110: { /* '110' */
            return KnxDatapointMainType_DPT_4_BYTE_SIGNED_VALUE
        }
        case 111: { /* '111' */
            return KnxDatapointMainType_DPT_4_BYTE_SIGNED_VALUE
        }
        case 112: { /* '112' */
            return KnxDatapointMainType_DPT_4_BYTE_SIGNED_VALUE
        }
        case 113: { /* '113' */
            return KnxDatapointMainType_DPT_4_BYTE_SIGNED_VALUE
        }
        case 114: { /* '114' */
            return KnxDatapointMainType_DPT_4_BYTE_SIGNED_VALUE
        }
        case 115: { /* '115' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 116: { /* '116' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 117: { /* '117' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 118: { /* '118' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 119: { /* '119' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 12: { /* '12' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 120: { /* '120' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 121: { /* '121' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 122: { /* '122' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 123: { /* '123' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 124: { /* '124' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 125: { /* '125' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 126: { /* '126' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 127: { /* '127' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 128: { /* '128' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 129: { /* '129' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 13: { /* '13' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 130: { /* '130' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 131: { /* '131' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 132: { /* '132' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 133: { /* '133' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 134: { /* '134' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 135: { /* '135' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 136: { /* '136' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 137: { /* '137' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 138: { /* '138' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 139: { /* '139' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 14: { /* '14' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 140: { /* '140' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 141: { /* '141' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 142: { /* '142' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 143: { /* '143' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 144: { /* '144' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 145: { /* '145' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 146: { /* '146' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 147: { /* '147' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 148: { /* '148' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 149: { /* '149' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 15: { /* '15' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 150: { /* '150' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 151: { /* '151' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 152: { /* '152' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 153: { /* '153' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 154: { /* '154' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 155: { /* '155' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 156: { /* '156' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 157: { /* '157' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 158: { /* '158' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 159: { /* '159' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 16: { /* '16' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 160: { /* '160' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 161: { /* '161' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 162: { /* '162' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 163: { /* '163' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 164: { /* '164' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 165: { /* '165' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 166: { /* '166' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 167: { /* '167' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 168: { /* '168' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 169: { /* '169' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 17: { /* '17' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 170: { /* '170' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 171: { /* '171' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 172: { /* '172' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 173: { /* '173' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 174: { /* '174' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 175: { /* '175' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 176: { /* '176' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 177: { /* '177' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 178: { /* '178' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 179: { /* '179' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 18: { /* '18' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 180: { /* '180' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 181: { /* '181' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 182: { /* '182' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 183: { /* '183' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 184: { /* '184' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 185: { /* '185' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 186: { /* '186' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 187: { /* '187' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 188: { /* '188' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 189: { /* '189' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 19: { /* '19' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 190: { /* '190' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 191: { /* '191' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 192: { /* '192' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 193: { /* '193' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 194: { /* '194' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 195: { /* '195' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 196: { /* '196' */
            return KnxDatapointMainType_DPT_4_BYTE_FLOAT_VALUE
        }
        case 197: { /* '197' */
            return KnxDatapointMainType_DPT_ENTRANCE_ACCESS
        }
        case 198: { /* '198' */
            return KnxDatapointMainType_DPT_CHARACTER_STRING
        }
        case 199: { /* '199' */
            return KnxDatapointMainType_DPT_CHARACTER_STRING
        }
        case 2: { /* '2' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 20: { /* '20' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 200: { /* '200' */
            return KnxDatapointMainType_DPT_SCENE_NUMBER
        }
        case 201: { /* '201' */
            return KnxDatapointMainType_DPT_SCENE_CONTROL
        }
        case 202: { /* '202' */
            return KnxDatapointMainType_DPT_DATE_TIME
        }
        case 203: { /* '203' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 204: { /* '204' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 205: { /* '205' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 206: { /* '206' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 207: { /* '207' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 208: { /* '208' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 209: { /* '209' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 21: { /* '21' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 210: { /* '210' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 211: { /* '211' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 212: { /* '212' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 213: { /* '213' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 214: { /* '214' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 215: { /* '215' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 216: { /* '216' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 217: { /* '217' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 218: { /* '218' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 219: { /* '219' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 22: { /* '22' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 220: { /* '220' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 221: { /* '221' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 222: { /* '222' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 223: { /* '223' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 224: { /* '224' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 225: { /* '225' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 226: { /* '226' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 227: { /* '227' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 228: { /* '228' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 229: { /* '229' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 23: { /* '23' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 230: { /* '230' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 231: { /* '231' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 232: { /* '232' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 233: { /* '233' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 234: { /* '234' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 235: { /* '235' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 236: { /* '236' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 237: { /* '237' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 238: { /* '238' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 239: { /* '239' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 24: { /* '24' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 240: { /* '240' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 241: { /* '241' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 242: { /* '242' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 243: { /* '243' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 244: { /* '244' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 245: { /* '245' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 246: { /* '246' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 247: { /* '247' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 248: { /* '248' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 249: { /* '249' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 25: { /* '25' */
            return KnxDatapointMainType_DPT_1_BIT_CONTROLLED
        }
        case 250: { /* '250' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 251: { /* '251' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 252: { /* '252' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 253: { /* '253' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 254: { /* '254' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 255: { /* '255' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 256: { /* '256' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 257: { /* '257' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 258: { /* '258' */
            return KnxDatapointMainType_DPT_1_BYTE
        }
        case 259: { /* '259' */
            return KnxDatapointMainType_DPT_8_BIT_SET
        }
        case 26: { /* '26' */
            return KnxDatapointMainType_DPT_1_BIT_CONTROLLED
        }
        case 260: { /* '260' */
            return KnxDatapointMainType_DPT_8_BIT_SET
        }
        case 261: { /* '261' */
            return KnxDatapointMainType_DPT_8_BIT_SET
        }
        case 262: { /* '262' */
            return KnxDatapointMainType_DPT_8_BIT_SET
        }
        case 263: { /* '263' */
            return KnxDatapointMainType_DPT_8_BIT_SET
        }
        case 264: { /* '264' */
            return KnxDatapointMainType_DPT_8_BIT_SET
        }
        case 265: { /* '265' */
            return KnxDatapointMainType_DPT_8_BIT_SET
        }
        case 266: { /* '266' */
            return KnxDatapointMainType_DPT_8_BIT_SET
        }
        case 267: { /* '267' */
            return KnxDatapointMainType_DPT_8_BIT_SET
        }
        case 268: { /* '268' */
            return KnxDatapointMainType_DPT_8_BIT_SET
        }
        case 269: { /* '269' */
            return KnxDatapointMainType_DPT_8_BIT_SET
        }
        case 27: { /* '27' */
            return KnxDatapointMainType_DPT_1_BIT_CONTROLLED
        }
        case 270: { /* '270' */
            return KnxDatapointMainType_DPT_8_BIT_SET
        }
        case 271: { /* '271' */
            return KnxDatapointMainType_DPT_8_BIT_SET
        }
        case 272: { /* '272' */
            return KnxDatapointMainType_DPT_8_BIT_SET
        }
        case 273: { /* '273' */
            return KnxDatapointMainType_DPT_16_BIT_SET
        }
        case 274: { /* '274' */
            return KnxDatapointMainType_DPT_16_BIT_SET
        }
        case 275: { /* '275' */
            return KnxDatapointMainType_DPT_16_BIT_SET
        }
        case 276: { /* '276' */
            return KnxDatapointMainType_DPT_16_BIT_SET
        }
        case 277: { /* '277' */
            return KnxDatapointMainType_DPT_16_BIT_SET
        }
        case 278: { /* '278' */
            return KnxDatapointMainType_DPT_16_BIT_SET
        }
        case 279: { /* '279' */
            return KnxDatapointMainType_DPT_2_BIT_SET
        }
        case 28: { /* '28' */
            return KnxDatapointMainType_DPT_1_BIT_CONTROLLED
        }
        case 280: { /* '280' */
            return KnxDatapointMainType_DPT_2_BIT_SET
        }
        case 281: { /* '281' */
            return KnxDatapointMainType_DPT_2_BIT_SET
        }
        case 282: { /* '282' */
            return KnxDatapointMainType_DPT_2_BIT_SET
        }
        case 283: { /* '283' */
            return KnxDatapointMainType_DPT_2_NIBBLE_SET
        }
        case 284: { /* '284' */
            return KnxDatapointMainType_DPT_8_BIT_SET_2
        }
        case 285: { /* '285' */
            return KnxDatapointMainType_DPT_32_BIT_SET
        }
        case 286: { /* '286' */
            return KnxDatapointMainType_DPT_ELECTRICAL_ENERGY
        }
        case 287: { /* '287' */
            return KnxDatapointMainType_DPT_ELECTRICAL_ENERGY
        }
        case 288: { /* '288' */
            return KnxDatapointMainType_DPT_ELECTRICAL_ENERGY
        }
        case 289: { /* '289' */
            return KnxDatapointMainType_DPT_24_TIMES_CHANNEL_ACTIVATION
        }
        case 29: { /* '29' */
            return KnxDatapointMainType_DPT_1_BIT_CONTROLLED
        }
        case 290: { /* '290' */
            return KnxDatapointMainType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM
        }
        case 291: { /* '291' */
            return KnxDatapointMainType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM
        }
        case 292: { /* '292' */
            return KnxDatapointMainType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM
        }
        case 293: { /* '293' */
            return KnxDatapointMainType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM
        }
        case 294: { /* '294' */
            return KnxDatapointMainType_DPT_DATAPOINT_TYPE_VERSION
        }
        case 295: { /* '295' */
            return KnxDatapointMainType_DPT_ALARM_INFO
        }
        case 296: { /* '296' */
            return KnxDatapointMainType_DPT_3X_2_BYTE_FLOAT_VALUE
        }
        case 297: { /* '297' */
            return KnxDatapointMainType_DPT_3X_2_BYTE_FLOAT_VALUE
        }
        case 298: { /* '298' */
            return KnxDatapointMainType_DPT_SCALING_SPEED
        }
        case 299: { /* '299' */
            return KnxDatapointMainType_DPT_SCALING_SPEED
        }
        case 3: { /* '3' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 30: { /* '30' */
            return KnxDatapointMainType_DPT_1_BIT_CONTROLLED
        }
        case 300: { /* '300' */
            return KnxDatapointMainType_DPT_4_1_1_BYTE_COMBINED_INFORMATION
        }
        case 301: { /* '301' */
            return KnxDatapointMainType_DPT_MBUS_ADDRESS
        }
        case 302: { /* '302' */
            return KnxDatapointMainType_DPT_3_BYTE_COLOUR_RGB
        }
        case 303: { /* '303' */
            return KnxDatapointMainType_DPT_LANGUAGE_CODE_ISO_639_1
        }
        case 304: { /* '304' */
            return KnxDatapointMainType_DPT_SIGNED_VALUE_WITH_CLASSIFICATION_AND_VALIDITY
        }
        case 305: { /* '305' */
            return KnxDatapointMainType_DPT_PRIORITISED_MODE_CONTROL
        }
        case 306: { /* '306' */
            return KnxDatapointMainType_DPT_CONFIGURATION_DIAGNOSTICS_16_BIT
        }
        case 307: { /* '307' */
            return KnxDatapointMainType_DPT_CONFIGURATION_DIAGNOSTICS_8_BIT
        }
        case 308: { /* '308' */
            return KnxDatapointMainType_DPT_POSITIONS
        }
        case 309: { /* '309' */
            return KnxDatapointMainType_DPT_STATUS_32_BIT
        }
        case 31: { /* '31' */
            return KnxDatapointMainType_DPT_1_BIT_CONTROLLED
        }
        case 310: { /* '310' */
            return KnxDatapointMainType_DPT_STATUS_48_BIT
        }
        case 311: { /* '311' */
            return KnxDatapointMainType_DPT_CONVERTER_STATUS
        }
        case 312: { /* '312' */
            return KnxDatapointMainType_DPT_CONVERTER_TEST_RESULT
        }
        case 313: { /* '313' */
            return KnxDatapointMainType_DPT_BATTERY_INFORMATION
        }
        case 314: { /* '314' */
            return KnxDatapointMainType_DPT_BRIGHTNESS_COLOUR_TEMPERATURE_TRANSITION
        }
        case 315: { /* '315' */
            return KnxDatapointMainType_DPT_STATUS_24_BIT
        }
        case 316: { /* '316' */
            return KnxDatapointMainType_DPT_COLOUR_RGBW
        }
        case 317: { /* '317' */
            return KnxDatapointMainType_DPT_RELATIVE_CONTROL_RGBW
        }
        case 318: { /* '318' */
            return KnxDatapointMainType_DPT_RELATIVE_CONTROL_RGB
        }
        case 319: { /* '319' */
            return KnxDatapointMainType_DPT_F32F32
        }
        case 32: { /* '32' */
            return KnxDatapointMainType_DPT_1_BIT_CONTROLLED
        }
        case 320: { /* '320' */
            return KnxDatapointMainType_DPT_F16F16F16F16
        }
        case 321: { /* '321' */
            return KnxDatapointMainType_DPT_F16F16F16F16
        }
        case 33: { /* '33' */
            return KnxDatapointMainType_DPT_1_BIT_CONTROLLED
        }
        case 34: { /* '34' */
            return KnxDatapointMainType_DPT_1_BIT_CONTROLLED
        }
        case 35: { /* '35' */
            return KnxDatapointMainType_DPT_1_BIT_CONTROLLED
        }
        case 36: { /* '36' */
            return KnxDatapointMainType_DPT_1_BIT_CONTROLLED
        }
        case 37: { /* '37' */
            return KnxDatapointMainType_DPT_3_BIT_CONTROLLED
        }
        case 38: { /* '38' */
            return KnxDatapointMainType_DPT_3_BIT_CONTROLLED
        }
        case 39: { /* '39' */
            return KnxDatapointMainType_DPT_CHARACTER
        }
        case 4: { /* '4' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 40: { /* '40' */
            return KnxDatapointMainType_DPT_CHARACTER
        }
        case 41: { /* '41' */
            return KnxDatapointMainType_DPT_8_BIT_UNSIGNED_VALUE
        }
        case 42: { /* '42' */
            return KnxDatapointMainType_DPT_8_BIT_UNSIGNED_VALUE
        }
        case 43: { /* '43' */
            return KnxDatapointMainType_DPT_8_BIT_UNSIGNED_VALUE
        }
        case 44: { /* '44' */
            return KnxDatapointMainType_DPT_8_BIT_UNSIGNED_VALUE
        }
        case 45: { /* '45' */
            return KnxDatapointMainType_DPT_8_BIT_UNSIGNED_VALUE
        }
        case 46: { /* '46' */
            return KnxDatapointMainType_DPT_8_BIT_UNSIGNED_VALUE
        }
        case 47: { /* '47' */
            return KnxDatapointMainType_DPT_8_BIT_UNSIGNED_VALUE
        }
        case 48: { /* '48' */
            return KnxDatapointMainType_DPT_8_BIT_SIGNED_VALUE
        }
        case 49: { /* '49' */
            return KnxDatapointMainType_DPT_8_BIT_SIGNED_VALUE
        }
        case 5: { /* '5' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 50: { /* '50' */
            return KnxDatapointMainType_DPT_8_BIT_SIGNED_VALUE
        }
        case 51: { /* '51' */
            return KnxDatapointMainType_DPT_2_BYTE_UNSIGNED_VALUE
        }
        case 52: { /* '52' */
            return KnxDatapointMainType_DPT_2_BYTE_UNSIGNED_VALUE
        }
        case 53: { /* '53' */
            return KnxDatapointMainType_DPT_2_BYTE_UNSIGNED_VALUE
        }
        case 54: { /* '54' */
            return KnxDatapointMainType_DPT_2_BYTE_UNSIGNED_VALUE
        }
        case 55: { /* '55' */
            return KnxDatapointMainType_DPT_2_BYTE_UNSIGNED_VALUE
        }
        case 56: { /* '56' */
            return KnxDatapointMainType_DPT_2_BYTE_UNSIGNED_VALUE
        }
        case 57: { /* '57' */
            return KnxDatapointMainType_DPT_2_BYTE_UNSIGNED_VALUE
        }
        case 58: { /* '58' */
            return KnxDatapointMainType_DPT_2_BYTE_UNSIGNED_VALUE
        }
        case 59: { /* '59' */
            return KnxDatapointMainType_DPT_2_BYTE_UNSIGNED_VALUE
        }
        case 6: { /* '6' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 60: { /* '60' */
            return KnxDatapointMainType_DPT_2_BYTE_UNSIGNED_VALUE
        }
        case 61: { /* '61' */
            return KnxDatapointMainType_DPT_2_BYTE_UNSIGNED_VALUE
        }
        case 62: { /* '62' */
            return KnxDatapointMainType_DPT_2_BYTE_UNSIGNED_VALUE
        }
        case 63: { /* '63' */
            return KnxDatapointMainType_DPT_2_BYTE_SIGNED_VALUE
        }
        case 64: { /* '64' */
            return KnxDatapointMainType_DPT_2_BYTE_SIGNED_VALUE
        }
        case 65: { /* '65' */
            return KnxDatapointMainType_DPT_2_BYTE_SIGNED_VALUE
        }
        case 66: { /* '66' */
            return KnxDatapointMainType_DPT_2_BYTE_SIGNED_VALUE
        }
        case 67: { /* '67' */
            return KnxDatapointMainType_DPT_2_BYTE_SIGNED_VALUE
        }
        case 68: { /* '68' */
            return KnxDatapointMainType_DPT_2_BYTE_SIGNED_VALUE
        }
        case 69: { /* '69' */
            return KnxDatapointMainType_DPT_2_BYTE_SIGNED_VALUE
        }
        case 7: { /* '7' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 70: { /* '70' */
            return KnxDatapointMainType_DPT_2_BYTE_SIGNED_VALUE
        }
        case 71: { /* '71' */
            return KnxDatapointMainType_DPT_2_BYTE_SIGNED_VALUE
        }
        case 72: { /* '72' */
            return KnxDatapointMainType_DPT_2_BYTE_SIGNED_VALUE
        }
        case 73: { /* '73' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 74: { /* '74' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 75: { /* '75' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 76: { /* '76' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 77: { /* '77' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 78: { /* '78' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 79: { /* '79' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 8: { /* '8' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 80: { /* '80' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 81: { /* '81' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 82: { /* '82' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 83: { /* '83' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 84: { /* '84' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 85: { /* '85' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 86: { /* '86' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 87: { /* '87' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 88: { /* '88' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 89: { /* '89' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 9: { /* '9' */
            return KnxDatapointMainType_DPT_1_BIT
        }
        case 90: { /* '90' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 91: { /* '91' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 92: { /* '92' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 93: { /* '93' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 94: { /* '94' */
            return KnxDatapointMainType_DPT_2_BYTE_FLOAT_VALUE
        }
        case 95: { /* '95' */
            return KnxDatapointMainType_DPT_TIME
        }
        case 96: { /* '96' */
            return KnxDatapointMainType_DPT_DATE
        }
        case 97: { /* '97' */
            return KnxDatapointMainType_DPT_4_BYTE_UNSIGNED_VALUE
        }
        case 98: { /* '98' */
            return KnxDatapointMainType_DPT_4_BYTE_UNSIGNED_VALUE
        }
        case 99: { /* '99' */
            return KnxDatapointMainType_DPT_4_BYTE_UNSIGNED_VALUE
        }
        default: {
            return 0
        }
    }
}
func KnxDatapointTypeByValue(value uint32) KnxDatapointType {
    switch value {
        case 0:
            return KnxDatapointType_DPT_UNKNOWN
        case 1:
            return KnxDatapointType_DPT_Switch
        case 10:
            return KnxDatapointType_DPT_Start
        case 100:
            return KnxDatapointType_DPT_LongTimePeriod_Hrs
        case 101:
            return KnxDatapointType_DPT_VolumeLiquid_Litre
        case 102:
            return KnxDatapointType_DPT_Volume_m_3
        case 103:
            return KnxDatapointType_DPT_Value_4_Count
        case 104:
            return KnxDatapointType_DPT_FlowRate_m3h
        case 105:
            return KnxDatapointType_DPT_ActiveEnergy
        case 106:
            return KnxDatapointType_DPT_ApparantEnergy
        case 107:
            return KnxDatapointType_DPT_ReactiveEnergy
        case 108:
            return KnxDatapointType_DPT_ActiveEnergy_kWh
        case 109:
            return KnxDatapointType_DPT_ApparantEnergy_kVAh
        case 11:
            return KnxDatapointType_DPT_State
        case 110:
            return KnxDatapointType_DPT_ReactiveEnergy_kVARh
        case 111:
            return KnxDatapointType_DPT_ActiveEnergy_MWh
        case 112:
            return KnxDatapointType_DPT_LongDeltaTimeSec
        case 113:
            return KnxDatapointType_DPT_DeltaVolumeLiquid_Litre
        case 114:
            return KnxDatapointType_DPT_DeltaVolume_m_3
        case 115:
            return KnxDatapointType_DPT_Value_Acceleration
        case 116:
            return KnxDatapointType_DPT_Value_Acceleration_Angular
        case 117:
            return KnxDatapointType_DPT_Value_Activation_Energy
        case 118:
            return KnxDatapointType_DPT_Value_Activity
        case 119:
            return KnxDatapointType_DPT_Value_Mol
        case 12:
            return KnxDatapointType_DPT_Invert
        case 120:
            return KnxDatapointType_DPT_Value_Amplitude
        case 121:
            return KnxDatapointType_DPT_Value_AngleRad
        case 122:
            return KnxDatapointType_DPT_Value_AngleDeg
        case 123:
            return KnxDatapointType_DPT_Value_Angular_Momentum
        case 124:
            return KnxDatapointType_DPT_Value_Angular_Velocity
        case 125:
            return KnxDatapointType_DPT_Value_Area
        case 126:
            return KnxDatapointType_DPT_Value_Capacitance
        case 127:
            return KnxDatapointType_DPT_Value_Charge_DensitySurface
        case 128:
            return KnxDatapointType_DPT_Value_Charge_DensityVolume
        case 129:
            return KnxDatapointType_DPT_Value_Compressibility
        case 13:
            return KnxDatapointType_DPT_DimSendStyle
        case 130:
            return KnxDatapointType_DPT_Value_Conductance
        case 131:
            return KnxDatapointType_DPT_Value_Electrical_Conductivity
        case 132:
            return KnxDatapointType_DPT_Value_Density
        case 133:
            return KnxDatapointType_DPT_Value_Electric_Charge
        case 134:
            return KnxDatapointType_DPT_Value_Electric_Current
        case 135:
            return KnxDatapointType_DPT_Value_Electric_CurrentDensity
        case 136:
            return KnxDatapointType_DPT_Value_Electric_DipoleMoment
        case 137:
            return KnxDatapointType_DPT_Value_Electric_Displacement
        case 138:
            return KnxDatapointType_DPT_Value_Electric_FieldStrength
        case 139:
            return KnxDatapointType_DPT_Value_Electric_Flux
        case 14:
            return KnxDatapointType_DPT_InputSource
        case 140:
            return KnxDatapointType_DPT_Value_Electric_FluxDensity
        case 141:
            return KnxDatapointType_DPT_Value_Electric_Polarization
        case 142:
            return KnxDatapointType_DPT_Value_Electric_Potential
        case 143:
            return KnxDatapointType_DPT_Value_Electric_PotentialDifference
        case 144:
            return KnxDatapointType_DPT_Value_ElectromagneticMoment
        case 145:
            return KnxDatapointType_DPT_Value_Electromotive_Force
        case 146:
            return KnxDatapointType_DPT_Value_Energy
        case 147:
            return KnxDatapointType_DPT_Value_Force
        case 148:
            return KnxDatapointType_DPT_Value_Frequency
        case 149:
            return KnxDatapointType_DPT_Value_Angular_Frequency
        case 15:
            return KnxDatapointType_DPT_Reset
        case 150:
            return KnxDatapointType_DPT_Value_Heat_Capacity
        case 151:
            return KnxDatapointType_DPT_Value_Heat_FlowRate
        case 152:
            return KnxDatapointType_DPT_Value_Heat_Quantity
        case 153:
            return KnxDatapointType_DPT_Value_Impedance
        case 154:
            return KnxDatapointType_DPT_Value_Length
        case 155:
            return KnxDatapointType_DPT_Value_Light_Quantity
        case 156:
            return KnxDatapointType_DPT_Value_Luminance
        case 157:
            return KnxDatapointType_DPT_Value_Luminous_Flux
        case 158:
            return KnxDatapointType_DPT_Value_Luminous_Intensity
        case 159:
            return KnxDatapointType_DPT_Value_Magnetic_FieldStrength
        case 16:
            return KnxDatapointType_DPT_Ack
        case 160:
            return KnxDatapointType_DPT_Value_Magnetic_Flux
        case 161:
            return KnxDatapointType_DPT_Value_Magnetic_FluxDensity
        case 162:
            return KnxDatapointType_DPT_Value_Magnetic_Moment
        case 163:
            return KnxDatapointType_DPT_Value_Magnetic_Polarization
        case 164:
            return KnxDatapointType_DPT_Value_Magnetization
        case 165:
            return KnxDatapointType_DPT_Value_MagnetomotiveForce
        case 166:
            return KnxDatapointType_DPT_Value_Mass
        case 167:
            return KnxDatapointType_DPT_Value_MassFlux
        case 168:
            return KnxDatapointType_DPT_Value_Momentum
        case 169:
            return KnxDatapointType_DPT_Value_Phase_AngleRad
        case 17:
            return KnxDatapointType_DPT_Trigger
        case 170:
            return KnxDatapointType_DPT_Value_Phase_AngleDeg
        case 171:
            return KnxDatapointType_DPT_Value_Power
        case 172:
            return KnxDatapointType_DPT_Value_Power_Factor
        case 173:
            return KnxDatapointType_DPT_Value_Pressure
        case 174:
            return KnxDatapointType_DPT_Value_Reactance
        case 175:
            return KnxDatapointType_DPT_Value_Resistance
        case 176:
            return KnxDatapointType_DPT_Value_Resistivity
        case 177:
            return KnxDatapointType_DPT_Value_SelfInductance
        case 178:
            return KnxDatapointType_DPT_Value_SolidAngle
        case 179:
            return KnxDatapointType_DPT_Value_Sound_Intensity
        case 18:
            return KnxDatapointType_DPT_Occupancy
        case 180:
            return KnxDatapointType_DPT_Value_Speed
        case 181:
            return KnxDatapointType_DPT_Value_Stress
        case 182:
            return KnxDatapointType_DPT_Value_Surface_Tension
        case 183:
            return KnxDatapointType_DPT_Value_Common_Temperature
        case 184:
            return KnxDatapointType_DPT_Value_Absolute_Temperature
        case 185:
            return KnxDatapointType_DPT_Value_TemperatureDifference
        case 186:
            return KnxDatapointType_DPT_Value_Thermal_Capacity
        case 187:
            return KnxDatapointType_DPT_Value_Thermal_Conductivity
        case 188:
            return KnxDatapointType_DPT_Value_ThermoelectricPower
        case 189:
            return KnxDatapointType_DPT_Value_Time
        case 19:
            return KnxDatapointType_DPT_Window_Door
        case 190:
            return KnxDatapointType_DPT_Value_Torque
        case 191:
            return KnxDatapointType_DPT_Value_Volume
        case 192:
            return KnxDatapointType_DPT_Value_Volume_Flux
        case 193:
            return KnxDatapointType_DPT_Value_Weight
        case 194:
            return KnxDatapointType_DPT_Value_Work
        case 195:
            return KnxDatapointType_DPT_Volume_Flux_Meter
        case 196:
            return KnxDatapointType_DPT_Volume_Flux_ls
        case 197:
            return KnxDatapointType_DPT_Access_Data
        case 198:
            return KnxDatapointType_DPT_String_ASCII
        case 199:
            return KnxDatapointType_DPT_String_8859_1
        case 2:
            return KnxDatapointType_DPT_Bool
        case 20:
            return KnxDatapointType_DPT_LogicalFunction
        case 200:
            return KnxDatapointType_DPT_SceneNumber
        case 201:
            return KnxDatapointType_DPT_SceneControl
        case 202:
            return KnxDatapointType_DPT_DateTime
        case 203:
            return KnxDatapointType_DPT_SCLOMode
        case 204:
            return KnxDatapointType_DPT_BuildingMode
        case 205:
            return KnxDatapointType_DPT_OccMode
        case 206:
            return KnxDatapointType_DPT_Priority
        case 207:
            return KnxDatapointType_DPT_LightApplicationMode
        case 208:
            return KnxDatapointType_DPT_ApplicationArea
        case 209:
            return KnxDatapointType_DPT_AlarmClassType
        case 21:
            return KnxDatapointType_DPT_Scene_AB
        case 210:
            return KnxDatapointType_DPT_PSUMode
        case 211:
            return KnxDatapointType_DPT_ErrorClass_System
        case 212:
            return KnxDatapointType_DPT_ErrorClass_HVAC
        case 213:
            return KnxDatapointType_DPT_Time_Delay
        case 214:
            return KnxDatapointType_DPT_Beaufort_Wind_Force_Scale
        case 215:
            return KnxDatapointType_DPT_SensorSelect
        case 216:
            return KnxDatapointType_DPT_ActuatorConnectType
        case 217:
            return KnxDatapointType_DPT_Cloud_Cover
        case 218:
            return KnxDatapointType_DPT_PowerReturnMode
        case 219:
            return KnxDatapointType_DPT_FuelType
        case 22:
            return KnxDatapointType_DPT_ShutterBlinds_Mode
        case 220:
            return KnxDatapointType_DPT_BurnerType
        case 221:
            return KnxDatapointType_DPT_HVACMode
        case 222:
            return KnxDatapointType_DPT_DHWMode
        case 223:
            return KnxDatapointType_DPT_LoadPriority
        case 224:
            return KnxDatapointType_DPT_HVACContrMode
        case 225:
            return KnxDatapointType_DPT_HVACEmergMode
        case 226:
            return KnxDatapointType_DPT_ChangeoverMode
        case 227:
            return KnxDatapointType_DPT_ValveMode
        case 228:
            return KnxDatapointType_DPT_DamperMode
        case 229:
            return KnxDatapointType_DPT_HeaterMode
        case 23:
            return KnxDatapointType_DPT_DayNight
        case 230:
            return KnxDatapointType_DPT_FanMode
        case 231:
            return KnxDatapointType_DPT_MasterSlaveMode
        case 232:
            return KnxDatapointType_DPT_StatusRoomSetp
        case 233:
            return KnxDatapointType_DPT_Metering_DeviceType
        case 234:
            return KnxDatapointType_DPT_HumDehumMode
        case 235:
            return KnxDatapointType_DPT_EnableHCStage
        case 236:
            return KnxDatapointType_DPT_ADAType
        case 237:
            return KnxDatapointType_DPT_BackupMode
        case 238:
            return KnxDatapointType_DPT_StartSynchronization
        case 239:
            return KnxDatapointType_DPT_Behaviour_Lock_Unlock
        case 24:
            return KnxDatapointType_DPT_Heat_Cool
        case 240:
            return KnxDatapointType_DPT_Behaviour_Bus_Power_Up_Down
        case 241:
            return KnxDatapointType_DPT_DALI_Fade_Time
        case 242:
            return KnxDatapointType_DPT_BlinkingMode
        case 243:
            return KnxDatapointType_DPT_LightControlMode
        case 244:
            return KnxDatapointType_DPT_SwitchPBModel
        case 245:
            return KnxDatapointType_DPT_PBAction
        case 246:
            return KnxDatapointType_DPT_DimmPBModel
        case 247:
            return KnxDatapointType_DPT_SwitchOnMode
        case 248:
            return KnxDatapointType_DPT_LoadTypeSet
        case 249:
            return KnxDatapointType_DPT_LoadTypeDetected
        case 25:
            return KnxDatapointType_DPT_Switch_Control
        case 250:
            return KnxDatapointType_DPT_Converter_Test_Control
        case 251:
            return KnxDatapointType_DPT_SABExcept_Behaviour
        case 252:
            return KnxDatapointType_DPT_SABBehaviour_Lock_Unlock
        case 253:
            return KnxDatapointType_DPT_SSSBMode
        case 254:
            return KnxDatapointType_DPT_BlindsControlMode
        case 255:
            return KnxDatapointType_DPT_CommMode
        case 256:
            return KnxDatapointType_DPT_AddInfoTypes
        case 257:
            return KnxDatapointType_DPT_RF_ModeSelect
        case 258:
            return KnxDatapointType_DPT_RF_FilterSelect
        case 259:
            return KnxDatapointType_DPT_StatusGen
        case 26:
            return KnxDatapointType_DPT_Bool_Control
        case 260:
            return KnxDatapointType_DPT_Device_Control
        case 261:
            return KnxDatapointType_DPT_ForceSign
        case 262:
            return KnxDatapointType_DPT_ForceSignCool
        case 263:
            return KnxDatapointType_DPT_StatusRHC
        case 264:
            return KnxDatapointType_DPT_StatusSDHWC
        case 265:
            return KnxDatapointType_DPT_FuelTypeSet
        case 266:
            return KnxDatapointType_DPT_StatusRCC
        case 267:
            return KnxDatapointType_DPT_StatusAHU
        case 268:
            return KnxDatapointType_DPT_CombinedStatus_RTSM
        case 269:
            return KnxDatapointType_DPT_LightActuatorErrorInfo
        case 27:
            return KnxDatapointType_DPT_Enable_Control
        case 270:
            return KnxDatapointType_DPT_RF_ModeInfo
        case 271:
            return KnxDatapointType_DPT_RF_FilterInfo
        case 272:
            return KnxDatapointType_DPT_Channel_Activation_8
        case 273:
            return KnxDatapointType_DPT_StatusDHWC
        case 274:
            return KnxDatapointType_DPT_StatusRHCC
        case 275:
            return KnxDatapointType_DPT_CombinedStatus_HVA
        case 276:
            return KnxDatapointType_DPT_CombinedStatus_RTC
        case 277:
            return KnxDatapointType_DPT_Media
        case 278:
            return KnxDatapointType_DPT_Channel_Activation_16
        case 279:
            return KnxDatapointType_DPT_OnOffAction
        case 28:
            return KnxDatapointType_DPT_Ramp_Control
        case 280:
            return KnxDatapointType_DPT_Alarm_Reaction
        case 281:
            return KnxDatapointType_DPT_UpDown_Action
        case 282:
            return KnxDatapointType_DPT_HVAC_PB_Action
        case 283:
            return KnxDatapointType_DPT_DoubleNibble
        case 284:
            return KnxDatapointType_DPT_SceneInfo
        case 285:
            return KnxDatapointType_DPT_CombinedInfoOnOff
        case 286:
            return KnxDatapointType_DPT_ActiveEnergy_V64
        case 287:
            return KnxDatapointType_DPT_ApparantEnergy_V64
        case 288:
            return KnxDatapointType_DPT_ReactiveEnergy_V64
        case 289:
            return KnxDatapointType_DPT_Channel_Activation_24
        case 29:
            return KnxDatapointType_DPT_Alarm_Control
        case 290:
            return KnxDatapointType_DPT_HVACModeNext
        case 291:
            return KnxDatapointType_DPT_DHWModeNext
        case 292:
            return KnxDatapointType_DPT_OccModeNext
        case 293:
            return KnxDatapointType_DPT_BuildingModeNext
        case 294:
            return KnxDatapointType_DPT_Version
        case 295:
            return KnxDatapointType_DPT_AlarmInfo
        case 296:
            return KnxDatapointType_DPT_TempRoomSetpSetF16_3
        case 297:
            return KnxDatapointType_DPT_TempRoomSetpSetShiftF16_3
        case 298:
            return KnxDatapointType_DPT_Scaling_Speed
        case 299:
            return KnxDatapointType_DPT_Scaling_Step_Time
        case 3:
            return KnxDatapointType_DPT_Enable
        case 30:
            return KnxDatapointType_DPT_BinaryValue_Control
        case 300:
            return KnxDatapointType_DPT_MeteringValue
        case 301:
            return KnxDatapointType_DPT_MBus_Address
        case 302:
            return KnxDatapointType_DPT_Colour_RGB
        case 303:
            return KnxDatapointType_DPT_LanguageCodeAlpha2_ASCII
        case 304:
            return KnxDatapointType_DPT_Tariff_ActiveEnergy
        case 305:
            return KnxDatapointType_DPT_Prioritised_Mode_Control
        case 306:
            return KnxDatapointType_DPT_DALI_Control_Gear_Diagnostic
        case 307:
            return KnxDatapointType_DPT_DALI_Diagnostics
        case 308:
            return KnxDatapointType_DPT_CombinedPosition
        case 309:
            return KnxDatapointType_DPT_StatusSAB
        case 31:
            return KnxDatapointType_DPT_Step_Control
        case 310:
            return KnxDatapointType_DPT_Colour_xyY
        case 311:
            return KnxDatapointType_DPT_Converter_Status
        case 312:
            return KnxDatapointType_DPT_Converter_Test_Result
        case 313:
            return KnxDatapointType_DPT_Battery_Info
        case 314:
            return KnxDatapointType_DPT_Brightness_Colour_Temperature_Transition
        case 315:
            return KnxDatapointType_DPT_Brightness_Colour_Temperature_Control
        case 316:
            return KnxDatapointType_DPT_Colour_RGBW
        case 317:
            return KnxDatapointType_DPT_Relative_Control_RGBW
        case 318:
            return KnxDatapointType_DPT_Relative_Control_RGB
        case 319:
            return KnxDatapointType_DPT_GeographicalLocation
        case 32:
            return KnxDatapointType_DPT_Direction1_Control
        case 320:
            return KnxDatapointType_DPT_TempRoomSetpSetF16_4
        case 321:
            return KnxDatapointType_DPT_TempRoomSetpSetShiftF16_4
        case 33:
            return KnxDatapointType_DPT_Direction2_Control
        case 34:
            return KnxDatapointType_DPT_Start_Control
        case 35:
            return KnxDatapointType_DPT_State_Control
        case 36:
            return KnxDatapointType_DPT_Invert_Control
        case 37:
            return KnxDatapointType_DPT_Control_Dimming
        case 38:
            return KnxDatapointType_DPT_Control_Blinds
        case 39:
            return KnxDatapointType_DPT_Char_ASCII
        case 4:
            return KnxDatapointType_DPT_Ramp
        case 40:
            return KnxDatapointType_DPT_Char_8859_1
        case 41:
            return KnxDatapointType_DPT_Scaling
        case 42:
            return KnxDatapointType_DPT_Angle
        case 43:
            return KnxDatapointType_DPT_Percent_U8
        case 44:
            return KnxDatapointType_DPT_DecimalFactor
        case 45:
            return KnxDatapointType_DPT_Tariff
        case 46:
            return KnxDatapointType_DPT_Value_1_Ucount
        case 47:
            return KnxDatapointType_DPT_FanStage
        case 48:
            return KnxDatapointType_DPT_Percent_V8
        case 49:
            return KnxDatapointType_DPT_Value_1_Count
        case 5:
            return KnxDatapointType_DPT_Alarm
        case 50:
            return KnxDatapointType_DPT_Status_Mode3
        case 51:
            return KnxDatapointType_DPT_Value_2_Ucount
        case 52:
            return KnxDatapointType_DPT_TimePeriodMsec
        case 53:
            return KnxDatapointType_DPT_TimePeriod10Msec
        case 54:
            return KnxDatapointType_DPT_TimePeriod100Msec
        case 55:
            return KnxDatapointType_DPT_TimePeriodSec
        case 56:
            return KnxDatapointType_DPT_TimePeriodMin
        case 57:
            return KnxDatapointType_DPT_TimePeriodHrs
        case 58:
            return KnxDatapointType_DPT_PropDataType
        case 59:
            return KnxDatapointType_DPT_Length_mm
        case 6:
            return KnxDatapointType_DPT_BinaryValue
        case 60:
            return KnxDatapointType_DPT_UElCurrentmA
        case 61:
            return KnxDatapointType_DPT_Brightness
        case 62:
            return KnxDatapointType_DPT_Absolute_Colour_Temperature
        case 63:
            return KnxDatapointType_DPT_Value_2_Count
        case 64:
            return KnxDatapointType_DPT_DeltaTimeMsec
        case 65:
            return KnxDatapointType_DPT_DeltaTime10Msec
        case 66:
            return KnxDatapointType_DPT_DeltaTime100Msec
        case 67:
            return KnxDatapointType_DPT_DeltaTimeSec
        case 68:
            return KnxDatapointType_DPT_DeltaTimeMin
        case 69:
            return KnxDatapointType_DPT_DeltaTimeHrs
        case 7:
            return KnxDatapointType_DPT_Step
        case 70:
            return KnxDatapointType_DPT_Percent_V16
        case 71:
            return KnxDatapointType_DPT_Rotation_Angle
        case 72:
            return KnxDatapointType_DPT_Length_m
        case 73:
            return KnxDatapointType_DPT_Value_Temp
        case 74:
            return KnxDatapointType_DPT_Value_Tempd
        case 75:
            return KnxDatapointType_DPT_Value_Tempa
        case 76:
            return KnxDatapointType_DPT_Value_Lux
        case 77:
            return KnxDatapointType_DPT_Value_Wsp
        case 78:
            return KnxDatapointType_DPT_Value_Pres
        case 79:
            return KnxDatapointType_DPT_Value_Humidity
        case 8:
            return KnxDatapointType_DPT_UpDown
        case 80:
            return KnxDatapointType_DPT_Value_AirQuality
        case 81:
            return KnxDatapointType_DPT_Value_AirFlow
        case 82:
            return KnxDatapointType_DPT_Value_Time1
        case 83:
            return KnxDatapointType_DPT_Value_Time2
        case 84:
            return KnxDatapointType_DPT_Value_Volt
        case 85:
            return KnxDatapointType_DPT_Value_Curr
        case 86:
            return KnxDatapointType_DPT_PowerDensity
        case 87:
            return KnxDatapointType_DPT_KelvinPerPercent
        case 88:
            return KnxDatapointType_DPT_Power
        case 89:
            return KnxDatapointType_DPT_Value_Volume_Flow
        case 9:
            return KnxDatapointType_DPT_OpenClose
        case 90:
            return KnxDatapointType_DPT_Rain_Amount
        case 91:
            return KnxDatapointType_DPT_Value_Temp_F
        case 92:
            return KnxDatapointType_DPT_Value_Wsp_kmh
        case 93:
            return KnxDatapointType_DPT_Value_Absolute_Humidity
        case 94:
            return KnxDatapointType_DPT_Concentration_ygm3
        case 95:
            return KnxDatapointType_DPT_TimeOfDay
        case 96:
            return KnxDatapointType_DPT_Date
        case 97:
            return KnxDatapointType_DPT_Value_4_Ucount
        case 98:
            return KnxDatapointType_DPT_LongTimePeriod_Sec
        case 99:
            return KnxDatapointType_DPT_LongTimePeriod_Min
    }
    return 0
}

func KnxDatapointTypeByName(value string) KnxDatapointType {
    switch value {
    case "DPT_UNKNOWN":
        return KnxDatapointType_DPT_UNKNOWN
    case "DPT_Switch":
        return KnxDatapointType_DPT_Switch
    case "DPT_Start":
        return KnxDatapointType_DPT_Start
    case "DPT_LongTimePeriod_Hrs":
        return KnxDatapointType_DPT_LongTimePeriod_Hrs
    case "DPT_VolumeLiquid_Litre":
        return KnxDatapointType_DPT_VolumeLiquid_Litre
    case "DPT_Volume_m_3":
        return KnxDatapointType_DPT_Volume_m_3
    case "DPT_Value_4_Count":
        return KnxDatapointType_DPT_Value_4_Count
    case "DPT_FlowRate_m3h":
        return KnxDatapointType_DPT_FlowRate_m3h
    case "DPT_ActiveEnergy":
        return KnxDatapointType_DPT_ActiveEnergy
    case "DPT_ApparantEnergy":
        return KnxDatapointType_DPT_ApparantEnergy
    case "DPT_ReactiveEnergy":
        return KnxDatapointType_DPT_ReactiveEnergy
    case "DPT_ActiveEnergy_kWh":
        return KnxDatapointType_DPT_ActiveEnergy_kWh
    case "DPT_ApparantEnergy_kVAh":
        return KnxDatapointType_DPT_ApparantEnergy_kVAh
    case "DPT_State":
        return KnxDatapointType_DPT_State
    case "DPT_ReactiveEnergy_kVARh":
        return KnxDatapointType_DPT_ReactiveEnergy_kVARh
    case "DPT_ActiveEnergy_MWh":
        return KnxDatapointType_DPT_ActiveEnergy_MWh
    case "DPT_LongDeltaTimeSec":
        return KnxDatapointType_DPT_LongDeltaTimeSec
    case "DPT_DeltaVolumeLiquid_Litre":
        return KnxDatapointType_DPT_DeltaVolumeLiquid_Litre
    case "DPT_DeltaVolume_m_3":
        return KnxDatapointType_DPT_DeltaVolume_m_3
    case "DPT_Value_Acceleration":
        return KnxDatapointType_DPT_Value_Acceleration
    case "DPT_Value_Acceleration_Angular":
        return KnxDatapointType_DPT_Value_Acceleration_Angular
    case "DPT_Value_Activation_Energy":
        return KnxDatapointType_DPT_Value_Activation_Energy
    case "DPT_Value_Activity":
        return KnxDatapointType_DPT_Value_Activity
    case "DPT_Value_Mol":
        return KnxDatapointType_DPT_Value_Mol
    case "DPT_Invert":
        return KnxDatapointType_DPT_Invert
    case "DPT_Value_Amplitude":
        return KnxDatapointType_DPT_Value_Amplitude
    case "DPT_Value_AngleRad":
        return KnxDatapointType_DPT_Value_AngleRad
    case "DPT_Value_AngleDeg":
        return KnxDatapointType_DPT_Value_AngleDeg
    case "DPT_Value_Angular_Momentum":
        return KnxDatapointType_DPT_Value_Angular_Momentum
    case "DPT_Value_Angular_Velocity":
        return KnxDatapointType_DPT_Value_Angular_Velocity
    case "DPT_Value_Area":
        return KnxDatapointType_DPT_Value_Area
    case "DPT_Value_Capacitance":
        return KnxDatapointType_DPT_Value_Capacitance
    case "DPT_Value_Charge_DensitySurface":
        return KnxDatapointType_DPT_Value_Charge_DensitySurface
    case "DPT_Value_Charge_DensityVolume":
        return KnxDatapointType_DPT_Value_Charge_DensityVolume
    case "DPT_Value_Compressibility":
        return KnxDatapointType_DPT_Value_Compressibility
    case "DPT_DimSendStyle":
        return KnxDatapointType_DPT_DimSendStyle
    case "DPT_Value_Conductance":
        return KnxDatapointType_DPT_Value_Conductance
    case "DPT_Value_Electrical_Conductivity":
        return KnxDatapointType_DPT_Value_Electrical_Conductivity
    case "DPT_Value_Density":
        return KnxDatapointType_DPT_Value_Density
    case "DPT_Value_Electric_Charge":
        return KnxDatapointType_DPT_Value_Electric_Charge
    case "DPT_Value_Electric_Current":
        return KnxDatapointType_DPT_Value_Electric_Current
    case "DPT_Value_Electric_CurrentDensity":
        return KnxDatapointType_DPT_Value_Electric_CurrentDensity
    case "DPT_Value_Electric_DipoleMoment":
        return KnxDatapointType_DPT_Value_Electric_DipoleMoment
    case "DPT_Value_Electric_Displacement":
        return KnxDatapointType_DPT_Value_Electric_Displacement
    case "DPT_Value_Electric_FieldStrength":
        return KnxDatapointType_DPT_Value_Electric_FieldStrength
    case "DPT_Value_Electric_Flux":
        return KnxDatapointType_DPT_Value_Electric_Flux
    case "DPT_InputSource":
        return KnxDatapointType_DPT_InputSource
    case "DPT_Value_Electric_FluxDensity":
        return KnxDatapointType_DPT_Value_Electric_FluxDensity
    case "DPT_Value_Electric_Polarization":
        return KnxDatapointType_DPT_Value_Electric_Polarization
    case "DPT_Value_Electric_Potential":
        return KnxDatapointType_DPT_Value_Electric_Potential
    case "DPT_Value_Electric_PotentialDifference":
        return KnxDatapointType_DPT_Value_Electric_PotentialDifference
    case "DPT_Value_ElectromagneticMoment":
        return KnxDatapointType_DPT_Value_ElectromagneticMoment
    case "DPT_Value_Electromotive_Force":
        return KnxDatapointType_DPT_Value_Electromotive_Force
    case "DPT_Value_Energy":
        return KnxDatapointType_DPT_Value_Energy
    case "DPT_Value_Force":
        return KnxDatapointType_DPT_Value_Force
    case "DPT_Value_Frequency":
        return KnxDatapointType_DPT_Value_Frequency
    case "DPT_Value_Angular_Frequency":
        return KnxDatapointType_DPT_Value_Angular_Frequency
    case "DPT_Reset":
        return KnxDatapointType_DPT_Reset
    case "DPT_Value_Heat_Capacity":
        return KnxDatapointType_DPT_Value_Heat_Capacity
    case "DPT_Value_Heat_FlowRate":
        return KnxDatapointType_DPT_Value_Heat_FlowRate
    case "DPT_Value_Heat_Quantity":
        return KnxDatapointType_DPT_Value_Heat_Quantity
    case "DPT_Value_Impedance":
        return KnxDatapointType_DPT_Value_Impedance
    case "DPT_Value_Length":
        return KnxDatapointType_DPT_Value_Length
    case "DPT_Value_Light_Quantity":
        return KnxDatapointType_DPT_Value_Light_Quantity
    case "DPT_Value_Luminance":
        return KnxDatapointType_DPT_Value_Luminance
    case "DPT_Value_Luminous_Flux":
        return KnxDatapointType_DPT_Value_Luminous_Flux
    case "DPT_Value_Luminous_Intensity":
        return KnxDatapointType_DPT_Value_Luminous_Intensity
    case "DPT_Value_Magnetic_FieldStrength":
        return KnxDatapointType_DPT_Value_Magnetic_FieldStrength
    case "DPT_Ack":
        return KnxDatapointType_DPT_Ack
    case "DPT_Value_Magnetic_Flux":
        return KnxDatapointType_DPT_Value_Magnetic_Flux
    case "DPT_Value_Magnetic_FluxDensity":
        return KnxDatapointType_DPT_Value_Magnetic_FluxDensity
    case "DPT_Value_Magnetic_Moment":
        return KnxDatapointType_DPT_Value_Magnetic_Moment
    case "DPT_Value_Magnetic_Polarization":
        return KnxDatapointType_DPT_Value_Magnetic_Polarization
    case "DPT_Value_Magnetization":
        return KnxDatapointType_DPT_Value_Magnetization
    case "DPT_Value_MagnetomotiveForce":
        return KnxDatapointType_DPT_Value_MagnetomotiveForce
    case "DPT_Value_Mass":
        return KnxDatapointType_DPT_Value_Mass
    case "DPT_Value_MassFlux":
        return KnxDatapointType_DPT_Value_MassFlux
    case "DPT_Value_Momentum":
        return KnxDatapointType_DPT_Value_Momentum
    case "DPT_Value_Phase_AngleRad":
        return KnxDatapointType_DPT_Value_Phase_AngleRad
    case "DPT_Trigger":
        return KnxDatapointType_DPT_Trigger
    case "DPT_Value_Phase_AngleDeg":
        return KnxDatapointType_DPT_Value_Phase_AngleDeg
    case "DPT_Value_Power":
        return KnxDatapointType_DPT_Value_Power
    case "DPT_Value_Power_Factor":
        return KnxDatapointType_DPT_Value_Power_Factor
    case "DPT_Value_Pressure":
        return KnxDatapointType_DPT_Value_Pressure
    case "DPT_Value_Reactance":
        return KnxDatapointType_DPT_Value_Reactance
    case "DPT_Value_Resistance":
        return KnxDatapointType_DPT_Value_Resistance
    case "DPT_Value_Resistivity":
        return KnxDatapointType_DPT_Value_Resistivity
    case "DPT_Value_SelfInductance":
        return KnxDatapointType_DPT_Value_SelfInductance
    case "DPT_Value_SolidAngle":
        return KnxDatapointType_DPT_Value_SolidAngle
    case "DPT_Value_Sound_Intensity":
        return KnxDatapointType_DPT_Value_Sound_Intensity
    case "DPT_Occupancy":
        return KnxDatapointType_DPT_Occupancy
    case "DPT_Value_Speed":
        return KnxDatapointType_DPT_Value_Speed
    case "DPT_Value_Stress":
        return KnxDatapointType_DPT_Value_Stress
    case "DPT_Value_Surface_Tension":
        return KnxDatapointType_DPT_Value_Surface_Tension
    case "DPT_Value_Common_Temperature":
        return KnxDatapointType_DPT_Value_Common_Temperature
    case "DPT_Value_Absolute_Temperature":
        return KnxDatapointType_DPT_Value_Absolute_Temperature
    case "DPT_Value_TemperatureDifference":
        return KnxDatapointType_DPT_Value_TemperatureDifference
    case "DPT_Value_Thermal_Capacity":
        return KnxDatapointType_DPT_Value_Thermal_Capacity
    case "DPT_Value_Thermal_Conductivity":
        return KnxDatapointType_DPT_Value_Thermal_Conductivity
    case "DPT_Value_ThermoelectricPower":
        return KnxDatapointType_DPT_Value_ThermoelectricPower
    case "DPT_Value_Time":
        return KnxDatapointType_DPT_Value_Time
    case "DPT_Window_Door":
        return KnxDatapointType_DPT_Window_Door
    case "DPT_Value_Torque":
        return KnxDatapointType_DPT_Value_Torque
    case "DPT_Value_Volume":
        return KnxDatapointType_DPT_Value_Volume
    case "DPT_Value_Volume_Flux":
        return KnxDatapointType_DPT_Value_Volume_Flux
    case "DPT_Value_Weight":
        return KnxDatapointType_DPT_Value_Weight
    case "DPT_Value_Work":
        return KnxDatapointType_DPT_Value_Work
    case "DPT_Volume_Flux_Meter":
        return KnxDatapointType_DPT_Volume_Flux_Meter
    case "DPT_Volume_Flux_ls":
        return KnxDatapointType_DPT_Volume_Flux_ls
    case "DPT_Access_Data":
        return KnxDatapointType_DPT_Access_Data
    case "DPT_String_ASCII":
        return KnxDatapointType_DPT_String_ASCII
    case "DPT_String_8859_1":
        return KnxDatapointType_DPT_String_8859_1
    case "DPT_Bool":
        return KnxDatapointType_DPT_Bool
    case "DPT_LogicalFunction":
        return KnxDatapointType_DPT_LogicalFunction
    case "DPT_SceneNumber":
        return KnxDatapointType_DPT_SceneNumber
    case "DPT_SceneControl":
        return KnxDatapointType_DPT_SceneControl
    case "DPT_DateTime":
        return KnxDatapointType_DPT_DateTime
    case "DPT_SCLOMode":
        return KnxDatapointType_DPT_SCLOMode
    case "DPT_BuildingMode":
        return KnxDatapointType_DPT_BuildingMode
    case "DPT_OccMode":
        return KnxDatapointType_DPT_OccMode
    case "DPT_Priority":
        return KnxDatapointType_DPT_Priority
    case "DPT_LightApplicationMode":
        return KnxDatapointType_DPT_LightApplicationMode
    case "DPT_ApplicationArea":
        return KnxDatapointType_DPT_ApplicationArea
    case "DPT_AlarmClassType":
        return KnxDatapointType_DPT_AlarmClassType
    case "DPT_Scene_AB":
        return KnxDatapointType_DPT_Scene_AB
    case "DPT_PSUMode":
        return KnxDatapointType_DPT_PSUMode
    case "DPT_ErrorClass_System":
        return KnxDatapointType_DPT_ErrorClass_System
    case "DPT_ErrorClass_HVAC":
        return KnxDatapointType_DPT_ErrorClass_HVAC
    case "DPT_Time_Delay":
        return KnxDatapointType_DPT_Time_Delay
    case "DPT_Beaufort_Wind_Force_Scale":
        return KnxDatapointType_DPT_Beaufort_Wind_Force_Scale
    case "DPT_SensorSelect":
        return KnxDatapointType_DPT_SensorSelect
    case "DPT_ActuatorConnectType":
        return KnxDatapointType_DPT_ActuatorConnectType
    case "DPT_Cloud_Cover":
        return KnxDatapointType_DPT_Cloud_Cover
    case "DPT_PowerReturnMode":
        return KnxDatapointType_DPT_PowerReturnMode
    case "DPT_FuelType":
        return KnxDatapointType_DPT_FuelType
    case "DPT_ShutterBlinds_Mode":
        return KnxDatapointType_DPT_ShutterBlinds_Mode
    case "DPT_BurnerType":
        return KnxDatapointType_DPT_BurnerType
    case "DPT_HVACMode":
        return KnxDatapointType_DPT_HVACMode
    case "DPT_DHWMode":
        return KnxDatapointType_DPT_DHWMode
    case "DPT_LoadPriority":
        return KnxDatapointType_DPT_LoadPriority
    case "DPT_HVACContrMode":
        return KnxDatapointType_DPT_HVACContrMode
    case "DPT_HVACEmergMode":
        return KnxDatapointType_DPT_HVACEmergMode
    case "DPT_ChangeoverMode":
        return KnxDatapointType_DPT_ChangeoverMode
    case "DPT_ValveMode":
        return KnxDatapointType_DPT_ValveMode
    case "DPT_DamperMode":
        return KnxDatapointType_DPT_DamperMode
    case "DPT_HeaterMode":
        return KnxDatapointType_DPT_HeaterMode
    case "DPT_DayNight":
        return KnxDatapointType_DPT_DayNight
    case "DPT_FanMode":
        return KnxDatapointType_DPT_FanMode
    case "DPT_MasterSlaveMode":
        return KnxDatapointType_DPT_MasterSlaveMode
    case "DPT_StatusRoomSetp":
        return KnxDatapointType_DPT_StatusRoomSetp
    case "DPT_Metering_DeviceType":
        return KnxDatapointType_DPT_Metering_DeviceType
    case "DPT_HumDehumMode":
        return KnxDatapointType_DPT_HumDehumMode
    case "DPT_EnableHCStage":
        return KnxDatapointType_DPT_EnableHCStage
    case "DPT_ADAType":
        return KnxDatapointType_DPT_ADAType
    case "DPT_BackupMode":
        return KnxDatapointType_DPT_BackupMode
    case "DPT_StartSynchronization":
        return KnxDatapointType_DPT_StartSynchronization
    case "DPT_Behaviour_Lock_Unlock":
        return KnxDatapointType_DPT_Behaviour_Lock_Unlock
    case "DPT_Heat_Cool":
        return KnxDatapointType_DPT_Heat_Cool
    case "DPT_Behaviour_Bus_Power_Up_Down":
        return KnxDatapointType_DPT_Behaviour_Bus_Power_Up_Down
    case "DPT_DALI_Fade_Time":
        return KnxDatapointType_DPT_DALI_Fade_Time
    case "DPT_BlinkingMode":
        return KnxDatapointType_DPT_BlinkingMode
    case "DPT_LightControlMode":
        return KnxDatapointType_DPT_LightControlMode
    case "DPT_SwitchPBModel":
        return KnxDatapointType_DPT_SwitchPBModel
    case "DPT_PBAction":
        return KnxDatapointType_DPT_PBAction
    case "DPT_DimmPBModel":
        return KnxDatapointType_DPT_DimmPBModel
    case "DPT_SwitchOnMode":
        return KnxDatapointType_DPT_SwitchOnMode
    case "DPT_LoadTypeSet":
        return KnxDatapointType_DPT_LoadTypeSet
    case "DPT_LoadTypeDetected":
        return KnxDatapointType_DPT_LoadTypeDetected
    case "DPT_Switch_Control":
        return KnxDatapointType_DPT_Switch_Control
    case "DPT_Converter_Test_Control":
        return KnxDatapointType_DPT_Converter_Test_Control
    case "DPT_SABExcept_Behaviour":
        return KnxDatapointType_DPT_SABExcept_Behaviour
    case "DPT_SABBehaviour_Lock_Unlock":
        return KnxDatapointType_DPT_SABBehaviour_Lock_Unlock
    case "DPT_SSSBMode":
        return KnxDatapointType_DPT_SSSBMode
    case "DPT_BlindsControlMode":
        return KnxDatapointType_DPT_BlindsControlMode
    case "DPT_CommMode":
        return KnxDatapointType_DPT_CommMode
    case "DPT_AddInfoTypes":
        return KnxDatapointType_DPT_AddInfoTypes
    case "DPT_RF_ModeSelect":
        return KnxDatapointType_DPT_RF_ModeSelect
    case "DPT_RF_FilterSelect":
        return KnxDatapointType_DPT_RF_FilterSelect
    case "DPT_StatusGen":
        return KnxDatapointType_DPT_StatusGen
    case "DPT_Bool_Control":
        return KnxDatapointType_DPT_Bool_Control
    case "DPT_Device_Control":
        return KnxDatapointType_DPT_Device_Control
    case "DPT_ForceSign":
        return KnxDatapointType_DPT_ForceSign
    case "DPT_ForceSignCool":
        return KnxDatapointType_DPT_ForceSignCool
    case "DPT_StatusRHC":
        return KnxDatapointType_DPT_StatusRHC
    case "DPT_StatusSDHWC":
        return KnxDatapointType_DPT_StatusSDHWC
    case "DPT_FuelTypeSet":
        return KnxDatapointType_DPT_FuelTypeSet
    case "DPT_StatusRCC":
        return KnxDatapointType_DPT_StatusRCC
    case "DPT_StatusAHU":
        return KnxDatapointType_DPT_StatusAHU
    case "DPT_CombinedStatus_RTSM":
        return KnxDatapointType_DPT_CombinedStatus_RTSM
    case "DPT_LightActuatorErrorInfo":
        return KnxDatapointType_DPT_LightActuatorErrorInfo
    case "DPT_Enable_Control":
        return KnxDatapointType_DPT_Enable_Control
    case "DPT_RF_ModeInfo":
        return KnxDatapointType_DPT_RF_ModeInfo
    case "DPT_RF_FilterInfo":
        return KnxDatapointType_DPT_RF_FilterInfo
    case "DPT_Channel_Activation_8":
        return KnxDatapointType_DPT_Channel_Activation_8
    case "DPT_StatusDHWC":
        return KnxDatapointType_DPT_StatusDHWC
    case "DPT_StatusRHCC":
        return KnxDatapointType_DPT_StatusRHCC
    case "DPT_CombinedStatus_HVA":
        return KnxDatapointType_DPT_CombinedStatus_HVA
    case "DPT_CombinedStatus_RTC":
        return KnxDatapointType_DPT_CombinedStatus_RTC
    case "DPT_Media":
        return KnxDatapointType_DPT_Media
    case "DPT_Channel_Activation_16":
        return KnxDatapointType_DPT_Channel_Activation_16
    case "DPT_OnOffAction":
        return KnxDatapointType_DPT_OnOffAction
    case "DPT_Ramp_Control":
        return KnxDatapointType_DPT_Ramp_Control
    case "DPT_Alarm_Reaction":
        return KnxDatapointType_DPT_Alarm_Reaction
    case "DPT_UpDown_Action":
        return KnxDatapointType_DPT_UpDown_Action
    case "DPT_HVAC_PB_Action":
        return KnxDatapointType_DPT_HVAC_PB_Action
    case "DPT_DoubleNibble":
        return KnxDatapointType_DPT_DoubleNibble
    case "DPT_SceneInfo":
        return KnxDatapointType_DPT_SceneInfo
    case "DPT_CombinedInfoOnOff":
        return KnxDatapointType_DPT_CombinedInfoOnOff
    case "DPT_ActiveEnergy_V64":
        return KnxDatapointType_DPT_ActiveEnergy_V64
    case "DPT_ApparantEnergy_V64":
        return KnxDatapointType_DPT_ApparantEnergy_V64
    case "DPT_ReactiveEnergy_V64":
        return KnxDatapointType_DPT_ReactiveEnergy_V64
    case "DPT_Channel_Activation_24":
        return KnxDatapointType_DPT_Channel_Activation_24
    case "DPT_Alarm_Control":
        return KnxDatapointType_DPT_Alarm_Control
    case "DPT_HVACModeNext":
        return KnxDatapointType_DPT_HVACModeNext
    case "DPT_DHWModeNext":
        return KnxDatapointType_DPT_DHWModeNext
    case "DPT_OccModeNext":
        return KnxDatapointType_DPT_OccModeNext
    case "DPT_BuildingModeNext":
        return KnxDatapointType_DPT_BuildingModeNext
    case "DPT_Version":
        return KnxDatapointType_DPT_Version
    case "DPT_AlarmInfo":
        return KnxDatapointType_DPT_AlarmInfo
    case "DPT_TempRoomSetpSetF16_3":
        return KnxDatapointType_DPT_TempRoomSetpSetF16_3
    case "DPT_TempRoomSetpSetShiftF16_3":
        return KnxDatapointType_DPT_TempRoomSetpSetShiftF16_3
    case "DPT_Scaling_Speed":
        return KnxDatapointType_DPT_Scaling_Speed
    case "DPT_Scaling_Step_Time":
        return KnxDatapointType_DPT_Scaling_Step_Time
    case "DPT_Enable":
        return KnxDatapointType_DPT_Enable
    case "DPT_BinaryValue_Control":
        return KnxDatapointType_DPT_BinaryValue_Control
    case "DPT_MeteringValue":
        return KnxDatapointType_DPT_MeteringValue
    case "DPT_MBus_Address":
        return KnxDatapointType_DPT_MBus_Address
    case "DPT_Colour_RGB":
        return KnxDatapointType_DPT_Colour_RGB
    case "DPT_LanguageCodeAlpha2_ASCII":
        return KnxDatapointType_DPT_LanguageCodeAlpha2_ASCII
    case "DPT_Tariff_ActiveEnergy":
        return KnxDatapointType_DPT_Tariff_ActiveEnergy
    case "DPT_Prioritised_Mode_Control":
        return KnxDatapointType_DPT_Prioritised_Mode_Control
    case "DPT_DALI_Control_Gear_Diagnostic":
        return KnxDatapointType_DPT_DALI_Control_Gear_Diagnostic
    case "DPT_DALI_Diagnostics":
        return KnxDatapointType_DPT_DALI_Diagnostics
    case "DPT_CombinedPosition":
        return KnxDatapointType_DPT_CombinedPosition
    case "DPT_StatusSAB":
        return KnxDatapointType_DPT_StatusSAB
    case "DPT_Step_Control":
        return KnxDatapointType_DPT_Step_Control
    case "DPT_Colour_xyY":
        return KnxDatapointType_DPT_Colour_xyY
    case "DPT_Converter_Status":
        return KnxDatapointType_DPT_Converter_Status
    case "DPT_Converter_Test_Result":
        return KnxDatapointType_DPT_Converter_Test_Result
    case "DPT_Battery_Info":
        return KnxDatapointType_DPT_Battery_Info
    case "DPT_Brightness_Colour_Temperature_Transition":
        return KnxDatapointType_DPT_Brightness_Colour_Temperature_Transition
    case "DPT_Brightness_Colour_Temperature_Control":
        return KnxDatapointType_DPT_Brightness_Colour_Temperature_Control
    case "DPT_Colour_RGBW":
        return KnxDatapointType_DPT_Colour_RGBW
    case "DPT_Relative_Control_RGBW":
        return KnxDatapointType_DPT_Relative_Control_RGBW
    case "DPT_Relative_Control_RGB":
        return KnxDatapointType_DPT_Relative_Control_RGB
    case "DPT_GeographicalLocation":
        return KnxDatapointType_DPT_GeographicalLocation
    case "DPT_Direction1_Control":
        return KnxDatapointType_DPT_Direction1_Control
    case "DPT_TempRoomSetpSetF16_4":
        return KnxDatapointType_DPT_TempRoomSetpSetF16_4
    case "DPT_TempRoomSetpSetShiftF16_4":
        return KnxDatapointType_DPT_TempRoomSetpSetShiftF16_4
    case "DPT_Direction2_Control":
        return KnxDatapointType_DPT_Direction2_Control
    case "DPT_Start_Control":
        return KnxDatapointType_DPT_Start_Control
    case "DPT_State_Control":
        return KnxDatapointType_DPT_State_Control
    case "DPT_Invert_Control":
        return KnxDatapointType_DPT_Invert_Control
    case "DPT_Control_Dimming":
        return KnxDatapointType_DPT_Control_Dimming
    case "DPT_Control_Blinds":
        return KnxDatapointType_DPT_Control_Blinds
    case "DPT_Char_ASCII":
        return KnxDatapointType_DPT_Char_ASCII
    case "DPT_Ramp":
        return KnxDatapointType_DPT_Ramp
    case "DPT_Char_8859_1":
        return KnxDatapointType_DPT_Char_8859_1
    case "DPT_Scaling":
        return KnxDatapointType_DPT_Scaling
    case "DPT_Angle":
        return KnxDatapointType_DPT_Angle
    case "DPT_Percent_U8":
        return KnxDatapointType_DPT_Percent_U8
    case "DPT_DecimalFactor":
        return KnxDatapointType_DPT_DecimalFactor
    case "DPT_Tariff":
        return KnxDatapointType_DPT_Tariff
    case "DPT_Value_1_Ucount":
        return KnxDatapointType_DPT_Value_1_Ucount
    case "DPT_FanStage":
        return KnxDatapointType_DPT_FanStage
    case "DPT_Percent_V8":
        return KnxDatapointType_DPT_Percent_V8
    case "DPT_Value_1_Count":
        return KnxDatapointType_DPT_Value_1_Count
    case "DPT_Alarm":
        return KnxDatapointType_DPT_Alarm
    case "DPT_Status_Mode3":
        return KnxDatapointType_DPT_Status_Mode3
    case "DPT_Value_2_Ucount":
        return KnxDatapointType_DPT_Value_2_Ucount
    case "DPT_TimePeriodMsec":
        return KnxDatapointType_DPT_TimePeriodMsec
    case "DPT_TimePeriod10Msec":
        return KnxDatapointType_DPT_TimePeriod10Msec
    case "DPT_TimePeriod100Msec":
        return KnxDatapointType_DPT_TimePeriod100Msec
    case "DPT_TimePeriodSec":
        return KnxDatapointType_DPT_TimePeriodSec
    case "DPT_TimePeriodMin":
        return KnxDatapointType_DPT_TimePeriodMin
    case "DPT_TimePeriodHrs":
        return KnxDatapointType_DPT_TimePeriodHrs
    case "DPT_PropDataType":
        return KnxDatapointType_DPT_PropDataType
    case "DPT_Length_mm":
        return KnxDatapointType_DPT_Length_mm
    case "DPT_BinaryValue":
        return KnxDatapointType_DPT_BinaryValue
    case "DPT_UElCurrentmA":
        return KnxDatapointType_DPT_UElCurrentmA
    case "DPT_Brightness":
        return KnxDatapointType_DPT_Brightness
    case "DPT_Absolute_Colour_Temperature":
        return KnxDatapointType_DPT_Absolute_Colour_Temperature
    case "DPT_Value_2_Count":
        return KnxDatapointType_DPT_Value_2_Count
    case "DPT_DeltaTimeMsec":
        return KnxDatapointType_DPT_DeltaTimeMsec
    case "DPT_DeltaTime10Msec":
        return KnxDatapointType_DPT_DeltaTime10Msec
    case "DPT_DeltaTime100Msec":
        return KnxDatapointType_DPT_DeltaTime100Msec
    case "DPT_DeltaTimeSec":
        return KnxDatapointType_DPT_DeltaTimeSec
    case "DPT_DeltaTimeMin":
        return KnxDatapointType_DPT_DeltaTimeMin
    case "DPT_DeltaTimeHrs":
        return KnxDatapointType_DPT_DeltaTimeHrs
    case "DPT_Step":
        return KnxDatapointType_DPT_Step
    case "DPT_Percent_V16":
        return KnxDatapointType_DPT_Percent_V16
    case "DPT_Rotation_Angle":
        return KnxDatapointType_DPT_Rotation_Angle
    case "DPT_Length_m":
        return KnxDatapointType_DPT_Length_m
    case "DPT_Value_Temp":
        return KnxDatapointType_DPT_Value_Temp
    case "DPT_Value_Tempd":
        return KnxDatapointType_DPT_Value_Tempd
    case "DPT_Value_Tempa":
        return KnxDatapointType_DPT_Value_Tempa
    case "DPT_Value_Lux":
        return KnxDatapointType_DPT_Value_Lux
    case "DPT_Value_Wsp":
        return KnxDatapointType_DPT_Value_Wsp
    case "DPT_Value_Pres":
        return KnxDatapointType_DPT_Value_Pres
    case "DPT_Value_Humidity":
        return KnxDatapointType_DPT_Value_Humidity
    case "DPT_UpDown":
        return KnxDatapointType_DPT_UpDown
    case "DPT_Value_AirQuality":
        return KnxDatapointType_DPT_Value_AirQuality
    case "DPT_Value_AirFlow":
        return KnxDatapointType_DPT_Value_AirFlow
    case "DPT_Value_Time1":
        return KnxDatapointType_DPT_Value_Time1
    case "DPT_Value_Time2":
        return KnxDatapointType_DPT_Value_Time2
    case "DPT_Value_Volt":
        return KnxDatapointType_DPT_Value_Volt
    case "DPT_Value_Curr":
        return KnxDatapointType_DPT_Value_Curr
    case "DPT_PowerDensity":
        return KnxDatapointType_DPT_PowerDensity
    case "DPT_KelvinPerPercent":
        return KnxDatapointType_DPT_KelvinPerPercent
    case "DPT_Power":
        return KnxDatapointType_DPT_Power
    case "DPT_Value_Volume_Flow":
        return KnxDatapointType_DPT_Value_Volume_Flow
    case "DPT_OpenClose":
        return KnxDatapointType_DPT_OpenClose
    case "DPT_Rain_Amount":
        return KnxDatapointType_DPT_Rain_Amount
    case "DPT_Value_Temp_F":
        return KnxDatapointType_DPT_Value_Temp_F
    case "DPT_Value_Wsp_kmh":
        return KnxDatapointType_DPT_Value_Wsp_kmh
    case "DPT_Value_Absolute_Humidity":
        return KnxDatapointType_DPT_Value_Absolute_Humidity
    case "DPT_Concentration_ygm3":
        return KnxDatapointType_DPT_Concentration_ygm3
    case "DPT_TimeOfDay":
        return KnxDatapointType_DPT_TimeOfDay
    case "DPT_Date":
        return KnxDatapointType_DPT_Date
    case "DPT_Value_4_Ucount":
        return KnxDatapointType_DPT_Value_4_Ucount
    case "DPT_LongTimePeriod_Sec":
        return KnxDatapointType_DPT_LongTimePeriod_Sec
    case "DPT_LongTimePeriod_Min":
        return KnxDatapointType_DPT_LongTimePeriod_Min
    }
    return 0
}

func CastKnxDatapointType(structType interface{}) KnxDatapointType {
    castFunc := func(typ interface{}) KnxDatapointType {
        if sKnxDatapointType, ok := typ.(KnxDatapointType); ok {
            return sKnxDatapointType
        }
        return 0
    }
    return castFunc(structType)
}

func (m KnxDatapointType) LengthInBits() uint16 {
    return 32
}

func (m KnxDatapointType) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KnxDatapointTypeParse(io *utils.ReadBuffer) (KnxDatapointType, error) {
    val, err := io.ReadUint32(32)
    if err != nil {
        return 0, nil
    }
    return KnxDatapointTypeByValue(val), nil
}

func (e KnxDatapointType) Serialize(io utils.WriteBuffer) error {
    err := io.WriteUint32(32, uint32(e))
    return err
}

func (e KnxDatapointType) String() string {
    switch e {
    case KnxDatapointType_DPT_UNKNOWN:
        return "DPT_UNKNOWN"
    case KnxDatapointType_DPT_Switch:
        return "DPT_Switch"
    case KnxDatapointType_DPT_Start:
        return "DPT_Start"
    case KnxDatapointType_DPT_LongTimePeriod_Hrs:
        return "DPT_LongTimePeriod_Hrs"
    case KnxDatapointType_DPT_VolumeLiquid_Litre:
        return "DPT_VolumeLiquid_Litre"
    case KnxDatapointType_DPT_Volume_m_3:
        return "DPT_Volume_m_3"
    case KnxDatapointType_DPT_Value_4_Count:
        return "DPT_Value_4_Count"
    case KnxDatapointType_DPT_FlowRate_m3h:
        return "DPT_FlowRate_m3h"
    case KnxDatapointType_DPT_ActiveEnergy:
        return "DPT_ActiveEnergy"
    case KnxDatapointType_DPT_ApparantEnergy:
        return "DPT_ApparantEnergy"
    case KnxDatapointType_DPT_ReactiveEnergy:
        return "DPT_ReactiveEnergy"
    case KnxDatapointType_DPT_ActiveEnergy_kWh:
        return "DPT_ActiveEnergy_kWh"
    case KnxDatapointType_DPT_ApparantEnergy_kVAh:
        return "DPT_ApparantEnergy_kVAh"
    case KnxDatapointType_DPT_State:
        return "DPT_State"
    case KnxDatapointType_DPT_ReactiveEnergy_kVARh:
        return "DPT_ReactiveEnergy_kVARh"
    case KnxDatapointType_DPT_ActiveEnergy_MWh:
        return "DPT_ActiveEnergy_MWh"
    case KnxDatapointType_DPT_LongDeltaTimeSec:
        return "DPT_LongDeltaTimeSec"
    case KnxDatapointType_DPT_DeltaVolumeLiquid_Litre:
        return "DPT_DeltaVolumeLiquid_Litre"
    case KnxDatapointType_DPT_DeltaVolume_m_3:
        return "DPT_DeltaVolume_m_3"
    case KnxDatapointType_DPT_Value_Acceleration:
        return "DPT_Value_Acceleration"
    case KnxDatapointType_DPT_Value_Acceleration_Angular:
        return "DPT_Value_Acceleration_Angular"
    case KnxDatapointType_DPT_Value_Activation_Energy:
        return "DPT_Value_Activation_Energy"
    case KnxDatapointType_DPT_Value_Activity:
        return "DPT_Value_Activity"
    case KnxDatapointType_DPT_Value_Mol:
        return "DPT_Value_Mol"
    case KnxDatapointType_DPT_Invert:
        return "DPT_Invert"
    case KnxDatapointType_DPT_Value_Amplitude:
        return "DPT_Value_Amplitude"
    case KnxDatapointType_DPT_Value_AngleRad:
        return "DPT_Value_AngleRad"
    case KnxDatapointType_DPT_Value_AngleDeg:
        return "DPT_Value_AngleDeg"
    case KnxDatapointType_DPT_Value_Angular_Momentum:
        return "DPT_Value_Angular_Momentum"
    case KnxDatapointType_DPT_Value_Angular_Velocity:
        return "DPT_Value_Angular_Velocity"
    case KnxDatapointType_DPT_Value_Area:
        return "DPT_Value_Area"
    case KnxDatapointType_DPT_Value_Capacitance:
        return "DPT_Value_Capacitance"
    case KnxDatapointType_DPT_Value_Charge_DensitySurface:
        return "DPT_Value_Charge_DensitySurface"
    case KnxDatapointType_DPT_Value_Charge_DensityVolume:
        return "DPT_Value_Charge_DensityVolume"
    case KnxDatapointType_DPT_Value_Compressibility:
        return "DPT_Value_Compressibility"
    case KnxDatapointType_DPT_DimSendStyle:
        return "DPT_DimSendStyle"
    case KnxDatapointType_DPT_Value_Conductance:
        return "DPT_Value_Conductance"
    case KnxDatapointType_DPT_Value_Electrical_Conductivity:
        return "DPT_Value_Electrical_Conductivity"
    case KnxDatapointType_DPT_Value_Density:
        return "DPT_Value_Density"
    case KnxDatapointType_DPT_Value_Electric_Charge:
        return "DPT_Value_Electric_Charge"
    case KnxDatapointType_DPT_Value_Electric_Current:
        return "DPT_Value_Electric_Current"
    case KnxDatapointType_DPT_Value_Electric_CurrentDensity:
        return "DPT_Value_Electric_CurrentDensity"
    case KnxDatapointType_DPT_Value_Electric_DipoleMoment:
        return "DPT_Value_Electric_DipoleMoment"
    case KnxDatapointType_DPT_Value_Electric_Displacement:
        return "DPT_Value_Electric_Displacement"
    case KnxDatapointType_DPT_Value_Electric_FieldStrength:
        return "DPT_Value_Electric_FieldStrength"
    case KnxDatapointType_DPT_Value_Electric_Flux:
        return "DPT_Value_Electric_Flux"
    case KnxDatapointType_DPT_InputSource:
        return "DPT_InputSource"
    case KnxDatapointType_DPT_Value_Electric_FluxDensity:
        return "DPT_Value_Electric_FluxDensity"
    case KnxDatapointType_DPT_Value_Electric_Polarization:
        return "DPT_Value_Electric_Polarization"
    case KnxDatapointType_DPT_Value_Electric_Potential:
        return "DPT_Value_Electric_Potential"
    case KnxDatapointType_DPT_Value_Electric_PotentialDifference:
        return "DPT_Value_Electric_PotentialDifference"
    case KnxDatapointType_DPT_Value_ElectromagneticMoment:
        return "DPT_Value_ElectromagneticMoment"
    case KnxDatapointType_DPT_Value_Electromotive_Force:
        return "DPT_Value_Electromotive_Force"
    case KnxDatapointType_DPT_Value_Energy:
        return "DPT_Value_Energy"
    case KnxDatapointType_DPT_Value_Force:
        return "DPT_Value_Force"
    case KnxDatapointType_DPT_Value_Frequency:
        return "DPT_Value_Frequency"
    case KnxDatapointType_DPT_Value_Angular_Frequency:
        return "DPT_Value_Angular_Frequency"
    case KnxDatapointType_DPT_Reset:
        return "DPT_Reset"
    case KnxDatapointType_DPT_Value_Heat_Capacity:
        return "DPT_Value_Heat_Capacity"
    case KnxDatapointType_DPT_Value_Heat_FlowRate:
        return "DPT_Value_Heat_FlowRate"
    case KnxDatapointType_DPT_Value_Heat_Quantity:
        return "DPT_Value_Heat_Quantity"
    case KnxDatapointType_DPT_Value_Impedance:
        return "DPT_Value_Impedance"
    case KnxDatapointType_DPT_Value_Length:
        return "DPT_Value_Length"
    case KnxDatapointType_DPT_Value_Light_Quantity:
        return "DPT_Value_Light_Quantity"
    case KnxDatapointType_DPT_Value_Luminance:
        return "DPT_Value_Luminance"
    case KnxDatapointType_DPT_Value_Luminous_Flux:
        return "DPT_Value_Luminous_Flux"
    case KnxDatapointType_DPT_Value_Luminous_Intensity:
        return "DPT_Value_Luminous_Intensity"
    case KnxDatapointType_DPT_Value_Magnetic_FieldStrength:
        return "DPT_Value_Magnetic_FieldStrength"
    case KnxDatapointType_DPT_Ack:
        return "DPT_Ack"
    case KnxDatapointType_DPT_Value_Magnetic_Flux:
        return "DPT_Value_Magnetic_Flux"
    case KnxDatapointType_DPT_Value_Magnetic_FluxDensity:
        return "DPT_Value_Magnetic_FluxDensity"
    case KnxDatapointType_DPT_Value_Magnetic_Moment:
        return "DPT_Value_Magnetic_Moment"
    case KnxDatapointType_DPT_Value_Magnetic_Polarization:
        return "DPT_Value_Magnetic_Polarization"
    case KnxDatapointType_DPT_Value_Magnetization:
        return "DPT_Value_Magnetization"
    case KnxDatapointType_DPT_Value_MagnetomotiveForce:
        return "DPT_Value_MagnetomotiveForce"
    case KnxDatapointType_DPT_Value_Mass:
        return "DPT_Value_Mass"
    case KnxDatapointType_DPT_Value_MassFlux:
        return "DPT_Value_MassFlux"
    case KnxDatapointType_DPT_Value_Momentum:
        return "DPT_Value_Momentum"
    case KnxDatapointType_DPT_Value_Phase_AngleRad:
        return "DPT_Value_Phase_AngleRad"
    case KnxDatapointType_DPT_Trigger:
        return "DPT_Trigger"
    case KnxDatapointType_DPT_Value_Phase_AngleDeg:
        return "DPT_Value_Phase_AngleDeg"
    case KnxDatapointType_DPT_Value_Power:
        return "DPT_Value_Power"
    case KnxDatapointType_DPT_Value_Power_Factor:
        return "DPT_Value_Power_Factor"
    case KnxDatapointType_DPT_Value_Pressure:
        return "DPT_Value_Pressure"
    case KnxDatapointType_DPT_Value_Reactance:
        return "DPT_Value_Reactance"
    case KnxDatapointType_DPT_Value_Resistance:
        return "DPT_Value_Resistance"
    case KnxDatapointType_DPT_Value_Resistivity:
        return "DPT_Value_Resistivity"
    case KnxDatapointType_DPT_Value_SelfInductance:
        return "DPT_Value_SelfInductance"
    case KnxDatapointType_DPT_Value_SolidAngle:
        return "DPT_Value_SolidAngle"
    case KnxDatapointType_DPT_Value_Sound_Intensity:
        return "DPT_Value_Sound_Intensity"
    case KnxDatapointType_DPT_Occupancy:
        return "DPT_Occupancy"
    case KnxDatapointType_DPT_Value_Speed:
        return "DPT_Value_Speed"
    case KnxDatapointType_DPT_Value_Stress:
        return "DPT_Value_Stress"
    case KnxDatapointType_DPT_Value_Surface_Tension:
        return "DPT_Value_Surface_Tension"
    case KnxDatapointType_DPT_Value_Common_Temperature:
        return "DPT_Value_Common_Temperature"
    case KnxDatapointType_DPT_Value_Absolute_Temperature:
        return "DPT_Value_Absolute_Temperature"
    case KnxDatapointType_DPT_Value_TemperatureDifference:
        return "DPT_Value_TemperatureDifference"
    case KnxDatapointType_DPT_Value_Thermal_Capacity:
        return "DPT_Value_Thermal_Capacity"
    case KnxDatapointType_DPT_Value_Thermal_Conductivity:
        return "DPT_Value_Thermal_Conductivity"
    case KnxDatapointType_DPT_Value_ThermoelectricPower:
        return "DPT_Value_ThermoelectricPower"
    case KnxDatapointType_DPT_Value_Time:
        return "DPT_Value_Time"
    case KnxDatapointType_DPT_Window_Door:
        return "DPT_Window_Door"
    case KnxDatapointType_DPT_Value_Torque:
        return "DPT_Value_Torque"
    case KnxDatapointType_DPT_Value_Volume:
        return "DPT_Value_Volume"
    case KnxDatapointType_DPT_Value_Volume_Flux:
        return "DPT_Value_Volume_Flux"
    case KnxDatapointType_DPT_Value_Weight:
        return "DPT_Value_Weight"
    case KnxDatapointType_DPT_Value_Work:
        return "DPT_Value_Work"
    case KnxDatapointType_DPT_Volume_Flux_Meter:
        return "DPT_Volume_Flux_Meter"
    case KnxDatapointType_DPT_Volume_Flux_ls:
        return "DPT_Volume_Flux_ls"
    case KnxDatapointType_DPT_Access_Data:
        return "DPT_Access_Data"
    case KnxDatapointType_DPT_String_ASCII:
        return "DPT_String_ASCII"
    case KnxDatapointType_DPT_String_8859_1:
        return "DPT_String_8859_1"
    case KnxDatapointType_DPT_Bool:
        return "DPT_Bool"
    case KnxDatapointType_DPT_LogicalFunction:
        return "DPT_LogicalFunction"
    case KnxDatapointType_DPT_SceneNumber:
        return "DPT_SceneNumber"
    case KnxDatapointType_DPT_SceneControl:
        return "DPT_SceneControl"
    case KnxDatapointType_DPT_DateTime:
        return "DPT_DateTime"
    case KnxDatapointType_DPT_SCLOMode:
        return "DPT_SCLOMode"
    case KnxDatapointType_DPT_BuildingMode:
        return "DPT_BuildingMode"
    case KnxDatapointType_DPT_OccMode:
        return "DPT_OccMode"
    case KnxDatapointType_DPT_Priority:
        return "DPT_Priority"
    case KnxDatapointType_DPT_LightApplicationMode:
        return "DPT_LightApplicationMode"
    case KnxDatapointType_DPT_ApplicationArea:
        return "DPT_ApplicationArea"
    case KnxDatapointType_DPT_AlarmClassType:
        return "DPT_AlarmClassType"
    case KnxDatapointType_DPT_Scene_AB:
        return "DPT_Scene_AB"
    case KnxDatapointType_DPT_PSUMode:
        return "DPT_PSUMode"
    case KnxDatapointType_DPT_ErrorClass_System:
        return "DPT_ErrorClass_System"
    case KnxDatapointType_DPT_ErrorClass_HVAC:
        return "DPT_ErrorClass_HVAC"
    case KnxDatapointType_DPT_Time_Delay:
        return "DPT_Time_Delay"
    case KnxDatapointType_DPT_Beaufort_Wind_Force_Scale:
        return "DPT_Beaufort_Wind_Force_Scale"
    case KnxDatapointType_DPT_SensorSelect:
        return "DPT_SensorSelect"
    case KnxDatapointType_DPT_ActuatorConnectType:
        return "DPT_ActuatorConnectType"
    case KnxDatapointType_DPT_Cloud_Cover:
        return "DPT_Cloud_Cover"
    case KnxDatapointType_DPT_PowerReturnMode:
        return "DPT_PowerReturnMode"
    case KnxDatapointType_DPT_FuelType:
        return "DPT_FuelType"
    case KnxDatapointType_DPT_ShutterBlinds_Mode:
        return "DPT_ShutterBlinds_Mode"
    case KnxDatapointType_DPT_BurnerType:
        return "DPT_BurnerType"
    case KnxDatapointType_DPT_HVACMode:
        return "DPT_HVACMode"
    case KnxDatapointType_DPT_DHWMode:
        return "DPT_DHWMode"
    case KnxDatapointType_DPT_LoadPriority:
        return "DPT_LoadPriority"
    case KnxDatapointType_DPT_HVACContrMode:
        return "DPT_HVACContrMode"
    case KnxDatapointType_DPT_HVACEmergMode:
        return "DPT_HVACEmergMode"
    case KnxDatapointType_DPT_ChangeoverMode:
        return "DPT_ChangeoverMode"
    case KnxDatapointType_DPT_ValveMode:
        return "DPT_ValveMode"
    case KnxDatapointType_DPT_DamperMode:
        return "DPT_DamperMode"
    case KnxDatapointType_DPT_HeaterMode:
        return "DPT_HeaterMode"
    case KnxDatapointType_DPT_DayNight:
        return "DPT_DayNight"
    case KnxDatapointType_DPT_FanMode:
        return "DPT_FanMode"
    case KnxDatapointType_DPT_MasterSlaveMode:
        return "DPT_MasterSlaveMode"
    case KnxDatapointType_DPT_StatusRoomSetp:
        return "DPT_StatusRoomSetp"
    case KnxDatapointType_DPT_Metering_DeviceType:
        return "DPT_Metering_DeviceType"
    case KnxDatapointType_DPT_HumDehumMode:
        return "DPT_HumDehumMode"
    case KnxDatapointType_DPT_EnableHCStage:
        return "DPT_EnableHCStage"
    case KnxDatapointType_DPT_ADAType:
        return "DPT_ADAType"
    case KnxDatapointType_DPT_BackupMode:
        return "DPT_BackupMode"
    case KnxDatapointType_DPT_StartSynchronization:
        return "DPT_StartSynchronization"
    case KnxDatapointType_DPT_Behaviour_Lock_Unlock:
        return "DPT_Behaviour_Lock_Unlock"
    case KnxDatapointType_DPT_Heat_Cool:
        return "DPT_Heat_Cool"
    case KnxDatapointType_DPT_Behaviour_Bus_Power_Up_Down:
        return "DPT_Behaviour_Bus_Power_Up_Down"
    case KnxDatapointType_DPT_DALI_Fade_Time:
        return "DPT_DALI_Fade_Time"
    case KnxDatapointType_DPT_BlinkingMode:
        return "DPT_BlinkingMode"
    case KnxDatapointType_DPT_LightControlMode:
        return "DPT_LightControlMode"
    case KnxDatapointType_DPT_SwitchPBModel:
        return "DPT_SwitchPBModel"
    case KnxDatapointType_DPT_PBAction:
        return "DPT_PBAction"
    case KnxDatapointType_DPT_DimmPBModel:
        return "DPT_DimmPBModel"
    case KnxDatapointType_DPT_SwitchOnMode:
        return "DPT_SwitchOnMode"
    case KnxDatapointType_DPT_LoadTypeSet:
        return "DPT_LoadTypeSet"
    case KnxDatapointType_DPT_LoadTypeDetected:
        return "DPT_LoadTypeDetected"
    case KnxDatapointType_DPT_Switch_Control:
        return "DPT_Switch_Control"
    case KnxDatapointType_DPT_Converter_Test_Control:
        return "DPT_Converter_Test_Control"
    case KnxDatapointType_DPT_SABExcept_Behaviour:
        return "DPT_SABExcept_Behaviour"
    case KnxDatapointType_DPT_SABBehaviour_Lock_Unlock:
        return "DPT_SABBehaviour_Lock_Unlock"
    case KnxDatapointType_DPT_SSSBMode:
        return "DPT_SSSBMode"
    case KnxDatapointType_DPT_BlindsControlMode:
        return "DPT_BlindsControlMode"
    case KnxDatapointType_DPT_CommMode:
        return "DPT_CommMode"
    case KnxDatapointType_DPT_AddInfoTypes:
        return "DPT_AddInfoTypes"
    case KnxDatapointType_DPT_RF_ModeSelect:
        return "DPT_RF_ModeSelect"
    case KnxDatapointType_DPT_RF_FilterSelect:
        return "DPT_RF_FilterSelect"
    case KnxDatapointType_DPT_StatusGen:
        return "DPT_StatusGen"
    case KnxDatapointType_DPT_Bool_Control:
        return "DPT_Bool_Control"
    case KnxDatapointType_DPT_Device_Control:
        return "DPT_Device_Control"
    case KnxDatapointType_DPT_ForceSign:
        return "DPT_ForceSign"
    case KnxDatapointType_DPT_ForceSignCool:
        return "DPT_ForceSignCool"
    case KnxDatapointType_DPT_StatusRHC:
        return "DPT_StatusRHC"
    case KnxDatapointType_DPT_StatusSDHWC:
        return "DPT_StatusSDHWC"
    case KnxDatapointType_DPT_FuelTypeSet:
        return "DPT_FuelTypeSet"
    case KnxDatapointType_DPT_StatusRCC:
        return "DPT_StatusRCC"
    case KnxDatapointType_DPT_StatusAHU:
        return "DPT_StatusAHU"
    case KnxDatapointType_DPT_CombinedStatus_RTSM:
        return "DPT_CombinedStatus_RTSM"
    case KnxDatapointType_DPT_LightActuatorErrorInfo:
        return "DPT_LightActuatorErrorInfo"
    case KnxDatapointType_DPT_Enable_Control:
        return "DPT_Enable_Control"
    case KnxDatapointType_DPT_RF_ModeInfo:
        return "DPT_RF_ModeInfo"
    case KnxDatapointType_DPT_RF_FilterInfo:
        return "DPT_RF_FilterInfo"
    case KnxDatapointType_DPT_Channel_Activation_8:
        return "DPT_Channel_Activation_8"
    case KnxDatapointType_DPT_StatusDHWC:
        return "DPT_StatusDHWC"
    case KnxDatapointType_DPT_StatusRHCC:
        return "DPT_StatusRHCC"
    case KnxDatapointType_DPT_CombinedStatus_HVA:
        return "DPT_CombinedStatus_HVA"
    case KnxDatapointType_DPT_CombinedStatus_RTC:
        return "DPT_CombinedStatus_RTC"
    case KnxDatapointType_DPT_Media:
        return "DPT_Media"
    case KnxDatapointType_DPT_Channel_Activation_16:
        return "DPT_Channel_Activation_16"
    case KnxDatapointType_DPT_OnOffAction:
        return "DPT_OnOffAction"
    case KnxDatapointType_DPT_Ramp_Control:
        return "DPT_Ramp_Control"
    case KnxDatapointType_DPT_Alarm_Reaction:
        return "DPT_Alarm_Reaction"
    case KnxDatapointType_DPT_UpDown_Action:
        return "DPT_UpDown_Action"
    case KnxDatapointType_DPT_HVAC_PB_Action:
        return "DPT_HVAC_PB_Action"
    case KnxDatapointType_DPT_DoubleNibble:
        return "DPT_DoubleNibble"
    case KnxDatapointType_DPT_SceneInfo:
        return "DPT_SceneInfo"
    case KnxDatapointType_DPT_CombinedInfoOnOff:
        return "DPT_CombinedInfoOnOff"
    case KnxDatapointType_DPT_ActiveEnergy_V64:
        return "DPT_ActiveEnergy_V64"
    case KnxDatapointType_DPT_ApparantEnergy_V64:
        return "DPT_ApparantEnergy_V64"
    case KnxDatapointType_DPT_ReactiveEnergy_V64:
        return "DPT_ReactiveEnergy_V64"
    case KnxDatapointType_DPT_Channel_Activation_24:
        return "DPT_Channel_Activation_24"
    case KnxDatapointType_DPT_Alarm_Control:
        return "DPT_Alarm_Control"
    case KnxDatapointType_DPT_HVACModeNext:
        return "DPT_HVACModeNext"
    case KnxDatapointType_DPT_DHWModeNext:
        return "DPT_DHWModeNext"
    case KnxDatapointType_DPT_OccModeNext:
        return "DPT_OccModeNext"
    case KnxDatapointType_DPT_BuildingModeNext:
        return "DPT_BuildingModeNext"
    case KnxDatapointType_DPT_Version:
        return "DPT_Version"
    case KnxDatapointType_DPT_AlarmInfo:
        return "DPT_AlarmInfo"
    case KnxDatapointType_DPT_TempRoomSetpSetF16_3:
        return "DPT_TempRoomSetpSetF16_3"
    case KnxDatapointType_DPT_TempRoomSetpSetShiftF16_3:
        return "DPT_TempRoomSetpSetShiftF16_3"
    case KnxDatapointType_DPT_Scaling_Speed:
        return "DPT_Scaling_Speed"
    case KnxDatapointType_DPT_Scaling_Step_Time:
        return "DPT_Scaling_Step_Time"
    case KnxDatapointType_DPT_Enable:
        return "DPT_Enable"
    case KnxDatapointType_DPT_BinaryValue_Control:
        return "DPT_BinaryValue_Control"
    case KnxDatapointType_DPT_MeteringValue:
        return "DPT_MeteringValue"
    case KnxDatapointType_DPT_MBus_Address:
        return "DPT_MBus_Address"
    case KnxDatapointType_DPT_Colour_RGB:
        return "DPT_Colour_RGB"
    case KnxDatapointType_DPT_LanguageCodeAlpha2_ASCII:
        return "DPT_LanguageCodeAlpha2_ASCII"
    case KnxDatapointType_DPT_Tariff_ActiveEnergy:
        return "DPT_Tariff_ActiveEnergy"
    case KnxDatapointType_DPT_Prioritised_Mode_Control:
        return "DPT_Prioritised_Mode_Control"
    case KnxDatapointType_DPT_DALI_Control_Gear_Diagnostic:
        return "DPT_DALI_Control_Gear_Diagnostic"
    case KnxDatapointType_DPT_DALI_Diagnostics:
        return "DPT_DALI_Diagnostics"
    case KnxDatapointType_DPT_CombinedPosition:
        return "DPT_CombinedPosition"
    case KnxDatapointType_DPT_StatusSAB:
        return "DPT_StatusSAB"
    case KnxDatapointType_DPT_Step_Control:
        return "DPT_Step_Control"
    case KnxDatapointType_DPT_Colour_xyY:
        return "DPT_Colour_xyY"
    case KnxDatapointType_DPT_Converter_Status:
        return "DPT_Converter_Status"
    case KnxDatapointType_DPT_Converter_Test_Result:
        return "DPT_Converter_Test_Result"
    case KnxDatapointType_DPT_Battery_Info:
        return "DPT_Battery_Info"
    case KnxDatapointType_DPT_Brightness_Colour_Temperature_Transition:
        return "DPT_Brightness_Colour_Temperature_Transition"
    case KnxDatapointType_DPT_Brightness_Colour_Temperature_Control:
        return "DPT_Brightness_Colour_Temperature_Control"
    case KnxDatapointType_DPT_Colour_RGBW:
        return "DPT_Colour_RGBW"
    case KnxDatapointType_DPT_Relative_Control_RGBW:
        return "DPT_Relative_Control_RGBW"
    case KnxDatapointType_DPT_Relative_Control_RGB:
        return "DPT_Relative_Control_RGB"
    case KnxDatapointType_DPT_GeographicalLocation:
        return "DPT_GeographicalLocation"
    case KnxDatapointType_DPT_Direction1_Control:
        return "DPT_Direction1_Control"
    case KnxDatapointType_DPT_TempRoomSetpSetF16_4:
        return "DPT_TempRoomSetpSetF16_4"
    case KnxDatapointType_DPT_TempRoomSetpSetShiftF16_4:
        return "DPT_TempRoomSetpSetShiftF16_4"
    case KnxDatapointType_DPT_Direction2_Control:
        return "DPT_Direction2_Control"
    case KnxDatapointType_DPT_Start_Control:
        return "DPT_Start_Control"
    case KnxDatapointType_DPT_State_Control:
        return "DPT_State_Control"
    case KnxDatapointType_DPT_Invert_Control:
        return "DPT_Invert_Control"
    case KnxDatapointType_DPT_Control_Dimming:
        return "DPT_Control_Dimming"
    case KnxDatapointType_DPT_Control_Blinds:
        return "DPT_Control_Blinds"
    case KnxDatapointType_DPT_Char_ASCII:
        return "DPT_Char_ASCII"
    case KnxDatapointType_DPT_Ramp:
        return "DPT_Ramp"
    case KnxDatapointType_DPT_Char_8859_1:
        return "DPT_Char_8859_1"
    case KnxDatapointType_DPT_Scaling:
        return "DPT_Scaling"
    case KnxDatapointType_DPT_Angle:
        return "DPT_Angle"
    case KnxDatapointType_DPT_Percent_U8:
        return "DPT_Percent_U8"
    case KnxDatapointType_DPT_DecimalFactor:
        return "DPT_DecimalFactor"
    case KnxDatapointType_DPT_Tariff:
        return "DPT_Tariff"
    case KnxDatapointType_DPT_Value_1_Ucount:
        return "DPT_Value_1_Ucount"
    case KnxDatapointType_DPT_FanStage:
        return "DPT_FanStage"
    case KnxDatapointType_DPT_Percent_V8:
        return "DPT_Percent_V8"
    case KnxDatapointType_DPT_Value_1_Count:
        return "DPT_Value_1_Count"
    case KnxDatapointType_DPT_Alarm:
        return "DPT_Alarm"
    case KnxDatapointType_DPT_Status_Mode3:
        return "DPT_Status_Mode3"
    case KnxDatapointType_DPT_Value_2_Ucount:
        return "DPT_Value_2_Ucount"
    case KnxDatapointType_DPT_TimePeriodMsec:
        return "DPT_TimePeriodMsec"
    case KnxDatapointType_DPT_TimePeriod10Msec:
        return "DPT_TimePeriod10Msec"
    case KnxDatapointType_DPT_TimePeriod100Msec:
        return "DPT_TimePeriod100Msec"
    case KnxDatapointType_DPT_TimePeriodSec:
        return "DPT_TimePeriodSec"
    case KnxDatapointType_DPT_TimePeriodMin:
        return "DPT_TimePeriodMin"
    case KnxDatapointType_DPT_TimePeriodHrs:
        return "DPT_TimePeriodHrs"
    case KnxDatapointType_DPT_PropDataType:
        return "DPT_PropDataType"
    case KnxDatapointType_DPT_Length_mm:
        return "DPT_Length_mm"
    case KnxDatapointType_DPT_BinaryValue:
        return "DPT_BinaryValue"
    case KnxDatapointType_DPT_UElCurrentmA:
        return "DPT_UElCurrentmA"
    case KnxDatapointType_DPT_Brightness:
        return "DPT_Brightness"
    case KnxDatapointType_DPT_Absolute_Colour_Temperature:
        return "DPT_Absolute_Colour_Temperature"
    case KnxDatapointType_DPT_Value_2_Count:
        return "DPT_Value_2_Count"
    case KnxDatapointType_DPT_DeltaTimeMsec:
        return "DPT_DeltaTimeMsec"
    case KnxDatapointType_DPT_DeltaTime10Msec:
        return "DPT_DeltaTime10Msec"
    case KnxDatapointType_DPT_DeltaTime100Msec:
        return "DPT_DeltaTime100Msec"
    case KnxDatapointType_DPT_DeltaTimeSec:
        return "DPT_DeltaTimeSec"
    case KnxDatapointType_DPT_DeltaTimeMin:
        return "DPT_DeltaTimeMin"
    case KnxDatapointType_DPT_DeltaTimeHrs:
        return "DPT_DeltaTimeHrs"
    case KnxDatapointType_DPT_Step:
        return "DPT_Step"
    case KnxDatapointType_DPT_Percent_V16:
        return "DPT_Percent_V16"
    case KnxDatapointType_DPT_Rotation_Angle:
        return "DPT_Rotation_Angle"
    case KnxDatapointType_DPT_Length_m:
        return "DPT_Length_m"
    case KnxDatapointType_DPT_Value_Temp:
        return "DPT_Value_Temp"
    case KnxDatapointType_DPT_Value_Tempd:
        return "DPT_Value_Tempd"
    case KnxDatapointType_DPT_Value_Tempa:
        return "DPT_Value_Tempa"
    case KnxDatapointType_DPT_Value_Lux:
        return "DPT_Value_Lux"
    case KnxDatapointType_DPT_Value_Wsp:
        return "DPT_Value_Wsp"
    case KnxDatapointType_DPT_Value_Pres:
        return "DPT_Value_Pres"
    case KnxDatapointType_DPT_Value_Humidity:
        return "DPT_Value_Humidity"
    case KnxDatapointType_DPT_UpDown:
        return "DPT_UpDown"
    case KnxDatapointType_DPT_Value_AirQuality:
        return "DPT_Value_AirQuality"
    case KnxDatapointType_DPT_Value_AirFlow:
        return "DPT_Value_AirFlow"
    case KnxDatapointType_DPT_Value_Time1:
        return "DPT_Value_Time1"
    case KnxDatapointType_DPT_Value_Time2:
        return "DPT_Value_Time2"
    case KnxDatapointType_DPT_Value_Volt:
        return "DPT_Value_Volt"
    case KnxDatapointType_DPT_Value_Curr:
        return "DPT_Value_Curr"
    case KnxDatapointType_DPT_PowerDensity:
        return "DPT_PowerDensity"
    case KnxDatapointType_DPT_KelvinPerPercent:
        return "DPT_KelvinPerPercent"
    case KnxDatapointType_DPT_Power:
        return "DPT_Power"
    case KnxDatapointType_DPT_Value_Volume_Flow:
        return "DPT_Value_Volume_Flow"
    case KnxDatapointType_DPT_OpenClose:
        return "DPT_OpenClose"
    case KnxDatapointType_DPT_Rain_Amount:
        return "DPT_Rain_Amount"
    case KnxDatapointType_DPT_Value_Temp_F:
        return "DPT_Value_Temp_F"
    case KnxDatapointType_DPT_Value_Wsp_kmh:
        return "DPT_Value_Wsp_kmh"
    case KnxDatapointType_DPT_Value_Absolute_Humidity:
        return "DPT_Value_Absolute_Humidity"
    case KnxDatapointType_DPT_Concentration_ygm3:
        return "DPT_Concentration_ygm3"
    case KnxDatapointType_DPT_TimeOfDay:
        return "DPT_TimeOfDay"
    case KnxDatapointType_DPT_Date:
        return "DPT_Date"
    case KnxDatapointType_DPT_Value_4_Ucount:
        return "DPT_Value_4_Ucount"
    case KnxDatapointType_DPT_LongTimePeriod_Sec:
        return "DPT_LongTimePeriod_Sec"
    case KnxDatapointType_DPT_LongTimePeriod_Min:
        return "DPT_LongTimePeriod_Min"
    }
    return ""
}
