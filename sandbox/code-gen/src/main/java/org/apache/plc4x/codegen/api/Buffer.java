package org.apache.plc4x.codegen.api;

/**
 * The Interface as described in {@link org.apache.plc4x.codegen.util.BufferUtil}.
 */
public interface Buffer {

    Integer readUint8();

    Integer readUint16();

    Long readUint32();
}
