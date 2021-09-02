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

#include "BoostConnection.h"
#include <org/apache/plc4x/cpp/api/exceptions/PlcConnectionException.h>
#include <org/apache/plc4x/cpp/utils/logger/DbgTrace.h>
#include <org/apache/plc4x/cpp/utils/logger/BLogger.h>
#include <org/apache/plc4x/cpp/utils/logger/ExLog.h>
#include <boost/lexical_cast.hpp>
#include <boost/format.hpp>

#define CON_DEFAULT_IPADDRESS   "127.0.0.1"
#define CON_DEFAULT_PORT_       0
#define CON_DEFAULT_TIMEOUT     2000 // Milliseconds

namespace org
{
    namespace apache
    {
        namespace plc4x
        {
            namespace cpp
            {
                namespace base
                {
                    namespace connection
                    {

                        using namespace boost;
                        using namespace boost::signals2;
                        using namespace org::apache::plc4x::cpp::utils;
                        using namespace org::apache::plc4x::cpp::api::exceptions;

                        BoostConnection::BoostConnection() :
                            _pSocket(NULL),
                            _pWorker(NULL),
                            _pReceiveDispatcher_Thread(NULL),
                            _isReceiving(false),
                            _pClientThread(NULL),
                            _uiConnectionTimeout(CON_DEFAULT_TIMEOUT)
                        {
                        }

                        BoostConnection::~BoostConnection()
                        {
                            if (_pSocket != NULL)
                            {
                                try
                                {
                                    this->close();
                                }
                                catch (std::exception& ex)
                                {
                                    ExLog(1, ex, "BoostConnection::~BoostConnection() failed");
                                }
                            }
                        }

                        void BoostConnection::setIPAddress(string strIPAddress)
                        {
                            if (strIPAddress != _strIPAddress)
                            {
                                if (_bConnected)
                                {
                                    _strErrorMessage = "Attempt to change the IP-Address " + strIPAddress + "of an already open connection";
                                    BOOST_THROW_EXCEPTION(PlcConnectionException(_strErrorMessage));
                                }
                                else
                                {
                                    try
                                    {
                                        if ((strIPAddress.length() >= 7) && (strIPAddress != _strIPAddress))
                                        {
                                            _strIPAddress = strIPAddress;
                                            resolveEndpoint();

                                            ip::tcp::endpoint ep(ip::address::from_string(strIPAddress), _usPort);
                                            _endpoint = ep;
                                        }
                                    }
                                    catch (std::exception& ex)
                                    {
                                        BOOST_THROW_EXCEPTION(ExLog(1, ex, "BoostConnection::setIP_Address(string strIP_Address) failed"));
                                    }
                                }
                            }
                        }

                        void BoostConnection::setPort(unsigned short usPort)
                        {
                            if (_usPort != usPort)
                            {
                                if (_bConnected)
                                {
                                    BOOST_THROW_EXCEPTION(PlcConnectionException("Attempt to change the port of an already open communication"));
                                }
                                else
                                {
                                    try
                                    {
                                        _usPort = usPort;
                                        resolveEndpoint();

                                        ip::tcp::endpoint ep(ip::address::from_string(_strIPAddress), _usPort);
                                        _endpoint = ep;
                                    }
                                    catch (std::exception& ex)
                                    {
                                        BOOST_THROW_EXCEPTION(ExLog(1, ex, "BoostConnection::setPort(unsigned short port) failed"));
                                    }
                                }
                            }
                        }

                        void BoostConnection::setEndpoint(tcp::endpoint endpoint)
                        {
                            if (endpoint != _endpoint)
                            {
                                if (_bConnected)
                                {
                                    _strErrorMessage = "Attempt to change the endpoint of an already open communication";
                                    _strErrorMessage = _strErrorMessage + "Setting Endpoint " + endpoint.address().to_string() + ":" + boost::lexical_cast<string>((unsigned int)endpoint.port()) + " for BoostConnection object failed";
                                    BOOST_THROW_EXCEPTION(PlcConnectionException(_strErrorMessage));
                                }
                                else
                                {
                                    try
                                    {
                                        _endpoint = endpoint;
                                        resolveEndpoint();
                                    }
                                    catch (std::exception& ex)
                                    {
                                        BOOST_THROW_EXCEPTION(ExLog(1, ex, "BoostConnection::setEndpoint(tcp::endpoint endpoint) failed"));
                                    }
                                }
                            }
                        }

