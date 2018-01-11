package org.apache.plc4x.java.api.messages.items;

import org.apache.plc4x.java.api.types.ResponseCode;

public abstract class ResponseItem<REQUEST_TYPE> {

    private final REQUEST_TYPE requestItem;

    private final ResponseCode responseCode;

    public ResponseItem(REQUEST_TYPE requestItem, ResponseCode responseCode) {
        this.requestItem = requestItem;
        this.responseCode = responseCode;
    }

    public REQUEST_TYPE getRequestItem() {
        return requestItem;
    }

    public ResponseCode getResponseCode() {
        return responseCode;
    }
}
