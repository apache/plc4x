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

package org.apache.plc4x.java.tools.eventpump.config;

import java.io.IOException;

/**
 * Configuration for a single tag.
 * <p>
 * This class supports both simple string format (just the address) and
 * extended format (address + optional transformation expression).
 * <p>
 * Simple format (YAML):
 * <pre>
 * tags:
 *   temperature: "MAIN.temperature"
 * </pre>
 * <p>
 * Extended format (YAML):
 * <pre>
 * tags:
 *   temperature:
 *     address: "MAIN.temperature"
 *     transform: "value * 1.8 + 32"
 * </pre>
 */
public class TagConfiguration {

    private String address;

    private String transform;

    /**
     * Creates an empty tag configuration.
     */
    public TagConfiguration() {
    }

    /**
     * Create a tag configuration with just an address.
     *
     * @param address The tag address
     */
    public TagConfiguration(String address) {
        this.address = address;
    }

    /**
     * Create a tag configuration with address and transform.
     *
     * @param address The tag address
     * @param transform The transformation expression (can be null)
     */
    public TagConfiguration(String address, String transform) {
        this.address = address;
        this.transform = transform;
    }

    /**
     * Get the tag address.
     *
     * @return The tag address
     */
    public String getAddress() {
        return address;
    }

    /**
     * Set the tag address.
     *
     * @param address The tag address
     */
    public void setAddress(String address) {
        this.address = address;
    }

    /**
     * Get the transformation expression.
     *
     * @return The transformation expression, or null if no transformation
     */
    public String getTransform() {
        return transform;
    }

    /**
     * Set the transformation expression.
     *
     * @param transform The transformation expression
     */
    public void setTransform(String transform) {
        this.transform = transform;
    }

    /**
     * Check if this tag has a transformation.
     *
     * @return true if a transformation is configured
     */
    public boolean hasTransform() {
        return transform != null && !transform.trim().isEmpty();
    }

    @Override
    public String toString() {
        if (hasTransform()) {
            return "TagConfiguration{address='" + address + "', transform='" + transform + "'}";
        } else {
            return "TagConfiguration{address='" + address + "'}";
        }
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;

        TagConfiguration that = (TagConfiguration) o;

        if (address != null ? !address.equals(that.address) : that.address != null) return false;
        return transform != null ? transform.equals(that.transform) : that.transform == null;
    }

    @Override
    public int hashCode() {
        int result = address != null ? address.hashCode() : 0;
        result = 31 * result + (transform != null ? transform.hashCode() : 0);
        return result;
    }

    /**
     * A tag in either format: just the address, or an object with address and transform.
     */
    static TagConfiguration fromTree(Object node, String path) throws IOException {
        if (ConfigNode.isScalar(node)) {
            return new TagConfiguration(ConfigNode.scalarToString(node, path));
        }
        if (node == null) {
            throw new IOException(path + ": expected an address or an object with an address, got nothing");
        }
        ConfigNode object = ConfigNode.object(node, path, "address", "transform");
        return new TagConfiguration(object.string("address"), object.string("transform"));
    }

    /**
     * The simple format when there's nothing but the address, the extended one otherwise.
     */
    Object toTree() {
        if (transform == null) {
            return address;
        }
        return new ConfigNode.Writer()
            .put("address", address)
            .put("transform", transform)
            .tree();
    }

}
