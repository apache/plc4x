package org.apache.plc4x.java.base.connection;

import org.apache.commons.lang3.NotImplementedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.channels.SelectableChannel;
import java.nio.channels.SelectionKey;
import java.nio.channels.Selector;
import java.nio.channels.spi.AbstractSelectionKey;

class SerialSelectionKey extends AbstractSelectionKey {

    private static final Logger logger = LoggerFactory.getLogger(SerialSelectionKey.class);

    final SelectableChannel channel;
    final Selector selector;
    int index;
    private volatile int interestOps;
    private int readyOps;

    SerialSelectionKey(SelectableChannel channel, Selector selector, int interestOps) {
        this.channel = channel;
        this.selector = selector;
        this.interestOps = interestOps;
    }

    public void setIndex(int index) {
        this.index = index;
    }

    @Override public SelectableChannel channel() {
        return this.channel;
    }

    @Override public Selector selector() {
        return this.selector;
    }

    @Override public int interestOps() {
        return this.interestOps;
    }

    @Override public SelectionKey interestOps(int ops) {
        this.interestOps = ops;
        return this;
    }

    @Override public int readyOps() {
        final int ops = this.readyOps;
        this.readyOps = 0;
        // Reset events for this here in Selector
        ((SerialPollingSelector) selector).removeEvent(this);
        logger.debug("Returning ready operation {}", ops);
        return ops;
    }

    public void addReadyOp(int event) {
        readyOps = readyOps | event;
        logger.debug("Adding event {} to ready ops, now having ready ops {}", event, readyOps);
    }
}
