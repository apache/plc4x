/*
 * StandardFirmata sketch for VirtualAVR integration testing.
 * This is a minimal Firmata implementation that works with avr8js simulator.
 *
 * Based on StandardFirmata by Firmata (LGPL v2.1+)
 * Simplified for testing purposes.
 *
 * Debug output: The built-in LED (pin 13) blinks when Firmata commands are received.
 * Analog changes are indicated by PWM on pin 9 (proportional to A0 value).
 */

#include <Firmata.h>

/* Debug LED for visual feedback */
#define DEBUG_LED 13
#define DEBUG_PWM 9

/* Command counters for debugging */
volatile unsigned int commandCount = 0;
volatile unsigned int analogChangeCount = 0;

/* analog inputs */
int analogInputsToReport = 0; // bitwise array to store pin reporting
int lastAnalogValues[16];

/* digital input ports */
byte reportPINs[TOTAL_PORTS];       // 1 = report this port, 0 = silence
byte previousPINs[TOTAL_PORTS];     // previous 8 bits sent

/* pins configuration */
byte pinConfig[TOTAL_PINS];         // configuration of every pin
byte portConfigInputs[TOTAL_PORTS]; // each bit: 1 = pin in INPUT, 0 = anything else

/* timer variables */
unsigned long currentMillis;        // store the current value from millis()
unsigned long previousMillis;       // for comparison with currentMillis
unsigned int samplingInterval = 19; // how often to sample analog inputs (in ms)

/* Debug: Toggle LED to indicate activity */
void debugBlink() {
  static bool ledState = false;
  ledState = !ledState;
  digitalWrite(DEBUG_LED, ledState ? HIGH : LOW);
  commandCount++;
}

/* Debug: Show analog value on PWM pin */
void debugAnalog(int value) {
  // Map 0-1023 to 0-255 for PWM
  analogWrite(DEBUG_PWM, value >> 2);
  analogChangeCount++;
}

void outputPort(byte portNumber, byte portValue, byte forceSend)
{
  portValue = portValue & portConfigInputs[portNumber];
  if (forceSend || previousPINs[portNumber] != portValue) {
    Firmata.sendDigitalPort(portNumber, portValue);
    previousPINs[portNumber] = portValue;
  }
}

void checkDigitalInputs(void)
{
  if (TOTAL_PORTS > 0 && reportPINs[0]) outputPort(0, readPort(0, portConfigInputs[0]), false);
  if (TOTAL_PORTS > 1 && reportPINs[1]) outputPort(1, readPort(1, portConfigInputs[1]), false);
  if (TOTAL_PORTS > 2 && reportPINs[2]) outputPort(2, readPort(2, portConfigInputs[2]), false);
}

void setPinModeCallback(byte pin, int mode)
{
  debugBlink();  // Visual feedback for received command
  if (pin < TOTAL_PINS) {
    pinConfig[pin] = mode;

    if (mode == PIN_MODE_INPUT || mode == PIN_MODE_PULLUP) {
      portConfigInputs[pin / 8] |= (1 << (pin & 7));
      if (mode == PIN_MODE_PULLUP) {
        digitalWrite(pin, HIGH);
      }
      pinMode(pin, INPUT);
    } else if (mode == PIN_MODE_OUTPUT) {
      portConfigInputs[pin / 8] &= ~(1 << (pin & 7));
      pinMode(pin, OUTPUT);
    } else if (mode == PIN_MODE_ANALOG) {
      portConfigInputs[pin / 8] &= ~(1 << (pin & 7));
    } else if (mode == PIN_MODE_PWM) {
      portConfigInputs[pin / 8] &= ~(1 << (pin & 7));
      pinMode(pin, OUTPUT);
    }
  }
}

void setDigitalPinValueCallback(byte pin, int value)
{
  debugBlink();  // Visual feedback for received command
  if (pin < TOTAL_PINS && pinConfig[pin] == PIN_MODE_OUTPUT) {
    digitalWrite(pin, value);
  }
}

void analogWriteCallback(byte pin, int value)
{
  debugBlink();  // Visual feedback for received command
  if (pin < TOTAL_PINS) {
    switch (pinConfig[pin]) {
      case PIN_MODE_PWM:
        analogWrite(pin, value);
        break;
    }
  }
}

void digitalWriteCallback(byte port, int value)
{
  debugBlink();  // Visual feedback for received command
  byte pin, lastPin, mask = 1, pinWriteMask = 0;

  if (port < TOTAL_PORTS) {
    lastPin = port * 8 + 8;
    if (lastPin > TOTAL_PINS) lastPin = TOTAL_PINS;

    for (pin = port * 8; pin < lastPin; pin++) {
      if (pinConfig[pin] == PIN_MODE_OUTPUT) {
        pinWriteMask |= mask;
      }
      mask = mask << 1;
    }
    writePort(port, (byte)value, pinWriteMask);
  }
}

void reportAnalogCallback(byte analogPin, int value)
{
  debugBlink();  // Visual feedback for received command
  if (analogPin < TOTAL_ANALOG_PINS) {
    if (value == 0) {
      analogInputsToReport &= ~(1 << analogPin);
    } else {
      analogInputsToReport |= (1 << analogPin);
    }
  }
}

