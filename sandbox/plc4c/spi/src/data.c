
#include <assert.h>
#include <plc4c/data.h>
#include <plc4c/spi/types_private.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>

plc4c_data *plc4c_data_create_boolean_data(bool b) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_BOOL;
  data->size = sizeof(bool);
  data->data.boolean_value = b;
  return data;
}

plc4c_data *plc4c_data_create_char_data(char c) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_CHAR;
  data->size = sizeof(char);
  data->data.char_value = c;
  return data;
}

plc4c_data *plc4c_data_create_uchar_data(unsigned char uc) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_UCHAR;
  data->size = sizeof(unsigned char);
  data->data.uchar_value = uc;
  return data;
}

plc4c_data *plc4c_data_create_short_data(short s) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_SHORT;
  data->size = sizeof(short);
  data->data.short_value = s;
  return data;
}

plc4c_data *plc4c_data_create_ushort_data(unsigned short us) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_USHORT;
  data->size = sizeof(unsigned short);
  data->data.ushort_value = us;
  return data;
}

plc4c_data *plc4c_data_create_int_data(int i) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_INT;
  data->size = sizeof(int);
  data->data.int_value = i;
  return data;
}

plc4c_data *plc4c_data_create_uint_data(unsigned int ui) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_UINT;
  data->size = sizeof(unsigned int);
  data->data.uint_value = ui;
  return data;
}

plc4c_data *plc4c_data_create_string_data(unsigned int size, char *s) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_STRING_POINTER;
  data->size = size;
  data->data.pstring_value = s;
  return data;
}

plc4c_data *plc4c_data_create_constant_string_data(unsigned int size, char *s) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_CONSTANT_STRING;
  data->size = size;
  data->data.const_string_value = s;
  return data;
}

plc4c_data *plc4c_data_create_void_pointer_data(void *v) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_VOID_POINTER;
  data->size = 0;
  data->data.pvoid_value = v;
  return data;
}

plc4c_data *plc4c_data_create_float_data(float f) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_FLOAT;
  data->size = sizeof(float);
  data->data.float_value = f;
}

void plc4c_data_printf(plc4c_data *data) {
  switch (data->data_type) {
    case PLC4C_CHAR:
      printf("%c", data->data.char_value);
      break;
    case PLC4C_UCHAR:
      printf("%c", data->data.uchar_value);
      break;
    case PLC4C_SHORT:
      printf("%d", data->data.short_value);
      break;
    case PLC4C_USHORT:
      printf("%d", data->data.ushort_value);
      break;
    case PLC4C_INT:
      printf("%d", data->data.int_value);
      break;
    case PLC4C_UINT:
      printf("%iu", data->data.uint_value);
      break;
    case PLC4C_FLOAT:
      printf("%f", data->data.float_value);
      break;
    case PLC4C_STRING_POINTER:
      printf("%s", data->data.pstring_value);
      break;
    case PLC4C_CONSTANT_STRING:
      printf("%s", data->data.const_string_value);
      break;
    case PLC4C_VOID_POINTER:
      if (data->custom_printf != NULL) {
        data->custom_printf(data);
      } else {
        printf("%p", data->data.pvoid_value);
      }
      break;
    default:
      break;
  }
}

void plc4c_data_set_custom_destroy(
    plc4c_data *data, plc4c_data_custom_destroy data_custom_destroy) {
  data->custom_destroy = data_custom_destroy;
}

void plc4c_data_set_custom_printf(plc4c_data *data,
                                  plc4c_data_custom_printf data_custom_printf) {
  data->custom_printf = data_custom_printf;
}

void plc4c_data_destroy(plc4c_data *data) {
  assert(data != NULL);
  if (data->custom_destroy != NULL) {
    data->custom_destroy(data);
  } else {
    switch (data->data_type) {
      case PLC4C_VOID_POINTER:
        free(data->data.pvoid_value);
        break;
      case PLC4C_STRING_POINTER:
        free(data->data.pstring_value);
        break;
      default:
        break;
    }
  }
  free(data);
}
