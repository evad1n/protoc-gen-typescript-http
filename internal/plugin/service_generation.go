package plugin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/evad1n/protoc-gen-typescript-http/internal/codegen"
	"github.com/evad1n/protoc-gen-typescript-http/internal/httprule"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type serviceGenerator struct {
	pkg     protoreflect.FullName
	service protoreflect.ServiceDescriptor
}

func GenerateServiceHeader(f *codegen.File) {
	f.Write("type RequestType = {")
	f.Write(indentBy(1), "path: string;")
	f.Write(indentBy(1), "method: string;")
	f.Write(indentBy(1), "body: string | null;")
	f.Write("};")
	f.Write()
	f.Write("type RequestHandler = (request: RequestType, meta: { service: string, method: string }) => Promise<unknown>;")
	f.Write()
}

func (s serviceGenerator) Generate(f *codegen.File) error {
	s.generateInterface(f)
	return s.generateClient(f)
}

func (s serviceGenerator) generateInterface(f *codegen.File) {
	commentGenerator{descriptor: s.service}.generateLeading(f, 0)
	f.Write("export interface ", descriptorTypeName(s.service), " {")
	rangeMethods(s.service.Methods(), func(method protoreflect.MethodDescriptor) {
		if !supportedMethod(method) {
			return
		}
		commentGenerator{descriptor: method}.generateLeading(f, 1)
		input := typeFromMessage(s.pkg, method.Input())
		output := typeFromMessage(s.pkg, method.Output())

		inputName := suffixName(input.Reference(), REQUEST_SUFFIX)
		if _, ok := GetWellKnownType(method.Input()); ok {
			inputName = input.Reference()
		}
		outputName := suffixName(output.Reference(), RESPONSE_SUFFIX)
		if _, ok := GetWellKnownType(method.Output()); ok {
			outputName = output.Reference()
		}

		f.Write(indentBy(1), method.Name(), "(request: ", inputName, "): Promise<", outputName, ">;")
	})
	f.Write("}")
	f.Write()
}

func (s serviceGenerator) generateClient(f *codegen.File) error {
	f.Write(
		"export function create",
		descriptorTypeName(s.service),
		"Client(",
		"\n",
		indentBy(1),
		"handler: RequestHandler",
		"\n",
		"): ",
		descriptorTypeName(s.service),
		" {",
	)
	f.Write(indentBy(1), "return {")
	var methodErr error
	rangeMethods(s.service.Methods(), func(method protoreflect.MethodDescriptor) {
		if err := s.generateMethod(f, method); err != nil {
			methodErr = fmt.Errorf("generate method %s: %w", method.Name(), err)
		}
	})
	if methodErr != nil {
		return methodErr
	}
	f.Write(indentBy(1), "};")
	f.Write("}")
	return nil
}

func (s serviceGenerator) generateMethod(f *codegen.File, method protoreflect.MethodDescriptor) error {
	outputType := typeFromMessage(s.pkg, method.Output())
	httpRule, ok := httprule.Get(method)
	if !ok {
		return nil
	}
	rule, err := httprule.ParseRule(httpRule)
	if err != nil {
		return fmt.Errorf("parse http rule: %w", err)
	}
	logV("generating method:", method.FullName(), httpRule)
	f.Write(indentBy(2), method.Name(), "(request) { // eslint-disable-line @typescript-eslint/no-unused-vars")
	s.generateMethodPathValidation(f, method, rule)
	s.generateMethodPath(f, method, rule)
	s.generateMethodBody(f, method, rule)
	s.generateMethodQuery(f, method, rule)
	f.Write(indentBy(3), "let uri = path;")
	f.Write(indentBy(3), "if (queryParams.length > 0) {")
	f.Write(indentBy(4), "uri += `?${queryParams.join(\"&\")}`")
	f.Write(indentBy(3), "}")
	f.Write(indentBy(3), "return handler({")
	f.Write(indentBy(4), "path: uri,")
	f.Write(indentBy(4), "method: ", strconv.Quote(rule.Method), ",")
	f.Write(indentBy(4), "body,")
	f.Write(indentBy(3), "}, {")
	f.Write(indentBy(4), "service: \"", method.Parent().Name(), "\",")
	f.Write(indentBy(4), "method: \"", method.Name(), "\",")
	f.Write(indentBy(3), "}) as Promise<", outputType.Reference(), ">;")
	f.Write(indentBy(2), "},")
	return nil
}

