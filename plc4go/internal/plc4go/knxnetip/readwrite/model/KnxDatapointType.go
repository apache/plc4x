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

type KnxDatapointType string

type IKnxDatapointType interface {
    FormatName() string
    MainNumber() uint16
    SubNumber() uint16
    Serialize(io utils.WriteBuffer) error
}

const(
    KnxDatapointType_DPT_Switch KnxDatapointType = "DPT_Switch"
    KnxDatapointType_DPT_Bool KnxDatapointType = "DPT_Bool"
    KnxDatapointType_DPT_Enable KnxDatapointType = "DPT_Enable"
    KnxDatapointType_DPT_Ramp KnxDatapointType = "DPT_Ramp"
    KnxDatapointType_DPT_Alarm KnxDatapointType = "DPT_Alarm"
    KnxDatapointType_DPT_BinaryValue KnxDatapointType = "DPT_BinaryValue"
    KnxDatapointType_DPT_Step KnxDatapointType = "DPT_Step"
    KnxDatapointType_DPT_UpDown KnxDatapointType = "DPT_UpDown"
    KnxDatapointType_DPT_OpenClose KnxDatapointType = "DPT_OpenClose"
    KnxDatapointType_DPT_Start KnxDatapointType = "DPT_Start"
    KnxDatapointType_DPT_State KnxDatapointType = "DPT_State"
    KnxDatapointType_DPT_Invert KnxDatapointType = "DPT_Invert"
    KnxDatapointType_DPT_DimSendStyle KnxDatapointType = "DPT_DimSendStyle"
    KnxDatapointType_DPT_InputSource KnxDatapointType = "DPT_InputSource"
    KnxDatapointType_DPT_Reset KnxDatapointType = "DPT_Reset"
    KnxDatapointType_DPT_Ack KnxDatapointType = "DPT_Ack"
    KnxDatapointType_DPT_Trigger KnxDatapointType = "DPT_Trigger"
    KnxDatapointType_DPT_Occupancy KnxDatapointType = "DPT_Occupancy"
    KnxDatapointType_DPT_Window_Door KnxDatapointType = "DPT_Window_Door"
    KnxDatapointType_DPT_LogicalFunction KnxDatapointType = "DPT_LogicalFunction"
    KnxDatapointType_DPT_Scene_AB KnxDatapointType = "DPT_Scene_AB"
    KnxDatapointType_DPT_ShutterBlinds_Mode KnxDatapointType = "DPT_ShutterBlinds_Mode"
    KnxDatapointType_DPT_Heat_Cool KnxDatapointType = "DPT_Heat_Cool"
    KnxDatapointType_DPT_Switch_Control KnxDatapointType = "DPT_Switch_Control"
    KnxDatapointType_DPT_Bool_Control KnxDatapointType = "DPT_Bool_Control"
    KnxDatapointType_DPT_Enable_Control KnxDatapointType = "DPT_Enable_Control"
    KnxDatapointType_DPT_Ramp_Control KnxDatapointType = "DPT_Ramp_Control"
    KnxDatapointType_DPT_Alarm_Control KnxDatapointType = "DPT_Alarm_Control"
    KnxDatapointType_DPT_BinaryValue_Control KnxDatapointType = "DPT_BinaryValue_Control"
    KnxDatapointType_DPT_Step_Control KnxDatapointType = "DPT_Step_Control"
    KnxDatapointType_DPT_Direction1_Control KnxDatapointType = "DPT_Direction1_Control"
    KnxDatapointType_DPT_Direction2_Control KnxDatapointType = "DPT_Direction2_Control"
    KnxDatapointType_DPT_Start_Control KnxDatapointType = "DPT_Start_Control"
    KnxDatapointType_DPT_State_Control KnxDatapointType = "DPT_State_Control"
    KnxDatapointType_DPT_Invert_Control KnxDatapointType = "DPT_Invert_Control"
    KnxDatapointType_DPT_Control_Dimming KnxDatapointType = "DPT_Control_Dimming"
    KnxDatapointType_DPT_Control_Blinds KnxDatapointType = "DPT_Control_Blinds"
    KnxDatapointType_DPT_Char_ASCII KnxDatapointType = "DPT_Char_ASCII"
    KnxDatapointType_DPT_Char_8859_1 KnxDatapointType = "DPT_Char_8859_1"
    KnxDatapointType_DPT_Scaling KnxDatapointType = "DPT_Scaling"
    KnxDatapointType_DPT_Angle KnxDatapointType = "DPT_Angle"
    KnxDatapointType_DPT_Percent_U8 KnxDatapointType = "DPT_Percent_U8"
    KnxDatapointType_DPT_DecimalFactor KnxDatapointType = "DPT_DecimalFactor"
    KnxDatapointType_DPT_Tariff KnxDatapointType = "DPT_Tariff"
    KnxDatapointType_DPT_Value_1_Ucount KnxDatapointType = "DPT_Value_1_Ucount"
    KnxDatapointType_DPT_Percent_V8 KnxDatapointType = "DPT_Percent_V8"
    KnxDatapointType_DPT_Value_1_Count KnxDatapointType = "DPT_Value_1_Count"
    KnxDatapointType_DPT_Status_Mode3 KnxDatapointType = "DPT_Status_Mode3"
    KnxDatapointType_DPT_Value_2_Ucount KnxDatapointType = "DPT_Value_2_Ucount"
    KnxDatapointType_DPT_TimePeriodMsec KnxDatapointType = "DPT_TimePeriodMsec"
    KnxDatapointType_DPT_TimePeriod10MSec KnxDatapointType = "DPT_TimePeriod10MSec"
    KnxDatapointType_DPT_TimePeriod100MSec KnxDatapointType = "DPT_TimePeriod100MSec"
    KnxDatapointType_DPT_TimePeriodSec KnxDatapointType = "DPT_TimePeriodSec"
    KnxDatapointType_DPT_TimePeriodMin KnxDatapointType = "DPT_TimePeriodMin"
    KnxDatapointType_DPT_TimePeriodHrs KnxDatapointType = "DPT_TimePeriodHrs"
    KnxDatapointType_DPT_PropDataType KnxDatapointType = "DPT_PropDataType"
    KnxDatapointType_DPT_Length_mm KnxDatapointType = "DPT_Length_mm"
    KnxDatapointType_DPT_UElCurrentmA KnxDatapointType = "DPT_UElCurrentmA"
    KnxDatapointType_DPT_Brightness KnxDatapointType = "DPT_Brightness"
    KnxDatapointType_DPT_Value_2_Count KnxDatapointType = "DPT_Value_2_Count"
    KnxDatapointType_DPT_DeltaTimeMsec KnxDatapointType = "DPT_DeltaTimeMsec"
    KnxDatapointType_DPT_DeltaTime10MSec KnxDatapointType = "DPT_DeltaTime10MSec"
    KnxDatapointType_DPT_DeltaTime100MSec KnxDatapointType = "DPT_DeltaTime100MSec"
    KnxDatapointType_DPT_DeltaTimeSec KnxDatapointType = "DPT_DeltaTimeSec"
    KnxDatapointType_DPT_DeltaTimeMin KnxDatapointType = "DPT_DeltaTimeMin"
    KnxDatapointType_DPT_DeltaTimeHrs KnxDatapointType = "DPT_DeltaTimeHrs"
    KnxDatapointType_DPT_Percent_V16 KnxDatapointType = "DPT_Percent_V16"
    KnxDatapointType_DPT_Rotation_Angle KnxDatapointType = "DPT_Rotation_Angle"
    KnxDatapointType_DPT_Value_Temp KnxDatapointType = "DPT_Value_Temp"
    KnxDatapointType_DPT_Value_Tempd KnxDatapointType = "DPT_Value_Tempd"
    KnxDatapointType_DPT_Value_Tempa KnxDatapointType = "DPT_Value_Tempa"
    KnxDatapointType_DPT_Value_Lux KnxDatapointType = "DPT_Value_Lux"
    KnxDatapointType_DPT_Value_Wsp KnxDatapointType = "DPT_Value_Wsp"
    KnxDatapointType_DPT_Value_Pres KnxDatapointType = "DPT_Value_Pres"
    KnxDatapointType_DPT_Value_Humidity KnxDatapointType = "DPT_Value_Humidity"
    KnxDatapointType_DPT_Value_AirQuality KnxDatapointType = "DPT_Value_AirQuality"
    KnxDatapointType_DPT_Value_Time1 KnxDatapointType = "DPT_Value_Time1"
    KnxDatapointType_DPT_Value_Time2 KnxDatapointType = "DPT_Value_Time2"
    KnxDatapointType_DPT_Value_Volt KnxDatapointType = "DPT_Value_Volt"
    KnxDatapointType_DPT_Value_Curr KnxDatapointType = "DPT_Value_Curr"
    KnxDatapointType_DPT_PowerDensity KnxDatapointType = "DPT_PowerDensity"
    KnxDatapointType_DPT_KelvinPerPercent KnxDatapointType = "DPT_KelvinPerPercent"
    KnxDatapointType_DPT_Power KnxDatapointType = "DPT_Power"
    KnxDatapointType_DPT_Value_Volume_Flow KnxDatapointType = "DPT_Value_Volume_Flow"
    KnxDatapointType_DPT_Rain_Amount KnxDatapointType = "DPT_Rain_Amount"
    KnxDatapointType_DPT_Value_Temp_F KnxDatapointType = "DPT_Value_Temp_F"
    KnxDatapointType_DPT_Value_Wsp_kmh KnxDatapointType = "DPT_Value_Wsp_kmh"
    KnxDatapointType_DPT_TimeOfDay KnxDatapointType = "DPT_TimeOfDay"
    KnxDatapointType_DPT_Date KnxDatapointType = "DPT_Date"
    KnxDatapointType_DPT_Value_4_Ucount KnxDatapointType = "DPT_Value_4_Ucount"
    KnxDatapointType_DPT_Value_4_Count KnxDatapointType = "DPT_Value_4_Count"
    KnxDatapointType_DPT_FlowRate_m3h KnxDatapointType = "DPT_FlowRate_m3h"
    KnxDatapointType_DPT_ActiveEnergy KnxDatapointType = "DPT_ActiveEnergy"
    KnxDatapointType_DPT_ApparantEnergy KnxDatapointType = "DPT_ApparantEnergy"
    KnxDatapointType_DPT_ReactiveEnergy KnxDatapointType = "DPT_ReactiveEnergy"
    KnxDatapointType_DPT_ActiveEnergy_kWh KnxDatapointType = "DPT_ActiveEnergy_kWh"
    KnxDatapointType_DPT_ApparantEnergy_kVAh KnxDatapointType = "DPT_ApparantEnergy_kVAh"
    KnxDatapointType_DPT_ReactiveEnergy_kVARh KnxDatapointType = "DPT_ReactiveEnergy_kVARh"
    KnxDatapointType_DPT_LongDeltaTimeSec KnxDatapointType = "DPT_LongDeltaTimeSec"
    KnxDatapointType_DPT_Value_Acceleration KnxDatapointType = "DPT_Value_Acceleration"
    KnxDatapointType_DPT_Value_Acceleration_Angular KnxDatapointType = "DPT_Value_Acceleration_Angular"
    KnxDatapointType_DPT_Value_Activation_Energy KnxDatapointType = "DPT_Value_Activation_Energy"
    KnxDatapointType_DPT_Value_Activity KnxDatapointType = "DPT_Value_Activity"
    KnxDatapointType_DPT_Value_Mol KnxDatapointType = "DPT_Value_Mol"
    KnxDatapointType_DPT_Value_Amplitude KnxDatapointType = "DPT_Value_Amplitude"
    KnxDatapointType_DPT_Value_AngleRad KnxDatapointType = "DPT_Value_AngleRad"
    KnxDatapointType_DPT_Value_AngleDeg KnxDatapointType = "DPT_Value_AngleDeg"
    KnxDatapointType_DPT_Value_Angular_Momentum KnxDatapointType = "DPT_Value_Angular_Momentum"
    KnxDatapointType_DPT_Value_Angular_Velocity KnxDatapointType = "DPT_Value_Angular_Velocity"
    KnxDatapointType_DPT_Value_Area KnxDatapointType = "DPT_Value_Area"
    KnxDatapointType_DPT_Value_Capacitance KnxDatapointType = "DPT_Value_Capacitance"
    KnxDatapointType_DPT_Value_Charge_DensitySurface KnxDatapointType = "DPT_Value_Charge_DensitySurface"
    KnxDatapointType_DPT_Value_Charge_DensityVolume KnxDatapointType = "DPT_Value_Charge_DensityVolume"
    KnxDatapointType_DPT_Value_Compressibility KnxDatapointType = "DPT_Value_Compressibility"
    KnxDatapointType_DPT_Value_Conductance KnxDatapointType = "DPT_Value_Conductance"
    KnxDatapointType_DPT_Value_Electrical_Conductivity KnxDatapointType = "DPT_Value_Electrical_Conductivity"
    KnxDatapointType_DPT_Value_Density KnxDatapointType = "DPT_Value_Density"
    KnxDatapointType_DPT_Value_Electric_Charge KnxDatapointType = "DPT_Value_Electric_Charge"
    KnxDatapointType_DPT_Value_Electric_Current KnxDatapointType = "DPT_Value_Electric_Current"
    KnxDatapointType_DPT_Value_Electric_CurrentDensity KnxDatapointType = "DPT_Value_Electric_CurrentDensity"
    KnxDatapointType_DPT_Value_Electric_DipoleMoment KnxDatapointType = "DPT_Value_Electric_DipoleMoment"
    KnxDatapointType_DPT_Value_Electric_Displacement KnxDatapointType = "DPT_Value_Electric_Displacement"
    KnxDatapointType_DPT_Value_Electric_FieldStrength KnxDatapointType = "DPT_Value_Electric_FieldStrength"
    KnxDatapointType_DPT_Value_Electric_Flux KnxDatapointType = "DPT_Value_Electric_Flux"
    KnxDatapointType_DPT_Value_Electric_FluxDensity KnxDatapointType = "DPT_Value_Electric_FluxDensity"
    KnxDatapointType_DPT_Value_Electric_Polarization KnxDatapointType = "DPT_Value_Electric_Polarization"
    KnxDatapointType_DPT_Value_Electric_Potential KnxDatapointType = "DPT_Value_Electric_Potential"
    KnxDatapointType_DPT_Value_Electric_PotentialDifference KnxDatapointType = "DPT_Value_Electric_PotentialDifference"
    KnxDatapointType_DPT_Value_ElectromagneticMoment KnxDatapointType = "DPT_Value_ElectromagneticMoment"
    KnxDatapointType_DPT_Value_Electromotive_Force KnxDatapointType = "DPT_Value_Electromotive_Force"
    KnxDatapointType_DPT_Value_Energy KnxDatapointType = "DPT_Value_Energy"
    KnxDatapointType_DPT_Value_Force KnxDatapointType = "DPT_Value_Force"
    KnxDatapointType_DPT_Value_Frequency KnxDatapointType = "DPT_Value_Frequency"
    KnxDatapointType_DPT_Value_Angular_Frequency KnxDatapointType = "DPT_Value_Angular_Frequency"
    KnxDatapointType_DPT_Value_Heat_Capacity KnxDatapointType = "DPT_Value_Heat_Capacity"
    KnxDatapointType_DPT_Value_Heat_FlowRate KnxDatapointType = "DPT_Value_Heat_FlowRate"
    KnxDatapointType_DPT_Value_Heat_Quantity KnxDatapointType = "DPT_Value_Heat_Quantity"
    KnxDatapointType_DPT_Value_Impedance KnxDatapointType = "DPT_Value_Impedance"
    KnxDatapointType_DPT_Value_Length KnxDatapointType = "DPT_Value_Length"
    KnxDatapointType_DPT_Value_Light_Quantity KnxDatapointType = "DPT_Value_Light_Quantity"
    KnxDatapointType_DPT_Value_Luminance KnxDatapointType = "DPT_Value_Luminance"
    KnxDatapointType_DPT_Value_Luminous_Flux KnxDatapointType = "DPT_Value_Luminous_Flux"
    KnxDatapointType_DPT_Value_Luminous_Intensity KnxDatapointType = "DPT_Value_Luminous_Intensity"
    KnxDatapointType_DPT_Value_Magnetic_FieldStrength KnxDatapointType = "DPT_Value_Magnetic_FieldStrength"
    KnxDatapointType_DPT_Value_Magnetic_Flux KnxDatapointType = "DPT_Value_Magnetic_Flux"
    KnxDatapointType_DPT_Value_Magnetic_FluxDensity KnxDatapointType = "DPT_Value_Magnetic_FluxDensity"
    KnxDatapointType_DPT_Value_Magnetic_Moment KnxDatapointType = "DPT_Value_Magnetic_Moment"
    KnxDatapointType_DPT_Value_Magnetic_Polarization KnxDatapointType = "DPT_Value_Magnetic_Polarization"
    KnxDatapointType_DPT_Value_Magnetization KnxDatapointType = "DPT_Value_Magnetization"
    KnxDatapointType_DPT_Value_MagnetomotiveForce KnxDatapointType = "DPT_Value_MagnetomotiveForce"
    KnxDatapointType_DPT_Value_Mass KnxDatapointType = "DPT_Value_Mass"
    KnxDatapointType_DPT_Value_MassFlux KnxDatapointType = "DPT_Value_MassFlux"
    KnxDatapointType_DPT_Value_Momentum KnxDatapointType = "DPT_Value_Momentum"
    KnxDatapointType_DPT_Value_Phase_AngleRad KnxDatapointType = "DPT_Value_Phase_AngleRad"
    KnxDatapointType_DPT_Value_Phase_AngleDeg KnxDatapointType = "DPT_Value_Phase_AngleDeg"
    KnxDatapointType_DPT_Value_Power KnxDatapointType = "DPT_Value_Power"
    KnxDatapointType_DPT_Value_Power_Factor KnxDatapointType = "DPT_Value_Power_Factor"
    KnxDatapointType_DPT_Value_Pressure KnxDatapointType = "DPT_Value_Pressure"
    KnxDatapointType_DPT_Value_Reactance KnxDatapointType = "DPT_Value_Reactance"
    KnxDatapointType_DPT_Value_Resistance KnxDatapointType = "DPT_Value_Resistance"
    KnxDatapointType_DPT_Value_Resistivity KnxDatapointType = "DPT_Value_Resistivity"
    KnxDatapointType_DPT_Value_SelfInductance KnxDatapointType = "DPT_Value_SelfInductance"
    KnxDatapointType_DPT_Value_SolidAngle KnxDatapointType = "DPT_Value_SolidAngle"
    KnxDatapointType_DPT_Value_Sound_Intensity KnxDatapointType = "DPT_Value_Sound_Intensity"
    KnxDatapointType_DPT_Value_Speed KnxDatapointType = "DPT_Value_Speed"
    KnxDatapointType_DPT_Value_Stress KnxDatapointType = "DPT_Value_Stress"
    KnxDatapointType_DPT_Value_Surface_Tension KnxDatapointType = "DPT_Value_Surface_Tension"
    KnxDatapointType_DPT_Value_Common_Temperature KnxDatapointType = "DPT_Value_Common_Temperature"
    KnxDatapointType_DPT_Value_Absolute_Temperature KnxDatapointType = "DPT_Value_Absolute_Temperature"
    KnxDatapointType_DPT_Value_TemperatureDifference KnxDatapointType = "DPT_Value_TemperatureDifference"
    KnxDatapointType_DPT_Value_Thermal_Capacity KnxDatapointType = "DPT_Value_Thermal_Capacity"
    KnxDatapointType_DPT_Value_Thermal_Conductivity KnxDatapointType = "DPT_Value_Thermal_Conductivity"
    KnxDatapointType_DPT_Value_ThermoelectricPower KnxDatapointType = "DPT_Value_ThermoelectricPower"
    KnxDatapointType_DPT_Value_Time KnxDatapointType = "DPT_Value_Time"
    KnxDatapointType_DPT_Value_Torque KnxDatapointType = "DPT_Value_Torque"
    KnxDatapointType_DPT_Value_Volume KnxDatapointType = "DPT_Value_Volume"
    KnxDatapointType_DPT_Value_Volume_Flux KnxDatapointType = "DPT_Value_Volume_Flux"
    KnxDatapointType_DPT_Value_Weight KnxDatapointType = "DPT_Value_Weight"
    KnxDatapointType_DPT_Value_Work KnxDatapointType = "DPT_Value_Work"
    KnxDatapointType_DPT_Access_Data KnxDatapointType = "DPT_Access_Data"
    KnxDatapointType_DPT_String_ASCII KnxDatapointType = "DPT_String_ASCII"
    KnxDatapointType_DPT_String_8859_1 KnxDatapointType = "DPT_String_8859_1"
    KnxDatapointType_DPT_SceneNumber KnxDatapointType = "DPT_SceneNumber"
    KnxDatapointType_DPT_SceneControl KnxDatapointType = "DPT_SceneControl"
    KnxDatapointType_DPT_DateTime KnxDatapointType = "DPT_DateTime"
    KnxDatapointType_DPT_SCLOMode KnxDatapointType = "DPT_SCLOMode"
    KnxDatapointType_DPT_BuildingMode KnxDatapointType = "DPT_BuildingMode"
    KnxDatapointType_DPT_OccMode KnxDatapointType = "DPT_OccMode"
    KnxDatapointType_DPT_Priority KnxDatapointType = "DPT_Priority"
    KnxDatapointType_DPT_LightApplicationMode KnxDatapointType = "DPT_LightApplicationMode"
    KnxDatapointType_DPT_ApplicationArea KnxDatapointType = "DPT_ApplicationArea"
    KnxDatapointType_DPT_AlarmClassType KnxDatapointType = "DPT_AlarmClassType"
    KnxDatapointType_DPT_PSUMode KnxDatapointType = "DPT_PSUMode"
    KnxDatapointType_DPT_ErrorClass_System KnxDatapointType = "DPT_ErrorClass_System"
    KnxDatapointType_DPT_ErrorClass_HVAC KnxDatapointType = "DPT_ErrorClass_HVAC"
    KnxDatapointType_DPT_Time_Delay KnxDatapointType = "DPT_Time_Delay"
    KnxDatapointType_DPT_Beaufort_Wind_Force_Scale KnxDatapointType = "DPT_Beaufort_Wind_Force_Scale"
    KnxDatapointType_DPT_SensorSelect KnxDatapointType = "DPT_SensorSelect"
    KnxDatapointType_DPT_ActuatorConnectType KnxDatapointType = "DPT_ActuatorConnectType"
    KnxDatapointType_DPT_FuelType KnxDatapointType = "DPT_FuelType"
    KnxDatapointType_DPT_BurnerType KnxDatapointType = "DPT_BurnerType"
    KnxDatapointType_DPT_HVACMode KnxDatapointType = "DPT_HVACMode"
    KnxDatapointType_DPT_DHWMode KnxDatapointType = "DPT_DHWMode"
    KnxDatapointType_DPT_LoadPriority KnxDatapointType = "DPT_LoadPriority"
    KnxDatapointType_DPT_HVACContrMode KnxDatapointType = "DPT_HVACContrMode"
    KnxDatapointType_DPT_HVACEmergMode KnxDatapointType = "DPT_HVACEmergMode"
    KnxDatapointType_DPT_ChangeoverMode KnxDatapointType = "DPT_ChangeoverMode"
    KnxDatapointType_DPT_ValveMode KnxDatapointType = "DPT_ValveMode"
    KnxDatapointType_DPT_DamperMode KnxDatapointType = "DPT_DamperMode"
    KnxDatapointType_DPT_HeaterMode KnxDatapointType = "DPT_HeaterMode"
    KnxDatapointType_DPT_FanMode KnxDatapointType = "DPT_FanMode"
    KnxDatapointType_DPT_MasterSlaveMode KnxDatapointType = "DPT_MasterSlaveMode"
    KnxDatapointType_DPT_StatusRoomSetp KnxDatapointType = "DPT_StatusRoomSetp"
    KnxDatapointType_DPT_ADAType KnxDatapointType = "DPT_ADAType"
    KnxDatapointType_DPT_BackupMode KnxDatapointType = "DPT_BackupMode"
    KnxDatapointType_DPT_StartSynchronization KnxDatapointType = "DPT_StartSynchronization"
    KnxDatapointType_DPT_Behaviour_Lock_Unlock KnxDatapointType = "DPT_Behaviour_Lock_Unlock"
    KnxDatapointType_DPT_Behaviour_Bus_Power_Up_Down KnxDatapointType = "DPT_Behaviour_Bus_Power_Up_Down"
    KnxDatapointType_DPT_DALI_Fade_Time KnxDatapointType = "DPT_DALI_Fade_Time"
    KnxDatapointType_DPT_BlinkingMode KnxDatapointType = "DPT_BlinkingMode"
    KnxDatapointType_DPT_LightControlMode KnxDatapointType = "DPT_LightControlMode"
    KnxDatapointType_DPT_SwitchPBModel KnxDatapointType = "DPT_SwitchPBModel"
    KnxDatapointType_DPT_PBAction KnxDatapointType = "DPT_PBAction"
    KnxDatapointType_DPT_DimmPBModel KnxDatapointType = "DPT_DimmPBModel"
    KnxDatapointType_DPT_SwitchOnMode KnxDatapointType = "DPT_SwitchOnMode"
    KnxDatapointType_DPT_LoadTypeSet KnxDatapointType = "DPT_LoadTypeSet"
    KnxDatapointType_DPT_LoadTypeDetected KnxDatapointType = "DPT_LoadTypeDetected"
    KnxDatapointType_DPT_SABExceptBehaviour KnxDatapointType = "DPT_SABExceptBehaviour"
    KnxDatapointType_DPT_SABBehaviour_Lock_Unlock KnxDatapointType = "DPT_SABBehaviour_Lock_Unlock"
    KnxDatapointType_DPT_SSSBMode KnxDatapointType = "DPT_SSSBMode"
    KnxDatapointType_DPT_BlindsControlMode KnxDatapointType = "DPT_BlindsControlMode"
    KnxDatapointType_DPT_CommMode KnxDatapointType = "DPT_CommMode"
    KnxDatapointType_DPT_AddInfoTypes KnxDatapointType = "DPT_AddInfoTypes"
    KnxDatapointType_DPT_RF_ModeSelect KnxDatapointType = "DPT_RF_ModeSelect"
    KnxDatapointType_DPT_RF_FilterSelect KnxDatapointType = "DPT_RF_FilterSelect"
    KnxDatapointType_DPT_StatusGen KnxDatapointType = "DPT_StatusGen"
    KnxDatapointType_DPT_Device_Control KnxDatapointType = "DPT_Device_Control"
    KnxDatapointType_DPT_ForceSign KnxDatapointType = "DPT_ForceSign"
    KnxDatapointType_DPT_ForceSignCool KnxDatapointType = "DPT_ForceSignCool"
    KnxDatapointType_DPT_StatusRHC KnxDatapointType = "DPT_StatusRHC"
    KnxDatapointType_DPT_StatusSDHWC KnxDatapointType = "DPT_StatusSDHWC"
    KnxDatapointType_DPT_FuelTypeSet KnxDatapointType = "DPT_FuelTypeSet"
    KnxDatapointType_DPT_StatusRCC KnxDatapointType = "DPT_StatusRCC"
    KnxDatapointType_DPT_StatusAHU KnxDatapointType = "DPT_StatusAHU"
    KnxDatapointType_DPT_LightActuatorErrorInfo KnxDatapointType = "DPT_LightActuatorErrorInfo"
    KnxDatapointType_DPT_RF_ModeInfo KnxDatapointType = "DPT_RF_ModeInfo"
    KnxDatapointType_DPT_RF_FilterInfo KnxDatapointType = "DPT_RF_FilterInfo"
    KnxDatapointType_DPT_Channel_Activation_8 KnxDatapointType = "DPT_Channel_Activation_8"
    KnxDatapointType_DPT_StatusDHWC KnxDatapointType = "DPT_StatusDHWC"
    KnxDatapointType_DPT_StatusRHCC KnxDatapointType = "DPT_StatusRHCC"
    KnxDatapointType_DPT_Media KnxDatapointType = "DPT_Media"
    KnxDatapointType_DPT_Channel_Activation_16 KnxDatapointType = "DPT_Channel_Activation_16"
    KnxDatapointType_DPT_OnOff_Action KnxDatapointType = "DPT_OnOff_Action"
    KnxDatapointType_DPT_Alarm_Reaction KnxDatapointType = "DPT_Alarm_Reaction"
    KnxDatapointType_DPT_UpDown_Action KnxDatapointType = "DPT_UpDown_Action"
    KnxDatapointType_DPT_HVAC_PB_Action KnxDatapointType = "DPT_HVAC_PB_Action"
    KnxDatapointType_DPT_VarString_8859_1 KnxDatapointType = "DPT_VarString_8859_1"
    KnxDatapointType_DPT_DoubleNibble KnxDatapointType = "DPT_DoubleNibble"
    KnxDatapointType_DPT_SceneInfo KnxDatapointType = "DPT_SceneInfo"
    KnxDatapointType_DPT_CombinedInfoOnOff KnxDatapointType = "DPT_CombinedInfoOnOff"
    KnxDatapointType_DPT_UTF_8 KnxDatapointType = "DPT_UTF_8"
    KnxDatapointType_DPT_ActiveEnergy_V64 KnxDatapointType = "DPT_ActiveEnergy_V64"
    KnxDatapointType_DPT_ApparantEnergy_V64 KnxDatapointType = "DPT_ApparantEnergy_V64"
    KnxDatapointType_DPT_ReactiveEnergy_V64 KnxDatapointType = "DPT_ReactiveEnergy_V64"
    KnxDatapointType_DPT_Channel_Activation_24 KnxDatapointType = "DPT_Channel_Activation_24"
    KnxDatapointType_DPT_PB_Action_HVAC_Extended KnxDatapointType = "DPT_PB_Action_HVAC_Extended"
    KnxDatapointType_DPT_Heat_Cool_Z KnxDatapointType = "DPT_Heat_Cool_Z"
    KnxDatapointType_DPT_BinaryValue_Z KnxDatapointType = "DPT_BinaryValue_Z"
    KnxDatapointType_DPT_HVACMode_Z KnxDatapointType = "DPT_HVACMode_Z"
    KnxDatapointType_DPT_DHWMode_Z KnxDatapointType = "DPT_DHWMode_Z"
    KnxDatapointType_DPT_HVACContrMode_Z KnxDatapointType = "DPT_HVACContrMode_Z"
    KnxDatapointType_DPT_EnablH_Cstage_Z_DPT_EnablH_CStage KnxDatapointType = "DPT_EnablH_Cstage_Z_DPT_EnablH_CStage"
    KnxDatapointType_DPT_BuildingMode_Z KnxDatapointType = "DPT_BuildingMode_Z"
    KnxDatapointType_DPT_OccMode_Z KnxDatapointType = "DPT_OccMode_Z"
    KnxDatapointType_DPT_HVACEmergMode_Z KnxDatapointType = "DPT_HVACEmergMode_Z"
    KnxDatapointType_DPT_RelValue_Z KnxDatapointType = "DPT_RelValue_Z"
    KnxDatapointType_DPT_UCountValue8_Z KnxDatapointType = "DPT_UCountValue8_Z"
    KnxDatapointType_DPT_TimePeriodMsec_Z KnxDatapointType = "DPT_TimePeriodMsec_Z"
    KnxDatapointType_DPT_TimePeriod10Msec_Z KnxDatapointType = "DPT_TimePeriod10Msec_Z"
    KnxDatapointType_DPT_TimePeriod100Msec_Z KnxDatapointType = "DPT_TimePeriod100Msec_Z"
    KnxDatapointType_DPT_TimePeriodSec_Z KnxDatapointType = "DPT_TimePeriodSec_Z"
    KnxDatapointType_DPT_TimePeriodMin_Z KnxDatapointType = "DPT_TimePeriodMin_Z"
    KnxDatapointType_DPT_TimePeriodHrs_Z KnxDatapointType = "DPT_TimePeriodHrs_Z"
    KnxDatapointType_DPT_UFlowRateLiter_h_Z KnxDatapointType = "DPT_UFlowRateLiter_h_Z"
    KnxDatapointType_DPT_UCountValue16_Z KnxDatapointType = "DPT_UCountValue16_Z"
    KnxDatapointType_DPT_UElCurrentyA_Z KnxDatapointType = "DPT_UElCurrentyA_Z"
    KnxDatapointType_DPT_PowerKW_Z KnxDatapointType = "DPT_PowerKW_Z"
    KnxDatapointType_DPT_AtmPressureAbs_Z KnxDatapointType = "DPT_AtmPressureAbs_Z"
    KnxDatapointType_DPT_PercentU16_Z KnxDatapointType = "DPT_PercentU16_Z"
    KnxDatapointType_DPT_HVACAirQual_Z KnxDatapointType = "DPT_HVACAirQual_Z"
    KnxDatapointType_DPT_WindSpeed_Z_DPT_WindSpeed KnxDatapointType = "DPT_WindSpeed_Z_DPT_WindSpeed"
    KnxDatapointType_DPT_SunIntensity_Z KnxDatapointType = "DPT_SunIntensity_Z"
    KnxDatapointType_DPT_HVACAirFlowAbs_Z KnxDatapointType = "DPT_HVACAirFlowAbs_Z"
    KnxDatapointType_DPT_RelSignedValue_Z KnxDatapointType = "DPT_RelSignedValue_Z"
    KnxDatapointType_DPT_DeltaTimeMsec_Z KnxDatapointType = "DPT_DeltaTimeMsec_Z"
    KnxDatapointType_DPT_DeltaTime10Msec_Z KnxDatapointType = "DPT_DeltaTime10Msec_Z"
    KnxDatapointType_DPT_DeltaTime100Msec_Z KnxDatapointType = "DPT_DeltaTime100Msec_Z"
    KnxDatapointType_DPT_DeltaTimeSec_Z KnxDatapointType = "DPT_DeltaTimeSec_Z"
    KnxDatapointType_DPT_DeltaTimeMin_Z KnxDatapointType = "DPT_DeltaTimeMin_Z"
    KnxDatapointType_DPT_DeltaTimeHrs_Z KnxDatapointType = "DPT_DeltaTimeHrs_Z"
    KnxDatapointType_DPT_Percent_V16_Z KnxDatapointType = "DPT_Percent_V16_Z"
    KnxDatapointType_DPT_TempHVACAbs_Z KnxDatapointType = "DPT_TempHVACAbs_Z"
    KnxDatapointType_DPT_TempHVACRel_Z KnxDatapointType = "DPT_TempHVACRel_Z"
    KnxDatapointType_DPT_HVACAirFlowRel_Z KnxDatapointType = "DPT_HVACAirFlowRel_Z"
    KnxDatapointType_DPT_HVACModeNext KnxDatapointType = "DPT_HVACModeNext"
    KnxDatapointType_DPT_DHWModeNext KnxDatapointType = "DPT_DHWModeNext"
    KnxDatapointType_DPT_OccModeNext KnxDatapointType = "DPT_OccModeNext"
    KnxDatapointType_DPT_BuildingModeNext KnxDatapointType = "DPT_BuildingModeNext"
    KnxDatapointType_DPT_StatusBUC KnxDatapointType = "DPT_StatusBUC"
    KnxDatapointType_DPT_LockSign KnxDatapointType = "DPT_LockSign"
    KnxDatapointType_DPT_ValueDemBOC KnxDatapointType = "DPT_ValueDemBOC"
    KnxDatapointType_DPT_ActPosDemAbs KnxDatapointType = "DPT_ActPosDemAbs"
    KnxDatapointType_DPT_StatusAct KnxDatapointType = "DPT_StatusAct"
    KnxDatapointType_DPT_StatusLightingActuator KnxDatapointType = "DPT_StatusLightingActuator"
    KnxDatapointType_DPT_StatusHPM KnxDatapointType = "DPT_StatusHPM"
    KnxDatapointType_DPT_TempRoomDemAbs KnxDatapointType = "DPT_TempRoomDemAbs"
    KnxDatapointType_DPT_StatusCPM KnxDatapointType = "DPT_StatusCPM"
    KnxDatapointType_DPT_StatusWTC KnxDatapointType = "DPT_StatusWTC"
    KnxDatapointType_DPT_TempFlowWaterDemAbs KnxDatapointType = "DPT_TempFlowWaterDemAbs"
    KnxDatapointType_DPT_EnergyDemWater KnxDatapointType = "DPT_EnergyDemWater"
    KnxDatapointType_DPT_TempRoomSetpSetShift3 KnxDatapointType = "DPT_TempRoomSetpSetShift3"
    KnxDatapointType_DPT_TempRoomSetpSet3 KnxDatapointType = "DPT_TempRoomSetpSet3"
    KnxDatapointType_DPT_TempRoomSetpSet4 KnxDatapointType = "DPT_TempRoomSetpSet4"
    KnxDatapointType_DPT_TempDHWSetpSet4 KnxDatapointType = "DPT_TempDHWSetpSet4"
    KnxDatapointType_DPT_TempRoomSetpSetShift4 KnxDatapointType = "DPT_TempRoomSetpSetShift4"
    KnxDatapointType_DPT_PowerFlowWaterDemHPM KnxDatapointType = "DPT_PowerFlowWaterDemHPM"
    KnxDatapointType_DPT_PowerFlowWaterDemCPM KnxDatapointType = "DPT_PowerFlowWaterDemCPM"
    KnxDatapointType_DPT_StatusBOC KnxDatapointType = "DPT_StatusBOC"
    KnxDatapointType_DPT_StatusCC KnxDatapointType = "DPT_StatusCC"
    KnxDatapointType_DPT_SpecHeatProd KnxDatapointType = "DPT_SpecHeatProd"
    KnxDatapointType_DPT_Version KnxDatapointType = "DPT_Version"
    KnxDatapointType_DPT_VolumeLiter_Z KnxDatapointType = "DPT_VolumeLiter_Z"
    KnxDatapointType_DPT_FlowRate_m3h_Z KnxDatapointType = "DPT_FlowRate_m3h_Z"
    KnxDatapointType_DPT_AlarmInfo KnxDatapointType = "DPT_AlarmInfo"
    KnxDatapointType_DPT_TempHVACAbsNext KnxDatapointType = "DPT_TempHVACAbsNext"
    KnxDatapointType_DPT_SerNum KnxDatapointType = "DPT_SerNum"
    KnxDatapointType_DPT_TempRoomSetpSetF163 KnxDatapointType = "DPT_TempRoomSetpSetF163"
    KnxDatapointType_DPT_TempRoomSetpSetShiftF163 KnxDatapointType = "DPT_TempRoomSetpSetShiftF163"
    KnxDatapointType_DPT_EnergyDemAir KnxDatapointType = "DPT_EnergyDemAir"
    KnxDatapointType_DPT_TempSupply_AirSetpSet KnxDatapointType = "DPT_TempSupply_AirSetpSet"
    KnxDatapointType_DPT_ScalingSpeed KnxDatapointType = "DPT_ScalingSpeed"
    KnxDatapointType_DPT_Scaling_Step_Time KnxDatapointType = "DPT_Scaling_Step_Time"
    KnxDatapointType_DPT_TariffNext KnxDatapointType = "DPT_TariffNext"
    KnxDatapointType_DPT_MeteringValue KnxDatapointType = "DPT_MeteringValue"
    KnxDatapointType_DPT_MBus_Address KnxDatapointType = "DPT_MBus_Address"
    KnxDatapointType_DPT_Locale_ASCII KnxDatapointType = "DPT_Locale_ASCII"
    KnxDatapointType_DPT_Colour_RGB KnxDatapointType = "DPT_Colour_RGB"
    KnxDatapointType_DPT_LanguageCodeAlpha2_ASCII KnxDatapointType = "DPT_LanguageCodeAlpha2_ASCII"
    KnxDatapointType_DPT_RegionCodeAlpha2_ASCII KnxDatapointType = "DPT_RegionCodeAlpha2_ASCII"
    KnxDatapointType_DPT_Tariff_ActiveEnergy KnxDatapointType = "DPT_Tariff_ActiveEnergy"
    KnxDatapointType_DPT_Prioritised_Mode_Control KnxDatapointType = "DPT_Prioritised_Mode_Control"
    KnxDatapointType_DPT_DALI_Control_Gear_Diagnostic KnxDatapointType = "DPT_DALI_Control_Gear_Diagnostic"
    KnxDatapointType_DPT_SceneConfig KnxDatapointType = "DPT_SceneConfig"
    KnxDatapointType_DPT_DALI_Diagnostics KnxDatapointType = "DPT_DALI_Diagnostics"
    KnxDatapointType_DPT_FlaggedScaling KnxDatapointType = "DPT_FlaggedScaling"
    KnxDatapointType_DPT_CombinedPosition KnxDatapointType = "DPT_CombinedPosition"
    KnxDatapointType_DPT_StatusSAB KnxDatapointType = "DPT_StatusSAB"
)


