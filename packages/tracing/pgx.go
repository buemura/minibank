package tracing

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// PgxTracer implements pgx.QueryTracer to create OpenTelemetry spans for database queries.
type PgxTracer struct {
	tracer trace.Tracer
}

func NewPgxTracer() *PgxTracer {
	return &PgxTracer{
		tracer: otel.Tracer("pgx"),
	}
}

type traceQueryCtxKey struct{}

func (t *PgxTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ctx, span := t.tracer.Start(ctx, "db.query",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.statement", data.SQL),
	)
	return context.WithValue(ctx, traceQueryCtxKey{}, span)
}

func (t *PgxTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	span, ok := ctx.Value(traceQueryCtxKey{}).(trace.Span)
	if !ok {
		return
	}
	if data.Err != nil {
		span.RecordError(data.Err)
		span.SetStatus(codes.Error, data.Err.Error())
	}
	span.SetAttributes(
		attribute.String("db.command_tag", data.CommandTag.String()),
	)
	span.End()
}
