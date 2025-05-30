package plugin

import (
	"slices"

	"github.com/evad1n/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type messageGenerator struct {
	pkg     protoreflect.FullName
	message protoreflect.MessageDescriptor
}

func (m messageGenerator) Generate(f *codegen.File) {
	commentGenerator{descriptor: m.message}.generateLeading(f, 0)

	f.Print("export type ", scopedDescriptorTypeName(m.pkg, m.message), " = {")

	rangeFields(m.message, func(field protoreflect.FieldDescriptor) {
		commentGenerator{descriptor: field}.generateLeading(f, 1)
		fieldType := typeFromField(m.pkg, field)

		behaviors := getFieldBehaviors(field)

		if slices.Contains(behaviors, annotations.FieldBehavior_OPTIONAL) {
			f.Print(indentBy(1), field.JSONName(), "?: ", fieldType.Reference(), ";")
		} else if slices.Contains(behaviors, annotations.FieldBehavior_REQUIRED) {
			f.Print(indentBy(1), field.JSONName(), ": ", fieldType.Reference(), ";")
		} else if field.ContainingOneof() == nil && !field.HasOptionalKeyword() {
			f.Print(indentBy(1), field.JSONName(), ": ", fieldType.Reference(), ";")
		} else {
			f.Print(indentBy(1), field.JSONName(), "?: ", fieldType.Reference(), ";")
		}
	})

	f.Print("};")
	f.Print()
}
