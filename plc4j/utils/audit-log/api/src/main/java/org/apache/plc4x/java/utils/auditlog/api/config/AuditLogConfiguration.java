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

package org.apache.plc4x.java.utils.auditlog.api.config;


import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;

/**
 * Configuration for the audit log system.
 * This configuration can be passed via connection string parameters.
 */
public class AuditLogConfiguration implements Configuration {

    /**
     * Path to the file where an audit log file will be created.
     * If not specified or empty, audit logging is disabled.
     */
    @ConfigurationParameter("audit-log-file")
    @Description("Path to the audit-log file where a all audit-log output will be written.")
    public String auditLogFile;

    /**
     * Checks if audit logging is enabled (i.e., audit log file is configured)
     * @return true if audit logging is enabled, false otherwise
     */
    public boolean isEnabled() {
        return auditLogFile != null && !auditLogFile.trim().isEmpty();
    }

}
