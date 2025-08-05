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
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import static io.netty.buffer.Unpooled.wrappedBuffer;
import org.apache.plc4x.java.api.messages.PlcMetadataKeys;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.spi.metadata.DefaultMetadata;
import org.apache.plc4x.java.api.metadata.time.TimeSource;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesChangeDrivenPush;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesChangeDrivenSubscribeResponse;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesPush;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesSubscribeResponse;
import org.apache.plc4x.java.s7.readwrite.utils.StaticHelper;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.Charset;
import java.time.*;
import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.IntStream;
import org.apache.plc4x.java.s7.readwrite.tag.S7SubscriptionTag;
import org.apache.plc4x.java.s7.readwrite.tag.S7Tag;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;

public class S7CyclicEvent extends S7EventBase implements S7Event {

    public enum Fields {
        TYPE,
        JOBID,
        TIMESTAMP,
        ITEMSCOUNT,
        REQUEST,
        MAP,
        RETURNCODE_,
        TRANSPORTSIZE_,
        DATA_
    }
    private final PlcSubscriptionRequest request;
    private final Map<String, Object> map;
    public S7CyclicEvent(PlcSubscriptionRequest request, short jobid, S7PayloadUserDataItemCyclicServicesPush event) {
        super(Instant.now(), new DefaultMetadata.Builder()
            .put(PlcMetadataKeys.TIMESTAMP_SOURCE, TimeSource.ASSUMPTION)
            .build()
        );
        this.map = new HashMap<>();
        this.request = request;
        map.put(Fields.TYPE.name(), "CYCEVENT");
        map.put(Fields.TIMESTAMP.name(), getTimestamp());
        map.put(Fields.JOBID.name(), jobid);
        map.put(Fields.ITEMSCOUNT.name(), event.getItemsCount());
        int[] n = new int[1];

        request.getTagNames().forEach(tagname -> {
            int i = n[0];
            map.put(Fields.RETURNCODE_.name() + i, event.getItems().get(i).getReturnCode().getValue());
            map.put(Fields.TRANSPORTSIZE_.name() + i, event.getItems().get(i).getTransportSize().getValue());
            map.put(tagname, dataToPlcValue(tagname, request, event.getItems().get(i).getData()));
            n[0]++;
        });
    }
    public S7CyclicEvent(PlcSubscriptionRequest request, short jobid, S7PayloadUserDataItemCyclicServicesChangeDrivenPush event) {
        super(Instant.now(), new DefaultMetadata.Builder()
            .put(PlcMetadataKeys.TIMESTAMP_SOURCE, TimeSource.ASSUMPTION)
            .build()
        );
        this.map = new HashMap<>();
        this.request = request;
        map.put(Fields.TYPE.name(), "CYCEVENT");
        map.put(Fields.TIMESTAMP.name(), getTimestamp());
        map.put(Fields.JOBID.name(), jobid);
        map.put(Fields.ITEMSCOUNT.name(), event.getItemsCount());
        int[] n = new int[1];
       
        request.getTagNames().forEach(tagname -> {
            int i = n[0];
            map.put(Fields.RETURNCODE_.name() + i, event.getItems().get(i).getReturnCode().getValue());
            map.put(Fields.TRANSPORTSIZE_.name() + i, event.getItems().get(i).getTransportSize().getValue());
            map.put(tagname, dataToPlcValue(tagname, request, event.getItems().get(i).getData()));
            n[0]++;
        });
       
    }
    public S7CyclicEvent(PlcSubscriptionRequest request, short jobid, S7PayloadUserDataItemCyclicServicesSubscribeResponse event) {
        super(Instant.now(), new DefaultMetadata.Builder()
            .put(PlcMetadataKeys.TIMESTAMP_SOURCE, TimeSource.ASSUMPTION)
            .build()
        );
        this.map = new HashMap<>();
        this.request = request;
        map.put(Fields.TYPE.name(), "CYCEVENT");
        map.put(Fields.TIMESTAMP.name(), getTimestamp());
        map.put(Fields.JOBID.name(), jobid);
        map.put(Fields.ITEMSCOUNT.name(), event.getItemsCount());
        int[] n = new int[1];
        request.getTagNames().forEach(tagname -> {
            int i = n[0];
            map.put(Fields.RETURNCODE_.name() + i, event.getItems().get(i).getReturnCode().getValue());
            map.put(Fields.TRANSPORTSIZE_.name() + i, event.getItems().get(i).getTransportSize().getValue());
            map.put(tagname, dataToPlcValue(tagname, request, event.getItems().get(i).getData()));
            n[0]++;
        });
    }
    public S7CyclicEvent(PlcSubscriptionRequest request, short jobid, S7PayloadUserDataItemCyclicServicesChangeDrivenSubscribeResponse event) {
        super(Instant.now(), new DefaultMetadata.Builder()
            .put(PlcMetadataKeys.TIMESTAMP_SOURCE, TimeSource.ASSUMPTION)
            .build()
        );
        this.map = new HashMap<>();
        this.request = request;
        map.put(Fields.TYPE.name(), "CYCEVENT");
        map.put(Fields.TIMESTAMP.name(), getTimestamp());
        map.put(Fields.JOBID.name(), jobid);
        map.put(Fields.ITEMSCOUNT.name(), event.getItemsCount());
        int[] n = new int[1];
       
        request.getTagNames().forEach(tagname -> {
            int i = n[0];
            map.put(Fields.RETURNCODE_.name() + i, event.getItems().get(i).getReturnCode().getValue());
            map.put(Fields.TRANSPORTSIZE_.name() + i, event.getItems().get(i).getTransportSize().getValue());
            map.put(tagname, dataToPlcValue(tagname, request, event.getItems().get(i).getData()));
            n[0]++;
        });
    }
    @Override
    public Map<String, Object> getMap() {
        return this.map;
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
        if (request.getTagNames().contains(name)) {
            PlcValue plcvalue = (PlcValue) map.get(name);
            plcvalue.getRaw();
            return plcvalue;
        }
        return null;
    }
    @Override
    public int getNumberOfValues(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }
    @Override
    public Object getObject(String name) {
        if ("REQUEST".equals(name)) return request;
        return null;
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
        return isValidBoolean(name, 0);
    }
    @Override
    public boolean isValidBoolean(String name, int index) {
        try {
            getBoolean(name, index);
            return true;
        } catch (Exception ex) {
            return false;
        }
    }
    @Override
    public Boolean getBoolean(String name) {
        return getBoolean(name, 0);
    }
    @Override
    public Boolean getBoolean(String name, int index) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        return byteBuf.getBoolean(index);
    }
    @Override
    public Collection<Boolean> getAllBooleans(String name) {
        throw new UnsupportedOperationException("Not supported yet.");
    }
    @Override
    public boolean isValidByte(String name) {
        return isValidByte(name, 0);
    }
    @Override
    public boolean isValidByte(String name, int index) {
        try {
            getByte(name, index);
            return true;
        } catch (Exception ex) {
            return false;
        }
    }
    @Override
    public Byte getByte(String name) {
        return getByte(name, 0);
    }
    @Override
    public Byte getByte(String name, int index) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int pos = index * Byte.BYTES;
        return byteBuf.getByte(pos);
    }
    @Override
    public Collection<Byte> getAllBytes(String name) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        byte[] array = (byte[]) map.get(name);
        return IntStream.range(0, array.length).mapToObj(i -> array[i]).collect(Collectors.toList());
    }
    @Override
    public boolean isValidShort(String name) {
        return isValidShort(name, 0);
    }
    @Override
    public boolean isValidShort(String name, int index) {
        try {
            getShort(name, index);
            return true;
        } catch (Exception ex) {
            return false;
        }
    }
    @Override
    public Short getShort(String name) {
        return getShort(name, 0);
    }
    @Override
    public Short getShort(String name, int index) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int pos = index * Short.BYTES;
        return byteBuf.getShort(pos);
    }
    @Override
    public Collection<Short> getAllShorts(String name) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        List<Short> list = new ArrayList<>();
        while (byteBuf.isReadable(Short.BYTES)) {
            list.add(byteBuf.readShort());
        }
        return list;
    }
    @Override
    public boolean isValidInteger(String name) {
        return isValidInteger(name, 0);
    }
    @Override
    public boolean isValidInteger(String name, int index) {
        try {
            getInteger(name, index);
            return true;
        } catch (Exception ex) {
            return false;
        }
    }
    @Override
    public Integer getInteger(String name) {
        return getInteger(name, 0);
    }
    @Override
    public Integer getInteger(String name, int index) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int pos = index * Integer.BYTES;
        return byteBuf.getInt(pos);
    }
    @Override
    public Collection<Integer> getAllIntegers(String name) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        List<Integer> list = new ArrayList<>();
        while (byteBuf.isReadable(Integer.BYTES)) {
            list.add(byteBuf.readInt());
        }
        return list;
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
        return isValidLong(name, 0);
    }
    @Override
    public boolean isValidLong(String name, int index) {
        try {
            getLong(name, index);
            return true;
        } catch (Exception ex) {
            return false;
        }
    }
    @Override
    public Long getLong(String name) {
        return getLong(name, 0);
    }
    @Override
    public Long getLong(String name, int index) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int pos = index * Long.BYTES;
        return byteBuf.getLong(pos);
    }
    @Override
    public Collection<Long> getAllLongs(String name) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        List<Long> list = new ArrayList<>();
        while (byteBuf.isReadable(Long.BYTES)) {
            list.add(byteBuf.readLong());
        }
        return list;
    }
    @Override
    public boolean isValidFloat(String name) {
        return isValidFloat(name, 0);
    }
    @Override
    public boolean isValidFloat(String name, int index) {
        try {
            getFloat(name, index);
            return true;
        } catch (Exception ex) {
            return false;
        }
    }
    @Override
    public Float getFloat(String name) {
        return getFloat(name, 0);
    }
    @Override
    public Float getFloat(String name, int index) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int pos = index * Float.BYTES;
        return byteBuf.getFloat(pos);
    }
    @Override
    public Collection<Float> getAllFloats(String name) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        List<Float> list = new ArrayList<>();
        while (byteBuf.isReadable(Float.BYTES)) {
            list.add(byteBuf.readFloat());
        }
        return list;
    }
    @Override
    public boolean isValidDouble(String name) {
        return isValidDouble(name, 0);
    }
    @Override
    public boolean isValidDouble(String name, int index) {
        try {
            getDouble(name, index);
            return true;
        } catch (Exception ex) {
            return false;
        }
    }
    @Override
    public Double getDouble(String name) {
        return getDouble(name, 0);
    }
    @Override
    public Double getDouble(String name, int index) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int pos = index * Double.BYTES;
        return byteBuf.getDouble(pos);
    }
    @Override
    public Collection<Double> getAllDoubles(String name) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        List<Double> list = new ArrayList<>();
        while (byteBuf.isReadable(Double.BYTES)) {
            list.add(byteBuf.readDouble());
        }
        return list;
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
        return isValidString(name, 0);
    }
    @Override
    public boolean isValidString(String name, int index) {
        try {
            getString(name, index);
            return true;
        } catch (Exception ex) {
            return false;
        }
    }
    @Override
    public String getString(String name) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        return byteBuf.toString(Charset.defaultCharset());
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
        return isValidTime(name, 0);
    }
    @Override
    public boolean isValidTime(String name, int index) {
        try {
            getTime(name, index);
            return true;
        } catch (Exception ex) {
            return false;
        }
    }
    @Override
    public LocalTime getTime(String name) {
        return getTime(name, 0);
    }
    /*
     * In S7, data type TIME occupies one double word.
     * The value is in milliseconds (ms).
     */
    @Override
    public LocalTime getTime(String name, int index) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int pos = index * Integer.BYTES;
        int value = byteBuf.getInt(pos);
        Duration dr = StaticHelper.s7TimeToDuration(value);
        return LocalTime.of(dr.toHoursPart(), dr.toMinutesPart(), dr.toSecondsPart(), dr.toNanosPart());
    }
    @Override
    public Collection<LocalTime> getAllTimes(String name) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int nitems = (byteBuf.capacity() / Integer.BYTES);
        List<LocalTime> items = new ArrayList<>();
        for (int i = 0; i < nitems; i++) {
            items.add(getTime(name, i));
        }
        return items;
    }
    @Override
    public boolean isValidDate(String name) {
        return isValidDate(name, 0);
    }
    @Override
    public boolean isValidDate(String name, int index) {
        try {
            getDate(name, index);
            return true;
        } catch (Exception ex) {
            return false;
        }
    }
    @Override
    public LocalDate getDate(String name) {
        return getDate(name, 0);
    }
    @Override
    public LocalDate getDate(String name, int index) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int pos = index * Short.BYTES;
        short value = byteBuf.getShort(pos);
        return StaticHelper.s7DateToLocalDate(value);
    }
    @Override
    public Collection<LocalDate> getAllDates(String name) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int nItems = (byteBuf.capacity() / Short.BYTES);
        List<LocalDate> items = new ArrayList<>();
        for (int i = 0; i < nItems; i++) {
            items.add(getDate(name, i));
        }
        return items;
    }
    @Override
    public boolean isValidDateTime(String name) {
        return isValidDateTime(name, 0);
    }
    @Override
    public boolean isValidDateTime(String name, int index) {
        try {
            getDateTime(name, index);
            return true;
        } catch (Exception ex) {
            return false;
        }
    }
    @Override
    public LocalDateTime getDateTime(String name) {
        return getDateTime(name, 0);
    }
    @Override
    public LocalDateTime getDateTime(String name, int index) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int pos = index * Long.BYTES;
        return StaticHelper.s7DateTimeToLocalDateTime(byteBuf.slice(pos, Long.BYTES));
    }
    @Override
    public Collection<LocalDateTime> getAllDateTimes(String name) {
        if (!(map.get(name) instanceof byte[])) {
            throw new UnsupportedOperationException("Field is not a buffer of bytes. Required byte[] type.");
        }
        ByteBuf byteBuf = Unpooled.wrappedBuffer((byte[]) map.get(name));
        int nitems = (byteBuf.capacity() / Long.BYTES);
        List<LocalDateTime> items = new ArrayList<>();
        for (int i = 0; i < nitems; i++) {
            items.add(getDateTime(name, i));
        }
        return items;
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
    @Override
    public boolean equals(Object obj) {
        if (this == obj) {
            return true;
        }
        if (obj == null) {
            return false;
        }
        if (getClass() != obj.getClass()) {
            return false;
        }
        final S7CyclicEvent other = (S7CyclicEvent) obj;
       
        for (String tag:request.getTagNames()) {
            final PlcValue othervalue = other.getPlcValue(tag);
            if (othervalue == null) return false;
            final PlcValue localvalue = (PlcValue) getPlcValue(tag);
            if (Arrays.equals(localvalue.getRaw(), othervalue.getRaw()) == false){
                return false;
            }
        };
       
        return true;
    }
    
        
    private static PlcValue dataToPlcValue(String tagname, PlcSubscriptionRequest request, List<Short> data){
        
        int[] i = new int[1];
       
        final byte[] buffer = new byte[data.size()];
        data.forEach( b -> {
            buffer[i[0]] = b.byteValue();
            i[0]++;
        });
               
        ByteBuf bb = wrappedBuffer(buffer);
       
       
        final DefaultPlcSubscriptionTag  dpst = (DefaultPlcSubscriptionTag) request.getTag(tagname);
        final S7SubscriptionTag subTag = (S7SubscriptionTag) dpst.getTag();
        final S7Tag[] s7Tags = subTag.getS7Tags();
       
        PlcValue plcValue = null;
       
        switch(s7Tags[0].getDataType()){
            case BOOL:
                Boolean[] bools = new Boolean[s7Tags[0].getNumberOfElements()];
                for (int iter = 0; iter < s7Tags[0].getNumberOfElements(); iter++) {
                    bools[iter] = bb.readBoolean();
                }
                plcValue = DefaultPlcValueHandler.of(s7Tags[0], bools);
                break;
            case BYTE:
                // TODO: This looks suspicious
                Byte[] bytes = new Byte[bb.capacity()];

                for (byte b:bb.array()) {
                    b = Byte.valueOf(bb.readByte());
                }
                plcValue = DefaultPlcValueHandler.of(s7Tags[0], bytes);
                break;
            case WORD:
                break;
            case DWORD:
                break;
            case LWORD:
                break;
            case INT:
                Short[] shorts = new Short[s7Tags[0].getNumberOfElements()];
                for (int iter = 0; iter < s7Tags[0].getNumberOfElements(); iter ++) {
                    shorts[iter] = bb.readShort();
                }
                plcValue = DefaultPlcValueHandler.of(s7Tags[0], shorts);
                break;
            case UINT:
                break;
            case SINT:
                break;
            case USINT:
                break;
            case DINT:
                // TODO: This looks suspicious
                Integer[] integers = new Integer[bb.capacity() / Integer.SIZE];
                for (Integer di:integers) {
                    di = Integer.valueOf(bb.readInt());
                }
                plcValue = DefaultPlcValueHandler.of(s7Tags[0], integers);
                break;
            case UDINT:
                break;
            case LINT:
                // TODO: This looks suspicious
                Long[] longs = new Long[bb.capacity() / Long.SIZE];
                for (Long l:longs) {
                    l = bb.readLong();
                }
                plcValue = DefaultPlcValueHandler.of(s7Tags[0], longs);
                break;
            case ULINT:
                break;
            case REAL:
                // TODO: This looks suspicious
                Float[] floats = new Float[bb.capacity() / Float.SIZE];
                for (Float f:floats) {
                    f = bb.readFloat();
                }
                plcValue = DefaultPlcValueHandler.of(s7Tags[0], floats);
                break;
            case LREAL:
                // TODO: This looks suspicious
                Double[] doubles = new Double[bb.capacity() / Double.SIZE];
                for (Double d:doubles) {
                    d = bb.readDouble();
                }
                plcValue = DefaultPlcValueHandler.of(s7Tags[0], doubles);
                break;
            case CHAR:
                break;
            case WCHAR:
                break;
            case STRING:
                break;
            case WSTRING:
                break;
            case S5TIME:
                break;
            case TIME:
                break;
            case LTIME:
                break;
            case DATE:
                break;
            case TIME_OF_DAY:
                break;
            case TOD:
                break;
            case LTIME_OF_DAY:
                break;
            case LTOD:
                break;
            case DATE_AND_TIME:
                break;
            case DT:
                break;
            case DATE_AND_LTIME:
                break;
            case LDT:
                break;
            case DTL:
                break;
        }
       
        return plcValue;
    }
}
