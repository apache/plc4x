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
package org.apache.plc4x.nifi.record;

import java.io.IOException;
import java.io.OutputStream;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.atomic.AtomicReference;

import org.apache.nifi.flowfile.FlowFile;
import org.apache.nifi.flowfile.attributes.CoreAttributes;
import org.apache.nifi.logging.ComponentLog;
import org.apache.nifi.processor.ProcessSession;
import org.apache.nifi.serialization.RecordSetWriter;
import org.apache.nifi.serialization.RecordSetWriterFactory;
import org.apache.nifi.serialization.WriteResult;
import org.apache.nifi.serialization.record.Record;
import org.apache.nifi.serialization.record.RecordSchema;
import org.apache.nifi.serialization.record.RecordSet;
import org.apache.plc4x.java.api.messages.PlcReadResponse;

public class RecordPlc4xWriter implements Plc4xWriter {

    private final RecordSetWriterFactory recordSetWriterFactory;
    private final AtomicReference<WriteResult> writeResultRef;
    private final Map<String, String> originalAttributes;
    private String mimeType;

    private RecordSet fullRecordSet;
    private RecordSchema writeSchema;


    public RecordPlc4xWriter(RecordSetWriterFactory recordSetWriterFactory, Map<String, String> originalAttributes) {
        this.recordSetWriterFactory = recordSetWriterFactory;
        this.writeResultRef = new AtomicReference<>();
        this.originalAttributes = originalAttributes;
    }

    @Override
    public long writePlcReadResponse(PlcReadResponse response, OutputStream outputStream, ComponentLog logger, 
                Plc4xReadResponseRowCallback callback, RecordSchema recordSchema, String timestampFieldName) throws Exception {
        
                    if (fullRecordSet == null) {
            fullRecordSet = new Plc4xReadResponseRecordSetWithCallback(response, callback, recordSchema, timestampFieldName);
            writeSchema = recordSetWriterFactory.getSchema(originalAttributes, fullRecordSet.getSchema());
        }
        Map<String, String> empty = new HashMap<>();
        try (final RecordSetWriter resultSetWriter = recordSetWriterFactory.createWriter(logger, writeSchema, outputStream, empty)) {
            writeResultRef.set(resultSetWriter.write(fullRecordSet));
            if (mimeType == null) {
                mimeType = resultSetWriter.getMimeType();
            }
            return writeResultRef.get().getRecordCount();
        } catch (final Exception e) {
            throw new IOException(e);
        }
    }

    @Override
    public long writePlcReadResponse(PlcReadResponse response, OutputStream outputStream, ComponentLog logger, 
            Plc4xReadResponseRowCallback callback, RecordSchema recordSchema, FlowFile originalFlowFile, String timestampFieldName) throws Exception {
        
                if (fullRecordSet == null) {
            fullRecordSet = new Plc4xReadResponseRecordSetWithCallback(response, callback, recordSchema, timestampFieldName);
            writeSchema = recordSetWriterFactory.getSchema(originalAttributes, fullRecordSet.getSchema());
        }

        RecordSetWriter resultSetWriter = null;
        try {
            if (originalFlowFile != null) {
                try {
                    resultSetWriter = recordSetWriterFactory.createWriter(logger, writeSchema, outputStream, originalFlowFile);
                } catch (final Exception e) {
                    throw new IOException(e);
                }
            } else {
                resultSetWriter = recordSetWriterFactory.createWriter(logger, writeSchema, outputStream, Collections.emptyMap());
            }

            writeResultRef.set(resultSetWriter.write(fullRecordSet));
            if (mimeType == null) {
                mimeType = resultSetWriter.getMimeType();
            }
            return writeResultRef.get().getRecordCount();
        } finally {
            if (resultSetWriter != null) {
                resultSetWriter.close();
            }
        }
    }


    @Override
    public void writeEmptyPlcReadResponse(OutputStream outputStream, ComponentLog logger) throws IOException {
        Map<String, String> empty = new HashMap<>();
        try (final RecordSetWriter resultSetWriter = recordSetWriterFactory.createWriter(logger, writeSchema, outputStream, empty)) {
            mimeType = resultSetWriter.getMimeType();
            resultSetWriter.beginRecordSet();
            resultSetWriter.finishRecordSet();
        } catch (final Exception e) {
            throw new IOException(e);
        }
    }

    @Override
    public void writeEmptyPlcReadResponse(OutputStream outputStream, ComponentLog logger, FlowFile originalFlowFile) throws IOException {
        try (final RecordSetWriter resultSetWriter = recordSetWriterFactory.createWriter(logger, writeSchema, outputStream, originalFlowFile)) {
            mimeType = resultSetWriter.getMimeType();
            resultSetWriter.beginRecordSet();
            resultSetWriter.finishRecordSet();
        } catch (final Exception e) {
            throw new IOException(e);
        }
    }

    @Override
    public String getMimeType() {
        return mimeType;
    }

    @Override
    public Map<String, String> getAttributesToAdd() {
        Map<String, String> attributesToAdd = new HashMap<>();
        attributesToAdd.put(CoreAttributes.MIME_TYPE.key(), mimeType);
        // Add any attributes from the record writer (if present)
        final WriteResult result = writeResultRef.get();
        if (result != null) {
            if (result.getAttributes() != null) {
                attributesToAdd.putAll(result.getAttributes());
            }
            attributesToAdd.put("record.count", String.valueOf(result.getRecordCount()));
        }
        return attributesToAdd;
    }

    @Override
    public void updateCounters(ProcessSession session) {
        final WriteResult result = writeResultRef.get();
        if (result != null) {
            session.adjustCounter("Records Written", result.getRecordCount(), false);
        }
    }

    private static class Plc4xReadResponseRecordSetWithCallback extends Plc4xReadResponseRecordSet {
        private final Plc4xReadResponseRowCallback callback;

        public Plc4xReadResponseRecordSetWithCallback(final PlcReadResponse readResponse, Plc4xReadResponseRowCallback callback, 
                RecordSchema recordSchema, String timestampFieldName) {

            super(readResponse, recordSchema, timestampFieldName);
            this.callback = callback;
        }

        @Override
        public Record next() throws IOException {
            if (hasMoreRows()) {
                PlcReadResponse response = getReadResponse();
                final Record record = createRecord(response);
                setMoreRows(false);
                if (callback != null) {
                    callback.processRow(response);
                }
                return record;
            } else {
                return null;
            }
        }
    }

    public RecordSchema getRecordSchema() {
        try {
            return this.fullRecordSet.getSchema();
        } catch (IOException e) {
            return null;
        }
    }

}