func (e KnxDatapointType) FormatName() string {
    switch e  {
        case "DPT_ADAType": { /* 'DPT_ADAType' */
            return "N8"
        }
        case "DPT_Access_Data": { /* 'DPT_Access_Data' */
            return "U4U4U4U4U4U4B4N4"
        }
        case "DPT_Ack": { /* 'DPT_Ack' */
            return "B1"
        }
        case "DPT_ActPosDemAbs": { /* 'DPT_ActPosDemAbs' */
            return "U8B8ActuatorPositionDemand"
        }
        case "DPT_ActiveEnergy": { /* 'DPT_ActiveEnergy' */
            return "V32"
        }
        case "DPT_ActiveEnergy_V64": { /* 'DPT_ActiveEnergy_V64' */
            return "V64"
        }
        case "DPT_ActiveEnergy_kWh": { /* 'DPT_ActiveEnergy_kWh' */
            return "V32"
        }
        case "DPT_ActuatorConnectType": { /* 'DPT_ActuatorConnectType' */
            return "N8"
        }
        case "DPT_AddInfoTypes": { /* 'DPT_AddInfoTypes' */
            return "N8"
        }
        case "DPT_Alarm": { /* 'DPT_Alarm' */
            return "B1"
        }
        case "DPT_AlarmClassType": { /* 'DPT_AlarmClassType' */
            return "N8"
        }
        case "DPT_AlarmInfo": { /* 'DPT_AlarmInfo' */
            return "U8N8N8N8B8B8"
        }
        case "DPT_Alarm_Control": { /* 'DPT_Alarm_Control' */
            return "B2"
        }
        case "DPT_Alarm_Reaction": { /* 'DPT_Alarm_Reaction' */
            return "N2"
        }
        case "DPT_Angle": { /* 'DPT_Angle' */
            return "U8"
        }
        case "DPT_ApparantEnergy": { /* 'DPT_ApparantEnergy' */
            return "V32"
        }
        case "DPT_ApparantEnergy_V64": { /* 'DPT_ApparantEnergy_V64' */
            return "V64"
        }
        case "DPT_ApparantEnergy_kVAh": { /* 'DPT_ApparantEnergy_kVAh' */
            return "V32"
        }
        case "DPT_ApplicationArea": { /* 'DPT_ApplicationArea' */
            return "N8"
        }
        case "DPT_AtmPressureAbs_Z": { /* 'DPT_AtmPressureAbs_Z' */
            return "U16Z8AtmPressure"
        }
        case "DPT_BackupMode": { /* 'DPT_BackupMode' */
            return "N8"
        }
        case "DPT_Beaufort_Wind_Force_Scale": { /* 'DPT_Beaufort_Wind_Force_Scale' */
            return "N8"
        }
        case "DPT_Behaviour_Bus_Power_Up_Down": { /* 'DPT_Behaviour_Bus_Power_Up_Down' */
            return "N8"
        }
        case "DPT_Behaviour_Lock_Unlock": { /* 'DPT_Behaviour_Lock_Unlock' */
            return "N8"
        }
        case "DPT_BinaryValue": { /* 'DPT_BinaryValue' */
            return "B1"
        }
        case "DPT_BinaryValue_Control": { /* 'DPT_BinaryValue_Control' */
            return "B2"
        }
        case "DPT_BinaryValue_Z": { /* 'DPT_BinaryValue_Z' */
            return "B1Z8BinaryValueZ"
        }
        case "DPT_BlindsControlMode": { /* 'DPT_BlindsControlMode' */
            return "N8"
        }
        case "DPT_BlinkingMode": { /* 'DPT_BlinkingMode' */
            return "N8"
        }
        case "DPT_Bool": { /* 'DPT_Bool' */
            return "B1"
        }
        case "DPT_Bool_Control": { /* 'DPT_Bool_Control' */
            return "B2"
        }
        case "DPT_Brightness": { /* 'DPT_Brightness' */
            return "U16"
        }
        case "DPT_BuildingMode": { /* 'DPT_BuildingMode' */
            return "N8"
        }
        case "DPT_BuildingModeNext": { /* 'DPT_BuildingModeNext' */
            return "U16N8BuildingModeAndTimeDelay"
        }
        case "DPT_BuildingMode_Z": { /* 'DPT_BuildingMode_Z' */
            return "N8Z8BuildingMode"
        }
        case "DPT_BurnerType": { /* 'DPT_BurnerType' */
            return "N8"
        }
        case "DPT_ChangeoverMode": { /* 'DPT_ChangeoverMode' */
            return "N8"
        }
        case "DPT_Channel_Activation_16": { /* 'DPT_Channel_Activation_16' */
            return "B16"
        }
        case "DPT_Channel_Activation_24": { /* 'DPT_Channel_Activation_24' */
            return "B24"
        }
        case "DPT_Channel_Activation_8": { /* 'DPT_Channel_Activation_8' */
            return "B8"
        }
        case "DPT_Char_8859_1": { /* 'DPT_Char_8859_1' */
            return "A8_8859_1"
        }
        case "DPT_Char_ASCII": { /* 'DPT_Char_ASCII' */
            return "A8_ASCII"
        }
        case "DPT_Colour_RGB": { /* 'DPT_Colour_RGB' */
            return "U8U8U8"
        }
        case "DPT_CombinedInfoOnOff": { /* 'DPT_CombinedInfoOnOff' */
            return "B32"
        }
        case "DPT_CombinedPosition": { /* 'DPT_CombinedPosition' */
            return "U8U8B8"
        }
        case "DPT_CommMode": { /* 'DPT_CommMode' */
            return "N8"
        }
        case "DPT_Control_Blinds": { /* 'DPT_Control_Blinds' */
            return "B1U3"
        }
        case "DPT_Control_Dimming": { /* 'DPT_Control_Dimming' */
            return "B1U3"
        }
        case "DPT_DALI_Control_Gear_Diagnostic": { /* 'DPT_DALI_Control_Gear_Diagnostic' */
            return "B10U6"
        }
        case "DPT_DALI_Diagnostics": { /* 'DPT_DALI_Diagnostics' */
            return "B2U6"
        }
        case "DPT_DALI_Fade_Time": { /* 'DPT_DALI_Fade_Time' */
            return "N8"
        }
        case "DPT_DHWMode": { /* 'DPT_DHWMode' */
            return "N8"
        }
        case "DPT_DHWModeNext": { /* 'DPT_DHWModeNext' */
            return "U16N8DhwModeAndTimeDelay"
        }
        case "DPT_DHWMode_Z": { /* 'DPT_DHWMode_Z' */
            return "N8Z8DhwMode"
        }
        case "DPT_DamperMode": { /* 'DPT_DamperMode' */
            return "N8"
        }
        case "DPT_Date": { /* 'DPT_Date' */
            return "r3N5r4N4r1U7"
        }
        case "DPT_DateTime": { /* 'DPT_DateTime' */
            return "U8r4U4r3U5U3U5r2U6r2U6B16"
        }
        case "DPT_DecimalFactor": { /* 'DPT_DecimalFactor' */
            return "U8"
        }
        case "DPT_DeltaTime100MSec": { /* 'DPT_DeltaTime100MSec' */
            return "V16"
        }
        case "DPT_DeltaTime100Msec_Z": { /* 'DPT_DeltaTime100Msec_Z' */
            return "V16Z8DeltaTime"
        }
        case "DPT_DeltaTime10MSec": { /* 'DPT_DeltaTime10MSec' */
            return "V16"
        }
        case "DPT_DeltaTime10Msec_Z": { /* 'DPT_DeltaTime10Msec_Z' */
            return "V16Z8DeltaTime"
        }
        case "DPT_DeltaTimeHrs": { /* 'DPT_DeltaTimeHrs' */
            return "V16"
        }
        case "DPT_DeltaTimeHrs_Z": { /* 'DPT_DeltaTimeHrs_Z' */
            return "V16Z8DeltaTime"
        }
        case "DPT_DeltaTimeMin": { /* 'DPT_DeltaTimeMin' */
            return "V16"
        }
        case "DPT_DeltaTimeMin_Z": { /* 'DPT_DeltaTimeMin_Z' */
            return "V16Z8DeltaTime"
        }
        case "DPT_DeltaTimeMsec": { /* 'DPT_DeltaTimeMsec' */
            return "V16"
        }
        case "DPT_DeltaTimeMsec_Z": { /* 'DPT_DeltaTimeMsec_Z' */
            return "V16Z8DeltaTime"
        }
        case "DPT_DeltaTimeSec": { /* 'DPT_DeltaTimeSec' */
            return "V16"
        }
        case "DPT_DeltaTimeSec_Z": { /* 'DPT_DeltaTimeSec_Z' */
            return "V16Z8DeltaTime"
        }
        case "DPT_Device_Control": { /* 'DPT_Device_Control' */
            return "B8"
        }
        case "DPT_DimSendStyle": { /* 'DPT_DimSendStyle' */
            return "B1"
        }
        case "DPT_DimmPBModel": { /* 'DPT_DimmPBModel' */
            return "N8"
        }
        case "DPT_Direction1_Control": { /* 'DPT_Direction1_Control' */
            return "B2"
        }
        case "DPT_Direction2_Control": { /* 'DPT_Direction2_Control' */
            return "B2"
        }
        case "DPT_DoubleNibble": { /* 'DPT_DoubleNibble' */
            return "U4U4"
        }
        case "DPT_EnablH_Cstage_Z_DPT_EnablH_CStage": { /* 'DPT_EnablH_Cstage_Z_DPT_EnablH_CStage' */
            return "N8Z8EnableHeatingOrCoolingStage"
        }
        case "DPT_Enable": { /* 'DPT_Enable' */
            return "B1"
        }
        case "DPT_Enable_Control": { /* 'DPT_Enable_Control' */
            return "B2"
        }
        case "DPT_EnergyDemAir": { /* 'DPT_EnergyDemAir' */
            return "V8N8N8"
        }
        case "DPT_EnergyDemWater": { /* 'DPT_EnergyDemWater' */
            return "U8N8"
        }
        case "DPT_ErrorClass_HVAC": { /* 'DPT_ErrorClass_HVAC' */
            return "N8"
        }
        case "DPT_ErrorClass_System": { /* 'DPT_ErrorClass_System' */
            return "N8"
        }
        case "DPT_FanMode": { /* 'DPT_FanMode' */
            return "N8"
        }
        case "DPT_FlaggedScaling": { /* 'DPT_FlaggedScaling' */
            return "U8r7B1"
        }
        case "DPT_FlowRate_m3h": { /* 'DPT_FlowRate_m3h' */
            return "V32"
        }
        case "DPT_FlowRate_m3h_Z": { /* 'DPT_FlowRate_m3h_Z' */
            return "V32Z8FlowRate"
        }
        case "DPT_ForceSign": { /* 'DPT_ForceSign' */
            return "B8"
        }
        case "DPT_ForceSignCool": { /* 'DPT_ForceSignCool' */
            return "B8"
        }
        case "DPT_FuelType": { /* 'DPT_FuelType' */
            return "N8"
        }
        case "DPT_FuelTypeSet": { /* 'DPT_FuelTypeSet' */
            return "B8"
        }
        case "DPT_HVACAirFlowAbs_Z": { /* 'DPT_HVACAirFlowAbs_Z' */
            return "U16Z8HvacAirFlow"
        }
        case "DPT_HVACAirFlowRel_Z": { /* 'DPT_HVACAirFlowRel_Z' */
            return "V16Z8RelSignedValue"
        }
        case "DPT_HVACAirQual_Z": { /* 'DPT_HVACAirQual_Z' */
            return "U16Z8HvacAirQuality"
        }
        case "DPT_HVACContrMode": { /* 'DPT_HVACContrMode' */
            return "N8"
        }
        case "DPT_HVACContrMode_Z": { /* 'DPT_HVACContrMode_Z' */
            return "N8Z8HvacControllingMode"
        }
        case "DPT_HVACEmergMode": { /* 'DPT_HVACEmergMode' */
            return "N8"
        }
        case "DPT_HVACEmergMode_Z": { /* 'DPT_HVACEmergMode_Z' */
            return "N8Z8EmergencyMode"
        }
        case "DPT_HVACMode": { /* 'DPT_HVACMode' */
            return "N8"
        }
        case "DPT_HVACModeNext": { /* 'DPT_HVACModeNext' */
            return "U16N8HvacModeAndTimeDelay"
        }
        case "DPT_HVACMode_Z": { /* 'DPT_HVACMode_Z' */
            return "N8Z8HvacOperatingMode"
        }
        case "DPT_HVAC_PB_Action": { /* 'DPT_HVAC_PB_Action' */
            return "N2"
        }
        case "DPT_Heat_Cool": { /* 'DPT_Heat_Cool' */
            return "B1"
        }
        case "DPT_Heat_Cool_Z": { /* 'DPT_Heat_Cool_Z' */
            return "B1Z8HeatingOrCoolingZ"
        }
        case "DPT_HeaterMode": { /* 'DPT_HeaterMode' */
            return "N8"
        }
        case "DPT_InputSource": { /* 'DPT_InputSource' */
            return "B1"
        }
        case "DPT_Invert": { /* 'DPT_Invert' */
            return "B1"
        }
        case "DPT_Invert_Control": { /* 'DPT_Invert_Control' */
            return "B2"
        }
        case "DPT_KelvinPerPercent": { /* 'DPT_KelvinPerPercent' */
            return "F16"
        }
        case "DPT_LanguageCodeAlpha2_ASCII": { /* 'DPT_LanguageCodeAlpha2_ASCII' */
            return "A8A8Language"
        }
        case "DPT_Length_mm": { /* 'DPT_Length_mm' */
            return "U16"
        }
        case "DPT_LightActuatorErrorInfo": { /* 'DPT_LightActuatorErrorInfo' */
            return "B8"
        }
        case "DPT_LightApplicationMode": { /* 'DPT_LightApplicationMode' */
            return "N8"
        }
        case "DPT_LightControlMode": { /* 'DPT_LightControlMode' */
            return "N8"
        }
        case "DPT_LoadPriority": { /* 'DPT_LoadPriority' */
            return "N8"
        }
        case "DPT_LoadTypeDetected": { /* 'DPT_LoadTypeDetected' */
            return "N8"
        }
        case "DPT_LoadTypeSet": { /* 'DPT_LoadTypeSet' */
            return "N8"
        }
        case "DPT_Locale_ASCII": { /* 'DPT_Locale_ASCII' */
            return "A8A8A8A8"
        }
        case "DPT_LockSign": { /* 'DPT_LockSign' */
            return "U8B8LockingSignal"
        }
        case "DPT_LogicalFunction": { /* 'DPT_LogicalFunction' */
            return "B1"
        }
        case "DPT_LongDeltaTimeSec": { /* 'DPT_LongDeltaTimeSec' */
            return "V32"
        }
        case "DPT_MBus_Address": { /* 'DPT_MBus_Address' */
            return "U16U32U8N8"
        }
        case "DPT_MasterSlaveMode": { /* 'DPT_MasterSlaveMode' */
            return "N8"
        }
        case "DPT_Media": { /* 'DPT_Media' */
            return "B16"
        }
        case "DPT_MeteringValue": { /* 'DPT_MeteringValue' */
            return "V32N8Z8"
        }
        case "DPT_OccMode": { /* 'DPT_OccMode' */
            return "N8"
        }
        case "DPT_OccModeNext": { /* 'DPT_OccModeNext' */
            return "U16N8OccupancyModeAndTimeDelay"
        }
        case "DPT_OccMode_Z": { /* 'DPT_OccMode_Z' */
            return "N8Z8OccupancyMode"
        }
        case "DPT_Occupancy": { /* 'DPT_Occupancy' */
            return "B1"
        }
        case "DPT_OnOff_Action": { /* 'DPT_OnOff_Action' */
            return "N2"
        }
        case "DPT_OpenClose": { /* 'DPT_OpenClose' */
            return "B1"
        }
        case "DPT_PBAction": { /* 'DPT_PBAction' */
            return "N8"
        }
        case "DPT_PB_Action_HVAC_Extended": { /* 'DPT_PB_Action_HVAC_Extended' */
            return "N3"
        }
        case "DPT_PSUMode": { /* 'DPT_PSUMode' */
            return "N8"
        }
        case "DPT_PercentU16_Z": { /* 'DPT_PercentU16_Z' */
            return "U16Z8PercentValue"
        }
        case "DPT_Percent_U8": { /* 'DPT_Percent_U8' */
            return "U8"
        }
        case "DPT_Percent_V16": { /* 'DPT_Percent_V16' */
            return "V16"
        }
        case "DPT_Percent_V16_Z": { /* 'DPT_Percent_V16_Z' */
            return "V16Z8RelSignedValue"
        }
        case "DPT_Percent_V8": { /* 'DPT_Percent_V8' */
            return "V8"
        }
        case "DPT_Power": { /* 'DPT_Power' */
            return "F16"
        }
        case "DPT_PowerDensity": { /* 'DPT_PowerDensity' */
            return "F16"
        }
        case "DPT_PowerFlowWaterDemCPM": { /* 'DPT_PowerFlowWaterDemCPM' */
            return "V16U8B8ChilledWater"
        }
        case "DPT_PowerFlowWaterDemHPM": { /* 'DPT_PowerFlowWaterDemHPM' */
            return "V16U8B8Heat"
        }
        case "DPT_PowerKW_Z": { /* 'DPT_PowerKW_Z' */
            return "U16Z8Power"
        }
        case "DPT_Prioritised_Mode_Control": { /* 'DPT_Prioritised_Mode_Control' */
            return "B1N3N4"
        }
        case "DPT_Priority": { /* 'DPT_Priority' */
            return "N8"
        }
        case "DPT_PropDataType": { /* 'DPT_PropDataType' */
            return "U16"
        }
        case "DPT_RF_FilterInfo": { /* 'DPT_RF_FilterInfo' */
            return "B8"
        }
        case "DPT_RF_FilterSelect": { /* 'DPT_RF_FilterSelect' */
            return "N8"
        }
        case "DPT_RF_ModeInfo": { /* 'DPT_RF_ModeInfo' */
            return "B8"
        }
        case "DPT_RF_ModeSelect": { /* 'DPT_RF_ModeSelect' */
            return "N8"
        }
        case "DPT_Rain_Amount": { /* 'DPT_Rain_Amount' */
            return "F16"
        }
        case "DPT_Ramp": { /* 'DPT_Ramp' */
            return "B1"
        }
        case "DPT_Ramp_Control": { /* 'DPT_Ramp_Control' */
            return "B2"
        }
        case "DPT_ReactiveEnergy": { /* 'DPT_ReactiveEnergy' */
            return "V32"
        }
        case "DPT_ReactiveEnergy_V64": { /* 'DPT_ReactiveEnergy_V64' */
            return "V64"
        }
        case "DPT_ReactiveEnergy_kVARh": { /* 'DPT_ReactiveEnergy_kVARh' */
            return "V32"
        }
        case "DPT_RegionCodeAlpha2_ASCII": { /* 'DPT_RegionCodeAlpha2_ASCII' */
            return "A8A8Region"
        }
        case "DPT_RelSignedValue_Z": { /* 'DPT_RelSignedValue_Z' */
            return "V8Z8RelSignedValue"
        }
        case "DPT_RelValue_Z": { /* 'DPT_RelValue_Z' */
            return "U8Z8Rel"
        }
        case "DPT_Reset": { /* 'DPT_Reset' */
            return "B1"
        }
        case "DPT_Rotation_Angle": { /* 'DPT_Rotation_Angle' */
            return "V16"
        }
        case "DPT_SABBehaviour_Lock_Unlock": { /* 'DPT_SABBehaviour_Lock_Unlock' */
            return "N8"
        }
        case "DPT_SABExceptBehaviour": { /* 'DPT_SABExceptBehaviour' */
            return "N8"
        }
        case "DPT_SCLOMode": { /* 'DPT_SCLOMode' */
            return "N8"
        }
        case "DPT_SSSBMode": { /* 'DPT_SSSBMode' */
            return "N8"
        }
        case "DPT_Scaling": { /* 'DPT_Scaling' */
            return "U8"
        }
        case "DPT_ScalingSpeed": { /* 'DPT_ScalingSpeed' */
            return "U16U8Scaling"
        }
        case "DPT_Scaling_Step_Time": { /* 'DPT_Scaling_Step_Time' */
            return "U16U8Scaling"
        }
        case "DPT_SceneConfig": { /* 'DPT_SceneConfig' */
            return "B2U6"
        }
        case "DPT_SceneControl": { /* 'DPT_SceneControl' */
            return "B1r1U6"
        }
        case "DPT_SceneInfo": { /* 'DPT_SceneInfo' */
            return "r1b1U6"
        }
        case "DPT_SceneNumber": { /* 'DPT_SceneNumber' */
            return "r2U6"
        }
        case "DPT_Scene_AB": { /* 'DPT_Scene_AB' */
            return "B1"
        }
        case "DPT_SensorSelect": { /* 'DPT_SensorSelect' */
            return "N8"
        }
        case "DPT_SerNum": { /* 'DPT_SerNum' */
            return "N16U32"
        }
        case "DPT_ShutterBlinds_Mode": { /* 'DPT_ShutterBlinds_Mode' */
            return "B1"
        }
        case "DPT_SpecHeatProd": { /* 'DPT_SpecHeatProd' */
            return "U16U8N8B8"
        }
        case "DPT_Start": { /* 'DPT_Start' */
            return "B1"
        }
        case "DPT_StartSynchronization": { /* 'DPT_StartSynchronization' */
            return "N8"
        }
        case "DPT_Start_Control": { /* 'DPT_Start_Control' */
            return "B2"
        }
        case "DPT_State": { /* 'DPT_State' */
            return "B1"
        }
        case "DPT_State_Control": { /* 'DPT_State_Control' */
            return "B2"
        }
        case "DPT_StatusAHU": { /* 'DPT_StatusAHU' */
            return "B8"
        }
        case "DPT_StatusAct": { /* 'DPT_StatusAct' */
            return "U8B8ActuatorPositionStatus"
        }
        case "DPT_StatusBOC": { /* 'DPT_StatusBOC' */
            return "V16U8B16Boiler"
        }
        case "DPT_StatusBUC": { /* 'DPT_StatusBUC' */
            return "U8B8StatusBurnerController"
        }
        case "DPT_StatusCC": { /* 'DPT_StatusCC' */
            return "V16U8B16Chiller"
        }
        case "DPT_StatusCPM": { /* 'DPT_StatusCPM' */
            return "V16B8ColdWaterProducerManagerStatus"
        }
        case "DPT_StatusDHWC": { /* 'DPT_StatusDHWC' */
            return "B16"
        }
        case "DPT_StatusGen": { /* 'DPT_StatusGen' */
            return "B8"
        }
        case "DPT_StatusHPM": { /* 'DPT_StatusHPM' */
            return "V16B8HeatProducerManagerStatus"
        }
        case "DPT_StatusLightingActuator": { /* 'DPT_StatusLightingActuator' */
            return "U8B8StatusLightingActuator"
        }
        case "DPT_StatusRCC": { /* 'DPT_StatusRCC' */
            return "B8"
        }
        case "DPT_StatusRHC": { /* 'DPT_StatusRHC' */
            return "B8"
        }
        case "DPT_StatusRHCC": { /* 'DPT_StatusRHCC' */
            return "B16"
        }
        case "DPT_StatusRoomSetp": { /* 'DPT_StatusRoomSetp' */
            return "N8"
        }
        case "DPT_StatusSAB": { /* 'DPT_StatusSAB' */
            return "U8U8B16"
        }
        case "DPT_StatusSDHWC": { /* 'DPT_StatusSDHWC' */
            return "B8"
        }
        case "DPT_StatusWTC": { /* 'DPT_StatusWTC' */
            return "V16B8WaterTemperatureControllerStatus"
        }
        case "DPT_Status_Mode3": { /* 'DPT_Status_Mode3' */
            return "B5N3"
        }
        case "DPT_Step": { /* 'DPT_Step' */
            return "B1"
        }
        case "DPT_Step_Control": { /* 'DPT_Step_Control' */
            return "B2"
        }
        case "DPT_String_8859_1": { /* 'DPT_String_8859_1' */
            return "A112_8859_1"
        }
        case "DPT_String_ASCII": { /* 'DPT_String_ASCII' */
            return "A112_ASCII"
        }
        case "DPT_SunIntensity_Z": { /* 'DPT_SunIntensity_Z' */
            return "U16Z8SunIntensity"
        }
        case "DPT_Switch": { /* 'DPT_Switch' */
            return "B1"
        }
        case "DPT_SwitchOnMode": { /* 'DPT_SwitchOnMode' */
            return "N8"
        }
        case "DPT_SwitchPBModel": { /* 'DPT_SwitchPBModel' */
            return "N8"
        }
        case "DPT_Switch_Control": { /* 'DPT_Switch_Control' */
            return "B2"
        }
        case "DPT_Tariff": { /* 'DPT_Tariff' */
            return "U8"
        }
        case "DPT_TariffNext": { /* 'DPT_TariffNext' */
            return "U16U8TariffNext"
        }
        case "DPT_Tariff_ActiveEnergy": { /* 'DPT_Tariff_ActiveEnergy' */
            return "V32U8B8"
        }
        case "DPT_TempDHWSetpSet4": { /* 'DPT_TempDHWSetpSet4' */
            return "V16V16V16V16DhwtTemperature"
        }
        case "DPT_TempFlowWaterDemAbs": { /* 'DPT_TempFlowWaterDemAbs' */
            return "V16B16"
        }
        case "DPT_TempHVACAbsNext": { /* 'DPT_TempHVACAbsNext' */
            return "U16V16"
        }
        case "DPT_TempHVACAbs_Z": { /* 'DPT_TempHVACAbs_Z' */
            return "V16Z8RelSignedValue"
        }
        case "DPT_TempHVACRel_Z": { /* 'DPT_TempHVACRel_Z' */
            return "V16Z8RelSignedValue"
        }
        case "DPT_TempRoomDemAbs": { /* 'DPT_TempRoomDemAbs' */
            return "V16B8RoomTemperatureDemand"
        }
        case "DPT_TempRoomSetpSet3": { /* 'DPT_TempRoomSetpSet3' */
            return "V16V16V16RoomTemperature"
        }
        case "DPT_TempRoomSetpSet4": { /* 'DPT_TempRoomSetpSet4' */
            return "V16V16V16V16RoomTemperature"
        }
        case "DPT_TempRoomSetpSetF163": { /* 'DPT_TempRoomSetpSetF163' */
            return "F16F16F16"
        }
        case "DPT_TempRoomSetpSetShift3": { /* 'DPT_TempRoomSetpSetShift3' */
            return "V16V16V16RoomTemperatureShift"
        }
        case "DPT_TempRoomSetpSetShift4": { /* 'DPT_TempRoomSetpSetShift4' */
            return "V16V16V16V16RoomTemperatureShift"
        }
        case "DPT_TempRoomSetpSetShiftF163": { /* 'DPT_TempRoomSetpSetShiftF163' */
            return "F16F16F16"
        }
        case "DPT_TempSupply_AirSetpSet": { /* 'DPT_TempSupply_AirSetpSet' */
            return "V16V16N8N8"
        }
        case "DPT_TimeOfDay": { /* 'DPT_TimeOfDay' */
            return "N3N5r2N6r2N6"
        }
        case "DPT_TimePeriod100MSec": { /* 'DPT_TimePeriod100MSec' */
            return "U16"
        }
        case "DPT_TimePeriod100Msec_Z": { /* 'DPT_TimePeriod100Msec_Z' */
            return "U16Z8TimePeriod"
        }
        case "DPT_TimePeriod10MSec": { /* 'DPT_TimePeriod10MSec' */
            return "U16"
        }
        case "DPT_TimePeriod10Msec_Z": { /* 'DPT_TimePeriod10Msec_Z' */
            return "U16Z8TimePeriod"
        }
        case "DPT_TimePeriodHrs": { /* 'DPT_TimePeriodHrs' */
            return "U16"
        }
        case "DPT_TimePeriodHrs_Z": { /* 'DPT_TimePeriodHrs_Z' */
            return "U16Z8TimePeriod"
        }
        case "DPT_TimePeriodMin": { /* 'DPT_TimePeriodMin' */
            return "U16"
        }
        case "DPT_TimePeriodMin_Z": { /* 'DPT_TimePeriodMin_Z' */
            return "U16Z8TimePeriod"
        }
        case "DPT_TimePeriodMsec": { /* 'DPT_TimePeriodMsec' */
            return "U16"
        }
        case "DPT_TimePeriodMsec_Z": { /* 'DPT_TimePeriodMsec_Z' */
            return "U16Z8TimePeriod"
        }
        case "DPT_TimePeriodSec": { /* 'DPT_TimePeriodSec' */
            return "U16"
        }
        case "DPT_TimePeriodSec_Z": { /* 'DPT_TimePeriodSec_Z' */
            return "U16Z8TimePeriod"
        }
        case "DPT_Time_Delay": { /* 'DPT_Time_Delay' */
            return "N8"
        }
        case "DPT_Trigger": { /* 'DPT_Trigger' */
            return "B1"
        }
        case "DPT_UCountValue16_Z": { /* 'DPT_UCountValue16_Z' */
            return "U16Z8Counter"
        }
        case "DPT_UCountValue8_Z": { /* 'DPT_UCountValue8_Z' */
            return "U8Z8Counter"
        }
        case "DPT_UElCurrentmA": { /* 'DPT_UElCurrentmA' */
            return "U16"
        }
        case "DPT_UElCurrentyA_Z": { /* 'DPT_UElCurrentyA_Z' */
            return "U16Z8ElectricCurrent"
        }
        case "DPT_UFlowRateLiter_h_Z": { /* 'DPT_UFlowRateLiter_h_Z' */
            return "U16Z8FlowRate"
        }
        case "DPT_UTF_8": { /* 'DPT_UTF_8' */
            return "An_UTF_8"
        }
        case "DPT_UpDown": { /* 'DPT_UpDown' */
            return "B1"
        }
        case "DPT_UpDown_Action": { /* 'DPT_UpDown_Action' */
            return "N2"
        }
        case "DPT_ValueDemBOC": { /* 'DPT_ValueDemBOC' */
            return "U8B8BoilerControllerDemandSignal"
        }
        case "DPT_Value_1_Count": { /* 'DPT_Value_1_Count' */
            return "V8"
        }
        case "DPT_Value_1_Ucount": { /* 'DPT_Value_1_Ucount' */
            return "U8"
        }
        case "DPT_Value_2_Count": { /* 'DPT_Value_2_Count' */
            return "V16"
        }
        case "DPT_Value_2_Ucount": { /* 'DPT_Value_2_Ucount' */
            return "U16"
        }
        case "DPT_Value_4_Count": { /* 'DPT_Value_4_Count' */
            return "V32"
        }
        case "DPT_Value_4_Ucount": { /* 'DPT_Value_4_Ucount' */
            return "U32"
        }
        case "DPT_Value_Absolute_Temperature": { /* 'DPT_Value_Absolute_Temperature' */
            return "F32"
        }
        case "DPT_Value_Acceleration": { /* 'DPT_Value_Acceleration' */
            return "F32"
        }
        case "DPT_Value_Acceleration_Angular": { /* 'DPT_Value_Acceleration_Angular' */
            return "F32"
        }
        case "DPT_Value_Activation_Energy": { /* 'DPT_Value_Activation_Energy' */
            return "F32"
        }
        case "DPT_Value_Activity": { /* 'DPT_Value_Activity' */
            return "F32"
        }
        case "DPT_Value_AirQuality": { /* 'DPT_Value_AirQuality' */
            return "F16"
        }
        case "DPT_Value_Amplitude": { /* 'DPT_Value_Amplitude' */
            return "F32"
        }
        case "DPT_Value_AngleDeg": { /* 'DPT_Value_AngleDeg' */
            return "F32"
        }
        case "DPT_Value_AngleRad": { /* 'DPT_Value_AngleRad' */
            return "F32"
        }
        case "DPT_Value_Angular_Frequency": { /* 'DPT_Value_Angular_Frequency' */
            return "F32"
        }
        case "DPT_Value_Angular_Momentum": { /* 'DPT_Value_Angular_Momentum' */
            return "F32"
        }
        case "DPT_Value_Angular_Velocity": { /* 'DPT_Value_Angular_Velocity' */
            return "F32"
        }
        case "DPT_Value_Area": { /* 'DPT_Value_Area' */
            return "F32"
        }
        case "DPT_Value_Capacitance": { /* 'DPT_Value_Capacitance' */
            return "F32"
        }
        case "DPT_Value_Charge_DensitySurface": { /* 'DPT_Value_Charge_DensitySurface' */
            return "F32"
        }
        case "DPT_Value_Charge_DensityVolume": { /* 'DPT_Value_Charge_DensityVolume' */
            return "F32"
        }
        case "DPT_Value_Common_Temperature": { /* 'DPT_Value_Common_Temperature' */
            return "F32"
        }
        case "DPT_Value_Compressibility": { /* 'DPT_Value_Compressibility' */
            return "F32"
        }
        case "DPT_Value_Conductance": { /* 'DPT_Value_Conductance' */
            return "F32"
        }
        case "DPT_Value_Curr": { /* 'DPT_Value_Curr' */
            return "F16"
        }
        case "DPT_Value_Density": { /* 'DPT_Value_Density' */
            return "F32"
        }
        case "DPT_Value_Electric_Charge": { /* 'DPT_Value_Electric_Charge' */
            return "F32"
        }
        case "DPT_Value_Electric_Current": { /* 'DPT_Value_Electric_Current' */
            return "F32"
        }
        case "DPT_Value_Electric_CurrentDensity": { /* 'DPT_Value_Electric_CurrentDensity' */
            return "F32"
        }
        case "DPT_Value_Electric_DipoleMoment": { /* 'DPT_Value_Electric_DipoleMoment' */
            return "F32"
        }
        case "DPT_Value_Electric_Displacement": { /* 'DPT_Value_Electric_Displacement' */
            return "F32"
        }
        case "DPT_Value_Electric_FieldStrength": { /* 'DPT_Value_Electric_FieldStrength' */
            return "F32"
        }
        case "DPT_Value_Electric_Flux": { /* 'DPT_Value_Electric_Flux' */
            return "F32"
        }
        case "DPT_Value_Electric_FluxDensity": { /* 'DPT_Value_Electric_FluxDensity' */
            return "F32"
        }
        case "DPT_Value_Electric_Polarization": { /* 'DPT_Value_Electric_Polarization' */
            return "F32"
        }
        case "DPT_Value_Electric_Potential": { /* 'DPT_Value_Electric_Potential' */
            return "F32"
        }
        case "DPT_Value_Electric_PotentialDifference": { /* 'DPT_Value_Electric_PotentialDifference' */
            return "F32"
        }
        case "DPT_Value_Electrical_Conductivity": { /* 'DPT_Value_Electrical_Conductivity' */
            return "F32"
        }
        case "DPT_Value_ElectromagneticMoment": { /* 'DPT_Value_ElectromagneticMoment' */
            return "F32"
        }
        case "DPT_Value_Electromotive_Force": { /* 'DPT_Value_Electromotive_Force' */
            return "F32"
        }
        case "DPT_Value_Energy": { /* 'DPT_Value_Energy' */
            return "F32"
        }
        case "DPT_Value_Force": { /* 'DPT_Value_Force' */
            return "F32"
        }
        case "DPT_Value_Frequency": { /* 'DPT_Value_Frequency' */
            return "F32"
        }
        case "DPT_Value_Heat_Capacity": { /* 'DPT_Value_Heat_Capacity' */
            return "F32"
        }
        case "DPT_Value_Heat_FlowRate": { /* 'DPT_Value_Heat_FlowRate' */
            return "F32"
        }
        case "DPT_Value_Heat_Quantity": { /* 'DPT_Value_Heat_Quantity' */
            return "F32"
        }
        case "DPT_Value_Humidity": { /* 'DPT_Value_Humidity' */
            return "F16"
        }
        case "DPT_Value_Impedance": { /* 'DPT_Value_Impedance' */
            return "F32"
        }
        case "DPT_Value_Length": { /* 'DPT_Value_Length' */
            return "F32"
        }
        case "DPT_Value_Light_Quantity": { /* 'DPT_Value_Light_Quantity' */
            return "F32"
        }
        case "DPT_Value_Luminance": { /* 'DPT_Value_Luminance' */
            return "F32"
        }
        case "DPT_Value_Luminous_Flux": { /* 'DPT_Value_Luminous_Flux' */
            return "F32"
        }
        case "DPT_Value_Luminous_Intensity": { /* 'DPT_Value_Luminous_Intensity' */
            return "F32"
        }
        case "DPT_Value_Lux": { /* 'DPT_Value_Lux' */
            return "F16"
        }
        case "DPT_Value_Magnetic_FieldStrength": { /* 'DPT_Value_Magnetic_FieldStrength' */
            return "F32"
        }
        case "DPT_Value_Magnetic_Flux": { /* 'DPT_Value_Magnetic_Flux' */
            return "F32"
        }
        case "DPT_Value_Magnetic_FluxDensity": { /* 'DPT_Value_Magnetic_FluxDensity' */
            return "F32"
        }
        case "DPT_Value_Magnetic_Moment": { /* 'DPT_Value_Magnetic_Moment' */
            return "F32"
        }
        case "DPT_Value_Magnetic_Polarization": { /* 'DPT_Value_Magnetic_Polarization' */
            return "F32"
        }
        case "DPT_Value_Magnetization": { /* 'DPT_Value_Magnetization' */
            return "F32"
        }
        case "DPT_Value_MagnetomotiveForce": { /* 'DPT_Value_MagnetomotiveForce' */
            return "F32"
        }
        case "DPT_Value_Mass": { /* 'DPT_Value_Mass' */
            return "F32"
        }
        case "DPT_Value_MassFlux": { /* 'DPT_Value_MassFlux' */
            return "F32"
        }
        case "DPT_Value_Mol": { /* 'DPT_Value_Mol' */
            return "F32"
        }
        case "DPT_Value_Momentum": { /* 'DPT_Value_Momentum' */
            return "F32"
        }
        case "DPT_Value_Phase_AngleDeg": { /* 'DPT_Value_Phase_AngleDeg' */
            return "F32"
        }
        case "DPT_Value_Phase_AngleRad": { /* 'DPT_Value_Phase_AngleRad' */
            return "F32"
        }
        case "DPT_Value_Power": { /* 'DPT_Value_Power' */
            return "F32"
        }
        case "DPT_Value_Power_Factor": { /* 'DPT_Value_Power_Factor' */
            return "F32"
        }
        case "DPT_Value_Pres": { /* 'DPT_Value_Pres' */
            return "F16"
        }
        case "DPT_Value_Pressure": { /* 'DPT_Value_Pressure' */
            return "F32"
        }
        case "DPT_Value_Reactance": { /* 'DPT_Value_Reactance' */
            return "F32"
        }
        case "DPT_Value_Resistance": { /* 'DPT_Value_Resistance' */
            return "F32"
        }
        case "DPT_Value_Resistivity": { /* 'DPT_Value_Resistivity' */
            return "F32"
        }
        case "DPT_Value_SelfInductance": { /* 'DPT_Value_SelfInductance' */
            return "F32"
        }
        case "DPT_Value_SolidAngle": { /* 'DPT_Value_SolidAngle' */
            return "F32"
        }
        case "DPT_Value_Sound_Intensity": { /* 'DPT_Value_Sound_Intensity' */
            return "F32"
        }
        case "DPT_Value_Speed": { /* 'DPT_Value_Speed' */
            return "F32"
        }
        case "DPT_Value_Stress": { /* 'DPT_Value_Stress' */
            return "F32"
        }
        case "DPT_Value_Surface_Tension": { /* 'DPT_Value_Surface_Tension' */
            return "F32"
        }
        case "DPT_Value_Temp": { /* 'DPT_Value_Temp' */
            return "F16"
        }
        case "DPT_Value_Temp_F": { /* 'DPT_Value_Temp_F' */
            return "F16"
        }
        case "DPT_Value_Tempa": { /* 'DPT_Value_Tempa' */
            return "F16"
        }
        case "DPT_Value_Tempd": { /* 'DPT_Value_Tempd' */
            return "F16"
        }
        case "DPT_Value_TemperatureDifference": { /* 'DPT_Value_TemperatureDifference' */
            return "F32"
        }
        case "DPT_Value_Thermal_Capacity": { /* 'DPT_Value_Thermal_Capacity' */
            return "F32"
        }
        case "DPT_Value_Thermal_Conductivity": { /* 'DPT_Value_Thermal_Conductivity' */
            return "F32"
        }
        case "DPT_Value_ThermoelectricPower": { /* 'DPT_Value_ThermoelectricPower' */
            return "F32"
        }
        case "DPT_Value_Time": { /* 'DPT_Value_Time' */
            return "F32"
        }
        case "DPT_Value_Time1": { /* 'DPT_Value_Time1' */
            return "F16"
        }
        case "DPT_Value_Time2": { /* 'DPT_Value_Time2' */
            return "F16"
        }
        case "DPT_Value_Torque": { /* 'DPT_Value_Torque' */
            return "F32"
        }
        case "DPT_Value_Volt": { /* 'DPT_Value_Volt' */
            return "F16"
        }
        case "DPT_Value_Volume": { /* 'DPT_Value_Volume' */
            return "F32"
        }
        case "DPT_Value_Volume_Flow": { /* 'DPT_Value_Volume_Flow' */
            return "F16"
        }
        case "DPT_Value_Volume_Flux": { /* 'DPT_Value_Volume_Flux' */
            return "F32"
        }
        case "DPT_Value_Weight": { /* 'DPT_Value_Weight' */
            return "F32"
        }
        case "DPT_Value_Work": { /* 'DPT_Value_Work' */
            return "F32"
        }
        case "DPT_Value_Wsp": { /* 'DPT_Value_Wsp' */
            return "F16"
        }
        case "DPT_Value_Wsp_kmh": { /* 'DPT_Value_Wsp_kmh' */
            return "F16"
        }
        case "DPT_ValveMode": { /* 'DPT_ValveMode' */
            return "N8"
        }
        case "DPT_VarString_8859_1": { /* 'DPT_VarString_8859_1' */
            return "An_8859_1"
        }
        case "DPT_Version": { /* 'DPT_Version' */
            return "U5U5U6"
        }
        case "DPT_VolumeLiter_Z": { /* 'DPT_VolumeLiter_Z' */
            return "V32Z8VolumeLiter"
        }
        case "DPT_WindSpeed_Z_DPT_WindSpeed": { /* 'DPT_WindSpeed_Z_DPT_WindSpeed' */
            return "U16Z8WindSpeed"
        }
        case "DPT_Window_Door": { /* 'DPT_Window_Door' */
            return "B1"
        }
        default: {
            return ""
        }
    }
}

