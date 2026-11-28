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

package gsdml

import (
	"bufio"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

const (
	// deviceProfileIdentification and deviceProfileClass mark a GSDML device profile, the check
	// plc4j's ProfinetConfiguration.getGsdProfile makes.
	deviceProfileIdentification = "PROFINET Device Profile"
	deviceProfileClass          = "Device"
)

// Parse reads a GSDML document.
func Parse(reader io.Reader) (*Profile, error) {
	var profile Profile
	if err := newDecoder(reader).Decode(&profile); err != nil {
		return nil, errors.Wrap(err, "error parsing the GSDML document")
	}
	return &profile, nil
}

// ParseFile reads a GSDML file.
func ParseFile(path string) (*Profile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	profile, err := Parse(bufio.NewReaderSize(file, 1<<16))
	if err != nil {
		return nil, errors.Wrapf(err, "in %s", path)
	}
	return profile, nil
}

// IsDeviceProfile reports whether the document is a Profinet device description.
func (p *Profile) IsDeviceProfile() bool {
	return strings.TrimSpace(p.ProfileHeader.ProfileIdentification) == deviceProfileIdentification &&
		strings.TrimSpace(p.ProfileHeader.ProfileClassID) == deviceProfileClass
}

// Identity is the vendor and device id the profile describes.
func (p *Profile) Identity() (vendorID, deviceID uint16, err error) {
	return parseIdentity(p.ProfileBody.DeviceIdentity.VendorID, p.ProfileBody.DeviceIdentity.DeviceID)
}

func parseIdentity(vendor, device string) (uint16, uint16, error) {
	vendorID, err := ParseNumber(vendor)
	if err != nil || vendorID > 0xFFFF {
		return 0, 0, errors.Errorf("bad VendorID %q", vendor)
	}
	deviceID, err := ParseNumber(device)
	if err != nil || deviceID > 0xFFFF {
		return 0, 0, errors.Errorf("bad DeviceID %q", device)
	}
	return uint16(vendorID), uint16(deviceID), nil
}

// ParseNumber reads a GSDML number: hexadecimal with 0x, decimal otherwise.
func ParseNumber(text string) (uint64, error) {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "0x") || strings.HasPrefix(text, "0X") {
		return strconv.ParseUint(text[2:], 16, 64)
	}
	return strconv.ParseUint(text, 10, 64)
}

// peekIdentity streams a file only as far as its DeviceIdentity, so a directory full of large
// GSDML files can be searched without parsing each one whole.
func peekIdentity(path string) (vendorID, deviceID uint16, ok bool) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, false
	}
	defer func() { _ = file.Close() }()
	decoder := newDecoder(bufio.NewReader(file))
	sawProfile := false
	for {
		token, err := decoder.Token()
		if err != nil {
			return 0, 0, false
		}
		start, isStart := token.(xml.StartElement)
		if !isStart {
			continue
		}
		switch start.Name.Local {
		case "ISO15745Profile":
			sawProfile = true
		case "DeviceIdentity":
			if !sawProfile {
				return 0, 0, false
			}
			var vendor, device string
			for _, attribute := range start.Attr {
				switch attribute.Name.Local {
				case "VendorID":
					vendor = attribute.Value
				case "DeviceID":
					device = attribute.Value
				}
			}
			vendorID, deviceID, err := parseIdentity(vendor, device)
			return vendorID, deviceID, err == nil
		case "ApplicationProcess":
			// DeviceIdentity comes before this; a file without one isn't a device profile.
			return 0, 0, false
		}
	}
}

// ExpandHome turns a leading ~ into the user's home directory, as plc4j does for gsd-directory.
func ExpandHome(path string) string {
	if path == "~" || strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[1:])
		}
	}
	return path
}

// FindInDirectory returns the device profile for a vendor and device id from a directory of GSDML
// files (*.xml). When several match - a vendor ships a new GSDML version beside the old one - the
// last by file name wins, which with the standard GSDML-Vx.y-Vendor-Device-yyyymmdd.xml naming is
// the newest date of the newest schema version. plc4j returns whichever the directory lists first.
func FindInDirectory(directory string, vendorID, deviceID uint16) (*Profile, string, error) {
	directory = ExpandHome(directory)
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, "", errors.Wrapf(err, "the GSDML directory %s can't be read", directory)
	}
	match := ""
	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".xml") {
			continue
		}
		path := filepath.Join(directory, entry.Name())
		if vendor, device, ok := peekIdentity(path); ok && vendor == vendorID && device == deviceID && entry.Name() > filepath.Base(match) {
			match = path
		}
	}
	if match == "" {
		return nil, "", errors.Errorf("no GSDML file in %s describes vendor 0x%04X device 0x%04X", directory, vendorID, deviceID)
	}
	profile, err := ParseFile(match)
	if err != nil {
		return nil, "", err
	}
	if !profile.IsDeviceProfile() {
		return nil, "", errors.Errorf("%s is not a Profinet device profile", match)
	}
	return profile, match, nil
}

// newDecoder is an XML decoder that also reads ISO-8859-1, the encoding most GSDML files declare.
func newDecoder(reader io.Reader) *xml.Decoder {
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charsetReader
	return decoder
}

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch strings.ToLower(strings.TrimSpace(charset)) {
	case "utf-8", "utf8", "us-ascii", "ascii":
		return input, nil
	case "iso-8859-1", "iso8859-1", "latin1", "latin-1", "l1":
		return &latin1Reader{source: bufio.NewReader(input)}, nil
	case "windows-1252", "cp1252":
		// Close enough for the names and texts in a GSDML file: the two only differ in 0x80-0x9F.
		return &latin1Reader{source: bufio.NewReader(input)}, nil
	}
	return nil, errors.Errorf("unsupported GSDML encoding %q", charset)
}

// latin1Reader turns ISO-8859-1 bytes into UTF-8: each byte is the code point of the same value.
type latin1Reader struct {
	source  *bufio.Reader
	pending []byte
}

func (r *latin1Reader) Read(p []byte) (int, error) {
	n := 0
	for n < len(p) {
		if len(r.pending) > 0 {
			copied := copy(p[n:], r.pending)
			r.pending = r.pending[copied:]
			n += copied
			continue
		}
		b, err := r.source.ReadByte()
		if err != nil {
			if n > 0 {
				return n, nil
			}
			return 0, err
		}
		if b < 0x80 {
			p[n] = b
			n++
			continue
		}
		r.pending = []byte{0xC0 | b>>6, 0x80 | b&0x3F}
	}
	return n, nil
}
