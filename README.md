# protoc-gen-endpoint

protoc-gen-endpoint is a protobuf plugin that allows generating extra
information for google.api.http.

# Usage

```protobuf
service Server {
    rpc Foo(Bar) returns (Baz) {
        option (google.api.http) = {
            get: "/v1/foo"
        };
        option (proto.endpoint) = {
            unauthenticated: true
        };
    }
}
```

Once `protoc-gen-endpoint` is installed in $PATH, the program can be used with
protoc via the flag: `--endpoint_out`

# License

BSD 3 clause
