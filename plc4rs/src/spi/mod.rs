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

pub mod config;
mod tcp;
mod udp;

use bytes::BytesMut;
use std::fmt::Debug;
use crate::types::Result;

// Re-export implementations
pub use tcp::TcpTransport;
pub use udp::UdpTransport;

/// Transport trait defining the interface for all transport implementations
#[async_trait::async_trait]
pub trait Transport: Debug + Send + Sync {
    /// Connect to the target
    async fn connect(&mut self) -> Result<()>;
    
    /// Read data from the transport
    async fn read(&mut self, buffer: &mut BytesMut) -> Result<usize>;
    
    /// Write data to the transport
    async fn write(&mut self, data: &[u8]) -> Result<usize>;
    
    /// Close the transport connection
    async fn close(&mut self) -> Result<()>;
}
