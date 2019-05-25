package org.apache.plc4x.plugins.codegenerator.example;

import org.apache.plc4x.codegen.api.Buffer;

public class TPKTPacketFactory {

    public TPKTPacket read(Buffer buffer) {
        assert buffer.readUint8() == TPKTPacket.PROTOCOL_ID;
        assert buffer.readUint8() == TPKTPacket.RESERVED_1;
        int len = buffer.readUint16();
        final Context context = new Context();
        context.put("payloadLength", len - 4);
        COTPPacket cotpPacket = COTPPacket.read(buffer, context);
        return new TPKTPacket(cotpPacket);
    }

}
