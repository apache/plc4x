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

package firmata

import (
	"context"
	"fmt"
	"runtime/debug"
	"slices"
	"sync"

	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/firmata/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
)

// unknownAnalogValue is what a subscription reports for a pin the board hasn't sampled yet. Ported
// from plc4j's FirmataConnection, which fills the gaps in a multi-pin event with -1.
const unknownAnalogValue = int16(-1)

// Connection is a connection to a firmata board. It is the port of plc4j's FirmataConnection
// together with the FirmataDriverContext it used to keep its pin bookkeeping in.
//
// Firmata is push-only: after a one-time system-reset handshake the board emits unsolicited
// analog-IO and digital-IO messages whenever a pin it was told to report on changes. The connection
// caches the latest value per pin and hands changes to the subscribers, which filter them against
// what their handles actually cover. There is no read path at all - plc4j's FirmataConnection
// implements writing and subscribing and nothing else, and so does this one.
type Connection struct {
	_default.DefaultConnection

	configuration Configuration
	messageCodec  spi.MessageCodec
	options       map[string][]string

	connectionId string
	tracer       tracer.Tracer

	// pinMutex guards the two pin-mode registries below, which the writer and the subscribers
	// both claim pins in.
	pinMutex sync.Mutex
	// digitalPins is what every digital pin we have touched has been configured as. A pin can only
	// ever have one mode, so a pin claimed as an output can't later be subscribed to (and the
	// other way around) - the same rule plc4j's FirmataConnection enforces.
	digitalPins map[uint8]readWriteModel.PinMode
	// analogPins is the set of analog pins reporting has been switched on for.
	analogPins map[uint8]readWriteModel.PinMode

	// valueMutex guards the value caches, which the codec worker writes and the subscribers read.
	valueMutex sync.RWMutex
	// analogValues is the latest sample of every analog pin the board has reported.
	analogValues map[uint8]int16
	// digitalValues is the latest state of every digital pin the board has reported.
	digitalValues map[uint8]bool

	// subscribersMutex guards the subscriber list, which is written by Subscribe and read by the
	// worker delivering incoming messages.
	subscribersMutex sync.RWMutex
	subscribers      []*Subscriber

	// lifecycleMutex guards the worker's cancel function so Close can be called concurrently with,
	// and without, a preceding Connect.
	lifecycleMutex sync.Mutex
	workerCancel   context.CancelFunc

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

var (
	_ spi.TransportInstanceExposer = (*Connection)(nil)
)

func NewConnection(configuration Configuration, messageCodec spi.MessageCodec, connectionOptions map[string][]string, tagHandler spi.PlcTagHandler, _options ...options.WithOption) *Connection {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	// Clipping the option list is what keeps the write and the subscription builders from racing:
	// both hand their own logger option downstream with append, and on a slice which still has
	// spare capacity those two appends write into the very same backing array - a torn interface
	// value is the failure mode, as ExtractCustomLogger reads the slot while the other append
	// overwrites it.
	_options = slices.Clip(_options)
	connection := &Connection{
		configuration: configuration,
		messageCodec:  messageCodec,
		options:       connectionOptions,
		digitalPins:   map[uint8]readWriteModel.PinMode{},
		analogPins:    map[uint8]readWriteModel.PinMode{},
		analogValues:  map[uint8]int16{},
		digitalValues: map[uint8]bool{},
		log:           customLogger,
		_options:      _options,
	}
	if traceEnabledOption, ok := connectionOptions["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = tracer.NewTracer(connection.connectionId, _options...)
		}
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		append(_options,
			_default.WithPlcTagHandler(tagHandler),
			_default.WithPlcValueHandler(NewValueHandler(_options...)),
		)...,
	)
	return connection
}

func (c *Connection) GetConnectionId() string {
	return c.connectionId
}

func (c *Connection) IsTraceEnabled() bool {
	return c.tracer != nil
}

func (c *Connection) GetTracer() tracer.Tracer {
	return c.tracer
}

