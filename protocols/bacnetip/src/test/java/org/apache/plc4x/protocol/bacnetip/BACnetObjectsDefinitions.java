/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.protocol.bacnetip;

import java.util.*;
import java.util.stream.Collectors;

import static org.apache.plc4x.protocol.bacnetip.BACnetObjectsDefinitions.ConformanceCode.Type.O;
import static org.apache.plc4x.protocol.bacnetip.BACnetObjectsDefinitions.ConformanceCode.Type.*;

public class BACnetObjectsDefinitions {

    /**
     * All bacnet objects
     */
    static List<BacNetObject> bacNetObjects;


    /**
     * object name (in upper case enum form) to {@link BacNetObject}
     */
    static Map<String, BacNetObject> objectNameToBacNetObjectMap;

    /**
     * a map which property is contained in which objects(name)
     */
    static Map<String, List<String>> propertyToObjectNamesMap;

    /**
     * a map which property-type-combinations is contained in which objects(name)
     */
    static Map<PropertyTypeCombination, List<String>> propertyTypeCombinationToObjectNameMap;

    /**
     * a map which maps the property id to a property name
     */
    static Map<String, String> propertyIdToPropertyNameMap;

    /**
     * a map which maps property to type variants
     */
    static Map<String, Set<String>> propertyToPropertyTypesMaps;

    static Map<PropertyTypeCombination, Integer> propertyTypeCombinationCount;

    static {
        createBacnetObjectsList();
        createObjectNameToBacNetObjectMap();
        createPropertyToObjectNameMap();
        createPropertyTypeCombinationToObjectNameMap();
        createPropertyTypeCombinationCount();
        createPropertyIdToPropertyNameMap();
        createPropertyToPropertyTypesMaps();
    }

    private static void createBacnetObjectsList() {
        bacNetObjects = List.of(
            createAnalogInput(),
            createAnalogOutput(),
            createAnalogValue(),
            createAveraging(),
            createBinarInput(),
            createBinaryOutput(),
            createBinaryValue(),
            createCalendar(),
            createCommand(),
            createDevice(),
            createEventEnrollment(),
            createFile(),
            createGroup(),
            createLifeSafetyPoint(),
            createLifeSafetyZone(),
            createLoop(),
            createMultiStateInput(),
            createMultiStateOutput(),
            createMultiStateValue(),
            createNotificationClass(),
            createProgram(),
            createPulseConverter(),
            createSchedule(),
            createTrendLog(),
            createAccessDoor(),
            createEventLog(),
            createLoadControl(),
            createStructuredView(),
            createTrendLogMultiple(),
            createAccessPoints(),
            createAccessZone(),
            createAccessUser(),
            createAccessRights(),
            createAcessCredential(),
            createCredentialDataInput(),
            createCharacterStringValue(),
            createDateTimeValue(),
            createLargeAnalogValue(),
            createBitStringValue(),
            createOctetStringValue(),
            createTimeValue(),
            createIntegerValue(),
            createPostiveIntegerValue(),
            createDateValue(),
            createDateTimePatternValue(),
            createTimePatternValue(),
            createDatePatternValue(),
            createNetworkSecurity(),
            createGlobalGroup(),
            createNotificationForwarder(),
            createAlertEnrollment(),
            createChannel(),
            createLightingOutput(),
            createBinaryLightingOutput(),
            createNetworkPort(),
            createTimer(),
            createElevatorGroup(),
            createLift(),
            createEscalator(),
            createAccumulator()
        );
    }

