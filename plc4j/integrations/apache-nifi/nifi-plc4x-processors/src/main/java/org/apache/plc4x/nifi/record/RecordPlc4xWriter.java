package org.apache.plc4x.nifi.record;

import java.io.IOException;
import java.io.OutputStream;
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
	public long writePlcReadResponse(PlcReadResponse response, OutputStream outputStream, ComponentLog logger, Plc4xReadResponseRowCallback callback) throws Exception {
		if (fullRecordSet == null) {
            fullRecordSet = new Plc4xReadResponseRecordSetWithCallback(response, callback);
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
	public long writePlcReadResponse(PlcReadResponse response, OutputStream outputStream, ComponentLog logger, Plc4xReadResponseRowCallback callback, FlowFile originalFlowFile) throws Exception {
		if (fullRecordSet == null) {	
            fullRecordSet = new Plc4xReadResponseRecordSetWithCallback(response, callback);
            writeSchema = recordSetWriterFactory.getSchema(originalAttributes, fullRecordSet.getSchema());
         }
        try (final RecordSetWriter resultSetWriter = recordSetWriterFactory.createWriter(logger, writeSchema, outputStream, originalFlowFile)) {
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
        public Plc4xReadResponseRecordSetWithCallback(final PlcReadResponse readResponse, Plc4xReadResponseRowCallback callback) throws IOException {
            super(readResponse);
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

}
