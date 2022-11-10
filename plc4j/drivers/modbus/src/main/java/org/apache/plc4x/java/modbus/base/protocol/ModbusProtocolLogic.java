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
package org.apache.plc4x.java.modbus.base.protocol;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.modbus.base.tag.*;
import org.apache.plc4x.java.modbus.readwrite.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcList;

import java.time.Duration;
import java.util.Arrays;
import java.util.BitSet;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;

public abstract class ModbusProtocolLogic<T extends ModbusADU> extends Plc4xProtocolBase<T> {

    protected final DriverType driverType;
    protected Duration requestTimeout;
    protected short unitIdentifier;
    protected RequestTransactionManager tm;
    protected final AtomicInteger transactionIdentifierGenerator = new AtomicInteger(1);
    protected final static int FC_EXTENDED_REGISTERS_GROUP_HEADER_LENGTH = 2;
    protected final static int FC_EXTENDED_REGISTERS_FILE_RECORD_LENGTH = 10000;

    public ModbusProtocolLogic(DriverType driverType) {
        this.driverType = driverType;
    }

    @Override
    public void close(ConversationContext<T> context) {
        // Nothing to do here ...
    }

    @Override
    protected void decode(ConversationContext<T> context, T msg) throws Exception {
        super.decode(context, msg);
    }

    protected PlcResponseCode getErrorCode(ModbusPDUError errorResponse) {
        switch (errorResponse.getExceptionCode()) {
            case ILLEGAL_FUNCTION:
                return PlcResponseCode.UNSUPPORTED;
            case ILLEGAL_DATA_ADDRESS:
                return PlcResponseCode.INVALID_ADDRESS;
            case ILLEGAL_DATA_VALUE:
                return PlcResponseCode.INVALID_DATA;
            case SLAVE_DEVICE_FAILURE:
                return PlcResponseCode.REMOTE_ERROR;
            case ACKNOWLEDGE:
                return PlcResponseCode.OK;
            case SLAVE_DEVICE_BUSY:
                return PlcResponseCode.REMOTE_BUSY;
            case NEGATIVE_ACKNOWLEDGE:
                return PlcResponseCode.REMOTE_ERROR;
            case MEMORY_PARITY_ERROR:
                return PlcResponseCode.INTERNAL_ERROR;
            case GATEWAY_PATH_UNAVAILABLE:
                return PlcResponseCode.INTERNAL_ERROR;
            case GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND:
                return PlcResponseCode.REMOTE_ERROR;
            default:
                // This generally implies that something went wrong which we didn't anticipate.
                return PlcResponseCode.INTERNAL_ERROR;
        }
    }

    protected ModbusPDU getReadRequestPdu(PlcTag tag) {
        if (tag instanceof ModbusTagDiscreteInput) {
            ModbusTagDiscreteInput discreteInput = (ModbusTagDiscreteInput) tag;
            return new ModbusPDUReadDiscreteInputsRequest(discreteInput.getAddress(), discreteInput.getNumberOfElements());
        } else if (tag instanceof ModbusTagCoil) {
            ModbusTagCoil coil = (ModbusTagCoil) tag;
            return new ModbusPDUReadCoilsRequest(coil.getAddress(), coil.getNumberOfElements());
        } else if (tag instanceof ModbusTagInputRegister) {
            ModbusTagInputRegister inputRegister = (ModbusTagInputRegister) tag;
            return new ModbusPDUReadInputRegistersRequest(inputRegister.getAddress(), Math.max(inputRegister.getLengthWords(), 1));
        } else if (tag instanceof ModbusTagHoldingRegister) {
            ModbusTagHoldingRegister holdingRegister = (ModbusTagHoldingRegister) tag;
            return new ModbusPDUReadHoldingRegistersRequest(holdingRegister.getAddress(), Math.max(holdingRegister.getLengthWords(), 1));
        } else if (tag instanceof ModbusTagExtendedRegister) {
            ModbusTagExtendedRegister extendedRegister = (ModbusTagExtendedRegister) tag;
            int group1Address = extendedRegister.getAddress() % 10000;
            int group2Address = 0;
            int group1Quantity;
            int group2Quantity;
            short group1FileNumber = (short) (Math.floor((float) extendedRegister.getAddress() / 10000) + 1);
            short group2FileNumber;
            List<ModbusPDUReadFileRecordRequestItem> itemArray;

            if ((group1Address + extendedRegister.getLengthWords()) <= FC_EXTENDED_REGISTERS_FILE_RECORD_LENGTH) {
                // If request doesn't span file records, use a single group
                group1Quantity = extendedRegister.getLengthWords();
                ModbusPDUReadFileRecordRequestItem group1 =
                    new ModbusPDUReadFileRecordRequestItem((short) 6, group1FileNumber, group1Address, group1Quantity);
                itemArray = Collections.singletonList(group1);
            } else {
                // If it doesn't span a file record. e.g. 609998[10] request 2 words in first group and 8 in second.
                group1Quantity = FC_EXTENDED_REGISTERS_FILE_RECORD_LENGTH - group1Address;
                group2Quantity = extendedRegister.getLengthWords() - group1Quantity;
                group2FileNumber = (short) (group1FileNumber + 1);
                ModbusPDUReadFileRecordRequestItem group1 =
                    new ModbusPDUReadFileRecordRequestItem((short) 6, group1FileNumber, group1Address, group1Quantity);
                ModbusPDUReadFileRecordRequestItem group2 =
                    new ModbusPDUReadFileRecordRequestItem((short) 6, group2FileNumber, group2Address, group2Quantity);
                itemArray = Arrays.asList(group1, group2);
            }
            return new ModbusPDUReadFileRecordRequest(itemArray);
        }
        throw new PlcRuntimeException("Unsupported read tag type " + tag.getClass().getName());
    }

