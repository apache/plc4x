/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.sandbox.java.dynamic.s7.actions;

import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.commons.scxml2.ActionExecutionContext;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.base.messages.InternalPlcReadRequest;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.java.base.messages.items.*;
import org.apache.plc4x.sandbox.java.dynamic.actions.ReceiveResponseAction;
import org.apache.plc4x.sandbox.java.dynamic.s7.types.DataTransportErrorCode;
import org.apache.plc4x.sandbox.java.dynamic.s7.utils.S7Field;
import org.jdom2.Document;
import org.jdom2.Element;
import org.jdom2.Namespace;
import org.jdom2.filter.Filters;
import org.jdom2.xpath.XPathExpression;
import org.jdom2.xpath.XPathFactory;

import javax.xml.bind.DatatypeConverter;
import java.math.BigInteger;
import java.nio.ByteBuffer;
import java.util.*;

public class S7DecodeReadResponseAction extends ReceiveResponseAction {

    @Override
    protected void processResponse(Document message, ActionExecutionContext ctx, PlcRequestContainer container) {
        InternalPlcReadRequest readRequest = (InternalPlcReadRequest) container.getRequest();
        LinkedHashSet<String> fieldNames = readRequest.getFieldNames();

        List<Namespace> namespaces = message.getRootElement().getNamespacesInScope();
        XPathFactory xPathFactory = XPathFactory.instance();
        XPathExpression<Element> xpath = xPathFactory.compile(
            "/s7:TpktMessage/userData/userData/s7:S7ResponseMessage/payloads/payload/s7:S7ResponsePayloadReadVar/item", Filters.element(), null, namespaces);
        List<Element> items = xpath.evaluate(message);

        // If the sizes don't match, something is wrong.
        if(items.size() != fieldNames.size()) {
            container.getResponseFuture().completeExceptionally(
                new PlcProtocolException("Response item size doesn't match request item size."));
            return;
        }

        Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> responseItems = new HashMap<>(items.size());
        int i = 0;
        for (String fieldName : fieldNames) {
            S7Field field = (S7Field) readRequest.getField(fieldName);

            Element item = items.get(i);
            String returnCode = item.getChild("returnCode").getTextTrim();
            byte returnCodeValue = (byte) (Short.valueOf(returnCode) & 0xFF);
            PlcResponseCode responseCode = null;
            BaseDefaultFieldItem fieldItem = null;
            switch(DataTransportErrorCode.valueOf(returnCodeValue)) {
                case RESERVED:
                    responseCode = PlcResponseCode.INTERNAL_ERROR;
                    break;
                case ACCESS_DENIED:
                    responseCode = PlcResponseCode.ACCESS_DENIED;
                    break;
                case DATA_TYPE_NOT_SUPPORTED:
                    responseCode = PlcResponseCode.INVALID_DATATYPE;
                    break;
                case INVALID_ADDRESS:
                    responseCode = PlcResponseCode.INVALID_ADDRESS;
                    break;
                case NOT_FOUND:
                    responseCode = PlcResponseCode.NOT_FOUND;
                    break;
                case OK: {
                    responseCode = PlcResponseCode.OK;
                    // Convert the hex encoded payload into a real byte array.
                    byte[] data = DatatypeConverter.parseHexBinary(item.getChild("data").getTextTrim());

                    // Depending on the field type in the request, interpret the data in the response accordingly.
                    switch (field.getDataType()) {
                        case BOOL: {
                            BitSet bits = BitSet.valueOf(data);
                            Boolean[] values = new Boolean[field.getNumElements()];
                            for (int bitNr = 0; bitNr < field.getNumElements(); bitNr++) {
                                values[bitNr] = bits.get(bitNr);
                            }
                            fieldItem = new DefaultBooleanFieldItem(values);
                            break;
                        }
                        // -----------------------------------------
                        // Bit strings
                        // In Tia and Step7 these types are marked as bit-strings
                        // which is a sequence of bits. Therefore the result will
                        // have more items than the request requested.
                        // -----------------------------------------
                        case BYTE: {
                            BitSet bits = BitSet.valueOf(data);
                            Boolean[] values = new Boolean[field.getNumElements() * 8];
                            for (int bitNr = 0; bitNr < field.getNumElements() * 8; bitNr++) {
                                values[bitNr] = bits.get(bitNr);
                            }
                            fieldItem = new DefaultBooleanFieldItem(values);
                            break;
                        }
                        case WORD: {
                            BitSet bits = BitSet.valueOf(data);
                            Boolean[] values = new Boolean[field.getNumElements() * 16];
                            for (int bitNr = 0; bitNr < field.getNumElements() * 16; bitNr++) {
                                values[bitNr] = bits.get(bitNr);
                            }
                            fieldItem = new DefaultBooleanFieldItem(values);
                            break;
                        }
                        case DWORD: {
                            BitSet bits = BitSet.valueOf(data);
                            Boolean[] values = new Boolean[field.getNumElements() * 32];
                            for (int bitNr = 0; bitNr < field.getNumElements() * 32; bitNr++) {
                                values[bitNr] = bits.get(bitNr);
                            }
                            fieldItem = new DefaultBooleanFieldItem(values);
                            break;
                        }
                        case LWORD: {
                            BitSet bits = BitSet.valueOf(data);
                            Boolean[] values = new Boolean[field.getNumElements() * 64];
                            for (int bitNr = 0; bitNr < field.getNumElements() * 64; bitNr++) {
                                values[bitNr] = bits.get(bitNr);
                            }
                            fieldItem = new DefaultBooleanFieldItem(values);
                            break;
                        }

                        // -----------------------------------------
                        // Integers
                        // -----------------------------------------
                        // (Signed) Small Int (8 bit)
                        case SINT: {
                            Byte[] values = new Byte[field.getNumElements()];
                            for (int valueNr = 0; valueNr < field.getNumElements(); valueNr++) {
                                values[valueNr] = data[valueNr];
                            }
                            fieldItem = new DefaultByteFieldItem(values);
                            break;
                        }
                        // Unsigned Small Int (8 bit)
                        case USINT: {
                            Short[] values = new Short[field.getNumElements()];
                            for (int valueNr = 0; valueNr < field.getNumElements(); valueNr++) {
                                values[valueNr] = (short) (data[valueNr] & 0xff);
                            }
                            fieldItem = new DefaultShortFieldItem(values);
                            break;
                        }
                        // Signed Int (16 bit)
                        case INT: {
                            ByteBuffer buffer = ByteBuffer.wrap(data);
                            Short[] values = new Short[field.getNumElements()];
                            for (int valueNr = 0; valueNr < field.getNumElements(); valueNr++) {
                                values[valueNr] = buffer.getShort();
                            }
                            fieldItem = new DefaultShortFieldItem(values);
                            break;
                        }
                        // Unsigned Int (16 bit)
                        case UINT: {
                            ByteBuffer buffer = ByteBuffer.wrap(data);
                            Integer[] values = new Integer[field.getNumElements()];
                            for (int valueNr = 0; valueNr < field.getNumElements(); valueNr++) {
                                values[valueNr] = (buffer.getShort() & 0xFFFF);
                            }
                            fieldItem = new DefaultIntegerFieldItem(values);
                            break;
                        }
                        // Double Precision Int (32 bit)
                        case DINT: {
                            ByteBuffer buffer = ByteBuffer.wrap(data);
                            Integer[] values = new Integer[field.getNumElements()];
                            for (int valueNr = 0; valueNr < field.getNumElements(); valueNr++) {
                                values[valueNr] = buffer.getInt();
                            }
                            fieldItem = new DefaultIntegerFieldItem(values);
                            break;
                        }
                        // Unsigned Double Precision Int (32 bit)
                        case UDINT: {
                            ByteBuffer buffer = ByteBuffer.wrap(data);
                            Long[] values = new Long[field.getNumElements()];
                            for (int valueNr = 0; valueNr < field.getNumElements(); valueNr++) {
                                values[valueNr] = (buffer.getInt() & 0xFFFFFFFFL);
                            }
                            fieldItem = new DefaultLongFieldItem(values);
                            break;
                        }
                        // Quadrupal Precision Int (64 bit)
                        case LINT: {
                            ByteBuffer buffer = ByteBuffer.wrap(data);
                            Long[] values = new Long[field.getNumElements()];
                            for (int valueNr = 0; valueNr < field.getNumElements(); valueNr++) {
                                values[valueNr] = buffer.getLong();
                            }
                            fieldItem = new DefaultLongFieldItem(values);
                            break;
                        }
                        // Unsigned Quadrupal Precision Int (64 bit)
                        case ULINT: {
                            BigInteger[] values = new BigInteger[field.getNumElements()];
                            for (int valueNr = 0; valueNr < field.getNumElements(); valueNr++) {
                                byte[] biBytes = new byte[]{data[valueNr * 8], data[(valueNr * 8) + 1],
                                    data[(valueNr * 8) + 2], data[(valueNr * 8) + 3], data[(valueNr * 8) + 4],
                                    data[(valueNr * 8) + 5], data[(valueNr * 8) + 6], data[(valueNr * 8) + 7]};
                                values[valueNr] = new BigInteger(biBytes);
                            }
                            fieldItem = new DefaultBigIntegerFieldItem(values);
                            break;
                        }

                        // -----------------------------------------
                        // Reals
                        // -----------------------------------------
                        // (32 bit)
                        case REAL: {
                            ByteBuffer buffer = ByteBuffer.wrap(data);
                            Float[] values = new Float[field.getNumElements()];
                            for (int valueNr = 0; valueNr < field.getNumElements(); valueNr++) {
                                values[valueNr] = buffer.getFloat();
                            }
                            fieldItem = new DefaultFloatFieldItem(values);
                            break;
                        }
                        // (64 bit)
                        case LREAL: {
                            ByteBuffer buffer = ByteBuffer.wrap(data);
                            Double[] values = new Double[field.getNumElements()];
                            for (int valueNr = 0; valueNr < field.getNumElements(); valueNr++) {
                                values[valueNr] = buffer.getDouble();
                            }
                            fieldItem = new DefaultDoubleFieldItem(values);
                            break;
                        }

                        // -----------------------------------------
                        // Durations
                        // -----------------------------------------
                        // IEC time
                        case TIME:
                            break;
                        case LTIME:
                            break;

                        // -----------------------------------------
                        // Date
                        // -----------------------------------------
                        // IEC date (yyyy-m-d)
                        case DATE:
                            break;

                        // -----------------------------------------
                        // Time of day
                        // -----------------------------------------
                        // Time (hh:mm:ss.S)
                        case TIME_OF_DAY:
                            break;

                        // -----------------------------------------
                        // Date and time of day
                        // -----------------------------------------
                        case DATE_AND_TIME:
                            break;

                        // -----------------------------------------
                        // ASCII Strings
                        // -----------------------------------------
                        // Single-byte character
                        case CHAR:
                            break;
                        // Double-byte character
                        case WCHAR:
                            break;
                        // Variable-length single-byte character string
                        case STRING:
                            break;
                        // Variable-length double-byte character string
                        case WSTRING:
                            break;
                    }
                }
            }
            responseItems.put(fieldName, new ImmutablePair<>(responseCode, fieldItem));
            i++;
        }

        PlcReadResponse response = new DefaultPlcReadResponse(readRequest, responseItems);
        container.getResponseFuture().complete(response);
    }

}
