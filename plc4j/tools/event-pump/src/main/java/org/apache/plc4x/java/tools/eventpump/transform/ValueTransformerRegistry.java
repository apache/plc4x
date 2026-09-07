/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.tools.eventpump.transform;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Registry for managing ValueTransformer implementations.
 * <p>
 * This allows for pluggable expression evaluators. By default, the SimpleExpressionEvaluator
 * is registered, but applications can register their own transformers (e.g., SpEL-based).
 * <p>
 * Thread-safe instance-based implementation. Use {@link #createDefault()} to create a registry
 * with standard transformers pre-registered, or create an empty instance with the constructor.
 * <p>
 * Example usage:
 * <pre>
 * // Create a default registry with standard transformers
 * ValueTransformerRegistry registry = ValueTransformerRegistry.createDefault();
 *
 * // Use it in a TagBatch
 * TagBatch batch = TagBatch.builder()
 *     .withTransformerRegistry(registry)
 *     ...
 *     .build();
 *
 * // Or create a custom registry
 * ValueTransformerRegistry customRegistry = new ValueTransformerRegistry();
 * customRegistry.register(new MyCustomTransformer());
 * </pre>
 */
public class ValueTransformerRegistry {

    private static final Logger LOGGER = LoggerFactory.getLogger(ValueTransformerRegistry.class);

    private final Map<String, ValueTransformer> transformers = new ConcurrentHashMap<>();
    private volatile String defaultTransformerName = "simple";

    /**
     * Create an empty registry with no transformers registered.
     * Use {@link #createDefault()} if you want standard transformers pre-registered.
     */
    public ValueTransformerRegistry() {
        // Empty constructor for custom registries
    }

    /**
     * Create a registry with standard transformers pre-registered.
     * This includes:
     * - SimpleExpressionEvaluator (as "simple", set as default)
     * - SpelValueTransformer (as "spel", if Spring Expression is on the classpath)
     *
     * @return A new registry with standard transformers
     */
    public static ValueTransformerRegistry createDefault() {
        ValueTransformerRegistry registry = new ValueTransformerRegistry();

        // Register the default simple evaluator
        registry.register(new SimpleExpressionEvaluator());

        return registry;
    }


    /**
     * Register a transformer implementation.
     *
     * @param transformer The transformer to register
     * @throws IllegalArgumentException if transformer is null or name is already registered
     */
    public void register(ValueTransformer transformer) {
        if (transformer == null) {
            throw new IllegalArgumentException("Transformer cannot be null");
        }

        String name = transformer.getName();
        if (name == null || name.trim().isEmpty()) {
            throw new IllegalArgumentException("Transformer name cannot be null or empty");
        }

        ValueTransformer existing = transformers.putIfAbsent(name, transformer);
        if (existing != null) {
            LOGGER.warn("Transformer '{}' is already registered, keeping existing", name);
        } else {
            LOGGER.info("Registered value transformer: {}", name);
        }
    }

    /**
     * Get a transformer by name.
     *
     * @param name The transformer name
     * @return The transformer, or null if not found
     */
    public ValueTransformer get(String name) {
        return transformers.get(name);
    }

    /**
     * Get the default transformer.
     *
     * @return The default transformer
     * @throws IllegalStateException if the default transformer is not registered
     */
    public ValueTransformer getDefault() {
        ValueTransformer transformer = transformers.get(defaultTransformerName);
        if (transformer == null) {
            throw new IllegalStateException("Default transformer '" + defaultTransformerName + "' is not registered");
        }
        return transformer;
    }

    /**
     * Set the default transformer by name.
     *
     * @param name The name of the transformer to use as default
     * @throws IllegalArgumentException if the transformer is not registered
     */
    public void setDefault(String name) {
        if (!transformers.containsKey(name)) {
            throw new IllegalArgumentException("Transformer '" + name + "' is not registered");
        }
        this.defaultTransformerName = name;
        LOGGER.info("Set default value transformer to: {}", name);
    }

    /**
     * Unregister a transformer.
     *
     * @param name The name of the transformer to unregister
     * @return true if the transformer was removed, false if it didn't exist
     * @throws IllegalStateException if trying to unregister the default transformer
     */
    public boolean unregister(String name) {
        if (name.equals(defaultTransformerName)) {
            throw new IllegalStateException("Cannot unregister the default transformer");
        }

        ValueTransformer removed = transformers.remove(name);
        if (removed != null) {
            LOGGER.info("Unregistered value transformer: {}", name);
            return true;
        }
        return false;
    }

    /**
     * Get all registered transformer names.
     *
     * @return A set of transformer names
     */
    public java.util.Set<String> getRegisteredNames() {
        return java.util.Collections.unmodifiableSet(transformers.keySet());
    }
}
