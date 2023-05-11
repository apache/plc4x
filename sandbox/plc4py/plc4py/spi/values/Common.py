from plc4py.spi.generation import WriteBuffer


class Serializable:

    def serialize(self, write_buffer: WriteBuffer):
        pass