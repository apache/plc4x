/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package gsdml reads GSDML files, the XML device descriptions (ISO 15745 profiles) a Profinet IO
// controller needs to know a device's modules, submodules, I/O data layout and parameters.
//
// It is the Go counterpart of plc4j profinet-ng's gsdml package, covering what a controller uses:
// device identity, device access points, modules, submodules and their I/O data and parameter
// records, value lists and texts. Unlike plc4j it also resolves submodules a module only
// references (GSDML 2.25 and later: SubmoduleList plus UseableSubmodules), which is how large
// modular stations such as WAGO 750 describe their I/O.
package gsdml

import (
	"encoding/xml"
)

// Profile is the document root, ISO15745Profile.
type Profile struct {
	XMLName       xml.Name      `xml:"ISO15745Profile"`
	ProfileHeader ProfileHeader `xml:"ProfileHeader"`
	ProfileBody   ProfileBody   `xml:"ProfileBody"`
}

type ProfileHeader struct {
	ProfileIdentification string `xml:"ProfileIdentification"`
	ProfileRevision       string `xml:"ProfileRevision"`
	ProfileName           string `xml:"ProfileName"`
	ProfileSource         string `xml:"ProfileSource"`
	ProfileClassID        string `xml:"ProfileClassID"`
}

type ProfileBody struct {
	DeviceIdentity     DeviceIdentity     `xml:"DeviceIdentity"`
	DeviceFunction     DeviceFunction     `xml:"DeviceFunction"`
	ApplicationProcess ApplicationProcess `xml:"ApplicationProcess"`
}

type DeviceIdentity struct {
	VendorID   string  `xml:"VendorID,attr"`
	DeviceID   string  `xml:"DeviceID,attr"`
	InfoText   TextRef `xml:"InfoText"`
	VendorName Value   `xml:"VendorName"`
}

type DeviceFunction struct {
	Family Family `xml:"Family"`
}

type Family struct {
	MainFamily    string `xml:"MainFamily,attr"`
	ProductFamily string `xml:"ProductFamily,attr"`
}

// TextRef points into the external text list.
type TextRef struct {
	TextID string `xml:"TextId,attr"`
}

type Value struct {
	Value string `xml:"Value,attr"`
}

type ApplicationProcess struct {
	DeviceAccessPoints []DeviceAccessPointItem `xml:"DeviceAccessPointList>DeviceAccessPointItem"`
	Modules            []ModuleItem            `xml:"ModuleList>ModuleItem"`
	Submodules         []SubmoduleItem         `xml:"SubmoduleList>SubmoduleItem"`
	ValueItems         []ValueItem             `xml:"ValueList>ValueItem"`
	ExternalTextList   ExternalTextList        `xml:"ExternalTextList"`
}

// DeviceAccessPointItem (DAP) is the head module in slot 0: the device's interface to the bus.
type DeviceAccessPointItem struct {
	ID                    string          `xml:"ID,attr"`
	PNIOVersion           string          `xml:"PNIO_Version,attr"`
	PhysicalSlots         string          `xml:"PhysicalSlots,attr"`
	ModuleIdentNumber     string          `xml:"ModuleIdentNumber,attr"`
	MinDeviceInterval     int             `xml:"MinDeviceInterval,attr"`
	DNSCompatibleName     string          `xml:"DNS_CompatibleName,attr"`
	FixedInSlots          string          `xml:"FixedInSlots,attr"`
	ObjectUUIDLocalIndex  int             `xml:"ObjectUUID_LocalIndex,attr"`
	DeviceAccessSupported bool            `xml:"DeviceAccessSupported,attr"`
	ModuleInfo            ModuleInfo      `xml:"ModuleInfo"`
	IOConfigData          IOConfigData    `xml:"IOConfigData"`
	UseableModules        []ModuleItemRef `xml:"UseableModules>ModuleItemRef"`
	VirtualSubmodules     []SubmoduleItem `xml:"VirtualSubmoduleList>VirtualSubmoduleItem"`
	SystemDefined         SystemDefined   `xml:"SystemDefinedSubmoduleList"`
}

type ModuleInfo struct {
	Name            TextRef `xml:"Name"`
	InfoText        TextRef `xml:"InfoText"`
	VendorName      Value   `xml:"VendorName"`
	OrderNumber     Value   `xml:"OrderNumber"`
	HardwareRelease Value   `xml:"HardwareRelease"`
	SoftwareRelease Value   `xml:"SoftwareRelease"`
}

type IOConfigData struct {
	MaxInputLength  int `xml:"MaxInputLength,attr"`
	MaxOutputLength int `xml:"MaxOutputLength,attr"`
}

type ModuleItemRef struct {
	ModuleItemTarget string `xml:"ModuleItemTarget,attr"`
	AllowedInSlots   string `xml:"AllowedInSlots,attr"`
	FixedInSlots     string `xml:"FixedInSlots,attr"`
	UsedInSlots      string `xml:"UsedInSlots,attr"`
}

// SystemDefined are the interface and port submodules of a DAP.
type SystemDefined struct {
	Interfaces []InterfaceSubmoduleItem `xml:"InterfaceSubmoduleItem"`
	Ports      []PortSubmoduleItem      `xml:"PortSubmoduleItem"`
}

type InterfaceSubmoduleItem struct {
	ID                   string `xml:"ID,attr"`
	SubmoduleIdentNumber string `xml:"SubmoduleIdentNumber,attr"`
	SubslotNumber        int    `xml:"SubslotNumber,attr"`
	TextID               string `xml:"TextId,attr"`
	SupportedRTClasses   string `xml:"SupportedRT_Classes,attr"`
	SupportedProtocols   string `xml:"SupportedProtocols,attr"`
}

