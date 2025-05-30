package plugin

import (
	"strings"

	"github.com/evad1n/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type commentGenerator struct {
	descriptor protoreflect.Descriptor
}

func (c commentGenerator) generateLeading(f *codegen.File, indent int) {
	loc := c.descriptor.ParentFile().SourceLocations().ByDescriptor(c.descriptor)
	lines := strings.Split(loc.LeadingComments, "\n")
	f.Print(indentBy(indent), "/**")
	for _, line := range lines {
		if line == "" {
			continue
		}
		f.Print(indentBy(indent), " * ", strings.TrimSpace(line))
	}
	if field, ok := c.descriptor.(protoreflect.FieldDescriptor); ok {
		if behaviorComment := fieldBehaviorComment(field); len(behaviorComment) > 0 {
			f.Print(indentBy(indent), " * ")
			f.Print(indentBy(indent), " * ", behaviorComment)
		}
	}
	f.Print(indentBy(indent), " */")
}

func fieldBehaviorComment(field protoreflect.FieldDescriptor) string {
	behaviors := getFieldBehaviors(field)
	if len(behaviors) == 0 {
		return ""
	}

	behaviorStrings := make([]string, 0, len(behaviors))
	for _, b := range behaviors {
		behaviorStrings = append(behaviorStrings, b.String())
	}
	return "Behaviors: " + strings.Join(behaviorStrings, ", ")
}

func getFieldBehaviors(field protoreflect.FieldDescriptor) []annotations.FieldBehavior {
	if behaviors, ok := proto.GetExtension(
		field.Options(), annotations.E_FieldBehavior,
	).([]annotations.FieldBehavior); ok {
		return behaviors
	}
	return nil
}
