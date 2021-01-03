package org.apache.plc4x.java.utils.connectionpool;

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcField;

import java.util.LinkedHashSet;
import java.util.List;
import java.util.concurrent.CompletableFuture;

/**
 * TODO write comment
 *
 * @author julian
 * Created by julian on 24.02.20
 */
public class CachedReadRequest implements PlcReadRequest {

    private final CachedPlcConnection parent;
    private final PlcReadRequest innerRequest;

    public CachedReadRequest(CachedPlcConnection parent, PlcReadRequest innerRequest) {
        this.parent = parent;
        this.innerRequest = innerRequest;
    }

    @Override
    public CompletableFuture<? extends PlcReadResponse> execute() {
        // Only allowed if connection is still active
        return parent.execute(innerRequest);
    }

    @Override
    public int getNumberOfFields() {
        return innerRequest.getNumberOfFields();
    }

    @Override
    public LinkedHashSet<String> getFieldNames() {
        return innerRequest.getFieldNames();
    }

    @Override
    public PlcField getField(String s) {
        return innerRequest.getField(s);
    }

    @Override
    public List<PlcField> getFields() {
        return innerRequest.getFields();
    }
}
