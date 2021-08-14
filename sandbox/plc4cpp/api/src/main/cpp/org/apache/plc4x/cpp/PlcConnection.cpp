/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
 
#include "PlcConnection.h"
#include "Logger.h"

org.apache.plc4x.java.api

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

				PlcConnection::PlcConnection() :
					_bOpen(false),
					_pSendBuffer(NULL),
					_pRecvBuffer(NULL),
					_iSendBufferSize(0),
					_iRecvBufferSize(0),
				{
					setSendBufferSize(DEFAULT_BUFFER_SIZE);
					setRecvBufferSize(DEFAULT_BUFFER_SIZE);
					_strErrorMessage = "";
				}

				PlcConnection::~PlcConnection()
				{
					if (_pSendBuffer != NULL) delete[] _pSendBuffer;
					if (_pRecvBuffer != NULL) delete[] _pRecvBuffer;
				}

				PlcField PlcConnection::prepareField(String fieldQuery) throws PlcInvalidFieldException 
				{
						BOOST_THROW_EXCEPTION(PlcRuntimeException("Parse method is not implemented for this connection / driver"));
				}
				
				void PlcConnection::setSendBufferSize(int iBufferSize)
				{
					if (iBlockSize != _iSendBufferSize)
					{
						_iSendBufferSize = iBlockSize;
						if (_pSendBuffer != NULL) delete[] _pSendBuffer;
						_pSendBuffer = new unsigned char[_iBlockSizeSend];
					}
				}

				void PlcConnection::setRecvBufferSize(int iBufferSize)
				{
					if (iBlockSize != _iRecvBufferSize)
					{
						_iRecvBufferSize = iBufferSize;
						if (_pRecvBuffer != NULL) delete[] _pRecvBuffer;
						_pRecvBuffer = new unsigned char[_iBlockSizeReceive];
					}
				}
			}
	    }
    }
  }
}
