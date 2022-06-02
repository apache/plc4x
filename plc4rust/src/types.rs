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

// [type ModbusPDUReadFileRecordRequestItem
//     [simple     uint 8     referenceType]
//     [simple     uint 16    fileNumber   ]
//     [simple     uint 16    recordNumber ]
//     [simple     uint 16    recordLength ]
// ]
// Necessary imports
// use std::io::{Error, ErrorKind, Read, Write};
// use crate::{Message, NoOption, ReadBuffer, WriteBuffer};
#[macro_export]
macro_rules! plc4x_type {
    (type $name:ident $([simple $type:ty : $parameter:ident]),*) => {
        #[derive(PartialEq,Eq,Clone,Debug)]
        pub struct $name {
            $(
                $parameter: $type,
            )*
        }

        impl Message for $name {
            type M = $name;
            type P = NoOption;

            fn get_length_in_bits(&self) -> u32 {
                todo!();
            }

            fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, std::io::Error> {
                use paste::paste;
                paste! {
                    $(
                        writer.[<write_ $type>](self.$parameter)?;
                    )*
                }
                Ok(0)
                // let mut size = writer.write_u8(self.reference_type)?;
                // size += writer.write_u16(self.file_number)?;
                // size += writer.write_u16(self.record_number)?;
                // size += writer.write_u16(self.record_length)?;
                // Ok(size)
            }

            fn parse<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::P>) -> Result<Self::M, std::io::Error> {
                use paste::paste;
                paste! {
                    $(
                        let $parameter = reader.[<read_ $type>]()?;
                    )*
                }

                Ok(Self::M {
                    $(
                        $parameter,
                    )*
                })
            }
        }
    }
}

#[cfg(test)]
mod test {

    use std::io::{Error, ErrorKind, Read, Write};
    use crate::{Message, NoOption, ReadBuffer, WriteBuffer};

    #[test]
    fn test() {
        // If this doesnt throw, its good...
        plc4x_type!
        [type ModbusPDUReadFileRecordRequestItem
            [simple u8 : referenceType],
            [simple u16 : fileNumber],
            [simple u16 : recordNumber],
            [simple u16 : recordLength]
        ];
    }
}




// #[derive(PartialEq,Eq,Clone,Debug)]
// pub struct ModbusPDUReadFileRecordRequestItem {
//     reference_type: u8,
//     file_number: u16,
//     record_number: u16,
//     record_length: u16,
// }
//
// impl Message for ModbusPDUReadFileRecordRequestItem {
//     type M = ModbusPDUReadFileRecordRequestItem;
//     type P = NoOption;
//
//     fn get_length_in_bits(&self) -> u32 {
//         56*8
//     }
//
//     fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, std::io::Error> {
//         let mut size = writer.write_u8(self.reference_type)?;
//         size += writer.write_u16(self.file_number)?;
//         size += writer.write_u16(self.record_number)?;
//         size += writer.write_u16(self.record_length)?;
//         Ok(size)
//     }
//
//     fn parse<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::P>) -> Result<Self::M, std::io::Error> {
//         let reference_type = reader.read_u8()?;
//         let file_number = reader.read_u16()?;
//         let record_number = reader.read_u16()?;
//         let record_length = reader.read_u16()?;
//
//         Ok(Self::M {
//             reference_type,
//             file_number,
//             record_number,
//             record_length,
//         })
//     }
// }