    protected ModbusPDU getWriteRequestPdu(PlcTag tag, PlcValue plcValue) {
        if (tag instanceof ModbusTagCoil) {
            ModbusTagCoil coil = (ModbusTagCoil) tag;
            ModbusPDUWriteMultipleCoilsRequest request =
                new ModbusPDUWriteMultipleCoilsRequest(coil.getAddress(), coil.getNumberOfElements(),
                    fromPlcValue(tag, plcValue));
            if (request.getQuantity() == coil.getNumberOfElements()) {
                return request;
            } else {
                throw new PlcRuntimeException("Number of requested bytes (" + request.getQuantity() +
                    ") doesn't match number of requested addresses (" + coil.getNumberOfElements() + ")");
            }
        } else if (tag instanceof ModbusTagHoldingRegister) {
            ModbusTagHoldingRegister holdingRegister = (ModbusTagHoldingRegister) tag;
            ModbusPDUWriteMultipleHoldingRegistersRequest request =
                new ModbusPDUWriteMultipleHoldingRegistersRequest(holdingRegister.getAddress(),
                    holdingRegister.getLengthWords(), fromPlcValue(tag, plcValue));
            if (request.getValue().length == holdingRegister.getLengthWords() * 2) {
                return request;
            } else {
                throw new PlcRuntimeException("Number of requested values (" + request.getValue().length / 2 +
                    ") doesn't match number of requested addresses (" + holdingRegister.getLengthWords() + ")");
            }
        } else if (tag instanceof ModbusTagExtendedRegister) {
            ModbusTagExtendedRegister extendedRegister = (ModbusTagExtendedRegister) tag;
            int group1Address = extendedRegister.getAddress() % FC_EXTENDED_REGISTERS_FILE_RECORD_LENGTH;
            int group2Address = 0;
            int group1Quantity;
            int group2Quantity;
            byte[] plcValue1, plcValue2;
            short group1FileNumber = (short)
                (Math.floor((float) extendedRegister.getAddress() / FC_EXTENDED_REGISTERS_FILE_RECORD_LENGTH) + 1);
            short group2FileNumber;
            List<ModbusPDUWriteFileRecordRequestItem> itemArray;
            if ((group1Address + extendedRegister.getLengthWords()) <= FC_EXTENDED_REGISTERS_FILE_RECORD_LENGTH) {
                //If request doesn't span file records, use a single group
                group1Quantity = extendedRegister.getLengthWords();
                ModbusPDUWriteFileRecordRequestItem group1 = new ModbusPDUWriteFileRecordRequestItem(
                    (short) 6, group1FileNumber, group1Address, fromPlcValue(tag, plcValue));
                itemArray = Collections.singletonList(group1);
            } else {
                //If it doesn't span a file record. e.g. 609998[10] request 2 words in first group and 8 in second.
                group1Quantity = FC_EXTENDED_REGISTERS_FILE_RECORD_LENGTH - group1Address;
                group2Quantity = extendedRegister.getLengthWords() - group1Quantity;
                group2FileNumber = (short) (group1FileNumber + 1);

                plcValue1 = ArrayUtils.subarray(fromPlcValue(tag, plcValue), 0, group1Quantity);
                plcValue2 = ArrayUtils.subarray(
                    fromPlcValue(tag, plcValue), group1Quantity, fromPlcValue(tag, plcValue).length);
                ModbusPDUWriteFileRecordRequestItem group1 = new ModbusPDUWriteFileRecordRequestItem(
                    (short) 6, group1FileNumber, group1Address, plcValue1);
                ModbusPDUWriteFileRecordRequestItem group2 = new ModbusPDUWriteFileRecordRequestItem(
                    (short) 6, group2FileNumber, group2Address, plcValue2);
                itemArray = Arrays.asList(group1, group2);
            }
            return new ModbusPDUWriteFileRecordRequest(itemArray);
        }
        throw new PlcRuntimeException("Unsupported write tag type " + tag.getClass().getName());
    }

