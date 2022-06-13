# gRPC connection with SSL/TLS in Go

## Types of gRPC connections
There are 3 types of gRPC connections:
```
    1. Insecure 
        *  Plaintext data, no encryption 
    2. Server-side TLS
        *  Encrypted data. Only Server needs to provide its certificate to client.
    3. Mutual TLS
        *  Encrypted data. Both server and client need to provide certificates to each other.
```

## grpc.WithPerRPCCredentials 

golang package implement PerRPCCredentials to reach the token verification. you can use this method to verity your request.

```go
type PerRPCCredentials interface {
    GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
    RequireTransportSecurity() bool
} 
```
## Build

```
./scripts/genTLS.sh
make
```