package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/event"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SSEHandler struct {
	manager    *event.Manager
	userClient userpb.UserServiceClient
}

func NewSSEHandler(manager *event.Manager, userClient userpb.UserServiceClient) *SSEHandler {
	return &SSEHandler{
		manager,
		userClient,
	}
}

func (h *SSEHandler) HandleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "yêu cầu user_id để kết nối", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if _, err := h.userClient.GetUserPublicById(ctx, &userpb.GetOneRequest{Id: userID}); err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				http.Error(w, st.Message(), http.StatusBadRequest)
			default:
				http.Error(w, st.Message(), http.StatusInternalServerError)
			}
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := event.NewUserReceiverEvent(userID)

	h.manager.Register <- user

	defer func() {
		h.manager.Unregister <- user
	}()

	fmt.Fprintf(w, "data: {\"event\":\"connected\",\"message\":\"SSE connection established\"}\n\n")
	w.(http.Flusher).Flush()

	clientGone := r.Context().Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message := <-user.Send:
			fmt.Fprintf(w, "data: %s\n\n", string(message))
			w.(http.Flusher).Flush()

		case <-user.Done:
			return

		case <-clientGone:
			return

		case <-ticker.C:
			fmt.Fprintf(w, ": keep-alive\n\n")
			w.(http.Flusher).Flush()
		}
	}
}
