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

	// Collect comment lines
	commentLines := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		commentLines = append(commentLines, strings.TrimSpace(line))
	}

	var behaviorComment string
	if field, ok := c.descriptor.(protoreflect.FieldDescriptor); ok {
		behaviorComment = fieldBehaviorComment(field)
	}

	// If there are no comments and no behaviors, do not write anything
	if len(commentLines) == 0 && len(behaviorComment) == 0 {
		return
	}

	f.Write(indentBy(indent), "/**")
	for _, line := range commentLines {
		f.Write(indentBy(indent), " * ", line)
	}
	if len(behaviorComment) > 0 {
		if len(commentLines) > 0 {
			f.Write(indentBy(indent), " * ")
		}
		f.Write(indentBy(indent), " * ", behaviorComment)
	}
	f.Write(indentBy(indent), " */")
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
