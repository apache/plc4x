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
            "errors"
            "github.com/apache/plc4x/plc4go/internal/plc4go/spi/values"
            "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
            api "github.com/apache/plc4x/plc4go/pkg/plc4go/values"
)

func KnxDatapointParse(io *utils.ReadBuffer, datapointType IKnxDatapointType) (api.PlcValue, error) {
    switch {
        case datapointType == KnxDatapointType.DPT_Switch: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_Bool: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_Enable: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_Ramp: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_Alarm: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_BinaryValue: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_Step: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_UpDown: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_OpenClose: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_Start: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_State: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_Invert: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_DimSendStyle: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_InputSource: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_Reset: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_Ack: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_Trigger: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_Occupancy: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_Window_Door: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_LogicalFunction: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_Scene_AB: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_ShutterBlinds_Mode: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_DayNight: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_Heat_Cool: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_Switch_Control: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            Control, _ControlErr := io.ReadBit()
            if _ControlErr != nil {
                return nil, errors.New("Error parsing 'Control' field " + _ControlErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Control)

            // Simple Field (On)
            On, _OnErr := io.ReadBit()
            if _OnErr != nil {
                return nil, errors.New("Error parsing 'On' field " + _OnErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(On)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Bool_Control: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            Control, _ControlErr := io.ReadBit()
            if _ControlErr != nil {
                return nil, errors.New("Error parsing 'Control' field " + _ControlErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Control)

            // Simple Field (True)
            True, _TrueErr := io.ReadBit()
            if _TrueErr != nil {
                return nil, errors.New("Error parsing 'True' field " + _TrueErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(True)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Enable_Control: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            Control, _ControlErr := io.ReadBit()
            if _ControlErr != nil {
                return nil, errors.New("Error parsing 'Control' field " + _ControlErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Control)

            // Simple Field (Enable)
            Enable, _EnableErr := io.ReadBit()
            if _EnableErr != nil {
                return nil, errors.New("Error parsing 'Enable' field " + _EnableErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Enable)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Ramp_Control: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            Control, _ControlErr := io.ReadBit()
            if _ControlErr != nil {
                return nil, errors.New("Error parsing 'Control' field " + _ControlErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Control)

            // Simple Field (Ramp)
            Ramp, _RampErr := io.ReadBit()
            if _RampErr != nil {
                return nil, errors.New("Error parsing 'Ramp' field " + _RampErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Ramp)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Alarm_Control: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            Control, _ControlErr := io.ReadBit()
            if _ControlErr != nil {
                return nil, errors.New("Error parsing 'Control' field " + _ControlErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Control)

            // Simple Field (Alarm)
            Alarm, _AlarmErr := io.ReadBit()
            if _AlarmErr != nil {
                return nil, errors.New("Error parsing 'Alarm' field " + _AlarmErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Alarm)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_BinaryValue_Control: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            Control, _ControlErr := io.ReadBit()
            if _ControlErr != nil {
                return nil, errors.New("Error parsing 'Control' field " + _ControlErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Control)

            // Simple Field (High)
            High, _HighErr := io.ReadBit()
            if _HighErr != nil {
                return nil, errors.New("Error parsing 'High' field " + _HighErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(High)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Step_Control: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            Control, _ControlErr := io.ReadBit()
            if _ControlErr != nil {
                return nil, errors.New("Error parsing 'Control' field " + _ControlErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Control)

            // Simple Field (Increase)
            Increase, _IncreaseErr := io.ReadBit()
            if _IncreaseErr != nil {
                return nil, errors.New("Error parsing 'Increase' field " + _IncreaseErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Increase)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Direction1_Control: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            Control, _ControlErr := io.ReadBit()
            if _ControlErr != nil {
                return nil, errors.New("Error parsing 'Control' field " + _ControlErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Control)

            // Simple Field (Down)
            Down, _DownErr := io.ReadBit()
            if _DownErr != nil {
                return nil, errors.New("Error parsing 'Down' field " + _DownErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Down)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Direction2_Control: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            Control, _ControlErr := io.ReadBit()
            if _ControlErr != nil {
                return nil, errors.New("Error parsing 'Control' field " + _ControlErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Control)

            // Simple Field (Close)
            Close, _CloseErr := io.ReadBit()
            if _CloseErr != nil {
                return nil, errors.New("Error parsing 'Close' field " + _CloseErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Close)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Start_Control: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            Control, _ControlErr := io.ReadBit()
            if _ControlErr != nil {
                return nil, errors.New("Error parsing 'Control' field " + _ControlErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Control)

            // Simple Field (Start)
            Start, _StartErr := io.ReadBit()
            if _StartErr != nil {
                return nil, errors.New("Error parsing 'Start' field " + _StartErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Start)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_State_Control: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            Control, _ControlErr := io.ReadBit()
            if _ControlErr != nil {
                return nil, errors.New("Error parsing 'Control' field " + _ControlErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Control)

            // Simple Field (Active)
            Active, _ActiveErr := io.ReadBit()
            if _ActiveErr != nil {
                return nil, errors.New("Error parsing 'Active' field " + _ActiveErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Active)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Invert_Control: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            Control, _ControlErr := io.ReadBit()
            if _ControlErr != nil {
                return nil, errors.New("Error parsing 'Control' field " + _ControlErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Control)

            // Simple Field (Inverted)
            Inverted, _InvertedErr := io.ReadBit()
            if _InvertedErr != nil {
                return nil, errors.New("Error parsing 'Inverted' field " + _InvertedErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Inverted)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Control_Dimming: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Increase)
            Increase, _IncreaseErr := io.ReadBit()
            if _IncreaseErr != nil {
                return nil, errors.New("Error parsing 'Increase' field " + _IncreaseErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Increase)

            // Simple Field (Stepcode)
            Stepcode, _StepcodeErr := io.ReadUint8(3)
            if _StepcodeErr != nil {
                return nil, errors.New("Error parsing 'Stepcode' field " + _StepcodeErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Stepcode)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Control_Blinds: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Down)
            Down, _DownErr := io.ReadBit()
            if _DownErr != nil {
                return nil, errors.New("Error parsing 'Down' field " + _DownErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Down)

            // Simple Field (Stepcode)
            Stepcode, _StepcodeErr := io.ReadUint8(3)
            if _StepcodeErr != nil {
                return nil, errors.New("Error parsing 'Stepcode' field " + _StepcodeErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Stepcode)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Char_ASCII: // STRING

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadString(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSTRING(value), nil
        case datapointType == KnxDatapointType.DPT_Char_8859_1: // STRING

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadString(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSTRING(value), nil
        case datapointType == KnxDatapointType.DPT_Scaling: // USINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case datapointType == KnxDatapointType.DPT_Angle: // USINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case datapointType == KnxDatapointType.DPT_Percent_U8: // USINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case datapointType == KnxDatapointType.DPT_DecimalFactor: // USINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case datapointType == KnxDatapointType.DPT_Tariff: // USINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case datapointType == KnxDatapointType.DPT_Value_1_Ucount: // USINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case datapointType == KnxDatapointType.DPT_FanStage: // USINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case datapointType == KnxDatapointType.DPT_Percent_V8: // SINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSINT(value), nil
        case datapointType == KnxDatapointType.DPT_Value_1_Count: // SINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSINT(value), nil
        case datapointType == KnxDatapointType.DPT_Status_Mode3: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (StatusA)
            StatusA, _StatusAErr := io.ReadBit()
            if _StatusAErr != nil {
                return nil, errors.New("Error parsing 'StatusA' field " + _StatusAErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(StatusA)

            // Simple Field (StatusB)
            StatusB, _StatusBErr := io.ReadBit()
            if _StatusBErr != nil {
                return nil, errors.New("Error parsing 'StatusB' field " + _StatusBErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(StatusB)

            // Simple Field (StatusC)
            StatusC, _StatusCErr := io.ReadBit()
            if _StatusCErr != nil {
                return nil, errors.New("Error parsing 'StatusC' field " + _StatusCErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(StatusC)

            // Simple Field (StatusD)
            StatusD, _StatusDErr := io.ReadBit()
            if _StatusDErr != nil {
                return nil, errors.New("Error parsing 'StatusD' field " + _StatusDErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(StatusD)

            // Simple Field (StatusE)
            StatusE, _StatusEErr := io.ReadBit()
            if _StatusEErr != nil {
                return nil, errors.New("Error parsing 'StatusE' field " + _StatusEErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(StatusE)

            // Simple Field (Mode)
            Mode, _ModeErr := io.ReadUint8(3)
            if _ModeErr != nil {
                return nil, errors.New("Error parsing 'Mode' field " + _ModeErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Mode)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Value_2_Ucount: // UINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUINT(value), nil
        case datapointType == KnxDatapointType.DPT_TimePeriodMsec: // UINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUINT(value), nil
        case datapointType == KnxDatapointType.DPT_TimePeriod10Msec: // UINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUINT(value), nil
        case datapointType == KnxDatapointType.DPT_TimePeriod100Msec: // UINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUINT(value), nil
        case datapointType == KnxDatapointType.DPT_TimePeriodSec: // UINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUINT(value), nil
        case datapointType == KnxDatapointType.DPT_TimePeriodMin: // UINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUINT(value), nil
        case datapointType == KnxDatapointType.DPT_TimePeriodHrs: // UINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUINT(value), nil
        case datapointType == KnxDatapointType.DPT_PropDataType: // UINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUINT(value), nil
        case datapointType == KnxDatapointType.DPT_Length_mm: // UINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUINT(value), nil
        case datapointType == KnxDatapointType.DPT_UElCurrentmA: // UINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUINT(value), nil
        case datapointType == KnxDatapointType.DPT_Brightness: // UINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUINT(value), nil
        case datapointType == KnxDatapointType.DPT_Absolute_Colour_Temperature: // UINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUINT(value), nil
        case datapointType == KnxDatapointType.DPT_Value_2_Count: // INT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcINT(value), nil
        case datapointType == KnxDatapointType.DPT_DeltaTimeMsec: // INT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcINT(value), nil
        case datapointType == KnxDatapointType.DPT_DeltaTime10Msec: // INT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcINT(value), nil
        case datapointType == KnxDatapointType.DPT_DeltaTime100Msec: // INT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcINT(value), nil
        case datapointType == KnxDatapointType.DPT_DeltaTimeSec: // INT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcINT(value), nil
        case datapointType == KnxDatapointType.DPT_DeltaTimeMin: // INT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcINT(value), nil
        case datapointType == KnxDatapointType.DPT_DeltaTimeHrs: // INT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcINT(value), nil
        case datapointType == KnxDatapointType.DPT_Percent_V16: // INT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcINT(value), nil
        case datapointType == KnxDatapointType.DPT_Rotation_Angle: // INT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcINT(value), nil
        case datapointType == KnxDatapointType.DPT_Length_m: // INT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcINT(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Temp: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Tempd: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Tempa: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Lux: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Wsp: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Pres: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Humidity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_AirQuality: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_AirFlow: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Time1: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Time2: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Volt: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Curr: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_PowerDensity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_KelvinPerPercent: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Power: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Volume_Flow: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Rain_Amount: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Temp_F: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Wsp_kmh: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Absolute_Humidity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Concentration_ygm3: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_TimeOfDay: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Day)
            Day, _DayErr := io.ReadUint8(3)
            if _DayErr != nil {
                return nil, errors.New("Error parsing 'Day' field " + _DayErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Day)

            // Simple Field (Hour)
            Hour, _HourErr := io.ReadUint8(5)
            if _HourErr != nil {
                return nil, errors.New("Error parsing 'Hour' field " + _HourErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Hour)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Minutes)
            Minutes, _MinutesErr := io.ReadUint8(6)
            if _MinutesErr != nil {
                return nil, errors.New("Error parsing 'Minutes' field " + _MinutesErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Minutes)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Seconds)
            Seconds, _SecondsErr := io.ReadUint8(6)
            if _SecondsErr != nil {
                return nil, errors.New("Error parsing 'Seconds' field " + _SecondsErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Seconds)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Date: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(3); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (DayOfMonth)
            DayOfMonth, _DayOfMonthErr := io.ReadUint8(5)
            if _DayOfMonthErr != nil {
                return nil, errors.New("Error parsing 'DayOfMonth' field " + _DayOfMonthErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(DayOfMonth)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Month)
            Month, _MonthErr := io.ReadUint8(4)
            if _MonthErr != nil {
                return nil, errors.New("Error parsing 'Month' field " + _MonthErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Month)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(1); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Year)
            Year, _YearErr := io.ReadUint8(7)
            if _YearErr != nil {
                return nil, errors.New("Error parsing 'Year' field " + _YearErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Year)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Value_4_Ucount: // UDINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUDINT(value), nil
        case datapointType == KnxDatapointType.DPT_LongTimePeriod_Sec: // UDINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUDINT(value), nil
        case datapointType == KnxDatapointType.DPT_LongTimePeriod_Min: // UDINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUDINT(value), nil
        case datapointType == KnxDatapointType.DPT_LongTimePeriod_Hrs: // UDINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUDINT(value), nil
        case datapointType == KnxDatapointType.DPT_VolumeLiquid_Litre: // UDINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUDINT(value), nil
        case datapointType == KnxDatapointType.DPT_Volume_m_3: // UDINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUDINT(value), nil
        case datapointType == KnxDatapointType.DPT_Value_4_Count: // DINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case datapointType == KnxDatapointType.DPT_FlowRate_m3h: // DINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case datapointType == KnxDatapointType.DPT_ActiveEnergy: // DINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case datapointType == KnxDatapointType.DPT_ApparantEnergy: // DINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case datapointType == KnxDatapointType.DPT_ReactiveEnergy: // DINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case datapointType == KnxDatapointType.DPT_ActiveEnergy_kWh: // DINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case datapointType == KnxDatapointType.DPT_ApparantEnergy_kVAh: // DINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case datapointType == KnxDatapointType.DPT_ReactiveEnergy_kVARh: // DINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case datapointType == KnxDatapointType.DPT_ActiveEnergy_MWh: // DINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case datapointType == KnxDatapointType.DPT_LongDeltaTimeSec: // DINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case datapointType == KnxDatapointType.DPT_DeltaVolumeLiquid_Litre: // DINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case datapointType == KnxDatapointType.DPT_DeltaVolume_m_3: // DINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Acceleration: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Acceleration_Angular: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Activation_Energy: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Activity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Mol: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Amplitude: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_AngleRad: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_AngleDeg: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Angular_Momentum: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Angular_Velocity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Area: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Capacitance: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Charge_DensitySurface: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Charge_DensityVolume: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Compressibility: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Conductance: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Electrical_Conductivity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Density: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Electric_Charge: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Electric_Current: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Electric_CurrentDensity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Electric_DipoleMoment: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Electric_Displacement: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Electric_FieldStrength: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Electric_Flux: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Electric_FluxDensity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Electric_Polarization: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Electric_Potential: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Electric_PotentialDifference: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_ElectromagneticMoment: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Electromotive_Force: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Energy: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Force: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Frequency: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Angular_Frequency: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Heat_Capacity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Heat_FlowRate: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Heat_Quantity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Impedance: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Length: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Light_Quantity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Luminance: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Luminous_Flux: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Luminous_Intensity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Magnetic_FieldStrength: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Magnetic_Flux: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Magnetic_FluxDensity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Magnetic_Moment: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Magnetic_Polarization: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Magnetization: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_MagnetomotiveForce: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Mass: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_MassFlux: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Momentum: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Phase_AngleRad: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Phase_AngleDeg: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Power: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Power_Factor: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Pressure: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Reactance: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Resistance: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Resistivity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_SelfInductance: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_SolidAngle: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Sound_Intensity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Speed: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Stress: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Surface_Tension: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Common_Temperature: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Absolute_Temperature: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_TemperatureDifference: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Thermal_Capacity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Thermal_Conductivity: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_ThermoelectricPower: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Time: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Torque: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Volume: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Volume_Flux: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Weight: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Value_Work: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Volume_Flux_Meter: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Volume_Flux_ls: // REAL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case datapointType == KnxDatapointType.DPT_Access_Data: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Value1)
            Value1, _Value1Err := io.ReadUint8(4)
            if _Value1Err != nil {
                return nil, errors.New("Error parsing 'Value1' field " + _Value1Err.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Value1)

            // Simple Field (Value1)
            Value1, _Value1Err := io.ReadUint8(4)
            if _Value1Err != nil {
                return nil, errors.New("Error parsing 'Value1' field " + _Value1Err.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Value1)

            // Simple Field (Value1)
            Value1, _Value1Err := io.ReadUint8(4)
            if _Value1Err != nil {
                return nil, errors.New("Error parsing 'Value1' field " + _Value1Err.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Value1)

            // Simple Field (Value1)
            Value1, _Value1Err := io.ReadUint8(4)
            if _Value1Err != nil {
                return nil, errors.New("Error parsing 'Value1' field " + _Value1Err.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Value1)

            // Simple Field (Value1)
            Value1, _Value1Err := io.ReadUint8(4)
            if _Value1Err != nil {
                return nil, errors.New("Error parsing 'Value1' field " + _Value1Err.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Value1)

            // Simple Field (Value1)
            Value1, _Value1Err := io.ReadUint8(4)
            if _Value1Err != nil {
                return nil, errors.New("Error parsing 'Value1' field " + _Value1Err.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Value1)

            // Simple Field (DetectionError)
            DetectionError, _DetectionErrorErr := io.ReadBit()
            if _DetectionErrorErr != nil {
                return nil, errors.New("Error parsing 'DetectionError' field " + _DetectionErrorErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(DetectionError)

            // Simple Field (Permission)
            Permission, _PermissionErr := io.ReadBit()
            if _PermissionErr != nil {
                return nil, errors.New("Error parsing 'Permission' field " + _PermissionErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Permission)

            // Simple Field (ReadDirection)
            ReadDirection, _ReadDirectionErr := io.ReadBit()
            if _ReadDirectionErr != nil {
                return nil, errors.New("Error parsing 'ReadDirection' field " + _ReadDirectionErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ReadDirection)

            // Simple Field (EncryptionOfAccessInformation)
            EncryptionOfAccessInformation, _EncryptionOfAccessInformationErr := io.ReadBit()
            if _EncryptionOfAccessInformationErr != nil {
                return nil, errors.New("Error parsing 'EncryptionOfAccessInformation' field " + _EncryptionOfAccessInformationErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(EncryptionOfAccessInformation)

            // Simple Field (IndexOfAccessIdentificationCode)
            IndexOfAccessIdentificationCode, _IndexOfAccessIdentificationCodeErr := io.ReadUint8(4)
            if _IndexOfAccessIdentificationCodeErr != nil {
                return nil, errors.New("Error parsing 'IndexOfAccessIdentificationCode' field " + _IndexOfAccessIdentificationCodeErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(IndexOfAccessIdentificationCode)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_String_ASCII: // STRING

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadString(112)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSTRING(value), nil
        case datapointType == KnxDatapointType.DPT_String_8859_1: // STRING

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadString(112)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSTRING(value), nil
        case datapointType == KnxDatapointType.DPT_SceneNumber: // USINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(6)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case datapointType == KnxDatapointType.DPT_SceneControl: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (LearnTheSceneCorrespondingToTheFieldSceneNumber)
            LearnTheSceneCorrespondingToTheFieldSceneNumber, _LearnTheSceneCorrespondingToTheFieldSceneNumberErr := io.ReadBit()
            if _LearnTheSceneCorrespondingToTheFieldSceneNumberErr != nil {
                return nil, errors.New("Error parsing 'LearnTheSceneCorrespondingToTheFieldSceneNumber' field " + _LearnTheSceneCorrespondingToTheFieldSceneNumberErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(LearnTheSceneCorrespondingToTheFieldSceneNumber)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(1); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (SceneNumber)
            SceneNumber, _SceneNumberErr := io.ReadUint8(6)
            if _SceneNumberErr != nil {
                return nil, errors.New("Error parsing 'SceneNumber' field " + _SceneNumberErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(SceneNumber)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_DateTime: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Year)
            Year, _YearErr := io.ReadUint8(8)
            if _YearErr != nil {
                return nil, errors.New("Error parsing 'Year' field " + _YearErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Year)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Month)
            Month, _MonthErr := io.ReadUint8(4)
            if _MonthErr != nil {
                return nil, errors.New("Error parsing 'Month' field " + _MonthErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Month)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(3); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Dayofmonth)
            Dayofmonth, _DayofmonthErr := io.ReadUint8(5)
            if _DayofmonthErr != nil {
                return nil, errors.New("Error parsing 'Dayofmonth' field " + _DayofmonthErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Dayofmonth)

            // Simple Field (Dayofweek)
            Dayofweek, _DayofweekErr := io.ReadUint8(3)
            if _DayofweekErr != nil {
                return nil, errors.New("Error parsing 'Dayofweek' field " + _DayofweekErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Dayofweek)

            // Simple Field (Hourofday)
            Hourofday, _HourofdayErr := io.ReadUint8(5)
            if _HourofdayErr != nil {
                return nil, errors.New("Error parsing 'Hourofday' field " + _HourofdayErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Hourofday)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Minutes)
            Minutes, _MinutesErr := io.ReadUint8(6)
            if _MinutesErr != nil {
                return nil, errors.New("Error parsing 'Minutes' field " + _MinutesErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Minutes)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Seconds)
            Seconds, _SecondsErr := io.ReadUint8(6)
            if _SecondsErr != nil {
                return nil, errors.New("Error parsing 'Seconds' field " + _SecondsErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Seconds)

            // Simple Field (Fault)
            Fault, _FaultErr := io.ReadBit()
            if _FaultErr != nil {
                return nil, errors.New("Error parsing 'Fault' field " + _FaultErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Fault)

            // Simple Field (WorkingDay)
            WorkingDay, _WorkingDayErr := io.ReadBit()
            if _WorkingDayErr != nil {
                return nil, errors.New("Error parsing 'WorkingDay' field " + _WorkingDayErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(WorkingDay)

            // Simple Field (NoWd)
            NoWd, _NoWdErr := io.ReadBit()
            if _NoWdErr != nil {
                return nil, errors.New("Error parsing 'NoWd' field " + _NoWdErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(NoWd)

            // Simple Field (NoYear)
            NoYear, _NoYearErr := io.ReadBit()
            if _NoYearErr != nil {
                return nil, errors.New("Error parsing 'NoYear' field " + _NoYearErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(NoYear)

            // Simple Field (NoDate)
            NoDate, _NoDateErr := io.ReadBit()
            if _NoDateErr != nil {
                return nil, errors.New("Error parsing 'NoDate' field " + _NoDateErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(NoDate)

            // Simple Field (NoDayOfWeek)
            NoDayOfWeek, _NoDayOfWeekErr := io.ReadBit()
            if _NoDayOfWeekErr != nil {
                return nil, errors.New("Error parsing 'NoDayOfWeek' field " + _NoDayOfWeekErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(NoDayOfWeek)

            // Simple Field (NoTime)
            NoTime, _NoTimeErr := io.ReadBit()
            if _NoTimeErr != nil {
                return nil, errors.New("Error parsing 'NoTime' field " + _NoTimeErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(NoTime)

            // Simple Field (StandardSummerTime)
            StandardSummerTime, _StandardSummerTimeErr := io.ReadBit()
            if _StandardSummerTimeErr != nil {
                return nil, errors.New("Error parsing 'StandardSummerTime' field " + _StandardSummerTimeErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(StandardSummerTime)

            // Simple Field (QualityOfClock)
            QualityOfClock, _QualityOfClockErr := io.ReadBit()
            if _QualityOfClockErr != nil {
                return nil, errors.New("Error parsing 'QualityOfClock' field " + _QualityOfClockErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(QualityOfClock)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_SCLOMode: // STRING
        case datapointType == KnxDatapointType.DPT_BuildingMode: // STRING
        case datapointType == KnxDatapointType.DPT_OccMode: // STRING
        case datapointType == KnxDatapointType.DPT_Priority: // STRING
        case datapointType == KnxDatapointType.DPT_LightApplicationMode: // STRING
        case datapointType == KnxDatapointType.DPT_ApplicationArea: // STRING
        case datapointType == KnxDatapointType.DPT_AlarmClassType: // STRING
        case datapointType == KnxDatapointType.DPT_PSUMode: // STRING
        case datapointType == KnxDatapointType.DPT_ErrorClass_System: // STRING
        case datapointType == KnxDatapointType.DPT_ErrorClass_HVAC: // STRING
        case datapointType == KnxDatapointType.DPT_Time_Delay: // STRING
        case datapointType == KnxDatapointType.DPT_Beaufort_Wind_Force_Scale: // STRING
        case datapointType == KnxDatapointType.DPT_SensorSelect: // STRING
        case datapointType == KnxDatapointType.DPT_ActuatorConnectType: // STRING
        case datapointType == KnxDatapointType.DPT_Cloud_Cover: // STRING
        case datapointType == KnxDatapointType.DPT_PowerReturnMode: // STRING
        case datapointType == KnxDatapointType.DPT_FuelType: // STRING
        case datapointType == KnxDatapointType.DPT_BurnerType: // STRING
        case datapointType == KnxDatapointType.DPT_HVACMode: // STRING
        case datapointType == KnxDatapointType.DPT_DHWMode: // STRING
        case datapointType == KnxDatapointType.DPT_LoadPriority: // STRING
        case datapointType == KnxDatapointType.DPT_HVACContrMode: // STRING
        case datapointType == KnxDatapointType.DPT_HVACEmergMode: // STRING
        case datapointType == KnxDatapointType.DPT_ChangeoverMode: // STRING
        case datapointType == KnxDatapointType.DPT_ValveMode: // STRING
        case datapointType == KnxDatapointType.DPT_DamperMode: // STRING
        case datapointType == KnxDatapointType.DPT_HeaterMode: // STRING
        case datapointType == KnxDatapointType.DPT_FanMode: // STRING
        case datapointType == KnxDatapointType.DPT_MasterSlaveMode: // STRING
        case datapointType == KnxDatapointType.DPT_StatusRoomSetp: // STRING
        case datapointType == KnxDatapointType.DPT_Metering_DeviceType: // STRING
        case datapointType == KnxDatapointType.DPT_HumDehumMode: // STRING
        case datapointType == KnxDatapointType.DPT_EnableHCStage: // STRING
        case datapointType == KnxDatapointType.DPT_ADAType: // STRING
        case datapointType == KnxDatapointType.DPT_BackupMode: // STRING
        case datapointType == KnxDatapointType.DPT_StartSynchronization: // STRING
        case datapointType == KnxDatapointType.DPT_Behaviour_Lock_Unlock: // STRING
        case datapointType == KnxDatapointType.DPT_Behaviour_Bus_Power_Up_Down: // STRING
        case datapointType == KnxDatapointType.DPT_DALI_Fade_Time: // STRING
        case datapointType == KnxDatapointType.DPT_BlinkingMode: // STRING
        case datapointType == KnxDatapointType.DPT_LightControlMode: // STRING
        case datapointType == KnxDatapointType.DPT_SwitchPBModel: // STRING
        case datapointType == KnxDatapointType.DPT_PBAction: // STRING
        case datapointType == KnxDatapointType.DPT_DimmPBModel: // STRING
        case datapointType == KnxDatapointType.DPT_SwitchOnMode: // STRING
        case datapointType == KnxDatapointType.DPT_LoadTypeSet: // STRING
        case datapointType == KnxDatapointType.DPT_LoadTypeDetected: // STRING
        case datapointType == KnxDatapointType.DPT_Converter_Test_Control: // STRING
        case datapointType == KnxDatapointType.DPT_SABExcept_Behaviour: // STRING
        case datapointType == KnxDatapointType.DPT_SABBehaviour_Lock_Unlock: // STRING
        case datapointType == KnxDatapointType.DPT_SSSBMode: // STRING
        case datapointType == KnxDatapointType.DPT_BlindsControlMode: // STRING
        case datapointType == KnxDatapointType.DPT_CommMode: // STRING
        case datapointType == KnxDatapointType.DPT_AddInfoTypes: // STRING
        case datapointType == KnxDatapointType.DPT_RF_ModeSelect: // STRING
        case datapointType == KnxDatapointType.DPT_RF_FilterSelect: // STRING
        case datapointType == KnxDatapointType.DPT_StatusGen: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(3); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (AlarmStatusOfCorrespondingDatapointIsNotAcknowledged)
            AlarmStatusOfCorrespondingDatapointIsNotAcknowledged, _AlarmStatusOfCorrespondingDatapointIsNotAcknowledgedErr := io.ReadBit()
            if _AlarmStatusOfCorrespondingDatapointIsNotAcknowledgedErr != nil {
                return nil, errors.New("Error parsing 'AlarmStatusOfCorrespondingDatapointIsNotAcknowledged' field " + _AlarmStatusOfCorrespondingDatapointIsNotAcknowledgedErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(AlarmStatusOfCorrespondingDatapointIsNotAcknowledged)

            // Simple Field (CorrespondingDatapointIsInAlarm)
            CorrespondingDatapointIsInAlarm, _CorrespondingDatapointIsInAlarmErr := io.ReadBit()
            if _CorrespondingDatapointIsInAlarmErr != nil {
                return nil, errors.New("Error parsing 'CorrespondingDatapointIsInAlarm' field " + _CorrespondingDatapointIsInAlarmErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(CorrespondingDatapointIsInAlarm)

            // Simple Field (CorrespondingDatapointMainValueIsOverridden)
            CorrespondingDatapointMainValueIsOverridden, _CorrespondingDatapointMainValueIsOverriddenErr := io.ReadBit()
            if _CorrespondingDatapointMainValueIsOverriddenErr != nil {
                return nil, errors.New("Error parsing 'CorrespondingDatapointMainValueIsOverridden' field " + _CorrespondingDatapointMainValueIsOverriddenErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(CorrespondingDatapointMainValueIsOverridden)

            // Simple Field (CorrespondingDatapointMainValueIsCorruptedDueToFailure)
            CorrespondingDatapointMainValueIsCorruptedDueToFailure, _CorrespondingDatapointMainValueIsCorruptedDueToFailureErr := io.ReadBit()
            if _CorrespondingDatapointMainValueIsCorruptedDueToFailureErr != nil {
                return nil, errors.New("Error parsing 'CorrespondingDatapointMainValueIsCorruptedDueToFailure' field " + _CorrespondingDatapointMainValueIsCorruptedDueToFailureErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(CorrespondingDatapointMainValueIsCorruptedDueToFailure)

            // Simple Field (CorrespondingDatapointValueIsOutOfService)
            CorrespondingDatapointValueIsOutOfService, _CorrespondingDatapointValueIsOutOfServiceErr := io.ReadBit()
            if _CorrespondingDatapointValueIsOutOfServiceErr != nil {
                return nil, errors.New("Error parsing 'CorrespondingDatapointValueIsOutOfService' field " + _CorrespondingDatapointValueIsOutOfServiceErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(CorrespondingDatapointValueIsOutOfService)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Device_Control: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (VerifyModeIsOn)
            VerifyModeIsOn, _VerifyModeIsOnErr := io.ReadBit()
            if _VerifyModeIsOnErr != nil {
                return nil, errors.New("Error parsing 'VerifyModeIsOn' field " + _VerifyModeIsOnErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(VerifyModeIsOn)

            // Simple Field (ADatagramWithTheOwnIndividualAddressAsSourceAddressHasBeenReceived)
            ADatagramWithTheOwnIndividualAddressAsSourceAddressHasBeenReceived, _ADatagramWithTheOwnIndividualAddressAsSourceAddressHasBeenReceivedErr := io.ReadBit()
            if _ADatagramWithTheOwnIndividualAddressAsSourceAddressHasBeenReceivedErr != nil {
                return nil, errors.New("Error parsing 'ADatagramWithTheOwnIndividualAddressAsSourceAddressHasBeenReceived' field " + _ADatagramWithTheOwnIndividualAddressAsSourceAddressHasBeenReceivedErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ADatagramWithTheOwnIndividualAddressAsSourceAddressHasBeenReceived)

            // Simple Field (TheUserApplicationIsStopped)
            TheUserApplicationIsStopped, _TheUserApplicationIsStoppedErr := io.ReadBit()
            if _TheUserApplicationIsStoppedErr != nil {
                return nil, errors.New("Error parsing 'TheUserApplicationIsStopped' field " + _TheUserApplicationIsStoppedErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(TheUserApplicationIsStopped)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_ForceSign: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Roomhmax)
            Roomhmax, _RoomhmaxErr := io.ReadBit()
            if _RoomhmaxErr != nil {
                return nil, errors.New("Error parsing 'Roomhmax' field " + _RoomhmaxErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Roomhmax)

            // Simple Field (Roomhconf)
            Roomhconf, _RoomhconfErr := io.ReadBit()
            if _RoomhconfErr != nil {
                return nil, errors.New("Error parsing 'Roomhconf' field " + _RoomhconfErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Roomhconf)

            // Simple Field (Dhwlegio)
            Dhwlegio, _DhwlegioErr := io.ReadBit()
            if _DhwlegioErr != nil {
                return nil, errors.New("Error parsing 'Dhwlegio' field " + _DhwlegioErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Dhwlegio)

            // Simple Field (Dhwnorm)
            Dhwnorm, _DhwnormErr := io.ReadBit()
            if _DhwnormErr != nil {
                return nil, errors.New("Error parsing 'Dhwnorm' field " + _DhwnormErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Dhwnorm)

            // Simple Field (Overrun)
            Overrun, _OverrunErr := io.ReadBit()
            if _OverrunErr != nil {
                return nil, errors.New("Error parsing 'Overrun' field " + _OverrunErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Overrun)

            // Simple Field (Oversupply)
            Oversupply, _OversupplyErr := io.ReadBit()
            if _OversupplyErr != nil {
                return nil, errors.New("Error parsing 'Oversupply' field " + _OversupplyErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Oversupply)

            // Simple Field (Protection)
            Protection, _ProtectionErr := io.ReadBit()
            if _ProtectionErr != nil {
                return nil, errors.New("Error parsing 'Protection' field " + _ProtectionErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Protection)

            // Simple Field (Forcerequest)
            Forcerequest, _ForcerequestErr := io.ReadBit()
            if _ForcerequestErr != nil {
                return nil, errors.New("Error parsing 'Forcerequest' field " + _ForcerequestErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Forcerequest)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_ForceSignCool: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_StatusRHC: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Summermode)
            Summermode, _SummermodeErr := io.ReadBit()
            if _SummermodeErr != nil {
                return nil, errors.New("Error parsing 'Summermode' field " + _SummermodeErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Summermode)

            // Simple Field (Statusstopoptim)
            Statusstopoptim, _StatusstopoptimErr := io.ReadBit()
            if _StatusstopoptimErr != nil {
                return nil, errors.New("Error parsing 'Statusstopoptim' field " + _StatusstopoptimErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Statusstopoptim)

            // Simple Field (Statusstartoptim)
            Statusstartoptim, _StatusstartoptimErr := io.ReadBit()
            if _StatusstartoptimErr != nil {
                return nil, errors.New("Error parsing 'Statusstartoptim' field " + _StatusstartoptimErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Statusstartoptim)

            // Simple Field (Statusmorningboost)
            Statusmorningboost, _StatusmorningboostErr := io.ReadBit()
            if _StatusmorningboostErr != nil {
                return nil, errors.New("Error parsing 'Statusmorningboost' field " + _StatusmorningboostErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Statusmorningboost)

            // Simple Field (Tempreturnlimit)
            Tempreturnlimit, _TempreturnlimitErr := io.ReadBit()
            if _TempreturnlimitErr != nil {
                return nil, errors.New("Error parsing 'Tempreturnlimit' field " + _TempreturnlimitErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Tempreturnlimit)

            // Simple Field (Tempflowlimit)
            Tempflowlimit, _TempflowlimitErr := io.ReadBit()
            if _TempflowlimitErr != nil {
                return nil, errors.New("Error parsing 'Tempflowlimit' field " + _TempflowlimitErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Tempflowlimit)

            // Simple Field (Satuseco)
            Satuseco, _SatusecoErr := io.ReadBit()
            if _SatusecoErr != nil {
                return nil, errors.New("Error parsing 'Satuseco' field " + _SatusecoErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Satuseco)

            // Simple Field (Fault)
            Fault, _FaultErr := io.ReadBit()
            if _FaultErr != nil {
                return nil, errors.New("Error parsing 'Fault' field " + _FaultErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Fault)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_StatusSDHWC: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Solarloadsufficient)
            Solarloadsufficient, _SolarloadsufficientErr := io.ReadBit()
            if _SolarloadsufficientErr != nil {
                return nil, errors.New("Error parsing 'Solarloadsufficient' field " + _SolarloadsufficientErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Solarloadsufficient)

            // Simple Field (Sdhwloadactive)
            Sdhwloadactive, _SdhwloadactiveErr := io.ReadBit()
            if _SdhwloadactiveErr != nil {
                return nil, errors.New("Error parsing 'Sdhwloadactive' field " + _SdhwloadactiveErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Sdhwloadactive)

            // Simple Field (Fault)
            Fault, _FaultErr := io.ReadBit()
            if _FaultErr != nil {
                return nil, errors.New("Error parsing 'Fault' field " + _FaultErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Fault)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_FuelTypeSet: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Solidstate)
            Solidstate, _SolidstateErr := io.ReadBit()
            if _SolidstateErr != nil {
                return nil, errors.New("Error parsing 'Solidstate' field " + _SolidstateErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Solidstate)

            // Simple Field (Gas)
            Gas, _GasErr := io.ReadBit()
            if _GasErr != nil {
                return nil, errors.New("Error parsing 'Gas' field " + _GasErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Gas)

            // Simple Field (Oil)
            Oil, _OilErr := io.ReadBit()
            if _OilErr != nil {
                return nil, errors.New("Error parsing 'Oil' field " + _OilErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Oil)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_StatusRCC: // BOOL

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case datapointType == KnxDatapointType.DPT_StatusAHU: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Cool)
            Cool, _CoolErr := io.ReadBit()
            if _CoolErr != nil {
                return nil, errors.New("Error parsing 'Cool' field " + _CoolErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Cool)

            // Simple Field (Heat)
            Heat, _HeatErr := io.ReadBit()
            if _HeatErr != nil {
                return nil, errors.New("Error parsing 'Heat' field " + _HeatErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Heat)

            // Simple Field (Fanactive)
            Fanactive, _FanactiveErr := io.ReadBit()
            if _FanactiveErr != nil {
                return nil, errors.New("Error parsing 'Fanactive' field " + _FanactiveErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Fanactive)

            // Simple Field (Fault)
            Fault, _FaultErr := io.ReadBit()
            if _FaultErr != nil {
                return nil, errors.New("Error parsing 'Fault' field " + _FaultErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Fault)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_CombinedStatus_RTSM: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(3); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (StatusOfHvacModeUser)
            StatusOfHvacModeUser, _StatusOfHvacModeUserErr := io.ReadBit()
            if _StatusOfHvacModeUserErr != nil {
                return nil, errors.New("Error parsing 'StatusOfHvacModeUser' field " + _StatusOfHvacModeUserErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(StatusOfHvacModeUser)

            // Simple Field (StatusOfComfortProlongationUser)
            StatusOfComfortProlongationUser, _StatusOfComfortProlongationUserErr := io.ReadBit()
            if _StatusOfComfortProlongationUserErr != nil {
                return nil, errors.New("Error parsing 'StatusOfComfortProlongationUser' field " + _StatusOfComfortProlongationUserErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(StatusOfComfortProlongationUser)

            // Simple Field (EffectiveValueOfTheComfortPushButton)
            EffectiveValueOfTheComfortPushButton, _EffectiveValueOfTheComfortPushButtonErr := io.ReadBit()
            if _EffectiveValueOfTheComfortPushButtonErr != nil {
                return nil, errors.New("Error parsing 'EffectiveValueOfTheComfortPushButton' field " + _EffectiveValueOfTheComfortPushButtonErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(EffectiveValueOfTheComfortPushButton)

            // Simple Field (EffectiveValueOfThePresenceStatus)
            EffectiveValueOfThePresenceStatus, _EffectiveValueOfThePresenceStatusErr := io.ReadBit()
            if _EffectiveValueOfThePresenceStatusErr != nil {
                return nil, errors.New("Error parsing 'EffectiveValueOfThePresenceStatus' field " + _EffectiveValueOfThePresenceStatusErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(EffectiveValueOfThePresenceStatus)

            // Simple Field (EffectiveValueOfTheWindowStatus)
            EffectiveValueOfTheWindowStatus, _EffectiveValueOfTheWindowStatusErr := io.ReadBit()
            if _EffectiveValueOfTheWindowStatusErr != nil {
                return nil, errors.New("Error parsing 'EffectiveValueOfTheWindowStatus' field " + _EffectiveValueOfTheWindowStatusErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(EffectiveValueOfTheWindowStatus)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_LightActuatorErrorInfo: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(1); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Overheat)
            Overheat, _OverheatErr := io.ReadBit()
            if _OverheatErr != nil {
                return nil, errors.New("Error parsing 'Overheat' field " + _OverheatErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Overheat)

            // Simple Field (Lampfailure)
            Lampfailure, _LampfailureErr := io.ReadBit()
            if _LampfailureErr != nil {
                return nil, errors.New("Error parsing 'Lampfailure' field " + _LampfailureErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Lampfailure)

            // Simple Field (Defectiveload)
            Defectiveload, _DefectiveloadErr := io.ReadBit()
            if _DefectiveloadErr != nil {
                return nil, errors.New("Error parsing 'Defectiveload' field " + _DefectiveloadErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Defectiveload)

            // Simple Field (Underload)
            Underload, _UnderloadErr := io.ReadBit()
            if _UnderloadErr != nil {
                return nil, errors.New("Error parsing 'Underload' field " + _UnderloadErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Underload)

            // Simple Field (Overcurrent)
            Overcurrent, _OvercurrentErr := io.ReadBit()
            if _OvercurrentErr != nil {
                return nil, errors.New("Error parsing 'Overcurrent' field " + _OvercurrentErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Overcurrent)

            // Simple Field (Undervoltage)
            Undervoltage, _UndervoltageErr := io.ReadBit()
            if _UndervoltageErr != nil {
                return nil, errors.New("Error parsing 'Undervoltage' field " + _UndervoltageErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Undervoltage)

            // Simple Field (Loaddetectionerror)
            Loaddetectionerror, _LoaddetectionerrorErr := io.ReadBit()
            if _LoaddetectionerrorErr != nil {
                return nil, errors.New("Error parsing 'Loaddetectionerror' field " + _LoaddetectionerrorErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Loaddetectionerror)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_RF_ModeInfo: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (BibatSlave)
            BibatSlave, _BibatSlaveErr := io.ReadBit()
            if _BibatSlaveErr != nil {
                return nil, errors.New("Error parsing 'BibatSlave' field " + _BibatSlaveErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(BibatSlave)

            // Simple Field (BibatMaster)
            BibatMaster, _BibatMasterErr := io.ReadBit()
            if _BibatMasterErr != nil {
                return nil, errors.New("Error parsing 'BibatMaster' field " + _BibatMasterErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(BibatMaster)

            // Simple Field (Asynchronous)
            Asynchronous, _AsynchronousErr := io.ReadBit()
            if _AsynchronousErr != nil {
                return nil, errors.New("Error parsing 'Asynchronous' field " + _AsynchronousErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Asynchronous)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_RF_FilterInfo: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Doa)
            Doa, _DoaErr := io.ReadBit()
            if _DoaErr != nil {
                return nil, errors.New("Error parsing 'Doa' field " + _DoaErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Doa)

            // Simple Field (KnxSn)
            KnxSn, _KnxSnErr := io.ReadBit()
            if _KnxSnErr != nil {
                return nil, errors.New("Error parsing 'KnxSn' field " + _KnxSnErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(KnxSn)

            // Simple Field (DoaAndKnxSn)
            DoaAndKnxSn, _DoaAndKnxSnErr := io.ReadBit()
            if _DoaAndKnxSnErr != nil {
                return nil, errors.New("Error parsing 'DoaAndKnxSn' field " + _DoaAndKnxSnErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(DoaAndKnxSn)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Channel_Activation_8: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel1)
            ActivationStateOfChannel1, _ActivationStateOfChannel1Err := io.ReadBit()
            if _ActivationStateOfChannel1Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel1' field " + _ActivationStateOfChannel1Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel1)

            // Simple Field (ActivationStateOfChannel2)
            ActivationStateOfChannel2, _ActivationStateOfChannel2Err := io.ReadBit()
            if _ActivationStateOfChannel2Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel2' field " + _ActivationStateOfChannel2Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel2)

            // Simple Field (ActivationStateOfChannel3)
            ActivationStateOfChannel3, _ActivationStateOfChannel3Err := io.ReadBit()
            if _ActivationStateOfChannel3Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel3' field " + _ActivationStateOfChannel3Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel3)

            // Simple Field (ActivationStateOfChannel4)
            ActivationStateOfChannel4, _ActivationStateOfChannel4Err := io.ReadBit()
            if _ActivationStateOfChannel4Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel4' field " + _ActivationStateOfChannel4Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel4)

            // Simple Field (ActivationStateOfChannel5)
            ActivationStateOfChannel5, _ActivationStateOfChannel5Err := io.ReadBit()
            if _ActivationStateOfChannel5Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel5' field " + _ActivationStateOfChannel5Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel5)

            // Simple Field (ActivationStateOfChannel6)
            ActivationStateOfChannel6, _ActivationStateOfChannel6Err := io.ReadBit()
            if _ActivationStateOfChannel6Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel6' field " + _ActivationStateOfChannel6Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel6)

            // Simple Field (ActivationStateOfChannel7)
            ActivationStateOfChannel7, _ActivationStateOfChannel7Err := io.ReadBit()
            if _ActivationStateOfChannel7Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel7' field " + _ActivationStateOfChannel7Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel7)

            // Simple Field (ActivationStateOfChannel8)
            ActivationStateOfChannel8, _ActivationStateOfChannel8Err := io.ReadBit()
            if _ActivationStateOfChannel8Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel8' field " + _ActivationStateOfChannel8Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel8)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_StatusDHWC: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Tempoptimshiftactive)
            Tempoptimshiftactive, _TempoptimshiftactiveErr := io.ReadBit()
            if _TempoptimshiftactiveErr != nil {
                return nil, errors.New("Error parsing 'Tempoptimshiftactive' field " + _TempoptimshiftactiveErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Tempoptimshiftactive)

            // Simple Field (Solarenergysupport)
            Solarenergysupport, _SolarenergysupportErr := io.ReadBit()
            if _SolarenergysupportErr != nil {
                return nil, errors.New("Error parsing 'Solarenergysupport' field " + _SolarenergysupportErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Solarenergysupport)

            // Simple Field (Solarenergyonly)
            Solarenergyonly, _SolarenergyonlyErr := io.ReadBit()
            if _SolarenergyonlyErr != nil {
                return nil, errors.New("Error parsing 'Solarenergyonly' field " + _SolarenergyonlyErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Solarenergyonly)

            // Simple Field (Otherenergysourceactive)
            Otherenergysourceactive, _OtherenergysourceactiveErr := io.ReadBit()
            if _OtherenergysourceactiveErr != nil {
                return nil, errors.New("Error parsing 'Otherenergysourceactive' field " + _OtherenergysourceactiveErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Otherenergysourceactive)

            // Simple Field (Dhwpushactive)
            Dhwpushactive, _DhwpushactiveErr := io.ReadBit()
            if _DhwpushactiveErr != nil {
                return nil, errors.New("Error parsing 'Dhwpushactive' field " + _DhwpushactiveErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Dhwpushactive)

            // Simple Field (Legioprotactive)
            Legioprotactive, _LegioprotactiveErr := io.ReadBit()
            if _LegioprotactiveErr != nil {
                return nil, errors.New("Error parsing 'Legioprotactive' field " + _LegioprotactiveErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Legioprotactive)

            // Simple Field (Dhwloadactive)
            Dhwloadactive, _DhwloadactiveErr := io.ReadBit()
            if _DhwloadactiveErr != nil {
                return nil, errors.New("Error parsing 'Dhwloadactive' field " + _DhwloadactiveErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Dhwloadactive)

            // Simple Field (Fault)
            Fault, _FaultErr := io.ReadBit()
            if _FaultErr != nil {
                return nil, errors.New("Error parsing 'Fault' field " + _FaultErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Fault)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_StatusRHCC: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(1); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Overheatalarm)
            Overheatalarm, _OverheatalarmErr := io.ReadBit()
            if _OverheatalarmErr != nil {
                return nil, errors.New("Error parsing 'Overheatalarm' field " + _OverheatalarmErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Overheatalarm)

            // Simple Field (Frostalarm)
            Frostalarm, _FrostalarmErr := io.ReadBit()
            if _FrostalarmErr != nil {
                return nil, errors.New("Error parsing 'Frostalarm' field " + _FrostalarmErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Frostalarm)

            // Simple Field (Dewpointstatus)
            Dewpointstatus, _DewpointstatusErr := io.ReadBit()
            if _DewpointstatusErr != nil {
                return nil, errors.New("Error parsing 'Dewpointstatus' field " + _DewpointstatusErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Dewpointstatus)

            // Simple Field (Coolingdisabled)
            Coolingdisabled, _CoolingdisabledErr := io.ReadBit()
            if _CoolingdisabledErr != nil {
                return nil, errors.New("Error parsing 'Coolingdisabled' field " + _CoolingdisabledErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Coolingdisabled)

            // Simple Field (Statusprecool)
            Statusprecool, _StatusprecoolErr := io.ReadBit()
            if _StatusprecoolErr != nil {
                return nil, errors.New("Error parsing 'Statusprecool' field " + _StatusprecoolErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Statusprecool)

            // Simple Field (Statusecoc)
            Statusecoc, _StatusecocErr := io.ReadBit()
            if _StatusecocErr != nil {
                return nil, errors.New("Error parsing 'Statusecoc' field " + _StatusecocErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Statusecoc)

            // Simple Field (Heatcoolmode)
            Heatcoolmode, _HeatcoolmodeErr := io.ReadBit()
            if _HeatcoolmodeErr != nil {
                return nil, errors.New("Error parsing 'Heatcoolmode' field " + _HeatcoolmodeErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Heatcoolmode)

            // Simple Field (Heatingdiabled)
            Heatingdiabled, _HeatingdiabledErr := io.ReadBit()
            if _HeatingdiabledErr != nil {
                return nil, errors.New("Error parsing 'Heatingdiabled' field " + _HeatingdiabledErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Heatingdiabled)

            // Simple Field (Statusstopoptim)
            Statusstopoptim, _StatusstopoptimErr := io.ReadBit()
            if _StatusstopoptimErr != nil {
                return nil, errors.New("Error parsing 'Statusstopoptim' field " + _StatusstopoptimErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Statusstopoptim)

            // Simple Field (Statusstartoptim)
            Statusstartoptim, _StatusstartoptimErr := io.ReadBit()
            if _StatusstartoptimErr != nil {
                return nil, errors.New("Error parsing 'Statusstartoptim' field " + _StatusstartoptimErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Statusstartoptim)

            // Simple Field (Statusmorningboosth)
            Statusmorningboosth, _StatusmorningboosthErr := io.ReadBit()
            if _StatusmorningboosthErr != nil {
                return nil, errors.New("Error parsing 'Statusmorningboosth' field " + _StatusmorningboosthErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Statusmorningboosth)

            // Simple Field (Tempflowreturnlimit)
            Tempflowreturnlimit, _TempflowreturnlimitErr := io.ReadBit()
            if _TempflowreturnlimitErr != nil {
                return nil, errors.New("Error parsing 'Tempflowreturnlimit' field " + _TempflowreturnlimitErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Tempflowreturnlimit)

            // Simple Field (Tempflowlimit)
            Tempflowlimit, _TempflowlimitErr := io.ReadBit()
            if _TempflowlimitErr != nil {
                return nil, errors.New("Error parsing 'Tempflowlimit' field " + _TempflowlimitErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Tempflowlimit)

            // Simple Field (Statusecoh)
            Statusecoh, _StatusecohErr := io.ReadBit()
            if _StatusecohErr != nil {
                return nil, errors.New("Error parsing 'Statusecoh' field " + _StatusecohErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Statusecoh)

            // Simple Field (Fault)
            Fault, _FaultErr := io.ReadBit()
            if _FaultErr != nil {
                return nil, errors.New("Error parsing 'Fault' field " + _FaultErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Fault)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_CombinedStatus_HVA: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (CalibrationMode)
            CalibrationMode, _CalibrationModeErr := io.ReadBit()
            if _CalibrationModeErr != nil {
                return nil, errors.New("Error parsing 'CalibrationMode' field " + _CalibrationModeErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(CalibrationMode)

            // Simple Field (LockedPosition)
            LockedPosition, _LockedPositionErr := io.ReadBit()
            if _LockedPositionErr != nil {
                return nil, errors.New("Error parsing 'LockedPosition' field " + _LockedPositionErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(LockedPosition)

            // Simple Field (ForcedPosition)
            ForcedPosition, _ForcedPositionErr := io.ReadBit()
            if _ForcedPositionErr != nil {
                return nil, errors.New("Error parsing 'ForcedPosition' field " + _ForcedPositionErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ForcedPosition)

            // Simple Field (ManuaOperationOverridden)
            ManuaOperationOverridden, _ManuaOperationOverriddenErr := io.ReadBit()
            if _ManuaOperationOverriddenErr != nil {
                return nil, errors.New("Error parsing 'ManuaOperationOverridden' field " + _ManuaOperationOverriddenErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ManuaOperationOverridden)

            // Simple Field (ServiceMode)
            ServiceMode, _ServiceModeErr := io.ReadBit()
            if _ServiceModeErr != nil {
                return nil, errors.New("Error parsing 'ServiceMode' field " + _ServiceModeErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ServiceMode)

            // Simple Field (ValveKick)
            ValveKick, _ValveKickErr := io.ReadBit()
            if _ValveKickErr != nil {
                return nil, errors.New("Error parsing 'ValveKick' field " + _ValveKickErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ValveKick)

            // Simple Field (Overload)
            Overload, _OverloadErr := io.ReadBit()
            if _OverloadErr != nil {
                return nil, errors.New("Error parsing 'Overload' field " + _OverloadErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Overload)

            // Simple Field (ShortCircuit)
            ShortCircuit, _ShortCircuitErr := io.ReadBit()
            if _ShortCircuitErr != nil {
                return nil, errors.New("Error parsing 'ShortCircuit' field " + _ShortCircuitErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ShortCircuit)

            // Simple Field (CurrentValvePosition)
            CurrentValvePosition, _CurrentValvePositionErr := io.ReadBit()
            if _CurrentValvePositionErr != nil {
                return nil, errors.New("Error parsing 'CurrentValvePosition' field " + _CurrentValvePositionErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(CurrentValvePosition)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_CombinedStatus_RTC: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (CoolingModeEnabled)
            CoolingModeEnabled, _CoolingModeEnabledErr := io.ReadBit()
            if _CoolingModeEnabledErr != nil {
                return nil, errors.New("Error parsing 'CoolingModeEnabled' field " + _CoolingModeEnabledErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(CoolingModeEnabled)

            // Simple Field (HeatingModeEnabled)
            HeatingModeEnabled, _HeatingModeEnabledErr := io.ReadBit()
            if _HeatingModeEnabledErr != nil {
                return nil, errors.New("Error parsing 'HeatingModeEnabled' field " + _HeatingModeEnabledErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(HeatingModeEnabled)

            // Simple Field (AdditionalHeatingCoolingStage2Stage)
            AdditionalHeatingCoolingStage2Stage, _AdditionalHeatingCoolingStage2StageErr := io.ReadBit()
            if _AdditionalHeatingCoolingStage2StageErr != nil {
                return nil, errors.New("Error parsing 'AdditionalHeatingCoolingStage2Stage' field " + _AdditionalHeatingCoolingStage2StageErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(AdditionalHeatingCoolingStage2Stage)

            // Simple Field (ControllerInactive)
            ControllerInactive, _ControllerInactiveErr := io.ReadBit()
            if _ControllerInactiveErr != nil {
                return nil, errors.New("Error parsing 'ControllerInactive' field " + _ControllerInactiveErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ControllerInactive)

            // Simple Field (OverheatAlarm)
            OverheatAlarm, _OverheatAlarmErr := io.ReadBit()
            if _OverheatAlarmErr != nil {
                return nil, errors.New("Error parsing 'OverheatAlarm' field " + _OverheatAlarmErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(OverheatAlarm)

            // Simple Field (FrostAlarm)
            FrostAlarm, _FrostAlarmErr := io.ReadBit()
            if _FrostAlarmErr != nil {
                return nil, errors.New("Error parsing 'FrostAlarm' field " + _FrostAlarmErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(FrostAlarm)

            // Simple Field (DewPointStatus)
            DewPointStatus, _DewPointStatusErr := io.ReadBit()
            if _DewPointStatusErr != nil {
                return nil, errors.New("Error parsing 'DewPointStatus' field " + _DewPointStatusErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(DewPointStatus)

            // Simple Field (ActiveMode)
            ActiveMode, _ActiveModeErr := io.ReadBit()
            if _ActiveModeErr != nil {
                return nil, errors.New("Error parsing 'ActiveMode' field " + _ActiveModeErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActiveMode)

            // Simple Field (GeneralFailureInformation)
            GeneralFailureInformation, _GeneralFailureInformationErr := io.ReadBit()
            if _GeneralFailureInformationErr != nil {
                return nil, errors.New("Error parsing 'GeneralFailureInformation' field " + _GeneralFailureInformationErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(GeneralFailureInformation)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Media: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint16(10); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (KnxIp)
            KnxIp, _KnxIpErr := io.ReadBit()
            if _KnxIpErr != nil {
                return nil, errors.New("Error parsing 'KnxIp' field " + _KnxIpErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(KnxIp)

            // Simple Field (Rf)
            Rf, _RfErr := io.ReadBit()
            if _RfErr != nil {
                return nil, errors.New("Error parsing 'Rf' field " + _RfErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Rf)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(1); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Pl110)
            Pl110, _Pl110Err := io.ReadBit()
            if _Pl110Err != nil {
                return nil, errors.New("Error parsing 'Pl110' field " + _Pl110Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Pl110)

            // Simple Field (Tp1)
            Tp1, _Tp1Err := io.ReadBit()
            if _Tp1Err != nil {
                return nil, errors.New("Error parsing 'Tp1' field " + _Tp1Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Tp1)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(1); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Channel_Activation_16: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel1)
            ActivationStateOfChannel1, _ActivationStateOfChannel1Err := io.ReadBit()
            if _ActivationStateOfChannel1Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel1' field " + _ActivationStateOfChannel1Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel1)

            // Simple Field (ActivationStateOfChannel2)
            ActivationStateOfChannel2, _ActivationStateOfChannel2Err := io.ReadBit()
            if _ActivationStateOfChannel2Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel2' field " + _ActivationStateOfChannel2Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel2)

            // Simple Field (ActivationStateOfChannel3)
            ActivationStateOfChannel3, _ActivationStateOfChannel3Err := io.ReadBit()
            if _ActivationStateOfChannel3Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel3' field " + _ActivationStateOfChannel3Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel3)

            // Simple Field (ActivationStateOfChannel4)
            ActivationStateOfChannel4, _ActivationStateOfChannel4Err := io.ReadBit()
            if _ActivationStateOfChannel4Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel4' field " + _ActivationStateOfChannel4Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel4)

            // Simple Field (ActivationStateOfChannel5)
            ActivationStateOfChannel5, _ActivationStateOfChannel5Err := io.ReadBit()
            if _ActivationStateOfChannel5Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel5' field " + _ActivationStateOfChannel5Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel5)

            // Simple Field (ActivationStateOfChannel6)
            ActivationStateOfChannel6, _ActivationStateOfChannel6Err := io.ReadBit()
            if _ActivationStateOfChannel6Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel6' field " + _ActivationStateOfChannel6Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel6)

            // Simple Field (ActivationStateOfChannel7)
            ActivationStateOfChannel7, _ActivationStateOfChannel7Err := io.ReadBit()
            if _ActivationStateOfChannel7Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel7' field " + _ActivationStateOfChannel7Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel7)

            // Simple Field (ActivationStateOfChannel8)
            ActivationStateOfChannel8, _ActivationStateOfChannel8Err := io.ReadBit()
            if _ActivationStateOfChannel8Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel8' field " + _ActivationStateOfChannel8Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel8)

            // Simple Field (ActivationStateOfChannel9)
            ActivationStateOfChannel9, _ActivationStateOfChannel9Err := io.ReadBit()
            if _ActivationStateOfChannel9Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel9' field " + _ActivationStateOfChannel9Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel9)

            // Simple Field (ActivationStateOfChannel10)
            ActivationStateOfChannel10, _ActivationStateOfChannel10Err := io.ReadBit()
            if _ActivationStateOfChannel10Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel10' field " + _ActivationStateOfChannel10Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel10)

            // Simple Field (ActivationStateOfChannel11)
            ActivationStateOfChannel11, _ActivationStateOfChannel11Err := io.ReadBit()
            if _ActivationStateOfChannel11Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel11' field " + _ActivationStateOfChannel11Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel11)

            // Simple Field (ActivationStateOfChannel12)
            ActivationStateOfChannel12, _ActivationStateOfChannel12Err := io.ReadBit()
            if _ActivationStateOfChannel12Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel12' field " + _ActivationStateOfChannel12Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel12)

            // Simple Field (ActivationStateOfChannel13)
            ActivationStateOfChannel13, _ActivationStateOfChannel13Err := io.ReadBit()
            if _ActivationStateOfChannel13Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel13' field " + _ActivationStateOfChannel13Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel13)

            // Simple Field (ActivationStateOfChannel14)
            ActivationStateOfChannel14, _ActivationStateOfChannel14Err := io.ReadBit()
            if _ActivationStateOfChannel14Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel14' field " + _ActivationStateOfChannel14Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel14)

            // Simple Field (ActivationStateOfChannel15)
            ActivationStateOfChannel15, _ActivationStateOfChannel15Err := io.ReadBit()
            if _ActivationStateOfChannel15Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel15' field " + _ActivationStateOfChannel15Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel15)

            // Simple Field (ActivationStateOfChannel16)
            ActivationStateOfChannel16, _ActivationStateOfChannel16Err := io.ReadBit()
            if _ActivationStateOfChannel16Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel16' field " + _ActivationStateOfChannel16Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel16)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_OnOffAction: // STRING
        case datapointType == KnxDatapointType.DPT_Alarm_Reaction: // STRING
        case datapointType == KnxDatapointType.DPT_UpDown_Action: // STRING
        case datapointType == KnxDatapointType.DPT_HVAC_PB_Action: // STRING
        case datapointType == KnxDatapointType.DPT_DoubleNibble: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Busy)
            Busy, _BusyErr := io.ReadUint8(4)
            if _BusyErr != nil {
                return nil, errors.New("Error parsing 'Busy' field " + _BusyErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Busy)

            // Simple Field (Nak)
            Nak, _NakErr := io.ReadUint8(4)
            if _NakErr != nil {
                return nil, errors.New("Error parsing 'Nak' field " + _NakErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Nak)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_SceneInfo: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(1); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (SceneIsInactive)
            SceneIsInactive, _SceneIsInactiveErr := io.ReadBit()
            if _SceneIsInactiveErr != nil {
                return nil, errors.New("Error parsing 'SceneIsInactive' field " + _SceneIsInactiveErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(SceneIsInactive)

            // Simple Field (Scenenumber)
            Scenenumber, _ScenenumberErr := io.ReadUint8(6)
            if _ScenenumberErr != nil {
                return nil, errors.New("Error parsing 'Scenenumber' field " + _ScenenumberErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Scenenumber)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_CombinedInfoOnOff: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput16)
            MaskBitInfoOnOffOutput16, _MaskBitInfoOnOffOutput16Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput16Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput16' field " + _MaskBitInfoOnOffOutput16Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput16)

            // Simple Field (MaskBitInfoOnOffOutput15)
            MaskBitInfoOnOffOutput15, _MaskBitInfoOnOffOutput15Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput15Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput15' field " + _MaskBitInfoOnOffOutput15Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput15)

            // Simple Field (MaskBitInfoOnOffOutput14)
            MaskBitInfoOnOffOutput14, _MaskBitInfoOnOffOutput14Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput14Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput14' field " + _MaskBitInfoOnOffOutput14Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput14)

            // Simple Field (MaskBitInfoOnOffOutput13)
            MaskBitInfoOnOffOutput13, _MaskBitInfoOnOffOutput13Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput13Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput13' field " + _MaskBitInfoOnOffOutput13Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput13)

            // Simple Field (MaskBitInfoOnOffOutput12)
            MaskBitInfoOnOffOutput12, _MaskBitInfoOnOffOutput12Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput12Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput12' field " + _MaskBitInfoOnOffOutput12Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput12)

            // Simple Field (MaskBitInfoOnOffOutput11)
            MaskBitInfoOnOffOutput11, _MaskBitInfoOnOffOutput11Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput11Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput11' field " + _MaskBitInfoOnOffOutput11Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput11)

            // Simple Field (MaskBitInfoOnOffOutput10)
            MaskBitInfoOnOffOutput10, _MaskBitInfoOnOffOutput10Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput10Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput10' field " + _MaskBitInfoOnOffOutput10Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput10)

            // Simple Field (MaskBitInfoOnOffOutput9)
            MaskBitInfoOnOffOutput9, _MaskBitInfoOnOffOutput9Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput9Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput9' field " + _MaskBitInfoOnOffOutput9Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput9)

            // Simple Field (MaskBitInfoOnOffOutput8)
            MaskBitInfoOnOffOutput8, _MaskBitInfoOnOffOutput8Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput8Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput8' field " + _MaskBitInfoOnOffOutput8Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput8)

            // Simple Field (MaskBitInfoOnOffOutput7)
            MaskBitInfoOnOffOutput7, _MaskBitInfoOnOffOutput7Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput7Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput7' field " + _MaskBitInfoOnOffOutput7Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput7)

            // Simple Field (MaskBitInfoOnOffOutput6)
            MaskBitInfoOnOffOutput6, _MaskBitInfoOnOffOutput6Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput6Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput6' field " + _MaskBitInfoOnOffOutput6Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput6)

            // Simple Field (MaskBitInfoOnOffOutput5)
            MaskBitInfoOnOffOutput5, _MaskBitInfoOnOffOutput5Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput5Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput5' field " + _MaskBitInfoOnOffOutput5Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput5)

            // Simple Field (MaskBitInfoOnOffOutput4)
            MaskBitInfoOnOffOutput4, _MaskBitInfoOnOffOutput4Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput4Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput4' field " + _MaskBitInfoOnOffOutput4Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput4)

            // Simple Field (MaskBitInfoOnOffOutput3)
            MaskBitInfoOnOffOutput3, _MaskBitInfoOnOffOutput3Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput3Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput3' field " + _MaskBitInfoOnOffOutput3Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput3)

            // Simple Field (MaskBitInfoOnOffOutput2)
            MaskBitInfoOnOffOutput2, _MaskBitInfoOnOffOutput2Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput2Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput2' field " + _MaskBitInfoOnOffOutput2Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput2)

            // Simple Field (MaskBitInfoOnOffOutput1)
            MaskBitInfoOnOffOutput1, _MaskBitInfoOnOffOutput1Err := io.ReadBit()
            if _MaskBitInfoOnOffOutput1Err != nil {
                return nil, errors.New("Error parsing 'MaskBitInfoOnOffOutput1' field " + _MaskBitInfoOnOffOutput1Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MaskBitInfoOnOffOutput1)

            // Simple Field (InfoOnOffOutput16)
            InfoOnOffOutput16, _InfoOnOffOutput16Err := io.ReadBit()
            if _InfoOnOffOutput16Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput16' field " + _InfoOnOffOutput16Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput16)

            // Simple Field (InfoOnOffOutput15)
            InfoOnOffOutput15, _InfoOnOffOutput15Err := io.ReadBit()
            if _InfoOnOffOutput15Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput15' field " + _InfoOnOffOutput15Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput15)

            // Simple Field (InfoOnOffOutput14)
            InfoOnOffOutput14, _InfoOnOffOutput14Err := io.ReadBit()
            if _InfoOnOffOutput14Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput14' field " + _InfoOnOffOutput14Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput14)

            // Simple Field (InfoOnOffOutput13)
            InfoOnOffOutput13, _InfoOnOffOutput13Err := io.ReadBit()
            if _InfoOnOffOutput13Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput13' field " + _InfoOnOffOutput13Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput13)

            // Simple Field (InfoOnOffOutput12)
            InfoOnOffOutput12, _InfoOnOffOutput12Err := io.ReadBit()
            if _InfoOnOffOutput12Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput12' field " + _InfoOnOffOutput12Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput12)

            // Simple Field (InfoOnOffOutput11)
            InfoOnOffOutput11, _InfoOnOffOutput11Err := io.ReadBit()
            if _InfoOnOffOutput11Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput11' field " + _InfoOnOffOutput11Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput11)

            // Simple Field (InfoOnOffOutput10)
            InfoOnOffOutput10, _InfoOnOffOutput10Err := io.ReadBit()
            if _InfoOnOffOutput10Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput10' field " + _InfoOnOffOutput10Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput10)

            // Simple Field (InfoOnOffOutput9)
            InfoOnOffOutput9, _InfoOnOffOutput9Err := io.ReadBit()
            if _InfoOnOffOutput9Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput9' field " + _InfoOnOffOutput9Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput9)

            // Simple Field (InfoOnOffOutput8)
            InfoOnOffOutput8, _InfoOnOffOutput8Err := io.ReadBit()
            if _InfoOnOffOutput8Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput8' field " + _InfoOnOffOutput8Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput8)

            // Simple Field (InfoOnOffOutput7)
            InfoOnOffOutput7, _InfoOnOffOutput7Err := io.ReadBit()
            if _InfoOnOffOutput7Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput7' field " + _InfoOnOffOutput7Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput7)

            // Simple Field (InfoOnOffOutput6)
            InfoOnOffOutput6, _InfoOnOffOutput6Err := io.ReadBit()
            if _InfoOnOffOutput6Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput6' field " + _InfoOnOffOutput6Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput6)

            // Simple Field (InfoOnOffOutput5)
            InfoOnOffOutput5, _InfoOnOffOutput5Err := io.ReadBit()
            if _InfoOnOffOutput5Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput5' field " + _InfoOnOffOutput5Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput5)

            // Simple Field (InfoOnOffOutput4)
            InfoOnOffOutput4, _InfoOnOffOutput4Err := io.ReadBit()
            if _InfoOnOffOutput4Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput4' field " + _InfoOnOffOutput4Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput4)

            // Simple Field (InfoOnOffOutput3)
            InfoOnOffOutput3, _InfoOnOffOutput3Err := io.ReadBit()
            if _InfoOnOffOutput3Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput3' field " + _InfoOnOffOutput3Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput3)

            // Simple Field (InfoOnOffOutput2)
            InfoOnOffOutput2, _InfoOnOffOutput2Err := io.ReadBit()
            if _InfoOnOffOutput2Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput2' field " + _InfoOnOffOutput2Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput2)

            // Simple Field (InfoOnOffOutput1)
            InfoOnOffOutput1, _InfoOnOffOutput1Err := io.ReadBit()
            if _InfoOnOffOutput1Err != nil {
                return nil, errors.New("Error parsing 'InfoOnOffOutput1' field " + _InfoOnOffOutput1Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(InfoOnOffOutput1)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_ActiveEnergy_V64: // LINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt64(64)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcLINT(value), nil
        case datapointType == KnxDatapointType.DPT_ApparantEnergy_V64: // LINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt64(64)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcLINT(value), nil
        case datapointType == KnxDatapointType.DPT_ReactiveEnergy_V64: // LINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt64(64)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcLINT(value), nil
        case datapointType == KnxDatapointType.DPT_Channel_Activation_24: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel1)
            ActivationStateOfChannel1, _ActivationStateOfChannel1Err := io.ReadBit()
            if _ActivationStateOfChannel1Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel1' field " + _ActivationStateOfChannel1Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel1)

            // Simple Field (ActivationStateOfChannel2)
            ActivationStateOfChannel2, _ActivationStateOfChannel2Err := io.ReadBit()
            if _ActivationStateOfChannel2Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel2' field " + _ActivationStateOfChannel2Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel2)

            // Simple Field (ActivationStateOfChannel3)
            ActivationStateOfChannel3, _ActivationStateOfChannel3Err := io.ReadBit()
            if _ActivationStateOfChannel3Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel3' field " + _ActivationStateOfChannel3Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel3)

            // Simple Field (ActivationStateOfChannel4)
            ActivationStateOfChannel4, _ActivationStateOfChannel4Err := io.ReadBit()
            if _ActivationStateOfChannel4Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel4' field " + _ActivationStateOfChannel4Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel4)

            // Simple Field (ActivationStateOfChannel5)
            ActivationStateOfChannel5, _ActivationStateOfChannel5Err := io.ReadBit()
            if _ActivationStateOfChannel5Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel5' field " + _ActivationStateOfChannel5Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel5)

            // Simple Field (ActivationStateOfChannel6)
            ActivationStateOfChannel6, _ActivationStateOfChannel6Err := io.ReadBit()
            if _ActivationStateOfChannel6Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel6' field " + _ActivationStateOfChannel6Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel6)

            // Simple Field (ActivationStateOfChannel7)
            ActivationStateOfChannel7, _ActivationStateOfChannel7Err := io.ReadBit()
            if _ActivationStateOfChannel7Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel7' field " + _ActivationStateOfChannel7Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel7)

            // Simple Field (ActivationStateOfChannel8)
            ActivationStateOfChannel8, _ActivationStateOfChannel8Err := io.ReadBit()
            if _ActivationStateOfChannel8Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel8' field " + _ActivationStateOfChannel8Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel8)

            // Simple Field (ActivationStateOfChannel9)
            ActivationStateOfChannel9, _ActivationStateOfChannel9Err := io.ReadBit()
            if _ActivationStateOfChannel9Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel9' field " + _ActivationStateOfChannel9Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel9)

            // Simple Field (ActivationStateOfChannel10)
            ActivationStateOfChannel10, _ActivationStateOfChannel10Err := io.ReadBit()
            if _ActivationStateOfChannel10Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel10' field " + _ActivationStateOfChannel10Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel10)

            // Simple Field (ActivationStateOfChannel11)
            ActivationStateOfChannel11, _ActivationStateOfChannel11Err := io.ReadBit()
            if _ActivationStateOfChannel11Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel11' field " + _ActivationStateOfChannel11Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel11)

            // Simple Field (ActivationStateOfChannel12)
            ActivationStateOfChannel12, _ActivationStateOfChannel12Err := io.ReadBit()
            if _ActivationStateOfChannel12Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel12' field " + _ActivationStateOfChannel12Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel12)

            // Simple Field (ActivationStateOfChannel13)
            ActivationStateOfChannel13, _ActivationStateOfChannel13Err := io.ReadBit()
            if _ActivationStateOfChannel13Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel13' field " + _ActivationStateOfChannel13Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel13)

            // Simple Field (ActivationStateOfChannel14)
            ActivationStateOfChannel14, _ActivationStateOfChannel14Err := io.ReadBit()
            if _ActivationStateOfChannel14Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel14' field " + _ActivationStateOfChannel14Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel14)

            // Simple Field (ActivationStateOfChannel15)
            ActivationStateOfChannel15, _ActivationStateOfChannel15Err := io.ReadBit()
            if _ActivationStateOfChannel15Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel15' field " + _ActivationStateOfChannel15Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel15)

            // Simple Field (ActivationStateOfChannel16)
            ActivationStateOfChannel16, _ActivationStateOfChannel16Err := io.ReadBit()
            if _ActivationStateOfChannel16Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel16' field " + _ActivationStateOfChannel16Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel16)

            // Simple Field (ActivationStateOfChannel17)
            ActivationStateOfChannel17, _ActivationStateOfChannel17Err := io.ReadBit()
            if _ActivationStateOfChannel17Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel17' field " + _ActivationStateOfChannel17Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel17)

            // Simple Field (ActivationStateOfChannel18)
            ActivationStateOfChannel18, _ActivationStateOfChannel18Err := io.ReadBit()
            if _ActivationStateOfChannel18Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel18' field " + _ActivationStateOfChannel18Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel18)

            // Simple Field (ActivationStateOfChannel19)
            ActivationStateOfChannel19, _ActivationStateOfChannel19Err := io.ReadBit()
            if _ActivationStateOfChannel19Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel19' field " + _ActivationStateOfChannel19Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel19)

            // Simple Field (ActivationStateOfChannel20)
            ActivationStateOfChannel20, _ActivationStateOfChannel20Err := io.ReadBit()
            if _ActivationStateOfChannel20Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel20' field " + _ActivationStateOfChannel20Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel20)

            // Simple Field (ActivationStateOfChannel21)
            ActivationStateOfChannel21, _ActivationStateOfChannel21Err := io.ReadBit()
            if _ActivationStateOfChannel21Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel21' field " + _ActivationStateOfChannel21Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel21)

            // Simple Field (ActivationStateOfChannel22)
            ActivationStateOfChannel22, _ActivationStateOfChannel22Err := io.ReadBit()
            if _ActivationStateOfChannel22Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel22' field " + _ActivationStateOfChannel22Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel22)

            // Simple Field (ActivationStateOfChannel23)
            ActivationStateOfChannel23, _ActivationStateOfChannel23Err := io.ReadBit()
            if _ActivationStateOfChannel23Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel23' field " + _ActivationStateOfChannel23Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel23)

            // Simple Field (ActivationStateOfChannel24)
            ActivationStateOfChannel24, _ActivationStateOfChannel24Err := io.ReadBit()
            if _ActivationStateOfChannel24Err != nil {
                return nil, errors.New("Error parsing 'ActivationStateOfChannel24' field " + _ActivationStateOfChannel24Err.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActivationStateOfChannel24)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_HVACModeNext: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (DelayTimeMin)
            DelayTimeMin, _DelayTimeMinErr := io.ReadUint16(16)
            if _DelayTimeMinErr != nil {
                return nil, errors.New("Error parsing 'DelayTimeMin' field " + _DelayTimeMinErr.Error())
            }
            _map["Struct"] = values.NewPlcUINT(DelayTimeMin)

            // Simple Field (HvacMode)
            HvacMode, _HvacModeErr := io.ReadUint8(8)
            if _HvacModeErr != nil {
                return nil, errors.New("Error parsing 'HvacMode' field " + _HvacModeErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(HvacMode)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_DHWModeNext: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (DelayTimeMin)
            DelayTimeMin, _DelayTimeMinErr := io.ReadUint16(16)
            if _DelayTimeMinErr != nil {
                return nil, errors.New("Error parsing 'DelayTimeMin' field " + _DelayTimeMinErr.Error())
            }
            _map["Struct"] = values.NewPlcUINT(DelayTimeMin)

            // Simple Field (DhwMode)
            DhwMode, _DhwModeErr := io.ReadUint8(8)
            if _DhwModeErr != nil {
                return nil, errors.New("Error parsing 'DhwMode' field " + _DhwModeErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(DhwMode)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_OccModeNext: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (DelayTimeMin)
            DelayTimeMin, _DelayTimeMinErr := io.ReadUint16(16)
            if _DelayTimeMinErr != nil {
                return nil, errors.New("Error parsing 'DelayTimeMin' field " + _DelayTimeMinErr.Error())
            }
            _map["Struct"] = values.NewPlcUINT(DelayTimeMin)

            // Simple Field (OccupancyMode)
            OccupancyMode, _OccupancyModeErr := io.ReadUint8(8)
            if _OccupancyModeErr != nil {
                return nil, errors.New("Error parsing 'OccupancyMode' field " + _OccupancyModeErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(OccupancyMode)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_BuildingModeNext: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (DelayTimeMin)
            DelayTimeMin, _DelayTimeMinErr := io.ReadUint16(16)
            if _DelayTimeMinErr != nil {
                return nil, errors.New("Error parsing 'DelayTimeMin' field " + _DelayTimeMinErr.Error())
            }
            _map["Struct"] = values.NewPlcUINT(DelayTimeMin)

            // Simple Field (BuildingMode)
            BuildingMode, _BuildingModeErr := io.ReadUint8(8)
            if _BuildingModeErr != nil {
                return nil, errors.New("Error parsing 'BuildingMode' field " + _BuildingModeErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(BuildingMode)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Version: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (MagicNumber)
            MagicNumber, _MagicNumberErr := io.ReadUint8(5)
            if _MagicNumberErr != nil {
                return nil, errors.New("Error parsing 'MagicNumber' field " + _MagicNumberErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(MagicNumber)

            // Simple Field (VersionNumber)
            VersionNumber, _VersionNumberErr := io.ReadUint8(5)
            if _VersionNumberErr != nil {
                return nil, errors.New("Error parsing 'VersionNumber' field " + _VersionNumberErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(VersionNumber)

            // Simple Field (RevisionNumber)
            RevisionNumber, _RevisionNumberErr := io.ReadUint8(6)
            if _RevisionNumberErr != nil {
                return nil, errors.New("Error parsing 'RevisionNumber' field " + _RevisionNumberErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(RevisionNumber)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_AlarmInfo: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (LogNumber)
            LogNumber, _LogNumberErr := io.ReadUint8(8)
            if _LogNumberErr != nil {
                return nil, errors.New("Error parsing 'LogNumber' field " + _LogNumberErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(LogNumber)

            // Simple Field (AlarmPriority)
            AlarmPriority, _AlarmPriorityErr := io.ReadUint8(8)
            if _AlarmPriorityErr != nil {
                return nil, errors.New("Error parsing 'AlarmPriority' field " + _AlarmPriorityErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(AlarmPriority)

            // Simple Field (ApplicationArea)
            ApplicationArea, _ApplicationAreaErr := io.ReadUint8(8)
            if _ApplicationAreaErr != nil {
                return nil, errors.New("Error parsing 'ApplicationArea' field " + _ApplicationAreaErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(ApplicationArea)

            // Simple Field (ErrorClass)
            ErrorClass, _ErrorClassErr := io.ReadUint8(8)
            if _ErrorClassErr != nil {
                return nil, errors.New("Error parsing 'ErrorClass' field " + _ErrorClassErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(ErrorClass)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (ErrorcodeSup)
            ErrorcodeSup, _ErrorcodeSupErr := io.ReadBit()
            if _ErrorcodeSupErr != nil {
                return nil, errors.New("Error parsing 'ErrorcodeSup' field " + _ErrorcodeSupErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ErrorcodeSup)

            // Simple Field (AlarmtextSup)
            AlarmtextSup, _AlarmtextSupErr := io.ReadBit()
            if _AlarmtextSupErr != nil {
                return nil, errors.New("Error parsing 'AlarmtextSup' field " + _AlarmtextSupErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(AlarmtextSup)

            // Simple Field (TimestampSup)
            TimestampSup, _TimestampSupErr := io.ReadBit()
            if _TimestampSupErr != nil {
                return nil, errors.New("Error parsing 'TimestampSup' field " + _TimestampSupErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(TimestampSup)

            // Simple Field (AckSup)
            AckSup, _AckSupErr := io.ReadBit()
            if _AckSupErr != nil {
                return nil, errors.New("Error parsing 'AckSup' field " + _AckSupErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(AckSup)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Locked)
            Locked, _LockedErr := io.ReadBit()
            if _LockedErr != nil {
                return nil, errors.New("Error parsing 'Locked' field " + _LockedErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Locked)

            // Simple Field (Alarmunack)
            Alarmunack, _AlarmunackErr := io.ReadBit()
            if _AlarmunackErr != nil {
                return nil, errors.New("Error parsing 'Alarmunack' field " + _AlarmunackErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Alarmunack)

            // Simple Field (Inalarm)
            Inalarm, _InalarmErr := io.ReadBit()
            if _InalarmErr != nil {
                return nil, errors.New("Error parsing 'Inalarm' field " + _InalarmErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Inalarm)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_TempRoomSetpSetF16_3: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Tempsetpcomf)
            Tempsetpcomf, _TempsetpcomfErr := io.ReadFloat32(true, 4, 11)
            if _TempsetpcomfErr != nil {
                return nil, errors.New("Error parsing 'Tempsetpcomf' field " + _TempsetpcomfErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(Tempsetpcomf)

            // Simple Field (Tempsetpstdby)
            Tempsetpstdby, _TempsetpstdbyErr := io.ReadFloat32(true, 4, 11)
            if _TempsetpstdbyErr != nil {
                return nil, errors.New("Error parsing 'Tempsetpstdby' field " + _TempsetpstdbyErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(Tempsetpstdby)

            // Simple Field (Tempsetpeco)
            Tempsetpeco, _TempsetpecoErr := io.ReadFloat32(true, 4, 11)
            if _TempsetpecoErr != nil {
                return nil, errors.New("Error parsing 'Tempsetpeco' field " + _TempsetpecoErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(Tempsetpeco)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_TempRoomSetpSetShiftF16_3: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Tempsetpshiftcomf)
            Tempsetpshiftcomf, _TempsetpshiftcomfErr := io.ReadFloat32(true, 4, 11)
            if _TempsetpshiftcomfErr != nil {
                return nil, errors.New("Error parsing 'Tempsetpshiftcomf' field " + _TempsetpshiftcomfErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(Tempsetpshiftcomf)

            // Simple Field (Tempsetpshiftstdby)
            Tempsetpshiftstdby, _TempsetpshiftstdbyErr := io.ReadFloat32(true, 4, 11)
            if _TempsetpshiftstdbyErr != nil {
                return nil, errors.New("Error parsing 'Tempsetpshiftstdby' field " + _TempsetpshiftstdbyErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(Tempsetpshiftstdby)

            // Simple Field (Tempsetpshifteco)
            Tempsetpshifteco, _TempsetpshiftecoErr := io.ReadFloat32(true, 4, 11)
            if _TempsetpshiftecoErr != nil {
                return nil, errors.New("Error parsing 'Tempsetpshifteco' field " + _TempsetpshiftecoErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(Tempsetpshifteco)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Scaling_Speed: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (TimePeriod)
            TimePeriod, _TimePeriodErr := io.ReadUint16(16)
            if _TimePeriodErr != nil {
                return nil, errors.New("Error parsing 'TimePeriod' field " + _TimePeriodErr.Error())
            }
            _map["Struct"] = values.NewPlcUINT(TimePeriod)

            // Simple Field (Percent)
            Percent, _PercentErr := io.ReadUint8(8)
            if _PercentErr != nil {
                return nil, errors.New("Error parsing 'Percent' field " + _PercentErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Percent)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Scaling_Step_Time: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (TimePeriod)
            TimePeriod, _TimePeriodErr := io.ReadUint16(16)
            if _TimePeriodErr != nil {
                return nil, errors.New("Error parsing 'TimePeriod' field " + _TimePeriodErr.Error())
            }
            _map["Struct"] = values.NewPlcUINT(TimePeriod)

            // Simple Field (Percent)
            Percent, _PercentErr := io.ReadUint8(8)
            if _PercentErr != nil {
                return nil, errors.New("Error parsing 'Percent' field " + _PercentErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Percent)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_MeteringValue: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Countval)
            Countval, _CountvalErr := io.ReadInt32(32)
            if _CountvalErr != nil {
                return nil, errors.New("Error parsing 'Countval' field " + _CountvalErr.Error())
            }
            _map["Struct"] = values.NewPlcDINT(Countval)

            // Simple Field (Valinffield)
            Valinffield, _ValinffieldErr := io.ReadUint8(8)
            if _ValinffieldErr != nil {
                return nil, errors.New("Error parsing 'Valinffield' field " + _ValinffieldErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Valinffield)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(3); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Alarmunack)
            Alarmunack, _AlarmunackErr := io.ReadBit()
            if _AlarmunackErr != nil {
                return nil, errors.New("Error parsing 'Alarmunack' field " + _AlarmunackErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Alarmunack)

            // Simple Field (Inalarm)
            Inalarm, _InalarmErr := io.ReadBit()
            if _InalarmErr != nil {
                return nil, errors.New("Error parsing 'Inalarm' field " + _InalarmErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Inalarm)

            // Simple Field (Overridden)
            Overridden, _OverriddenErr := io.ReadBit()
            if _OverriddenErr != nil {
                return nil, errors.New("Error parsing 'Overridden' field " + _OverriddenErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Overridden)

            // Simple Field (Fault)
            Fault, _FaultErr := io.ReadBit()
            if _FaultErr != nil {
                return nil, errors.New("Error parsing 'Fault' field " + _FaultErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Fault)

            // Simple Field (Outofservice)
            Outofservice, _OutofserviceErr := io.ReadBit()
            if _OutofserviceErr != nil {
                return nil, errors.New("Error parsing 'Outofservice' field " + _OutofserviceErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Outofservice)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_MBus_Address: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Manufactid)
            Manufactid, _ManufactidErr := io.ReadUint16(16)
            if _ManufactidErr != nil {
                return nil, errors.New("Error parsing 'Manufactid' field " + _ManufactidErr.Error())
            }
            _map["Struct"] = values.NewPlcUINT(Manufactid)

            // Simple Field (Identnumber)
            Identnumber, _IdentnumberErr := io.ReadUint32(32)
            if _IdentnumberErr != nil {
                return nil, errors.New("Error parsing 'Identnumber' field " + _IdentnumberErr.Error())
            }
            _map["Struct"] = values.NewPlcUDINT(Identnumber)

            // Simple Field (Version)
            Version, _VersionErr := io.ReadUint8(8)
            if _VersionErr != nil {
                return nil, errors.New("Error parsing 'Version' field " + _VersionErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Version)

            // Simple Field (Medium)
            Medium, _MediumErr := io.ReadUint8(8)
            if _MediumErr != nil {
                return nil, errors.New("Error parsing 'Medium' field " + _MediumErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Medium)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Colour_RGB: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (R)
            R, _RErr := io.ReadUint8(8)
            if _RErr != nil {
                return nil, errors.New("Error parsing 'R' field " + _RErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(R)

            // Simple Field (G)
            G, _GErr := io.ReadUint8(8)
            if _GErr != nil {
                return nil, errors.New("Error parsing 'G' field " + _GErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(G)

            // Simple Field (B)
            B, _BErr := io.ReadUint8(8)
            if _BErr != nil {
                return nil, errors.New("Error parsing 'B' field " + _BErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(B)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_LanguageCodeAlpha2_ASCII: // STRING

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadString(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSTRING(value), nil
        case datapointType == KnxDatapointType.DPT_Tariff_ActiveEnergy: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Activeelectricalenergy)
            Activeelectricalenergy, _ActiveelectricalenergyErr := io.ReadInt32(32)
            if _ActiveelectricalenergyErr != nil {
                return nil, errors.New("Error parsing 'Activeelectricalenergy' field " + _ActiveelectricalenergyErr.Error())
            }
            _map["Struct"] = values.NewPlcDINT(Activeelectricalenergy)

            // Simple Field (Tariff)
            Tariff, _TariffErr := io.ReadUint8(8)
            if _TariffErr != nil {
                return nil, errors.New("Error parsing 'Tariff' field " + _TariffErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Tariff)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Electricalengergyvalidity)
            Electricalengergyvalidity, _ElectricalengergyvalidityErr := io.ReadBit()
            if _ElectricalengergyvalidityErr != nil {
                return nil, errors.New("Error parsing 'Electricalengergyvalidity' field " + _ElectricalengergyvalidityErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Electricalengergyvalidity)

            // Simple Field (Tariffvalidity)
            Tariffvalidity, _TariffvalidityErr := io.ReadBit()
            if _TariffvalidityErr != nil {
                return nil, errors.New("Error parsing 'Tariffvalidity' field " + _TariffvalidityErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Tariffvalidity)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Prioritised_Mode_Control: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (DeactivationOfPriority)
            DeactivationOfPriority, _DeactivationOfPriorityErr := io.ReadBit()
            if _DeactivationOfPriorityErr != nil {
                return nil, errors.New("Error parsing 'DeactivationOfPriority' field " + _DeactivationOfPriorityErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(DeactivationOfPriority)

            // Simple Field (PriorityLevel)
            PriorityLevel, _PriorityLevelErr := io.ReadUint8(3)
            if _PriorityLevelErr != nil {
                return nil, errors.New("Error parsing 'PriorityLevel' field " + _PriorityLevelErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(PriorityLevel)

            // Simple Field (ModeLevel)
            ModeLevel, _ModeLevelErr := io.ReadUint8(4)
            if _ModeLevelErr != nil {
                return nil, errors.New("Error parsing 'ModeLevel' field " + _ModeLevelErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(ModeLevel)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_DALI_Control_Gear_Diagnostic: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (ConvertorError)
            ConvertorError, _ConvertorErrorErr := io.ReadBit()
            if _ConvertorErrorErr != nil {
                return nil, errors.New("Error parsing 'ConvertorError' field " + _ConvertorErrorErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ConvertorError)

            // Simple Field (BallastFailure)
            BallastFailure, _BallastFailureErr := io.ReadBit()
            if _BallastFailureErr != nil {
                return nil, errors.New("Error parsing 'BallastFailure' field " + _BallastFailureErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(BallastFailure)

            // Simple Field (LampFailure)
            LampFailure, _LampFailureErr := io.ReadBit()
            if _LampFailureErr != nil {
                return nil, errors.New("Error parsing 'LampFailure' field " + _LampFailureErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(LampFailure)

            // Simple Field (ReadOrResponse)
            ReadOrResponse, _ReadOrResponseErr := io.ReadBit()
            if _ReadOrResponseErr != nil {
                return nil, errors.New("Error parsing 'ReadOrResponse' field " + _ReadOrResponseErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ReadOrResponse)

            // Simple Field (AddressIndicator)
            AddressIndicator, _AddressIndicatorErr := io.ReadBit()
            if _AddressIndicatorErr != nil {
                return nil, errors.New("Error parsing 'AddressIndicator' field " + _AddressIndicatorErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(AddressIndicator)

            // Simple Field (DaliDeviceAddressOrDaliGroupAddress)
            DaliDeviceAddressOrDaliGroupAddress, _DaliDeviceAddressOrDaliGroupAddressErr := io.ReadUint8(6)
            if _DaliDeviceAddressOrDaliGroupAddressErr != nil {
                return nil, errors.New("Error parsing 'DaliDeviceAddressOrDaliGroupAddress' field " + _DaliDeviceAddressOrDaliGroupAddressErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(DaliDeviceAddressOrDaliGroupAddress)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_DALI_Diagnostics: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (BallastFailure)
            BallastFailure, _BallastFailureErr := io.ReadBit()
            if _BallastFailureErr != nil {
                return nil, errors.New("Error parsing 'BallastFailure' field " + _BallastFailureErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(BallastFailure)

            // Simple Field (LampFailure)
            LampFailure, _LampFailureErr := io.ReadBit()
            if _LampFailureErr != nil {
                return nil, errors.New("Error parsing 'LampFailure' field " + _LampFailureErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(LampFailure)

            // Simple Field (DeviceAddress)
            DeviceAddress, _DeviceAddressErr := io.ReadUint8(6)
            if _DeviceAddressErr != nil {
                return nil, errors.New("Error parsing 'DeviceAddress' field " + _DeviceAddressErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(DeviceAddress)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_CombinedPosition: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (HeightPosition)
            HeightPosition, _HeightPositionErr := io.ReadUint8(8)
            if _HeightPositionErr != nil {
                return nil, errors.New("Error parsing 'HeightPosition' field " + _HeightPositionErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(HeightPosition)

            // Simple Field (SlatsPosition)
            SlatsPosition, _SlatsPositionErr := io.ReadUint8(8)
            if _SlatsPositionErr != nil {
                return nil, errors.New("Error parsing 'SlatsPosition' field " + _SlatsPositionErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(SlatsPosition)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (ValidityHeightPosition)
            ValidityHeightPosition, _ValidityHeightPositionErr := io.ReadBit()
            if _ValidityHeightPositionErr != nil {
                return nil, errors.New("Error parsing 'ValidityHeightPosition' field " + _ValidityHeightPositionErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ValidityHeightPosition)

            // Simple Field (ValiditySlatsPosition)
            ValiditySlatsPosition, _ValiditySlatsPositionErr := io.ReadBit()
            if _ValiditySlatsPositionErr != nil {
                return nil, errors.New("Error parsing 'ValiditySlatsPosition' field " + _ValiditySlatsPositionErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ValiditySlatsPosition)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_StatusSAB: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (HeightPosition)
            HeightPosition, _HeightPositionErr := io.ReadUint8(8)
            if _HeightPositionErr != nil {
                return nil, errors.New("Error parsing 'HeightPosition' field " + _HeightPositionErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(HeightPosition)

            // Simple Field (SlatsPosition)
            SlatsPosition, _SlatsPositionErr := io.ReadUint8(8)
            if _SlatsPositionErr != nil {
                return nil, errors.New("Error parsing 'SlatsPosition' field " + _SlatsPositionErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(SlatsPosition)

            // Simple Field (UpperEndPosReached)
            UpperEndPosReached, _UpperEndPosReachedErr := io.ReadBit()
            if _UpperEndPosReachedErr != nil {
                return nil, errors.New("Error parsing 'UpperEndPosReached' field " + _UpperEndPosReachedErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(UpperEndPosReached)

            // Simple Field (LowerEndPosReached)
            LowerEndPosReached, _LowerEndPosReachedErr := io.ReadBit()
            if _LowerEndPosReachedErr != nil {
                return nil, errors.New("Error parsing 'LowerEndPosReached' field " + _LowerEndPosReachedErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(LowerEndPosReached)

            // Simple Field (LowerPredefPosReachedTypHeight100PercentSlatsAngle100Percent)
            LowerPredefPosReachedTypHeight100PercentSlatsAngle100Percent, _LowerPredefPosReachedTypHeight100PercentSlatsAngle100PercentErr := io.ReadBit()
            if _LowerPredefPosReachedTypHeight100PercentSlatsAngle100PercentErr != nil {
                return nil, errors.New("Error parsing 'LowerPredefPosReachedTypHeight100PercentSlatsAngle100Percent' field " + _LowerPredefPosReachedTypHeight100PercentSlatsAngle100PercentErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(LowerPredefPosReachedTypHeight100PercentSlatsAngle100Percent)

            // Simple Field (TargetPosDrive)
            TargetPosDrive, _TargetPosDriveErr := io.ReadBit()
            if _TargetPosDriveErr != nil {
                return nil, errors.New("Error parsing 'TargetPosDrive' field " + _TargetPosDriveErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(TargetPosDrive)

            // Simple Field (RestrictionOfTargetHeightPosPosCanNotBeReached)
            RestrictionOfTargetHeightPosPosCanNotBeReached, _RestrictionOfTargetHeightPosPosCanNotBeReachedErr := io.ReadBit()
            if _RestrictionOfTargetHeightPosPosCanNotBeReachedErr != nil {
                return nil, errors.New("Error parsing 'RestrictionOfTargetHeightPosPosCanNotBeReached' field " + _RestrictionOfTargetHeightPosPosCanNotBeReachedErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(RestrictionOfTargetHeightPosPosCanNotBeReached)

            // Simple Field (RestrictionOfSlatsHeightPosPosCanNotBeReached)
            RestrictionOfSlatsHeightPosPosCanNotBeReached, _RestrictionOfSlatsHeightPosPosCanNotBeReachedErr := io.ReadBit()
            if _RestrictionOfSlatsHeightPosPosCanNotBeReachedErr != nil {
                return nil, errors.New("Error parsing 'RestrictionOfSlatsHeightPosPosCanNotBeReached' field " + _RestrictionOfSlatsHeightPosPosCanNotBeReachedErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(RestrictionOfSlatsHeightPosPosCanNotBeReached)

            // Simple Field (AtLeastOneOfTheInputsWindRainFrostAlarmIsInAlarm)
            AtLeastOneOfTheInputsWindRainFrostAlarmIsInAlarm, _AtLeastOneOfTheInputsWindRainFrostAlarmIsInAlarmErr := io.ReadBit()
            if _AtLeastOneOfTheInputsWindRainFrostAlarmIsInAlarmErr != nil {
                return nil, errors.New("Error parsing 'AtLeastOneOfTheInputsWindRainFrostAlarmIsInAlarm' field " + _AtLeastOneOfTheInputsWindRainFrostAlarmIsInAlarmErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(AtLeastOneOfTheInputsWindRainFrostAlarmIsInAlarm)

            // Simple Field (UpDownPositionIsForcedByMoveupdownforcedInput)
            UpDownPositionIsForcedByMoveupdownforcedInput, _UpDownPositionIsForcedByMoveupdownforcedInputErr := io.ReadBit()
            if _UpDownPositionIsForcedByMoveupdownforcedInputErr != nil {
                return nil, errors.New("Error parsing 'UpDownPositionIsForcedByMoveupdownforcedInput' field " + _UpDownPositionIsForcedByMoveupdownforcedInputErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(UpDownPositionIsForcedByMoveupdownforcedInput)

            // Simple Field (MovementIsLockedEGByDevicelockedInput)
            MovementIsLockedEGByDevicelockedInput, _MovementIsLockedEGByDevicelockedInputErr := io.ReadBit()
            if _MovementIsLockedEGByDevicelockedInputErr != nil {
                return nil, errors.New("Error parsing 'MovementIsLockedEGByDevicelockedInput' field " + _MovementIsLockedEGByDevicelockedInputErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(MovementIsLockedEGByDevicelockedInput)

            // Simple Field (ActuatorSetvalueIsLocallyOverriddenEGViaALocalUserInterface)
            ActuatorSetvalueIsLocallyOverriddenEGViaALocalUserInterface, _ActuatorSetvalueIsLocallyOverriddenEGViaALocalUserInterfaceErr := io.ReadBit()
            if _ActuatorSetvalueIsLocallyOverriddenEGViaALocalUserInterfaceErr != nil {
                return nil, errors.New("Error parsing 'ActuatorSetvalueIsLocallyOverriddenEGViaALocalUserInterface' field " + _ActuatorSetvalueIsLocallyOverriddenEGViaALocalUserInterfaceErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ActuatorSetvalueIsLocallyOverriddenEGViaALocalUserInterface)

            // Simple Field (GeneralFailureOfTheActuatorOrTheDrive)
            GeneralFailureOfTheActuatorOrTheDrive, _GeneralFailureOfTheActuatorOrTheDriveErr := io.ReadBit()
            if _GeneralFailureOfTheActuatorOrTheDriveErr != nil {
                return nil, errors.New("Error parsing 'GeneralFailureOfTheActuatorOrTheDrive' field " + _GeneralFailureOfTheActuatorOrTheDriveErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(GeneralFailureOfTheActuatorOrTheDrive)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(3); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (ValidityHeightPos)
            ValidityHeightPos, _ValidityHeightPosErr := io.ReadBit()
            if _ValidityHeightPosErr != nil {
                return nil, errors.New("Error parsing 'ValidityHeightPos' field " + _ValidityHeightPosErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ValidityHeightPos)

            // Simple Field (ValiditySlatsPos)
            ValiditySlatsPos, _ValiditySlatsPosErr := io.ReadBit()
            if _ValiditySlatsPosErr != nil {
                return nil, errors.New("Error parsing 'ValiditySlatsPos' field " + _ValiditySlatsPosErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ValiditySlatsPos)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Colour_xyY: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (XAxis)
            XAxis, _XAxisErr := io.ReadUint16(16)
            if _XAxisErr != nil {
                return nil, errors.New("Error parsing 'XAxis' field " + _XAxisErr.Error())
            }
            _map["Struct"] = values.NewPlcUINT(XAxis)

            // Simple Field (YAxis)
            YAxis, _YAxisErr := io.ReadUint16(16)
            if _YAxisErr != nil {
                return nil, errors.New("Error parsing 'YAxis' field " + _YAxisErr.Error())
            }
            _map["Struct"] = values.NewPlcUINT(YAxis)

            // Simple Field (Brightness)
            Brightness, _BrightnessErr := io.ReadUint8(8)
            if _BrightnessErr != nil {
                return nil, errors.New("Error parsing 'Brightness' field " + _BrightnessErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Brightness)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (ValidityXy)
            ValidityXy, _ValidityXyErr := io.ReadBit()
            if _ValidityXyErr != nil {
                return nil, errors.New("Error parsing 'ValidityXy' field " + _ValidityXyErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ValidityXy)

            // Simple Field (ValidityBrightness)
            ValidityBrightness, _ValidityBrightnessErr := io.ReadBit()
            if _ValidityBrightnessErr != nil {
                return nil, errors.New("Error parsing 'ValidityBrightness' field " + _ValidityBrightnessErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ValidityBrightness)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Converter_Status: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (ConverterModeAccordingToTheDaliConverterStateMachine)
            ConverterModeAccordingToTheDaliConverterStateMachine, _ConverterModeAccordingToTheDaliConverterStateMachineErr := io.ReadUint8(4)
            if _ConverterModeAccordingToTheDaliConverterStateMachineErr != nil {
                return nil, errors.New("Error parsing 'ConverterModeAccordingToTheDaliConverterStateMachine' field " + _ConverterModeAccordingToTheDaliConverterStateMachineErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(ConverterModeAccordingToTheDaliConverterStateMachine)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Hs)
            Hs, _HsErr := io.ReadBit()
            if _HsErr != nil {
                return nil, errors.New("Error parsing 'Hs' field " + _HsErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Hs)

            // Simple Field (Hs)
            Hs, _HsErr := io.ReadBit()
            if _HsErr != nil {
                return nil, errors.New("Error parsing 'Hs' field " + _HsErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Hs)

            // Simple Field (FunctionTestPending)
            FunctionTestPending, _FunctionTestPendingErr := io.ReadUint8(2)
            if _FunctionTestPendingErr != nil {
                return nil, errors.New("Error parsing 'FunctionTestPending' field " + _FunctionTestPendingErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(FunctionTestPending)

            // Simple Field (DurationTestPending)
            DurationTestPending, _DurationTestPendingErr := io.ReadUint8(2)
            if _DurationTestPendingErr != nil {
                return nil, errors.New("Error parsing 'DurationTestPending' field " + _DurationTestPendingErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(DurationTestPending)

            // Simple Field (PartialDurationTestPending)
            PartialDurationTestPending, _PartialDurationTestPendingErr := io.ReadUint8(2)
            if _PartialDurationTestPendingErr != nil {
                return nil, errors.New("Error parsing 'PartialDurationTestPending' field " + _PartialDurationTestPendingErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(PartialDurationTestPending)

            // Simple Field (ConverterFailure)
            ConverterFailure, _ConverterFailureErr := io.ReadUint8(2)
            if _ConverterFailureErr != nil {
                return nil, errors.New("Error parsing 'ConverterFailure' field " + _ConverterFailureErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(ConverterFailure)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Converter_Test_Result: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Ltrf)
            Ltrf, _LtrfErr := io.ReadUint8(4)
            if _LtrfErr != nil {
                return nil, errors.New("Error parsing 'Ltrf' field " + _LtrfErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Ltrf)

            // Simple Field (Ltrd)
            Ltrd, _LtrdErr := io.ReadUint8(4)
            if _LtrdErr != nil {
                return nil, errors.New("Error parsing 'Ltrd' field " + _LtrdErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Ltrd)

            // Simple Field (Ltrp)
            Ltrp, _LtrpErr := io.ReadUint8(4)
            if _LtrpErr != nil {
                return nil, errors.New("Error parsing 'Ltrp' field " + _LtrpErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Ltrp)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Sf)
            Sf, _SfErr := io.ReadUint8(2)
            if _SfErr != nil {
                return nil, errors.New("Error parsing 'Sf' field " + _SfErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Sf)

            // Simple Field (Sd)
            Sd, _SdErr := io.ReadUint8(2)
            if _SdErr != nil {
                return nil, errors.New("Error parsing 'Sd' field " + _SdErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Sd)

            // Simple Field (Sp)
            Sp, _SpErr := io.ReadUint8(2)
            if _SpErr != nil {
                return nil, errors.New("Error parsing 'Sp' field " + _SpErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Sp)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Ldtr)
            Ldtr, _LdtrErr := io.ReadUint16(16)
            if _LdtrErr != nil {
                return nil, errors.New("Error parsing 'Ldtr' field " + _LdtrErr.Error())
            }
            _map["Struct"] = values.NewPlcUINT(Ldtr)

            // Simple Field (Lpdtr)
            Lpdtr, _LpdtrErr := io.ReadUint8(8)
            if _LpdtrErr != nil {
                return nil, errors.New("Error parsing 'Lpdtr' field " + _LpdtrErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Lpdtr)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Battery_Info: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (BatteryFailure)
            BatteryFailure, _BatteryFailureErr := io.ReadBit()
            if _BatteryFailureErr != nil {
                return nil, errors.New("Error parsing 'BatteryFailure' field " + _BatteryFailureErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(BatteryFailure)

            // Simple Field (BatteryDurationFailure)
            BatteryDurationFailure, _BatteryDurationFailureErr := io.ReadBit()
            if _BatteryDurationFailureErr != nil {
                return nil, errors.New("Error parsing 'BatteryDurationFailure' field " + _BatteryDurationFailureErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(BatteryDurationFailure)

            // Simple Field (BatteryFullyCharged)
            BatteryFullyCharged, _BatteryFullyChargedErr := io.ReadBit()
            if _BatteryFullyChargedErr != nil {
                return nil, errors.New("Error parsing 'BatteryFullyCharged' field " + _BatteryFullyChargedErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(BatteryFullyCharged)

            // Simple Field (BatteryChargeLevel)
            BatteryChargeLevel, _BatteryChargeLevelErr := io.ReadUint8(8)
            if _BatteryChargeLevelErr != nil {
                return nil, errors.New("Error parsing 'BatteryChargeLevel' field " + _BatteryChargeLevelErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(BatteryChargeLevel)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Brightness_Colour_Temperature_Transition: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Ms)
            Ms, _MsErr := io.ReadUint16(16)
            if _MsErr != nil {
                return nil, errors.New("Error parsing 'Ms' field " + _MsErr.Error())
            }
            _map["Struct"] = values.NewPlcUINT(Ms)

            // Simple Field (TemperatureK)
            TemperatureK, _TemperatureKErr := io.ReadUint16(16)
            if _TemperatureKErr != nil {
                return nil, errors.New("Error parsing 'TemperatureK' field " + _TemperatureKErr.Error())
            }
            _map["Struct"] = values.NewPlcUINT(TemperatureK)

            // Simple Field (Percent)
            Percent, _PercentErr := io.ReadUint8(8)
            if _PercentErr != nil {
                return nil, errors.New("Error parsing 'Percent' field " + _PercentErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(Percent)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (ValidityOfTheTimePeriod)
            ValidityOfTheTimePeriod, _ValidityOfTheTimePeriodErr := io.ReadBit()
            if _ValidityOfTheTimePeriodErr != nil {
                return nil, errors.New("Error parsing 'ValidityOfTheTimePeriod' field " + _ValidityOfTheTimePeriodErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ValidityOfTheTimePeriod)

            // Simple Field (ValidityOfTheAbsoluteColourTemperature)
            ValidityOfTheAbsoluteColourTemperature, _ValidityOfTheAbsoluteColourTemperatureErr := io.ReadBit()
            if _ValidityOfTheAbsoluteColourTemperatureErr != nil {
                return nil, errors.New("Error parsing 'ValidityOfTheAbsoluteColourTemperature' field " + _ValidityOfTheAbsoluteColourTemperatureErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ValidityOfTheAbsoluteColourTemperature)

            // Simple Field (ValidityOfTheAbsoluteBrightness)
            ValidityOfTheAbsoluteBrightness, _ValidityOfTheAbsoluteBrightnessErr := io.ReadBit()
            if _ValidityOfTheAbsoluteBrightnessErr != nil {
                return nil, errors.New("Error parsing 'ValidityOfTheAbsoluteBrightness' field " + _ValidityOfTheAbsoluteBrightnessErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ValidityOfTheAbsoluteBrightness)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Brightness_Colour_Temperature_Control: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Cct)
            Cct, _CctErr := io.ReadBit()
            if _CctErr != nil {
                return nil, errors.New("Error parsing 'Cct' field " + _CctErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Cct)

            // Simple Field (StepCodeColourTemperature)
            StepCodeColourTemperature, _StepCodeColourTemperatureErr := io.ReadUint8(3)
            if _StepCodeColourTemperatureErr != nil {
                return nil, errors.New("Error parsing 'StepCodeColourTemperature' field " + _StepCodeColourTemperatureErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(StepCodeColourTemperature)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Cb)
            Cb, _CbErr := io.ReadBit()
            if _CbErr != nil {
                return nil, errors.New("Error parsing 'Cb' field " + _CbErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Cb)

            // Simple Field (StepCodeBrightness)
            StepCodeBrightness, _StepCodeBrightnessErr := io.ReadUint8(3)
            if _StepCodeBrightnessErr != nil {
                return nil, errors.New("Error parsing 'StepCodeBrightness' field " + _StepCodeBrightnessErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(StepCodeBrightness)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (CctAndStepCodeColourValidity)
            CctAndStepCodeColourValidity, _CctAndStepCodeColourValidityErr := io.ReadBit()
            if _CctAndStepCodeColourValidityErr != nil {
                return nil, errors.New("Error parsing 'CctAndStepCodeColourValidity' field " + _CctAndStepCodeColourValidityErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(CctAndStepCodeColourValidity)

            // Simple Field (CbAndStepCodeBrightnessValidity)
            CbAndStepCodeBrightnessValidity, _CbAndStepCodeBrightnessValidityErr := io.ReadBit()
            if _CbAndStepCodeBrightnessValidityErr != nil {
                return nil, errors.New("Error parsing 'CbAndStepCodeBrightnessValidity' field " + _CbAndStepCodeBrightnessValidityErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(CbAndStepCodeBrightnessValidity)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Colour_RGBW: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (ColourLevelRed)
            ColourLevelRed, _ColourLevelRedErr := io.ReadUint8(8)
            if _ColourLevelRedErr != nil {
                return nil, errors.New("Error parsing 'ColourLevelRed' field " + _ColourLevelRedErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(ColourLevelRed)

            // Simple Field (ColourLevelGreen)
            ColourLevelGreen, _ColourLevelGreenErr := io.ReadUint8(8)
            if _ColourLevelGreenErr != nil {
                return nil, errors.New("Error parsing 'ColourLevelGreen' field " + _ColourLevelGreenErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(ColourLevelGreen)

            // Simple Field (ColourLevelBlue)
            ColourLevelBlue, _ColourLevelBlueErr := io.ReadUint8(8)
            if _ColourLevelBlueErr != nil {
                return nil, errors.New("Error parsing 'ColourLevelBlue' field " + _ColourLevelBlueErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(ColourLevelBlue)

            // Simple Field (ColourLevelWhite)
            ColourLevelWhite, _ColourLevelWhiteErr := io.ReadUint8(8)
            if _ColourLevelWhiteErr != nil {
                return nil, errors.New("Error parsing 'ColourLevelWhite' field " + _ColourLevelWhiteErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(ColourLevelWhite)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Mr)
            Mr, _MrErr := io.ReadBit()
            if _MrErr != nil {
                return nil, errors.New("Error parsing 'Mr' field " + _MrErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Mr)

            // Simple Field (Mg)
            Mg, _MgErr := io.ReadBit()
            if _MgErr != nil {
                return nil, errors.New("Error parsing 'Mg' field " + _MgErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Mg)

            // Simple Field (Mb)
            Mb, _MbErr := io.ReadBit()
            if _MbErr != nil {
                return nil, errors.New("Error parsing 'Mb' field " + _MbErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Mb)

            // Simple Field (Mw)
            Mw, _MwErr := io.ReadBit()
            if _MwErr != nil {
                return nil, errors.New("Error parsing 'Mw' field " + _MwErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Mw)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Relative_Control_RGBW: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Maskcw)
            Maskcw, _MaskcwErr := io.ReadBit()
            if _MaskcwErr != nil {
                return nil, errors.New("Error parsing 'Maskcw' field " + _MaskcwErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Maskcw)

            // Simple Field (Maskcb)
            Maskcb, _MaskcbErr := io.ReadBit()
            if _MaskcbErr != nil {
                return nil, errors.New("Error parsing 'Maskcb' field " + _MaskcbErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Maskcb)

            // Simple Field (Maskcg)
            Maskcg, _MaskcgErr := io.ReadBit()
            if _MaskcgErr != nil {
                return nil, errors.New("Error parsing 'Maskcg' field " + _MaskcgErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Maskcg)

            // Simple Field (Maskcr)
            Maskcr, _MaskcrErr := io.ReadBit()
            if _MaskcrErr != nil {
                return nil, errors.New("Error parsing 'Maskcr' field " + _MaskcrErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Maskcr)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Cw)
            Cw, _CwErr := io.ReadBit()
            if _CwErr != nil {
                return nil, errors.New("Error parsing 'Cw' field " + _CwErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Cw)

            // Simple Field (StepCodeColourWhite)
            StepCodeColourWhite, _StepCodeColourWhiteErr := io.ReadUint8(3)
            if _StepCodeColourWhiteErr != nil {
                return nil, errors.New("Error parsing 'StepCodeColourWhite' field " + _StepCodeColourWhiteErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(StepCodeColourWhite)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Cb)
            Cb, _CbErr := io.ReadBit()
            if _CbErr != nil {
                return nil, errors.New("Error parsing 'Cb' field " + _CbErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Cb)

            // Simple Field (StepCodeColourBlue)
            StepCodeColourBlue, _StepCodeColourBlueErr := io.ReadUint8(3)
            if _StepCodeColourBlueErr != nil {
                return nil, errors.New("Error parsing 'StepCodeColourBlue' field " + _StepCodeColourBlueErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(StepCodeColourBlue)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Cg)
            Cg, _CgErr := io.ReadBit()
            if _CgErr != nil {
                return nil, errors.New("Error parsing 'Cg' field " + _CgErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Cg)

            // Simple Field (StepCodeColourGreen)
            StepCodeColourGreen, _StepCodeColourGreenErr := io.ReadUint8(3)
            if _StepCodeColourGreenErr != nil {
                return nil, errors.New("Error parsing 'StepCodeColourGreen' field " + _StepCodeColourGreenErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(StepCodeColourGreen)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Cr)
            Cr, _CrErr := io.ReadBit()
            if _CrErr != nil {
                return nil, errors.New("Error parsing 'Cr' field " + _CrErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Cr)

            // Simple Field (StepCodeColourRed)
            StepCodeColourRed, _StepCodeColourRedErr := io.ReadUint8(3)
            if _StepCodeColourRedErr != nil {
                return nil, errors.New("Error parsing 'StepCodeColourRed' field " + _StepCodeColourRedErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(StepCodeColourRed)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_Relative_Control_RGB: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Cb)
            Cb, _CbErr := io.ReadBit()
            if _CbErr != nil {
                return nil, errors.New("Error parsing 'Cb' field " + _CbErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Cb)

            // Simple Field (StepCodeColourBlue)
            StepCodeColourBlue, _StepCodeColourBlueErr := io.ReadUint8(3)
            if _StepCodeColourBlueErr != nil {
                return nil, errors.New("Error parsing 'StepCodeColourBlue' field " + _StepCodeColourBlueErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(StepCodeColourBlue)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Cg)
            Cg, _CgErr := io.ReadBit()
            if _CgErr != nil {
                return nil, errors.New("Error parsing 'Cg' field " + _CgErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Cg)

            // Simple Field (StepCodeColourGreen)
            StepCodeColourGreen, _StepCodeColourGreenErr := io.ReadUint8(3)
            if _StepCodeColourGreenErr != nil {
                return nil, errors.New("Error parsing 'StepCodeColourGreen' field " + _StepCodeColourGreenErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(StepCodeColourGreen)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Cr)
            Cr, _CrErr := io.ReadBit()
            if _CrErr != nil {
                return nil, errors.New("Error parsing 'Cr' field " + _CrErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(Cr)

            // Simple Field (StepCodeColourRed)
            StepCodeColourRed, _StepCodeColourRedErr := io.ReadUint8(3)
            if _StepCodeColourRedErr != nil {
                return nil, errors.New("Error parsing 'StepCodeColourRed' field " + _StepCodeColourRedErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(StepCodeColourRed)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_GeographicalLocation: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (Longitude)
            Longitude, _LongitudeErr := io.ReadFloat32(true, 8, 23)
            if _LongitudeErr != nil {
                return nil, errors.New("Error parsing 'Longitude' field " + _LongitudeErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(Longitude)

            // Simple Field (Latitude)
            Latitude, _LatitudeErr := io.ReadFloat32(true, 8, 23)
            if _LatitudeErr != nil {
                return nil, errors.New("Error parsing 'Latitude' field " + _LatitudeErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(Latitude)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_TempRoomSetpSetF16_4: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (RoomTemperatureSetpointComfort)
            RoomTemperatureSetpointComfort, _RoomTemperatureSetpointComfortErr := io.ReadFloat32(true, 4, 11)
            if _RoomTemperatureSetpointComfortErr != nil {
                return nil, errors.New("Error parsing 'RoomTemperatureSetpointComfort' field " + _RoomTemperatureSetpointComfortErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(RoomTemperatureSetpointComfort)

            // Simple Field (RoomTemperatureSetpointStandby)
            RoomTemperatureSetpointStandby, _RoomTemperatureSetpointStandbyErr := io.ReadFloat32(true, 4, 11)
            if _RoomTemperatureSetpointStandbyErr != nil {
                return nil, errors.New("Error parsing 'RoomTemperatureSetpointStandby' field " + _RoomTemperatureSetpointStandbyErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(RoomTemperatureSetpointStandby)

            // Simple Field (RoomTemperatureSetpointEconomy)
            RoomTemperatureSetpointEconomy, _RoomTemperatureSetpointEconomyErr := io.ReadFloat32(true, 4, 11)
            if _RoomTemperatureSetpointEconomyErr != nil {
                return nil, errors.New("Error parsing 'RoomTemperatureSetpointEconomy' field " + _RoomTemperatureSetpointEconomyErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(RoomTemperatureSetpointEconomy)

            // Simple Field (RoomTemperatureSetpointBuildingProtection)
            RoomTemperatureSetpointBuildingProtection, _RoomTemperatureSetpointBuildingProtectionErr := io.ReadFloat32(true, 4, 11)
            if _RoomTemperatureSetpointBuildingProtectionErr != nil {
                return nil, errors.New("Error parsing 'RoomTemperatureSetpointBuildingProtection' field " + _RoomTemperatureSetpointBuildingProtectionErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(RoomTemperatureSetpointBuildingProtection)
            return values.NewPlcStruct(_map), nil
        case datapointType == KnxDatapointType.DPT_TempRoomSetpSetShiftF16_4: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (RoomTemperatureSetpointShiftComfort)
            RoomTemperatureSetpointShiftComfort, _RoomTemperatureSetpointShiftComfortErr := io.ReadFloat32(true, 4, 11)
            if _RoomTemperatureSetpointShiftComfortErr != nil {
                return nil, errors.New("Error parsing 'RoomTemperatureSetpointShiftComfort' field " + _RoomTemperatureSetpointShiftComfortErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(RoomTemperatureSetpointShiftComfort)

            // Simple Field (RoomTemperatureSetpointShiftStandby)
            RoomTemperatureSetpointShiftStandby, _RoomTemperatureSetpointShiftStandbyErr := io.ReadFloat32(true, 4, 11)
            if _RoomTemperatureSetpointShiftStandbyErr != nil {
                return nil, errors.New("Error parsing 'RoomTemperatureSetpointShiftStandby' field " + _RoomTemperatureSetpointShiftStandbyErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(RoomTemperatureSetpointShiftStandby)

            // Simple Field (RoomTemperatureSetpointShiftEconomy)
            RoomTemperatureSetpointShiftEconomy, _RoomTemperatureSetpointShiftEconomyErr := io.ReadFloat32(true, 4, 11)
            if _RoomTemperatureSetpointShiftEconomyErr != nil {
                return nil, errors.New("Error parsing 'RoomTemperatureSetpointShiftEconomy' field " + _RoomTemperatureSetpointShiftEconomyErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(RoomTemperatureSetpointShiftEconomy)

            // Simple Field (RoomTemperatureSetpointShiftBuildingProtection)
            RoomTemperatureSetpointShiftBuildingProtection, _RoomTemperatureSetpointShiftBuildingProtectionErr := io.ReadFloat32(true, 4, 11)
            if _RoomTemperatureSetpointShiftBuildingProtectionErr != nil {
                return nil, errors.New("Error parsing 'RoomTemperatureSetpointShiftBuildingProtection' field " + _RoomTemperatureSetpointShiftBuildingProtectionErr.Error())
            }
            _map["Struct"] = values.NewPlcREAL(RoomTemperatureSetpointShiftBuildingProtection)
            return values.NewPlcStruct(_map), nil
    }
    return nil, errors.New("unsupported type")
}

func KnxDatapointSerialize(io *utils.WriteBuffer, value api.PlcValue, datapointType IKnxDatapointType) error {
    switch {
        case datapointType == KnxDatapointType.DPT_Switch: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Bool: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Enable: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Ramp: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Alarm: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_BinaryValue: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Step: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_UpDown: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_OpenClose: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Start: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_State: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Invert: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_DimSendStyle: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_InputSource: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Reset: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Ack: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Trigger: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Occupancy: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Window_Door: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_LogicalFunction: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Scene_AB: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_ShutterBlinds_Mode: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_DayNight: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Heat_Cool: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Switch_Control: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Control' field " + _err.Error())
            }

            // Simple Field (On)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'On' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Bool_Control: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Control' field " + _err.Error())
            }

            // Simple Field (True)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'True' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Enable_Control: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Control' field " + _err.Error())
            }

            // Simple Field (Enable)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Enable' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Ramp_Control: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Control' field " + _err.Error())
            }

            // Simple Field (Ramp)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Ramp' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Alarm_Control: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Control' field " + _err.Error())
            }

            // Simple Field (Alarm)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Alarm' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_BinaryValue_Control: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Control' field " + _err.Error())
            }

            // Simple Field (High)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'High' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Step_Control: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Control' field " + _err.Error())
            }

            // Simple Field (Increase)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Increase' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Direction1_Control: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Control' field " + _err.Error())
            }

            // Simple Field (Down)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Down' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Direction2_Control: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Control' field " + _err.Error())
            }

            // Simple Field (Close)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Close' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Start_Control: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Control' field " + _err.Error())
            }

            // Simple Field (Start)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Start' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_State_Control: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Control' field " + _err.Error())
            }

            // Simple Field (Active)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Active' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Invert_Control: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Control)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Control' field " + _err.Error())
            }

            // Simple Field (Inverted)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Inverted' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Control_Dimming: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Increase)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Increase' field " + _err.Error())
            }

            // Simple Field (Stepcode)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Stepcode' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Control_Blinds: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Down)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Down' field " + _err.Error())
            }

            // Simple Field (Stepcode)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Stepcode' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Char_ASCII: // STRING

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x0)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteString(8, "ASCII", value.GetString()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Char_8859_1: // STRING

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x0)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteString(8, "ISO-8859-1", value.GetString()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Scaling: // USINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Angle: // USINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Percent_U8: // USINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_DecimalFactor: // USINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Tariff: // USINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_1_Ucount: // USINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_FanStage: // USINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Percent_V8: // SINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt8(8, value.GetInt8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_1_Count: // SINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt8(8, value.GetInt8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Status_Mode3: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (StatusA)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'StatusA' field " + _err.Error())
            }

            // Simple Field (StatusB)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'StatusB' field " + _err.Error())
            }

            // Simple Field (StatusC)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'StatusC' field " + _err.Error())
            }

            // Simple Field (StatusD)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'StatusD' field " + _err.Error())
            }

            // Simple Field (StatusE)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'StatusE' field " + _err.Error())
            }

            // Simple Field (Mode)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Mode' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_2_Ucount: // UINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_TimePeriodMsec: // UINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_TimePeriod10Msec: // UINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_TimePeriod100Msec: // UINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_TimePeriodSec: // UINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_TimePeriodMin: // UINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_TimePeriodHrs: // UINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_PropDataType: // UINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Length_mm: // UINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_UElCurrentmA: // UINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Brightness: // UINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Absolute_Colour_Temperature: // UINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_2_Count: // INT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_DeltaTimeMsec: // INT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_DeltaTime10Msec: // INT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_DeltaTime100Msec: // INT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_DeltaTimeSec: // INT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_DeltaTimeMin: // INT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_DeltaTimeHrs: // INT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Percent_V16: // INT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Rotation_Angle: // INT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Length_m: // INT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Temp: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Tempd: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Tempa: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Lux: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Wsp: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Pres: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Humidity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_AirQuality: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_AirFlow: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Time1: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Time2: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Volt: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Curr: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_PowerDensity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_KelvinPerPercent: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Power: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Volume_Flow: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Rain_Amount: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Temp_F: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Wsp_kmh: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Absolute_Humidity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Concentration_ygm3: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_TimeOfDay: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Day)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Day' field " + _err.Error())
            }

            // Simple Field (Hour)
            if _err := io.WriteUint8(5, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Hour' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(2, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Minutes)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Minutes' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(2, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Seconds)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Seconds' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Date: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(3, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (DayOfMonth)
            if _err := io.WriteUint8(5, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'DayOfMonth' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Month)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Month' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(1, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Year)
            if _err := io.WriteUint8(7, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Year' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_4_Ucount: // UDINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_LongTimePeriod_Sec: // UDINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_LongTimePeriod_Min: // UDINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_LongTimePeriod_Hrs: // UDINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_VolumeLiquid_Litre: // UDINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Volume_m_3: // UDINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_4_Count: // DINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_FlowRate_m3h: // DINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_ActiveEnergy: // DINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_ApparantEnergy: // DINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_ReactiveEnergy: // DINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_ActiveEnergy_kWh: // DINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_ApparantEnergy_kVAh: // DINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_ReactiveEnergy_kVARh: // DINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_ActiveEnergy_MWh: // DINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_LongDeltaTimeSec: // DINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_DeltaVolumeLiquid_Litre: // DINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_DeltaVolume_m_3: // DINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Acceleration: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Acceleration_Angular: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Activation_Energy: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Activity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Mol: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Amplitude: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_AngleRad: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_AngleDeg: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Angular_Momentum: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Angular_Velocity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Area: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Capacitance: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Charge_DensitySurface: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Charge_DensityVolume: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Compressibility: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Conductance: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Electrical_Conductivity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Density: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Electric_Charge: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Electric_Current: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Electric_CurrentDensity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Electric_DipoleMoment: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Electric_Displacement: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Electric_FieldStrength: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Electric_Flux: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Electric_FluxDensity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Electric_Polarization: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Electric_Potential: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Electric_PotentialDifference: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_ElectromagneticMoment: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Electromotive_Force: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Energy: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Force: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Frequency: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Angular_Frequency: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Heat_Capacity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Heat_FlowRate: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Heat_Quantity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Impedance: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Length: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Light_Quantity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Luminance: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Luminous_Flux: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Luminous_Intensity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Magnetic_FieldStrength: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Magnetic_Flux: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Magnetic_FluxDensity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Magnetic_Moment: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Magnetic_Polarization: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Magnetization: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_MagnetomotiveForce: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Mass: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_MassFlux: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Momentum: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Phase_AngleRad: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Phase_AngleDeg: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Power: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Power_Factor: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Pressure: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Reactance: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Resistance: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Resistivity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_SelfInductance: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_SolidAngle: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Sound_Intensity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Speed: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Stress: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Surface_Tension: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Common_Temperature: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Absolute_Temperature: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_TemperatureDifference: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Thermal_Capacity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Thermal_Conductivity: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_ThermoelectricPower: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Time: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Torque: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Volume: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Volume_Flux: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Weight: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Value_Work: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Volume_Flux_Meter: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Volume_Flux_ls: // REAL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Access_Data: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Value1)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Value1' field " + _err.Error())
            }

            // Simple Field (Value1)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Value1' field " + _err.Error())
            }

            // Simple Field (Value1)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Value1' field " + _err.Error())
            }

            // Simple Field (Value1)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Value1' field " + _err.Error())
            }

            // Simple Field (Value1)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Value1' field " + _err.Error())
            }

            // Simple Field (Value1)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Value1' field " + _err.Error())
            }

            // Simple Field (DetectionError)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'DetectionError' field " + _err.Error())
            }

            // Simple Field (Permission)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Permission' field " + _err.Error())
            }

            // Simple Field (ReadDirection)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ReadDirection' field " + _err.Error())
            }

            // Simple Field (EncryptionOfAccessInformation)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'EncryptionOfAccessInformation' field " + _err.Error())
            }

            // Simple Field (IndexOfAccessIdentificationCode)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'IndexOfAccessIdentificationCode' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_String_ASCII: // STRING

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x0)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteString(112, "ASCII", value.GetString()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_String_8859_1: // STRING

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x0)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteString(112, "ISO-8859-1", value.GetString()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_SceneNumber: // USINT

            // Simple Field (value)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_SceneControl: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (LearnTheSceneCorrespondingToTheFieldSceneNumber)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'LearnTheSceneCorrespondingToTheFieldSceneNumber' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(1, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (SceneNumber)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'SceneNumber' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_DateTime: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Year)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Year' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Month)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Month' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(3, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Dayofmonth)
            if _err := io.WriteUint8(5, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Dayofmonth' field " + _err.Error())
            }

            // Simple Field (Dayofweek)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Dayofweek' field " + _err.Error())
            }

            // Simple Field (Hourofday)
            if _err := io.WriteUint8(5, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Hourofday' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(2, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Minutes)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Minutes' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(2, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Seconds)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Seconds' field " + _err.Error())
            }

            // Simple Field (Fault)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Fault' field " + _err.Error())
            }

            // Simple Field (WorkingDay)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'WorkingDay' field " + _err.Error())
            }

            // Simple Field (NoWd)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'NoWd' field " + _err.Error())
            }

            // Simple Field (NoYear)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'NoYear' field " + _err.Error())
            }

            // Simple Field (NoDate)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'NoDate' field " + _err.Error())
            }

            // Simple Field (NoDayOfWeek)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'NoDayOfWeek' field " + _err.Error())
            }

            // Simple Field (NoTime)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'NoTime' field " + _err.Error())
            }

            // Simple Field (StandardSummerTime)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'StandardSummerTime' field " + _err.Error())
            }

            // Simple Field (QualityOfClock)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'QualityOfClock' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_SCLOMode: // STRING
        case datapointType == KnxDatapointType.DPT_BuildingMode: // STRING
        case datapointType == KnxDatapointType.DPT_OccMode: // STRING
        case datapointType == KnxDatapointType.DPT_Priority: // STRING
        case datapointType == KnxDatapointType.DPT_LightApplicationMode: // STRING
        case datapointType == KnxDatapointType.DPT_ApplicationArea: // STRING
        case datapointType == KnxDatapointType.DPT_AlarmClassType: // STRING
        case datapointType == KnxDatapointType.DPT_PSUMode: // STRING
        case datapointType == KnxDatapointType.DPT_ErrorClass_System: // STRING
        case datapointType == KnxDatapointType.DPT_ErrorClass_HVAC: // STRING
        case datapointType == KnxDatapointType.DPT_Time_Delay: // STRING
        case datapointType == KnxDatapointType.DPT_Beaufort_Wind_Force_Scale: // STRING
        case datapointType == KnxDatapointType.DPT_SensorSelect: // STRING
        case datapointType == KnxDatapointType.DPT_ActuatorConnectType: // STRING
        case datapointType == KnxDatapointType.DPT_Cloud_Cover: // STRING
        case datapointType == KnxDatapointType.DPT_PowerReturnMode: // STRING
        case datapointType == KnxDatapointType.DPT_FuelType: // STRING
        case datapointType == KnxDatapointType.DPT_BurnerType: // STRING
        case datapointType == KnxDatapointType.DPT_HVACMode: // STRING
        case datapointType == KnxDatapointType.DPT_DHWMode: // STRING
        case datapointType == KnxDatapointType.DPT_LoadPriority: // STRING
        case datapointType == KnxDatapointType.DPT_HVACContrMode: // STRING
        case datapointType == KnxDatapointType.DPT_HVACEmergMode: // STRING
        case datapointType == KnxDatapointType.DPT_ChangeoverMode: // STRING
        case datapointType == KnxDatapointType.DPT_ValveMode: // STRING
        case datapointType == KnxDatapointType.DPT_DamperMode: // STRING
        case datapointType == KnxDatapointType.DPT_HeaterMode: // STRING
        case datapointType == KnxDatapointType.DPT_FanMode: // STRING
        case datapointType == KnxDatapointType.DPT_MasterSlaveMode: // STRING
        case datapointType == KnxDatapointType.DPT_StatusRoomSetp: // STRING
        case datapointType == KnxDatapointType.DPT_Metering_DeviceType: // STRING
        case datapointType == KnxDatapointType.DPT_HumDehumMode: // STRING
        case datapointType == KnxDatapointType.DPT_EnableHCStage: // STRING
        case datapointType == KnxDatapointType.DPT_ADAType: // STRING
        case datapointType == KnxDatapointType.DPT_BackupMode: // STRING
        case datapointType == KnxDatapointType.DPT_StartSynchronization: // STRING
        case datapointType == KnxDatapointType.DPT_Behaviour_Lock_Unlock: // STRING
        case datapointType == KnxDatapointType.DPT_Behaviour_Bus_Power_Up_Down: // STRING
        case datapointType == KnxDatapointType.DPT_DALI_Fade_Time: // STRING
        case datapointType == KnxDatapointType.DPT_BlinkingMode: // STRING
        case datapointType == KnxDatapointType.DPT_LightControlMode: // STRING
        case datapointType == KnxDatapointType.DPT_SwitchPBModel: // STRING
        case datapointType == KnxDatapointType.DPT_PBAction: // STRING
        case datapointType == KnxDatapointType.DPT_DimmPBModel: // STRING
        case datapointType == KnxDatapointType.DPT_SwitchOnMode: // STRING
        case datapointType == KnxDatapointType.DPT_LoadTypeSet: // STRING
        case datapointType == KnxDatapointType.DPT_LoadTypeDetected: // STRING
        case datapointType == KnxDatapointType.DPT_Converter_Test_Control: // STRING
        case datapointType == KnxDatapointType.DPT_SABExcept_Behaviour: // STRING
        case datapointType == KnxDatapointType.DPT_SABBehaviour_Lock_Unlock: // STRING
        case datapointType == KnxDatapointType.DPT_SSSBMode: // STRING
        case datapointType == KnxDatapointType.DPT_BlindsControlMode: // STRING
        case datapointType == KnxDatapointType.DPT_CommMode: // STRING
        case datapointType == KnxDatapointType.DPT_AddInfoTypes: // STRING
        case datapointType == KnxDatapointType.DPT_RF_ModeSelect: // STRING
        case datapointType == KnxDatapointType.DPT_RF_FilterSelect: // STRING
        case datapointType == KnxDatapointType.DPT_StatusGen: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(3, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (AlarmStatusOfCorrespondingDatapointIsNotAcknowledged)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'AlarmStatusOfCorrespondingDatapointIsNotAcknowledged' field " + _err.Error())
            }

            // Simple Field (CorrespondingDatapointIsInAlarm)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'CorrespondingDatapointIsInAlarm' field " + _err.Error())
            }

            // Simple Field (CorrespondingDatapointMainValueIsOverridden)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'CorrespondingDatapointMainValueIsOverridden' field " + _err.Error())
            }

            // Simple Field (CorrespondingDatapointMainValueIsCorruptedDueToFailure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'CorrespondingDatapointMainValueIsCorruptedDueToFailure' field " + _err.Error())
            }

            // Simple Field (CorrespondingDatapointValueIsOutOfService)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'CorrespondingDatapointValueIsOutOfService' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Device_Control: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (VerifyModeIsOn)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'VerifyModeIsOn' field " + _err.Error())
            }

            // Simple Field (ADatagramWithTheOwnIndividualAddressAsSourceAddressHasBeenReceived)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ADatagramWithTheOwnIndividualAddressAsSourceAddressHasBeenReceived' field " + _err.Error())
            }

            // Simple Field (TheUserApplicationIsStopped)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'TheUserApplicationIsStopped' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_ForceSign: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Roomhmax)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Roomhmax' field " + _err.Error())
            }

            // Simple Field (Roomhconf)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Roomhconf' field " + _err.Error())
            }

            // Simple Field (Dhwlegio)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Dhwlegio' field " + _err.Error())
            }

            // Simple Field (Dhwnorm)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Dhwnorm' field " + _err.Error())
            }

            // Simple Field (Overrun)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Overrun' field " + _err.Error())
            }

            // Simple Field (Oversupply)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Oversupply' field " + _err.Error())
            }

            // Simple Field (Protection)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Protection' field " + _err.Error())
            }

            // Simple Field (Forcerequest)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Forcerequest' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_ForceSignCool: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_StatusRHC: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Summermode)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Summermode' field " + _err.Error())
            }

            // Simple Field (Statusstopoptim)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Statusstopoptim' field " + _err.Error())
            }

            // Simple Field (Statusstartoptim)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Statusstartoptim' field " + _err.Error())
            }

            // Simple Field (Statusmorningboost)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Statusmorningboost' field " + _err.Error())
            }

            // Simple Field (Tempreturnlimit)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Tempreturnlimit' field " + _err.Error())
            }

            // Simple Field (Tempflowlimit)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Tempflowlimit' field " + _err.Error())
            }

            // Simple Field (Satuseco)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Satuseco' field " + _err.Error())
            }

            // Simple Field (Fault)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Fault' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_StatusSDHWC: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Solarloadsufficient)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Solarloadsufficient' field " + _err.Error())
            }

            // Simple Field (Sdhwloadactive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Sdhwloadactive' field " + _err.Error())
            }

            // Simple Field (Fault)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Fault' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_FuelTypeSet: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Solidstate)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Solidstate' field " + _err.Error())
            }

            // Simple Field (Gas)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Gas' field " + _err.Error())
            }

            // Simple Field (Oil)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Oil' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_StatusRCC: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_StatusAHU: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Cool)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Cool' field " + _err.Error())
            }

            // Simple Field (Heat)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Heat' field " + _err.Error())
            }

            // Simple Field (Fanactive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Fanactive' field " + _err.Error())
            }

            // Simple Field (Fault)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Fault' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_CombinedStatus_RTSM: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(3, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (StatusOfHvacModeUser)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'StatusOfHvacModeUser' field " + _err.Error())
            }

            // Simple Field (StatusOfComfortProlongationUser)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'StatusOfComfortProlongationUser' field " + _err.Error())
            }

            // Simple Field (EffectiveValueOfTheComfortPushButton)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'EffectiveValueOfTheComfortPushButton' field " + _err.Error())
            }

            // Simple Field (EffectiveValueOfThePresenceStatus)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'EffectiveValueOfThePresenceStatus' field " + _err.Error())
            }

            // Simple Field (EffectiveValueOfTheWindowStatus)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'EffectiveValueOfTheWindowStatus' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_LightActuatorErrorInfo: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(1, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Overheat)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Overheat' field " + _err.Error())
            }

            // Simple Field (Lampfailure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Lampfailure' field " + _err.Error())
            }

            // Simple Field (Defectiveload)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Defectiveload' field " + _err.Error())
            }

            // Simple Field (Underload)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Underload' field " + _err.Error())
            }

            // Simple Field (Overcurrent)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Overcurrent' field " + _err.Error())
            }

            // Simple Field (Undervoltage)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Undervoltage' field " + _err.Error())
            }

            // Simple Field (Loaddetectionerror)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Loaddetectionerror' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_RF_ModeInfo: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (BibatSlave)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'BibatSlave' field " + _err.Error())
            }

            // Simple Field (BibatMaster)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'BibatMaster' field " + _err.Error())
            }

            // Simple Field (Asynchronous)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Asynchronous' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_RF_FilterInfo: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Doa)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Doa' field " + _err.Error())
            }

            // Simple Field (KnxSn)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'KnxSn' field " + _err.Error())
            }

            // Simple Field (DoaAndKnxSn)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'DoaAndKnxSn' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Channel_Activation_8: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel1)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel1' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel2)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel2' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel3)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel3' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel4)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel4' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel5)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel5' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel6)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel6' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel7)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel7' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel8)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel8' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_StatusDHWC: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Tempoptimshiftactive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Tempoptimshiftactive' field " + _err.Error())
            }

            // Simple Field (Solarenergysupport)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Solarenergysupport' field " + _err.Error())
            }

            // Simple Field (Solarenergyonly)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Solarenergyonly' field " + _err.Error())
            }

            // Simple Field (Otherenergysourceactive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Otherenergysourceactive' field " + _err.Error())
            }

            // Simple Field (Dhwpushactive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Dhwpushactive' field " + _err.Error())
            }

            // Simple Field (Legioprotactive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Legioprotactive' field " + _err.Error())
            }

            // Simple Field (Dhwloadactive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Dhwloadactive' field " + _err.Error())
            }

            // Simple Field (Fault)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Fault' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_StatusRHCC: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(1, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Overheatalarm)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Overheatalarm' field " + _err.Error())
            }

            // Simple Field (Frostalarm)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Frostalarm' field " + _err.Error())
            }

            // Simple Field (Dewpointstatus)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Dewpointstatus' field " + _err.Error())
            }

            // Simple Field (Coolingdisabled)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Coolingdisabled' field " + _err.Error())
            }

            // Simple Field (Statusprecool)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Statusprecool' field " + _err.Error())
            }

            // Simple Field (Statusecoc)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Statusecoc' field " + _err.Error())
            }

            // Simple Field (Heatcoolmode)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Heatcoolmode' field " + _err.Error())
            }

            // Simple Field (Heatingdiabled)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Heatingdiabled' field " + _err.Error())
            }

            // Simple Field (Statusstopoptim)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Statusstopoptim' field " + _err.Error())
            }

            // Simple Field (Statusstartoptim)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Statusstartoptim' field " + _err.Error())
            }

            // Simple Field (Statusmorningboosth)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Statusmorningboosth' field " + _err.Error())
            }

            // Simple Field (Tempflowreturnlimit)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Tempflowreturnlimit' field " + _err.Error())
            }

            // Simple Field (Tempflowlimit)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Tempflowlimit' field " + _err.Error())
            }

            // Simple Field (Statusecoh)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Statusecoh' field " + _err.Error())
            }

            // Simple Field (Fault)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Fault' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_CombinedStatus_HVA: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (CalibrationMode)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'CalibrationMode' field " + _err.Error())
            }

            // Simple Field (LockedPosition)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'LockedPosition' field " + _err.Error())
            }

            // Simple Field (ForcedPosition)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ForcedPosition' field " + _err.Error())
            }

            // Simple Field (ManuaOperationOverridden)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ManuaOperationOverridden' field " + _err.Error())
            }

            // Simple Field (ServiceMode)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ServiceMode' field " + _err.Error())
            }

            // Simple Field (ValveKick)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ValveKick' field " + _err.Error())
            }

            // Simple Field (Overload)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Overload' field " + _err.Error())
            }

            // Simple Field (ShortCircuit)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ShortCircuit' field " + _err.Error())
            }

            // Simple Field (CurrentValvePosition)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'CurrentValvePosition' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_CombinedStatus_RTC: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (CoolingModeEnabled)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'CoolingModeEnabled' field " + _err.Error())
            }

            // Simple Field (HeatingModeEnabled)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'HeatingModeEnabled' field " + _err.Error())
            }

            // Simple Field (AdditionalHeatingCoolingStage2Stage)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'AdditionalHeatingCoolingStage2Stage' field " + _err.Error())
            }

            // Simple Field (ControllerInactive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ControllerInactive' field " + _err.Error())
            }

            // Simple Field (OverheatAlarm)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'OverheatAlarm' field " + _err.Error())
            }

            // Simple Field (FrostAlarm)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'FrostAlarm' field " + _err.Error())
            }

            // Simple Field (DewPointStatus)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'DewPointStatus' field " + _err.Error())
            }

            // Simple Field (ActiveMode)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActiveMode' field " + _err.Error())
            }

            // Simple Field (GeneralFailureInformation)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'GeneralFailureInformation' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Media: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint16(10, uint16(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (KnxIp)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'KnxIp' field " + _err.Error())
            }

            // Simple Field (Rf)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Rf' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(1, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Pl110)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Pl110' field " + _err.Error())
            }

            // Simple Field (Tp1)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Tp1' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(1, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Channel_Activation_16: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel1)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel1' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel2)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel2' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel3)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel3' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel4)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel4' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel5)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel5' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel6)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel6' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel7)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel7' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel8)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel8' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel9)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel9' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel10)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel10' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel11)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel11' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel12)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel12' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel13)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel13' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel14)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel14' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel15)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel15' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel16)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel16' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_OnOffAction: // STRING
        case datapointType == KnxDatapointType.DPT_Alarm_Reaction: // STRING
        case datapointType == KnxDatapointType.DPT_UpDown_Action: // STRING
        case datapointType == KnxDatapointType.DPT_HVAC_PB_Action: // STRING
        case datapointType == KnxDatapointType.DPT_DoubleNibble: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Busy)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Busy' field " + _err.Error())
            }

            // Simple Field (Nak)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Nak' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_SceneInfo: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(1, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (SceneIsInactive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'SceneIsInactive' field " + _err.Error())
            }

            // Simple Field (Scenenumber)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Scenenumber' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_CombinedInfoOnOff: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput16)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput16' field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput15)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput15' field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput14)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput14' field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput13)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput13' field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput12)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput12' field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput11)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput11' field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput10)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput10' field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput9)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput9' field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput8)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput8' field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput7)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput7' field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput6)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput6' field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput5)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput5' field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput4)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput4' field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput3)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput3' field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput2)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput2' field " + _err.Error())
            }

            // Simple Field (MaskBitInfoOnOffOutput1)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MaskBitInfoOnOffOutput1' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput16)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput16' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput15)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput15' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput14)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput14' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput13)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput13' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput12)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput12' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput11)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput11' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput10)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput10' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput9)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput9' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput8)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput8' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput7)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput7' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput6)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput6' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput5)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput5' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput4)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput4' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput3)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput3' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput2)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput2' field " + _err.Error())
            }

            // Simple Field (InfoOnOffOutput1)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'InfoOnOffOutput1' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_ActiveEnergy_V64: // LINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt64(64, value.GetInt64()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_ApparantEnergy_V64: // LINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt64(64, value.GetInt64()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_ReactiveEnergy_V64: // LINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteInt64(64, value.GetInt64()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Channel_Activation_24: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel1)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel1' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel2)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel2' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel3)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel3' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel4)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel4' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel5)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel5' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel6)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel6' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel7)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel7' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel8)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel8' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel9)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel9' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel10)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel10' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel11)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel11' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel12)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel12' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel13)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel13' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel14)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel14' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel15)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel15' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel16)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel16' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel17)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel17' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel18)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel18' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel19)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel19' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel20)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel20' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel21)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel21' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel22)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel22' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel23)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel23' field " + _err.Error())
            }

            // Simple Field (ActivationStateOfChannel24)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActivationStateOfChannel24' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_HVACModeNext: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (DelayTimeMin)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'DelayTimeMin' field " + _err.Error())
            }

            // Simple Field (HvacMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'HvacMode' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_DHWModeNext: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (DelayTimeMin)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'DelayTimeMin' field " + _err.Error())
            }

            // Simple Field (DhwMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'DhwMode' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_OccModeNext: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (DelayTimeMin)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'DelayTimeMin' field " + _err.Error())
            }

            // Simple Field (OccupancyMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'OccupancyMode' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_BuildingModeNext: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (DelayTimeMin)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'DelayTimeMin' field " + _err.Error())
            }

            // Simple Field (BuildingMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'BuildingMode' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Version: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (MagicNumber)
            if _err := io.WriteUint8(5, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'MagicNumber' field " + _err.Error())
            }

            // Simple Field (VersionNumber)
            if _err := io.WriteUint8(5, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'VersionNumber' field " + _err.Error())
            }

            // Simple Field (RevisionNumber)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'RevisionNumber' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_AlarmInfo: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (LogNumber)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'LogNumber' field " + _err.Error())
            }

            // Simple Field (AlarmPriority)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'AlarmPriority' field " + _err.Error())
            }

            // Simple Field (ApplicationArea)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'ApplicationArea' field " + _err.Error())
            }

            // Simple Field (ErrorClass)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'ErrorClass' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (ErrorcodeSup)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ErrorcodeSup' field " + _err.Error())
            }

            // Simple Field (AlarmtextSup)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'AlarmtextSup' field " + _err.Error())
            }

            // Simple Field (TimestampSup)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'TimestampSup' field " + _err.Error())
            }

            // Simple Field (AckSup)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'AckSup' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Locked)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Locked' field " + _err.Error())
            }

            // Simple Field (Alarmunack)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Alarmunack' field " + _err.Error())
            }

            // Simple Field (Inalarm)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Inalarm' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_TempRoomSetpSetF16_3: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Tempsetpcomf)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'Tempsetpcomf' field " + _err.Error())
            }

            // Simple Field (Tempsetpstdby)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'Tempsetpstdby' field " + _err.Error())
            }

            // Simple Field (Tempsetpeco)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'Tempsetpeco' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_TempRoomSetpSetShiftF16_3: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Tempsetpshiftcomf)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'Tempsetpshiftcomf' field " + _err.Error())
            }

            // Simple Field (Tempsetpshiftstdby)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'Tempsetpshiftstdby' field " + _err.Error())
            }

            // Simple Field (Tempsetpshifteco)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'Tempsetpshifteco' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Scaling_Speed: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (TimePeriod)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'TimePeriod' field " + _err.Error())
            }

            // Simple Field (Percent)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Percent' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Scaling_Step_Time: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (TimePeriod)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'TimePeriod' field " + _err.Error())
            }

            // Simple Field (Percent)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Percent' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_MeteringValue: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Countval)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'Countval' field " + _err.Error())
            }

            // Simple Field (Valinffield)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Valinffield' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(3, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Alarmunack)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Alarmunack' field " + _err.Error())
            }

            // Simple Field (Inalarm)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Inalarm' field " + _err.Error())
            }

            // Simple Field (Overridden)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Overridden' field " + _err.Error())
            }

            // Simple Field (Fault)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Fault' field " + _err.Error())
            }

            // Simple Field (Outofservice)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Outofservice' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_MBus_Address: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Manufactid)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'Manufactid' field " + _err.Error())
            }

            // Simple Field (Identnumber)
            if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
                return errors.New("Error serializing 'Identnumber' field " + _err.Error())
            }

            // Simple Field (Version)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Version' field " + _err.Error())
            }

            // Simple Field (Medium)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Medium' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Colour_RGB: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (R)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'R' field " + _err.Error())
            }

            // Simple Field (G)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'G' field " + _err.Error())
            }

            // Simple Field (B)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'B' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_LanguageCodeAlpha2_ASCII: // STRING

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x0)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteString(16, "ASCII", value.GetString()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Tariff_ActiveEnergy: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Activeelectricalenergy)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'Activeelectricalenergy' field " + _err.Error())
            }

            // Simple Field (Tariff)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Tariff' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Electricalengergyvalidity)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Electricalengergyvalidity' field " + _err.Error())
            }

            // Simple Field (Tariffvalidity)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Tariffvalidity' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Prioritised_Mode_Control: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (DeactivationOfPriority)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'DeactivationOfPriority' field " + _err.Error())
            }

            // Simple Field (PriorityLevel)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'PriorityLevel' field " + _err.Error())
            }

            // Simple Field (ModeLevel)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'ModeLevel' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_DALI_Control_Gear_Diagnostic: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (ConvertorError)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ConvertorError' field " + _err.Error())
            }

            // Simple Field (BallastFailure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'BallastFailure' field " + _err.Error())
            }

            // Simple Field (LampFailure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'LampFailure' field " + _err.Error())
            }

            // Simple Field (ReadOrResponse)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ReadOrResponse' field " + _err.Error())
            }

            // Simple Field (AddressIndicator)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'AddressIndicator' field " + _err.Error())
            }

            // Simple Field (DaliDeviceAddressOrDaliGroupAddress)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'DaliDeviceAddressOrDaliGroupAddress' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_DALI_Diagnostics: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (BallastFailure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'BallastFailure' field " + _err.Error())
            }

            // Simple Field (LampFailure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'LampFailure' field " + _err.Error())
            }

            // Simple Field (DeviceAddress)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'DeviceAddress' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_CombinedPosition: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (HeightPosition)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'HeightPosition' field " + _err.Error())
            }

            // Simple Field (SlatsPosition)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'SlatsPosition' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (ValidityHeightPosition)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ValidityHeightPosition' field " + _err.Error())
            }

            // Simple Field (ValiditySlatsPosition)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ValiditySlatsPosition' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_StatusSAB: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (HeightPosition)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'HeightPosition' field " + _err.Error())
            }

            // Simple Field (SlatsPosition)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'SlatsPosition' field " + _err.Error())
            }

            // Simple Field (UpperEndPosReached)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'UpperEndPosReached' field " + _err.Error())
            }

            // Simple Field (LowerEndPosReached)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'LowerEndPosReached' field " + _err.Error())
            }

            // Simple Field (LowerPredefPosReachedTypHeight100PercentSlatsAngle100Percent)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'LowerPredefPosReachedTypHeight100PercentSlatsAngle100Percent' field " + _err.Error())
            }

            // Simple Field (TargetPosDrive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'TargetPosDrive' field " + _err.Error())
            }

            // Simple Field (RestrictionOfTargetHeightPosPosCanNotBeReached)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'RestrictionOfTargetHeightPosPosCanNotBeReached' field " + _err.Error())
            }

            // Simple Field (RestrictionOfSlatsHeightPosPosCanNotBeReached)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'RestrictionOfSlatsHeightPosPosCanNotBeReached' field " + _err.Error())
            }

            // Simple Field (AtLeastOneOfTheInputsWindRainFrostAlarmIsInAlarm)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'AtLeastOneOfTheInputsWindRainFrostAlarmIsInAlarm' field " + _err.Error())
            }

            // Simple Field (UpDownPositionIsForcedByMoveupdownforcedInput)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'UpDownPositionIsForcedByMoveupdownforcedInput' field " + _err.Error())
            }

            // Simple Field (MovementIsLockedEGByDevicelockedInput)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'MovementIsLockedEGByDevicelockedInput' field " + _err.Error())
            }

            // Simple Field (ActuatorSetvalueIsLocallyOverriddenEGViaALocalUserInterface)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ActuatorSetvalueIsLocallyOverriddenEGViaALocalUserInterface' field " + _err.Error())
            }

            // Simple Field (GeneralFailureOfTheActuatorOrTheDrive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'GeneralFailureOfTheActuatorOrTheDrive' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(3, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (ValidityHeightPos)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ValidityHeightPos' field " + _err.Error())
            }

            // Simple Field (ValiditySlatsPos)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ValiditySlatsPos' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Colour_xyY: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (XAxis)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'XAxis' field " + _err.Error())
            }

            // Simple Field (YAxis)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'YAxis' field " + _err.Error())
            }

            // Simple Field (Brightness)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Brightness' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (ValidityXy)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ValidityXy' field " + _err.Error())
            }

            // Simple Field (ValidityBrightness)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ValidityBrightness' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Converter_Status: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (ConverterModeAccordingToTheDaliConverterStateMachine)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'ConverterModeAccordingToTheDaliConverterStateMachine' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(2, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Hs)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Hs' field " + _err.Error())
            }

            // Simple Field (Hs)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Hs' field " + _err.Error())
            }

            // Simple Field (FunctionTestPending)
            if _err := io.WriteUint8(2, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'FunctionTestPending' field " + _err.Error())
            }

            // Simple Field (DurationTestPending)
            if _err := io.WriteUint8(2, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'DurationTestPending' field " + _err.Error())
            }

            // Simple Field (PartialDurationTestPending)
            if _err := io.WriteUint8(2, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'PartialDurationTestPending' field " + _err.Error())
            }

            // Simple Field (ConverterFailure)
            if _err := io.WriteUint8(2, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'ConverterFailure' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Converter_Test_Result: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Ltrf)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Ltrf' field " + _err.Error())
            }

            // Simple Field (Ltrd)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Ltrd' field " + _err.Error())
            }

            // Simple Field (Ltrp)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Ltrp' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Sf)
            if _err := io.WriteUint8(2, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Sf' field " + _err.Error())
            }

            // Simple Field (Sd)
            if _err := io.WriteUint8(2, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Sd' field " + _err.Error())
            }

            // Simple Field (Sp)
            if _err := io.WriteUint8(2, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Sp' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(2, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Ldtr)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'Ldtr' field " + _err.Error())
            }

            // Simple Field (Lpdtr)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Lpdtr' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Battery_Info: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (BatteryFailure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'BatteryFailure' field " + _err.Error())
            }

            // Simple Field (BatteryDurationFailure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'BatteryDurationFailure' field " + _err.Error())
            }

            // Simple Field (BatteryFullyCharged)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'BatteryFullyCharged' field " + _err.Error())
            }

            // Simple Field (BatteryChargeLevel)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'BatteryChargeLevel' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Brightness_Colour_Temperature_Transition: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Ms)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'Ms' field " + _err.Error())
            }

            // Simple Field (TemperatureK)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'TemperatureK' field " + _err.Error())
            }

            // Simple Field (Percent)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'Percent' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (ValidityOfTheTimePeriod)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ValidityOfTheTimePeriod' field " + _err.Error())
            }

            // Simple Field (ValidityOfTheAbsoluteColourTemperature)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ValidityOfTheAbsoluteColourTemperature' field " + _err.Error())
            }

            // Simple Field (ValidityOfTheAbsoluteBrightness)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ValidityOfTheAbsoluteBrightness' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Brightness_Colour_Temperature_Control: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Cct)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Cct' field " + _err.Error())
            }

            // Simple Field (StepCodeColourTemperature)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'StepCodeColourTemperature' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Cb)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Cb' field " + _err.Error())
            }

            // Simple Field (StepCodeBrightness)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'StepCodeBrightness' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (CctAndStepCodeColourValidity)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'CctAndStepCodeColourValidity' field " + _err.Error())
            }

            // Simple Field (CbAndStepCodeBrightnessValidity)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'CbAndStepCodeBrightnessValidity' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Colour_RGBW: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (ColourLevelRed)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'ColourLevelRed' field " + _err.Error())
            }

            // Simple Field (ColourLevelGreen)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'ColourLevelGreen' field " + _err.Error())
            }

            // Simple Field (ColourLevelBlue)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'ColourLevelBlue' field " + _err.Error())
            }

            // Simple Field (ColourLevelWhite)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'ColourLevelWhite' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Mr)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Mr' field " + _err.Error())
            }

            // Simple Field (Mg)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Mg' field " + _err.Error())
            }

            // Simple Field (Mb)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Mb' field " + _err.Error())
            }

            // Simple Field (Mw)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Mw' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Relative_Control_RGBW: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Maskcw)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Maskcw' field " + _err.Error())
            }

            // Simple Field (Maskcb)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Maskcb' field " + _err.Error())
            }

            // Simple Field (Maskcg)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Maskcg' field " + _err.Error())
            }

            // Simple Field (Maskcr)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Maskcr' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Cw)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Cw' field " + _err.Error())
            }

            // Simple Field (StepCodeColourWhite)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'StepCodeColourWhite' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Cb)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Cb' field " + _err.Error())
            }

            // Simple Field (StepCodeColourBlue)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'StepCodeColourBlue' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Cg)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Cg' field " + _err.Error())
            }

            // Simple Field (StepCodeColourGreen)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'StepCodeColourGreen' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Cr)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Cr' field " + _err.Error())
            }

            // Simple Field (StepCodeColourRed)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'StepCodeColourRed' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_Relative_Control_RGB: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Cb)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Cb' field " + _err.Error())
            }

            // Simple Field (StepCodeColourBlue)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'StepCodeColourBlue' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Cg)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Cg' field " + _err.Error())
            }

            // Simple Field (StepCodeColourGreen)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'StepCodeColourGreen' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Cr)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'Cr' field " + _err.Error())
            }

            // Simple Field (StepCodeColourRed)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'StepCodeColourRed' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_GeographicalLocation: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (Longitude)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'Longitude' field " + _err.Error())
            }

            // Simple Field (Latitude)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'Latitude' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_TempRoomSetpSetF16_4: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (RoomTemperatureSetpointComfort)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'RoomTemperatureSetpointComfort' field " + _err.Error())
            }

            // Simple Field (RoomTemperatureSetpointStandby)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'RoomTemperatureSetpointStandby' field " + _err.Error())
            }

            // Simple Field (RoomTemperatureSetpointEconomy)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'RoomTemperatureSetpointEconomy' field " + _err.Error())
            }

            // Simple Field (RoomTemperatureSetpointBuildingProtection)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'RoomTemperatureSetpointBuildingProtection' field " + _err.Error())
            }
        case datapointType == KnxDatapointType.DPT_TempRoomSetpSetShiftF16_4: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (RoomTemperatureSetpointShiftComfort)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'RoomTemperatureSetpointShiftComfort' field " + _err.Error())
            }

            // Simple Field (RoomTemperatureSetpointShiftStandby)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'RoomTemperatureSetpointShiftStandby' field " + _err.Error())
            }

            // Simple Field (RoomTemperatureSetpointShiftEconomy)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'RoomTemperatureSetpointShiftEconomy' field " + _err.Error())
            }

            // Simple Field (RoomTemperatureSetpointShiftBuildingProtection)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'RoomTemperatureSetpointShiftBuildingProtection' field " + _err.Error())
            }
        default:

            return errors.New("unsupported type")
    }
    return nil
}

