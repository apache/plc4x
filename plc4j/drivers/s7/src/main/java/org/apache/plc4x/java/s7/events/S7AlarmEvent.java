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
package org.apache.plc4x.java.s7.events;

import java.util.Collections;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.readwrite.*;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.*;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class S7AlarmEvent extends S7EventBase {


    public enum Fields {
        MAP,

        TYPE,
        ASSOCIATED_VALUES,
        TIMESTAMP,
        TIMESTAMP_GOING,
        TIMESTAMP_COMING,

        EVENT_ID,
        EVENT_STATE,
        STATE,
        ACKSTATE_GOING,
        ACKSTATE_COMING,
        EVENT_GOING,
        EVENT_COMING,
        EVENT_LAST_CHANGE,

        SIG,
        SIG_1,
        SIG_2,
        SIG_3,
        SIG_4,
        SIG_5,
        SIG_6,
        SIG_7,
        SIG_8,

        SIG_STATE,
        SIG_1_STATE,
        SIG_2_STATE,
        SIG_3_STATE,
        SIG_4_STATE,
        SIG_5_STATE,
        SIG_6_STATE,
        SIG_7_STATE,
        SIG_8_STATE,

        SIG_DATA,
        SIG_1_DATA,
        SIG_2_DATA,
        SIG_3_DATA,
        SIG_4_DATA,
        SIG_5_DATA,
        SIG_6_DATA,
        SIG_7_DATA,
        SIG_8_DATA,

        SIG_1_DATA_GOING,
        SIG_2_DATA_GOING,
        SIG_3_DATA_GOING,
        SIG_4_DATA_GOING,
        SIG_5_DATA_GOING,
        SIG_6_DATA_GOING,
        SIG_7_DATA_GOING,
        SIG_8_DATA_GOING,
        SIG_1_DATA_COMING,
        SIG_2_DATA_COMING,
        SIG_3_DATA_COMING,
        SIG_4_DATA_COMING,
        SIG_5_DATA_COMING,
        SIG_6_DATA_COMING,
        SIG_7_DATA_COMING,
        SIG_8_DATA_COMING,

        SIG_1_DATA_STATUS,
        SIG_2_DATA_STATUS,
        SIG_3_DATA_STATUS,
        SIG_4_DATA_STATUS,
        SIG_5_DATA_STATUS,
        SIG_6_DATA_STATUS,
        SIG_7_DATA_STATUS,
        SIG_8_DATA_STATUS,

        SIG_1_DATA_SIZE,
        SIG_2_DATA_SIZE,
        SIG_3_DATA_SIZE,
        SIG_4_DATA_SIZE,
        SIG_5_DATA_SIZE,
        SIG_6_DATA_SIZE,
        SIG_7_DATA_SIZE,
        SIG_8_DATA_SIZE,

        SIG_1_DATA_LENGTH,
        SIG_2_DATA_LENGTH,
        SIG_3_DATA_LENGTH,
        SIG_4_DATA_LENGTH,
        SIG_5_DATA_LENGTH,
        SIG_6_DATA_LENGTH,
        SIG_7_DATA_LENGTH,
        SIG_8_DATA_LENGTH,

    }

    private final Map<String, Object> map;

    S7AlarmEvent(Instant timestamp, Map<String, Object> obj) {
        super(timestamp);
        this.map = obj;
    }


    @Override
    public Map<String, Object> getMap() {
        return map;
    }

    @Override
    public PlcReadRequest getRequest() {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public PlcValue getAsPlcValue() {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public PlcValue getPlcValue(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public int getNumberOfValues(String name) {
        return (int) map.get(Fields.ASSOCIATED_VALUES.name());
    }

    @Override
    public Object getObject(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Object getObject(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<Object> getAllObjects(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidBoolean(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidBoolean(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Boolean getBoolean(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Boolean getBoolean(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<Boolean> getAllBooleans(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidByte(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidByte(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Byte getByte(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Byte getByte(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<Byte> getAllBytes(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidShort(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidShort(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Short getShort(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Short getShort(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<Short> getAllShorts(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidInteger(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidInteger(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Integer getInteger(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Integer getInteger(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<Integer> getAllIntegers(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidBigInteger(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidBigInteger(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public BigInteger getBigInteger(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public BigInteger getBigInteger(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<BigInteger> getAllBigIntegers(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidLong(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidLong(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Long getLong(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Long getLong(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<Long> getAllLongs(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidFloat(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidFloat(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Float getFloat(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Float getFloat(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<Float> getAllFloats(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidDouble(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidDouble(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Double getDouble(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Double getDouble(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<Double> getAllDoubles(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidBigDecimal(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidBigDecimal(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public BigDecimal getBigDecimal(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public BigDecimal getBigDecimal(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<BigDecimal> getAllBigDecimals(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidString(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidString(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public String getString(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public String getString(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<String> getAllStrings(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidTime(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidTime(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public LocalTime getTime(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public LocalTime getTime(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<LocalTime> getAllTimes(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidDate(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidDate(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public LocalDate getDate(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public LocalDate getDate(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<LocalDate> getAllDates(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidDateTime(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public boolean isValidDateTime(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public LocalDateTime getDateTime(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public LocalDateTime getDateTime(String name, int index) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<LocalDateTime> getAllDateTimes(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public Collection<String> getTagNames() {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public PlcTag getTag(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public PlcResponseCode getResponseCode(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    public static S7AlarmEvent of(Object obj) {
        if (obj instanceof S7PayloadAlarmAckInd) {
            AlarmMessageAckPushType msg = ((S7PayloadAlarmAckInd) obj).getAlarmMessage();
            DateAndTime dt = msg.getTimeStamp();
            int year = (dt.getYear() >= 90) ? dt.getYear() + 1900 : dt.getYear() + 2000;
            LocalDateTime ldt = LocalDateTime.of(year,
                dt.getMonth(),
                dt.getDay(),
                dt.getHour(),
                dt.getMinutes(),
                dt.getSeconds(),
                dt.getMsec() * 1000000);
            Instant timeStamp = ldt.toInstant(ZoneOffset.UTC);

            Map<String, Object> map = new HashMap<>();
            map.put(S7SysEvent.Fields.TIMESTAMP.name(), timeStamp);

            List<AlarmMessageAckObjectPushType> items = msg.getMessageObjects();
            for (AlarmMessageAckObjectPushType item : items) {
                map.put(Fields.EVENT_ID.name(), item.getEventId());
                map.put(Fields.TYPE.name(), "ALARMACK_IND");
                map.put(Fields.ASSOCIATED_VALUES.name(), item.getNumberOfValues());

                map.put(Fields.SIG_1_DATA_GOING.name(), item.getAckStateGoing().getSIG_1());
                map.put(Fields.SIG_2_DATA_GOING.name(), item.getAckStateGoing().getSIG_2());
                map.put(Fields.SIG_3_DATA_GOING.name(), item.getAckStateGoing().getSIG_3());
                map.put(Fields.SIG_4_DATA_GOING.name(), item.getAckStateGoing().getSIG_4());
                map.put(Fields.SIG_5_DATA_GOING.name(), item.getAckStateGoing().getSIG_5());
                map.put(Fields.SIG_6_DATA_GOING.name(), item.getAckStateGoing().getSIG_6());
                map.put(Fields.SIG_7_DATA_GOING.name(), item.getAckStateGoing().getSIG_7());
                map.put(Fields.SIG_8_DATA_GOING.name(), item.getAckStateGoing().getSIG_8());

                map.put(Fields.SIG_1_DATA_COMING.name(), item.getAckStateComing().getSIG_1());
                map.put(Fields.SIG_2_DATA_COMING.name(), item.getAckStateComing().getSIG_2());
                map.put(Fields.SIG_3_DATA_COMING.name(), item.getAckStateComing().getSIG_3());
                map.put(Fields.SIG_4_DATA_COMING.name(), item.getAckStateComing().getSIG_4());
                map.put(Fields.SIG_5_DATA_COMING.name(), item.getAckStateComing().getSIG_5());
                map.put(Fields.SIG_6_DATA_COMING.name(), item.getAckStateComing().getSIG_6());
                map.put(Fields.SIG_7_DATA_COMING.name(), item.getAckStateComing().getSIG_7());
                map.put(Fields.SIG_8_DATA_COMING.name(), item.getAckStateComing().getSIG_8());
            }
            return new S7AlarmEvent(timeStamp, map);
        } else {

            AlarmMessagePushType msg = null;

            if (obj instanceof S7PayloadAlarm8) {
                msg = ((S7PayloadAlarm8) obj).getAlarmMessage();
            } else if (obj instanceof S7PayloadNotify) {
                msg = ((S7PayloadNotify) obj).getAlarmMessage();
            } else if (obj instanceof S7PayloadAlarmSQ) {
                msg = ((S7PayloadAlarmSQ) obj).getAlarmMessage();
            } else if (obj instanceof S7PayloadAlarmS) {
                msg = ((S7PayloadAlarmS) obj).getAlarmMessage();
            } else if (obj instanceof S7PayloadNotify8) {
                msg = ((S7PayloadNotify8) obj).getAlarmMessage();
            } else {
                throw new PlcRuntimeException("Unsupported type: " + obj.getClass().getName());
            }

            DateAndTime dt = msg.getTimeStamp();
            int year = (dt.getYear() >= 90) ? dt.getYear() + 1900 : dt.getYear() + 2000;
            LocalDateTime ldt = LocalDateTime.of(year,
                dt.getMonth(),
                dt.getDay(),
                dt.getHour(),
                dt.getMinutes(),
                dt.getSeconds(),
                dt.getMsec() * 1000000);
            Instant timeStamp = ldt.toInstant(ZoneOffset.UTC);

            Map<String, Object> map = new HashMap<>();
            map.put(S7SysEvent.Fields.TIMESTAMP.name(), timeStamp);

            List<AlarmMessageObjectPushType> items = msg.getMessageObjects();
            for (AlarmMessageObjectPushType item : items) {
                map.put(Fields.EVENT_ID.name(), item.getEventId());

                if (obj instanceof S7PayloadAlarm8) {
                    map.put(Fields.TYPE.name(), "ALARM8");
                }
                if (obj instanceof S7PayloadNotify) {
                    map.put(Fields.TYPE.name(), "NOTIFY");
                }
                if (obj instanceof S7PayloadAlarmSQ) {
                    map.put(Fields.TYPE.name(), "ALARMSQ");
                }
                if (obj instanceof S7PayloadAlarmS) {
                    map.put(Fields.TYPE.name(), "ALARMS");
                }
                if (obj instanceof S7PayloadNotify8) {
                    map.put(Fields.TYPE.name(), "NOTIFY8");
                }

                map.put(Fields.ASSOCIATED_VALUES.name(), item.getNumberOfValues());

                map.put(Fields.SIG_1.name(), item.getEventState().getSIG_1());
                map.put(Fields.SIG_2.name(), item.getEventState().getSIG_2());
                map.put(Fields.SIG_3.name(), item.getEventState().getSIG_3());
                map.put(Fields.SIG_4.name(), item.getEventState().getSIG_4());
                map.put(Fields.SIG_5.name(), item.getEventState().getSIG_5());
                map.put(Fields.SIG_6.name(), item.getEventState().getSIG_6());
                map.put(Fields.SIG_7.name(), item.getEventState().getSIG_7());
                map.put(Fields.SIG_8.name(), item.getEventState().getSIG_8());


                map.put(Fields.SIG_1_STATE.name(), item.getLocalState().getSIG_1());
                map.put(Fields.SIG_2_STATE.name(), item.getLocalState().getSIG_2());
                map.put(Fields.SIG_3_STATE.name(), item.getLocalState().getSIG_3());
                map.put(Fields.SIG_4_STATE.name(), item.getLocalState().getSIG_4());
                map.put(Fields.SIG_5_STATE.name(), item.getLocalState().getSIG_5());
                map.put(Fields.SIG_6_STATE.name(), item.getLocalState().getSIG_6());
                map.put(Fields.SIG_7_STATE.name(), item.getLocalState().getSIG_7());
                map.put(Fields.SIG_8_STATE.name(), item.getLocalState().getSIG_8());

                map.put(Fields.SIG_1_DATA_GOING.name(), item.getAckStateGoing().getSIG_1());
                map.put(Fields.SIG_2_DATA_GOING.name(), item.getAckStateGoing().getSIG_2());
                map.put(Fields.SIG_3_DATA_GOING.name(), item.getAckStateGoing().getSIG_3());
                map.put(Fields.SIG_4_DATA_GOING.name(), item.getAckStateGoing().getSIG_4());
                map.put(Fields.SIG_5_DATA_GOING.name(), item.getAckStateGoing().getSIG_5());
                map.put(Fields.SIG_6_DATA_GOING.name(), item.getAckStateGoing().getSIG_6());
                map.put(Fields.SIG_7_DATA_GOING.name(), item.getAckStateGoing().getSIG_7());
                map.put(Fields.SIG_8_DATA_GOING.name(), item.getAckStateGoing().getSIG_8());

                map.put(Fields.SIG_1_DATA_COMING.name(), item.getAckStateComing().getSIG_1());
                map.put(Fields.SIG_2_DATA_COMING.name(), item.getAckStateComing().getSIG_2());
                map.put(Fields.SIG_3_DATA_COMING.name(), item.getAckStateComing().getSIG_3());
                map.put(Fields.SIG_4_DATA_COMING.name(), item.getAckStateComing().getSIG_4());
                map.put(Fields.SIG_5_DATA_COMING.name(), item.getAckStateComing().getSIG_5());
                map.put(Fields.SIG_6_DATA_COMING.name(), item.getAckStateComing().getSIG_6());
                map.put(Fields.SIG_7_DATA_COMING.name(), item.getAckStateComing().getSIG_7());
                map.put(Fields.SIG_8_DATA_COMING.name(), item.getAckStateComing().getSIG_8());

                List<AssociatedValueType> values = item.getAssociatedValues();
                int i = 1;
                int j = 0;
                for (AssociatedValueType value : values) {
                    map.put("SIG_" + i + "_DATA_STATUS", value.getReturnCode().getValue());
                    map.put("SIG_" + i + "_DATA_SIZE", value.getTransportSize().getValue());
                    map.put("SIG_" + i + "_DATA_LENGTH", value.getValueLength());
                    byte[] data = new byte[value.getData().size()];
                    j = 0;
                    for (short s : value.getData()) {
                        data[j] = (byte) s;
                        j++;
                    }
                    map.put("SIG_" + i + "_DATA", data);
                    i++;
                }
            }

            return new S7AlarmEvent(timeStamp, map);
        }

    }

}
