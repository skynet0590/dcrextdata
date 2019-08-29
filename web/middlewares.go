package web

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
)

const (
	ctxSyncDataType = iota
)

func syncDataType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxSyncDataType,
			chi.URLParam(r, "dataType"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getSyncTypeCtx retrieves the syncDataType data from the request context.
// If not set, the return value is an empty string.
func getSyncDataTypeCtx(r *http.Request) string {
	syncType, ok := r.Context().Value(ctxSyncDataType).(string)
	if !ok {
		log.Trace("sync type not set")
		return ""
	}
	return syncType
}
