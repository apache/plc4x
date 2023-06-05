<!--
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
  -->

# Profinet

## Connection

In general a connection is ethernet based and therefore needs only a MAC address.
We could simply send out every packet on every network device and listen to all for incoming packets, but that would be quite inefficient.

During discovery, we get information like MAC address, IP address, Vendor-Id, Device-ID, but in theory the mac address should be enough.
So I am thinking of two options: 

1. The connection string contains an ip address
2. The connection string doesn't contain an ip address

If an IP address is present, we could theoretically simply use that ip address and try to find a local network device able to reach this.
However on a system there might be multiple devices able to reach that target. 
In my case my primary internet connection goes via Wi-Fi and a secondary goes via cable (So I can pass the full network to a VM and still be online via Wi-Fi). 
We want to use the fastest option.

So I was thinking of sending an ARP (Address Resolution Protocol) request for resolving the remote address and to take the ip address returned first. 

As soon as we have found the local network device then we can send a PN_DCP 