type PortSubmoduleItem struct {
	ID                   string `xml:"ID,attr"`
	SubmoduleIdentNumber string `xml:"SubmoduleIdentNumber,attr"`
	SubslotNumber        int    `xml:"SubslotNumber,attr"`
	TextID               string `xml:"TextId,attr"`
	MaxPortRxDelay       int    `xml:"MaxPortRxDelay,attr"`
	MaxPortTxDelay       int    `xml:"MaxPortTxDelay,attr"`
}

// ModuleItem is a module that goes in a slot.
type ModuleItem struct {
	ID                string             `xml:"ID,attr"`
	ModuleIdentNumber string             `xml:"ModuleIdentNumber,attr"`
	PhysicalSubslots  string             `xml:"PhysicalSubslots,attr"`
	ModuleInfo        ModuleInfo         `xml:"ModuleInfo"`
	VirtualSubmodules []SubmoduleItem    `xml:"VirtualSubmoduleList>VirtualSubmoduleItem"`
	UseableSubmodules []SubmoduleItemRef `xml:"UseableSubmodules>SubmoduleItemRef"`
}

// SubmoduleItemRef names a submodule of the SubmoduleList a module can carry.
type SubmoduleItemRef struct {
	SubmoduleItemTarget string `xml:"SubmoduleItemTarget,attr"`
	AllowedInSubslots   string `xml:"AllowedInSubslots,attr"`
	UsedInSubslots      string `xml:"UsedInSubslots,attr"`
	FixedInSubslots     string `xml:"FixedInSubslots,attr"`
}

// SubmoduleItem is a VirtualSubmoduleItem or a SubmoduleList SubmoduleItem: the unit that has I/O
// data and parameters.
type SubmoduleItem struct {
	ID                   string           `xml:"ID,attr"`
	SubmoduleIdentNumber string           `xml:"SubmoduleIdentNumber,attr"`
	FixedInSubslots      string           `xml:"FixedInSubslots,attr"`
	MayIssueProcessAlarm bool             `xml:"MayIssueProcessAlarm,attr"`
	WriteableIMRecords   string           `xml:"Writeable_IM_Records,attr"`
	IOData               IOData           `xml:"IOData"`
	Records              []RecordDataItem `xml:"RecordDataList>ParameterRecordDataItem"`
	ModuleInfo           ModuleInfo       `xml:"ModuleInfo"`
}

// IOData is a submodule's cyclic data. IOPS and IOCS (provider and consumer status) default to one
// byte each when the attributes are absent.
type IOData struct {
	IOPSLength *int        `xml:"IOPS_Length,attr"`
	IOCSLength *int        `xml:"IOCS_Length,attr"`
	Inputs     []IODataSet `xml:"Input"`
	Outputs    []IODataSet `xml:"Output"`
}

type IODataSet struct {
	Consistency string     `xml:"Consistency,attr"`
	DataItems   []DataItem `xml:"DataItem"`
}

// DataItem is one value in a submodule's I/O data, in the order they are laid out.
type DataItem struct {
	DataType     string        `xml:"DataType,attr"`
	TextID       string        `xml:"TextId,attr"`
	Length       int           `xml:"Length,attr"`
	UseAsBits    bool          `xml:"UseAsBits,attr"`
	BitDataItems []BitDataItem `xml:"BitDataItem"`
}

type BitDataItem struct {
	BitOffset int    `xml:"BitOffset,attr"`
	TextID    string `xml:"TextId,attr"`
}

// RecordDataItem is a parameter record written to the submodule while connecting.
type RecordDataItem struct {
	Index            int           `xml:"Index,attr"`
	Length           int           `xml:"Length,attr"`
	TransferSequence int           `xml:"TransferSequence,attr"`
	Name             TextRef       `xml:"Name"`
	Consts           []RecordConst `xml:"Const"`
	Refs             []RecordRef   `xml:"Ref"`
}

// RecordConst is fixed record content: comma separated bytes at a byte offset.
type RecordConst struct {
	Data       string `xml:"Data,attr"`
	ByteOffset int    `xml:"ByteOffset,attr"`
}

// RecordRef is one parameter in a record, with its default value.
type RecordRef struct {
	ID              string `xml:"ID,attr"`
	DataType        string `xml:"DataType,attr"`
	ByteOffset      int    `xml:"ByteOffset,attr"`
	BitOffset       int    `xml:"BitOffset,attr"`
	BitLength       int    `xml:"BitLength,attr"`
	Length          int    `xml:"Length,attr"`
	DefaultValue    string `xml:"DefaultValue,attr"`
	AllowedValues   string `xml:"AllowedValues,attr"`
	Changeable      *bool  `xml:"Changeable,attr"`
	Visible         *bool  `xml:"Visible,attr"`
	TextID          string `xml:"TextId,attr"`
	ValueItemTarget string `xml:"ValueItemTarget,attr"`
}

// ValueItem names the values a parameter can take.
type ValueItem struct {
	ID          string   `xml:"ID,attr"`
	Assignments []Assign `xml:"Assignments>Assign"`
}

type Assign struct {
	Content string `xml:"Content,attr"`
	TextID  string `xml:"TextId,attr"`
}

type ExternalTextList struct {
	PrimaryLanguage []Text `xml:"PrimaryLanguage>Text"`
}

type Text struct {
	TextID string `xml:"TextId,attr"`
	Value  string `xml:"Value,attr"`
}
