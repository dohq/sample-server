package main

// Result is json response struct.
type Result struct {
	Message string            `json:"message"`
	Envs    map[string]string `json:"envs,omitempty"`
}
