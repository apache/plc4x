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

package umas

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// browse runs one query and insists on a result: a browse which never delivers would hang its caller.
func browse(t *testing.T, connection *Connection, queries map[string]string) apiModel.PlcBrowseResponse {
	t.Helper()
	builder := connection.BrowseRequestBuilder()
	for queryName, query := range queries {
		builder = builder.AddQuery(queryName, query)
	}
	request, err := builder.Build()
	require.NoError(t, err)
	select {
	case result := <-request.Execute(testutils.TestContext(t)):
		require.NotNil(t, result)
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		return result.GetResponse()
	case <-time.After(5 * time.Second):
		t.Fatal("the browse never delivered a result")
		return nil
	}
}

// dataTypeName is the type name of a browse item. It is not on the apiModel.PlcBrowseItem interface,
// only on the SPI's implementation, so it has to be asked for through that.
func dataTypeName(t *testing.T, item apiModel.PlcBrowseItem) string {
	t.Helper()
	defaultItem, ok := item.(*spiModel.DefaultPlcBrowseItem)
	require.True(t, ok, "%T is not a DefaultPlcBrowseItem", item)
	return defaultItem.GetDataTypeName()
}

// itemsByName indexes a browse result so the assertions don't depend on the order.
func itemsByName(items []apiModel.PlcBrowseItem) map[string]apiModel.PlcBrowseItem {
	byName := map[string]apiModel.PlcBrowseItem{}
	for _, item := range items {
		byName[item.GetName()] = item
	}
	return byName
}

func TestBrowser_ListsEverySymbol(t *testing.T) {
	connection, _ := newConnectedConnection(t)

	response := browse(t, connection, map[string]string{"all": "*"})
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("all"))

	items := itemsByName(response.GetQueryResults("all"))
	require.Len(t, items, len(defaultFixture().symbols))

	real := items["g_r32"]
	require.NotNil(t, real)
	assert.Equal(t, apiValues.REAL, real.GetTag().GetValueType())
	assert.Equal(t, "REAL", dataTypeName(t, real))
	assert.True(t, real.IsReadable())
	assert.True(t, real.IsWritable())
	// Subscribing works because it is emulated by polling the read path.
	assert.True(t, real.IsSubscribable())
	// The tag a browse hands out has to be usable as a read address.
	_, err := NewTagHandler().ParseTag(real.GetTag().GetAddressString())
	assert.NoError(t, err)

	assert.Equal(t, apiValues.STRING, items["g_string"].GetTag().GetValueType())
	assert.Equal(t, apiValues.INT, items["g_b16"].GetTag().GetValueType())
}

// A struct symbol reports its members as children, with the name the project gave the type.
func TestBrowser_DescribesAStructSymbol(t *testing.T) {
	connection, _ := newConnectedConnection(t)

	items := itemsByName(browse(t, connection, map[string]string{"all": "*"}).GetQueryResults("all"))
	plant := items["g_plant"]
	require.NotNil(t, plant)
	assert.Equal(t, apiValues.Struct, plant.GetTag().GetValueType())
	assert.Equal(t, "MY_STRUCT", dataTypeName(t, plant))

	children := plant.GetChildren()
	require.Len(t, children, 2)
	require.Contains(t, children, "meta")
	assert.Equal(t, apiValues.DINT, children["meta"].GetTag().GetValueType())
	require.Contains(t, children, "r32")
	assert.Equal(t, apiValues.REAL, children["r32"].GetTag().GetValueType())
}

// An array symbol reports the type of its elements plus its bounds, so a caller sees "array of DINT"
// rather than "unknown".
func TestBrowser_DescribesAnArraySymbol(t *testing.T) {
	connection, _ := newConnectedConnection(t)

	items := itemsByName(browse(t, connection, map[string]string{"all": "*"}).GetQueryResults("all"))
	array := items["g_arrInt"]
	require.NotNil(t, array)
	assert.Equal(t, apiValues.DINT, array.GetTag().GetValueType(), "the element type")
	assert.Equal(t, "MY_ARRAY", dataTypeName(t, array))

	arrayInfo := array.GetTag().GetArrayInfo()
	require.Len(t, arrayInfo, 1)
	assert.Equal(t, uint32(0), arrayInfo[0].GetLowerBound())
	assert.Equal(t, uint32(9), arrayInfo[0].GetUpperBound())
}

// A query filters by name, which is where this driver goes beyond plc4j (whose browse ignores the
// query and answers with the whole table every time).
func TestBrowser_FiltersByQuery(t *testing.T) {
	connection, _ := newConnectedConnection(t)

	response := browse(t, connection, map[string]string{
		"strings": "g_str*",
		"none":    "h_*",
		"exact":   "g_r32",
	})

	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("strings"))
	strings := response.GetQueryResults("strings")
	require.Len(t, strings, 1)
	assert.Equal(t, "g_string", strings[0].GetName())

	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("none"))
	assert.Empty(t, response.GetQueryResults("none"))

	exact := response.GetQueryResults("exact")
	require.Len(t, exact, 1)
	assert.Equal(t, "g_r32", exact[0].GetName())
}

