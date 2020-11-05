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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

type KnxDatapointTypes string

type IKnxDatapointTypes interface {
    spi.Message
    FormatName() string
    MainNumber() uint16
    SubNumber() uint16
    Serialize(io utils.WriteBuffer) error
}

const(
    KnxDatapointTypes_DPT_Switch KnxDatapointTypes = "DPT_Switch"
    KnxDatapointTypes_DPT_Bool KnxDatapointTypes = "DPT_Bool"
    KnxDatapointTypes_DPT_Enable KnxDatapointTypes = "DPT_Enable"
    KnxDatapointTypes_DPT_Ramp KnxDatapointTypes = "DPT_Ramp"
    KnxDatapointTypes_DPT_Alarm KnxDatapointTypes = "DPT_Alarm"
    KnxDatapointTypes_DPT_BinaryValue KnxDatapointTypes = "DPT_BinaryValue"
    KnxDatapointTypes_DPT_Step KnxDatapointTypes = "DPT_Step"
    KnxDatapointTypes_DPT_UpDown KnxDatapointTypes = "DPT_UpDown"
    KnxDatapointTypes_DPT_OpenClose KnxDatapointTypes = "DPT_OpenClose"
    KnxDatapointTypes_DPT_Start KnxDatapointTypes = "DPT_Start"
    KnxDatapointTypes_DPT_State KnxDatapointTypes = "DPT_State"
    KnxDatapointTypes_DPT_Invert KnxDatapointTypes = "DPT_Invert"
    KnxDatapointTypes_DPT_DimSendStyle KnxDatapointTypes = "DPT_DimSendStyle"
    KnxDatapointTypes_DPT_InputSource KnxDatapointTypes = "DPT_InputSource"
    KnxDatapointTypes_DPT_Reset KnxDatapointTypes = "DPT_Reset"
    KnxDatapointTypes_DPT_Ack KnxDatapointTypes = "DPT_Ack"
    KnxDatapointTypes_DPT_Trigger KnxDatapointTypes = "DPT_Trigger"
    KnxDatapointTypes_DPT_Occupancy KnxDatapointTypes = "DPT_Occupancy"
    KnxDatapointTypes_DPT_Window_Door KnxDatapointTypes = "DPT_Window_Door"
    KnxDatapointTypes_DPT_LogicalFunction KnxDatapointTypes = "DPT_LogicalFunction"
    KnxDatapointTypes_DPT_Scene_AB KnxDatapointTypes = "DPT_Scene_AB"
    KnxDatapointTypes_DPT_ShutterBlinds_Mode KnxDatapointTypes = "DPT_ShutterBlinds_Mode"
    KnxDatapointTypes_DPT_Heat_Cool KnxDatapointTypes = "DPT_Heat_Cool"
    KnxDatapointTypes_DPT_Switch_Control KnxDatapointTypes = "DPT_Switch_Control"
    KnxDatapointTypes_DPT_Bool_Control KnxDatapointTypes = "DPT_Bool_Control"
    KnxDatapointTypes_DPT_Enable_Control KnxDatapointTypes = "DPT_Enable_Control"
    KnxDatapointTypes_DPT_Ramp_Control KnxDatapointTypes = "DPT_Ramp_Control"
    KnxDatapointTypes_DPT_Alarm_Control KnxDatapointTypes = "DPT_Alarm_Control"
    KnxDatapointTypes_DPT_BinaryValue_Control KnxDatapointTypes = "DPT_BinaryValue_Control"
    KnxDatapointTypes_DPT_Step_Control KnxDatapointTypes = "DPT_Step_Control"
    KnxDatapointTypes_DPT_Direction1_Control KnxDatapointTypes = "DPT_Direction1_Control"
    KnxDatapointTypes_DPT_Direction2_Control KnxDatapointTypes = "DPT_Direction2_Control"
    KnxDatapointTypes_DPT_Start_Control KnxDatapointTypes = "DPT_Start_Control"
    KnxDatapointTypes_DPT_State_Control KnxDatapointTypes = "DPT_State_Control"
    KnxDatapointTypes_DPT_Invert_Control KnxDatapointTypes = "DPT_Invert_Control"
    KnxDatapointTypes_DPT_Control_Dimming KnxDatapointTypes = "DPT_Control_Dimming"
    KnxDatapointTypes_DPT_Control_Blinds KnxDatapointTypes = "DPT_Control_Blinds"
    KnxDatapointTypes_DPT_Char_ASCII KnxDatapointTypes = "DPT_Char_ASCII"
    KnxDatapointTypes_DPT_Char_8859_1 KnxDatapointTypes = "DPT_Char_8859_1"
    KnxDatapointTypes_DPT_Scaling KnxDatapointTypes = "DPT_Scaling"
    KnxDatapointTypes_DPT_Angle KnxDatapointTypes = "DPT_Angle"
    KnxDatapointTypes_DPT_Percent_U8 KnxDatapointTypes = "DPT_Percent_U8"
    KnxDatapointTypes_DPT_DecimalFactor KnxDatapointTypes = "DPT_DecimalFactor"
    KnxDatapointTypes_DPT_Tariff KnxDatapointTypes = "DPT_Tariff"
    KnxDatapointTypes_DPT_Value_1_Ucount KnxDatapointTypes = "DPT_Value_1_Ucount"
    KnxDatapointTypes_DPT_Percent_V8 KnxDatapointTypes = "DPT_Percent_V8"
    KnxDatapointTypes_DPT_Value_1_Count KnxDatapointTypes = "DPT_Value_1_Count"
    KnxDatapointTypes_DPT_Status_Mode3 KnxDatapointTypes = "DPT_Status_Mode3"
    KnxDatapointTypes_DPT_Value_2_Ucount KnxDatapointTypes = "DPT_Value_2_Ucount"
    KnxDatapointTypes_DPT_TimePeriodMsec KnxDatapointTypes = "DPT_TimePeriodMsec"
    KnxDatapointTypes_DPT_TimePeriod10MSec KnxDatapointTypes = "DPT_TimePeriod10MSec"
    KnxDatapointTypes_DPT_TimePeriod100MSec KnxDatapointTypes = "DPT_TimePeriod100MSec"
    KnxDatapointTypes_DPT_TimePeriodSec KnxDatapointTypes = "DPT_TimePeriodSec"
    KnxDatapointTypes_DPT_TimePeriodMin KnxDatapointTypes = "DPT_TimePeriodMin"
    KnxDatapointTypes_DPT_TimePeriodHrs KnxDatapointTypes = "DPT_TimePeriodHrs"
    KnxDatapointTypes_DPT_PropDataType KnxDatapointTypes = "DPT_PropDataType"
    KnxDatapointTypes_DPT_Length_mm KnxDatapointTypes = "DPT_Length_mm"
    KnxDatapointTypes_DPT_UElCurrentmA KnxDatapointTypes = "DPT_UElCurrentmA"
    KnxDatapointTypes_DPT_Brightness KnxDatapointTypes = "DPT_Brightness"
    KnxDatapointTypes_DPT_Value_2_Count KnxDatapointTypes = "DPT_Value_2_Count"
    KnxDatapointTypes_DPT_DeltaTimeMsec KnxDatapointTypes = "DPT_DeltaTimeMsec"
    KnxDatapointTypes_DPT_DeltaTime10MSec KnxDatapointTypes = "DPT_DeltaTime10MSec"
    KnxDatapointTypes_DPT_DeltaTime100MSec KnxDatapointTypes = "DPT_DeltaTime100MSec"
    KnxDatapointTypes_DPT_DeltaTimeSec KnxDatapointTypes = "DPT_DeltaTimeSec"
    KnxDatapointTypes_DPT_DeltaTimeMin KnxDatapointTypes = "DPT_DeltaTimeMin"
    KnxDatapointTypes_DPT_DeltaTimeHrs KnxDatapointTypes = "DPT_DeltaTimeHrs"
    KnxDatapointTypes_DPT_Percent_V16 KnxDatapointTypes = "DPT_Percent_V16"
    KnxDatapointTypes_DPT_Rotation_Angle KnxDatapointTypes = "DPT_Rotation_Angle"
    KnxDatapointTypes_DPT_Value_Temp KnxDatapointTypes = "DPT_Value_Temp"
    KnxDatapointTypes_DPT_Value_Tempd KnxDatapointTypes = "DPT_Value_Tempd"
    KnxDatapointTypes_DPT_Value_Tempa KnxDatapointTypes = "DPT_Value_Tempa"
    KnxDatapointTypes_DPT_Value_Lux KnxDatapointTypes = "DPT_Value_Lux"
    KnxDatapointTypes_DPT_Value_Wsp KnxDatapointTypes = "DPT_Value_Wsp"
    KnxDatapointTypes_DPT_Value_Pres KnxDatapointTypes = "DPT_Value_Pres"
    KnxDatapointTypes_DPT_Value_Humidity KnxDatapointTypes = "DPT_Value_Humidity"
    KnxDatapointTypes_DPT_Value_AirQuality KnxDatapointTypes = "DPT_Value_AirQuality"
    KnxDatapointTypes_DPT_Value_Time1 KnxDatapointTypes = "DPT_Value_Time1"
    KnxDatapointTypes_DPT_Value_Time2 KnxDatapointTypes = "DPT_Value_Time2"
    KnxDatapointTypes_DPT_Value_Volt KnxDatapointTypes = "DPT_Value_Volt"
    KnxDatapointTypes_DPT_Value_Curr KnxDatapointTypes = "DPT_Value_Curr"
    KnxDatapointTypes_DPT_PowerDensity KnxDatapointTypes = "DPT_PowerDensity"
    KnxDatapointTypes_DPT_KelvinPerPercent KnxDatapointTypes = "DPT_KelvinPerPercent"
    KnxDatapointTypes_DPT_Power KnxDatapointTypes = "DPT_Power"
    KnxDatapointTypes_DPT_Value_Volume_Flow KnxDatapointTypes = "DPT_Value_Volume_Flow"
    KnxDatapointTypes_DPT_Rain_Amount KnxDatapointTypes = "DPT_Rain_Amount"
    KnxDatapointTypes_DPT_Value_Temp_F KnxDatapointTypes = "DPT_Value_Temp_F"
    KnxDatapointTypes_DPT_Value_Wsp_kmh KnxDatapointTypes = "DPT_Value_Wsp_kmh"
    KnxDatapointTypes_DPT_TimeOfDay KnxDatapointTypes = "DPT_TimeOfDay"
    KnxDatapointTypes_DPT_Date KnxDatapointTypes = "DPT_Date"
    KnxDatapointTypes_DPT_Value_4_Ucount KnxDatapointTypes = "DPT_Value_4_Ucount"
    KnxDatapointTypes_DPT_Value_4_Count KnxDatapointTypes = "DPT_Value_4_Count"
    KnxDatapointTypes_DPT_FlowRate_m3h KnxDatapointTypes = "DPT_FlowRate_m3h"
    KnxDatapointTypes_DPT_ActiveEnergy KnxDatapointTypes = "DPT_ActiveEnergy"
    KnxDatapointTypes_DPT_ApparantEnergy KnxDatapointTypes = "DPT_ApparantEnergy"
    KnxDatapointTypes_DPT_ReactiveEnergy KnxDatapointTypes = "DPT_ReactiveEnergy"
    KnxDatapointTypes_DPT_ActiveEnergy_kWh KnxDatapointTypes = "DPT_ActiveEnergy_kWh"
    KnxDatapointTypes_DPT_ApparantEnergy_kVAh KnxDatapointTypes = "DPT_ApparantEnergy_kVAh"
    KnxDatapointTypes_DPT_ReactiveEnergy_kVARh KnxDatapointTypes = "DPT_ReactiveEnergy_kVARh"
    KnxDatapointTypes_DPT_LongDeltaTimeSec KnxDatapointTypes = "DPT_LongDeltaTimeSec"
    KnxDatapointTypes_DPT_Value_Acceleration KnxDatapointTypes = "DPT_Value_Acceleration"
    KnxDatapointTypes_DPT_Value_Acceleration_Angular KnxDatapointTypes = "DPT_Value_Acceleration_Angular"
    KnxDatapointTypes_DPT_Value_Activation_Energy KnxDatapointTypes = "DPT_Value_Activation_Energy"
    KnxDatapointTypes_DPT_Value_Activity KnxDatapointTypes = "DPT_Value_Activity"
    KnxDatapointTypes_DPT_Value_Mol KnxDatapointTypes = "DPT_Value_Mol"
    KnxDatapointTypes_DPT_Value_Amplitude KnxDatapointTypes = "DPT_Value_Amplitude"
    KnxDatapointTypes_DPT_Value_AngleRad KnxDatapointTypes = "DPT_Value_AngleRad"
    KnxDatapointTypes_DPT_Value_AngleDeg KnxDatapointTypes = "DPT_Value_AngleDeg"
    KnxDatapointTypes_DPT_Value_Angular_Momentum KnxDatapointTypes = "DPT_Value_Angular_Momentum"
    KnxDatapointTypes_DPT_Value_Angular_Velocity KnxDatapointTypes = "DPT_Value_Angular_Velocity"
    KnxDatapointTypes_DPT_Value_Area KnxDatapointTypes = "DPT_Value_Area"
    KnxDatapointTypes_DPT_Value_Capacitance KnxDatapointTypes = "DPT_Value_Capacitance"
    KnxDatapointTypes_DPT_Value_Charge_DensitySurface KnxDatapointTypes = "DPT_Value_Charge_DensitySurface"
    KnxDatapointTypes_DPT_Value_Charge_DensityVolume KnxDatapointTypes = "DPT_Value_Charge_DensityVolume"
    KnxDatapointTypes_DPT_Value_Compressibility KnxDatapointTypes = "DPT_Value_Compressibility"
    KnxDatapointTypes_DPT_Value_Conductance KnxDatapointTypes = "DPT_Value_Conductance"
    KnxDatapointTypes_DPT_Value_Electrical_Conductivity KnxDatapointTypes = "DPT_Value_Electrical_Conductivity"
    KnxDatapointTypes_DPT_Value_Density KnxDatapointTypes = "DPT_Value_Density"
    KnxDatapointTypes_DPT_Value_Electric_Charge KnxDatapointTypes = "DPT_Value_Electric_Charge"
    KnxDatapointTypes_DPT_Value_Electric_Current KnxDatapointTypes = "DPT_Value_Electric_Current"
    KnxDatapointTypes_DPT_Value_Electric_CurrentDensity KnxDatapointTypes = "DPT_Value_Electric_CurrentDensity"
    KnxDatapointTypes_DPT_Value_Electric_DipoleMoment KnxDatapointTypes = "DPT_Value_Electric_DipoleMoment"
    KnxDatapointTypes_DPT_Value_Electric_Displacement KnxDatapointTypes = "DPT_Value_Electric_Displacement"
    KnxDatapointTypes_DPT_Value_Electric_FieldStrength KnxDatapointTypes = "DPT_Value_Electric_FieldStrength"
    KnxDatapointTypes_DPT_Value_Electric_Flux KnxDatapointTypes = "DPT_Value_Electric_Flux"
    KnxDatapointTypes_DPT_Value_Electric_FluxDensity KnxDatapointTypes = "DPT_Value_Electric_FluxDensity"
    KnxDatapointTypes_DPT_Value_Electric_Polarization KnxDatapointTypes = "DPT_Value_Electric_Polarization"
    KnxDatapointTypes_DPT_Value_Electric_Potential KnxDatapointTypes = "DPT_Value_Electric_Potential"
    KnxDatapointTypes_DPT_Value_Electric_PotentialDifference KnxDatapointTypes = "DPT_Value_Electric_PotentialDifference"
    KnxDatapointTypes_DPT_Value_ElectromagneticMoment KnxDatapointTypes = "DPT_Value_ElectromagneticMoment"
    KnxDatapointTypes_DPT_Value_Electromotive_Force KnxDatapointTypes = "DPT_Value_Electromotive_Force"
    KnxDatapointTypes_DPT_Value_Energy KnxDatapointTypes = "DPT_Value_Energy"
    KnxDatapointTypes_DPT_Value_Force KnxDatapointTypes = "DPT_Value_Force"
    KnxDatapointTypes_DPT_Value_Frequency KnxDatapointTypes = "DPT_Value_Frequency"
    KnxDatapointTypes_DPT_Value_Angular_Frequency KnxDatapointTypes = "DPT_Value_Angular_Frequency"
    KnxDatapointTypes_DPT_Value_Heat_Capacity KnxDatapointTypes = "DPT_Value_Heat_Capacity"
    KnxDatapointTypes_DPT_Value_Heat_FlowRate KnxDatapointTypes = "DPT_Value_Heat_FlowRate"
    KnxDatapointTypes_DPT_Value_Heat_Quantity KnxDatapointTypes = "DPT_Value_Heat_Quantity"
    KnxDatapointTypes_DPT_Value_Impedance KnxDatapointTypes = "DPT_Value_Impedance"
    KnxDatapointTypes_DPT_Value_Length KnxDatapointTypes = "DPT_Value_Length"
    KnxDatapointTypes_DPT_Value_Light_Quantity KnxDatapointTypes = "DPT_Value_Light_Quantity"
    KnxDatapointTypes_DPT_Value_Luminance KnxDatapointTypes = "DPT_Value_Luminance"
    KnxDatapointTypes_DPT_Value_Luminous_Flux KnxDatapointTypes = "DPT_Value_Luminous_Flux"
    KnxDatapointTypes_DPT_Value_Luminous_Intensity KnxDatapointTypes = "DPT_Value_Luminous_Intensity"
    KnxDatapointTypes_DPT_Value_Magnetic_FieldStrength KnxDatapointTypes = "DPT_Value_Magnetic_FieldStrength"
    KnxDatapointTypes_DPT_Value_Magnetic_Flux KnxDatapointTypes = "DPT_Value_Magnetic_Flux"
    KnxDatapointTypes_DPT_Value_Magnetic_FluxDensity KnxDatapointTypes = "DPT_Value_Magnetic_FluxDensity"
    KnxDatapointTypes_DPT_Value_Magnetic_Moment KnxDatapointTypes = "DPT_Value_Magnetic_Moment"
    KnxDatapointTypes_DPT_Value_Magnetic_Polarization KnxDatapointTypes = "DPT_Value_Magnetic_Polarization"
    KnxDatapointTypes_DPT_Value_Magnetization KnxDatapointTypes = "DPT_Value_Magnetization"
    KnxDatapointTypes_DPT_Value_MagnetomotiveForce KnxDatapointTypes = "DPT_Value_MagnetomotiveForce"
    KnxDatapointTypes_DPT_Value_Mass KnxDatapointTypes = "DPT_Value_Mass"
    KnxDatapointTypes_DPT_Value_MassFlux KnxDatapointTypes = "DPT_Value_MassFlux"
    KnxDatapointTypes_DPT_Value_Momentum KnxDatapointTypes = "DPT_Value_Momentum"
    KnxDatapointTypes_DPT_Value_Phase_AngleRad KnxDatapointTypes = "DPT_Value_Phase_AngleRad"
    KnxDatapointTypes_DPT_Value_Phase_AngleDeg KnxDatapointTypes = "DPT_Value_Phase_AngleDeg"
    KnxDatapointTypes_DPT_Value_Power KnxDatapointTypes = "DPT_Value_Power"
    KnxDatapointTypes_DPT_Value_Power_Factor KnxDatapointTypes = "DPT_Value_Power_Factor"
    KnxDatapointTypes_DPT_Value_Pressure KnxDatapointTypes = "DPT_Value_Pressure"
    KnxDatapointTypes_DPT_Value_Reactance KnxDatapointTypes = "DPT_Value_Reactance"
    KnxDatapointTypes_DPT_Value_Resistance KnxDatapointTypes = "DPT_Value_Resistance"
    KnxDatapointTypes_DPT_Value_Resistivity KnxDatapointTypes = "DPT_Value_Resistivity"
    KnxDatapointTypes_DPT_Value_SelfInductance KnxDatapointTypes = "DPT_Value_SelfInductance"
    KnxDatapointTypes_DPT_Value_SolidAngle KnxDatapointTypes = "DPT_Value_SolidAngle"
    KnxDatapointTypes_DPT_Value_Sound_Intensity KnxDatapointTypes = "DPT_Value_Sound_Intensity"
    KnxDatapointTypes_DPT_Value_Speed KnxDatapointTypes = "DPT_Value_Speed"
    KnxDatapointTypes_DPT_Value_Stress KnxDatapointTypes = "DPT_Value_Stress"
    KnxDatapointTypes_DPT_Value_Surface_Tension KnxDatapointTypes = "DPT_Value_Surface_Tension"
    KnxDatapointTypes_DPT_Value_Common_Temperature KnxDatapointTypes = "DPT_Value_Common_Temperature"
    KnxDatapointTypes_DPT_Value_Absolute_Temperature KnxDatapointTypes = "DPT_Value_Absolute_Temperature"
    KnxDatapointTypes_DPT_Value_TemperatureDifference KnxDatapointTypes = "DPT_Value_TemperatureDifference"
    KnxDatapointTypes_DPT_Value_Thermal_Capacity KnxDatapointTypes = "DPT_Value_Thermal_Capacity"
    KnxDatapointTypes_DPT_Value_Thermal_Conductivity KnxDatapointTypes = "DPT_Value_Thermal_Conductivity"
    KnxDatapointTypes_DPT_Value_ThermoelectricPower KnxDatapointTypes = "DPT_Value_ThermoelectricPower"
    KnxDatapointTypes_DPT_Value_Time KnxDatapointTypes = "DPT_Value_Time"
    KnxDatapointTypes_DPT_Value_Torque KnxDatapointTypes = "DPT_Value_Torque"
    KnxDatapointTypes_DPT_Value_Volume KnxDatapointTypes = "DPT_Value_Volume"
    KnxDatapointTypes_DPT_Value_Volume_Flux KnxDatapointTypes = "DPT_Value_Volume_Flux"
    KnxDatapointTypes_DPT_Value_Weight KnxDatapointTypes = "DPT_Value_Weight"
    KnxDatapointTypes_DPT_Value_Work KnxDatapointTypes = "DPT_Value_Work"
    KnxDatapointTypes_DPT_Access_Data KnxDatapointTypes = "DPT_Access_Data"
    KnxDatapointTypes_DPT_String_ASCII KnxDatapointTypes = "DPT_String_ASCII"
    KnxDatapointTypes_DPT_String_8859_1 KnxDatapointTypes = "DPT_String_8859_1"
    KnxDatapointTypes_DPT_SceneNumber KnxDatapointTypes = "DPT_SceneNumber"
    KnxDatapointTypes_DPT_SceneControl KnxDatapointTypes = "DPT_SceneControl"
    KnxDatapointTypes_DPT_DateTime KnxDatapointTypes = "DPT_DateTime"
    KnxDatapointTypes_DPT_SCLOMode KnxDatapointTypes = "DPT_SCLOMode"
    KnxDatapointTypes_DPT_BuildingMode KnxDatapointTypes = "DPT_BuildingMode"
    KnxDatapointTypes_DPT_OccMode KnxDatapointTypes = "DPT_OccMode"
    KnxDatapointTypes_DPT_Priority KnxDatapointTypes = "DPT_Priority"
    KnxDatapointTypes_DPT_LightApplicationMode KnxDatapointTypes = "DPT_LightApplicationMode"
    KnxDatapointTypes_DPT_ApplicationArea KnxDatapointTypes = "DPT_ApplicationArea"
    KnxDatapointTypes_DPT_AlarmClassType KnxDatapointTypes = "DPT_AlarmClassType"
    KnxDatapointTypes_DPT_PSUMode KnxDatapointTypes = "DPT_PSUMode"
    KnxDatapointTypes_DPT_ErrorClass_System KnxDatapointTypes = "DPT_ErrorClass_System"
    KnxDatapointTypes_DPT_ErrorClass_HVAC KnxDatapointTypes = "DPT_ErrorClass_HVAC"
    KnxDatapointTypes_DPT_Time_Delay KnxDatapointTypes = "DPT_Time_Delay"
    KnxDatapointTypes_DPT_Beaufort_Wind_Force_Scale KnxDatapointTypes = "DPT_Beaufort_Wind_Force_Scale"
    KnxDatapointTypes_DPT_SensorSelect KnxDatapointTypes = "DPT_SensorSelect"
    KnxDatapointTypes_DPT_ActuatorConnectType KnxDatapointTypes = "DPT_ActuatorConnectType"
    KnxDatapointTypes_DPT_FuelType KnxDatapointTypes = "DPT_FuelType"
    KnxDatapointTypes_DPT_BurnerType KnxDatapointTypes = "DPT_BurnerType"
    KnxDatapointTypes_DPT_HVACMode KnxDatapointTypes = "DPT_HVACMode"
    KnxDatapointTypes_DPT_DHWMode KnxDatapointTypes = "DPT_DHWMode"
    KnxDatapointTypes_DPT_LoadPriority KnxDatapointTypes = "DPT_LoadPriority"
    KnxDatapointTypes_DPT_HVACContrMode KnxDatapointTypes = "DPT_HVACContrMode"
    KnxDatapointTypes_DPT_HVACEmergMode KnxDatapointTypes = "DPT_HVACEmergMode"
    KnxDatapointTypes_DPT_ChangeoverMode KnxDatapointTypes = "DPT_ChangeoverMode"
    KnxDatapointTypes_DPT_ValveMode KnxDatapointTypes = "DPT_ValveMode"
    KnxDatapointTypes_DPT_DamperMode KnxDatapointTypes = "DPT_DamperMode"
    KnxDatapointTypes_DPT_HeaterMode KnxDatapointTypes = "DPT_HeaterMode"
    KnxDatapointTypes_DPT_FanMode KnxDatapointTypes = "DPT_FanMode"
    KnxDatapointTypes_DPT_MasterSlaveMode KnxDatapointTypes = "DPT_MasterSlaveMode"
    KnxDatapointTypes_DPT_StatusRoomSetp KnxDatapointTypes = "DPT_StatusRoomSetp"
    KnxDatapointTypes_DPT_ADAType KnxDatapointTypes = "DPT_ADAType"
    KnxDatapointTypes_DPT_BackupMode KnxDatapointTypes = "DPT_BackupMode"
    KnxDatapointTypes_DPT_StartSynchronization KnxDatapointTypes = "DPT_StartSynchronization"
    KnxDatapointTypes_DPT_Behaviour_Lock_Unlock KnxDatapointTypes = "DPT_Behaviour_Lock_Unlock"
    KnxDatapointTypes_DPT_Behaviour_Bus_Power_Up_Down KnxDatapointTypes = "DPT_Behaviour_Bus_Power_Up_Down"
    KnxDatapointTypes_DPT_DALI_Fade_Time KnxDatapointTypes = "DPT_DALI_Fade_Time"
    KnxDatapointTypes_DPT_BlinkingMode KnxDatapointTypes = "DPT_BlinkingMode"
    KnxDatapointTypes_DPT_LightControlMode KnxDatapointTypes = "DPT_LightControlMode"
    KnxDatapointTypes_DPT_SwitchPBModel KnxDatapointTypes = "DPT_SwitchPBModel"
    KnxDatapointTypes_DPT_PBAction KnxDatapointTypes = "DPT_PBAction"
    KnxDatapointTypes_DPT_DimmPBModel KnxDatapointTypes = "DPT_DimmPBModel"
    KnxDatapointTypes_DPT_SwitchOnMode KnxDatapointTypes = "DPT_SwitchOnMode"
    KnxDatapointTypes_DPT_LoadTypeSet KnxDatapointTypes = "DPT_LoadTypeSet"
    KnxDatapointTypes_DPT_LoadTypeDetected KnxDatapointTypes = "DPT_LoadTypeDetected"
    KnxDatapointTypes_DPT_SABExceptBehaviour KnxDatapointTypes = "DPT_SABExceptBehaviour"
    KnxDatapointTypes_DPT_SABBehaviour_Lock_Unlock KnxDatapointTypes = "DPT_SABBehaviour_Lock_Unlock"
    KnxDatapointTypes_DPT_SSSBMode KnxDatapointTypes = "DPT_SSSBMode"
    KnxDatapointTypes_DPT_BlindsControlMode KnxDatapointTypes = "DPT_BlindsControlMode"
    KnxDatapointTypes_DPT_CommMode KnxDatapointTypes = "DPT_CommMode"
    KnxDatapointTypes_DPT_AddInfoTypes KnxDatapointTypes = "DPT_AddInfoTypes"
    KnxDatapointTypes_DPT_RF_ModeSelect KnxDatapointTypes = "DPT_RF_ModeSelect"
    KnxDatapointTypes_DPT_RF_FilterSelect KnxDatapointTypes = "DPT_RF_FilterSelect"
    KnxDatapointTypes_DPT_StatusGen KnxDatapointTypes = "DPT_StatusGen"
    KnxDatapointTypes_DPT_Device_Control KnxDatapointTypes = "DPT_Device_Control"
    KnxDatapointTypes_DPT_ForceSign KnxDatapointTypes = "DPT_ForceSign"
    KnxDatapointTypes_DPT_ForceSignCool KnxDatapointTypes = "DPT_ForceSignCool"
    KnxDatapointTypes_DPT_StatusRHC KnxDatapointTypes = "DPT_StatusRHC"
    KnxDatapointTypes_DPT_StatusSDHWC KnxDatapointTypes = "DPT_StatusSDHWC"
    KnxDatapointTypes_DPT_FuelTypeSet KnxDatapointTypes = "DPT_FuelTypeSet"
    KnxDatapointTypes_DPT_StatusRCC KnxDatapointTypes = "DPT_StatusRCC"
    KnxDatapointTypes_DPT_StatusAHU KnxDatapointTypes = "DPT_StatusAHU"
    KnxDatapointTypes_DPT_LightActuatorErrorInfo KnxDatapointTypes = "DPT_LightActuatorErrorInfo"
    KnxDatapointTypes_DPT_RF_ModeInfo KnxDatapointTypes = "DPT_RF_ModeInfo"
    KnxDatapointTypes_DPT_RF_FilterInfo KnxDatapointTypes = "DPT_RF_FilterInfo"
    KnxDatapointTypes_DPT_Channel_Activation_8 KnxDatapointTypes = "DPT_Channel_Activation_8"
    KnxDatapointTypes_DPT_StatusDHWC KnxDatapointTypes = "DPT_StatusDHWC"
    KnxDatapointTypes_DPT_StatusRHCC KnxDatapointTypes = "DPT_StatusRHCC"
    KnxDatapointTypes_DPT_Media KnxDatapointTypes = "DPT_Media"
    KnxDatapointTypes_DPT_Channel_Activation_16 KnxDatapointTypes = "DPT_Channel_Activation_16"
    KnxDatapointTypes_DPT_OnOff_Action KnxDatapointTypes = "DPT_OnOff_Action"
    KnxDatapointTypes_DPT_Alarm_Reaction KnxDatapointTypes = "DPT_Alarm_Reaction"
    KnxDatapointTypes_DPT_UpDown_Action KnxDatapointTypes = "DPT_UpDown_Action"
    KnxDatapointTypes_DPT_HVAC_PB_Action KnxDatapointTypes = "DPT_HVAC_PB_Action"
    KnxDatapointTypes_DPT_VarString_8859_1 KnxDatapointTypes = "DPT_VarString_8859_1"
    KnxDatapointTypes_DPT_DoubleNibble KnxDatapointTypes = "DPT_DoubleNibble"
    KnxDatapointTypes_DPT_SceneInfo KnxDatapointTypes = "DPT_SceneInfo"
    KnxDatapointTypes_DPT_CombinedInfoOnOff KnxDatapointTypes = "DPT_CombinedInfoOnOff"
    KnxDatapointTypes_DPT_UTF_8 KnxDatapointTypes = "DPT_UTF_8"
    KnxDatapointTypes_DPT_ActiveEnergy_V64 KnxDatapointTypes = "DPT_ActiveEnergy_V64"
    KnxDatapointTypes_DPT_ApparantEnergy_V64 KnxDatapointTypes = "DPT_ApparantEnergy_V64"
    KnxDatapointTypes_DPT_ReactiveEnergy_V64 KnxDatapointTypes = "DPT_ReactiveEnergy_V64"
    KnxDatapointTypes_DPT_Channel_Activation_24 KnxDatapointTypes = "DPT_Channel_Activation_24"
    KnxDatapointTypes_DPT_PB_Action_HVAC_Extended KnxDatapointTypes = "DPT_PB_Action_HVAC_Extended"
    KnxDatapointTypes_DPT_Heat_Cool_Z KnxDatapointTypes = "DPT_Heat_Cool_Z"
    KnxDatapointTypes_DPT_BinaryValue_Z KnxDatapointTypes = "DPT_BinaryValue_Z"
    KnxDatapointTypes_DPT_HVACMode_Z KnxDatapointTypes = "DPT_HVACMode_Z"
    KnxDatapointTypes_DPT_DHWMode_Z KnxDatapointTypes = "DPT_DHWMode_Z"
    KnxDatapointTypes_DPT_HVACContrMode_Z KnxDatapointTypes = "DPT_HVACContrMode_Z"
    KnxDatapointTypes_DPT_EnablH_Cstage_Z_DPT_EnablH_CStage KnxDatapointTypes = "DPT_EnablH_Cstage_Z_DPT_EnablH_CStage"
    KnxDatapointTypes_DPT_BuildingMode_Z KnxDatapointTypes = "DPT_BuildingMode_Z"
    KnxDatapointTypes_DPT_OccMode_Z KnxDatapointTypes = "DPT_OccMode_Z"
    KnxDatapointTypes_DPT_HVACEmergMode_Z KnxDatapointTypes = "DPT_HVACEmergMode_Z"
    KnxDatapointTypes_DPT_RelValue_Z KnxDatapointTypes = "DPT_RelValue_Z"
    KnxDatapointTypes_DPT_UCountValue8_Z KnxDatapointTypes = "DPT_UCountValue8_Z"
    KnxDatapointTypes_DPT_TimePeriodMsec_Z KnxDatapointTypes = "DPT_TimePeriodMsec_Z"
    KnxDatapointTypes_DPT_TimePeriod10Msec_Z KnxDatapointTypes = "DPT_TimePeriod10Msec_Z"
    KnxDatapointTypes_DPT_TimePeriod100Msec_Z KnxDatapointTypes = "DPT_TimePeriod100Msec_Z"
    KnxDatapointTypes_DPT_TimePeriodSec_Z KnxDatapointTypes = "DPT_TimePeriodSec_Z"
    KnxDatapointTypes_DPT_TimePeriodMin_Z KnxDatapointTypes = "DPT_TimePeriodMin_Z"
    KnxDatapointTypes_DPT_TimePeriodHrs_Z KnxDatapointTypes = "DPT_TimePeriodHrs_Z"
    KnxDatapointTypes_DPT_UFlowRateLiter_h_Z KnxDatapointTypes = "DPT_UFlowRateLiter_h_Z"
    KnxDatapointTypes_DPT_UCountValue16_Z KnxDatapointTypes = "DPT_UCountValue16_Z"
    KnxDatapointTypes_DPT_UElCurrentyA_Z KnxDatapointTypes = "DPT_UElCurrentyA_Z"
    KnxDatapointTypes_DPT_PowerKW_Z KnxDatapointTypes = "DPT_PowerKW_Z"
    KnxDatapointTypes_DPT_AtmPressureAbs_Z KnxDatapointTypes = "DPT_AtmPressureAbs_Z"
    KnxDatapointTypes_DPT_PercentU16_Z KnxDatapointTypes = "DPT_PercentU16_Z"
    KnxDatapointTypes_DPT_HVACAirQual_Z KnxDatapointTypes = "DPT_HVACAirQual_Z"
    KnxDatapointTypes_DPT_WindSpeed_Z_DPT_WindSpeed KnxDatapointTypes = "DPT_WindSpeed_Z_DPT_WindSpeed"
    KnxDatapointTypes_DPT_SunIntensity_Z KnxDatapointTypes = "DPT_SunIntensity_Z"
    KnxDatapointTypes_DPT_HVACAirFlowAbs_Z KnxDatapointTypes = "DPT_HVACAirFlowAbs_Z"
    KnxDatapointTypes_DPT_RelSignedValue_Z KnxDatapointTypes = "DPT_RelSignedValue_Z"
    KnxDatapointTypes_DPT_DeltaTimeMsec_Z KnxDatapointTypes = "DPT_DeltaTimeMsec_Z"
    KnxDatapointTypes_DPT_DeltaTime10Msec_Z KnxDatapointTypes = "DPT_DeltaTime10Msec_Z"
    KnxDatapointTypes_DPT_DeltaTime100Msec_Z KnxDatapointTypes = "DPT_DeltaTime100Msec_Z"
    KnxDatapointTypes_DPT_DeltaTimeSec_Z KnxDatapointTypes = "DPT_DeltaTimeSec_Z"
    KnxDatapointTypes_DPT_DeltaTimeMin_Z KnxDatapointTypes = "DPT_DeltaTimeMin_Z"
    KnxDatapointTypes_DPT_DeltaTimeHrs_Z KnxDatapointTypes = "DPT_DeltaTimeHrs_Z"
    KnxDatapointTypes_DPT_Percent_V16_Z KnxDatapointTypes = "DPT_Percent_V16_Z"
    KnxDatapointTypes_DPT_TempHVACAbs_Z KnxDatapointTypes = "DPT_TempHVACAbs_Z"
    KnxDatapointTypes_DPT_TempHVACRel_Z KnxDatapointTypes = "DPT_TempHVACRel_Z"
    KnxDatapointTypes_DPT_HVACAirFlowRel_Z KnxDatapointTypes = "DPT_HVACAirFlowRel_Z"
    KnxDatapointTypes_DPT_HVACModeNext KnxDatapointTypes = "DPT_HVACModeNext"
    KnxDatapointTypes_DPT_DHWModeNext KnxDatapointTypes = "DPT_DHWModeNext"
    KnxDatapointTypes_DPT_OccModeNext KnxDatapointTypes = "DPT_OccModeNext"
    KnxDatapointTypes_DPT_BuildingModeNext KnxDatapointTypes = "DPT_BuildingModeNext"
    KnxDatapointTypes_DPT_StatusBUC KnxDatapointTypes = "DPT_StatusBUC"
    KnxDatapointTypes_DPT_LockSign KnxDatapointTypes = "DPT_LockSign"
    KnxDatapointTypes_DPT_ValueDemBOC KnxDatapointTypes = "DPT_ValueDemBOC"
    KnxDatapointTypes_DPT_ActPosDemAbs KnxDatapointTypes = "DPT_ActPosDemAbs"
    KnxDatapointTypes_DPT_StatusAct KnxDatapointTypes = "DPT_StatusAct"
    KnxDatapointTypes_DPT_StatusLightingActuator KnxDatapointTypes = "DPT_StatusLightingActuator"
    KnxDatapointTypes_DPT_StatusHPM KnxDatapointTypes = "DPT_StatusHPM"
    KnxDatapointTypes_DPT_TempRoomDemAbs KnxDatapointTypes = "DPT_TempRoomDemAbs"
    KnxDatapointTypes_DPT_StatusCPM KnxDatapointTypes = "DPT_StatusCPM"
    KnxDatapointTypes_DPT_StatusWTC KnxDatapointTypes = "DPT_StatusWTC"
    KnxDatapointTypes_DPT_TempFlowWaterDemAbs KnxDatapointTypes = "DPT_TempFlowWaterDemAbs"
    KnxDatapointTypes_DPT_EnergyDemWater KnxDatapointTypes = "DPT_EnergyDemWater"
    KnxDatapointTypes_DPT_TempRoomSetpSetShift3 KnxDatapointTypes = "DPT_TempRoomSetpSetShift3"
    KnxDatapointTypes_DPT_TempRoomSetpSet3 KnxDatapointTypes = "DPT_TempRoomSetpSet3"
    KnxDatapointTypes_DPT_TempRoomSetpSet4 KnxDatapointTypes = "DPT_TempRoomSetpSet4"
    KnxDatapointTypes_DPT_TempDHWSetpSet4 KnxDatapointTypes = "DPT_TempDHWSetpSet4"
    KnxDatapointTypes_DPT_TempRoomSetpSetShift4 KnxDatapointTypes = "DPT_TempRoomSetpSetShift4"
    KnxDatapointTypes_DPT_PowerFlowWaterDemHPM KnxDatapointTypes = "DPT_PowerFlowWaterDemHPM"
    KnxDatapointTypes_DPT_PowerFlowWaterDemCPM KnxDatapointTypes = "DPT_PowerFlowWaterDemCPM"
    KnxDatapointTypes_DPT_StatusBOC KnxDatapointTypes = "DPT_StatusBOC"
    KnxDatapointTypes_DPT_StatusCC KnxDatapointTypes = "DPT_StatusCC"
    KnxDatapointTypes_DPT_SpecHeatProd KnxDatapointTypes = "DPT_SpecHeatProd"
    KnxDatapointTypes_DPT_Version KnxDatapointTypes = "DPT_Version"
    KnxDatapointTypes_DPT_VolumeLiter_Z KnxDatapointTypes = "DPT_VolumeLiter_Z"
    KnxDatapointTypes_DPT_FlowRate_m3h_Z KnxDatapointTypes = "DPT_FlowRate_m3h_Z"
    KnxDatapointTypes_DPT_AlarmInfo KnxDatapointTypes = "DPT_AlarmInfo"
    KnxDatapointTypes_DPT_TempHVACAbsNext KnxDatapointTypes = "DPT_TempHVACAbsNext"
    KnxDatapointTypes_DPT_SerNum KnxDatapointTypes = "DPT_SerNum"
    KnxDatapointTypes_DPT_TempRoomSetpSetF163 KnxDatapointTypes = "DPT_TempRoomSetpSetF163"
    KnxDatapointTypes_DPT_TempRoomSetpSetShiftF163 KnxDatapointTypes = "DPT_TempRoomSetpSetShiftF163"
    KnxDatapointTypes_DPT_EnergyDemAir KnxDatapointTypes = "DPT_EnergyDemAir"
    KnxDatapointTypes_DPT_TempSupply_AirSetpSet KnxDatapointTypes = "DPT_TempSupply_AirSetpSet"
    KnxDatapointTypes_DPT_ScalingSpeed KnxDatapointTypes = "DPT_ScalingSpeed"
    KnxDatapointTypes_DPT_Scaling_Step_Time KnxDatapointTypes = "DPT_Scaling_Step_Time"
    KnxDatapointTypes_DPT_TariffNext KnxDatapointTypes = "DPT_TariffNext"
    KnxDatapointTypes_DPT_MeteringValue KnxDatapointTypes = "DPT_MeteringValue"
    KnxDatapointTypes_DPT_MBus_Address KnxDatapointTypes = "DPT_MBus_Address"
    KnxDatapointTypes_DPT_Locale_ASCII KnxDatapointTypes = "DPT_Locale_ASCII"
    KnxDatapointTypes_DPT_Colour_RGB KnxDatapointTypes = "DPT_Colour_RGB"
    KnxDatapointTypes_DPT_LanguageCodeAlpha2_ASCII KnxDatapointTypes = "DPT_LanguageCodeAlpha2_ASCII"
    KnxDatapointTypes_DPT_RegionCodeAlpha2_ASCII KnxDatapointTypes = "DPT_RegionCodeAlpha2_ASCII"
    KnxDatapointTypes_DPT_Tariff_ActiveEnergy KnxDatapointTypes = "DPT_Tariff_ActiveEnergy"
    KnxDatapointTypes_DPT_Prioritised_Mode_Control KnxDatapointTypes = "DPT_Prioritised_Mode_Control"
    KnxDatapointTypes_DPT_DALI_Control_Gear_Diagnostic KnxDatapointTypes = "DPT_DALI_Control_Gear_Diagnostic"
    KnxDatapointTypes_DPT_SceneConfig KnxDatapointTypes = "DPT_SceneConfig"
    KnxDatapointTypes_DPT_DALI_Diagnostics KnxDatapointTypes = "DPT_DALI_Diagnostics"
    KnxDatapointTypes_DPT_FlaggedScaling KnxDatapointTypes = "DPT_FlaggedScaling"
    KnxDatapointTypes_DPT_CombinedPosition KnxDatapointTypes = "DPT_CombinedPosition"
    KnxDatapointTypes_DPT_StatusSAB KnxDatapointTypes = "DPT_StatusSAB"
)


