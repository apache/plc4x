package org.apache.plc4x.java.profinet.device;

import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.profinet.readwrite.DceRpc_Packet;

public interface ProfinetCallable {
    void handle(DceRpc_Packet packet) throws PlcException;

    DceRpc_Packet create() throws PlcException;
}
