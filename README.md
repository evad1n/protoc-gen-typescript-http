# protoc-gen-typescript-http

Generates Typescript types and service clients from protobuf definitions
annotated with
[http rules](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto).
The generated types follow the
[canonical JSON encoding](https://developers.google.com/protocol-buffers/docs/proto3#json).

**Experimental**: This library is under active development and breaking changes
to config files, APIs and generated code are expected between releases.

## Using the plugin

For examples of correctly annotated protobuf defintions and the generated code,
look at [examples](./examples).

### Install the plugin

```bash
go install github.com/evad1n/protoc-gen-typescript-http@latest
```

### Invocation

```bash
protoc
  --typescript-http_out [OUTPUT DIR] \
  [.proto files ...]
```

### With `buf` and `buf.gen.yaml`

> https://buf.build/docs/configuration/v2/buf-gen-yaml/

```yml
plugins:
  - local: protoc-gen-typescript-http
    out: src/gen
    opt:
      - verbose=true
```

### Options

- `verbose` - print some extra information when running
- `requestTypeSuffix` - Suffix for generating request types. Default = `__Request`
- `responseTypeSuffix` - Suffix for generating response types Default = `__Response`

______________________________________________________________________

The generated clients can be used with any HTTP client that returns a Promise
containing JSON data.

```typescript
const rootUrl = "...";

type Request = {
  path: string,
  method: string,
  body: string | null
}

function fetchRequestHandler({path, method, body}: Request) {
  return fetch(rootUrl + path, {method, body}).then(response => response.json())
}

export function siteClient() {
  return createShipperServiceClient(fetchRequestHandler);
}
```
