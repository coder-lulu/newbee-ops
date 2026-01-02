package jsonconv

import "testing"

func TestParseStringSlice(t *testing.T) {
    tests := []struct{ in string; want int; ok bool }{
        {"[\"ssh\",\"rdp\"]", 2, true},
        {"[]", 0, true},
        {"", 0, true},
        {"not json", 0, false},
    }
    for _, tt := range tests {
        got, err := ParseStringSlice(tt.in)
        if (err == nil) != tt.ok { t.Fatalf("ParseStringSlice(%q) err=%v", tt.in, err) }
        if len(got) != tt.want { t.Fatalf("len=%d, want %d", len(got), tt.want) }
    }
}

func TestParseStringIntMap(t *testing.T) {
    m, err := ParseStringIntMap("{\"ssh\":22,\"rdp\":3389}")
    if err != nil { t.Fatalf("err: %v", err) }
    if m["ssh"] != 22 || m["rdp"] != 3389 { t.Fatalf("unexpected %v", m) }
}

func TestParseStringMap(t *testing.T) {
    m, err := ParseStringMap("{\"env\":\"prod\"}")
    if err != nil { t.Fatalf("err: %v", err) }
    if m["env"] != "prod" { t.Fatalf("unexpected %v", m) }
}

func TestMustMarshal(t *testing.T) {
    s := MustMarshal([]string{"ssh"})
    if s == "" { t.Fatalf("want non-empty json") }
}

func TestValidate(t *testing.T) {
    if !ValidateCapabilities([]string{"ssh","rdp"}) { t.Fatalf("caps valid expected") }
    if ValidateCapabilities([]string{"foo"}) { t.Fatalf("caps invalid expected") }
    if !ValidatePorts(map[string]int{"ssh":22}) { t.Fatalf("ports valid expected") }
    if ValidatePorts(map[string]int{"ssh":70000}) { t.Fatalf("port range invalid expected") }
    if ValidatePorts(map[string]int{"foo":22}) { t.Fatalf("protocol key invalid expected") }
}
