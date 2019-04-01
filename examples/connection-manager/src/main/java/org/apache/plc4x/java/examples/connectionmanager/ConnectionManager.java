package org.apache.plc4x.java.examples.connectionmanager;

import java.net.InetSocketAddress;
import java.net.Socket;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.logging.Level;
import java.util.logging.Logger;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.s7.S7PlcDriver;

public class ConnectionManager implements Runnable {

    public static final String PLC4JTYPE_SIEMENS = "Siemens S7";

    private final AtomicBoolean running = new AtomicBoolean(false);
    private final Logger logger = Logger.getLogger(this.getClass().getName());
    private Thread worker;

    private PlcConnection plcConnection;

    private String addressString;
    private String plcType;
    private boolean isAlive;

    // checks for the "liveness" of the connection and restarts connection after timeout of 1 minute
    public ConnectionManager(String addressString, String plcType) {
        this.addressString = addressString;
        this.plcType = plcType;
    }

    @Override
    public void run() {
        running.set(true);
        while (running.get()) {
            try {
                // initialize for the first time
                if (plcConnection == null) {
                    logger.log(Level.FINEST, "driver is not initialized");
                    plcConnection = initPLC(addressString);
                    isAlive = true;
                }

                // String host = connection2host(addressString);
                // boolean ping = ping(host, 102, 1000);
                boolean ping = channelPingCheck(1000, "%M1.2:BOOL");
                if (!ping) {
                    logger.log(Level.FINEST, "simple check of connection failed");
                    isAlive = false;
                    // will throw exception after timeout
                    plcConnection = initPLC(addressString);
                    isAlive = true;
                    // just to be shure we restart the collector

                }
            } catch (PlcConnectionException e) {
                logger.log(Level.WARNING, "error connecting with driver", e);
            } catch (Exception e) {
                logger.log(Level.FINEST, "Error in Connection Manager", e);
            }
            // wait 1 minute
            sleep(60000);
        }
    }

    private void sleep(int sleeptime) {
        try {
            Thread.sleep(sleeptime);
        } catch (InterruptedException e1) {
            logger.log(Level.WARNING, "Thread was interrupted", e1);
            Thread.currentThread().interrupt();
        }
    }

    void start() {
        worker = new Thread(this, "plc4j-connectionManager");
        worker.start();
    }

    void stop() {
        running.set(false);
    }

    PlcConnection getPlcConnection() throws PlcConnectionException {
        if (plcConnection == null) {
            logger.log(Level.FINEST, "driver is not initialized");
            throw new PlcConnectionException("driver not initialized...please wait");
        }
        if (!isAlive) {
            logger.log(Level.FINEST, "reconnecting driver");
            throw new PlcConnectionException("driver is not initialized...reconnecting");
        }

        return plcConnection;
    }

    private PlcConnection initPLC(String addressString) throws PlcConnectionException {
        String driverName = plcType;
        PlcConnection plcConnection = null;

        switch (driverName) {
            case PLC4JTYPE_SIEMENS:
                plcConnection = new PlcDriverManager(S7PlcDriver.class.getClassLoader())
                    .getConnection(addressString);
                break;
            default:
                plcConnection = new PlcDriverManager(S7PlcDriver.class.getClassLoader())
                    .getConnection(addressString);
                break;
        }

        // can we ommit this?
        if (plcConnection == null) {
            throw new PlcConnectionException("error initializing driver");
        }

        // Check if this connection support reading of data.
        if (!plcConnection.getMetadata().canRead()) {
            logger.log(Level.SEVERE, "This connection doesn't support reading.");
            throw new IllegalAccessError("cant read from Driver");
        }
        return plcConnection;
    }

    private boolean channelPingCheck(int timeout) {
        String variable = "%M1.2:BOOL";
        return channelPingCheck(timeout, variable);
    }

    private boolean channelPingCheck(int timeout, String variable) {

//      String variable = "%M1.2:BOOL";

//      boolean expectedResult = true;

        PlcReadRequest.Builder builder = plcConnection.readRequestBuilder();

        builder.addItem(variable, variable);

        PlcReadRequest readRequest = builder.build();

        PlcReadResponse result = null;
        try {
            result = readRequest.execute().get(timeout, TimeUnit.MILLISECONDS);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        } catch (ExecutionException | TimeoutException e) {
            return false;
        }

        if (result == null) {
            return false;
        }

        Object content;
        content = result.getObject(variable);
        if (content == null) {
            return false;
        }
        //we could compare against the real value of the object here, but then we need to be more specific about the variable type
        return true;

    }

    public String connection2host(String input) {
        String host = null;
        Pattern S7_URI_PATTERN = Pattern
            .compile("^s7://(?<host>.*)/(?<rack>\\d{1,4})/(?<slot>\\d{1,4})(?<params>\\?.*)?");

        Matcher matcher = S7_URI_PATTERN.matcher(input);
        if (!matcher.matches()) {
            logger.log(Level.SEVERE,
                "Connection url doesn't match the format 's7://{host|ip}/{rack}/{slot}'");
        } else {
            host = matcher.group("host");
        }
        return host;
    }

    public boolean ping(String host, int port, int timeout) {
        Socket s = null;
        try {
            s = new Socket();
            s.connect(new InetSocketAddress(host, port), timeout);
            return true;
        } catch (Exception e) {
            return false;
        } finally {
            if (s != null) {
                try {
                    s.close();
                } catch (Exception e) {
                }
            }
        }
    }

}
