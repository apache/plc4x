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

using namespace org::apache::plc4x::cpp::api::metadata;
using namespace org::apache::plc4x::cpp::api::messages;

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
					/**
					* Interface defining the most basic methods a PLC4X connection should support.
					* This generally handles the connection establishment itself and the parsing of
					* field address strings to the platform dependent PlcField instances.
					*/
					class PlcConnection
					{
						public:
							/**
							* Establishes the connection to the remote PLC.
							* @throws PlcConnectionException if the connection attempt failed
							*/
							virtual void connect() = 0;

							/**
							* Indicates if the connection is established to a remote PLC.
							* @return {@code true} if connected, {@code false} otherwise
							*/
							virtual bool isConnected() = 0;

							/**
							 * Closes the connection to the remote PLC.
							 * @throws Exception if shutting down the connection failed
							 * Todo: implement java.lang.Autocloseable
							 */
						
								virtual void close() = 0;

							/**
							 * Provides connection metadata.
							 */
							virtual PlcConnectionMetadata* getMetadata() const = 0;

							/**
							 * Obtain read request builder.
							 * @throws PlcUnsupportedOperationException if the connection does not support reading
							 * Todo: implement pendant for java Builder pattern
							 */
							virtual PlcReadRequest* readRequestBuilder() = 0;

							/**
							 * Obtain write request builder.
							 * @throws PlcUnsupportedOperationException if the connection does not support writing
							  * Todo: implement pendant for java Builder pattern (PlcWriteRequest.Builder)
							 */
							virtual PlcWriteRequest* writeRequestBuilder() = 0;

							/**
							 * Obtain subscription request builder.
							 * @throws PlcUnsupportedOperationException if the connection does not support subscription
							 * Todo: implement pendant for java Builder pattern (PlcSubscriptionRequest.Builder)
							 */
							virtual PlcSubscriptionRequest* subscriptionRequestBuilder() = 0;

							/**
							 * Obtain unsubscription request builder.
							 * @throws PlcUnsupportedOperationException if the connection does not support subscription
							 * Todo: implement pendant for java Builder pattern (PlcSubscriptionRequest.Builder)

							 */
							virtual PlcUnsubscriptionRequest* unsubscriptionRequestBuilder() = 0;

						private:
					};
				}
			}
		}
	}
}

#endif