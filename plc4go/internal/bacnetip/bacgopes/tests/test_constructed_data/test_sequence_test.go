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

package test_constructed_data

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

func TestEmptySequence(t *testing.T) {
	t.Run("test_empty_sequence", func(t *testing.T) {
		seq, err := NewEmptySequence(NoKWArgs())
		require.NoError(t, err)
		t.Logf("%#v", seq)

		// encode it in a tag list
		tagList := NewTagList(nil)
		err = seq.Encode(tagList)
		require.NoError(t, err)
		t.Logf("%v", tagList)

		// create another sequence and decode the tag list
		seq, err = NewEmptySequence(NoKWArgs())
		require.NoError(t, err)
		err = seq.Decode(tagList)
		require.NoError(t, err)
		t.Logf("%#v", seq)
	})
	t.Run("test_no_elements", func(t *testing.T) {
		// create another sequence and decode the tag list
		_, err := NewEmptySequence(NKW(KnownKey("some_element"), nil))
		require.Error(t, err)
	})
}

func TestSimpleSequence(t *testing.T) {
	t.Run("test_missing_element", func(t *testing.T) {
		// create a sequence with a missing required element
		seq, err := NewSimpleSequence(NoKWArgs())
		require.NoError(t, err)

		// encode it in a tag list
		tagList := NewTagList(nil)
		err = seq.Encode(tagList)
		require.Error(t, err)
	})
	t.Run("test_wrong_type", func(t *testing.T) {
		// create a sequence with wrong element value type
		seq, err := NewSimpleSequence(NKW(KnownKey("hydrogen"), 12))
		require.NoError(t, err)

		// encode it in a tag list
		tagList := NewTagList(nil)
		err = seq.Encode(tagList)
		require.Error(t, err)
	})
	t.Run("test_codec", func(t *testing.T) {
		// create a sequence
		seq, err := NewSimpleSequence(NKW(KnownKey("hydrogen"), false))
		require.NoError(t, err)

		// encode it in a tag list
		tagList := NewTagList(nil)
		err = seq.Encode(tagList)

		// create another sequence and decode the taglist
		seq, err = NewSimpleSequence(NoKWArgs())
		require.NoError(t, err)
		err = seq.Decode(tagList)
		require.NoError(t, err)
	})
}

func TestCompoundSequence1(t *testing.T) {
	t.Run("test_missing_element", func(t *testing.T) {
		// create a sequence with a missing required element
		seq, err := NewCompoundSequence1(NoKWArgs())
		require.NoError(t, err)

		// encode it in a tag list
		tagList := NewTagList(nil)
		err = seq.Encode(tagList)
		require.Error(t, err)

		// create a sequence with a missing required element
		seq, err = NewCompoundSequence1(NKW(KnownKey("hydrogen"), true))
		require.NoError(t, err)

		// encode it in a tag list
		tagList = NewTagList(nil)
		err = seq.Encode(tagList)
		require.Error(t, err)

		// create a sequence with a missing required element
		seq, err = NewCompoundSequence1(NKW(KnownKey("helium"), 2))
		require.NoError(t, err)

		// encode it in a tag list
		tagList = NewTagList(nil)
		err = seq.Encode(tagList)
		require.Error(t, err)

	})
	t.Run("test_codec", func(t *testing.T) {
		// create a sequence
		seq, err := NewCompoundSequence1(NKW(KnownKey("hydrogen"), true, KnownKey("helium"), 2))
		require.NoError(t, err)

		// encode it in a tag list
		tagList := NewTagList(nil)
		err = seq.Encode(tagList)

		// create another sequence and decode the taglist
		seq, err = NewCompoundSequence1(NoKWArgs())
		require.NoError(t, err)
		err = seq.Decode(tagList)
		require.NoError(t, err)
	})
}

func TestCompoundSequence2(t *testing.T) {
	t.Run("test_missing_element", func(t *testing.T) {
		// create a sequence with a missing required element
		seq, err := NewCompoundSequence2(NoKWArgs())
		require.NoError(t, err)

		// encode it in a tag list
		tagList := NewTagList(nil)
		err = seq.Encode(tagList)
		require.Error(t, err)

		// create a sequence with a missing required element
		seq, err = NewCompoundSequence2(NKW(KnownKey("lithium"), true))
		require.NoError(t, err)

		// encode it in a tag list
		tagList = NewTagList(nil)
		err = seq.Encode(tagList)
		require.Error(t, err)
	})
	t.Run("test_codec_1", func(t *testing.T) {
		// create a sequence
		seq, err := NewCompoundSequence2(NKW(KnownKey("beryllium"), 2))
		require.NoError(t, err)

		// encode it in a tag list
		tagList := NewTagList(nil)
		err = seq.Encode(tagList)

		// create another sequence and decode the taglist
		seq, err = NewCompoundSequence2(NoKWArgs())
		require.NoError(t, err)
		err = seq.Decode(tagList)
		require.NoError(t, err)
	})
	t.Run("test_codec_2", func(t *testing.T) {
		// create a sequence
		seq, err := NewCompoundSequence2(NKW(KnownKey("lithium"), true, KnownKey("beryllium"), 2))
		require.NoError(t, err)

		// encode it in a tag list
		tagList := NewTagList(nil)
		err = seq.Encode(tagList)

		// create another sequence and decode the taglist
		seq, err = NewCompoundSequence2(NoKWArgs())
		require.NoError(t, err)
		err = seq.Decode(tagList)
		require.NoError(t, err)
	})
}
