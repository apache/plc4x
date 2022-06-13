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

#include "BLogger.h"
#include <boost/log/utility/setup/common_attributes.hpp>
#include <boost/filesystem.hpp>
#include <boost/filesystem/fstream.hpp>
#include <boost/regex.hpp>
#include <boost/smart_ptr/shared_ptr.hpp>
#include <boost/log/sinks/text_file_backend.hpp>
#include <boost/log/core/core.hpp>
#include <boost/log/attributes.hpp>
#include <boost/log/attributes/scoped_attribute.hpp>
#include <boost/lexical_cast.hpp>
#include <boost/thread.hpp>
#include <boost/log/sinks/text_file_backend.hpp>
#include <boost/log/utility/setup/file.hpp>
#include <boost/log/sources/severity_logger.hpp>
#include <boost/log/sources/record_ostream.hpp>
#include <boost/date_time/posix_time/posix_time.hpp>
#include <boost/log/support/date_time.hpp>

#include <ctime>

namespace org
{
	namespace apache
	{
		namespace plc4x
		{
			namespace cpp
			{
				namespace utils
				{
					namespace fs = boost::filesystem;
					namespace logging = boost::log;
					namespace attrs = boost::log::attributes;

					const char* BLogger::_pDefaultFileTitle = "plc4x";
					const char* BLogger::_pDefaultFileExtension = ".log";

					BLogger* BLogger::m_pInstance = NULL;

					BLogger* BLogger::Instance()
					{
						if (!m_pInstance)
						{
							m_pInstance = new BLogger;
						}

					return m_pInstance;
					}

					void BLogger::Release()
					{
						if (m_pInstance != NULL)
						{
							delete m_pInstance;
							m_pInstance = NULL;
						}

					return;
					}

					BLogger::BLogger() :
						_iMaxFileSize(10 * 1024 * 1024),			/* default max. filesize = 10 MB */
						_uiMaxNumberOfFiles(5),
						_strFileTitle(_pDefaultFileTitle),
						_strFileExtension(_pDefaultFileExtension),
						_logFilter (BoostLogger::trivial::severity >= BoostLogger::trivial::trace),
						_severityLevel(BoostLogger::trivial::trace),
						_uiFlushCounter(0)
					{
						_iSeverityLevelNumber=getNumberfromSeverityLevel(_severityLevel);

						_uiMaxFlushEntrys[0] = 20 ; // trace
						_uiMaxFlushEntrys[1] = 10 ; // debug
						_uiMaxFlushEntrys[2] = 20 ; // info // 0 -> write once
						_uiMaxFlushEntrys[3] = 0 ; // warning // 0 -> write once
						_uiMaxFlushEntrys[4] = 0 ; // error // 0 -> write once
						_uiMaxFlushEntrys[5] = 0 ; // fatal // 0 -> write once
						
                        Initialize();

                        return;
					}

					BLogger::~BLogger()
					{
						_uiFlushCounter = _uiMaxFlushEntrys[_iSeverityLevelNumber];
						FlushLog();

                        return;
					}

					void BLogger::setFileTitle(string strFileTitle)
					{
						if (strFileTitle != _strFileTitle)
						{
							_strFileTitle = strFileTitle;
                        }

                        return;
					}

					void BLogger::setFileExtension(string strFileExtension)
					{
						if (strFileExtension != _strFileExtension)
						{
							_strFileExtension = strFileExtension;
						}

                        return;
					}

					void BLogger::setMaxFileSize(int iMaxFileSize)
					{
						if (iMaxFileSize != _iMaxFileSize)
						{
							_iMaxFileSize = iMaxFileSize; /* MB */
						}

                        return;
					}

					void BLogger::setMaxNumberOfFiles(unsigned int uiMaxNumberOfFiles)
					{
						if (uiMaxNumberOfFiles != _uiMaxNumberOfFiles)
						{
							_uiMaxNumberOfFiles = uiMaxNumberOfFiles;
						}

                        return;
					}

