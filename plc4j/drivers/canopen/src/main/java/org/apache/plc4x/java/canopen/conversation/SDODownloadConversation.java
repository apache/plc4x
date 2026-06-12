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
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.canopen.CANOpenConnection;
import org.apache.plc4x.java.canopen.readwrite.CANOpenDataType;
import org.apache.plc4x.java.canopen.readwrite.CANOpenFrame;
import org.apache.plc4x.java.canopen.readwrite.CANOpenPayload;
import org.apache.plc4x.java.canopen.readwrite.CANOpenSDORequest;
import org.apache.plc4x.java.canopen.readwrite.CANOpenSDOResponse;
import org.apache.plc4x.java.canopen.readwrite.CANOpenService;
import org.apache.plc4x.java.canopen.readwrite.DataItem;
import org.apache.plc4x.java.canopen.readwrite.IndexAddress;
import org.apache.plc4x.java.canopen.readwrite.SDOAbort;
import org.apache.plc4x.java.canopen.readwrite.SDOAbortRequest;
import org.apache.plc4x.java.canopen.readwrite.SDOAbortResponse;
import org.apache.plc4x.java.canopen.readwrite.SDOInitiateDownloadRequest;
import org.apache.plc4x.java.canopen.readwrite.SDOInitiateDownloadResponse;
import org.apache.plc4x.java.canopen.readwrite.SDOInitiateExpeditedUploadResponse;
import org.apache.plc4x.java.canopen.readwrite.SDOInitiateSegmentedUploadResponse;
import org.apache.plc4x.java.canopen.readwrite.SDORequestCommand;
import org.apache.plc4x.java.canopen.readwrite.SDOResponse;
import org.apache.plc4x.java.canopen.readwrite.SDOResponseCommand;
import org.apache.plc4x.java.canopen.readwrite.SDOSegmentDownloadRequest;
import org.apache.plc4x.java.canopen.readwrite.SDOSegmentDownloadResponse;
import org.apache.plc4x.java.canopen.transport.CANOpenAbortException;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;

import java.time.Duration;
import java.util.concurrent.CompletableFuture;

/**
 * Performs an SDO download (write) — expedited (≤4 byte value in one frame) or
 * segmented (announce size, then push 7-byte data segments until done).
 */
public class SDODownloadConversation {

    private final CANConversation conversation;
    private final int nodeId;
    private final int answerNodeId;
    private final IndexAddress indexAddress;
    private final byte[] data;
    private final Duration timeout;

    public SDODownloadConversation(CANConversation conversation, int nodeId, int answerNodeId,
                                   IndexAddress indexAddress, PlcValue value, CANOpenDataType type,
                                   Duration timeout) {
        this.conversation = conversation;
        this.nodeId = nodeId;
        this.answerNodeId = answerNodeId;
        this.indexAddress = indexAddress;
        this.timeout = timeout;
        try {
            int byteLen = DataItem.getLengthInBytes(value, type, value.getLength());
            WriteBufferByteBased buffer = new WriteBufferByteBased(new byte[byteLen], CANOpenConnection.LITTLE_ENDIAN_OPTIONS);
            DataItem.staticSerialize(buffer, value, type, value.getLength());
            data = buffer.getBytes();
        } catch (BufferException e) {
            throw new PlcRuntimeException("Could not serialize SDO write payload", e);
        }
    }

