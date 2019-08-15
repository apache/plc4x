package org.apache.plc4x.logstash;

import co.elastic.logstash.api.Configuration;
import co.elastic.logstash.api.Context;
import co.elastic.logstash.api.Input;
import co.elastic.logstash.api.LogstashPlugin;
import co.elastic.logstash.api.PluginConfigSpec;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;

import java.util.Arrays;
import java.util.Collection;
import java.util.Collections;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CountDownLatch;
import java.util.function.Consumer;

// class name must match plugin name
@LogstashPlugin(name="plc4x_input")
public class Plc4xInput implements Input {

    public static final PluginConfigSpec<Map<String, Object>> FIELDS_CONFIG =
            PluginConfigSpec.hashSetting("fields");

    public static final PluginConfigSpec<String> CONNECTION_STRING_CONFIG =
            PluginConfigSpec.requiredStringSetting("connection_string");
    private final String connectionString;
    private final Map<String, Object> fields;

    private String id;
    private PlcConnection plcConnection;

    // all plugins must provide a constructor that accepts id, Configuration, and Context
    public Plc4xInput(String id, Configuration config, Context context) {
        // constructors should validate configuration options
        this.id = id;
        fields = config.get(FIELDS_CONFIG);
        connectionString = config.get(CONNECTION_STRING_CONFIG);
    }

    @Override
    public void start(Consumer<Map<String, Object>> consumer) {

        // The start method should push Map<String, Object> instances to the supplied QueueWriter
        // instance. Those will be converted to Event instances later in the Logstash event
        // processing pipeline.
        //
        // Inputs that operate on unbounded streams of data or that poll indefinitely for new
        // events should loop indefinitely until they receive a stop request. Inputs that produce
        // a finite sequence of events should loop until that sequence is exhausted or until they
        // receive a stop request, whichever comes first.
        // Establish a connection to the plc using the url provided as first argument
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection(connectionString)) {

            // Check if this connection support reading of data.
            if (!plcConnection.getMetadata().canRead()) {
                System.err.println("This connection doesn't support reading.");
                return;
            }

            // Create a new read request:
            PlcReadRequest.Builder builder = plcConnection.readRequestBuilder();
            for (String key: fields.keySet()
                 ) {
                builder.addItem(key, fields.get(key).toString());
            }
            PlcReadRequest readRequest = builder.build();

            PlcReadResponse syncResponse = readRequest.execute().get();
        } catch (PlcConnectionException e) {
            e.printStackTrace();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    @Override
    public void stop() {
        stopped = true; // set flag to request cooperative stop of input
    }

    @Override
    public void awaitStop() throws InterruptedException {
        done.await(); // blocks until input has stopped
    }

    @Override
    public Collection<PluginConfigSpec<?>> configSchema() {
        // should return a list of all configuration options for this plugin
        return Arrays.asList(FIELDS_CONFIG, CONNECTION_STRING_CONFIG);
    }

    @Override
    public String getId() {
        return this.id;
    }
}
