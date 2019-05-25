package org.apache.plc4x.plugins.codegenerator.example;

import org.apache.plc4x.codegen.api.Buffer;

/**
 * Example... for
 * <code>
 *   [type 'TPKTPacket'
 *     [const    uint8      'protocolId' '0x3']
 *     [reserved uint8      '0x0']
 *     [field    uint16     'len']
 *     [field    COTPPacket 'payload'  {payloadLength: 'len - 4'}]
 * ]
 * </code>
 */
public class TPKTPacket {

    public static final short PROTOCOL_ID = 0x3;
    public static final short RESERVED_1 = 0x0;

    private final COTPPacket payload;

    public TPKTPacket(COTPPacket payload) {
        this.payload = payload;
    }

    public void write(Buffer buffer) {
        buffer.writeUint8(PROTOCOL_ID);
        buffer.writeUint8(RESERVED_1);
        buffer.writeUint16(payload.getSize() + 4);
        payload.write(buffer);
    }

}