func (e KnxDatapointType) MainNumber() uint16 {
    switch e  {
        case "DPT_ADAType": { /* 'DPT_ADAType' */
            return 20
        }
        case "DPT_Access_Data": { /* 'DPT_Access_Data' */
            return 15
        }
        case "DPT_Ack": { /* 'DPT_Ack' */
            return 1
        }
        case "DPT_ActPosDemAbs": { /* 'DPT_ActPosDemAbs' */
            return 207
        }
        case "DPT_ActiveEnergy": { /* 'DPT_ActiveEnergy' */
            return 13
        }
        case "DPT_ActiveEnergy_V64": { /* 'DPT_ActiveEnergy_V64' */
            return 29
        }
        case "DPT_ActiveEnergy_kWh": { /* 'DPT_ActiveEnergy_kWh' */
            return 13
        }
        case "DPT_ActuatorConnectType": { /* 'DPT_ActuatorConnectType' */
            return 20
        }
        case "DPT_AddInfoTypes": { /* 'DPT_AddInfoTypes' */
            return 20
        }
        case "DPT_Alarm": { /* 'DPT_Alarm' */
            return 1
        }
        case "DPT_AlarmClassType": { /* 'DPT_AlarmClassType' */
            return 20
        }
        case "DPT_AlarmInfo": { /* 'DPT_AlarmInfo' */
            return 219
        }
        case "DPT_Alarm_Control": { /* 'DPT_Alarm_Control' */
            return 2
        }
        case "DPT_Alarm_Reaction": { /* 'DPT_Alarm_Reaction' */
            return 23
        }
        case "DPT_Angle": { /* 'DPT_Angle' */
            return 5
        }
        case "DPT_ApparantEnergy": { /* 'DPT_ApparantEnergy' */
            return 13
        }
        case "DPT_ApparantEnergy_V64": { /* 'DPT_ApparantEnergy_V64' */
            return 29
        }
        case "DPT_ApparantEnergy_kVAh": { /* 'DPT_ApparantEnergy_kVAh' */
            return 13
        }
        case "DPT_ApplicationArea": { /* 'DPT_ApplicationArea' */
            return 20
        }
        case "DPT_AtmPressureAbs_Z": { /* 'DPT_AtmPressureAbs_Z' */
            return 203
        }
        case "DPT_BackupMode": { /* 'DPT_BackupMode' */
            return 20
        }
        case "DPT_Beaufort_Wind_Force_Scale": { /* 'DPT_Beaufort_Wind_Force_Scale' */
            return 20
        }
        case "DPT_Behaviour_Bus_Power_Up_Down": { /* 'DPT_Behaviour_Bus_Power_Up_Down' */
            return 20
        }
        case "DPT_Behaviour_Lock_Unlock": { /* 'DPT_Behaviour_Lock_Unlock' */
            return 20
        }
        case "DPT_BinaryValue": { /* 'DPT_BinaryValue' */
            return 1
        }
        case "DPT_BinaryValue_Control": { /* 'DPT_BinaryValue_Control' */
            return 2
        }
        case "DPT_BinaryValue_Z": { /* 'DPT_BinaryValue_Z' */
            return 200
        }
        case "DPT_BlindsControlMode": { /* 'DPT_BlindsControlMode' */
            return 20
        }
        case "DPT_BlinkingMode": { /* 'DPT_BlinkingMode' */
            return 20
        }
        case "DPT_Bool": { /* 'DPT_Bool' */
            return 1
        }
        case "DPT_Bool_Control": { /* 'DPT_Bool_Control' */
            return 2
        }
        case "DPT_Brightness": { /* 'DPT_Brightness' */
            return 7
        }
        case "DPT_BuildingMode": { /* 'DPT_BuildingMode' */
            return 20
        }
        case "DPT_BuildingModeNext": { /* 'DPT_BuildingModeNext' */
            return 206
        }
        case "DPT_BuildingMode_Z": { /* 'DPT_BuildingMode_Z' */
            return 201
        }
        case "DPT_BurnerType": { /* 'DPT_BurnerType' */
            return 20
        }
        case "DPT_ChangeoverMode": { /* 'DPT_ChangeoverMode' */
            return 20
        }
        case "DPT_Channel_Activation_16": { /* 'DPT_Channel_Activation_16' */
            return 22
        }
        case "DPT_Channel_Activation_24": { /* 'DPT_Channel_Activation_24' */
            return 30
        }
        case "DPT_Channel_Activation_8": { /* 'DPT_Channel_Activation_8' */
            return 21
        }
        case "DPT_Char_8859_1": { /* 'DPT_Char_8859_1' */
            return 4
        }
        case "DPT_Char_ASCII": { /* 'DPT_Char_ASCII' */
            return 4
        }
        case "DPT_Colour_RGB": { /* 'DPT_Colour_RGB' */
            return 232
        }
        case "DPT_CombinedInfoOnOff": { /* 'DPT_CombinedInfoOnOff' */
            return 27
        }
        case "DPT_CombinedPosition": { /* 'DPT_CombinedPosition' */
            return 240
        }
        case "DPT_CommMode": { /* 'DPT_CommMode' */
            return 20
        }
        case "DPT_Control_Blinds": { /* 'DPT_Control_Blinds' */
            return 3
        }
        case "DPT_Control_Dimming": { /* 'DPT_Control_Dimming' */
            return 3
        }
        case "DPT_DALI_Control_Gear_Diagnostic": { /* 'DPT_DALI_Control_Gear_Diagnostic' */
            return 237
        }
        case "DPT_DALI_Diagnostics": { /* 'DPT_DALI_Diagnostics' */
            return 238
        }
        case "DPT_DALI_Fade_Time": { /* 'DPT_DALI_Fade_Time' */
            return 20
        }
        case "DPT_DHWMode": { /* 'DPT_DHWMode' */
            return 20
        }
        case "DPT_DHWModeNext": { /* 'DPT_DHWModeNext' */
            return 206
        }
        case "DPT_DHWMode_Z": { /* 'DPT_DHWMode_Z' */
            return 201
        }
        case "DPT_DamperMode": { /* 'DPT_DamperMode' */
            return 20
        }
        case "DPT_Date": { /* 'DPT_Date' */
            return 11
        }
        case "DPT_DateTime": { /* 'DPT_DateTime' */
            return 19
        }
        case "DPT_DecimalFactor": { /* 'DPT_DecimalFactor' */
            return 5
        }
        case "DPT_DeltaTime100MSec": { /* 'DPT_DeltaTime100MSec' */
            return 8
        }
        case "DPT_DeltaTime100Msec_Z": { /* 'DPT_DeltaTime100Msec_Z' */
            return 205
        }
        case "DPT_DeltaTime10MSec": { /* 'DPT_DeltaTime10MSec' */
            return 8
        }
        case "DPT_DeltaTime10Msec_Z": { /* 'DPT_DeltaTime10Msec_Z' */
            return 205
        }
        case "DPT_DeltaTimeHrs": { /* 'DPT_DeltaTimeHrs' */
            return 8
        }
        case "DPT_DeltaTimeHrs_Z": { /* 'DPT_DeltaTimeHrs_Z' */
            return 205
        }
        case "DPT_DeltaTimeMin": { /* 'DPT_DeltaTimeMin' */
            return 8
        }
        case "DPT_DeltaTimeMin_Z": { /* 'DPT_DeltaTimeMin_Z' */
            return 205
        }
        case "DPT_DeltaTimeMsec": { /* 'DPT_DeltaTimeMsec' */
            return 8
        }
        case "DPT_DeltaTimeMsec_Z": { /* 'DPT_DeltaTimeMsec_Z' */
            return 205
        }
        case "DPT_DeltaTimeSec": { /* 'DPT_DeltaTimeSec' */
            return 8
        }
        case "DPT_DeltaTimeSec_Z": { /* 'DPT_DeltaTimeSec_Z' */
            return 205
        }
        case "DPT_Device_Control": { /* 'DPT_Device_Control' */
            return 21
        }
        case "DPT_DimSendStyle": { /* 'DPT_DimSendStyle' */
            return 1
        }
        case "DPT_DimmPBModel": { /* 'DPT_DimmPBModel' */
            return 20
        }
        case "DPT_Direction1_Control": { /* 'DPT_Direction1_Control' */
            return 2
        }
        case "DPT_Direction2_Control": { /* 'DPT_Direction2_Control' */
            return 2
        }
        case "DPT_DoubleNibble": { /* 'DPT_DoubleNibble' */
            return 25
        }
        case "DPT_EnablH_Cstage_Z_DPT_EnablH_CStage": { /* 'DPT_EnablH_Cstage_Z_DPT_EnablH_CStage' */
            return 201
        }
        case "DPT_Enable": { /* 'DPT_Enable' */
            return 1
        }
        case "DPT_Enable_Control": { /* 'DPT_Enable_Control' */
            return 2
        }
        case "DPT_EnergyDemAir": { /* 'DPT_EnergyDemAir' */
            return 223
        }
        case "DPT_EnergyDemWater": { /* 'DPT_EnergyDemWater' */
            return 211
        }
        case "DPT_ErrorClass_HVAC": { /* 'DPT_ErrorClass_HVAC' */
            return 20
        }
        case "DPT_ErrorClass_System": { /* 'DPT_ErrorClass_System' */
            return 20
        }
        case "DPT_FanMode": { /* 'DPT_FanMode' */
            return 20
        }
        case "DPT_FlaggedScaling": { /* 'DPT_FlaggedScaling' */
            return 239
        }
        case "DPT_FlowRate_m3h": { /* 'DPT_FlowRate_m3h' */
            return 13
        }
        case "DPT_FlowRate_m3h_Z": { /* 'DPT_FlowRate_m3h_Z' */
            return 218
        }
        case "DPT_ForceSign": { /* 'DPT_ForceSign' */
            return 21
        }
        case "DPT_ForceSignCool": { /* 'DPT_ForceSignCool' */
            return 21
        }
        case "DPT_FuelType": { /* 'DPT_FuelType' */
            return 20
        }
        case "DPT_FuelTypeSet": { /* 'DPT_FuelTypeSet' */
            return 21
        }
        case "DPT_HVACAirFlowAbs_Z": { /* 'DPT_HVACAirFlowAbs_Z' */
            return 203
        }
        case "DPT_HVACAirFlowRel_Z": { /* 'DPT_HVACAirFlowRel_Z' */
            return 205
        }
        case "DPT_HVACAirQual_Z": { /* 'DPT_HVACAirQual_Z' */
            return 203
        }
        case "DPT_HVACContrMode": { /* 'DPT_HVACContrMode' */
            return 20
        }
        case "DPT_HVACContrMode_Z": { /* 'DPT_HVACContrMode_Z' */
            return 201
        }
        case "DPT_HVACEmergMode": { /* 'DPT_HVACEmergMode' */
            return 20
        }
        case "DPT_HVACEmergMode_Z": { /* 'DPT_HVACEmergMode_Z' */
            return 201
        }
        case "DPT_HVACMode": { /* 'DPT_HVACMode' */
            return 20
        }
        case "DPT_HVACModeNext": { /* 'DPT_HVACModeNext' */
            return 206
        }
        case "DPT_HVACMode_Z": { /* 'DPT_HVACMode_Z' */
            return 201
        }
        case "DPT_HVAC_PB_Action": { /* 'DPT_HVAC_PB_Action' */
            return 23
        }
        case "DPT_Heat_Cool": { /* 'DPT_Heat_Cool' */
            return 1
        }
        case "DPT_Heat_Cool_Z": { /* 'DPT_Heat_Cool_Z' */
            return 200
        }
        case "DPT_HeaterMode": { /* 'DPT_HeaterMode' */
            return 20
        }
        case "DPT_InputSource": { /* 'DPT_InputSource' */
            return 1
        }
        case "DPT_Invert": { /* 'DPT_Invert' */
            return 1
        }
        case "DPT_Invert_Control": { /* 'DPT_Invert_Control' */
            return 2
        }
        case "DPT_KelvinPerPercent": { /* 'DPT_KelvinPerPercent' */
            return 9
        }
        case "DPT_LanguageCodeAlpha2_ASCII": { /* 'DPT_LanguageCodeAlpha2_ASCII' */
            return 234
        }
        case "DPT_Length_mm": { /* 'DPT_Length_mm' */
            return 7
        }
        case "DPT_LightActuatorErrorInfo": { /* 'DPT_LightActuatorErrorInfo' */
            return 21
        }
        case "DPT_LightApplicationMode": { /* 'DPT_LightApplicationMode' */
            return 20
        }
        case "DPT_LightControlMode": { /* 'DPT_LightControlMode' */
            return 20
        }
        case "DPT_LoadPriority": { /* 'DPT_LoadPriority' */
            return 20
        }
        case "DPT_LoadTypeDetected": { /* 'DPT_LoadTypeDetected' */
            return 20
        }
        case "DPT_LoadTypeSet": { /* 'DPT_LoadTypeSet' */
            return 20
        }
        case "DPT_Locale_ASCII": { /* 'DPT_Locale_ASCII' */
            return 231
        }
        case "DPT_LockSign": { /* 'DPT_LockSign' */
            return 207
        }
        case "DPT_LogicalFunction": { /* 'DPT_LogicalFunction' */
            return 1
        }
        case "DPT_LongDeltaTimeSec": { /* 'DPT_LongDeltaTimeSec' */
            return 13
        }
        case "DPT_MBus_Address": { /* 'DPT_MBus_Address' */
            return 230
        }
        case "DPT_MasterSlaveMode": { /* 'DPT_MasterSlaveMode' */
            return 20
        }
        case "DPT_Media": { /* 'DPT_Media' */
            return 22
        }
        case "DPT_MeteringValue": { /* 'DPT_MeteringValue' */
            return 229
        }
        case "DPT_OccMode": { /* 'DPT_OccMode' */
            return 20
        }
        case "DPT_OccModeNext": { /* 'DPT_OccModeNext' */
            return 206
        }
        case "DPT_OccMode_Z": { /* 'DPT_OccMode_Z' */
            return 201
        }
        case "DPT_Occupancy": { /* 'DPT_Occupancy' */
            return 1
        }
        case "DPT_OnOff_Action": { /* 'DPT_OnOff_Action' */
            return 23
        }
        case "DPT_OpenClose": { /* 'DPT_OpenClose' */
            return 1
        }
        case "DPT_PBAction": { /* 'DPT_PBAction' */
            return 20
        }
        case "DPT_PB_Action_HVAC_Extended": { /* 'DPT_PB_Action_HVAC_Extended' */
            return 31
        }
        case "DPT_PSUMode": { /* 'DPT_PSUMode' */
            return 20
        }
        case "DPT_PercentU16_Z": { /* 'DPT_PercentU16_Z' */
            return 203
        }
        case "DPT_Percent_U8": { /* 'DPT_Percent_U8' */
            return 5
        }
        case "DPT_Percent_V16": { /* 'DPT_Percent_V16' */
            return 8
        }
        case "DPT_Percent_V16_Z": { /* 'DPT_Percent_V16_Z' */
            return 205
        }
        case "DPT_Percent_V8": { /* 'DPT_Percent_V8' */
            return 6
        }
        case "DPT_Power": { /* 'DPT_Power' */
            return 9
        }
        case "DPT_PowerDensity": { /* 'DPT_PowerDensity' */
            return 9
        }
        case "DPT_PowerFlowWaterDemCPM": { /* 'DPT_PowerFlowWaterDemCPM' */
            return 214
        }
        case "DPT_PowerFlowWaterDemHPM": { /* 'DPT_PowerFlowWaterDemHPM' */
            return 214
        }
        case "DPT_PowerKW_Z": { /* 'DPT_PowerKW_Z' */
            return 203
        }
        case "DPT_Prioritised_Mode_Control": { /* 'DPT_Prioritised_Mode_Control' */
            return 236
        }
        case "DPT_Priority": { /* 'DPT_Priority' */
            return 20
        }
        case "DPT_PropDataType": { /* 'DPT_PropDataType' */
            return 7
        }
        case "DPT_RF_FilterInfo": { /* 'DPT_RF_FilterInfo' */
            return 21
        }
        case "DPT_RF_FilterSelect": { /* 'DPT_RF_FilterSelect' */
            return 20
        }
        case "DPT_RF_ModeInfo": { /* 'DPT_RF_ModeInfo' */
            return 21
        }
        case "DPT_RF_ModeSelect": { /* 'DPT_RF_ModeSelect' */
            return 20
        }
        case "DPT_Rain_Amount": { /* 'DPT_Rain_Amount' */
            return 9
        }
        case "DPT_Ramp": { /* 'DPT_Ramp' */
            return 1
        }
        case "DPT_Ramp_Control": { /* 'DPT_Ramp_Control' */
            return 2
        }
        case "DPT_ReactiveEnergy": { /* 'DPT_ReactiveEnergy' */
            return 13
        }
        case "DPT_ReactiveEnergy_V64": { /* 'DPT_ReactiveEnergy_V64' */
            return 29
        }
        case "DPT_ReactiveEnergy_kVARh": { /* 'DPT_ReactiveEnergy_kVARh' */
            return 13
        }
        case "DPT_RegionCodeAlpha2_ASCII": { /* 'DPT_RegionCodeAlpha2_ASCII' */
            return 234
        }
        case "DPT_RelSignedValue_Z": { /* 'DPT_RelSignedValue_Z' */
            return 204
        }
        case "DPT_RelValue_Z": { /* 'DPT_RelValue_Z' */
            return 202
        }
        case "DPT_Reset": { /* 'DPT_Reset' */
            return 1
        }
        case "DPT_Rotation_Angle": { /* 'DPT_Rotation_Angle' */
            return 8
        }
        case "DPT_SABBehaviour_Lock_Unlock": { /* 'DPT_SABBehaviour_Lock_Unlock' */
            return 20
        }
        case "DPT_SABExceptBehaviour": { /* 'DPT_SABExceptBehaviour' */
            return 20
        }
        case "DPT_SCLOMode": { /* 'DPT_SCLOMode' */
            return 20
        }
        case "DPT_SSSBMode": { /* 'DPT_SSSBMode' */
            return 20
        }
        case "DPT_Scaling": { /* 'DPT_Scaling' */
            return 5
        }
        case "DPT_ScalingSpeed": { /* 'DPT_ScalingSpeed' */
            return 225
        }
        case "DPT_Scaling_Step_Time": { /* 'DPT_Scaling_Step_Time' */
            return 225
        }
        case "DPT_SceneConfig": { /* 'DPT_SceneConfig' */
            return 238
        }
        case "DPT_SceneControl": { /* 'DPT_SceneControl' */
            return 18
        }
        case "DPT_SceneInfo": { /* 'DPT_SceneInfo' */
            return 26
        }
        case "DPT_SceneNumber": { /* 'DPT_SceneNumber' */
            return 17
        }
        case "DPT_Scene_AB": { /* 'DPT_Scene_AB' */
            return 1
        }
        case "DPT_SensorSelect": { /* 'DPT_SensorSelect' */
            return 20
        }
        case "DPT_SerNum": { /* 'DPT_SerNum' */
            return 221
        }
        case "DPT_ShutterBlinds_Mode": { /* 'DPT_ShutterBlinds_Mode' */
            return 1
        }
        case "DPT_SpecHeatProd": { /* 'DPT_SpecHeatProd' */
            return 216
        }
        case "DPT_Start": { /* 'DPT_Start' */
            return 1
        }
        case "DPT_StartSynchronization": { /* 'DPT_StartSynchronization' */
            return 20
        }
        case "DPT_Start_Control": { /* 'DPT_Start_Control' */
            return 2
        }
        case "DPT_State": { /* 'DPT_State' */
            return 1
        }
        case "DPT_State_Control": { /* 'DPT_State_Control' */
            return 2
        }
        case "DPT_StatusAHU": { /* 'DPT_StatusAHU' */
            return 21
        }
        case "DPT_StatusAct": { /* 'DPT_StatusAct' */
            return 207
        }
        case "DPT_StatusBOC": { /* 'DPT_StatusBOC' */
            return 215
        }
        case "DPT_StatusBUC": { /* 'DPT_StatusBUC' */
            return 207
        }
        case "DPT_StatusCC": { /* 'DPT_StatusCC' */
            return 215
        }
        case "DPT_StatusCPM": { /* 'DPT_StatusCPM' */
            return 209
        }
        case "DPT_StatusDHWC": { /* 'DPT_StatusDHWC' */
            return 22
        }
        case "DPT_StatusGen": { /* 'DPT_StatusGen' */
            return 21
        }
        case "DPT_StatusHPM": { /* 'DPT_StatusHPM' */
            return 209
        }
        case "DPT_StatusLightingActuator": { /* 'DPT_StatusLightingActuator' */
            return 207
        }
        case "DPT_StatusRCC": { /* 'DPT_StatusRCC' */
            return 21
        }
        case "DPT_StatusRHC": { /* 'DPT_StatusRHC' */
            return 21
        }
        case "DPT_StatusRHCC": { /* 'DPT_StatusRHCC' */
            return 22
        }
        case "DPT_StatusRoomSetp": { /* 'DPT_StatusRoomSetp' */
            return 20
        }
        case "DPT_StatusSAB": { /* 'DPT_StatusSAB' */
            return 241
        }
        case "DPT_StatusSDHWC": { /* 'DPT_StatusSDHWC' */
            return 21
        }
        case "DPT_StatusWTC": { /* 'DPT_StatusWTC' */
            return 209
        }
        case "DPT_Status_Mode3": { /* 'DPT_Status_Mode3' */
            return 6
        }
        case "DPT_Step": { /* 'DPT_Step' */
            return 1
        }
        case "DPT_Step_Control": { /* 'DPT_Step_Control' */
            return 2
        }
        case "DPT_String_8859_1": { /* 'DPT_String_8859_1' */
            return 16
        }
        case "DPT_String_ASCII": { /* 'DPT_String_ASCII' */
            return 16
        }
        case "DPT_SunIntensity_Z": { /* 'DPT_SunIntensity_Z' */
            return 203
        }
        case "DPT_Switch": { /* 'DPT_Switch' */
            return 1
        }
        case "DPT_SwitchOnMode": { /* 'DPT_SwitchOnMode' */
            return 20
        }
        case "DPT_SwitchPBModel": { /* 'DPT_SwitchPBModel' */
            return 20
        }
        case "DPT_Switch_Control": { /* 'DPT_Switch_Control' */
            return 2
        }
        case "DPT_Tariff": { /* 'DPT_Tariff' */
            return 5
        }
        case "DPT_TariffNext": { /* 'DPT_TariffNext' */
            return 225
        }
        case "DPT_Tariff_ActiveEnergy": { /* 'DPT_Tariff_ActiveEnergy' */
            return 235
        }
        case "DPT_TempDHWSetpSet4": { /* 'DPT_TempDHWSetpSet4' */
            return 213
        }
        case "DPT_TempFlowWaterDemAbs": { /* 'DPT_TempFlowWaterDemAbs' */
            return 210
        }
        case "DPT_TempHVACAbsNext": { /* 'DPT_TempHVACAbsNext' */
            return 220
        }
        case "DPT_TempHVACAbs_Z": { /* 'DPT_TempHVACAbs_Z' */
            return 205
        }
        case "DPT_TempHVACRel_Z": { /* 'DPT_TempHVACRel_Z' */
            return 205
        }
        case "DPT_TempRoomDemAbs": { /* 'DPT_TempRoomDemAbs' */
            return 209
        }
        case "DPT_TempRoomSetpSet3": { /* 'DPT_TempRoomSetpSet3' */
            return 212
        }
        case "DPT_TempRoomSetpSet4": { /* 'DPT_TempRoomSetpSet4' */
            return 213
        }
        case "DPT_TempRoomSetpSetF163": { /* 'DPT_TempRoomSetpSetF163' */
            return 222
        }
        case "DPT_TempRoomSetpSetShift3": { /* 'DPT_TempRoomSetpSetShift3' */
            return 212
        }
        case "DPT_TempRoomSetpSetShift4": { /* 'DPT_TempRoomSetpSetShift4' */
            return 213
        }
        case "DPT_TempRoomSetpSetShiftF163": { /* 'DPT_TempRoomSetpSetShiftF163' */
            return 222
        }
        case "DPT_TempSupply_AirSetpSet": { /* 'DPT_TempSupply_AirSetpSet' */
            return 224
        }
        case "DPT_TimeOfDay": { /* 'DPT_TimeOfDay' */
            return 10
        }
        case "DPT_TimePeriod100MSec": { /* 'DPT_TimePeriod100MSec' */
            return 7
        }
        case "DPT_TimePeriod100Msec_Z": { /* 'DPT_TimePeriod100Msec_Z' */
            return 203
        }
        case "DPT_TimePeriod10MSec": { /* 'DPT_TimePeriod10MSec' */
            return 7
        }
        case "DPT_TimePeriod10Msec_Z": { /* 'DPT_TimePeriod10Msec_Z' */
            return 203
        }
        case "DPT_TimePeriodHrs": { /* 'DPT_TimePeriodHrs' */
            return 7
        }
        case "DPT_TimePeriodHrs_Z": { /* 'DPT_TimePeriodHrs_Z' */
            return 203
        }
        case "DPT_TimePeriodMin": { /* 'DPT_TimePeriodMin' */
            return 7
        }
        case "DPT_TimePeriodMin_Z": { /* 'DPT_TimePeriodMin_Z' */
            return 203
        }
        case "DPT_TimePeriodMsec": { /* 'DPT_TimePeriodMsec' */
            return 7
        }
        case "DPT_TimePeriodMsec_Z": { /* 'DPT_TimePeriodMsec_Z' */
            return 203
        }
        case "DPT_TimePeriodSec": { /* 'DPT_TimePeriodSec' */
            return 7
        }
        case "DPT_TimePeriodSec_Z": { /* 'DPT_TimePeriodSec_Z' */
            return 203
        }
        case "DPT_Time_Delay": { /* 'DPT_Time_Delay' */
            return 20
        }
        case "DPT_Trigger": { /* 'DPT_Trigger' */
            return 1
        }
        case "DPT_UCountValue16_Z": { /* 'DPT_UCountValue16_Z' */
            return 203
        }
        case "DPT_UCountValue8_Z": { /* 'DPT_UCountValue8_Z' */
            return 202
        }
        case "DPT_UElCurrentmA": { /* 'DPT_UElCurrentmA' */
            return 7
        }
        case "DPT_UElCurrentyA_Z": { /* 'DPT_UElCurrentyA_Z' */
            return 203
        }
        case "DPT_UFlowRateLiter_h_Z": { /* 'DPT_UFlowRateLiter_h_Z' */
            return 203
        }
        case "DPT_UTF_8": { /* 'DPT_UTF_8' */
            return 28
        }
        case "DPT_UpDown": { /* 'DPT_UpDown' */
            return 1
        }
        case "DPT_UpDown_Action": { /* 'DPT_UpDown_Action' */
            return 23
        }
        case "DPT_ValueDemBOC": { /* 'DPT_ValueDemBOC' */
            return 207
        }
        case "DPT_Value_1_Count": { /* 'DPT_Value_1_Count' */
            return 6
        }
        case "DPT_Value_1_Ucount": { /* 'DPT_Value_1_Ucount' */
            return 5
        }
        case "DPT_Value_2_Count": { /* 'DPT_Value_2_Count' */
            return 8
        }
        case "DPT_Value_2_Ucount": { /* 'DPT_Value_2_Ucount' */
            return 7
        }
        case "DPT_Value_4_Count": { /* 'DPT_Value_4_Count' */
            return 13
        }
        case "DPT_Value_4_Ucount": { /* 'DPT_Value_4_Ucount' */
            return 12
        }
        case "DPT_Value_Absolute_Temperature": { /* 'DPT_Value_Absolute_Temperature' */
            return 14
        }
        case "DPT_Value_Acceleration": { /* 'DPT_Value_Acceleration' */
            return 14
        }
        case "DPT_Value_Acceleration_Angular": { /* 'DPT_Value_Acceleration_Angular' */
            return 14
        }
        case "DPT_Value_Activation_Energy": { /* 'DPT_Value_Activation_Energy' */
            return 14
        }
        case "DPT_Value_Activity": { /* 'DPT_Value_Activity' */
            return 14
        }
        case "DPT_Value_AirQuality": { /* 'DPT_Value_AirQuality' */
            return 9
        }
        case "DPT_Value_Amplitude": { /* 'DPT_Value_Amplitude' */
            return 14
        }
        case "DPT_Value_AngleDeg": { /* 'DPT_Value_AngleDeg' */
            return 14
        }
        case "DPT_Value_AngleRad": { /* 'DPT_Value_AngleRad' */
            return 14
        }
        case "DPT_Value_Angular_Frequency": { /* 'DPT_Value_Angular_Frequency' */
            return 14
        }
        case "DPT_Value_Angular_Momentum": { /* 'DPT_Value_Angular_Momentum' */
            return 14
        }
        case "DPT_Value_Angular_Velocity": { /* 'DPT_Value_Angular_Velocity' */
            return 14
        }
        case "DPT_Value_Area": { /* 'DPT_Value_Area' */
            return 14
        }
        case "DPT_Value_Capacitance": { /* 'DPT_Value_Capacitance' */
            return 14
        }
        case "DPT_Value_Charge_DensitySurface": { /* 'DPT_Value_Charge_DensitySurface' */
            return 14
        }
        case "DPT_Value_Charge_DensityVolume": { /* 'DPT_Value_Charge_DensityVolume' */
            return 14
        }
        case "DPT_Value_Common_Temperature": { /* 'DPT_Value_Common_Temperature' */
            return 14
        }
        case "DPT_Value_Compressibility": { /* 'DPT_Value_Compressibility' */
            return 14
        }
        case "DPT_Value_Conductance": { /* 'DPT_Value_Conductance' */
            return 14
        }
        case "DPT_Value_Curr": { /* 'DPT_Value_Curr' */
            return 9
        }
        case "DPT_Value_Density": { /* 'DPT_Value_Density' */
            return 14
        }
        case "DPT_Value_Electric_Charge": { /* 'DPT_Value_Electric_Charge' */
            return 14
        }
        case "DPT_Value_Electric_Current": { /* 'DPT_Value_Electric_Current' */
            return 14
        }
        case "DPT_Value_Electric_CurrentDensity": { /* 'DPT_Value_Electric_CurrentDensity' */
            return 14
        }
        case "DPT_Value_Electric_DipoleMoment": { /* 'DPT_Value_Electric_DipoleMoment' */
            return 14
        }
        case "DPT_Value_Electric_Displacement": { /* 'DPT_Value_Electric_Displacement' */
            return 14
        }
        case "DPT_Value_Electric_FieldStrength": { /* 'DPT_Value_Electric_FieldStrength' */
            return 14
        }
        case "DPT_Value_Electric_Flux": { /* 'DPT_Value_Electric_Flux' */
            return 14
        }
        case "DPT_Value_Electric_FluxDensity": { /* 'DPT_Value_Electric_FluxDensity' */
            return 14
        }
        case "DPT_Value_Electric_Polarization": { /* 'DPT_Value_Electric_Polarization' */
            return 14
        }
        case "DPT_Value_Electric_Potential": { /* 'DPT_Value_Electric_Potential' */
            return 14
        }
        case "DPT_Value_Electric_PotentialDifference": { /* 'DPT_Value_Electric_PotentialDifference' */
            return 14
        }
        case "DPT_Value_Electrical_Conductivity": { /* 'DPT_Value_Electrical_Conductivity' */
            return 14
        }
        case "DPT_Value_ElectromagneticMoment": { /* 'DPT_Value_ElectromagneticMoment' */
            return 14
        }
        case "DPT_Value_Electromotive_Force": { /* 'DPT_Value_Electromotive_Force' */
            return 14
        }
        case "DPT_Value_Energy": { /* 'DPT_Value_Energy' */
            return 14
        }
        case "DPT_Value_Force": { /* 'DPT_Value_Force' */
            return 14
        }
        case "DPT_Value_Frequency": { /* 'DPT_Value_Frequency' */
            return 14
        }
        case "DPT_Value_Heat_Capacity": { /* 'DPT_Value_Heat_Capacity' */
            return 14
        }
        case "DPT_Value_Heat_FlowRate": { /* 'DPT_Value_Heat_FlowRate' */
            return 14
        }
        case "DPT_Value_Heat_Quantity": { /* 'DPT_Value_Heat_Quantity' */
            return 14
        }
        case "DPT_Value_Humidity": { /* 'DPT_Value_Humidity' */
            return 9
        }
        case "DPT_Value_Impedance": { /* 'DPT_Value_Impedance' */
            return 14
        }
        case "DPT_Value_Length": { /* 'DPT_Value_Length' */
            return 14
        }
        case "DPT_Value_Light_Quantity": { /* 'DPT_Value_Light_Quantity' */
            return 14
        }
        case "DPT_Value_Luminance": { /* 'DPT_Value_Luminance' */
            return 14
        }
        case "DPT_Value_Luminous_Flux": { /* 'DPT_Value_Luminous_Flux' */
            return 14
        }
        case "DPT_Value_Luminous_Intensity": { /* 'DPT_Value_Luminous_Intensity' */
            return 14
        }
        case "DPT_Value_Lux": { /* 'DPT_Value_Lux' */
            return 9
        }
        case "DPT_Value_Magnetic_FieldStrength": { /* 'DPT_Value_Magnetic_FieldStrength' */
            return 14
        }
        case "DPT_Value_Magnetic_Flux": { /* 'DPT_Value_Magnetic_Flux' */
            return 14
        }
        case "DPT_Value_Magnetic_FluxDensity": { /* 'DPT_Value_Magnetic_FluxDensity' */
            return 14
        }
        case "DPT_Value_Magnetic_Moment": { /* 'DPT_Value_Magnetic_Moment' */
            return 14
        }
        case "DPT_Value_Magnetic_Polarization": { /* 'DPT_Value_Magnetic_Polarization' */
            return 14
        }
        case "DPT_Value_Magnetization": { /* 'DPT_Value_Magnetization' */
            return 14
        }
        case "DPT_Value_MagnetomotiveForce": { /* 'DPT_Value_MagnetomotiveForce' */
            return 14
        }
        case "DPT_Value_Mass": { /* 'DPT_Value_Mass' */
            return 14
        }
        case "DPT_Value_MassFlux": { /* 'DPT_Value_MassFlux' */
            return 14
        }
        case "DPT_Value_Mol": { /* 'DPT_Value_Mol' */
            return 14
        }
        case "DPT_Value_Momentum": { /* 'DPT_Value_Momentum' */
            return 14
        }
        case "DPT_Value_Phase_AngleDeg": { /* 'DPT_Value_Phase_AngleDeg' */
            return 14
        }
        case "DPT_Value_Phase_AngleRad": { /* 'DPT_Value_Phase_AngleRad' */
            return 14
        }
        case "DPT_Value_Power": { /* 'DPT_Value_Power' */
            return 14
        }
        case "DPT_Value_Power_Factor": { /* 'DPT_Value_Power_Factor' */
            return 14
        }
        case "DPT_Value_Pres": { /* 'DPT_Value_Pres' */
            return 9
        }
        case "DPT_Value_Pressure": { /* 'DPT_Value_Pressure' */
            return 14
        }
        case "DPT_Value_Reactance": { /* 'DPT_Value_Reactance' */
            return 14
        }
        case "DPT_Value_Resistance": { /* 'DPT_Value_Resistance' */
            return 14
        }
        case "DPT_Value_Resistivity": { /* 'DPT_Value_Resistivity' */
            return 14
        }
        case "DPT_Value_SelfInductance": { /* 'DPT_Value_SelfInductance' */
            return 14
        }
        case "DPT_Value_SolidAngle": { /* 'DPT_Value_SolidAngle' */
            return 14
        }
        case "DPT_Value_Sound_Intensity": { /* 'DPT_Value_Sound_Intensity' */
            return 14
        }
        case "DPT_Value_Speed": { /* 'DPT_Value_Speed' */
            return 14
        }
        case "DPT_Value_Stress": { /* 'DPT_Value_Stress' */
            return 14
        }
        case "DPT_Value_Surface_Tension": { /* 'DPT_Value_Surface_Tension' */
            return 14
        }
        case "DPT_Value_Temp": { /* 'DPT_Value_Temp' */
            return 9
        }
        case "DPT_Value_Temp_F": { /* 'DPT_Value_Temp_F' */
            return 9
        }
        case "DPT_Value_Tempa": { /* 'DPT_Value_Tempa' */
            return 9
        }
        case "DPT_Value_Tempd": { /* 'DPT_Value_Tempd' */
            return 9
        }
        case "DPT_Value_TemperatureDifference": { /* 'DPT_Value_TemperatureDifference' */
            return 14
        }
        case "DPT_Value_Thermal_Capacity": { /* 'DPT_Value_Thermal_Capacity' */
            return 14
        }
        case "DPT_Value_Thermal_Conductivity": { /* 'DPT_Value_Thermal_Conductivity' */
            return 14
        }
        case "DPT_Value_ThermoelectricPower": { /* 'DPT_Value_ThermoelectricPower' */
            return 14
        }
        case "DPT_Value_Time": { /* 'DPT_Value_Time' */
            return 14
        }
        case "DPT_Value_Time1": { /* 'DPT_Value_Time1' */
            return 9
        }
        case "DPT_Value_Time2": { /* 'DPT_Value_Time2' */
            return 9
        }
        case "DPT_Value_Torque": { /* 'DPT_Value_Torque' */
            return 14
        }
        case "DPT_Value_Volt": { /* 'DPT_Value_Volt' */
            return 9
        }
        case "DPT_Value_Volume": { /* 'DPT_Value_Volume' */
            return 14
        }
        case "DPT_Value_Volume_Flow": { /* 'DPT_Value_Volume_Flow' */
            return 9
        }
        case "DPT_Value_Volume_Flux": { /* 'DPT_Value_Volume_Flux' */
            return 14
        }
        case "DPT_Value_Weight": { /* 'DPT_Value_Weight' */
            return 14
        }
        case "DPT_Value_Work": { /* 'DPT_Value_Work' */
            return 14
        }
        case "DPT_Value_Wsp": { /* 'DPT_Value_Wsp' */
            return 9
        }
        case "DPT_Value_Wsp_kmh": { /* 'DPT_Value_Wsp_kmh' */
            return 9
        }
        case "DPT_ValveMode": { /* 'DPT_ValveMode' */
            return 20
        }
        case "DPT_VarString_8859_1": { /* 'DPT_VarString_8859_1' */
            return 24
        }
        case "DPT_Version": { /* 'DPT_Version' */
            return 217
        }
        case "DPT_VolumeLiter_Z": { /* 'DPT_VolumeLiter_Z' */
            return 218
        }
        case "DPT_WindSpeed_Z_DPT_WindSpeed": { /* 'DPT_WindSpeed_Z_DPT_WindSpeed' */
            return 203
        }
        case "DPT_Window_Door": { /* 'DPT_Window_Door' */
            return 1
        }
        default: {
            return 0
        }
    }
}