void reportDigitalCallback(byte port, int value)
{
  debugBlink();  // Visual feedback for received command
  if (port < TOTAL_PORTS) {
    reportPINs[port] = (byte)value;
    if (value) outputPort(port, readPort(port, portConfigInputs[port]), true);
  }
}

void sysexCallback(byte command, byte argc, byte *argv)
{
  switch (command) {
    case SAMPLING_INTERVAL:
      if (argc > 1) {
        samplingInterval = argv[0] + (argv[1] << 7);
        if (samplingInterval < 10) {
          samplingInterval = 10;
        }
      }
      break;
    case CAPABILITY_QUERY:
      Firmata.write(START_SYSEX);
      Firmata.write(CAPABILITY_RESPONSE);
      for (byte pin = 0; pin < TOTAL_PINS; pin++) {
        if (IS_PIN_DIGITAL(pin)) {
          Firmata.write((byte)PIN_MODE_INPUT);
          Firmata.write(1);
          Firmata.write((byte)PIN_MODE_PULLUP);
          Firmata.write(1);
          Firmata.write((byte)PIN_MODE_OUTPUT);
          Firmata.write(1);
        }
        if (IS_PIN_ANALOG(pin)) {
          Firmata.write(PIN_MODE_ANALOG);
          Firmata.write(10); // 10-bit resolution
        }
        if (IS_PIN_PWM(pin)) {
          Firmata.write(PIN_MODE_PWM);
          Firmata.write(8); // 8-bit resolution
        }
        Firmata.write(127);
      }
      Firmata.write(END_SYSEX);
      break;
    case PIN_STATE_QUERY:
      if (argc > 0) {
        byte pin = argv[0];
        Firmata.write(START_SYSEX);
        Firmata.write(PIN_STATE_RESPONSE);
        Firmata.write(pin);
        if (pin < TOTAL_PINS) {
          Firmata.write((byte)pinConfig[pin]);
          if (pinConfig[pin] == PIN_MODE_OUTPUT) {
            Firmata.write((byte)digitalRead(pin));
          }
        }
        Firmata.write(END_SYSEX);
      }
      break;
    case ANALOG_MAPPING_QUERY:
      Firmata.write(START_SYSEX);
      Firmata.write(ANALOG_MAPPING_RESPONSE);
      for (byte pin = 0; pin < TOTAL_PINS; pin++) {
        Firmata.write(IS_PIN_ANALOG(pin) ? PIN_TO_ANALOG(pin) : 127);
      }
      Firmata.write(END_SYSEX);
      break;
  }
}

void systemResetCallback()
{
  for (byte i = 0; i < TOTAL_PORTS; i++) {
    reportPINs[i] = false;
    portConfigInputs[i] = 0;
    previousPINs[i] = 0;
  }
  for (byte i = 0; i < TOTAL_PINS; i++) {
    pinConfig[i] = PIN_MODE_OUTPUT;
  }
  for (byte i = 0; i < TOTAL_ANALOG_PINS; i++) {
    lastAnalogValues[i] = 0;
  }
  analogInputsToReport = 0;
}

void setup()
{
  // Initialize debug pins
  pinMode(DEBUG_LED, OUTPUT);
  pinMode(DEBUG_PWM, OUTPUT);
  digitalWrite(DEBUG_LED, LOW);

  // Blink LED 3 times to indicate startup
  for (int i = 0; i < 3; i++) {
    digitalWrite(DEBUG_LED, HIGH);
    delay(100);
    digitalWrite(DEBUG_LED, LOW);
    delay(100);
  }

  Firmata.setFirmwareVersion(FIRMATA_FIRMWARE_MAJOR_VERSION, FIRMATA_FIRMWARE_MINOR_VERSION);

  Firmata.attach(ANALOG_MESSAGE, analogWriteCallback);
  Firmata.attach(DIGITAL_MESSAGE, digitalWriteCallback);
  Firmata.attach(REPORT_ANALOG, reportAnalogCallback);
  Firmata.attach(REPORT_DIGITAL, reportDigitalCallback);
  Firmata.attach(SET_PIN_MODE, setPinModeCallback);
  Firmata.attach(SET_DIGITAL_PIN_VALUE, setDigitalPinValueCallback);
  Firmata.attach(START_SYSEX, sysexCallback);
  Firmata.attach(SYSTEM_RESET, systemResetCallback);

  Firmata.begin(57600);
  systemResetCallback();
}

void loop()
{
  byte pin, analogPin;

  checkDigitalInputs();

  currentMillis = millis();
  if (currentMillis - previousMillis >= samplingInterval) {
    previousMillis = currentMillis;

    for (pin = 0; pin < TOTAL_ANALOG_PINS; pin++) {
      if (analogInputsToReport & (1 << pin)) {
        int value = analogRead(pin);
        if (value != lastAnalogValues[pin]) {
          Firmata.sendAnalog(pin, value);
          lastAnalogValues[pin] = value;
          // Debug: Show analog value changes on PWM pin 9
          if (pin == 0) {
            debugAnalog(value);
          }
        }
      }
    }
  }

  while (Firmata.available()) {
    Firmata.processInput();
  }
}
