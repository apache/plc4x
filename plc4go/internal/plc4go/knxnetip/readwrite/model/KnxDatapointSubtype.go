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
	DatapointType() KnxDatapointType
	Text() string
	Serialize(io utils.WriteBuffer) error
}

const (
	KnxDatapointSubtype_DPST_Switch                                   KnxDatapointSubtype = 10001
	KnxDatapointSubtype_DPST_Bool                                     KnxDatapointSubtype = 10002
	KnxDatapointSubtype_DPST_Enable                                   KnxDatapointSubtype = 10003
	KnxDatapointSubtype_DPST_Ramp                                     KnxDatapointSubtype = 10004
	KnxDatapointSubtype_DPST_Alarm                                    KnxDatapointSubtype = 10005
	KnxDatapointSubtype_DPST_BinaryValue                              KnxDatapointSubtype = 10006
	KnxDatapointSubtype_DPST_Step                                     KnxDatapointSubtype = 10007
	KnxDatapointSubtype_DPST_UpDown                                   KnxDatapointSubtype = 10008
	KnxDatapointSubtype_DPST_OpenClose                                KnxDatapointSubtype = 10009
	KnxDatapointSubtype_DPST_Start                                    KnxDatapointSubtype = 10010
	KnxDatapointSubtype_DPST_State                                    KnxDatapointSubtype = 10011
	KnxDatapointSubtype_DPST_Invert                                   KnxDatapointSubtype = 10012
	KnxDatapointSubtype_DPST_DimSendStyle                             KnxDatapointSubtype = 10013
	KnxDatapointSubtype_DPST_InputSource                              KnxDatapointSubtype = 10014
	KnxDatapointSubtype_DPST_Reset                                    KnxDatapointSubtype = 10015
	KnxDatapointSubtype_DPST_Ack                                      KnxDatapointSubtype = 10016
	KnxDatapointSubtype_DPST_Trigger                                  KnxDatapointSubtype = 10017
	KnxDatapointSubtype_DPST_Occupancy                                KnxDatapointSubtype = 10018
	KnxDatapointSubtype_DPST_Window_Door                              KnxDatapointSubtype = 10019
	KnxDatapointSubtype_DPST_LogicalFunction                          KnxDatapointSubtype = 10021
	KnxDatapointSubtype_DPST_Scene_AB                                 KnxDatapointSubtype = 10022
	KnxDatapointSubtype_DPST_ShutterBlinds_Mode                       KnxDatapointSubtype = 10023
	KnxDatapointSubtype_DPST_DayNight                                 KnxDatapointSubtype = 10024
	KnxDatapointSubtype_DPST_Heat_Cool                                KnxDatapointSubtype = 10100
	KnxDatapointSubtype_DPST_Switch_Control                           KnxDatapointSubtype = 20001
	KnxDatapointSubtype_DPST_Bool_Control                             KnxDatapointSubtype = 20002
	KnxDatapointSubtype_DPST_Enable_Control                           KnxDatapointSubtype = 20003
	KnxDatapointSubtype_DPST_Ramp_Control                             KnxDatapointSubtype = 20004
	KnxDatapointSubtype_DPST_Alarm_Control                            KnxDatapointSubtype = 20005
	KnxDatapointSubtype_DPST_BinaryValue_Control                      KnxDatapointSubtype = 20006
	KnxDatapointSubtype_DPST_Step_Control                             KnxDatapointSubtype = 20007
	KnxDatapointSubtype_DPST_Direction1_Control                       KnxDatapointSubtype = 20008
	KnxDatapointSubtype_DPST_Direction2_Control                       KnxDatapointSubtype = 20009
	KnxDatapointSubtype_DPST_Start_Control                            KnxDatapointSubtype = 20010
	KnxDatapointSubtype_DPST_State_Control                            KnxDatapointSubtype = 20011
	KnxDatapointSubtype_DPST_Invert_Control                           KnxDatapointSubtype = 20012
	KnxDatapointSubtype_DPST_Control_Dimming                          KnxDatapointSubtype = 30007
	KnxDatapointSubtype_DPST_Control_Blinds                           KnxDatapointSubtype = 30008
	KnxDatapointSubtype_DPST_Char_ASCII                               KnxDatapointSubtype = 40001
	KnxDatapointSubtype_DPST_Char_8859_1                              KnxDatapointSubtype = 40002
	KnxDatapointSubtype_DPST_Scaling                                  KnxDatapointSubtype = 50001
	KnxDatapointSubtype_DPST_Angle                                    KnxDatapointSubtype = 50003
	KnxDatapointSubtype_DPST_Percent_U8                               KnxDatapointSubtype = 50004
	KnxDatapointSubtype_DPST_DecimalFactor                            KnxDatapointSubtype = 50005
	KnxDatapointSubtype_DPST_Tariff                                   KnxDatapointSubtype = 50006
	KnxDatapointSubtype_DPST_Value_1_Ucount                           KnxDatapointSubtype = 50010
	KnxDatapointSubtype_DPST_FanStage                                 KnxDatapointSubtype = 50100
	KnxDatapointSubtype_DPST_Percent_V8                               KnxDatapointSubtype = 60001
	KnxDatapointSubtype_DPST_Value_1_Count                            KnxDatapointSubtype = 60010
	KnxDatapointSubtype_DPST_Status_Mode3                             KnxDatapointSubtype = 60020
	KnxDatapointSubtype_DPST_Value_2_Ucount                           KnxDatapointSubtype = 70001
	KnxDatapointSubtype_DPST_TimePeriodMsec                           KnxDatapointSubtype = 70002
	KnxDatapointSubtype_DPST_TimePeriod10Msec                         KnxDatapointSubtype = 70003
	KnxDatapointSubtype_DPST_TimePeriod100Msec                        KnxDatapointSubtype = 70004
	KnxDatapointSubtype_DPST_TimePeriodSec                            KnxDatapointSubtype = 70005
	KnxDatapointSubtype_DPST_TimePeriodMin                            KnxDatapointSubtype = 70006
	KnxDatapointSubtype_DPST_TimePeriodHrs                            KnxDatapointSubtype = 70007
	KnxDatapointSubtype_DPST_PropDataType                             KnxDatapointSubtype = 70010
	KnxDatapointSubtype_DPST_Length_mm                                KnxDatapointSubtype = 70011
	KnxDatapointSubtype_DPST_UElCurrentmA                             KnxDatapointSubtype = 70012
	KnxDatapointSubtype_DPST_Brightness                               KnxDatapointSubtype = 70013
	KnxDatapointSubtype_DPST_Absolute_Colour_Temperature              KnxDatapointSubtype = 70600
	KnxDatapointSubtype_DPST_Value_2_Count                            KnxDatapointSubtype = 80001
	KnxDatapointSubtype_DPST_DeltaTimeMsec                            KnxDatapointSubtype = 80002
	KnxDatapointSubtype_DPST_DeltaTime10Msec                          KnxDatapointSubtype = 80003
	KnxDatapointSubtype_DPST_DeltaTime100Msec                         KnxDatapointSubtype = 80004
	KnxDatapointSubtype_DPST_DeltaTimeSec                             KnxDatapointSubtype = 80005
	KnxDatapointSubtype_DPST_DeltaTimeMin                             KnxDatapointSubtype = 80006
	KnxDatapointSubtype_DPST_DeltaTimeHrs                             KnxDatapointSubtype = 80007
	KnxDatapointSubtype_DPST_Percent_V16                              KnxDatapointSubtype = 80010
	KnxDatapointSubtype_DPST_Rotation_Angle                           KnxDatapointSubtype = 80011
	KnxDatapointSubtype_DPST_Length_m                                 KnxDatapointSubtype = 80012
	KnxDatapointSubtype_DPST_Value_Temp                               KnxDatapointSubtype = 90001
	KnxDatapointSubtype_DPST_Value_Tempd                              KnxDatapointSubtype = 90002
	KnxDatapointSubtype_DPST_Value_Tempa                              KnxDatapointSubtype = 90003
	KnxDatapointSubtype_DPST_Value_Lux                                KnxDatapointSubtype = 90004
	KnxDatapointSubtype_DPST_Value_Wsp                                KnxDatapointSubtype = 90005
	KnxDatapointSubtype_DPST_Value_Pres                               KnxDatapointSubtype = 90006
	KnxDatapointSubtype_DPST_Value_Humidity                           KnxDatapointSubtype = 90007
	KnxDatapointSubtype_DPST_Value_AirQuality                         KnxDatapointSubtype = 90008
	KnxDatapointSubtype_DPST_Value_AirFlow                            KnxDatapointSubtype = 90009
	KnxDatapointSubtype_DPST_Value_Time1                              KnxDatapointSubtype = 90010
	KnxDatapointSubtype_DPST_Value_Time2                              KnxDatapointSubtype = 90011
	KnxDatapointSubtype_DPST_Value_Volt                               KnxDatapointSubtype = 90020
	KnxDatapointSubtype_DPST_Value_Curr                               KnxDatapointSubtype = 90021
	KnxDatapointSubtype_DPST_PowerDensity                             KnxDatapointSubtype = 90022
	KnxDatapointSubtype_DPST_KelvinPerPercent                         KnxDatapointSubtype = 90023
	KnxDatapointSubtype_DPST_Power                                    KnxDatapointSubtype = 90024
	KnxDatapointSubtype_DPST_Value_Volume_Flow                        KnxDatapointSubtype = 90025
	KnxDatapointSubtype_DPST_Rain_Amount                              KnxDatapointSubtype = 90026
	KnxDatapointSubtype_DPST_Value_Temp_F                             KnxDatapointSubtype = 90027
	KnxDatapointSubtype_DPST_Value_Wsp_kmh                            KnxDatapointSubtype = 90028
	KnxDatapointSubtype_DPST_Value_Absolute_Humidity                  KnxDatapointSubtype = 90029
	KnxDatapointSubtype_DPST_Concentration_ygm3                       KnxDatapointSubtype = 90030
	KnxDatapointSubtype_DPST_TimeOfDay                                KnxDatapointSubtype = 100001
	KnxDatapointSubtype_DPST_Date                                     KnxDatapointSubtype = 110001
	KnxDatapointSubtype_DPST_Value_4_Ucount                           KnxDatapointSubtype = 120001
	KnxDatapointSubtype_DPST_LongTimePeriod_Sec                       KnxDatapointSubtype = 120100
	KnxDatapointSubtype_DPST_LongTimePeriod_Min                       KnxDatapointSubtype = 120101
	KnxDatapointSubtype_DPST_LongTimePeriod_Hrs                       KnxDatapointSubtype = 120102
	KnxDatapointSubtype_DPST_VolumeLiquid_Litre                       KnxDatapointSubtype = 121200
	KnxDatapointSubtype_DPST_Volume_m_3                               KnxDatapointSubtype = 121201
	KnxDatapointSubtype_DPST_Value_4_Count                            KnxDatapointSubtype = 130001
	KnxDatapointSubtype_DPST_FlowRate_m3h                             KnxDatapointSubtype = 130002
	KnxDatapointSubtype_DPST_ActiveEnergy                             KnxDatapointSubtype = 130010
	KnxDatapointSubtype_DPST_ApparantEnergy                           KnxDatapointSubtype = 130011
	KnxDatapointSubtype_DPST_ReactiveEnergy                           KnxDatapointSubtype = 130012
	KnxDatapointSubtype_DPST_ActiveEnergy_kWh                         KnxDatapointSubtype = 130013
	KnxDatapointSubtype_DPST_ApparantEnergy_kVAh                      KnxDatapointSubtype = 130014
	KnxDatapointSubtype_DPST_ReactiveEnergy_kVARh                     KnxDatapointSubtype = 130015
	KnxDatapointSubtype_DPST_ActiveEnergy_MWh                         KnxDatapointSubtype = 130016
	KnxDatapointSubtype_DPST_LongDeltaTimeSec                         KnxDatapointSubtype = 130100
	KnxDatapointSubtype_DPST_DeltaVolumeLiquid_Litre                  KnxDatapointSubtype = 131200
	KnxDatapointSubtype_DPST_DeltaVolume_m_3                          KnxDatapointSubtype = 131201
	KnxDatapointSubtype_DPST_Value_Acceleration                       KnxDatapointSubtype = 140000
	KnxDatapointSubtype_DPST_Value_Acceleration_Angular               KnxDatapointSubtype = 140001
	KnxDatapointSubtype_DPST_Value_Activation_Energy                  KnxDatapointSubtype = 140002
	KnxDatapointSubtype_DPST_Value_Activity                           KnxDatapointSubtype = 140003
	KnxDatapointSubtype_DPST_Value_Mol                                KnxDatapointSubtype = 140004
	KnxDatapointSubtype_DPST_Value_Amplitude                          KnxDatapointSubtype = 140005
	KnxDatapointSubtype_DPST_Value_AngleRad                           KnxDatapointSubtype = 140006
	KnxDatapointSubtype_DPST_Value_AngleDeg                           KnxDatapointSubtype = 140007
	KnxDatapointSubtype_DPST_Value_Angular_Momentum                   KnxDatapointSubtype = 140008
	KnxDatapointSubtype_DPST_Value_Angular_Velocity                   KnxDatapointSubtype = 140009
	KnxDatapointSubtype_DPST_Value_Area                               KnxDatapointSubtype = 140010
	KnxDatapointSubtype_DPST_Value_Capacitance                        KnxDatapointSubtype = 140011
	KnxDatapointSubtype_DPST_Value_Charge_DensitySurface              KnxDatapointSubtype = 140012
	KnxDatapointSubtype_DPST_Value_Charge_DensityVolume               KnxDatapointSubtype = 140013
	KnxDatapointSubtype_DPST_Value_Compressibility                    KnxDatapointSubtype = 140014
	KnxDatapointSubtype_DPST_Value_Conductance                        KnxDatapointSubtype = 140015
	KnxDatapointSubtype_DPST_Value_Electrical_Conductivity            KnxDatapointSubtype = 140016
	KnxDatapointSubtype_DPST_Value_Density                            KnxDatapointSubtype = 140017
	KnxDatapointSubtype_DPST_Value_Electric_Charge                    KnxDatapointSubtype = 140018
	KnxDatapointSubtype_DPST_Value_Electric_Current                   KnxDatapointSubtype = 140019
	KnxDatapointSubtype_DPST_Value_Electric_CurrentDensity            KnxDatapointSubtype = 140020
	KnxDatapointSubtype_DPST_Value_Electric_DipoleMoment              KnxDatapointSubtype = 140021
	KnxDatapointSubtype_DPST_Value_Electric_Displacement              KnxDatapointSubtype = 140022
	KnxDatapointSubtype_DPST_Value_Electric_FieldStrength             KnxDatapointSubtype = 140023
	KnxDatapointSubtype_DPST_Value_Electric_Flux                      KnxDatapointSubtype = 140024
	KnxDatapointSubtype_DPST_Value_Electric_FluxDensity               KnxDatapointSubtype = 140025
	KnxDatapointSubtype_DPST_Value_Electric_Polarization              KnxDatapointSubtype = 140026
	KnxDatapointSubtype_DPST_Value_Electric_Potential                 KnxDatapointSubtype = 140027
	KnxDatapointSubtype_DPST_Value_Electric_PotentialDifference       KnxDatapointSubtype = 140028
	KnxDatapointSubtype_DPST_Value_ElectromagneticMoment              KnxDatapointSubtype = 140029
	KnxDatapointSubtype_DPST_Value_Electromotive_Force                KnxDatapointSubtype = 140030
	KnxDatapointSubtype_DPST_Value_Energy                             KnxDatapointSubtype = 140031
	KnxDatapointSubtype_DPST_Value_Force                              KnxDatapointSubtype = 140032
	KnxDatapointSubtype_DPST_Value_Frequency                          KnxDatapointSubtype = 140033
	KnxDatapointSubtype_DPST_Value_Angular_Frequency                  KnxDatapointSubtype = 140034
	KnxDatapointSubtype_DPST_Value_Heat_Capacity                      KnxDatapointSubtype = 140035
	KnxDatapointSubtype_DPST_Value_Heat_FlowRate                      KnxDatapointSubtype = 140036
	KnxDatapointSubtype_DPST_Value_Heat_Quantity                      KnxDatapointSubtype = 140037
	KnxDatapointSubtype_DPST_Value_Impedance                          KnxDatapointSubtype = 140038
	KnxDatapointSubtype_DPST_Value_Length                             KnxDatapointSubtype = 140039
	KnxDatapointSubtype_DPST_Value_Light_Quantity                     KnxDatapointSubtype = 140040
	KnxDatapointSubtype_DPST_Value_Luminance                          KnxDatapointSubtype = 140041
	KnxDatapointSubtype_DPST_Value_Luminous_Flux                      KnxDatapointSubtype = 140042
	KnxDatapointSubtype_DPST_Value_Luminous_Intensity                 KnxDatapointSubtype = 140043
	KnxDatapointSubtype_DPST_Value_Magnetic_FieldStrength             KnxDatapointSubtype = 140044
	KnxDatapointSubtype_DPST_Value_Magnetic_Flux                      KnxDatapointSubtype = 140045
	KnxDatapointSubtype_DPST_Value_Magnetic_FluxDensity               KnxDatapointSubtype = 140046
	KnxDatapointSubtype_DPST_Value_Magnetic_Moment                    KnxDatapointSubtype = 140047
	KnxDatapointSubtype_DPST_Value_Magnetic_Polarization              KnxDatapointSubtype = 140048
	KnxDatapointSubtype_DPST_Value_Magnetization                      KnxDatapointSubtype = 140049
	KnxDatapointSubtype_DPST_Value_MagnetomotiveForce                 KnxDatapointSubtype = 140050
	KnxDatapointSubtype_DPST_Value_Mass                               KnxDatapointSubtype = 140051
	KnxDatapointSubtype_DPST_Value_MassFlux                           KnxDatapointSubtype = 140052
	KnxDatapointSubtype_DPST_Value_Momentum                           KnxDatapointSubtype = 140053
	KnxDatapointSubtype_DPST_Value_Phase_AngleRad                     KnxDatapointSubtype = 140054
	KnxDatapointSubtype_DPST_Value_Phase_AngleDeg                     KnxDatapointSubtype = 140055
	KnxDatapointSubtype_DPST_Value_Power                              KnxDatapointSubtype = 140056
	KnxDatapointSubtype_DPST_Value_Power_Factor                       KnxDatapointSubtype = 140057
	KnxDatapointSubtype_DPST_Value_Pressure                           KnxDatapointSubtype = 140058
	KnxDatapointSubtype_DPST_Value_Reactance                          KnxDatapointSubtype = 140059
	KnxDatapointSubtype_DPST_Value_Resistance                         KnxDatapointSubtype = 140060
	KnxDatapointSubtype_DPST_Value_Resistivity                        KnxDatapointSubtype = 140061
	KnxDatapointSubtype_DPST_Value_SelfInductance                     KnxDatapointSubtype = 140062
	KnxDatapointSubtype_DPST_Value_SolidAngle                         KnxDatapointSubtype = 140063
	KnxDatapointSubtype_DPST_Value_Sound_Intensity                    KnxDatapointSubtype = 140064
	KnxDatapointSubtype_DPST_Value_Speed                              KnxDatapointSubtype = 140065
	KnxDatapointSubtype_DPST_Value_Stress                             KnxDatapointSubtype = 140066
	KnxDatapointSubtype_DPST_Value_Surface_Tension                    KnxDatapointSubtype = 140067
	KnxDatapointSubtype_DPST_Value_Common_Temperature                 KnxDatapointSubtype = 140068
	KnxDatapointSubtype_DPST_Value_Absolute_Temperature               KnxDatapointSubtype = 140069
	KnxDatapointSubtype_DPST_Value_TemperatureDifference              KnxDatapointSubtype = 140070
	KnxDatapointSubtype_DPST_Value_Thermal_Capacity                   KnxDatapointSubtype = 140071
	KnxDatapointSubtype_DPST_Value_Thermal_Conductivity               KnxDatapointSubtype = 140072
	KnxDatapointSubtype_DPST_Value_ThermoelectricPower                KnxDatapointSubtype = 140073
	KnxDatapointSubtype_DPST_Value_Time                               KnxDatapointSubtype = 140074
	KnxDatapointSubtype_DPST_Value_Torque                             KnxDatapointSubtype = 140075
	KnxDatapointSubtype_DPST_Value_Volume                             KnxDatapointSubtype = 140076
	KnxDatapointSubtype_DPST_Value_Volume_Flux                        KnxDatapointSubtype = 140077
	KnxDatapointSubtype_DPST_Value_Weight                             KnxDatapointSubtype = 140078
	KnxDatapointSubtype_DPST_Value_Work                               KnxDatapointSubtype = 140079
	KnxDatapointSubtype_DPST_Volume_Flux_Meter                        KnxDatapointSubtype = 141200
	KnxDatapointSubtype_DPST_Volume_Flux_ls                           KnxDatapointSubtype = 141201
	KnxDatapointSubtype_DPST_Access_Data                              KnxDatapointSubtype = 150000
	KnxDatapointSubtype_DPST_String_ASCII                             KnxDatapointSubtype = 160000
	KnxDatapointSubtype_DPST_String_8859_1                            KnxDatapointSubtype = 160001
	KnxDatapointSubtype_DPST_SceneNumber                              KnxDatapointSubtype = 170001
	KnxDatapointSubtype_DPST_SceneControl                             KnxDatapointSubtype = 180001
	KnxDatapointSubtype_DPST_DateTime                                 KnxDatapointSubtype = 190001
	KnxDatapointSubtype_DPST_SCLOMode                                 KnxDatapointSubtype = 200001
	KnxDatapointSubtype_DPST_BuildingMode                             KnxDatapointSubtype = 200002
	KnxDatapointSubtype_DPST_OccMode                                  KnxDatapointSubtype = 200003
	KnxDatapointSubtype_DPST_Priority                                 KnxDatapointSubtype = 200004
	KnxDatapointSubtype_DPST_LightApplicationMode                     KnxDatapointSubtype = 200005
	KnxDatapointSubtype_DPST_ApplicationArea                          KnxDatapointSubtype = 200006
	KnxDatapointSubtype_DPST_AlarmClassType                           KnxDatapointSubtype = 200007
	KnxDatapointSubtype_DPST_PSUMode                                  KnxDatapointSubtype = 200008
	KnxDatapointSubtype_DPST_ErrorClass_System                        KnxDatapointSubtype = 200011
	KnxDatapointSubtype_DPST_ErrorClass_HVAC                          KnxDatapointSubtype = 200012
	KnxDatapointSubtype_DPST_Time_Delay                               KnxDatapointSubtype = 200013
	KnxDatapointSubtype_DPST_Beaufort_Wind_Force_Scale                KnxDatapointSubtype = 200014
	KnxDatapointSubtype_DPST_SensorSelect                             KnxDatapointSubtype = 200017
	KnxDatapointSubtype_DPST_ActuatorConnectType                      KnxDatapointSubtype = 200020
	KnxDatapointSubtype_DPST_Cloud_Cover                              KnxDatapointSubtype = 200021
	KnxDatapointSubtype_DPST_PowerReturnMode                          KnxDatapointSubtype = 200022
	KnxDatapointSubtype_DPST_FuelType                                 KnxDatapointSubtype = 200100
	KnxDatapointSubtype_DPST_BurnerType                               KnxDatapointSubtype = 200101
	KnxDatapointSubtype_DPST_HVACMode                                 KnxDatapointSubtype = 200102
	KnxDatapointSubtype_DPST_DHWMode                                  KnxDatapointSubtype = 200103
	KnxDatapointSubtype_DPST_LoadPriority                             KnxDatapointSubtype = 200104
	KnxDatapointSubtype_DPST_HVACContrMode                            KnxDatapointSubtype = 200105
	KnxDatapointSubtype_DPST_HVACEmergMode                            KnxDatapointSubtype = 200106
	KnxDatapointSubtype_DPST_ChangeoverMode                           KnxDatapointSubtype = 200107
	KnxDatapointSubtype_DPST_ValveMode                                KnxDatapointSubtype = 200108
	KnxDatapointSubtype_DPST_DamperMode                               KnxDatapointSubtype = 200109
	KnxDatapointSubtype_DPST_HeaterMode                               KnxDatapointSubtype = 200110
	KnxDatapointSubtype_DPST_FanMode                                  KnxDatapointSubtype = 200111
	KnxDatapointSubtype_DPST_MasterSlaveMode                          KnxDatapointSubtype = 200112
	KnxDatapointSubtype_DPST_StatusRoomSetp                           KnxDatapointSubtype = 200113
	KnxDatapointSubtype_DPST_Metering_DeviceType                      KnxDatapointSubtype = 200114
	KnxDatapointSubtype_DPST_HumDehumMode                             KnxDatapointSubtype = 200115
	KnxDatapointSubtype_DPST_EnableHCStage                            KnxDatapointSubtype = 200116
	KnxDatapointSubtype_DPST_ADAType                                  KnxDatapointSubtype = 200120
	KnxDatapointSubtype_DPST_BackupMode                               KnxDatapointSubtype = 200121
	KnxDatapointSubtype_DPST_StartSynchronization                     KnxDatapointSubtype = 200122
	KnxDatapointSubtype_DPST_Behaviour_Lock_Unlock                    KnxDatapointSubtype = 200600
	KnxDatapointSubtype_DPST_Behaviour_Bus_Power_Up_Down              KnxDatapointSubtype = 200601
	KnxDatapointSubtype_DPST_DALI_Fade_Time                           KnxDatapointSubtype = 200602
	KnxDatapointSubtype_DPST_BlinkingMode                             KnxDatapointSubtype = 200603
	KnxDatapointSubtype_DPST_LightControlMode                         KnxDatapointSubtype = 200604
	KnxDatapointSubtype_DPST_SwitchPBModel                            KnxDatapointSubtype = 200605
	KnxDatapointSubtype_DPST_PBAction                                 KnxDatapointSubtype = 200606
	KnxDatapointSubtype_DPST_DimmPBModel                              KnxDatapointSubtype = 200607
	KnxDatapointSubtype_DPST_SwitchOnMode                             KnxDatapointSubtype = 200608
	KnxDatapointSubtype_DPST_LoadTypeSet                              KnxDatapointSubtype = 200609
	KnxDatapointSubtype_DPST_LoadTypeDetected                         KnxDatapointSubtype = 200610
	KnxDatapointSubtype_DPST_Converter_Test_Control                   KnxDatapointSubtype = 200611
	KnxDatapointSubtype_DPST_SABExcept_Behaviour                      KnxDatapointSubtype = 200801
	KnxDatapointSubtype_DPST_SABBehaviour_Lock_Unlock                 KnxDatapointSubtype = 200802
	KnxDatapointSubtype_DPST_SSSBMode                                 KnxDatapointSubtype = 200803
	KnxDatapointSubtype_DPST_BlindsControlMode                        KnxDatapointSubtype = 200804
	KnxDatapointSubtype_DPST_CommMode                                 KnxDatapointSubtype = 201000
	KnxDatapointSubtype_DPST_AddInfoTypes                             KnxDatapointSubtype = 201001
	KnxDatapointSubtype_DPST_RF_ModeSelect                            KnxDatapointSubtype = 201002
	KnxDatapointSubtype_DPST_RF_FilterSelect                          KnxDatapointSubtype = 201003
	KnxDatapointSubtype_DPST_StatusGen                                KnxDatapointSubtype = 210001
	KnxDatapointSubtype_DPST_Device_Control                           KnxDatapointSubtype = 210002
	KnxDatapointSubtype_DPST_ForceSign                                KnxDatapointSubtype = 210100
	KnxDatapointSubtype_DPST_ForceSignCool                            KnxDatapointSubtype = 210101
	KnxDatapointSubtype_DPST_StatusRHC                                KnxDatapointSubtype = 210102
	KnxDatapointSubtype_DPST_StatusSDHWC                              KnxDatapointSubtype = 210103
	KnxDatapointSubtype_DPST_FuelTypeSet                              KnxDatapointSubtype = 210104
	KnxDatapointSubtype_DPST_StatusRCC                                KnxDatapointSubtype = 210105
	KnxDatapointSubtype_DPST_StatusAHU                                KnxDatapointSubtype = 210106
	KnxDatapointSubtype_DPST_CombinedStatus_RTSM                      KnxDatapointSubtype = 210107
	KnxDatapointSubtype_DPST_LightActuatorErrorInfo                   KnxDatapointSubtype = 210601
	KnxDatapointSubtype_DPST_RF_ModeInfo                              KnxDatapointSubtype = 211000
	KnxDatapointSubtype_DPST_RF_FilterInfo                            KnxDatapointSubtype = 211001
	KnxDatapointSubtype_DPST_Channel_Activation_8                     KnxDatapointSubtype = 211010
	KnxDatapointSubtype_DPST_StatusDHWC                               KnxDatapointSubtype = 220100
	KnxDatapointSubtype_DPST_StatusRHCC                               KnxDatapointSubtype = 220101
	KnxDatapointSubtype_DPST_CombinedStatus_HVA                       KnxDatapointSubtype = 220102
	KnxDatapointSubtype_DPST_CombinedStatus_RTC                       KnxDatapointSubtype = 220103
	KnxDatapointSubtype_DPST_Media                                    KnxDatapointSubtype = 221000
	KnxDatapointSubtype_DPST_Channel_Activation_16                    KnxDatapointSubtype = 221010
	KnxDatapointSubtype_DPST_OnOffAction                              KnxDatapointSubtype = 230001
	KnxDatapointSubtype_DPST_Alarm_Reaction                           KnxDatapointSubtype = 230002
	KnxDatapointSubtype_DPST_UpDown_Action                            KnxDatapointSubtype = 230003
	KnxDatapointSubtype_DPST_HVAC_PB_Action                           KnxDatapointSubtype = 230102
	KnxDatapointSubtype_DPST_DoubleNibble                             KnxDatapointSubtype = 251000
	KnxDatapointSubtype_DPST_SceneInfo                                KnxDatapointSubtype = 260001
	KnxDatapointSubtype_DPST_CombinedInfoOnOff                        KnxDatapointSubtype = 270001
	KnxDatapointSubtype_DPST_ActiveEnergy_V64                         KnxDatapointSubtype = 290010
	KnxDatapointSubtype_DPST_ApparantEnergy_V64                       KnxDatapointSubtype = 290011
	KnxDatapointSubtype_DPST_ReactiveEnergy_V64                       KnxDatapointSubtype = 290012
	KnxDatapointSubtype_DPST_Channel_Activation_24                    KnxDatapointSubtype = 301010
	KnxDatapointSubtype_DPST_HVACModeNext                             KnxDatapointSubtype = 2060100
	KnxDatapointSubtype_DPST_DHWModeNext                              KnxDatapointSubtype = 2060102
	KnxDatapointSubtype_DPST_OccModeNext                              KnxDatapointSubtype = 2060104
	KnxDatapointSubtype_DPST_BuildingModeNext                         KnxDatapointSubtype = 2060105
	KnxDatapointSubtype_DPST_Version                                  KnxDatapointSubtype = 2170001
	KnxDatapointSubtype_DPST_AlarmInfo                                KnxDatapointSubtype = 2190001
	KnxDatapointSubtype_DPST_TempRoomSetpSetF16_3                     KnxDatapointSubtype = 2220100
	KnxDatapointSubtype_DPST_TempRoomSetpSetShiftF16_3                KnxDatapointSubtype = 2220101
	KnxDatapointSubtype_DPST_Scaling_Speed                            KnxDatapointSubtype = 2250001
	KnxDatapointSubtype_DPST_Scaling_Step_Time                        KnxDatapointSubtype = 2250002
	KnxDatapointSubtype_DPST_MeteringValue                            KnxDatapointSubtype = 2290001
	KnxDatapointSubtype_DPST_MBus_Address                             KnxDatapointSubtype = 2301000
	KnxDatapointSubtype_DPST_Colour_RGB                               KnxDatapointSubtype = 2320600
	KnxDatapointSubtype_DPST_LanguageCodeAlpha2_ASCII                 KnxDatapointSubtype = 2340001
	KnxDatapointSubtype_DPST_Tariff_ActiveEnergy                      KnxDatapointSubtype = 2350001
	KnxDatapointSubtype_DPST_Prioritised_Mode_Control                 KnxDatapointSubtype = 2360001
	KnxDatapointSubtype_DPST_DALI_Control_Gear_Diagnostic             KnxDatapointSubtype = 2370600
	KnxDatapointSubtype_DPST_DALI_Diagnostics                         KnxDatapointSubtype = 2380600
	KnxDatapointSubtype_DPST_CombinedPosition                         KnxDatapointSubtype = 2400800
	KnxDatapointSubtype_DPST_StatusSAB                                KnxDatapointSubtype = 2410800
	KnxDatapointSubtype_DPST_Colour_xyY                               KnxDatapointSubtype = 2420600
	KnxDatapointSubtype_DPST_Converter_Status                         KnxDatapointSubtype = 2440600
	KnxDatapointSubtype_DPST_Converter_Test_Result                    KnxDatapointSubtype = 2450600
	KnxDatapointSubtype_DPST_Battery_Info                             KnxDatapointSubtype = 2460600
	KnxDatapointSubtype_DPST_Brightness_Colour_Temperature_Transition KnxDatapointSubtype = 2490600
	KnxDatapointSubtype_DPST_Brightness_Colour_Temperature_Control    KnxDatapointSubtype = 2500600
	KnxDatapointSubtype_DPST_Colour_RGBW                              KnxDatapointSubtype = 2510600
	KnxDatapointSubtype_DPST_Relative_Control_RGBW                    KnxDatapointSubtype = 2520600
	KnxDatapointSubtype_DPST_Relative_Control_RGB                     KnxDatapointSubtype = 2540600
	KnxDatapointSubtype_DPST_GeographicalLocation                     KnxDatapointSubtype = 2550001
	KnxDatapointSubtype_DPST_TempRoomSetpSetF16_4                     KnxDatapointSubtype = 2750100
	KnxDatapointSubtype_DPST_TempRoomSetpSetShiftF16_4                KnxDatapointSubtype = 2750101
)