    public void execute(CompletableFuture<PlcResponseCode> receiver) {
        if (data.length > 4) {
            // Segmented path — first announce the size, then push 7-byte chunks.
            SDOInitiateSegmentedUploadResponse size = new SDOInitiateSegmentedUploadResponse((long) data.length);
            CompletableFuture<CANOpenFrame> initResp = conversation.expect(
                ConversationPredicates.sdoTransmitFrom(answerNodeId), timeout);
            conversation.sendToWire(createFrame(new CANOpenSDORequest(SDORequestCommand.INITIATE_DOWNLOAD,
                new SDOInitiateDownloadRequest(false, true, indexAddress, size))));

            initResp.whenComplete((frame, error) -> {
                if (error != null) {
                    receiver.completeExceptionally(error);
                    return;
                }
                SDOResponse response = ((CANOpenSDOResponse) frame.getPayload()).getResponse();
                if (response instanceof SDOAbortResponse) {
                    SDOAbort abort = ((SDOAbortResponse) response).getAbort();
                    receiver.completeExceptionally(new CANOpenAbortException(
                        "Could not initiate download", abort.getCode()));
                    return;
                }
                if (!(response instanceof SDOInitiateDownloadResponse)) {
                    receiver.complete(PlcResponseCode.REMOTE_ERROR);
                    return;
                }
                SDOInitiateDownloadResponse initResponse = (SDOInitiateDownloadResponse) response;
                if (!initResponse.getAddress().equals(indexAddress)) {
                    // Mismatched address — abort and report.
                    conversation.sendToWire(createFrame(new CANOpenSDORequest(SDORequestCommand.ABORT,
                        new SDOAbortRequest(new SDOAbort(indexAddress, 1000L)))));
                    receiver.complete(PlcResponseCode.REMOTE_ERROR);
                    return;
                }
                put(receiver, false, 0);
            });
            return;
        }

        // Expedited path — single frame carries up to 4 bytes of payload.
        SDOInitiateDownloadRequest rq = new SDOInitiateDownloadRequest(
            true, true, indexAddress, new SDOInitiateExpeditedUploadResponse(data));

        CompletableFuture<CANOpenFrame> resp = conversation.expect(
            ConversationPredicates.sdoTransmitFrom(answerNodeId), timeout);
        conversation.sendToWire(createFrame(new CANOpenSDORequest(SDORequestCommand.INITIATE_DOWNLOAD, rq)));

        resp.whenComplete((frame, error) -> {
            if (error != null) {
                receiver.completeExceptionally(error);
                return;
            }
            SDOResponse response = ((CANOpenSDOResponse) frame.getPayload()).getResponse();
            if (response instanceof SDOAbortResponse) {
                SDOAbort abort = ((SDOAbortResponse) response).getAbort();
                receiver.completeExceptionally(new CANOpenAbortException(
                    "Could not initiate download", abort.getCode()));
                return;
            }
            if (response instanceof SDOInitiateDownloadResponse
                && ((SDOInitiateDownloadResponse) response).getCommand() == SDOResponseCommand.INITIATE_DOWNLOAD) {
                receiver.complete(PlcResponseCode.OK);
            } else {
                receiver.complete(PlcResponseCode.REMOTE_ERROR);
            }
        });
    }

    private void put(CompletableFuture<PlcResponseCode> receiver, boolean toggle, int offset) {
        int remaining = data.length - offset;
        byte[] segment = new byte[Math.min(remaining, 7)];
        System.arraycopy(data, offset, segment, 0, segment.length);
        boolean last = remaining <= 7;

        CompletableFuture<CANOpenFrame> resp = conversation.expect(
            ConversationPredicates.sdoTransmitFrom(answerNodeId), timeout);
        conversation.sendToWire(createFrame(new CANOpenSDORequest(SDORequestCommand.SEGMENT_DOWNLOAD,
            new SDOSegmentDownloadRequest(toggle, last, segment))));

        resp.whenComplete((frame, error) -> {
            if (error != null) {
                receiver.completeExceptionally(error);
                return;
            }
            SDOResponse response = ((CANOpenSDOResponse) frame.getPayload()).getResponse();
            if (response instanceof SDOAbortResponse) {
                SDOAbort abort = ((SDOAbortResponse) response).getAbort();
                receiver.completeExceptionally(new CANOpenAbortException(
                    "Remote aborted segmented download", abort.getCode()));
                return;
            }
            if (!(response instanceof SDOSegmentDownloadResponse)) {
                receiver.complete(PlcResponseCode.REMOTE_ERROR);
                return;
            }
            SDOSegmentDownloadResponse seg = (SDOSegmentDownloadResponse) response;
            if (seg.getToggle() != toggle) {
                receiver.complete(PlcResponseCode.REMOTE_ERROR);
                return;
            }
            int newOffset = offset + segment.length;
            if (newOffset >= data.length) {
                receiver.complete(PlcResponseCode.OK);
            } else {
                put(receiver, !toggle, newOffset);
            }
        });
    }

    private CANOpenFrame createFrame(CANOpenPayload payload) {
        return new CANOpenFrame((short) nodeId, CANOpenService.RECEIVE_SDO, payload);
    }

}
