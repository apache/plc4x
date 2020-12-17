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
package org.apache.plc4x.test.driver.model;

import org.apache.plc4x.test.model.Location;
import org.apache.plc4x.test.model.LocationAware;

import java.util.List;
import java.util.Optional;

public class Testcase implements LocationAware {

    private final String name;
    private final String description;
    private final List<TestStep> steps;
    private Location location;

    public Testcase(String name, String description, List<TestStep> steps) {
        this.name = name;
        this.description = description;
        this.steps = steps;
    }

    public String getName() {
        return name;
    }

    public String getDescription() {
        return description;
    }

    public List<TestStep> getSteps() {
        return steps;
    }

    public void setLocation(Location location) {
        this.location = location;
    }

    @Override
    public Optional<Location> getLocation() {
        return Optional.ofNullable(location);
    }

}
