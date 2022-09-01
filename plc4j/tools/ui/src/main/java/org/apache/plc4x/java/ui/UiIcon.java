/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.ui;

import org.kordamp.ikonli.Ikon;

public enum UiIcon implements Ikon {

    FOLDER("ui-icon-folder", '\ue815');

    private String description;
    private char icon;

    UiIcon(String description, char icon) {
        this.description = description;
        this.icon = icon;
    }

    @Override
    public String getDescription() {
        return description;
    }

    @Override
    public int getCode() {
        return icon;
    }

    public static UiIcon findByDescription(String description) {
        for (UiIcon icon : values()) {
            if (icon.description.equals(description)) {
                return icon;
            }
        }
        throw new IllegalArgumentException("Icon not supported: " + description);
    }

}
