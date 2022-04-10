from abc import abstractmethod
from dataclasses import dataclass
from typing import Type

import pluggy

import plc4py
from plc4py.api.PlcConnection import PlcConnection
from plc4py.api.messages.PlcRequest import ReadRequestBuilder


@dataclass
class ModbusConnection(PlcConnection):
    """A hook implementation namespace."""

    def connect(self):
        """
        Establishes the connection to the remote PLC.
        """
        pass

    def is_connected(self) -> bool:
        """
        Indicates if the connection is established to a remote PLC.
        :return: True if connection, False otherwise
        """
        pass

    def close(self) -> None:
        """
        Closes the connection to the remote PLC.
        :return:
        """
        pass

    def read_request_builder(self) -> ReadRequestBuilder:
        """
        :return: read request builder.
        """
        pass


class ModbusConnectionLoader:

    @staticmethod
    @plc4py.hookimpl
    def get_connection() -> Type[ModbusConnection]:
        return ModbusConnection

    @staticmethod
    @plc4py.hookimpl
    def key() -> str:
        return "modbus"

