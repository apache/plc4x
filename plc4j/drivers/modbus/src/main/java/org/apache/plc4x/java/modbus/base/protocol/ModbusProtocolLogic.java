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
import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcRawByteArray;

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
    protected PlcTag pingAddress;
    protected ModbusByteOrder defaultPayloadByteOrder;

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

    protected short getUnitId(PlcTag tag) {
        if (tag instanceof ModbusTag) {
            Short unitId = ((ModbusTag) tag).getUnitId();
            return unitId == null ? unitIdentifier : unitId;
        }

        return unitIdentifier;
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
            ModbusByteOrder byteOrder = defaultPayloadByteOrder;
            if(coil.getByteOrder() != null) {
                byteOrder = coil.getByteOrder();
            }
            ModbusPDUWriteMultipleCoilsRequest request =
                new ModbusPDUWriteMultipleCoilsRequest(coil.getAddress(), coil.getNumberOfElements(),
                    fromPlcValue(tag, plcValue, byteOrder));
            if (request.getQuantity() == coil.getNumberOfElements()) {
                return request;
            } else {
                throw new PlcRuntimeException("Number of requested bytes (" + request.getQuantity() +
                    ") doesn't match number of requested addresses (" + coil.getNumberOfElements() + ")");
            }
        } else if (tag instanceof ModbusTagHoldingRegister) {
            ModbusTagHoldingRegister holdingRegister = (ModbusTagHoldingRegister) tag;
            ModbusByteOrder byteOrder = defaultPayloadByteOrder;
            if(holdingRegister.getByteOrder() != null) {
                byteOrder = holdingRegister.getByteOrder();
            }
            byte[] bytes = fromPlcValue(tag, plcValue, byteOrder);
            if(bytes.length == 2) {
                int value = ((bytes[0] & 0xFF) << 8) | (bytes[1] & 0xFF);
                return new ModbusPDUWriteSingleRegisterRequest(holdingRegister.getAddress(), value);
            } else {
                ModbusPDUWriteMultipleHoldingRegistersRequest request =
                    new ModbusPDUWriteMultipleHoldingRegistersRequest(holdingRegister.getAddress(),
                        holdingRegister.getLengthWords(), bytes);
                if (request.getValue().length == holdingRegister.getLengthWords() * 2) {
                    return request;
                } else {
                    throw new PlcRuntimeException("Number of requested values (" + request.getValue().length / 2 +
                        ") doesn't match number of requested addresses (" + holdingRegister.getLengthWords() + ")");
                }
            }
        } else if (tag instanceof ModbusTagExtendedRegister) {
            ModbusTagExtendedRegister extendedRegister = (ModbusTagExtendedRegister) tag;
            ModbusByteOrder byteOrder = defaultPayloadByteOrder;
            if(extendedRegister.getByteOrder() != null) {
                byteOrder = extendedRegister.getByteOrder();
            }
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
                    (short) 6, group1FileNumber, group1Address, fromPlcValue(tag, plcValue, byteOrder));
                itemArray = Collections.singletonList(group1);
            } else {
                //If it doesn't span a file record. e.g. 609998[10] request 2 words in first group and 8 in second.
                group1Quantity = FC_EXTENDED_REGISTERS_FILE_RECORD_LENGTH - group1Address;
                group2Quantity = extendedRegister.getLengthWords() - group1Quantity;
                group2FileNumber = (short) (group1FileNumber + 1);

                plcValue1 = ArrayUtils.subarray(fromPlcValue(tag, plcValue, byteOrder), 0, group1Quantity);
                plcValue2 = ArrayUtils.subarray(
                    fromPlcValue(tag, plcValue, byteOrder), group1Quantity, fromPlcValue(tag, plcValue, byteOrder).length);
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

    protected PlcValue toPlcValue(ModbusPDU request, ModbusPDU response, ModbusDataType dataType, ModbusByteOrder byteOrder) throws ParseException {
        if (request instanceof ModbusPDUReadDiscreteInputsRequest) {
            if (!(response instanceof ModbusPDUReadDiscreteInputsResponse)) {
                throw new PlcRuntimeException("Unexpected response type. " +
                    "Expected ModbusPDUReadDiscreteInputsResponse, but got " + response.getClass().getName());
            }
            ModbusPDUReadDiscreteInputsResponse resp = (ModbusPDUReadDiscreteInputsResponse) response;
            return new PlcRawByteArray(resp.getValue());
        } else if (request instanceof ModbusPDUReadCoilsRequest) {
            if (!(response instanceof ModbusPDUReadCoilsResponse)) {
                throw new PlcRuntimeException("Unexpected response type. " +
                    "Expected ModbusPDUReadCoilsResponse, but got " + response.getClass().getName());
            }
            ModbusPDUReadCoilsResponse resp = (ModbusPDUReadCoilsResponse) response;
            return new PlcRawByteArray(resp.getValue());
        } else if (request instanceof ModbusPDUReadInputRegistersRequest) {
            if (!(response instanceof ModbusPDUReadInputRegistersResponse)) {
                throw new PlcRuntimeException("Unexpected response type. " +
                    "Expected ModbusPDUReadInputRegistersResponse, but got " + response.getClass().getName());
            }
            ModbusPDUReadInputRegistersResponse resp = (ModbusPDUReadInputRegistersResponse) response;
            return new PlcRawByteArray(resp.getValue());
        } else if (request instanceof ModbusPDUReadHoldingRegistersRequest) {
            if (!(response instanceof ModbusPDUReadHoldingRegistersResponse)) {
                throw new PlcRuntimeException("Unexpected response type. " +
                    "Expected ModbusPDUReadHoldingRegistersResponse, but got " + response.getClass().getName());
            }
            ModbusPDUReadHoldingRegistersResponse resp = (ModbusPDUReadHoldingRegistersResponse) response;
            return new PlcRawByteArray(resp.getValue());
        } else if (request instanceof ModbusPDUReadFileRecordRequest) {
            if (!(response instanceof ModbusPDUReadFileRecordResponse)) {
                throw new PlcRuntimeException("Unexpected response type. " +
                    "Expected ModbusPDUReadFileRecordResponse, but got " + response.getClass().getName());
            }
            ModbusPDUReadFileRecordResponse resp = (ModbusPDUReadFileRecordResponse) response;
            // TODO: This is an over-simplification ...
            return new PlcRawByteArray(resp.getItems().get(0).getData());
        }
        return null;
    }

    protected byte[] fromPlcValue(PlcTag tag, PlcValue plcValue, ModbusByteOrder byteOrder) {
        ModbusDataType tagDataType = ((ModbusTag) tag).getDataType();
        try {
            if (tag instanceof ModbusTagCoil) {
                // If it's a single value, just cast that to a number.
                if(plcValue instanceof PlcBOOL) {
                    byte byteValue = (byte) (plcValue.getBoolean() ? 1 : 0);
                    return new byte[]{byteValue};
                }
                // If it's a List, convert the booleans in the list into an array of bytes.
                else if(plcValue instanceof PlcList) {
                    PlcList valueList = (PlcList) plcValue;
                    WriteBufferByteBased wb = getWriteBuffer(((plcValue.getLength() - 1) / 8) + 1, byteOrder);
                    int paddingBits = 8 - (plcValue.getLength() % 8);
                    if(paddingBits < 8) {
                        for(int i = 0; i < paddingBits; i++) {
                            wb.writeBit(false);
                        }
                    }
                    for(int i = 0; i < plcValue.getLength(); i++) {
                        // We need to serialize the bits in reverse order for them to end in the right coils.
                        PlcValue value = valueList.getIndex((plcValue.getLength() - 1) - i);
                        if(!(value instanceof PlcBOOL)) {
                            throw new PlcRuntimeException("Expecting only BOOL values when writing coils.");
                        }
                        PlcBOOL boolValue = (PlcBOOL) value;
                        wb.writeBit(boolValue.getBoolean());
                    }

                    // Do the byte-swapping, if needed.
                    byte[] bytes = wb.getBytes();
                    if(byteOrder == ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP || byteOrder == ModbusByteOrder.LITTLE_ENDIAN_BYTE_SWAP) {
                        bytes = byteSwap(bytes);
                    }
                    // Reverse the bytes to have the "unfinished bytes" at the end.
                    ArrayUtils.reverse(bytes);
                    return bytes;
                }
                else {
                    throw new PlcRuntimeException("Expecting only BOOL or List values when writing coils.");
               }
            }
            else if (plcValue instanceof PlcList) {
                WriteBufferByteBased writeBuffer = getWriteBuffer(DataItem.getLengthInBytes(plcValue, tagDataType, plcValue.getLength(), byteOrder == ModbusByteOrder.BIG_ENDIAN), byteOrder);
                DataItem.staticSerialize(writeBuffer, plcValue, tagDataType, plcValue.getLength(), byteOrder == ModbusByteOrder.BIG_ENDIAN);
                byte[] data = writeBuffer.getBytes();
                if(byteOrder == ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP || byteOrder == ModbusByteOrder.LITTLE_ENDIAN_BYTE_SWAP) {
                    data = byteSwap(data);
                }
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
                WriteBufferByteBased writeBuffer = getWriteBuffer(DataItem.getLengthInBytes(plcValue, tagDataType, plcValue.getLength(), byteOrder == ModbusByteOrder.BIG_ENDIAN), byteOrder);
                DataItem.staticSerialize(writeBuffer, plcValue, tagDataType, plcValue.getLength(), byteOrder == ModbusByteOrder.BIG_ENDIAN);
                byte[] bytes = writeBuffer.getBytes();
                if(byteOrder == ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP || byteOrder == ModbusByteOrder.LITTLE_ENDIAN_BYTE_SWAP) {
                    bytes = byteSwap(bytes);
                }
                return bytes;
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

    /*protected PlcValue readCoilBooleanList(int count, byte[] data) throws ParseException {
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
    }*/

    private WriteBufferByteBased getWriteBuffer(int size, ModbusByteOrder byteOrder) {
        switch (byteOrder) {
            case LITTLE_ENDIAN: {
                // [4, 3, 2, 1]
                // [8, 7, 6, 5, 4, 3, 2, 1]
                return new WriteBufferByteBased(size, ByteOrder.LITTLE_ENDIAN);
            }
            case BIG_ENDIAN_BYTE_SWAP: {
                // [2, 1, 4, 3]
                // [2, 1, 4, 3, 6, 5, 8, 7]
                return new WriteBufferByteBased(size, ByteOrder.BIG_ENDIAN);
            }
            case LITTLE_ENDIAN_BYTE_SWAP: {
                // [3, 4, 1, 2]
                // [7, 8, 5, 6, 3, 4, 1, 2]
                return new WriteBufferByteBased(size, ByteOrder.LITTLE_ENDIAN);
            }
            default:
                // 16909060
                // [1, 2, 3, 4]
                // 72623859790382856
                // [1, 2, 3, 4, 5, 6, 7, 8]
                return new WriteBufferByteBased(size, ByteOrder.BIG_ENDIAN);
        }
    }

    public static byte[] byteSwap(byte[] in) {
        byte[] out = new byte[in.length];
        for(int i = 0; i < out.length; i += 2) {
            out[i] = in[i + 1];
            out[i + 1] = in[i];
        }
        return out;
    }

}
