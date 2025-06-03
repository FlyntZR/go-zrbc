package interceptor

import (
	"context"
	"fmt"
	"time"

	"go-zrbc/pkg/xlog"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
)

var (
	rpcReqQPS = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "rpc_request_qps_total",
		Help: "The total number of processed events",
	}, []string{
		"service_id",
		"method",
		"err_status",
	})

	httpResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "rpc_resptime",
		Help: "The total number of processed events",
	}, []string{
		"service_id",
		"method",
	})
)

type Interceptor struct {
	ServiceId string
}

func NewInterceptor(ServiceId string) *Interceptor {
	return &Interceptor{ServiceId: ServiceId}
}

func (inter *Interceptor) GrpcLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		guuid := uuid.New().String()
		ctx = context.WithValue(ctx, "uuid", guuid)

		defer func() {
			rpcReqQPS.With(prometheus.Labels{
				"service_id": inter.ServiceId,
				"method":     info.FullMethod,
				"err_status": fmt.Sprintf("%v", err),
			}).Inc()

			httpResponseTime.With(prometheus.Labels{
				"service_id": inter.ServiceId,
				"method":     info.FullMethod,
			}).Observe(float64(time.Since(start) / time.Millisecond))

			xlog.Debugf("method:%s uuid(%v),  takes(%v ms), req(%v), response(%+v), err(%v)\n",
				info.FullMethod, guuid, int64(time.Since(start)/time.Millisecond), req, resp, err)
		}()

		resp, err = handler(ctx, req)

		return resp, err
	}
}

func (inter *Interceptor) GrpcLoggerStream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		start := time.Now()
		guuid := uuid.New().String()

		defer func() {
			rpcReqQPS.With(prometheus.Labels{
				"service_id": inter.ServiceId,
				"method":     info.FullMethod,
				"err":        fmt.Sprintf("%v", err),
			}).Inc()

			httpResponseTime.With(prometheus.Labels{
				"service_id": inter.ServiceId,
				"method":     info.FullMethod,
			}).Observe(float64(time.Since(start) / time.Millisecond))

			xlog.Debugf("method:%s uuid(%v),  takes(%v ms), err(%v)\n",
				info.FullMethod, guuid, int64(time.Since(start)/time.Millisecond), err)
		}()

		err = handler(srv, ss)
		return err
	}
}

func (inter *Interceptor) GrpcRecover() grpc.UnaryServerInterceptor {
	var customFunc recovery.RecoveryHandlerFunc
	customFunc = func(p interface{}) (err error) {
		xlog.Debug("server err: %+v\n", p)
		return nil
	}

	recovery.WithRecoveryHandler(customFunc)
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(customFunc),
	}

	return recovery.UnaryServerInterceptor(opts...)
}

func (inter *Interceptor) GrpcRecoverStream() grpc.StreamServerInterceptor {
	var customFunc recovery.RecoveryHandlerFunc
	customFunc = func(p interface{}) (err error) {
		xlog.Debug("server err: %+v\n", p)
		return nil
	}

	recovery.WithRecoveryHandler(customFunc)
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(customFunc),
	}

	return recovery.StreamServerInterceptor(opts...)
}
