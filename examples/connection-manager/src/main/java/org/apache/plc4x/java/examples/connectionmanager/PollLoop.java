package org.apache.plc4x.java.examples.connectionmanager;

import static org.apache.plc4x.java.examples.connectionmanager.ConnectionManager.PLC4JTYPE_SIEMENS;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.logging.Level;
import java.util.logging.Logger;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;

public class PollLoop {

    public static final Logger logger = Logger.getLogger("PollLoop.class");

    public static PlcConnection plcConnection;
    public static ConnectionManager connectionManager;

    public static void main(String[] args) throws Exception {
        
        Collector collector;

        List<String> variables = new ArrayList<>();
        variables.add("%M89:REAL"); // currentSpeedInRpm	DriveVariables	Real	%MD89
        variables.add("%Q20:REAL"); // temperatureInCelsius	SensorVariables	Real	%QD20
        variables.add("%I58.0:BOOL"); // switchingStateOfCapSensor SensorVariable Bool %I58.0
        variables.add("%Q25:WORD"); // distanceInMm	SensorVariables	Word	%QW25
        variables.add("%Q82:REAL"); // currentDrivePercent	DriveVariables	Real	%QD82
        variables.add("%M86:INT");  // driveSetFreqInPercent	DriveVariables	Int	%MW86

        variables.add("%I66:WORD"); // rawValueOfTempSensor	SensorVariables	Word	%IW66

        variables.add("%I58:WORD"); // capSensorWord SensorVariables	Word	%IW58

        connectionManager = new ConnectionManager("s7://192.168.100.49/0/1",
            PLC4JTYPE_SIEMENS);
        connectionManager.start();

        collector = new Collector(variables, Integer.parseInt("1000"));
        collector.start();

        try {
            Thread.sleep(100000);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }

    public static class Collector implements Runnable {

        private List<String> variables;
        private int samplingRate;

        private final AtomicBoolean running = new AtomicBoolean(false);

        private Thread worker;

        int incrementalSleepTime = 0;

        public Collector(List<String> variables, int samplingRate) {
            this.variables = variables;
            this.samplingRate = samplingRate;
        }

        @Override
        public void run() {
            running.set(true);
            int timeout = samplingRate;

            while (running.get()) {
                // if connection is not initialized at the beginning
                try {
                    plcConnection = connectionManager.getPlcConnection();
                } catch (PlcConnectionException e) {
                    sleep(timeout);
                    incrementalSleep();
                    continue;
                }

                // Create a new read request:
                // Give the single item requested the alias name "value"
                PlcReadRequest.Builder builder = plcConnection.readRequestBuilder();
                for (int i = 0; i < variables.size(); i++) {
                    builder.addItem(variables.get(i), variables.get(i));
                }

                PlcReadRequest readRequest = builder.build();

                // Read synchronously ...
                logger.log(Level.FINEST, "Synchronous request ...");
                PlcReadResponse syncResponse = null;
                try {
                    // fetch the response, thread pause until 1s is elapsed
                    syncResponse = readRequest.execute().get(timeout, TimeUnit.MILLISECONDS);
                    incrementalSleepTime = 250;
                } catch (ExecutionException e) {
                    logger.log(Level.SEVERE, "Error getting response", e);
                    incrementalSleep();
                    continue;
                } catch (InterruptedException e) {
                    logger.log(Level.WARNING, "Thread was interrupted" + e);
                    Thread.currentThread().interrupt();
                    // when a timeout occurs we want to stop the loop and start from beginning, 1s has already elapsed
                } catch (TimeoutException e) {
                    logger.log(Level.SEVERE, "Timeout Exception", e);
                    incrementalSleep();
                    continue;
                }

                Object[] event = response2Event(syncResponse, variables);
                String logoutput = Arrays.toString(event);
                logger.log(Level.FINEST, logoutput);
                Object[] events = response2Event(syncResponse, variables);
                System.out.println(Arrays.toString(events));
                // when sucessfull wait 1s to fetch next result
                sleep(timeout);
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

        private void incrementalSleep() {
            if (incrementalSleepTime < 60000) {
                incrementalSleepTime += 250;
            }
            try {
                Thread.sleep(incrementalSleepTime);
            } catch (InterruptedException e1) {
                logger.log(Level.WARNING, "Thread was interrupted" + e1);
                Thread.currentThread().interrupt();
            }
        }

        public void start() {
            worker = new Thread(this, "plc4j-collector");
            worker.start();
        }

        public void stop() {
            running.set(false);
        }
    }

    public static Object[] response2Event(PlcReadResponse response, List<String> variables) {
        // field names are returned in sorted order we do not want that
//    List<String> fieldNames = new ArrayList<>(response.getFieldNames());
        List<String> fieldNames = variables;
        Object[] event = new Object[fieldNames.size() + 1];

        event[0] = System.currentTimeMillis();

        for (int i = 0; i < fieldNames.size(); i++) {
            if (response.getResponseCode(fieldNames.get(i)) == PlcResponseCode.OK) {
                Object value = response.getObject(fieldNames.get(i));
                value = convertBoolean(value);
                event[i + 1] = value.toString();
            }

            // Something went wrong, to output an error message instead.
            else {
                System.out.println(
                    "Error[" + fieldNames.get(i) + "]: " + response
                        .getResponseCode(fieldNames.get(i))
                        .name());
            }
        }
        return event;
    }

    /**
     * converts an array of boolean into  a more readable "0101" String
     */
    public static Object convertBoolean(Object input) {
        StringBuilder stringBuilder = new StringBuilder();

        if (input instanceof ArrayList) {
            ArrayList<Boolean> list = (ArrayList<Boolean>) input;

            for (Boolean boolElem : list) {
                String tmp2;
                boolean tmp = Boolean.parseBoolean(boolElem.toString());
                tmp2 = tmp ? "1" : "0";

                stringBuilder.append(tmp2);
            }
            return stringBuilder.toString();
        } else {
            return input;
        }
    }
}
