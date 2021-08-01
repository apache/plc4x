/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
#ifndef _PLC_CONNECTION
#define _PLC_CONNECTION

#include "metadata/PlcConnectionMetadata.h"
#include "../api/messages/PlcReadRequest.h"
#include "../api/messages/PlcWriteRequest.h"
#include "../api/messages/PlcSubscriptionRequest.h"
#include "../api/messages/PlcUnsubscriptionRequest.h"

#include <org/apache/plc4x/cpp/utils/logger/DatDmp.h>
#include <boost/signals2.hpp>
#include <boost/thread.hpp>

#define DEFAULT_BUFFER_SIZE	16384

namespace org
{
  namespace apache
  {
    namespace plc4x
    {
	    namespace cpp
	    {
			namespace api
			{

				using namespace boost::signals2;
                using namespace org::apache::plc4x::cpp::api::metadata;
                using namespace org::apache::plc4x::cpp::api::messages;
                using namespace org::apache::plc4x::cpp::utils;

				/**
				* Interface defining the most basic methods a PLC4X connection should support.
				* This generally handles the connection establishment itself and the parsing of
				* field address strings to the platform dependent PlcField instances.
				*/
				
				class PlcConnection
				{
				public:

					PlcConnection();
					~PlcConnection();
					
					
					/**
					* Establishes the connection to the remote PLC.
					* @throws PlcConnectionException if the connection attempt failed
					*/
					virtual void connect() = 0;
					
					/**
					* Indicates if the connection is established to a remote PLC.
					* @return {@code true} if connected, {@code false} otherwise
					*/
					
					inline bool isConnected() const { return _bConnected; }
					
					/**
					 * Closes the connection to the remote PLC.
					 * @throws Exception if shutting down the connection failed
					 */
					virtual void close() = 0;

					/**
					 * Parse a fieldQuery for the given connection type.
					 *
					 * @throws PlcRuntimeException If the string cannot be parsed
					 */
                    PlcField* prepareField(std::string strFieldQuery);

					/**
					 * Provides connection metadata.
					 */
					virtual PlcConnectionMetadata* getMetadata() = 0;

					/**
					 * Execute a ping query against a remote device to check the availability of the connection.
					 *
					 * @return CompletableFuture that is completed successfully (Void) or unsuccessfully with an PlcException.
					 */
					virtual void ping() = 0;

					/**
					 * Obtain read request builder.
					 * @throws PlcUnsupportedOperationException if the connection does not support reading
					 */
					virtual PlcReadRequest::Builder* readRequestBuilder() = 0;

					/**
					 * Obtain write request builder.
					 * @throws PlcUnsupportedOperationException if the connection does not support writing
					 */
					//virtual void Builder* writeRequestBuilder() = 0;

					/**
					 * Obtain subscription request builder.
					 * @throws PlcUnsupportedOperationException if the connection does not support subscription
					 */
					//virtual void Builder* subscriptionRequestBuilder() = 0;

					/**
					 * Obtain unsubscription request builder.
					 * @throws PlcUnsupportedOperationException if the connection does not support subscription
					 */
					//virtual void Builder* unsubscriptionRequestBuilder() = 0;

					
					virtual bool send(unsigned char* pBytesToSend, int iNumBytesToSend) = 0;
					

					inline int getSendBufferSize() const { return _iSendBufferSize; }
					void setSendBufferSize(int iBufferSize);
				
					inline int getRecvBufferSize() const { return _iRecvBufferSize; }
					void setRecvBufferSize(int iBlockSize);

				protected:

					virtual void onReceive(const boost::system::error_code& errorCode, std::size_t szRecvBytes) = 0;

				protected:
					bool            _bConnected;
					unsigned char*  _pSendBuffer;
					int             _iSendBufferSize;
					unsigned char*  _pRecvBuffer;
					int             _iRecvBufferSize;
					std::string     _strErrorMessage;

					boost::mutex              _mtxComLock;
					boost::condition_variable _condComHandler;

                    DatDmp _dump;

				};
			}
		}
    }
  }
}

#endif