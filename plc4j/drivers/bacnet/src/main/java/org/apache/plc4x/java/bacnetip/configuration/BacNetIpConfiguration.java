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
package org.apache.plc4x.java.bacnetip.configuration;

import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.Description;

public class BacNetIpConfiguration implements PlcConnectionConfiguration {

    // Path to a single EDE file.
    @ConfigurationParameter("ede-file-path")
    @Description("Path to the location of a single EDE file, that contains the descriptor for the target device.")
    private String edeFilePath;

    // Path to a directory containing many EDE files.
    @ConfigurationParameter("ede-directory-path")
    @Description("Path to the directory used for storing multiple EDE files. These files contain the descriptors for the possible target devices.")
    private String edeDirectoryPath;

    public String getEdeFilePath() {
        return edeFilePath;
    }

    public void setEdeFilePath(String edeFilePath) {
        this.edeFilePath = edeFilePath;
    }

    public String getEdeDirectoryPath() {
        return edeDirectoryPath;
    }

    public void setEdeDirectoryPath(String edeDirectoryPath) {
        this.edeDirectoryPath = edeDirectoryPath;
    }

}
