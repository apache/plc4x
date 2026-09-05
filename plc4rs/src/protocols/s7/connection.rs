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

use nom::number::complete::{be_u16, be_u8};
use nom::IResult;

#[derive(Debug, Clone)]
pub struct COTPConnectionRequest {
    dst_ref: u16,
    src_ref: u16,
    class: u8,
    parameters: Vec<ConnectionParameter>,
}

#[derive(Debug, Clone)]
pub enum ConnectionParameter {
    TpduSize(u8),
    SrcTsap(u16),
    DstTsap(u16),
}

#[derive(Debug, Clone)]
pub struct COTPConnectionResponse {
    dst_ref: u16,
    src_ref: u16,
    class: u8,
    parameters: Vec<ConnectionParameter>,
}

impl COTPConnectionRequest {
    pub fn parse(input: &[u8]) -> IResult<&[u8], Self> {
        let (input, dst_ref) = be_u16(input)?;
        let (input, src_ref) = be_u16(input)?;
        let (input, class) = be_u8(input)?;
        let (input, parameters) = ConnectionParameter::parse_all(input)?;

        Ok((
            input,
            COTPConnectionRequest {
                dst_ref,
                src_ref,
                class,
                parameters,
            },
        ))
    }
}

impl COTPConnectionResponse {
    pub fn parse(input: &[u8]) -> IResult<&[u8], Self> {
        let (input, dst_ref) = be_u16(input)?;
        let (input, src_ref) = be_u16(input)?;
        let (input, class) = be_u8(input)?;
        let (input, parameters) = ConnectionParameter::parse_all(input)?;

        Ok((
            input,
            COTPConnectionResponse {
                dst_ref,
                src_ref,
                class,
                parameters,
            },
        ))
    }
}

impl ConnectionParameter {
    fn parse_all(input: &[u8]) -> IResult<&[u8], Vec<Self>> {
        let mut parameters = Vec::new();
        let mut remaining = input;

        while !remaining.is_empty() {
            let (input, param) = Self::parse(remaining)?;
            parameters.push(param);
            remaining = input;
        }

        Ok((remaining, parameters))
    }

    fn parse(input: &[u8]) -> IResult<&[u8], Self> {
        let (input, param_code) = be_u8(input)?;
        let (input, param_length) = be_u8(input)?;

        match param_code {
            0xC0 => {
                let (input, size) = be_u8(input)?;
                Ok((input, ConnectionParameter::TpduSize(size)))
            }
            0xC1 => {
                let (input, tsap) = be_u16(input)?;
                Ok((input, ConnectionParameter::SrcTsap(tsap)))
            }
            0xC2 => {
                let (input, tsap) = be_u16(input)?;
                Ok((input, ConnectionParameter::DstTsap(tsap)))
            }
            _ => Err(nom::Err::Error(nom::error::Error::new(
                input,
                nom::error::ErrorKind::Tag,
            ))),
        }
    }
}
