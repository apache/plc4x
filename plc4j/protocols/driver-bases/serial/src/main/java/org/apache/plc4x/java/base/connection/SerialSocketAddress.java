package org.apache.plc4x.java.base.connection;

import java.net.SocketAddress;
import java.util.Objects;

public class SerialSocketAddress extends SocketAddress {

    private final String identifier;

    public SerialSocketAddress(String identifier) {
        this.identifier = identifier;
    }

    public String getIdentifier() {
        return identifier;
    }

    @Override public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        SerialSocketAddress that = (SerialSocketAddress) o;
        return Objects.equals(identifier, that.identifier);
    }

    @Override public int hashCode() {
        return Objects.hash(identifier);
    }
}
