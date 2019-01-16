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
package org.apache.plc4x.java.s7.netty.strategies;

import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.messages.S7ResponseMessage;
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.S7Payload;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.*;
import org.hamcrest.core.IsCollectionContaining;
import org.junit.Before;
import org.junit.Test;

import java.util.*;

import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.lessThan;
import static org.hamcrest.collection.IsCollectionWithSize.hasSize;
import static org.hamcrest.core.Is.is;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.hamcrest.core.IsNull.nullValue;
import static org.junit.Assert.assertThat;

public class DefaultS7MessageProcessorTest {

    private S7MessageProcessor SUT;

    @Before
    public void setUp() {
        SUT = new DefaultS7MessageProcessor();
    }

    /**
     * In this test both the size of the request as well as the estimated response will be well within
     * the bounds set by the PDU size parameter, so we are expecting the processor to not change anything.
     *
     * @throws PlcException something went wrong.
     */
    @Test
    public void readMessageSimple() throws PlcException {
        S7RequestMessage request = createReadMessage(
            Collections.singletonList(
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 1, (short) 0, (byte) 0)));
        Collection<S7RequestMessage> processedRequests = SUT.processRequest(request, 250);

        assertThat(processedRequests, notNullValue());
        assertThat(processedRequests, hasSize(1));

