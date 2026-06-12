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
package org.apache.plc4x.java.canopen.conversation;

import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.canopen.CANOpenConnection;
import org.apache.plc4x.java.canopen.readwrite.CANOpenFrame;
import org.apache.plc4x.java.canopen.readwrite.CANOpenPayload;
import org.apache.plc4x.java.canopen.readwrite.CANOpenSDORequest;
import org.apache.plc4x.java.canopen.readwrite.CANOpenSDOResponse;
import org.apache.plc4x.java.canopen.readwrite.CANOpenService;
import org.apache.plc4x.java.canopen.readwrite.CANOpenDataType;
import org.apache.plc4x.java.canopen.readwrite.DataItem;
import org.apache.plc4x.java.canopen.readwrite.IndexAddress;
import org.apache.plc4x.java.canopen.readwrite.SDOAbort;
import org.apache.plc4x.java.canopen.readwrite.SDOAbortRequest;
import org.apache.plc4x.java.canopen.readwrite.SDOAbortResponse;
import org.apache.plc4x.java.canopen.readwrite.SDOInitiateExpeditedUploadResponse;
import org.apache.plc4x.java.canopen.readwrite.SDOInitiateSegmentedUploadResponse;
import org.apache.plc4x.java.canopen.readwrite.SDOInitiateUploadRequest;
import org.apache.plc4x.java.canopen.readwrite.SDOInitiateUploadResponse;
import org.apache.plc4x.java.canopen.readwrite.SDORequestCommand;
import org.apache.plc4x.java.canopen.readwrite.SDOResponse;
import org.apache.plc4x.java.canopen.readwrite.SDOSegmentUploadRequest;
import org.apache.plc4x.java.canopen.readwrite.SDOSegmentUploadResponse;
import org.apache.plc4x.java.canopen.transport.CANOpenAbortException;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.concurrent.CompletableFuture;

/**
 * Performs an SDO upload (read) — either expedited (≤4 byte value in one frame) or
 * segmented (response says "N bytes, fetch in segments").
 */
public class SDOUploadConversation {

    private static final Logger LOGGER = LoggerFactory.getLogger(SDOUploadConversation.class);

    private final CANConversation conversation;
    private final int nodeId;
    private final int answerNodeId;
    private final IndexAddress address;
    private final CANOpenDataType type;
    private final Duration timeout;

    public SDOUploadConversation(CANConversation conversation, int nodeId, int answerNodeId,
                                 IndexAddress address, CANOpenDataType type, Duration timeout) {
        this.conversation = conversation;
        this.nodeId = nodeId;
        this.answerNodeId = answerNodeId;
        this.address = address;
        this.type = type;
        this.timeout = timeout;
    }

    public void execute(CompletableFuture<PlcValue> receiver) {
        SDOInitiateUploadRequest rq = new SDOInitiateUploadRequest(address);
        CompletableFuture<CANOpenFrame> initialResponse = conversation.expect(
            ConversationPredicates.sdoTransmitFrom(answerNodeId), timeout);
        conversation.sendToWire(createFrame(new CANOpenSDORequest(SDORequestCommand.INITIATE_UPLOAD, rq)));

        initialResponse.whenComplete((frame, error) -> {
            if (error != null) {
                receiver.completeExceptionally(error);
                return;
            }
            try {
                SDOResponse response = ((CANOpenSDOResponse) frame.getPayload()).getResponse();
                if (response instanceof SDOAbortResponse) {
                    SDOAbort abort = ((SDOAbortResponse) response).getAbort();
                    receiver.completeExceptionally(new CANOpenAbortException(
                        "Could not initiate upload", abort.getCode()));
                    return;
                }
                if (!(response instanceof SDOInitiateUploadResponse)) {
                    receiver.completeExceptionally(new PlcException(
                        "Unexpected SDO response: " + response.getClass().getSimpleName()));
                    return;
                }
                SDOInitiateUploadResponse upload = (SDOInitiateUploadResponse) response;
                if (!upload.getAddress().equals(address)) {
                    receiver.completeExceptionally(new PlcException(
                        "SDO upload response for the wrong address: " + upload.getAddress()));
                    return;
                }
                handle(receiver, upload);
            } catch (Exception e) {
                receiver.completeExceptionally(e);
            }
        });
    }

