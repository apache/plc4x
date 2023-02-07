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
package org.apache.plc4x.java.canopen.api.conversation.canopen;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.canopen.transport.CANOpenAbortException;
import org.apache.plc4x.java.canopen.readwrite.*;
import org.apache.plc4x.java.spi.generation.ByteOrder;

import java.util.concurrent.CompletableFuture;

import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;

public class SDODownloadConversation extends CANOpenConversationBase {

    private final CANConversation delegate;
    private final IndexAddress indexAddress;
    private final byte[] data;

    public SDODownloadConversation(CANConversation delegate, int nodeId, int answerNodeId, IndexAddress indexAddress, PlcValue value, CANOpenDataType type) {
        super(delegate, nodeId, answerNodeId);
        this.delegate = delegate;
        this.indexAddress = indexAddress;

        try {
            WriteBufferByteBased writeBuffer = new WriteBufferByteBased(DataItem.getLengthInBytes(value, type, null), ByteOrder.LITTLE_ENDIAN);
            DataItem.staticSerialize(writeBuffer, value, type, null, ByteOrder.LITTLE_ENDIAN);
            data = writeBuffer.getBytes();
        } catch (SerializationException e) {
            throw new PlcRuntimeException("Could not serialize data", e);
        }
    }

    public void execute(CompletableFuture<PlcResponseCode> receiver) {
        if (data.length > 4) {
            // segmented
            SDOInitiateSegmentedUploadResponse size = new SDOInitiateSegmentedUploadResponse(data.length);
            delegate.send(createFrame(new SDOInitiateDownloadRequest(false, true, indexAddress, size)))
                .check(new NodeIdPredicate(answerNodeId))
                .onTimeout(receiver::completeExceptionally)
                .onError((response, error) -> receiver.completeExceptionally(error))
                .unwrap(CANOpenFrame::getPayload)
                .only(CANOpenSDOResponse.class)
                .unwrap(CANOpenSDOResponse::getResponse)
                .check(new TypeOrAbortPredicate<>(SDOInitiateDownloadResponse.class))
                .unwrap(payload -> unwrap(SDOInitiateDownloadResponse.class, payload))
                .handle(either -> {
                    if (either.isLeft()) {
                        receiver.completeExceptionally(new CANOpenAbortException("Could not initiate upload", either.getLeft().getCode()));
                    } else {
                        SDOInitiateDownloadResponse response = either.get();
                        if (response.getAddress().equals(indexAddress)) {
                            put(data, receiver, false, 0);
                        } else {
                            // TODO find proper error code in spec
                            SDOAbort abort = new SDOAbort(indexAddress, 1000);
                            delegate.sendToWire(createFrame(new SDOAbortRequest(abort)));
                            receiver.complete(PlcResponseCode.REMOTE_ERROR);
                        }
                    }
                });

            return;
        }

        // expedited
        SDOInitiateDownloadRequest rq = new SDOInitiateDownloadRequest(
            true, true,
            indexAddress,
            new SDOInitiateExpeditedUploadResponse(data)
        );

        delegate.send(createFrame(rq))
            .check(new NodeIdPredicate(answerNodeId))
            .onTimeout(receiver::completeExceptionally)
            .unwrap(CANOpenFrame::getPayload)
            .only(CANOpenSDOResponse.class)
            .onError((response, error) -> onError(receiver, response, error))
            .unwrap(CANOpenSDOResponse::getResponse)
            .check(new TypeOrAbortPredicate<>(SDOInitiateDownloadResponse.class))
            .unwrap(payload -> unwrap(SDOInitiateDownloadResponse.class, payload))
            .handle(either -> {
                if (either.isLeft()) {
                    receiver.completeExceptionally(new CANOpenAbortException("Could not initiate upload", either.getLeft().getCode()));
                } else {
                    SDOInitiateDownloadResponse response = either.get();
                    if (response.getCommand() == SDOResponseCommand.INITIATE_DOWNLOAD) {
                        receiver.complete(PlcResponseCode.OK);
                    } else {
                        receiver.complete(PlcResponseCode.REMOTE_ERROR);
                    }
                }
            });
    }

    private void put(byte[] data, CompletableFuture<PlcResponseCode> receiver, boolean toggle, int offset) {
        int remaining = data.length - offset;
        byte[] segment = new byte[Math.min(remaining, 7)];
        System.arraycopy(data, offset, segment, 0, segment.length);

        delegate.send(createFrame(new SDOSegmentDownloadRequest(toggle, remaining <= 7, segment)))
            .check(new NodeIdPredicate(answerNodeId))
            .onTimeout(receiver::completeExceptionally)
            .unwrap(CANOpenFrame::getPayload)
            .only(CANOpenSDOResponse.class)
            .onError((response, error) -> onError(receiver, response, error))
            .unwrap(CANOpenSDOResponse::getResponse)
            .check(new TypeOrAbortPredicate<>(SDOSegmentDownloadResponse.class))
            .unwrap(payload -> unwrap(SDOSegmentDownloadResponse.class, payload))
            .handle(either -> {
                if (either.isLeft()) {
                    return;
                } else {
                    SDOSegmentDownloadResponse response = either.get();
                    if (response.getToggle() != toggle) {
                        receiver.complete(PlcResponseCode.REMOTE_ERROR);
                        return;
                    }

                    if (offset + segment.length == data.length) {
                        receiver.complete(PlcResponseCode.OK);
                    } else {
                        put(data, receiver, !toggle, offset + segment.length);
                    }
                }
            });
    }

}