func (e KnxDatapointSubtype) DatapointType() KnxDatapointType {
	switch e {
	case 100001:
		{ /* '100001' */
			return KnxDatapointType_DPT_TIME
		}
	case 10001:
		{ /* '10001' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10002:
		{ /* '10002' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10003:
		{ /* '10003' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10004:
		{ /* '10004' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10005:
		{ /* '10005' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10006:
		{ /* '10006' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10007:
		{ /* '10007' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10008:
		{ /* '10008' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10009:
		{ /* '10009' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10010:
		{ /* '10010' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10011:
		{ /* '10011' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10012:
		{ /* '10012' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10013:
		{ /* '10013' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10014:
		{ /* '10014' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10015:
		{ /* '10015' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10016:
		{ /* '10016' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10017:
		{ /* '10017' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10018:
		{ /* '10018' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10019:
		{ /* '10019' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10021:
		{ /* '10021' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10022:
		{ /* '10022' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10023:
		{ /* '10023' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10024:
		{ /* '10024' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 10100:
		{ /* '10100' */
			return KnxDatapointType_DPT_1_BIT
		}
	case 110001:
		{ /* '110001' */
			return KnxDatapointType_DPT_DATE
		}
	case 120001:
		{ /* '120001' */
			return KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE
		}
	case 120100:
		{ /* '120100' */
			return KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE
		}
	case 120101:
		{ /* '120101' */
			return KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE
		}
	case 120102:
		{ /* '120102' */
			return KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE
		}
	case 121200:
		{ /* '121200' */
			return KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE
		}
	case 121201:
		{ /* '121201' */
			return KnxDatapointType_DPT_4_BYTE_UNSIGNED_VALUE
		}
	case 130001:
		{ /* '130001' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 130002:
		{ /* '130002' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 130010:
		{ /* '130010' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 130011:
		{ /* '130011' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 130012:
		{ /* '130012' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 130013:
		{ /* '130013' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 130014:
		{ /* '130014' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 130015:
		{ /* '130015' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 130016:
		{ /* '130016' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 130100:
		{ /* '130100' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 131200:
		{ /* '131200' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 131201:
		{ /* '131201' */
			return KnxDatapointType_DPT_4_BYTE_SIGNED_VALUE
		}
	case 140000:
		{ /* '140000' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140001:
		{ /* '140001' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140002:
		{ /* '140002' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140003:
		{ /* '140003' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140004:
		{ /* '140004' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140005:
		{ /* '140005' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140006:
		{ /* '140006' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140007:
		{ /* '140007' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140008:
		{ /* '140008' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140009:
		{ /* '140009' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140010:
		{ /* '140010' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140011:
		{ /* '140011' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140012:
		{ /* '140012' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140013:
		{ /* '140013' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140014:
		{ /* '140014' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140015:
		{ /* '140015' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140016:
		{ /* '140016' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140017:
		{ /* '140017' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140018:
		{ /* '140018' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140019:
		{ /* '140019' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140020:
		{ /* '140020' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140021:
		{ /* '140021' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140022:
		{ /* '140022' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140023:
		{ /* '140023' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140024:
		{ /* '140024' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140025:
		{ /* '140025' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140026:
		{ /* '140026' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140027:
		{ /* '140027' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140028:
		{ /* '140028' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140029:
		{ /* '140029' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140030:
		{ /* '140030' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140031:
		{ /* '140031' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140032:
		{ /* '140032' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140033:
		{ /* '140033' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140034:
		{ /* '140034' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140035:
		{ /* '140035' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140036:
		{ /* '140036' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140037:
		{ /* '140037' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140038:
		{ /* '140038' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140039:
		{ /* '140039' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140040:
		{ /* '140040' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140041:
		{ /* '140041' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140042:
		{ /* '140042' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140043:
		{ /* '140043' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140044:
		{ /* '140044' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140045:
		{ /* '140045' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140046:
		{ /* '140046' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140047:
		{ /* '140047' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140048:
		{ /* '140048' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140049:
		{ /* '140049' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140050:
		{ /* '140050' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140051:
		{ /* '140051' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140052:
		{ /* '140052' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140053:
		{ /* '140053' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140054:
		{ /* '140054' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140055:
		{ /* '140055' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140056:
		{ /* '140056' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140057:
		{ /* '140057' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140058:
		{ /* '140058' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140059:
		{ /* '140059' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140060:
		{ /* '140060' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140061:
		{ /* '140061' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140062:
		{ /* '140062' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140063:
		{ /* '140063' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140064:
		{ /* '140064' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140065:
		{ /* '140065' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140066:
		{ /* '140066' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140067:
		{ /* '140067' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140068:
		{ /* '140068' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140069:
		{ /* '140069' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140070:
		{ /* '140070' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140071:
		{ /* '140071' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140072:
		{ /* '140072' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140073:
		{ /* '140073' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140074:
		{ /* '140074' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140075:
		{ /* '140075' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140076:
		{ /* '140076' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140077:
		{ /* '140077' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140078:
		{ /* '140078' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 140079:
		{ /* '140079' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 141200:
		{ /* '141200' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 141201:
		{ /* '141201' */
			return KnxDatapointType_DPT_4_BYTE_FLOAT_VALUE
		}
	case 150000:
		{ /* '150000' */
			return KnxDatapointType_DPT_ENTRANCE_ACCESS
		}
	case 160000:
		{ /* '160000' */
			return KnxDatapointType_DPT_CHARACTER_STRING
		}
	case 160001:
		{ /* '160001' */
			return KnxDatapointType_DPT_CHARACTER_STRING
		}
	case 170001:
		{ /* '170001' */
			return KnxDatapointType_DPT_SCENE_NUMBER
		}
	case 180001:
		{ /* '180001' */
			return KnxDatapointType_DPT_SCENE_CONTROL
		}
	case 190001:
		{ /* '190001' */
			return KnxDatapointType_DPT_DATE_TIME
		}
	case 200001:
		{ /* '200001' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200002:
		{ /* '200002' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200003:
		{ /* '200003' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200004:
		{ /* '200004' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200005:
		{ /* '200005' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200006:
		{ /* '200006' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200007:
		{ /* '200007' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200008:
		{ /* '200008' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 20001:
		{ /* '20001' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 200011:
		{ /* '200011' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200012:
		{ /* '200012' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200013:
		{ /* '200013' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200014:
		{ /* '200014' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200017:
		{ /* '200017' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 20002:
		{ /* '20002' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 200020:
		{ /* '200020' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200021:
		{ /* '200021' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200022:
		{ /* '200022' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 20003:
		{ /* '20003' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 20004:
		{ /* '20004' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 20005:
		{ /* '20005' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 20006:
		{ /* '20006' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 20007:
		{ /* '20007' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 20008:
		{ /* '20008' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 20009:
		{ /* '20009' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 20010:
		{ /* '20010' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 200100:
		{ /* '200100' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200101:
		{ /* '200101' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200102:
		{ /* '200102' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200103:
		{ /* '200103' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200104:
		{ /* '200104' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200105:
		{ /* '200105' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200106:
		{ /* '200106' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200107:
		{ /* '200107' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200108:
		{ /* '200108' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200109:
		{ /* '200109' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 20011:
		{ /* '20011' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 200110:
		{ /* '200110' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200111:
		{ /* '200111' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200112:
		{ /* '200112' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200113:
		{ /* '200113' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200114:
		{ /* '200114' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200115:
		{ /* '200115' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200116:
		{ /* '200116' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 20012:
		{ /* '20012' */
			return KnxDatapointType_DPT_1_BIT_CONTROLLED
		}
	case 200120:
		{ /* '200120' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200121:
		{ /* '200121' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200122:
		{ /* '200122' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200600:
		{ /* '200600' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200601:
		{ /* '200601' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200602:
		{ /* '200602' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200603:
		{ /* '200603' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200604:
		{ /* '200604' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200605:
		{ /* '200605' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200606:
		{ /* '200606' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200607:
		{ /* '200607' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200608:
		{ /* '200608' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200609:
		{ /* '200609' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200610:
		{ /* '200610' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200611:
		{ /* '200611' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200801:
		{ /* '200801' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200802:
		{ /* '200802' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200803:
		{ /* '200803' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 200804:
		{ /* '200804' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 201000:
		{ /* '201000' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 201001:
		{ /* '201001' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 201002:
		{ /* '201002' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 201003:
		{ /* '201003' */
			return KnxDatapointType_DPT_1_BYTE
		}
	case 2060100:
		{ /* '2060100' */
			return KnxDatapointType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM
		}
	case 2060102:
		{ /* '2060102' */
			return KnxDatapointType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM
		}
	case 2060104:
		{ /* '2060104' */
			return KnxDatapointType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM
		}
	case 2060105:
		{ /* '2060105' */
			return KnxDatapointType_DPT_16_BIT_UNSIGNED_VALUE_AND_8_BIT_ENUM
		}
	case 210001:
		{ /* '210001' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 210002:
		{ /* '210002' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 210100:
		{ /* '210100' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 210101:
		{ /* '210101' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 210102:
		{ /* '210102' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 210103:
		{ /* '210103' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 210104:
		{ /* '210104' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 210105:
		{ /* '210105' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 210106:
		{ /* '210106' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 210107:
		{ /* '210107' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 210601:
		{ /* '210601' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 211000:
		{ /* '211000' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 211001:
		{ /* '211001' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 211010:
		{ /* '211010' */
			return KnxDatapointType_DPT_8_BIT_SET
		}
	case 2170001:
		{ /* '2170001' */
			return KnxDatapointType_DPT_DATAPOINT_TYPE_VERSION
		}
	case 2190001:
		{ /* '2190001' */
			return KnxDatapointType_DPT_ALARM_INFO
		}
	case 220100:
		{ /* '220100' */
			return KnxDatapointType_DPT_16_BIT_SET
		}
	case 220101:
		{ /* '220101' */
			return KnxDatapointType_DPT_16_BIT_SET
		}
	case 220102:
		{ /* '220102' */
			return KnxDatapointType_DPT_16_BIT_SET
		}
	case 220103:
		{ /* '220103' */
			return KnxDatapointType_DPT_16_BIT_SET
		}
	case 221000:
		{ /* '221000' */
			return KnxDatapointType_DPT_16_BIT_SET
		}
	case 221010:
		{ /* '221010' */
			return KnxDatapointType_DPT_16_BIT_SET
		}
	case 2220100:
		{ /* '2220100' */
			return KnxDatapointType_DPT_3X_2_BYTE_FLOAT_VALUE
		}
	case 2220101:
		{ /* '2220101' */
			return KnxDatapointType_DPT_3X_2_BYTE_FLOAT_VALUE
		}
	case 2250001:
		{ /* '2250001' */
			return KnxDatapointType_DPT_SCALING_SPEED
		}
	case 2250002:
		{ /* '2250002' */
			return KnxDatapointType_DPT_SCALING_SPEED
		}
	case 2290001:
		{ /* '2290001' */
			return KnxDatapointType_DPT_4_1_1_BYTE_COMBINED_INFORMATION
		}
	case 230001:
		{ /* '230001' */
			return KnxDatapointType_DPT_2_BIT_SET
		}
	case 230002:
		{ /* '230002' */
			return KnxDatapointType_DPT_2_BIT_SET
		}
	case 230003:
		{ /* '230003' */
			return KnxDatapointType_DPT_2_BIT_SET
		}
	case 2301000:
		{ /* '2301000' */
			return KnxDatapointType_DPT_MBUS_ADDRESS
		}
	case 230102:
		{ /* '230102' */
			return KnxDatapointType_DPT_2_BIT_SET
		}
	case 2320600:
		{ /* '2320600' */
			return KnxDatapointType_DPT_3_BYTE_COLOUR_RGB
		}
	case 2340001:
		{ /* '2340001' */
			return KnxDatapointType_DPT_LANGUAGE_CODE_ISO_639_1
		}
	case 2350001:
		{ /* '2350001' */
			return KnxDatapointType_DPT_SIGNED_VALUE_WITH_CLASSIFICATION_AND_VALIDITY
		}
	case 2360001:
		{ /* '2360001' */
			return KnxDatapointType_DPT_PRIORITISED_MODE_CONTROL
		}
	case 2370600:
		{ /* '2370600' */
			return KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_16_BIT
		}
	case 2380600:
		{ /* '2380600' */
			return KnxDatapointType_DPT_CONFIGURATION_DIAGNOSTICS_8_BIT
		}
	case 2400800:
		{ /* '2400800' */
			return KnxDatapointType_DPT_POSITIONS
		}
	case 2410800:
		{ /* '2410800' */
			return KnxDatapointType_DPT_STATUS_32_BIT
		}
	case 2420600:
		{ /* '2420600' */
			return KnxDatapointType_DPT_STATUS_48_BIT
		}
	case 2440600:
		{ /* '2440600' */
			return KnxDatapointType_DPT_CONVERTER_STATUS
		}
	case 2450600:
		{ /* '2450600' */
			return KnxDatapointType_DPT_CONVERTER_TEST_RESULT
		}
	case 2460600:
		{ /* '2460600' */
			return KnxDatapointType_DPT_BATTERY_INFORMATION
		}
	case 2490600:
		{ /* '2490600' */
			return KnxDatapointType_DPT_BRIGHTNESS_COLOUR_TEMPERATURE_TRANSITION
		}
	case 2500600:
		{ /* '2500600' */
			return KnxDatapointType_DPT_STATUS_24_BIT
		}
	case 251000:
		{ /* '251000' */
			return KnxDatapointType_DPT_2_NIBBLE_SET
		}
	case 2510600:
		{ /* '2510600' */
			return KnxDatapointType_DPT_COLOUR_RGBW
		}
	case 2520600:
		{ /* '2520600' */
			return KnxDatapointType_DPT_RELATIVE_CONTROL_RGBW
		}
	case 2540600:
		{ /* '2540600' */
			return KnxDatapointType_DPT_RELATIVE_CONTROL_RGB
		}
	case 2550001:
		{ /* '2550001' */
			return KnxDatapointType_DPT_F32F32
		}
	case 260001:
		{ /* '260001' */
			return KnxDatapointType_DPT_8_BIT_SET_2
		}
	case 270001:
		{ /* '270001' */
			return KnxDatapointType_DPT_32_BIT_SET
		}
	case 2750100:
		{ /* '2750100' */
			return KnxDatapointType_DPT_F16F16F16F16
		}
	case 2750101:
		{ /* '2750101' */
			return KnxDatapointType_DPT_F16F16F16F16
		}
	case 290010:
		{ /* '290010' */
			return KnxDatapointType_DPT_ELECTRICAL_ENERGY
		}
	case 290011:
		{ /* '290011' */
			return KnxDatapointType_DPT_ELECTRICAL_ENERGY
		}
	case 290012:
		{ /* '290012' */
			return KnxDatapointType_DPT_ELECTRICAL_ENERGY
		}
	case 30007:
		{ /* '30007' */
			return KnxDatapointType_DPT_3_BIT_CONTROLLED
		}
	case 30008:
		{ /* '30008' */
			return KnxDatapointType_DPT_3_BIT_CONTROLLED
		}
	case 301010:
		{ /* '301010' */
			return KnxDatapointType_DPT_24_TIMES_CHANNEL_ACTIVATION
		}
	case 40001:
		{ /* '40001' */
			return KnxDatapointType_DPT_CHARACTER
		}
	case 40002:
		{ /* '40002' */
			return KnxDatapointType_DPT_CHARACTER
		}
	case 50001:
		{ /* '50001' */
			return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
		}
	case 50003:
		{ /* '50003' */
			return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
		}
	case 50004:
		{ /* '50004' */
			return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
		}
	case 50005:
		{ /* '50005' */
			return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
		}
	case 50006:
		{ /* '50006' */
			return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
		}
	case 50010:
		{ /* '50010' */
			return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
		}
	case 50100:
		{ /* '50100' */
			return KnxDatapointType_DPT_8_BIT_UNSIGNED_VALUE
		}
	case 60001:
		{ /* '60001' */
			return KnxDatapointType_DPT_8_BIT_SIGNED_VALUE
		}
	case 60010:
		{ /* '60010' */
			return KnxDatapointType_DPT_8_BIT_SIGNED_VALUE
		}
	case 60020:
		{ /* '60020' */
			return KnxDatapointType_DPT_8_BIT_SIGNED_VALUE
		}
	case 70001:
		{ /* '70001' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 70002:
		{ /* '70002' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 70003:
		{ /* '70003' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 70004:
		{ /* '70004' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 70005:
		{ /* '70005' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 70006:
		{ /* '70006' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 70007:
		{ /* '70007' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 70010:
		{ /* '70010' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 70011:
		{ /* '70011' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 70012:
		{ /* '70012' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 70013:
		{ /* '70013' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 70600:
		{ /* '70600' */
			return KnxDatapointType_DPT_2_BYTE_UNSIGNED_VALUE
		}
	case 80001:
		{ /* '80001' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 80002:
		{ /* '80002' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 80003:
		{ /* '80003' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 80004:
		{ /* '80004' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 80005:
		{ /* '80005' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 80006:
		{ /* '80006' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 80007:
		{ /* '80007' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 80010:
		{ /* '80010' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 80011:
		{ /* '80011' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 80012:
		{ /* '80012' */
			return KnxDatapointType_DPT_2_BYTE_SIGNED_VALUE
		}
	case 90001:
		{ /* '90001' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90002:
		{ /* '90002' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90003:
		{ /* '90003' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90004:
		{ /* '90004' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90005:
		{ /* '90005' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90006:
		{ /* '90006' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90007:
		{ /* '90007' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90008:
		{ /* '90008' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90009:
		{ /* '90009' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90010:
		{ /* '90010' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90011:
		{ /* '90011' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90020:
		{ /* '90020' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90021:
		{ /* '90021' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90022:
		{ /* '90022' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90023:
		{ /* '90023' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90024:
		{ /* '90024' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90025:
		{ /* '90025' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90026:
		{ /* '90026' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90027:
		{ /* '90027' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90028:
		{ /* '90028' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90029:
		{ /* '90029' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	case 90030:
		{ /* '90030' */
			return KnxDatapointType_DPT_2_BYTE_FLOAT_VALUE
		}
	default:
		{
			return 0
		}
	}
}

