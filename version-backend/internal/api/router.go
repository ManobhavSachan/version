package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"version-backend/internal/api/handlers"
	"version-backend/internal/api/middleware"
	"version-backend/internal/db"

	"github.com/gorilla/mux"
)

// Router represents the HTTP router
type Router struct {
	*mux.Router
	startTime time.Time
	db        *db.DB
}

// enableCORS adds CORS middleware to allow frontend requests
func (r *Router) enableCORS() {
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if req.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, req)
		})
	})
}

// NewRouter creates a new HTTP router with all routes configured
func NewRouter(db *db.DB) *Router {
	r := mux.NewRouter()
	router := &Router{
		Router:    r,
		startTime: time.Now(),
		db:        db,
	}

	// Enable CORS for all routes
	router.enableCORS()

	// Add middleware
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)
	r.Use(middleware.WithDB(db))

	// Welcome page with ASCII art
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		welcome := `
                    ___           ___           ___                        ___           ___     
      ___          /  /\         /  /\         /  /\           ___        /  /\         /  /\    
     /  /\        /  /::\       /  /::\       /  /::\         /__/\      /  /::\       /  /::|   
    /  /:/       /  /:/\:\     /  /:/\:\     /__/:/\:\        \__\:\    /  /:/\:\     /  /:|:|   
   /  /:/       /  /::\ \:\   /  /::\ \:\   _\_ \:\ \:\       /  /::\  /  /:/  \:\   /  /:/|:|__ 
  /__/:/  ___  /__/:/\:\ \:\ /__/:/\:\_\:\ /__/\ \:\ \:\   __/  /:/\/ /__/:/ \__\:\ /__/:/ |:| /\
  |  |:| /  /\ \  \:\ \:\_\/ \__\/~|::\/:/ \  \:\ \:\_\/  /__/\/:/~~  \  \:\ /  /:/ \__\/  |:|/:/
  |  |:|/  /:/  \  \:\ \:\      |  |:|::/   \  \:\_\:\    \  \::/      \  \:\  /:/      |  |:/:/ 
  |__|:|__/:/    \  \:\_\/      |  |:|\/     \  \:\/:/     \  \:\       \  \:\/:/       |__|::/  
   \__\::::/      \  \:\        |__|:|~       \  \::/       \__\/        \  \::/        /__/:/   
       ~~~~        \__\/         \__\|         \__\/                      \__\/         \__\/      
                                                                    
 =================================================================
                   System Information Monitor
 =================================================================

 Available Endpoints:
 -------------------
 GET /api/latest_data  -> Returns system information
                         • OS Version
                         • Osquery Version
                         • Installed Applications

 System Status:
 -------------
 GET /health          -> Basic health check
 GET /status          -> Detailed system status

 =================================================================
 Server running on port 7070 | Made with ❤️  using Go & Osquery
 =================================================================
`
		fmt.Fprint(w, welcome)
	}).Methods(http.MethodGet)

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/latest_data", handlers.GetLatestData).Methods(http.MethodGet)

	// Health check - simple endpoint for load balancers
	r.HandleFunc("/health", router.handleHealth).Methods(http.MethodGet)

	// Status - detailed system status
	r.HandleFunc("/status", router.handleStatus).Methods(http.MethodGet)

	return router
}

// Run starts the HTTP server
func (r *Router) Run(addr string) error {
	return http.ListenAndServe(addr, r)
}

// handleHealth handles the basic health check endpoint
func (r *Router) handleHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check database connection
	err := r.db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "unhealthy",
			"error":  "Database connection failed",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

// handleStatus handles the detailed status endpoint
func (r *Router) handleStatus(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Check database connection
	dbStatus := "connected"
	if err := r.db.Ping(); err != nil {
		dbStatus = "disconnected"
	}

	// Get latest system info timestamp
	var lastUpdate string
	err := r.db.Get(&lastUpdate, "SELECT updated_at FROM system_info ORDER BY updated_at DESC LIMIT 1")
	if err != nil {
		lastUpdate = "no data collected yet"
	}

	status := map[string]interface{}{
		"status": map[string]interface{}{
			"uptime":     time.Since(r.startTime).String(),
			"startTime":  r.startTime.Format(time.RFC3339),
			"goroutines": runtime.NumGoroutine(),
			"memory": map[string]interface{}{
				"alloc":      fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024),
				"totalAlloc": fmt.Sprintf("%.2f MB", float64(m.TotalAlloc)/1024/1024),
				"sys":        fmt.Sprintf("%.2f MB", float64(m.Sys)/1024/1024),
				"numGC":      m.NumGC,
			},
		},
		"database": map[string]interface{}{
			"status":           dbStatus,
			"lastDataUpdate":   lastUpdate,
			"maxOpenConns":     r.db.Stats().MaxOpenConnections,
			"openConnections":  r.db.Stats().OpenConnections,
			"inUseConnections": r.db.Stats().InUse,
			"idleConnections":  r.db.Stats().Idle,
		},
		"build": map[string]interface{}{
			"goVersion": runtime.Version(),
			"os":        runtime.GOOS,
			"arch":      runtime.GOARCH,
			"cpus":      runtime.NumCPU(),
		},
	}

	json.NewEncoder(w).Encode(status)
}