func (c *Connection) GetConnection() plc4go.PlcConnection {
	return c
}

func (c *Connection) GetMessageCodec() spi.MessageCodec {
	return c.messageCodec
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Connect / close
////////////////////////////////////////////////////////////////////////////////////////////////////

// Connect opens the transport, starts listening for the board's unsolicited messages and then runs
// the handshake plc4j's FirmataConnection runs: a system reset, answered by the board with its
// protocol version followed by a report-firmware message.
func (c *Connection) Connect(ctx context.Context) error {
	c.log.Trace().Msg("Connecting")
	if err := c.DefaultConnection.Connect(ctx); err != nil {
		return errors.Wrap(err, "error connecting default connection")
	}
	c.startIncomingWorker()
	if err := c.handshake(ctx); err != nil {
		if closeErr := c.Close(); closeErr != nil {
			c.log.Debug().Err(closeErr).Msg("error closing connection after a failed handshake")
		}
		return errors.Wrap(err, "error during the firmata connect handshake")
	}
	return nil
}

func (c *Connection) Close() error {
	c.log.Trace().Msg("Closing")
	c.lifecycleMutex.Lock()
	workerCancel := c.workerCancel
	c.workerCancel = nil
	c.lifecycleMutex.Unlock()
	if workerCancel != nil {
		workerCancel()
	}
	err := c.DefaultConnection.Close()
	c.wg.Wait()

	c.subscribersMutex.Lock()
	subscribers := c.subscribers
	c.subscribers = nil
	c.subscribersMutex.Unlock()
	for _, subscriber := range subscribers {
		subscriber.Close()
	}
	return err
}

// handshake resets the board and waits for it to report its firmware, which is how plc4j's
// FirmataConnection decides a board is really there. The protocol-version message the board sends
// first isn't waited for: it carries nothing this driver needs and arrives before the reset is even
// acknowledged.
func (c *Connection) handshake(ctx context.Context) error {
	requestCtx, cancelRequest := context.WithTimeout(ctx, c.configuration.requestTimeout)
	defer cancelRequest()

	reportChan := make(chan readWriteModel.SysexCommandReportFirmwareResponse, 1)
	errChan := make(chan error, 1)
	systemReset := readWriteModel.NewFirmataMessageCommand(readWriteModel.NewFirmataCommandSystemReset())
	if err := c.messageCodec.SendRequest(requestCtx, "handshake", systemReset,
		func(message spi.Message) bool {
			_, ok := firmwareReport(message)
			return ok
		},
		func(message spi.Message) error {
			report, _ := firmwareReport(message)
			select {
			case reportChan <- report:
			default:
				c.log.Warn().Msg("failed to hand on the report-firmware response")
			}
			return nil
		},
		func(err error) error {
			select {
			case errChan <- err:
			default:
				c.log.Warn().Msg("failed to hand on the handshake error")
			}
			return nil
		}); err != nil {
		return errors.Wrap(err, "error sending the system reset")
	}

	// A board only volunteers its firmware report when it really restarts: StandardFirmata's
	// systemResetCallback just clears the pin state, printFirmwareVersion is called from setup().
	// Over serial that is masked - opening the port toggles DTR, which resets the AVR - but nothing
	// resets a WiFi or Ethernet board when a TCP connection is opened, so waiting for a volunteered
	// report would time out on exactly the transport this driver newly speaks. Asking for the report
	// explicitly (0xF0 0x79 0xF7) covers that case: over TCP the board is up and answers, and over
	// serial the request may well be lost in the bootloader the reset dropped the board into - which
	// costs nothing, as that is precisely the case where the report is volunteered. A board which
	// answers both simply answers twice, and the second answer has no expectation left to match.
	// plc4j's FirmataConnection.onConnect only ever waits, which is the gap.
	reportFirmwareRequest := readWriteModel.NewFirmataMessageCommand(
		readWriteModel.NewFirmataCommandSysex(readWriteModel.NewSysexCommandReportFirmwareRequest()))
	if err := c.messageCodec.Send(requestCtx, "handshake", reportFirmwareRequest); err != nil {
		return errors.Wrap(err, "error requesting the firmware report")
	}

	select {
	case report := <-reportChan:
		c.log.Info().
			Uint8("majorVersion", report.GetMajorVersion()).
			Uint8("minorVersion", report.GetMinorVersion()).
			Str("name", string(report.GetFileName())).
			Msg("Connected to a firmata host")
		return nil
	case err := <-errChan:
		return errors.Wrap(err, "error waiting for the report-firmware response")
	case <-requestCtx.Done():
		return requestCtx.Err()
	}
}

// firmwareReport picks the report-firmware response out of an incoming message, which the board
// wraps in a sysex command inside a system message.
func firmwareReport(message spi.Message) (readWriteModel.SysexCommandReportFirmwareResponse, bool) {
	command, ok := message.(readWriteModel.FirmataMessageCommand)
	if !ok {
		return nil, false
	}
	sysex, ok := command.GetCommand().(readWriteModel.FirmataCommandSysex)
	if !ok {
		return nil, false
	}
	report, ok := sysex.GetCommand().(readWriteModel.SysexCommandReportFirmwareResponse)
	return report, ok
}

// startIncomingWorker pumps the messages the codec couldn't match to an expectation - which after
// the handshake is every message the board sends - into the value caches and on to the subscribers.
func (c *Connection) startIncomingWorker() {
	workerCtx, workerCancel := context.WithCancel(context.Background())
	c.lifecycleMutex.Lock()
	c.workerCancel = workerCancel
	c.lifecycleMutex.Unlock()

	incomingMessageChannel := c.messageCodec.GetDefaultIncomingMessageChannel()
	c.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				c.log.Error().
					Str("stack", string(debug.Stack())).
					Interface("err", err).
					Msg("panic-ed")
			}
		}()
		for {
			select {
			case <-workerCtx.Done():
				c.log.Info().Msg("Ending incoming message transfer")
				return
			case message := <-incomingMessageChannel:
				c.handleIncomingMessage(message)
			}
		}
	})
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Incoming messages
////////////////////////////////////////////////////////////////////////////////////////////////////

