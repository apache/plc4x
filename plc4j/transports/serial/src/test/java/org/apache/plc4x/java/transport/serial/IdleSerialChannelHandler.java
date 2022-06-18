package org.apache.plc4x.java.transport.serial;

import io.netty.buffer.ByteBuf;

class IdleSerialChannelHandler extends SerialChannelHandler {

    public IdleSerialChannelHandler() {
        super(null);
    }

    @Override
    public boolean open() {
        return true;
    }

    @Override
    String getIdentifier() {
        return "Emulated Com Port";
    }

    @Override
    void registerSelectionKey(SerialSelectionKey selectionKey) {
        /*
         * We can ignore registration as we will neven trigger events
         */
    }

    @Override
    public void close() {
    }

    @Override
    public int read(ByteBuf buf) {
        return 0;
    }

    @Override
    public int write(ByteBuf buf) {
        throw new UnsupportedOperationException();
    }
}