                        /**---------------------------------------------------------------------
                            * Service Callback receive TCP Packets
                            *---------------------------------------------------------------------*/
                        void BoostConnection::dispatchReceive()
                        {
                            _strErrorMessage = "Attempt to use BoostConnection::dispatchReceive of BoostConnection";
                            LOG_INFO(_strErrorMessage); DEBUG_TRACE1("%s\n", _strErrorMessage.c_str());
                            BOOST_THROW_EXCEPTION(PlcConnectionException(_strErrorMessage));
                        }

                        /**---------------------------------------------------------------------
                        Connect TCP Socket    
                        *---------------------------------------------------------------------*/
                        void BoostConnection::connect()
                        {
                            try
                            {

                                resolveEndpoint();

                                _endpoint = *_iterEndPoint;

                                BoostConnection::close();

                                _io_service.reset();

                                if (_pSocket == NULL) _pSocket = new tcp::socket(_io_service);

                                _pSocket->async_connect(_endpoint, boost::bind(&BoostConnection::onConnect, this, asio::placeholders::error, _iterEndPoint));

                                _pWorker = new io_service::work(_io_service);
                                _pClientThread = new boost::thread(boost::bind(&asio::io_service::run, &_io_service));

                                boost::mutex::scoped_lock lock(_mtxComLock);
                                _condComHandler.timed_wait(lock, boost::posix_time::milliseconds(_uiConnectionTimeout));

                                LOG_BOOL("Open BoostCommunication", _bConnected);

                                if (_bConnected)
                                {
                                    _pReceiveDispatcher_Thread = new boost::thread(boost::bind(&BoostConnection::dispatchReceive, this));
                                }
                                else
                                {
                                    BOOST_THROW_EXCEPTION(PlcConnectionException("Timeout while opening BoostConnection"));
                                }

                            }
                            catch (std::exception& ex)
                            {
                                BoostConnection::close();

                                BOOST_THROW_EXCEPTION(ExLog(1, ex, "BoostConnection::connect() failed"));
                            }
                        }

                        /**---------------------------------------------------------------------
                          Close TCP Socket
                        *---------------------------------------------------------------------*/
                        void BoostConnection::close()
                        {

                            try
                            {
                                _bConnected = false;

                                if ((_pClientThread != NULL) && (boost::this_thread::get_id() != _pClientThread->get_id()))
                                {
                                    if (_pSocket != NULL)
                                    {
                                        _io_service.post(boost::bind(&asio::ip::tcp::socket::close, _pSocket));
                                    }
                                    _io_service.post(boost::bind(&io_service::stop, &_io_service));

                                }
                                else
                                {
                                    if (_pSocket != NULL)
                                    {
                                        _pSocket->close();
                                    }
                                    _io_service.stop();
                                }

                                clearRessources();

                            }
                            catch (std::exception& ex)
                            {
                                clearRessources();

                                _bConnected = false;

                                BOOST_THROW_EXCEPTION(ExLog(1, ex, "BoostConnection::close() failed"));
                            }

                            return;
                        }

                        /**---------------------------------------------------------------------
                        Send Bytes via TCP Socket
                            *---------------------------------------------------------------------*/
                        bool BoostConnection::send(unsigned char* pBytesToSend, int iBytesToSend)
                        {
                            bool bSuccess = false;

                            boost::lock_guard<boost::mutex> lock(_mutexSend);

                            if (_bConnected)
                            {
                                try
                                {
                                    _pSocket->async_send(asio::buffer(pBytesToSend, iBytesToSend),
                                        boost::bind(&BoostConnection::onSend, this, boost::asio::placeholders::error, boost::asio::placeholders::bytes_transferred)
                                    );

                                    _dump.setLeadingNewLine(true);
                                    std::string strTelegram = "\nBoostConnection TCP Telegram Send" + _dump.Show(pBytesToSend, iBytesToSend, "  ** ");
                                    LOG_TRACE(strTelegram.c_str()); DEBUG_TRACE1("%s", strTelegram.c_str());

                                    bSuccess = true;
                                }
                                catch (std::exception& ex)
                                {
                                    bSuccess = false;
                                    BOOST_THROW_EXCEPTION(ExLog(1, ex, "BoostConnection::send() failed"));
                                }
                            }
                            return bSuccess;
                        }


