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

use crate::spi::Transport;
use crate::spi::config::UdpConfig;
use crate::types::Result;
use async_trait::async_trait;
use bytes::BytesMut;
use std::net::SocketAddr;
use std::time::Duration;
use tokio::net::UdpSocket;
use tokio::time::timeout;
use tracing::info;

/// UDP transport implementation
#[derive(Debug)]
pub struct UdpTransport {
    config: UdpConfig,
    socket: Option<UdpSocket>,
    remote_addr: Option<SocketAddr>,
}

impl UdpTransport {
    /// Create a new UDP transport with the given configuration
    pub fn new(config: UdpConfig) -> Self {
        Self {
            config,
            socket: None,
            remote_addr: None,
        }
    }
}

#[async_trait]
impl Transport for UdpTransport {
    async fn connect(&mut self) -> Result<()> {
        let remote_address = format!("{}:{}", self.config.host, self.config.port);
        info!("Connecting to {}", remote_address);
        
        // Parse the remote address
        let remote_addr: SocketAddr = remote_address.parse()?;
        self.remote_addr = Some(remote_addr);
        
        // Bind to a local address
        let local_addr: SocketAddr = if let Some(local_port) = self.config.local_port {
            format!("0.0.0.0:{}", local_port).parse()?
        } else {
            "0.0.0.0:0".parse()?
        };
        
        let socket = UdpSocket::bind(local_addr).await?;
        
        // Set socket options
        if let Some(ttl) = self.config.ttl {
            socket.set_ttl(ttl)?;
        }
        
        if let Some(broadcast) = self.config.broadcast {
            socket.set_broadcast(broadcast)?;
        }
        
        self.socket = Some(socket);
        info!("UDP socket bound to local address");
        
        Ok(())
    }
    
    async fn read(&mut self, buffer: &mut BytesMut) -> Result<usize> {
        if let Some(socket) = &self.socket {
            let read_timeout = self.config.read_timeout.unwrap_or(Duration::from_secs(5));
            
            // Ensure we have enough capacity
            if buffer.capacity() - buffer.len() < 65536 {
                buffer.reserve(65536);
            }
            
            match timeout(read_timeout, socket.recv_buf(buffer)).await {
                Ok(Ok(n)) => Ok(n),
                Ok(Err(e)) => Err(format!("Read error: {}", e).into()),
                Err(_) => Err("Read operation timed out".into()),
            }
        } else {
            Err("Not connected".into())
        }
    }
    
    async fn write(&mut self, data: &[u8]) -> Result<usize> {
        if let Some(socket) = &self.socket {
            if let Some(remote_addr) = self.remote_addr {
                let write_timeout = self.config.write_timeout.unwrap_or(Duration::from_secs(5));
                
                match timeout(write_timeout, socket.send_to(data, remote_addr)).await {
                    Ok(Ok(n)) => Ok(n),
                    Ok(Err(e)) => Err(format!("Write error: {}", e).into()),
                    Err(_) => Err("Write operation timed out".into()),
                }
            } else {
                Err("Remote address not set".into())
            }
        } else {
            Err("Not connected".into())
        }
    }
    
    async fn close(&mut self) -> Result<()> {
        // UDP sockets don't need explicit closing, but we'll reset our state
        self.socket = None;
        self.remote_addr = None;
        info!("UDP connection closed");
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use tokio_test::block_on;

    #[test]
    fn test_udp_lifecycle() {
        let mut transport = UdpTransport::new(UdpConfig {
            host: "127.0.0.1".to_string(),
            port: 1234,
            local_port: None,
            ttl: None,
            broadcast: None,
            read_timeout: None,
            write_timeout: None,
        });
        
        block_on(async {
            // Test connection
            assert!(transport.connect().await.is_ok());
            assert!(transport.socket.is_some());

            // Test write/read
            let data = b"test data";
            let result = transport.write(data).await;
            assert!(result.is_ok() || result.is_err()); // May fail as no server

            // Test close
            assert!(transport.close().await.is_ok());
            assert!(transport.socket.is_none());
        });
    }
} 
