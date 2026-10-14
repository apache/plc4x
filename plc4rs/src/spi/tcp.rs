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
use crate::spi::config::TcpConfig;
use crate::types::Result;
use async_trait::async_trait;
use bytes::BytesMut;
use std::time::Duration;
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::TcpStream;
use tokio::time::timeout;
use tracing::{debug, error, info};

/// TCP transport implementation
#[derive(Debug)]
pub struct TcpTransport {
    config: TcpConfig,
    stream: Option<TcpStream>,
}

impl TcpTransport {
    /// Create a new TCP transport with the given configuration
    pub fn new(config: TcpConfig) -> Self {
        Self {
            config,
            stream: None,
        }
    }
}

#[async_trait]
impl Transport for TcpTransport {
    async fn connect(&mut self) -> Result<()> {
        let address = format!("{}:{}", self.config.host, self.config.port);
        info!("Connecting to {}", address);
        
        let mut retry_count = 0;
        let max_retries = self.config.retry_count.unwrap_or(3);
        
        while retry_count < max_retries {
            match timeout(
                self.config.connect_timeout.unwrap_or(Duration::from_secs(5)),
                TcpStream::connect(&address),
            ).await {
                Ok(Ok(stream)) => {
                    // Configure the stream
                    if let Some(nodelay) = self.config.nodelay {
                        stream.set_nodelay(nodelay)?;
                    }
                    
                    if let Some(ttl) = self.config.ttl {
                        stream.set_ttl(ttl)?;
                    }
                    
                    self.stream = Some(stream);
                    info!("Connected to {}", address);
                    return Ok(());
                }
                Ok(Err(e)) => {
                    error!("Failed to connect: {}", e);
                    retry_count += 1;
                    if retry_count < max_retries {
                        let backoff = Duration::from_millis(100 * 2u64.pow(retry_count));
                        debug!("Retrying in {:?} (attempt {}/{})", backoff, retry_count + 1, max_retries);
                        tokio::time::sleep(backoff).await;
                    }
                }
                Err(_) => {
                    error!("Connection timed out");
                    retry_count += 1;
                    if retry_count < max_retries {
                        let backoff = Duration::from_millis(100 * 2u64.pow(retry_count));
                        debug!("Retrying in {:?} (attempt {}/{})", backoff, retry_count + 1, max_retries);
                        tokio::time::sleep(backoff).await;
                    }
                }
            }
        }
        
        Err("Failed to connect after maximum retries".into())
    }
    
    async fn read(&mut self, buffer: &mut BytesMut) -> Result<usize> {
        if let Some(stream) = &mut self.stream {
            let read_timeout = self.config.read_timeout.unwrap_or(Duration::from_secs(5));
            match timeout(read_timeout, stream.read_buf(buffer)).await {
                Ok(Ok(n)) => Ok(n),
                Ok(Err(e)) => Err(format!("Read error: {}", e).into()),
                Err(_) => Err("Read operation timed out".into()),
            }
        } else {
            Err("Not connected".into())
        }
    }
    
    async fn write(&mut self, data: &[u8]) -> Result<usize> {
        if let Some(stream) = &mut self.stream {
            let write_timeout = self.config.write_timeout.unwrap_or(Duration::from_secs(5));
            match timeout(write_timeout, stream.write(data)).await {
                Ok(Ok(n)) => Ok(n),
                Ok(Err(e)) => Err(format!("Write error: {}", e).into()),
                Err(_) => Err("Write operation timed out".into()),
            }
        } else {
            Err("Not connected".into())
        }
    }
    
    async fn close(&mut self) -> Result<()> {
        if let Some(stream) = &mut self.stream {
            stream.shutdown().await?;
            self.stream = None;
            info!("Connection closed");
        }
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use tokio_test::block_on;

    #[test]
    fn test_tcp_transport() {
        let mut transport = TcpTransport::new(TcpConfig {
            base: Default::default(),
            no_delay: true,
            keep_alive: true,
            host: "127.0.0.1".to_string(),
            port: 102,
            retry_count: None,
            connect_timeout: None,
            read_timeout: None,
            write_timeout: None,
            nodelay: None,
            ttl: None,
        });
        
        // Test connection
        block_on(async {
            assert!(transport.connect().await.is_err()); // Should fail as no server is running
            assert!(transport.stream.is_none());
        });
    }
} 
