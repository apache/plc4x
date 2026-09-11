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
    public enum KnxInterfaceObjectType : ushort
    {
        OT_UNKNOWN = 0,
        OT_GENERAL = 1,
        OT_DEVICE = 2,
        OT_ADDRESS_TABLE = 3,
        OT_ASSOCIATION_TABLE = 4,
        OT_APPLICATION_PROGRAM = 5,
        OT_INTERACE_PROGRAM = 6,
        OT_EIBOBJECT_ASSOCIATATION_TABLE = 7,
        OT_ROUTER = 8,
        OT_LTE_ADDRESS_ROUTING_TABLE = 9,
        OT_CEMI_SERVER = 10,
        OT_GROUP_OBJECT_TABLE = 11,
        OT_POLLING_MASTER = 12,
        OT_KNXIP_PARAMETER = 13,
        OT_FILE_SERVER = 14,
        OT_SECURITY = 15,
        OT_RF_MEDIUM = 16,
        OT_INDOOR_BRIGHTNESS_SENSOR = 17,
        OT_INDOOR_LUMINANCE_SENSOR = 18,
        OT_LIGHT_SWITCHING_ACTUATOR_BASIC = 19,
        OT_DIMMING_ACTUATOR_BASIC = 20,
        OT_DIMMING_SENSOR_BASIC = 21,
        OT_SWITCHING_SENSOR_BASIC = 22,
        OT_SUNBLIND_ACTUATOR_BASIC = 23,
        OT_SUNBLIND_SENSOR_BASIC = 24,
    }

    public static class KnxInterfaceObjectTypeExtensions
    {
        public static string GetCode(this KnxInterfaceObjectType value) => value switch
        {
            KnxInterfaceObjectType.OT_UNKNOWN => "U",
            KnxInterfaceObjectType.OT_GENERAL => "G",
            KnxInterfaceObjectType.OT_DEVICE => "0",
            KnxInterfaceObjectType.OT_ADDRESS_TABLE => "1",
            KnxInterfaceObjectType.OT_ASSOCIATION_TABLE => "2",
            KnxInterfaceObjectType.OT_APPLICATION_PROGRAM => "3",
            KnxInterfaceObjectType.OT_INTERACE_PROGRAM => "4",
            KnxInterfaceObjectType.OT_EIBOBJECT_ASSOCIATATION_TABLE => "5",
            KnxInterfaceObjectType.OT_ROUTER => "6",
            KnxInterfaceObjectType.OT_LTE_ADDRESS_ROUTING_TABLE => "7",
            KnxInterfaceObjectType.OT_CEMI_SERVER => "8",
            KnxInterfaceObjectType.OT_GROUP_OBJECT_TABLE => "9",
            KnxInterfaceObjectType.OT_POLLING_MASTER => "10",
            KnxInterfaceObjectType.OT_KNXIP_PARAMETER => "11",
            KnxInterfaceObjectType.OT_FILE_SERVER => "13",
            KnxInterfaceObjectType.OT_SECURITY => "17",
            KnxInterfaceObjectType.OT_RF_MEDIUM => "19",
            KnxInterfaceObjectType.OT_INDOOR_BRIGHTNESS_SENSOR => "409",
            KnxInterfaceObjectType.OT_INDOOR_LUMINANCE_SENSOR => "410",
            KnxInterfaceObjectType.OT_LIGHT_SWITCHING_ACTUATOR_BASIC => "417",
            KnxInterfaceObjectType.OT_DIMMING_ACTUATOR_BASIC => "418",
            KnxInterfaceObjectType.OT_DIMMING_SENSOR_BASIC => "420",
            KnxInterfaceObjectType.OT_SWITCHING_SENSOR_BASIC => "421",
            KnxInterfaceObjectType.OT_SUNBLIND_ACTUATOR_BASIC => "800",
            KnxInterfaceObjectType.OT_SUNBLIND_SENSOR_BASIC => "801",
            _ => default,
        };
        public static string GetName(this KnxInterfaceObjectType value) => value switch
        {
            KnxInterfaceObjectType.OT_UNKNOWN => "Unknown Interface Object Type",
            KnxInterfaceObjectType.OT_GENERAL => "General Interface Object Type",
            KnxInterfaceObjectType.OT_DEVICE => "Device Object",
            KnxInterfaceObjectType.OT_ADDRESS_TABLE => "Addresstable Object",
            KnxInterfaceObjectType.OT_ASSOCIATION_TABLE => "Associationtable Object",
            KnxInterfaceObjectType.OT_APPLICATION_PROGRAM => "Applicationprogram Object",
            KnxInterfaceObjectType.OT_INTERACE_PROGRAM => "Interfaceprogram Object",
            KnxInterfaceObjectType.OT_EIBOBJECT_ASSOCIATATION_TABLE => "KNX-Object Associationtable Object",
            KnxInterfaceObjectType.OT_ROUTER => "Router Object",
            KnxInterfaceObjectType.OT_LTE_ADDRESS_ROUTING_TABLE => "LTE Address Routing Table Object",
            KnxInterfaceObjectType.OT_CEMI_SERVER => "cEMI Server Object",
            KnxInterfaceObjectType.OT_GROUP_OBJECT_TABLE => "Group Object Table Object",
            KnxInterfaceObjectType.OT_POLLING_MASTER => "Polling Master",
            KnxInterfaceObjectType.OT_KNXIP_PARAMETER => "KNXnet/IP Parameter Object",
            KnxInterfaceObjectType.OT_FILE_SERVER => "File Server Object",
            KnxInterfaceObjectType.OT_SECURITY => "Security Object",
            KnxInterfaceObjectType.OT_RF_MEDIUM => "RF Medium Object",
            KnxInterfaceObjectType.OT_INDOOR_BRIGHTNESS_SENSOR => "Indoor Brightness Sensor",
            KnxInterfaceObjectType.OT_INDOOR_LUMINANCE_SENSOR => "Indoor Luminance Sensor",
            KnxInterfaceObjectType.OT_LIGHT_SWITCHING_ACTUATOR_BASIC => "Light Switching Actuator Basic",
            KnxInterfaceObjectType.OT_DIMMING_ACTUATOR_BASIC => "Dimming Actuator Basic",
            KnxInterfaceObjectType.OT_DIMMING_SENSOR_BASIC => "Dimming   Sensor Basic",
            KnxInterfaceObjectType.OT_SWITCHING_SENSOR_BASIC => "Switching Sensor Basic",
            KnxInterfaceObjectType.OT_SUNBLIND_ACTUATOR_BASIC => "Sunblind Actuator Basic",
            KnxInterfaceObjectType.OT_SUNBLIND_SENSOR_BASIC => "Sunblind Sensor Basic",
            _ => default,
        };
    }
}
