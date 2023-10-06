#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#
from abc import ABC, abstractmethod
from dataclasses import dataclass
from typing import TypeVar, Generic, Callable, Any, List

from plc4py.api.listener import EventListener
from plc4py.spi.Plc4xBaseProtocol import Plc4xBaseProtocol
from plc4py.spi.context.DriverContext import DriverContext
from plc4py.spi.messages.ChannelMessage import ChannelMessage
from plc4py.utils.GenericTypes import ByteOrder


class Configuration:
    pass


class ChannelPipeline:
    pass


T = TypeVar("T", bound=ChannelMessage)


class ProtocolStackConfigurator(Generic[T], ABC):
    """
    Provides an easy interface for configuring protocol stacks
    """

    @abstractmethod
    def configure_pipeline(
        self,
        configuration: Configuration,
        pipeline: ChannelPipeline,
        passive: bool = False,
        listeners: List[EventListener] = [],
    ) -> Plc4xBaseProtocol:
        pass


U = TypeVar("U", bound=ChannelMessage)


@dataclass
class SingleProtocolStackConfigurator(ProtocolStackConfigurator[U], ABC):
    base_message_class: U
    byte_order: ByteOrder
    protocol_class: Plc4xBaseProtocol
    driver_context_class: DriverContext
    packet_size_estimate: Callable[..., int]
    corrupt_packet_remover: Callable[..., Any]

    def configure_pipeline(
        self,
        configuration: Configuration,
        pipeline: ChannelPipeline,
        passive: bool = False,
        listeners: List[EventListener] = [],
    ) -> Plc4xBaseProtocol:
        pass
