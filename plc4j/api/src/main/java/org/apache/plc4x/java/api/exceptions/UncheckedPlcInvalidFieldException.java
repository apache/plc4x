package org.apache.plc4x.java.api.exceptions;

/**
 * Wrapper for {@link PlcInvalidFieldException} as this is used in lambdas and stream apis.
 *
 * @see java.io.UncheckedIOException
 */
public class UncheckedPlcInvalidFieldException extends RuntimeException {

    private PlcInvalidFieldException wrappedException;

    public UncheckedPlcInvalidFieldException(PlcInvalidFieldException wrappedException) {
        this.wrappedException = wrappedException;
    }

    public PlcInvalidFieldException getWrappedException() {
        return wrappedException;
    }
}
