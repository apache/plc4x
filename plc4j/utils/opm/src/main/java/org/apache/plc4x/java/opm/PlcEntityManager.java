package org.apache.plc4x.java.opm;

import net.bytebuddy.ByteBuddy;
import net.bytebuddy.description.modifier.Visibility;
import net.bytebuddy.implementation.MethodDelegation;
import net.bytebuddy.implementation.bind.annotation.*;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.Optional;
import java.util.concurrent.Callable;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

import static net.bytebuddy.matcher.ElementMatchers.any;

/**
 * Plc4x equivalent of Jpas EntityManager for implementing Object-Plc-Mapping.
 * This means that calls to a plc can be done by using plain POJOs with Annotations.
 *
 * First, the necessary annotations are {@link PlcEntity} and {@link PlcField}.
 * For a class to be useable as PlcEntity it needs
 * <ul>
 *     <li>be non-final (as proxiing has to be used in case of {@link #connect(Class)}</li>
 * </ul>
 */
public class PlcEntityManager {

    private static final Logger logger = LoggerFactory.getLogger(PlcEntityManager.class);

    private final PlcDriverManager driverManager;

    public PlcEntityManager() {
        this.driverManager = new PlcDriverManager();
    }

    public PlcEntityManager(PlcDriverManager driverManager) {
        this.driverManager = driverManager;
    }

    public <T> T read(Class<T> clazz) throws OPMException {
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

            PlcReadRequest.Builder requestBuilder = reader.readRequestBuilder();

            // Do the necessary queries for all fields
            // HashMap<ReadRequestItem<?>, Field> requestItems = new HashMap<>();
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

                requestBuilder.addItem(field.getName(), query);
            }

