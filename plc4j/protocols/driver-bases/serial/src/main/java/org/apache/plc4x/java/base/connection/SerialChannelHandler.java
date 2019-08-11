package org.apache.plc4x.java.base.connection;

import com.fazecast.jSerialComm.SerialPort;
import com.fazecast.jSerialComm.SerialPortDataListener;
import com.fazecast.jSerialComm.SerialPortEvent;
import io.netty.channel.jsc.JSerialCommDeviceAddress;

import java.net.SocketAddress;
import java.nio.channels.SelectionKey;

/**
 * This is a wrapper mostly for testing {@link SerialChannel}, @{@link SerialPollingSelector},
 * {@link SerialSelectionKey}, @{@link SerialSelectorProvider} and @{@link SerialSocketChannel}.
 */
public abstract class SerialChannelHandler {

    private final SocketAddress address;

    public SerialChannelHandler(SocketAddress address) {
        this.address = address;
    }

    abstract boolean open();

    abstract String getIdentifier();

    /**
     * This method registers the Callback to the SelectionKey / {@link java.nio.channels.Selector}
     * which is necessary to notify the {@link java.nio.channels.Selector} about
     * available data.
     */
    abstract void registerSelectionKey(SerialSelectionKey selectionKey);

    public abstract void close();

    public static class DummyHandler extends SerialChannelHandler {

        public static final DummyHandler INSTANCE = new DummyHandler(null);

        private SerialSelectionKey selectionKey;

        public DummyHandler(SocketAddress address) {
            super(address);
        }

        @Override public boolean open() {
            return true;
        }

        @Override public String getIdentifier() {
            return null;
        }

        @Override public void registerSelectionKey(SerialSelectionKey selectionKey) {
            this.selectionKey = selectionKey;
        }

        @Override public void close() {
            // NOOP
        }

        public void fireEvent(int readyOp) {
            ((SerialPollingSelector) this.selectionKey.selector())
                .addEvent(new SerialPollingSelector.SelectorEvent(this.selectionKey, readyOp));
        }
    }


    public static class SerialPortHandler extends SerialChannelHandler {

        private SerialPort comPort;

        public SerialPortHandler(SocketAddress address) {
            super(address);
            comPort = SerialPort.getCommPort(((SerialSocketAddress) address).getIdentifier());
        }

        @Override public boolean open() {
            return comPort.openPort();
        }

        @Override public String getIdentifier() {
            return comPort.getDescriptivePortName();
        }

        @Override public void registerSelectionKey(SerialSelectionKey selectionKey) {
            comPort.addDataListener(new SerialPortDataListener() {
                @Override
                public int getListeningEvents() {
                    return SerialPort.LISTENING_EVENT_DATA_AVAILABLE;
                }

                @Override
                public void serialEvent(SerialPortEvent event) {
                    ((SerialPollingSelector) selectionKey.selector())
                        .addEvent(new SerialPollingSelector.SelectorEvent(selectionKey, SelectionKey.OP_READ));
                }
            });
        }

        @Override public void close() {
            this.comPort.closePort();
        }
    }
}