// handleIncomingMessage caches what a board reported and hands on what actually changed. Only
// changes are published, the way plc4j's FirmataConnection does it: the board keeps sampling and
// resends a value it has already sent.
func (c *Connection) handleIncomingMessage(message spi.Message) {
	switch typed := message.(type) {
	case readWriteModel.FirmataMessageAnalogIO:
		pin := typed.GetPin()
		value, err := decodeAnalogValue(typed.GetData())
		if err != nil {
			c.log.Debug().Err(err).Uint8("pin", pin).Msg("discarding a malformed analog-IO message")
			return
		}
		if !c.updateAnalogValue(pin, value) {
			return
		}
		for _, subscriber := range c.activeSubscribers() {
			subscriber.handleAnalogUpdate(pin)
		}
	case readWriteModel.FirmataMessageDigitalIO:
		changedPins, err := c.updateDigitalValues(typed.GetPinBlock(), typed.GetData())
		if err != nil {
			c.log.Debug().Err(err).Uint8("pinBlock", typed.GetPinBlock()).Msg("discarding a malformed digital-IO message")
			return
		}
		if len(changedPins) == 0 {
			return
		}
		for _, subscriber := range c.activeSubscribers() {
			subscriber.handleDigitalUpdate(changedPins)
		}
	default:
		c.log.Debug().Type("message", message).Msg("Unexpected firmata message")
	}
}

// updateAnalogValue caches a sample and reports whether it differs from the one before it.
func (c *Connection) updateAnalogValue(pin uint8, value int16) bool {
	c.valueMutex.Lock()
	defer c.valueMutex.Unlock()
	if previous, ok := c.analogValues[pin]; ok && previous == value {
		return false
	}
	c.analogValues[pin] = value
	return true
}