    protected PlcValue toPlcValue(ModbusPDU request, ModbusPDU response, ModbusDataType dataType) throws ParseException {
        short tagDataTypeSize = dataType.getDataTypeSize();

        if (request instanceof ModbusPDUReadDiscreteInputsRequest) {
            if (!(response instanceof ModbusPDUReadDiscreteInputsResponse)) {
                throw new PlcRuntimeException("Unexpected response type. " +
                    "Expected ModbusPDUReadDiscreteInputsResponse, but got " + response.getClass().getName());
            }
            ModbusPDUReadDiscreteInputsRequest req = (ModbusPDUReadDiscreteInputsRequest) request;
            ModbusPDUReadDiscreteInputsResponse resp = (ModbusPDUReadDiscreteInputsResponse) response;
            return readCoilBooleanList(req.getQuantity(), resp.getValue());
        } else if (request instanceof ModbusPDUReadCoilsRequest) {
            if (!(response instanceof ModbusPDUReadCoilsResponse)) {
                throw new PlcRuntimeException("Unexpected response type. " +
                    "Expected ModbusPDUReadCoilsResponse, but got " + response.getClass().getName());
            }
            ModbusPDUReadCoilsRequest req = (ModbusPDUReadCoilsRequest) request;
            ModbusPDUReadCoilsResponse resp = (ModbusPDUReadCoilsResponse) response;
            return readCoilBooleanList(req.getQuantity(), resp.getValue());
        } else if (request instanceof ModbusPDUReadInputRegistersRequest) {
            if (!(response instanceof ModbusPDUReadInputRegistersResponse)) {
                throw new PlcRuntimeException("Unexpected response type. " +
                    "Expected ModbusPDUReadInputRegistersResponse, but got " + response.getClass().getName());
            }
            ModbusPDUReadInputRegistersRequest req = (ModbusPDUReadInputRegistersRequest) request;
            ModbusPDUReadInputRegistersResponse resp = (ModbusPDUReadInputRegistersResponse) response;
            ReadBuffer io = new ReadBufferByteBased(resp.getValue());
            if (tagDataTypeSize < 2) {
                io.readByte();
            }
            return DataItem.staticParse(io, dataType, Math.max(Math.round(req.getQuantity() / (tagDataTypeSize / 2.0f)), 1));
        } else if (request instanceof ModbusPDUReadHoldingRegistersRequest) {
            if (!(response instanceof ModbusPDUReadHoldingRegistersResponse)) {
                throw new PlcRuntimeException("Unexpected response type. " +
                    "Expected ModbusPDUReadHoldingRegistersResponse, but got " + response.getClass().getName());
            }
            ModbusPDUReadHoldingRegistersRequest req = (ModbusPDUReadHoldingRegistersRequest) request;
            ModbusPDUReadHoldingRegistersResponse resp = (ModbusPDUReadHoldingRegistersResponse) response;
            ReadBuffer io = new ReadBufferByteBased(resp.getValue());
            if ((dataType != ModbusDataType.STRING) && tagDataTypeSize < 2) {
                io.readByte();
            }
            return DataItem.staticParse(io, dataType, Math.max(Math.round(req.getQuantity() / (tagDataTypeSize / 2.0f)), 1));
        } else if (request instanceof ModbusPDUReadFileRecordRequest) {
            if (!(response instanceof ModbusPDUReadFileRecordResponse)) {
                throw new PlcRuntimeException("Unexpected response type. " +
                    "Expected ModbusPDUReadFileRecordResponse, but got " + response.getClass().getName());
            }
            ModbusPDUReadFileRecordRequest req = (ModbusPDUReadFileRecordRequest) request;
            ModbusPDUReadFileRecordResponse resp = (ModbusPDUReadFileRecordResponse) response;
            ReadBuffer io;
            short dataLength;

            if (resp.getItems().size() == 2 && resp.getItems().size() == req.getItems().size()) {
                //If request was split over file records, two groups in response should be received.
                io = new ReadBufferByteBased(ArrayUtils.addAll(resp.getItems().get(0).getData(), resp.getItems().get(1).getData()));
                dataLength = (short) (resp.getItems().get(0).getLengthInBytes() + resp.getItems().get(1).getLengthInBytes() - (2 * FC_EXTENDED_REGISTERS_GROUP_HEADER_LENGTH));
            } else if (resp.getItems().size() == 1 && resp.getItems().size() == req.getItems().size()) {
                //If request was within a single file record, one group should be received.
                io = new ReadBufferByteBased(resp.getItems().get(0).getData());
                dataLength = (short) (resp.getItems().get(0).getLengthInBytes() - FC_EXTENDED_REGISTERS_GROUP_HEADER_LENGTH);
            } else {
                throw new PlcRuntimeException("Unexpected number of groups in response. " +
                    "Expected " + req.getItems().size() + ", but got " + resp.getItems().size());
            }
            if (tagDataTypeSize < 2) {
                io.readByte();
            }
            return DataItem.staticParse(io, dataType, Math.round(Math.max(dataLength / 2.0f, 1) / Math.max(tagDataTypeSize / 2.0f, 1)));
        }
        return null;
    }

