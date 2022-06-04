use std::iter;
use tokio::net::{TcpListener, TcpSocket, TcpStream};
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use plc4rust::modbus::{ModbusPDU, ModbusADU, ModbusADUOptions};
use plc4rust::modbus::modbus_pdu::ModbusPDUReadCoilsRequest;
use plc4rust::modbus::modbus_adu::ModbusTcpADU;
use plc4rust::{Endianess, Message};
use plc4rust::write_buffer::WriteBuffer;
use tokio::select;
use tokio::pin;
use tokio::runtime::Handle;
use plc4rust::Endianess::BigEndian;
use plc4rust::read_buffer::ReadBuffer;

#[tokio::main]
async fn main() {
    match TcpStream::connect("127.0.0.1:502").await {
        Ok(mut stream) => {




            println!("Connetion established!");
            let send_future = send(&mut stream);

            // loop {
            //     select! {
            //         val = send_future => {
            //             println!("Was sent!")
            //         }
            //         val = receiver => {
            //             println!("Received?!")
            //         }
            //     }
            // }
            send_future.await;

            let receiver = async {
                let mut buf = [0; 1024];
                // In a loop, read data from the socket and write the data back.
                loop {
                    let n = match stream.read(&mut buf).await {
                        // socket closed
                        Ok(n) if n == 0 => return,
                        Ok(n) => n,
                        Err(e) => {
                            eprintln!("failed to read from socket; err = {:?}", e);
                            return;
                        }
                    };

                    println!("We received bytes: {:?}", buf.as_slice());

                    // Deserialize
                    let mut read_buffer = ReadBuffer::new(BigEndian, buf.as_slice());
                    let response = ModbusADU::parse(&mut read_buffer, Some(ModbusADUOptions {
                        driver_type: plc4rust::modbus::DriverType::MODBUS_TCP,
                        response: true
                    }));

                    println!("Parsed the response: {:?}", response);
                }
            };
            println!("Waiting for bytes...");
            receiver.await;
        }
        Err(err) => {
            println!("No connection possible :(");
            println!("{}", err);
        }
    }
}

async fn send(stream: &mut TcpStream) {
    let request = create_request();

    let mut bytes: Vec<u8> = vec![];
    let mut write_buffer = WriteBuffer::new(Endianess::BigEndian, &mut bytes);

    request.serialize(&mut write_buffer).unwrap();

    println!("Write bytes...");
    stream.write(bytes.as_slice()).await;
    println!("Bytes written...")
}

fn create_request() -> ModbusADU {
    let pdu = ModbusPDU::ModbusPDUReadCoilsRequest(
        ModbusPDUReadCoilsRequest {
            startingAddress: 0,
            quantity: 1,
        }
    );

    let adu = ModbusADU::ModbusTcpADU(ModbusTcpADU {
        transaction_identifier: 1,
        protocol_identifier: 0x0000,
        unit_identifier: 1,
        pdu,
    });

    return adu;
}