// updateDigitalValues caches the state of the 8 pins of one port and reports which of them changed.
func (c *Connection) updateDigitalValues(pinBlock uint8, data []int8) ([]uint8, error) {
	port, err := decodeDigitalPortValue(data)
	if err != nil {
		return nil, err
	}
	firstPin := int(pinBlock) * 8
	c.valueMutex.Lock()
	defer c.valueMutex.Unlock()
	var changedPins []uint8
	for i := 0; i < 8; i++ {
		pin := firstPin + i
		if pin > maxDigitalPin {
			break
		}
		value := port&(1<<i) != 0
		// A pin nobody has heard from yet counts as low, the way plc4j's all-zero BitSet does.
		// Anything else would make the first report of a port look like a change of all 8 of its
		// pins and wake every subscriber of the port up for nothing.
		if c.digitalValues[uint8(pin)] == value {
			continue
		}
		c.digitalValues[uint8(pin)] = value
		changedPins = append(changedPins, uint8(pin))
	}
	return changedPins, nil
}

// analogValue is the latest sample of a pin, or -1 for a pin the board hasn't reported yet (plc4j
// FirmataConnection.publishAnalogEvents).
func (c *Connection) analogValue(pin uint8) int16 {
	c.valueMutex.RLock()
	defer c.valueMutex.RUnlock()
	if value, ok := c.analogValues[pin]; ok {
		return value
	}
	return unknownAnalogValue
}

// digitalValue is the latest state of a pin, low for a pin the board hasn't reported yet.
func (c *Connection) digitalValue(pin uint8) bool {
	c.valueMutex.RLock()
	defer c.valueMutex.RUnlock()
	return c.digitalValues[pin]
}

// decodeAnalogValue reassembles the 14 bit sample an analog-IO message carries: the low 7 bits of
// the first byte are the least significant part, the low 7 bits of the second byte the most
// significant one.
func decodeAnalogValue(data []int8) (int16, error) {
	if len(data) < 2 {
		return 0, errors.Errorf("an analog-IO message carries 2 data bytes, got %d", len(data))
	}
	return int16(uint16(data[0])&0x7F | (uint16(data[1])&0x7F)<<7), nil
}

