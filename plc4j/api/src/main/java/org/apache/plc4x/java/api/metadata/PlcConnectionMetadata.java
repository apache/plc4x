package org.apache.plc4x.java.api.metadata;

/**
 * Information about connection capabilities.
 * This includes connection and driver specific metadata.
 */
public interface PlcConnectionMetadata {

    /**
     * Indicates that the connection supports reading.
     */
    boolean canRead();

    /**
     * Indicates that the connection supports writing.
     */
    boolean canWrite();

    /**
     * Indicates that the connection supports subscription.
     */
    boolean canSubscribe();

}