func (e KnxDatapointTypes) FormatName() string {
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

func (e KnxDatapointTypes) MainNumber() uint16 {
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

func (e KnxDatapointTypes) SubNumber() uint16 {
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
func KnxDatapointTypesValueOf(value string) KnxDatapointTypes {
    switch value {
        case "DPT_ADAType":
            return KnxDatapointTypes_DPT_ADAType
        case "DPT_Access_Data":
            return KnxDatapointTypes_DPT_Access_Data
        case "DPT_Ack":
            return KnxDatapointTypes_DPT_Ack
        case "DPT_ActPosDemAbs":
            return KnxDatapointTypes_DPT_ActPosDemAbs
        case "DPT_ActiveEnergy":
            return KnxDatapointTypes_DPT_ActiveEnergy
        case "DPT_ActiveEnergy_V64":
            return KnxDatapointTypes_DPT_ActiveEnergy_V64
        case "DPT_ActiveEnergy_kWh":
            return KnxDatapointTypes_DPT_ActiveEnergy_kWh
        case "DPT_ActuatorConnectType":
            return KnxDatapointTypes_DPT_ActuatorConnectType
        case "DPT_AddInfoTypes":
            return KnxDatapointTypes_DPT_AddInfoTypes
        case "DPT_Alarm":
            return KnxDatapointTypes_DPT_Alarm
        case "DPT_AlarmClassType":
            return KnxDatapointTypes_DPT_AlarmClassType
        case "DPT_AlarmInfo":
            return KnxDatapointTypes_DPT_AlarmInfo
        case "DPT_Alarm_Control":
            return KnxDatapointTypes_DPT_Alarm_Control
        case "DPT_Alarm_Reaction":
            return KnxDatapointTypes_DPT_Alarm_Reaction
        case "DPT_Angle":
            return KnxDatapointTypes_DPT_Angle
        case "DPT_ApparantEnergy":
            return KnxDatapointTypes_DPT_ApparantEnergy
        case "DPT_ApparantEnergy_V64":
            return KnxDatapointTypes_DPT_ApparantEnergy_V64
        case "DPT_ApparantEnergy_kVAh":
            return KnxDatapointTypes_DPT_ApparantEnergy_kVAh
        case "DPT_ApplicationArea":
            return KnxDatapointTypes_DPT_ApplicationArea
        case "DPT_AtmPressureAbs_Z":
            return KnxDatapointTypes_DPT_AtmPressureAbs_Z
        case "DPT_BackupMode":
            return KnxDatapointTypes_DPT_BackupMode
        case "DPT_Beaufort_Wind_Force_Scale":
            return KnxDatapointTypes_DPT_Beaufort_Wind_Force_Scale
        case "DPT_Behaviour_Bus_Power_Up_Down":
            return KnxDatapointTypes_DPT_Behaviour_Bus_Power_Up_Down
        case "DPT_Behaviour_Lock_Unlock":
            return KnxDatapointTypes_DPT_Behaviour_Lock_Unlock
        case "DPT_BinaryValue":
            return KnxDatapointTypes_DPT_BinaryValue
        case "DPT_BinaryValue_Control":
            return KnxDatapointTypes_DPT_BinaryValue_Control
        case "DPT_BinaryValue_Z":
            return KnxDatapointTypes_DPT_BinaryValue_Z
        case "DPT_BlindsControlMode":
            return KnxDatapointTypes_DPT_BlindsControlMode
        case "DPT_BlinkingMode":
            return KnxDatapointTypes_DPT_BlinkingMode
        case "DPT_Bool":
            return KnxDatapointTypes_DPT_Bool
        case "DPT_Bool_Control":
            return KnxDatapointTypes_DPT_Bool_Control
        case "DPT_Brightness":
            return KnxDatapointTypes_DPT_Brightness
        case "DPT_BuildingMode":
            return KnxDatapointTypes_DPT_BuildingMode
        case "DPT_BuildingModeNext":
            return KnxDatapointTypes_DPT_BuildingModeNext
        case "DPT_BuildingMode_Z":
            return KnxDatapointTypes_DPT_BuildingMode_Z
        case "DPT_BurnerType":
            return KnxDatapointTypes_DPT_BurnerType
        case "DPT_ChangeoverMode":
            return KnxDatapointTypes_DPT_ChangeoverMode
        case "DPT_Channel_Activation_16":
            return KnxDatapointTypes_DPT_Channel_Activation_16
        case "DPT_Channel_Activation_24":
            return KnxDatapointTypes_DPT_Channel_Activation_24
        case "DPT_Channel_Activation_8":
            return KnxDatapointTypes_DPT_Channel_Activation_8
        case "DPT_Char_8859_1":
            return KnxDatapointTypes_DPT_Char_8859_1
        case "DPT_Char_ASCII":
            return KnxDatapointTypes_DPT_Char_ASCII
        case "DPT_Colour_RGB":
            return KnxDatapointTypes_DPT_Colour_RGB
        case "DPT_CombinedInfoOnOff":
            return KnxDatapointTypes_DPT_CombinedInfoOnOff
        case "DPT_CombinedPosition":
            return KnxDatapointTypes_DPT_CombinedPosition
        case "DPT_CommMode":
            return KnxDatapointTypes_DPT_CommMode
        case "DPT_Control_Blinds":
            return KnxDatapointTypes_DPT_Control_Blinds
        case "DPT_Control_Dimming":
            return KnxDatapointTypes_DPT_Control_Dimming
        case "DPT_DALI_Control_Gear_Diagnostic":
            return KnxDatapointTypes_DPT_DALI_Control_Gear_Diagnostic
        case "DPT_DALI_Diagnostics":
            return KnxDatapointTypes_DPT_DALI_Diagnostics
        case "DPT_DALI_Fade_Time":
            return KnxDatapointTypes_DPT_DALI_Fade_Time
        case "DPT_DHWMode":
            return KnxDatapointTypes_DPT_DHWMode
        case "DPT_DHWModeNext":
            return KnxDatapointTypes_DPT_DHWModeNext
        case "DPT_DHWMode_Z":
            return KnxDatapointTypes_DPT_DHWMode_Z
        case "DPT_DamperMode":
            return KnxDatapointTypes_DPT_DamperMode
        case "DPT_Date":
            return KnxDatapointTypes_DPT_Date
        case "DPT_DateTime":
            return KnxDatapointTypes_DPT_DateTime
        case "DPT_DecimalFactor":
            return KnxDatapointTypes_DPT_DecimalFactor
        case "DPT_DeltaTime100MSec":
            return KnxDatapointTypes_DPT_DeltaTime100MSec
        case "DPT_DeltaTime100Msec_Z":
            return KnxDatapointTypes_DPT_DeltaTime100Msec_Z
        case "DPT_DeltaTime10MSec":
            return KnxDatapointTypes_DPT_DeltaTime10MSec
        case "DPT_DeltaTime10Msec_Z":
            return KnxDatapointTypes_DPT_DeltaTime10Msec_Z
        case "DPT_DeltaTimeHrs":
            return KnxDatapointTypes_DPT_DeltaTimeHrs
        case "DPT_DeltaTimeHrs_Z":
            return KnxDatapointTypes_DPT_DeltaTimeHrs_Z
        case "DPT_DeltaTimeMin":
            return KnxDatapointTypes_DPT_DeltaTimeMin
        case "DPT_DeltaTimeMin_Z":
            return KnxDatapointTypes_DPT_DeltaTimeMin_Z
        case "DPT_DeltaTimeMsec":
            return KnxDatapointTypes_DPT_DeltaTimeMsec
        case "DPT_DeltaTimeMsec_Z":
            return KnxDatapointTypes_DPT_DeltaTimeMsec_Z
        case "DPT_DeltaTimeSec":
            return KnxDatapointTypes_DPT_DeltaTimeSec
        case "DPT_DeltaTimeSec_Z":
            return KnxDatapointTypes_DPT_DeltaTimeSec_Z
        case "DPT_Device_Control":
            return KnxDatapointTypes_DPT_Device_Control
        case "DPT_DimSendStyle":
            return KnxDatapointTypes_DPT_DimSendStyle
        case "DPT_DimmPBModel":
            return KnxDatapointTypes_DPT_DimmPBModel
        case "DPT_Direction1_Control":
            return KnxDatapointTypes_DPT_Direction1_Control
        case "DPT_Direction2_Control":
            return KnxDatapointTypes_DPT_Direction2_Control
        case "DPT_DoubleNibble":
            return KnxDatapointTypes_DPT_DoubleNibble
        case "DPT_EnablH_Cstage_Z_DPT_EnablH_CStage":
            return KnxDatapointTypes_DPT_EnablH_Cstage_Z_DPT_EnablH_CStage
        case "DPT_Enable":
            return KnxDatapointTypes_DPT_Enable
        case "DPT_Enable_Control":
            return KnxDatapointTypes_DPT_Enable_Control
        case "DPT_EnergyDemAir":
            return KnxDatapointTypes_DPT_EnergyDemAir
        case "DPT_EnergyDemWater":
            return KnxDatapointTypes_DPT_EnergyDemWater
        case "DPT_ErrorClass_HVAC":
            return KnxDatapointTypes_DPT_ErrorClass_HVAC
        case "DPT_ErrorClass_System":
            return KnxDatapointTypes_DPT_ErrorClass_System
        case "DPT_FanMode":
            return KnxDatapointTypes_DPT_FanMode
        case "DPT_FlaggedScaling":
            return KnxDatapointTypes_DPT_FlaggedScaling
        case "DPT_FlowRate_m3h":
            return KnxDatapointTypes_DPT_FlowRate_m3h
        case "DPT_FlowRate_m3h_Z":
            return KnxDatapointTypes_DPT_FlowRate_m3h_Z
        case "DPT_ForceSign":
            return KnxDatapointTypes_DPT_ForceSign
        case "DPT_ForceSignCool":
            return KnxDatapointTypes_DPT_ForceSignCool
        case "DPT_FuelType":
            return KnxDatapointTypes_DPT_FuelType
        case "DPT_FuelTypeSet":
            return KnxDatapointTypes_DPT_FuelTypeSet
        case "DPT_HVACAirFlowAbs_Z":
            return KnxDatapointTypes_DPT_HVACAirFlowAbs_Z
        case "DPT_HVACAirFlowRel_Z":
            return KnxDatapointTypes_DPT_HVACAirFlowRel_Z
        case "DPT_HVACAirQual_Z":
            return KnxDatapointTypes_DPT_HVACAirQual_Z
        case "DPT_HVACContrMode":
            return KnxDatapointTypes_DPT_HVACContrMode
        case "DPT_HVACContrMode_Z":
            return KnxDatapointTypes_DPT_HVACContrMode_Z
        case "DPT_HVACEmergMode":
            return KnxDatapointTypes_DPT_HVACEmergMode
        case "DPT_HVACEmergMode_Z":
            return KnxDatapointTypes_DPT_HVACEmergMode_Z
        case "DPT_HVACMode":
            return KnxDatapointTypes_DPT_HVACMode
        case "DPT_HVACModeNext":
            return KnxDatapointTypes_DPT_HVACModeNext
        case "DPT_HVACMode_Z":
            return KnxDatapointTypes_DPT_HVACMode_Z
        case "DPT_HVAC_PB_Action":
            return KnxDatapointTypes_DPT_HVAC_PB_Action
        case "DPT_Heat_Cool":
            return KnxDatapointTypes_DPT_Heat_Cool
        case "DPT_Heat_Cool_Z":
            return KnxDatapointTypes_DPT_Heat_Cool_Z
        case "DPT_HeaterMode":
            return KnxDatapointTypes_DPT_HeaterMode
        case "DPT_InputSource":
            return KnxDatapointTypes_DPT_InputSource
        case "DPT_Invert":
            return KnxDatapointTypes_DPT_Invert
        case "DPT_Invert_Control":
            return KnxDatapointTypes_DPT_Invert_Control
        case "DPT_KelvinPerPercent":
            return KnxDatapointTypes_DPT_KelvinPerPercent
        case "DPT_LanguageCodeAlpha2_ASCII":
            return KnxDatapointTypes_DPT_LanguageCodeAlpha2_ASCII
        case "DPT_Length_mm":
            return KnxDatapointTypes_DPT_Length_mm
        case "DPT_LightActuatorErrorInfo":
            return KnxDatapointTypes_DPT_LightActuatorErrorInfo
        case "DPT_LightApplicationMode":
            return KnxDatapointTypes_DPT_LightApplicationMode
        case "DPT_LightControlMode":
            return KnxDatapointTypes_DPT_LightControlMode
        case "DPT_LoadPriority":
            return KnxDatapointTypes_DPT_LoadPriority
        case "DPT_LoadTypeDetected":
            return KnxDatapointTypes_DPT_LoadTypeDetected
        case "DPT_LoadTypeSet":
            return KnxDatapointTypes_DPT_LoadTypeSet
        case "DPT_Locale_ASCII":
            return KnxDatapointTypes_DPT_Locale_ASCII
        case "DPT_LockSign":
            return KnxDatapointTypes_DPT_LockSign
        case "DPT_LogicalFunction":
            return KnxDatapointTypes_DPT_LogicalFunction
        case "DPT_LongDeltaTimeSec":
            return KnxDatapointTypes_DPT_LongDeltaTimeSec
        case "DPT_MBus_Address":
            return KnxDatapointTypes_DPT_MBus_Address
        case "DPT_MasterSlaveMode":
            return KnxDatapointTypes_DPT_MasterSlaveMode
        case "DPT_Media":
            return KnxDatapointTypes_DPT_Media
        case "DPT_MeteringValue":
            return KnxDatapointTypes_DPT_MeteringValue
        case "DPT_OccMode":
            return KnxDatapointTypes_DPT_OccMode
        case "DPT_OccModeNext":
            return KnxDatapointTypes_DPT_OccModeNext
        case "DPT_OccMode_Z":
            return KnxDatapointTypes_DPT_OccMode_Z
        case "DPT_Occupancy":
            return KnxDatapointTypes_DPT_Occupancy
        case "DPT_OnOff_Action":
            return KnxDatapointTypes_DPT_OnOff_Action
        case "DPT_OpenClose":
            return KnxDatapointTypes_DPT_OpenClose
        case "DPT_PBAction":
            return KnxDatapointTypes_DPT_PBAction
        case "DPT_PB_Action_HVAC_Extended":
            return KnxDatapointTypes_DPT_PB_Action_HVAC_Extended
        case "DPT_PSUMode":
            return KnxDatapointTypes_DPT_PSUMode
        case "DPT_PercentU16_Z":
            return KnxDatapointTypes_DPT_PercentU16_Z
        case "DPT_Percent_U8":
            return KnxDatapointTypes_DPT_Percent_U8
        case "DPT_Percent_V16":
            return KnxDatapointTypes_DPT_Percent_V16
        case "DPT_Percent_V16_Z":
            return KnxDatapointTypes_DPT_Percent_V16_Z
        case "DPT_Percent_V8":
            return KnxDatapointTypes_DPT_Percent_V8
        case "DPT_Power":
            return KnxDatapointTypes_DPT_Power
        case "DPT_PowerDensity":
            return KnxDatapointTypes_DPT_PowerDensity
        case "DPT_PowerFlowWaterDemCPM":
            return KnxDatapointTypes_DPT_PowerFlowWaterDemCPM
        case "DPT_PowerFlowWaterDemHPM":
            return KnxDatapointTypes_DPT_PowerFlowWaterDemHPM
        case "DPT_PowerKW_Z":
            return KnxDatapointTypes_DPT_PowerKW_Z
        case "DPT_Prioritised_Mode_Control":
            return KnxDatapointTypes_DPT_Prioritised_Mode_Control
        case "DPT_Priority":
            return KnxDatapointTypes_DPT_Priority
        case "DPT_PropDataType":
            return KnxDatapointTypes_DPT_PropDataType
        case "DPT_RF_FilterInfo":
            return KnxDatapointTypes_DPT_RF_FilterInfo
        case "DPT_RF_FilterSelect":
            return KnxDatapointTypes_DPT_RF_FilterSelect
        case "DPT_RF_ModeInfo":
            return KnxDatapointTypes_DPT_RF_ModeInfo
        case "DPT_RF_ModeSelect":
            return KnxDatapointTypes_DPT_RF_ModeSelect
        case "DPT_Rain_Amount":
            return KnxDatapointTypes_DPT_Rain_Amount
        case "DPT_Ramp":
            return KnxDatapointTypes_DPT_Ramp
        case "DPT_Ramp_Control":
            return KnxDatapointTypes_DPT_Ramp_Control
        case "DPT_ReactiveEnergy":
            return KnxDatapointTypes_DPT_ReactiveEnergy
        case "DPT_ReactiveEnergy_V64":
            return KnxDatapointTypes_DPT_ReactiveEnergy_V64
        case "DPT_ReactiveEnergy_kVARh":
            return KnxDatapointTypes_DPT_ReactiveEnergy_kVARh
        case "DPT_RegionCodeAlpha2_ASCII":
            return KnxDatapointTypes_DPT_RegionCodeAlpha2_ASCII
        case "DPT_RelSignedValue_Z":
            return KnxDatapointTypes_DPT_RelSignedValue_Z
        case "DPT_RelValue_Z":
            return KnxDatapointTypes_DPT_RelValue_Z
        case "DPT_Reset":
            return KnxDatapointTypes_DPT_Reset
        case "DPT_Rotation_Angle":
            return KnxDatapointTypes_DPT_Rotation_Angle
        case "DPT_SABBehaviour_Lock_Unlock":
            return KnxDatapointTypes_DPT_SABBehaviour_Lock_Unlock
        case "DPT_SABExceptBehaviour":
            return KnxDatapointTypes_DPT_SABExceptBehaviour
        case "DPT_SCLOMode":
            return KnxDatapointTypes_DPT_SCLOMode
        case "DPT_SSSBMode":
            return KnxDatapointTypes_DPT_SSSBMode
        case "DPT_Scaling":
            return KnxDatapointTypes_DPT_Scaling
        case "DPT_ScalingSpeed":
            return KnxDatapointTypes_DPT_ScalingSpeed
        case "DPT_Scaling_Step_Time":
            return KnxDatapointTypes_DPT_Scaling_Step_Time
        case "DPT_SceneConfig":
            return KnxDatapointTypes_DPT_SceneConfig
        case "DPT_SceneControl":
            return KnxDatapointTypes_DPT_SceneControl
        case "DPT_SceneInfo":
            return KnxDatapointTypes_DPT_SceneInfo
        case "DPT_SceneNumber":
            return KnxDatapointTypes_DPT_SceneNumber
        case "DPT_Scene_AB":
            return KnxDatapointTypes_DPT_Scene_AB
        case "DPT_SensorSelect":
            return KnxDatapointTypes_DPT_SensorSelect
        case "DPT_SerNum":
            return KnxDatapointTypes_DPT_SerNum
        case "DPT_ShutterBlinds_Mode":
            return KnxDatapointTypes_DPT_ShutterBlinds_Mode
        case "DPT_SpecHeatProd":
            return KnxDatapointTypes_DPT_SpecHeatProd
        case "DPT_Start":
            return KnxDatapointTypes_DPT_Start
        case "DPT_StartSynchronization":
            return KnxDatapointTypes_DPT_StartSynchronization
        case "DPT_Start_Control":
            return KnxDatapointTypes_DPT_Start_Control
        case "DPT_State":
            return KnxDatapointTypes_DPT_State
        case "DPT_State_Control":
            return KnxDatapointTypes_DPT_State_Control
        case "DPT_StatusAHU":
            return KnxDatapointTypes_DPT_StatusAHU
        case "DPT_StatusAct":
            return KnxDatapointTypes_DPT_StatusAct
        case "DPT_StatusBOC":
            return KnxDatapointTypes_DPT_StatusBOC
        case "DPT_StatusBUC":
            return KnxDatapointTypes_DPT_StatusBUC
        case "DPT_StatusCC":
            return KnxDatapointTypes_DPT_StatusCC
        case "DPT_StatusCPM":
            return KnxDatapointTypes_DPT_StatusCPM
        case "DPT_StatusDHWC":
            return KnxDatapointTypes_DPT_StatusDHWC
        case "DPT_StatusGen":
            return KnxDatapointTypes_DPT_StatusGen
        case "DPT_StatusHPM":
            return KnxDatapointTypes_DPT_StatusHPM
        case "DPT_StatusLightingActuator":
            return KnxDatapointTypes_DPT_StatusLightingActuator
        case "DPT_StatusRCC":
            return KnxDatapointTypes_DPT_StatusRCC
        case "DPT_StatusRHC":
            return KnxDatapointTypes_DPT_StatusRHC
        case "DPT_StatusRHCC":
            return KnxDatapointTypes_DPT_StatusRHCC
        case "DPT_StatusRoomSetp":
            return KnxDatapointTypes_DPT_StatusRoomSetp
        case "DPT_StatusSAB":
            return KnxDatapointTypes_DPT_StatusSAB
        case "DPT_StatusSDHWC":
            return KnxDatapointTypes_DPT_StatusSDHWC
        case "DPT_StatusWTC":
            return KnxDatapointTypes_DPT_StatusWTC
        case "DPT_Status_Mode3":
            return KnxDatapointTypes_DPT_Status_Mode3
        case "DPT_Step":
            return KnxDatapointTypes_DPT_Step
        case "DPT_Step_Control":
            return KnxDatapointTypes_DPT_Step_Control
        case "DPT_String_8859_1":
            return KnxDatapointTypes_DPT_String_8859_1
        case "DPT_String_ASCII":
            return KnxDatapointTypes_DPT_String_ASCII
        case "DPT_SunIntensity_Z":
            return KnxDatapointTypes_DPT_SunIntensity_Z
        case "DPT_Switch":
            return KnxDatapointTypes_DPT_Switch
        case "DPT_SwitchOnMode":
            return KnxDatapointTypes_DPT_SwitchOnMode
        case "DPT_SwitchPBModel":
            return KnxDatapointTypes_DPT_SwitchPBModel
        case "DPT_Switch_Control":
            return KnxDatapointTypes_DPT_Switch_Control
        case "DPT_Tariff":
            return KnxDatapointTypes_DPT_Tariff
        case "DPT_TariffNext":
            return KnxDatapointTypes_DPT_TariffNext
        case "DPT_Tariff_ActiveEnergy":
            return KnxDatapointTypes_DPT_Tariff_ActiveEnergy
        case "DPT_TempDHWSetpSet4":
            return KnxDatapointTypes_DPT_TempDHWSetpSet4
        case "DPT_TempFlowWaterDemAbs":
            return KnxDatapointTypes_DPT_TempFlowWaterDemAbs
        case "DPT_TempHVACAbsNext":
            return KnxDatapointTypes_DPT_TempHVACAbsNext
        case "DPT_TempHVACAbs_Z":
            return KnxDatapointTypes_DPT_TempHVACAbs_Z
        case "DPT_TempHVACRel_Z":
            return KnxDatapointTypes_DPT_TempHVACRel_Z
        case "DPT_TempRoomDemAbs":
            return KnxDatapointTypes_DPT_TempRoomDemAbs
        case "DPT_TempRoomSetpSet3":
            return KnxDatapointTypes_DPT_TempRoomSetpSet3
        case "DPT_TempRoomSetpSet4":
            return KnxDatapointTypes_DPT_TempRoomSetpSet4
        case "DPT_TempRoomSetpSetF163":
            return KnxDatapointTypes_DPT_TempRoomSetpSetF163
        case "DPT_TempRoomSetpSetShift3":
            return KnxDatapointTypes_DPT_TempRoomSetpSetShift3
        case "DPT_TempRoomSetpSetShift4":
            return KnxDatapointTypes_DPT_TempRoomSetpSetShift4
        case "DPT_TempRoomSetpSetShiftF163":
            return KnxDatapointTypes_DPT_TempRoomSetpSetShiftF163
        case "DPT_TempSupply_AirSetpSet":
            return KnxDatapointTypes_DPT_TempSupply_AirSetpSet
        case "DPT_TimeOfDay":
            return KnxDatapointTypes_DPT_TimeOfDay
        case "DPT_TimePeriod100MSec":
            return KnxDatapointTypes_DPT_TimePeriod100MSec
        case "DPT_TimePeriod100Msec_Z":
            return KnxDatapointTypes_DPT_TimePeriod100Msec_Z
        case "DPT_TimePeriod10MSec":
            return KnxDatapointTypes_DPT_TimePeriod10MSec
        case "DPT_TimePeriod10Msec_Z":
            return KnxDatapointTypes_DPT_TimePeriod10Msec_Z
        case "DPT_TimePeriodHrs":
            return KnxDatapointTypes_DPT_TimePeriodHrs
        case "DPT_TimePeriodHrs_Z":
            return KnxDatapointTypes_DPT_TimePeriodHrs_Z
        case "DPT_TimePeriodMin":
            return KnxDatapointTypes_DPT_TimePeriodMin
        case "DPT_TimePeriodMin_Z":
            return KnxDatapointTypes_DPT_TimePeriodMin_Z
        case "DPT_TimePeriodMsec":
            return KnxDatapointTypes_DPT_TimePeriodMsec
        case "DPT_TimePeriodMsec_Z":
            return KnxDatapointTypes_DPT_TimePeriodMsec_Z
        case "DPT_TimePeriodSec":
            return KnxDatapointTypes_DPT_TimePeriodSec
        case "DPT_TimePeriodSec_Z":
            return KnxDatapointTypes_DPT_TimePeriodSec_Z
        case "DPT_Time_Delay":
            return KnxDatapointTypes_DPT_Time_Delay
        case "DPT_Trigger":
            return KnxDatapointTypes_DPT_Trigger
        case "DPT_UCountValue16_Z":
            return KnxDatapointTypes_DPT_UCountValue16_Z
        case "DPT_UCountValue8_Z":
            return KnxDatapointTypes_DPT_UCountValue8_Z
        case "DPT_UElCurrentmA":
            return KnxDatapointTypes_DPT_UElCurrentmA
        case "DPT_UElCurrentyA_Z":
            return KnxDatapointTypes_DPT_UElCurrentyA_Z
        case "DPT_UFlowRateLiter_h_Z":
            return KnxDatapointTypes_DPT_UFlowRateLiter_h_Z
        case "DPT_UTF_8":
            return KnxDatapointTypes_DPT_UTF_8
        case "DPT_UpDown":
            return KnxDatapointTypes_DPT_UpDown
        case "DPT_UpDown_Action":
            return KnxDatapointTypes_DPT_UpDown_Action
        case "DPT_ValueDemBOC":
            return KnxDatapointTypes_DPT_ValueDemBOC
        case "DPT_Value_1_Count":
            return KnxDatapointTypes_DPT_Value_1_Count
        case "DPT_Value_1_Ucount":
            return KnxDatapointTypes_DPT_Value_1_Ucount
        case "DPT_Value_2_Count":
            return KnxDatapointTypes_DPT_Value_2_Count
        case "DPT_Value_2_Ucount":
            return KnxDatapointTypes_DPT_Value_2_Ucount
        case "DPT_Value_4_Count":
            return KnxDatapointTypes_DPT_Value_4_Count
        case "DPT_Value_4_Ucount":
            return KnxDatapointTypes_DPT_Value_4_Ucount
        case "DPT_Value_Absolute_Temperature":
            return KnxDatapointTypes_DPT_Value_Absolute_Temperature
        case "DPT_Value_Acceleration":
            return KnxDatapointTypes_DPT_Value_Acceleration
        case "DPT_Value_Acceleration_Angular":
            return KnxDatapointTypes_DPT_Value_Acceleration_Angular
        case "DPT_Value_Activation_Energy":
            return KnxDatapointTypes_DPT_Value_Activation_Energy
        case "DPT_Value_Activity":
            return KnxDatapointTypes_DPT_Value_Activity
        case "DPT_Value_AirQuality":
            return KnxDatapointTypes_DPT_Value_AirQuality
        case "DPT_Value_Amplitude":
            return KnxDatapointTypes_DPT_Value_Amplitude
        case "DPT_Value_AngleDeg":
            return KnxDatapointTypes_DPT_Value_AngleDeg
        case "DPT_Value_AngleRad":
            return KnxDatapointTypes_DPT_Value_AngleRad
        case "DPT_Value_Angular_Frequency":
            return KnxDatapointTypes_DPT_Value_Angular_Frequency
        case "DPT_Value_Angular_Momentum":
            return KnxDatapointTypes_DPT_Value_Angular_Momentum
        case "DPT_Value_Angular_Velocity":
            return KnxDatapointTypes_DPT_Value_Angular_Velocity
        case "DPT_Value_Area":
            return KnxDatapointTypes_DPT_Value_Area
        case "DPT_Value_Capacitance":
            return KnxDatapointTypes_DPT_Value_Capacitance
        case "DPT_Value_Charge_DensitySurface":
            return KnxDatapointTypes_DPT_Value_Charge_DensitySurface
        case "DPT_Value_Charge_DensityVolume":
            return KnxDatapointTypes_DPT_Value_Charge_DensityVolume
        case "DPT_Value_Common_Temperature":
            return KnxDatapointTypes_DPT_Value_Common_Temperature
        case "DPT_Value_Compressibility":
            return KnxDatapointTypes_DPT_Value_Compressibility
        case "DPT_Value_Conductance":
            return KnxDatapointTypes_DPT_Value_Conductance
        case "DPT_Value_Curr":
            return KnxDatapointTypes_DPT_Value_Curr
        case "DPT_Value_Density":
            return KnxDatapointTypes_DPT_Value_Density
        case "DPT_Value_Electric_Charge":
            return KnxDatapointTypes_DPT_Value_Electric_Charge
        case "DPT_Value_Electric_Current":
            return KnxDatapointTypes_DPT_Value_Electric_Current
        case "DPT_Value_Electric_CurrentDensity":
            return KnxDatapointTypes_DPT_Value_Electric_CurrentDensity
        case "DPT_Value_Electric_DipoleMoment":
            return KnxDatapointTypes_DPT_Value_Electric_DipoleMoment
        case "DPT_Value_Electric_Displacement":
            return KnxDatapointTypes_DPT_Value_Electric_Displacement
        case "DPT_Value_Electric_FieldStrength":
            return KnxDatapointTypes_DPT_Value_Electric_FieldStrength
        case "DPT_Value_Electric_Flux":
            return KnxDatapointTypes_DPT_Value_Electric_Flux
        case "DPT_Value_Electric_FluxDensity":
            return KnxDatapointTypes_DPT_Value_Electric_FluxDensity
        case "DPT_Value_Electric_Polarization":
            return KnxDatapointTypes_DPT_Value_Electric_Polarization
        case "DPT_Value_Electric_Potential":
            return KnxDatapointTypes_DPT_Value_Electric_Potential
        case "DPT_Value_Electric_PotentialDifference":
            return KnxDatapointTypes_DPT_Value_Electric_PotentialDifference
        case "DPT_Value_Electrical_Conductivity":
            return KnxDatapointTypes_DPT_Value_Electrical_Conductivity
        case "DPT_Value_ElectromagneticMoment":
            return KnxDatapointTypes_DPT_Value_ElectromagneticMoment
        case "DPT_Value_Electromotive_Force":
            return KnxDatapointTypes_DPT_Value_Electromotive_Force
        case "DPT_Value_Energy":
            return KnxDatapointTypes_DPT_Value_Energy
        case "DPT_Value_Force":
            return KnxDatapointTypes_DPT_Value_Force
        case "DPT_Value_Frequency":
            return KnxDatapointTypes_DPT_Value_Frequency
        case "DPT_Value_Heat_Capacity":
            return KnxDatapointTypes_DPT_Value_Heat_Capacity
        case "DPT_Value_Heat_FlowRate":
            return KnxDatapointTypes_DPT_Value_Heat_FlowRate
        case "DPT_Value_Heat_Quantity":
            return KnxDatapointTypes_DPT_Value_Heat_Quantity
        case "DPT_Value_Humidity":
            return KnxDatapointTypes_DPT_Value_Humidity
        case "DPT_Value_Impedance":
            return KnxDatapointTypes_DPT_Value_Impedance
        case "DPT_Value_Length":
            return KnxDatapointTypes_DPT_Value_Length
        case "DPT_Value_Light_Quantity":
            return KnxDatapointTypes_DPT_Value_Light_Quantity
        case "DPT_Value_Luminance":
            return KnxDatapointTypes_DPT_Value_Luminance
        case "DPT_Value_Luminous_Flux":
            return KnxDatapointTypes_DPT_Value_Luminous_Flux
        case "DPT_Value_Luminous_Intensity":
            return KnxDatapointTypes_DPT_Value_Luminous_Intensity
        case "DPT_Value_Lux":
            return KnxDatapointTypes_DPT_Value_Lux
        case "DPT_Value_Magnetic_FieldStrength":
            return KnxDatapointTypes_DPT_Value_Magnetic_FieldStrength
        case "DPT_Value_Magnetic_Flux":
            return KnxDatapointTypes_DPT_Value_Magnetic_Flux
        case "DPT_Value_Magnetic_FluxDensity":
            return KnxDatapointTypes_DPT_Value_Magnetic_FluxDensity
        case "DPT_Value_Magnetic_Moment":
            return KnxDatapointTypes_DPT_Value_Magnetic_Moment
        case "DPT_Value_Magnetic_Polarization":
            return KnxDatapointTypes_DPT_Value_Magnetic_Polarization
        case "DPT_Value_Magnetization":
            return KnxDatapointTypes_DPT_Value_Magnetization
        case "DPT_Value_MagnetomotiveForce":
            return KnxDatapointTypes_DPT_Value_MagnetomotiveForce
        case "DPT_Value_Mass":
            return KnxDatapointTypes_DPT_Value_Mass
        case "DPT_Value_MassFlux":
            return KnxDatapointTypes_DPT_Value_MassFlux
        case "DPT_Value_Mol":
            return KnxDatapointTypes_DPT_Value_Mol
        case "DPT_Value_Momentum":
            return KnxDatapointTypes_DPT_Value_Momentum
        case "DPT_Value_Phase_AngleDeg":
            return KnxDatapointTypes_DPT_Value_Phase_AngleDeg
        case "DPT_Value_Phase_AngleRad":
            return KnxDatapointTypes_DPT_Value_Phase_AngleRad
        case "DPT_Value_Power":
            return KnxDatapointTypes_DPT_Value_Power
        case "DPT_Value_Power_Factor":
            return KnxDatapointTypes_DPT_Value_Power_Factor
        case "DPT_Value_Pres":
            return KnxDatapointTypes_DPT_Value_Pres
        case "DPT_Value_Pressure":
            return KnxDatapointTypes_DPT_Value_Pressure
        case "DPT_Value_Reactance":
            return KnxDatapointTypes_DPT_Value_Reactance
        case "DPT_Value_Resistance":
            return KnxDatapointTypes_DPT_Value_Resistance
        case "DPT_Value_Resistivity":
            return KnxDatapointTypes_DPT_Value_Resistivity
        case "DPT_Value_SelfInductance":
            return KnxDatapointTypes_DPT_Value_SelfInductance
        case "DPT_Value_SolidAngle":
            return KnxDatapointTypes_DPT_Value_SolidAngle
        case "DPT_Value_Sound_Intensity":
            return KnxDatapointTypes_DPT_Value_Sound_Intensity
        case "DPT_Value_Speed":
            return KnxDatapointTypes_DPT_Value_Speed
        case "DPT_Value_Stress":
            return KnxDatapointTypes_DPT_Value_Stress
        case "DPT_Value_Surface_Tension":
            return KnxDatapointTypes_DPT_Value_Surface_Tension
        case "DPT_Value_Temp":
            return KnxDatapointTypes_DPT_Value_Temp
        case "DPT_Value_Temp_F":
            return KnxDatapointTypes_DPT_Value_Temp_F
        case "DPT_Value_Tempa":
            return KnxDatapointTypes_DPT_Value_Tempa
        case "DPT_Value_Tempd":
            return KnxDatapointTypes_DPT_Value_Tempd
        case "DPT_Value_TemperatureDifference":
            return KnxDatapointTypes_DPT_Value_TemperatureDifference
        case "DPT_Value_Thermal_Capacity":
            return KnxDatapointTypes_DPT_Value_Thermal_Capacity
        case "DPT_Value_Thermal_Conductivity":
            return KnxDatapointTypes_DPT_Value_Thermal_Conductivity
        case "DPT_Value_ThermoelectricPower":
            return KnxDatapointTypes_DPT_Value_ThermoelectricPower
        case "DPT_Value_Time":
            return KnxDatapointTypes_DPT_Value_Time
        case "DPT_Value_Time1":
            return KnxDatapointTypes_DPT_Value_Time1
        case "DPT_Value_Time2":
            return KnxDatapointTypes_DPT_Value_Time2
        case "DPT_Value_Torque":
            return KnxDatapointTypes_DPT_Value_Torque
        case "DPT_Value_Volt":
            return KnxDatapointTypes_DPT_Value_Volt
        case "DPT_Value_Volume":
            return KnxDatapointTypes_DPT_Value_Volume
        case "DPT_Value_Volume_Flow":
            return KnxDatapointTypes_DPT_Value_Volume_Flow
        case "DPT_Value_Volume_Flux":
            return KnxDatapointTypes_DPT_Value_Volume_Flux
        case "DPT_Value_Weight":
            return KnxDatapointTypes_DPT_Value_Weight
        case "DPT_Value_Work":
            return KnxDatapointTypes_DPT_Value_Work
        case "DPT_Value_Wsp":
            return KnxDatapointTypes_DPT_Value_Wsp
        case "DPT_Value_Wsp_kmh":
            return KnxDatapointTypes_DPT_Value_Wsp_kmh
        case "DPT_ValveMode":
            return KnxDatapointTypes_DPT_ValveMode
        case "DPT_VarString_8859_1":
            return KnxDatapointTypes_DPT_VarString_8859_1
        case "DPT_Version":
            return KnxDatapointTypes_DPT_Version
        case "DPT_VolumeLiter_Z":
            return KnxDatapointTypes_DPT_VolumeLiter_Z
        case "DPT_WindSpeed_Z_DPT_WindSpeed":
            return KnxDatapointTypes_DPT_WindSpeed_Z_DPT_WindSpeed
        case "DPT_Window_Door":
            return KnxDatapointTypes_DPT_Window_Door
    }
    return ""
}

func CastKnxDatapointTypes(structType interface{}) KnxDatapointTypes {
    castFunc := func(typ interface{}) KnxDatapointTypes {
        if sKnxDatapointTypes, ok := typ.(KnxDatapointTypes); ok {
            return sKnxDatapointTypes
        }
        return ""
    }
    return castFunc(structType)
}

func (m KnxDatapointTypes) LengthInBits() uint16 {
    return 0
}

func (m KnxDatapointTypes) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KnxDatapointTypesParse(io *utils.ReadBuffer) (KnxDatapointTypes, error) {
    // TODO: Implement ...
    return "", nil
}

func (e KnxDatapointTypes) Serialize(io utils.WriteBuffer) error {
    // TODO: Implement ...
    return nil
}
