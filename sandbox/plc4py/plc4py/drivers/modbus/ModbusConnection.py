from dataclasses import dataclass
from typing import Type

from plc4py import hookimpl
from plc4py.api.PlcConnection import PlcConnection


@dataclass
class ModbusConnection(PlcConnection):
    """A hook implementation namespace."""


class ModbusConnectionLoader:

    @staticmethod
    @hookimpl
    def get_type() -> Type[ModbusConnection]:
        return ModbusConnection

    @staticmethod
    @hookimpl
    def key() -> str:
        return "modbus"

