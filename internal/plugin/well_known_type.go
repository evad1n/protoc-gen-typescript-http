package plugin

import (
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	wellKnownPrefix = "google.protobuf."
)

// WellKnownType represents a well-known type in the Google Protocol Buffers ecosystem. Such as google.protobuf.Timestamp.
type WellKnownType string

// https://developers.google.com/protocol-buffers/docs/reference/google.protobuf
const (
	WellKnownAny       WellKnownType = "google.protobuf.Any"
	WellKnownDuration  WellKnownType = "google.protobuf.Duration"
	WellKnownEmpty     WellKnownType = "google.protobuf.Empty"
	WellKnownFieldMask WellKnownType = "google.protobuf.FieldMask"
	WellKnownStruct    WellKnownType = "google.protobuf.Struct"
	WellKnownTimestamp WellKnownType = "google.protobuf.Timestamp"

	// Wrapper types.
	WellKnownFloatValue  WellKnownType = "google.protobuf.FloatValue"
	WellKnownInt64Value  WellKnownType = "google.protobuf.Int64Value"
	WellKnownInt32Value  WellKnownType = "google.protobuf.Int32Value"
	WellKnownUInt64Value WellKnownType = "google.protobuf.UInt64Value"
	WellKnownUInt32Value WellKnownType = "google.protobuf.UInt32Value"
	WellKnownBytesValue  WellKnownType = "google.protobuf.BytesValue"
	WellKnownDoubleValue WellKnownType = "google.protobuf.DoubleValue"
	WellKnownBoolValue   WellKnownType = "google.protobuf.BoolValue"
	WellKnownStringValue WellKnownType = "google.protobuf.StringValue"

	// Descriptor types.
	WellKnownValue     WellKnownType = "google.protobuf.Value"
	WellKnownNullValue WellKnownType = "google.protobuf.NullValue"
	WellKnownListValue WellKnownType = "google.protobuf.ListValue"
)

func IsWellKnownType(desc protoreflect.Descriptor) bool {
	switch desc.(type) {
	case protoreflect.MessageDescriptor, protoreflect.EnumDescriptor:
		return strings.HasPrefix(string(desc.FullName()), wellKnownPrefix)
	default:
		return false
	}
}

func GetWellKnownType(desc protoreflect.Descriptor) (WellKnownType, bool) {
	if !IsWellKnownType(desc) {
		return "", false
	}
	return WellKnownType(desc.FullName()), true
}

func (wkt WellKnownType) Name() string {
	return "wellKnown" + strings.TrimPrefix(string(wkt), wellKnownPrefix)
}

func (wkt WellKnownType) TypeDeclaration() string {
	var w writer
	switch wkt {
	case WellKnownAny:
		w.Write("/**")
		w.Write(" * If the Any contains a value that has a special JSON mapping,")
		w.Write(" * it will be converted as follows:")
		w.Write(" * {\"@type\": xxx, \"value\": yyy}.")
		w.Write(" * Otherwise, the value will be converted into a JSON object,")
		w.Write(" * and the \"@type\" field will be inserted to indicate the actual data type.")
		w.Write(" */")
		w.Write("interface ", wkt.Name(), " {")
		w.Write("  ", "\"@type\": string;")
		w.Write("  [key: string]: unknown;")
		w.Write("}")
	case WellKnownDuration:
		w.Write("/**")
		w.Write(" * Generated output always contains 0, 3, 6, or 9 fractional digits,")
		w.Write(" * depending on required precision, followed by the suffix \"s\".")
		w.Write(" * Accepted are any fractional digits (also none) as long as they fit")
		w.Write(" * into nano-seconds precision and the suffix \"s\" is required.")
		w.Write(" */")
		w.Write("type ", wkt.Name(), " = string;")
	case WellKnownEmpty:
		w.Write("/**")
		w.Write(" * An empty JSON object")
		w.Write(" */")
		w.Write("type ", wkt.Name(), " = Record<never, never>;")
	case WellKnownTimestamp:
		w.Write("/**")
		w.Write(" * Encoded using RFC 3339, where generated output will always be Z-normalized")
		w.Write(" * and uses 0, 3, 6 or 9 fractional digits.")
		w.Write(" * Offsets other than \"Z\" are also accepted.")
		w.Write(" */")
		w.Write("type ", wkt.Name(), " = string;")
	case WellKnownFieldMask:
		w.Write("/**")
		w.Write(" * In JSON, a field mask is encoded as a single string where paths are")
		w.Write(" * separated by a comma. Fields name in each path are converted")
		w.Write(" * to/from lower-camel naming conventions.")
		w.Write(" * As an example, consider the following message declarations:")
		w.Write(" *")
		w.Write(" *     message Profile {")
		w.Write(" *       User user = 1;")
		w.Write(" *       Photo photo = 2;")
		w.Write(" *     }")
		w.Write(" *     message User {")
		w.Write(" *       string display_name = 1;")
		w.Write(" *       string address = 2;")
		w.Write(" *     }")
		w.Write(" *")
		w.Write(" * In proto a field mask for `Profile` may look as such:")
		w.Write(" *")
		w.Write(" *     mask {")
		w.Write(" *       paths: \"user.display_name\"")
		w.Write(" *       paths: \"photo\"")
		w.Write(" *     }")
		w.Write(" *")
		w.Write(" * In JSON, the same mask is represented as below:")
		w.Write(" *")
		w.Write(" *     {")
		w.Write(" *       mask: \"user.displayName,photo\"")
		w.Write(" *     }")
		w.Write(" */")
		w.Write("type ", wkt.Name(), " = string;")
	case WellKnownFloatValue,
		WellKnownDoubleValue,
		WellKnownInt64Value,
		WellKnownInt32Value,
		WellKnownUInt64Value,
		WellKnownUInt32Value:
		w.Write("type ", wkt.Name(), " = number | null;")
	case WellKnownBytesValue, WellKnownStringValue:
		w.Write("type ", wkt.Name(), " = string | null;")
	case WellKnownBoolValue:
		w.Write("type ", wkt.Name(), " = boolean | null;")
	case WellKnownStruct:
		w.Write("/**")
		w.Write(" * Any JSON value.")
		w.Write(" */")
		w.Write("type ", wkt.Name(), " = Record<string, unknown>;")
	case WellKnownValue:
		w.Write("type ", wkt.Name(), " = unknown;")
	case WellKnownNullValue:
		w.Write("type ", wkt.Name(), " = null;")
	case WellKnownListValue:
		w.Write("type ", wkt.Name(), " = ", WellKnownValue.Name(), "[];")
	default:
		w.Write("/**")
		w.Write(" * No mapping for this well known type is generated, yet.")
		w.Write(" */")
		w.Write("type ", wkt.Name(), " = unknown;")
	}
	return w.String()
}

type writer struct {
	b strings.Builder
}

func (w *writer) Write(ss ...string) {
	for _, s := range ss {
		// strings.Builder never returns an error, so safe to ignore
		_, _ = w.b.WriteString(s)
	}
	_, _ = w.b.WriteString("\n")
}

func (w *writer) String() string {
	return w.b.String()
}
