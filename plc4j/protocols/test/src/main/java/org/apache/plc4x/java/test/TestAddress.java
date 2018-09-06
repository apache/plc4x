/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.test;

import org.apache.plc4x.java.api.model.Address;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Test address for accessing values in virtual devices.
 *
 */
class TestAddress implements Address {
    private static final Pattern ADDRESS_PATTERN = Pattern.compile("^\\w+$");

    private final String value;

    public static final TestAddress RANDOM = new TestAddress("random");

    private TestAddress(String value) {
        this.value = value;
    }

    public static boolean isValid(String addressString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        return matcher.matches();
    }

    public static TestAddress of(String addressString) {
        return new TestAddress(addressString);
    }

    @Override
    public int hashCode() {
        return value.hashCode();
    }

    @Override
    public boolean equals(Object o) {
        if (o == this)
            return true;

        if (!(o instanceof TestAddress))
            return false;

        TestAddress that = (TestAddress) o;

        return this.value.equals(that.value);
    }

    @Override
    public String toString() {
        return value;
    }
}
