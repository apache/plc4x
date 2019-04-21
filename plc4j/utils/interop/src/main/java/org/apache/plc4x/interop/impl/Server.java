package org.apache.plc4x.interop.impl;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.interop.InteropServer;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;
import org.apache.plc4x.java.base.messages.items.DefaultLongFieldItem;
import org.apache.plc4x.java.mock.MockDevice;
import org.apache.plc4x.java.mock.PlcMockConnection;
import org.apache.thrift.server.TServer;
import org.apache.thrift.server.TSimpleServer;
import org.apache.thrift.transport.TServerSocket;
import org.apache.thrift.transport.TServerTransport;

import java.util.Collection;
import java.util.function.Consumer;

public class Server {

    public static void main(String[] args) throws PlcConnectionException {
        final PlcDriverManager driverManager = new PlcDriverManager();

        // Do some mocking
        final PlcMockConnection mockConnection = (PlcMockConnection) driverManager.getConnection("mock:a");

        mockConnection.setDevice(new MyMockDevice());

        final Handler handler = new Handler(driverManager);
        final InteropServer.Processor<Handler> processor = new InteropServer.Processor<>(handler);

        try {
            TServerTransport serverTransport = new TServerSocket(9090);
            TServer server = new TSimpleServer(new TServer.Args(serverTransport).processor(processor));

            System.out.println("Starting the simple server...");
            server.serve();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private static class MyMockDevice implements MockDevice {
        @Override public Pair<PlcResponseCode, BaseDefaultFieldItem> read(String fieldQuery) {
            return Pair.of(PlcResponseCode.OK, new DefaultLongFieldItem(100L));
        }

        @Override public PlcResponseCode write(String fieldQuery, Object value) {
            return null;
        }

        @Override public Pair<PlcResponseCode, PlcSubscriptionHandle> subscribe(String fieldQuery) {
            return null;
        }

        @Override public void unsubscribe() {

        }

        @Override public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
            return null;
        }

        @Override public void unregister(PlcConsumerRegistration registration) {

        }
    }
}
