// Code generated by protoc-gen-typescript-http. DO NOT EDIT.
/* eslint-disable camelcase */
// @ts-nocheck

type wellKnownDoubleValue = number | null;

type wellKnownFloatValue = number | null;

/**
 * Generated output always contains 0, 3, 6, or 9 fractional digits,
 * depending on required precision, followed by the suffix "s".
 * Accepted are any fractional digits (also none) as long as they fit
 * into nano-seconds precision and the suffix "s" is required.
 */
type wellKnownDuration = string;

type wellKnownBoolValue = boolean | null;

type wellKnownBytesValue = string | null;

type wellKnownUInt32Value = number | null;

type wellKnownUInt64Value = number | null;

type wellKnownListValue = wellKnownValue[];

type wellKnownInt64Value = number | null;

type wellKnownStringValue = string | null;

/**
 * In JSON, a field mask is encoded as a single string where paths are
 * separated by a comma. Fields name in each path are converted
 * to/from lower-camel naming conventions.
 * As an example, consider the following message declarations:
 *
 *     message Profile {
 *       User user = 1;
 *       Photo photo = 2;
 *     }
 *     message User {
 *       string display_name = 1;
 *       string address = 2;
 *     }
 *
 * In proto a field mask for `Profile` may look as such:
 *
 *     mask {
 *       paths: "user.display_name"
 *       paths: "photo"
 *     }
 *
 * In JSON, the same mask is represented as below:
 *
 *     {
 *       mask: "user.displayName,photo"
 *     }
 */
type wellKnownFieldMask = string;

type wellKnownValue = unknown;

type wellKnownNullValue = null;

type wellKnownInt32Value = number | null;

/**
 * If the Any contains a value that has a special JSON mapping,
 * it will be converted as follows:
 * {"@type": xxx, "value": yyy}.
 * Otherwise, the value will be converted into a JSON object,
 * and the "@type" field will be inserted to indicate the actual data type.
 */
interface wellKnownAny {
  "@type": string;
  [key: string]: unknown;
}

/**
 * An empty JSON object
 */
type wellKnownEmpty = Record<never, never>;

/**
 * Any JSON value.
 */
type wellKnownStruct = Record<string, unknown>;

/**
 * Enum
 */
export type einrideexamplesyntaxv1_Enum =
  /**
   * ENUM_UNSPECIFIED
   */
  | "ENUM_UNSPECIFIED"
  /**
   * ENUM_ONE
   */
  | "ENUM_ONE"
  /**
   * ENUM_TWO
   */
  | "ENUM_TWO";

/**
 * NestedEnum
 */
export type einrideexamplesyntaxv1_Message_NestedEnum =
  /**
   * NESTEDENUM_UNSPECIFIED
   */
  "NESTEDENUM_UNSPECIFIED";
/**
 * Message
 */
export type Message = {
  forwardedMessage: einrideexamplesyntaxv1_Message;
  forwardedEnum: einrideexamplesyntaxv1_Enum;
};

/**
 * Message
 */
