//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

using System;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.spi.drivers
{
    /// <summary>
    /// Raised when a protocol message cannot be encoded, decoded or transmitted.
    /// </summary>
    public class MessageCodecException : Exception
    {
        public MessageCodecException(string message) : base(message)
        {
        }

        public MessageCodecException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }

    /// <summary>
    /// A protocol message that can report its own wire length and serialize itself.
    /// </summary>
    /// <remarks>
    /// Generated message classes are intended to implement this. Note that the C#
    /// code generator does not yet emit serialization logic (see the port notes in
    /// the PLC4Net README), so today this is implemented by hand-written messages
    /// and by test fixtures.
    /// </remarks>
    public interface IMessage
    {
        /// <summary>Length of this message once serialized, in bytes.</summary>
        int GetLengthInBytes();

        /// <summary>Writes this message into the supplied buffer.</summary>
        void Serialize(WriteBuffer writeBuffer);
    }

    /// <summary>
    /// Turns a transport's byte stream into whole protocol messages and back.
    /// </summary>
    /// <remarks>
    /// Mirrors the Java SPI3 <c>MessageCodecBase</c>. The framing contract is the
    /// three abstract members: a driver states the smallest header it can act on,
    /// derives the total frame length from that header, and parses a complete frame.
    /// Everything else — partial reads, several messages arriving in one packet,
    /// waiting for the rest of a frame — is handled here.
    /// </remarks>
    /// <typeparam name="TMessage">the message type this codec handles</typeparam>
    public abstract class MessageCodecBase<TMessage> where TMessage : class
    {
        private readonly string _protocolName;
        private readonly ITransportInstance _transportInstance;
        private readonly Action<TMessage> _messageHandler;

        protected MessageCodecBase(
            string protocolName,
            ITransportInstance transportInstance,
            Action<TMessage> messageHandler)
        {
            _protocolName = protocolName ?? throw new ArgumentNullException(nameof(protocolName));
            _transportInstance = transportInstance ?? throw new ArgumentNullException(nameof(transportInstance));
            _messageHandler = messageHandler ?? throw new ArgumentNullException(nameof(messageHandler));
        }

        protected ITransportInstance TransportInstance => _transportInstance;

        /// <summary>
        /// Smallest number of bytes needed before <see cref="CalculateTotalMessageSize"/>
        /// can say anything useful.
        /// </summary>
        protected abstract int GetMinimumHeaderSize();

        /// <summary>
        /// Total frame length in bytes, derived from the peeked header. Return a negative
        /// value to signal "cannot tell yet, wait for more bytes".
        /// </summary>
        protected abstract int CalculateTotalMessageSize(byte[] header, int availableBytes);

        /// <summary>Parses exactly one complete frame.</summary>
        protected abstract TMessage ParseMessage(ReadBuffer readBuffer);

        public void Send(TMessage message)
        {
            if (message == null)
            {
                throw new ArgumentNullException(nameof(message));
            }
            if (!(message is IMessage serializable))
            {
                throw new MessageCodecException(
                    $"{typeof(TMessage).Name} does not implement IMessage and cannot be serialized.");
            }

            try
            {
                var writeBuffer = new WriteBuffer();
                serializable.Serialize(writeBuffer);
                _transportInstance.Write(writeBuffer.GetBytes());
            }
            catch (TransportException e)
            {
                throw new MessageCodecException($"Failed to send {_protocolName} message", e);
            }
            catch (ParseException e)
            {
                throw new MessageCodecException($"Failed to serialize {_protocolName} message", e);
            }
        }

        /// <summary>
        /// Drains as many complete messages as the buffer currently holds, handing each
        /// to the message handler. Returns without consuming anything when the next
        /// frame has not fully arrived.
        /// </summary>
        public void ProcessIncomingData()
        {
            try
            {
                while (true)
                {
                    var minimumHeaderSize = GetMinimumHeaderSize();

                    var availableBytes = _transportInstance.GetNumBytesAvailable();
                    if (availableBytes < minimumHeaderSize)
                    {
                        return;
                    }

                    var header = _transportInstance.PeekReadableBytes(minimumHeaderSize);
                    var totalMessageSize = CalculateTotalMessageSize(header, availableBytes);
                    if (totalMessageSize < 0)
                    {
                        return;
                    }

                    if (availableBytes < totalMessageSize)
                    {
                        // The frame is still in flight; leave the bytes buffered.
                        return;
                    }

                    var messageBytes = _transportInstance.Read(totalMessageSize);
                    var message = ParseMessage(new ReadBuffer(messageBytes));
                    _messageHandler(message);
                }
            }
            catch (TransportException e)
            {
                throw new MessageCodecException($"Failed to receive {_protocolName} message", e);
            }
            catch (ParseException e)
            {
                throw new MessageCodecException($"Failed to parse {_protocolName} message", e);
            }
        }

        public bool IsOpen => _transportInstance.IsOpen;

        public void Close()
        {
            try
            {
                _transportInstance.Close();
            }
            catch (TransportException e)
            {
                throw new MessageCodecException("Failed to close transport", e);
            }
        }
    }
}
