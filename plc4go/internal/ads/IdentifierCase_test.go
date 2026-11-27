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

package ads

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/internal/ads/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
)

// IEC 61131-3 identifiers are case insensitive. A TwinCAT 3 device stores symbol names as they
// were declared ("MAIN.myBool"), a TwinCAT 2 device stores every one of them in upper case, so an
// address can only be written to match one of the two unless resolution ignores case.
//
// The device here is a TwinCAT 2 one: one function block symbol, upper cased throughout.
func connectionWithUpperCaseSymbol() *Connection {
	boolType := driverModel.NewAdsDataTypeTableEntryBuilder().
		WithSize(1).
		WithMainName("BOOL").
		WithSecondaryName("BOOL").
		MustBuild()
	member := driverModel.NewAdsDataTypeTableEntryBuilder().
		WithSize(1).
		WithOffset(4).
		WithMainName("MYBOOL").
		WithSecondaryName("BOOL").
		MustBuild()
	block := driverModel.NewAdsDataTypeTableEntryBuilder().
		WithSize(8).
		WithMainName("FB_TEST").
		WithSecondaryName("FB_TEST").
		WithNumChildren(1).
		WithChildren(member).
		MustBuild()
	symbol := driverModel.NewAdsSymbolTableEntryBuilder().
		WithGroup(0x4020).
		WithOffset(100).
		WithSize(8).
		WithName("MAIN.FBTEST").
		WithDataTypeName("FB_TEST").
		MustBuild()

	return &Connection{
		log: zerolog.Nop(),
		driverContext: &DriverContext{
			// The fallback is only reached on a device that reported itself as TwinCAT 2.
			isTwinCat2: true,
			dataTypeTable: map[string]driverModel.AdsDataTypeTableEntry{
				"BOOL": boolType, "FB_TEST": block,
			},
			symbolTable: map[string]driverModel.AdsSymbolTableEntry{"MAIN.FBTEST": symbol},
		},
	}
}

// Every spelling of the same address resolves to the same place: the member sits 4 bytes into a
// symbol at offset 100.
func TestIdentifierCase_symbolAndMemberResolveWhateverTheSpelling(t *testing.T) {
	for _, address := range []string{
		"MAIN.FBTEST.MYBOOL",
		"MAIN.fbTest.myBool",
		"main.fbtest.mybool",
		"MaIn.FbTeSt.MyBoOl",
	} {
		t.Run(address, func(t *testing.T) {
			tag, err := connectionWithUpperCaseSymbol().resolveSymbolicTag(context.Background(),
				model.SymbolicPlcTag{SymbolicAddress: address})
			require.NoError(t, err)
			assert.Equal(t, uint32(0x4020), tag.IndexGroup)
			assert.Equal(t, uint32(104), tag.IndexOffset)
		})
	}
}

// A name that differs by more than case is still not found.
func TestIdentifierCase_unknownNamesAreStillUnresolved(t *testing.T) {
	_, err := connectionWithUpperCaseSymbol().resolveSymbolicTag(context.Background(),
		model.SymbolicPlcTag{SymbolicAddress: "MAIN.FBTEST.MYBOOLEAN"})
	assert.Error(t, err)

	_, err = connectionWithUpperCaseSymbol().resolveSymbolicTag(context.Background(),
		model.SymbolicPlcTag{SymbolicAddress: "MAIN.OTHER"})
	assert.Error(t, err)
}

// An exact match is preferred, so a table that happens to hold two spellings keeps resolving the
// exact one to itself. IEC 61131-3 does not allow this within one scope, but the preference is
// what keeps every address that resolves today resolving to the same symbol.
func TestIdentifierCase_exactMatchWins(t *testing.T) {
	connection := connectionWithUpperCaseSymbol()
	lowerCased := driverModel.NewAdsSymbolTableEntryBuilder().
		WithGroup(0x4030).
		WithOffset(200).
		WithSize(8).
		WithName("MAIN.fbtest").
		WithDataTypeName("FB_TEST").
		MustBuild()
	connection.driverContext.symbolTable["MAIN.fbtest"] = lowerCased

	entry, ok := connection.driverContext.lookupSymbol("MAIN.fbtest")
	require.True(t, ok)
	assert.Equal(t, uint32(0x4030), entry.GetGroup())

	entry, ok = connection.driverContext.lookupSymbol("MAIN.FBTEST")
	require.True(t, ok)
	assert.Equal(t, uint32(0x4020), entry.GetGroup())
}

func TestIdentifierCase_childByName(t *testing.T) {
	driverContext := connectionWithUpperCaseSymbol().driverContext
	block, ok := driverContext.lookupDataType("FB_TEST")
	require.True(t, ok)

	for _, name := range []string{"MYBOOL", "myBool", "mybool"} {
		child, found := driverContext.childByName(block, name)
		require.True(t, found, name)
		assert.Equal(t, "MYBOOL", child.GetMainName())
	}

	_, found := driverContext.childByName(block, "MYBOOLEAN")
	assert.False(t, found)
}

// On a device that is not TwinCAT 2 nothing falls back: the table is matched as it is
// written, which is what TwinCAT 3 does today.
func TestIdentifierCase_noFallbackWhenNotTwinCat2(t *testing.T) {
	connection := connectionWithUpperCaseSymbol()
	connection.driverContext.isTwinCat2 = false

	_, err := connection.resolveSymbolicTag(context.Background(),
		model.SymbolicPlcTag{SymbolicAddress: "main.fbtest.mybool"})
	assert.Error(t, err, "a lower case spelling must not resolve on TwinCAT 3")

	tag, err := connection.resolveSymbolicTag(context.Background(),
		model.SymbolicPlcTag{SymbolicAddress: "MAIN.FBTEST.MYBOOL"})
	require.NoError(t, err, "the spelling the table holds still resolves")
	assert.Equal(t, uint32(104), tag.IndexOffset)
}
