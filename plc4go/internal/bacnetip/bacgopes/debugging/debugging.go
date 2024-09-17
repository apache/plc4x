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
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"time"
	"unsafe"
)

var _isDebuggingActive = false
var _debug = CreateDebugPrinter()

func Btox(data []byte, sep string) string {
	hexString := hex.EncodeToString(data)
	if sep != "" {
		pairs := make([]string, len(hexString)/2)
		for i := 0; i < len(hexString)-1; i += 2 {
			pairs[i/2] = hexString[i : i+2]
		}
		hexString = strings.Join(pairs, ".")
	}
	return hexString
}

func Xtob(hexString string) ([]byte, error) {
	compile, err := regexp.Compile("[^0-9a-fA-F]")
	if err != nil {
		return nil, err
	}
	replaceAll := compile.ReplaceAll([]byte(hexString), nil)
	decodeString, err := hex.DecodeString(string(replaceAll))
	if err != nil {
		return nil, err
	}
	return decodeString, nil
}

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

func (d *DebugContents) Format(s fmt.State, v rune) {
	if d == nil {
		if _debug != nil {
			_debug("nil debug attempt") // TODO: this should no happen
		}
		return
	}
	switch v {
	case 's', 'v':
		// TODO: check if that hacky hacky makes sense
		if len(d.debuggables) > 0 {
			debuggable := d.debuggables[0]
			_, _ = fmt.Fprintf(s, "<%s at %p>", QualifiedTypeName(debuggable), debuggable)
		}
	case 'r':
		// TODO: check if that hacky hacky makes sense
		if len(d.debuggables) > 0 {
			debuggable := d.debuggables[0]
			_, _ = fmt.Fprintf(s, "<%s at %p>\n", QualifiedTypeName(debuggable), debuggable)
			_, _ = fmt.Fprintf(s, "    <%s at %p>\n", QualifiedTypeName(debuggable), debuggable) // TODO: why is this duplicated?
		}
		d.PrintDebugContents(2, s, nil)
	}
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
				keyPrintVerb := verbForType(key, 'r')
				elemPrintVerb := verbForType(elem, 'r')
				_, _ = fmt.Fprintf(file, "%s"+string(keyPrintVerb)+" : "+string(elemPrintVerb)+"\n", strings.Repeat("    ", indent), key, elem)
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
			printVerb := verbForType(value, 'r')
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
	typeNameString := "bacgopes." + typeOf.String()
	if customProjectName != "" {
		typeNameString = strings.ReplaceAll(typeNameString, projectName, customProjectName)
	}
	return typeNameString
}

// TODO: migrate comp debug logger to here...
type LoggingFormatter struct {
	// TODO: implement me
}

type DefaultRFormatter struct {
	header        string
	extraPrinters []DebugContentPrinter
}

func NewDefaultRFormatter(extraPrinters ...DebugContentPrinter) *DefaultRFormatter {
	pc, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("oh no")
	}
	dir := path.Dir(file)
	rootIndex := strings.Index(dir, "bacgopes")
	dir = dir[rootIndex:]
	dirPrefix := path.Base(dir) + "_"
	base := path.Base(file)
	prefix := strings.TrimSuffix(base, ".go")
	prefix = strings.TrimPrefix(prefix, dirPrefix)
	qualifier := strings.ReplaceAll(dirPrefix, "_", ".")
	header := fmt.Sprintf("<%s at 0x%x>", "bacgopes."+qualifier+prefix, pc)
	if customProjectName != "" {
		header = strings.ReplaceAll(header, projectName, customProjectName)
	}
	return &DefaultRFormatter{
		header:        header,
		extraPrinters: extraPrinters,
	}
}

func (d *DefaultRFormatter) Format(s fmt.State, v rune) {
	if d.header == "" && len(d.extraPrinters) == 0 {
		panic("misconfiguration")
	}
	switch v {
	case 'r':
		_, _ = s.Write([]byte(d.header))
		if len(d.extraPrinters) > 0 {
			_, _ = s.Write([]byte("\n"))
		}
		for _, printer := range d.extraPrinters {
			printer.PrintDebugContents(1, s, nil)
		}
	case 'v', 's':
		_, _ = s.Write([]byte(d.header))
	}
}

func StructName() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	dir := path.Dir(file)
	rootIndex := strings.Index(dir, "bacgopes")
	dir = dir[rootIndex:]
	dirPrefix := path.Base(dir) + "_"
	base := path.Base(file)
	prefix := strings.TrimSuffix(base, ".go")
	return strings.TrimPrefix(prefix, dirPrefix)
}

