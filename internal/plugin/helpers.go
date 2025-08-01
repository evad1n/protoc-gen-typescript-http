package plugin

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func scopedDescriptorTypeName(pkg protoreflect.FullName, desc protoreflect.Descriptor) string {
	name := string(desc.Name())
	var prefix string
	if desc.Parent() != desc.ParentFile() {
		prefix = descriptorTypeName(desc.Parent()) + "_"
	}
	if desc.ParentFile().Package() != pkg {
		prefix = packagePrefix(desc.ParentFile().Package()) + prefix
	}
	return prefix + name
}

func descriptorTypeName(desc protoreflect.Descriptor) string {
	name := string(desc.Name())
	var prefix string
	if desc.Parent() != desc.ParentFile() {
		prefix = descriptorTypeName(desc.Parent()) + "_"
	}
	return prefix + name
}

func packagePrefix(pkg protoreflect.FullName) string {
	return strings.Join(strings.Split(string(pkg), "."), "") + "_"
}

func rangeFields(message protoreflect.MessageDescriptor, f func(field protoreflect.FieldDescriptor)) {
	for i := 0; i < message.Fields().Len(); i++ {
		f(message.Fields().Get(i))
	}
}

func rangeMethods(methods protoreflect.MethodDescriptors, f func(method protoreflect.MethodDescriptor)) {
	for i := 0; i < methods.Len(); i++ {
		f(methods.Get(i))
	}
}

func rangeEnumValues(enum protoreflect.EnumDescriptor, f func(value protoreflect.EnumValueDescriptor, last bool)) {
	for i := 0; i < enum.Values().Len(); i++ {
		if i == enum.Values().Len()-1 {
			f(enum.Values().Get(i), true)
		} else {
			f(enum.Values().Get(i), false)
		}
	}
}

// indentBy is a utility function that returns a string with `n` levels of indentation, where each level is represented by two spaces.
func indentBy(n int) string {
	return strings.Repeat("  ", n)
}

const (
	REQUEST_SUFFIX  = "__Request"
	RESPONSE_SUFFIX = "__Response"
)

// suffixName is a utility function that appends a suffix to a name if it does not already end with that suffix.
func suffixName(name string, suffix string) string {
	if strings.HasSuffix(name, suffix) {
		return name
	}

	return fmt.Sprintf("%s%s", name, suffix)
}
