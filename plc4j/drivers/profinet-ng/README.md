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

PR IO generally works via raw Ethernet communication, so in general the Source and Target MAC address are the only bits of information needed. 
However, does the Profinet family of protocols also have a number of assisting protocols. 
The probably most important one is the Profinet CM (Context Manager) protocol, which is needed for establishing connections.
This protocol is based on DCE/RPC messages, which are transported via UDP.
In order to do this, an IP address is needed.

Therefore, the simplest connection string for a Profinet device will be:

    profinet://192.168.24.31

The driver will do an ARP lookup in order to get the MAC address of the remote device and in order to select the right device.

While most PN devices will be configured by the engineering tools and hereby also will be assigned an IP address, some devices need to be configured by the Profinet Master on connection establishment.
For these devices we will have to come up with a mechanism, where the connection string will contain a mac-address and an ip-address as option. 
The driver would be responsible for configuring the device before establishing a connection. 
A connection string would look like this in this case:

    profinet://60:6d:3c:3d:a9:a3?ip-address=192.168.24.31

## GSD Device Profiles

A Profinet device is described by a GSD file. 
A GSD file is an XML file containing a description of the device.
Only with this information can we provide browsing functionality, verify requests and reduce the amount of information passed to the driver to an absolute minimum.

There are multiple coordinates needed to identify which resources a PN device offers.

The most important ones are the vendor-id and the device-id. This tuple uniquely identifies a PN device or a family of devices.

When connecting to a remote device, we execute a PN-DCP Identification request is sent to the device. 
This is the same type of request used in the Auto-Discovery of PLC4X to find PN devices. 
The main difference in this case however is that we don't send it to the broadcast MAC address, but send it to the devices MAC address instead.
In the response we can get the vendor-id and the device-id along with a number of additional information (IP settings, name of the device, type of device).
With this information we then iterate over the GSD files in our "gsd directory" (Which defaults to "~/.gsd"). 
As soon as we found an XML file with a matching pair of vendor-id and product-id we use that device profile for this connection.

A PN device always had one fixed component, the so-called DAP (Device Access Point). 
This is always located in "slot 0" and you can think of this as the IO component that provides the connection between the network and the devices, using the Profinet protocol.
Most device profiles only have one DAP and in this case, we simply use that. 
However, there are some devices where the manufacturer updated the device over time and didn't give the new device a new product id.
In this case, device profiles can contain a variety of DAPs. 

If we are connecting to such a device it is important to know which DAP the devices provides.

Next to the DAP in Slot 0 come the actual devices. 
In some devices they are integrated into the same housing, only logically having multiple slots.
This is the case for my Advantec Adam modules belong to this category the same way my Siemens Simodode Pro V PN device does.

Others however allow understanding this concept a bit better though. 
So-called Bus-Couplers for example provide PN access to IO devices. 
Here you have the bus-coupler in slot 0 (which is the left-most element) and then you add on IO modules by sliding them onto the right. 
When adding a new module to the existing, the new module is always added to the right.

In PN terms each of these IO devices take a slot. 
The left-most IO-device (next to the bus-coupler) is in slot 1, the next one in 2 and so on.

As you can imagine that there is a great difference in how these modules work and what information they provide or are able to accept. 
An 8-times digital Input definitely has a different set of information as an output or an analog input.
These differences are modeled in so-called "Modules" in the GSD files. 
As depending on the used DAP there might be differences in which modules are supported in which slot the DAP element contains a list of "usable modules".
This describes which types of modules are allowed on which slot.

In case of a Wago bus-coupler for example this looks like this:

    <UseableModules>
		<ModuleItemRef ModuleItemTarget="0606_0000" AllowedInSlots="1..250"/>
		<ModuleItemRef ModuleItemTarget="0610_0000" AllowedInSlots="1..250"/>
		<ModuleItemRef ModuleItemTarget="0611_0000" AllowedInSlots="1..250"/>
		<ModuleItemRef ModuleItemTarget="00DI_02DIA_0000" AllowedInSlots="1..250"/>
		<ModuleItemRef ModuleItemTarget="0400_0000" AllowedInSlots="1..250"/>
		<ModuleItemRef ModuleItemTarget="0401_0000" AllowedInSlots="1..250"/>
		<ModuleItemRef ModuleItemTarget="0402_0000" AllowedInSlots="1..250"/>
		<ModuleItemRef ModuleItemTarget="0403_0000" AllowedInSlots="1..250"/>
		<ModuleItemRef ModuleItemTarget="0405_0000" AllowedInSlots="1..250"/>
		<ModuleItemRef ModuleItemTarget="0406_0000" AllowedInSlots="1..250"/>
        ...
    <UseableModules>

You can see that the bus-coupler generally allows adding 250 devices. 
Each entry in this list being one IO module Product Wago has to offer.
"0606_0000" being a "Power Supply 24 V DC, max. 1.0 A; Ex i; diagnostics" and "0403_0000" being a "4-channel digital input; 24 V DC; 0.2 ms input filter; 2- to 3-conductor connection; high-side switching"

## Auto-Configuring the connection 

So in order to know which type of module is present in which slot, we need to ask the PN device.

This is done by reading the "RealIdentificationData" from the device using PN-CM (DCE/RPC, UDP) as part of the connection establishment process.

This information gives us detailed information on how many slots are being used, which slot uses which module and which sub-slots this device has and what subModules these map to.
From the information on which module is used in slot 0 we also know which DAP is being used.

With this information can we fully provide PLC4X browse functionality, and we are able to validate request before sending anything to the device. 


