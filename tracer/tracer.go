package tracer

type contextKey string

const IrisContextKey contextKey = "IrisContextKey"
const PanicContextKey contextKey = "PanicContextKey"
const TracingRequestKey contextKey = "TracingRequestKey"
