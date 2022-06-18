package org.apache.plc4x.java.transport.serial;

import io.netty.buffer.ByteBuf;

import java.net.SocketAddress;

public class DummyHandler extends SerialChannelHandler {

    public static final DummyHandler INSTANCE = new DummyHandler(null);

    private SerialSelectionKey selectionKey;

    public DummyHandler(SocketAddress address) {
        super(address);
    }

    @Override
    public boolean open() {
        return true;
    }

    @Override
    public String getIdentifier() {
        return null;
    }

    @Override
    public void registerSelectionKey(SerialSelectionKey selectionKey) {
        this.selectionKey = selectionKey;
    }

    @Override
    public void close() {
        // NOOP
    }

    @Override
    public int read(ByteBuf buf) {
        buf.writeByte(1);
        return 1;
    }

    @Override
    public int write(ByteBuf buf) {
        System.out.println("Haha i wrote something");
        return 1;
    }

    public void fireEvent(int readyOp) {
        ((SerialPollingSelector) this.selectionKey.selector())
            .addEvent(new SerialPollingSelector.SelectorEvent(this.selectionKey, readyOp));
    }
}
