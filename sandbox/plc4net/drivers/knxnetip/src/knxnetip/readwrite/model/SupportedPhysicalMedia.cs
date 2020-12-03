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

    public static class SupportedPhysicalMediaInfo
    {

        public static bool KnxSupport(this SupportedPhysicalMedia value)
        {
            switch (value)
            {
                case SupportedPhysicalMedia.OTHER: { /* '0x00' */
                    return true;
                }
                case SupportedPhysicalMedia.OIL_METER: { /* '0x01' */
                    return true;
                }
                case SupportedPhysicalMedia.ELECTRICITY_METER: { /* '0x02' */
                    return true;
                }
                case SupportedPhysicalMedia.GAS_METER: { /* '0x03' */
                    return true;
                }
                case SupportedPhysicalMedia.HEAT_METER: { /* '0x04' */
                    return true;
                }
                case SupportedPhysicalMedia.STEAM_METER: { /* '0x05' */
                    return true;
                }
                case SupportedPhysicalMedia.WARM_WATER_METER: { /* '0x06' */
                    return true;
                }
                case SupportedPhysicalMedia.WATER_METER: { /* '0x07' */
                    return true;
                }
                case SupportedPhysicalMedia.HEAT_COST_ALLOCATOR: { /* '0x08' */
                    return true;
                }
                case SupportedPhysicalMedia.COMPRESSED_AIR: { /* '0x09' */
                    return false;
                }
                case SupportedPhysicalMedia.COOLING_LOAD_METER_INLET: { /* '0x0A' */
                    return true;
                }
                case SupportedPhysicalMedia.COOLING_LOAD_METER_OUTLET: { /* '0x0B' */
                    return true;
                }
                case SupportedPhysicalMedia.HEAT_INLET: { /* '0x0C' */
                    return true;
                }
                case SupportedPhysicalMedia.HEAT_AND_COOL: { /* '0x0D' */
                    return true;
                }
                case SupportedPhysicalMedia.BUS_OR_SYSTEM: { /* '0x0E' */
                    return false;
                }
                case SupportedPhysicalMedia.UNKNOWN_DEVICE_TYPE: { /* '0x0F' */
                    return false;
                }
                case SupportedPhysicalMedia.BREAKER: { /* '0x20' */
                    return true;
                }
                case SupportedPhysicalMedia.VALVE: { /* '0x21' */
                    return true;
                }
                case SupportedPhysicalMedia.WASTE_WATER_METER: { /* '0x28' */
                    return true;
                }
                case SupportedPhysicalMedia.GARBAGE: { /* '0x29' */
                    return true;
                }
                case SupportedPhysicalMedia.RADIO_CONVERTER: { /* '0x37' */
                    return false;
                }
                default: {
                    return false;
                }
            }
        }

        public static string Description(this SupportedPhysicalMedia value)
        {
            switch (value)
            {
                case SupportedPhysicalMedia.OTHER: { /* '0x00' */
                    return "used_for_undefined_physical_medium";
                }
                case SupportedPhysicalMedia.OIL_METER: { /* '0x01' */
                    return "measures_volume_of_oil";
                }
                case SupportedPhysicalMedia.ELECTRICITY_METER: { /* '0x02' */
                    return "measures_electric_energy";
                }
                case SupportedPhysicalMedia.GAS_METER: { /* '0x03' */
                    return "measures_volume_of_gaseous_energy";
                }
                case SupportedPhysicalMedia.HEAT_METER: { /* '0x04' */
                    return "heat_energy_measured_in_outlet_pipe";
                }
                case SupportedPhysicalMedia.STEAM_METER: { /* '0x05' */
                    return "measures_weight_of_hot_steam";
                }
                case SupportedPhysicalMedia.WARM_WATER_METER: { /* '0x06' */
                    return "measured_heated_water_volume";
                }
                case SupportedPhysicalMedia.WATER_METER: { /* '0x07' */
                    return "measured_water_volume";
                }
                case SupportedPhysicalMedia.HEAT_COST_ALLOCATOR: { /* '0x08' */
                    return "measured_relative_cumulated_heat_consumption";
                }
                case SupportedPhysicalMedia.COMPRESSED_AIR: { /* '0x09' */
                    return "measures_weight_of_compressed_air";
                }
                case SupportedPhysicalMedia.COOLING_LOAD_METER_INLET: { /* '0x0A' */
                    return "cooling_energy_measured_in_inlet_pipe";
                }
                case SupportedPhysicalMedia.COOLING_LOAD_METER_OUTLET: { /* '0x0B' */
                    return "cooling_energy_measured_in_outlet_pipe";
                }
                case SupportedPhysicalMedia.HEAT_INLET: { /* '0x0C' */
                    return "heat_energy_measured_in_inlet_pipe";
                }
                case SupportedPhysicalMedia.HEAT_AND_COOL: { /* '0x0D' */
                    return "measures_both_heat_and_cool";
                }
                case SupportedPhysicalMedia.BUS_OR_SYSTEM: { /* '0x0E' */
                    return "no_meter";
                }
                case SupportedPhysicalMedia.UNKNOWN_DEVICE_TYPE: { /* '0x0F' */
                    return "used_for_undefined_physical_medium";
                }
                case SupportedPhysicalMedia.BREAKER: { /* '0x20' */
                    return "status_of_electric_energy_supply";
                }
                case SupportedPhysicalMedia.VALVE: { /* '0x21' */
                    return "status_of_supply_of_Gas_or_water";
                }
                case SupportedPhysicalMedia.WASTE_WATER_METER: { /* '0x28' */
                    return "measured_volume_of_disposed_water";
                }
                case SupportedPhysicalMedia.GARBAGE: { /* '0x29' */
                    return "measured_weight_of_disposed_rubbish";
                }
                case SupportedPhysicalMedia.RADIO_CONVERTER: { /* '0x37' */
                    return "enables_the_radio_transmission_of_a_meter_without_a_radio_interface";
                }
                default: {
                    return "";
                }
            }
        }
    }

}