					void BLogger::setSeverityLogLevel(severity_level severityLevel)
					{
						_severityLevel = severityLevel;

						_iSeverityLevelNumber = getNumberfromSeverityLevel(_severityLevel);

						_logFilter = BoostLogger::trivial::severity >= severityLevel;

						BoostLogger::core::get()->set_filter(_logFilter);

                        return;
					}

					
					void BLogger::clearLogDir()
					{
						const boost::regex logFilter( "(.*\\" + _strFileExtension + ")" );

						fs::path phPathFile(_strFileTitle);
						fs::path phDir = phPathFile.parent_path();
						if (fs::exists(phDir))
						{
							fs::directory_iterator it{phDir};
							while (it != boost::filesystem::directory_iterator{})
							{
								if (fs::is_regular_file(*it))
								{
									std::string logFilename = it->path().filename().string();

									if( boost::regex_match( logFilename, logFilter ) )
									{
										boost::filesystem::remove(it->path());
									}
								}
								it++;
							}
						}

                        return;
					}

					void BLogger::Initialize()
					{

						if (_strFileTitle.length() == 0) _strFileTitle = _pDefaultFileTitle;
						if (_strFileExtension.length() == 0) _strFileExtension = _pDefaultFileExtension;
						string strLogFileName = _strFileTitle + "_%m%d%Y_%H%M%S_%N" + _strFileExtension;

						boost::shared_ptr< logging::core > core = logging::core::get();
						core->flush();
						core->remove_all_sinks();
						core->add_global_attribute("TimeStamp", attrs::local_clock());

						boost::shared_ptr< sinks::text_file_backend > backend =
								boost::make_shared< sinks::text_file_backend >
						(
								keywords::file_name = strLogFileName.c_str(),
								keywords::rotation_size = _iMaxFileSize * 1024 * 1024,
								keywords::time_based_rotation = sinks::file::rotation_at_time_point(0, 0, 0)
						);

						fs::path phDirPath(strLogFileName);

						backend->set_file_collector(sinks::file::make_collector
								(
									keywords::target = phDirPath.parent_path().c_str(),
									keywords::max_size = (_uiMaxNumberOfFiles - 1)  * _iMaxFileSize * 1024 * 1024,
									keywords::min_free_space = 50 * 1024 * 1024
								));

						typedef sinks::synchronous_sink< sinks::text_file_backend > sink_t;
						boost::shared_ptr< sink_t > sink(new sink_t(backend));

						sink->set_formatter
							(
									expr::stream
									<< expr::format_date_time< boost::posix_time::ptime >("TimeStamp", "[%d.%m. %H:%M:%S.%f] ")
									<< " <" << BoostLogger::trivial::severity
									<< "> : " << expr::message
							);

						core->add_sink(sink);

                        return;
					}

					void  BLogger::Log(const char* pMessage, severity_level logLevel /*= info*/)
					{
						try
						{
							if (isLog(logLevel))
							{
								_mtxLock.lock();

								BOOST_LOG_SEV(_lg, severity_level::info) << "ThreadID 0x" << boost::this_thread::get_id() << " " << pMessage;

								if(logLevel == severity_level::fatal)
									ForceFlushLog();
								else
									FlushLog();

								_mtxLock.unlock();
							}
						}
						catch(...)
						{
							_mtxLock.unlock();
						}

                        return;
					}

					void  BLogger::Log(string strMessage, severity_level logLevel /* = info */)
					{
						try
						{
							if (isLog(logLevel))
							{
								_mtxLock.lock();

								BOOST_LOG_SEV(_lg, logLevel) << "ThreadID 0x" << boost::this_thread::get_id() << " " << strMessage;

								if(logLevel == severity_level::fatal)
									ForceFlushLog();
								else
									FlushLog();

								_mtxLock.unlock();
							}
						}
						catch(...)
						{
							_mtxLock.unlock();
						}

                        return;
					}

					void  BLogger::LogBool(bool b, severity_level logLevel /*= info*/)
					{

						try
						{
							if (isLog(logLevel))
							{
								_mtxLock.lock();

								string strBoolMessage = b ? _strTrueMessage : _strFalseMessage;
								BOOST_LOG_SEV(_lg, logLevel) << "ThreadID 0x" << boost::this_thread::get_id() << " " << strBoolMessage;

								if(logLevel == severity_level::fatal)
									ForceFlushLog();
								else
									FlushLog();

								_mtxLock.unlock();
							}
						}
						catch(...)
						{
							_mtxLock.unlock();
						}

                        return;
					}

