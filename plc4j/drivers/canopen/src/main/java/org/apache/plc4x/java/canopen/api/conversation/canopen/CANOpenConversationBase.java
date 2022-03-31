/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.canopen.api.conversation.canopen;

import io.vavr.control.Either;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.canopen.readwrite.CANOpenFrame;
import org.apache.plc4x.java.canopen.readwrite.*;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;

import java.util.concurrent.CompletableFuture;
import java.util.function.Predicate;

public abstract class CANOpenConversationBase {

    protected final CANConversation delegate;
    protected final int nodeId;
    protected final int answerNodeId;

    public CANOpenConversationBase(CANConversation delegate, int nodeId, int answerNodeId) {
        this.delegate = delegate;
        this.nodeId = nodeId;
        this.answerNodeId = answerNodeId;
    }

    protected PlcValue decodeFrom(byte[] data, CANOpenDataType type, int length) throws ParseException {
        return DataItem.staticParse(new ReadBufferByteBased(data, ByteOrder.LITTLE_ENDIAN), type, length);
    }

    protected <T> void onError(CompletableFuture<T> receiver, CANOpenSDOResponse response, Throwable error) {
        if (error != null) {
            receiver.completeExceptionally(error);
            return;
        }

        if (response.getResponse() instanceof SDOAbortResponse) {
            SDOAbortResponse abort = (SDOAbortResponse) response.getResponse();
            SDOAbort sdoAbort = abort.getAbort();
            receiver.completeExceptionally(new PlcException("Could not read value. Remote party reported code " + sdoAbort.getCode()));
        }
    }

    protected <X extends SDOResponse> Either<SDOAbort, X> unwrap(Class<X> payload, SDOResponse response) {
        if (response instanceof SDOAbortResponse) {
            return Either.left(((SDOAbortResponse) response).getAbort());
        }
        if (payload.isInstance(response)) {
            return Either.right((X) response);
        }
        throw new RuntimeException("Unexpected payload kind " + response);
    }

    protected CANOpenFrame createFrame(SDORequest rq) {
        return new CANOpenFrame((short) nodeId, CANOpenService.RECEIVE_SDO, new CANOpenSDORequest(rq.getCommand(), rq));
    }

    static class NodeIdPredicate implements Predicate<CANOpenFrame> {

        private final int nodeId;

        NodeIdPredicate(int nodeId) {
            this.nodeId = nodeId;
        }

        @Override
        public boolean test(CANOpenFrame frame) {
            return frame.getNodeId() == nodeId && frame.getService() == CANOpenService.TRANSMIT_SDO;
        }

        @Override
        public String toString() {
            return "NodeIdPredicate [" + nodeId + "]";
        }
    }

    static class TypePredicate<T, X> implements Predicate<X> {

        private final Class<T> type;

        public TypePredicate(Class<T> type) {
            this.type = type;
        }

        @Override
        public boolean test(X value) {
            return type.isInstance(value);
        }

        public String toString() {
            return "Type [" + type + "]";
        }
    }

    static class TypeOrAbortPredicate<T extends SDOResponse> extends TypePredicate<T, SDOResponse> {

        public TypeOrAbortPredicate(Class<T> type) {
            super(type);
        }

        @Override
        public boolean test(SDOResponse response) {
            return super.test(response) || response instanceof SDOAbortResponse;
        }
    }
}