    protected byte[] fromPlcValue(PlcTag tag, PlcValue plcValue) {
        ModbusDataType tagDataType = ((ModbusTag) tag).getDataType();
        try {
            if (plcValue instanceof PlcList) {
                WriteBufferByteBased writeBuffer = new WriteBufferByteBased(DataItem.getLengthInBytes(plcValue, tagDataType, plcValue.getLength()));
                DataItem.staticSerialize(writeBuffer, plcValue, tagDataType, plcValue.getLength(), ByteOrder.BIG_ENDIAN);
                byte[] data = writeBuffer.getData();
                if (((ModbusTag) tag).getDataType() == ModbusDataType.BOOL) {
                    //Reverse Bits in each byte as
                    //they should be ordered like this: 8 7 6 5 4 3 2 1 | 0 0 0 0 0 0 0 9
                    byte[] bytes = new byte[data.length];
                    for (int i = 0; i < data.length; i++) {
                        bytes[i] = reverseBitsOfByte(data[i]);
                    }
                    return bytes;
                }
                return data;
            } else {
                WriteBufferByteBased writeBuffer = new WriteBufferByteBased(DataItem.getLengthInBytes(plcValue, tagDataType, plcValue.getLength()));
                DataItem.staticSerialize(writeBuffer, plcValue, tagDataType, plcValue.getLength(), ByteOrder.BIG_ENDIAN);
                return writeBuffer.getData();
            }
        } catch (SerializationException e) {
            throw new PlcRuntimeException("Unable to parse PlcValue :- " + e);
        }

    }

    protected byte reverseBitsOfByte(byte b) {
        BitSet bits = BitSet.valueOf(new byte[]{b});
        BitSet reverse = BitSet.valueOf(new byte[]{(byte) 0xFF});
        for (int j = 0; j < 8; j++) {
            reverse.set(j, bits.get(7 - j));
        }
        //toByteArray returns an empty array if all the bits are set to 0.
        return Arrays.copyOf(reverse.toByteArray(), 1)[0];
    }

    protected PlcValue readCoilBooleanList(int count, byte[] data) throws ParseException {
        ReadBuffer io = new ReadBufferByteBased(data);
        if (count == 1) {
            // Skip the first 7 bits.
            io.readInt(7);
            return new PlcBOOL(io.readBit());
        }

        int numFullBytes = count / 8;
        int numBitsIncompleteByte = count - (numFullBytes * 8);
        PlcValue[] values = new PlcValue[count];
        for (int i = 0; i < numFullBytes; i++) {
            values[(i*8)+7] = new PlcBOOL(io.readBit());
            values[(i*8)+6] = new PlcBOOL(io.readBit());
            values[(i*8)+5] = new PlcBOOL(io.readBit());
            values[(i*8)+4] = new PlcBOOL(io.readBit());
            values[(i*8)+3] = new PlcBOOL(io.readBit());
            values[(i*8)+2] = new PlcBOOL(io.readBit());
            values[(i*8)+1] = new PlcBOOL(io.readBit());
            values[(i*8)] = new PlcBOOL(io.readBit());
        }
        if(numBitsIncompleteByte > 0) {
            io.readInt(8 - numBitsIncompleteByte);
            for (int i = 1; i <= numBitsIncompleteByte; i++) {
                values[(numFullBytes*8)+(numBitsIncompleteByte - i)] = new PlcBOOL(io.readBit());
            }
        }
        return new PlcList(Arrays.asList(values));
    }

}