        Optional<VarParameter> parameter = processedRequests.iterator().next().getParameter(VarParameter.class);
        assertThat(parameter.isPresent(), is(true));
        VarParameter varParameter = parameter.get();
        assertThat(varParameter.getItems(), hasSize(1));
    }

    /**
     * In this test both the size of the request as well as the estimated response will be well within
     * the bounds set by the PDU size parameter, so we are expecting the processor to not change anything.
     *
     * @throws PlcException something went wrong.
     */
    @Test
    public void readMessageMultipleItems() throws PlcException {
        S7RequestMessage request = createReadMessage(
            Arrays.asList(
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 1, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 2, (short) 0, (byte) 0)));
        Collection<S7RequestMessage> processedRequests = SUT.processRequest(request, 250);

        assertThat(processedRequests, notNullValue());
        assertThat(processedRequests, hasSize(1));

        Optional<VarParameter> parameter = processedRequests.iterator().next().getParameter(VarParameter.class);
        assertThat(parameter.isPresent(), is(true));
        VarParameter varParameter = parameter.get();
        assertThat(varParameter.getItems(), hasSize(2));
    }

    /**
     * The maximum number of request items to fit into a PDU with the size 250 is usually 19. So this request
     * should just fit into one message and therefore we expect the processor to leave the message unchanged.
     *
     * @throws PlcException something went wrong.
     */
    @Test
    public void readMessageMultipleItemsWithLargeRequestSize() throws PlcException {
        S7RequestMessage request = createReadMessage(
            Arrays.asList(
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 1, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 2, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 3, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 4, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 5, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 6, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 7, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 8, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 9, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 10, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 11, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 12, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 13, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 14, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 15, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 16, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 17, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 18, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 19, (short) 0, (byte) 0)));
        Collection<S7RequestMessage> processedRequests = SUT.processRequest(request, 250);

        assertThat(processedRequests, notNullValue());
        assertThat(processedRequests, hasSize(1));

        Optional<VarParameter> parameter = processedRequests.iterator().next().getParameter(VarParameter.class);
        assertThat(parameter.isPresent(), is(true));
        VarParameter varParameter = parameter.get();
        assertThat(varParameter.getItems(), hasSize(19));
    }

    /**
     * In this request, the request size itself exceeds the bounds set by the PDU size parameter,
     * therefore the processor is expected to split the one request with two items up into two requests
     * with each one item.
     *
     * @throws PlcException something went wrong.
     */
    @Test
    public void readMessageMultipleItemsWithTooLargeRequestSize() throws PlcException {
        S7RequestMessage request = createReadMessage(
            Arrays.asList(
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 1, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 2, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 3, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 4, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 5, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 6, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 7, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 8, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 9, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 10, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 11, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 12, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 13, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 14, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 15, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 16, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 17, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 18, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 19, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 20, (short) 0, (byte) 0)));
        Collection<S7RequestMessage> processedRequests = SUT.processRequest(request, 250);

        assertThat(processedRequests, notNullValue());
        assertThat(processedRequests, hasSize(2));

        int totalItems = 0;
        for (S7RequestMessage requestMessage : processedRequests) {
            Optional<VarParameter> parameter = requestMessage.getParameter(VarParameter.class);
            assertThat(parameter.isPresent(), is(true));
            VarParameter varParameter = parameter.get();
            int numItems = varParameter.getItems().size();
            // From calculation and byte counting we know that with a pdu-size of 250 a read message can't have
            // more than 19 items.
            assertThat(numItems, lessThan(20));
            totalItems += numItems;
        }
        // In total 20 items should have been found.
        assertThat(totalItems, is(20));
    }

    /**
     * In this request, the request size itself is way within the bounds set by the PDU size parameter,
     * however the estimated size of the response would exceed this greatly, therefore the processor is
     * expected to split the one request with two items up into two requests with each one item.
     *
     * @throws PlcException something went wrong.
     */
    @Test
    public void readMessageCompositeWithLargeResponseSize() throws PlcException {
        S7RequestMessage request = createReadMessage(
            Arrays.asList(
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 200, (short) 1, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 200, (short) 2, (short) 0, (byte) 0)));
        Collection<S7RequestMessage> processedRequests = SUT.processRequest(request, 256);

        assertThat(processedRequests, notNullValue());
        assertThat(processedRequests, hasSize(2));

        for (S7RequestMessage requestMessage : processedRequests) {
            Optional<VarParameter> parameter = requestMessage.getParameter(VarParameter.class);
            assertThat(parameter.isPresent(), is(true));
            VarParameter varParameter = parameter.get();
            assertThat(varParameter.getItems(), hasSize(1));
        }
    }

    /**
     * In this request, the request size itself is way within the bounds set by the PDU size parameter,
     * however the estimated size of the response would exceed this greatly. In contrast to the test above
     * a single item is too large for being sent at all. So in this case the item should be split up into
     * multiple element and joined together after all sub-messages are returned.
     *
     * @throws PlcException something went wrong.
     */
    @Test
    public void readMessageWithTooLargeResponseSize() throws PlcException {
        S7RequestMessage request = createReadMessage(
            Collections.singletonList(
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.REAL, (short) 400, (short) 1, (short) 0, (byte) 0)));
        Collection<S7RequestMessage> processedRequests = SUT.processRequest(request, 256);

        assertThat(processedRequests, notNullValue());
        assertThat(processedRequests, hasSize(7));

        int totalNumItems = 0;
        for (S7RequestMessage requestMessage : processedRequests) {
            Optional<VarParameter> parameter = requestMessage.getParameter(VarParameter.class);
            assertThat(parameter.isPresent(), is(true));
            VarParameter varParameter = parameter.get();
            assertThat(varParameter.getItems(), hasSize(1));
            totalNumItems += ((S7AnyVarParameterItem) varParameter.getItems().iterator().next()).getNumElements();
        }
        assertThat(totalNumItems, equalTo(400));
    }

    /**
     * In this request, we only send one single element to one single field. Nothing should be changed.
     *
     * @throws PlcException something went wrong.
     */
    @Test
    public void writeMessageSimple() throws PlcException {
        S7RequestMessage request = createWriteMessage(
            Collections.singletonList(
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 1, (short) 0, (byte) 0)),
            Collections.singletonList(
                new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[] {0x00}))
            );
        Collection<S7RequestMessage> processedRequests = SUT.processRequest(request, 250);

        assertThat(processedRequests, notNullValue());
        assertThat(processedRequests, hasSize(1));

        S7RequestMessage processedRequest = processedRequests.iterator().next();

        // Check the parameter
        Optional<VarParameter> parameter = processedRequest.getParameter(VarParameter.class);
        assertThat(parameter.isPresent(), is(true));
        VarParameter varParameter = parameter.get();
        assertThat(varParameter.getItems(), hasSize(1));

        // Check the payload
        Optional<VarPayload> payload = processedRequest.getPayload(VarPayload.class);
        assertThat(payload.isPresent(), is(true));
        VarPayload varPayload = payload.get();
        assertThat(varPayload.getItems(), hasSize(1));
    }

    /**
     * In this request, we send an array of bit elements to a single field, the request should be broken
     * up into multiple single element write messages as the S7 doesn't seem to like writing of arrays.
     *
     * @throws PlcException something went wrong.
     */
    @Test
    public void writeMessageSimpleBitArray() throws PlcException {
        S7RequestMessage request = createWriteMessage(
            Collections.singletonList(
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BOOL, (short) 10, (short) 1, (short) 0, (byte) 0)),
            Collections.singletonList(
                new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BIT, new byte[] {
                    (byte) 0xAA, (byte) 0x02}))
        );
        Collection<S7RequestMessage> processedRequests = SUT.processRequest(request, 250);

        // Initialize a set of expected fields.
        Set<String> expectedFields = new HashSet<>(10);
        for(int i = 0; i < 10; i++) {
            expectedFields.add(i / 8 + "/" + i % 8);
        }

        // We are expecting to receive 10 messages as we had an array of 10 items.
        assertThat(processedRequests, hasSize(10));
        // Process all the messages and check each one
        for (S7RequestMessage processedRequest : processedRequests) {
            Optional<VarParameter> parameterOptional = processedRequest.getParameter(VarParameter.class);
            assertThat(parameterOptional.isPresent(), is(true));
            VarParameter varParameter = parameterOptional.get();

            Optional<VarPayload> payloadOptional = processedRequest.getPayload(VarPayload.class);
            assertThat(payloadOptional.isPresent(), is(true));
            VarPayload varPayload = payloadOptional.get();

            assertThat(varParameter.getItems(), hasSize(1));
            assertThat(varPayload.getItems(), hasSize(1));

            VarParameterItem parameterItem = varParameter.getItems().iterator().next();
            assertThat(parameterItem.getAddressingMode(), is(VariableAddressingMode.S7ANY));
            S7AnyVarParameterItem s7AnyParameterItem = (S7AnyVarParameterItem) parameterItem;
            assertThat(s7AnyParameterItem.getMemoryArea(), is(MemoryArea.DATA_BLOCKS));
            assertThat(s7AnyParameterItem.getDataType(), is(TransportSize.BOOL));
            assertThat(s7AnyParameterItem.getNumElements(), is(1));
            String fieldString = s7AnyParameterItem.getByteOffset() + "/" + s7AnyParameterItem.getBitOffset();
            assertThat(expectedFields, IsCollectionContaining.hasItem(fieldString));

            VarPayloadItem payloadItem = varPayload.getItems().iterator().next();
            // We are expecting that the payload is simply "the field of the byte + 1".
            assertThat(payloadItem.getData().length, is(1));
            int value = (s7AnyParameterItem.getByteOffset() * 8) + s7AnyParameterItem.getBitOffset();
            byte expectedValue = (value % 2 == 0) ? (byte) 0x00 : (byte) 0x01;
            assertThat(payloadItem.getData()[0], is(expectedValue));

            // Remove the used field from the list of available ones.
            expectedFields.remove(fieldString);
        }

        // In the end all fields should have been used.
        assertThat(expectedFields, hasSize(0));
    }

    /**
     * In this request, we send an array of elements to a single field, the request should be broken
     * up into multiple single element write messages as the S7 doesn't seem to like writing of arrays.
     *
     * @throws PlcException something went wrong.
     */
    @Test
    public void writeMessageSimpleByteArray() throws PlcException {
        S7RequestMessage request = createWriteMessage(
            Collections.singletonList(
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, 10, 1, (short) 0, (byte) 0)),
            Collections.singletonList(
                new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[] {
                    0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A}))
        );
        Collection<S7RequestMessage> processedRequests = SUT.processRequest(request, 250);

        // Initialize a set of expected fields.
        Set<Integer> expectedFields = new HashSet<>(10);
        for(int i = 0; i < 10; i++) {
            expectedFields.add(i);
        }

        // We are expecting to receive 10 messages as we had an array of 10 items.
        assertThat(processedRequests, hasSize(10));
        // Process all the messages and check each one
        for (S7RequestMessage processedRequest : processedRequests) {
            Optional<VarParameter> parameterOptional = processedRequest.getParameter(VarParameter.class);
            assertThat(parameterOptional.isPresent(), is(true));
            VarParameter varParameter = parameterOptional.get();

            Optional<VarPayload> payloadOptional = processedRequest.getPayload(VarPayload.class);
            assertThat(payloadOptional.isPresent(), is(true));
            VarPayload varPayload = payloadOptional.get();

            assertThat(varParameter.getItems(), hasSize(1));
            assertThat(varPayload.getItems(), hasSize(1));

            VarParameterItem parameterItem = varParameter.getItems().iterator().next();
            assertThat(parameterItem.getAddressingMode(), is(VariableAddressingMode.S7ANY));
            S7AnyVarParameterItem s7AnyParameterItem = (S7AnyVarParameterItem) parameterItem;
            assertThat(s7AnyParameterItem.getMemoryArea(), is(MemoryArea.DATA_BLOCKS));
            assertThat(s7AnyParameterItem.getDataType(), is(TransportSize.BYTE));
            assertThat(s7AnyParameterItem.getNumElements(), is(1));
            // Check the field is in the expected range and hasn't been used yet.
            assertThat(expectedFields.contains(s7AnyParameterItem.getByteOffset()), is(true));
            assertThat(s7AnyParameterItem.getBitOffset(), is((byte) 0));

            VarPayloadItem payloadItem = varPayload.getItems().iterator().next();
            // We are expecting that the payload is simply "the field of the byte + 1".
            assertThat(payloadItem.getData().length, is(1));
            byte expectedValue = (byte) ((byte) s7AnyParameterItem.getByteOffset() + (byte) 1);
            assertThat(payloadItem.getData()[0], is(expectedValue));

            // Remove the used field from the list of available ones.
            expectedFields.remove(s7AnyParameterItem.getByteOffset());
        }

        // In the end all fields should have been used.
        assertThat(expectedFields, hasSize(0));
    }

    /**
     * In this request, we send an array of elements to a single field, the request should be broken
     * up into multiple single element write messages as the S7 doesn't seem to like writing of arrays.
     *
     * @throws PlcException something went wrong.
     */
    @Test
    public void writeMessageSimpleDwordArray() throws PlcException {
        S7RequestMessage request = createWriteMessage(
            Collections.singletonList(
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.DWORD, (short) 10, (short) 1, (short) 0, (byte) 0)),
            Collections.singletonList(
                new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[] {
                    0x00, 0x00, 0x00, 0x01,
                    0x00, 0x00, 0x00, 0x02,
                    0x00, 0x00, 0x00, 0x03,
                    0x00, 0x00, 0x00, 0x04,
                    0x00, 0x00, 0x00, 0x05,
                    0x00, 0x00, 0x00, 0x06,
                    0x00, 0x00, 0x00, 0x07,
                    0x00, 0x00, 0x00, 0x08,
                    0x00, 0x00, 0x00, 0x09,
                    0x00, 0x00, 0x00, 0x0A})));
        Collection<S7RequestMessage> processedRequests = SUT.processRequest(request, 250);

        // Initialize a set of expected fields.
        Set<Integer> expectedFields = new HashSet<>(10);
        for(int i = 0; i < 10; i++) {
            expectedFields.add((i * 4));
        }

        // We are expecting to receive 10 messages as we had an array of 10 items.
        assertThat(processedRequests, hasSize(10));
        // Process all the messages and check each one
        for (S7RequestMessage processedRequest : processedRequests) {
            Optional<VarParameter> parameterOptional = processedRequest.getParameter(VarParameter.class);
            assertThat(parameterOptional.isPresent(), is(true));
            VarParameter varParameter = parameterOptional.get();

            Optional<VarPayload> payloadOptional = processedRequest.getPayload(VarPayload.class);
            assertThat(payloadOptional.isPresent(), is(true));
            VarPayload varPayload = payloadOptional.get();

            assertThat(varParameter.getItems(), hasSize(1));
            assertThat(varPayload.getItems(), hasSize(1));

            VarParameterItem parameterItem = varParameter.getItems().iterator().next();
            assertThat(parameterItem.getAddressingMode(), is(VariableAddressingMode.S7ANY));
            S7AnyVarParameterItem s7AnyParameterItem = (S7AnyVarParameterItem) parameterItem;
            assertThat(s7AnyParameterItem.getMemoryArea(), is(MemoryArea.DATA_BLOCKS));
            assertThat(s7AnyParameterItem.getDataType(), is(TransportSize.DWORD));
            assertThat(s7AnyParameterItem.getNumElements(), is(1));
            // Check the field is in the expected range and hasn't been used yet.
            assertThat(expectedFields.contains(s7AnyParameterItem.getByteOffset()), is(true));
            assertThat(s7AnyParameterItem.getBitOffset(), is((byte) 0));

            VarPayloadItem payloadItem = varPayload.getItems().iterator().next();
            assertThat(payloadItem.getData().length, is(4));
            // We are expecting that the payload is the index number of the dword item in the original array.
            int expectedValue = (s7AnyParameterItem.getByteOffset() / 4) + 1;
            int actualValue = (payloadItem.getData()[0] << 32) + (payloadItem.getData()[1] << 16) + (payloadItem.getData()[2] << 8) + payloadItem.getData()[3];
            assertThat(actualValue, is(expectedValue));

            // Remove the used field from the list of available ones.
            expectedFields.remove(s7AnyParameterItem.getByteOffset());
        }

        // In the end all fields should have been used.
        assertThat(expectedFields, hasSize(0));
    }

    /**
     * In this test, we are writing multiple independent items in one message. This has to be split up.
     *
     * @throws PlcException something went wrong.
     */
    @Test
    public void writeMessageMultipleItems() throws PlcException {
        S7RequestMessage request = createWriteMessage(
            Arrays.asList(
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BOOL, (short) 1, (short) 1, (short) 0, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 1, (short) 1, (short) 1, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.DWORD, (short) 1, (short) 1, (short) 2, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.REAL, (short) 1, (short) 1, (short) 5, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.INT, (short) 1, (short) 1, (short) 9, (byte) 0)),
            Arrays.asList(
                new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BIT, new byte[] {0x01}),
                new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[] {0x02}),
                new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[] {0x00, 0x00, 0x00, 0x00}),
                new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[] {0x00, 0x00, 0x00, 0x00}),
                new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[] {0x00, 0x00, 0x00, 0x00})
            )
        );
        Collection<S7RequestMessage> processedRequests = SUT.processRequest(request, 250);

        assertThat(processedRequests, notNullValue());
        assertThat(processedRequests, hasSize(5));

        for (S7RequestMessage processedRequest : processedRequests) {
            // Check the parameter
            Optional<VarParameter> parameter = processedRequest.getParameter(VarParameter.class);
            assertThat(parameter.isPresent(), is(true));
            VarParameter varParameter = parameter.get();
            assertThat(varParameter.getItems(), hasSize(1));

            // Check the payload
            Optional<VarPayload> payload = processedRequest.getPayload(VarPayload.class);
            assertThat(payload.isPresent(), is(true));
            VarPayload varPayload = payload.get();
            assertThat(varPayload.getItems(), hasSize(1));
        }
    }

    /**
     * In this test, we are writing multiple independent array items in one message. This has to be split up both.
     * regarding the independent items, but also regarding the array items.
     *
     * @throws PlcException something went wrong.
     */
    @Test
    public void writeMessageMultipleArrayItems() throws PlcException {
        S7RequestMessage request = createWriteMessage(
            Arrays.asList(
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.BYTE, (short) 2, (short) 1, (short) 1, (byte) 0),
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    TransportSize.DWORD, (short) 2, (short) 1, (short) 2, (byte) 0)),
            Arrays.asList(
                new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[] {
                    0x01, 0x02}),
                new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[] {
                    0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x04})
            )
        );
        Collection<S7RequestMessage> processedRequests = SUT.processRequest(request, 250);

        assertThat(processedRequests, notNullValue());
        assertThat(processedRequests, hasSize(4));

        for (S7RequestMessage processedRequest : processedRequests) {
            // Check the parameter
            Optional<VarParameter> parameter = processedRequest.getParameter(VarParameter.class);
            assertThat(parameter.isPresent(), is(true));
            VarParameter varParameter = parameter.get();
            assertThat(varParameter.getItems(), hasSize(1));

            // Check the payload
            Optional<VarPayload> payload = processedRequest.getPayload(VarPayload.class);
            assertThat(payload.isPresent(), is(true));
            VarPayload varPayload = payload.get();
            assertThat(varPayload.getItems(), hasSize(1));
        }
    }

    /**
     * This test handles the special case in which a response is part of a single request message.
     * This means that it is immediatly finished and is hereby immediatly processed.
     *
     * @throws PlcException
     */
    @Test
    public void processSimpleMessageResponse() throws PlcException {
        S7RequestMessage requestMessage = new S7RequestMessage(MessageType.JOB, (short) 1, Collections.emptyList(), Collections.emptyList(), null);
        S7ResponseMessage responseMessage = new S7ResponseMessage(MessageType.JOB, (short) 1, Collections.emptyList(), Collections.emptyList(), (byte) 0x00, (byte) 0x00);
        S7ResponseMessage processedResponse = SUT.processResponse(requestMessage, responseMessage);
        // In this case the response should be returned unchanged.
        assertThat(processedResponse, is(responseMessage));
    }

    /**
     * This test handles the special case in which a response is part of a single request message.
     * This means that it is immediately finished and is hereby immediately processed.
     *
     * @throws PlcException
     */
    @Test
    public void processCompositeMessageReadResponse() throws PlcException {
        S7RequestMessage originalRequestMessage = new S7RequestMessage(MessageType.JOB, (short) 1,
            Collections.singletonList(
                new VarParameter(ParameterType.READ_VAR,
                    Arrays.asList(
                        new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                            TransportSize.BYTE, (short) 1, (short) 1, (short) 2, (byte) 0),
                        new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                            TransportSize.BYTE, (short) 1, (short) 3, (short) 4, (byte) 0))
                )
            ),
            Collections.emptyList(),
            null);
        DefaultS7MessageProcessor.S7CompositeRequestMessage compositeRequestMessage =
            new DefaultS7MessageProcessor.S7CompositeRequestMessage(originalRequestMessage);

        S7RequestMessage fragment1RequestMessage = new S7RequestMessage(MessageType.JOB, (short) 1,
            Collections.singletonList(
                new VarParameter(ParameterType.READ_VAR,
                    Arrays.asList(
                        new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                            TransportSize.BYTE, (short) 1, (short) 1, (short) 2, (byte) 0))
                )
            ),
            Collections.emptyList(),
            compositeRequestMessage);
        compositeRequestMessage.addRequestMessage(fragment1RequestMessage);

        S7RequestMessage fragment2RequestMessage = new S7RequestMessage(MessageType.JOB, (short) 3,
            Collections.singletonList(
                new VarParameter(ParameterType.READ_VAR,
                    Arrays.asList(
                        new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                            TransportSize.BYTE, (short) 1, (short) 3, (short) 4, (byte) 0))
                )
            ),
            Collections.emptyList(),
            compositeRequestMessage);
        compositeRequestMessage.addRequestMessage(fragment2RequestMessage);

        // Virtually add a response for the first response.
        fragment1RequestMessage.setAcknowledged(true);
        S7ResponseMessage fragment1ResponseMessage = new S7ResponseMessage(MessageType.JOB, (short) 2,
            Collections.singletonList(
                new VarParameter(ParameterType.READ_VAR, new LinkedList<>(Collections.singletonList(
                    new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                        TransportSize.BYTE, (short) 1, (short) 1, (short) 2, (byte) 0))))),
            Collections.singletonList(
                new VarPayload(ParameterType.READ_VAR, new LinkedList<>(Collections.singletonList(
                    new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[]{0x42}))))),
            (byte) 0x00, (byte) 0x00);
        S7ResponseMessage processedResponse = SUT.processResponse(fragment1RequestMessage, fragment1ResponseMessage);
        // As only one of the two requests is responded, the result should be null.
        assertThat(processedResponse, nullValue());

        // Virtually add a response for the second response.
        fragment2RequestMessage.setAcknowledged(true);
        S7ResponseMessage fragment2ResponseMessage = new S7ResponseMessage(MessageType.JOB, (short) 3,
            Collections.singletonList(
                new VarParameter(ParameterType.READ_VAR, new LinkedList<>(Collections.singletonList(
                    new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                        TransportSize.BYTE, (short) 1, (short) 3, (short) 4, (byte) 0))))),
            Collections.singletonList(
                new VarPayload(ParameterType.READ_VAR, new LinkedList<>(Collections.singletonList(
                    new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[]{0x23}))))),
            (byte) 0x00, (byte) 0x00);
        // This time we expect all messages of the composite to be acknowledged and the processResponse should
        // return a merged version of the individual responses content.
        processedResponse = SUT.processResponse(fragment2RequestMessage, fragment2ResponseMessage);
        // As this is the last request being responded, the result should be not null this time.
        assertThat(processedResponse, notNullValue());

        // Check the content.
        assertThat(processedResponse.getParameters(), hasSize(1));
        assertThat(processedResponse.getParameter(VarParameter.class).isPresent(), is(true));
        VarParameter varParameter = processedResponse.getParameter(VarParameter.class).get();
        assertThat(varParameter.getItems(), hasSize(2));

        assertThat(processedResponse.getPayloads(), hasSize(1));
        assertThat(processedResponse.getPayload(VarPayload.class).isPresent(), is(true));

        VarPayload varPayload = processedResponse.getPayload(VarPayload.class).get();
        assertThat(varPayload.getItems(), hasSize(2));
    }

    /**
     * This test handles the special case in which a response is part of a single request message.
     * This means that it is immediately finished and is hereby immediately processed.
     *
     * @throws PlcException
     */
    @Test
    public void processCompositeMessageMergedReadResponse() throws PlcException {
        // Produce a composite request, consisting of an original request reading 4*10=40 items.
        // Where this is split up into 4 sub-messages with each 10 items.
        DefaultS7MessageProcessor.S7CompositeRequestMessage compositeRequestMessage = createCompositeReadMessage(4, 10);

        List<S7RequestMessage> requestMessages = compositeRequestMessage.getRequestMessages();
        // Generate dummy responses for each sub request.
        List<S7ResponseMessage> responseMessages = createResponseMessages(requestMessages);

        // Iterate over all pairs of requests and responses and have the processor process that.
        S7ResponseMessage processedResponse = null;
        int size = requestMessages.size();
        for (int i = 0; i < size; i++) {
            S7RequestMessage requestMessage = requestMessages.get(i);
            S7ResponseMessage responseMessage = responseMessages.get(i);
            requestMessage.setAcknowledged(true);

            processedResponse = SUT.processResponse(requestMessage, responseMessage);
            // Only after processing the least response, should the processor return something not null.
            if(i < size - 1) {
                assertThat(processedResponse, nullValue());
            }
        }

        // Check the content.
        assertThat(processedResponse, notNullValue());
        assertThat(processedResponse.getParameters(), hasSize(1));
        assertThat(processedResponse.getParameter(VarParameter.class).isPresent(), is(true));
        VarParameter varParameter = processedResponse.getParameter(VarParameter.class).get();
        assertThat(varParameter.getItems(), hasSize(1));

        assertThat(processedResponse.getPayloads(), hasSize(1));
        assertThat(processedResponse.getPayload(VarPayload.class).isPresent(), is(true));

        VarPayload varPayload = processedResponse.getPayload(VarPayload.class).get();
        assertThat(varPayload.getItems(), hasSize(1));
    }

    /**
     * This test handles the special case in which a response is part of a single request message.
     * This means that it is immediately finished and is hereby immediately processed.
     *
     * @throws PlcException
     */
    @Test
    public void processCompositeMessageWriteResponse() throws PlcException {
        S7RequestMessage originalRequestMessage = new S7RequestMessage(MessageType.JOB, (short) 1,
            Collections.singletonList(
                new VarParameter(ParameterType.WRITE_VAR, new LinkedList<>(Arrays.asList(
                    new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                        TransportSize.BYTE, (short) 1, (short) 1, (short) 2, (byte) 0),
                    new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                        TransportSize.BYTE, (short) 1, (short) 1, (short) 3, (byte) 0))))),
            Collections.emptyList(), null);
        DefaultS7MessageProcessor.S7CompositeRequestMessage compositeRequestMessage =
            new DefaultS7MessageProcessor.S7CompositeRequestMessage(originalRequestMessage);

        S7RequestMessage fragment1RequestMessage = new S7RequestMessage(MessageType.JOB, (short) 2,
            Collections.singletonList(
                new VarParameter(ParameterType.WRITE_VAR, new LinkedList<>(Collections.singletonList(
                    new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                        TransportSize.BYTE, (short) 1, (short) 1, (short) 2, (byte) 0))))),
            Collections.emptyList(), compositeRequestMessage);
        compositeRequestMessage.addRequestMessage(fragment1RequestMessage);
        S7RequestMessage fragment2RequestMessage = new S7RequestMessage(MessageType.JOB, (short) 3,
            Collections.singletonList(
                new VarParameter(ParameterType.WRITE_VAR, new LinkedList<>(Collections.singletonList(
                    new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                        TransportSize.BYTE, (short) 1, (short) 1, (short) 3, (byte) 0))))),
            Collections.emptyList(), compositeRequestMessage);
        compositeRequestMessage.addRequestMessage(fragment2RequestMessage);

        // Virtually add a response for the first response.
        fragment1RequestMessage.setAcknowledged(true);
        S7ResponseMessage fragment1ResponseMessage = new S7ResponseMessage(MessageType.JOB, (short) 2,
            Collections.singletonList(
                new VarParameter(ParameterType.WRITE_VAR, new LinkedList<>(Collections.singletonList(
                    new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                        TransportSize.BYTE, (short) 1, (short) 1, (short) 2, (byte) 0))))),
            Collections.singletonList(
                new VarPayload(ParameterType.WRITE_VAR, new LinkedList<>(Collections.singletonList(
                    new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[]{0x42}))))),
            (byte) 0x00, (byte) 0x00);
        S7ResponseMessage processedResponse = SUT.processResponse(fragment1RequestMessage, fragment1ResponseMessage);
        // As only one of the two requests is responded, the result should be null.
        assertThat(processedResponse, nullValue());

        // Virtually add a response for the second response.
        fragment2RequestMessage.setAcknowledged(true);
        S7ResponseMessage fragment2ResponseMessage = new S7ResponseMessage(MessageType.JOB, (short) 3,
            Collections.singletonList(
                new VarParameter(ParameterType.WRITE_VAR, new LinkedList<>(Collections.singletonList(
                    new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                        TransportSize.BYTE, (short) 1, (short) 1, (short) 3, (byte) 0))))),
            Collections.singletonList(
                new VarPayload(ParameterType.WRITE_VAR, new LinkedList<>(Collections.singletonList(
                    new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[]{0x23}))))),
            (byte) 0x00, (byte) 0x00);

        // This time we expect all messages of the composite to be acknowledged and the processResponse should
        // return a merged version of the individual responses content.
        processedResponse = SUT.processResponse(fragment2RequestMessage, fragment2ResponseMessage);

        // As this is the last request being responded, the result should be not null this time.
        assertThat(processedResponse, notNullValue());

        // Check the content.
        assertThat(processedResponse.getParameters(), hasSize(1));
        assertThat(processedResponse.getParameter(VarParameter.class).isPresent(), is(true));
        VarParameter varParameter = processedResponse.getParameter(VarParameter.class).get();
        assertThat(varParameter.getItems(), hasSize(2));

        assertThat(processedResponse.getPayloads(), hasSize(1));
        assertThat(processedResponse.getPayload(VarPayload.class).isPresent(), is(true));
        VarPayload varPayload = processedResponse.getPayload(VarPayload.class).get();
        assertThat(varPayload.getItems(), hasSize(2));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Helpers
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private S7RequestMessage createReadMessage(List<VarParameterItem> parameterItems) {
        return new S7RequestMessage(MessageType.JOB, (short) 42,
            Collections.singletonList(
                new VarParameter(ParameterType.READ_VAR, parameterItems)),
            Collections.emptyList(), null);
    }

    private S7RequestMessage createWriteMessage(List<VarParameterItem> parameterItems,
                                                List<VarPayloadItem> payloadItems) {
        return new S7RequestMessage(MessageType.JOB, (short) 42,
            Collections.singletonList(
                new VarParameter(ParameterType.WRITE_VAR, parameterItems)),
            Collections.singletonList(
                new VarPayload(ParameterType.WRITE_VAR, payloadItems)),
            null);
    }

    private DefaultS7MessageProcessor.S7CompositeRequestMessage createCompositeReadMessage(int numSubMessages, int numItemsPerSubMessage) {
        int totalItems = numSubMessages * numItemsPerSubMessage;
        S7RequestMessage originalRequestMessage = new S7RequestMessage(MessageType.JOB, (short) 1,
            Collections.singletonList(
                new VarParameter(ParameterType.READ_VAR,
                    Collections.singletonList(
                        new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                            TransportSize.BYTE, (short) totalItems, (short) 1, (short) 0, (byte) 0))
                )
            ),
            Collections.emptyList(),
            null);
        DefaultS7MessageProcessor.S7CompositeRequestMessage compositeRequestMessage =
            new DefaultS7MessageProcessor.S7CompositeRequestMessage(originalRequestMessage);

        int curPos = 0;
        for(int i = 0; i < numSubMessages; i++) {
            S7RequestMessage subMessage = new S7RequestMessage(MessageType.JOB, (short) 1,
                Collections.singletonList(
                    new VarParameter(ParameterType.READ_VAR,
                        Collections.singletonList(
                            new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                                TransportSize.BYTE, (short) numItemsPerSubMessage, (short) 1, (short) curPos, (byte) 0))
                    )
                ),
                Collections.emptyList(),
                compositeRequestMessage);
            curPos += numItemsPerSubMessage;
            compositeRequestMessage.addRequestMessage(subMessage);
        }

        return compositeRequestMessage;
    }

    private DefaultS7MessageProcessor.S7CompositeRequestMessage createCompositeWriteMessage() {
        return null;
    }

    private List<S7ResponseMessage> createResponseMessages(Collection<S7RequestMessage> requests) {
        List<S7ResponseMessage> responses = new ArrayList<>(requests.size());
        byte counter = 0;
        for (S7RequestMessage request : requests) {
            List<S7Payload> payloads = new ArrayList<>(request.getParameters().size());
            for (S7Parameter parameter : request.getParameters()) {
                VarParameter varParameter = (VarParameter) parameter;
                List<VarPayloadItem> varPayloadItems = new ArrayList<>(varParameter.getItems().size());
                for (VarParameterItem item : varParameter.getItems()) {
                    S7AnyVarParameterItem anyVarParameterItem = (S7AnyVarParameterItem) item;
                    byte[] data = new byte[
                        anyVarParameterItem.getDataType().getSizeInBytes() * anyVarParameterItem.getNumElements()];
                    for(int i = 0; i < data.length; i++) {
                        data[i] = counter++;
                    }
                    VarPayloadItem payloadItem = new VarPayloadItem(DataTransportErrorCode.OK,
                        anyVarParameterItem.getDataType().getDataTransportSize(), data);
                    varPayloadItems.add(payloadItem);
                }
                S7Payload payload = new VarPayload(parameter.getType(), varPayloadItems);
                payloads.add(payload);
            }
            S7ResponseMessage response = new S7ResponseMessage(request.getMessageType(), request.getTpduReference(),
                request.getParameters(), payloads, (byte) 0xFF, (byte) 0xFF);
            responses.add(response);
        }
        return responses;
    }

}