    private static BacNetObject createAccumulator() {
        return b("Accumulator",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "Unsigned", c(R, 1)),
            p("Description", "CharacterString", c(O)),
            p("Device_Type", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Scale", "BACnetScale", c(R)),
            p("Units", "BACnetEngineeringUnits", c(R)),
            p("Prescale", "BACnetPrescale", c(O)),
            p("Max_Pres_Value", "Unsigned", c(R)),
            p("Value_Change_Time", "BACnetDateTime", c(O, 2)),
            p("Value_Before_Change", "Unsigned", c(O, 2, 3)),
            p("Value_Set", "Unsigned", c(O, 2, 3)),
            p("Logging_Record", "BACnetAccumulatorRecord", c(O)),
            p("Logging_Object", "BACnetObjectIdentifier", c(O)),
            p("Pulse_Rate", "Unsigned", c(O, 1, 4, 7)),
            p("High_Limit", "Unsigned", c(O, 4, 6)),
            p("Low_Limit", "Unsigned", c(O, 4, 6)),
            p("Limit_Monitoring_Interval", "Unsigned", c(O, 4, 7)),
            p("Notification_Class", "Unsigned", c(O, 4, 6)),
            p("Time_Delay", "Unsigned", c(O, 4, 6)),
            p("Limit_Enable", "BACnetLimitEnable", c(O, 4, 6)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Notify_Type", "BACnetNotifyType", c(O, 4, 6)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 4, 6)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 4, 6)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 6)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 6, 8)),
            p("Time_Delay_Normal", "Unsigned", c(O, 6)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 9)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Fault_High_Limit", "Unsigned", c(O, 10)),
            p("Fault_Low_Limit", "Unsigned", c(O, 10)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createEscalator() {
        return b("Escalator",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Elevator_Group", "BACnetObjectIdentifier", c(R)),
            p("Group_ID", "Unsigned8", c(R)),
            p("Installation_ID", "Unsigned8", c(R)),
            p("Power_Mode", "BOOLEAN", c(O)),
            p("Operation_Direction", "BACnetEscalatorOperationDirection", c(R)),
            p("Escalator_Mode", "BACnetEscalatorMode", c(O)),
            p("Energy_Meter", "REAL", c(O)),
            p("Energy_Meter_Ref", "BACnetDeviceObjectReference", c(O)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Fault_Signals", "BACnetLIST of BACnetEscalatorFault", c(O)),
            p("Passenger_Alarm", "BOOLEAN", c(R)),
            p("Time_Delay", "Unsigned", c(O, 1, 2)),
            p("Time_Delay_Normal", "Unsigned", c(O, 2)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 1, 2)),
            p("Notification_Class", "Unsigned", c(O, 1, 2)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 1, 2)),
            p("Event_State", "BACnetEventState", c(O, 1, 2)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 1, 2)),
            p("Notify_Type", "BACnetNotifyType", c(O, 1, 2)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 1, 2)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 2)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 2)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 2, 3)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 2)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 4)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createLift() {
        return b("Lift",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Elevator_Group", "BACnetObjectIdentifier", c(R)),
            p("Group_ID", "Unsigned8", c(R)),
            p("Installation_ID", "Unsigned8", c(R)),
            p("Floor_Text", "BACnetARRAY[N] of CharacterString", c(O)),
            p("Car_Door_Text", "BACnetARRAY[N] of CharacterString", c(O)),
            p("Assigned_Landing_Calls", "BACnetARRAY[N] of BACnetAssignedLandingCalls", c(O)),
            p("Making_Car_Call", "BACnetARRAY[N] of Unsigned8", c(O)),
            p("Registered_Car_Call", "BACnetARRAY[N] of BACnetLiftCarCallList", c(O)),
            p("Car_Position", "Unsigned8", c(R)),
            p("Car_Moving_Direction", "BACnetLiftCarDirection", c(R)),
            p("Car_Assigned_Direction", "BACnetLiftCarDirection", c(O)),
            p("Car_Door_Status", "BACnetARRAY[N] of BACnetDoorStatus", c(R)),
            p("Car_Door_Command", "BACnetARRAY[N] of BACnetLiftCarDoorCommand", c(O)),
            p("Car_Door_Zone", "BOOLEAN", c(O)),
            p("Car_Mode", "BACnetLiftCarMode", c(O)),
            p("Car_Load", "REAL", c(O)),
            p("Car_Load_Units", "BACnetEngineeringUnits", c(O, 1)),
            p("Next_Stopping_Floor", "Unsigned8", c(O)),
            p("Passenger_Alarm", "BOOLEAN", c(R)),
            p("Time_Delay", "Unsigned", c(O, 2, 3)),
            p("Time_Delay_Normal", "Unsigned", c(O, 3)),
            p("Energy_Meter", "REAL", c(O)),
            p("Energy_Meter_Ref", "BACnetDeviceObjectReference", c(O)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Car_Drive_Status", "BACnetLiftCarDriveStatus", c(O)),
            p("Fault_Signals", "BACnetLIST of BACnetLiftFault", c(R)),
            p("Landing_Door_Status", "BACnetARRAY[N] of BACnetLandingDoorStatus", c(O)),
            p("Higher_Deck", "BACnetObjectIdentifier", c(O)),
            p("Lower_Deck", "BACnetObjectIdentifier", c(O)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 2, 3)),
            p("Notification_Class", "Unsigned", c(O, 2, 3)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 2, 3)),
            p("Event_State", "BACnetEventState", c(O, 2, 3)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 2, 3)),
            p("Notify_Type", "BACnetNotifyType", c(O, 2, 3)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 2, 3)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 3)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 3)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 3)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 3, 4)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 5)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createElevatorGroup() {
        return b("Elevator Group",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Machine_Room_ID", "BACnetObjectIdentifier", c(R)),
            p("Group_ID", "Unsigned8", c(R)),
            p("Group_Members", "BACnetARRAY[N] of BACnetObjectIdentifier", c(R)),
            p("Group_Mode", "BACnetLiftGroupMode", c(O, 1)),
            p("Landing_Calls", "BACnetLIST of BACnetLandingCallStatus", c(O, 1)),
            p("Landing_Call_Control", "BACnetLandingCallStatus", c(O, 1)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createTimer() {
        return b("Timer",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "Unsigned", c(R, 1)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(O, 3)),
            p("Reliability", "BACnetReliability", c(O, 1)),
            p("Out_Of_Service", "BOOLEAN", c(O)),
            p("Timer_State", "BACnetTimerState", c(R, 2)),
            p("Timer_Running", "BOOLEAN", c(R, 2)),
            p("Update_Time", "BACnetDateTime", c(O, 3)),
            p("Last_State_Change", "BACnetTimerTransition", c(O)),
            p("Expiration_Time", "BACnetDateTime", c(O)),
            p("Initial_Timeout", "Unsigned", c(O)),
            p("Default_Timeout", "Unsigned", c(O)),
            p("Min_Pres_Value", "Unsigned", c(O, 4)),
            p("Max_Pres_Value", "Unsigned", c(O, 4)),
            p("Resolution", "Unsigned", c(O)),
            p("State_Change_Values", "BACnetARRAY[7] of BACnetTimerStateChangeValue", c(O, 5)),
            p("List_Of_Object_Property_References", "BACnetLIST of BACnetDeviceObjectPropertyReference", c(O, 5)),
            p("Priority_For_Writing", "Unsigned(1..16)", c(O, 5)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 3, 6)),
            p("Notification_Class", "Unsigned", c(O, 3, 6)),
            p("Time_Delay", "Unsigned", c(O, 3, 6)),
            p("Time_Delay_Normal", "Unsigned", c(O, 3, 6)),
            p("Alarm_Values", "BACnetLIST of BACnetTimerState", c(O, 3, 6)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 3, 6)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 3, 6)),
            p("Notify_Type", "BACnetNotifyType", c(O, 3, 6)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 3, 6)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 6)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 6, 7)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 8)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createNetworkPort() {
        return b("Network Port",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Reliability", "BACnetReliability", c(R)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Network_Type", "BACnetNetworkType", c(R)),
            p("Protocol_Level", "BACnetProtocolLevel", c(R)),
            p("Reference_Port", "Unsigned", c(O)),
            p("Network_Number", "Unsigned16", c(R, 1)),
            p("Network_Number_Quality", "BACnetNetworkNumberQuality", c(R)),
            p("Changes_Pending", "BOOLEAN", c(R)),
            p("Command", "BACnetNetworkPortCommand", c(O, 2)),
            p("MAC_Address", "OCTET STRING", c(O, 3)),
            p("APDU_Length", "Unsigned", c(R)),
            p("Link_Speed", "REAL", c(R)),
            p("Link_Speeds", "BACnetARRAY[N] of REAL", c(O, 4)),
            p("Link_Speed_Autonegotiate", "BOOLEAN", c(O)),
            p("Network_Interface_Name", "CharacterString", c(O)),
            p("BACnet_IP_Mode", "BACnetIPMode", c(O, 5)),
            p("IP_Address", "OCTET STRING", c(O, 6)),
            p("BACnet_IP_UDP_Port", "Unsigned16", c(O, 5)),
            p("IP_Subnet_Mask", "OCTET STRING", c(O, 6)),
            p("IP_Default_Gateway", "OCTET STRING", c(O, 6)),
            p("BACnet_IP_Multicast_Address", "OCTET STRING", c(O, 7)),
            p("IP_DNS_Server", "BACnetARRAY[N] of OCTET STRING", c(O, 6)),
            p("IP_DHCP_Enable", "BOOLEAN", c(O, 8)),
            p("IP_DHCP_Lease_Time", "Unsigned", c(O)),
            p("IP_DHCP_Lease_Time_Remaining", "Unsigned", c(O)),
            p("IP_DHCP_Server", "OCTET STRING", c(O)),
            p("BACnet_IP_NAT_Traversal", "BOOLEAN", c(O, 9)),
            p("BACnet_IP_Global_Address", "BACnetHostNPort", c(O, 10)),
            p("BBMD_Broadcast_Distribution_Table", "BACnetLIST of BACnetBDTEntry", c(O, 11)),
            p("BBMD_Accept_FD_Registrations", "BOOLEAN", c(O, 11)),
            p("BBMD_Foreign_Device_Table", "BACnetLIST of BACnetFDTEntry", c(O, 12)),
            p("FD_BBMD_Address", "BACnetHostNPort", c(O, 13)),
            p("FD_Subscription_Lifetime", "Unsigned16", c(O, 13)),
            p("BACnet_IPv6_Mode", "BACnetIPMode", c(O, 14)),
            p("IPv6_Address", "OCTET STRING", c(O, 15)),
            p("IPv6_Prefix_Length", "Unsigned8", c(O, 15)),
            p("BACnet_IPv6_UDP_Port", "Unsigned16", c(O, 14)),
            p("IPv6_Default_Gateway", "OCTET STRING", c(O, 15)),
            p("BACnet_IPv6_Multicast_Address", "OCTET STRING", c(O, 14)),
            p("IPv6_DNS_Server", "BACnetARRAY[N] of OCTET STRING", c(O, 15)),
            p("IPv6_Auto_Addressing_Enable", "BOOLEAN", c(O, 16)),
            p("IPv6_DHCP_Lease_Time", "Unsigned", c(O)),
            p("IPv6_DHCP_Lease_Time_Remaining", "Unsigned", c(O)),
            p("IPv6_DHCP_Server", "OCTET STRING", c(O)),
            p("IPv6_Zone_Index", "CharacterString", c(O, 17)),
            p("Max_Master", "Unsigned8(0..127)", c(O, 18)),
            p("Max_Info_Frames", "Unsigned8", c(O, 18)),
            p("Slave_Proxy_Enable", "BOOLEAN", c(O, 19)),
            p("Manual_Slave_Address_Binding", "BACnetLIST of BACnetAddressBinding", c(O, 19)),
            p("Auto_Slave_Discovery", "BOOLEAN", c(O, 20)),
            p("Slave_Address_Binding", "BACnetLIST of BACnetAddressBinding", c(O, 21)),
            p("Virtual_MAC_Address_Table", "BACnetLIST of BACnetVMACEntry", c(O, 22)),
            p("Routing_Table", "BACnetLIST of BACnetRouterEntry", c(O)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 23, 24)),
            p("Notification_Class", "Unsigned", c(O, 23, 24)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 23, 24)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 23, 24)),
            p("Notify_Type", "BACnetNotifyType", c(O, 23, 24)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 23, 24)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 24)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 24)),
            p("Event_State", "BACnetEventState", c(O, 23)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createBinaryLightingOutput() {
        return b("Binary Lighting Output",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "BACnetBinaryLightingPV", c(W)),
            p("Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Blink_Warn_Enable", "BOOLEAN", c(R)),
            p("Egress_Time", "Unsigned", c(R)),
            p("Egress_Active", "BOOLEAN", c(R)),
            p("Feedback_Value", "BACnetBinaryLightingPV", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(R)),
            p("Relinquish_Default", "BACnetBinaryLightingPV", c(R)),
            p("Power", "REAL", c(O)),
            p("Polarity", "BACnetPolarity", c(O)),
            p("Elapsed_Active_Time", "Unsigned32", c(O, 1)),
            p("Time_Of_Active_Time_Reset", "BACnetDateTime", c(O, 1)),
            p("Strike_Count", "Unsigned", c(O, 2)),
            p("Time_Of_Strike_Count_Reset", "BACnetDateTime", c(O, 2)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 3, 4)),
            p("Notification_Class", "Unsigned", c(O, 3, 4)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 3, 4)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 3, 4)),
            p("Notify_Type", "BACnetNotifyType", c(O, 3, 4)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 3, 4)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 4)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 4)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 5)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(R)),
            p("Value_Source", "BACnetValueSource", c(O, 6, 8, 10)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 7, 9)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 7, 9)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 9)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createLightingOutput() {
        return b("Lighting Output",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "REAL", c(W)),
            p("Tracking_Value", "REAL", c(R)),
            p("Lighting_Command", "BACnetLightingCommand", c(W)),
            p("In_Progress", "BACnetLightingInProgress", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Blink_Warn_Enable", "BOOLEAN", c(R)),
            p("Egress_Time", "Unsigned", c(R)),
            p("Egress_Active", "BOOLEAN", c(R)),
            p("Default_Fade_Time", "Unsigned", c(R)),
            p("Default_Ramp_Rate", "REAL", c(R)),
            p("Default_Step_Increment", "REAL", c(R)),
            p("Transition", "BACnetLightingTransition", c(O)),
            p("Feedback_Value", "REAL", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(R)),
            p("Relinquish_Default", "REAL", c(R)),
            p("Power", "REAL", c(O)),
            p("Instantaneous_Power", "REAL", c(O)),
            p("Min_Actual_Value", "REAL", c(O, 1)),
            p("Max_Actual_Value", "REAL", c(O, 1)),
            p("Lighting_Command_Default_Priority", "Unsigned", c(R)),
            p("COV_Increment", "REAL", c(O, 2)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 3)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(R)),
            p("Value_Source", "BACnetValueSource", c(O, 4, 6, 8)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 5, 7)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 5, 7)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 7)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createChannel() {
        return b("Channel",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "BACnetChannelValue", c(W)),
            p("Last_Priority", "Unsigned", c(R)),
            p("Write_Status", "BACnetWriteStatus", c(R)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("List_Of_Object_Property_References", "BACnetARRAY[N] of BACnetDeviceObjectPropertyReference", c(W, 1)),
            p("Execution_Delay", "BACnetARRAY[N] of Unsigned", c(O, 1)),
            p("Allow_Group_Delay_Inhibit", "BOOLEAN", c(O)),
            p("Channel_Number", "Unsigned16", c(W)),
            p("Control_Groups", "BACnetARRAY[N] of Unsigned32", c(W)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 2, 3)),
            p("Notification_Class", "Unsigned", c(O, 2, 3)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 2, 3)),
            p("Event_State", "BACnetEventState", c(O, 2, 3)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 2, 3)),
            p("Notify_Type", "BACnetNotifyType", c(O, 2, 3)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 2, 3)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 3)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 3)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 4)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Value_Source", "BACnetValueSource", c(O, 5, 6, 7)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createAlertEnrollment() {
        return b("Alert Enrollment",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "BACnetObjectIdentifier", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Event_Detection_Enable", "BOOLEAN", c(R)),
            p("Notification_Class", "Unsigned", c(R)),
            p("Event_Enable", "BACnetEventTransitionBits", c(R)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(R)),
            p("Notify_Type", "BACnetNotifyType", c(R)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(R)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createNotificationForwarder() {
        return b("Notification Forwarder",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Reliability", "BACnetReliability", c(R)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Recipient_List", "BACnetLIST of BACnetDestination", c(R)),
            p("Subscribed_Recipients", "BACnetLIST of BACnetEventNotificationSubscription", c(W)),
            p("Process_Identifier_Filter", "BACnetProcessIdSelection", c(R)),
            p("Port_Filter", "BACnetARRAY[N] of BACnetPortPermission", c(O, 1)),
            p("Local_Forwarding_Only", "BOOLEAN", c(R)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createGlobalGroup() {
        return b("Global Group",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Group_Members", "BACnetARRAY[N] of BACnetDeviceObjectPropertyReference", c(R)),
            p("Group_Member_Names", "BACnetARRAY[N] of CharacterString", c(O)),
            p("Present_Value", "BACnetARRAY[N] of BACnetPropertyAccessResult", c(R)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Member_Status_Flags", "BACnetStatusFlags", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_of_Service", "BOOLEAN", c(R)),
            p("Update_Interval", "Unsigned", c(O)),
            p("Requested_Update_Interval", "Unsigned", c(O)),
            p("COV_Resubscription_Interval", "Unsigned", c(O)),
            p("Client_COV_Increment", "BACnetClientCOV", c(O)),
            p("Time_Delay", "Unsigned", c(O, 1, 4)),
            p("Notification_Class", "Unsigned", c(O, 1, 4)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 1, 4)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 1, 4)),
            p("Notify_Type", "BACnetNotifyType", c(O, 1, 4)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 1, 4)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 4)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 4)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 1, 4)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 4)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 4, 5)),
            p("Time_Delay_Normal", "Unsigned", c(O, 4)),
            p("COVU_Period", "Unsigned", c(O, 2)),
            p("COVU_Recipients", "BACnetLIST of BACnetRecipient", c(O, 2)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 6)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createNetworkSecurity() {
        return b("Network Security",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Base_Device_Security_Policy", "BACnetSecurityLevel", c(W)),
            p("Network_Access_Security_Policies", "BACnetARRAY[N] of BACnetNetworkSecurityPolicy", c(W)),
            p("Security_Time_Window", "Unsigned", c(W)),
            p("Packet_Reorder_Time", "Unsigned", c(W)),
            p("Distribution_Key_Revision", "Unsigned8", c(R)),
            p("Key_Sets", "BACnetARRAY[2] of BACnetSecurityKeySet", c(R)),
            p("Last_Key_Server", "BACnetAddressBinding", c(W)),
            p("Security_PDU_Timeout", "Unsigned16", c(W)),
            p("Update_Key_Set_Timeout", "Unsigned16", c(R)),
            p("Supported_Security_Algorithms", "BACnetLIST of Unsigned8", c(R)),
            p("Do_Not_Hide", "BOOLEAN", c(W)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createDatePatternValue() {
        return b("Date Pattern Value",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "Date", c(R, 1)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(O)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(O, 2)),
            p("Relinquish_Default", "Date", c(O, 2)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 3)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 4, 5)),
            p("Notification_Class", "Unsigned", c(O, 4, 5)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 4, 5)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 4, 5)),
            p("Notify_Type", "BACnetNotifyType", c(O, 4, 5)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 4, 5)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(O, 2)),
            p("Value_Source", "BACnetValueSource", c(O, 6, 8, 10)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 7, 9)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 7, 9)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 9)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createTimePatternValue() {
        return b("Time Pattern Value",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "Time", c(R, 1)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(O)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(O, 2)),
            p("Relinquish_Default", "Time", c(O, 2)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 3)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 4, 5)),
            p("Notification_Class", "Unsigned", c(O, 4, 5)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 4, 5)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 4, 5)),
            p("Notify_Type", "BACnetNotifyType", c(O, 4, 5)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 4, 5)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(O, 2)),
            p("Value_Source", "BACnetValueSource", c(O, 6, 8, 10)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 7, 9)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 7, 9)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 9)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createDateTimePatternValue() {
        return b("DateTime Pattern Value",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "BACnetDateTime", c(R, 1)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(O)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(O)),
            p("Is_UTC", "BOOLEAN", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(O, 2)),
            p("Relinquish_Default", "BACnetDateTime", c(O, 2)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 3)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 4, 5)),
            p("Notification_Class", "Unsigned", c(O, 4, 5)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 4, 5)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 4, 5)),
            p("Notify_Type", "BACnetNotifyType", c(O, 4, 5)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 4, 5)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(O, 2)),
            p("Value_Source", "BACnetValueSource", c(O, 6, 8, 10)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 7, 9)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 7, 9)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 9)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createDateValue() {
        return b("Date Value",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "Date", c(R, 1)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(O)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(O, 2)),
            p("Relinquish_Default", "Date", c(O, 2)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 3)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 4, 5)),
            p("Notification_Class", "Unsigned", c(O, 4, 5)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 4, 5)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 4, 5)),
            p("Notify_Type", "BACnetNotifyType", c(O, 4, 5)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 4, 5)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(O, 2)),
            p("Value_Source", "BACnetValueSource", c(O, 6, 8, 10)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 7, 9)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 7, 9)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 9)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createPostiveIntegerValue() {
        return b("Positive Integer Value",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "Unsigned", c(R, 1)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(O, 4)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(O)),
            p("Units", "BACnetEngineeringUnits", c(R)),
            p("Priority_Array", "BACnetPriorityArray", c(O, 2)),
            p("Relinquish_Default", "Unsigned", c(O, 2)),
            p("COV_Increment", "Unsigned", c(O, 3)),
            p("Time_Delay", "Unsigned", c(O, 4, 6)),
            p("Notification_Class", "Unsigned", c(O, 4, 6)),
            p("High_Limit", "Unsigned", c(O, 4, 6)),
            p("Low_Limit", "Unsigned", c(O, 4, 6)),
            p("Deadband", "Unsigned", c(O, 4, 6)),
            p("Limit_Enable", "BACnetLimitEnable", c(O, 4, 6)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Notify_Type", "BACnetNotifyType", c(O, 4, 6)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 4, 6)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 4, 6)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 6)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 6, 7)),
            p("Time_Delay_Normal", "Unsigned", c(O, 6)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 8)),
            p("Min_Pres_Value", "Unsigned", c(O)),
            p("Max_Pres_Value", "Unsigned", c(O)),
            p("Resolution", "Unsigned", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Fault_High_Limit", "Unsigned", c(O, 9)),
            p("Fault_Low_Limit", "Unsigned", c(O, 9)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(O, 2)),
            p("Value_Source", "BACnetValueSource", c(O, 10, 12, 14)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 11, 13)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 11, 13)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 13)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createIntegerValue() {
        return b("Integer Value",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "INTEGER", c(R, 1)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(O, 4)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(O)),
            p("Units", "BACnetEngineeringUnits", c(R)),
            p("Priority_Array", "BACnetPriorityArray", c(O, 2)),
            p("Relinquish_Default", "INTEGER", c(O, 2)),
            p("COV_Increment", "Unsigned", c(O, 3)),
            p("Time_Delay", "Unsigned", c(O, 4, 6)),
            p("Notification_Class", "Unsigned", c(O, 4, 6)),
            p("High_Limit", "INTEGER", c(O, 4, 6)),
            p("Low_Limit", "INTEGER", c(O, 4, 6)),
            p("Deadband", "Unsigned", c(O, 4, 6)),
            p("Limit_Enable", "BACnetLimitEnable", c(O, 4, 6)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Notify_Type", "BACnetNotifyType", c(O, 4, 6)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 4, 6)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 4, 6)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 6)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 6, 7)),
            p("Time_Delay_Normal", "Unsigned", c(O, 6)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 8)),
            p("Min_Pres_Value", "INTEGER", c(O)),
            p("Max_Pres_Value", "INTEGER", c(O)),
            p("Resolution", "INTEGER", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Fault_High_Limit", "INTEGER", c(O, 9)),
            p("Fault_Low_Limit", "INTEGER", c(O, 9)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(O, 2)),
            p("Value_Source", "BACnetValueSource", c(O, 10, 12, 14)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 11, 13)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 11, 13)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 13)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createTimeValue() {
        return b("Time Value",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "Time", c(R, 1)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(O)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(O, 2)),
            p("Relinquish_Default", "Time", c(O, 2)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 3)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 4, 5)),
            p("Notification_Class", "Unsigned", c(O, 4, 5)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 4, 5)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 4, 5)),
            p("Notify_Type", "BACnetNotifyType", c(O, 4, 5)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 4, 5)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(O, 2)),
            p("Value_Source", "BACnetValueSource", c(O, 6, 8, 10)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 7, 9)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 7, 9)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 9)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createOctetStringValue() {
        return b("OctetString Value",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "OCTET STRING", c(R, 1)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(O)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(O, 2)),
            p("Relinquish_Default", "OCTET STRING", c(O, 2)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 3)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(O, 2)),
            p("Value_Source", "BACnetValueSource", c(O, 4, 6, 8)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 5, 7)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 5, 7)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 7)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createBitStringValue() {
        return b("BitString Value",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "BIT STRING", c(R, 1)),
            p("Bit_Text", "BACnetARRAY[N] of CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(O, 3)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(O, 2)),
            p("Relinquish_Default", "BIT STRING", c(O, 2)),
            p("Time_Delay", "Unsigned", c(O, 3, 5)),
            p("Notification_Class", "Unsigned", c(O, 3, 5)),
            p("Alarm_Values", "BACnetARRAY[N] of BIT STRING", c(O, 3, 5)),
            p("Bit_Mask", "BIT STRING", c(O, 3, 5)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 3, 5)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 3, 5)),
            p("Notify_Type", "BACnetNotifyType", c(O, 3, 5)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 3, 5)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 3, 5)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 5)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 5, 6)),
            p("Time_Delay_Normal", "Unsigned", c(O, 5)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 7)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(O, 2)),
            p("Value_Source", "BACnetValueSource", c(O)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createLargeAnalogValue() {
        return b("Large Analog Value",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "Double", c(R, 1)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(O, 4)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(O)),
            p("Units", "BACnetEngineeringUnits", c(R)),
            p("Priority_Array", "BACnetPriorityArray", c(O, 2)),
            p("Relinquish_Default", "Double", c(O, 2)),
            p("COV_Increment", "Double", c(O, 3)),
            p("Time_Delay", "Unsigned", c(O, 4, 6)),
            p("Notification_Class", "Unsigned", c(O, 4, 6)),
            p("High_Limit", "Double", c(O, 4, 6)),
            p("Low_Limit", "Double", c(O, 4, 6)),
            p("Deadband", "Double", c(O, 4, 6)),
            p("Limit_Enable", "BACnetLimitEnable", c(O, 4, 6)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Notify_Type", "BACnetNotifyType", c(O, 4, 6)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 4, 6)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 4, 6)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 6)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 6, 7)),
            p("Time_Delay_Normal", "Unsigned", c(O, 6)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 8)),
            p("Min_Pres_Value", "Double", c(O)),
            p("Max_Pres_Value", "Double", c(O)),
            p("Resolution", "Double", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Fault_High_Limit", "Double", c(O, 9)),
            p("Fault_Low_Limit", "Double", c(O, 9)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(O, 2)),
            p("Value_Source", "BACnetValueSource", c(O, 10, 12, 14)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 11, 13)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 11, 13)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 13)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createDateTimeValue() {
        return b("DateTime Value",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "BACnetDateTime", c(R, 1)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(O)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(O, 2)),
            p("Relinquish_Default", "BACnetDateTime", c(O, 2)),
            p("Is_UTC", "BOOLEAN", c(O)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 3)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 4, 5)),
            p("Notification_Class", "Unsigned", c(O, 4, 5)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 4, 5)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 4, 5)),
            p("Notify_Type", "BACnetNotifyType", c(O, 4, 5)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 4, 5)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(O, 2)),
            p("Value_Source", "BACnetValueSource", c(O, 6, 8, 10)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 7, 9)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 7, 9)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 9)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createCharacterStringValue() {
        return b("CharacterString Value",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "CharacterString", c(R, 1)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(O, 3)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(O, 2)),
            p("Relinquish_Default", "CharacterString", c(O, 2)),
            p("Time_Delay", "Unsigned", c(O, 3, 5)),
            p("Notification_Class", "Unsigned", c(O, 3, 5)),
            p("Alarm_Values", "BACnetARRAY[N] of BACnetOptionalCharacterString", c(O, 3, 5)),
            p("Fault_Values", "BACnetARRAY[N] of BACnetOptionalCharacterString", c(O, 7)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 3, 5)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 3, 5)),
            p("Notify_Type", "BACnetNotifyType", c(O, 3, 5)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 3, 5)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 3, 5)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 5)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 5, 6)),
            p("Time_Delay_Normal", "Unsigned", c(O, 5)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 7)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(O, 2)),
            p("Value_Source", "BACnetValueSource", c(O, 8, 10, 12)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 9, 11)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 9, 11)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 11)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createCredentialDataInput() {
        return b("Credential Data Input",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "BACnetAuthenticationFactor", c(R, 1)),
            p("Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Reliability", "BACnetReliability", c(R, 1)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Supported_Formats", "BACnetARRAY[N] of BACnetAuthenticationFactorFormat", c(R)),
            p("Supported_Format_Classes", "BACnetARRAY[N] of Unsigned", c(O, 2)),
            p("Update_Time", "BACnetTimeStamp", c(R)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 3, 4)),
            p("Notification_Class", "Unsigned", c(O, 3, 4)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 3, 4)),
            p("Event_State", "BACnetEventState", c(O, 3, 4)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 3, 4)),
            p("Notify_Type", "BACnetNotifyType", c(O, 3, 4)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 3, 4)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 4)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 4)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createAcessCredential() {
        return b("Access Credential",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Global_Identifier", "Unsigned32", c(W)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Reliability", "BACnetReliability", c(R)),
            p("Credential_Status", "BACnetBinaryPV", c(R)),
            p("Reason_For_Disable", "BACnetLIST of BACnetAccessCredentialDisableReason", c(R)),
            p("Authentication_Factors", "BACnetARRAY[N] of BACnetCredentialAuthenticationFactor", c(R)),
            p("Activation_Time", "BACnetDateTime", c(R)),
            p("Expiration_Time", "BACnetDateTime", c(R)),
            p("Credential_Disable", "BACnetAccessCredentialDisable", c(R)),
            p("Days_Remaining", "INTEGER", c(O, 1)),
            p("Uses_Remaining", "INTEGER", c(O)),
            p("Absentee_Limit", "Unsigned", c(O, 1)),
            p("Belongs_To", "BACnetDeviceObjectReference", c(O)),
            p("Assigned_Access_Rights", "BACnetARRAY[N] of BACnetAssignedAccessRights", c(R)),
            p("Last_Access_Point", "BACnetDeviceObjectReference", c(O)),
            p("Last_Access_Event", "BACnetAccessEvent", c(O)),
            p("Last_Use_Time", "BACnetDateTime", c(O)),
            p("Trace_Flag", "BOOLEAN", c(O)),
            p("Threat_Authority", "BACnetAccessThreatLevel", c(O)),
            p("Extended_Time_Enable", "BOOLEAN", c(O)),
            p("Authorization_Exemptions", "BACnetLIST of BACnetAuthorizationExemption", c(O)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createAccessRights() {
        return b("Access Rights",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Global_Identifier", "Unsigned32", c(W)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Reliability", "BACnetReliability", c(R)),
            p("Enable", "BOOLEAN", c(R)),
            p("Negative_Access_Rules", "BACnetARRAY[N] of BACnetAccessRule", c(R)),
            p("Positive_Access_Rules", "BACnetARRAY[N] of BACnetAccessRule", c(R)),
            p("Accompaniment", "BACnetDeviceObjectReference", c(O)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createAccessUser() {
        return b("Access User",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Global_Identifier", "Unsigned32", c(W)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Reliability", "BACnetReliability", c(R)),
            p("User_Type", "BACnetAccessUserType", c(R)),
            p("User_Name", "CharacterString", c(O)),
            p("User_External_Identifier", "CharacterString", c(O)),
            p("User_Information_Reference", "CharacterString", c(O)),
            p("Members", "BACnetLIST of BACnetDeviceObjectReference", c(O)),
            p("Member_Of", "BACnetLIST of BACnetDeviceObjectReference", c(O)),
            p("Credentials", "BACnetLIST of BACnetDeviceObjectReference", c(R)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createAccessZone() {
        return b("Access Zone",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Global_Identifier", "Unsigned32", c(W)),
            p("Occupancy_State", "BACnetAccessZoneOccupancyState", c(R)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(R, 1)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Occupancy_Count", "Unsigned", c(O, 1, 3, 4)),
            p("Occupancy_Count_Enable", "BOOLEAN", c(O, 3, 4)),
            p("Adjust_Value", "INTEGER", c(O, 3, 4, 5)),
            p("Occupancy_Upper_Limit", "Unsigned", c(O)),
            p("Occupancy_Lower_Limit", "Unsigned", c(O)),
            p("Credentials_In_Zone", "BACnetLIST of BACnetDeviceObjectReference", c(O)),
            p("Last_Credential_Added", "BACnetDeviceObjectReference", c(O)),
            p("Last_Credential_Added_Time", "BACnetDateTime", c(O)),
            p("Last_Credential_Removed", "BACnetDeviceObjectReference", c(O)),
            p("Last_Credential_Removed_Time", "BACnetDateTime", c(O)),
            p("Passback_Mode", "BACnetAccessPassbackMode", c(O)),
            p("Passback_Timeout", "Unsigned", c(O, 2)),
            p("Entry_Points", "BACnetLIST of BACnetDeviceObjectReference", c(R)),
            p("Exit_Points", "BACnetLIST of BACnetDeviceObjectReference", c(R)),
            p("Time_Delay", "Unsigned", c(O, 3, 7)),
            p("Notification_Class", "Unsigned", c(O, 3, 7)),
            p("Alarm_Values", "BACnetLIST of BACnetAccessZoneOccupancyState", c(O, 3, 7)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 3, 7)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 3, 7)),
            p("Notify_Type", "BACnetNotifyType", c(O, 3, 7)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 3, 7)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 7)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 7)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 3, 7)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 7)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 7, 8)),
            p("Time_Delay_Normal", "Unsigned", c(O, 7)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createAccessPoints() {
        return b("Access Point",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(R)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Authentication_Status", "BACnetAuthenticationStatus", c(R)),
            p("Active_Authentication_Policy", "Unsigned", c(R)),
            p("Number_Of_Authentication_Policies", "Unsigned", c(R)),
            p("Authentication_Policy_List", "BACnetARRAY[N] of BACnetAuthenticationPolicy", c(O, 1)),
            p("Authentication_Policy_Names", "BACnetARRAY[N] of CharacterString", c(O, 1)),
            p("Authorization_Mode", "BACnetAuthorizationMode", c(R)),
            p("Verification_Time", "Unsigned", c(O)),
            p("Lockout", "BOOLEAN", c(O, 2)),
            p("Lockout_Relinquish_Time", "Unsigned", c(O)),
            p("Failed_Attempts", "Unsigned", c(O)),
            p("Failed_Attempt_Events", "BACnetLIST of BACnetAccessEvent", c(O)),
            p("Max_Failed_Attempts", "Unsigned", c(O, 3)),
            p("Failed_Attempts_Time", "Unsigned", c(O, 3)),
            p("Threat_Level", "BACnetAccessThreatLevel", c(O)),
            p("Occupancy_Upper_Limit_Enforced", "BOOLEAN", c(O)),
            p("Occupancy_Lower_Limit_Enforced", "BOOLEAN", c(O)),
            p("Occupancy_Count_Adjust", "BOOLEAN", c(O)),
            p("Accompaniment_Time", "Unsigned", c(O)),
            p("Access_Event", "BACnetAccessEvent", c(R)),
            p("Access_Event_Tag", "Unsigned", c(R)),
            p("Access_Event_Time", "BACnetTimeStamp", c(R)),
            p("Access_Event_Credential", "BACnetDeviceObjectReference", c(R)),
            p("Access_Event_Authentication_Factor", "BACnetAuthenticationFactor", c(O)),
            p("Access_Doors", "BACnetARRAY[N] of BACnetDeviceObjectReference", c(R)),
            p("Priority_For_Writing", "Unsigned(1..16)", c(R)),
            p("Muster_Point", "BOOLEAN", c(O)),
            p("Zone_To", "BACnetDeviceObjectReference", c(O)),
            p("Zone_From", "BACnetDeviceObjectReference", c(O)),
            p("Notification_Class", "Unsigned", c(O, 4, 6)),
            p("Transaction_Notification_Class", "Unsigned", c(O)),
            p("Access_Alarm_Events", "BACnetLIST of BACnetAccessEvent", c(O, 4, 6)),
            p("Access_Transaction_Events", "BACnetLIST of BACnetAccessEvent", c(O, 4, 6)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Notify_Type", "BACnetNotifyType", c(O, 4, 6)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 4, 6)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 4, 6)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 6)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 6, 7)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createTrendLogMultiple() {
        return b("Trend Log Multiple",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Enable", "BOOLEAN", c(W)),
            p("Start_Time", "BACnetDateTime", c(O, 1)),
            p("Stop_Time", "BACnetDateTime", c(O, 1)),
            p("Log_DeviceObjectProperty", "BACnetARRAY[N] of BACnetDeviceObjectPropertyReference", c(R)),
            p("Logging_Type", "BACnetLoggingType", c(R)),
            p("Log_Interval", "Unsigned", c(R, 2)),
            p("Align_Intervals", "BOOLEAN", c(O, 3)),
            p("Interval_Offset", "Unsigned", c(O, 3)),
            p("Trigger", "BOOLEAN", c(O)),
            p("Stop_When_Full", "BOOLEAN", c(R)),
            p("Buffer_Size", "Unsigned32", c(R)),
            p("Log_Buffer", "BACnetLIST of BACnetLogMultipleRecord", c(R)),
            p("Record_Count", "Unsigned32", c(W)),
            p("Total_Record_Count", "Unsigned32", c(R)),
            p("Notification_Threshold", "Unsigned32", c(O, 4, 6)),
            p("Records_Since_Notification", "Unsigned32", c(O, 4, 6)),
            p("Last_Notify_Record", "Unsigned32", c(O, 4, 6)),
            p("Notification_Class", "Unsigned", c(O, 4, 6)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Notify_Type", "BACnetNotifyType", c(O, 4, 6)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 4, 6)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 4, 6)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 6)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 6, 7)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 8)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createStructuredView() {
        return b("Structured View",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Node_Type", "BACnetNodeType", c(R)),
            p("Node_Subtype", "CharacterString", c(O)),
            p("Subordinate_List", "BACnetARRAY[N] of BACnetDeviceObjectReference", c(R)),
            p("Subordinate_Annotations", "BACnetARRAY[N] of CharacterString", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Subordinate_Tags", "BACnetARRAY[N] of BACnetNameValueCollection", c(O)),
            p("Subordinate_Node_Types", "BACnetARRAY[N] of BACnetNodeType", c(O)),
            p("Subordinate_Relationships", "BACnetARRAY[N] of BACnetRelationship", c(O)),
            p("Default_Subordinate_Relationship", "BACnetRelationship", c(O)),
            p("Represents", "BACnetDeviceObjectReference", c(O)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createLoadControl() {
        return b("Load Control",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "BACnetShedState", c(R)),
            p("State_Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Requested_Shed_Level", "BACnetShedLevel", c(W)),
            p("Start_Time", "BACnetDateTime", c(W)),
            p("Shed_Duration", "Unsigned", c(W)),
            p("Duty_Window", "Unsigned", c(W)),
            p("Enable", "BOOLEAN", c(W)),
            p("Full_Duty_Baseline", "REAL", c(O)),
            p("Expected_Shed_Level", "BACnetShedLevel", c(R)),
            p("Actual_Shed_Level", "BACnetShedLevel", c(R)),
            p("Shed_Levels", "BACnetARRAY[N] of Unsigned", c(W, 1)),
            p("Shed_Level_Descriptions", "BACnetARRAY[N] of CharacterString", c(R)),
            p("Notification_Class", "Unsigned", c(O, 2, 4)),
            p("Time_Delay", "Unsigned", c(O, 2, 4)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 2, 4)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 2, 4)),
            p("Notify_Type", "BACnetNotifyType", c(O, 2, 4)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 2, 4)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 4)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 4)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 2, 4)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 4)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 4, 5)),
            p("Time_Delay_Normal", "Unsigned", c(O, 4)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 6)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Value_Source", "BACnetValueSource", c(O, 7, 8, 9)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createEventLog() {
        return b("Event Log",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Enable", "BOOLEAN", c(W)),
            p("Start_Time", "BACnetDateTime", c(O, 1, 2)),
            p("Stop_Time", "BACnetDateTime", c(O, 1, 2)),
            p("Stop_When_Full", "BOOLEAN", c(R)),
            p("Buffer_Size", "Unsigned32", c(R)),
            p("Log_Buffer", "BACnetLIST of BACnetEventLogRecord", c(R)),
            p("Record_Count", "Unsigned32", c(W)),
            p("Total_Record_Count", "Unsigned32", c(R)),
            p("Notification_Threshold", "Unsigned32", c(O, 3, 5)),
            p("Records_Since_Notification", "Unsigned32", c(O, 3, 5)),
            p("Last_Notify_Record", "Unsigned32", c(O, 3, 5)),
            p("Notification_Class", "Unsigned", c(O, 3, 5)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 3, 5)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 3, 5)),
            p("Notify_Type", "BACnetNotifyType", c(O, 3, 5)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 3, 5)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 3, 5)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 5)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 5, 6)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 7)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createAccessDoor() {
        return b("Access Door",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "BACnetDoorValue", c(W)),
            p("Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(R)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Priority_Array", "BACnetPriorityArray", c(R)),
            p("Relinquish_Default", "BACnetDoorValue", c(R)),
            p("Door_Status", "BACnetDoorStatus", c(O, 1, 2)),
            p("Lock_Status", "BACnetLockStatus", c(O, 1)),
            p("Secured_Status", "BACnetDoorSecuredStatus", c(O)),
            p("Door_Members", "BACnetARRAY[N] of BACnetDeviceObjectReference", c(O)),
            p("Door_Pulse_Time", "Unsigned", c(R)),
            p("Door_Extended_Pulse_Time", "Unsigned", c(R)),
            p("Door_Unlock_Delay_Time", "Unsigned", c(O)),
            p("Door_Open_Too_Long_Time", "Unsigned", c(R)),
            p("Door_Alarm_State", "BACnetDoorAlarmState", c(O, 1, 3)),
            p("Masked_Alarm_Values", "BACnetLIST of BACnetDoorAlarmState", c(O)),
            p("Maintenance_Required", "BACnetMaintenance", c(O)),
            p("Time_Delay", "Unsigned", c(O, 3, 5)),
            p("Notification_Class", "Unsigned", c(O, 3, 5)),
            p("Alarm_Values", "BACnetLIST of BACnetDoorAlarmState", c(O, 3, 5)),
            p("Fault_Values", "BACnetLIST of BACnetDoorAlarmState", c(O)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 3, 5)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 3, 5)),
            p("Notify_Type", "BACnetNotifyType", c(O, 3, 5)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 3, 5)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 3, 5)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 5)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 5, 6)),
            p("Time_Delay_Normal", "Unsigned", c(O, 5)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(R)),
            p("Value_Source", "BACnetValueSource", c(O, 7, 9, 11)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 8, 10)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 8, 10)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 10)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createTrendLog() {
        return b("Trend Log",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Enable", "BOOLEAN", c(W)),
            p("Start_Time", "BACnetDateTime", c(O, 1, 2)),
            p("Stop_Time", "BACnetDateTime", c(O, 1, 2)),
            p("Log_DeviceObjectProperty", "BACnetDeviceObjectPropertyReference", c(O, 8)),
            p("Log_Interval", "Unsigned", c(O, 1, 3)),
            p("COV_Resubscription_Interval", "Unsigned", c(O)),
            p("Client_COV_Increment", "BACnetClientCOV", c(O)),
            p("Stop_When_Full", "BOOLEAN", c(R)),
            p("Buffer_Size", "Unsigned32", c(R)),
            p("Log_Buffer", "BACnetLIST of BACnetLogRecord", c(R)),
            p("Record_Count", "Unsigned32", c(W)),
            p("Total_Record_Count", "Unsigned32", c(R)),
            p("Logging_Type", "BACnetLoggingType", c(R)),
            p("Align_Intervals", "BOOLEAN", c(O, 5)),
            p("Interval_Offset", "Unsigned", c(O, 5)),
            p("Trigger", "BOOLEAN", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Notification_Threshold", "Unsigned32", c(O, 4, 7)),
            p("Records_Since_Notification", "Unsigned32", c(O, 4, 7)),
            p("Last_Notify_Record", "Unsigned32", c(O, 4, 7)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Notification_Class", "Unsigned", c(O, 4, 7)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 4, 7)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 4, 7)),
            p("Notify_Type", "BACnetNotifyType", c(O, 4, 7)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 4, 7)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 7)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 7)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 4, 7)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 7)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 7, 9)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 10)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createSchedule() {
        return b("Schedule",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "Any", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Effective_Period", "BACnetDateRange", c(R)),
            p("Weekly_Schedule", "BACnetARRAY[7] of BACnetDailySchedule", c(O, 1)),
            p("Exception_Schedule", "BACnetARRAY[N] of BACnetSpecialEvent", c(O, 1)),
            p("Schedule_Default", "Any", c(R)),
            p("List_Of_Object_Property_References", "BACnetLIST of BACnetDeviceObjectPropertyReference", c(R)),
            p("Priority_For_Writing", "Unsigned(1..16)", c(R)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Reliability", "BACnetReliability", c(R)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 2, 3)),
            p("Notification_Class", "Unsigned", c(O, 2, 3)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 2, 3)),
            p("Event_State", "BACnetEventState", c(O, 2, 3)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 2, 3)),
            p("Notify_Type", "BACnetNotifyType", c(O, 2, 3)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 2, 3)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 3)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 3)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createPulseConverter() {
        return b("Pulse Converter",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "REAL", c(R, 1)),
            p("Input_Reference", "BACnetObjectPropertyReference", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Units", "BACnetEngineeringUnits", c(R)),
            p("Scale_Factor", "REAL", c(R)),
            p("Adjust_Value", "REAL", c(W)),
            p("Count", "Unsigned", c(R)),
            p("Update_Time", "BACnetDateTime", c(R)),
            p("Count_Change_Time", "BACnetDateTime", c(R)),
            p("Count_Before_Change", "Unsigned", c(R)),
            p("COV_Increment", "REAL", c(O, 2)),
            p("COV_Period", "Unsigned", c(O, 2)),
            p("Notification_Class", "Unsigned", c(O, 3, 5)),
            p("Time_Delay", "Unsigned", c(O, 3, 5)),
            p("High_Limit", "REAL", c(O, 3, 5)),
            p("Low_Limit", "REAL", c(O, 3, 5)),
            p("Deadband", "REAL", c(O, 3, 5)),
            p("Limit_Enable", "BACnetLimitEnable", c(O, 3, 5)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 3, 5)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 3, 5)),
            p("Notify_Type", "BACnetNotifyType", c(O, 3, 5)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 3, 5)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 3, 5)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 5)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 5, 6)),
            p("Time_Delay_Normal", "Unsigned", c(O, 5)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 7)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createProgram() {
        return b("Program",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Program_State", "BACnetProgramState", c(R)),
            p("Program_Change", "BACnetProgramRequest", c(W)),
            p("Reason_For_Halt", "BACnetProgramError", c(O, 1)),
            p("Description_Of_Halt", "CharacterString", c(O, 1)),
            p("Program_Location", "CharacterString", c(O)),
            p("Description", "CharacterString", c(O)),
            p("Instance_Of", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 2, 3)),
            p("Notification_Class", "Unsigned", c(O, 2, 3)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 2, 3)),
            p("Event_State", "BACnetEventState", c(O, 2, 3)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 2, 3)),
            p("Notify_Type", "BACnetNotifyType", c(O, 2, 3)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 2, 3)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 3)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 3)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 4)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createNotificationClass() {
        return b("Notification Class",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Notification_Class", "Unsigned", c(R)),
            p("Priority", "BACnetARRAY[3] of Unsigned", c(R)),
            p("Ack_Required", "BACnetEventTransitionBits", c(R)),
            p("Recipient_List", "BACnetLIST of BACnetDestination", c(R)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Status_Flags", "BACnetStatusFlags", c(O, 1)),
            p("Event_State", "BACnetEventState", c(O, 1)),
            p("Reliability", "BACnetReliability", c(O, 1)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 1, 2)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 1, 2)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 1, 2)),
            p("Notify_Type", "BACnetNotifyType", c(O, 1, 2)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 1, 2)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 2)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 2)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 3)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createMultiStateValue() {
        return b("Multi-state Value",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "Unsigned", c(R, 1)),
            p("Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Number_Of_States", "Unsigned", c(R)),
            p("State_Text", "BACnetARRAY[N] of CharacterString", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(O, 3)),
            p("Relinquish_Default", "Unsigned", c(O, 3)),
            p("Time_Delay", "Unsigned", c(O, 4, 6)),
            p("Notification_Class", "Unsigned", c(O, 4, 6)),
            p("Alarm_Values", "BACnetLIST of Unsigned", c(O, 4, 6)),
            p("Fault_Values", "BACnetLIST of Unsigned", c(O, 8)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Notify_Type", "BACnetNotifyType", c(O, 4, 6)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 4, 6)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 4, 6)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 6)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 6, 7)),
            p("Time_Delay_Normal", "Unsigned", c(O, 6)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 8)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(O, 3)),
            p("Value_Source", "BACnetValueSource", c(O, 9, 11, 13)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 10, 12)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 10, 12)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 12)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createMultiStateOutput() {
        return b("Multi-state Output",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "Unsigned", c(W)),
            p("Description", "CharacterString", c(O)),
            p("Device_Type", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Number_Of_States", "Unsigned", c(R)),
            p("State_Text", "BACnetARRAY[N] of CharacterString", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(R)),
            p("Relinquish_Default", "Unsigned", c(R)),
            p("Time_Delay", "Unsigned", c(O, 1, 3)),
            p("Notification_Class", "Unsigned", c(O, 1, 3)),
            p("Feedback_Value", "Unsigned", c(O, 1)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 1, 3)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 1, 3)),
            p("Notify_Type", "BACnetNotifyType", c(O, 1, 3)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 1, 3)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 3)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 3)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 1, 3)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 3)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 3, 4)),
            p("Time_Delay_Normal", "Unsigned", c(O, 3)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 5)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Interface_Value", "BACnetOptionalUnsigned", c(O)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(R)),
            p("Value_Source", "BACnetValueSource", c(O, 6, 8, 10)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 7, 9)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 7, 9)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 9)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createMultiStateInput() {
        return b("Multi-state Input",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "Unsigned", c(R, 1)),
            p("Description", "CharacterString", c(O)),
            p("Device_Type", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Number_Of_States", "Unsigned", c(R)),
            p("State_Text", "BACnetARRAY[N] of CharacterString", c(O)),
            p("Time_Delay", "Unsigned", c(O, 3, 5)),
            p("Notification_Class", "Unsigned", c(O, 3, 5)),
            p("Alarm_Values", "BACnetLIST of Unsigned", c(O, 3, 5)),
            p("Fault_Values", "BACnetLIST of Unsigned", c(O, 7)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 3, 5)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 3, 5)),
            p("Notify_Type", "BACnetNotifyType", c(O, 3, 5)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 3, 5)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 3, 5)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 5)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 5, 6)),
            p("Time_Delay_Normal", "Unsigned", c(O, 5)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 7)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Interface_Value", "BACnetOptionalUnsigned", c(O)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createLoop() {
        return b("Loop",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "REAL", c(R, 7)),
            p("Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(O, 7)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Update_Interval", "Unsigned", c(O)),
            p("Output_Units", "BACnetEngineeringUnits", c(R)),
            p("Manipulated_Variable_Reference", "BACnetObjectPropertyReference", c(R)),
            p("Controlled_Variable_Reference", "BACnetObjectPropertyReference", c(R)),
            p("Controlled_Variable_Value", "REAL", c(R)),
            p("Controlled_Variable_Units", "BACnetEngineeringUnits", c(R)),
            p("Setpoint_Reference", "BACnetSetpointReference", c(R)),
            p("Setpoint", "REAL", c(R)),
            p("Action", "BACnetAction", c(R)),
            p("Proportional_Constant", "REAL", c(O, 1)),
            p("Proportional_Constant_Units", "BACnetEngineeringUnits", c(O, 1)),
            p("Integral_Constant", "REAL", c(O, 2)),
            p("Integral_Constant_Units", "BACnetEngineeringUnits", c(O, 2)),
            p("Derivative_Constant", "REAL", c(O, 3)),
            p("Derivative_Constant_Units", "BACnetEngineeringUnits", c(O, 3)),
            p("Bias", "REAL", c(O)),
            p("Maximum_Output", "REAL", c(O)),
            p("Minimum_Output", "REAL", c(O)),
            p("Priority_For_Writing", "Unsigned(1..16)", c(R)),
            p("COV_Increment", "REAL", c(O, 4)),
            p("Time_Delay", "Unsigned", c(O, 5, 8)),
            p("Notification_Class", "Unsigned", c(O, 5, 8)),
            p("Error_Limit", "REAL", c(O, 5, 8)),
            p("Deadband", "REAL", c(O, 5, 8)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 5, 8)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 5, 8)),
            p("Notify_Type", "BACnetNotifyType", c(O, 5, 8)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 5, 8)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 8)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 8)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 5, 8)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 8)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 8, 9)),
            p("Time_Delay_Normal", "Unsigned", c(O, 8)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 10)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Low_Diff_Limit", "BACnetOptionalREAL", c(O)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createLifeSafetyZone() {
        return b("Life Safety Zone",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "BACnetLifeSafetyState", c(R)),
            p("Tracking_Value", "BACnetLifeSafetyState", c(R, 1)),
            p("Description", "CharacterString", c(O)),
            p("Device_Type", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(R, 1)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Mode", "BACnetLifeSafetyMode", c(W)),
            p("Accepted_Modes", "BACnetLIST of BACnetLifeSafetyMode", c(R)),
            p("Time_Delay", "Unsigned", c(O, 2, 4)),
            p("Notification_Class", "Unsigned", c(O, 2, 4)),
            p("Life_Safety_Alarm_Values", "BACnetLIST of BACnetLifeSafetyState", c(O, 2, 4)),
            p("Alarm_Values", "BACnetLIST of BACnetLifeSafetyState", c(O, 2, 4)),
            p("Fault_Values", "BACnetLIST of BACnetLifeSafetyState", c(O)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 2, 4)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 2, 4)),
            p("Notify_Type", "BACnetNotifyType", c(O, 2, 4)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 2, 4)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 4)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 4)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 2, 4)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 4)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 4, 5)),
            p("Time_Delay_Normal", "Unsigned", c(O, 4)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O)),
            p("Silenced", "BACnetSilencedState", c(R)),
            p("Operation_Expected", "BACnetLifeSafetyOperation", c(R)),
            p("Maintenance_Required", "BOOLEAN", c(O)),
            p("Zone_Members", "BACnetLIST of BACnetDeviceObjectReference", c(R)),
            p("Member_Of", "BACnetLIST of BACnetDeviceObjectReference", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Value_Source", "BACnetValueSource", c(O, 6, 7, 8)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createLifeSafetyPoint() {
        return b("Life Safety Point",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "BACnetLifeSafetyState", c(R)),
            p("Tracking_Value", "BACnetLifeSafetyState", c(R, 1)),
            p("Description", "CharacterString", c(O)),
            p("Device_Type", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(R, 1)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Mode", "BACnetLifeSafetyMode", c(W)),
            p("Accepted_Modes", "BACnetLIST of BACnetLifeSafetyMode", c(R)),
            p("Time_Delay", "Unsigned", c(O, 2, 5)),
            p("Notification_Class", "Unsigned", c(O, 2, 5)),
            p("Life_Safety_Alarm_Values", "BACnetLIST of BACnetLifeSafetyState", c(O, 2, 5)),
            p("Alarm_Values", "BACnetLIST of BACnetLifeSafetyState", c(O, 2, 5)),
            p("Fault_Values", "BACnetLIST of BACnetLifeSafetyState", c(O)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 2, 5)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 2, 5)),
            p("Notify_Type", "BACnetNotifyType", c(O, 2, 5)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 2, 5)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 2, 5)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 5)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 5, 6)),
            p("Time_Delay_Normal", "Unsigned", c(O, 5)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O)),
            p("Silenced", "BACnetSilencedState", c(R)),
            p("Operation_Expected", "BACnetLifeSafetyOperation", c(R)),
            p("Maintenance_Required", "BACnetMaintenance", c(O)),
            p("Setting", "Unsigned8", c(O)),
            p("Direct_Reading", "REAL", c(O, 3)),
            p("Units", "BACnetEngineeringUnits", c(O, 3)),
            p("Member_Of", "BACnetLIST of BACnetDeviceObjectReference", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Value_Source", "BACnetValueSource", c(O, 7, 8, 9)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createGroup() {
        return b("Group",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("List_Of_Group_Members", "BACnetLIST of ReadAccessSpecification", c(R)),
            p("Present_Value", "BACnetLIST of ReadAccessResult", c(R)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createFile() {
        return b("File",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("File_Type", "CharacterString", c(R)),
            p("File_Size", "Unsigned", c(R, 1)),
            p("Modification_Date", "BACnetDateTime", c(R)),
            p("Archive", "BOOLEAN", c(W)),
            p("Read_Only", "BOOLEAN", c(R)),
            p("File_Access_Method", "BACnetFileAccessMethod", c(R)),
            p("Record_Count", "Unsigned", c(O, 2)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createEventEnrollment() {
        return b("Event Enrollment",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Event_Type", "BACnetEventType", c(R)),
            p("Notify_Type", "BACnetNotifyType", c(R)),
            p("Event_Parameters", "BACnetEventParameter", c(R)),
            p("Object_Property_Reference", "BACnetDeviceObjectPropertyReference", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Event_Enable", "BACnetEventTransitionBits", c(R)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(R)),
            p("Notification_Class", "Unsigned", c(R)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(R)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O)),
            p("Event_Detection_Enable", "BOOLEAN", c(R)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 2)),
            p("Time_Delay_Normal", "Unsigned", c(R)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Reliability", "BACnetReliability", c(R)),
            p("Fault_Type", "BACnetFaultType", c(O, 3)),
            p("Fault_Parameters", "BACnetFaultParameter", c(O, 3)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(R)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createDevice() {
        return b("Device",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("System_Status", "BACnetDeviceStatus", c(R)),
            p("Vendor_Name", "CharacterString", c(R)),
            p("Vendor_Identifier", "Unsigned16", c(R)),
            p("Model_Name", "CharacterString", c(R)),
            p("Firmware_Revision", "CharacterString", c(R)),
            p("Application_Software_Version", "CharacterString", c(R)),
            p("Location", "CharacterString", c(O)),
            p("Description", "CharacterString", c(O)),
            p("Protocol_Version", "Unsigned", c(R)),
            p("Protocol_Revision", "Unsigned", c(R)),
            p("Protocol_Services_Supported", "BACnetServicesSupported", c(R)),
            p("Protocol_Object_Types_Supported", "BACnetObjectTypesSupported", c(R)),
            p("Object_List", "BACnetARRAY[N] of BACnetObjectIdentifier", c(R)),
            p("Structured_Object_List", "BACnetARRAY[N] of BACnetObjectIdentifier", c(O)),
            p("Max_APDU_Length_Accepted", "Unsigned", c(R)),
            p("Segmentation_Supported", "BACnetSegmentation", c(R)),
            p("Max_Segments_Accepted", "Unsigned", c(O, 1)),
            p("VT_Classes_Supported", "BACnetLIST of BACnetVTClass", c(O, 2)),
            p("Active_VT_Sessions", "BACnetLIST of BACnetVTSession", c(O, 2)),
            p("Local_Time", "Time", c(O, 3, 4, 15)),
            p("Local_Date", "Date", c(O, 3, 4, 15)),
            p("UTC_Offset", "INTEGER", c(O, 4)),
            p("Daylight_Savings_Status", "BOOLEAN", c(O, 4)),
            p("APDU_Segment_Timeout", "Unsigned", c(O, 1)),
            p("APDU_Timeout", "Unsigned", c(R)),
            p("Number_Of_APDU_Retries", "Unsigned", c(R)),
            p("Time_Synchronization_Recipients", "BACnetLIST of BACnetRecipient", c(O, 5)),
            p("Max_Master", "Unsigned(0..127)", c(O, 6)),
            p("Max_Info_Frames", "Unsigned", c(O, 6)),
            p("Device_Address_Binding", "BACnetLIST of BACnetAddressBinding", c(R)),
            p("Database_Revision", "Unsigned", c(R)),
            p("Configuration_Files", "BACnetARRAY[N] of BACnetObjectIdentifier", c(O, 7)),
            p("Last_Restore_Time", "BACnetTimeStamp", c(O, 7)),
            p("Backup_Failure_Timeout", "Unsigned16", c(O, 8)),
            p("Backup_Preparation_Time", "Unsigned16", c(O, 16)),
            p("Restore_Preparation_Time", "Unsigned16", c(O, 16)),
            p("Restore_Completion_Time", "Unsigned16", c(O, 16)),
            p("Backup_And_Restore_State", "BACnetBackupState", c(O, 7)),
            p("Active_COV_Subscriptions", "BACnetLIST of BACnetCOVSubscription", c(O, 9)),
            p("Last_Restart_Reason", "BACnetRestartReason", c(O, 13)),
            p("Time_Of_Device_Restart", "BACnetTimeStamp", c(O, 13)),
            p("Restart_Notification_Recipients", "BACnetLIST of BACnetRecipient", c(O, 17)),
            p("UTC_Time_Synchronization_Recipients", "BACnetLIST of BACnetRecipient", c(O, 5)),
            p("Time_Synchronization_Interval", "Unsigned", c(O, 5)),
            p("Align_Intervals", "BOOLEAN", c(O, 5)),
            p("Interval_Offset", "Unsigned", c(O, 5)),
            p("Serial_Number", "CharacterString", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Status_Flags", "BACnetStatusFlags", c(O, 18)),
            p("Event_State", "BACnetEventState", c(O, 18)),
            p("Reliability", "BACnetReliability", c(O, 18)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 18, 19)),
            p("Notification_Class", "Unsigned", c(O, 18, 19)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 18, 19)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 18, 19)),
            p("Notify_Type", "BACnetNotifyType", c(O, 18, 19)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 18, 19)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 19)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 19)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 20)),
            p("Active_COV_Multiple_Subscriptions", "BACnetLIST of BACnetCOVMultipleSubscription", c(O, 21)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Deployed_Profile_Location", "CharacterString", c(O, 22)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createCommand() {
        return b("Command",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "Unsigned", c(W)),
            p("In_Process", "BOOLEAN", c(R)),
            p("All_Writes_Successful", "BOOLEAN", c(R)),
            p("Action", "BACnetARRAY[N] of BACnetActionList", c(R)),
            p("Action_Text", "BACnetARRAY[N] of CharacterString", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Status_Flags", "BACnetStatusFlags", c(O, 1)),
            p("Event_State", "BACnetEventState", c(O, 1)),
            p("Reliability", "BACnetReliability", c(O, 1)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 1, 2)),
            p("Notification_Class", "Unsigned", c(O, 1, 2)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 1, 2)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 1, 2)),
            p("Notify_Type", "BACnetNotifyType", c(O, 1, 2)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 1, 2)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 2)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 2)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 3)),
            p("Value_Source", "BACnetValueSource", c(O, 4, 5, 6)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createCalendar() {
        return b("Calendar",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Description", "CharacterString", c(O)),
            p("Present_Value", "BOOLEAN", c(R)),
            p("Date_List", "BACnetLIST of BACnetCalendarEntry", c(R)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createBinaryValue() {
        return b("Binary Value",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "BACnetBinaryPV", c(R, 1)),
            p("Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Inactive_Text", "CharacterString", c(O, 2)),
            p("Active_Text", "CharacterString", c(O, 2)),
            p("Change_Of_State_Time", "BACnetDateTime", c(O, 3)),
            p("Change_Of_State_Count", "Unsigned", c(O, 3)),
            p("Time_Of_State_Count_Reset", "BACnetDateTime", c(O, 3)),
            p("Elapsed_Active_Time", "Unsigned32", c(O, 4)),
            p("Time_Of_Active_Time_Reset", "BACnetDateTime", c(O, 4)),
            p("Minimum_Off_Time", "Unsigned32", c(O)),
            p("Minimum_On_Time", "Unsigned32", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(O, 5)),
            p("Relinquish_Default", "BACnetBinaryPV", c(O, 5)),
            p("Time_Delay", "Unsigned", c(O, 6, 8)),
            p("Notification_Class", "Unsigned", c(O, 6, 8)),
            p("Alarm_Value", "BACnetBinaryPV", c(O, 6, 8)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 6, 8)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 6, 8)),
            p("Notify_Type", "BACnetNotifyType", c(O, 6, 8)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 6, 8)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 8)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 8)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 6, 8)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 8)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 8, 9)),
            p("Time_Delay_Normal", "Unsigned", c(O, 8)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 10)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(O, 5)),
            p("Value_Source", "BACnetValueSource", c(O, 11, 13, 15)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 12, 14)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 12, 14)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 14)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createBinaryOutput() {
        return b("Binary Output",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "BACnetBinaryPV", c(W)),
            p("Description", "CharacterString", c(O)),
            p("Device_Type", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Polarity", "BACnetPolarity", c(R)),
            p("Inactive_Text", "CharacterString", c(O, 1)),
            p("Active_Text", "CharacterString", c(O, 1)),
            p("Change_Of_State_Time", "BACnetDateTime", c(O, 2)),
            p("Change_Of_State_Count", "Unsigned", c(O, 2)),
            p("Time_Of_State_Count_Reset", "BACnetDateTime", c(O, 2)),
            p("Elapsed_Active_Time", "Unsigned32", c(O, 3)),
            p("Time_Of_Active_Time_Reset", "BACnetDateTime", c(O, 3)),
            p("Minimum_Off_Time", "Unsigned32", c(O)),
            p("Minimum_On_Time", "Unsigned32", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(R)),
            p("Relinquish_Default", "BACnetBinaryPV", c(R)),
            p("Time_Delay", "Unsigned", c(O, 4, 6)),
            p("Notification_Class", "Unsigned", c(O, 4, 6)),
            p("Feedback_Value", "BACnetBinaryPV", c(O, 4)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 4, 6)),
            p("Notify_Type", "BACnetNotifyType", c(O, 4, 6)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 4, 6)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 4, 6)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 6)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 6, 7)),
            p("Time_Delay_Normal", "Unsigned", c(O, 6)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 8)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Interface_Value", "BACnetOptionalBinaryPV", c(O)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(R)),
            p("Value_Source", "BACnetValueSource", c(O, 9, 11, 13)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 10, 12)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 10, 12)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 12)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createBinarInput() {
        return b("Binary Input",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "BACnetBinaryPV", c(R, 1)),
            p("Description", "CharacterString", c(O)),
            p("Device_Type", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Polarity", "BACnetPolarity", c(R)),
            p("Inactive_Text", "CharacterString", c(O, 2)),
            p("Active_Text", "CharacterString", c(O, 2)),
            p("Change_Of_State_Time", "BACnetDateTime", c(O, 3)),
            p("Change_Of_State_Count", "Unsigned", c(O, 3)),
            p("Time_Of_State_Count_Reset", "BACnetDateTime", c(O, 3)),
            p("Elapsed_Active_Time", "Unsigned32", c(O, 4)),
            p("Time_Of_Active_Time_Reset", "BACnetDateTime", c(O, 4)),
            p("Time_Delay", "Unsigned", c(O, 5, 7)),
            p("Notification_Class", "Unsigned", c(O, 5, 7)),
            p("Alarm_Value", "BACnetBinaryPV", c(O, 5, 7)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 5, 7)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 5, 7)),
            p("Notify_Type", "BACnetNotifyType", c(O, 5, 7)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 5, 7)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 7)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 7)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 5, 7)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 7)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 7, 8)),
            p("Time_Delay_Normal", "Unsigned", c(O, 7)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 9)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Interface_Value", "BACnetOptionalBinaryPV", c(O)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createAveraging() {
        return b("Averaging",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Minimum_Value", "REAL", c(R)),
            p("Minimum_Value_Timestamp", "BACnetDateTime", c(O)),
            p("Average_Value", "REAL", c(R)),
            p("Variance_Value", "REAL", c(O)),
            p("Maximum_Value", "REAL", c(R)),
            p("Maximum_Value_Timestamp", "BACnetDateTime", c(O)),
            p("Description", "CharacterString", c(O)),
            p("Attempted_Samples", "Unsigned", c(W, 1)),
            p("Valid_Samples", "Unsigned", c(R)),
            p("Object_Property_Reference", "BACnetDeviceObjectPropertyReference", c(R, 1)),
            p("Window_Interval", "Unsigned", c(W, 1)),
            p("Window_Samples", "Unsigned", c(W, 1)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createAnalogValue() {
        return b("Analog Value",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "REAL", c(R, 4)),
            p("Description", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Units", "BACnetEngineeringUnits", c(R)),
            p("Priority_Array", "BACnetPriorityArray", c(O, 1)),
            p("Relinquish_Default", "REAL", c(O, 1)),
            p("COV_Increment", "REAL", c(O, 2)),
            p("Time_Delay", "Unsigned", c(O, 3, 6)),
            p("Notification_Class", "Unsigned", c(O, 3, 6)),
            p("High_Limit", "REAL", c(O, 3, 6)),
            p("Low_Limit", "REAL", c(O, 3, 6)),
            p("Deadband", "REAL", c(O, 3, 6)),
            p("Limit_Enable", "BACnetLimitEnable", c(O, 3, 6)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 3, 6)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 3, 6)),
            p("Notify_Type", "BACnetNotifyType", c(O, 3, 6)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 3, 6)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 6)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 3, 6)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 6)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 6, 7)),
            p("Time_Delay_Normal", "Unsigned", c(O, 6)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 8)),
            p("Min_Pres_Value", "REAL", c(O)),
            p("Max_Pres_Value", "REAL", c(O)),
            p("Resolution", "REAL", c(O)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Fault_High_Limit", "REAL", c(O, 9)),
            p("Fault_Low_Limit", "REAL", c(O, 9)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(O, 1)),
            p("Value_Source", "BACnetValueSource", c(O, 10, 12, 14)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 11, 13)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 11, 13)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 13)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createAnalogOutput() {
        return b("Analog Output",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "REAL", c(W)),
            p("Description", "CharacterString", c(O)),
            p("Device_Type", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Units", "BACnetEngineeringUnits", c(R)),
            p("Min_Pres_Value", "REAL", c(O)),
            p("Max_Pres_Value", "REAL", c(O)),
            p("Resolution", "REAL", c(O)),
            p("Priority_Array", "BACnetPriorityArray", c(R)),
            p("Relinquish_Default", "REAL", c(R)),
            p("COV_Increment", "REAL", c(O, 1)),
            p("Time_Delay", "Unsigned", c(O, 2, 4)),
            p("Notification_Class", "Unsigned", c(O, 2, 4)),
            p("High_Limit", "REAL", c(O, 2, 4)),
            p("Low_Limit", "REAL", c(O, 2, 4)),
            p("Deadband", "REAL", c(O, 2, 4)),
            p("Limit_Enable", "BACnetLimitEnable", c(O, 2, 4)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 2, 4)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 2, 4)),
            p("Notify_Type", "BACnetNotifyType", c(O, 2, 4)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 2, 4)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 4)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 4)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 2, 4)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 4)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 4, 5)),
            p("Time_Delay_Normal", "Unsigned", c(O, 4)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 6)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Interface_Value", "BACnetOptionalREAL", c(O)),
            p("Current_Command_Priority", "BACnetOptionalUnsigned", c(R)),
            p("Value_Source", "BACnetValueSource", c(O, 7, 9, 11)),
            p("Value_Source_Array", "BACnetARRAY[16] of BACnetValueSource", c(O, 8, 10)),
            p("Last_Command_Time", "BACnetTimeStamp", c(O, 8, 10)),
            p("Command_Time_Array", "BACnetARRAY[16] of BACnetTimeStamp", c(O, 10)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static BacNetObject createAnalogInput() {
        return b("Analog Input",
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Identifier", "BACnetObjectIdentifier", c(R)),
            p("Object_Name", "CharacterString", c(R)),
            p("Object_Type", "BACnetObjectType", c(R)),
            p("Present_Value", "REAL", c(R, 1)),
            p("Description", "CharacterString", c(O)),
            p("Device_Type", "CharacterString", c(O)),
            p("Status_Flags", "BACnetStatusFlags", c(R)),
            p("Event_State", "BACnetEventState", c(R)),
            p("Reliability", "BACnetReliability", c(O)),
            p("Out_Of_Service", "BOOLEAN", c(R)),
            p("Update_Interval", "Unsigned", c(O)),
            p("Units", "BACnetEngineeringUnits", c(R)),
            p("Min_Pres_Value", "REAL", c(O)),
            p("Max_Pres_Value", "REAL", c(O)),
            p("Resolution", "REAL", c(O)),
            p("COV_Increment", "REAL", c(O, 2)),
            p("Time_Delay", "Unsigned", c(O, 3, 5)),
            p("Notification_Class", "Unsigned", c(O, 3, 5)),
            p("High_Limit", "REAL", c(O, 3, 5)),
            p("Low_Limit", "REAL", c(O, 3, 5)),
            p("Deadband", "REAL", c(O, 3, 5)),
            p("Limit_Enable", "BACnetLimitEnable", c(O, 3, 5)),
            p("Event_Enable", "BACnetEventTransitionBits", c(O, 3, 5)),
            p("Acked_Transitions", "BACnetEventTransitionBits", c(O, 3, 5)),
            p("Notify_Type", "BACnetNotifyType", c(O, 3, 5)),
            p("Event_Time_Stamps", "BACnetARRAY[3] of BACnetTimeStamp", c(O, 3, 5)),
            p("Event_Message_Texts", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Message_Texts_Config", "BACnetARRAY[3] of CharacterString", c(O, 5)),
            p("Event_Detection_Enable", "BOOLEAN", c(O, 3, 5)),
            p("Event_Algorithm_Inhibit_Ref", "BACnetObjectPropertyReference", c(O, 5)),
            p("Event_Algorithm_Inhibit", "BOOLEAN", c(O, 5, 6)),
            p("Time_Delay_Normal", "Unsigned", c(O, 5)),
            p("Reliability_Evaluation_Inhibit", "BOOLEAN", c(O, 7)),
            p("Property_List", "BACnetARRAY[N] of BACnetPropertyIdentifier", c(R)),
            p("Interface_Value", "BACnetOptionalREAL", c(O)),
            p("Fault_High_Limit", "REAL", c(O, 8)),
            p("Fault_Low_Limit", "REAL", c(O, 8)),
            p("Tags", "BACnetARRAY[N] of BACnetNameValue", c(O)),
            p("Profile_Location", "CharacterString", c(O)),
            p("Profile_Name", "CharacterString", c(O))
        );
    }

    private static void createObjectNameToBacNetObjectMap() {
        objectNameToBacNetObjectMap = bacNetObjects.stream().collect(Collectors.toMap(BACnetObjectsDefinitions::mapObjectNameToEnumName, bacNetObject -> bacNetObject));
    }

    static String mapObjectNameToEnumName(BacNetObject bacNetObject) {
        return mapObjectNameToEnumName(bacNetObject.name);
    }

    static String mapObjectNameToEnumName(String name) {
        String upperCase = name.toUpperCase();
        String minusReplaced = upperCase.replaceAll("-", "_");
        String spacesReplaced = minusReplaced.replaceAll(" ", "_");
        String mappedName = spacesReplaced;
        switch (spacesReplaced) {
            case "DATE_PATTERN_VALUE":
                mappedName = "DATEPATTERN_VALUE";
                break;
            case "DATETIME_PATTERN_VALUE":
                mappedName = "DATETIMEPATTERN_VALUE";
                break;
            case "TIME_PATTERN_VALUE":
                mappedName = "TIMEPATTERN_VALUE";
                break;
        }
        return mappedName;
    }

    private static void createPropertyToObjectNameMap() {
        propertyToObjectNamesMap = new HashMap<>();
        bacNetObjects.forEach(bacNetObject -> {
            String bacNetObjectName = bacNetObject.name;
            bacNetObject.properties.forEach(bacNetProperty -> {
                String propertyIdentifier = bacNetProperty.propertyIdentifier;
                propertyToObjectNamesMap.putIfAbsent(propertyIdentifier, new LinkedList<>());
                propertyToObjectNamesMap.get(propertyIdentifier).add(bacNetObjectName);
            });
        });
    }

    private static void createPropertyTypeCombinationToObjectNameMap() {
        propertyTypeCombinationToObjectNameMap = new HashMap<>();
        bacNetObjects.forEach(bacNetObject -> {
            String bacNetObjectName = bacNetObject.name;
            bacNetObject.properties.forEach(bacNetProperty -> {
                PropertyTypeCombination propertyTypeCombination = new PropertyTypeCombination(bacNetProperty);
                propertyTypeCombinationToObjectNameMap.putIfAbsent(propertyTypeCombination, new LinkedList<>());
                propertyTypeCombinationToObjectNameMap.get(propertyTypeCombination).add(bacNetObjectName);
            });
        });
    }

    private static void createPropertyTypeCombinationCount() {
        propertyTypeCombinationCount = new HashMap<>();
        bacNetObjects.forEach(bacNetObject -> {
            bacNetObject.properties.forEach(bacNetProperty -> {
                PropertyTypeCombination propertyTypeCombination = new PropertyTypeCombination(bacNetProperty);
                propertyTypeCombinationCount.putIfAbsent(propertyTypeCombination, 0);
                int count = propertyTypeCombinationCount.get(propertyTypeCombination);
                propertyTypeCombinationCount.put(propertyTypeCombination, ++count);
            });
        });
    }

    static void createPropertyIdToPropertyNameMap() {
        propertyIdToPropertyNameMap = new HashMap<>();
        propertyToObjectNamesMap.keySet().forEach(propertyName ->
            propertyIdToPropertyNameMap.put(mapPropertyNameToEnumName(propertyName), propertyName)
        );
    }

    static void createPropertyToPropertyTypesMaps() {
        propertyToPropertyTypesMaps = new HashMap<>();
        propertyTypeCombinationToObjectNameMap.keySet().forEach(propertyTypeCombination -> {
            propertyToPropertyTypesMaps.putIfAbsent(propertyTypeCombination.propertyIdentifier, new HashSet<>());
            propertyToPropertyTypesMaps.get(propertyTypeCombination.propertyIdentifier).add(propertyTypeCombination.propertyDataType);
        });
    }

    static String mapPropertyNameToEnumName(String propertyName) {
        String upperCase = propertyName.toUpperCase();
        String minusReplaced = upperCase.replaceAll("-", "_");
        String spacesReplaced = minusReplaced.replaceAll(" ", "_");
        String mappedName = spacesReplaced;
        switch (spacesReplaced) {
            case "TODO":
                mappedName = "TODO";
                break;
        }
        return mappedName;
    }


    static BacNetObject b(String name, BacNetProperty... properties) {
        return new BacNetObject(name, List.of(properties));
    }

    static BacNetProperty p(String propertyIdentifier, String propertyDataType, ConformanceCode conformanceCode) {
        return new BacNetProperty(propertyIdentifier, propertyDataType, conformanceCode);
    }

    static ConformanceCode c(ConformanceCode.Type type, Integer... footnotes) {
        return new ConformanceCode(type, List.of(footnotes));
    }

    static class BacNetObject {
        final String name;
        final List<BacNetProperty> properties;

        BacNetObject(String name, List<BacNetProperty> properties) {
            this.name = name;
            this.properties = properties;
        }
    }

    static class BacNetProperty {
        final String propertyIdentifier;
        final String propertyDataType;
        final ConformanceCode conformanceCode;

        BacNetProperty(String propertyIdentifier, String propertyDataType, ConformanceCode conformanceCode) {
            this.propertyIdentifier = propertyIdentifier;
            this.propertyDataType = propertyDataType;
            this.conformanceCode = conformanceCode;
        }
    }

    static class ConformanceCode {
        final Type type;
        final List<Integer> additionalFootnotes;

        public ConformanceCode(Type type, List<Integer> additionalFootnotes) {
            this.type = type;
            this.additionalFootnotes = additionalFootnotes;
        }

        enum Type {
            R, W, O
        }
    }

    static class PropertyTypeCombination implements Comparable<PropertyTypeCombination> {
        final String propertyIdentifier;
        final String propertyDataType;

        PropertyTypeCombination(BacNetProperty bacNetProperty) {
            propertyIdentifier = bacNetProperty.propertyIdentifier;
            propertyDataType = bacNetProperty.propertyDataType;
        }

        PropertyTypeCombination(String propertyIdentifier, String propertyDataType) {
            this.propertyIdentifier = propertyIdentifier;
            this.propertyDataType = propertyDataType;
        }

        @Override
        public boolean equals(Object o) {
            if (this == o) return true;
            if (o == null || getClass() != o.getClass()) return false;
            PropertyTypeCombination that = (PropertyTypeCombination) o;
            return propertyIdentifier.equals(that.propertyIdentifier) && propertyDataType.equals(that.propertyDataType);
        }

        @Override
        public int hashCode() {
            return Objects.hash(propertyIdentifier, propertyDataType);
        }

        @Override
        public String toString() {
            return propertyIdentifier + ':' + propertyDataType;
        }


        @Override
        public int compareTo(PropertyTypeCombination propertyTypeCombination) {
            return this.toString().compareTo(propertyTypeCombination.toString());
        }
    }
}
