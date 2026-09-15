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

use plc4rs::spi::{
    TcpTransport, UdpTransport,
    config::{TransportConfig, TcpConfig},
};
use std::time::Duration;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // TCP Example
    let tcp_config = TcpConfig {
        base: TransportConfig {
            connect_timeout: Duration::from_secs(10),
            read_timeout: Duration::from_secs(2),
            write_timeout: Duration::from_secs(2),
            buffer_size: 1024,
        },
        no_delay: true,
        keep_alive: true,
    };

    let mut tcp = TcpTransport::new_with_config("192.168.1.1".into(), 102, tcp_config);
    tcp.connect().await?;
    
    let data = b"Hello PLC";
    tcp.write(data).await?;
    
    let mut buffer = vec![0u8; 1024];
    let len = tcp.read(&mut buffer).await?;
    println!("TCP Received: {:?}", &buffer[..len]);
    
    tcp.close().await?;

    // UDP Example
    let udp_config = TransportConfig {
        connect_timeout: Duration::from_secs(5),
        read_timeout: Duration::from_secs(1),
        write_timeout: Duration::from_secs(1),
        buffer_size: 1024,
    };

    let mut udp = UdpTransport::new_with_config("192.168.1.1".into(), 102, udp_config);
    udp.connect().await?;
    
    udp.write(data).await?;
    let len = udp.read(&mut buffer).await?;
    println!("UDP Received: {:?}", &buffer[..len]);
    
    udp.close().await?;
    
    Ok(())
} 
