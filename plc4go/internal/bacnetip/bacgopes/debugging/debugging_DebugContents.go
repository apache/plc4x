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

package debugging

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strings"
	"unsafe"
)

type Debuggable interface {
	GetDebugAttr(attr string) any
}

type DebugContentPrinter interface {
	PrintDebugContents(indent int, file io.Writer, _ids []uintptr)
}

type DebugContents struct {
	debuggables        []Debuggable
	debuggableContents map[Debuggable][]string
	extraPrinters      []DebugContentPrinter
}

func NewDebugContents(debuggable Debuggable, contents ...string) *DebugContents {
	return &DebugContents{
		debuggables: []Debuggable{debuggable},
		debuggableContents: map[Debuggable][]string{
			debuggable: contents,
		},
	}
}

var _ DebugContentPrinter = (*DebugContents)(nil)

func (d *DebugContents) AddExtraPrinters(printers ...DebugContentPrinter) {
	d.extraPrinters = append(d.extraPrinters, printers...)
}

func (d *DebugContents) AddDebugContents(debuggable Debuggable, contents ...string) {
	d.debuggables = append(d.debuggables, debuggable)
	d.debuggableContents[debuggable] = contents
}

func (d *DebugContents) StructHeader() []byte {
	return []byte(fmt.Sprintf("<%s object at %p>", d.LeafNameOrFallback(), d))
}

func (d *DebugContents) Format(s fmt.State, v rune) {
	if d == nil {
		if _debug != nil {
			_debug("nil debug attempt") // TODO: this should no happen
		}
		return
	}
	switch v {
	case 's', 'v':
		_, _ = s.Write(d.StructHeader())
	case 'r':
		_, _ = s.Write(d.StructHeader())
		_, _ = s.Write([]byte{'\n'})
		_, _ = s.Write(append([]byte("    "), d.StructHeader()...)) // TODO: why is this duplicated?
		_, _ = s.Write([]byte{'\n'})
		d.PrintDebugContents(2, s, nil)
	}
}

func (d *DebugContents) LeafNameOrFallback() string {
	if len(d.debuggables) > 0 {
		debuggable := d.debuggables[0]
		if leafNameSupplier, ok := debuggable.(interface{ GetLeafName() string }); ok {
			return leafNameSupplier.GetLeafName()
		}
		return QualifiedTypeName(debuggable)
	}
	return StructName()
}

