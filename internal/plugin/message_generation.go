package plugin

import (
	"fmt"
	"slices"

	"github.com/evad1n/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	generatedMessagesRegistry = make(map[string]protoreflect.MessageDescriptor) // To avoid generating the same message multiple times
)

type messageGenerator struct {
	pkg            protoreflect.FullName
	message        protoreflect.MessageDescriptor
	usedInRequest  bool
	usedInResponse bool
}

func (m messageGenerator) Generate(f *codegen.File) {
	if !m.usedInRequest && !m.usedInResponse {
		m.generateDefaultType(f)
	}

	if m.usedInRequest {
		m.generateRequestType(f)
	}

	if m.usedInResponse {
		m.generateResponseType(f)
	}
}

func (m messageGenerator) generateDefaultType(f *codegen.File) {
	commentGenerator{descriptor: m.message}.generateLeading(f, 0)

	f.Write("export type ", scopedDescriptorTypeName(m.pkg, m.message), " = {")

	rangeFields(m.message, func(field protoreflect.FieldDescriptor) {
		commentGenerator{descriptor: field}.generateLeading(f, 1)
		fieldType := typeFromField(m.pkg, field)

		f.Write(indentBy(1), field.JSONName(), ": ", fieldType.Reference(), ";")
	})

	f.Write("};")
	f.Write()
}

func (m messageGenerator) generateRequestType(f *codegen.File) {
	commentGenerator{descriptor: m.message}.generateLeading(f, 0)

	// We are generating inline, so keep track of additional messages that need to be generated and generate at the end
	additionalMessagesToGenerate := make(map[protoreflect.FullName]messageGenerator)

	typeName := suffixName(scopedDescriptorTypeName(m.pkg, m.message), REQUEST_SUFFIX)
	if existingType, ok := generatedMessagesRegistry[typeName]; ok {
		// If the request type already exists, we can skip generating it.
		log(fmt.Sprintf("skipping request type %s, already exists as %s", typeName, existingType))
		return
	}
	generatedMessagesRegistry[typeName] = m.message

	f.Write("export type ", typeName, " = {")

	rangeFields(m.message, func(field protoreflect.FieldDescriptor) {
		if !getFieldShouldGenerate(field, true) {
			return
		}

		commentGenerator{descriptor: field}.generateLeading(f, 1)

		var fieldTypeName string

		if field.Kind() == protoreflect.MessageKind {
			message := field.Message()

			messageRequiresDiscrimination := getMessageRequiresDiscrimination(message, 0, make(map[protoreflect.FullName]bool))

			if messageRequiresDiscrimination {
				fieldTypeName = suffixName(scopedDescriptorTypeName(m.pkg, message), REQUEST_SUFFIX)

				// If the field is a message we need to recurse
				messageGen := messageGenerator{
					pkg:            m.pkg,
					message:        message,
					usedInRequest:  true,
					usedInResponse: m.usedInResponse,
				}
				additionalMessagesToGenerate[message.FullName()] = messageGen
			} else {
				fieldTypeName = typeFromField(m.pkg, field).Reference()
			}
		} else {
			fieldTypeName = typeFromField(m.pkg, field).Reference()
		}

		fieldCardinalitySymbol := getFieldCardinalitySymbol(field, true)

		f.Write(indentBy(1), field.JSONName(), fieldCardinalitySymbol, ": ", fieldTypeName, ";")
	})

	f.Write("};")
	f.Write()

	// Generate additional messages that were collected
	for _, additionalMessage := range additionalMessagesToGenerate {
		additionalMessage.generateRequestType(f)
	}
}

