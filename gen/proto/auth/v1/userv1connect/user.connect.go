// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: proto/auth/v1/user.proto

package userv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	http "net/http"
	strings "strings"
	v1 "task-manager/gen/proto/auth/v1"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion0_1_0

const (
	// AuthServiceName is the fully-qualified name of the AuthService service.
	AuthServiceName = "proto.auth.v1.AuthService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// AuthServiceSignUpProcedure is the fully-qualified name of the AuthService's SignUp RPC.
	AuthServiceSignUpProcedure = "/proto.auth.v1.AuthService/SignUp"
	// AuthServiceSignInProcedure is the fully-qualified name of the AuthService's SignIn RPC.
	AuthServiceSignInProcedure = "/proto.auth.v1.AuthService/SignIn"
)

// AuthServiceClient is a client for the proto.auth.v1.AuthService service.
type AuthServiceClient interface {
	SignUp(context.Context, *connect.Request[v1.SignUpRequest]) (*connect.Response[v1.SignUpResponse], error)
	SignIn(context.Context, *connect.Request[v1.SignInRequest]) (*connect.Response[v1.SignInResponse], error)
}

// NewAuthServiceClient constructs a client for the proto.auth.v1.AuthService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAuthServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) AuthServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &authServiceClient{
		signUp: connect.NewClient[v1.SignUpRequest, v1.SignUpResponse](
			httpClient,
			baseURL+AuthServiceSignUpProcedure,
			opts...,
		),
		signIn: connect.NewClient[v1.SignInRequest, v1.SignInResponse](
			httpClient,
			baseURL+AuthServiceSignInProcedure,
			opts...,
		),
	}
}

// authServiceClient implements AuthServiceClient.
type authServiceClient struct {
	signUp *connect.Client[v1.SignUpRequest, v1.SignUpResponse]
	signIn *connect.Client[v1.SignInRequest, v1.SignInResponse]
}

// SignUp calls proto.auth.v1.AuthService.SignUp.
func (c *authServiceClient) SignUp(ctx context.Context, req *connect.Request[v1.SignUpRequest]) (*connect.Response[v1.SignUpResponse], error) {
	return c.signUp.CallUnary(ctx, req)
}

// SignIn calls proto.auth.v1.AuthService.SignIn.
func (c *authServiceClient) SignIn(ctx context.Context, req *connect.Request[v1.SignInRequest]) (*connect.Response[v1.SignInResponse], error) {
	return c.signIn.CallUnary(ctx, req)
}

// AuthServiceHandler is an implementation of the proto.auth.v1.AuthService service.
type AuthServiceHandler interface {
	SignUp(context.Context, *connect.Request[v1.SignUpRequest]) (*connect.Response[v1.SignUpResponse], error)
	SignIn(context.Context, *connect.Request[v1.SignInRequest]) (*connect.Response[v1.SignInResponse], error)
}

// NewAuthServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAuthServiceHandler(svc AuthServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	authServiceSignUpHandler := connect.NewUnaryHandler(
		AuthServiceSignUpProcedure,
		svc.SignUp,
		opts...,
	)
	authServiceSignInHandler := connect.NewUnaryHandler(
		AuthServiceSignInProcedure,
		svc.SignIn,
		opts...,
	)
	return "/proto.auth.v1.AuthService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AuthServiceSignUpProcedure:
			authServiceSignUpHandler.ServeHTTP(w, r)
		case AuthServiceSignInProcedure:
			authServiceSignInHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAuthServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAuthServiceHandler struct{}

func (UnimplementedAuthServiceHandler) SignUp(context.Context, *connect.Request[v1.SignUpRequest]) (*connect.Response[v1.SignUpResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("proto.auth.v1.AuthService.SignUp is not implemented"))
}

func (UnimplementedAuthServiceHandler) SignIn(context.Context, *connect.Request[v1.SignInRequest]) (*connect.Response[v1.SignInResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("proto.auth.v1.AuthService.SignIn is not implemented"))
}