func (d *DebugContents) PrintDebugContents(indent int, file io.Writer, _ids []uintptr) {
	if indent == 0 {
		indent = 1
	}
	if file == nil {
		file = os.Stderr
	}

	// loop through the debuggables and look for debugContents
	attrToDebuggable := map[string]Debuggable{} // Reverse map to know where to get the attr from
	var attrs []string
	cids := map[uintptr]struct{}{
		uintptr(unsafe.Pointer(d)): {},
	}
	var ownFn []DebugContentPrinter // TODO: make them member and add a add method
	for _, debuggable := range d.debuggables {
		debugContents, ok := d.debuggableContents[debuggable]
		if !ok {
			continue
		}

		// already seen it?
		if _, ok := cids[uintptr(unsafe.Pointer(&debuggable))]; ok {
			// continue // TODO: seems the logic is broken: e.g. RecurringTask is seen but it should be a own entity
		}
		cids[uintptr(unsafe.Pointer(&debuggable))] = struct{}{}

		for _, attr := range debugContents {
			if !slices.Contains(attrs, attr) {
				cleanAttr := attr
				for strings.HasSuffix(cleanAttr, "-") || strings.HasSuffix(cleanAttr, "*") || strings.HasSuffix(cleanAttr, "+") {
					cleanAttr = strings.TrimSuffix(cleanAttr, "-")
					cleanAttr = strings.TrimSuffix(cleanAttr, "*")
					cleanAttr = strings.TrimSuffix(cleanAttr, "+")
				}
				attrToDebuggable[cleanAttr] = debuggable
				attrs = append(attrs, attr)
			}
		}
	}

	// a bit of debugging
	if _debug != nil {
		_debug("    - attrs: %r", attrs)
		//_debug("    - ownFn: %r", ownFn)
	}

	for _, printer := range d.extraPrinters {
		printer.PrintDebugContents(indent, file, _ids)
	}
	for _, attr := range attrs {
		// assume you're going deep, but not into lists and dicts
		goDeep := true
		goListDict := false
		goHexed := false

		// attribute list might want to go deep
		if strings.HasSuffix(attr, "-") {
			goDeep = false
			attr = attr[:len(attr)-1]
		} else if strings.HasSuffix(attr, "*") {
			goHexed = true
			attr = attr[:len(attr)-1]
		} else if strings.HasSuffix(attr, "+") {
			goDeep = false
			goListDict = true
			attr = attr[:len(attr)-1]
			if strings.HasSuffix(attr, "+") {
				goDeep = true
				attr = attr[:len(attr)-1]
			}
		}

		debuggable := attrToDebuggable[attr]
		if debuggable == nil {
			panic("attr misconfiguration " + attr)
		}
		value := debuggable.GetDebugAttr(attr)

		// skip nil
		if isNil(value) {
			continue
		}

		// standard output
		if goListDict && isList(value) {
			list, ok := toList(value)
			if !ok {
				panic("impossible")
			}
			_, _ = fmt.Fprintf(file, "%s%s = [\n", strings.Repeat("    ", indent), attr)
			indent += 1
			for i, elem := range list {
				_, _ = fmt.Fprintf(file, "%s[%d] %r\n", strings.Repeat("    ", indent), i, elem)
				if deepDebugContent, ok := elem.(interface {
					DebugContents(int, io.Writer, []uintptr)
				}); goDeep && ok {
					if slices.Contains(_ids, uintptr(unsafe.Pointer(&deepDebugContent))) {
						_ids = append(_ids, uintptr(unsafe.Pointer(&deepDebugContent)))
						deepDebugContent.DebugContents(indent+1, file, _ids)
					}
				}
			}
			indent -= 1
			_, _ = fmt.Fprintf(file, "%s    ]\n", strings.Repeat("    ", indent))
		} else if goListDict && isDict(value) {
			_map, ok := toMap(value)
			if !ok {
				panic("impossible")
			}
			_, _ = fmt.Fprintf(file, "%s%s = {\n", strings.Repeat("    ", indent), attr)
			indent += 1
			for key, elem := range _map {
				keyPrintVerb := VerbForType(key, 'r')
				elemPrintVerb := VerbForType(elem, 'r')
				_, _ = fmt.Fprintf(file, "%s%"+string(keyPrintVerb)+" : %"+string(elemPrintVerb)+"\n", strings.Repeat("    ", indent), key, elem)
				if deepDebugContent, ok := elem.(interface {
					DebugContents(int, io.Writer, []uintptr)
				}); goDeep && ok {
					if slices.Contains(_ids, uintptr(unsafe.Pointer(&deepDebugContent))) {
						_ids = append(_ids, uintptr(unsafe.Pointer(&deepDebugContent)))
						deepDebugContent.DebugContents(indent+1, file, _ids)
					}
				}
			}
			indent -= 1
			_, _ = fmt.Fprintf(file, "%s    }\n", strings.Repeat("    ", indent))
		} else if goHexed && isString(value) { // TODO: add support
			panic("add support")
		} else {
			printVerb := VerbForType(value, 'r')
			_, _ = fmt.Fprintf(file, "%s%s = %"+string(printVerb)+"\n", strings.Repeat("    ", indent), attr, value)

			// go nested if it is debuggable
			if deepDebugContent, ok := value.(DebugContentPrinter); goDeep && ok {
				if slices.Contains(_ids, uintptr(unsafe.Pointer(&deepDebugContent))) {
					_ids = append(_ids, uintptr(unsafe.Pointer(&deepDebugContent)))
					deepDebugContent.PrintDebugContents(indent+1, file, _ids)
				}
			}
		}
	}

	// go through the functions
	slices.Reverse(ownFn)
	for _, printer := range ownFn {
		printer.PrintDebugContents(indent, file, _ids)
	}
}

func isList(input any) bool {
	reflectValue := reflect.ValueOf(input)
	switch reflectValue.Kind() {
	case reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}

func toList(input any) (list []any, ok bool) {
	reflectValue := reflect.ValueOf(input)
	switch reflectValue.Kind() {
	case reflect.Slice, reflect.Array:
	default:
		return nil, false
	}
	list = make([]any, reflectValue.Len())
	for i := 0; i < reflectValue.Len(); i++ {
		elem := reflectValue.Index(i).Interface()
		list[i] = elem
	}
	return list, true
}

func isDict(input any) bool {
	reflectValue := reflect.ValueOf(input)
	switch reflectValue.Kind() {
	case reflect.Map:
		return true
	default:
		return false
	}
}

func toMap(input any) (_map map[string]any, ok bool) {
	reflectValue := reflect.ValueOf(input)
	switch reflectValue.Kind() {
	case reflect.Map:
	default:
		return nil, false
	}
	_map = make(map[string]any)
	for _, k := range reflectValue.MapKeys() {
		value := reflectValue.MapIndex(k)
		_map[k.String()] = value.Interface() //TODO: check if we further need to convert map key
	}
	return _map, true
}

func isString(input any) bool {
	_, ok := input.(string)
	return ok
}

func TypeName(anything any) string {
	typeOf := reflect.TypeOf(anything)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	return typeOf.Name()
}

func QualifiedTypeName(anything any) string {
	typeOf := reflect.TypeOf(anything)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	typeNameString := projectName + "." + typeOf.String()
	if customProjectName != "" {
		typeNameString = strings.ReplaceAll(typeNameString, projectName, customProjectName)
	}
	return typeNameString
}
