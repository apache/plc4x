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

import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.canopen.api.segmentation.accumulator.ByteStorage;
import org.apache.plc4x.java.canopen.transport.CANOpenAbortException;
import org.apache.plc4x.java.canopen.readwrite.CANOpenFrame;
import org.apache.plc4x.java.canopen.readwrite.*;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;
import java.util.function.BiConsumer;

public class SDOUploadConversation extends CANOpenConversationBase {

    private final Logger logger = LoggerFactory.getLogger(SDOUploadConversation.class);
    private final IndexAddress address;
    private final CANOpenDataType type;

    public SDOUploadConversation(CANConversation delegate, int nodeId, int answerNodeId, IndexAddress address, CANOpenDataType type) {
        super(delegate, nodeId, answerNodeId);
        this.address = address;
        this.type = type;
    }

    public void execute(CompletableFuture<PlcValue> receiver) {
        SDOInitiateUploadRequest rq = new SDOInitiateUploadRequest(address);

        delegate.send(createFrame(rq))
            .check(new NodeIdPredicate(answerNodeId))
            .onTimeout(receiver::completeExceptionally)
            .unwrap(CANOpenFrame::getPayload)
            .only(CANOpenSDOResponse.class)
            .onError((payload, error) -> onError(receiver, payload, error))
            .unwrap(CANOpenSDOResponse::getResponse)
            .check(new TypeOrAbortPredicate<>(SDOInitiateUploadResponse.class))
            .unwrap(payload -> unwrap(SDOInitiateUploadResponse.class, payload))
            .check(either -> either.isLeft() || either.get().getAddress().equals(address))
            .handle(either -> {
                if (either.isLeft()) {
                    SDOAbort abort = either.getLeft();
                    receiver.completeExceptionally(new CANOpenAbortException("Could not complete operation", abort.getCode()));
                } else {
                    handle(receiver, either.get());
                }
            });
    }

    private void handle(CompletableFuture<PlcValue> receiver, SDOInitiateUploadResponse answer) {
        BiConsumer<Integer, byte[]> valueCallback = (length, bytes) -> {
            try {
                final PlcValue decodedValue = decodeFrom(bytes, type, length);
                receiver.complete(decodedValue);
            } catch (ArrayIndexOutOfBoundsException | ParseException e) {
                receiver.completeExceptionally(e);
            }
        };

        if (answer.getExpedited() && answer.getIndicated() && answer.getPayload() instanceof SDOInitiateExpeditedUploadResponse) {
            SDOInitiateExpeditedUploadResponse payload = (SDOInitiateExpeditedUploadResponse) answer.getPayload();
            valueCallback.accept(payload.getData().length, payload.getData());
        } else if (answer.getPayload() instanceof SDOInitiateSegmentedUploadResponse) {
            logger.debug("Beginning of segmented operation for address {}/{}", Integer.toHexString(address.getIndex()), Integer.toHexString(address.getSubindex()));
            ByteStorage.SDOUploadStorage storage = new ByteStorage.SDOUploadStorage();
            storage.append(answer);

            SDOInitiateSegmentedUploadResponse segment = (SDOInitiateSegmentedUploadResponse) answer.getPayload();
            fetch(storage, valueCallback, receiver, false, Long.valueOf(segment.getBytes()).intValue());
        } else {
            receiver.completeExceptionally(new PlcException("Unsupported SDO operation kind."));
        }
    }

    private void fetch(ByteStorage.SDOUploadStorage storage, BiConsumer<Integer, byte[]> valueCallback, CompletableFuture<PlcValue> receiver, boolean toggle, int size) {
        logger.info("Request next data block for address {}/{}", Integer.toHexString(address.getIndex()), Integer.toHexString(address.getSubindex()));
        delegate.send(createFrame(new SDOSegmentUploadRequest(toggle)))
            .check(new NodeIdPredicate(answerNodeId))
            .onTimeout(receiver::completeExceptionally)
            .unwrap(CANOpenFrame::getPayload)
            .only(CANOpenSDOResponse.class)
            .onError((payload, error) -> onError(receiver, payload, error))
            .unwrap(CANOpenSDOResponse::getResponse)
            .check(new TypeOrAbortPredicate<>(SDOSegmentUploadResponse.class))
            .unwrap(payload -> unwrap(SDOSegmentUploadResponse.class, payload))
            .handle(either -> {
                if (either.isLeft()) {
                    SDOAbort abort = either.getLeft();
                    receiver.completeExceptionally(new CANOpenAbortException("Could not complete operation", abort.getCode()));
                } else {
                    SDOSegmentUploadResponse response = either.get();
                    if (response.getToggle() != toggle) {
                        // TODO find proper error code in specs
                        receiver.completeExceptionally(new CANOpenAbortException("Remote operation failed", 1000));
                        SDOAbort abort = new SDOAbort(address, 1000);
                        delegate.sendToWire(createFrame(new SDOAbortRequest(abort)));
                        return;
                    }

                    storage.append(either.get());
                    if (response.getLast()) {
                        // validate size
                        logger.trace("Completed reading of data from {}/{}, collected {}, wanted {}", Integer.toHexString(address.getIndex()), Integer.toHexString(address.getSubindex()), storage.size(), size);
                        valueCallback.accept(Long.valueOf(size).intValue(), storage.get());
                    } else {
                        logger.trace("Continue reading of data from {}/{}, collected {}, wanted {}", Integer.toHexString(address.getIndex()), Integer.toHexString(address.getSubindex()), storage.size(), size);
                        fetch(storage, valueCallback, receiver, !toggle, size);
                    }
                }
            });
    }

}