					void  BLogger::LogBool(string strMessage, bool b, severity_level logLevel /*= info*/)
					{
						try
						{
							if (isLog(logLevel))
							{
								_mtxLock.lock();

								string strBoolMessage = b ? _strTrueMessage : _strFalseMessage;
								BOOST_LOG_SEV(_lg, logLevel) << "ThreadID 0x" << boost::this_thread::get_id() << " " << strMessage << ' ' << strBoolMessage;

								if(logLevel == severity_level::fatal)
									ForceFlushLog();
								else
									FlushLog();

								_mtxLock.unlock();
							}
						}
						catch(...)
						{
							_mtxLock.unlock();
						}

                        return;
					}

					void BLogger::LogHexDump(unsigned char* pBinaryData, unsigned int uiNumBytes, const char* pDumpPrefix /*= NULL*/, severity_level logLevel /*= info*/)
					{
						try
						{
							if (isLog(logLevel))
							{
								_mtxLock.lock();

								DatDmp dump;
								dump.setLeadingNewLine(true);
								dump.setPrefix(pDumpPrefix);
								string strDumpMessage = dump.Show(pBinaryData, uiNumBytes);
								BOOST_LOG_SEV(_lg, logLevel) << "ThreadID 0x" << boost::this_thread::get_id() << " " << strDumpMessage;

                                if (logLevel == severity_level::fatal)
                                {
                                    ForceFlushLog();
                                }
                                else
                                {
                                    FlushLog();
                                }

								_mtxLock.unlock();
							}
						}
						catch(...)
						{
							_mtxLock.unlock();
						}

                        return;
					}

					void BLogger::FlushLog(void)
					{
						if(_uiFlushCounter >= _uiMaxFlushEntrys[_iSeverityLevelNumber])
						{
							BOOST_LOG_SEV(_lg, boost::log::trivial::severity_level::info) << "ThreadID 0x" << boost::this_thread::get_id() << " Flush LogEntries";
							logging::core::get()->flush();
							_uiFlushCounter = 0;
						}
						else
						{
							_uiFlushCounter++;
						}

                        return;
					}

					void BLogger::ForceFlushLog(void)
					{
						unsigned int ui = _uiMaxFlushEntrys[_iSeverityLevelNumber];
						_uiMaxFlushEntrys[_iSeverityLevelNumber] = 0;
						FlushLog();
						_uiMaxFlushEntrys[_iSeverityLevelNumber] = ui;

                        return;
					}

					void BLogger::setSeverityLogLevelNumber(int iSeverityLevelNumber)
					{

						_iSeverityLevelNumber = iSeverityLevelNumber;

						_severityLevel = getSeverityLevelfromNumber(_iSeverityLevelNumber);

						setSeverityLogLevel(_severityLevel);

                        return;
					}

					bool BLogger::isLog(severity_level logLevel)
					{
						bool bResult = false;

						try
						{
							int ilogLevel = getNumberfromSeverityLevel(logLevel);

							if (_iSeverityLevelNumber <= ilogLevel)
							{
								bResult = true;
							}
						}
						catch(...)
						{
						}

						return bResult;
					}

					severity_level BLogger::getSeverityLevelfromNumber(int iSeverityLevelNumber)
					{
						severity_level logLevel = boost::log::trivial::severity_level::trace;

						switch (iSeverityLevelNumber)
						{
							case 0: logLevel = boost::log::trivial::severity_level::trace; break;
							case 1: logLevel = boost::log::trivial::severity_level::debug; break;
							case 2: logLevel = boost::log::trivial::severity_level::info; break;
							case 3: logLevel = boost::log::trivial::severity_level::warning; break;
							case 4: logLevel = boost::log::trivial::severity_level::error; break;
							case 5: logLevel = boost::log::trivial::severity_level::fatal; break;
							default: logLevel = boost::log::trivial::severity_level::trace; break;
						}

						return logLevel;
					}

					int BLogger::getNumberfromSeverityLevel(severity_level severityLevel)
					{
						int iSeverityLevelNumber = 0;

						switch (severityLevel)
						{
							case boost::log::trivial::severity_level::trace: iSeverityLevelNumber = 0; break;
							case boost::log::trivial::severity_level::debug: iSeverityLevelNumber = 1; break;
							case boost::log::trivial::severity_level::info: iSeverityLevelNumber = 2; break;
							case boost::log::trivial::severity_level::warning: iSeverityLevelNumber = 3; break;
							case boost::log::trivial::severity_level::error: iSeverityLevelNumber = 4; break;
							case boost::log::trivial::severity_level::fatal: iSeverityLevelNumber = 5; break;
						}

						return iSeverityLevelNumber;
					}
				}
			}
		}
	}
}
