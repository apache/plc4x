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
package org.apache.plc4x.java.s7.readwrite.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.s7.readwrite.AlarmType;
import org.apache.plc4x.java.s7.readwrite.EventType;
import org.apache.plc4x.java.s7.readwrite.types.S7SubscriptionFieldType;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class S7SubscriptionField implements PlcField {

    //Event Subscription 
    private static final Pattern EVENT_SUBSCRIPTION_TYPE_PATTERN =
        //Pattern.compile("(^MODE)|(^SYS)|(^USR)|(^ALM_S)|(^ALM_8)");
        Pattern.compile("(^MODE)|(^SYS)|(^USR)|(^ALM)");
    //Event ack
    private static final Pattern EVENT_ALARM_ACK_PATTERN =
        Pattern.compile("(^ACK:)(((?:,{0,1})(16#[0-9a-fA-F]{8}))+)");

    //Query alarms from PLC.
    //TODO: Query SCAN 
    private static final Pattern EVENT_ALARM_QUERY_PATTERN =
        Pattern.compile("(^QUERY:)((ALARM_S)|(ALARM_8))");

    private final String address;
    private final S7SubscriptionFieldType fieldType;
    private final EventType eventtype;
    private final S7Field s7field;
    private final ArrayList<Integer> ackAlarms;
    private final AlarmType alarmQueryType;

    public S7SubscriptionField(String address, S7SubscriptionFieldType fieldType, EventType eventtype) {
        this.address = address;
        this.fieldType = fieldType;
        this.eventtype = eventtype;
        this.s7field = null;
        this.ackAlarms = null;
        this.alarmQueryType = null;
    }

    public S7SubscriptionField(String address, S7SubscriptionFieldType fieldType, ArrayList<Integer> ackAlarms) {
        this.address = address;
        this.fieldType = fieldType;
        this.eventtype = null;
        this.s7field = null;
        this.ackAlarms = ackAlarms;
        this.alarmQueryType = null;
    }

    public S7SubscriptionField(String address, S7SubscriptionFieldType fieldType, AlarmType alarmQueryType) {
        this.address = address;
        this.fieldType = fieldType;
        this.eventtype = null;
        this.s7field = null;
        this.ackAlarms = null;
        this.alarmQueryType = alarmQueryType;
    }

    public S7SubscriptionField(String address, S7SubscriptionFieldType fieldType, S7Field s7field) {
        this.address = address;
        this.fieldType = fieldType;
        this.eventtype = null;
        this.s7field = s7field;
        this.ackAlarms = null;
        this.alarmQueryType = null;
    }

    @Override
    public String getAddressString() {
        return address;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.RAW_BYTE_ARRAY;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return Collections.emptyList();
    }

    public S7SubscriptionFieldType getFieldType() {
        return fieldType;
    }

    public EventType getEventType() {
        return eventtype;
    }

    public S7Field getS7field() {
        return s7field;
    }

    public ArrayList<Integer> getAckAlarms() {
        return ackAlarms;
    }

    public AlarmType getAlarmQueryType() {
        return alarmQueryType;
    }

    public static boolean matches(String fieldString) {
        return EVENT_SUBSCRIPTION_TYPE_PATTERN.matcher(fieldString).matches() ||
            EVENT_ALARM_ACK_PATTERN.matcher(fieldString).matches() ||
            EVENT_ALARM_QUERY_PATTERN.matcher(fieldString).matches() ||
            S7Field.matches(fieldString);
    }

    public static S7SubscriptionField of(String fieldString) {
        {
            Matcher matcher = EVENT_SUBSCRIPTION_TYPE_PATTERN.matcher(fieldString);
            if (matcher.matches()) {
                return new S7SubscriptionField(fieldString, S7SubscriptionFieldType.EVENT_SUBSCRIPTION,
                    EventType.valueOf(fieldString));
            }
        }

        {
            //TODO: Actually only ALARM_S (SIG_1)
            Matcher matcher = EVENT_ALARM_ACK_PATTERN.matcher(fieldString);
            if (matcher.matches()) {
                String EventIds = matcher.group(2);
                String[] arrStrEventId = EventIds.split(",");
                ArrayList<Integer> arrEventId = new ArrayList<>();
                for (String EventId : arrStrEventId) {
                    EventId = EventId.replaceAll("16#", "");
                    arrEventId.add(Integer.parseInt(EventId, 16));
                }
                return new S7SubscriptionField(fieldString, S7SubscriptionFieldType.ALARM_ACK, arrEventId);
            }
        }

        {
            //TODO: Support for ALARM_8            
            Matcher matcher = EVENT_ALARM_QUERY_PATTERN.matcher(fieldString);
            if (matcher.matches()) {
                return new S7SubscriptionField(fieldString, S7SubscriptionFieldType.ALARM_QUERY, AlarmType.ALARM_S);
            }
        }

        {
            if (S7Field.matches(fieldString)) {
                S7Field s7field = S7Field.of(fieldString);
                switch (s7field.getDataType()) {
                    case BYTE:

                    default:

                }
                return new S7SubscriptionField(fieldString, S7SubscriptionFieldType.CYCLIC_SUBSCRIPTION,
                    s7field);
            }
        }
        throw new PlcInvalidFieldException("Unable to parse address: " + fieldString);
    }

}
