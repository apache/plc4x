package readwrite

import (

)

type Message interface {
    lengthInBytes() int16
    lengthInBits() int16
}