type DebugPrinter = func(format string, a ...any)

const projectName = "bacgopes"

var customProjectName = os.Getenv("BACGOPES_DEBUG_CUSTOM_PROJECT_NAME")
var customReplaces map[string]string

func init() {
	customReplaces = map[string]string{}
	for _, replace := range strings.Split(os.Getenv("BACGOPES_DEBUG_CUSTOM_REPLACES"), ",") {
		if replace == "" || !strings.Contains(replace, "=") {
			continue
		}
		kv := strings.SplitN(replace, "=", 2)
		if len(kv) != 2 {
			panic("invalid replace " + replace)
		}
		customReplaces[kv[0]] = kv[1]
	}
}

func CreateDebugPrinter() DebugPrinter {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return nil
	}
	dir := path.Dir(file)
	rootIndex := strings.Index(dir, "bacgopes")
	dir = dir[rootIndex:]
	qualifier := strings.ReplaceAll(dir, "/", ".")
	dirPrefix := path.Base(dir) + "_"

	bacgopesDebug := os.Getenv("BACGOPES_DEBUG")
	if strings.Contains(bacgopesDebug, qualifier) {
		_isDebuggingActive = true
		return func(format string, a ...any) {
			pc, file, _, ok := runtime.Caller(1)
			if !ok {
				return
			}
			base := path.Base(file)
			prefix := strings.TrimSuffix(base, ".go")
			if !strings.HasPrefix(prefix, dirPrefix) && !strings.Contains(prefix, "tests") && false { // TODO: disabled for now as it makes more trouble for the rest
				// Attach the fuction name // TODO: check if that makes sense, only a workaround for bind at the moment
				details := runtime.FuncForPC(pc)
				name := details.Name()
				name = name[strings.LastIndex(name, ".")+1:]
				prefix = strings.ToLower(name)
			}
			prefix = strings.TrimPrefix(prefix, dirPrefix)
			formatString := "DEBUG:" + qualifier + "." + prefix + ":" + format + "\n"
			formatString = cleanupFormatString(formatString)
			if customProjectName != "" {
				formatString = strings.ReplaceAll(formatString, projectName, customProjectName)
			}
			formatString = fixVerbs(formatString, a...)
			output := fmt.Sprintf(formatString, a...)
			if strings.HasSuffix(output, "\n\n") { // TODO: another hacky workaround
				output = strings.TrimSuffix(output, "\n")
			}
			for k, v := range customReplaces {
				output = strings.ReplaceAll(output, k, v)
			}
			_, _ = os.Stdout.Write([]byte(output))
		}
	}
	return nil
}

func IsDebuggingActive() bool {
	return _isDebuggingActive
}

func cleanupFormatString(s string) string {
	// TODO: investigate via comm.comm is happening
	s = strings.ReplaceAll(s, ".comm.comm:", ".comm:")
	s = strings.ReplaceAll(s, "pdu.comm_PCI:", "comm.PCI:")
	s = strings.ReplaceAll(s, "pdu.comm_PDUData:", "comm.PDUData:")
	s = strings.ReplaceAll(s, "DEBUG:"+projectName+".tests", "DEBUG:tests")
	return s
}

func fixVerbs(formatString string, values ...any) string {
	length := len(formatString)
	verbNumber := -1
	for i, r := range formatString {
		switch r {
		case '%':
			nextIndex := i + 1
			if nextIndex >= length {
				continue
			}
			followRune := formatString[nextIndex]
			if followRune != '%' {
				verbNumber++
			}
			if followRune == 'r' && verbNumber < len(values) { // TODO: this completely breaks at indexed verbs... better fix assap
				runes := []rune(formatString)
				runes[nextIndex] = verbForType(values[verbNumber], 'r')
				formatString = string(runes)
			}
		}
	}
	return formatString
}

func verbForType(value any, printVerb rune) rune {
	switch value.(type) {
	case string:
		printVerb = 's'
	case bool:
		printVerb = 't'
	case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, uintptr:
		printVerb = 'd'
	case float32, float64:
		printVerb = 'f'
	case complex64, complex128:
		printVerb = 'v' // TODO: what is it meant for?
	case time.Time, time.Duration:
		printVerb = 's'
	case []byte:
		printVerb = 'v'
	}
	return printVerb
}

// clone from comp to avoid circular dependencies // TODO: maybe move Btox somewhere else or come up with something smarter there
func isNil(v interface{}) bool {
	if v == nil {
		return true
	}
	valueOf := reflect.ValueOf(v)
	switch valueOf.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		return valueOf.IsNil()
	default:
		return false
	}
}
