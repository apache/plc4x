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
package org.apache.plc4x.java.canopen.api.segmentation.accumulator;

import org.apache.plc4x.java.canopen.readwrite.*;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.function.Function;

public class ByteStorage<T> implements Storage<T, byte[]> {

    private final List<byte[]> segments = new ArrayList<>();
    private final Function<T, byte[]> extractor;
    private long size = 0;

    public ByteStorage(Function<T, byte[]> extractor) {
        this.extractor = extractor;
    }

    @Override
    public void append(T frame) {
        segments.add(extractor.apply(frame));
        size += segments.get(segments.size() - 1).length;
    }

    public long size() {
        return size;
    }

    @Override
    public byte[] get() {
        Optional<byte[]> collect = segments.stream().reduce((b1, b2) -> {
            byte[] combined = new byte[b1.length + b2.length];
            System.arraycopy(b1, 0, combined, 0, b1.length);
            System.arraycopy(b2, 0, combined, b1.length, b2.length);
            return combined;
        });
        return collect.orElse(new byte[0]);
    }

    public static class SDOUploadStorage extends ByteStorage<SDOResponse> {
        public SDOUploadStorage() {
            super((sdoResponse -> {
                if (sdoResponse instanceof SDOSegmentUploadResponse) {
                    return ((SDOSegmentUploadResponse) sdoResponse).getData();
                }
                if (sdoResponse instanceof SDOInitiateUploadResponse) {
                    SDOInitiateUploadResponse initiate = (SDOInitiateUploadResponse) sdoResponse;

                    if (initiate.getPayload() instanceof SDOInitiateExpeditedUploadResponse) {
                        return ((SDOInitiateExpeditedUploadResponse) initiate.getPayload()).getData();
                    }
                }
                return new byte[0];
            }));
        }
    }

    public static class SDODownloadStorage extends ByteStorage<SDORequest> {
        public SDODownloadStorage() {
            super((sdoRequest -> {
                if (sdoRequest instanceof SDOSegmentDownloadRequest) {
                    return ((SDOSegmentDownloadRequest) sdoRequest).getData();
                }
                if (sdoRequest instanceof  SDOInitiateDownloadRequest) {
                    SDOInitiateDownloadRequest initiate = (SDOInitiateDownloadRequest) sdoRequest;

                    if (initiate.getPayload() instanceof SDOInitiateExpeditedUploadResponse) {
                        return ((SDOInitiateExpeditedUploadResponse) initiate.getPayload()).getData();
                    }
                }
                return new byte[0];
            }));
        }
    }

}
