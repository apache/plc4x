use tokio::net::{TcpListener, TcpStream};
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use plc4rust::Endianess::BigEndian;
use plc4rust::modbus::{DriverType, ModbusADU, ModbusADUOptions};
use plc4rust::read_buffer::ReadBuffer;
use plc4rust::{Message, modbus};
use plc4rust::modbus::modbus_adu::ModbusTcpADU;
use plc4rust::modbus::modbus_pdu::ModbusPDUReadCoilsResponse;
use plc4rust::write_buffer::WriteBuffer;

#[tokio::main]
async fn main() {
    println!("Before...");
    if let Ok(mut tcp_listener) = TcpListener::bind("0.0.0.0:502").await {
        while let Ok((mut tcp_stream, _socket_addr)) = tcp_listener.accept().await {
            tokio::spawn(async move {
                let mut buf = [0; 1024];
                // In a loop, read data from the socket and write the data back.
                loop {
                    let n = match tcp_stream.read(&mut buf).await {
                        // socket closed
                        Ok(n) if n == 0 => return,
                        Ok(n) => n,
                        Err(e) => {
                            eprintln!("failed to read from socket; err = {:?}", e);
                            return;
                        }
                    };
                    println!("We received {:?}", &buf);
                    let mut read_buffer = ReadBuffer::new(BigEndian, buf.as_slice());
                    let request = ModbusADU::parse(&mut read_buffer, Some(ModbusADUOptions {
                        driver_type: DriverType::MODBUS_TCP,
                        response: false
                    }));

                    println!("Decoded request: {:?}", request);

                    // Returning response
                    let mut bytes: Vec<u8> = vec![];
                    let mut write_buffer = WriteBuffer::new(BigEndian, &mut bytes);

                    let response = if let Ok(modbus::modbus_adu::ModbusADU::ModbusTcpADU(tcp_request)) = request {
                        let pdu = match tcp_request.pdu {
                            modbus::modbus_pdu::ModbusPDU::ModbusPDUReadCoilsRequest(msg) => {
                                modbus::modbus_pdu::ModbusPDU::ModbusPDUReadCoilsResponse(ModbusPDUReadCoilsResponse {
                                    value: vec![1, 2, 3]
                                })
                            }
                            _ => {
                                panic!("Unable to deser pdu")
                            }
                        };
                        ModbusADU::ModbusTcpADU(
                                ModbusTcpADU {
                                    transaction_identifier: tcp_request.transaction_identifier,
                                    protocol_identifier: tcp_request.protocol_identifier,
                                    unit_identifier: tcp_request.unit_identifier,
                                    pdu
                                }
                            )
                    } else {
                        panic!("Waah?!");
                    };

                    println!("Sending response {:?}", response);
                    response.serialize(&mut write_buffer);

                    println!("Sending bytes {:?}", &bytes);
                    tcp_stream.write_all(bytes.as_slice()).await;

                    // // Write the data back
                    // if let Err(e) = tcp_stream.write_all(&buf[0..n]).await {
                    //     eprintln!("failed to write to socket; err = {:?}", e);
                    //     return;
                    // }
                }
            });
        }
    }
    println!("After")
}
