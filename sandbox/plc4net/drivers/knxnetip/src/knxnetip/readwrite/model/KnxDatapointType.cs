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

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{

    public enum KnxDatapointType
    {

        DPT_Switch,
        DPT_Bool,
        DPT_Enable,
        DPT_Ramp,
        DPT_Alarm,
        DPT_BinaryValue,
        DPT_Step,
        DPT_UpDown,
        DPT_OpenClose,
        DPT_Start,
        DPT_State,
        DPT_Invert,
        DPT_DimSendStyle,
        DPT_InputSource,
        DPT_Reset,
        DPT_Ack,
        DPT_Trigger,
        DPT_Occupancy,
        DPT_Window_Door,
        DPT_LogicalFunction,
        DPT_Scene_AB,
        DPT_ShutterBlinds_Mode,
        DPT_Heat_Cool,
        DPT_Switch_Control,
        DPT_Bool_Control,
        DPT_Enable_Control,
        DPT_Ramp_Control,
        DPT_Alarm_Control,
        DPT_BinaryValue_Control,
        DPT_Step_Control,
        DPT_Direction1_Control,
        DPT_Direction2_Control,
        DPT_Start_Control,
        DPT_State_Control,
        DPT_Invert_Control,
        DPT_Control_Dimming,
        DPT_Control_Blinds,
        DPT_Char_ASCII,
        DPT_Char_8859_1,
        DPT_Scaling,
        DPT_Angle,
        DPT_Percent_U8,
        DPT_DecimalFactor,
        DPT_Tariff,
        DPT_Value_1_Ucount,
        DPT_Percent_V8,
        DPT_Value_1_Count,
        DPT_Status_Mode3,
        DPT_Value_2_Ucount,
        DPT_TimePeriodMsec,
        DPT_TimePeriod10MSec,
        DPT_TimePeriod100MSec,
        DPT_TimePeriodSec,
        DPT_TimePeriodMin,
        DPT_TimePeriodHrs,
        DPT_PropDataType,
        DPT_Length_mm,
        DPT_UElCurrentmA,
        DPT_Brightness,
        DPT_Value_2_Count,
        DPT_DeltaTimeMsec,
        DPT_DeltaTime10MSec,
        DPT_DeltaTime100MSec,
        DPT_DeltaTimeSec,
        DPT_DeltaTimeMin,
        DPT_DeltaTimeHrs,
        DPT_Percent_V16,
        DPT_Rotation_Angle,
        DPT_Value_Temp,
        DPT_Value_Tempd,
        DPT_Value_Tempa,
        DPT_Value_Lux,
        DPT_Value_Wsp,
        DPT_Value_Pres,
        DPT_Value_Humidity,
        DPT_Value_AirQuality,
        DPT_Value_Time1,
        DPT_Value_Time2,
        DPT_Value_Volt,
        DPT_Value_Curr,
        DPT_PowerDensity,
        DPT_KelvinPerPercent,
        DPT_Power,
        DPT_Value_Volume_Flow,
        DPT_Rain_Amount,
        DPT_Value_Temp_F,
        DPT_Value_Wsp_kmh,
        DPT_TimeOfDay,
        DPT_Date,
        DPT_Value_4_Ucount,
        DPT_Value_4_Count,
        DPT_FlowRate_m3h,
        DPT_ActiveEnergy,
        DPT_ApparantEnergy,
        DPT_ReactiveEnergy,
        DPT_ActiveEnergy_kWh,
        DPT_ApparantEnergy_kVAh,
        DPT_ReactiveEnergy_kVARh,
        DPT_LongDeltaTimeSec,
        DPT_Value_Acceleration,
        DPT_Value_Acceleration_Angular,
        DPT_Value_Activation_Energy,
        DPT_Value_Activity,
        DPT_Value_Mol,
        DPT_Value_Amplitude,
        DPT_Value_AngleRad,
        DPT_Value_AngleDeg,
        DPT_Value_Angular_Momentum,
        DPT_Value_Angular_Velocity,
        DPT_Value_Area,
        DPT_Value_Capacitance,
        DPT_Value_Charge_DensitySurface,
        DPT_Value_Charge_DensityVolume,
        DPT_Value_Compressibility,
        DPT_Value_Conductance,
        DPT_Value_Electrical_Conductivity,
        DPT_Value_Density,
        DPT_Value_Electric_Charge,
        DPT_Value_Electric_Current,
        DPT_Value_Electric_CurrentDensity,
        DPT_Value_Electric_DipoleMoment,
        DPT_Value_Electric_Displacement,
        DPT_Value_Electric_FieldStrength,
        DPT_Value_Electric_Flux,
        DPT_Value_Electric_FluxDensity,
        DPT_Value_Electric_Polarization,
        DPT_Value_Electric_Potential,
        DPT_Value_Electric_PotentialDifference,
        DPT_Value_ElectromagneticMoment,
        DPT_Value_Electromotive_Force,
        DPT_Value_Energy,
        DPT_Value_Force,
        DPT_Value_Frequency,
        DPT_Value_Angular_Frequency,
        DPT_Value_Heat_Capacity,
        DPT_Value_Heat_FlowRate,
        DPT_Value_Heat_Quantity,
        DPT_Value_Impedance,
        DPT_Value_Length,
        DPT_Value_Light_Quantity,
        DPT_Value_Luminance,
        DPT_Value_Luminous_Flux,
        DPT_Value_Luminous_Intensity,
        DPT_Value_Magnetic_FieldStrength,
        DPT_Value_Magnetic_Flux,
        DPT_Value_Magnetic_FluxDensity,
        DPT_Value_Magnetic_Moment,
        DPT_Value_Magnetic_Polarization,
        DPT_Value_Magnetization,
        DPT_Value_MagnetomotiveForce,
        DPT_Value_Mass,
        DPT_Value_MassFlux,
        DPT_Value_Momentum,
        DPT_Value_Phase_AngleRad,
        DPT_Value_Phase_AngleDeg,
        DPT_Value_Power,
        DPT_Value_Power_Factor,
        DPT_Value_Pressure,
        DPT_Value_Reactance,
        DPT_Value_Resistance,
        DPT_Value_Resistivity,
        DPT_Value_SelfInductance,
        DPT_Value_SolidAngle,
        DPT_Value_Sound_Intensity,
        DPT_Value_Speed,
        DPT_Value_Stress,
        DPT_Value_Surface_Tension,
        DPT_Value_Common_Temperature,
        DPT_Value_Absolute_Temperature,
        DPT_Value_TemperatureDifference,
        DPT_Value_Thermal_Capacity,
        DPT_Value_Thermal_Conductivity,
        DPT_Value_ThermoelectricPower,
        DPT_Value_Time,
        DPT_Value_Torque,
        DPT_Value_Volume,
        DPT_Value_Volume_Flux,
        DPT_Value_Weight,
        DPT_Value_Work,
        DPT_Access_Data,
        DPT_String_ASCII,
        DPT_String_8859_1,
        DPT_SceneNumber,
        DPT_SceneControl,
        DPT_DateTime,
        DPT_SCLOMode,
        DPT_BuildingMode,
        DPT_OccMode,
        DPT_Priority,
        DPT_LightApplicationMode,
        DPT_ApplicationArea,
        DPT_AlarmClassType,
        DPT_PSUMode,
        DPT_ErrorClass_System,
        DPT_ErrorClass_HVAC,
        DPT_Time_Delay,
        DPT_Beaufort_Wind_Force_Scale,
        DPT_SensorSelect,
        DPT_ActuatorConnectType,
        DPT_FuelType,
        DPT_BurnerType,
        DPT_HVACMode,
        DPT_DHWMode,
        DPT_LoadPriority,
        DPT_HVACContrMode,
        DPT_HVACEmergMode,
        DPT_ChangeoverMode,
        DPT_ValveMode,
        DPT_DamperMode,
        DPT_HeaterMode,
        DPT_FanMode,
        DPT_MasterSlaveMode,
        DPT_StatusRoomSetp,
        DPT_ADAType,
        DPT_BackupMode,
        DPT_StartSynchronization,
        DPT_Behaviour_Lock_Unlock,
        DPT_Behaviour_Bus_Power_Up_Down,
        DPT_DALI_Fade_Time,
        DPT_BlinkingMode,
        DPT_LightControlMode,
        DPT_SwitchPBModel,
        DPT_PBAction,
        DPT_DimmPBModel,
        DPT_SwitchOnMode,
        DPT_LoadTypeSet,
        DPT_LoadTypeDetected,
        DPT_SABExceptBehaviour,
        DPT_SABBehaviour_Lock_Unlock,
        DPT_SSSBMode,
        DPT_BlindsControlMode,
        DPT_CommMode,
        DPT_AddInfoTypes,
        DPT_RF_ModeSelect,
        DPT_RF_FilterSelect,
        DPT_StatusGen,
        DPT_Device_Control,
        DPT_ForceSign,
        DPT_ForceSignCool,
        DPT_StatusRHC,
        DPT_StatusSDHWC,
        DPT_FuelTypeSet,
        DPT_StatusRCC,
        DPT_StatusAHU,
        DPT_LightActuatorErrorInfo,
        DPT_RF_ModeInfo,
        DPT_RF_FilterInfo,
        DPT_Channel_Activation_8,
        DPT_StatusDHWC,
        DPT_StatusRHCC,
        DPT_Media,
        DPT_Channel_Activation_16,
        DPT_OnOff_Action,
        DPT_Alarm_Reaction,
        DPT_UpDown_Action,
        DPT_HVAC_PB_Action,
        DPT_VarString_8859_1,
        DPT_DoubleNibble,
        DPT_SceneInfo,
        DPT_CombinedInfoOnOff,
        DPT_UTF_8,
        DPT_ActiveEnergy_V64,
        DPT_ApparantEnergy_V64,
        DPT_ReactiveEnergy_V64,
        DPT_Channel_Activation_24,
        DPT_PB_Action_HVAC_Extended,
        DPT_Heat_Cool_Z,
        DPT_BinaryValue_Z,
        DPT_HVACMode_Z,
        DPT_DHWMode_Z,
        DPT_HVACContrMode_Z,
        DPT_EnablH_Cstage_Z_DPT_EnablH_CStage,
        DPT_BuildingMode_Z,
        DPT_OccMode_Z,
        DPT_HVACEmergMode_Z,
        DPT_RelValue_Z,
        DPT_UCountValue8_Z,
        DPT_TimePeriodMsec_Z,
        DPT_TimePeriod10Msec_Z,
        DPT_TimePeriod100Msec_Z,
        DPT_TimePeriodSec_Z,
        DPT_TimePeriodMin_Z,
        DPT_TimePeriodHrs_Z,
        DPT_UFlowRateLiter_h_Z,
        DPT_UCountValue16_Z,
        DPT_UElCurrentyA_Z,
        DPT_PowerKW_Z,
        DPT_AtmPressureAbs_Z,
        DPT_PercentU16_Z,
        DPT_HVACAirQual_Z,
        DPT_WindSpeed_Z_DPT_WindSpeed,
        DPT_SunIntensity_Z,
        DPT_HVACAirFlowAbs_Z,
        DPT_RelSignedValue_Z,
        DPT_DeltaTimeMsec_Z,
        DPT_DeltaTime10Msec_Z,
        DPT_DeltaTime100Msec_Z,
        DPT_DeltaTimeSec_Z,
        DPT_DeltaTimeMin_Z,
        DPT_DeltaTimeHrs_Z,
        DPT_Percent_V16_Z,
        DPT_TempHVACAbs_Z,
        DPT_TempHVACRel_Z,
        DPT_HVACAirFlowRel_Z,
        DPT_HVACModeNext,
        DPT_DHWModeNext,
        DPT_OccModeNext,
        DPT_BuildingModeNext,
        DPT_StatusBUC,
        DPT_LockSign,
        DPT_ValueDemBOC,
        DPT_ActPosDemAbs,
        DPT_StatusAct,
        DPT_StatusLightingActuator,
        DPT_StatusHPM,
        DPT_TempRoomDemAbs,
        DPT_StatusCPM,
        DPT_StatusWTC,
        DPT_TempFlowWaterDemAbs,
        DPT_EnergyDemWater,
        DPT_TempRoomSetpSetShift3,
        DPT_TempRoomSetpSet3,
        DPT_TempRoomSetpSet4,
        DPT_TempDHWSetpSet4,
        DPT_TempRoomSetpSetShift4,
        DPT_PowerFlowWaterDemHPM,
        DPT_PowerFlowWaterDemCPM,
        DPT_StatusBOC,
        DPT_StatusCC,
        DPT_SpecHeatProd,
        DPT_Version,
        DPT_VolumeLiter_Z,
        DPT_FlowRate_m3h_Z,
        DPT_AlarmInfo,
        DPT_TempHVACAbsNext,
        DPT_SerNum,
        DPT_TempRoomSetpSetF163,
        DPT_TempRoomSetpSetShiftF163,
        DPT_EnergyDemAir,
        DPT_TempSupply_AirSetpSet,
        DPT_ScalingSpeed,
        DPT_Scaling_Step_Time,
        DPT_TariffNext,
        DPT_MeteringValue,
        DPT_MBus_Address,
        DPT_Locale_ASCII,
        DPT_Colour_RGB,
        DPT_LanguageCodeAlpha2_ASCII,
        DPT_RegionCodeAlpha2_ASCII,
        DPT_Tariff_ActiveEnergy,
        DPT_Prioritised_Mode_Control,
        DPT_DALI_Control_Gear_Diagnostic,
        DPT_SceneConfig,
        DPT_DALI_Diagnostics,
        DPT_FlaggedScaling,
        DPT_CombinedPosition,
        DPT_StatusSAB,

    }

    public static class KnxDatapointTypeInfo
    {

        public static string FormatName(this KnxDatapointType value)
        {
            switch (value)
            {
                case KnxDatapointType.DPT_ADAType: { /* 'DPT_ADAType' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Access_Data: { /* 'DPT_Access_Data' */
                    return "U4U4U4U4U4U4B4N4";
                }
                case KnxDatapointType.DPT_Ack: { /* 'DPT_Ack' */
                    return "B1";
                }
                case KnxDatapointType.DPT_ActPosDemAbs: { /* 'DPT_ActPosDemAbs' */
                    return "U8B8ActuatorPositionDemand";
                }
                case KnxDatapointType.DPT_ActiveEnergy: { /* 'DPT_ActiveEnergy' */
                    return "V32";
                }
                case KnxDatapointType.DPT_ActiveEnergy_V64: { /* 'DPT_ActiveEnergy_V64' */
                    return "V64";
                }
                case KnxDatapointType.DPT_ActiveEnergy_kWh: { /* 'DPT_ActiveEnergy_kWh' */
                    return "V32";
                }
                case KnxDatapointType.DPT_ActuatorConnectType: { /* 'DPT_ActuatorConnectType' */
                    return "N8";
                }
                case KnxDatapointType.DPT_AddInfoTypes: { /* 'DPT_AddInfoTypes' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Alarm: { /* 'DPT_Alarm' */
                    return "B1";
                }
                case KnxDatapointType.DPT_AlarmClassType: { /* 'DPT_AlarmClassType' */
                    return "N8";
                }
                case KnxDatapointType.DPT_AlarmInfo: { /* 'DPT_AlarmInfo' */
                    return "U8N8N8N8B8B8";
                }
                case KnxDatapointType.DPT_Alarm_Control: { /* 'DPT_Alarm_Control' */
                    return "B2";
                }
                case KnxDatapointType.DPT_Alarm_Reaction: { /* 'DPT_Alarm_Reaction' */
                    return "N2";
                }
                case KnxDatapointType.DPT_Angle: { /* 'DPT_Angle' */
                    return "U8";
                }
                case KnxDatapointType.DPT_ApparantEnergy: { /* 'DPT_ApparantEnergy' */
                    return "V32";
                }
                case KnxDatapointType.DPT_ApparantEnergy_V64: { /* 'DPT_ApparantEnergy_V64' */
                    return "V64";
                }
                case KnxDatapointType.DPT_ApparantEnergy_kVAh: { /* 'DPT_ApparantEnergy_kVAh' */
                    return "V32";
                }
                case KnxDatapointType.DPT_ApplicationArea: { /* 'DPT_ApplicationArea' */
                    return "N8";
                }
                case KnxDatapointType.DPT_AtmPressureAbs_Z: { /* 'DPT_AtmPressureAbs_Z' */
                    return "U16Z8AtmPressure";
                }
                case KnxDatapointType.DPT_BackupMode: { /* 'DPT_BackupMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Beaufort_Wind_Force_Scale: { /* 'DPT_Beaufort_Wind_Force_Scale' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Behaviour_Bus_Power_Up_Down: { /* 'DPT_Behaviour_Bus_Power_Up_Down' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Behaviour_Lock_Unlock: { /* 'DPT_Behaviour_Lock_Unlock' */
                    return "N8";
                }
                case KnxDatapointType.DPT_BinaryValue: { /* 'DPT_BinaryValue' */
                    return "B1";
                }
                case KnxDatapointType.DPT_BinaryValue_Control: { /* 'DPT_BinaryValue_Control' */
                    return "B2";
                }
                case KnxDatapointType.DPT_BinaryValue_Z: { /* 'DPT_BinaryValue_Z' */
                    return "B1Z8BinaryValueZ";
                }
                case KnxDatapointType.DPT_BlindsControlMode: { /* 'DPT_BlindsControlMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_BlinkingMode: { /* 'DPT_BlinkingMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Bool: { /* 'DPT_Bool' */
                    return "B1";
                }
                case KnxDatapointType.DPT_Bool_Control: { /* 'DPT_Bool_Control' */
                    return "B2";
                }
                case KnxDatapointType.DPT_Brightness: { /* 'DPT_Brightness' */
                    return "U16";
                }
                case KnxDatapointType.DPT_BuildingMode: { /* 'DPT_BuildingMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_BuildingModeNext: { /* 'DPT_BuildingModeNext' */
                    return "U16N8BuildingModeAndTimeDelay";
                }
                case KnxDatapointType.DPT_BuildingMode_Z: { /* 'DPT_BuildingMode_Z' */
                    return "N8Z8BuildingMode";
                }
                case KnxDatapointType.DPT_BurnerType: { /* 'DPT_BurnerType' */
                    return "N8";
                }
                case KnxDatapointType.DPT_ChangeoverMode: { /* 'DPT_ChangeoverMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Channel_Activation_16: { /* 'DPT_Channel_Activation_16' */
                    return "B16";
                }
                case KnxDatapointType.DPT_Channel_Activation_24: { /* 'DPT_Channel_Activation_24' */
                    return "B24";
                }
                case KnxDatapointType.DPT_Channel_Activation_8: { /* 'DPT_Channel_Activation_8' */
                    return "B8";
                }
                case KnxDatapointType.DPT_Char_8859_1: { /* 'DPT_Char_8859_1' */
                    return "A8_8859_1";
                }
                case KnxDatapointType.DPT_Char_ASCII: { /* 'DPT_Char_ASCII' */
                    return "A8_ASCII";
                }
                case KnxDatapointType.DPT_Colour_RGB: { /* 'DPT_Colour_RGB' */
                    return "U8U8U8";
                }
                case KnxDatapointType.DPT_CombinedInfoOnOff: { /* 'DPT_CombinedInfoOnOff' */
                    return "B32";
                }
                case KnxDatapointType.DPT_CombinedPosition: { /* 'DPT_CombinedPosition' */
                    return "U8U8B8";
                }
                case KnxDatapointType.DPT_CommMode: { /* 'DPT_CommMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Control_Blinds: { /* 'DPT_Control_Blinds' */
                    return "B1U3";
                }
                case KnxDatapointType.DPT_Control_Dimming: { /* 'DPT_Control_Dimming' */
                    return "B1U3";
                }
                case KnxDatapointType.DPT_DALI_Control_Gear_Diagnostic: { /* 'DPT_DALI_Control_Gear_Diagnostic' */
                    return "B10U6";
                }
                case KnxDatapointType.DPT_DALI_Diagnostics: { /* 'DPT_DALI_Diagnostics' */
                    return "B2U6";
                }
                case KnxDatapointType.DPT_DALI_Fade_Time: { /* 'DPT_DALI_Fade_Time' */
                    return "N8";
                }
                case KnxDatapointType.DPT_DHWMode: { /* 'DPT_DHWMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_DHWModeNext: { /* 'DPT_DHWModeNext' */
                    return "U16N8DhwModeAndTimeDelay";
                }
                case KnxDatapointType.DPT_DHWMode_Z: { /* 'DPT_DHWMode_Z' */
                    return "N8Z8DhwMode";
                }
                case KnxDatapointType.DPT_DamperMode: { /* 'DPT_DamperMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Date: { /* 'DPT_Date' */
                    return "r3N5r4N4r1U7";
                }
                case KnxDatapointType.DPT_DateTime: { /* 'DPT_DateTime' */
                    return "U8r4U4r3U5U3U5r2U6r2U6B16";
                }
                case KnxDatapointType.DPT_DecimalFactor: { /* 'DPT_DecimalFactor' */
                    return "U8";
                }
                case KnxDatapointType.DPT_DeltaTime100MSec: { /* 'DPT_DeltaTime100MSec' */
                    return "V16";
                }
                case KnxDatapointType.DPT_DeltaTime100Msec_Z: { /* 'DPT_DeltaTime100Msec_Z' */
                    return "V16Z8DeltaTime";
                }
                case KnxDatapointType.DPT_DeltaTime10MSec: { /* 'DPT_DeltaTime10MSec' */
                    return "V16";
                }
                case KnxDatapointType.DPT_DeltaTime10Msec_Z: { /* 'DPT_DeltaTime10Msec_Z' */
                    return "V16Z8DeltaTime";
                }
                case KnxDatapointType.DPT_DeltaTimeHrs: { /* 'DPT_DeltaTimeHrs' */
                    return "V16";
                }
                case KnxDatapointType.DPT_DeltaTimeHrs_Z: { /* 'DPT_DeltaTimeHrs_Z' */
                    return "V16Z8DeltaTime";
                }
                case KnxDatapointType.DPT_DeltaTimeMin: { /* 'DPT_DeltaTimeMin' */
                    return "V16";
                }
                case KnxDatapointType.DPT_DeltaTimeMin_Z: { /* 'DPT_DeltaTimeMin_Z' */
                    return "V16Z8DeltaTime";
                }
                case KnxDatapointType.DPT_DeltaTimeMsec: { /* 'DPT_DeltaTimeMsec' */
                    return "V16";
                }
                case KnxDatapointType.DPT_DeltaTimeMsec_Z: { /* 'DPT_DeltaTimeMsec_Z' */
                    return "V16Z8DeltaTime";
                }
                case KnxDatapointType.DPT_DeltaTimeSec: { /* 'DPT_DeltaTimeSec' */
                    return "V16";
                }
                case KnxDatapointType.DPT_DeltaTimeSec_Z: { /* 'DPT_DeltaTimeSec_Z' */
                    return "V16Z8DeltaTime";
                }
                case KnxDatapointType.DPT_Device_Control: { /* 'DPT_Device_Control' */
                    return "B8";
                }
                case KnxDatapointType.DPT_DimSendStyle: { /* 'DPT_DimSendStyle' */
                    return "B1";
                }
                case KnxDatapointType.DPT_DimmPBModel: { /* 'DPT_DimmPBModel' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Direction1_Control: { /* 'DPT_Direction1_Control' */
                    return "B2";
                }
                case KnxDatapointType.DPT_Direction2_Control: { /* 'DPT_Direction2_Control' */
                    return "B2";
                }
                case KnxDatapointType.DPT_DoubleNibble: { /* 'DPT_DoubleNibble' */
                    return "U4U4";
                }
                case KnxDatapointType.DPT_EnablH_Cstage_Z_DPT_EnablH_CStage: { /* 'DPT_EnablH_Cstage_Z_DPT_EnablH_CStage' */
                    return "N8Z8EnableHeatingOrCoolingStage";
                }
                case KnxDatapointType.DPT_Enable: { /* 'DPT_Enable' */
                    return "B1";
                }
                case KnxDatapointType.DPT_Enable_Control: { /* 'DPT_Enable_Control' */
                    return "B2";
                }
                case KnxDatapointType.DPT_EnergyDemAir: { /* 'DPT_EnergyDemAir' */
                    return "V8N8N8";
                }
                case KnxDatapointType.DPT_EnergyDemWater: { /* 'DPT_EnergyDemWater' */
                    return "U8N8";
                }
                case KnxDatapointType.DPT_ErrorClass_HVAC: { /* 'DPT_ErrorClass_HVAC' */
                    return "N8";
                }
                case KnxDatapointType.DPT_ErrorClass_System: { /* 'DPT_ErrorClass_System' */
                    return "N8";
                }
                case KnxDatapointType.DPT_FanMode: { /* 'DPT_FanMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_FlaggedScaling: { /* 'DPT_FlaggedScaling' */
                    return "U8r7B1";
                }
                case KnxDatapointType.DPT_FlowRate_m3h: { /* 'DPT_FlowRate_m3h' */
                    return "V32";
                }
                case KnxDatapointType.DPT_FlowRate_m3h_Z: { /* 'DPT_FlowRate_m3h_Z' */
                    return "V32Z8FlowRate";
                }
                case KnxDatapointType.DPT_ForceSign: { /* 'DPT_ForceSign' */
                    return "B8";
                }
                case KnxDatapointType.DPT_ForceSignCool: { /* 'DPT_ForceSignCool' */
                    return "B8";
                }
                case KnxDatapointType.DPT_FuelType: { /* 'DPT_FuelType' */
                    return "N8";
                }
                case KnxDatapointType.DPT_FuelTypeSet: { /* 'DPT_FuelTypeSet' */
                    return "B8";
                }
                case KnxDatapointType.DPT_HVACAirFlowAbs_Z: { /* 'DPT_HVACAirFlowAbs_Z' */
                    return "U16Z8HvacAirFlow";
                }
                case KnxDatapointType.DPT_HVACAirFlowRel_Z: { /* 'DPT_HVACAirFlowRel_Z' */
                    return "V16Z8RelSignedValue";
                }
                case KnxDatapointType.DPT_HVACAirQual_Z: { /* 'DPT_HVACAirQual_Z' */
                    return "U16Z8HvacAirQuality";
                }
                case KnxDatapointType.DPT_HVACContrMode: { /* 'DPT_HVACContrMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_HVACContrMode_Z: { /* 'DPT_HVACContrMode_Z' */
                    return "N8Z8HvacControllingMode";
                }
                case KnxDatapointType.DPT_HVACEmergMode: { /* 'DPT_HVACEmergMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_HVACEmergMode_Z: { /* 'DPT_HVACEmergMode_Z' */
                    return "N8Z8EmergencyMode";
                }
                case KnxDatapointType.DPT_HVACMode: { /* 'DPT_HVACMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_HVACModeNext: { /* 'DPT_HVACModeNext' */
                    return "U16N8HvacModeAndTimeDelay";
                }
                case KnxDatapointType.DPT_HVACMode_Z: { /* 'DPT_HVACMode_Z' */
                    return "N8Z8HvacOperatingMode";
                }
                case KnxDatapointType.DPT_HVAC_PB_Action: { /* 'DPT_HVAC_PB_Action' */
                    return "N2";
                }
                case KnxDatapointType.DPT_Heat_Cool: { /* 'DPT_Heat_Cool' */
                    return "B1";
                }
                case KnxDatapointType.DPT_Heat_Cool_Z: { /* 'DPT_Heat_Cool_Z' */
                    return "B1Z8HeatingOrCoolingZ";
                }
                case KnxDatapointType.DPT_HeaterMode: { /* 'DPT_HeaterMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_InputSource: { /* 'DPT_InputSource' */
                    return "B1";
                }
                case KnxDatapointType.DPT_Invert: { /* 'DPT_Invert' */
                    return "B1";
                }
                case KnxDatapointType.DPT_Invert_Control: { /* 'DPT_Invert_Control' */
                    return "B2";
                }
                case KnxDatapointType.DPT_KelvinPerPercent: { /* 'DPT_KelvinPerPercent' */
                    return "F16";
                }
                case KnxDatapointType.DPT_LanguageCodeAlpha2_ASCII: { /* 'DPT_LanguageCodeAlpha2_ASCII' */
                    return "A8A8Language";
                }
                case KnxDatapointType.DPT_Length_mm: { /* 'DPT_Length_mm' */
                    return "U16";
                }
                case KnxDatapointType.DPT_LightActuatorErrorInfo: { /* 'DPT_LightActuatorErrorInfo' */
                    return "B8";
                }
                case KnxDatapointType.DPT_LightApplicationMode: { /* 'DPT_LightApplicationMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_LightControlMode: { /* 'DPT_LightControlMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_LoadPriority: { /* 'DPT_LoadPriority' */
                    return "N8";
                }
                case KnxDatapointType.DPT_LoadTypeDetected: { /* 'DPT_LoadTypeDetected' */
                    return "N8";
                }
                case KnxDatapointType.DPT_LoadTypeSet: { /* 'DPT_LoadTypeSet' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Locale_ASCII: { /* 'DPT_Locale_ASCII' */
                    return "A8A8A8A8";
                }
                case KnxDatapointType.DPT_LockSign: { /* 'DPT_LockSign' */
                    return "U8B8LockingSignal";
                }
                case KnxDatapointType.DPT_LogicalFunction: { /* 'DPT_LogicalFunction' */
                    return "B1";
                }
                case KnxDatapointType.DPT_LongDeltaTimeSec: { /* 'DPT_LongDeltaTimeSec' */
                    return "V32";
                }
                case KnxDatapointType.DPT_MBus_Address: { /* 'DPT_MBus_Address' */
                    return "U16U32U8N8";
                }
                case KnxDatapointType.DPT_MasterSlaveMode: { /* 'DPT_MasterSlaveMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Media: { /* 'DPT_Media' */
                    return "B16";
                }
                case KnxDatapointType.DPT_MeteringValue: { /* 'DPT_MeteringValue' */
                    return "V32N8Z8";
                }
                case KnxDatapointType.DPT_OccMode: { /* 'DPT_OccMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_OccModeNext: { /* 'DPT_OccModeNext' */
                    return "U16N8OccupancyModeAndTimeDelay";
                }
                case KnxDatapointType.DPT_OccMode_Z: { /* 'DPT_OccMode_Z' */
                    return "N8Z8OccupancyMode";
                }
                case KnxDatapointType.DPT_Occupancy: { /* 'DPT_Occupancy' */
                    return "B1";
                }
                case KnxDatapointType.DPT_OnOff_Action: { /* 'DPT_OnOff_Action' */
                    return "N2";
                }
                case KnxDatapointType.DPT_OpenClose: { /* 'DPT_OpenClose' */
                    return "B1";
                }
                case KnxDatapointType.DPT_PBAction: { /* 'DPT_PBAction' */
                    return "N8";
                }
                case KnxDatapointType.DPT_PB_Action_HVAC_Extended: { /* 'DPT_PB_Action_HVAC_Extended' */
                    return "N3";
                }
                case KnxDatapointType.DPT_PSUMode: { /* 'DPT_PSUMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_PercentU16_Z: { /* 'DPT_PercentU16_Z' */
                    return "U16Z8PercentValue";
                }
                case KnxDatapointType.DPT_Percent_U8: { /* 'DPT_Percent_U8' */
                    return "U8";
                }
                case KnxDatapointType.DPT_Percent_V16: { /* 'DPT_Percent_V16' */
                    return "V16";
                }
                case KnxDatapointType.DPT_Percent_V16_Z: { /* 'DPT_Percent_V16_Z' */
                    return "V16Z8RelSignedValue";
                }
                case KnxDatapointType.DPT_Percent_V8: { /* 'DPT_Percent_V8' */
                    return "V8";
                }
                case KnxDatapointType.DPT_Power: { /* 'DPT_Power' */
                    return "F16";
                }
                case KnxDatapointType.DPT_PowerDensity: { /* 'DPT_PowerDensity' */
                    return "F16";
                }
                case KnxDatapointType.DPT_PowerFlowWaterDemCPM: { /* 'DPT_PowerFlowWaterDemCPM' */
                    return "V16U8B8ChilledWater";
                }
                case KnxDatapointType.DPT_PowerFlowWaterDemHPM: { /* 'DPT_PowerFlowWaterDemHPM' */
                    return "V16U8B8Heat";
                }
                case KnxDatapointType.DPT_PowerKW_Z: { /* 'DPT_PowerKW_Z' */
                    return "U16Z8Power";
                }
                case KnxDatapointType.DPT_Prioritised_Mode_Control: { /* 'DPT_Prioritised_Mode_Control' */
                    return "B1N3N4";
                }
                case KnxDatapointType.DPT_Priority: { /* 'DPT_Priority' */
                    return "N8";
                }
                case KnxDatapointType.DPT_PropDataType: { /* 'DPT_PropDataType' */
                    return "U16";
                }
                case KnxDatapointType.DPT_RF_FilterInfo: { /* 'DPT_RF_FilterInfo' */
                    return "B8";
                }
                case KnxDatapointType.DPT_RF_FilterSelect: { /* 'DPT_RF_FilterSelect' */
                    return "N8";
                }
                case KnxDatapointType.DPT_RF_ModeInfo: { /* 'DPT_RF_ModeInfo' */
                    return "B8";
                }
                case KnxDatapointType.DPT_RF_ModeSelect: { /* 'DPT_RF_ModeSelect' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Rain_Amount: { /* 'DPT_Rain_Amount' */
                    return "F16";
                }
                case KnxDatapointType.DPT_Ramp: { /* 'DPT_Ramp' */
                    return "B1";
                }
                case KnxDatapointType.DPT_Ramp_Control: { /* 'DPT_Ramp_Control' */
                    return "B2";
                }
                case KnxDatapointType.DPT_ReactiveEnergy: { /* 'DPT_ReactiveEnergy' */
                    return "V32";
                }
                case KnxDatapointType.DPT_ReactiveEnergy_V64: { /* 'DPT_ReactiveEnergy_V64' */
                    return "V64";
                }
                case KnxDatapointType.DPT_ReactiveEnergy_kVARh: { /* 'DPT_ReactiveEnergy_kVARh' */
                    return "V32";
                }
                case KnxDatapointType.DPT_RegionCodeAlpha2_ASCII: { /* 'DPT_RegionCodeAlpha2_ASCII' */
                    return "A8A8Region";
                }
                case KnxDatapointType.DPT_RelSignedValue_Z: { /* 'DPT_RelSignedValue_Z' */
                    return "V8Z8RelSignedValue";
                }
                case KnxDatapointType.DPT_RelValue_Z: { /* 'DPT_RelValue_Z' */
                    return "U8Z8Rel";
                }
                case KnxDatapointType.DPT_Reset: { /* 'DPT_Reset' */
                    return "B1";
                }
                case KnxDatapointType.DPT_Rotation_Angle: { /* 'DPT_Rotation_Angle' */
                    return "V16";
                }
                case KnxDatapointType.DPT_SABBehaviour_Lock_Unlock: { /* 'DPT_SABBehaviour_Lock_Unlock' */
                    return "N8";
                }
                case KnxDatapointType.DPT_SABExceptBehaviour: { /* 'DPT_SABExceptBehaviour' */
                    return "N8";
                }
                case KnxDatapointType.DPT_SCLOMode: { /* 'DPT_SCLOMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_SSSBMode: { /* 'DPT_SSSBMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Scaling: { /* 'DPT_Scaling' */
                    return "U8";
                }
                case KnxDatapointType.DPT_ScalingSpeed: { /* 'DPT_ScalingSpeed' */
                    return "U16U8Scaling";
                }
                case KnxDatapointType.DPT_Scaling_Step_Time: { /* 'DPT_Scaling_Step_Time' */
                    return "U16U8Scaling";
                }
                case KnxDatapointType.DPT_SceneConfig: { /* 'DPT_SceneConfig' */
                    return "B2U6";
                }
                case KnxDatapointType.DPT_SceneControl: { /* 'DPT_SceneControl' */
                    return "B1r1U6";
                }
                case KnxDatapointType.DPT_SceneInfo: { /* 'DPT_SceneInfo' */
                    return "r1b1U6";
                }
                case KnxDatapointType.DPT_SceneNumber: { /* 'DPT_SceneNumber' */
                    return "r2U6";
                }
                case KnxDatapointType.DPT_Scene_AB: { /* 'DPT_Scene_AB' */
                    return "B1";
                }
                case KnxDatapointType.DPT_SensorSelect: { /* 'DPT_SensorSelect' */
                    return "N8";
                }
                case KnxDatapointType.DPT_SerNum: { /* 'DPT_SerNum' */
                    return "N16U32";
                }
                case KnxDatapointType.DPT_ShutterBlinds_Mode: { /* 'DPT_ShutterBlinds_Mode' */
                    return "B1";
                }
                case KnxDatapointType.DPT_SpecHeatProd: { /* 'DPT_SpecHeatProd' */
                    return "U16U8N8B8";
                }
                case KnxDatapointType.DPT_Start: { /* 'DPT_Start' */
                    return "B1";
                }
                case KnxDatapointType.DPT_StartSynchronization: { /* 'DPT_StartSynchronization' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Start_Control: { /* 'DPT_Start_Control' */
                    return "B2";
                }
                case KnxDatapointType.DPT_State: { /* 'DPT_State' */
                    return "B1";
                }
                case KnxDatapointType.DPT_State_Control: { /* 'DPT_State_Control' */
                    return "B2";
                }
                case KnxDatapointType.DPT_StatusAHU: { /* 'DPT_StatusAHU' */
                    return "B8";
                }
                case KnxDatapointType.DPT_StatusAct: { /* 'DPT_StatusAct' */
                    return "U8B8ActuatorPositionStatus";
                }
                case KnxDatapointType.DPT_StatusBOC: { /* 'DPT_StatusBOC' */
                    return "V16U8B16Boiler";
                }
                case KnxDatapointType.DPT_StatusBUC: { /* 'DPT_StatusBUC' */
                    return "U8B8StatusBurnerController";
                }
                case KnxDatapointType.DPT_StatusCC: { /* 'DPT_StatusCC' */
                    return "V16U8B16Chiller";
                }
                case KnxDatapointType.DPT_StatusCPM: { /* 'DPT_StatusCPM' */
                    return "V16B8ColdWaterProducerManagerStatus";
                }
                case KnxDatapointType.DPT_StatusDHWC: { /* 'DPT_StatusDHWC' */
                    return "B16";
                }
                case KnxDatapointType.DPT_StatusGen: { /* 'DPT_StatusGen' */
                    return "B8";
                }
                case KnxDatapointType.DPT_StatusHPM: { /* 'DPT_StatusHPM' */
                    return "V16B8HeatProducerManagerStatus";
                }
                case KnxDatapointType.DPT_StatusLightingActuator: { /* 'DPT_StatusLightingActuator' */
                    return "U8B8StatusLightingActuator";
                }
                case KnxDatapointType.DPT_StatusRCC: { /* 'DPT_StatusRCC' */
                    return "B8";
                }
                case KnxDatapointType.DPT_StatusRHC: { /* 'DPT_StatusRHC' */
                    return "B8";
                }
                case KnxDatapointType.DPT_StatusRHCC: { /* 'DPT_StatusRHCC' */
                    return "B16";
                }
                case KnxDatapointType.DPT_StatusRoomSetp: { /* 'DPT_StatusRoomSetp' */
                    return "N8";
                }
                case KnxDatapointType.DPT_StatusSAB: { /* 'DPT_StatusSAB' */
                    return "U8U8B16";
                }
                case KnxDatapointType.DPT_StatusSDHWC: { /* 'DPT_StatusSDHWC' */
                    return "B8";
                }
                case KnxDatapointType.DPT_StatusWTC: { /* 'DPT_StatusWTC' */
                    return "V16B8WaterTemperatureControllerStatus";
                }
                case KnxDatapointType.DPT_Status_Mode3: { /* 'DPT_Status_Mode3' */
                    return "B5N3";
                }
                case KnxDatapointType.DPT_Step: { /* 'DPT_Step' */
                    return "B1";
                }
                case KnxDatapointType.DPT_Step_Control: { /* 'DPT_Step_Control' */
                    return "B2";
                }
                case KnxDatapointType.DPT_String_8859_1: { /* 'DPT_String_8859_1' */
                    return "A112_8859_1";
                }
                case KnxDatapointType.DPT_String_ASCII: { /* 'DPT_String_ASCII' */
                    return "A112_ASCII";
                }
                case KnxDatapointType.DPT_SunIntensity_Z: { /* 'DPT_SunIntensity_Z' */
                    return "U16Z8SunIntensity";
                }
                case KnxDatapointType.DPT_Switch: { /* 'DPT_Switch' */
                    return "B1";
                }
                case KnxDatapointType.DPT_SwitchOnMode: { /* 'DPT_SwitchOnMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_SwitchPBModel: { /* 'DPT_SwitchPBModel' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Switch_Control: { /* 'DPT_Switch_Control' */
                    return "B2";
                }
                case KnxDatapointType.DPT_Tariff: { /* 'DPT_Tariff' */
                    return "U8";
                }
                case KnxDatapointType.DPT_TariffNext: { /* 'DPT_TariffNext' */
                    return "U16U8TariffNext";
                }
                case KnxDatapointType.DPT_Tariff_ActiveEnergy: { /* 'DPT_Tariff_ActiveEnergy' */
                    return "V32U8B8";
                }
                case KnxDatapointType.DPT_TempDHWSetpSet4: { /* 'DPT_TempDHWSetpSet4' */
                    return "V16V16V16V16DhwtTemperature";
                }
                case KnxDatapointType.DPT_TempFlowWaterDemAbs: { /* 'DPT_TempFlowWaterDemAbs' */
                    return "V16B16";
                }
                case KnxDatapointType.DPT_TempHVACAbsNext: { /* 'DPT_TempHVACAbsNext' */
                    return "U16V16";
                }
                case KnxDatapointType.DPT_TempHVACAbs_Z: { /* 'DPT_TempHVACAbs_Z' */
                    return "V16Z8RelSignedValue";
                }
                case KnxDatapointType.DPT_TempHVACRel_Z: { /* 'DPT_TempHVACRel_Z' */
                    return "V16Z8RelSignedValue";
                }
                case KnxDatapointType.DPT_TempRoomDemAbs: { /* 'DPT_TempRoomDemAbs' */
                    return "V16B8RoomTemperatureDemand";
                }
                case KnxDatapointType.DPT_TempRoomSetpSet3: { /* 'DPT_TempRoomSetpSet3' */
                    return "V16V16V16RoomTemperature";
                }
                case KnxDatapointType.DPT_TempRoomSetpSet4: { /* 'DPT_TempRoomSetpSet4' */
                    return "V16V16V16V16RoomTemperature";
                }
                case KnxDatapointType.DPT_TempRoomSetpSetF163: { /* 'DPT_TempRoomSetpSetF163' */
                    return "F16F16F16";
                }
                case KnxDatapointType.DPT_TempRoomSetpSetShift3: { /* 'DPT_TempRoomSetpSetShift3' */
                    return "V16V16V16RoomTemperatureShift";
                }
                case KnxDatapointType.DPT_TempRoomSetpSetShift4: { /* 'DPT_TempRoomSetpSetShift4' */
                    return "V16V16V16V16RoomTemperatureShift";
                }
                case KnxDatapointType.DPT_TempRoomSetpSetShiftF163: { /* 'DPT_TempRoomSetpSetShiftF163' */
                    return "F16F16F16";
                }
                case KnxDatapointType.DPT_TempSupply_AirSetpSet: { /* 'DPT_TempSupply_AirSetpSet' */
                    return "V16V16N8N8";
                }
                case KnxDatapointType.DPT_TimeOfDay: { /* 'DPT_TimeOfDay' */
                    return "N3N5r2N6r2N6";
                }
                case KnxDatapointType.DPT_TimePeriod100MSec: { /* 'DPT_TimePeriod100MSec' */
                    return "U16";
                }
                case KnxDatapointType.DPT_TimePeriod100Msec_Z: { /* 'DPT_TimePeriod100Msec_Z' */
                    return "U16Z8TimePeriod";
                }
                case KnxDatapointType.DPT_TimePeriod10MSec: { /* 'DPT_TimePeriod10MSec' */
                    return "U16";
                }
                case KnxDatapointType.DPT_TimePeriod10Msec_Z: { /* 'DPT_TimePeriod10Msec_Z' */
                    return "U16Z8TimePeriod";
                }
                case KnxDatapointType.DPT_TimePeriodHrs: { /* 'DPT_TimePeriodHrs' */
                    return "U16";
                }
                case KnxDatapointType.DPT_TimePeriodHrs_Z: { /* 'DPT_TimePeriodHrs_Z' */
                    return "U16Z8TimePeriod";
                }
                case KnxDatapointType.DPT_TimePeriodMin: { /* 'DPT_TimePeriodMin' */
                    return "U16";
                }
                case KnxDatapointType.DPT_TimePeriodMin_Z: { /* 'DPT_TimePeriodMin_Z' */
                    return "U16Z8TimePeriod";
                }
                case KnxDatapointType.DPT_TimePeriodMsec: { /* 'DPT_TimePeriodMsec' */
                    return "U16";
                }
                case KnxDatapointType.DPT_TimePeriodMsec_Z: { /* 'DPT_TimePeriodMsec_Z' */
                    return "U16Z8TimePeriod";
                }
                case KnxDatapointType.DPT_TimePeriodSec: { /* 'DPT_TimePeriodSec' */
                    return "U16";
                }
                case KnxDatapointType.DPT_TimePeriodSec_Z: { /* 'DPT_TimePeriodSec_Z' */
                    return "U16Z8TimePeriod";
                }
                case KnxDatapointType.DPT_Time_Delay: { /* 'DPT_Time_Delay' */
                    return "N8";
                }
                case KnxDatapointType.DPT_Trigger: { /* 'DPT_Trigger' */
                    return "B1";
                }
                case KnxDatapointType.DPT_UCountValue16_Z: { /* 'DPT_UCountValue16_Z' */
                    return "U16Z8Counter";
                }
                case KnxDatapointType.DPT_UCountValue8_Z: { /* 'DPT_UCountValue8_Z' */
                    return "U8Z8Counter";
                }
                case KnxDatapointType.DPT_UElCurrentmA: { /* 'DPT_UElCurrentmA' */
                    return "U16";
                }
                case KnxDatapointType.DPT_UElCurrentyA_Z: { /* 'DPT_UElCurrentyA_Z' */
                    return "U16Z8ElectricCurrent";
                }
                case KnxDatapointType.DPT_UFlowRateLiter_h_Z: { /* 'DPT_UFlowRateLiter_h_Z' */
                    return "U16Z8FlowRate";
                }
                case KnxDatapointType.DPT_UTF_8: { /* 'DPT_UTF_8' */
                    return "An_UTF_8";
                }
                case KnxDatapointType.DPT_UpDown: { /* 'DPT_UpDown' */
                    return "B1";
                }
                case KnxDatapointType.DPT_UpDown_Action: { /* 'DPT_UpDown_Action' */
                    return "N2";
                }
                case KnxDatapointType.DPT_ValueDemBOC: { /* 'DPT_ValueDemBOC' */
                    return "U8B8BoilerControllerDemandSignal";
                }
                case KnxDatapointType.DPT_Value_1_Count: { /* 'DPT_Value_1_Count' */
                    return "V8";
                }
                case KnxDatapointType.DPT_Value_1_Ucount: { /* 'DPT_Value_1_Ucount' */
                    return "U8";
                }
                case KnxDatapointType.DPT_Value_2_Count: { /* 'DPT_Value_2_Count' */
                    return "V16";
                }
                case KnxDatapointType.DPT_Value_2_Ucount: { /* 'DPT_Value_2_Ucount' */
                    return "U16";
                }
                case KnxDatapointType.DPT_Value_4_Count: { /* 'DPT_Value_4_Count' */
                    return "V32";
                }
                case KnxDatapointType.DPT_Value_4_Ucount: { /* 'DPT_Value_4_Ucount' */
                    return "U32";
                }
                case KnxDatapointType.DPT_Value_Absolute_Temperature: { /* 'DPT_Value_Absolute_Temperature' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Acceleration: { /* 'DPT_Value_Acceleration' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Acceleration_Angular: { /* 'DPT_Value_Acceleration_Angular' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Activation_Energy: { /* 'DPT_Value_Activation_Energy' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Activity: { /* 'DPT_Value_Activity' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_AirQuality: { /* 'DPT_Value_AirQuality' */
                    return "F16";
                }
                case KnxDatapointType.DPT_Value_Amplitude: { /* 'DPT_Value_Amplitude' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_AngleDeg: { /* 'DPT_Value_AngleDeg' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_AngleRad: { /* 'DPT_Value_AngleRad' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Angular_Frequency: { /* 'DPT_Value_Angular_Frequency' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Angular_Momentum: { /* 'DPT_Value_Angular_Momentum' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Angular_Velocity: { /* 'DPT_Value_Angular_Velocity' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Area: { /* 'DPT_Value_Area' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Capacitance: { /* 'DPT_Value_Capacitance' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Charge_DensitySurface: { /* 'DPT_Value_Charge_DensitySurface' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Charge_DensityVolume: { /* 'DPT_Value_Charge_DensityVolume' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Common_Temperature: { /* 'DPT_Value_Common_Temperature' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Compressibility: { /* 'DPT_Value_Compressibility' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Conductance: { /* 'DPT_Value_Conductance' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Curr: { /* 'DPT_Value_Curr' */
                    return "F16";
                }
                case KnxDatapointType.DPT_Value_Density: { /* 'DPT_Value_Density' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Electric_Charge: { /* 'DPT_Value_Electric_Charge' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Electric_Current: { /* 'DPT_Value_Electric_Current' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Electric_CurrentDensity: { /* 'DPT_Value_Electric_CurrentDensity' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Electric_DipoleMoment: { /* 'DPT_Value_Electric_DipoleMoment' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Electric_Displacement: { /* 'DPT_Value_Electric_Displacement' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Electric_FieldStrength: { /* 'DPT_Value_Electric_FieldStrength' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Electric_Flux: { /* 'DPT_Value_Electric_Flux' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Electric_FluxDensity: { /* 'DPT_Value_Electric_FluxDensity' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Electric_Polarization: { /* 'DPT_Value_Electric_Polarization' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Electric_Potential: { /* 'DPT_Value_Electric_Potential' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Electric_PotentialDifference: { /* 'DPT_Value_Electric_PotentialDifference' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Electrical_Conductivity: { /* 'DPT_Value_Electrical_Conductivity' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_ElectromagneticMoment: { /* 'DPT_Value_ElectromagneticMoment' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Electromotive_Force: { /* 'DPT_Value_Electromotive_Force' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Energy: { /* 'DPT_Value_Energy' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Force: { /* 'DPT_Value_Force' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Frequency: { /* 'DPT_Value_Frequency' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Heat_Capacity: { /* 'DPT_Value_Heat_Capacity' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Heat_FlowRate: { /* 'DPT_Value_Heat_FlowRate' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Heat_Quantity: { /* 'DPT_Value_Heat_Quantity' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Humidity: { /* 'DPT_Value_Humidity' */
                    return "F16";
                }
                case KnxDatapointType.DPT_Value_Impedance: { /* 'DPT_Value_Impedance' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Length: { /* 'DPT_Value_Length' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Light_Quantity: { /* 'DPT_Value_Light_Quantity' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Luminance: { /* 'DPT_Value_Luminance' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Luminous_Flux: { /* 'DPT_Value_Luminous_Flux' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Luminous_Intensity: { /* 'DPT_Value_Luminous_Intensity' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Lux: { /* 'DPT_Value_Lux' */
                    return "F16";
                }
                case KnxDatapointType.DPT_Value_Magnetic_FieldStrength: { /* 'DPT_Value_Magnetic_FieldStrength' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Magnetic_Flux: { /* 'DPT_Value_Magnetic_Flux' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Magnetic_FluxDensity: { /* 'DPT_Value_Magnetic_FluxDensity' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Magnetic_Moment: { /* 'DPT_Value_Magnetic_Moment' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Magnetic_Polarization: { /* 'DPT_Value_Magnetic_Polarization' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Magnetization: { /* 'DPT_Value_Magnetization' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_MagnetomotiveForce: { /* 'DPT_Value_MagnetomotiveForce' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Mass: { /* 'DPT_Value_Mass' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_MassFlux: { /* 'DPT_Value_MassFlux' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Mol: { /* 'DPT_Value_Mol' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Momentum: { /* 'DPT_Value_Momentum' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Phase_AngleDeg: { /* 'DPT_Value_Phase_AngleDeg' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Phase_AngleRad: { /* 'DPT_Value_Phase_AngleRad' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Power: { /* 'DPT_Value_Power' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Power_Factor: { /* 'DPT_Value_Power_Factor' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Pres: { /* 'DPT_Value_Pres' */
                    return "F16";
                }
                case KnxDatapointType.DPT_Value_Pressure: { /* 'DPT_Value_Pressure' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Reactance: { /* 'DPT_Value_Reactance' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Resistance: { /* 'DPT_Value_Resistance' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Resistivity: { /* 'DPT_Value_Resistivity' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_SelfInductance: { /* 'DPT_Value_SelfInductance' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_SolidAngle: { /* 'DPT_Value_SolidAngle' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Sound_Intensity: { /* 'DPT_Value_Sound_Intensity' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Speed: { /* 'DPT_Value_Speed' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Stress: { /* 'DPT_Value_Stress' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Surface_Tension: { /* 'DPT_Value_Surface_Tension' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Temp: { /* 'DPT_Value_Temp' */
                    return "F16";
                }
                case KnxDatapointType.DPT_Value_Temp_F: { /* 'DPT_Value_Temp_F' */
                    return "F16";
                }
                case KnxDatapointType.DPT_Value_Tempa: { /* 'DPT_Value_Tempa' */
                    return "F16";
                }
                case KnxDatapointType.DPT_Value_Tempd: { /* 'DPT_Value_Tempd' */
                    return "F16";
                }
                case KnxDatapointType.DPT_Value_TemperatureDifference: { /* 'DPT_Value_TemperatureDifference' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Thermal_Capacity: { /* 'DPT_Value_Thermal_Capacity' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Thermal_Conductivity: { /* 'DPT_Value_Thermal_Conductivity' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_ThermoelectricPower: { /* 'DPT_Value_ThermoelectricPower' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Time: { /* 'DPT_Value_Time' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Time1: { /* 'DPT_Value_Time1' */
                    return "F16";
                }
                case KnxDatapointType.DPT_Value_Time2: { /* 'DPT_Value_Time2' */
                    return "F16";
                }
                case KnxDatapointType.DPT_Value_Torque: { /* 'DPT_Value_Torque' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Volt: { /* 'DPT_Value_Volt' */
                    return "F16";
                }
                case KnxDatapointType.DPT_Value_Volume: { /* 'DPT_Value_Volume' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Volume_Flow: { /* 'DPT_Value_Volume_Flow' */
                    return "F16";
                }
                case KnxDatapointType.DPT_Value_Volume_Flux: { /* 'DPT_Value_Volume_Flux' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Weight: { /* 'DPT_Value_Weight' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Work: { /* 'DPT_Value_Work' */
                    return "F32";
                }
                case KnxDatapointType.DPT_Value_Wsp: { /* 'DPT_Value_Wsp' */
                    return "F16";
                }
                case KnxDatapointType.DPT_Value_Wsp_kmh: { /* 'DPT_Value_Wsp_kmh' */
                    return "F16";
                }
                case KnxDatapointType.DPT_ValveMode: { /* 'DPT_ValveMode' */
                    return "N8";
                }
                case KnxDatapointType.DPT_VarString_8859_1: { /* 'DPT_VarString_8859_1' */
                    return "An_8859_1";
                }
                case KnxDatapointType.DPT_Version: { /* 'DPT_Version' */
                    return "U5U5U6";
                }
                case KnxDatapointType.DPT_VolumeLiter_Z: { /* 'DPT_VolumeLiter_Z' */
                    return "V32Z8VolumeLiter";
                }
                case KnxDatapointType.DPT_WindSpeed_Z_DPT_WindSpeed: { /* 'DPT_WindSpeed_Z_DPT_WindSpeed' */
                    return "U16Z8WindSpeed";
                }
                case KnxDatapointType.DPT_Window_Door: { /* 'DPT_Window_Door' */
                    return "B1";
                }
                default: {
                    return "";
                }
            }
        }

        public static ushort MainNumber(this KnxDatapointType value)
        {
            switch (value)
            {
                case KnxDatapointType.DPT_ADAType: { /* 'DPT_ADAType' */
                    return 20;
                }
                case KnxDatapointType.DPT_Access_Data: { /* 'DPT_Access_Data' */
                    return 15;
                }
                case KnxDatapointType.DPT_Ack: { /* 'DPT_Ack' */
                    return 1;
                }
                case KnxDatapointType.DPT_ActPosDemAbs: { /* 'DPT_ActPosDemAbs' */
                    return 207;
                }
                case KnxDatapointType.DPT_ActiveEnergy: { /* 'DPT_ActiveEnergy' */
                    return 13;
                }
                case KnxDatapointType.DPT_ActiveEnergy_V64: { /* 'DPT_ActiveEnergy_V64' */
                    return 29;
                }
                case KnxDatapointType.DPT_ActiveEnergy_kWh: { /* 'DPT_ActiveEnergy_kWh' */
                    return 13;
                }
                case KnxDatapointType.DPT_ActuatorConnectType: { /* 'DPT_ActuatorConnectType' */
                    return 20;
                }
                case KnxDatapointType.DPT_AddInfoTypes: { /* 'DPT_AddInfoTypes' */
                    return 20;
                }
                case KnxDatapointType.DPT_Alarm: { /* 'DPT_Alarm' */
                    return 1;
                }
                case KnxDatapointType.DPT_AlarmClassType: { /* 'DPT_AlarmClassType' */
                    return 20;
                }
                case KnxDatapointType.DPT_AlarmInfo: { /* 'DPT_AlarmInfo' */
                    return 219;
                }
                case KnxDatapointType.DPT_Alarm_Control: { /* 'DPT_Alarm_Control' */
                    return 2;
                }
                case KnxDatapointType.DPT_Alarm_Reaction: { /* 'DPT_Alarm_Reaction' */
                    return 23;
                }
                case KnxDatapointType.DPT_Angle: { /* 'DPT_Angle' */
                    return 5;
                }
                case KnxDatapointType.DPT_ApparantEnergy: { /* 'DPT_ApparantEnergy' */
                    return 13;
                }
                case KnxDatapointType.DPT_ApparantEnergy_V64: { /* 'DPT_ApparantEnergy_V64' */
                    return 29;
                }
                case KnxDatapointType.DPT_ApparantEnergy_kVAh: { /* 'DPT_ApparantEnergy_kVAh' */
                    return 13;
                }
                case KnxDatapointType.DPT_ApplicationArea: { /* 'DPT_ApplicationArea' */
                    return 20;
                }
                case KnxDatapointType.DPT_AtmPressureAbs_Z: { /* 'DPT_AtmPressureAbs_Z' */
                    return 203;
                }
                case KnxDatapointType.DPT_BackupMode: { /* 'DPT_BackupMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_Beaufort_Wind_Force_Scale: { /* 'DPT_Beaufort_Wind_Force_Scale' */
                    return 20;
                }
                case KnxDatapointType.DPT_Behaviour_Bus_Power_Up_Down: { /* 'DPT_Behaviour_Bus_Power_Up_Down' */
                    return 20;
                }
                case KnxDatapointType.DPT_Behaviour_Lock_Unlock: { /* 'DPT_Behaviour_Lock_Unlock' */
                    return 20;
                }
                case KnxDatapointType.DPT_BinaryValue: { /* 'DPT_BinaryValue' */
                    return 1;
                }
                case KnxDatapointType.DPT_BinaryValue_Control: { /* 'DPT_BinaryValue_Control' */
                    return 2;
                }
                case KnxDatapointType.DPT_BinaryValue_Z: { /* 'DPT_BinaryValue_Z' */
                    return 200;
                }
                case KnxDatapointType.DPT_BlindsControlMode: { /* 'DPT_BlindsControlMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_BlinkingMode: { /* 'DPT_BlinkingMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_Bool: { /* 'DPT_Bool' */
                    return 1;
                }
                case KnxDatapointType.DPT_Bool_Control: { /* 'DPT_Bool_Control' */
                    return 2;
                }
                case KnxDatapointType.DPT_Brightness: { /* 'DPT_Brightness' */
                    return 7;
                }
                case KnxDatapointType.DPT_BuildingMode: { /* 'DPT_BuildingMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_BuildingModeNext: { /* 'DPT_BuildingModeNext' */
                    return 206;
                }
                case KnxDatapointType.DPT_BuildingMode_Z: { /* 'DPT_BuildingMode_Z' */
                    return 201;
                }
                case KnxDatapointType.DPT_BurnerType: { /* 'DPT_BurnerType' */
                    return 20;
                }
                case KnxDatapointType.DPT_ChangeoverMode: { /* 'DPT_ChangeoverMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_Channel_Activation_16: { /* 'DPT_Channel_Activation_16' */
                    return 22;
                }
                case KnxDatapointType.DPT_Channel_Activation_24: { /* 'DPT_Channel_Activation_24' */
                    return 30;
                }
                case KnxDatapointType.DPT_Channel_Activation_8: { /* 'DPT_Channel_Activation_8' */
                    return 21;
                }
                case KnxDatapointType.DPT_Char_8859_1: { /* 'DPT_Char_8859_1' */
                    return 4;
                }
                case KnxDatapointType.DPT_Char_ASCII: { /* 'DPT_Char_ASCII' */
                    return 4;
                }
                case KnxDatapointType.DPT_Colour_RGB: { /* 'DPT_Colour_RGB' */
                    return 232;
                }
                case KnxDatapointType.DPT_CombinedInfoOnOff: { /* 'DPT_CombinedInfoOnOff' */
                    return 27;
                }
                case KnxDatapointType.DPT_CombinedPosition: { /* 'DPT_CombinedPosition' */
                    return 240;
                }
                case KnxDatapointType.DPT_CommMode: { /* 'DPT_CommMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_Control_Blinds: { /* 'DPT_Control_Blinds' */
                    return 3;
                }
                case KnxDatapointType.DPT_Control_Dimming: { /* 'DPT_Control_Dimming' */
                    return 3;
                }
                case KnxDatapointType.DPT_DALI_Control_Gear_Diagnostic: { /* 'DPT_DALI_Control_Gear_Diagnostic' */
                    return 237;
                }
                case KnxDatapointType.DPT_DALI_Diagnostics: { /* 'DPT_DALI_Diagnostics' */
                    return 238;
                }
                case KnxDatapointType.DPT_DALI_Fade_Time: { /* 'DPT_DALI_Fade_Time' */
                    return 20;
                }
                case KnxDatapointType.DPT_DHWMode: { /* 'DPT_DHWMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_DHWModeNext: { /* 'DPT_DHWModeNext' */
                    return 206;
                }
                case KnxDatapointType.DPT_DHWMode_Z: { /* 'DPT_DHWMode_Z' */
                    return 201;
                }
                case KnxDatapointType.DPT_DamperMode: { /* 'DPT_DamperMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_Date: { /* 'DPT_Date' */
                    return 11;
                }
                case KnxDatapointType.DPT_DateTime: { /* 'DPT_DateTime' */
                    return 19;
                }
                case KnxDatapointType.DPT_DecimalFactor: { /* 'DPT_DecimalFactor' */
                    return 5;
                }
                case KnxDatapointType.DPT_DeltaTime100MSec: { /* 'DPT_DeltaTime100MSec' */
                    return 8;
                }
                case KnxDatapointType.DPT_DeltaTime100Msec_Z: { /* 'DPT_DeltaTime100Msec_Z' */
                    return 205;
                }
                case KnxDatapointType.DPT_DeltaTime10MSec: { /* 'DPT_DeltaTime10MSec' */
                    return 8;
                }
                case KnxDatapointType.DPT_DeltaTime10Msec_Z: { /* 'DPT_DeltaTime10Msec_Z' */
                    return 205;
                }
                case KnxDatapointType.DPT_DeltaTimeHrs: { /* 'DPT_DeltaTimeHrs' */
                    return 8;
                }
                case KnxDatapointType.DPT_DeltaTimeHrs_Z: { /* 'DPT_DeltaTimeHrs_Z' */
                    return 205;
                }
                case KnxDatapointType.DPT_DeltaTimeMin: { /* 'DPT_DeltaTimeMin' */
                    return 8;
                }
                case KnxDatapointType.DPT_DeltaTimeMin_Z: { /* 'DPT_DeltaTimeMin_Z' */
                    return 205;
                }
                case KnxDatapointType.DPT_DeltaTimeMsec: { /* 'DPT_DeltaTimeMsec' */
                    return 8;
                }
                case KnxDatapointType.DPT_DeltaTimeMsec_Z: { /* 'DPT_DeltaTimeMsec_Z' */
                    return 205;
                }
                case KnxDatapointType.DPT_DeltaTimeSec: { /* 'DPT_DeltaTimeSec' */
                    return 8;
                }
                case KnxDatapointType.DPT_DeltaTimeSec_Z: { /* 'DPT_DeltaTimeSec_Z' */
                    return 205;
                }
                case KnxDatapointType.DPT_Device_Control: { /* 'DPT_Device_Control' */
                    return 21;
                }
                case KnxDatapointType.DPT_DimSendStyle: { /* 'DPT_DimSendStyle' */
                    return 1;
                }
                case KnxDatapointType.DPT_DimmPBModel: { /* 'DPT_DimmPBModel' */
                    return 20;
                }
                case KnxDatapointType.DPT_Direction1_Control: { /* 'DPT_Direction1_Control' */
                    return 2;
                }
                case KnxDatapointType.DPT_Direction2_Control: { /* 'DPT_Direction2_Control' */
                    return 2;
                }
                case KnxDatapointType.DPT_DoubleNibble: { /* 'DPT_DoubleNibble' */
                    return 25;
                }
                case KnxDatapointType.DPT_EnablH_Cstage_Z_DPT_EnablH_CStage: { /* 'DPT_EnablH_Cstage_Z_DPT_EnablH_CStage' */
                    return 201;
                }
                case KnxDatapointType.DPT_Enable: { /* 'DPT_Enable' */
                    return 1;
                }
                case KnxDatapointType.DPT_Enable_Control: { /* 'DPT_Enable_Control' */
                    return 2;
                }
                case KnxDatapointType.DPT_EnergyDemAir: { /* 'DPT_EnergyDemAir' */
                    return 223;
                }
                case KnxDatapointType.DPT_EnergyDemWater: { /* 'DPT_EnergyDemWater' */
                    return 211;
                }
                case KnxDatapointType.DPT_ErrorClass_HVAC: { /* 'DPT_ErrorClass_HVAC' */
                    return 20;
                }
                case KnxDatapointType.DPT_ErrorClass_System: { /* 'DPT_ErrorClass_System' */
                    return 20;
                }
                case KnxDatapointType.DPT_FanMode: { /* 'DPT_FanMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_FlaggedScaling: { /* 'DPT_FlaggedScaling' */
                    return 239;
                }
                case KnxDatapointType.DPT_FlowRate_m3h: { /* 'DPT_FlowRate_m3h' */
                    return 13;
                }
                case KnxDatapointType.DPT_FlowRate_m3h_Z: { /* 'DPT_FlowRate_m3h_Z' */
                    return 218;
                }
                case KnxDatapointType.DPT_ForceSign: { /* 'DPT_ForceSign' */
                    return 21;
                }
                case KnxDatapointType.DPT_ForceSignCool: { /* 'DPT_ForceSignCool' */
                    return 21;
                }
                case KnxDatapointType.DPT_FuelType: { /* 'DPT_FuelType' */
                    return 20;
                }
                case KnxDatapointType.DPT_FuelTypeSet: { /* 'DPT_FuelTypeSet' */
                    return 21;
                }
                case KnxDatapointType.DPT_HVACAirFlowAbs_Z: { /* 'DPT_HVACAirFlowAbs_Z' */
                    return 203;
                }
                case KnxDatapointType.DPT_HVACAirFlowRel_Z: { /* 'DPT_HVACAirFlowRel_Z' */
                    return 205;
                }
                case KnxDatapointType.DPT_HVACAirQual_Z: { /* 'DPT_HVACAirQual_Z' */
                    return 203;
                }
                case KnxDatapointType.DPT_HVACContrMode: { /* 'DPT_HVACContrMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_HVACContrMode_Z: { /* 'DPT_HVACContrMode_Z' */
                    return 201;
                }
                case KnxDatapointType.DPT_HVACEmergMode: { /* 'DPT_HVACEmergMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_HVACEmergMode_Z: { /* 'DPT_HVACEmergMode_Z' */
                    return 201;
                }
                case KnxDatapointType.DPT_HVACMode: { /* 'DPT_HVACMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_HVACModeNext: { /* 'DPT_HVACModeNext' */
                    return 206;
                }
                case KnxDatapointType.DPT_HVACMode_Z: { /* 'DPT_HVACMode_Z' */
                    return 201;
                }
                case KnxDatapointType.DPT_HVAC_PB_Action: { /* 'DPT_HVAC_PB_Action' */
                    return 23;
                }
                case KnxDatapointType.DPT_Heat_Cool: { /* 'DPT_Heat_Cool' */
                    return 1;
                }
                case KnxDatapointType.DPT_Heat_Cool_Z: { /* 'DPT_Heat_Cool_Z' */
                    return 200;
                }
                case KnxDatapointType.DPT_HeaterMode: { /* 'DPT_HeaterMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_InputSource: { /* 'DPT_InputSource' */
                    return 1;
                }
                case KnxDatapointType.DPT_Invert: { /* 'DPT_Invert' */
                    return 1;
                }
                case KnxDatapointType.DPT_Invert_Control: { /* 'DPT_Invert_Control' */
                    return 2;
                }
                case KnxDatapointType.DPT_KelvinPerPercent: { /* 'DPT_KelvinPerPercent' */
                    return 9;
                }
                case KnxDatapointType.DPT_LanguageCodeAlpha2_ASCII: { /* 'DPT_LanguageCodeAlpha2_ASCII' */
                    return 234;
                }
                case KnxDatapointType.DPT_Length_mm: { /* 'DPT_Length_mm' */
                    return 7;
                }
                case KnxDatapointType.DPT_LightActuatorErrorInfo: { /* 'DPT_LightActuatorErrorInfo' */
                    return 21;
                }
                case KnxDatapointType.DPT_LightApplicationMode: { /* 'DPT_LightApplicationMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_LightControlMode: { /* 'DPT_LightControlMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_LoadPriority: { /* 'DPT_LoadPriority' */
                    return 20;
                }
                case KnxDatapointType.DPT_LoadTypeDetected: { /* 'DPT_LoadTypeDetected' */
                    return 20;
                }
                case KnxDatapointType.DPT_LoadTypeSet: { /* 'DPT_LoadTypeSet' */
                    return 20;
                }
                case KnxDatapointType.DPT_Locale_ASCII: { /* 'DPT_Locale_ASCII' */
                    return 231;
                }
                case KnxDatapointType.DPT_LockSign: { /* 'DPT_LockSign' */
                    return 207;
                }
                case KnxDatapointType.DPT_LogicalFunction: { /* 'DPT_LogicalFunction' */
                    return 1;
                }
                case KnxDatapointType.DPT_LongDeltaTimeSec: { /* 'DPT_LongDeltaTimeSec' */
                    return 13;
                }
                case KnxDatapointType.DPT_MBus_Address: { /* 'DPT_MBus_Address' */
                    return 230;
                }
                case KnxDatapointType.DPT_MasterSlaveMode: { /* 'DPT_MasterSlaveMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_Media: { /* 'DPT_Media' */
                    return 22;
                }
                case KnxDatapointType.DPT_MeteringValue: { /* 'DPT_MeteringValue' */
                    return 229;
                }
                case KnxDatapointType.DPT_OccMode: { /* 'DPT_OccMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_OccModeNext: { /* 'DPT_OccModeNext' */
                    return 206;
                }
                case KnxDatapointType.DPT_OccMode_Z: { /* 'DPT_OccMode_Z' */
                    return 201;
                }
                case KnxDatapointType.DPT_Occupancy: { /* 'DPT_Occupancy' */
                    return 1;
                }
                case KnxDatapointType.DPT_OnOff_Action: { /* 'DPT_OnOff_Action' */
                    return 23;
                }
                case KnxDatapointType.DPT_OpenClose: { /* 'DPT_OpenClose' */
                    return 1;
                }
                case KnxDatapointType.DPT_PBAction: { /* 'DPT_PBAction' */
                    return 20;
                }
                case KnxDatapointType.DPT_PB_Action_HVAC_Extended: { /* 'DPT_PB_Action_HVAC_Extended' */
                    return 31;
                }
                case KnxDatapointType.DPT_PSUMode: { /* 'DPT_PSUMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_PercentU16_Z: { /* 'DPT_PercentU16_Z' */
                    return 203;
                }
                case KnxDatapointType.DPT_Percent_U8: { /* 'DPT_Percent_U8' */
                    return 5;
                }
                case KnxDatapointType.DPT_Percent_V16: { /* 'DPT_Percent_V16' */
                    return 8;
                }
                case KnxDatapointType.DPT_Percent_V16_Z: { /* 'DPT_Percent_V16_Z' */
                    return 205;
                }
                case KnxDatapointType.DPT_Percent_V8: { /* 'DPT_Percent_V8' */
                    return 6;
                }
                case KnxDatapointType.DPT_Power: { /* 'DPT_Power' */
                    return 9;
                }
                case KnxDatapointType.DPT_PowerDensity: { /* 'DPT_PowerDensity' */
                    return 9;
                }
                case KnxDatapointType.DPT_PowerFlowWaterDemCPM: { /* 'DPT_PowerFlowWaterDemCPM' */
                    return 214;
                }
                case KnxDatapointType.DPT_PowerFlowWaterDemHPM: { /* 'DPT_PowerFlowWaterDemHPM' */
                    return 214;
                }
                case KnxDatapointType.DPT_PowerKW_Z: { /* 'DPT_PowerKW_Z' */
                    return 203;
                }
                case KnxDatapointType.DPT_Prioritised_Mode_Control: { /* 'DPT_Prioritised_Mode_Control' */
                    return 236;
                }
                case KnxDatapointType.DPT_Priority: { /* 'DPT_Priority' */
                    return 20;
                }
                case KnxDatapointType.DPT_PropDataType: { /* 'DPT_PropDataType' */
                    return 7;
                }
                case KnxDatapointType.DPT_RF_FilterInfo: { /* 'DPT_RF_FilterInfo' */
                    return 21;
                }
                case KnxDatapointType.DPT_RF_FilterSelect: { /* 'DPT_RF_FilterSelect' */
                    return 20;
                }
                case KnxDatapointType.DPT_RF_ModeInfo: { /* 'DPT_RF_ModeInfo' */
                    return 21;
                }
                case KnxDatapointType.DPT_RF_ModeSelect: { /* 'DPT_RF_ModeSelect' */
                    return 20;
                }
                case KnxDatapointType.DPT_Rain_Amount: { /* 'DPT_Rain_Amount' */
                    return 9;
                }
                case KnxDatapointType.DPT_Ramp: { /* 'DPT_Ramp' */
                    return 1;
                }
                case KnxDatapointType.DPT_Ramp_Control: { /* 'DPT_Ramp_Control' */
                    return 2;
                }
                case KnxDatapointType.DPT_ReactiveEnergy: { /* 'DPT_ReactiveEnergy' */
                    return 13;
                }
                case KnxDatapointType.DPT_ReactiveEnergy_V64: { /* 'DPT_ReactiveEnergy_V64' */
                    return 29;
                }
                case KnxDatapointType.DPT_ReactiveEnergy_kVARh: { /* 'DPT_ReactiveEnergy_kVARh' */
                    return 13;
                }
                case KnxDatapointType.DPT_RegionCodeAlpha2_ASCII: { /* 'DPT_RegionCodeAlpha2_ASCII' */
                    return 234;
                }
                case KnxDatapointType.DPT_RelSignedValue_Z: { /* 'DPT_RelSignedValue_Z' */
                    return 204;
                }
                case KnxDatapointType.DPT_RelValue_Z: { /* 'DPT_RelValue_Z' */
                    return 202;
                }
                case KnxDatapointType.DPT_Reset: { /* 'DPT_Reset' */
                    return 1;
                }
                case KnxDatapointType.DPT_Rotation_Angle: { /* 'DPT_Rotation_Angle' */
                    return 8;
                }
                case KnxDatapointType.DPT_SABBehaviour_Lock_Unlock: { /* 'DPT_SABBehaviour_Lock_Unlock' */
                    return 20;
                }
                case KnxDatapointType.DPT_SABExceptBehaviour: { /* 'DPT_SABExceptBehaviour' */
                    return 20;
                }
                case KnxDatapointType.DPT_SCLOMode: { /* 'DPT_SCLOMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_SSSBMode: { /* 'DPT_SSSBMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_Scaling: { /* 'DPT_Scaling' */
                    return 5;
                }
                case KnxDatapointType.DPT_ScalingSpeed: { /* 'DPT_ScalingSpeed' */
                    return 225;
                }
                case KnxDatapointType.DPT_Scaling_Step_Time: { /* 'DPT_Scaling_Step_Time' */
                    return 225;
                }
                case KnxDatapointType.DPT_SceneConfig: { /* 'DPT_SceneConfig' */
                    return 238;
                }
                case KnxDatapointType.DPT_SceneControl: { /* 'DPT_SceneControl' */
                    return 18;
                }
                case KnxDatapointType.DPT_SceneInfo: { /* 'DPT_SceneInfo' */
                    return 26;
                }
                case KnxDatapointType.DPT_SceneNumber: { /* 'DPT_SceneNumber' */
                    return 17;
                }
                case KnxDatapointType.DPT_Scene_AB: { /* 'DPT_Scene_AB' */
                    return 1;
                }
                case KnxDatapointType.DPT_SensorSelect: { /* 'DPT_SensorSelect' */
                    return 20;
                }
                case KnxDatapointType.DPT_SerNum: { /* 'DPT_SerNum' */
                    return 221;
                }
                case KnxDatapointType.DPT_ShutterBlinds_Mode: { /* 'DPT_ShutterBlinds_Mode' */
                    return 1;
                }
                case KnxDatapointType.DPT_SpecHeatProd: { /* 'DPT_SpecHeatProd' */
                    return 216;
                }
                case KnxDatapointType.DPT_Start: { /* 'DPT_Start' */
                    return 1;
                }
                case KnxDatapointType.DPT_StartSynchronization: { /* 'DPT_StartSynchronization' */
                    return 20;
                }
                case KnxDatapointType.DPT_Start_Control: { /* 'DPT_Start_Control' */
                    return 2;
                }
                case KnxDatapointType.DPT_State: { /* 'DPT_State' */
                    return 1;
                }
                case KnxDatapointType.DPT_State_Control: { /* 'DPT_State_Control' */
                    return 2;
                }
                case KnxDatapointType.DPT_StatusAHU: { /* 'DPT_StatusAHU' */
                    return 21;
                }
                case KnxDatapointType.DPT_StatusAct: { /* 'DPT_StatusAct' */
                    return 207;
                }
                case KnxDatapointType.DPT_StatusBOC: { /* 'DPT_StatusBOC' */
                    return 215;
                }
                case KnxDatapointType.DPT_StatusBUC: { /* 'DPT_StatusBUC' */
                    return 207;
                }
                case KnxDatapointType.DPT_StatusCC: { /* 'DPT_StatusCC' */
                    return 215;
                }
                case KnxDatapointType.DPT_StatusCPM: { /* 'DPT_StatusCPM' */
                    return 209;
                }
                case KnxDatapointType.DPT_StatusDHWC: { /* 'DPT_StatusDHWC' */
                    return 22;
                }
                case KnxDatapointType.DPT_StatusGen: { /* 'DPT_StatusGen' */
                    return 21;
                }
                case KnxDatapointType.DPT_StatusHPM: { /* 'DPT_StatusHPM' */
                    return 209;
                }
                case KnxDatapointType.DPT_StatusLightingActuator: { /* 'DPT_StatusLightingActuator' */
                    return 207;
                }
                case KnxDatapointType.DPT_StatusRCC: { /* 'DPT_StatusRCC' */
                    return 21;
                }
                case KnxDatapointType.DPT_StatusRHC: { /* 'DPT_StatusRHC' */
                    return 21;
                }
                case KnxDatapointType.DPT_StatusRHCC: { /* 'DPT_StatusRHCC' */
                    return 22;
                }
                case KnxDatapointType.DPT_StatusRoomSetp: { /* 'DPT_StatusRoomSetp' */
                    return 20;
                }
                case KnxDatapointType.DPT_StatusSAB: { /* 'DPT_StatusSAB' */
                    return 241;
                }
                case KnxDatapointType.DPT_StatusSDHWC: { /* 'DPT_StatusSDHWC' */
                    return 21;
                }
                case KnxDatapointType.DPT_StatusWTC: { /* 'DPT_StatusWTC' */
                    return 209;
                }
                case KnxDatapointType.DPT_Status_Mode3: { /* 'DPT_Status_Mode3' */
                    return 6;
                }
                case KnxDatapointType.DPT_Step: { /* 'DPT_Step' */
                    return 1;
                }
                case KnxDatapointType.DPT_Step_Control: { /* 'DPT_Step_Control' */
                    return 2;
                }
                case KnxDatapointType.DPT_String_8859_1: { /* 'DPT_String_8859_1' */
                    return 16;
                }
                case KnxDatapointType.DPT_String_ASCII: { /* 'DPT_String_ASCII' */
                    return 16;
                }
                case KnxDatapointType.DPT_SunIntensity_Z: { /* 'DPT_SunIntensity_Z' */
                    return 203;
                }
                case KnxDatapointType.DPT_Switch: { /* 'DPT_Switch' */
                    return 1;
                }
                case KnxDatapointType.DPT_SwitchOnMode: { /* 'DPT_SwitchOnMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_SwitchPBModel: { /* 'DPT_SwitchPBModel' */
                    return 20;
                }
                case KnxDatapointType.DPT_Switch_Control: { /* 'DPT_Switch_Control' */
                    return 2;
                }
                case KnxDatapointType.DPT_Tariff: { /* 'DPT_Tariff' */
                    return 5;
                }
                case KnxDatapointType.DPT_TariffNext: { /* 'DPT_TariffNext' */
                    return 225;
                }
                case KnxDatapointType.DPT_Tariff_ActiveEnergy: { /* 'DPT_Tariff_ActiveEnergy' */
                    return 235;
                }
                case KnxDatapointType.DPT_TempDHWSetpSet4: { /* 'DPT_TempDHWSetpSet4' */
                    return 213;
                }
                case KnxDatapointType.DPT_TempFlowWaterDemAbs: { /* 'DPT_TempFlowWaterDemAbs' */
                    return 210;
                }
                case KnxDatapointType.DPT_TempHVACAbsNext: { /* 'DPT_TempHVACAbsNext' */
                    return 220;
                }
                case KnxDatapointType.DPT_TempHVACAbs_Z: { /* 'DPT_TempHVACAbs_Z' */
                    return 205;
                }
                case KnxDatapointType.DPT_TempHVACRel_Z: { /* 'DPT_TempHVACRel_Z' */
                    return 205;
                }
                case KnxDatapointType.DPT_TempRoomDemAbs: { /* 'DPT_TempRoomDemAbs' */
                    return 209;
                }
                case KnxDatapointType.DPT_TempRoomSetpSet3: { /* 'DPT_TempRoomSetpSet3' */
                    return 212;
                }
                case KnxDatapointType.DPT_TempRoomSetpSet4: { /* 'DPT_TempRoomSetpSet4' */
                    return 213;
                }
                case KnxDatapointType.DPT_TempRoomSetpSetF163: { /* 'DPT_TempRoomSetpSetF163' */
                    return 222;
                }
                case KnxDatapointType.DPT_TempRoomSetpSetShift3: { /* 'DPT_TempRoomSetpSetShift3' */
                    return 212;
                }
                case KnxDatapointType.DPT_TempRoomSetpSetShift4: { /* 'DPT_TempRoomSetpSetShift4' */
                    return 213;
                }
                case KnxDatapointType.DPT_TempRoomSetpSetShiftF163: { /* 'DPT_TempRoomSetpSetShiftF163' */
                    return 222;
                }
                case KnxDatapointType.DPT_TempSupply_AirSetpSet: { /* 'DPT_TempSupply_AirSetpSet' */
                    return 224;
                }
                case KnxDatapointType.DPT_TimeOfDay: { /* 'DPT_TimeOfDay' */
                    return 10;
                }
                case KnxDatapointType.DPT_TimePeriod100MSec: { /* 'DPT_TimePeriod100MSec' */
                    return 7;
                }
                case KnxDatapointType.DPT_TimePeriod100Msec_Z: { /* 'DPT_TimePeriod100Msec_Z' */
                    return 203;
                }
                case KnxDatapointType.DPT_TimePeriod10MSec: { /* 'DPT_TimePeriod10MSec' */
                    return 7;
                }
                case KnxDatapointType.DPT_TimePeriod10Msec_Z: { /* 'DPT_TimePeriod10Msec_Z' */
                    return 203;
                }
                case KnxDatapointType.DPT_TimePeriodHrs: { /* 'DPT_TimePeriodHrs' */
                    return 7;
                }
                case KnxDatapointType.DPT_TimePeriodHrs_Z: { /* 'DPT_TimePeriodHrs_Z' */
                    return 203;
                }
                case KnxDatapointType.DPT_TimePeriodMin: { /* 'DPT_TimePeriodMin' */
                    return 7;
                }
                case KnxDatapointType.DPT_TimePeriodMin_Z: { /* 'DPT_TimePeriodMin_Z' */
                    return 203;
                }
                case KnxDatapointType.DPT_TimePeriodMsec: { /* 'DPT_TimePeriodMsec' */
                    return 7;
                }
                case KnxDatapointType.DPT_TimePeriodMsec_Z: { /* 'DPT_TimePeriodMsec_Z' */
                    return 203;
                }
                case KnxDatapointType.DPT_TimePeriodSec: { /* 'DPT_TimePeriodSec' */
                    return 7;
                }
                case KnxDatapointType.DPT_TimePeriodSec_Z: { /* 'DPT_TimePeriodSec_Z' */
                    return 203;
                }
                case KnxDatapointType.DPT_Time_Delay: { /* 'DPT_Time_Delay' */
                    return 20;
                }
                case KnxDatapointType.DPT_Trigger: { /* 'DPT_Trigger' */
                    return 1;
                }
                case KnxDatapointType.DPT_UCountValue16_Z: { /* 'DPT_UCountValue16_Z' */
                    return 203;
                }
                case KnxDatapointType.DPT_UCountValue8_Z: { /* 'DPT_UCountValue8_Z' */
                    return 202;
                }
                case KnxDatapointType.DPT_UElCurrentmA: { /* 'DPT_UElCurrentmA' */
                    return 7;
                }
                case KnxDatapointType.DPT_UElCurrentyA_Z: { /* 'DPT_UElCurrentyA_Z' */
                    return 203;
                }
                case KnxDatapointType.DPT_UFlowRateLiter_h_Z: { /* 'DPT_UFlowRateLiter_h_Z' */
                    return 203;
                }
                case KnxDatapointType.DPT_UTF_8: { /* 'DPT_UTF_8' */
                    return 28;
                }
                case KnxDatapointType.DPT_UpDown: { /* 'DPT_UpDown' */
                    return 1;
                }
                case KnxDatapointType.DPT_UpDown_Action: { /* 'DPT_UpDown_Action' */
                    return 23;
                }
                case KnxDatapointType.DPT_ValueDemBOC: { /* 'DPT_ValueDemBOC' */
                    return 207;
                }
                case KnxDatapointType.DPT_Value_1_Count: { /* 'DPT_Value_1_Count' */
                    return 6;
                }
                case KnxDatapointType.DPT_Value_1_Ucount: { /* 'DPT_Value_1_Ucount' */
                    return 5;
                }
                case KnxDatapointType.DPT_Value_2_Count: { /* 'DPT_Value_2_Count' */
                    return 8;
                }
                case KnxDatapointType.DPT_Value_2_Ucount: { /* 'DPT_Value_2_Ucount' */
                    return 7;
                }
                case KnxDatapointType.DPT_Value_4_Count: { /* 'DPT_Value_4_Count' */
                    return 13;
                }
                case KnxDatapointType.DPT_Value_4_Ucount: { /* 'DPT_Value_4_Ucount' */
                    return 12;
                }
                case KnxDatapointType.DPT_Value_Absolute_Temperature: { /* 'DPT_Value_Absolute_Temperature' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Acceleration: { /* 'DPT_Value_Acceleration' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Acceleration_Angular: { /* 'DPT_Value_Acceleration_Angular' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Activation_Energy: { /* 'DPT_Value_Activation_Energy' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Activity: { /* 'DPT_Value_Activity' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_AirQuality: { /* 'DPT_Value_AirQuality' */
                    return 9;
                }
                case KnxDatapointType.DPT_Value_Amplitude: { /* 'DPT_Value_Amplitude' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_AngleDeg: { /* 'DPT_Value_AngleDeg' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_AngleRad: { /* 'DPT_Value_AngleRad' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Angular_Frequency: { /* 'DPT_Value_Angular_Frequency' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Angular_Momentum: { /* 'DPT_Value_Angular_Momentum' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Angular_Velocity: { /* 'DPT_Value_Angular_Velocity' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Area: { /* 'DPT_Value_Area' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Capacitance: { /* 'DPT_Value_Capacitance' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Charge_DensitySurface: { /* 'DPT_Value_Charge_DensitySurface' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Charge_DensityVolume: { /* 'DPT_Value_Charge_DensityVolume' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Common_Temperature: { /* 'DPT_Value_Common_Temperature' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Compressibility: { /* 'DPT_Value_Compressibility' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Conductance: { /* 'DPT_Value_Conductance' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Curr: { /* 'DPT_Value_Curr' */
                    return 9;
                }
                case KnxDatapointType.DPT_Value_Density: { /* 'DPT_Value_Density' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Electric_Charge: { /* 'DPT_Value_Electric_Charge' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Electric_Current: { /* 'DPT_Value_Electric_Current' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Electric_CurrentDensity: { /* 'DPT_Value_Electric_CurrentDensity' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Electric_DipoleMoment: { /* 'DPT_Value_Electric_DipoleMoment' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Electric_Displacement: { /* 'DPT_Value_Electric_Displacement' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Electric_FieldStrength: { /* 'DPT_Value_Electric_FieldStrength' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Electric_Flux: { /* 'DPT_Value_Electric_Flux' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Electric_FluxDensity: { /* 'DPT_Value_Electric_FluxDensity' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Electric_Polarization: { /* 'DPT_Value_Electric_Polarization' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Electric_Potential: { /* 'DPT_Value_Electric_Potential' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Electric_PotentialDifference: { /* 'DPT_Value_Electric_PotentialDifference' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Electrical_Conductivity: { /* 'DPT_Value_Electrical_Conductivity' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_ElectromagneticMoment: { /* 'DPT_Value_ElectromagneticMoment' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Electromotive_Force: { /* 'DPT_Value_Electromotive_Force' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Energy: { /* 'DPT_Value_Energy' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Force: { /* 'DPT_Value_Force' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Frequency: { /* 'DPT_Value_Frequency' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Heat_Capacity: { /* 'DPT_Value_Heat_Capacity' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Heat_FlowRate: { /* 'DPT_Value_Heat_FlowRate' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Heat_Quantity: { /* 'DPT_Value_Heat_Quantity' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Humidity: { /* 'DPT_Value_Humidity' */
                    return 9;
                }
                case KnxDatapointType.DPT_Value_Impedance: { /* 'DPT_Value_Impedance' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Length: { /* 'DPT_Value_Length' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Light_Quantity: { /* 'DPT_Value_Light_Quantity' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Luminance: { /* 'DPT_Value_Luminance' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Luminous_Flux: { /* 'DPT_Value_Luminous_Flux' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Luminous_Intensity: { /* 'DPT_Value_Luminous_Intensity' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Lux: { /* 'DPT_Value_Lux' */
                    return 9;
                }
                case KnxDatapointType.DPT_Value_Magnetic_FieldStrength: { /* 'DPT_Value_Magnetic_FieldStrength' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Magnetic_Flux: { /* 'DPT_Value_Magnetic_Flux' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Magnetic_FluxDensity: { /* 'DPT_Value_Magnetic_FluxDensity' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Magnetic_Moment: { /* 'DPT_Value_Magnetic_Moment' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Magnetic_Polarization: { /* 'DPT_Value_Magnetic_Polarization' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Magnetization: { /* 'DPT_Value_Magnetization' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_MagnetomotiveForce: { /* 'DPT_Value_MagnetomotiveForce' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Mass: { /* 'DPT_Value_Mass' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_MassFlux: { /* 'DPT_Value_MassFlux' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Mol: { /* 'DPT_Value_Mol' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Momentum: { /* 'DPT_Value_Momentum' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Phase_AngleDeg: { /* 'DPT_Value_Phase_AngleDeg' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Phase_AngleRad: { /* 'DPT_Value_Phase_AngleRad' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Power: { /* 'DPT_Value_Power' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Power_Factor: { /* 'DPT_Value_Power_Factor' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Pres: { /* 'DPT_Value_Pres' */
                    return 9;
                }
                case KnxDatapointType.DPT_Value_Pressure: { /* 'DPT_Value_Pressure' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Reactance: { /* 'DPT_Value_Reactance' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Resistance: { /* 'DPT_Value_Resistance' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Resistivity: { /* 'DPT_Value_Resistivity' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_SelfInductance: { /* 'DPT_Value_SelfInductance' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_SolidAngle: { /* 'DPT_Value_SolidAngle' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Sound_Intensity: { /* 'DPT_Value_Sound_Intensity' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Speed: { /* 'DPT_Value_Speed' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Stress: { /* 'DPT_Value_Stress' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Surface_Tension: { /* 'DPT_Value_Surface_Tension' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Temp: { /* 'DPT_Value_Temp' */
                    return 9;
                }
                case KnxDatapointType.DPT_Value_Temp_F: { /* 'DPT_Value_Temp_F' */
                    return 9;
                }
                case KnxDatapointType.DPT_Value_Tempa: { /* 'DPT_Value_Tempa' */
                    return 9;
                }
                case KnxDatapointType.DPT_Value_Tempd: { /* 'DPT_Value_Tempd' */
                    return 9;
                }
                case KnxDatapointType.DPT_Value_TemperatureDifference: { /* 'DPT_Value_TemperatureDifference' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Thermal_Capacity: { /* 'DPT_Value_Thermal_Capacity' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Thermal_Conductivity: { /* 'DPT_Value_Thermal_Conductivity' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_ThermoelectricPower: { /* 'DPT_Value_ThermoelectricPower' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Time: { /* 'DPT_Value_Time' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Time1: { /* 'DPT_Value_Time1' */
                    return 9;
                }
                case KnxDatapointType.DPT_Value_Time2: { /* 'DPT_Value_Time2' */
                    return 9;
                }
                case KnxDatapointType.DPT_Value_Torque: { /* 'DPT_Value_Torque' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Volt: { /* 'DPT_Value_Volt' */
                    return 9;
                }
                case KnxDatapointType.DPT_Value_Volume: { /* 'DPT_Value_Volume' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Volume_Flow: { /* 'DPT_Value_Volume_Flow' */
                    return 9;
                }
                case KnxDatapointType.DPT_Value_Volume_Flux: { /* 'DPT_Value_Volume_Flux' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Weight: { /* 'DPT_Value_Weight' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Work: { /* 'DPT_Value_Work' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Wsp: { /* 'DPT_Value_Wsp' */
                    return 9;
                }
                case KnxDatapointType.DPT_Value_Wsp_kmh: { /* 'DPT_Value_Wsp_kmh' */
                    return 9;
                }
                case KnxDatapointType.DPT_ValveMode: { /* 'DPT_ValveMode' */
                    return 20;
                }
                case KnxDatapointType.DPT_VarString_8859_1: { /* 'DPT_VarString_8859_1' */
                    return 24;
                }
                case KnxDatapointType.DPT_Version: { /* 'DPT_Version' */
                    return 217;
                }
                case KnxDatapointType.DPT_VolumeLiter_Z: { /* 'DPT_VolumeLiter_Z' */
                    return 218;
                }
                case KnxDatapointType.DPT_WindSpeed_Z_DPT_WindSpeed: { /* 'DPT_WindSpeed_Z_DPT_WindSpeed' */
                    return 203;
                }
                case KnxDatapointType.DPT_Window_Door: { /* 'DPT_Window_Door' */
                    return 1;
                }
                default: {
                    return 0;
                }
            }
        }

        public static ushort SubNumber(this KnxDatapointType value)
        {
            switch (value)
            {
                case KnxDatapointType.DPT_ADAType: { /* 'DPT_ADAType' */
                    return 120;
                }
                case KnxDatapointType.DPT_Access_Data: { /* 'DPT_Access_Data' */
                    return 0;
                }
                case KnxDatapointType.DPT_Ack: { /* 'DPT_Ack' */
                    return 16;
                }
                case KnxDatapointType.DPT_ActPosDemAbs: { /* 'DPT_ActPosDemAbs' */
                    return 104;
                }
                case KnxDatapointType.DPT_ActiveEnergy: { /* 'DPT_ActiveEnergy' */
                    return 10;
                }
                case KnxDatapointType.DPT_ActiveEnergy_V64: { /* 'DPT_ActiveEnergy_V64' */
                    return 10;
                }
                case KnxDatapointType.DPT_ActiveEnergy_kWh: { /* 'DPT_ActiveEnergy_kWh' */
                    return 13;
                }
                case KnxDatapointType.DPT_ActuatorConnectType: { /* 'DPT_ActuatorConnectType' */
                    return 20;
                }
                case KnxDatapointType.DPT_AddInfoTypes: { /* 'DPT_AddInfoTypes' */
                    return 1001;
                }
                case KnxDatapointType.DPT_Alarm: { /* 'DPT_Alarm' */
                    return 5;
                }
                case KnxDatapointType.DPT_AlarmClassType: { /* 'DPT_AlarmClassType' */
                    return 7;
                }
                case KnxDatapointType.DPT_AlarmInfo: { /* 'DPT_AlarmInfo' */
                    return 1;
                }
                case KnxDatapointType.DPT_Alarm_Control: { /* 'DPT_Alarm_Control' */
                    return 5;
                }
                case KnxDatapointType.DPT_Alarm_Reaction: { /* 'DPT_Alarm_Reaction' */
                    return 2;
                }
                case KnxDatapointType.DPT_Angle: { /* 'DPT_Angle' */
                    return 3;
                }
                case KnxDatapointType.DPT_ApparantEnergy: { /* 'DPT_ApparantEnergy' */
                    return 11;
                }
                case KnxDatapointType.DPT_ApparantEnergy_V64: { /* 'DPT_ApparantEnergy_V64' */
                    return 11;
                }
                case KnxDatapointType.DPT_ApparantEnergy_kVAh: { /* 'DPT_ApparantEnergy_kVAh' */
                    return 14;
                }
                case KnxDatapointType.DPT_ApplicationArea: { /* 'DPT_ApplicationArea' */
                    return 6;
                }
                case KnxDatapointType.DPT_AtmPressureAbs_Z: { /* 'DPT_AtmPressureAbs_Z' */
                    return 15;
                }
                case KnxDatapointType.DPT_BackupMode: { /* 'DPT_BackupMode' */
                    return 121;
                }
                case KnxDatapointType.DPT_Beaufort_Wind_Force_Scale: { /* 'DPT_Beaufort_Wind_Force_Scale' */
                    return 14;
                }
                case KnxDatapointType.DPT_Behaviour_Bus_Power_Up_Down: { /* 'DPT_Behaviour_Bus_Power_Up_Down' */
                    return 601;
                }
                case KnxDatapointType.DPT_Behaviour_Lock_Unlock: { /* 'DPT_Behaviour_Lock_Unlock' */
                    return 600;
                }
                case KnxDatapointType.DPT_BinaryValue: { /* 'DPT_BinaryValue' */
                    return 6;
                }
                case KnxDatapointType.DPT_BinaryValue_Control: { /* 'DPT_BinaryValue_Control' */
                    return 6;
                }
                case KnxDatapointType.DPT_BinaryValue_Z: { /* 'DPT_BinaryValue_Z' */
                    return 101;
                }
                case KnxDatapointType.DPT_BlindsControlMode: { /* 'DPT_BlindsControlMode' */
                    return 804;
                }
                case KnxDatapointType.DPT_BlinkingMode: { /* 'DPT_BlinkingMode' */
                    return 603;
                }
                case KnxDatapointType.DPT_Bool: { /* 'DPT_Bool' */
                    return 2;
                }
                case KnxDatapointType.DPT_Bool_Control: { /* 'DPT_Bool_Control' */
                    return 2;
                }
                case KnxDatapointType.DPT_Brightness: { /* 'DPT_Brightness' */
                    return 13;
                }
                case KnxDatapointType.DPT_BuildingMode: { /* 'DPT_BuildingMode' */
                    return 2;
                }
                case KnxDatapointType.DPT_BuildingModeNext: { /* 'DPT_BuildingModeNext' */
                    return 105;
                }
                case KnxDatapointType.DPT_BuildingMode_Z: { /* 'DPT_BuildingMode_Z' */
                    return 107;
                }
                case KnxDatapointType.DPT_BurnerType: { /* 'DPT_BurnerType' */
                    return 101;
                }
                case KnxDatapointType.DPT_ChangeoverMode: { /* 'DPT_ChangeoverMode' */
                    return 107;
                }
                case KnxDatapointType.DPT_Channel_Activation_16: { /* 'DPT_Channel_Activation_16' */
                    return 1010;
                }
                case KnxDatapointType.DPT_Channel_Activation_24: { /* 'DPT_Channel_Activation_24' */
                    return 1010;
                }
                case KnxDatapointType.DPT_Channel_Activation_8: { /* 'DPT_Channel_Activation_8' */
                    return 1010;
                }
                case KnxDatapointType.DPT_Char_8859_1: { /* 'DPT_Char_8859_1' */
                    return 2;
                }
                case KnxDatapointType.DPT_Char_ASCII: { /* 'DPT_Char_ASCII' */
                    return 1;
                }
                case KnxDatapointType.DPT_Colour_RGB: { /* 'DPT_Colour_RGB' */
                    return 600;
                }
                case KnxDatapointType.DPT_CombinedInfoOnOff: { /* 'DPT_CombinedInfoOnOff' */
                    return 1;
                }
                case KnxDatapointType.DPT_CombinedPosition: { /* 'DPT_CombinedPosition' */
                    return 800;
                }
                case KnxDatapointType.DPT_CommMode: { /* 'DPT_CommMode' */
                    return 1000;
                }
                case KnxDatapointType.DPT_Control_Blinds: { /* 'DPT_Control_Blinds' */
                    return 8;
                }
                case KnxDatapointType.DPT_Control_Dimming: { /* 'DPT_Control_Dimming' */
                    return 7;
                }
                case KnxDatapointType.DPT_DALI_Control_Gear_Diagnostic: { /* 'DPT_DALI_Control_Gear_Diagnostic' */
                    return 600;
                }
                case KnxDatapointType.DPT_DALI_Diagnostics: { /* 'DPT_DALI_Diagnostics' */
                    return 600;
                }
                case KnxDatapointType.DPT_DALI_Fade_Time: { /* 'DPT_DALI_Fade_Time' */
                    return 602;
                }
                case KnxDatapointType.DPT_DHWMode: { /* 'DPT_DHWMode' */
                    return 103;
                }
                case KnxDatapointType.DPT_DHWModeNext: { /* 'DPT_DHWModeNext' */
                    return 102;
                }
                case KnxDatapointType.DPT_DHWMode_Z: { /* 'DPT_DHWMode_Z' */
                    return 102;
                }
                case KnxDatapointType.DPT_DamperMode: { /* 'DPT_DamperMode' */
                    return 109;
                }
                case KnxDatapointType.DPT_Date: { /* 'DPT_Date' */
                    return 1;
                }
                case KnxDatapointType.DPT_DateTime: { /* 'DPT_DateTime' */
                    return 1;
                }
                case KnxDatapointType.DPT_DecimalFactor: { /* 'DPT_DecimalFactor' */
                    return 5;
                }
                case KnxDatapointType.DPT_DeltaTime100MSec: { /* 'DPT_DeltaTime100MSec' */
                    return 4;
                }
                case KnxDatapointType.DPT_DeltaTime100Msec_Z: { /* 'DPT_DeltaTime100Msec_Z' */
                    return 4;
                }
                case KnxDatapointType.DPT_DeltaTime10MSec: { /* 'DPT_DeltaTime10MSec' */
                    return 3;
                }
                case KnxDatapointType.DPT_DeltaTime10Msec_Z: { /* 'DPT_DeltaTime10Msec_Z' */
                    return 3;
                }
                case KnxDatapointType.DPT_DeltaTimeHrs: { /* 'DPT_DeltaTimeHrs' */
                    return 7;
                }
                case KnxDatapointType.DPT_DeltaTimeHrs_Z: { /* 'DPT_DeltaTimeHrs_Z' */
                    return 7;
                }
                case KnxDatapointType.DPT_DeltaTimeMin: { /* 'DPT_DeltaTimeMin' */
                    return 6;
                }
                case KnxDatapointType.DPT_DeltaTimeMin_Z: { /* 'DPT_DeltaTimeMin_Z' */
                    return 6;
                }
                case KnxDatapointType.DPT_DeltaTimeMsec: { /* 'DPT_DeltaTimeMsec' */
                    return 2;
                }
                case KnxDatapointType.DPT_DeltaTimeMsec_Z: { /* 'DPT_DeltaTimeMsec_Z' */
                    return 2;
                }
                case KnxDatapointType.DPT_DeltaTimeSec: { /* 'DPT_DeltaTimeSec' */
                    return 5;
                }
                case KnxDatapointType.DPT_DeltaTimeSec_Z: { /* 'DPT_DeltaTimeSec_Z' */
                    return 5;
                }
                case KnxDatapointType.DPT_Device_Control: { /* 'DPT_Device_Control' */
                    return 2;
                }
                case KnxDatapointType.DPT_DimSendStyle: { /* 'DPT_DimSendStyle' */
                    return 13;
                }
                case KnxDatapointType.DPT_DimmPBModel: { /* 'DPT_DimmPBModel' */
                    return 607;
                }
                case KnxDatapointType.DPT_Direction1_Control: { /* 'DPT_Direction1_Control' */
                    return 8;
                }
                case KnxDatapointType.DPT_Direction2_Control: { /* 'DPT_Direction2_Control' */
                    return 9;
                }
                case KnxDatapointType.DPT_DoubleNibble: { /* 'DPT_DoubleNibble' */
                    return 1000;
                }
                case KnxDatapointType.DPT_EnablH_Cstage_Z_DPT_EnablH_CStage: { /* 'DPT_EnablH_Cstage_Z_DPT_EnablH_CStage' */
                    return 105;
                }
                case KnxDatapointType.DPT_Enable: { /* 'DPT_Enable' */
                    return 3;
                }
                case KnxDatapointType.DPT_Enable_Control: { /* 'DPT_Enable_Control' */
                    return 3;
                }
                case KnxDatapointType.DPT_EnergyDemAir: { /* 'DPT_EnergyDemAir' */
                    return 100;
                }
                case KnxDatapointType.DPT_EnergyDemWater: { /* 'DPT_EnergyDemWater' */
                    return 100;
                }
                case KnxDatapointType.DPT_ErrorClass_HVAC: { /* 'DPT_ErrorClass_HVAC' */
                    return 12;
                }
                case KnxDatapointType.DPT_ErrorClass_System: { /* 'DPT_ErrorClass_System' */
                    return 11;
                }
                case KnxDatapointType.DPT_FanMode: { /* 'DPT_FanMode' */
                    return 111;
                }
                case KnxDatapointType.DPT_FlaggedScaling: { /* 'DPT_FlaggedScaling' */
                    return 1;
                }
                case KnxDatapointType.DPT_FlowRate_m3h: { /* 'DPT_FlowRate_m3h' */
                    return 2;
                }
                case KnxDatapointType.DPT_FlowRate_m3h_Z: { /* 'DPT_FlowRate_m3h_Z' */
                    return 2;
                }
                case KnxDatapointType.DPT_ForceSign: { /* 'DPT_ForceSign' */
                    return 100;
                }
                case KnxDatapointType.DPT_ForceSignCool: { /* 'DPT_ForceSignCool' */
                    return 101;
                }
                case KnxDatapointType.DPT_FuelType: { /* 'DPT_FuelType' */
                    return 100;
                }
                case KnxDatapointType.DPT_FuelTypeSet: { /* 'DPT_FuelTypeSet' */
                    return 104;
                }
                case KnxDatapointType.DPT_HVACAirFlowAbs_Z: { /* 'DPT_HVACAirFlowAbs_Z' */
                    return 104;
                }
                case KnxDatapointType.DPT_HVACAirFlowRel_Z: { /* 'DPT_HVACAirFlowRel_Z' */
                    return 102;
                }
                case KnxDatapointType.DPT_HVACAirQual_Z: { /* 'DPT_HVACAirQual_Z' */
                    return 100;
                }
                case KnxDatapointType.DPT_HVACContrMode: { /* 'DPT_HVACContrMode' */
                    return 105;
                }
                case KnxDatapointType.DPT_HVACContrMode_Z: { /* 'DPT_HVACContrMode_Z' */
                    return 104;
                }
                case KnxDatapointType.DPT_HVACEmergMode: { /* 'DPT_HVACEmergMode' */
                    return 106;
                }
                case KnxDatapointType.DPT_HVACEmergMode_Z: { /* 'DPT_HVACEmergMode_Z' */
                    return 109;
                }
                case KnxDatapointType.DPT_HVACMode: { /* 'DPT_HVACMode' */
                    return 102;
                }
                case KnxDatapointType.DPT_HVACModeNext: { /* 'DPT_HVACModeNext' */
                    return 100;
                }
                case KnxDatapointType.DPT_HVACMode_Z: { /* 'DPT_HVACMode_Z' */
                    return 100;
                }
                case KnxDatapointType.DPT_HVAC_PB_Action: { /* 'DPT_HVAC_PB_Action' */
                    return 102;
                }
                case KnxDatapointType.DPT_Heat_Cool: { /* 'DPT_Heat_Cool' */
                    return 100;
                }
                case KnxDatapointType.DPT_Heat_Cool_Z: { /* 'DPT_Heat_Cool_Z' */
                    return 100;
                }
                case KnxDatapointType.DPT_HeaterMode: { /* 'DPT_HeaterMode' */
                    return 110;
                }
                case KnxDatapointType.DPT_InputSource: { /* 'DPT_InputSource' */
                    return 14;
                }
                case KnxDatapointType.DPT_Invert: { /* 'DPT_Invert' */
                    return 12;
                }
                case KnxDatapointType.DPT_Invert_Control: { /* 'DPT_Invert_Control' */
                    return 12;
                }
                case KnxDatapointType.DPT_KelvinPerPercent: { /* 'DPT_KelvinPerPercent' */
                    return 23;
                }
                case KnxDatapointType.DPT_LanguageCodeAlpha2_ASCII: { /* 'DPT_LanguageCodeAlpha2_ASCII' */
                    return 1;
                }
                case KnxDatapointType.DPT_Length_mm: { /* 'DPT_Length_mm' */
                    return 11;
                }
                case KnxDatapointType.DPT_LightActuatorErrorInfo: { /* 'DPT_LightActuatorErrorInfo' */
                    return 601;
                }
                case KnxDatapointType.DPT_LightApplicationMode: { /* 'DPT_LightApplicationMode' */
                    return 5;
                }
                case KnxDatapointType.DPT_LightControlMode: { /* 'DPT_LightControlMode' */
                    return 604;
                }
                case KnxDatapointType.DPT_LoadPriority: { /* 'DPT_LoadPriority' */
                    return 104;
                }
                case KnxDatapointType.DPT_LoadTypeDetected: { /* 'DPT_LoadTypeDetected' */
                    return 610;
                }
                case KnxDatapointType.DPT_LoadTypeSet: { /* 'DPT_LoadTypeSet' */
                    return 609;
                }
                case KnxDatapointType.DPT_Locale_ASCII: { /* 'DPT_Locale_ASCII' */
                    return 1;
                }
                case KnxDatapointType.DPT_LockSign: { /* 'DPT_LockSign' */
                    return 101;
                }
                case KnxDatapointType.DPT_LogicalFunction: { /* 'DPT_LogicalFunction' */
                    return 21;
                }
                case KnxDatapointType.DPT_LongDeltaTimeSec: { /* 'DPT_LongDeltaTimeSec' */
                    return 100;
                }
                case KnxDatapointType.DPT_MBus_Address: { /* 'DPT_MBus_Address' */
                    return 1000;
                }
                case KnxDatapointType.DPT_MasterSlaveMode: { /* 'DPT_MasterSlaveMode' */
                    return 112;
                }
                case KnxDatapointType.DPT_Media: { /* 'DPT_Media' */
                    return 1000;
                }
                case KnxDatapointType.DPT_MeteringValue: { /* 'DPT_MeteringValue' */
                    return 1;
                }
                case KnxDatapointType.DPT_OccMode: { /* 'DPT_OccMode' */
                    return 3;
                }
                case KnxDatapointType.DPT_OccModeNext: { /* 'DPT_OccModeNext' */
                    return 104;
                }
                case KnxDatapointType.DPT_OccMode_Z: { /* 'DPT_OccMode_Z' */
                    return 108;
                }
                case KnxDatapointType.DPT_Occupancy: { /* 'DPT_Occupancy' */
                    return 18;
                }
                case KnxDatapointType.DPT_OnOff_Action: { /* 'DPT_OnOff_Action' */
                    return 1;
                }
                case KnxDatapointType.DPT_OpenClose: { /* 'DPT_OpenClose' */
                    return 9;
                }
                case KnxDatapointType.DPT_PBAction: { /* 'DPT_PBAction' */
                    return 606;
                }
                case KnxDatapointType.DPT_PB_Action_HVAC_Extended: { /* 'DPT_PB_Action_HVAC_Extended' */
                    return 101;
                }
                case KnxDatapointType.DPT_PSUMode: { /* 'DPT_PSUMode' */
                    return 8;
                }
                case KnxDatapointType.DPT_PercentU16_Z: { /* 'DPT_PercentU16_Z' */
                    return 17;
                }
                case KnxDatapointType.DPT_Percent_U8: { /* 'DPT_Percent_U8' */
                    return 4;
                }
                case KnxDatapointType.DPT_Percent_V16: { /* 'DPT_Percent_V16' */
                    return 10;
                }
                case KnxDatapointType.DPT_Percent_V16_Z: { /* 'DPT_Percent_V16_Z' */
                    return 17;
                }
                case KnxDatapointType.DPT_Percent_V8: { /* 'DPT_Percent_V8' */
                    return 1;
                }
                case KnxDatapointType.DPT_Power: { /* 'DPT_Power' */
                    return 24;
                }
                case KnxDatapointType.DPT_PowerDensity: { /* 'DPT_PowerDensity' */
                    return 22;
                }
                case KnxDatapointType.DPT_PowerFlowWaterDemCPM: { /* 'DPT_PowerFlowWaterDemCPM' */
                    return 101;
                }
                case KnxDatapointType.DPT_PowerFlowWaterDemHPM: { /* 'DPT_PowerFlowWaterDemHPM' */
                    return 100;
                }
                case KnxDatapointType.DPT_PowerKW_Z: { /* 'DPT_PowerKW_Z' */
                    return 14;
                }
                case KnxDatapointType.DPT_Prioritised_Mode_Control: { /* 'DPT_Prioritised_Mode_Control' */
                    return 1;
                }
                case KnxDatapointType.DPT_Priority: { /* 'DPT_Priority' */
                    return 4;
                }
                case KnxDatapointType.DPT_PropDataType: { /* 'DPT_PropDataType' */
                    return 10;
                }
                case KnxDatapointType.DPT_RF_FilterInfo: { /* 'DPT_RF_FilterInfo' */
                    return 1001;
                }
                case KnxDatapointType.DPT_RF_FilterSelect: { /* 'DPT_RF_FilterSelect' */
                    return 1003;
                }
                case KnxDatapointType.DPT_RF_ModeInfo: { /* 'DPT_RF_ModeInfo' */
                    return 1000;
                }
                case KnxDatapointType.DPT_RF_ModeSelect: { /* 'DPT_RF_ModeSelect' */
                    return 1002;
                }
                case KnxDatapointType.DPT_Rain_Amount: { /* 'DPT_Rain_Amount' */
                    return 26;
                }
                case KnxDatapointType.DPT_Ramp: { /* 'DPT_Ramp' */
                    return 4;
                }
                case KnxDatapointType.DPT_Ramp_Control: { /* 'DPT_Ramp_Control' */
                    return 4;
                }
                case KnxDatapointType.DPT_ReactiveEnergy: { /* 'DPT_ReactiveEnergy' */
                    return 12;
                }
                case KnxDatapointType.DPT_ReactiveEnergy_V64: { /* 'DPT_ReactiveEnergy_V64' */
                    return 12;
                }
                case KnxDatapointType.DPT_ReactiveEnergy_kVARh: { /* 'DPT_ReactiveEnergy_kVARh' */
                    return 15;
                }
                case KnxDatapointType.DPT_RegionCodeAlpha2_ASCII: { /* 'DPT_RegionCodeAlpha2_ASCII' */
                    return 2;
                }
                case KnxDatapointType.DPT_RelSignedValue_Z: { /* 'DPT_RelSignedValue_Z' */
                    return 1;
                }
                case KnxDatapointType.DPT_RelValue_Z: { /* 'DPT_RelValue_Z' */
                    return 1;
                }
                case KnxDatapointType.DPT_Reset: { /* 'DPT_Reset' */
                    return 15;
                }
                case KnxDatapointType.DPT_Rotation_Angle: { /* 'DPT_Rotation_Angle' */
                    return 11;
                }
                case KnxDatapointType.DPT_SABBehaviour_Lock_Unlock: { /* 'DPT_SABBehaviour_Lock_Unlock' */
                    return 802;
                }
                case KnxDatapointType.DPT_SABExceptBehaviour: { /* 'DPT_SABExceptBehaviour' */
                    return 801;
                }
                case KnxDatapointType.DPT_SCLOMode: { /* 'DPT_SCLOMode' */
                    return 1;
                }
                case KnxDatapointType.DPT_SSSBMode: { /* 'DPT_SSSBMode' */
                    return 803;
                }
                case KnxDatapointType.DPT_Scaling: { /* 'DPT_Scaling' */
                    return 1;
                }
                case KnxDatapointType.DPT_ScalingSpeed: { /* 'DPT_ScalingSpeed' */
                    return 1;
                }
                case KnxDatapointType.DPT_Scaling_Step_Time: { /* 'DPT_Scaling_Step_Time' */
                    return 2;
                }
                case KnxDatapointType.DPT_SceneConfig: { /* 'DPT_SceneConfig' */
                    return 1;
                }
                case KnxDatapointType.DPT_SceneControl: { /* 'DPT_SceneControl' */
                    return 1;
                }
                case KnxDatapointType.DPT_SceneInfo: { /* 'DPT_SceneInfo' */
                    return 1;
                }
                case KnxDatapointType.DPT_SceneNumber: { /* 'DPT_SceneNumber' */
                    return 1;
                }
                case KnxDatapointType.DPT_Scene_AB: { /* 'DPT_Scene_AB' */
                    return 22;
                }
                case KnxDatapointType.DPT_SensorSelect: { /* 'DPT_SensorSelect' */
                    return 17;
                }
                case KnxDatapointType.DPT_SerNum: { /* 'DPT_SerNum' */
                    return 1;
                }
                case KnxDatapointType.DPT_ShutterBlinds_Mode: { /* 'DPT_ShutterBlinds_Mode' */
                    return 23;
                }
                case KnxDatapointType.DPT_SpecHeatProd: { /* 'DPT_SpecHeatProd' */
                    return 100;
                }
                case KnxDatapointType.DPT_Start: { /* 'DPT_Start' */
                    return 10;
                }
                case KnxDatapointType.DPT_StartSynchronization: { /* 'DPT_StartSynchronization' */
                    return 122;
                }
                case KnxDatapointType.DPT_Start_Control: { /* 'DPT_Start_Control' */
                    return 10;
                }
                case KnxDatapointType.DPT_State: { /* 'DPT_State' */
                    return 11;
                }
                case KnxDatapointType.DPT_State_Control: { /* 'DPT_State_Control' */
                    return 11;
                }
                case KnxDatapointType.DPT_StatusAHU: { /* 'DPT_StatusAHU' */
                    return 106;
                }
                case KnxDatapointType.DPT_StatusAct: { /* 'DPT_StatusAct' */
                    return 105;
                }
                case KnxDatapointType.DPT_StatusBOC: { /* 'DPT_StatusBOC' */
                    return 100;
                }
                case KnxDatapointType.DPT_StatusBUC: { /* 'DPT_StatusBUC' */
                    return 100;
                }
                case KnxDatapointType.DPT_StatusCC: { /* 'DPT_StatusCC' */
                    return 101;
                }
                case KnxDatapointType.DPT_StatusCPM: { /* 'DPT_StatusCPM' */
                    return 102;
                }
                case KnxDatapointType.DPT_StatusDHWC: { /* 'DPT_StatusDHWC' */
                    return 100;
                }
                case KnxDatapointType.DPT_StatusGen: { /* 'DPT_StatusGen' */
                    return 1;
                }
                case KnxDatapointType.DPT_StatusHPM: { /* 'DPT_StatusHPM' */
                    return 100;
                }
                case KnxDatapointType.DPT_StatusLightingActuator: { /* 'DPT_StatusLightingActuator' */
                    return 600;
                }
                case KnxDatapointType.DPT_StatusRCC: { /* 'DPT_StatusRCC' */
                    return 105;
                }
                case KnxDatapointType.DPT_StatusRHC: { /* 'DPT_StatusRHC' */
                    return 102;
                }
                case KnxDatapointType.DPT_StatusRHCC: { /* 'DPT_StatusRHCC' */
                    return 101;
                }
                case KnxDatapointType.DPT_StatusRoomSetp: { /* 'DPT_StatusRoomSetp' */
                    return 113;
                }
                case KnxDatapointType.DPT_StatusSAB: { /* 'DPT_StatusSAB' */
                    return 800;
                }
                case KnxDatapointType.DPT_StatusSDHWC: { /* 'DPT_StatusSDHWC' */
                    return 103;
                }
                case KnxDatapointType.DPT_StatusWTC: { /* 'DPT_StatusWTC' */
                    return 103;
                }
                case KnxDatapointType.DPT_Status_Mode3: { /* 'DPT_Status_Mode3' */
                    return 20;
                }
                case KnxDatapointType.DPT_Step: { /* 'DPT_Step' */
                    return 7;
                }
                case KnxDatapointType.DPT_Step_Control: { /* 'DPT_Step_Control' */
                    return 7;
                }
                case KnxDatapointType.DPT_String_8859_1: { /* 'DPT_String_8859_1' */
                    return 1;
                }
                case KnxDatapointType.DPT_String_ASCII: { /* 'DPT_String_ASCII' */
                    return 0;
                }
                case KnxDatapointType.DPT_SunIntensity_Z: { /* 'DPT_SunIntensity_Z' */
                    return 102;
                }
                case KnxDatapointType.DPT_Switch: { /* 'DPT_Switch' */
                    return 1;
                }
                case KnxDatapointType.DPT_SwitchOnMode: { /* 'DPT_SwitchOnMode' */
                    return 608;
                }
                case KnxDatapointType.DPT_SwitchPBModel: { /* 'DPT_SwitchPBModel' */
                    return 605;
                }
                case KnxDatapointType.DPT_Switch_Control: { /* 'DPT_Switch_Control' */
                    return 1;
                }
                case KnxDatapointType.DPT_Tariff: { /* 'DPT_Tariff' */
                    return 6;
                }
                case KnxDatapointType.DPT_TariffNext: { /* 'DPT_TariffNext' */
                    return 3;
                }
                case KnxDatapointType.DPT_Tariff_ActiveEnergy: { /* 'DPT_Tariff_ActiveEnergy' */
                    return 1;
                }
                case KnxDatapointType.DPT_TempDHWSetpSet4: { /* 'DPT_TempDHWSetpSet4' */
                    return 101;
                }
                case KnxDatapointType.DPT_TempFlowWaterDemAbs: { /* 'DPT_TempFlowWaterDemAbs' */
                    return 100;
                }
                case KnxDatapointType.DPT_TempHVACAbsNext: { /* 'DPT_TempHVACAbsNext' */
                    return 100;
                }
                case KnxDatapointType.DPT_TempHVACAbs_Z: { /* 'DPT_TempHVACAbs_Z' */
                    return 100;
                }
                case KnxDatapointType.DPT_TempHVACRel_Z: { /* 'DPT_TempHVACRel_Z' */
                    return 101;
                }
                case KnxDatapointType.DPT_TempRoomDemAbs: { /* 'DPT_TempRoomDemAbs' */
                    return 101;
                }
                case KnxDatapointType.DPT_TempRoomSetpSet3: { /* 'DPT_TempRoomSetpSet3' */
                    return 101;
                }
                case KnxDatapointType.DPT_TempRoomSetpSet4: { /* 'DPT_TempRoomSetpSet4' */
                    return 100;
                }
                case KnxDatapointType.DPT_TempRoomSetpSetF163: { /* 'DPT_TempRoomSetpSetF163' */
                    return 100;
                }
                case KnxDatapointType.DPT_TempRoomSetpSetShift3: { /* 'DPT_TempRoomSetpSetShift3' */
                    return 100;
                }
                case KnxDatapointType.DPT_TempRoomSetpSetShift4: { /* 'DPT_TempRoomSetpSetShift4' */
                    return 102;
                }
                case KnxDatapointType.DPT_TempRoomSetpSetShiftF163: { /* 'DPT_TempRoomSetpSetShiftF163' */
                    return 101;
                }
                case KnxDatapointType.DPT_TempSupply_AirSetpSet: { /* 'DPT_TempSupply_AirSetpSet' */
                    return 100;
                }
                case KnxDatapointType.DPT_TimeOfDay: { /* 'DPT_TimeOfDay' */
                    return 1;
                }
                case KnxDatapointType.DPT_TimePeriod100MSec: { /* 'DPT_TimePeriod100MSec' */
                    return 4;
                }
                case KnxDatapointType.DPT_TimePeriod100Msec_Z: { /* 'DPT_TimePeriod100Msec_Z' */
                    return 4;
                }
                case KnxDatapointType.DPT_TimePeriod10MSec: { /* 'DPT_TimePeriod10MSec' */
                    return 3;
                }
                case KnxDatapointType.DPT_TimePeriod10Msec_Z: { /* 'DPT_TimePeriod10Msec_Z' */
                    return 3;
                }
                case KnxDatapointType.DPT_TimePeriodHrs: { /* 'DPT_TimePeriodHrs' */
                    return 7;
                }
                case KnxDatapointType.DPT_TimePeriodHrs_Z: { /* 'DPT_TimePeriodHrs_Z' */
                    return 7;
                }
                case KnxDatapointType.DPT_TimePeriodMin: { /* 'DPT_TimePeriodMin' */
                    return 6;
                }
                case KnxDatapointType.DPT_TimePeriodMin_Z: { /* 'DPT_TimePeriodMin_Z' */
                    return 6;
                }
                case KnxDatapointType.DPT_TimePeriodMsec: { /* 'DPT_TimePeriodMsec' */
                    return 2;
                }
                case KnxDatapointType.DPT_TimePeriodMsec_Z: { /* 'DPT_TimePeriodMsec_Z' */
                    return 2;
                }
                case KnxDatapointType.DPT_TimePeriodSec: { /* 'DPT_TimePeriodSec' */
                    return 5;
                }
                case KnxDatapointType.DPT_TimePeriodSec_Z: { /* 'DPT_TimePeriodSec_Z' */
                    return 5;
                }
                case KnxDatapointType.DPT_Time_Delay: { /* 'DPT_Time_Delay' */
                    return 13;
                }
                case KnxDatapointType.DPT_Trigger: { /* 'DPT_Trigger' */
                    return 17;
                }
                case KnxDatapointType.DPT_UCountValue16_Z: { /* 'DPT_UCountValue16_Z' */
                    return 12;
                }
                case KnxDatapointType.DPT_UCountValue8_Z: { /* 'DPT_UCountValue8_Z' */
                    return 2;
                }
                case KnxDatapointType.DPT_UElCurrentmA: { /* 'DPT_UElCurrentmA' */
                    return 12;
                }
                case KnxDatapointType.DPT_UElCurrentyA_Z: { /* 'DPT_UElCurrentyA_Z' */
                    return 13;
                }
                case KnxDatapointType.DPT_UFlowRateLiter_h_Z: { /* 'DPT_UFlowRateLiter_h_Z' */
                    return 11;
                }
                case KnxDatapointType.DPT_UTF_8: { /* 'DPT_UTF_8' */
                    return 1;
                }
                case KnxDatapointType.DPT_UpDown: { /* 'DPT_UpDown' */
                    return 8;
                }
                case KnxDatapointType.DPT_UpDown_Action: { /* 'DPT_UpDown_Action' */
                    return 3;
                }
                case KnxDatapointType.DPT_ValueDemBOC: { /* 'DPT_ValueDemBOC' */
                    return 102;
                }
                case KnxDatapointType.DPT_Value_1_Count: { /* 'DPT_Value_1_Count' */
                    return 10;
                }
                case KnxDatapointType.DPT_Value_1_Ucount: { /* 'DPT_Value_1_Ucount' */
                    return 10;
                }
                case KnxDatapointType.DPT_Value_2_Count: { /* 'DPT_Value_2_Count' */
                    return 1;
                }
                case KnxDatapointType.DPT_Value_2_Ucount: { /* 'DPT_Value_2_Ucount' */
                    return 1;
                }
                case KnxDatapointType.DPT_Value_4_Count: { /* 'DPT_Value_4_Count' */
                    return 1;
                }
                case KnxDatapointType.DPT_Value_4_Ucount: { /* 'DPT_Value_4_Ucount' */
                    return 1;
                }
                case KnxDatapointType.DPT_Value_Absolute_Temperature: { /* 'DPT_Value_Absolute_Temperature' */
                    return 69;
                }
                case KnxDatapointType.DPT_Value_Acceleration: { /* 'DPT_Value_Acceleration' */
                    return 0;
                }
                case KnxDatapointType.DPT_Value_Acceleration_Angular: { /* 'DPT_Value_Acceleration_Angular' */
                    return 1;
                }
                case KnxDatapointType.DPT_Value_Activation_Energy: { /* 'DPT_Value_Activation_Energy' */
                    return 2;
                }
                case KnxDatapointType.DPT_Value_Activity: { /* 'DPT_Value_Activity' */
                    return 3;
                }
                case KnxDatapointType.DPT_Value_AirQuality: { /* 'DPT_Value_AirQuality' */
                    return 8;
                }
                case KnxDatapointType.DPT_Value_Amplitude: { /* 'DPT_Value_Amplitude' */
                    return 5;
                }
                case KnxDatapointType.DPT_Value_AngleDeg: { /* 'DPT_Value_AngleDeg' */
                    return 7;
                }
                case KnxDatapointType.DPT_Value_AngleRad: { /* 'DPT_Value_AngleRad' */
                    return 6;
                }
                case KnxDatapointType.DPT_Value_Angular_Frequency: { /* 'DPT_Value_Angular_Frequency' */
                    return 34;
                }
                case KnxDatapointType.DPT_Value_Angular_Momentum: { /* 'DPT_Value_Angular_Momentum' */
                    return 8;
                }
                case KnxDatapointType.DPT_Value_Angular_Velocity: { /* 'DPT_Value_Angular_Velocity' */
                    return 9;
                }
                case KnxDatapointType.DPT_Value_Area: { /* 'DPT_Value_Area' */
                    return 10;
                }
                case KnxDatapointType.DPT_Value_Capacitance: { /* 'DPT_Value_Capacitance' */
                    return 11;
                }
                case KnxDatapointType.DPT_Value_Charge_DensitySurface: { /* 'DPT_Value_Charge_DensitySurface' */
                    return 12;
                }
                case KnxDatapointType.DPT_Value_Charge_DensityVolume: { /* 'DPT_Value_Charge_DensityVolume' */
                    return 13;
                }
                case KnxDatapointType.DPT_Value_Common_Temperature: { /* 'DPT_Value_Common_Temperature' */
                    return 68;
                }
                case KnxDatapointType.DPT_Value_Compressibility: { /* 'DPT_Value_Compressibility' */
                    return 14;
                }
                case KnxDatapointType.DPT_Value_Conductance: { /* 'DPT_Value_Conductance' */
                    return 15;
                }
                case KnxDatapointType.DPT_Value_Curr: { /* 'DPT_Value_Curr' */
                    return 21;
                }
                case KnxDatapointType.DPT_Value_Density: { /* 'DPT_Value_Density' */
                    return 17;
                }
                case KnxDatapointType.DPT_Value_Electric_Charge: { /* 'DPT_Value_Electric_Charge' */
                    return 18;
                }
                case KnxDatapointType.DPT_Value_Electric_Current: { /* 'DPT_Value_Electric_Current' */
                    return 19;
                }
                case KnxDatapointType.DPT_Value_Electric_CurrentDensity: { /* 'DPT_Value_Electric_CurrentDensity' */
                    return 20;
                }
                case KnxDatapointType.DPT_Value_Electric_DipoleMoment: { /* 'DPT_Value_Electric_DipoleMoment' */
                    return 21;
                }
                case KnxDatapointType.DPT_Value_Electric_Displacement: { /* 'DPT_Value_Electric_Displacement' */
                    return 22;
                }
                case KnxDatapointType.DPT_Value_Electric_FieldStrength: { /* 'DPT_Value_Electric_FieldStrength' */
                    return 23;
                }
                case KnxDatapointType.DPT_Value_Electric_Flux: { /* 'DPT_Value_Electric_Flux' */
                    return 24;
                }
                case KnxDatapointType.DPT_Value_Electric_FluxDensity: { /* 'DPT_Value_Electric_FluxDensity' */
                    return 25;
                }
                case KnxDatapointType.DPT_Value_Electric_Polarization: { /* 'DPT_Value_Electric_Polarization' */
                    return 26;
                }
                case KnxDatapointType.DPT_Value_Electric_Potential: { /* 'DPT_Value_Electric_Potential' */
                    return 27;
                }
                case KnxDatapointType.DPT_Value_Electric_PotentialDifference: { /* 'DPT_Value_Electric_PotentialDifference' */
                    return 28;
                }
                case KnxDatapointType.DPT_Value_Electrical_Conductivity: { /* 'DPT_Value_Electrical_Conductivity' */
                    return 16;
                }
                case KnxDatapointType.DPT_Value_ElectromagneticMoment: { /* 'DPT_Value_ElectromagneticMoment' */
                    return 29;
                }
                case KnxDatapointType.DPT_Value_Electromotive_Force: { /* 'DPT_Value_Electromotive_Force' */
                    return 30;
                }
                case KnxDatapointType.DPT_Value_Energy: { /* 'DPT_Value_Energy' */
                    return 31;
                }
                case KnxDatapointType.DPT_Value_Force: { /* 'DPT_Value_Force' */
                    return 32;
                }
                case KnxDatapointType.DPT_Value_Frequency: { /* 'DPT_Value_Frequency' */
                    return 33;
                }
                case KnxDatapointType.DPT_Value_Heat_Capacity: { /* 'DPT_Value_Heat_Capacity' */
                    return 35;
                }
                case KnxDatapointType.DPT_Value_Heat_FlowRate: { /* 'DPT_Value_Heat_FlowRate' */
                    return 36;
                }
                case KnxDatapointType.DPT_Value_Heat_Quantity: { /* 'DPT_Value_Heat_Quantity' */
                    return 37;
                }
                case KnxDatapointType.DPT_Value_Humidity: { /* 'DPT_Value_Humidity' */
                    return 7;
                }
                case KnxDatapointType.DPT_Value_Impedance: { /* 'DPT_Value_Impedance' */
                    return 38;
                }
                case KnxDatapointType.DPT_Value_Length: { /* 'DPT_Value_Length' */
                    return 39;
                }
                case KnxDatapointType.DPT_Value_Light_Quantity: { /* 'DPT_Value_Light_Quantity' */
                    return 40;
                }
                case KnxDatapointType.DPT_Value_Luminance: { /* 'DPT_Value_Luminance' */
                    return 41;
                }
                case KnxDatapointType.DPT_Value_Luminous_Flux: { /* 'DPT_Value_Luminous_Flux' */
                    return 42;
                }
                case KnxDatapointType.DPT_Value_Luminous_Intensity: { /* 'DPT_Value_Luminous_Intensity' */
                    return 43;
                }
                case KnxDatapointType.DPT_Value_Lux: { /* 'DPT_Value_Lux' */
                    return 4;
                }
                case KnxDatapointType.DPT_Value_Magnetic_FieldStrength: { /* 'DPT_Value_Magnetic_FieldStrength' */
                    return 44;
                }
                case KnxDatapointType.DPT_Value_Magnetic_Flux: { /* 'DPT_Value_Magnetic_Flux' */
                    return 45;
                }
                case KnxDatapointType.DPT_Value_Magnetic_FluxDensity: { /* 'DPT_Value_Magnetic_FluxDensity' */
                    return 46;
                }
                case KnxDatapointType.DPT_Value_Magnetic_Moment: { /* 'DPT_Value_Magnetic_Moment' */
                    return 47;
                }
                case KnxDatapointType.DPT_Value_Magnetic_Polarization: { /* 'DPT_Value_Magnetic_Polarization' */
                    return 48;
                }
                case KnxDatapointType.DPT_Value_Magnetization: { /* 'DPT_Value_Magnetization' */
                    return 49;
                }
                case KnxDatapointType.DPT_Value_MagnetomotiveForce: { /* 'DPT_Value_MagnetomotiveForce' */
                    return 50;
                }
                case KnxDatapointType.DPT_Value_Mass: { /* 'DPT_Value_Mass' */
                    return 51;
                }
                case KnxDatapointType.DPT_Value_MassFlux: { /* 'DPT_Value_MassFlux' */
                    return 52;
                }
                case KnxDatapointType.DPT_Value_Mol: { /* 'DPT_Value_Mol' */
                    return 4;
                }
                case KnxDatapointType.DPT_Value_Momentum: { /* 'DPT_Value_Momentum' */
                    return 53;
                }
                case KnxDatapointType.DPT_Value_Phase_AngleDeg: { /* 'DPT_Value_Phase_AngleDeg' */
                    return 55;
                }
                case KnxDatapointType.DPT_Value_Phase_AngleRad: { /* 'DPT_Value_Phase_AngleRad' */
                    return 54;
                }
                case KnxDatapointType.DPT_Value_Power: { /* 'DPT_Value_Power' */
                    return 56;
                }
                case KnxDatapointType.DPT_Value_Power_Factor: { /* 'DPT_Value_Power_Factor' */
                    return 57;
                }
                case KnxDatapointType.DPT_Value_Pres: { /* 'DPT_Value_Pres' */
                    return 6;
                }
                case KnxDatapointType.DPT_Value_Pressure: { /* 'DPT_Value_Pressure' */
                    return 58;
                }
                case KnxDatapointType.DPT_Value_Reactance: { /* 'DPT_Value_Reactance' */
                    return 59;
                }
                case KnxDatapointType.DPT_Value_Resistance: { /* 'DPT_Value_Resistance' */
                    return 60;
                }
                case KnxDatapointType.DPT_Value_Resistivity: { /* 'DPT_Value_Resistivity' */
                    return 61;
                }
                case KnxDatapointType.DPT_Value_SelfInductance: { /* 'DPT_Value_SelfInductance' */
                    return 62;
                }
                case KnxDatapointType.DPT_Value_SolidAngle: { /* 'DPT_Value_SolidAngle' */
                    return 63;
                }
                case KnxDatapointType.DPT_Value_Sound_Intensity: { /* 'DPT_Value_Sound_Intensity' */
                    return 64;
                }
                case KnxDatapointType.DPT_Value_Speed: { /* 'DPT_Value_Speed' */
                    return 65;
                }
                case KnxDatapointType.DPT_Value_Stress: { /* 'DPT_Value_Stress' */
                    return 66;
                }
                case KnxDatapointType.DPT_Value_Surface_Tension: { /* 'DPT_Value_Surface_Tension' */
                    return 67;
                }
                case KnxDatapointType.DPT_Value_Temp: { /* 'DPT_Value_Temp' */
                    return 1;
                }
                case KnxDatapointType.DPT_Value_Temp_F: { /* 'DPT_Value_Temp_F' */
                    return 27;
                }
                case KnxDatapointType.DPT_Value_Tempa: { /* 'DPT_Value_Tempa' */
                    return 3;
                }
                case KnxDatapointType.DPT_Value_Tempd: { /* 'DPT_Value_Tempd' */
                    return 2;
                }
                case KnxDatapointType.DPT_Value_TemperatureDifference: { /* 'DPT_Value_TemperatureDifference' */
                    return 70;
                }
                case KnxDatapointType.DPT_Value_Thermal_Capacity: { /* 'DPT_Value_Thermal_Capacity' */
                    return 71;
                }
                case KnxDatapointType.DPT_Value_Thermal_Conductivity: { /* 'DPT_Value_Thermal_Conductivity' */
                    return 72;
                }
                case KnxDatapointType.DPT_Value_ThermoelectricPower: { /* 'DPT_Value_ThermoelectricPower' */
                    return 73;
                }
                case KnxDatapointType.DPT_Value_Time: { /* 'DPT_Value_Time' */
                    return 74;
                }
                case KnxDatapointType.DPT_Value_Time1: { /* 'DPT_Value_Time1' */
                    return 10;
                }
                case KnxDatapointType.DPT_Value_Time2: { /* 'DPT_Value_Time2' */
                    return 11;
                }
                case KnxDatapointType.DPT_Value_Torque: { /* 'DPT_Value_Torque' */
                    return 75;
                }
                case KnxDatapointType.DPT_Value_Volt: { /* 'DPT_Value_Volt' */
                    return 20;
                }
                case KnxDatapointType.DPT_Value_Volume: { /* 'DPT_Value_Volume' */
                    return 76;
                }
                case KnxDatapointType.DPT_Value_Volume_Flow: { /* 'DPT_Value_Volume_Flow' */
                    return 25;
                }
                case KnxDatapointType.DPT_Value_Volume_Flux: { /* 'DPT_Value_Volume_Flux' */
                    return 77;
                }
                case KnxDatapointType.DPT_Value_Weight: { /* 'DPT_Value_Weight' */
                    return 78;
                }
                case KnxDatapointType.DPT_Value_Work: { /* 'DPT_Value_Work' */
                    return 79;
                }
                case KnxDatapointType.DPT_Value_Wsp: { /* 'DPT_Value_Wsp' */
                    return 5;
                }
                case KnxDatapointType.DPT_Value_Wsp_kmh: { /* 'DPT_Value_Wsp_kmh' */
                    return 28;
                }
                case KnxDatapointType.DPT_ValveMode: { /* 'DPT_ValveMode' */
                    return 108;
                }
                case KnxDatapointType.DPT_VarString_8859_1: { /* 'DPT_VarString_8859_1' */
                    return 1;
                }
                case KnxDatapointType.DPT_Version: { /* 'DPT_Version' */
                    return 1;
                }
                case KnxDatapointType.DPT_VolumeLiter_Z: { /* 'DPT_VolumeLiter_Z' */
                    return 1;
                }
                case KnxDatapointType.DPT_WindSpeed_Z_DPT_WindSpeed: { /* 'DPT_WindSpeed_Z_DPT_WindSpeed' */
                    return 101;
                }
                case KnxDatapointType.DPT_Window_Door: { /* 'DPT_Window_Door' */
                    return 19;
                }
                default: {
                    return 0;
                }
            }
        }
    }

}

