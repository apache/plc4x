//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

// Code generated from the mspec by plc4net-code-gen. DO NOT EDIT.

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{
    public enum SupportedPhysicalMedia : byte
    {
        OTHER = 0x00,
        OIL_METER = 0x01,
        ELECTRICITY_METER = 0x02,
        GAS_METER = 0x03,
        HEAT_METER = 0x04,
        STEAM_METER = 0x05,
        WARM_WATER_METER = 0x06,
        WATER_METER = 0x07,
        HEAT_COST_ALLOCATOR = 0x08,
        COMPRESSED_AIR = 0x09,
        COOLING_LOAD_METER_INLET = 0x0A,
        COOLING_LOAD_METER_OUTLET = 0x0B,
        HEAT_INLET = 0x0C,
        HEAT_AND_COOL = 0x0D,
        BUS_OR_SYSTEM = 0x0E,
        UNKNOWN_DEVICE_TYPE = 0x0F,
        BREAKER = 0x20,
        VALVE = 0x21,
        WASTE_WATER_METER = 0x28,
        GARBAGE = 0x29,
        RADIO_CONVERTER = 0x37,
    }

    public static class SupportedPhysicalMediaExtensions
    {
        public static string GetDescription(this SupportedPhysicalMedia value) => value switch
        {
            SupportedPhysicalMedia.OTHER => "used_for_undefined_physical_medium",
            SupportedPhysicalMedia.OIL_METER => "measures_volume_of_oil",
            SupportedPhysicalMedia.ELECTRICITY_METER => "measures_electric_energy",
            SupportedPhysicalMedia.GAS_METER => "measures_volume_of_gaseous_energy",
            SupportedPhysicalMedia.HEAT_METER => "heat_energy_measured_in_outlet_pipe",
            SupportedPhysicalMedia.STEAM_METER => "measures_weight_of_hot_steam",
            SupportedPhysicalMedia.WARM_WATER_METER => "measured_heated_water_volume",
            SupportedPhysicalMedia.WATER_METER => "measured_water_volume",
            SupportedPhysicalMedia.HEAT_COST_ALLOCATOR => "measured_relative_cumulated_heat_consumption",
            SupportedPhysicalMedia.COMPRESSED_AIR => "measures_weight_of_compressed_air",
            SupportedPhysicalMedia.COOLING_LOAD_METER_INLET => "cooling_energy_measured_in_inlet_pipe",
            SupportedPhysicalMedia.COOLING_LOAD_METER_OUTLET => "cooling_energy_measured_in_outlet_pipe",
            SupportedPhysicalMedia.HEAT_INLET => "heat_energy_measured_in_inlet_pipe",
            SupportedPhysicalMedia.HEAT_AND_COOL => "measures_both_heat_and_cool",
            SupportedPhysicalMedia.BUS_OR_SYSTEM => "no_meter",
            SupportedPhysicalMedia.UNKNOWN_DEVICE_TYPE => "used_for_undefined_physical_medium",
            SupportedPhysicalMedia.BREAKER => "status_of_electric_energy_supply",
            SupportedPhysicalMedia.VALVE => "status_of_supply_of_Gas_or_water",
            SupportedPhysicalMedia.WASTE_WATER_METER => "measured_volume_of_disposed_water",
            SupportedPhysicalMedia.GARBAGE => "measured_weight_of_disposed_rubbish",
            SupportedPhysicalMedia.RADIO_CONVERTER => "enables_the_radio_transmission_of_a_meter_without_a_radio_interface",
            _ => default,
        };
        public static bool GetKnxSupport(this SupportedPhysicalMedia value) => value switch
        {
            SupportedPhysicalMedia.OTHER => true,
            SupportedPhysicalMedia.OIL_METER => true,
            SupportedPhysicalMedia.ELECTRICITY_METER => true,
            SupportedPhysicalMedia.GAS_METER => true,
            SupportedPhysicalMedia.HEAT_METER => true,
            SupportedPhysicalMedia.STEAM_METER => true,
            SupportedPhysicalMedia.WARM_WATER_METER => true,
            SupportedPhysicalMedia.WATER_METER => true,
            SupportedPhysicalMedia.HEAT_COST_ALLOCATOR => true,
            SupportedPhysicalMedia.COMPRESSED_AIR => false,
            SupportedPhysicalMedia.COOLING_LOAD_METER_INLET => true,
            SupportedPhysicalMedia.COOLING_LOAD_METER_OUTLET => true,
            SupportedPhysicalMedia.HEAT_INLET => true,
            SupportedPhysicalMedia.HEAT_AND_COOL => true,
            SupportedPhysicalMedia.BUS_OR_SYSTEM => false,
            SupportedPhysicalMedia.UNKNOWN_DEVICE_TYPE => false,
            SupportedPhysicalMedia.BREAKER => true,
            SupportedPhysicalMedia.VALVE => true,
            SupportedPhysicalMedia.WASTE_WATER_METER => true,
            SupportedPhysicalMedia.GARBAGE => true,
            SupportedPhysicalMedia.RADIO_CONVERTER => false,
            _ => default,
        };
    }
}