func (e KnxDatapointType) SubNumber() uint16 {
    switch e  {
        case "DPT_ADAType": { /* 'DPT_ADAType' */
            return 120
        }
        case "DPT_Access_Data": { /* 'DPT_Access_Data' */
            return 0
        }
        case "DPT_Ack": { /* 'DPT_Ack' */
            return 16
        }
        case "DPT_ActPosDemAbs": { /* 'DPT_ActPosDemAbs' */
            return 104
        }
        case "DPT_ActiveEnergy": { /* 'DPT_ActiveEnergy' */
            return 10
        }
        case "DPT_ActiveEnergy_V64": { /* 'DPT_ActiveEnergy_V64' */
            return 10
        }
        case "DPT_ActiveEnergy_kWh": { /* 'DPT_ActiveEnergy_kWh' */
            return 13
        }
        case "DPT_ActuatorConnectType": { /* 'DPT_ActuatorConnectType' */
            return 20
        }
        case "DPT_AddInfoTypes": { /* 'DPT_AddInfoTypes' */
            return 1001
        }
        case "DPT_Alarm": { /* 'DPT_Alarm' */
            return 5
        }
        case "DPT_AlarmClassType": { /* 'DPT_AlarmClassType' */
            return 7
        }
        case "DPT_AlarmInfo": { /* 'DPT_AlarmInfo' */
            return 1
        }
        case "DPT_Alarm_Control": { /* 'DPT_Alarm_Control' */
            return 5
        }
        case "DPT_Alarm_Reaction": { /* 'DPT_Alarm_Reaction' */
            return 2
        }
        case "DPT_Angle": { /* 'DPT_Angle' */
            return 3
        }
        case "DPT_ApparantEnergy": { /* 'DPT_ApparantEnergy' */
            return 11
        }
        case "DPT_ApparantEnergy_V64": { /* 'DPT_ApparantEnergy_V64' */
            return 11
        }
        case "DPT_ApparantEnergy_kVAh": { /* 'DPT_ApparantEnergy_kVAh' */
            return 14
        }
        case "DPT_ApplicationArea": { /* 'DPT_ApplicationArea' */
            return 6
        }
        case "DPT_AtmPressureAbs_Z": { /* 'DPT_AtmPressureAbs_Z' */
            return 15
        }
        case "DPT_BackupMode": { /* 'DPT_BackupMode' */
            return 121
        }
        case "DPT_Beaufort_Wind_Force_Scale": { /* 'DPT_Beaufort_Wind_Force_Scale' */
            return 14
        }
        case "DPT_Behaviour_Bus_Power_Up_Down": { /* 'DPT_Behaviour_Bus_Power_Up_Down' */
            return 601
        }
        case "DPT_Behaviour_Lock_Unlock": { /* 'DPT_Behaviour_Lock_Unlock' */
            return 600
        }
        case "DPT_BinaryValue": { /* 'DPT_BinaryValue' */
            return 6
        }
        case "DPT_BinaryValue_Control": { /* 'DPT_BinaryValue_Control' */
            return 6
        }
        case "DPT_BinaryValue_Z": { /* 'DPT_BinaryValue_Z' */
            return 101
        }
        case "DPT_BlindsControlMode": { /* 'DPT_BlindsControlMode' */
            return 804
        }
        case "DPT_BlinkingMode": { /* 'DPT_BlinkingMode' */
            return 603
        }
        case "DPT_Bool": { /* 'DPT_Bool' */
            return 2
        }
        case "DPT_Bool_Control": { /* 'DPT_Bool_Control' */
            return 2
        }
        case "DPT_Brightness": { /* 'DPT_Brightness' */
            return 13
        }
        case "DPT_BuildingMode": { /* 'DPT_BuildingMode' */
            return 2
        }
        case "DPT_BuildingModeNext": { /* 'DPT_BuildingModeNext' */
            return 105
        }
        case "DPT_BuildingMode_Z": { /* 'DPT_BuildingMode_Z' */
            return 107
        }
        case "DPT_BurnerType": { /* 'DPT_BurnerType' */
            return 101
        }
        case "DPT_ChangeoverMode": { /* 'DPT_ChangeoverMode' */
            return 107
        }
        case "DPT_Channel_Activation_16": { /* 'DPT_Channel_Activation_16' */
            return 1010
        }
        case "DPT_Channel_Activation_24": { /* 'DPT_Channel_Activation_24' */
            return 1010
        }
        case "DPT_Channel_Activation_8": { /* 'DPT_Channel_Activation_8' */
            return 1010
        }
        case "DPT_Char_8859_1": { /* 'DPT_Char_8859_1' */
            return 2
        }
        case "DPT_Char_ASCII": { /* 'DPT_Char_ASCII' */
            return 1
        }
        case "DPT_Colour_RGB": { /* 'DPT_Colour_RGB' */
            return 600
        }
        case "DPT_CombinedInfoOnOff": { /* 'DPT_CombinedInfoOnOff' */
            return 1
        }
        case "DPT_CombinedPosition": { /* 'DPT_CombinedPosition' */
            return 800
        }
        case "DPT_CommMode": { /* 'DPT_CommMode' */
            return 1000
        }
        case "DPT_Control_Blinds": { /* 'DPT_Control_Blinds' */
            return 8
        }
        case "DPT_Control_Dimming": { /* 'DPT_Control_Dimming' */
            return 7
        }
        case "DPT_DALI_Control_Gear_Diagnostic": { /* 'DPT_DALI_Control_Gear_Diagnostic' */
            return 600
        }
        case "DPT_DALI_Diagnostics": { /* 'DPT_DALI_Diagnostics' */
            return 600
        }
        case "DPT_DALI_Fade_Time": { /* 'DPT_DALI_Fade_Time' */
            return 602
        }
        case "DPT_DHWMode": { /* 'DPT_DHWMode' */
            return 103
        }
        case "DPT_DHWModeNext": { /* 'DPT_DHWModeNext' */
            return 102
        }
        case "DPT_DHWMode_Z": { /* 'DPT_DHWMode_Z' */
            return 102
        }
        case "DPT_DamperMode": { /* 'DPT_DamperMode' */
            return 109
        }
        case "DPT_Date": { /* 'DPT_Date' */
            return 1
        }
        case "DPT_DateTime": { /* 'DPT_DateTime' */
            return 1
        }
        case "DPT_DecimalFactor": { /* 'DPT_DecimalFactor' */
            return 5
        }
        case "DPT_DeltaTime100MSec": { /* 'DPT_DeltaTime100MSec' */
            return 4
        }
        case "DPT_DeltaTime100Msec_Z": { /* 'DPT_DeltaTime100Msec_Z' */
            return 4
        }
        case "DPT_DeltaTime10MSec": { /* 'DPT_DeltaTime10MSec' */
            return 3
        }
        case "DPT_DeltaTime10Msec_Z": { /* 'DPT_DeltaTime10Msec_Z' */
            return 3
        }
        case "DPT_DeltaTimeHrs": { /* 'DPT_DeltaTimeHrs' */
            return 7
        }
        case "DPT_DeltaTimeHrs_Z": { /* 'DPT_DeltaTimeHrs_Z' */
            return 7
        }
        case "DPT_DeltaTimeMin": { /* 'DPT_DeltaTimeMin' */
            return 6
        }
        case "DPT_DeltaTimeMin_Z": { /* 'DPT_DeltaTimeMin_Z' */
            return 6
        }
        case "DPT_DeltaTimeMsec": { /* 'DPT_DeltaTimeMsec' */
            return 2
        }
        case "DPT_DeltaTimeMsec_Z": { /* 'DPT_DeltaTimeMsec_Z' */
            return 2
        }
        case "DPT_DeltaTimeSec": { /* 'DPT_DeltaTimeSec' */
            return 5
        }
        case "DPT_DeltaTimeSec_Z": { /* 'DPT_DeltaTimeSec_Z' */
            return 5
        }
        case "DPT_Device_Control": { /* 'DPT_Device_Control' */
            return 2
        }
        case "DPT_DimSendStyle": { /* 'DPT_DimSendStyle' */
            return 13
        }
        case "DPT_DimmPBModel": { /* 'DPT_DimmPBModel' */
            return 607
        }
        case "DPT_Direction1_Control": { /* 'DPT_Direction1_Control' */
            return 8
        }
        case "DPT_Direction2_Control": { /* 'DPT_Direction2_Control' */
            return 9
        }
        case "DPT_DoubleNibble": { /* 'DPT_DoubleNibble' */
            return 1000
        }
        case "DPT_EnablH_Cstage_Z_DPT_EnablH_CStage": { /* 'DPT_EnablH_Cstage_Z_DPT_EnablH_CStage' */
            return 105
        }
        case "DPT_Enable": { /* 'DPT_Enable' */
            return 3
        }
        case "DPT_Enable_Control": { /* 'DPT_Enable_Control' */
            return 3
        }
        case "DPT_EnergyDemAir": { /* 'DPT_EnergyDemAir' */
            return 100
        }
        case "DPT_EnergyDemWater": { /* 'DPT_EnergyDemWater' */
            return 100
        }
        case "DPT_ErrorClass_HVAC": { /* 'DPT_ErrorClass_HVAC' */
            return 12
        }
        case "DPT_ErrorClass_System": { /* 'DPT_ErrorClass_System' */
            return 11
        }
        case "DPT_FanMode": { /* 'DPT_FanMode' */
            return 111
        }
        case "DPT_FlaggedScaling": { /* 'DPT_FlaggedScaling' */
            return 1
        }
        case "DPT_FlowRate_m3h": { /* 'DPT_FlowRate_m3h' */
            return 2
        }
        case "DPT_FlowRate_m3h_Z": { /* 'DPT_FlowRate_m3h_Z' */
            return 2
        }
        case "DPT_ForceSign": { /* 'DPT_ForceSign' */
            return 100
        }
        case "DPT_ForceSignCool": { /* 'DPT_ForceSignCool' */
            return 101
        }
        case "DPT_FuelType": { /* 'DPT_FuelType' */
            return 100
        }
        case "DPT_FuelTypeSet": { /* 'DPT_FuelTypeSet' */
            return 104
        }
        case "DPT_HVACAirFlowAbs_Z": { /* 'DPT_HVACAirFlowAbs_Z' */
            return 104
        }
        case "DPT_HVACAirFlowRel_Z": { /* 'DPT_HVACAirFlowRel_Z' */
            return 102
        }
        case "DPT_HVACAirQual_Z": { /* 'DPT_HVACAirQual_Z' */
            return 100
        }
        case "DPT_HVACContrMode": { /* 'DPT_HVACContrMode' */
            return 105
        }
        case "DPT_HVACContrMode_Z": { /* 'DPT_HVACContrMode_Z' */
            return 104
        }
        case "DPT_HVACEmergMode": { /* 'DPT_HVACEmergMode' */
            return 106
        }
        case "DPT_HVACEmergMode_Z": { /* 'DPT_HVACEmergMode_Z' */
            return 109
        }
        case "DPT_HVACMode": { /* 'DPT_HVACMode' */
            return 102
        }
        case "DPT_HVACModeNext": { /* 'DPT_HVACModeNext' */
            return 100
        }
        case "DPT_HVACMode_Z": { /* 'DPT_HVACMode_Z' */
            return 100
        }
        case "DPT_HVAC_PB_Action": { /* 'DPT_HVAC_PB_Action' */
            return 102
        }
        case "DPT_Heat_Cool": { /* 'DPT_Heat_Cool' */
            return 100
        }
        case "DPT_Heat_Cool_Z": { /* 'DPT_Heat_Cool_Z' */
            return 100
        }
        case "DPT_HeaterMode": { /* 'DPT_HeaterMode' */
            return 110
        }
        case "DPT_InputSource": { /* 'DPT_InputSource' */
            return 14
        }
        case "DPT_Invert": { /* 'DPT_Invert' */
            return 12
        }
        case "DPT_Invert_Control": { /* 'DPT_Invert_Control' */
            return 12
        }
        case "DPT_KelvinPerPercent": { /* 'DPT_KelvinPerPercent' */
            return 23
        }
        case "DPT_LanguageCodeAlpha2_ASCII": { /* 'DPT_LanguageCodeAlpha2_ASCII' */
            return 1
        }
        case "DPT_Length_mm": { /* 'DPT_Length_mm' */
            return 11
        }
        case "DPT_LightActuatorErrorInfo": { /* 'DPT_LightActuatorErrorInfo' */
            return 601
        }
        case "DPT_LightApplicationMode": { /* 'DPT_LightApplicationMode' */
            return 5
        }
        case "DPT_LightControlMode": { /* 'DPT_LightControlMode' */
            return 604
        }
        case "DPT_LoadPriority": { /* 'DPT_LoadPriority' */
            return 104
        }
        case "DPT_LoadTypeDetected": { /* 'DPT_LoadTypeDetected' */
            return 610
        }
        case "DPT_LoadTypeSet": { /* 'DPT_LoadTypeSet' */
            return 609
        }
        case "DPT_Locale_ASCII": { /* 'DPT_Locale_ASCII' */
            return 1
        }
        case "DPT_LockSign": { /* 'DPT_LockSign' */
            return 101
        }
        case "DPT_LogicalFunction": { /* 'DPT_LogicalFunction' */
            return 21
        }
        case "DPT_LongDeltaTimeSec": { /* 'DPT_LongDeltaTimeSec' */
            return 100
        }
        case "DPT_MBus_Address": { /* 'DPT_MBus_Address' */
            return 1000
        }
        case "DPT_MasterSlaveMode": { /* 'DPT_MasterSlaveMode' */
            return 112
        }
        case "DPT_Media": { /* 'DPT_Media' */
            return 1000
        }
        case "DPT_MeteringValue": { /* 'DPT_MeteringValue' */
            return 1
        }
        case "DPT_OccMode": { /* 'DPT_OccMode' */
            return 3
        }
        case "DPT_OccModeNext": { /* 'DPT_OccModeNext' */
            return 104
        }
        case "DPT_OccMode_Z": { /* 'DPT_OccMode_Z' */
            return 108
        }
        case "DPT_Occupancy": { /* 'DPT_Occupancy' */
            return 18
        }
        case "DPT_OnOff_Action": { /* 'DPT_OnOff_Action' */
            return 1
        }
        case "DPT_OpenClose": { /* 'DPT_OpenClose' */
            return 9
        }
        case "DPT_PBAction": { /* 'DPT_PBAction' */
            return 606
        }
        case "DPT_PB_Action_HVAC_Extended": { /* 'DPT_PB_Action_HVAC_Extended' */
            return 101
        }
        case "DPT_PSUMode": { /* 'DPT_PSUMode' */
            return 8
        }
        case "DPT_PercentU16_Z": { /* 'DPT_PercentU16_Z' */
            return 17
        }
        case "DPT_Percent_U8": { /* 'DPT_Percent_U8' */
            return 4
        }
        case "DPT_Percent_V16": { /* 'DPT_Percent_V16' */
            return 10
        }
        case "DPT_Percent_V16_Z": { /* 'DPT_Percent_V16_Z' */
            return 17
        }
        case "DPT_Percent_V8": { /* 'DPT_Percent_V8' */
            return 1
        }
        case "DPT_Power": { /* 'DPT_Power' */
            return 24
        }
        case "DPT_PowerDensity": { /* 'DPT_PowerDensity' */
            return 22
        }
        case "DPT_PowerFlowWaterDemCPM": { /* 'DPT_PowerFlowWaterDemCPM' */
            return 101
        }
        case "DPT_PowerFlowWaterDemHPM": { /* 'DPT_PowerFlowWaterDemHPM' */
            return 100
        }
        case "DPT_PowerKW_Z": { /* 'DPT_PowerKW_Z' */
            return 14
        }
        case "DPT_Prioritised_Mode_Control": { /* 'DPT_Prioritised_Mode_Control' */
            return 1
        }
        case "DPT_Priority": { /* 'DPT_Priority' */
            return 4
        }
        case "DPT_PropDataType": { /* 'DPT_PropDataType' */
            return 10
        }
        case "DPT_RF_FilterInfo": { /* 'DPT_RF_FilterInfo' */
            return 1001
        }
        case "DPT_RF_FilterSelect": { /* 'DPT_RF_FilterSelect' */
            return 1003
        }
        case "DPT_RF_ModeInfo": { /* 'DPT_RF_ModeInfo' */
            return 1000
        }
        case "DPT_RF_ModeSelect": { /* 'DPT_RF_ModeSelect' */
            return 1002
        }
        case "DPT_Rain_Amount": { /* 'DPT_Rain_Amount' */
            return 26
        }
        case "DPT_Ramp": { /* 'DPT_Ramp' */
            return 4
        }
        case "DPT_Ramp_Control": { /* 'DPT_Ramp_Control' */
            return 4
        }
        case "DPT_ReactiveEnergy": { /* 'DPT_ReactiveEnergy' */
            return 12
        }
        case "DPT_ReactiveEnergy_V64": { /* 'DPT_ReactiveEnergy_V64' */
            return 12
        }
        case "DPT_ReactiveEnergy_kVARh": { /* 'DPT_ReactiveEnergy_kVARh' */
            return 15
        }
        case "DPT_RegionCodeAlpha2_ASCII": { /* 'DPT_RegionCodeAlpha2_ASCII' */
            return 2
        }
        case "DPT_RelSignedValue_Z": { /* 'DPT_RelSignedValue_Z' */
            return 1
        }
        case "DPT_RelValue_Z": { /* 'DPT_RelValue_Z' */
            return 1
        }
        case "DPT_Reset": { /* 'DPT_Reset' */
            return 15
        }
        case "DPT_Rotation_Angle": { /* 'DPT_Rotation_Angle' */
            return 11
        }
        case "DPT_SABBehaviour_Lock_Unlock": { /* 'DPT_SABBehaviour_Lock_Unlock' */
            return 802
        }
        case "DPT_SABExceptBehaviour": { /* 'DPT_SABExceptBehaviour' */
            return 801
        }
        case "DPT_SCLOMode": { /* 'DPT_SCLOMode' */
            return 1
        }
        case "DPT_SSSBMode": { /* 'DPT_SSSBMode' */
            return 803
        }
        case "DPT_Scaling": { /* 'DPT_Scaling' */
            return 1
        }
        case "DPT_ScalingSpeed": { /* 'DPT_ScalingSpeed' */
            return 1
        }
        case "DPT_Scaling_Step_Time": { /* 'DPT_Scaling_Step_Time' */
            return 2
        }
        case "DPT_SceneConfig": { /* 'DPT_SceneConfig' */
            return 1
        }
        case "DPT_SceneControl": { /* 'DPT_SceneControl' */
            return 1
        }
        case "DPT_SceneInfo": { /* 'DPT_SceneInfo' */
            return 1
        }
        case "DPT_SceneNumber": { /* 'DPT_SceneNumber' */
            return 1
        }
        case "DPT_Scene_AB": { /* 'DPT_Scene_AB' */
            return 22
        }
        case "DPT_SensorSelect": { /* 'DPT_SensorSelect' */
            return 17
        }
        case "DPT_SerNum": { /* 'DPT_SerNum' */
            return 1
        }
        case "DPT_ShutterBlinds_Mode": { /* 'DPT_ShutterBlinds_Mode' */
            return 23
        }
        case "DPT_SpecHeatProd": { /* 'DPT_SpecHeatProd' */
            return 100
        }
        case "DPT_Start": { /* 'DPT_Start' */
            return 10
        }
        case "DPT_StartSynchronization": { /* 'DPT_StartSynchronization' */
            return 122
        }
        case "DPT_Start_Control": { /* 'DPT_Start_Control' */
            return 10
        }
        case "DPT_State": { /* 'DPT_State' */
            return 11
        }
        case "DPT_State_Control": { /* 'DPT_State_Control' */
            return 11
        }
        case "DPT_StatusAHU": { /* 'DPT_StatusAHU' */
            return 106
        }
        case "DPT_StatusAct": { /* 'DPT_StatusAct' */
            return 105
        }
        case "DPT_StatusBOC": { /* 'DPT_StatusBOC' */
            return 100
        }
        case "DPT_StatusBUC": { /* 'DPT_StatusBUC' */
            return 100
        }
        case "DPT_StatusCC": { /* 'DPT_StatusCC' */
            return 101
        }
        case "DPT_StatusCPM": { /* 'DPT_StatusCPM' */
            return 102
        }
        case "DPT_StatusDHWC": { /* 'DPT_StatusDHWC' */
            return 100
        }
        case "DPT_StatusGen": { /* 'DPT_StatusGen' */
            return 1
        }
        case "DPT_StatusHPM": { /* 'DPT_StatusHPM' */
            return 100
        }
        case "DPT_StatusLightingActuator": { /* 'DPT_StatusLightingActuator' */
            return 600
        }
        case "DPT_StatusRCC": { /* 'DPT_StatusRCC' */
            return 105
        }
        case "DPT_StatusRHC": { /* 'DPT_StatusRHC' */
            return 102
        }
        case "DPT_StatusRHCC": { /* 'DPT_StatusRHCC' */
            return 101
        }
        case "DPT_StatusRoomSetp": { /* 'DPT_StatusRoomSetp' */
            return 113
        }
        case "DPT_StatusSAB": { /* 'DPT_StatusSAB' */
            return 800
        }
        case "DPT_StatusSDHWC": { /* 'DPT_StatusSDHWC' */
            return 103
        }
        case "DPT_StatusWTC": { /* 'DPT_StatusWTC' */
            return 103
        }
        case "DPT_Status_Mode3": { /* 'DPT_Status_Mode3' */
            return 20
        }
        case "DPT_Step": { /* 'DPT_Step' */
            return 7
        }
        case "DPT_Step_Control": { /* 'DPT_Step_Control' */
            return 7
        }
        case "DPT_String_8859_1": { /* 'DPT_String_8859_1' */
            return 1
        }
        case "DPT_String_ASCII": { /* 'DPT_String_ASCII' */
            return 0
        }
        case "DPT_SunIntensity_Z": { /* 'DPT_SunIntensity_Z' */
            return 102
        }
        case "DPT_Switch": { /* 'DPT_Switch' */
            return 1
        }
        case "DPT_SwitchOnMode": { /* 'DPT_SwitchOnMode' */
            return 608
        }
        case "DPT_SwitchPBModel": { /* 'DPT_SwitchPBModel' */
            return 605
        }
        case "DPT_Switch_Control": { /* 'DPT_Switch_Control' */
            return 1
        }
        case "DPT_Tariff": { /* 'DPT_Tariff' */
            return 6
        }
        case "DPT_TariffNext": { /* 'DPT_TariffNext' */
            return 3
        }
        case "DPT_Tariff_ActiveEnergy": { /* 'DPT_Tariff_ActiveEnergy' */
            return 1
        }
        case "DPT_TempDHWSetpSet4": { /* 'DPT_TempDHWSetpSet4' */
            return 101
        }
        case "DPT_TempFlowWaterDemAbs": { /* 'DPT_TempFlowWaterDemAbs' */
            return 100
        }
        case "DPT_TempHVACAbsNext": { /* 'DPT_TempHVACAbsNext' */
            return 100
        }
        case "DPT_TempHVACAbs_Z": { /* 'DPT_TempHVACAbs_Z' */
            return 100
        }
        case "DPT_TempHVACRel_Z": { /* 'DPT_TempHVACRel_Z' */
            return 101
        }
        case "DPT_TempRoomDemAbs": { /* 'DPT_TempRoomDemAbs' */
            return 101
        }
        case "DPT_TempRoomSetpSet3": { /* 'DPT_TempRoomSetpSet3' */
            return 101
        }
        case "DPT_TempRoomSetpSet4": { /* 'DPT_TempRoomSetpSet4' */
            return 100
        }
        case "DPT_TempRoomSetpSetF163": { /* 'DPT_TempRoomSetpSetF163' */
            return 100
        }
        case "DPT_TempRoomSetpSetShift3": { /* 'DPT_TempRoomSetpSetShift3' */
            return 100
        }
        case "DPT_TempRoomSetpSetShift4": { /* 'DPT_TempRoomSetpSetShift4' */
            return 102
        }
        case "DPT_TempRoomSetpSetShiftF163": { /* 'DPT_TempRoomSetpSetShiftF163' */
            return 101
        }
        case "DPT_TempSupply_AirSetpSet": { /* 'DPT_TempSupply_AirSetpSet' */
            return 100
        }
        case "DPT_TimeOfDay": { /* 'DPT_TimeOfDay' */
            return 1
        }
        case "DPT_TimePeriod100MSec": { /* 'DPT_TimePeriod100MSec' */
            return 4
        }
        case "DPT_TimePeriod100Msec_Z": { /* 'DPT_TimePeriod100Msec_Z' */
            return 4
        }
        case "DPT_TimePeriod10MSec": { /* 'DPT_TimePeriod10MSec' */
            return 3
        }
        case "DPT_TimePeriod10Msec_Z": { /* 'DPT_TimePeriod10Msec_Z' */
            return 3
        }
        case "DPT_TimePeriodHrs": { /* 'DPT_TimePeriodHrs' */
            return 7
        }
        case "DPT_TimePeriodHrs_Z": { /* 'DPT_TimePeriodHrs_Z' */
            return 7
        }
        case "DPT_TimePeriodMin": { /* 'DPT_TimePeriodMin' */
            return 6
        }
        case "DPT_TimePeriodMin_Z": { /* 'DPT_TimePeriodMin_Z' */
            return 6
        }
        case "DPT_TimePeriodMsec": { /* 'DPT_TimePeriodMsec' */
            return 2
        }
        case "DPT_TimePeriodMsec_Z": { /* 'DPT_TimePeriodMsec_Z' */
            return 2
        }
        case "DPT_TimePeriodSec": { /* 'DPT_TimePeriodSec' */
            return 5
        }
        case "DPT_TimePeriodSec_Z": { /* 'DPT_TimePeriodSec_Z' */
            return 5
        }
        case "DPT_Time_Delay": { /* 'DPT_Time_Delay' */
            return 13
        }
        case "DPT_Trigger": { /* 'DPT_Trigger' */
            return 17
        }
        case "DPT_UCountValue16_Z": { /* 'DPT_UCountValue16_Z' */
            return 12
        }
        case "DPT_UCountValue8_Z": { /* 'DPT_UCountValue8_Z' */
            return 2
        }
        case "DPT_UElCurrentmA": { /* 'DPT_UElCurrentmA' */
            return 12
        }
        case "DPT_UElCurrentyA_Z": { /* 'DPT_UElCurrentyA_Z' */
            return 13
        }
        case "DPT_UFlowRateLiter_h_Z": { /* 'DPT_UFlowRateLiter_h_Z' */
            return 11
        }
        case "DPT_UTF_8": { /* 'DPT_UTF_8' */
            return 1
        }
        case "DPT_UpDown": { /* 'DPT_UpDown' */
            return 8
        }
        case "DPT_UpDown_Action": { /* 'DPT_UpDown_Action' */
            return 3
        }
        case "DPT_ValueDemBOC": { /* 'DPT_ValueDemBOC' */
            return 102
        }
        case "DPT_Value_1_Count": { /* 'DPT_Value_1_Count' */
            return 10
        }
        case "DPT_Value_1_Ucount": { /* 'DPT_Value_1_Ucount' */
            return 10
        }
        case "DPT_Value_2_Count": { /* 'DPT_Value_2_Count' */
            return 1
        }
        case "DPT_Value_2_Ucount": { /* 'DPT_Value_2_Ucount' */
            return 1
        }
        case "DPT_Value_4_Count": { /* 'DPT_Value_4_Count' */
            return 1
        }
        case "DPT_Value_4_Ucount": { /* 'DPT_Value_4_Ucount' */
            return 1
        }
        case "DPT_Value_Absolute_Temperature": { /* 'DPT_Value_Absolute_Temperature' */
            return 69
        }
        case "DPT_Value_Acceleration": { /* 'DPT_Value_Acceleration' */
            return 0
        }
        case "DPT_Value_Acceleration_Angular": { /* 'DPT_Value_Acceleration_Angular' */
            return 1
        }
        case "DPT_Value_Activation_Energy": { /* 'DPT_Value_Activation_Energy' */
            return 2
        }
        case "DPT_Value_Activity": { /* 'DPT_Value_Activity' */
            return 3
        }
        case "DPT_Value_AirQuality": { /* 'DPT_Value_AirQuality' */
            return 8
        }
        case "DPT_Value_Amplitude": { /* 'DPT_Value_Amplitude' */
            return 5
        }
        case "DPT_Value_AngleDeg": { /* 'DPT_Value_AngleDeg' */
            return 7
        }
        case "DPT_Value_AngleRad": { /* 'DPT_Value_AngleRad' */
            return 6
        }
        case "DPT_Value_Angular_Frequency": { /* 'DPT_Value_Angular_Frequency' */
            return 34
        }
        case "DPT_Value_Angular_Momentum": { /* 'DPT_Value_Angular_Momentum' */
            return 8
        }
        case "DPT_Value_Angular_Velocity": { /* 'DPT_Value_Angular_Velocity' */
            return 9
        }
        case "DPT_Value_Area": { /* 'DPT_Value_Area' */
            return 10
        }
        case "DPT_Value_Capacitance": { /* 'DPT_Value_Capacitance' */
            return 11
        }
        case "DPT_Value_Charge_DensitySurface": { /* 'DPT_Value_Charge_DensitySurface' */
            return 12
        }
        case "DPT_Value_Charge_DensityVolume": { /* 'DPT_Value_Charge_DensityVolume' */
            return 13
        }
        case "DPT_Value_Common_Temperature": { /* 'DPT_Value_Common_Temperature' */
            return 68
        }
        case "DPT_Value_Compressibility": { /* 'DPT_Value_Compressibility' */
            return 14
        }
        case "DPT_Value_Conductance": { /* 'DPT_Value_Conductance' */
            return 15
        }
        case "DPT_Value_Curr": { /* 'DPT_Value_Curr' */
            return 21
        }
        case "DPT_Value_Density": { /* 'DPT_Value_Density' */
            return 17
        }
        case "DPT_Value_Electric_Charge": { /* 'DPT_Value_Electric_Charge' */
            return 18
        }
        case "DPT_Value_Electric_Current": { /* 'DPT_Value_Electric_Current' */
            return 19
        }
        case "DPT_Value_Electric_CurrentDensity": { /* 'DPT_Value_Electric_CurrentDensity' */
            return 20
        }
        case "DPT_Value_Electric_DipoleMoment": { /* 'DPT_Value_Electric_DipoleMoment' */
            return 21
        }
        case "DPT_Value_Electric_Displacement": { /* 'DPT_Value_Electric_Displacement' */
            return 22
        }
        case "DPT_Value_Electric_FieldStrength": { /* 'DPT_Value_Electric_FieldStrength' */
            return 23
        }
        case "DPT_Value_Electric_Flux": { /* 'DPT_Value_Electric_Flux' */
            return 24
        }
        case "DPT_Value_Electric_FluxDensity": { /* 'DPT_Value_Electric_FluxDensity' */
            return 25
        }
        case "DPT_Value_Electric_Polarization": { /* 'DPT_Value_Electric_Polarization' */
            return 26
        }
        case "DPT_Value_Electric_Potential": { /* 'DPT_Value_Electric_Potential' */
            return 27
        }
        case "DPT_Value_Electric_PotentialDifference": { /* 'DPT_Value_Electric_PotentialDifference' */
            return 28
        }
        case "DPT_Value_Electrical_Conductivity": { /* 'DPT_Value_Electrical_Conductivity' */
            return 16
        }
        case "DPT_Value_ElectromagneticMoment": { /* 'DPT_Value_ElectromagneticMoment' */
            return 29
        }
        case "DPT_Value_Electromotive_Force": { /* 'DPT_Value_Electromotive_Force' */
            return 30
        }
        case "DPT_Value_Energy": { /* 'DPT_Value_Energy' */
            return 31
        }
        case "DPT_Value_Force": { /* 'DPT_Value_Force' */
            return 32
        }
        case "DPT_Value_Frequency": { /* 'DPT_Value_Frequency' */
            return 33
        }
        case "DPT_Value_Heat_Capacity": { /* 'DPT_Value_Heat_Capacity' */
            return 35
        }
        case "DPT_Value_Heat_FlowRate": { /* 'DPT_Value_Heat_FlowRate' */
            return 36
        }
        case "DPT_Value_Heat_Quantity": { /* 'DPT_Value_Heat_Quantity' */
            return 37
        }
        case "DPT_Value_Humidity": { /* 'DPT_Value_Humidity' */
            return 7
        }
        case "DPT_Value_Impedance": { /* 'DPT_Value_Impedance' */
            return 38
        }
        case "DPT_Value_Length": { /* 'DPT_Value_Length' */
            return 39
        }
        case "DPT_Value_Light_Quantity": { /* 'DPT_Value_Light_Quantity' */
            return 40
        }
        case "DPT_Value_Luminance": { /* 'DPT_Value_Luminance' */
            return 41
        }
        case "DPT_Value_Luminous_Flux": { /* 'DPT_Value_Luminous_Flux' */
            return 42
        }
        case "DPT_Value_Luminous_Intensity": { /* 'DPT_Value_Luminous_Intensity' */
            return 43
        }
        case "DPT_Value_Lux": { /* 'DPT_Value_Lux' */
            return 4
        }
        case "DPT_Value_Magnetic_FieldStrength": { /* 'DPT_Value_Magnetic_FieldStrength' */
            return 44
        }
        case "DPT_Value_Magnetic_Flux": { /* 'DPT_Value_Magnetic_Flux' */
            return 45
        }
        case "DPT_Value_Magnetic_FluxDensity": { /* 'DPT_Value_Magnetic_FluxDensity' */
            return 46
        }
        case "DPT_Value_Magnetic_Moment": { /* 'DPT_Value_Magnetic_Moment' */
            return 47
        }
        case "DPT_Value_Magnetic_Polarization": { /* 'DPT_Value_Magnetic_Polarization' */
            return 48
        }
        case "DPT_Value_Magnetization": { /* 'DPT_Value_Magnetization' */
            return 49
        }
        case "DPT_Value_MagnetomotiveForce": { /* 'DPT_Value_MagnetomotiveForce' */
            return 50
        }
        case "DPT_Value_Mass": { /* 'DPT_Value_Mass' */
            return 51
        }
        case "DPT_Value_MassFlux": { /* 'DPT_Value_MassFlux' */
            return 52
        }
        case "DPT_Value_Mol": { /* 'DPT_Value_Mol' */
            return 4
        }
        case "DPT_Value_Momentum": { /* 'DPT_Value_Momentum' */
            return 53
        }
        case "DPT_Value_Phase_AngleDeg": { /* 'DPT_Value_Phase_AngleDeg' */
            return 55
        }
        case "DPT_Value_Phase_AngleRad": { /* 'DPT_Value_Phase_AngleRad' */
            return 54
        }
        case "DPT_Value_Power": { /* 'DPT_Value_Power' */
            return 56
        }
        case "DPT_Value_Power_Factor": { /* 'DPT_Value_Power_Factor' */
            return 57
        }
        case "DPT_Value_Pres": { /* 'DPT_Value_Pres' */
            return 6
        }
        case "DPT_Value_Pressure": { /* 'DPT_Value_Pressure' */
            return 58
        }
        case "DPT_Value_Reactance": { /* 'DPT_Value_Reactance' */
            return 59
        }
        case "DPT_Value_Resistance": { /* 'DPT_Value_Resistance' */
            return 60
        }
        case "DPT_Value_Resistivity": { /* 'DPT_Value_Resistivity' */
            return 61
        }
        case "DPT_Value_SelfInductance": { /* 'DPT_Value_SelfInductance' */
            return 62
        }
        case "DPT_Value_SolidAngle": { /* 'DPT_Value_SolidAngle' */
            return 63
        }
        case "DPT_Value_Sound_Intensity": { /* 'DPT_Value_Sound_Intensity' */
            return 64
        }
        case "DPT_Value_Speed": { /* 'DPT_Value_Speed' */
            return 65
        }
        case "DPT_Value_Stress": { /* 'DPT_Value_Stress' */
            return 66
        }
        case "DPT_Value_Surface_Tension": { /* 'DPT_Value_Surface_Tension' */
            return 67
        }
        case "DPT_Value_Temp": { /* 'DPT_Value_Temp' */
            return 1
        }
        case "DPT_Value_Temp_F": { /* 'DPT_Value_Temp_F' */
            return 27
        }
        case "DPT_Value_Tempa": { /* 'DPT_Value_Tempa' */
            return 3
        }
        case "DPT_Value_Tempd": { /* 'DPT_Value_Tempd' */
            return 2
        }
        case "DPT_Value_TemperatureDifference": { /* 'DPT_Value_TemperatureDifference' */
            return 70
        }
        case "DPT_Value_Thermal_Capacity": { /* 'DPT_Value_Thermal_Capacity' */
            return 71
        }
        case "DPT_Value_Thermal_Conductivity": { /* 'DPT_Value_Thermal_Conductivity' */
            return 72
        }
        case "DPT_Value_ThermoelectricPower": { /* 'DPT_Value_ThermoelectricPower' */
            return 73
        }
        case "DPT_Value_Time": { /* 'DPT_Value_Time' */
            return 74
        }
        case "DPT_Value_Time1": { /* 'DPT_Value_Time1' */
            return 10
        }
        case "DPT_Value_Time2": { /* 'DPT_Value_Time2' */
            return 11
        }
        case "DPT_Value_Torque": { /* 'DPT_Value_Torque' */
            return 75
        }
        case "DPT_Value_Volt": { /* 'DPT_Value_Volt' */
            return 20
        }
        case "DPT_Value_Volume": { /* 'DPT_Value_Volume' */
            return 76
        }
        case "DPT_Value_Volume_Flow": { /* 'DPT_Value_Volume_Flow' */
            return 25
        }
        case "DPT_Value_Volume_Flux": { /* 'DPT_Value_Volume_Flux' */
            return 77
        }
        case "DPT_Value_Weight": { /* 'DPT_Value_Weight' */
            return 78
        }
        case "DPT_Value_Work": { /* 'DPT_Value_Work' */
            return 79
        }
        case "DPT_Value_Wsp": { /* 'DPT_Value_Wsp' */
            return 5
        }
        case "DPT_Value_Wsp_kmh": { /* 'DPT_Value_Wsp_kmh' */
            return 28
        }
        case "DPT_ValveMode": { /* 'DPT_ValveMode' */
            return 108
        }
        case "DPT_VarString_8859_1": { /* 'DPT_VarString_8859_1' */
            return 1
        }
        case "DPT_Version": { /* 'DPT_Version' */
            return 1
        }
        case "DPT_VolumeLiter_Z": { /* 'DPT_VolumeLiter_Z' */
            return 1
        }
        case "DPT_WindSpeed_Z_DPT_WindSpeed": { /* 'DPT_WindSpeed_Z_DPT_WindSpeed' */
            return 101
        }
        case "DPT_Window_Door": { /* 'DPT_Window_Door' */
            return 19
        }
        default: {
            return 0
        }
    }
}
func KnxDatapointTypeValueOf(value string) KnxDatapointType {
    switch value {
        case "DPT_ADAType":
            return KnxDatapointType_DPT_ADAType
        case "DPT_Access_Data":
            return KnxDatapointType_DPT_Access_Data
        case "DPT_Ack":
            return KnxDatapointType_DPT_Ack
        case "DPT_ActPosDemAbs":
            return KnxDatapointType_DPT_ActPosDemAbs
        case "DPT_ActiveEnergy":
            return KnxDatapointType_DPT_ActiveEnergy
        case "DPT_ActiveEnergy_V64":
            return KnxDatapointType_DPT_ActiveEnergy_V64
        case "DPT_ActiveEnergy_kWh":
            return KnxDatapointType_DPT_ActiveEnergy_kWh
        case "DPT_ActuatorConnectType":
            return KnxDatapointType_DPT_ActuatorConnectType
        case "DPT_AddInfoTypes":
            return KnxDatapointType_DPT_AddInfoTypes
        case "DPT_Alarm":
            return KnxDatapointType_DPT_Alarm
        case "DPT_AlarmClassType":
            return KnxDatapointType_DPT_AlarmClassType
        case "DPT_AlarmInfo":
            return KnxDatapointType_DPT_AlarmInfo
        case "DPT_Alarm_Control":
            return KnxDatapointType_DPT_Alarm_Control
        case "DPT_Alarm_Reaction":
            return KnxDatapointType_DPT_Alarm_Reaction
        case "DPT_Angle":
            return KnxDatapointType_DPT_Angle
        case "DPT_ApparantEnergy":
            return KnxDatapointType_DPT_ApparantEnergy
        case "DPT_ApparantEnergy_V64":
            return KnxDatapointType_DPT_ApparantEnergy_V64
        case "DPT_ApparantEnergy_kVAh":
            return KnxDatapointType_DPT_ApparantEnergy_kVAh
        case "DPT_ApplicationArea":
            return KnxDatapointType_DPT_ApplicationArea
        case "DPT_AtmPressureAbs_Z":
            return KnxDatapointType_DPT_AtmPressureAbs_Z
        case "DPT_BackupMode":
            return KnxDatapointType_DPT_BackupMode
        case "DPT_Beaufort_Wind_Force_Scale":
            return KnxDatapointType_DPT_Beaufort_Wind_Force_Scale
        case "DPT_Behaviour_Bus_Power_Up_Down":
            return KnxDatapointType_DPT_Behaviour_Bus_Power_Up_Down
        case "DPT_Behaviour_Lock_Unlock":
            return KnxDatapointType_DPT_Behaviour_Lock_Unlock
        case "DPT_BinaryValue":
            return KnxDatapointType_DPT_BinaryValue
        case "DPT_BinaryValue_Control":
            return KnxDatapointType_DPT_BinaryValue_Control
        case "DPT_BinaryValue_Z":
            return KnxDatapointType_DPT_BinaryValue_Z
        case "DPT_BlindsControlMode":
            return KnxDatapointType_DPT_BlindsControlMode
        case "DPT_BlinkingMode":
            return KnxDatapointType_DPT_BlinkingMode
        case "DPT_Bool":
            return KnxDatapointType_DPT_Bool
        case "DPT_Bool_Control":
            return KnxDatapointType_DPT_Bool_Control
        case "DPT_Brightness":
            return KnxDatapointType_DPT_Brightness
        case "DPT_BuildingMode":
            return KnxDatapointType_DPT_BuildingMode
        case "DPT_BuildingModeNext":
            return KnxDatapointType_DPT_BuildingModeNext
        case "DPT_BuildingMode_Z":
            return KnxDatapointType_DPT_BuildingMode_Z
        case "DPT_BurnerType":
            return KnxDatapointType_DPT_BurnerType
        case "DPT_ChangeoverMode":
            return KnxDatapointType_DPT_ChangeoverMode
        case "DPT_Channel_Activation_16":
            return KnxDatapointType_DPT_Channel_Activation_16
        case "DPT_Channel_Activation_24":
            return KnxDatapointType_DPT_Channel_Activation_24
        case "DPT_Channel_Activation_8":
            return KnxDatapointType_DPT_Channel_Activation_8
        case "DPT_Char_8859_1":
            return KnxDatapointType_DPT_Char_8859_1
        case "DPT_Char_ASCII":
            return KnxDatapointType_DPT_Char_ASCII
        case "DPT_Colour_RGB":
            return KnxDatapointType_DPT_Colour_RGB
        case "DPT_CombinedInfoOnOff":
            return KnxDatapointType_DPT_CombinedInfoOnOff
        case "DPT_CombinedPosition":
            return KnxDatapointType_DPT_CombinedPosition
        case "DPT_CommMode":
            return KnxDatapointType_DPT_CommMode
        case "DPT_Control_Blinds":
            return KnxDatapointType_DPT_Control_Blinds
        case "DPT_Control_Dimming":
            return KnxDatapointType_DPT_Control_Dimming
        case "DPT_DALI_Control_Gear_Diagnostic":
            return KnxDatapointType_DPT_DALI_Control_Gear_Diagnostic
        case "DPT_DALI_Diagnostics":
            return KnxDatapointType_DPT_DALI_Diagnostics
        case "DPT_DALI_Fade_Time":
            return KnxDatapointType_DPT_DALI_Fade_Time
        case "DPT_DHWMode":
            return KnxDatapointType_DPT_DHWMode
        case "DPT_DHWModeNext":
            return KnxDatapointType_DPT_DHWModeNext
        case "DPT_DHWMode_Z":
            return KnxDatapointType_DPT_DHWMode_Z
        case "DPT_DamperMode":
            return KnxDatapointType_DPT_DamperMode
        case "DPT_Date":
            return KnxDatapointType_DPT_Date
        case "DPT_DateTime":
            return KnxDatapointType_DPT_DateTime
        case "DPT_DecimalFactor":
            return KnxDatapointType_DPT_DecimalFactor
        case "DPT_DeltaTime100MSec":
            return KnxDatapointType_DPT_DeltaTime100MSec
        case "DPT_DeltaTime100Msec_Z":
            return KnxDatapointType_DPT_DeltaTime100Msec_Z
        case "DPT_DeltaTime10MSec":
            return KnxDatapointType_DPT_DeltaTime10MSec
        case "DPT_DeltaTime10Msec_Z":
            return KnxDatapointType_DPT_DeltaTime10Msec_Z
        case "DPT_DeltaTimeHrs":
            return KnxDatapointType_DPT_DeltaTimeHrs
        case "DPT_DeltaTimeHrs_Z":
            return KnxDatapointType_DPT_DeltaTimeHrs_Z
        case "DPT_DeltaTimeMin":
            return KnxDatapointType_DPT_DeltaTimeMin
        case "DPT_DeltaTimeMin_Z":
            return KnxDatapointType_DPT_DeltaTimeMin_Z
        case "DPT_DeltaTimeMsec":
            return KnxDatapointType_DPT_DeltaTimeMsec
        case "DPT_DeltaTimeMsec_Z":
            return KnxDatapointType_DPT_DeltaTimeMsec_Z
        case "DPT_DeltaTimeSec":
            return KnxDatapointType_DPT_DeltaTimeSec
        case "DPT_DeltaTimeSec_Z":
            return KnxDatapointType_DPT_DeltaTimeSec_Z
        case "DPT_Device_Control":
            return KnxDatapointType_DPT_Device_Control
        case "DPT_DimSendStyle":
            return KnxDatapointType_DPT_DimSendStyle
        case "DPT_DimmPBModel":
            return KnxDatapointType_DPT_DimmPBModel
        case "DPT_Direction1_Control":
            return KnxDatapointType_DPT_Direction1_Control
        case "DPT_Direction2_Control":
            return KnxDatapointType_DPT_Direction2_Control
        case "DPT_DoubleNibble":
            return KnxDatapointType_DPT_DoubleNibble
        case "DPT_EnablH_Cstage_Z_DPT_EnablH_CStage":
            return KnxDatapointType_DPT_EnablH_Cstage_Z_DPT_EnablH_CStage
        case "DPT_Enable":
            return KnxDatapointType_DPT_Enable
        case "DPT_Enable_Control":
            return KnxDatapointType_DPT_Enable_Control
        case "DPT_EnergyDemAir":
            return KnxDatapointType_DPT_EnergyDemAir
        case "DPT_EnergyDemWater":
            return KnxDatapointType_DPT_EnergyDemWater
        case "DPT_ErrorClass_HVAC":
            return KnxDatapointType_DPT_ErrorClass_HVAC
        case "DPT_ErrorClass_System":
            return KnxDatapointType_DPT_ErrorClass_System
        case "DPT_FanMode":
            return KnxDatapointType_DPT_FanMode
        case "DPT_FlaggedScaling":
            return KnxDatapointType_DPT_FlaggedScaling
        case "DPT_FlowRate_m3h":
            return KnxDatapointType_DPT_FlowRate_m3h
        case "DPT_FlowRate_m3h_Z":
            return KnxDatapointType_DPT_FlowRate_m3h_Z
        case "DPT_ForceSign":
            return KnxDatapointType_DPT_ForceSign
        case "DPT_ForceSignCool":
            return KnxDatapointType_DPT_ForceSignCool
        case "DPT_FuelType":
            return KnxDatapointType_DPT_FuelType
        case "DPT_FuelTypeSet":
            return KnxDatapointType_DPT_FuelTypeSet
        case "DPT_HVACAirFlowAbs_Z":
            return KnxDatapointType_DPT_HVACAirFlowAbs_Z
        case "DPT_HVACAirFlowRel_Z":
            return KnxDatapointType_DPT_HVACAirFlowRel_Z
        case "DPT_HVACAirQual_Z":
            return KnxDatapointType_DPT_HVACAirQual_Z
        case "DPT_HVACContrMode":
            return KnxDatapointType_DPT_HVACContrMode
        case "DPT_HVACContrMode_Z":
            return KnxDatapointType_DPT_HVACContrMode_Z
        case "DPT_HVACEmergMode":
            return KnxDatapointType_DPT_HVACEmergMode
        case "DPT_HVACEmergMode_Z":
            return KnxDatapointType_DPT_HVACEmergMode_Z
        case "DPT_HVACMode":
            return KnxDatapointType_DPT_HVACMode
        case "DPT_HVACModeNext":
            return KnxDatapointType_DPT_HVACModeNext
        case "DPT_HVACMode_Z":
            return KnxDatapointType_DPT_HVACMode_Z
        case "DPT_HVAC_PB_Action":
            return KnxDatapointType_DPT_HVAC_PB_Action
        case "DPT_Heat_Cool":
            return KnxDatapointType_DPT_Heat_Cool
        case "DPT_Heat_Cool_Z":
            return KnxDatapointType_DPT_Heat_Cool_Z
        case "DPT_HeaterMode":
            return KnxDatapointType_DPT_HeaterMode
        case "DPT_InputSource":
            return KnxDatapointType_DPT_InputSource
        case "DPT_Invert":
            return KnxDatapointType_DPT_Invert
        case "DPT_Invert_Control":
            return KnxDatapointType_DPT_Invert_Control
        case "DPT_KelvinPerPercent":
            return KnxDatapointType_DPT_KelvinPerPercent
        case "DPT_LanguageCodeAlpha2_ASCII":
            return KnxDatapointType_DPT_LanguageCodeAlpha2_ASCII
        case "DPT_Length_mm":
            return KnxDatapointType_DPT_Length_mm
        case "DPT_LightActuatorErrorInfo":
            return KnxDatapointType_DPT_LightActuatorErrorInfo
        case "DPT_LightApplicationMode":
            return KnxDatapointType_DPT_LightApplicationMode
        case "DPT_LightControlMode":
            return KnxDatapointType_DPT_LightControlMode
        case "DPT_LoadPriority":
            return KnxDatapointType_DPT_LoadPriority
        case "DPT_LoadTypeDetected":
            return KnxDatapointType_DPT_LoadTypeDetected
        case "DPT_LoadTypeSet":
            return KnxDatapointType_DPT_LoadTypeSet
        case "DPT_Locale_ASCII":
            return KnxDatapointType_DPT_Locale_ASCII
        case "DPT_LockSign":
            return KnxDatapointType_DPT_LockSign
        case "DPT_LogicalFunction":
            return KnxDatapointType_DPT_LogicalFunction
        case "DPT_LongDeltaTimeSec":
            return KnxDatapointType_DPT_LongDeltaTimeSec
        case "DPT_MBus_Address":
            return KnxDatapointType_DPT_MBus_Address
        case "DPT_MasterSlaveMode":
            return KnxDatapointType_DPT_MasterSlaveMode
        case "DPT_Media":
            return KnxDatapointType_DPT_Media
        case "DPT_MeteringValue":
            return KnxDatapointType_DPT_MeteringValue
        case "DPT_OccMode":
            return KnxDatapointType_DPT_OccMode
        case "DPT_OccModeNext":
            return KnxDatapointType_DPT_OccModeNext
        case "DPT_OccMode_Z":
            return KnxDatapointType_DPT_OccMode_Z
        case "DPT_Occupancy":
            return KnxDatapointType_DPT_Occupancy
        case "DPT_OnOff_Action":
            return KnxDatapointType_DPT_OnOff_Action
        case "DPT_OpenClose":
            return KnxDatapointType_DPT_OpenClose
        case "DPT_PBAction":
            return KnxDatapointType_DPT_PBAction
        case "DPT_PB_Action_HVAC_Extended":
            return KnxDatapointType_DPT_PB_Action_HVAC_Extended
        case "DPT_PSUMode":
            return KnxDatapointType_DPT_PSUMode
        case "DPT_PercentU16_Z":
            return KnxDatapointType_DPT_PercentU16_Z
        case "DPT_Percent_U8":
            return KnxDatapointType_DPT_Percent_U8
        case "DPT_Percent_V16":
            return KnxDatapointType_DPT_Percent_V16
        case "DPT_Percent_V16_Z":
            return KnxDatapointType_DPT_Percent_V16_Z
        case "DPT_Percent_V8":
            return KnxDatapointType_DPT_Percent_V8
        case "DPT_Power":
            return KnxDatapointType_DPT_Power
        case "DPT_PowerDensity":
            return KnxDatapointType_DPT_PowerDensity
        case "DPT_PowerFlowWaterDemCPM":
            return KnxDatapointType_DPT_PowerFlowWaterDemCPM
        case "DPT_PowerFlowWaterDemHPM":
            return KnxDatapointType_DPT_PowerFlowWaterDemHPM
        case "DPT_PowerKW_Z":
            return KnxDatapointType_DPT_PowerKW_Z
        case "DPT_Prioritised_Mode_Control":
            return KnxDatapointType_DPT_Prioritised_Mode_Control
        case "DPT_Priority":
            return KnxDatapointType_DPT_Priority
        case "DPT_PropDataType":
            return KnxDatapointType_DPT_PropDataType
        case "DPT_RF_FilterInfo":
            return KnxDatapointType_DPT_RF_FilterInfo
        case "DPT_RF_FilterSelect":
            return KnxDatapointType_DPT_RF_FilterSelect
        case "DPT_RF_ModeInfo":
            return KnxDatapointType_DPT_RF_ModeInfo
        case "DPT_RF_ModeSelect":
            return KnxDatapointType_DPT_RF_ModeSelect
        case "DPT_Rain_Amount":
            return KnxDatapointType_DPT_Rain_Amount
        case "DPT_Ramp":
            return KnxDatapointType_DPT_Ramp
        case "DPT_Ramp_Control":
            return KnxDatapointType_DPT_Ramp_Control
        case "DPT_ReactiveEnergy":
            return KnxDatapointType_DPT_ReactiveEnergy
        case "DPT_ReactiveEnergy_V64":
            return KnxDatapointType_DPT_ReactiveEnergy_V64
        case "DPT_ReactiveEnergy_kVARh":
            return KnxDatapointType_DPT_ReactiveEnergy_kVARh
        case "DPT_RegionCodeAlpha2_ASCII":
            return KnxDatapointType_DPT_RegionCodeAlpha2_ASCII
        case "DPT_RelSignedValue_Z":
            return KnxDatapointType_DPT_RelSignedValue_Z
        case "DPT_RelValue_Z":
            return KnxDatapointType_DPT_RelValue_Z
        case "DPT_Reset":
            return KnxDatapointType_DPT_Reset
        case "DPT_Rotation_Angle":
            return KnxDatapointType_DPT_Rotation_Angle
        case "DPT_SABBehaviour_Lock_Unlock":
            return KnxDatapointType_DPT_SABBehaviour_Lock_Unlock
        case "DPT_SABExceptBehaviour":
            return KnxDatapointType_DPT_SABExceptBehaviour
        case "DPT_SCLOMode":
            return KnxDatapointType_DPT_SCLOMode
        case "DPT_SSSBMode":
            return KnxDatapointType_DPT_SSSBMode
        case "DPT_Scaling":
            return KnxDatapointType_DPT_Scaling
        case "DPT_ScalingSpeed":
            return KnxDatapointType_DPT_ScalingSpeed
        case "DPT_Scaling_Step_Time":
            return KnxDatapointType_DPT_Scaling_Step_Time
        case "DPT_SceneConfig":
            return KnxDatapointType_DPT_SceneConfig
        case "DPT_SceneControl":
            return KnxDatapointType_DPT_SceneControl
        case "DPT_SceneInfo":
            return KnxDatapointType_DPT_SceneInfo
        case "DPT_SceneNumber":
            return KnxDatapointType_DPT_SceneNumber
        case "DPT_Scene_AB":
            return KnxDatapointType_DPT_Scene_AB
        case "DPT_SensorSelect":
            return KnxDatapointType_DPT_SensorSelect
        case "DPT_SerNum":
            return KnxDatapointType_DPT_SerNum
        case "DPT_ShutterBlinds_Mode":
            return KnxDatapointType_DPT_ShutterBlinds_Mode
        case "DPT_SpecHeatProd":
            return KnxDatapointType_DPT_SpecHeatProd
        case "DPT_Start":
            return KnxDatapointType_DPT_Start
        case "DPT_StartSynchronization":
            return KnxDatapointType_DPT_StartSynchronization
        case "DPT_Start_Control":
            return KnxDatapointType_DPT_Start_Control
        case "DPT_State":
            return KnxDatapointType_DPT_State
        case "DPT_State_Control":
            return KnxDatapointType_DPT_State_Control
        case "DPT_StatusAHU":
            return KnxDatapointType_DPT_StatusAHU
        case "DPT_StatusAct":
            return KnxDatapointType_DPT_StatusAct
        case "DPT_StatusBOC":
            return KnxDatapointType_DPT_StatusBOC
        case "DPT_StatusBUC":
            return KnxDatapointType_DPT_StatusBUC
        case "DPT_StatusCC":
            return KnxDatapointType_DPT_StatusCC
        case "DPT_StatusCPM":
            return KnxDatapointType_DPT_StatusCPM
        case "DPT_StatusDHWC":
            return KnxDatapointType_DPT_StatusDHWC
        case "DPT_StatusGen":
            return KnxDatapointType_DPT_StatusGen
        case "DPT_StatusHPM":
            return KnxDatapointType_DPT_StatusHPM
        case "DPT_StatusLightingActuator":
            return KnxDatapointType_DPT_StatusLightingActuator
        case "DPT_StatusRCC":
            return KnxDatapointType_DPT_StatusRCC
        case "DPT_StatusRHC":
            return KnxDatapointType_DPT_StatusRHC
        case "DPT_StatusRHCC":
            return KnxDatapointType_DPT_StatusRHCC
        case "DPT_StatusRoomSetp":
            return KnxDatapointType_DPT_StatusRoomSetp
        case "DPT_StatusSAB":
            return KnxDatapointType_DPT_StatusSAB
        case "DPT_StatusSDHWC":
            return KnxDatapointType_DPT_StatusSDHWC
        case "DPT_StatusWTC":
            return KnxDatapointType_DPT_StatusWTC
        case "DPT_Status_Mode3":
            return KnxDatapointType_DPT_Status_Mode3
        case "DPT_Step":
            return KnxDatapointType_DPT_Step
        case "DPT_Step_Control":
            return KnxDatapointType_DPT_Step_Control
        case "DPT_String_8859_1":
            return KnxDatapointType_DPT_String_8859_1
        case "DPT_String_ASCII":
            return KnxDatapointType_DPT_String_ASCII
        case "DPT_SunIntensity_Z":
            return KnxDatapointType_DPT_SunIntensity_Z
        case "DPT_Switch":
            return KnxDatapointType_DPT_Switch
        case "DPT_SwitchOnMode":
            return KnxDatapointType_DPT_SwitchOnMode
        case "DPT_SwitchPBModel":
            return KnxDatapointType_DPT_SwitchPBModel
        case "DPT_Switch_Control":
            return KnxDatapointType_DPT_Switch_Control
        case "DPT_Tariff":
            return KnxDatapointType_DPT_Tariff
        case "DPT_TariffNext":
            return KnxDatapointType_DPT_TariffNext
        case "DPT_Tariff_ActiveEnergy":
            return KnxDatapointType_DPT_Tariff_ActiveEnergy
        case "DPT_TempDHWSetpSet4":
            return KnxDatapointType_DPT_TempDHWSetpSet4
        case "DPT_TempFlowWaterDemAbs":
            return KnxDatapointType_DPT_TempFlowWaterDemAbs
        case "DPT_TempHVACAbsNext":
            return KnxDatapointType_DPT_TempHVACAbsNext
        case "DPT_TempHVACAbs_Z":
            return KnxDatapointType_DPT_TempHVACAbs_Z
        case "DPT_TempHVACRel_Z":
            return KnxDatapointType_DPT_TempHVACRel_Z
        case "DPT_TempRoomDemAbs":
            return KnxDatapointType_DPT_TempRoomDemAbs
        case "DPT_TempRoomSetpSet3":
            return KnxDatapointType_DPT_TempRoomSetpSet3
        case "DPT_TempRoomSetpSet4":
            return KnxDatapointType_DPT_TempRoomSetpSet4
        case "DPT_TempRoomSetpSetF163":
            return KnxDatapointType_DPT_TempRoomSetpSetF163
        case "DPT_TempRoomSetpSetShift3":
            return KnxDatapointType_DPT_TempRoomSetpSetShift3
        case "DPT_TempRoomSetpSetShift4":
            return KnxDatapointType_DPT_TempRoomSetpSetShift4
        case "DPT_TempRoomSetpSetShiftF163":
            return KnxDatapointType_DPT_TempRoomSetpSetShiftF163
        case "DPT_TempSupply_AirSetpSet":
            return KnxDatapointType_DPT_TempSupply_AirSetpSet
        case "DPT_TimeOfDay":
            return KnxDatapointType_DPT_TimeOfDay
        case "DPT_TimePeriod100MSec":
            return KnxDatapointType_DPT_TimePeriod100MSec
        case "DPT_TimePeriod100Msec_Z":
            return KnxDatapointType_DPT_TimePeriod100Msec_Z
        case "DPT_TimePeriod10MSec":
            return KnxDatapointType_DPT_TimePeriod10MSec
        case "DPT_TimePeriod10Msec_Z":
            return KnxDatapointType_DPT_TimePeriod10Msec_Z
        case "DPT_TimePeriodHrs":
            return KnxDatapointType_DPT_TimePeriodHrs
        case "DPT_TimePeriodHrs_Z":
            return KnxDatapointType_DPT_TimePeriodHrs_Z
        case "DPT_TimePeriodMin":
            return KnxDatapointType_DPT_TimePeriodMin
        case "DPT_TimePeriodMin_Z":
            return KnxDatapointType_DPT_TimePeriodMin_Z
        case "DPT_TimePeriodMsec":
            return KnxDatapointType_DPT_TimePeriodMsec
        case "DPT_TimePeriodMsec_Z":
            return KnxDatapointType_DPT_TimePeriodMsec_Z
        case "DPT_TimePeriodSec":
            return KnxDatapointType_DPT_TimePeriodSec
        case "DPT_TimePeriodSec_Z":
            return KnxDatapointType_DPT_TimePeriodSec_Z
        case "DPT_Time_Delay":
            return KnxDatapointType_DPT_Time_Delay
        case "DPT_Trigger":
            return KnxDatapointType_DPT_Trigger
        case "DPT_UCountValue16_Z":
            return KnxDatapointType_DPT_UCountValue16_Z
        case "DPT_UCountValue8_Z":
            return KnxDatapointType_DPT_UCountValue8_Z
        case "DPT_UElCurrentmA":
            return KnxDatapointType_DPT_UElCurrentmA
        case "DPT_UElCurrentyA_Z":
            return KnxDatapointType_DPT_UElCurrentyA_Z
        case "DPT_UFlowRateLiter_h_Z":
            return KnxDatapointType_DPT_UFlowRateLiter_h_Z
        case "DPT_UTF_8":
            return KnxDatapointType_DPT_UTF_8
        case "DPT_UpDown":
            return KnxDatapointType_DPT_UpDown
        case "DPT_UpDown_Action":
            return KnxDatapointType_DPT_UpDown_Action
        case "DPT_ValueDemBOC":
            return KnxDatapointType_DPT_ValueDemBOC
        case "DPT_Value_1_Count":
            return KnxDatapointType_DPT_Value_1_Count
        case "DPT_Value_1_Ucount":
            return KnxDatapointType_DPT_Value_1_Ucount
        case "DPT_Value_2_Count":
            return KnxDatapointType_DPT_Value_2_Count
        case "DPT_Value_2_Ucount":
            return KnxDatapointType_DPT_Value_2_Ucount
        case "DPT_Value_4_Count":
            return KnxDatapointType_DPT_Value_4_Count
        case "DPT_Value_4_Ucount":
            return KnxDatapointType_DPT_Value_4_Ucount
        case "DPT_Value_Absolute_Temperature":
            return KnxDatapointType_DPT_Value_Absolute_Temperature
        case "DPT_Value_Acceleration":
            return KnxDatapointType_DPT_Value_Acceleration
        case "DPT_Value_Acceleration_Angular":
            return KnxDatapointType_DPT_Value_Acceleration_Angular
        case "DPT_Value_Activation_Energy":
            return KnxDatapointType_DPT_Value_Activation_Energy
        case "DPT_Value_Activity":
            return KnxDatapointType_DPT_Value_Activity
        case "DPT_Value_AirQuality":
            return KnxDatapointType_DPT_Value_AirQuality
        case "DPT_Value_Amplitude":
            return KnxDatapointType_DPT_Value_Amplitude
        case "DPT_Value_AngleDeg":
            return KnxDatapointType_DPT_Value_AngleDeg
        case "DPT_Value_AngleRad":
            return KnxDatapointType_DPT_Value_AngleRad
        case "DPT_Value_Angular_Frequency":
            return KnxDatapointType_DPT_Value_Angular_Frequency
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
        case "DPT_Value_Common_Temperature":
            return KnxDatapointType_DPT_Value_Common_Temperature
        case "DPT_Value_Compressibility":
            return KnxDatapointType_DPT_Value_Compressibility
        case "DPT_Value_Conductance":
            return KnxDatapointType_DPT_Value_Conductance
        case "DPT_Value_Curr":
            return KnxDatapointType_DPT_Value_Curr
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
        case "DPT_Value_Electric_FluxDensity":
            return KnxDatapointType_DPT_Value_Electric_FluxDensity
        case "DPT_Value_Electric_Polarization":
            return KnxDatapointType_DPT_Value_Electric_Polarization
        case "DPT_Value_Electric_Potential":
            return KnxDatapointType_DPT_Value_Electric_Potential
        case "DPT_Value_Electric_PotentialDifference":
            return KnxDatapointType_DPT_Value_Electric_PotentialDifference
        case "DPT_Value_Electrical_Conductivity":
            return KnxDatapointType_DPT_Value_Electrical_Conductivity
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
        case "DPT_Value_Heat_Capacity":
            return KnxDatapointType_DPT_Value_Heat_Capacity
        case "DPT_Value_Heat_FlowRate":
            return KnxDatapointType_DPT_Value_Heat_FlowRate
        case "DPT_Value_Heat_Quantity":
            return KnxDatapointType_DPT_Value_Heat_Quantity
        case "DPT_Value_Humidity":
            return KnxDatapointType_DPT_Value_Humidity
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
        case "DPT_Value_Lux":
            return KnxDatapointType_DPT_Value_Lux
        case "DPT_Value_Magnetic_FieldStrength":
            return KnxDatapointType_DPT_Value_Magnetic_FieldStrength
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
        case "DPT_Value_Mol":
            return KnxDatapointType_DPT_Value_Mol
        case "DPT_Value_Momentum":
            return KnxDatapointType_DPT_Value_Momentum
        case "DPT_Value_Phase_AngleDeg":
            return KnxDatapointType_DPT_Value_Phase_AngleDeg
        case "DPT_Value_Phase_AngleRad":
            return KnxDatapointType_DPT_Value_Phase_AngleRad
        case "DPT_Value_Power":
            return KnxDatapointType_DPT_Value_Power
        case "DPT_Value_Power_Factor":
            return KnxDatapointType_DPT_Value_Power_Factor
        case "DPT_Value_Pres":
            return KnxDatapointType_DPT_Value_Pres
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
        case "DPT_Value_Speed":
            return KnxDatapointType_DPT_Value_Speed
        case "DPT_Value_Stress":
            return KnxDatapointType_DPT_Value_Stress
        case "DPT_Value_Surface_Tension":
            return KnxDatapointType_DPT_Value_Surface_Tension
        case "DPT_Value_Temp":
            return KnxDatapointType_DPT_Value_Temp
        case "DPT_Value_Temp_F":
            return KnxDatapointType_DPT_Value_Temp_F
        case "DPT_Value_Tempa":
            return KnxDatapointType_DPT_Value_Tempa
        case "DPT_Value_Tempd":
            return KnxDatapointType_DPT_Value_Tempd
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
        case "DPT_Value_Time1":
            return KnxDatapointType_DPT_Value_Time1
        case "DPT_Value_Time2":
            return KnxDatapointType_DPT_Value_Time2
        case "DPT_Value_Torque":
            return KnxDatapointType_DPT_Value_Torque
        case "DPT_Value_Volt":
            return KnxDatapointType_DPT_Value_Volt
        case "DPT_Value_Volume":
            return KnxDatapointType_DPT_Value_Volume
        case "DPT_Value_Volume_Flow":
            return KnxDatapointType_DPT_Value_Volume_Flow
        case "DPT_Value_Volume_Flux":
            return KnxDatapointType_DPT_Value_Volume_Flux
        case "DPT_Value_Weight":
            return KnxDatapointType_DPT_Value_Weight
        case "DPT_Value_Work":
            return KnxDatapointType_DPT_Value_Work
        case "DPT_Value_Wsp":
            return KnxDatapointType_DPT_Value_Wsp
        case "DPT_Value_Wsp_kmh":
            return KnxDatapointType_DPT_Value_Wsp_kmh
        case "DPT_ValveMode":
            return KnxDatapointType_DPT_ValveMode
        case "DPT_VarString_8859_1":
            return KnxDatapointType_DPT_VarString_8859_1
        case "DPT_Version":
            return KnxDatapointType_DPT_Version
        case "DPT_VolumeLiter_Z":
            return KnxDatapointType_DPT_VolumeLiter_Z
        case "DPT_WindSpeed_Z_DPT_WindSpeed":
            return KnxDatapointType_DPT_WindSpeed_Z_DPT_WindSpeed
        case "DPT_Window_Door":
            return KnxDatapointType_DPT_Window_Door
    }
    return ""
}

