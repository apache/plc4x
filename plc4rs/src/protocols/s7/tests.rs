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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_tpkt_packet_parse() {
        let input = &[
            0x03, 0x00, 0x00, 0x1A, // TPKT header
            0x02, 0xF0, 0x80, // COTP header
            0x32, 0x01, 0x00, 0x00, // S7 header start
            0x00, 0x01, 0x00, 0x08, 0x00,
            0x00, // S7 header end
                  // ... payload ...
        ];

        let (remaining, packet) = TPKTPacket::parse(input).unwrap();
        assert_eq!(packet.protocol_id, 0x03);
        assert_eq!(packet.length, 26);
        assert!(remaining.is_empty());
    }

    #[test]
    fn test_s7_message_parse() {
        // Add S7 message parsing test
    }

    #[test]
    fn test_connection_request_parse() {
        let input = &[
            0x00, 0x0C, // Destination Reference
            0x00, 0x10, // Source Reference
            0x00, // Class
            0xC1, 0x02, 0x01, 0x00, // Source TSAP
            0xC2, 0x02, 0x01, 0x02, // Destination TSAP
        ];

        let (remaining, request) = COTPConnectionRequest::parse(input).unwrap();
        assert_eq!(request.dst_ref, 12);
        assert_eq!(request.src_ref, 16);
        assert_eq!(request.class, 0);
        assert!(remaining.is_empty());
    }

    #[test]
    fn test_s7_payload_parse() {
        let input = &[
            0x00, // Return code (Success)
            0x04, // Transport size (Word)
            0x00, 0x02, // Length
            0x12, 0x34, // Data
        ];

        let (remaining, payload) = S7Payload::parse(input, 6).unwrap();
        assert_eq!(payload.items.len(), 1);
        assert_eq!(payload.items[0].data, vec![0x12, 0x34]);
        assert!(remaining.is_empty());
    }

    #[test]
    fn test_parameter_item_parse() {
        let input = &[
            0x02, // VarSpec (Byte)
            0x0A, // Length
            0x10, // Syntax ID (S7Any)
            0x02, // Transport size
            0x00, 0x01, // DB number
            0x84, // Area (DataBlocks)
            0x00, 0x00, 0x00, 0x00, // Start address
        ];

        let (remaining, item) = ParameterItem::parse(input).unwrap();
        assert_eq!(item.addr_length, 10);
        assert_eq!(item.db_number, 1);
        assert!(remaining.is_empty());
    }
}
