/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.bacnetip.readwrite.utils;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.bacnetip.readwrite.*;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.List;

import static org.apache.plc4x.java.spi.generation.WithReaderWriterArgs.WithAdditionalStringRepresentation;

public class StaticHelper {

    public static final Logger LOGGER = LoggerFactory.getLogger(StaticHelper.class);

    public static BACnetPropertyIdentifier readPropertyIdentifier(ReadBuffer readBuffer, Long actualLength) throws ParseException {
        int bitsToRead = (int) (actualLength * 8);
        long readUnsignedLong = readBuffer.readUnsignedLong("propertyIdentifier", bitsToRead);
        if (!BACnetPropertyIdentifier.isDefined(readUnsignedLong)) {
            return BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE;
        }
        return BACnetPropertyIdentifier.enumForValue(readUnsignedLong);
    }

    public static void writePropertyIdentifier(WriteBuffer writeBuffer, BACnetPropertyIdentifier value) throws SerializationException {
        if (value == null || value == BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        int bitsToWrite;
        long valueValue = value.getValue();
        if (valueValue <= 0xffL) {
            bitsToWrite = 8;
        } else if (valueValue <= 0xffffL) {
            bitsToWrite = 16;
        } else if (valueValue <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong("propertyIdentifier", bitsToWrite, valueValue, WithAdditionalStringRepresentation(value.name()));
    }

    public static void writeProprietaryPropertyIdentifier(WriteBuffer writeBuffer, BACnetPropertyIdentifier baCnetPropertyIdentifier, long value) throws SerializationException {
        if (baCnetPropertyIdentifier != null && baCnetPropertyIdentifier != BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        int bitsToWrite;
        if (value <= 0xffL) {
            bitsToWrite = 8;
        } else if (value <= 0xffffL) {
            bitsToWrite = 16;
        } else if (value <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong("proprietaryPropertyIdentifier", bitsToWrite, value, WithAdditionalStringRepresentation(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE.name()));
    }

    public static Long readProprietaryPropertyIdentifier(ReadBuffer readBuffer, BACnetPropertyIdentifier value, Long actualLength) throws ParseException {
        if (value != null && value != BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE) {
            return 0L;
        }
        // We need to reset our reader to the position we read before
        readBuffer.reset((int) (readBuffer.getPos() - actualLength));
        int bitsToRead = (int) (actualLength * 8);
        return readBuffer.readUnsignedLong("proprietaryPropertyIdentifier", bitsToRead);
    }

    public static BACnetEventType readEventType(ReadBuffer readBuffer, Long actualLength) throws ParseException {
        int bitsToRead = (int) (actualLength * 8);
        int readUnsignedLong = readBuffer.readUnsignedInt("eventType", bitsToRead);
        if (!BACnetEventType.isDefined(readUnsignedLong)) {
            return BACnetEventType.VENDOR_PROPRIETARY_VALUE;
        }
        return BACnetEventType.enumForValue(readUnsignedLong);
    }

    public static void writeEventType(WriteBuffer writeBuffer, BACnetEventType value) throws SerializationException {
        if (value == null || value == BACnetEventType.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        int bitsToWrite;
        long valueValue = value.getValue();
        if (valueValue <= 0xffL) {
            bitsToWrite = 8;
        } else if (valueValue <= 0xffffL) {
            bitsToWrite = 16;
        } else if (valueValue <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong("eventType", bitsToWrite, valueValue, WithAdditionalStringRepresentation(value.name()));
    }

    public static Long readProprietaryEventType(ReadBuffer readBuffer, BACnetEventType value, Long actualLength) throws ParseException {
        if (value != null && value != BACnetEventType.VENDOR_PROPRIETARY_VALUE) {
            return 0L;
        }
        // We need to reset our reader to the position we read before
        readBuffer.reset((int) (readBuffer.getPos() - actualLength));
        int bitsToRead = (int) (actualLength * 8);
        return readBuffer.readUnsignedLong("proprietaryEventType", bitsToRead);
    }

    public static void writeProprietaryEventType(WriteBuffer writeBuffer, BACnetEventType eventType, long value) throws SerializationException {
        if (eventType != null && eventType != BACnetEventType.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        int bitsToWrite;
        if (value <= 0xffL) {
            bitsToWrite = 8;
        } else if (value <= 0xffffL) {
            bitsToWrite = 16;
        } else if (value <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong("proprietaryEventType", bitsToWrite, value, WithAdditionalStringRepresentation(BACnetEventType.VENDOR_PROPRIETARY_VALUE.name()));
    }

    public static BACnetEventState readEventState(ReadBuffer readBuffer, Long actualLength) throws ParseException {
        int bitsToRead = (int) (actualLength * 8);
        int readUnsignedLong = readBuffer.readUnsignedInt("eventState", bitsToRead);
        if (!BACnetEventState.isDefined(readUnsignedLong)) {
            return BACnetEventState.VENDOR_PROPRIETARY_VALUE;
        }
        return BACnetEventState.enumForValue(readUnsignedLong);
    }

    public static void writeEventState(WriteBuffer writeBuffer, BACnetEventState value) throws SerializationException {
        if (value == null || value == BACnetEventState.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        int bitsToWrite;
        long valueValue = value.getValue();
        if (valueValue <= 0xffL) {
            bitsToWrite = 8;
        } else if (valueValue <= 0xffffL) {
            bitsToWrite = 16;
        } else if (valueValue <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong("eventState", bitsToWrite, valueValue, WithAdditionalStringRepresentation(value.name()));
    }

    public static Long readProprietaryEventState(ReadBuffer readBuffer, BACnetEventState value, Long actualLength) throws ParseException {
        if (value != null && value != BACnetEventState.VENDOR_PROPRIETARY_VALUE) {
            return 0L;
        }
        // We need to reset our reader to the position we read before
        readBuffer.reset((int) (readBuffer.getPos() - actualLength));
        int bitsToRead = (int) (actualLength * 8);
        return readBuffer.readUnsignedLong("proprietaryEventState", bitsToRead);
    }

    public static void writeProprietaryEventState(WriteBuffer writeBuffer, BACnetEventState eventState, long value) throws SerializationException {
        if (eventState != null && eventState != BACnetEventState.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        int bitsToWrite;
        if (value <= 0xffL) {
            bitsToWrite = 8;
        } else if (value <= 0xffffL) {
            bitsToWrite = 16;
        } else if (value <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong("proprietaryEventState", bitsToWrite, value, WithAdditionalStringRepresentation(BACnetEventState.VENDOR_PROPRIETARY_VALUE.name()));
    }

    public static BACnetObjectType readObjectType(ReadBuffer readBuffer) throws ParseException {
        int readUnsignedLong = readBuffer.readUnsignedInt("objectType", 10);
        if (!BACnetObjectType.isDefined(readUnsignedLong)) {
            return BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
        }
        return BACnetObjectType.enumForValue(readUnsignedLong);
    }

    public static void writeObjectType(WriteBuffer writeBuffer, BACnetObjectType value) throws SerializationException {
        if (value == null || value == BACnetObjectType.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        writeBuffer.writeUnsignedLong("objectType", 10, value.getValue(), WithAdditionalStringRepresentation(value.name()));
    }

    public static Integer readProprietaryObjectType(ReadBuffer readBuffer, BACnetObjectType value) throws ParseException {
        if (value != null && value != BACnetObjectType.VENDOR_PROPRIETARY_VALUE) {
            return 0;
        }
        // We need to reset our reader to the position we read before
        // TODO: maybe we reset to much here because pos is byte based
        // we consume the leftover bits before we reset to avoid trouble // TODO: we really need bit precision on resetting
        readBuffer.readUnsignedInt(6);
        readBuffer.reset(readBuffer.getPos() - 2);
        return readBuffer.readUnsignedInt("proprietaryObjectType", 10);
    }

    public static void writeProprietaryObjectType(WriteBuffer writeBuffer, BACnetObjectType objectType, int value) throws SerializationException {
        if (objectType != null && objectType != BACnetObjectType.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        writeBuffer.writeUnsignedInt("proprietaryObjectType", 10, value, WithAdditionalStringRepresentation(BACnetObjectType.VENDOR_PROPRIETARY_VALUE.name()));
    }

    public static boolean isBACnetConstructedDataClosingTag(ReadBuffer readBuffer, boolean instantTerminate, int expectedTagNumber) {
        if (instantTerminate) {
            return true;
        }
        int oldPos = readBuffer.getPos();
        try {
            // TODO: add graceful exit if we know already that we are at the end (we might need to add available bytes to reader)
            int tagNumber = readBuffer.readUnsignedInt(4);
            boolean isContextTag = readBuffer.readBit();
            int tagValue = readBuffer.readUnsignedInt(3);

            boolean foundOurClosingTag = isContextTag && tagNumber == expectedTagNumber && tagValue == 0x7;
            LOGGER.debug("Checking at pos pos:{}: tagNumber:{}, isContextTag:{}, tagValue:{}, expectedTagNumber:{}. foundOurClosingTag:{}", oldPos, tagNumber, isContextTag, tagValue, expectedTagNumber, foundOurClosingTag);
            return foundOurClosingTag;
        } catch (ParseException e) {
            LOGGER.warn("Error reading termination bit", e);
            return true;
        } catch (ArrayIndexOutOfBoundsException e) {
            LOGGER.debug("Reached EOF at {}", oldPos, e);
            return true;
        } finally {
            readBuffer.reset(oldPos);
        }
    }

    public static BACnetDataType guessDataType(BACnetObjectType objectType, BACnetContextTagPropertyIdentifier propertyIdentifier) {
        switch (propertyIdentifier.getPropertyIdentifier()) {
            case ABSENTEE_LIMIT:
            case ACCEPTED_MODES:
            case ACCESS_ALARM_EVENTS:
            case ACCESS_DOORS:
            case ACCESS_EVENT:
            case ACCESS_EVENT_AUTHENTICATION_FACTOR:
            case ACCESS_EVENT_CREDENTIAL:
            case ACCESS_EVENT_TAG:
            case ACCESS_EVENT_TIME:
            case ACCESS_TRANSACTION_EVENTS:
            case ACCOMPANIMENT:
            case ACCOMPANIMENT_TIME:
            case ACK_REQUIRED:
            case ACKED_TRANSITIONS:
            case ACTION:
            case ACTION_TEXT:
            case ACTIVATION_TIME:
            case ACTIVE_AUTHENTICATION_POLICY:
            case ACTIVE_COV_MULTIPLE_SUBSCRIPTIONS:
            case ACTIVE_COV_SUBSCRIPTIONS:
            case ACTIVE_TEXT:
            case ACTIVE_VT_SESSIONS:
            case ACTUAL_SHED_LEVEL:
            case ADJUST_VALUE:
            case ALARM_VALUE:
            case ALARM_VALUES:
            case ALIGN_INTERVALS:
            case ALL:
            case ALL_WRITES_SUCCESSFUL:
            case ALLOW_GROUP_DELAY_INHIBIT:
            case APDU_LENGTH:
            case APDU_SEGMENT_TIMEOUT:
            case APDU_TIMEOUT:
            case APPLICATION_SOFTWARE_VERSION:
            case ARCHIVE:
            case ASSIGNED_ACCESS_RIGHTS:
            case ASSIGNED_LANDING_CALLS:
            case ATTEMPTED_SAMPLES:
            case AUTHENTICATION_FACTORS:
            case AUTHENTICATION_POLICY_LIST:
            case AUTHENTICATION_POLICY_NAMES:
            case AUTHENTICATION_STATUS:
            case AUTHORIZATION_EXEMPTIONS:
            case AUTHORIZATION_MODE:
            case AUTO_SLAVE_DISCOVERY:
            case AVERAGE_VALUE:
            case BACKUP_AND_RESTORE_STATE:
            case BACKUP_FAILURE_TIMEOUT:
            case BACKUP_PREPARATION_TIME:
            case BACNET_IP_GLOBAL_ADDRESS:
            case BACNET_IP_MODE:
            case BACNET_IP_MULTICAST_ADDRESS:
            case BACNET_IP_NAT_TRAVERSAL:
            case BACNET_IP_UDP_PORT:
            case BACNET_IPV6_MODE:
            case BACNET_IPV6_UDP_PORT:
            case BACNET_IPV6_MULTICAST_ADDRESS:
            case BASE_DEVICE_SECURITY_POLICY:
            case BBMD_ACCEPT_FD_REGISTRATIONS:
            case BBMD_BROADCAST_DISTRIBUTION_TABLE:
            case BBMD_FOREIGN_DEVICE_TABLE:
            case BELONGS_TO:
            case BIAS:
            case BIT_MASK:
            case BIT_TEXT:
            case BLINK_WARN_ENABLE:
            case BUFFER_SIZE:
            case CAR_ASSIGNED_DIRECTION:
            case CAR_DOOR_COMMAND:
            case CAR_DOOR_STATUS:
            case CAR_DOOR_TEXT:
            case CAR_DOOR_ZONE:
            case CAR_DRIVE_STATUS:
            case CAR_LOAD:
            case CAR_LOAD_UNITS:
            case CAR_MODE:
            case CAR_MOVING_DIRECTION:
            case CAR_POSITION:
            case CHANGE_OF_STATE_COUNT:
            case CHANGE_OF_STATE_TIME:
            case CHANGES_PENDING:
            case CHANNEL_NUMBER:
            case CLIENT_COV_INCREMENT:
            case COMMAND:
            case COMMAND_TIME_ARRAY:
            case CONFIGURATION_FILES:
            case CONTROL_GROUPS:
            case CONTROLLED_VARIABLE_REFERENCE:
            case CONTROLLED_VARIABLE_UNITS:
            case CONTROLLED_VARIABLE_VALUE:
            case COUNT:
            case COUNT_BEFORE_CHANGE:
            case COUNT_CHANGE_TIME:
            case COV_INCREMENT:
            case COV_PERIOD:
            case COV_RESUBSCRIPTION_INTERVAL:
            case COVU_PERIOD:
            case COVU_RECIPIENTS:
            case CREDENTIAL_DISABLE:
            case CREDENTIAL_STATUS:
            case CREDENTIALS:
            case CREDENTIALS_IN_ZONE:
            case CURRENT_COMMAND_PRIORITY:
            case DATABASE_REVISION:
            case DATE_LIST:
            case DAYLIGHT_SAVINGS_STATUS:
            case DAYS_REMAINING:
            case DEADBAND:
            case DEFAULT_FADE_TIME:
            case DEFAULT_RAMP_RATE:
            case DEFAULT_STEP_INCREMENT:
            case DEFAULT_SUBORDINATE_RELATIONSHIP:
            case DEFAULT_TIMEOUT:
            case DEPLOYED_PROFILE_LOCATION:
            case DERIVATIVE_CONSTANT:
            case DERIVATIVE_CONSTANT_UNITS:
            case DESCRIPTION:
            case DESCRIPTION_OF_HALT:
            case DEVICE_ADDRESS_BINDING:
            case DEVICE_TYPE:
            case DIRECT_READING:
            case DISTRIBUTION_KEY_REVISION:
            case DO_NOT_HIDE:
            case DOOR_ALARM_STATE:
            case DOOR_EXTENDED_PULSE_TIME:
            case DOOR_MEMBERS:
            case DOOR_OPEN_TOO_LONG_TIME:
            case DOOR_PULSE_TIME:
            case DOOR_STATUS:
            case DOOR_UNLOCK_DELAY_TIME:
            case DUTY_WINDOW:
            case EFFECTIVE_PERIOD:
            case EGRESS_ACTIVE:
            case EGRESS_TIME:
            case ELAPSED_ACTIVE_TIME:
            case ELEVATOR_GROUP:
            case ENABLE:
            case ENERGY_METER:
            case ENERGY_METER_REF:
            case ENTRY_POINTS:
            case ERROR_LIMIT:
            case ESCALATOR_MODE:
            case EVENT_ALGORITHM_INHIBIT:
            case EVENT_ALGORITHM_INHIBIT_REF:
            case EVENT_DETECTION_ENABLE:
            case EVENT_ENABLE:
            case EVENT_MESSAGE_TEXTS:
            case EVENT_MESSAGE_TEXTS_CONFIG:
            case EVENT_PARAMETERS:
            case EVENT_STATE:
            case EVENT_TIME_STAMPS:
            case EVENT_TYPE:
            case EXCEPTION_SCHEDULE:
            case EXECUTION_DELAY:
            case EXIT_POINTS:
            case EXPECTED_SHED_LEVEL:
            case EXPIRATION_TIME:
            case EXTENDED_TIME_ENABLE:
            case FAILED_ATTEMPT_EVENTS:
            case FAILED_ATTEMPTS:
            case FAILED_ATTEMPTS_TIME:
            case FAULT_HIGH_LIMIT:
            case FAULT_LOW_LIMIT:
            case FAULT_PARAMETERS:
            case FAULT_SIGNALS:
            case FAULT_TYPE:
            case FAULT_VALUES:
            case FD_BBMD_ADDRESS:
            case FD_SUBSCRIPTION_LIFETIME:
            case FEEDBACK_VALUE:
            case FILE_ACCESS_METHOD:
            case FILE_SIZE:
            case FILE_TYPE:
            case FIRMWARE_REVISION:
            case FLOOR_TEXT:
            case FULL_DUTY_BASELINE:
            case GLOBAL_IDENTIFIER:
            case GROUP_ID:
            case GROUP_MEMBER_NAMES:
            case GROUP_MEMBERS:
            case GROUP_MODE:
            case HIGH_LIMIT:
            case HIGHER_DECK:
            case IN_PROCESS:
            case IN_PROGRESS:
            case INACTIVE_TEXT:
            case INITIAL_TIMEOUT:
            case INPUT_REFERENCE:
            case INSTALLATION_ID:
            case INSTANCE_OF:
            case INSTANTANEOUS_POWER:
            case INTEGRAL_CONSTANT:
            case INTEGRAL_CONSTANT_UNITS:
            case INTERFACE_VALUE:
            case INTERVAL_OFFSET:
            case IP_ADDRESS:
            case IP_DEFAULT_GATEWAY:
            case IP_DHCP_ENABLE:
            case IP_DHCP_LEASE_TIME:
            case IP_DHCP_LEASE_TIME_REMAINING:
            case IP_DHCP_SERVER:
            case IP_DNS_SERVER:
            case IP_SUBNET_MASK:
            case IPV6_ADDRESS:
            case IPV6_AUTO_ADDRESSING_ENABLE:
            case IPV6_DEFAULT_GATEWAY:
            case IPV6_DHCP_LEASE_TIME:
            case IPV6_DHCP_LEASE_TIME_REMAINING:
            case IPV6_DHCP_SERVER:
            case IPV6_DNS_SERVER:
            case IPV6_PREFIX_LENGTH:
            case IPV6_ZONE_INDEX:
            case IS_UTC:
            case KEY_SETS:
            case LANDING_CALL_CONTROL:
            case LANDING_CALLS:
            case LANDING_DOOR_STATUS:
            case LAST_ACCESS_EVENT:
            case LAST_ACCESS_POINT:
            case LAST_COMMAND_TIME:
            case LAST_CREDENTIAL_ADDED:
            case LAST_CREDENTIAL_ADDED_TIME:
            case LAST_CREDENTIAL_REMOVED:
            case LAST_CREDENTIAL_REMOVED_TIME:
            case LAST_KEY_SERVER:
            case LAST_NOTIFY_RECORD:
            case LAST_PRIORITY:
            case LAST_RESTART_REASON:
            case LAST_RESTORE_TIME:
            case LAST_STATE_CHANGE:
            case LAST_USE_TIME:
            case LIFE_SAFETY_ALARM_VALUES:
            case LIGHTING_COMMAND:
            case LIGHTING_COMMAND_DEFAULT_PRIORITY:
            case LIMIT_ENABLE:
            case LIMIT_MONITORING_INTERVAL:
            case LINK_SPEED:
            case LINK_SPEED_AUTONEGOTIATE:
            case LINK_SPEEDS:
            case LIST_OF_GROUP_MEMBERS:
            case LIST_OF_OBJECT_PROPERTY_REFERENCES:
            case LOCAL_DATE:
            case LOCAL_FORWARDING_ONLY:
            case LOCAL_TIME:
            case LOCATION:
            case LOCK_STATUS:
            case LOCKOUT:
            case LOCKOUT_RELINQUISH_TIME:
            case LOG_BUFFER:
            case LOG_DEVICE_OBJECT_PROPERTY:
            case LOG_INTERVAL:
            case LOGGING_OBJECT:
            case LOGGING_RECORD:
            case LOGGING_TYPE:
            case LOW_DIFF_LIMIT:
            case LOW_LIMIT:
            case LOWER_DECK:
            case MAC_ADDRESS:
            case MACHINE_ROOM_ID:
            case MAINTENANCE_REQUIRED:
            case MAKING_CAR_CALL:
            case MANIPULATED_VARIABLE_REFERENCE:
            case MANUAL_SLAVE_ADDRESS_BINDING:
            case MASKED_ALARM_VALUES:
            case MAX_ACTUAL_VALUE:
            case MAX_APDU_LENGTH_ACCEPTED:
            case MAX_FAILED_ATTEMPTS:
            case MAX_INFO_FRAMES:
            case MAX_MASTER:
            case MAX_PRES_VALUE:
            case MAX_SEGMENTS_ACCEPTED:
            case MAXIMUM_OUTPUT:
            case MAXIMUM_VALUE:
            case MAXIMUM_VALUE_TIMESTAMP:
            case MEMBER_OF:
            case MEMBER_STATUS_FLAGS:
            case MEMBERS:
            case MIN_ACTUAL_VALUE:
            case MIN_PRES_VALUE:
            case MINIMUM_OFF_TIME:
            case MINIMUM_ON_TIME:
            case MINIMUM_OUTPUT:
            case MINIMUM_VALUE:
            case MINIMUM_VALUE_TIMESTAMP:
            case MODE:
            case MODEL_NAME:
            case MODIFICATION_DATE:
            case MUSTER_POINT:
            case NEGATIVE_ACCESS_RULES:
            case NETWORK_ACCESS_SECURITY_POLICIES:
            case NETWORK_INTERFACE_NAME:
            case NETWORK_NUMBER:
            case NETWORK_NUMBER_QUALITY:
            case NETWORK_TYPE:
            case NEXT_STOPPING_FLOOR:
            case NODE_SUBTYPE:
            case NODE_TYPE:
            case NOTIFICATION_CLASS:
            case NOTIFICATION_THRESHOLD:
            case NOTIFY_TYPE:
            case NUMBER_OF_APDU_RETRIES:
            case NUMBER_OF_AUTHENTICATION_POLICIES:
            case NUMBER_OF_STATES:
            case OBJECT_IDENTIFIER:
            case OBJECT_LIST:
            case OBJECT_NAME:
            case OBJECT_PROPERTY_REFERENCE:
            case OBJECT_TYPE:
            case OCCUPANCY_COUNT:
            case OCCUPANCY_COUNT_ADJUST:
            case OCCUPANCY_COUNT_ENABLE:
            case OCCUPANCY_LOWER_LIMIT:
            case OCCUPANCY_LOWER_LIMIT_ENFORCED:
            case OCCUPANCY_STATE:
            case OCCUPANCY_UPPER_LIMIT:
            case OCCUPANCY_UPPER_LIMIT_ENFORCED:
            case OPERATION_DIRECTION:
            case OPERATION_EXPECTED:
            case OPTIONAL:
            case OUT_OF_SERVICE:
            case OUTPUT_UNITS:
            case PACKET_REORDER_TIME:
            case PASSBACK_MODE:
            case PASSBACK_TIMEOUT:
            case PASSENGER_ALARM:
            case POLARITY:
            case PORT_FILTER:
            case POSITIVE_ACCESS_RULES:
            case POWER:
            case POWER_MODE:
            case PRESCALE:
            case PRESENT_VALUE:
            case PRIORITY:
            case PRIORITY_ARRAY:
            case PRIORITY_FOR_WRITING:
            case PROCESS_IDENTIFIER:
            case PROCESS_IDENTIFIER_FILTER:
            case PROFILE_LOCATION:
            case PROFILE_NAME:
            case PROGRAM_CHANGE:
            case PROGRAM_LOCATION:
            case PROGRAM_STATE:
            case PROPERTY_LIST:
            case PROPORTIONAL_CONSTANT:
            case PROPORTIONAL_CONSTANT_UNITS:
            case PROTOCOL_LEVEL:
            case PROTOCOL_CONFORMANCE_CLASS:
            case PROTOCOL_OBJECT_TYPES_SUPPORTED:
            case PROTOCOL_REVISION:
            case PROTOCOL_SERVICES_SUPPORTED:
            case PROTOCOL_VERSION:
            case PULSE_RATE:
            case READ_ONLY:
            case REASON_FOR_DISABLE:
            case REASON_FOR_HALT:
            case RECIPIENT_LIST:
            case RECORD_COUNT:
            case RECORDS_SINCE_NOTIFICATION:
            case REFERENCE_PORT:
            case REGISTERED_CAR_CALL:
            case RELIABILITY:
            case RELIABILITY_EVALUATION_INHIBIT:
            case RELINQUISH_DEFAULT:
            case REPRESENTS:
            case REQUESTED_SHED_LEVEL:
            case REQUESTED_UPDATE_INTERVAL:
            case REQUIRED:
            case RESOLUTION:
            case RESTART_NOTIFICATION_RECIPIENTS:
            case RESTORE_COMPLETION_TIME:
            case RESTORE_PREPARATION_TIME:
            case ROUTING_TABLE:
            case SCALE:
            case SCALE_FACTOR:
            case SCHEDULE_DEFAULT:
            case SECURED_STATUS:
            case SECURITY_PDU_TIMEOUT:
            case SECURITY_TIME_WINDOW:
            case SEGMENTATION_SUPPORTED:
            case SERIAL_NUMBER:
            case SETPOINT:
            case SETPOINT_REFERENCE:
            case SETTING:
            case SHED_DURATION:
            case SHED_LEVEL_DESCRIPTIONS:
            case SHED_LEVELS:
            case SILENCED:
            case SLAVE_ADDRESS_BINDING:
            case SLAVE_PROXY_ENABLE:
            case START_TIME:
            case STATE_CHANGE_VALUES:
            case STATE_DESCRIPTION:
            case STATE_TEXT:
            case STATUS_FLAGS:
            case STOP_TIME:
            case STOP_WHEN_FULL:
            case STRIKE_COUNT:
            case STRUCTURED_OBJECT_LIST:
            case SUBORDINATE_ANNOTATIONS:
            case SUBORDINATE_LIST:
            case SUBORDINATE_NODE_TYPES:
            case SUBORDINATE_RELATIONSHIPS:
            case SUBORDINATE_TAGS:
            case SUBSCRIBED_RECIPIENTS:
            case SUPPORTED_FORMAT_CLASSES:
            case SUPPORTED_FORMATS:
            case SUPPORTED_SECURITY_ALGORITHMS:
            case SYSTEM_STATUS:
            case TAGS:
            case THREAT_AUTHORITY:
            case THREAT_LEVEL:
            case TIME_DELAY:
            case TIME_DELAY_NORMAL:
            case TIME_OF_ACTIVE_TIME_RESET:
            case TIME_OF_DEVICE_RESTART:
            case TIME_OF_STATE_COUNT_RESET:
            case TIME_OF_STRIKE_COUNT_RESET:
            case TIME_SYNCHRONIZATION_INTERVAL:
            case TIME_SYNCHRONIZATION_RECIPIENTS:
            case TIMER_RUNNING:
            case TIMER_STATE:
            case TOTAL_RECORD_COUNT:
            case TRACE_FLAG:
            case TRACKING_VALUE:
            case TRANSACTION_NOTIFICATION_CLASS:
            case TRANSITION:
            case TRIGGER:
            case UNITS:
            case UPDATE_INTERVAL:
            case UPDATE_KEY_SET_TIMEOUT:
            case UPDATE_TIME:
            case USER_EXTERNAL_IDENTIFIER:
            case USER_INFORMATION_REFERENCE:
            case USER_NAME:
            case USER_TYPE:
            case USES_REMAINING:
            case UTC_OFFSET:
            case UTC_TIME_SYNCHRONIZATION_RECIPIENTS:
            case VALID_SAMPLES:
            case VALUE_BEFORE_CHANGE:
            case VALUE_CHANGE_TIME:
            case VALUE_SET:
            case VALUE_SOURCE:
            case VALUE_SOURCE_ARRAY:
            case VARIANCE_VALUE:
            case VENDOR_IDENTIFIER:
            case VENDOR_NAME:
            case VERIFICATION_TIME:
            case VIRTUAL_MAC_ADDRESS_TABLE:
            case VT_CLASSES_SUPPORTED:
            case WEEKLY_SCHEDULE:
            case WINDOW_INTERVAL:
            case WINDOW_SAMPLES:
            case WRITE_STATUS:
            case ZONE_FROM:
            case ZONE_MEMBERS:
            case ZONE_TO:
            case VENDOR_PROPRIETARY_VALUE:
        }
        switch (objectType) {
            case ACCESS_CREDENTIAL:
            case ACCESS_DOOR:
            case ACCESS_POINT:
            case ACCESS_RIGHTS:
            case ACCESS_USER:
            case ACCESS_ZONE:
            case ACCUMULATOR:
                break;
            case ALERT_ENROLLMENT:
                break;
            case ANALOG_INPUT:
            case ANALOG_OUTPUT:
                return BACnetDataType.REAL;
            case ANALOG_VALUE:
            case AVERAGING:
            case BINARY_INPUT:
            case BINARY_LIGHTING_OUTPUT:
            case BINARY_OUTPUT:
            case BINARY_VALUE:
                break;
            case BITSTRING_VALUE:
                return BACnetDataType.BIT_STRING;
            case CALENDAR:
            case CHANNEL:
            case CHARACTERSTRING_VALUE:
            case COMMAND:
            case CREDENTIAL_DATA_INPUT:
            case DATEPATTERN_VALUE:
            case DATE_VALUE:
            case DATETIMEPATTERN_VALUE:
            case DATETIME_VALUE:
            case DEVICE:
            case ELEVATOR_GROUP:
            case ESCALATOR:
            case EVENT_ENROLLMENT:
            case EVENT_LOG:
            case FILE:
            case GLOBAL_GROUP:
            case GROUP:
                break;
            case INTEGER_VALUE:
                return BACnetDataType.SIGNED_INTEGER;
            case LARGE_ANALOG_VALUE:
            case LIFE_SAFETY_POINT:
                // TODO: temporary
                return BACnetDataType.BACNET_PROPERTY_IDENTIFIER;
            case LIFE_SAFETY_ZONE:
                return BACnetDataType.BACNET_OBJECT_IDENTIFIER;
            case LIFT:
            case LIGHTING_OUTPUT:
                break;
            case LOAD_CONTROL:
                // TODO: temporary // FIXME: this is just so tags get consumed
                return BACnetDataType.ENUMERATED;
            case LOOP:
            case MULTI_STATE_INPUT:
            case MULTI_STATE_OUTPUT:
            case MULTI_STATE_VALUE:
            case NETWORK_PORT:
            case NETWORK_SECURITY:
            case NOTIFICATION_CLASS:
            case NOTIFICATION_FORWARDER:
            case OCTETSTRING_VALUE:
            case POSITIVE_INTEGER_VALUE:
            case PROGRAM:
            case PULSE_CONVERTER:
            case SCHEDULE:
                break;
            case STRUCTURED_VIEW:
                // TODO: temporary
                return BACnetDataType.BACNET_OBJECT_IDENTIFIER;
            case TIMEPATTERN_VALUE:
            case TIME_VALUE:
            case TIMER:
            case TREND_LOG:
            case TREND_LOG_MULTIPLE:
                return BACnetDataType.ENUMERATED;
        }
        return BACnetDataType.ENUMERATED;
    }

    public static long parseVarUint(byte[] data) {
        if (data.length == 0) {
            return 0;
        }
        return new BigInteger(data).longValue();
    }

    public static byte[] writeVarUint(long value) {
        return BigInteger.valueOf(value).toByteArray();
    }

    public static BACnetTagHeader createBACnetTagHeaderBalanced(boolean isContext, short id, long value) {
        TagClass tagClass = TagClass.APPLICATION_TAGS;
        if (isContext) {
            tagClass = TagClass.CONTEXT_SPECIFIC_TAGS;
        }

        byte tagNumber;
        Short extTagNumber = null;
        if (id <= 14) {
            tagNumber = (byte) id;
        } else {
            tagNumber = 0xF;
            extTagNumber = id;
        }

        byte lengthValueType;
        Short extLength = null;
        Integer extExtLength = null;
        Long extExtExtLength = null;
        if (value <= 4) {
            lengthValueType = (byte) value;
        } else {
            lengthValueType = 5;
            // Depending on the length, we will either write it as an 8 bit, 32 bit, or 64 bit integer
            if (value <= 253) {
                extLength = (short) value;
            } else if (value <= 65535) {
                extLength = 254;
                extExtLength = (int) value;
            } else {
                extLength = 255;
                extExtExtLength = value;
            }
        }

        return new BACnetTagHeader(tagNumber, tagClass, lengthValueType, extTagNumber, extLength, extExtLength, extExtExtLength);
    }

    public static BACnetApplicationTagNull createBACnetApplicationTagNull() {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.NULL.getValue(), 0);
        return new BACnetApplicationTagNull(header);
    }

    public static BACnetContextTagNull createBACnetContextTagNull(byte tagNumber) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 0);
        return new BACnetContextTagNull(header, (short) tagNumber, true);
    }

    public static BACnetOpeningTag createBACnetOpeningTag(short tagNum) {
        byte tagNumber;
        Short extTagNumber = null;
        if (tagNum <= 14) {
            tagNumber = (byte) tagNum;
        } else {
            tagNumber = 0xF;
            extTagNumber = tagNum;
        }
        BACnetTagHeader header = new BACnetTagHeader(tagNumber, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) 0x6, extTagNumber, null, null, null);
        return new BACnetOpeningTag(header, tagNum, 0x6L);
    }

    public static BACnetClosingTag createBACnetClosingTag(short tagNum) {
        byte tagNumber;
        Short extTagNumber = null;
        if (tagNum <= 14) {
            tagNumber = (byte) tagNum;
        } else {
            tagNumber = 0xF;
            extTagNumber = tagNum;
        }
        BACnetTagHeader header = new BACnetTagHeader(tagNumber, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) 0x7, extTagNumber, null, null, null);
        return new BACnetClosingTag(header, tagNum, 0x7L);
    }

    public static BACnetApplicationTagObjectIdentifier createBACnetApplicationTagObjectIdentifier(int objectType, long instance) {
        BACnetTagHeader header = new BACnetTagHeader((byte) BACnetDataType.SIGNED_INTEGER.getValue(), TagClass.APPLICATION_TAGS, (byte) 4, null, null, null, null);
        BACnetObjectType objectTypeEnum = BACnetObjectType.enumForValue(objectType);
        int proprietaryValue = 0;
        if (objectType >= 128 || !BACnetObjectType.isDefined(objectType)) {
            objectTypeEnum = BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
            proprietaryValue = objectType;
        }
        BACnetTagPayloadObjectIdentifier payload = new BACnetTagPayloadObjectIdentifier(objectTypeEnum, proprietaryValue, instance);
        return new BACnetApplicationTagObjectIdentifier(header, payload);
    }

    public static BACnetContextTagObjectIdentifier createBACnetContextTagObjectIdentifier(byte tagNum, int objectType, long instance) {
        BACnetTagHeader header = new BACnetTagHeader(tagNum, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) 4, null, null, null, null);
        BACnetObjectType objectTypeEnum = BACnetObjectType.enumForValue(objectType);
        int proprietaryValue = 0;
        if (objectType >= 128 || !BACnetObjectType.isDefined(objectType)) {
            objectTypeEnum = BACnetObjectType.VENDOR_PROPRIETARY_VALUE;
            proprietaryValue = objectType;
        }
        BACnetTagPayloadObjectIdentifier payload = new BACnetTagPayloadObjectIdentifier(objectTypeEnum, proprietaryValue, instance);
        return new BACnetContextTagObjectIdentifier(header, payload, (short) tagNum, true);
    }

    public static BACnetContextTagPropertyIdentifier createBACnetContextTagPropertyIdentifier(byte tagNum, int propertyType) {
        BACnetPropertyIdentifier propertyIdentifier = BACnetPropertyIdentifier.enumForValue(propertyType);
        int proprietaryValue = 0;
        if (!BACnetPropertyIdentifier.isDefined(propertyType)) {
            propertyIdentifier = BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE;
            proprietaryValue = propertyType;
        }
        BACnetTagHeader header = new BACnetTagHeader(tagNum, TagClass.CONTEXT_SPECIFIC_TAGS, (byte) requiredLength(propertyType), null, null, null, null);
        return new BACnetContextTagPropertyIdentifier(header, propertyIdentifier, proprietaryValue, (short) tagNum, true, 0L);
    }

    public static BACnetApplicationTagBoolean createBACnetApplicationTagBoolean(boolean value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.BOOLEAN.getValue(), value ? 1L : 0L);
        return new BACnetApplicationTagBoolean(header, new BACnetTagPayloadBoolean(value ? 1L : 0L));
    }

    public static BACnetContextTagBoolean createBACnetContextTagBoolean(byte tagNumber, boolean value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 1);
        return new BACnetContextTagBoolean(header, (short) (value ? 1 : 0), new BACnetTagPayloadBoolean(value ? 1L : 0L), (short) tagNumber, true);
    }

    public static BACnetApplicationTagUnsignedInteger createBACnetApplicationTagUnsignedInteger(long value) {
        Pair<Long, BACnetTagPayloadUnsignedInteger> lengthPayload = createUnsignedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.UNSIGNED_INTEGER.getValue(), lengthPayload.getLeft());
        return new BACnetApplicationTagUnsignedInteger(header, lengthPayload.getRight());
    }

    public static BACnetContextTagUnsignedInteger createBACnetContextTagUnsignedInteger(byte tagNumber, long value) {
        Pair<Long, BACnetTagPayloadUnsignedInteger> lengthPayload = createUnsignedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, lengthPayload.getLeft());
        return new BACnetContextTagUnsignedInteger(header, lengthPayload.getRight(), (short) tagNumber, true);
    }

    public static Pair<Long, BACnetTagPayloadUnsignedInteger> createUnsignedPayload(long value) {
        long length;
        Short valueUint8 = null;
        Integer valueUint16 = null;
        Long valueUint24 = null;
        Long valueUint32 = null;
        BigInteger valueUint40 = null;
        BigInteger valueUint48 = null;
        BigInteger valueUint56 = null;
        BigInteger valueUint64 = null;
        if (value < 0x100) {
            length = 1;
            valueUint8 = (short) value;
        } else if (value < 0x10000) {
            length = 2;
            valueUint16 = (int) value;
        } else if (value < 0x1000000) {
            length = 3;
            valueUint24 = value;
        } else {
            length = 4;
            valueUint32 = value;
        }
        BACnetTagPayloadUnsignedInteger payload = new BACnetTagPayloadUnsignedInteger(valueUint8, valueUint16, valueUint24, valueUint32, valueUint40, valueUint48, valueUint56, valueUint64, length);
        return Pair.of(length, payload);
    }

    public static BACnetApplicationTagSignedInteger createBACnetApplicationTagSignedInteger(long value) {
        Pair<Long, BACnetTagPayloadSignedInteger> lengthPayload = createSignedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.SIGNED_INTEGER.getValue(), lengthPayload.getLeft());
        return new BACnetApplicationTagSignedInteger(header, lengthPayload.getRight());
    }

    public static BACnetContextTagSignedInteger createBACnetContextTagSignedInteger(short tagNumber, long value) {
        Pair<Long, BACnetTagPayloadSignedInteger> lengthPayload = createSignedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, (byte) tagNumber, lengthPayload.getLeft());
        return new BACnetContextTagSignedInteger(header, lengthPayload.getRight(), tagNumber, true);
    }

    public static Pair<Long, BACnetTagPayloadSignedInteger> createSignedPayload(long value) {
        long length;
        Byte valueInt8 = null;
        Short valueInt16 = null;
        Integer valueInt24 = null;
        Integer valueInt32 = null;
        if (value < 0x100) {
            length = 1;
            valueInt8 = (byte) value;
        } else if (value < 0x10000) {
            length = 2;
            valueInt16 = (short) value;
        } else if (value < 0x1000000) {
            length = 3;
            valueInt24 = (int) value;
        } else {
            length = 4;
            valueInt32 = (int) value;
        }
        BACnetTagPayloadSignedInteger payload = new BACnetTagPayloadSignedInteger(valueInt8, valueInt16, valueInt24, valueInt32, null, null, null, null, length);
        return Pair.of(length, payload);
    }

    public static BACnetApplicationTagReal createBACnetApplicationTagReal(float value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.REAL.getValue(), 4);
        return new BACnetApplicationTagReal(header, new BACnetTagPayloadReal(value));
    }

    public static BACnetContextTagReal createBACnetContextTagReal(byte tagNumber, float value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 4);
        return new BACnetContextTagReal(header, new BACnetTagPayloadReal(value), (short) tagNumber, true);
    }