            // Build the request
            PlcReadRequest request;
            try {
                request = requestBuilder.build();
            } catch (PlcInvalidFieldException e) {
                throw new OPMException("Unable to parse one field request", e);
            }

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
            for (String field : request.getFieldNames()) {
                Object value = response.getObject(field);

                if (value == null) {
                    throw new OPMException("Unable to fetch value for field '" + field + "'");
                }

                // Fetch first value
                Field objectField = clazz.getDeclaredField(field);
                objectField.setAccessible(true);
                objectField.set(instance, value);
            }
            return instance;
        } catch (PlcConnectionException e) {
            throw new OPMException("Unable to get connection with url '" + source + "'", e);
        } catch (Exception e) {
            throw new OPMException("Unable to fetch PlcEntity " + clazz.getName(), e);
        }
    }

    /**
     * Returns a connected proxy.
     *
     * @param clazz
     * @param <T>
     * @return
     * @throws OPMException
     */
    public <T> T connect(Class<T> clazz) throws OPMException {
        PlcEntity annotation = clazz.getAnnotation(PlcEntity.class);
        if (annotation == null) {
            throw new OPMException("Need to be a PLC Entity, please add Annotation.");
        }
        try {
            // Use Byte Buddy to generate a subclassed proxy that delegates all PlcField Methods
            // to the intercept method
            T instance = new ByteBuddy()
                .subclass(clazz)
                .defineField("parent", Class.class, Visibility.PUBLIC)
                .method(any()).intercept(MethodDelegation.to(this))
                .make()
                .load(Thread.currentThread().getContextClassLoader())
                .getLoaded()
                .getConstructor()
                .newInstance();

            Field parent = instance.getClass().getDeclaredField("parent");
            parent.set(instance, clazz);
            return instance;
        } catch (NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException | NoSuchFieldException e) {
            throw new OPMException("Unable to instantiate Proxy", e);
        }
    }

    /**
     * Intersect "defined" methods
     *
     * @param o
     * @param m
     * @param c
     * @return
     * @throws OPMException
     */
    @RuntimeType
    public Object intercept(@This Object o, @Origin Method m, @SuperCall Callable<?> c, @Super Object that) throws OPMException {
        logger.trace("Invoked " + m.getName() + " fetch all values...");

        Field field = that.getClass().getDeclaredFields()[0];
        field.setAccessible(true);
        Object base;
        try {
            base = field.get(that);
        } catch (IllegalAccessException e) {
            throw new OPMException("...", e);
        }

        if (m.getName().startsWith("get") || m.getName().startsWith("is")) {
            // Fetch single value
            return fetchValueInternal(that, m);
        }

        // Fetch all values, than invoke method
        try {
            fetchAllValues(base, m);
            return c.call();
        } catch (Exception e) {
            throw new OPMException("Unbale to forward call", e);
        }
    }

    /**
     * Intersect abstract methods
     *
     * @param m
     * @param o
     * @return
     * @throws OPMException
     */
    @RuntimeType
    public Object interceptGetter(@Origin Method m, @This Object o, @Super Object that) throws OPMException {
        fetchValueInternal(that, m);

        // Finished
        return 1L;
    }

    private void fetchAllValues(Object o, Method m) throws OPMException {
        Class<?> baseClass;
        try {
            baseClass = (Class<?>) o.getClass().getDeclaredField("parent").get(o);
        } catch (IllegalAccessException | NoSuchFieldException e) {
            e.printStackTrace();
            throw new OPMException("...", e);
        }
        PlcEntity plcEntity = baseClass.getAnnotation(PlcEntity.class);
        logger.trace("For source: " + plcEntity.value());
        logger.trace("Using the DriverManager: " + driverManager);

        Optional<PlcReader> reader;
        try {
            reader = driverManager.getConnection(plcEntity.value()).getReader();
        } catch (PlcConnectionException e) {
            throw new OPMException("Unable to acquire connection", e);
        }

        if (reader.isPresent() == false) {
            throw new OPMException("Unable to generate Reader");
        }

        PlcReader plcReader = reader.get();

        // Assume to do the query here...
        PlcReadRequest.Builder builder = plcReader.readRequestBuilder();
        for (Field field : baseClass.getDeclaredFields()) {
            // Check if the field has an annotation
            PlcField plcField = field.getDeclaredAnnotation(PlcField.class);
            if (plcField != null) {
                logger.trace("Adding field " + field.getName() + " to request as " + plcField.value());
                builder.addItem(field.getName(), plcField.value());
            }

        }
        PlcReadRequest request = builder.build();

        PlcReadResponse<?> response = getPlcReadResponse(plcReader, request);

        // Fill all requested fields
        for (String fieldName : response.getFieldNames()) {
            logger.trace("Value for field " + fieldName + " is " + response.getObject(fieldName));
            logger.trace("Setting test value");
            try {
                Field field = baseClass.getDeclaredField(fieldName);
                field.setAccessible(true);
                if (field.getType().isPrimitive()) {
                    if (field.getType() == byte.class) {
                        field.set(o, (byte) response.getByte(fieldName));
                    } else if (field.getType() == int.class) {
                        field.set(o, (int) response.getInteger(fieldName));
                    } else if (field.getType() == long.class) {
                        field.set(o, (long) response.getLong(fieldName));
                    } else if (field.getType() == short.class) {
                        field.set(o, (short) response.getShort(fieldName));
                    }
                    // TODO this should fail on Short, Integer, ... because it always gets a Long
                } else if (field.getType().isAssignableFrom(response.getObject(fieldName).getClass())){
                    field.set(o, response.getObject(fieldName));
                } else {
                    logger.trace("Unable to assign return value {} to field {} with type {}", response.getObject(fieldName), fieldName, field.getType());
                }
            } catch (NoSuchFieldException | IllegalAccessException e) {
                e.printStackTrace();
            }
        }
    }

    private Object fetchValueInternal(Object o, Method m) throws OPMException {
        String s = m.getName().substring(3);
        // First char to lower
        String variable = s.substring(0, 1).toLowerCase().concat(s.substring(1));
        logger.trace("Variable: " + variable);
        PlcField annotation = null;
        try {
            annotation = m.getDeclaringClass().getDeclaredField(variable).getDeclaredAnnotation(PlcField.class);
        } catch (NoSuchFieldException e) {
            e.printStackTrace();
        }
        logger.trace("You wanted field: " + annotation.value());
        PlcEntity plcEntity = m.getDeclaringClass().getAnnotation(PlcEntity.class);
        logger.trace("For source: " + plcEntity.value());
        logger.trace("Using the DriverManager: " + driverManager);

        Optional<PlcReader> reader;
        try {
            reader = driverManager.getConnection(plcEntity.value()).getReader();
        } catch (PlcConnectionException e) {
            throw new OPMException("Unable to acquire connection", e);
        }

        if (reader.isPresent() == false) {
            throw new OPMException("Unable to generate Reader");
        }

        PlcReader plcReader = reader.get();

        // Assume to do the query here...
        PlcReadRequest request = reader.get().readRequestBuilder()
            .addItem(m.getName(), annotation.value())
            .build();

        PlcReadResponse<?> response;
        response = getPlcReadResponse(plcReader, request);

        Object responseObject = response.getObject(m.getName());

        if (responseObject.getClass().isAssignableFrom(m.getReturnType())) {
            return responseObject;
        } else {
                throw new OPMException("Unable to cast the PLC Object '" + responseObject +
                    "' to the expected method return type '" + m.getReturnType() + "'");
        }
    }

    /**
     * Fetch the request and do appropriate error handling
     * @param plcReader
     * @param request
     * @return
     * @throws OPMException
     */
    private PlcReadResponse<?> getPlcReadResponse(PlcReader plcReader, PlcReadRequest request) throws OPMException {
        PlcReadResponse<?> response;
        try {
            response = plcReader.read(request).get();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new OPMException("Exception during execution", e);
        } catch (ExecutionException e) {
            throw new OPMException("Exception during execution", e);
        }
        return response;
    }

}
