package org.apache.plc4x.java.utils.rawsockets.netty;

public class RawSocketException extends Exception {

    public RawSocketException() {
    }

    public RawSocketException(String message) {
        super(message);
    }

    public RawSocketException(String message, Throwable cause) {
        super(message, cause);
    }

}