                        void BoostConnection::ping()
                        {
                           try 
                           {
                               /*(Socket s = new Socket())
                               s.connect(address, PING_TIMEOUT_MS);
                               // TODO keep the address for a (timely) next request???
                               s.setReuseAddress(true);*/
                           }
                           catch (std::exception ex) 
                           {
                               BOOST_THROW_EXCEPTION(ExLog(1, ex, "Unable to ping remote host"));
                           }
                        }
                        /**---------------------------------------------------------------------
                        Free Ressources
                        *---------------------------------------------------------------------*/
                        void BoostConnection::clearRessources()
                        {
                            try
                            {

                                _io_service.reset();

                                if (_pReceiveDispatcher_Thread != NULL)
                                {
                                    if (boost::this_thread::get_id() != _pReceiveDispatcher_Thread->get_id())
                                    {
                                        _pReceiveDispatcher_Thread->interrupt();
                                        _pReceiveDispatcher_Thread->join();
                                        delete _pReceiveDispatcher_Thread; _pReceiveDispatcher_Thread = NULL;
                                    }
                                }

                                if (_pClientThread != NULL)
                                {
                                    if (boost::this_thread::get_id() != _pClientThread->get_id())
                                    {
                                        _pClientThread->interrupt();
                                        _pClientThread->join();
                                        delete _pClientThread; _pClientThread = NULL;
                                    }
                                }

                                if (_pWorker != NULL)
                                {
                                    delete _pWorker; _pWorker = NULL;
                                }

                                if (_pSocket != NULL)
                                {
                                    delete _pSocket; _pSocket = NULL;
                                }

                            }
                            catch (std::exception&)
                            {

                            }

                            return;
                        }

                        /**---------------------------------------------------------------------
                        Recieve next TCP Packets 
                        *---------------------------------------------------------------------*/
                        void BoostConnection::continueReceive()
                        {
                            _pSocket->async_receive(
                                asio::buffer(_pRecvBuffer, _iRecvBufferSize),
                                bind(&BoostConnection::onReceive, this, asio::placeholders::error, asio::placeholders::bytes_transferred)
                            );
                        }

                        /**---------------------------------------------------------------------
                            Convert endpoint to String
                        *---------------------------------------------------------------------*/
                        string BoostConnection::endPointToString(tcp::endpoint ep)
                        {
                            return string(ep.address().to_string() + ":" + boost::lexical_cast<std::string>((unsigned short)ep.port()));
                        }

                        /**---------------------------------------------------------------------
                        Service Callback for Async-Send Result by OnSend
                        *---------------------------------------------------------------------*/
                        void BoostConnection::onSend(const system::error_code& errorCode, std::size_t bytes_transferred)
                        {
                            if (errorCode != system::errc::success)
                            {
                                int iErrCode = errorCode.value();
                                /*if (iErrCode != NO_ERROR)
                                {
                                    // WSAECONNRESET = 10054, WSAECONNABORTED = 10053, EPIPE, ECONNABORTED indicates "broken Connection"
                                    if (iErrCode == WSAECONNRESET || iErrCode == WSAECONNABORTED || iErrCode == EPIPE)
                                    {
                                        BoostConnection::close();
                                    }

                                    std::stringstream stm;
                                    stm << boost::format("BoostConnection Send Error %i: %s") % iErrCode % errorCode.message();
                                    _strErrorMessage = stm.str();
                                    LOG_TRACE(_strErrorMessage.c_str()); DEBUG_TRACE1("%s", _strErrorMessage.c_str());
                                    // TODO: Signal
                                }*/
                            }
                        }