func (s serviceGenerator) generateMethodPathValidation(
	f *codegen.File,
	method protoreflect.MethodDescriptor,
	rule httprule.Rule,
) {

	// fmt.Fprintln(os.Stderr, "generateMethodPathValidation", rule.Method, rule.Template.Segments)
	for _, seg := range rule.Template.Segments {
		if seg.Kind != httprule.SegmentKindVariable {
			continue
		}
		fp := seg.Variable.FieldPath
		nullPath := nullPropagationPath(fp, method)
		protoPath := strings.Join(fp, ".")
		errMsg := "missing required field request." + protoPath
		f.Write(indentBy(3), "if (!request.", nullPath, ") {")
		f.Write(indentBy(4), "throw new Error(", strconv.Quote(errMsg), ");")
		f.Write(indentBy(3), "}")
	}
}

func (s serviceGenerator) generateMethodPath(
	f *codegen.File,
	method protoreflect.MethodDescriptor,
	rule httprule.Rule,
) {
	pathParts := make([]string, 0, len(rule.Template.Segments))
	for _, seg := range rule.Template.Segments {
		switch seg.Kind {
		case httprule.SegmentKindVariable:
			fieldPath := jsonPath(seg.Variable.FieldPath, method)
			pathParts = append(pathParts, "${request."+fieldPath+"}")
		case httprule.SegmentKindLiteral:
			pathParts = append(pathParts, seg.Literal)
		case httprule.SegmentKindMatchSingle: // TODO: Double check this and following case
			pathParts = append(pathParts, "*")
		case httprule.SegmentKindMatchMultiple:
			pathParts = append(pathParts, "**")
		}
	}
	path := strings.Join(pathParts, "/")
	if rule.Template.Verb != "" {
		path += ":" + rule.Template.Verb
	}
	f.Write(indentBy(3), "const path = `", path, "`; // eslint-disable-line quotes")
}

func (s serviceGenerator) generateMethodBody(
	f *codegen.File,
	method protoreflect.MethodDescriptor,
	rule httprule.Rule,
) {
	switch {
	case rule.Body == "":
		f.Write(indentBy(3), "const body = null;")
	case rule.Body == "*":
		f.Write(indentBy(3), "const body = JSON.stringify(request);")
	default:
		nullPath := nullPropagationPath(httprule.FieldPath{rule.Body}, method)
		f.Write(indentBy(3), "const body = JSON.stringify(request?.", nullPath, " ?? {});")
	}
}

func (s serviceGenerator) generateMethodQuery(
	f *codegen.File,
	method protoreflect.MethodDescriptor,
	rule httprule.Rule,
) {
	f.Write(indentBy(3), "const queryParams: string[] = [];")
	// nothing in query
	if rule.Body == "*" {
		return
	}
	// fields covered by path
	pathCovered := make(map[string]struct{})
	for _, segment := range rule.Template.Segments {
		if segment.Kind != httprule.SegmentKindVariable {
			continue
		}
		pathCovered[segment.Variable.FieldPath.String()] = struct{}{}
	}
	walkJSONLeafFields(method, func(path httprule.FieldPath, field protoreflect.FieldDescriptor) {
		if _, ok := pathCovered[path.String()]; ok {
			return
		}
		if rule.Body != "" && path[0] == rule.Body {
			return
		}
		nullPath := nullPropagationPath(path, method)
		jp := jsonPath(path, method)
		f.Write(indentBy(3), "if (request.", nullPath, ") {")
		switch {
		case field.IsList():
			f.Write(indentBy(4), "request.", jp, ".forEach((x) => {")
			f.Write(indentBy(5), "queryParams.push(`", jp, "=${encodeURIComponent(x.toString())}`)")
			f.Write(indentBy(4), "})")
		default:
			f.Write(indentBy(4), "queryParams.push(`", jp, "=${encodeURIComponent(request.", jp, ".toString())}`)")
		}
		f.Write(indentBy(3), "}")
	})
}

func supportedMethod(method protoreflect.MethodDescriptor) bool {
	_, ok := httprule.Get(method)
	return ok && !method.IsStreamingClient() && !method.IsStreamingServer()
}

func jsonPath(path httprule.FieldPath, method protoreflect.MethodDescriptor) string {
	return strings.Join(jsonPathSegments(path, method), ".")
}

func nullPropagationPath(path httprule.FieldPath, method protoreflect.MethodDescriptor) string {
	return strings.Join(jsonPathSegments(path, method), "?.")
}

func jsonPathSegments(path httprule.FieldPath, method protoreflect.MethodDescriptor) []string {
	message := method.Input()
	segs := make([]string, len(path))
	for i, p := range path {
		field := message.Fields().ByName(protoreflect.Name(p))
		if field == nil {
			err := fmt.Errorf("ERROR: (%s) field %q not found in message %q", method.FullName(), p, message.FullName())
			logV(err)
			addGenerationError(err)
		} else {
			segs[i] = field.JSONName()
			if i < len(path) {
				message = field.Message()
			}
		}
	}
	return segs
}
