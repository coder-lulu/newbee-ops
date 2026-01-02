package jsonconv

import "encoding/json"

var allowedProtocols = map[string]struct{}{
    "ssh": {},
    "telnet": {},
    "rdp": {},
    "vnc": {},
}

// ParseStringSlice parses a JSON array string into []string. Empty or invalid returns empty slice and error.
func ParseStringSlice(s string) ([]string, error) {
    if s == "" { return []string{}, nil }
    var out []string
    err := json.Unmarshal([]byte(s), &out)
    if err != nil { return []string{}, err }
    return out, nil
}

// ParseStringIntMap parses a JSON object string into map[string]int.
func ParseStringIntMap(s string) (map[string]int, error) {
    if s == "" { return map[string]int{}, nil }
    var out map[string]int
    err := json.Unmarshal([]byte(s), &out)
    if err != nil { return map[string]int{}, err }
    return out, nil
}

// ParseStringMap parses a JSON object string into map[string]string.
func ParseStringMap(s string) (map[string]string, error) {
    if s == "" { return map[string]string{}, nil }
    var out map[string]string
    err := json.Unmarshal([]byte(s), &out)
    if err != nil { return map[string]string{}, err }
    return out, nil
}

// MustMarshal marshals v to JSON string, returns empty string on error.
func MustMarshal(v interface{}) string {
    b, err := json.Marshal(v)
    if err != nil { return "" }
    return string(b)
}

// ValidateCapabilities checks that all caps are within allowed protocols.
func ValidateCapabilities(caps []string) bool {
    for _, c := range caps {
        if _, ok := allowedProtocols[c]; !ok { return false }
    }
    return true
}

// ValidatePorts checks that keys are allowed protocols and values are in 1..65535
func ValidatePorts(ports map[string]int) bool {
    for k, v := range ports {
        if _, ok := allowedProtocols[k]; !ok { return false }
        if v <= 0 || v > 65535 { return false }
    }
    return true
}
