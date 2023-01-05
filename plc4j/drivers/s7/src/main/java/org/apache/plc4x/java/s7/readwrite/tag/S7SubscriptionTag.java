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
package org.apache.plc4x.java.s7.readwrite.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.s7.readwrite.AlarmType;
import org.apache.plc4x.java.s7.readwrite.EventType;
import org.apache.plc4x.java.s7.readwrite.types.S7SubscriptionType;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class S7SubscriptionTag implements PlcTag {

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
    private final S7SubscriptionType subscriptionType;
    private final EventType eventtype;
    private final S7Tag s7Tag;
    private final ArrayList<Integer> ackAlarms;
    private final AlarmType alarmQueryType;

    public S7SubscriptionTag(String address, S7SubscriptionType subscriptionType, EventType eventtype) {
        this.address = address;
        this.subscriptionType = subscriptionType;
        this.eventtype = eventtype;
        this.s7Tag = null;
        this.ackAlarms = null;
        this.alarmQueryType = null;
    }

    public S7SubscriptionTag(String address, S7SubscriptionType subscriptionType, ArrayList<Integer> ackAlarms) {
        this.address = address;
        this.subscriptionType = subscriptionType;
        this.eventtype = null;
        this.s7Tag = null;
        this.ackAlarms = ackAlarms;
        this.alarmQueryType = null;
    }

    public S7SubscriptionTag(String address, S7SubscriptionType subscriptionType, AlarmType alarmQueryType) {
        this.address = address;
        this.subscriptionType = subscriptionType;
        this.eventtype = null;
        this.s7Tag = null;
        this.ackAlarms = null;
        this.alarmQueryType = alarmQueryType;
    }

    public S7SubscriptionTag(String address, S7SubscriptionType subscriptionType, S7Tag s7Tag) {
        this.address = address;
        this.subscriptionType = subscriptionType;
        this.eventtype = null;
        this.s7Tag = s7Tag;
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

    public S7SubscriptionType getTagType() {
        return subscriptionType;
    }

    public EventType getEventType() {
        return eventtype;
    }

    public S7Tag getS7Tag() {
        return s7Tag;
    }

    public ArrayList<Integer> getAckAlarms() {
        return ackAlarms;
    }

    public AlarmType getAlarmQueryType() {
        return alarmQueryType;
    }

    public static boolean matches(String tagString) {
        return EVENT_SUBSCRIPTION_TYPE_PATTERN.matcher(tagString).matches() ||
            EVENT_ALARM_ACK_PATTERN.matcher(tagString).matches() ||
            EVENT_ALARM_QUERY_PATTERN.matcher(tagString).matches() ||
            S7Tag.matches(tagString);
    }

    public static S7SubscriptionTag of(String tagString) {
        {
            Matcher matcher = EVENT_SUBSCRIPTION_TYPE_PATTERN.matcher(tagString);
            if (matcher.matches()) {
                return new S7SubscriptionTag(tagString, S7SubscriptionType.EVENT_SUBSCRIPTION,
                    EventType.valueOf(tagString));
            }
        }

        {
            //TODO: Actually only ALARM_S (SIG_1)
            Matcher matcher = EVENT_ALARM_ACK_PATTERN.matcher(tagString);
            if (matcher.matches()) {
                String EventIds = matcher.group(2);
                String[] arrStrEventId = EventIds.split(",");
                ArrayList<Integer> arrEventId = new ArrayList<>();
                for (String EventId : arrStrEventId) {
                    EventId = EventId.replaceAll("16#", "");
                    arrEventId.add(Integer.parseInt(EventId, 16));
                }
                return new S7SubscriptionTag(tagString, S7SubscriptionType.ALARM_ACK, arrEventId);
            }
        }

        {
            //TODO: Support for ALARM_8            
            Matcher matcher = EVENT_ALARM_QUERY_PATTERN.matcher(tagString);
            if (matcher.matches()) {
                return new S7SubscriptionTag(tagString, S7SubscriptionType.ALARM_QUERY, AlarmType.ALARM_S);
            }
        }

        {
            if (S7Tag.matches(tagString)) {
                S7Tag s7Tag = S7Tag.of(tagString);
                switch (s7Tag.getDataType()) {
                    case BYTE:

                    default:

                }
                return new S7SubscriptionTag(tagString, S7SubscriptionType.CYCLIC_SUBSCRIPTION,
                    s7Tag);
            }
        }
        throw new PlcInvalidTagException("Unable to parse address: " + tagString);
    }

}
