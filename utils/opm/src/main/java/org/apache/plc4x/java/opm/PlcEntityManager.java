package org.apache.plc4x.java.opm;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidAddressException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.api.messages.items.ReadResponseItem;

import java.lang.reflect.Field;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

/**
 * Manages Entities.
 */
public class PlcEntityManager {

    private final PlcDriverManager driverManager;

    public PlcEntityManager() {
        this.driverManager = new PlcDriverManager();
    }

    public PlcEntityManager(PlcDriverManager driverManager) {
        this.driverManager = driverManager;
    }

    public <T> T find(Class<T> clazz) throws OPMException {
        PlcEntity annotation = clazz.getAnnotation(PlcEntity.class);
        if (annotation == null) {
            throw new IllegalArgumentException("Given Class is no Plc Entity, i.e., not annotated with @PlcEntity");
        }
        // Check if default constructor exists
        try {
            clazz.getConstructor();
        } catch (NoSuchMethodException e) {
            throw new IllegalArgumentException("Cannot use PlcEntity without default constructor");
        }
        String source = annotation.value();

        PlcReader reader;

        try (PlcConnection connection = driverManager.getConnection(source)) {

            if (!connection.getReader().isPresent()) {
                throw new OPMException("Unable to get Reader for connection with url '" + source + "'");
            }

            reader = connection.getReader().get();

            // Do the necessary queries for all fields
            HashMap<ReadRequestItem<?>, Field> requestItems = new HashMap<>();
            for (Field field : clazz.getDeclaredFields()) {
                PlcField fieldAnnotation = field.getAnnotation(PlcField.class);
                if (fieldAnnotation == null) {
                    // Ignore that field
                    continue;
                }
                // Create the suitable Request
                String query = fieldAnnotation.value();
                Class<?> expectedType;
                if (field.getType().isPrimitive()) {
                    if (field.getType() == long.class) {
                        expectedType = Long.class;
                    } else {
                        throw new OPMException("Unable to work on fields with type " + field.getType());
                    }
                } else {
                    expectedType = field.getType();
                }

                ReadRequestItem<?> item;
                try {
                    item = new ReadRequestItem<>(expectedType, connection.parseAddress(query));
                } catch (PlcInvalidAddressException e) {
                    throw new OPMException("Unable to parse address '" + query + "'");
                }

                // Store the item
                requestItems.put(item, field);
            }

            // Build the request
            PlcReadRequest.Builder builder = new PlcReadRequest.Builder();
            requestItems.keySet().forEach(item -> builder.addItem(item));
            PlcReadRequest request = builder.build();

            // Perform the request
            PlcReadResponse response;
            try {
                response = reader.read(request).get(1_000, TimeUnit.MILLISECONDS);
            } catch (InterruptedException | ExecutionException e) {
                throw new OPMException("Request fetching not able", e);
            } catch (TimeoutException e) {
                throw new OPMException("Timeout during fetching values", e);
            }

            // Construct the Object
            T instance = clazz.getConstructor().newInstance();

            // Assign values to all fields
            for (Map.Entry<ReadRequestItem<?>, Field> entry : requestItems.entrySet()) {
                Optional<? extends ReadResponseItem<?>> responseItem = response.getValue(entry.getKey());

                if (!responseItem.isPresent()) {
                    throw new OPMException("Unable to fetch value for field '" + entry.getValue().getName() + "'");
                }

                // Fetch first value
                Object value = responseItem.get().getValues().get(0);

                Field field = entry.getValue();
                field.setAccessible(true);
                field.set(instance, value);
            }
            return instance;
        } catch (PlcConnectionException e) {
            throw new OPMException("Unable to get connection with url '" + source + "'", e);
        } catch (Exception e) {
            throw new OPMException("Unable to fetch PlcEntity " + clazz.getName(), e);
        }
    }

}