    private void handle(CompletableFuture<PlcValue> receiver, SDOInitiateUploadResponse answer) {
        if (answer.getExpedited() && answer.getIndicated()
            && answer.getPayload() instanceof SDOInitiateExpeditedUploadResponse) {
            SDOInitiateExpeditedUploadResponse payload = (SDOInitiateExpeditedUploadResponse) answer.getPayload();
            completeWithValue(receiver, payload.getData(), payload.getData().length);
        } else if (answer.getPayload() instanceof SDOInitiateSegmentedUploadResponse) {
            SDOInitiateSegmentedUploadResponse segment = (SDOInitiateSegmentedUploadResponse) answer.getPayload();
            int totalSize = (int) segment.getBytes();
            fetch(receiver, new byte[totalSize], 0, totalSize, false);
        } else {
            receiver.completeExceptionally(new PlcException("Unsupported SDO upload kind."));
        }
    }

    private void fetch(CompletableFuture<PlcValue> receiver, byte[] storage, int offset, int totalSize, boolean toggle) {
        CompletableFuture<CANOpenFrame> nextSegment = conversation.expect(
            ConversationPredicates.sdoTransmitFrom(answerNodeId), timeout);
        conversation.sendToWire(createFrame(new CANOpenSDORequest(SDORequestCommand.SEGMENT_UPLOAD,
            new SDOSegmentUploadRequest(toggle))));

        nextSegment.whenComplete((frame, error) -> {
            if (error != null) {
                receiver.completeExceptionally(error);
                return;
            }
            try {
                SDOResponse response = ((CANOpenSDOResponse) frame.getPayload()).getResponse();
                if (response instanceof SDOAbortResponse) {
                    SDOAbort abort = ((SDOAbortResponse) response).getAbort();
                    receiver.completeExceptionally(new CANOpenAbortException(
                        "Remote aborted segmented upload", abort.getCode()));
                    return;
                }
                if (!(response instanceof SDOSegmentUploadResponse)) {
                    receiver.completeExceptionally(new PlcException(
                        "Unexpected SDO segment response: " + response.getClass().getSimpleName()));
                    return;
                }
                SDOSegmentUploadResponse seg = (SDOSegmentUploadResponse) response;
                if (seg.getToggle() != toggle) {
                    // Toggle mismatch — abort and surface the error.
                    conversation.sendToWire(createFrame(new CANOpenSDORequest(SDORequestCommand.ABORT,
                        new SDOAbortRequest(new SDOAbort(address, 1000L)))));
                    receiver.completeExceptionally(new CANOpenAbortException("Toggle mismatch", 1000));
                    return;
                }
                byte[] data = seg.getData();
                int newOffset = offset + data.length;
                if (newOffset > storage.length) {
                    // More bytes than the initial response promised — clamp.
                    newOffset = storage.length;
                }
                System.arraycopy(data, 0, storage, offset, newOffset - offset);
                if (seg.getLast()) {
                    completeWithValue(receiver, storage, totalSize);
                } else {
                    fetch(receiver, storage, newOffset, totalSize, !toggle);
                }
            } catch (Exception e) {
                receiver.completeExceptionally(e);
            }
        });
    }

    private void completeWithValue(CompletableFuture<PlcValue> receiver, byte[] bytes, int length) {
        try {
            ReadBufferByteBased buffer = new ReadBufferByteBased(bytes, CANOpenConnection.LITTLE_ENDIAN_OPTIONS);
            PlcValue value = DataItem.staticParse(buffer, type, length);
            receiver.complete(value);
        } catch (BufferException | ArrayIndexOutOfBoundsException e) {
            receiver.completeExceptionally(e);
        }
    }

    private CANOpenFrame createFrame(CANOpenPayload payload) {
        return new CANOpenFrame((short) nodeId, CANOpenService.RECEIVE_SDO, payload);
    }

}
