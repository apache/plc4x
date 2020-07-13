package org.apache.plc4x.java.spi.internal;

import org.apache.plc4x.java.spi.ConversationContext;

import java.util.function.BooleanSupplier;

class DefaultContextHandler implements ConversationContext.ContextHandler {

    private final BooleanSupplier getDone;
    private final Runnable cancel;

    public DefaultContextHandler(BooleanSupplier getDone, Runnable cancel) {
        this.getDone = getDone;
        this.cancel = cancel;
    }

    @Override
    public boolean isDone() {
        return this.getDone.getAsBoolean();
    }

    @Override
    public void cancel() {
        this.cancel.run();
    }
}
