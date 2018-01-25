package org.apache.plc4x.java.utils.rawsockets;

public interface RawSocketListener {

    void packetReceived(byte[] rawData);

}
