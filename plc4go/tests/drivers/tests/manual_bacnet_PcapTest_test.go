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

package tests

import (
	"io"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/options/converter"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/pcap"
)

func TestBacnetDriverWithPcap(t *testing.T) {
	t.Skip() // Manual test don't check in un-skipped

	file := path.Join(os.TempDir(), "bacnet-stack-services.cap")
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		println("File does not exist... re-downloading")
		if err := download(file, "https://wiki.wireshark.org/Protocols/bacnet?action=AttachFile&do=get&target=bacnet-stack-services.cap"); err != nil {
			panic(err)
		}
	}
	optionsForTesting := testutils.EnrichOptionsWithOptionsForTesting(t)

	driverManager := plc4go.NewPlcDriverManager(converter.WithOptionToExternal(optionsForTesting...)...)
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(bacnetip.NewDriver(optionsForTesting...))
	driverManager.(spi.TransportAware).RegisterTransport(pcap.NewTransport(optionsForTesting...))
	result := <-driverManager.GetConnection("bacnet-ip:pcap://" + file + "?transport-type=udp&speed-factor=0")
	if result.GetErr() != nil {
		panic(result.GetErr())
	}
	connection := result.GetConnection()
	defer connection.Close()
	build, err := connection.SubscriptionRequestBuilder().
		AddEventTagAddress("furz", "*/*/*").
		AddPreRegisteredConsumer("furz", func(event apiModel.PlcSubscriptionEvent) {
			println(event)
		}).
		Build()
	require.NoError(t, err)
	requestResult := <-build.Execute()
	if requestResult.GetErr() != nil {
		panic(requestResult.GetErr())
	}
	t.Logf("got response %v", requestResult.GetResponse())

	for connection.IsConnected() {
		t.Log("Still sleeping")
		time.Sleep(time.Second)
	}
}

func download(dstPath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	// Create the file
	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			panic(err)
		}
	}(out)

	_, err = io.Copy(out, resp.Body)
	return err
}