// decodeDigitalPortValue reassembles the 8 pin states a digital-IO message carries: the low 7 bits
// of the first byte are the first 7 pins of the port, bit 0 of the second byte is the eighth.
func decodeDigitalPortValue(data []int8) (uint8, error) {
	if len(data) < 2 {
		return 0, errors.Errorf("a digital-IO message carries 2 data bytes, got %d", len(data))
	}
	port := uint8(data[0]) & 0x7F
	if uint8(data[1])&0x01 == 0x01 {
		port |= 0x80
	}
	return port, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Pin bookkeeping
////////////////////////////////////////////////////////////////////////////////////////////////////

// pinClaim undoes a claim whose messages never made it onto the wire. A claim is recorded before
// the board is told about it - the messages which configure the pins are what the caller sends -
// so a send which fails part-way would otherwise leave the driver believing pins are configured
// that the board has never heard of: the retry would find them claimed, emit no set-pin-mode and
// no report-digital, and report a success which does nothing.
//
// Only the pins this very claim added are given back, so a claim which found a pin already in the
// right mode leaves that pin (and whoever configured it) alone. A nil claim is a no-op.
type pinClaim func()

// rollback gives the newly claimed pins back, if there are any.
func (claim pinClaim) rollback() {
	if claim != nil {
		claim()
	}
}

// releaseDigitalPins takes the given pins back out of the digital registry.
func (c *Connection) releaseDigitalPins(pins []uint8) {
	if len(pins) == 0 {
		return
	}
	c.pinMutex.Lock()
	defer c.pinMutex.Unlock()
	for _, pin := range pins {
		delete(c.digitalPins, pin)
	}
}

// releaseAnalogPins takes the given pins back out of the analog registry.
func (c *Connection) releaseAnalogPins(pins []uint8) {
	if len(pins) == 0 {
		return
	}
	c.pinMutex.Lock()
	defer c.pinMutex.Unlock()
	for _, pin := range pins {
		delete(c.analogPins, pin)
	}
}

// claimOutputPins reserves a run of digital pins as outputs and reports the set-pin-mode messages
// which have to go out before a value can be written to them. A pin which is already an output
// needs no message; a pin which is something else cannot be written to at all. The returned claim
// hands the newly reserved pins back should those messages not make it out.
func (c *Connection) claimOutputPins(address uint8, quantity uint8) ([]readWriteModel.FirmataMessage, pinClaim, error) {
	c.pinMutex.Lock()
	defer c.pinMutex.Unlock()
	if err := c.checkDigitalPins(address, quantity, readWriteModel.PinMode_PinModeOutput); err != nil {
		return nil, nil, err
	}
	var messages []readWriteModel.FirmataMessage
	var claimedPins []uint8
	for pin := int(address); pin < int(address)+int(quantity); pin++ {
		if _, claimed := c.digitalPins[uint8(pin)]; claimed {
			continue
		}
		c.digitalPins[uint8(pin)] = readWriteModel.PinMode_PinModeOutput
		claimedPins = append(claimedPins, uint8(pin))
		messages = append(messages, readWriteModel.NewFirmataMessageCommand(
			readWriteModel.NewFirmataCommandSetPinMode(uint8(pin), readWriteModel.PinMode_PinModeOutput)))
	}
	return messages, func() { c.releaseDigitalPins(claimedPins) }, nil
}

// claimDigitalInputPins reserves a run of digital pins as inputs (or pullup inputs) and reports the
// messages which switch reporting on for them: a set-pin-mode per newly claimed pin plus a
// report-digital per port they sit in, as reporting is switched on per port of 8 pins rather than
// per pin.
func (c *Connection) claimDigitalInputPins(address uint8, quantity uint8, mode readWriteModel.PinMode) ([]readWriteModel.FirmataMessage, pinClaim, error) {
	c.pinMutex.Lock()
	defer c.pinMutex.Unlock()
	if err := c.checkDigitalPins(address, quantity, mode); err != nil {
		return nil, nil, err
	}
	var messages []readWriteModel.FirmataMessage
	var claimedPins []uint8
	var ports []uint8
	seenPorts := map[uint8]bool{}
	for pin := int(address); pin < int(address)+int(quantity); pin++ {
		if _, claimed := c.digitalPins[uint8(pin)]; claimed {
			// Already reporting in exactly this mode - checkDigitalPins made sure of that.
			continue
		}
		c.digitalPins[uint8(pin)] = mode
		claimedPins = append(claimedPins, uint8(pin))
		messages = append(messages, readWriteModel.NewFirmataMessageCommand(
			readWriteModel.NewFirmataCommandSetPinMode(uint8(pin), mode)))
		if port := uint8(pin) / 8; !seenPorts[port] {
			seenPorts[port] = true
			ports = append(ports, port)
		}
	}
	for _, port := range ports {
		messages = append(messages, readWriteModel.NewFirmataMessageSubscribeDigitalPinValue(port, true))
	}
	return messages, func() { c.releaseDigitalPins(claimedPins) }, nil
}

// claimAnalogInputPins switches reporting on for a run of analog pins. Analog pins need no
// set-pin-mode: asking the board to report an analog pin is enough (plc4j FirmataConnection
// buildSubscribeMessages).
func (c *Connection) claimAnalogInputPins(address uint8, quantity uint8) ([]readWriteModel.FirmataMessage, pinClaim, error) {
	c.pinMutex.Lock()
	defer c.pinMutex.Unlock()
	var messages []readWriteModel.FirmataMessage
	var claimedPins []uint8
	for pin := int(address); pin < int(address)+int(quantity); pin++ {
		if _, claimed := c.analogPins[uint8(pin)]; claimed {
			continue
		}
		c.analogPins[uint8(pin)] = readWriteModel.PinMode_PinModeInput
		claimedPins = append(claimedPins, uint8(pin))
		messages = append(messages, readWriteModel.NewFirmataMessageSubscribeAnalogPinValue(uint8(pin), true))
	}
	return messages, func() { c.releaseAnalogPins(claimedPins) }, nil
}

// checkDigitalPins refuses a run of pins in which any single pin is already claimed for something
// else. Checking the whole run before claiming any of it keeps a rejected tag from leaving half its
// pins reconfigured. Callers hold pinMutex.
func (c *Connection) checkDigitalPins(address uint8, quantity uint8, mode readWriteModel.PinMode) error {
	for pin := int(address); pin < int(address)+int(quantity); pin++ {
		existing, claimed := c.digitalPins[uint8(pin)]
		if claimed && existing != mode {
			return errors.Errorf("pin %d is already configured as %s and can't be used as %s",
				pin, existing.String(), mode.String())
		}
	}
	return nil
}

// sendAll ships a batch of messages in order. Firmata acknowledges nothing, so there is nothing to
// wait for in between.
func (c *Connection) sendAll(ctx context.Context, interactionInfo string, messages []readWriteModel.FirmataMessage) error {
	for _, message := range messages {
		if err := c.messageCodec.Send(ctx, interactionInfo, message); err != nil {
			return errors.Wrapf(err, "error sending %T", message)
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Subscribers
////////////////////////////////////////////////////////////////////////////////////////////////////

func (c *Connection) addSubscriber(subscriber *Subscriber) {
	c.subscribersMutex.Lock()
	defer c.subscribersMutex.Unlock()
	for _, existing := range c.subscribers {
		if existing == subscriber {
			return
		}
	}
	c.subscribers = append(c.subscribers, subscriber)
}

func (c *Connection) removeSubscriber(subscriber *Subscriber) {
	c.subscribersMutex.Lock()
	defer c.subscribersMutex.Unlock()
	for i, existing := range c.subscribers {
		if existing == subscriber {
			c.subscribers = append(c.subscribers[:i], c.subscribers[i+1:]...)
			return
		}
	}
}

// activeSubscribers snapshots the subscriber list so events can be delivered without holding the
// lock while calling into consumer code, which may subscribe or unsubscribe from inside a callback.
func (c *Connection) activeSubscribers() []*Subscriber {
	c.subscribersMutex.RLock()
	defer c.subscribersMutex.RUnlock()
	return append([]*Subscriber(nil), c.subscribers...)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Metadata and builders
////////////////////////////////////////////////////////////////////////////////////////////////////

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &_default.DefaultConnectionMetadata{
		// Firmata has no read path: the board pushes values for the pins it was told to report on
		// and answers no request for a current value. plc4j's FirmataDriver says the same by
		// implementing canWrite and canSubscribe and nothing else.
		ProvidesReading:     false,
		ProvidesWriting:     true,
		ProvidesSubscribing: true,
		ProvidesBrowsing:    false,
	}
}

func (c *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilder(
		c.GetPlcTagHandler(),
		c.GetPlcValueHandler(),
		NewWriter(c, append(c._options, options.WithCustomLogger(c.log))...),
	)
}

func (c *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcSubscriptionRequestBuilder(
		c.GetPlcTagHandler(),
		c.GetPlcValueHandler(),
		NewSubscriber(c, append(c._options, options.WithCustomLogger(c.log))...),
	)
}

// UnsubscriptionRequestBuilder hands out the default builder: the handles carry the subscriber they
// belong to, so nothing driver specific is needed to take them back.
func (c *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcUnsubscriptionRequestBuilder()
}

func (c *Connection) String() string {
	return fmt.Sprintf("firmata.Connection{requestTimeout: %s}", c.configuration.requestTimeout)
}
