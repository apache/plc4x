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

use std::time::Duration;

/// Common trait for all transport configurations
pub trait TransportConfig {
    /// Get the host address
    fn host(&self) -> &str;

    /// Get the port number
    fn port(&self) -> u16;
}

/// Configuration for TCP transport
#[derive(Debug, Clone)]
pub struct TcpConfig {
    pub host: String,
    pub port: u16,
    pub connect_timeout: Option<Duration>,
    pub read_timeout: Option<Duration>,
    pub write_timeout: Option<Duration>,
    pub nodelay: Option<bool>,
    pub ttl: Option<u32>,
    pub retry_count: Option<u32>,
}

impl TcpConfig {
    /// Create a new TCP configuration with default settings
    pub fn new(host: impl Into<String>, port: u16) -> Self {
        Self {
            host: host.into(),
            port,
            connect_timeout: Some(Duration::from_secs(5)),
            read_timeout: Some(Duration::from_secs(5)),
            write_timeout: Some(Duration::from_secs(5)),
            nodelay: Some(true),
            ttl: None,
            retry_count: Some(3),
        }
    }

    /// Set the connect timeout
    pub fn with_connect_timeout(mut self, timeout: Duration) -> Self {
        self.connect_timeout = Some(timeout);
        self
    }

    /// Set the read timeout
    pub fn with_read_timeout(mut self, timeout: Duration) -> Self {
        self.read_timeout = Some(timeout);
        self
    }

    /// Set the write timeout
    pub fn with_write_timeout(mut self, timeout: Duration) -> Self {
        self.write_timeout = Some(timeout);
        self
    }

    /// Set the TCP nodelay option
    pub fn with_nodelay(mut self, nodelay: bool) -> Self {
        self.nodelay = Some(nodelay);
        self
    }

    /// Set the TTL value
    pub fn with_ttl(mut self, ttl: u32) -> Self {
        self.ttl = Some(ttl);
        self
    }

    /// Set the retry count
    pub fn with_retry_count(mut self, count: u32) -> Self {
        self.retry_count = Some(count);
        self
    }
}

impl TransportConfig for TcpConfig {
    fn host(&self) -> &str {
        &self.host
    }

    fn port(&self) -> u16 {
        self.port
    }
}

/// Configuration for UDP transport
#[derive(Debug, Clone)]
pub struct UdpConfig {
    pub host: String,
    pub port: u16,
    pub local_port: Option<u16>,
    pub read_timeout: Option<Duration>,
    pub write_timeout: Option<Duration>,
    pub ttl: Option<u32>,
    pub broadcast: Option<bool>,
}

impl UdpConfig {
    /// Create a new UDP configuration with default settings
    pub fn new(host: impl Into<String>, port: u16) -> Self {
        Self {
            host: host.into(),
            port,
            local_port: None,
            read_timeout: Some(Duration::from_secs(5)),
            write_timeout: Some(Duration::from_secs(5)),
            ttl: None,
            broadcast: None,
        }
    }

    /// Set the local port to bind to
    pub fn with_local_port(mut self, port: u16) -> Self {
        self.local_port = Some(port);
        self
    }

    /// Set the read timeout
    pub fn with_read_timeout(mut self, timeout: Duration) -> Self {
        self.read_timeout = Some(timeout);
        self
    }

    /// Set the write timeout
    pub fn with_write_timeout(mut self, timeout: Duration) -> Self {
        self.write_timeout = Some(timeout);
        self
    }

    /// Set the TTL value
    pub fn with_ttl(mut self, ttl: u32) -> Self {
        self.ttl = Some(ttl);
        self
    }

    /// Set the broadcast option
    pub fn with_broadcast(mut self, broadcast: bool) -> Self {
        self.broadcast = Some(broadcast);
        self
    }
}

impl TransportConfig for UdpConfig {
    fn host(&self) -> &str {
        &self.host
    }

    fn port(&self) -> u16 {
        self.port
    }
}