export type einrideexamplesyntaxv1_Message = {
  /**
   * double
   */
  double: number;
  /**
   * float
   */
  float: number;
  /**
   * int32
   */
  int32: number;
  /**
   * int64
   */
  int64: number;
  /**
   * uint32
   */
  uint32: number;
  /**
   * uint64
   */
  uint64: number;
  /**
   * sint32
   */
  sint32: number;
  /**
   * sint64
   */
  sint64: number;
  /**
   * fixed32
   */
  fixed32: number;
  /**
   * fixed64
   */
  fixed64: number;
  /**
   * sfixed32
   */
  sfixed32: number;
  /**
   * sfixed64
   */
  sfixed64: number;
  /**
   * bool
   */
  bool: boolean;
  /**
   * string
   */
  string: string;
  /**
   * bytes
   */
  bytes: string;
  /**
   * enum
   */
  enum: einrideexamplesyntaxv1_Enum;
  /**
   * message
   */
  message: einrideexamplesyntaxv1_Message;
  /**
   * optional double
   */
  optionalDouble: number;
  /**
   * optional float
   */
  optionalFloat: number;
  /**
   * optional int32
   */
  optionalInt32: number;
  /**
   * optional int64
   */
  optionalInt64: number;
  /**
   * optional uint32
   */
  optionalUint32: number;
  /**
   * optional uint64
   */
  optionalUint64: number;
  /**
   * optional sint32
   */
  optionalSint32: number;
  /**
   * optional sint64
   */
  optionalSint64: number;
  /**
   * optional fixed32
   */
  optionalFixed32: number;
  /**
   * optional fixed64
   */
  optionalFixed64: number;
  /**
   * optional sfixed32
   */
  optionalSfixed32: number;
  /**
   * optional sfixed64
   */
  optionalSfixed64: number;
  /**
   * optional bool
   */
  optionalBool: boolean;
  /**
   * optional string
   */
  optionalString: string;
  /**
   * optional bytes
   */
  optionalBytes: string;
  /**
   * optional enum
   */
  optionalEnum: einrideexamplesyntaxv1_Enum;
  /**
   * optional message
   */
  optionalMessage: einrideexamplesyntaxv1_Message;
  /**
   * repeated_double
   */
  repeatedDouble: number[];
  /**
   * repeated_float
   */
  repeatedFloat: number[];
  /**
   * repeated_int32
   */
  repeatedInt32: number[];
  /**
   * repeated_int64
   */
  repeatedInt64: number[];
  /**
   * repeated_uint32
   */
  repeatedUint32: number[];
  /**
   * repeated_uint64
   */
  repeatedUint64: number[];
  /**
   * repeated_sint32
   */
  repeatedSint32: number[];
  /**
   * repeated_sint64
   */
  repeatedSint64: number[];
  /**
   * repeated_fixed32
   */
  repeatedFixed32: number[];
  /**
   * repeated_fixed64
   */
  repeatedFixed64: number[];
  /**
   * repeated_sfixed32
   */
  repeatedSfixed32: number[];
  /**
   * repeated_sfixed64
   */
  repeatedSfixed64: number[];
  /**
   * repeated_bool
   */
  repeatedBool: boolean[];
  /**
   * repeated_string
   */
  repeatedString: string[];
  /**
   * repeated_bytes
   */
  repeatedBytes: string[];
  /**
   * repeated_enum
   */
  repeatedEnum: einrideexamplesyntaxv1_Enum[];
  /**
   * repeated_message
   */
  repeatedMessage: einrideexamplesyntaxv1_Message[];
  /**
   * map_string_string
   */
  mapStringString: { [key: string]: string };
  /**
   * map_string_message
   */
  mapStringMessage: { [key: string]: einrideexamplesyntaxv1_Message };
  /**
   * oneof_string
   */
  oneofString: string;
  /**
   * oneof_enum
   */
  oneofEnum: einrideexamplesyntaxv1_Enum;
  /**
   * oneof_message1
   */
  oneofMessage1: einrideexamplesyntaxv1_Message;
  /**
   * oneof_message2
   */
  oneofMessage2: einrideexamplesyntaxv1_Message;
  /**
   * any
   */
  any: wellKnownAny;
  /**
   * repeated_any
   */
  repeatedAny: wellKnownAny[];
  /**
   * duration
   */
  duration: wellKnownDuration;
  /**
   * repeated_duration
   */
  repeatedDuration: wellKnownDuration[];
  /**
   * empty
   */
  empty: wellKnownEmpty;
  /**
   * repeated_empty
   */
  repeatedEmpty: wellKnownEmpty[];
  /**
   * field_mask
   */
  fieldMask: wellKnownFieldMask;
  /**
   * repeated_field_mask
   */
  repeatedFieldMask: wellKnownFieldMask[];
  /**
   * struct
   */
  struct: wellKnownStruct;
  /**
   * repeated_struct
   */
  repeatedStruct: wellKnownStruct[];
  /**
   * value
   */
  value: wellKnownValue;
  /**
   * repeated_value
   */
  repeatedValue: wellKnownValue[];
  /**
   * null_value
   */
  nullValue: wellKnownNullValue;
  /**
   * repeated_null_value
   */
  repeatedNullValue: wellKnownNullValue[];
  /**
   * list_value
   */
  listValue: wellKnownListValue;
  /**
   * repeated_list_value
   */
  repeatedListValue: wellKnownListValue[];
  /**
   * bool_value
   */
  boolValue: wellKnownBoolValue;
  /**
   * repeated_bool_value
   */
  repeatedBoolValue: wellKnownBoolValue[];
  /**
   * bytes_value
   */
  bytesValue: wellKnownBytesValue;
  /**
   * repeated_bytes_value
   */
  repeatedBytesValue: wellKnownBytesValue[];
  /**
   * double_value
   */
  doubleValue: wellKnownDoubleValue;
  /**
   * repeated_double_value
   */
  repeatedDoubleValue: wellKnownDoubleValue[];
  /**
   * float_value
   */
  floatValue: wellKnownFloatValue;
  /**
   * repeated_float_value
   */
  repeatedFloatValue: wellKnownFloatValue[];
  /**
   * int32_value
   */
  int32Value: wellKnownInt32Value;
  /**
   * repeated_int32_value
   */
  repeatedInt32Value: wellKnownInt32Value[];
  /**
   * int64_value
   */
  int64Value: wellKnownInt64Value;
  /**
   * repeated_int64_value
   */
  repeatedInt64Value: wellKnownInt64Value[];
  /**
   * uint32_value
   */
  uint32Value: wellKnownUInt32Value;
  /**
   * repeated_uint32_value
   */
  repeatedUint32Value: wellKnownUInt32Value[];
  /**
   * uint64_value
   */
  uint64Value: wellKnownUInt64Value;
  /**
   * repeated_uint64_value
   */
  repeatedUint64Value: wellKnownUInt64Value[];
  /**
   * string_value
   */
  stringValue: wellKnownUInt64Value;
  /**
   * repeated_string_value
   */
  repeatedStringValue: wellKnownStringValue[];
};

/**
 * NestedMessage
 */
export type einrideexamplesyntaxv1_Message_NestedMessage = {
  /**
   * nested_message.string
   */
  string: string;
};


// @@protoc_insertion_point(typescript-http-eof)
