//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package model

import (
    "encoding/xml"
    "errors"
    "io"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type ProjectInstallationIdentifier struct {
    ProjectNumber uint8
    InstallationNumber uint8

}

// The corresponding interface
type IProjectInstallationIdentifier interface {
    spi.Message
    Serialize(io utils.WriteBuffer) error
}


func NewProjectInstallationIdentifier(projectNumber uint8, installationNumber uint8) spi.Message {
    return &ProjectInstallationIdentifier{ProjectNumber: projectNumber, InstallationNumber: installationNumber}
}

func CastIProjectInstallationIdentifier(structType interface{}) IProjectInstallationIdentifier {
    castFunc := func(typ interface{}) IProjectInstallationIdentifier {
        if iProjectInstallationIdentifier, ok := typ.(IProjectInstallationIdentifier); ok {
            return iProjectInstallationIdentifier
        }
        return nil
    }
    return castFunc(structType)
}

func CastProjectInstallationIdentifier(structType interface{}) ProjectInstallationIdentifier {
    castFunc := func(typ interface{}) ProjectInstallationIdentifier {
        if sProjectInstallationIdentifier, ok := typ.(ProjectInstallationIdentifier); ok {
            return sProjectInstallationIdentifier
        }
        if sProjectInstallationIdentifier, ok := typ.(*ProjectInstallationIdentifier); ok {
            return *sProjectInstallationIdentifier
        }
        return ProjectInstallationIdentifier{}
    }
    return castFunc(structType)
}

func (m ProjectInstallationIdentifier) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Simple field (projectNumber)
    lengthInBits += 8

    // Simple field (installationNumber)
    lengthInBits += 8

    return lengthInBits
}

func (m ProjectInstallationIdentifier) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ProjectInstallationIdentifierParse(io *utils.ReadBuffer) (spi.Message, error) {

    // Simple Field (projectNumber)
    projectNumber, _projectNumberErr := io.ReadUint8(8)
    if _projectNumberErr != nil {
        return nil, errors.New("Error parsing 'projectNumber' field " + _projectNumberErr.Error())
    }

    // Simple Field (installationNumber)
    installationNumber, _installationNumberErr := io.ReadUint8(8)
    if _installationNumberErr != nil {
        return nil, errors.New("Error parsing 'installationNumber' field " + _installationNumberErr.Error())
    }

    // Create the instance
    return NewProjectInstallationIdentifier(projectNumber, installationNumber), nil
}

func (m ProjectInstallationIdentifier) Serialize(io utils.WriteBuffer) error {

    // Simple Field (projectNumber)
    projectNumber := uint8(m.ProjectNumber)
    _projectNumberErr := io.WriteUint8(8, (projectNumber))
    if _projectNumberErr != nil {
        return errors.New("Error serializing 'projectNumber' field " + _projectNumberErr.Error())
    }

    // Simple Field (installationNumber)
    installationNumber := uint8(m.InstallationNumber)
    _installationNumberErr := io.WriteUint8(8, (installationNumber))
    if _installationNumberErr != nil {
        return errors.New("Error serializing 'installationNumber' field " + _installationNumberErr.Error())
    }

    return nil
}

func (m *ProjectInstallationIdentifier) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    for {
        token, err := d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "projectNumber":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ProjectNumber = data
            case "installationNumber":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.InstallationNumber = data
            }
        }
    }
}

func (m ProjectInstallationIdentifier) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.ProjectInstallationIdentifier"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ProjectNumber, xml.StartElement{Name: xml.Name{Local: "projectNumber"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.InstallationNumber, xml.StartElement{Name: xml.Name{Local: "installationNumber"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

