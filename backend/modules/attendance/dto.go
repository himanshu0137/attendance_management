package attendance

// PunchRequest represents the request to record an attendance punch.
type PunchRequest struct {
	EmployeeID string `json:"employee_id"`
	PunchType  string `json:"punch_type"` // "IN" or "OUT"
	Location   string `json:"location"`
}

// PunchResponse represents a punch response.
type PunchResponse struct {
	ID         string `json:"id"`
	EmployeeID string `json:"employee_id"`
	PunchTime  int64  `json:"punch_time"`
	PunchType  string `json:"punch_type"`
	Location   string `json:"location"`
}

// ReportRequest represents the request for attendance report.
type ReportRequest struct {
	EmployeeID string `json:"employee_id"`
	FromDate   int64  `json:"from_date"`
	ToDate     int64  `json:"to_date"`
}

// ReportResponse represents an attendance report response.
type ReportResponse struct {
	ID           string  `json:"id"`
	EmployeeID   string  `json:"employee_id"`
	ReportDate   int64   `json:"report_date"`
	TotalHours   float64 `json:"total_hours"`
	Status       string  `json:"status"`
	Punches      []PunchResponse `json:"punches"`
}