package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func generateRequestID() string {
	bytes := make([]byte, 10)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func logWithRequest(r *http.Request, level slog.Level, msg string, args ...any) {
	requestID, _ := r.Context().Value("request_id").(string)
	if requestID == "" {
		requestID = "unknown"
	}

	baseArgs := []any{
		"request_id", requestID,
		"method", r.Method,
		"url", r.URL.String(),
		"remote_addr", r.RemoteAddr,
		"user_agent", r.Header.Get("User-Agent"),
	}

	allArgs := append(baseArgs, args...)
	logger.Log(context.Background(), level, msg, allArgs...)
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := generateRequestID()
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		r = r.WithContext(ctx)
		logWithRequest(r, slog.LevelInfo, "Received request")

		start := time.Now()
		next(w, r)
		duration := time.Since(start)
		logWithRequest(r, slog.LevelInfo, fmt.Sprintf("Request completed in %d ms", duration.Milliseconds()))
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprint("RECEIVED RESPONSE: ", r.URL.Path, r.Method, r.Body)))

}

func PostUploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errMsg := fmt.Sprintf("Method %s is not allowed for endpoint %s", r.Method, r.URL.String())
		logWithRequest(r, slog.LevelWarn, errMsg)
		http.Error(w, errMsg, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		errMsg := "Unable to parse multipart form"
		logWithRequest(r, slog.LevelError, errMsg, "error", err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		errMsg := "Unable to retrieve file from form"
		logWithRequest(r, slog.LevelError, errMsg, "error", err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	defer file.Close()

	logWithRequest(r, slog.LevelInfo, "Received file",
		"original_filename", header.Filename,
		"file_size_in_bytes", header.Size,
		"content_type", header.Header.Get("Content-Type"))

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)

	dst, err := os.Create("uploads/" + filename)
	if err != nil {
		errMsg := "Unable to create file"
		logWithRequest(r, slog.LevelError, errMsg,
			"filename", filename,
			"error", err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	bytesWritten, err := io.Copy(dst, file)
	if err != nil {
		errMsg := "Unable to save file"
		logWithRequest(r, slog.LevelError, errMsg,
			"filename", filename,
			"error", err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf(`{"status": "success", "filename": "%s", "size": %d}`, filename, bytesWritten)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func main() {

	http.HandleFunc("/check", get)
	http.HandleFunc("/upload_image", loggingMiddleware(PostUploadImage))

	if err := http.ListenAndServe(":8081", nil); err != nil {
		logger.Error("Server failed to start", "error", err.Error())
	}
}