func TestBrowser_HonorsAnInterceptor(t *testing.T) {
	connection, _ := newConnectedConnection(t)

	request, err := connection.BrowseRequestBuilder().AddQuery("all", "*").Build()
	require.NoError(t, err)

	var offered []string
	results := request.ExecuteWithInterceptor(testutils.TestContext(t), func(item apiModel.PlcBrowseItem) bool {
		offered = append(offered, item.GetName())
		return item.GetName() == "g_r32"
	})
	select {
	case result := <-results:
		require.NoError(t, result.GetErr())
		items := result.GetResponse().GetQueryResults("all")
		require.Len(t, items, 1)
		assert.Equal(t, "g_r32", items[0].GetName())
		assert.Len(t, offered, len(defaultFixture().symbols), "the interceptor sees every symbol")
	case <-time.After(5 * time.Second):
		t.Fatal("the browse never delivered a result")
	}
}

// A browse on a connection whose dictionary download failed retries it, which is what plc4j's browse
// does too.
func TestBrowser_RetriesTheDictionaryDownload(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	connectResult := connect(t, connection)
	runHandshakeWithoutDictionary(t, codec)
	// The download during connect is refused.
	codec.answerWith(t, codec.nextRequest(t), readWriteModel.NewUmasPDUErrorResponse(0, []byte{0x01}))
	requireConnected(t, connectResult)
	require.False(t, connection.session.hasSymbols())

	browseResult := make(chan apiModel.PlcBrowseResponse, 1)
	go func() { browseResult <- browse(t, connection, map[string]string{"all": "*"}) }()

	// This time the PLC answers, and the browse has the dictionary it needs.
	runDictionaryDownload(t, codec, defaultFixture())

	select {
	case response := <-browseResult:
		assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("all"))
		assert.Len(t, response.GetQueryResults("all"), len(defaultFixture().symbols))
	case <-time.After(10 * time.Second):
		t.Fatal("the browse never finished")
	}
	assert.True(t, connection.session.hasSymbols())
}

// A browse which still can't get the dictionary reports a remote error rather than an empty success.
func TestBrowser_WithoutADictionaryReportsAnError(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	connectResult := connect(t, connection)
	runHandshakeWithoutDictionary(t, codec)
	codec.answerWith(t, codec.nextRequest(t), readWriteModel.NewUmasPDUErrorResponse(0, []byte{0x01}))
	requireConnected(t, connectResult)

	browseResult := make(chan apiModel.PlcBrowseResponse, 1)
	go func() { browseResult <- browse(t, connection, map[string]string{"all": "*"}) }()
	codec.answerWith(t, codec.nextRequest(t), readWriteModel.NewUmasPDUErrorResponse(0, []byte{0x01}))

	select {
	case response := <-browseResult:
		assert.Equal(t, apiModel.PlcResponseCode_REMOTE_ERROR, response.GetResponseCode("all"))
		assert.Empty(t, response.GetQueryResults("all"))
	case <-time.After(10 * time.Second):
		t.Fatal("the browse never finished")
	}
}

// A dictionary that describes a type which contains itself must not send the browser into an endless
// recursion. Nothing stops a device from reporting one.
func TestBrowser_SurvivesASelfReferentialType(t *testing.T) {
	connection, _ := newConnectedConnection(t)

	// A struct type whose only member is of the type itself.
	connection.session.setStructType(customTypeIdBase+9, "RECURSIVE", []readWriteModel.UmasUDTDefinition{
		readWriteModel.NewUmasUDTDefinition(customTypeIdBase+9, 0, 0, 0, "self"),
	})
	connection.session.setSymbols([]readWriteModel.UmasUnlocatedVariableReference{
		readWriteModel.NewUmasUnlocatedVariableReference(customTypeIdBase+9, 2, 0, 0, 0, "g_recursive"),
	})

	items := browse(t, connection, map[string]string{"all": "*"}).GetQueryResults("all")
	require.Len(t, items, 1)
	assert.Equal(t, apiValues.Struct, items[0].GetTag().GetValueType())
	// The nesting stops somewhere rather than never.
	depth := 0
	item := items[0]
	for len(item.GetChildren()) > 0 && depth < maxTypeNestingDepth+5 {
		depth++
		for _, child := range item.GetChildren() {
			item = child
			break
		}
	}
	assert.LessOrEqual(t, depth, maxTypeNestingDepth)
}

func TestBrowser_RefusesAQueryOfTheWrongKind(t *testing.T) {
	connection, _ := newConnectedConnection(t)
	code, items := connection.browseQuery(testutils.TestContext(t), nil, foreignQuery{})
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_ADDRESS, code)
	assert.Empty(t, items)
}

// foreignQuery is a query object from some other driver, which this browser can make nothing of.
type foreignQuery struct{}

func (foreignQuery) GetQueryString() string { return "not ours" }

func TestDataTypeName(t *testing.T) {
	connection, _ := newConnectedConnection(t)
	assert.Equal(t, "REAL", connection.dataTypeName(typeIdReal))
	assert.Equal(t, "MY_STRUCT", connection.dataTypeName(customTypeIdBase+1))
	// A type id nothing knows about is reported as exactly that, not as some plausible type.
	assert.Equal(t, "UNKNOWN(999)", connection.dataTypeName(999))
}