func CastKnxDatapointType(structType interface{}) KnxDatapointType {
    castFunc := func(typ interface{}) KnxDatapointType {
        if sKnxDatapointType, ok := typ.(KnxDatapointType); ok {
            return sKnxDatapointType
        }
        return ""
    }
    return castFunc(structType)
}

func (m KnxDatapointType) LengthInBits() uint16 {
    return 0
}

func (m KnxDatapointType) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}


func (e KnxDatapointType) String() string {
    switch e {
    case KnxDatapointType_DPT_ADAType:
        return "DPT_ADAType"
    case KnxDatapointType_DPT_Access_Data:
        return "DPT_Access_Data"
    case KnxDatapointType_DPT_Ack:
        return "DPT_Ack"
    case KnxDatapointType_DPT_ActPosDemAbs:
        return "DPT_ActPosDemAbs"
    case KnxDatapointType_DPT_ActiveEnergy:
        return "DPT_ActiveEnergy"
    case KnxDatapointType_DPT_ActiveEnergy_V64:
        return "DPT_ActiveEnergy_V64"
    case KnxDatapointType_DPT_ActiveEnergy_kWh:
        return "DPT_ActiveEnergy_kWh"
    case KnxDatapointType_DPT_ActuatorConnectType:
        return "DPT_ActuatorConnectType"
    case KnxDatapointType_DPT_AddInfoTypes:
        return "DPT_AddInfoTypes"
    case KnxDatapointType_DPT_Alarm:
        return "DPT_Alarm"
    case KnxDatapointType_DPT_AlarmClassType:
        return "DPT_AlarmClassType"
    case KnxDatapointType_DPT_AlarmInfo:
        return "DPT_AlarmInfo"
    case KnxDatapointType_DPT_Alarm_Control:
        return "DPT_Alarm_Control"
    case KnxDatapointType_DPT_Alarm_Reaction:
        return "DPT_Alarm_Reaction"
    case KnxDatapointType_DPT_Angle:
        return "DPT_Angle"
    case KnxDatapointType_DPT_ApparantEnergy:
        return "DPT_ApparantEnergy"
    case KnxDatapointType_DPT_ApparantEnergy_V64:
        return "DPT_ApparantEnergy_V64"
    case KnxDatapointType_DPT_ApparantEnergy_kVAh:
        return "DPT_ApparantEnergy_kVAh"
    case KnxDatapointType_DPT_ApplicationArea:
        return "DPT_ApplicationArea"
    case KnxDatapointType_DPT_AtmPressureAbs_Z:
        return "DPT_AtmPressureAbs_Z"
    case KnxDatapointType_DPT_BackupMode:
        return "DPT_BackupMode"
    case KnxDatapointType_DPT_Beaufort_Wind_Force_Scale:
        return "DPT_Beaufort_Wind_Force_Scale"
    case KnxDatapointType_DPT_Behaviour_Bus_Power_Up_Down:
        return "DPT_Behaviour_Bus_Power_Up_Down"
    case KnxDatapointType_DPT_Behaviour_Lock_Unlock:
        return "DPT_Behaviour_Lock_Unlock"
    case KnxDatapointType_DPT_BinaryValue:
        return "DPT_BinaryValue"
    case KnxDatapointType_DPT_BinaryValue_Control:
        return "DPT_BinaryValue_Control"
    case KnxDatapointType_DPT_BinaryValue_Z:
        return "DPT_BinaryValue_Z"
    case KnxDatapointType_DPT_BlindsControlMode:
        return "DPT_BlindsControlMode"
    case KnxDatapointType_DPT_BlinkingMode:
        return "DPT_BlinkingMode"
    case KnxDatapointType_DPT_Bool:
        return "DPT_Bool"
    case KnxDatapointType_DPT_Bool_Control:
        return "DPT_Bool_Control"
    case KnxDatapointType_DPT_Brightness:
        return "DPT_Brightness"
    case KnxDatapointType_DPT_BuildingMode:
        return "DPT_BuildingMode"
    case KnxDatapointType_DPT_BuildingModeNext:
        return "DPT_BuildingModeNext"
    case KnxDatapointType_DPT_BuildingMode_Z:
        return "DPT_BuildingMode_Z"
    case KnxDatapointType_DPT_BurnerType:
        return "DPT_BurnerType"
    case KnxDatapointType_DPT_ChangeoverMode:
        return "DPT_ChangeoverMode"
    case KnxDatapointType_DPT_Channel_Activation_16:
        return "DPT_Channel_Activation_16"
    case KnxDatapointType_DPT_Channel_Activation_24:
        return "DPT_Channel_Activation_24"
    case KnxDatapointType_DPT_Channel_Activation_8:
        return "DPT_Channel_Activation_8"
    case KnxDatapointType_DPT_Char_8859_1:
        return "DPT_Char_8859_1"
    case KnxDatapointType_DPT_Char_ASCII:
        return "DPT_Char_ASCII"
    case KnxDatapointType_DPT_Colour_RGB:
        return "DPT_Colour_RGB"
    case KnxDatapointType_DPT_CombinedInfoOnOff:
        return "DPT_CombinedInfoOnOff"
    case KnxDatapointType_DPT_CombinedPosition:
        return "DPT_CombinedPosition"
    case KnxDatapointType_DPT_CommMode:
        return "DPT_CommMode"
    case KnxDatapointType_DPT_Control_Blinds:
        return "DPT_Control_Blinds"
    case KnxDatapointType_DPT_Control_Dimming:
        return "DPT_Control_Dimming"
    case KnxDatapointType_DPT_DALI_Control_Gear_Diagnostic:
        return "DPT_DALI_Control_Gear_Diagnostic"
    case KnxDatapointType_DPT_DALI_Diagnostics:
        return "DPT_DALI_Diagnostics"
    case KnxDatapointType_DPT_DALI_Fade_Time:
        return "DPT_DALI_Fade_Time"
    case KnxDatapointType_DPT_DHWMode:
        return "DPT_DHWMode"
    case KnxDatapointType_DPT_DHWModeNext:
        return "DPT_DHWModeNext"
    case KnxDatapointType_DPT_DHWMode_Z:
        return "DPT_DHWMode_Z"
    case KnxDatapointType_DPT_DamperMode:
        return "DPT_DamperMode"
    case KnxDatapointType_DPT_Date:
        return "DPT_Date"
    case KnxDatapointType_DPT_DateTime:
        return "DPT_DateTime"
    case KnxDatapointType_DPT_DecimalFactor:
        return "DPT_DecimalFactor"
    case KnxDatapointType_DPT_DeltaTime100MSec:
        return "DPT_DeltaTime100MSec"
    case KnxDatapointType_DPT_DeltaTime100Msec_Z:
        return "DPT_DeltaTime100Msec_Z"
    case KnxDatapointType_DPT_DeltaTime10MSec:
        return "DPT_DeltaTime10MSec"
    case KnxDatapointType_DPT_DeltaTime10Msec_Z:
        return "DPT_DeltaTime10Msec_Z"
    case KnxDatapointType_DPT_DeltaTimeHrs:
        return "DPT_DeltaTimeHrs"
    case KnxDatapointType_DPT_DeltaTimeHrs_Z:
        return "DPT_DeltaTimeHrs_Z"
    case KnxDatapointType_DPT_DeltaTimeMin:
        return "DPT_DeltaTimeMin"
    case KnxDatapointType_DPT_DeltaTimeMin_Z:
        return "DPT_DeltaTimeMin_Z"
    case KnxDatapointType_DPT_DeltaTimeMsec:
        return "DPT_DeltaTimeMsec"
    case KnxDatapointType_DPT_DeltaTimeMsec_Z:
        return "DPT_DeltaTimeMsec_Z"
    case KnxDatapointType_DPT_DeltaTimeSec:
        return "DPT_DeltaTimeSec"
    case KnxDatapointType_DPT_DeltaTimeSec_Z:
        return "DPT_DeltaTimeSec_Z"
    case KnxDatapointType_DPT_Device_Control:
        return "DPT_Device_Control"
    case KnxDatapointType_DPT_DimSendStyle:
        return "DPT_DimSendStyle"
    case KnxDatapointType_DPT_DimmPBModel:
        return "DPT_DimmPBModel"
    case KnxDatapointType_DPT_Direction1_Control:
        return "DPT_Direction1_Control"
    case KnxDatapointType_DPT_Direction2_Control:
        return "DPT_Direction2_Control"
    case KnxDatapointType_DPT_DoubleNibble:
        return "DPT_DoubleNibble"
    case KnxDatapointType_DPT_EnablH_Cstage_Z_DPT_EnablH_CStage:
        return "DPT_EnablH_Cstage_Z_DPT_EnablH_CStage"
    case KnxDatapointType_DPT_Enable:
        return "DPT_Enable"
    case KnxDatapointType_DPT_Enable_Control:
        return "DPT_Enable_Control"
    case KnxDatapointType_DPT_EnergyDemAir:
        return "DPT_EnergyDemAir"
    case KnxDatapointType_DPT_EnergyDemWater:
        return "DPT_EnergyDemWater"
    case KnxDatapointType_DPT_ErrorClass_HVAC:
        return "DPT_ErrorClass_HVAC"
    case KnxDatapointType_DPT_ErrorClass_System:
        return "DPT_ErrorClass_System"
    case KnxDatapointType_DPT_FanMode:
        return "DPT_FanMode"
    case KnxDatapointType_DPT_FlaggedScaling:
        return "DPT_FlaggedScaling"
    case KnxDatapointType_DPT_FlowRate_m3h:
        return "DPT_FlowRate_m3h"
    case KnxDatapointType_DPT_FlowRate_m3h_Z:
        return "DPT_FlowRate_m3h_Z"
    case KnxDatapointType_DPT_ForceSign:
        return "DPT_ForceSign"
    case KnxDatapointType_DPT_ForceSignCool:
        return "DPT_ForceSignCool"
    case KnxDatapointType_DPT_FuelType:
        return "DPT_FuelType"
    case KnxDatapointType_DPT_FuelTypeSet:
        return "DPT_FuelTypeSet"
    case KnxDatapointType_DPT_HVACAirFlowAbs_Z:
        return "DPT_HVACAirFlowAbs_Z"
    case KnxDatapointType_DPT_HVACAirFlowRel_Z:
        return "DPT_HVACAirFlowRel_Z"
    case KnxDatapointType_DPT_HVACAirQual_Z:
        return "DPT_HVACAirQual_Z"
    case KnxDatapointType_DPT_HVACContrMode:
        return "DPT_HVACContrMode"
    case KnxDatapointType_DPT_HVACContrMode_Z:
        return "DPT_HVACContrMode_Z"
    case KnxDatapointType_DPT_HVACEmergMode:
        return "DPT_HVACEmergMode"
    case KnxDatapointType_DPT_HVACEmergMode_Z:
        return "DPT_HVACEmergMode_Z"
    case KnxDatapointType_DPT_HVACMode:
        return "DPT_HVACMode"
    case KnxDatapointType_DPT_HVACModeNext:
        return "DPT_HVACModeNext"
    case KnxDatapointType_DPT_HVACMode_Z:
        return "DPT_HVACMode_Z"
    case KnxDatapointType_DPT_HVAC_PB_Action:
        return "DPT_HVAC_PB_Action"
    case KnxDatapointType_DPT_Heat_Cool:
        return "DPT_Heat_Cool"
    case KnxDatapointType_DPT_Heat_Cool_Z:
        return "DPT_Heat_Cool_Z"
    case KnxDatapointType_DPT_HeaterMode:
        return "DPT_HeaterMode"
    case KnxDatapointType_DPT_InputSource:
        return "DPT_InputSource"
    case KnxDatapointType_DPT_Invert:
        return "DPT_Invert"
    case KnxDatapointType_DPT_Invert_Control:
        return "DPT_Invert_Control"
    case KnxDatapointType_DPT_KelvinPerPercent:
        return "DPT_KelvinPerPercent"
    case KnxDatapointType_DPT_LanguageCodeAlpha2_ASCII:
        return "DPT_LanguageCodeAlpha2_ASCII"
    case KnxDatapointType_DPT_Length_mm:
        return "DPT_Length_mm"
    case KnxDatapointType_DPT_LightActuatorErrorInfo:
        return "DPT_LightActuatorErrorInfo"
    case KnxDatapointType_DPT_LightApplicationMode:
        return "DPT_LightApplicationMode"
    case KnxDatapointType_DPT_LightControlMode:
        return "DPT_LightControlMode"
    case KnxDatapointType_DPT_LoadPriority:
        return "DPT_LoadPriority"
    case KnxDatapointType_DPT_LoadTypeDetected:
        return "DPT_LoadTypeDetected"
    case KnxDatapointType_DPT_LoadTypeSet:
        return "DPT_LoadTypeSet"
    case KnxDatapointType_DPT_Locale_ASCII:
        return "DPT_Locale_ASCII"
    case KnxDatapointType_DPT_LockSign:
        return "DPT_LockSign"
    case KnxDatapointType_DPT_LogicalFunction:
        return "DPT_LogicalFunction"
    case KnxDatapointType_DPT_LongDeltaTimeSec:
        return "DPT_LongDeltaTimeSec"
    case KnxDatapointType_DPT_MBus_Address:
        return "DPT_MBus_Address"
    case KnxDatapointType_DPT_MasterSlaveMode:
        return "DPT_MasterSlaveMode"
    case KnxDatapointType_DPT_Media:
        return "DPT_Media"
    case KnxDatapointType_DPT_MeteringValue:
        return "DPT_MeteringValue"
    case KnxDatapointType_DPT_OccMode:
        return "DPT_OccMode"
    case KnxDatapointType_DPT_OccModeNext:
        return "DPT_OccModeNext"
    case KnxDatapointType_DPT_OccMode_Z:
        return "DPT_OccMode_Z"
    case KnxDatapointType_DPT_Occupancy:
        return "DPT_Occupancy"
    case KnxDatapointType_DPT_OnOff_Action:
        return "DPT_OnOff_Action"
    case KnxDatapointType_DPT_OpenClose:
        return "DPT_OpenClose"
    case KnxDatapointType_DPT_PBAction:
        return "DPT_PBAction"
    case KnxDatapointType_DPT_PB_Action_HVAC_Extended:
        return "DPT_PB_Action_HVAC_Extended"
    case KnxDatapointType_DPT_PSUMode:
        return "DPT_PSUMode"
    case KnxDatapointType_DPT_PercentU16_Z:
        return "DPT_PercentU16_Z"
    case KnxDatapointType_DPT_Percent_U8:
        return "DPT_Percent_U8"
    case KnxDatapointType_DPT_Percent_V16:
        return "DPT_Percent_V16"
    case KnxDatapointType_DPT_Percent_V16_Z:
        return "DPT_Percent_V16_Z"
    case KnxDatapointType_DPT_Percent_V8:
        return "DPT_Percent_V8"
    case KnxDatapointType_DPT_Power:
        return "DPT_Power"
    case KnxDatapointType_DPT_PowerDensity:
        return "DPT_PowerDensity"
    case KnxDatapointType_DPT_PowerFlowWaterDemCPM:
        return "DPT_PowerFlowWaterDemCPM"
    case KnxDatapointType_DPT_PowerFlowWaterDemHPM:
        return "DPT_PowerFlowWaterDemHPM"
    case KnxDatapointType_DPT_PowerKW_Z:
        return "DPT_PowerKW_Z"
    case KnxDatapointType_DPT_Prioritised_Mode_Control:
        return "DPT_Prioritised_Mode_Control"
    case KnxDatapointType_DPT_Priority:
        return "DPT_Priority"
    case KnxDatapointType_DPT_PropDataType:
        return "DPT_PropDataType"
    case KnxDatapointType_DPT_RF_FilterInfo:
        return "DPT_RF_FilterInfo"
    case KnxDatapointType_DPT_RF_FilterSelect:
        return "DPT_RF_FilterSelect"
    case KnxDatapointType_DPT_RF_ModeInfo:
        return "DPT_RF_ModeInfo"
    case KnxDatapointType_DPT_RF_ModeSelect:
        return "DPT_RF_ModeSelect"
    case KnxDatapointType_DPT_Rain_Amount:
        return "DPT_Rain_Amount"
    case KnxDatapointType_DPT_Ramp:
        return "DPT_Ramp"
    case KnxDatapointType_DPT_Ramp_Control:
        return "DPT_Ramp_Control"
    case KnxDatapointType_DPT_ReactiveEnergy:
        return "DPT_ReactiveEnergy"
    case KnxDatapointType_DPT_ReactiveEnergy_V64:
        return "DPT_ReactiveEnergy_V64"
    case KnxDatapointType_DPT_ReactiveEnergy_kVARh:
        return "DPT_ReactiveEnergy_kVARh"
    case KnxDatapointType_DPT_RegionCodeAlpha2_ASCII:
        return "DPT_RegionCodeAlpha2_ASCII"
    case KnxDatapointType_DPT_RelSignedValue_Z:
        return "DPT_RelSignedValue_Z"
    case KnxDatapointType_DPT_RelValue_Z:
        return "DPT_RelValue_Z"
    case KnxDatapointType_DPT_Reset:
        return "DPT_Reset"
    case KnxDatapointType_DPT_Rotation_Angle:
        return "DPT_Rotation_Angle"
    case KnxDatapointType_DPT_SABBehaviour_Lock_Unlock:
        return "DPT_SABBehaviour_Lock_Unlock"
    case KnxDatapointType_DPT_SABExceptBehaviour:
        return "DPT_SABExceptBehaviour"
    case KnxDatapointType_DPT_SCLOMode:
        return "DPT_SCLOMode"
    case KnxDatapointType_DPT_SSSBMode:
        return "DPT_SSSBMode"
    case KnxDatapointType_DPT_Scaling:
        return "DPT_Scaling"
    case KnxDatapointType_DPT_ScalingSpeed:
        return "DPT_ScalingSpeed"
    case KnxDatapointType_DPT_Scaling_Step_Time:
        return "DPT_Scaling_Step_Time"
    case KnxDatapointType_DPT_SceneConfig:
        return "DPT_SceneConfig"
    case KnxDatapointType_DPT_SceneControl:
        return "DPT_SceneControl"
    case KnxDatapointType_DPT_SceneInfo:
        return "DPT_SceneInfo"
    case KnxDatapointType_DPT_SceneNumber:
        return "DPT_SceneNumber"
    case KnxDatapointType_DPT_Scene_AB:
        return "DPT_Scene_AB"
    case KnxDatapointType_DPT_SensorSelect:
        return "DPT_SensorSelect"
    case KnxDatapointType_DPT_SerNum:
        return "DPT_SerNum"
    case KnxDatapointType_DPT_ShutterBlinds_Mode:
        return "DPT_ShutterBlinds_Mode"
    case KnxDatapointType_DPT_SpecHeatProd:
        return "DPT_SpecHeatProd"
    case KnxDatapointType_DPT_Start:
        return "DPT_Start"
    case KnxDatapointType_DPT_StartSynchronization:
        return "DPT_StartSynchronization"
    case KnxDatapointType_DPT_Start_Control:
        return "DPT_Start_Control"
    case KnxDatapointType_DPT_State:
        return "DPT_State"
    case KnxDatapointType_DPT_State_Control:
        return "DPT_State_Control"
    case KnxDatapointType_DPT_StatusAHU:
        return "DPT_StatusAHU"
    case KnxDatapointType_DPT_StatusAct:
        return "DPT_StatusAct"
    case KnxDatapointType_DPT_StatusBOC:
        return "DPT_StatusBOC"
    case KnxDatapointType_DPT_StatusBUC:
        return "DPT_StatusBUC"
    case KnxDatapointType_DPT_StatusCC:
        return "DPT_StatusCC"
    case KnxDatapointType_DPT_StatusCPM:
        return "DPT_StatusCPM"
    case KnxDatapointType_DPT_StatusDHWC:
        return "DPT_StatusDHWC"
    case KnxDatapointType_DPT_StatusGen:
        return "DPT_StatusGen"
    case KnxDatapointType_DPT_StatusHPM:
        return "DPT_StatusHPM"
    case KnxDatapointType_DPT_StatusLightingActuator:
        return "DPT_StatusLightingActuator"
    case KnxDatapointType_DPT_StatusRCC:
        return "DPT_StatusRCC"
    case KnxDatapointType_DPT_StatusRHC:
        return "DPT_StatusRHC"
    case KnxDatapointType_DPT_StatusRHCC:
        return "DPT_StatusRHCC"
    case KnxDatapointType_DPT_StatusRoomSetp:
        return "DPT_StatusRoomSetp"
    case KnxDatapointType_DPT_StatusSAB:
        return "DPT_StatusSAB"
    case KnxDatapointType_DPT_StatusSDHWC:
        return "DPT_StatusSDHWC"
    case KnxDatapointType_DPT_StatusWTC:
        return "DPT_StatusWTC"
    case KnxDatapointType_DPT_Status_Mode3:
        return "DPT_Status_Mode3"
    case KnxDatapointType_DPT_Step:
        return "DPT_Step"
    case KnxDatapointType_DPT_Step_Control:
        return "DPT_Step_Control"
    case KnxDatapointType_DPT_String_8859_1:
        return "DPT_String_8859_1"
    case KnxDatapointType_DPT_String_ASCII:
        return "DPT_String_ASCII"
    case KnxDatapointType_DPT_SunIntensity_Z:
        return "DPT_SunIntensity_Z"
    case KnxDatapointType_DPT_Switch:
        return "DPT_Switch"
    case KnxDatapointType_DPT_SwitchOnMode:
        return "DPT_SwitchOnMode"
    case KnxDatapointType_DPT_SwitchPBModel:
        return "DPT_SwitchPBModel"
    case KnxDatapointType_DPT_Switch_Control:
        return "DPT_Switch_Control"
    case KnxDatapointType_DPT_Tariff:
        return "DPT_Tariff"
    case KnxDatapointType_DPT_TariffNext:
        return "DPT_TariffNext"
    case KnxDatapointType_DPT_Tariff_ActiveEnergy:
        return "DPT_Tariff_ActiveEnergy"
    case KnxDatapointType_DPT_TempDHWSetpSet4:
        return "DPT_TempDHWSetpSet4"
    case KnxDatapointType_DPT_TempFlowWaterDemAbs:
        return "DPT_TempFlowWaterDemAbs"
    case KnxDatapointType_DPT_TempHVACAbsNext:
        return "DPT_TempHVACAbsNext"
    case KnxDatapointType_DPT_TempHVACAbs_Z:
        return "DPT_TempHVACAbs_Z"
    case KnxDatapointType_DPT_TempHVACRel_Z:
        return "DPT_TempHVACRel_Z"
    case KnxDatapointType_DPT_TempRoomDemAbs:
        return "DPT_TempRoomDemAbs"
    case KnxDatapointType_DPT_TempRoomSetpSet3:
        return "DPT_TempRoomSetpSet3"
    case KnxDatapointType_DPT_TempRoomSetpSet4:
        return "DPT_TempRoomSetpSet4"
    case KnxDatapointType_DPT_TempRoomSetpSetF163:
        return "DPT_TempRoomSetpSetF163"
    case KnxDatapointType_DPT_TempRoomSetpSetShift3:
        return "DPT_TempRoomSetpSetShift3"
    case KnxDatapointType_DPT_TempRoomSetpSetShift4:
        return "DPT_TempRoomSetpSetShift4"
    case KnxDatapointType_DPT_TempRoomSetpSetShiftF163:
        return "DPT_TempRoomSetpSetShiftF163"
    case KnxDatapointType_DPT_TempSupply_AirSetpSet:
        return "DPT_TempSupply_AirSetpSet"
    case KnxDatapointType_DPT_TimeOfDay:
        return "DPT_TimeOfDay"
    case KnxDatapointType_DPT_TimePeriod100MSec:
        return "DPT_TimePeriod100MSec"
    case KnxDatapointType_DPT_TimePeriod100Msec_Z:
        return "DPT_TimePeriod100Msec_Z"
    case KnxDatapointType_DPT_TimePeriod10MSec:
        return "DPT_TimePeriod10MSec"
    case KnxDatapointType_DPT_TimePeriod10Msec_Z:
        return "DPT_TimePeriod10Msec_Z"
    case KnxDatapointType_DPT_TimePeriodHrs:
        return "DPT_TimePeriodHrs"
    case KnxDatapointType_DPT_TimePeriodHrs_Z:
        return "DPT_TimePeriodHrs_Z"
    case KnxDatapointType_DPT_TimePeriodMin:
        return "DPT_TimePeriodMin"
    case KnxDatapointType_DPT_TimePeriodMin_Z:
        return "DPT_TimePeriodMin_Z"
    case KnxDatapointType_DPT_TimePeriodMsec:
        return "DPT_TimePeriodMsec"
    case KnxDatapointType_DPT_TimePeriodMsec_Z:
        return "DPT_TimePeriodMsec_Z"
    case KnxDatapointType_DPT_TimePeriodSec:
        return "DPT_TimePeriodSec"
    case KnxDatapointType_DPT_TimePeriodSec_Z:
        return "DPT_TimePeriodSec_Z"
    case KnxDatapointType_DPT_Time_Delay:
        return "DPT_Time_Delay"
    case KnxDatapointType_DPT_Trigger:
        return "DPT_Trigger"
    case KnxDatapointType_DPT_UCountValue16_Z:
        return "DPT_UCountValue16_Z"
    case KnxDatapointType_DPT_UCountValue8_Z:
        return "DPT_UCountValue8_Z"
    case KnxDatapointType_DPT_UElCurrentmA:
        return "DPT_UElCurrentmA"
    case KnxDatapointType_DPT_UElCurrentyA_Z:
        return "DPT_UElCurrentyA_Z"
    case KnxDatapointType_DPT_UFlowRateLiter_h_Z:
        return "DPT_UFlowRateLiter_h_Z"
    case KnxDatapointType_DPT_UTF_8:
        return "DPT_UTF_8"
    case KnxDatapointType_DPT_UpDown:
        return "DPT_UpDown"
    case KnxDatapointType_DPT_UpDown_Action:
        return "DPT_UpDown_Action"
    case KnxDatapointType_DPT_ValueDemBOC:
        return "DPT_ValueDemBOC"
    case KnxDatapointType_DPT_Value_1_Count:
        return "DPT_Value_1_Count"
    case KnxDatapointType_DPT_Value_1_Ucount:
        return "DPT_Value_1_Ucount"
    case KnxDatapointType_DPT_Value_2_Count:
        return "DPT_Value_2_Count"
    case KnxDatapointType_DPT_Value_2_Ucount:
        return "DPT_Value_2_Ucount"
    case KnxDatapointType_DPT_Value_4_Count:
        return "DPT_Value_4_Count"
    case KnxDatapointType_DPT_Value_4_Ucount:
        return "DPT_Value_4_Ucount"
    case KnxDatapointType_DPT_Value_Absolute_Temperature:
        return "DPT_Value_Absolute_Temperature"
    case KnxDatapointType_DPT_Value_Acceleration:
        return "DPT_Value_Acceleration"
    case KnxDatapointType_DPT_Value_Acceleration_Angular:
        return "DPT_Value_Acceleration_Angular"
    case KnxDatapointType_DPT_Value_Activation_Energy:
        return "DPT_Value_Activation_Energy"
    case KnxDatapointType_DPT_Value_Activity:
        return "DPT_Value_Activity"
    case KnxDatapointType_DPT_Value_AirQuality:
        return "DPT_Value_AirQuality"
    case KnxDatapointType_DPT_Value_Amplitude:
        return "DPT_Value_Amplitude"
    case KnxDatapointType_DPT_Value_AngleDeg:
        return "DPT_Value_AngleDeg"
    case KnxDatapointType_DPT_Value_AngleRad:
        return "DPT_Value_AngleRad"
    case KnxDatapointType_DPT_Value_Angular_Frequency:
        return "DPT_Value_Angular_Frequency"
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
    case KnxDatapointType_DPT_Value_Common_Temperature:
        return "DPT_Value_Common_Temperature"
    case KnxDatapointType_DPT_Value_Compressibility:
        return "DPT_Value_Compressibility"
    case KnxDatapointType_DPT_Value_Conductance:
        return "DPT_Value_Conductance"
    case KnxDatapointType_DPT_Value_Curr:
        return "DPT_Value_Curr"
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
    case KnxDatapointType_DPT_Value_Electric_FluxDensity:
        return "DPT_Value_Electric_FluxDensity"
    case KnxDatapointType_DPT_Value_Electric_Polarization:
        return "DPT_Value_Electric_Polarization"
    case KnxDatapointType_DPT_Value_Electric_Potential:
        return "DPT_Value_Electric_Potential"
    case KnxDatapointType_DPT_Value_Electric_PotentialDifference:
        return "DPT_Value_Electric_PotentialDifference"
    case KnxDatapointType_DPT_Value_Electrical_Conductivity:
        return "DPT_Value_Electrical_Conductivity"
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
    case KnxDatapointType_DPT_Value_Heat_Capacity:
        return "DPT_Value_Heat_Capacity"
    case KnxDatapointType_DPT_Value_Heat_FlowRate:
        return "DPT_Value_Heat_FlowRate"
    case KnxDatapointType_DPT_Value_Heat_Quantity:
        return "DPT_Value_Heat_Quantity"
    case KnxDatapointType_DPT_Value_Humidity:
        return "DPT_Value_Humidity"
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
    case KnxDatapointType_DPT_Value_Lux:
        return "DPT_Value_Lux"
    case KnxDatapointType_DPT_Value_Magnetic_FieldStrength:
        return "DPT_Value_Magnetic_FieldStrength"
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
    case KnxDatapointType_DPT_Value_Mol:
        return "DPT_Value_Mol"
    case KnxDatapointType_DPT_Value_Momentum:
        return "DPT_Value_Momentum"
    case KnxDatapointType_DPT_Value_Phase_AngleDeg:
        return "DPT_Value_Phase_AngleDeg"
    case KnxDatapointType_DPT_Value_Phase_AngleRad:
        return "DPT_Value_Phase_AngleRad"
    case KnxDatapointType_DPT_Value_Power:
        return "DPT_Value_Power"
    case KnxDatapointType_DPT_Value_Power_Factor:
        return "DPT_Value_Power_Factor"
    case KnxDatapointType_DPT_Value_Pres:
        return "DPT_Value_Pres"
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
    case KnxDatapointType_DPT_Value_Speed:
        return "DPT_Value_Speed"
    case KnxDatapointType_DPT_Value_Stress:
        return "DPT_Value_Stress"
    case KnxDatapointType_DPT_Value_Surface_Tension:
        return "DPT_Value_Surface_Tension"
    case KnxDatapointType_DPT_Value_Temp:
        return "DPT_Value_Temp"
    case KnxDatapointType_DPT_Value_Temp_F:
        return "DPT_Value_Temp_F"
    case KnxDatapointType_DPT_Value_Tempa:
        return "DPT_Value_Tempa"
    case KnxDatapointType_DPT_Value_Tempd:
        return "DPT_Value_Tempd"
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
    case KnxDatapointType_DPT_Value_Time1:
        return "DPT_Value_Time1"
    case KnxDatapointType_DPT_Value_Time2:
        return "DPT_Value_Time2"
    case KnxDatapointType_DPT_Value_Torque:
        return "DPT_Value_Torque"
    case KnxDatapointType_DPT_Value_Volt:
        return "DPT_Value_Volt"
    case KnxDatapointType_DPT_Value_Volume:
        return "DPT_Value_Volume"
    case KnxDatapointType_DPT_Value_Volume_Flow:
        return "DPT_Value_Volume_Flow"
    case KnxDatapointType_DPT_Value_Volume_Flux:
        return "DPT_Value_Volume_Flux"
    case KnxDatapointType_DPT_Value_Weight:
        return "DPT_Value_Weight"
    case KnxDatapointType_DPT_Value_Work:
        return "DPT_Value_Work"
    case KnxDatapointType_DPT_Value_Wsp:
        return "DPT_Value_Wsp"
    case KnxDatapointType_DPT_Value_Wsp_kmh:
        return "DPT_Value_Wsp_kmh"
    case KnxDatapointType_DPT_ValveMode:
        return "DPT_ValveMode"
    case KnxDatapointType_DPT_VarString_8859_1:
        return "DPT_VarString_8859_1"
    case KnxDatapointType_DPT_Version:
        return "DPT_Version"
    case KnxDatapointType_DPT_VolumeLiter_Z:
        return "DPT_VolumeLiter_Z"
    case KnxDatapointType_DPT_WindSpeed_Z_DPT_WindSpeed:
        return "DPT_WindSpeed_Z_DPT_WindSpeed"
    case KnxDatapointType_DPT_Window_Door:
        return "DPT_Window_Door"
    }
    return ""
}