func (m messageGenerator) generateResponseType(f *codegen.File) {
	commentGenerator{descriptor: m.message}.generateLeading(f, 0)

	// We are generating inline, so keep track of additional messages that need to be generated and generate at the end
	additionalMessagesToGenerate := make(map[protoreflect.FullName]messageGenerator)

	typeName := suffixName(scopedDescriptorTypeName(m.pkg, m.message), RESPONSE_SUFFIX)
	if existingType, ok := generatedMessagesRegistry[typeName]; ok {
		// If the response type already exists, we can skip generating it.
		log(fmt.Sprintf("skipping response type %s, already exists as %s", typeName, existingType))
		return
	}
	generatedMessagesRegistry[typeName] = m.message

	f.Write("export type ", typeName, " = {")

	rangeFields(m.message, func(field protoreflect.FieldDescriptor) {
		if !getFieldShouldGenerate(field, false) {
			return
		}

		commentGenerator{descriptor: field}.generateLeading(f, 1)

		var fieldTypeName string

		if field.Kind() == protoreflect.MessageKind {
			message := field.Message()

			messageRequiresDiscrimination := getMessageRequiresDiscrimination(message, 0, make(map[protoreflect.FullName]bool))

			if messageRequiresDiscrimination {
				fieldTypeName = suffixName(scopedDescriptorTypeName(m.pkg, message), RESPONSE_SUFFIX)

				// If the field is a message we need to recurse
				messageGen := messageGenerator{
					pkg:            m.pkg,
					message:        message,
					usedInRequest:  m.usedInRequest,
					usedInResponse: true,
				}
				additionalMessagesToGenerate[message.FullName()] = messageGen
			} else {
				fieldTypeName = typeFromField(m.pkg, field).Reference()
			}
		} else {
			fieldTypeName = typeFromField(m.pkg, field).Reference()
		}

		fieldCardinalitySymbol := getFieldCardinalitySymbol(field, false)

		f.Write(indentBy(1), field.JSONName(), fieldCardinalitySymbol, ": ", fieldTypeName, ";")
	})

	f.Write("};")
	f.Write()

	// Generate additional messages that were collected
	for _, additionalMessage := range additionalMessagesToGenerate {
		additionalMessage.generateResponseType(f)
	}
}

func getFieldShouldGenerate(field protoreflect.FieldDescriptor, isRequest bool) bool {
	behaviors := getFieldBehaviors(field)

	if isRequest {
		if slices.Contains(behaviors, annotations.FieldBehavior_OUTPUT_ONLY) {
			return false
		}

		return true
	}

	// Response
	if slices.Contains(behaviors, annotations.FieldBehavior_INPUT_ONLY) {
		return false
	}

	return true
}

func getFieldCardinalitySymbol(field protoreflect.FieldDescriptor, isRequest bool) string {
	behaviors := getFieldBehaviors(field)

	if field.ContainingOneof() != nil {
		return "?"
	}

	if isRequest {
		if slices.Contains(behaviors, annotations.FieldBehavior_OPTIONAL) {
			return "?"
		}
	}

	return ""
}

var behaviorsRequiringDiscrimination = []annotations.FieldBehavior{
	annotations.FieldBehavior_OUTPUT_ONLY,
	annotations.FieldBehavior_INPUT_ONLY,
	annotations.FieldBehavior_OPTIONAL,
}

// getMessageRequiresDiscrimination checks if any of a message's fields' behavior annotation suggests it should have different type definitions for the request and response. Works recursively through nested messages.
// The visited map prevents infinite recursion in case of cyclic message references.
func getMessageRequiresDiscrimination(message protoreflect.MessageDescriptor, depth int, visited map[protoreflect.FullName]bool) bool {
	if visited[message.FullName()] {
		return false
	}
	visited[message.FullName()] = true

	messageRequiresDiscrimination := false
	rangeFields(message, func(field protoreflect.FieldDescriptor) {
		// If it's a message field, we need to check its nested fields as well.
		if field.Kind() == protoreflect.MessageKind {
			nestedMessage := field.Message()

			if getMessageRequiresDiscrimination(nestedMessage, depth+1, visited) {
				messageRequiresDiscrimination = true
			}
			return
		}

		if getFieldRequiresRequestDiscrimination(field) {
			messageRequiresDiscrimination = true
		}
	})

	return messageRequiresDiscrimination
}

// getFieldRequiresRequestDiscrimination checks if a field's behavior annotation suggests it should have different type definitions for the request and response. Example: OUTPUT_ONLY would be a field that is not required in the request, but is present in the response.
func getFieldRequiresRequestDiscrimination(field protoreflect.FieldDescriptor) bool {
	behaviors := getFieldBehaviors(field)
	return slices.ContainsFunc(behaviors, func(b annotations.FieldBehavior) bool {
		return slices.Contains(behaviorsRequiringDiscrimination, b)
	})
}
