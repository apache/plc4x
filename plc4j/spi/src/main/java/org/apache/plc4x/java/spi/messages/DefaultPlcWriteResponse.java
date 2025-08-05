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
package org.apache.plc4x.java.spi.messages;

import java.util.Collections;
import java.util.Map.Entry;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.metadata.Metadata;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.util.Collection;
import java.util.Map;

import static org.apache.plc4x.java.spi.generation.WithReaderWriterArgs.WithAdditionalStringRepresentation;
import static org.apache.plc4x.java.spi.generation.WithReaderWriterArgs.WithRenderAsList;

public class DefaultPlcWriteResponse implements PlcWriteResponse, Serializable {

    private final PlcWriteRequest request;
    private final Map<String, PlcResponseCode> responseCodes;
    private final Map<String, Metadata> metadata;

    public DefaultPlcWriteResponse(PlcWriteRequest request,
                                   Map<String, PlcResponseCode> responseCodes) {
        this(request, responseCodes, Collections.emptyMap());
    }

    public DefaultPlcWriteResponse(PlcWriteRequest request,
                                   Map<String, PlcResponseCode> responseCodes,
                                   Map<String, Metadata> metadata) {
        this.request = request;
        this.responseCodes = responseCodes;
        this.metadata = metadata;
    }

    @Override
    public PlcWriteRequest getRequest() {
        return request;
    }

    @Override
    public Metadata getTagMetadata(String tag) {
        return metadata.getOrDefault(tag, Metadata.EMPTY);
    }

    @Override
    public Collection<String> getTagNames() {
        return request.getTagNames();
    }

    @Override
    public PlcTag getTag(String name) {
        return request.getTag(name);
    }

    @Override
    public PlcResponseCode getResponseCode(String name) {
        return responseCodes.get(name);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcWriteResponse");

        writeBuffer.pushContext("request");
        if (request instanceof Serializable) {
            ((Serializable) request).serialize(writeBuffer);
        }
        writeBuffer.popContext("request");

        writeBuffer.pushContext("responseCodes", WithRenderAsList(true));
        for (Map.Entry<String, PlcResponseCode> tagEntry : responseCodes.entrySet()) {
            String tagName = tagEntry.getKey();
            writeBuffer.pushContext(tagName);
            final PlcResponseCode tagResponseCode = tagEntry.getValue();
            writeBuffer.writeUnsignedByte("ResponseCode", 8, (byte) tagResponseCode.getValue(), WithAdditionalStringRepresentation(tagResponseCode.name()));
            writeBuffer.popContext(tagName);
        }
        writeBuffer.popContext("responseCodes");

        if (metadata != null && !metadata.isEmpty()) {
            writeBuffer.pushContext("metadata", WithRenderAsList(true));

            for (Entry<String, Metadata> entry : metadata.entrySet()) {
                if (entry.getValue() instanceof Serializable) {
                    writeBuffer.pushContext(entry.getKey());
                    ((Serializable) entry.getValue()).serialize(writeBuffer);
                    writeBuffer.popContext(entry.getKey());
                }
            }

            writeBuffer.popContext("metadata");
        }

        writeBuffer.popContext("PlcWriteResponse");
    }

}
