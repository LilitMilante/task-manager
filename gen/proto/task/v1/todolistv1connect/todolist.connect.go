// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: proto/task/v1/todolist.proto

package todolistv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	http "net/http"
	strings "strings"
	v1 "task-manager/gen/proto/task/v1"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion0_1_0

const (
	// TaskServiceName is the fully-qualified name of the TaskService service.
	TaskServiceName = "proto.task.v1.TaskService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// TaskServiceAddTaskProcedure is the fully-qualified name of the TaskService's AddTask RPC.
	TaskServiceAddTaskProcedure = "/proto.task.v1.TaskService/AddTask"
)

// TaskServiceClient is a client for the proto.task.v1.TaskService service.
type TaskServiceClient interface {
	AddTask(context.Context, *connect.Request[v1.AddTaskRequest]) (*connect.Response[v1.AddTaskResponse], error)
}

// NewTaskServiceClient constructs a client for the proto.task.v1.TaskService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewTaskServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) TaskServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &taskServiceClient{
		addTask: connect.NewClient[v1.AddTaskRequest, v1.AddTaskResponse](
			httpClient,
			baseURL+TaskServiceAddTaskProcedure,
			opts...,
		),
	}
}

// taskServiceClient implements TaskServiceClient.
type taskServiceClient struct {
	addTask *connect.Client[v1.AddTaskRequest, v1.AddTaskResponse]
}

// AddTask calls proto.task.v1.TaskService.AddTask.
func (c *taskServiceClient) AddTask(ctx context.Context, req *connect.Request[v1.AddTaskRequest]) (*connect.Response[v1.AddTaskResponse], error) {
	return c.addTask.CallUnary(ctx, req)
}

// TaskServiceHandler is an implementation of the proto.task.v1.TaskService service.
type TaskServiceHandler interface {
	AddTask(context.Context, *connect.Request[v1.AddTaskRequest]) (*connect.Response[v1.AddTaskResponse], error)
}

// NewTaskServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewTaskServiceHandler(svc TaskServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	taskServiceAddTaskHandler := connect.NewUnaryHandler(
		TaskServiceAddTaskProcedure,
		svc.AddTask,
		opts...,
	)
	return "/proto.task.v1.TaskService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case TaskServiceAddTaskProcedure:
			taskServiceAddTaskHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedTaskServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedTaskServiceHandler struct{}

func (UnimplementedTaskServiceHandler) AddTask(context.Context, *connect.Request[v1.AddTaskRequest]) (*connect.Response[v1.AddTaskResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("proto.task.v1.TaskService.AddTask is not implemented"))
}