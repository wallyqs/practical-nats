package kit

// Location is represents the latitude and longitude pair.
type Location struct {
	// Latitude is the latitude of the user making the request.
	Latitude float64 `json:"lat,omitempty"`

	// Longitude is the longitude of the user making the request.
	Longitude float64 `json:"lng,omitempty"`
}

// DriverAgentRequest is the request sent to the driver.
type DriverAgentRequest struct {
	// Type is the type of agent that is requested.
	Type string `json:"type,omitempty"`

	// Location is the location of the user that is being
	// served the request.
	Location *Location `json:"location,omitempty"`

	// RequestID is the ID from the request.
	RequestID string `json:"request_id,omitempty"`
}

// DriverAgentResponse is the response from the driver.
type DriverAgentResponse struct {
	// ID is the identifier of the driver that will accept
	// the request.
	ID string `json:"driver_id,omitempty"`

	// Error is included in case there was an error handling the
	// request.
	Error string `json:"error,omitempty"`
}