    public static BACnetApplicationTagDouble createBACnetApplicationTagDouble(double value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.DOUBLE.getValue(), 8);
        return new BACnetApplicationTagDouble(header, new BACnetTagPayloadDouble(value));
    }

    public static BACnetContextTagDouble createBACnetContextTagDouble(byte tagNumber, double value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 8);
        return new BACnetContextTagDouble(header, new BACnetTagPayloadDouble(value), (short) tagNumber, true);
    }

    public static BACnetApplicationTagOctetString createBACnetApplicationTagOctetString(String value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.OCTET_STRING.getValue(), value.length() + 1);
        return new BACnetApplicationTagOctetString(header, new BACnetTagPayloadOctetString(value, (long) value.length() + 1));
    }

    public static BACnetContextTagOctetString createBACnetContextTagOctetString(byte tagNumber, String value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, value.length() + 1);
        return new BACnetContextTagOctetString(header, new BACnetTagPayloadOctetString(value, (long) value.length() + 1), (short) tagNumber, true);
    }

    public static BACnetApplicationTagCharacterString createBACnetApplicationTagCharacterString(BACnetCharacterEncoding baCnetCharacterEncoding, String value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.CHARACTER_STRING.getValue(), value.length() + 1);
        return new BACnetApplicationTagCharacterString(header, new BACnetTagPayloadCharacterString(baCnetCharacterEncoding, value, (long) value.length() + 1));
    }

    public static BACnetContextTagCharacterString createBACnetContextTagCharacterString(byte tagNumber, BACnetCharacterEncoding baCnetCharacterEncoding, String value) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, value.length() + 1);
        return new BACnetContextTagCharacterString(header, new BACnetTagPayloadCharacterString(baCnetCharacterEncoding, value, (long) value.length() + 1), (short) tagNumber, true);
    }

    public static BACnetApplicationTagBitString createBACnetApplicationTagBitString(List<Boolean> value) {
        long numberOfBytesNeeded = (value.size() + 7) / 8;
        short unusedBits = (short) (8 - (value.size() % 8));
        if (unusedBits == 8) {
            unusedBits = 0;
        }
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.BIT_STRING.getValue(), numberOfBytesNeeded + 1);
        return new BACnetApplicationTagBitString(header, new BACnetTagPayloadBitString(unusedBits, value, new ArrayList<>(unusedBits), numberOfBytesNeeded + 1));
    }

    public static BACnetContextTagBitString createBACnetContextTagBitString(byte tagNumber, List<Boolean> value) {
        long numberOfBytesNeeded = (value.size() + 7) / 8;
        short unusedBits = (short) (8 - (value.size() % 8));
        if (unusedBits == 8) {
            unusedBits = 0;
        }
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, numberOfBytesNeeded + 1);
        return new BACnetContextTagBitString(header, new BACnetTagPayloadBitString(unusedBits, value, new ArrayList<>(unusedBits), numberOfBytesNeeded + 1), (short) tagNumber, true);
    }

    public static BACnetApplicationTagEnumerated createBACnetApplicationTagEnumerated(long value) {
        Pair<Long, BACnetTagPayloadEnumerated> lengthPayload = CreateEnumeratedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.ENUMERATED.getValue(), lengthPayload.getLeft());
        return new BACnetApplicationTagEnumerated(header, lengthPayload.getRight());
    }

    public static BACnetContextTagEnumerated createBACnetContextTagEnumerated(byte tagNumber, long value) {
        Pair<Long, BACnetTagPayloadEnumerated> lengthPayload = CreateEnumeratedPayload(value);
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, lengthPayload.getLeft());
        return new BACnetContextTagEnumerated(header, lengthPayload.getRight(), (short) tagNumber, true);
    }

    public static Pair<Long, BACnetTagPayloadEnumerated> CreateEnumeratedPayload(long value) {
        long length = requiredLength(value);
        byte[] data = writeVarUint(value);
        BACnetTagPayloadEnumerated payload = new BACnetTagPayloadEnumerated(data, length);
        return Pair.of(length, payload);
    }

    public static BACnetApplicationTagDate createBACnetApplicationTagDate(int year, short month, short dayOfMonth, short dayOfWeek) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.DATE.getValue(), 4);
        short yearMinus1900 = (short) (year - 1900);
        if (year == 0xFF) {
            yearMinus1900 = 0xFF;
        }
        return new BACnetApplicationTagDate(header, new BACnetTagPayloadDate(yearMinus1900, month, dayOfMonth, dayOfWeek));
    }

    public static BACnetContextTagDate createBACnetContextTagDate(byte tagNumber, int year, short month, short dayOfMonth, short dayOfWeek) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 4);
        short yearMinus1900 = (short) (year - 1900);
        if (year == 0xFF) {
            yearMinus1900 = 0xFF;
        }
        return new BACnetContextTagDate(header, new BACnetTagPayloadDate(yearMinus1900, month, dayOfMonth, dayOfWeek), (short) tagNumber, true);
    }

    public static BACnetApplicationTagTime createBACnetApplicationTagTime(short hour, short minute, short second, short fractional) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(false, BACnetDataType.TIME.getValue(), 4);
        return new BACnetApplicationTagTime(header, new BACnetTagPayloadTime(hour, minute, second, fractional));
    }

    public static BACnetContextTagTime createBACnetContextTagTime(byte tagNumber, short hour, short minute, short second, short fractional) {
        BACnetTagHeader header = createBACnetTagHeaderBalanced(true, tagNumber, 4);
        return new BACnetContextTagTime(header, new BACnetTagPayloadTime(hour, minute, second, fractional), (short) tagNumber, true);
    }

    private static long requiredLength(long value) {
        long length;
        if (value < 0x100)
            length = 1;
        else if (value < 0x10000)
            length = 2;
        else if (value < 0x1000000)
            length = 3;
        else
            length = 4;
        return length;
    }

    public static BACnetContextTagPropertyIdentifier dummyPropertyIdentifier() {
        return new BACnetContextTagPropertyIdentifier(null, BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE, 0L, (short) 0, true, 0L);
    }
}
