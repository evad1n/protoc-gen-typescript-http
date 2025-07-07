package plugin

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/evad1n/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type generatorOptions struct {
	verbose            bool
	requestTypeSuffix  string
	responseTypeSuffix string
}

func (o generatorOptions) String() string {
	var opts []string
	opts = append(opts, fmt.Sprintf("verbose=%v", o.verbose))
	return strings.Join(opts, ",")
}

var (
	options generatorOptions = generatorOptions{
		verbose:            false,
		requestTypeSuffix:  "__Request",
		responseTypeSuffix: "__Response",
	}
	generationErrors []error
)

func log(args ...any) {
	fmt.Fprint(os.Stderr, "protoc-gen-typescript-http: ")
	fmt.Fprintln(os.Stderr, args...)
}

// logV is a verbose log function that only logs if verbose mode is enabled.
func logV(args ...any) {
	if !options.verbose {
		return
	}
	log(args...)
}

func addGenerationError(err error) {
	generationErrors = append(generationErrors, err)
}

func Generate(request *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	opts, err := parseOptions(request.GetParameter())
	if err != nil {
		return nil, fmt.Errorf("parse options: %w", err)
	}

	options = opts

	logV("options:", options)

	logV("generating files for", len(request.GetFileToGenerate()), "files")
	for _, f := range request.GetFileToGenerate() {
		logV("generating file", f)
	}

	generate := make(map[string]struct{})
	registry, err := protodesc.NewFiles(&descriptorpb.FileDescriptorSet{
		File: request.GetProtoFile(),
	})
	if err != nil {
		return nil, fmt.Errorf("create proto registry: %w", err)
	}
	for _, f := range request.GetFileToGenerate() {
		generate[f] = struct{}{}
	}
	packageRegistry := make(map[protoreflect.FullName][]protoreflect.FileDescriptor)
	for _, f := range request.GetFileToGenerate() {
		file, err := registry.FindFileByPath(f)
		if err != nil {
			return nil, fmt.Errorf("find file %s: %w", f, err)
		}
		packageRegistry[file.Package()] = append(packageRegistry[file.Package()], file)
	}

	var res pluginpb.CodeGeneratorResponse
	for pkg, files := range packageRegistry {
		logV(fmt.Sprint(string(pkg), ":"))
		for _, file := range files {
			logV(fmt.Sprint(indentBy(1), file.Path()))
		}

		var index codegen.File
		indexPathElems := append(strings.Split(string(pkg), "."), "index.ts")
		if err := (packageGenerator{pkg: pkg, files: files}).Generate(&index); err != nil {
			return nil, fmt.Errorf("generate package '%s': %w", pkg, err)
		}
		index.Write()
		index.Write("// @@protoc_insertion_point(typescript-http-eof)")
		res.File = append(res.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(path.Join(indexPathElems...)),
			Content: proto.String(string(index.Content())),
		})
	}
	res.SupportedFeatures = proto.Uint64(uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL))

	if len(generationErrors) > 0 {
		log("generation errors:")
		for _, err := range generationErrors {
			log(err)
		}
		return nil, fmt.Errorf("encountered %d errors during generation", len(generationErrors))
	}

	return &res, nil
}

// Looks like `verbose=true,param=value`
func parseOptions(parameterString string) (generatorOptions, error) {
	opts := generatorOptions{}
	if parameterString == "" {
		return opts, nil
	}
	for _, opt := range strings.Split(parameterString, ",") {
		opt = strings.TrimSpace(opt)
		if opt == "" {
			continue
		}
		parts := strings.SplitN(opt, "=", 2)
		key := parts[0]
		val := "true"
		if len(parts) == 2 {
			val = parts[1]
		}
		switch key {
		case "verbose":
			opts.verbose = val == "true"
		case "requestTypeSuffix":
			if val == "" {
				return opts, fmt.Errorf("requestTypeSuffix cannot be empty")
			}
			opts.requestTypeSuffix = val
		case "responseTypeSuffix":
			if val == "" {
				return opts, fmt.Errorf("responseTypeSuffix cannot be empty")
			}
			opts.responseTypeSuffix = val
		default:
			return opts, fmt.Errorf("unknown option: %s", key)
		}
	}
	return opts, nil
}