                        /**---------------------------------------------------------------------
                        Service Callback , Receive Data via TCP
                       *---------------------------------------------------------------------*/
                        void BoostConnection::onReceive(const system::error_code& errorCode, std::size_t bytes_transferred)
                        {
                            DEBUG_TRACE0("BoostConnection::onReceive() called\n");
                            if (errorCode != system::errc::success && errorCode.value() != ENOENT)
                            {
                                int iErrCode = errorCode.value();
                                /*if (iErrCode != NO_ERROR)
                                {
                                    // WSAECONNRESET = 10054, WSAECONNABORTED = 10053, EPIPE, ECONNABORTED indicates "broken Connection"
                                    if (iErrCode == WSAECONNRESET || iErrCode == WSAECONNABORTED || iErrCode == EPIPE)
                                    {
                                        BoostConnection::close();
                                    }

                                    std::stringstream stm;
                                    stm << boost::format("BoostConnection Receive Error %i: %s") % iErrCode % errorCode.message();
                                    _strErrorMessage = stm.str();
                                    LOG_TRACE(_strErrorMessage.c_str()); DEBUG_TRACE1("%s", _strErrorMessage.c_str());
                                    // TODO: Signal
                                }*/
                            }
                        }

                        /**---------------------------------------------------------------------
                        Resolve Endpoint
                        *---------------------------------------------------------------------*/
                        void BoostConnection::resolveEndpoint()
                        {
                            try
                            {
                                tcp::resolver Resolver(_io_service);
                                tcp::resolver::query EndpointQuery(_strIPAddress, std::to_string(_usPort));
                                _iterEndPoint = Resolver.resolve(EndpointQuery);
                            }
                            catch (std::exception& ex)
                            {
                                BOOST_THROW_EXCEPTION(ExLog(1, ex, "BoostConnection::resolveEndpoint() failed"));
                            }
                        }

                        /**---------------------------------------------------------------------
                        Service Callback open Connection Result
                        *---------------------------------------------------------------------*/
                        void BoostConnection::onConnect(const system::error_code& errorCode, tcp::resolver::iterator endpointIterator)
                        {

                            // Setzen des IsOpen auf true oder false 
                            _bConnected = errorCode == system::errc::success;

                            ++endpointIterator;

                            if (!_bConnected && endpointIterator != tcp::resolver::iterator())	// connection failed, ..
                            {
                                try
                                {
                                    // .. try the next endpoint in the list.
                                    _pSocket->close();
                                    _endpoint = *endpointIterator;
                                    _pSocket->async_connect(_endpoint, boost::bind(&BoostConnection::onConnect, this, asio::placeholders::error, endpointIterator));
                                }
                                catch (std::exception& ex)
                                {
                                    ExLog(1, ex, "BoostConnection::onConnect() Handler failed");
                                }
                            }
                            else if (errorCode == system::errc::success)
                            {
                                bool bServiceStopped = _io_service.stopped();

                                try
                                {
                                    _pSocket->async_receive(
                                        asio::buffer(_pRecvBuffer, _iRecvBufferSize),
                                        bind(&BoostConnection::onReceive, this, asio::placeholders::error, asio::placeholders::bytes_transferred)
                                    );

                                }
                                catch (std::exception& ex)
                                {
                                    ExLog(1, ex, "BoostConnection::onConnect() async_receive failed");
                                }

                                if (bServiceStopped)
                                {
                                    LOG_INFO("Restarting Client Thread for BoostConnection OnConnect");

                                    try
                                    {
                                        if (_pClientThread != NULL) delete _pClientThread;
                                        _pClientThread = new boost::thread(bind(&asio::io_service::run, &_io_service));
                                    }
                                    catch (std::exception& ex)
                                    {
                                        ExLog(1, ex, "BoostConnection::onConnect() failed");
                                    }
                                }
                            }

                            _condComHandler.notify_all();

                        }
                    }
                }
            }
        }
    }
}
