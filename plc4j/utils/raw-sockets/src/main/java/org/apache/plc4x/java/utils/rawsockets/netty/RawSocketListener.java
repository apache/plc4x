package org.apache.plc4x.java.utils.rawsockets.netty;

public interface RawSocketListener {

    void packetReceived(byte[] rawData);

}
