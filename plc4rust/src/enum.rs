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

// [enum uint 8 ModbusErrorCode
//     ['1'    ILLEGAL_FUNCTION]
//     ['2'    ILLEGAL_DATA_ADDRESS]
//     ['3'    ILLEGAL_DATA_VALUE]
//     ['4'    SLAVE_DEVICE_FAILURE]
//     ['5'    ACKNOWLEDGE]
//     ['6'    SLAVE_DEVICE_BUSY]
//     ['7'    NEGATIVE_ACKNOWLEDGE]
//     ['8'    MEMORY_PARITY_ERROR]
//     ['10'   GATEWAY_PATH_UNAVAILABLE]
//     ['11'   GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND]
// ]

// Necessary imports
// use std::io::{Error, ErrorKind, Read, Write};
// use crate::{Message, NoOption, ReadBuffer, WriteBuffer};
#[macro_export]
macro_rules! plc4x_enum {
    (enum $type:ty : $name:ident $([$pat:expr => $branch:ident]),*) => {
        #[derive(Copy, Clone, PartialEq, Debug)]
        #[allow(non_camel_case_types)]
        pub enum $name {
            $(
                $branch,
            )* // Repeat for each `expr` passed to the macro
        }

        impl TryFrom<$type> for $name {
        type Error = ();

        fn try_from(value: $type) -> Result<Self, Self::Error> {
            match value {
                $(
                    $pat => {
                        Ok($name::$branch)
                    },
                )* // Repeat for each `expr` passed to the macro
                _ => {
                    Err(())
                }
            }
        }
        }

        impl Into<$type> for $name {
            fn into(self) -> $type {
                match self {
                    $(
                        $name::$branch => {
                            $pat
                        },
                    )*
                }
            }
        }

        impl Message for $name {
            type M = $name;
            type P = NoOption;

            fn get_length_in_bits(&self) -> u32 {
                todo!()
            }

            fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
                use paste::paste;
                paste! {
                    writer.[<write_ $type>]((*self).into())
                }
            }

            fn parse<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::P>) -> Result<Self::M, Error> {
                assert!(parameter.is_none());
                use paste::paste;
                paste! {
                    let result = reader.[<read_ $type>]()?;
                }
                match $name::try_from(result) {
                    Ok(result) => {
                        Ok(result)
                    }
                    Err(_) => {
                        Err(Error::new(ErrorKind::InvalidInput, format!("Cannot parse {}", result)))
                    }
                }
            }
        }
    };
    (enum $bits:expr => $type:ty : $name:ident $([$pat:expr => $branch:ident]),*) => {
        #[derive(Copy, Clone, PartialEq, Debug)]
        #[allow(non_camel_case_types)]
        pub enum $name {
            $(
            $branch,
            )* // Repeat for each `expr` passed to the macro
        }

        impl TryFrom<$type> for $name {
            type Error = ();

            fn try_from(value: $type) -> Result<Self, Self::Error> {
                match value {
                    $(
                        $pat => {
                        Ok($name::$branch)
                    },
                    )* // Repeat for each `expr` passed to the macro
                    _ => {
                    Err(())
                    }
                }
            }
        }

        impl Into<$type> for $name {
            fn into(self) -> $type {
                match self {
                    $(
                        $name::$branch => {
                        $pat
                    },
                    )*
                }
            }
        }

    impl Message for $name {
        type M = $name;
        type P = NoOption;

        fn get_length_in_bits(&self) -> u32 {
            todo!()
        }

        fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
            let x: $type = (*self).into();
            writer.write_u_n($bits, x as u64)
        }

        fn parse<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::P>) -> Result<Self::M, Error> {
            assert!(parameter.is_none());
            let result = reader.read_u_n($bits)?;
            match $name::try_from(result as $type) {
            Ok(result) => {
            Ok(result)
            }
            Err(_) => {
            Err(Error::new(ErrorKind::InvalidInput, format!("Cannot parse {}", result)))
            }
            }
        }
    }
    }
}

#[cfg(test)]
mod test {
    use std::io::{Error, ErrorKind, Read, Write};
    use crate::{Endianess, Message, NoOption, ReadBuffer, WriteBuffer};

    #[test]
    fn test_macro() {
        plc4x_enum!
        [enum u8 : ModbusErrorCode
            [1 =>   ILLEGAL_FUNCTION],
            [2 =>   ILLEGAL_DATA_ADDRESS],
            [3 =>   ILLEGAL_DATA_VALUE],
            [4 =>   SLAVE_DEVICE_FAILURE],
            [5 =>   ACKNOWLEDGE],
            [6 =>   SLAVE_DEVICE_BUSY],
            [7 =>   NEGATIVE_ACKNOWLEDGE],
            [8 =>   MEMORY_PARITY_ERROR],
            [10 =>  GATEWAY_PATH_UNAVAILABLE],
            [11 =>  GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND]
        ];

        plc4x_enum![enum u16 : ModbusErrorCode2
        [1 =>   ILLEGAL_FUNCTION],
        [2 =>   ILLEGAL_DATA_ADDRESS],
        [3 =>   ILLEGAL_DATA_VALUE],
        [4 =>   SLAVE_DEVICE_FAILURE],
        [5 =>   ACKNOWLEDGE],
        [6 =>   SLAVE_DEVICE_BUSY],
        [7 =>   NEGATIVE_ACKNOWLEDGE],
        [8 =>   MEMORY_PARITY_ERROR],
        [10 =>  GATEWAY_PATH_UNAVAILABLE],
        [11 =>  GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND]
    ];

        let x = ModbusErrorCode2::ILLEGAL_FUNCTION;
        let y: u16 = x.into();

        assert_eq!(1, y);
    }

    #[test]
    fn test_u7() {
        // [enum uint 7 ModbusDeviceInformationConformityLevel
        // ['0x01' BASIC_STREAM_ONLY   ]
        // ['0x02' REGULAR_STREAM_ONLY ]
        // ['0x03' EXTENDED_STREAM_ONLY]
        // ]
        plc4x_enum![enum 7 => u8 : ModbusDeviceInformationConformityLevel
            [0x01 => BASIC_STREAM_ONLY],
            [0x02 => REGULAR_STREAM_ONLY ],
            [0x03 => EXTENDED_STREAM_ONLY]
        ];

        let x = ModbusDeviceInformationConformityLevel::BASIC_STREAM_ONLY;
        let y: u8 = x.into();
        assert_eq!(0x01, y);

        let r: Vec<u8> = vec![];
        let mut reader = ReadBuffer::new(Endianess::BigEndian, &*r);

        let y = ModbusDeviceInformationConformityLevel::parse(&mut reader, None);
    }
}