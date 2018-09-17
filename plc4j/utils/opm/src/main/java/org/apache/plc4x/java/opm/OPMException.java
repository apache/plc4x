package org.apache.plc4x.java.opm;

/**
 * Default Exception.
 */
public class OPMException extends Exception {

    public OPMException(String message) {
        super(message);
    }

    public OPMException(String message, Throwable cause) {
        super(message, cause);
    }
}