func (e KnxDatapointSubtype) Text() string {
	switch e {
	case 100001:
		{ /* '100001' */
			return "time of day"
		}
	case 10001:
		{ /* '10001' */
			return "switch"
		}
	case 10002:
		{ /* '10002' */
			return "boolean"
		}
	case 10003:
		{ /* '10003' */
			return "enable"
		}
	case 10004:
		{ /* '10004' */
			return "ramp"
		}
	case 10005:
		{ /* '10005' */
			return "alarm"
		}
	case 10006:
		{ /* '10006' */
			return "binary value"
		}
	case 10007:
		{ /* '10007' */
			return "step"
		}
	case 10008:
		{ /* '10008' */
			return "up/down"
		}
	case 10009:
		{ /* '10009' */
			return "open/close"
		}
	case 10010:
		{ /* '10010' */
			return "start/stop"
		}
	case 10011:
		{ /* '10011' */
			return "state"
		}
	case 10012:
		{ /* '10012' */
			return "invert"
		}
	case 10013:
		{ /* '10013' */
			return "dim send style"
		}
	case 10014:
		{ /* '10014' */
			return "input source"
		}
	case 10015:
		{ /* '10015' */
			return "reset"
		}
	case 10016:
		{ /* '10016' */
			return "acknowledge"
		}
	case 10017:
		{ /* '10017' */
			return "trigger"
		}
	case 10018:
		{ /* '10018' */
			return "occupancy"
		}
	case 10019:
		{ /* '10019' */
			return "window/door"
		}
	case 10021:
		{ /* '10021' */
			return "logical function"
		}
	case 10022:
		{ /* '10022' */
			return "scene"
		}
	case 10023:
		{ /* '10023' */
			return "shutter/blinds mode"
		}
	case 10024:
		{ /* '10024' */
			return "day/night"
		}
	case 10100:
		{ /* '10100' */
			return "cooling/heating"
		}
	case 110001:
		{ /* '110001' */
			return "date"
		}
	case 120001:
		{ /* '120001' */
			return "counter pulses (unsigned)"
		}
	case 120100:
		{ /* '120100' */
			return "counter timesec (s)"
		}
	case 120101:
		{ /* '120101' */
			return "counter timemin (min)"
		}
	case 120102:
		{ /* '120102' */
			return "counter timehrs (h)"
		}
	case 121200:
		{ /* '121200' */
			return "volume liquid (l)"
		}
	case 121201:
		{ /* '121201' */
			return "volume (m³)"
		}
	case 130001:
		{ /* '130001' */
			return "counter pulses (signed)"
		}
	case 130002:
		{ /* '130002' */
			return "flow rate (m³/h)"
		}
	case 130010:
		{ /* '130010' */
			return "active energy (Wh)"
		}
	case 130011:
		{ /* '130011' */
			return "apparant energy (VAh)"
		}
	case 130012:
		{ /* '130012' */
			return "reactive energy (VARh)"
		}
	case 130013:
		{ /* '130013' */
			return "active energy (kWh)"
		}
	case 130014:
		{ /* '130014' */
			return "apparant energy (kVAh)"
		}
	case 130015:
		{ /* '130015' */
			return "reactive energy (kVARh)"
		}
	case 130016:
		{ /* '130016' */
			return "active energy (MWh)"
		}
	case 130100:
		{ /* '130100' */
			return "time lag (s)"
		}
	case 131200:
		{ /* '131200' */
			return "delta volume liquid (l)"
		}
	case 131201:
		{ /* '131201' */
			return "delta volume (m³)"
		}
	case 140000:
		{ /* '140000' */
			return "acceleration (m/s²)"
		}
	case 140001:
		{ /* '140001' */
			return "angular acceleration (rad/s²)"
		}
	case 140002:
		{ /* '140002' */
			return "activation energy (J/mol)"
		}
	case 140003:
		{ /* '140003' */
			return "radioactive activity (1/s)"
		}
	case 140004:
		{ /* '140004' */
			return "amount of substance (mol)"
		}
	case 140005:
		{ /* '140005' */
			return "amplitude"
		}
	case 140006:
		{ /* '140006' */
			return "angle (radiant)"
		}
	case 140007:
		{ /* '140007' */
			return "angle (degree)"
		}
	case 140008:
		{ /* '140008' */
			return "angular momentum (Js)"
		}
	case 140009:
		{ /* '140009' */
			return "angular velocity (rad/s)"
		}
	case 140010:
		{ /* '140010' */
			return "area (m*m)"
		}
	case 140011:
		{ /* '140011' */
			return "capacitance (F)"
		}
	case 140012:
		{ /* '140012' */
			return "flux density (C/m²)"
		}
	case 140013:
		{ /* '140013' */
			return "charge density (C/m³)"
		}
	case 140014:
		{ /* '140014' */
			return "compressibility (m²/N)"
		}
	case 140015:
		{ /* '140015' */
			return "conductance (S)"
		}
	case 140016:
		{ /* '140016' */
			return "conductivity  (S/m)"
		}
	case 140017:
		{ /* '140017' */
			return "density (kg/m³)"
		}
	case 140018:
		{ /* '140018' */
			return "electric charge (C)"
		}
	case 140019:
		{ /* '140019' */
			return "electric current (A)"
		}
	case 140020:
		{ /* '140020' */
			return "electric current density (A/m²)"
		}
	case 140021:
		{ /* '140021' */
			return "electric dipole moment (Cm)"
		}
	case 140022:
		{ /* '140022' */
			return "electric displacement (C/m²)"
		}
	case 140023:
		{ /* '140023' */
			return "electric field strength (V/m)"
		}
	case 140024:
		{ /* '140024' */
			return "electric flux (C)"
		}
	case 140025:
		{ /* '140025' */
			return "electric flux density (C/m²)"
		}
	case 140026:
		{ /* '140026' */
			return "electric polarization (C/m²)"
		}
	case 140027:
		{ /* '140027' */
			return "electric potential (V)"
		}
	case 140028:
		{ /* '140028' */
			return "electric potential difference (V)"
		}
	case 140029:
		{ /* '140029' */
			return "electromagnetic moment (Am²)"
		}
	case 140030:
		{ /* '140030' */
			return "electromotive force (V)"
		}
	case 140031:
		{ /* '140031' */
			return "energy (J)"
		}
	case 140032:
		{ /* '140032' */
			return "force (N)"
		}
	case 140033:
		{ /* '140033' */
			return "frequency (Hz)"
		}
	case 140034:
		{ /* '140034' */
			return "angular frequency (rad/s)"
		}
	case 140035:
		{ /* '140035' */
			return "heat capacity (J/K)"
		}
	case 140036:
		{ /* '140036' */
			return "heat flow rate (W)"
		}
	case 140037:
		{ /* '140037' */
			return "heat quantity"
		}
	case 140038:
		{ /* '140038' */
			return "impedance (Ω)"
		}
	case 140039:
		{ /* '140039' */
			return "length (m)"
		}
	case 140040:
		{ /* '140040' */
			return "light quantity (J)"
		}
	case 140041:
		{ /* '140041' */
			return "luminance (cd/m²)"
		}
	case 140042:
		{ /* '140042' */
			return "luminous flux (lm)"
		}
	case 140043:
		{ /* '140043' */
			return "luminous intensity (cd)"
		}
	case 140044:
		{ /* '140044' */
			return "magnetic field strength (A/m)"
		}
	case 140045:
		{ /* '140045' */
			return "magnetic flux (Wb)"
		}
	case 140046:
		{ /* '140046' */
			return "magnetic flux density (T)"
		}
	case 140047:
		{ /* '140047' */
			return "magnetic moment (Am²)"
		}
	case 140048:
		{ /* '140048' */
			return "magnetic polarization (T)"
		}
	case 140049:
		{ /* '140049' */
			return "magnetization (A/m)"
		}
	case 140050:
		{ /* '140050' */
			return "magnetomotive force (A)"
		}
	case 140051:
		{ /* '140051' */
			return "mass (kg)"
		}
	case 140052:
		{ /* '140052' */
			return "mass flux (kg/s)"
		}
	case 140053:
		{ /* '140053' */
			return "momentum (N/s)"
		}
	case 140054:
		{ /* '140054' */
			return "phase angle (rad)"
		}
	case 140055:
		{ /* '140055' */
			return "phase angle (°)"
		}
	case 140056:
		{ /* '140056' */
			return "power (W)"
		}
	case 140057:
		{ /* '140057' */
			return "power factor (cos Φ)"
		}
	case 140058:
		{ /* '140058' */
			return "pressure (Pa)"
		}
	case 140059:
		{ /* '140059' */
			return "reactance (Ω)"
		}
	case 140060:
		{ /* '140060' */
			return "resistance (Ω)"
		}
	case 140061:
		{ /* '140061' */
			return "resistivity (Ωm)"
		}
	case 140062:
		{ /* '140062' */
			return "self inductance (H)"
		}
	case 140063:
		{ /* '140063' */
			return "solid angle (sr)"
		}
	case 140064:
		{ /* '140064' */
			return "sound intensity (W/m²)"
		}
	case 140065:
		{ /* '140065' */
			return "speed (m/s)"
		}
	case 140066:
		{ /* '140066' */
			return "stress (Pa)"
		}
	case 140067:
		{ /* '140067' */
			return "surface tension (N/m)"
		}
	case 140068:
		{ /* '140068' */
			return "temperature (°C)"
		}
	case 140069:
		{ /* '140069' */
			return "temperature absolute (K)"
		}
	case 140070:
		{ /* '140070' */
			return "temperature difference (K)"
		}
	case 140071:
		{ /* '140071' */
			return "thermal capacity (J/K)"
		}
	case 140072:
		{ /* '140072' */
			return "thermal conductivity (W/mK)"
		}
	case 140073:
		{ /* '140073' */
			return "thermoelectric power (V/K)"
		}
	case 140074:
		{ /* '140074' */
			return "time (s)"
		}
	case 140075:
		{ /* '140075' */
			return "torque (Nm)"
		}
	case 140076:
		{ /* '140076' */
			return "volume (m³)"
		}
	case 140077:
		{ /* '140077' */
			return "volume flux (m³/s)"
		}
	case 140078:
		{ /* '140078' */
			return "weight (N)"
		}
	case 140079:
		{ /* '140079' */
			return "work (J)"
		}
	case 141200:
		{ /* '141200' */
			return "volume flux for meters (m³/h)"
		}
	case 141201:
		{ /* '141201' */
			return "volume flux for meters (1/ls)"
		}
	case 150000:
		{ /* '150000' */
			return "entrance access"
		}
	case 160000:
		{ /* '160000' */
			return "Character String (ASCII)"
		}
	case 160001:
		{ /* '160001' */
			return "Character String (ISO 8859-1)"
		}
	case 170001:
		{ /* '170001' */
			return "scene number"
		}
	case 180001:
		{ /* '180001' */
			return "scene control"
		}
	case 190001:
		{ /* '190001' */
			return "date time"
		}
	case 200001:
		{ /* '200001' */
			return "SCLO mode"
		}
	case 200002:
		{ /* '200002' */
			return "building mode"
		}
	case 200003:
		{ /* '200003' */
			return "occupied"
		}
	case 200004:
		{ /* '200004' */
			return "priority"
		}
	case 200005:
		{ /* '200005' */
			return "light application mode"
		}
	case 200006:
		{ /* '200006' */
			return "light application area"
		}
	case 200007:
		{ /* '200007' */
			return "alarm class"
		}
	case 200008:
		{ /* '200008' */
			return "PSU mode"
		}
	case 20001:
		{ /* '20001' */
			return "switch control"
		}
	case 200011:
		{ /* '200011' */
			return "system error class"
		}
	case 200012:
		{ /* '200012' */
			return "HVAC error class"
		}
	case 200013:
		{ /* '200013' */
			return "time delay"
		}
	case 200014:
		{ /* '200014' */
			return "wind force scale (0..12)"
		}
	case 200017:
		{ /* '200017' */
			return "sensor mode"
		}
	case 20002:
		{ /* '20002' */
			return "boolean control"
		}
	case 200020:
		{ /* '200020' */
			return "actuator connect type"
		}
	case 200021:
		{ /* '200021' */
			return "cloud cover"
		}
	case 200022:
		{ /* '200022' */
			return "power return mode"
		}
	case 20003:
		{ /* '20003' */
			return "enable control"
		}
	case 20004:
		{ /* '20004' */
			return "ramp control"
		}
	case 20005:
		{ /* '20005' */
			return "alarm control"
		}
	case 20006:
		{ /* '20006' */
			return "binary value control"
		}
	case 20007:
		{ /* '20007' */
			return "step control"
		}
	case 20008:
		{ /* '20008' */
			return "direction control 1"
		}
	case 20009:
		{ /* '20009' */
			return "direction control 2"
		}
	case 20010:
		{ /* '20010' */
			return "start control"
		}
	case 200100:
		{ /* '200100' */
			return "fuel type"
		}
	case 200101:
		{ /* '200101' */
			return "burner type"
		}
	case 200102:
		{ /* '200102' */
			return "HVAC mode"
		}
	case 200103:
		{ /* '200103' */
			return "DHW mode"
		}
	case 200104:
		{ /* '200104' */
			return "load priority"
		}
	case 200105:
		{ /* '200105' */
			return "HVAC control mode"
		}
	case 200106:
		{ /* '200106' */
			return "HVAC emergency mode"
		}
	case 200107:
		{ /* '200107' */
			return "changeover mode"
		}
	case 200108:
		{ /* '200108' */
			return "valve mode"
		}
	case 200109:
		{ /* '200109' */
			return "damper mode"
		}
	case 20011:
		{ /* '20011' */
			return "state control"
		}
	case 200110:
		{ /* '200110' */
			return "heater mode"
		}
	case 200111:
		{ /* '200111' */
			return "fan mode"
		}
	case 200112:
		{ /* '200112' */
			return "master/slave mode"
		}
	case 200113:
		{ /* '200113' */
			return "status room setpoint"
		}
	case 200114:
		{ /* '200114' */
			return "metering device type"
		}
	case 200115:
		{ /* '200115' */
			return "hum dehum mode"
		}
	case 200116:
		{ /* '200116' */
			return "enable H/C stage"
		}
	case 20012:
		{ /* '20012' */
			return "invert control"
		}
	case 200120:
		{ /* '200120' */
			return "ADA type"
		}
	case 200121:
		{ /* '200121' */
			return "backup mode"
		}
	case 200122:
		{ /* '200122' */
			return "start syncronization type"
		}
	case 200600:
		{ /* '200600' */
			return "behavior lock/unlock"
		}
	case 200601:
		{ /* '200601' */
			return "behavior bus power up/down"
		}
	case 200602:
		{ /* '200602' */
			return "dali fade time"
		}
	case 200603:
		{ /* '200603' */
			return "blink mode"
		}
	case 200604:
		{ /* '200604' */
			return "light control mode"
		}
	case 200605:
		{ /* '200605' */
			return "PB switch mode"
		}
	case 200606:
		{ /* '200606' */
			return "PB action mode"
		}
	case 200607:
		{ /* '200607' */
			return "PB dimm mode"
		}
	case 200608:
		{ /* '200608' */
			return "switch on mode"
		}
	case 200609:
		{ /* '200609' */
			return "load type"
		}
	case 200610:
		{ /* '200610' */
			return "load type detection"
		}
	case 200611:
		{ /* '200611' */
			return "converter test control"
		}
	case 200801:
		{ /* '200801' */
			return "SAB except behavior"
		}
	case 200802:
		{ /* '200802' */
			return "SAB behavior on lock/unlock"
		}
	case 200803:
		{ /* '200803' */
			return "SSSB mode"
		}
	case 200804:
		{ /* '200804' */
			return "blinds control mode"
		}
	case 201000:
		{ /* '201000' */
			return "communication mode"
		}
	case 201001:
		{ /* '201001' */
			return "additional information type"
		}
	case 201002:
		{ /* '201002' */
			return "RF mode selection"
		}
	case 201003:
		{ /* '201003' */
			return "RF filter mode selection"
		}
	case 2060100:
		{ /* '2060100' */
			return "time delay & HVAC mode"
		}
	case 2060102:
		{ /* '2060102' */
			return "time delay & DHW mode"
		}
	case 2060104:
		{ /* '2060104' */
			return "time delay & occupancy mode"
		}
	case 2060105:
		{ /* '2060105' */
			return "time delay & building mode"
		}
	case 210001:
		{ /* '210001' */
			return "general status"
		}
	case 210002:
		{ /* '210002' */
			return "device control"
		}
	case 210100:
		{ /* '210100' */
			return "forcing signal"
		}
	case 210101:
		{ /* '210101' */
			return "forcing signal cool"
		}
	case 210102:
		{ /* '210102' */
			return "room heating controller status"
		}
	case 210103:
		{ /* '210103' */
			return "solar DHW controller status"
		}
	case 210104:
		{ /* '210104' */
			return "fuel type set"
		}
	case 210105:
		{ /* '210105' */
			return "room cooling controller status"
		}
	case 210106:
		{ /* '210106' */
			return "ventilation controller status"
		}
	case 210107:
		{ /* '210107' */
			return "combined status RTSM"
		}
	case 210601:
		{ /* '210601' */
			return "lighting actuator error information"
		}
	case 211000:
		{ /* '211000' */
			return "RF communication mode info"
		}
	case 211001:
		{ /* '211001' */
			return "cEMI server supported RF filtering modes"
		}
	case 211010:
		{ /* '211010' */
			return "channel activation for 8 channels"
		}
	case 2170001:
		{ /* '2170001' */
			return "DPT version"
		}
	case 2190001:
		{ /* '2190001' */
			return "alarm info"
		}
	case 220100:
		{ /* '220100' */
			return "DHW controller status"
		}
	case 220101:
		{ /* '220101' */
			return "RHCC status"
		}
	case 220102:
		{ /* '220102' */
			return "combined status HVA"
		}
	case 220103:
		{ /* '220103' */
			return "combined status RTC"
		}
	case 221000:
		{ /* '221000' */
			return "media"
		}
	case 221010:
		{ /* '221010' */
			return "channel activation for 16 channels"
		}
	case 2220100:
		{ /* '2220100' */
			return "room temperature setpoint"
		}
	case 2220101:
		{ /* '2220101' */
			return "room temperature setpoint shift"
		}
	case 2250001:
		{ /* '2250001' */
			return "scaling speed"
		}
	case 2250002:
		{ /* '2250002' */
			return "scaling step time"
		}
	case 2290001:
		{ /* '2290001' */
			return "metering value (value,encoding,cmd)"
		}
	case 230001:
		{ /* '230001' */
			return "on/off action"
		}
	case 230002:
		{ /* '230002' */
			return "alarm reaction"
		}
	case 230003:
		{ /* '230003' */
			return "up/down action"
		}
	case 2301000:
		{ /* '2301000' */
			return "MBus address"
		}
	case 230102:
		{ /* '230102' */
			return "HVAC push button action"
		}
	case 2320600:
		{ /* '2320600' */
			return "RGB value 3x(0..255)"
		}
	case 2340001:
		{ /* '2340001' */
			return "language code (ASCII)"
		}
	case 2350001:
		{ /* '2350001' */
			return "electrical energy with tariff"
		}
	case 2360001:
		{ /* '2360001' */
			return "priority control"
		}
	case 2370600:
		{ /* '2370600' */
			return "diagnostic value"
		}
	case 2380600:
		{ /* '2380600' */
			return "diagnostic value"
		}
	case 2400800:
		{ /* '2400800' */
			return "combined position"
		}
	case 2410800:
		{ /* '2410800' */
			return "status sunblind & shutter actuator"
		}
	case 2420600:
		{ /* '2420600' */
			return "colour xyY"
		}
	case 2440600:
		{ /* '2440600' */
			return "DALI converter status"
		}
	case 2450600:
		{ /* '2450600' */
			return "DALI converter test result"
		}
	case 2460600:
		{ /* '2460600' */
			return "Battery Information"
		}
	case 2490600:
		{ /* '2490600' */
			return "brightness colour temperature transition"
		}
	case 2500600:
		{ /* '2500600' */
			return "brightness colour temperature control"
		}
	case 251000:
		{ /* '251000' */
			return "busy/nak repetitions"
		}
	case 2510600:
		{ /* '2510600' */
			return "RGBW value 4x(0..100%)"
		}
	case 2520600:
		{ /* '2520600' */
			return "RGBW relative control"
		}
	case 2540600:
		{ /* '2540600' */
			return "RGB relative control"
		}
	case 2550001:
		{ /* '2550001' */
			return "geographical location (longitude and latitude) expressed in degrees"
		}
	case 260001:
		{ /* '260001' */
			return "scene information"
		}
	case 270001:
		{ /* '270001' */
			return "bit-combined info on/off"
		}
	case 2750100:
		{ /* '2750100' */
			return "Temperature setpoint setting for 4 HVAC Modes"
		}
	case 2750101:
		{ /* '2750101' */
			return "Temperature setpoint shift setting for 4 HVAC Modes"
		}
	case 290010:
		{ /* '290010' */
			return "active energy (Wh)"
		}
	case 290011:
		{ /* '290011' */
			return "apparant energy (VAh)"
		}
	case 290012:
		{ /* '290012' */
			return "reactive energy (VARh)"
		}
	case 30007:
		{ /* '30007' */
			return "dimming control"
		}
	case 30008:
		{ /* '30008' */
			return "blind control"
		}
	case 301010:
		{ /* '301010' */
			return "activation state 0..23"
		}
	case 40001:
		{ /* '40001' */
			return "character (ASCII)"
		}
	case 40002:
		{ /* '40002' */
			return "character (ISO 8859-1)"
		}
	case 50001:
		{ /* '50001' */
			return "percentage (0..100%)"
		}
	case 50003:
		{ /* '50003' */
			return "angle (degrees)"
		}
	case 50004:
		{ /* '50004' */
			return "percentage (0..255%)"
		}
	case 50005:
		{ /* '50005' */
			return "ratio (0..255)"
		}
	case 50006:
		{ /* '50006' */
			return "tariff (0..255)"
		}
	case 50010:
		{ /* '50010' */
			return "counter pulses (0..255)"
		}
	case 50100:
		{ /* '50100' */
			return "fan stage (0..255)"
		}
	case 60001:
		{ /* '60001' */
			return "percentage (-128..127%)"
		}
	case 60010:
		{ /* '60010' */
			return "counter pulses (-128..127)"
		}
	case 60020:
		{ /* '60020' */
			return "status with mode"
		}
	case 70001:
		{ /* '70001' */
			return "pulses"
		}
	case 70002:
		{ /* '70002' */
			return "time (ms)"
		}
	case 70003:
		{ /* '70003' */
			return "time (10 ms)"
		}
	case 70004:
		{ /* '70004' */
			return "time (100 ms)"
		}
	case 70005:
		{ /* '70005' */
			return "time (s)"
		}
	case 70006:
		{ /* '70006' */
			return "time (min)"
		}
	case 70007:
		{ /* '70007' */
			return "time (h)"
		}
	case 70010:
		{ /* '70010' */
			return "property data type"
		}
	case 70011:
		{ /* '70011' */
			return "length (mm)"
		}
	case 70012:
		{ /* '70012' */
			return "current (mA)"
		}
	case 70013:
		{ /* '70013' */
			return "brightness (lux)"
		}
	case 70600:
		{ /* '70600' */
			return "absolute colour temperature (K)"
		}
	case 80001:
		{ /* '80001' */
			return "pulses difference"
		}
	case 80002:
		{ /* '80002' */
			return "time lag (ms)"
		}
	case 80003:
		{ /* '80003' */
			return "time lag(10 ms)"
		}
	case 80004:
		{ /* '80004' */
			return "time lag (100 ms)"
		}
	case 80005:
		{ /* '80005' */
			return "time lag (s)"
		}
	case 80006:
		{ /* '80006' */
			return "time lag (min)"
		}
	case 80007:
		{ /* '80007' */
			return "time lag (h)"
		}
	case 80010:
		{ /* '80010' */
			return "percentage difference (%)"
		}
	case 80011:
		{ /* '80011' */
			return "rotation angle (°)"
		}
	case 80012:
		{ /* '80012' */
			return "length (m)"
		}
	case 90001:
		{ /* '90001' */
			return "temperature (°C)"
		}
	case 90002:
		{ /* '90002' */
			return "temperature difference (K)"
		}
	case 90003:
		{ /* '90003' */
			return "kelvin/hour (K/h)"
		}
	case 90004:
		{ /* '90004' */
			return "lux (Lux)"
		}
	case 90005:
		{ /* '90005' */
			return "speed (m/s)"
		}
	case 90006:
		{ /* '90006' */
			return "pressure (Pa)"
		}
	case 90007:
		{ /* '90007' */
			return "humidity (%)"
		}
	case 90008:
		{ /* '90008' */
			return "parts/million (ppm)"
		}
	case 90009:
		{ /* '90009' */
			return "air flow (m³/h)"
		}
	case 90010:
		{ /* '90010' */
			return "time (s)"
		}
	case 90011:
		{ /* '90011' */
			return "time (ms)"
		}
	case 90020:
		{ /* '90020' */
			return "voltage (mV)"
		}
	case 90021:
		{ /* '90021' */
			return "current (mA)"
		}
	case 90022:
		{ /* '90022' */
			return "power denisity (W/m²)"
		}
	case 90023:
		{ /* '90023' */
			return "kelvin/percent (K/%)"
		}
	case 90024:
		{ /* '90024' */
			return "power (kW)"
		}
	case 90025:
		{ /* '90025' */
			return "volume flow (l/h)"
		}
	case 90026:
		{ /* '90026' */
			return "rain amount (l/m²)"
		}
	case 90027:
		{ /* '90027' */
			return "temperature (°F)"
		}
	case 90028:
		{ /* '90028' */
			return "wind speed (km/h)"
		}
	case 90029:
		{ /* '90029' */
			return "absolute humidity (g/m³)"
		}
	case 90030:
		{ /* '90030' */
			return "concentration (µg/m³)"
		}
	default:
		{
			return ""
		}
	}
}
func KnxDatapointSubtypeValueOf(value uint32) KnxDatapointSubtype {
	switch value {
	case 100001:
		return KnxDatapointSubtype_DPST_TimeOfDay
	case 10001:
		return KnxDatapointSubtype_DPST_Switch
	case 10002:
		return KnxDatapointSubtype_DPST_Bool
	case 10003:
		return KnxDatapointSubtype_DPST_Enable
	case 10004:
		return KnxDatapointSubtype_DPST_Ramp
	case 10005:
		return KnxDatapointSubtype_DPST_Alarm
	case 10006:
		return KnxDatapointSubtype_DPST_BinaryValue
	case 10007:
		return KnxDatapointSubtype_DPST_Step
	case 10008:
		return KnxDatapointSubtype_DPST_UpDown
	case 10009:
		return KnxDatapointSubtype_DPST_OpenClose
	case 10010:
		return KnxDatapointSubtype_DPST_Start
	case 10011:
		return KnxDatapointSubtype_DPST_State
	case 10012:
		return KnxDatapointSubtype_DPST_Invert
	case 10013:
		return KnxDatapointSubtype_DPST_DimSendStyle
	case 10014:
		return KnxDatapointSubtype_DPST_InputSource
	case 10015:
		return KnxDatapointSubtype_DPST_Reset
	case 10016:
		return KnxDatapointSubtype_DPST_Ack
	case 10017:
		return KnxDatapointSubtype_DPST_Trigger
	case 10018:
		return KnxDatapointSubtype_DPST_Occupancy
	case 10019:
		return KnxDatapointSubtype_DPST_Window_Door
	case 10021:
		return KnxDatapointSubtype_DPST_LogicalFunction
	case 10022:
		return KnxDatapointSubtype_DPST_Scene_AB
	case 10023:
		return KnxDatapointSubtype_DPST_ShutterBlinds_Mode
	case 10024:
		return KnxDatapointSubtype_DPST_DayNight
	case 10100:
		return KnxDatapointSubtype_DPST_Heat_Cool
	case 110001:
		return KnxDatapointSubtype_DPST_Date
	case 120001:
		return KnxDatapointSubtype_DPST_Value_4_Ucount
	case 120100:
		return KnxDatapointSubtype_DPST_LongTimePeriod_Sec
	case 120101:
		return KnxDatapointSubtype_DPST_LongTimePeriod_Min
	case 120102:
		return KnxDatapointSubtype_DPST_LongTimePeriod_Hrs
	case 121200:
		return KnxDatapointSubtype_DPST_VolumeLiquid_Litre
	case 121201:
		return KnxDatapointSubtype_DPST_Volume_m_3
	case 130001:
		return KnxDatapointSubtype_DPST_Value_4_Count
	case 130002:
		return KnxDatapointSubtype_DPST_FlowRate_m3h
	case 130010:
		return KnxDatapointSubtype_DPST_ActiveEnergy
	case 130011:
		return KnxDatapointSubtype_DPST_ApparantEnergy
	case 130012:
		return KnxDatapointSubtype_DPST_ReactiveEnergy
	case 130013:
		return KnxDatapointSubtype_DPST_ActiveEnergy_kWh
	case 130014:
		return KnxDatapointSubtype_DPST_ApparantEnergy_kVAh
	case 130015:
		return KnxDatapointSubtype_DPST_ReactiveEnergy_kVARh
	case 130016:
		return KnxDatapointSubtype_DPST_ActiveEnergy_MWh
	case 130100:
		return KnxDatapointSubtype_DPST_LongDeltaTimeSec
	case 131200:
		return KnxDatapointSubtype_DPST_DeltaVolumeLiquid_Litre
	case 131201:
		return KnxDatapointSubtype_DPST_DeltaVolume_m_3
	case 140000:
		return KnxDatapointSubtype_DPST_Value_Acceleration
	case 140001:
		return KnxDatapointSubtype_DPST_Value_Acceleration_Angular
	case 140002:
		return KnxDatapointSubtype_DPST_Value_Activation_Energy
	case 140003:
		return KnxDatapointSubtype_DPST_Value_Activity
	case 140004:
		return KnxDatapointSubtype_DPST_Value_Mol
	case 140005:
		return KnxDatapointSubtype_DPST_Value_Amplitude
	case 140006:
		return KnxDatapointSubtype_DPST_Value_AngleRad
	case 140007:
		return KnxDatapointSubtype_DPST_Value_AngleDeg
	case 140008:
		return KnxDatapointSubtype_DPST_Value_Angular_Momentum
	case 140009:
		return KnxDatapointSubtype_DPST_Value_Angular_Velocity
	case 140010:
		return KnxDatapointSubtype_DPST_Value_Area
	case 140011:
		return KnxDatapointSubtype_DPST_Value_Capacitance
	case 140012:
		return KnxDatapointSubtype_DPST_Value_Charge_DensitySurface
	case 140013:
		return KnxDatapointSubtype_DPST_Value_Charge_DensityVolume
	case 140014:
		return KnxDatapointSubtype_DPST_Value_Compressibility
	case 140015:
		return KnxDatapointSubtype_DPST_Value_Conductance
	case 140016:
		return KnxDatapointSubtype_DPST_Value_Electrical_Conductivity
	case 140017:
		return KnxDatapointSubtype_DPST_Value_Density
	case 140018:
		return KnxDatapointSubtype_DPST_Value_Electric_Charge
	case 140019:
		return KnxDatapointSubtype_DPST_Value_Electric_Current
	case 140020:
		return KnxDatapointSubtype_DPST_Value_Electric_CurrentDensity
	case 140021:
		return KnxDatapointSubtype_DPST_Value_Electric_DipoleMoment
	case 140022:
		return KnxDatapointSubtype_DPST_Value_Electric_Displacement
	case 140023:
		return KnxDatapointSubtype_DPST_Value_Electric_FieldStrength
	case 140024:
		return KnxDatapointSubtype_DPST_Value_Electric_Flux
	case 140025:
		return KnxDatapointSubtype_DPST_Value_Electric_FluxDensity
	case 140026:
		return KnxDatapointSubtype_DPST_Value_Electric_Polarization
	case 140027:
		return KnxDatapointSubtype_DPST_Value_Electric_Potential
	case 140028:
		return KnxDatapointSubtype_DPST_Value_Electric_PotentialDifference
	case 140029:
		return KnxDatapointSubtype_DPST_Value_ElectromagneticMoment
	case 140030:
		return KnxDatapointSubtype_DPST_Value_Electromotive_Force
	case 140031:
		return KnxDatapointSubtype_DPST_Value_Energy
	case 140032:
		return KnxDatapointSubtype_DPST_Value_Force
	case 140033:
		return KnxDatapointSubtype_DPST_Value_Frequency
	case 140034:
		return KnxDatapointSubtype_DPST_Value_Angular_Frequency
	case 140035:
		return KnxDatapointSubtype_DPST_Value_Heat_Capacity
	case 140036:
		return KnxDatapointSubtype_DPST_Value_Heat_FlowRate
	case 140037:
		return KnxDatapointSubtype_DPST_Value_Heat_Quantity
	case 140038:
		return KnxDatapointSubtype_DPST_Value_Impedance
	case 140039:
		return KnxDatapointSubtype_DPST_Value_Length
	case 140040:
		return KnxDatapointSubtype_DPST_Value_Light_Quantity
	case 140041:
		return KnxDatapointSubtype_DPST_Value_Luminance
	case 140042:
		return KnxDatapointSubtype_DPST_Value_Luminous_Flux
	case 140043:
		return KnxDatapointSubtype_DPST_Value_Luminous_Intensity
	case 140044:
		return KnxDatapointSubtype_DPST_Value_Magnetic_FieldStrength
	case 140045:
		return KnxDatapointSubtype_DPST_Value_Magnetic_Flux
	case 140046:
		return KnxDatapointSubtype_DPST_Value_Magnetic_FluxDensity
	case 140047:
		return KnxDatapointSubtype_DPST_Value_Magnetic_Moment
	case 140048:
		return KnxDatapointSubtype_DPST_Value_Magnetic_Polarization
	case 140049:
		return KnxDatapointSubtype_DPST_Value_Magnetization
	case 140050:
		return KnxDatapointSubtype_DPST_Value_MagnetomotiveForce
	case 140051:
		return KnxDatapointSubtype_DPST_Value_Mass
	case 140052:
		return KnxDatapointSubtype_DPST_Value_MassFlux
	case 140053:
		return KnxDatapointSubtype_DPST_Value_Momentum
	case 140054:
		return KnxDatapointSubtype_DPST_Value_Phase_AngleRad
	case 140055:
		return KnxDatapointSubtype_DPST_Value_Phase_AngleDeg
	case 140056:
		return KnxDatapointSubtype_DPST_Value_Power
	case 140057:
		return KnxDatapointSubtype_DPST_Value_Power_Factor
	case 140058:
		return KnxDatapointSubtype_DPST_Value_Pressure
	case 140059:
		return KnxDatapointSubtype_DPST_Value_Reactance
	case 140060:
		return KnxDatapointSubtype_DPST_Value_Resistance
	case 140061:
		return KnxDatapointSubtype_DPST_Value_Resistivity
	case 140062:
		return KnxDatapointSubtype_DPST_Value_SelfInductance
	case 140063:
		return KnxDatapointSubtype_DPST_Value_SolidAngle
	case 140064:
		return KnxDatapointSubtype_DPST_Value_Sound_Intensity
	case 140065:
		return KnxDatapointSubtype_DPST_Value_Speed
	case 140066:
		return KnxDatapointSubtype_DPST_Value_Stress
	case 140067:
		return KnxDatapointSubtype_DPST_Value_Surface_Tension
	case 140068:
		return KnxDatapointSubtype_DPST_Value_Common_Temperature
	case 140069:
		return KnxDatapointSubtype_DPST_Value_Absolute_Temperature
	case 140070:
		return KnxDatapointSubtype_DPST_Value_TemperatureDifference
	case 140071:
		return KnxDatapointSubtype_DPST_Value_Thermal_Capacity
	case 140072:
		return KnxDatapointSubtype_DPST_Value_Thermal_Conductivity
	case 140073:
		return KnxDatapointSubtype_DPST_Value_ThermoelectricPower
	case 140074:
		return KnxDatapointSubtype_DPST_Value_Time
	case 140075:
		return KnxDatapointSubtype_DPST_Value_Torque
	case 140076:
		return KnxDatapointSubtype_DPST_Value_Volume
	case 140077:
		return KnxDatapointSubtype_DPST_Value_Volume_Flux
	case 140078:
		return KnxDatapointSubtype_DPST_Value_Weight
	case 140079:
		return KnxDatapointSubtype_DPST_Value_Work
	case 141200:
		return KnxDatapointSubtype_DPST_Volume_Flux_Meter
	case 141201:
		return KnxDatapointSubtype_DPST_Volume_Flux_ls
	case 150000:
		return KnxDatapointSubtype_DPST_Access_Data
	case 160000:
		return KnxDatapointSubtype_DPST_String_ASCII
	case 160001:
		return KnxDatapointSubtype_DPST_String_8859_1
	case 170001:
		return KnxDatapointSubtype_DPST_SceneNumber
	case 180001:
		return KnxDatapointSubtype_DPST_SceneControl
	case 190001:
		return KnxDatapointSubtype_DPST_DateTime
	case 200001:
		return KnxDatapointSubtype_DPST_SCLOMode
	case 200002:
		return KnxDatapointSubtype_DPST_BuildingMode
	case 200003:
		return KnxDatapointSubtype_DPST_OccMode
	case 200004:
		return KnxDatapointSubtype_DPST_Priority
	case 200005:
		return KnxDatapointSubtype_DPST_LightApplicationMode
	case 200006:
		return KnxDatapointSubtype_DPST_ApplicationArea
	case 200007:
		return KnxDatapointSubtype_DPST_AlarmClassType
	case 200008:
		return KnxDatapointSubtype_DPST_PSUMode
	case 20001:
		return KnxDatapointSubtype_DPST_Switch_Control
	case 200011:
		return KnxDatapointSubtype_DPST_ErrorClass_System
	case 200012:
		return KnxDatapointSubtype_DPST_ErrorClass_HVAC
	case 200013:
		return KnxDatapointSubtype_DPST_Time_Delay
	case 200014:
		return KnxDatapointSubtype_DPST_Beaufort_Wind_Force_Scale
	case 200017:
		return KnxDatapointSubtype_DPST_SensorSelect
	case 20002:
		return KnxDatapointSubtype_DPST_Bool_Control
	case 200020:
		return KnxDatapointSubtype_DPST_ActuatorConnectType
	case 200021:
		return KnxDatapointSubtype_DPST_Cloud_Cover
	case 200022:
		return KnxDatapointSubtype_DPST_PowerReturnMode
	case 20003:
		return KnxDatapointSubtype_DPST_Enable_Control
	case 20004:
		return KnxDatapointSubtype_DPST_Ramp_Control
	case 20005:
		return KnxDatapointSubtype_DPST_Alarm_Control
	case 20006:
		return KnxDatapointSubtype_DPST_BinaryValue_Control
	case 20007:
		return KnxDatapointSubtype_DPST_Step_Control
	case 20008:
		return KnxDatapointSubtype_DPST_Direction1_Control
	case 20009:
		return KnxDatapointSubtype_DPST_Direction2_Control
	case 20010:
		return KnxDatapointSubtype_DPST_Start_Control
	case 200100:
		return KnxDatapointSubtype_DPST_FuelType
	case 200101:
		return KnxDatapointSubtype_DPST_BurnerType
	case 200102:
		return KnxDatapointSubtype_DPST_HVACMode
	case 200103:
		return KnxDatapointSubtype_DPST_DHWMode
	case 200104:
		return KnxDatapointSubtype_DPST_LoadPriority
	case 200105:
		return KnxDatapointSubtype_DPST_HVACContrMode
	case 200106:
		return KnxDatapointSubtype_DPST_HVACEmergMode
	case 200107:
		return KnxDatapointSubtype_DPST_ChangeoverMode
	case 200108:
		return KnxDatapointSubtype_DPST_ValveMode
	case 200109:
		return KnxDatapointSubtype_DPST_DamperMode
	case 20011:
		return KnxDatapointSubtype_DPST_State_Control
	case 200110:
		return KnxDatapointSubtype_DPST_HeaterMode
	case 200111:
		return KnxDatapointSubtype_DPST_FanMode
	case 200112:
		return KnxDatapointSubtype_DPST_MasterSlaveMode
	case 200113:
		return KnxDatapointSubtype_DPST_StatusRoomSetp
	case 200114:
		return KnxDatapointSubtype_DPST_Metering_DeviceType
	case 200115:
		return KnxDatapointSubtype_DPST_HumDehumMode
	case 200116:
		return KnxDatapointSubtype_DPST_EnableHCStage
	case 20012:
		return KnxDatapointSubtype_DPST_Invert_Control
	case 200120:
		return KnxDatapointSubtype_DPST_ADAType
	case 200121:
		return KnxDatapointSubtype_DPST_BackupMode
	case 200122:
		return KnxDatapointSubtype_DPST_StartSynchronization
	case 200600:
		return KnxDatapointSubtype_DPST_Behaviour_Lock_Unlock
	case 200601:
		return KnxDatapointSubtype_DPST_Behaviour_Bus_Power_Up_Down
	case 200602:
		return KnxDatapointSubtype_DPST_DALI_Fade_Time
	case 200603:
		return KnxDatapointSubtype_DPST_BlinkingMode
	case 200604:
		return KnxDatapointSubtype_DPST_LightControlMode
	case 200605:
		return KnxDatapointSubtype_DPST_SwitchPBModel
	case 200606:
		return KnxDatapointSubtype_DPST_PBAction
	case 200607:
		return KnxDatapointSubtype_DPST_DimmPBModel
	case 200608:
		return KnxDatapointSubtype_DPST_SwitchOnMode
	case 200609:
		return KnxDatapointSubtype_DPST_LoadTypeSet
	case 200610:
		return KnxDatapointSubtype_DPST_LoadTypeDetected
	case 200611:
		return KnxDatapointSubtype_DPST_Converter_Test_Control
	case 200801:
		return KnxDatapointSubtype_DPST_SABExcept_Behaviour
	case 200802:
		return KnxDatapointSubtype_DPST_SABBehaviour_Lock_Unlock
	case 200803:
		return KnxDatapointSubtype_DPST_SSSBMode
	case 200804:
		return KnxDatapointSubtype_DPST_BlindsControlMode
	case 201000:
		return KnxDatapointSubtype_DPST_CommMode
	case 201001:
		return KnxDatapointSubtype_DPST_AddInfoTypes
	case 201002:
		return KnxDatapointSubtype_DPST_RF_ModeSelect
	case 201003:
		return KnxDatapointSubtype_DPST_RF_FilterSelect
	case 2060100:
		return KnxDatapointSubtype_DPST_HVACModeNext
	case 2060102:
		return KnxDatapointSubtype_DPST_DHWModeNext
	case 2060104:
		return KnxDatapointSubtype_DPST_OccModeNext
	case 2060105:
		return KnxDatapointSubtype_DPST_BuildingModeNext
	case 210001:
		return KnxDatapointSubtype_DPST_StatusGen
	case 210002:
		return KnxDatapointSubtype_DPST_Device_Control
	case 210100:
		return KnxDatapointSubtype_DPST_ForceSign
	case 210101:
		return KnxDatapointSubtype_DPST_ForceSignCool
	case 210102:
		return KnxDatapointSubtype_DPST_StatusRHC
	case 210103:
		return KnxDatapointSubtype_DPST_StatusSDHWC
	case 210104:
		return KnxDatapointSubtype_DPST_FuelTypeSet
	case 210105:
		return KnxDatapointSubtype_DPST_StatusRCC
	case 210106:
		return KnxDatapointSubtype_DPST_StatusAHU
	case 210107:
		return KnxDatapointSubtype_DPST_CombinedStatus_RTSM
	case 210601:
		return KnxDatapointSubtype_DPST_LightActuatorErrorInfo
	case 211000:
		return KnxDatapointSubtype_DPST_RF_ModeInfo
	case 211001:
		return KnxDatapointSubtype_DPST_RF_FilterInfo
	case 211010:
		return KnxDatapointSubtype_DPST_Channel_Activation_8
	case 2170001:
		return KnxDatapointSubtype_DPST_Version
	case 2190001:
		return KnxDatapointSubtype_DPST_AlarmInfo
	case 220100:
		return KnxDatapointSubtype_DPST_StatusDHWC
	case 220101:
		return KnxDatapointSubtype_DPST_StatusRHCC
	case 220102:
		return KnxDatapointSubtype_DPST_CombinedStatus_HVA
	case 220103:
		return KnxDatapointSubtype_DPST_CombinedStatus_RTC
	case 221000:
		return KnxDatapointSubtype_DPST_Media
	case 221010:
		return KnxDatapointSubtype_DPST_Channel_Activation_16
	case 2220100:
		return KnxDatapointSubtype_DPST_TempRoomSetpSetF16_3
	case 2220101:
		return KnxDatapointSubtype_DPST_TempRoomSetpSetShiftF16_3
	case 2250001:
		return KnxDatapointSubtype_DPST_Scaling_Speed
	case 2250002:
		return KnxDatapointSubtype_DPST_Scaling_Step_Time
	case 2290001:
		return KnxDatapointSubtype_DPST_MeteringValue
	case 230001:
		return KnxDatapointSubtype_DPST_OnOffAction
	case 230002:
		return KnxDatapointSubtype_DPST_Alarm_Reaction
	case 230003:
		return KnxDatapointSubtype_DPST_UpDown_Action
	case 2301000:
		return KnxDatapointSubtype_DPST_MBus_Address
	case 230102:
		return KnxDatapointSubtype_DPST_HVAC_PB_Action
	case 2320600:
		return KnxDatapointSubtype_DPST_Colour_RGB
	case 2340001:
		return KnxDatapointSubtype_DPST_LanguageCodeAlpha2_ASCII
	case 2350001:
		return KnxDatapointSubtype_DPST_Tariff_ActiveEnergy
	case 2360001:
		return KnxDatapointSubtype_DPST_Prioritised_Mode_Control
	case 2370600:
		return KnxDatapointSubtype_DPST_DALI_Control_Gear_Diagnostic
	case 2380600:
		return KnxDatapointSubtype_DPST_DALI_Diagnostics
	case 2400800:
		return KnxDatapointSubtype_DPST_CombinedPosition
	case 2410800:
		return KnxDatapointSubtype_DPST_StatusSAB
	case 2420600:
		return KnxDatapointSubtype_DPST_Colour_xyY
	case 2440600:
		return KnxDatapointSubtype_DPST_Converter_Status
	case 2450600:
		return KnxDatapointSubtype_DPST_Converter_Test_Result
	case 2460600:
		return KnxDatapointSubtype_DPST_Battery_Info
	case 2490600:
		return KnxDatapointSubtype_DPST_Brightness_Colour_Temperature_Transition
	case 2500600:
		return KnxDatapointSubtype_DPST_Brightness_Colour_Temperature_Control
	case 251000:
		return KnxDatapointSubtype_DPST_DoubleNibble
	case 2510600:
		return KnxDatapointSubtype_DPST_Colour_RGBW
	case 2520600:
		return KnxDatapointSubtype_DPST_Relative_Control_RGBW
	case 2540600:
		return KnxDatapointSubtype_DPST_Relative_Control_RGB
	case 2550001:
		return KnxDatapointSubtype_DPST_GeographicalLocation
	case 260001:
		return KnxDatapointSubtype_DPST_SceneInfo
	case 270001:
		return KnxDatapointSubtype_DPST_CombinedInfoOnOff
	case 2750100:
		return KnxDatapointSubtype_DPST_TempRoomSetpSetF16_4
	case 2750101:
		return KnxDatapointSubtype_DPST_TempRoomSetpSetShiftF16_4
	case 290010:
		return KnxDatapointSubtype_DPST_ActiveEnergy_V64
	case 290011:
		return KnxDatapointSubtype_DPST_ApparantEnergy_V64
	case 290012:
		return KnxDatapointSubtype_DPST_ReactiveEnergy_V64
	case 30007:
		return KnxDatapointSubtype_DPST_Control_Dimming
	case 30008:
		return KnxDatapointSubtype_DPST_Control_Blinds
	case 301010:
		return KnxDatapointSubtype_DPST_Channel_Activation_24
	case 40001:
		return KnxDatapointSubtype_DPST_Char_ASCII
	case 40002:
		return KnxDatapointSubtype_DPST_Char_8859_1
	case 50001:
		return KnxDatapointSubtype_DPST_Scaling
	case 50003:
		return KnxDatapointSubtype_DPST_Angle
	case 50004:
		return KnxDatapointSubtype_DPST_Percent_U8
	case 50005:
		return KnxDatapointSubtype_DPST_DecimalFactor
	case 50006:
		return KnxDatapointSubtype_DPST_Tariff
	case 50010:
		return KnxDatapointSubtype_DPST_Value_1_Ucount
	case 50100:
		return KnxDatapointSubtype_DPST_FanStage
	case 60001:
		return KnxDatapointSubtype_DPST_Percent_V8
	case 60010:
		return KnxDatapointSubtype_DPST_Value_1_Count
	case 60020:
		return KnxDatapointSubtype_DPST_Status_Mode3
	case 70001:
		return KnxDatapointSubtype_DPST_Value_2_Ucount
	case 70002:
		return KnxDatapointSubtype_DPST_TimePeriodMsec
	case 70003:
		return KnxDatapointSubtype_DPST_TimePeriod10Msec
	case 70004:
		return KnxDatapointSubtype_DPST_TimePeriod100Msec
	case 70005:
		return KnxDatapointSubtype_DPST_TimePeriodSec
	case 70006:
		return KnxDatapointSubtype_DPST_TimePeriodMin
	case 70007:
		return KnxDatapointSubtype_DPST_TimePeriodHrs
	case 70010:
		return KnxDatapointSubtype_DPST_PropDataType
	case 70011:
		return KnxDatapointSubtype_DPST_Length_mm
	case 70012:
		return KnxDatapointSubtype_DPST_UElCurrentmA
	case 70013:
		return KnxDatapointSubtype_DPST_Brightness
	case 70600:
		return KnxDatapointSubtype_DPST_Absolute_Colour_Temperature
	case 80001:
		return KnxDatapointSubtype_DPST_Value_2_Count
	case 80002:
		return KnxDatapointSubtype_DPST_DeltaTimeMsec
	case 80003:
		return KnxDatapointSubtype_DPST_DeltaTime10Msec
	case 80004:
		return KnxDatapointSubtype_DPST_DeltaTime100Msec
	case 80005:
		return KnxDatapointSubtype_DPST_DeltaTimeSec
	case 80006:
		return KnxDatapointSubtype_DPST_DeltaTimeMin
	case 80007:
		return KnxDatapointSubtype_DPST_DeltaTimeHrs
	case 80010:
		return KnxDatapointSubtype_DPST_Percent_V16
	case 80011:
		return KnxDatapointSubtype_DPST_Rotation_Angle
	case 80012:
		return KnxDatapointSubtype_DPST_Length_m
	case 90001:
		return KnxDatapointSubtype_DPST_Value_Temp
	case 90002:
		return KnxDatapointSubtype_DPST_Value_Tempd
	case 90003:
		return KnxDatapointSubtype_DPST_Value_Tempa
	case 90004:
		return KnxDatapointSubtype_DPST_Value_Lux
	case 90005:
		return KnxDatapointSubtype_DPST_Value_Wsp
	case 90006:
		return KnxDatapointSubtype_DPST_Value_Pres
	case 90007:
		return KnxDatapointSubtype_DPST_Value_Humidity
	case 90008:
		return KnxDatapointSubtype_DPST_Value_AirQuality
	case 90009:
		return KnxDatapointSubtype_DPST_Value_AirFlow
	case 90010:
		return KnxDatapointSubtype_DPST_Value_Time1
	case 90011:
		return KnxDatapointSubtype_DPST_Value_Time2
	case 90020:
		return KnxDatapointSubtype_DPST_Value_Volt
	case 90021:
		return KnxDatapointSubtype_DPST_Value_Curr
	case 90022:
		return KnxDatapointSubtype_DPST_PowerDensity
	case 90023:
		return KnxDatapointSubtype_DPST_KelvinPerPercent
	case 90024:
		return KnxDatapointSubtype_DPST_Power
	case 90025:
		return KnxDatapointSubtype_DPST_Value_Volume_Flow
	case 90026:
		return KnxDatapointSubtype_DPST_Rain_Amount
	case 90027:
		return KnxDatapointSubtype_DPST_Value_Temp_F
	case 90028:
		return KnxDatapointSubtype_DPST_Value_Wsp_kmh
	case 90029:
		return KnxDatapointSubtype_DPST_Value_Absolute_Humidity
	case 90030:
		return KnxDatapointSubtype_DPST_Concentration_ygm3
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
	return KnxDatapointSubtypeValueOf(val), nil
}

func (e KnxDatapointSubtype) Serialize(io utils.WriteBuffer) error {
	err := io.WriteUint32(32, uint32(e))
	return err
}

func (e KnxDatapointSubtype) String() string {
	switch e {
	case KnxDatapointSubtype_DPST_TimeOfDay:
		return "DPST_TimeOfDay"
	case KnxDatapointSubtype_DPST_Switch:
		return "DPST_Switch"
	case KnxDatapointSubtype_DPST_Bool:
		return "DPST_Bool"
	case KnxDatapointSubtype_DPST_Enable:
		return "DPST_Enable"
	case KnxDatapointSubtype_DPST_Ramp:
		return "DPST_Ramp"
	case KnxDatapointSubtype_DPST_Alarm:
		return "DPST_Alarm"
	case KnxDatapointSubtype_DPST_BinaryValue:
		return "DPST_BinaryValue"
	case KnxDatapointSubtype_DPST_Step:
		return "DPST_Step"
	case KnxDatapointSubtype_DPST_UpDown:
		return "DPST_UpDown"
	case KnxDatapointSubtype_DPST_OpenClose:
		return "DPST_OpenClose"
	case KnxDatapointSubtype_DPST_Start:
		return "DPST_Start"
	case KnxDatapointSubtype_DPST_State:
		return "DPST_State"
	case KnxDatapointSubtype_DPST_Invert:
		return "DPST_Invert"
	case KnxDatapointSubtype_DPST_DimSendStyle:
		return "DPST_DimSendStyle"
	case KnxDatapointSubtype_DPST_InputSource:
		return "DPST_InputSource"
	case KnxDatapointSubtype_DPST_Reset:
		return "DPST_Reset"
	case KnxDatapointSubtype_DPST_Ack:
		return "DPST_Ack"
	case KnxDatapointSubtype_DPST_Trigger:
		return "DPST_Trigger"
	case KnxDatapointSubtype_DPST_Occupancy:
		return "DPST_Occupancy"
	case KnxDatapointSubtype_DPST_Window_Door:
		return "DPST_Window_Door"
	case KnxDatapointSubtype_DPST_LogicalFunction:
		return "DPST_LogicalFunction"
	case KnxDatapointSubtype_DPST_Scene_AB:
		return "DPST_Scene_AB"
	case KnxDatapointSubtype_DPST_ShutterBlinds_Mode:
		return "DPST_ShutterBlinds_Mode"
	case KnxDatapointSubtype_DPST_DayNight:
		return "DPST_DayNight"
	case KnxDatapointSubtype_DPST_Heat_Cool:
		return "DPST_Heat_Cool"
	case KnxDatapointSubtype_DPST_Date:
		return "DPST_Date"
	case KnxDatapointSubtype_DPST_Value_4_Ucount:
		return "DPST_Value_4_Ucount"
	case KnxDatapointSubtype_DPST_LongTimePeriod_Sec:
		return "DPST_LongTimePeriod_Sec"
	case KnxDatapointSubtype_DPST_LongTimePeriod_Min:
		return "DPST_LongTimePeriod_Min"
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
	case KnxDatapointSubtype_DPST_PSUMode:
		return "DPST_PSUMode"
	case KnxDatapointSubtype_DPST_Switch_Control:
		return "DPST_Switch_Control"
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
	case KnxDatapointSubtype_DPST_Bool_Control:
		return "DPST_Bool_Control"
	case KnxDatapointSubtype_DPST_ActuatorConnectType:
		return "DPST_ActuatorConnectType"
	case KnxDatapointSubtype_DPST_Cloud_Cover:
		return "DPST_Cloud_Cover"
	case KnxDatapointSubtype_DPST_PowerReturnMode:
		return "DPST_PowerReturnMode"
	case KnxDatapointSubtype_DPST_Enable_Control:
		return "DPST_Enable_Control"
	case KnxDatapointSubtype_DPST_Ramp_Control:
		return "DPST_Ramp_Control"
	case KnxDatapointSubtype_DPST_Alarm_Control:
		return "DPST_Alarm_Control"
	case KnxDatapointSubtype_DPST_BinaryValue_Control:
		return "DPST_BinaryValue_Control"
	case KnxDatapointSubtype_DPST_Step_Control:
		return "DPST_Step_Control"
	case KnxDatapointSubtype_DPST_Direction1_Control:
		return "DPST_Direction1_Control"
	case KnxDatapointSubtype_DPST_Direction2_Control:
		return "DPST_Direction2_Control"
	case KnxDatapointSubtype_DPST_Start_Control:
		return "DPST_Start_Control"
	case KnxDatapointSubtype_DPST_FuelType:
		return "DPST_FuelType"
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
	case KnxDatapointSubtype_DPST_State_Control:
		return "DPST_State_Control"
	case KnxDatapointSubtype_DPST_HeaterMode:
		return "DPST_HeaterMode"
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
	case KnxDatapointSubtype_DPST_Invert_Control:
		return "DPST_Invert_Control"
	case KnxDatapointSubtype_DPST_ADAType:
		return "DPST_ADAType"
	case KnxDatapointSubtype_DPST_BackupMode:
		return "DPST_BackupMode"
	case KnxDatapointSubtype_DPST_StartSynchronization:
		return "DPST_StartSynchronization"
	case KnxDatapointSubtype_DPST_Behaviour_Lock_Unlock:
		return "DPST_Behaviour_Lock_Unlock"
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
	case KnxDatapointSubtype_DPST_HVACModeNext:
		return "DPST_HVACModeNext"
	case KnxDatapointSubtype_DPST_DHWModeNext:
		return "DPST_DHWModeNext"
	case KnxDatapointSubtype_DPST_OccModeNext:
		return "DPST_OccModeNext"
	case KnxDatapointSubtype_DPST_BuildingModeNext:
		return "DPST_BuildingModeNext"
	case KnxDatapointSubtype_DPST_StatusGen:
		return "DPST_StatusGen"
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
	case KnxDatapointSubtype_DPST_RF_ModeInfo:
		return "DPST_RF_ModeInfo"
	case KnxDatapointSubtype_DPST_RF_FilterInfo:
		return "DPST_RF_FilterInfo"
	case KnxDatapointSubtype_DPST_Channel_Activation_8:
		return "DPST_Channel_Activation_8"
	case KnxDatapointSubtype_DPST_Version:
		return "DPST_Version"
	case KnxDatapointSubtype_DPST_AlarmInfo:
		return "DPST_AlarmInfo"
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
	case KnxDatapointSubtype_DPST_TempRoomSetpSetF16_3:
		return "DPST_TempRoomSetpSetF16_3"
	case KnxDatapointSubtype_DPST_TempRoomSetpSetShiftF16_3:
		return "DPST_TempRoomSetpSetShiftF16_3"
	case KnxDatapointSubtype_DPST_Scaling_Speed:
		return "DPST_Scaling_Speed"
	case KnxDatapointSubtype_DPST_Scaling_Step_Time:
		return "DPST_Scaling_Step_Time"
	case KnxDatapointSubtype_DPST_MeteringValue:
		return "DPST_MeteringValue"
	case KnxDatapointSubtype_DPST_OnOffAction:
		return "DPST_OnOffAction"
	case KnxDatapointSubtype_DPST_Alarm_Reaction:
		return "DPST_Alarm_Reaction"
	case KnxDatapointSubtype_DPST_UpDown_Action:
		return "DPST_UpDown_Action"
	case KnxDatapointSubtype_DPST_MBus_Address:
		return "DPST_MBus_Address"
	case KnxDatapointSubtype_DPST_HVAC_PB_Action:
		return "DPST_HVAC_PB_Action"
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
	case KnxDatapointSubtype_DPST_DoubleNibble:
		return "DPST_DoubleNibble"
	case KnxDatapointSubtype_DPST_Colour_RGBW:
		return "DPST_Colour_RGBW"
	case KnxDatapointSubtype_DPST_Relative_Control_RGBW:
		return "DPST_Relative_Control_RGBW"
	case KnxDatapointSubtype_DPST_Relative_Control_RGB:
		return "DPST_Relative_Control_RGB"
	case KnxDatapointSubtype_DPST_GeographicalLocation:
		return "DPST_GeographicalLocation"
	case KnxDatapointSubtype_DPST_SceneInfo:
		return "DPST_SceneInfo"
	case KnxDatapointSubtype_DPST_CombinedInfoOnOff:
		return "DPST_CombinedInfoOnOff"
	case KnxDatapointSubtype_DPST_TempRoomSetpSetF16_4:
		return "DPST_TempRoomSetpSetF16_4"
	case KnxDatapointSubtype_DPST_TempRoomSetpSetShiftF16_4:
		return "DPST_TempRoomSetpSetShiftF16_4"
	case KnxDatapointSubtype_DPST_ActiveEnergy_V64:
		return "DPST_ActiveEnergy_V64"
	case KnxDatapointSubtype_DPST_ApparantEnergy_V64:
		return "DPST_ApparantEnergy_V64"
	case KnxDatapointSubtype_DPST_ReactiveEnergy_V64:
		return "DPST_ReactiveEnergy_V64"
	case KnxDatapointSubtype_DPST_Control_Dimming:
		return "DPST_Control_Dimming"
	case KnxDatapointSubtype_DPST_Control_Blinds:
		return "DPST_Control_Blinds"
	case KnxDatapointSubtype_DPST_Channel_Activation_24:
		return "DPST_Channel_Activation_24"
	case KnxDatapointSubtype_DPST_Char_ASCII:
		return "DPST_Char_ASCII"
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
	}
	return ""
}
