package hostman

import "testing"

func TestStaticSource_Validate_Success_WithHost(t *testing.T) {
    h := "example.com"
    s := &StaticSource{ Name: "s1", Host: &h, Ip: "127.0.0.1" }
    if err := s.Validate(); err != nil {
        t.Fatalf("Validate unexpected error: %v", err)
    }
}

func TestStaticSource_Validate_Success_WithHosts(t *testing.T) {
    hs := []string{"a.example.com", "b.example.com"}
    s := &StaticSource{ Name: "s2", Hosts: &hs, Ip: "127.0.0.1" }
    if err := s.Validate(); err != nil {
        t.Fatalf("Validate unexpected error: %v", err)
    }
}

func TestStaticSource_Validate_Fails_WhenNeitherHostNorHosts(t *testing.T) {
    s := &StaticSource{ Name: "s3", Ip: "127.0.0.1" }
    if err := s.Validate(); err == nil {
        t.Fatalf("expected error when neither host nor hosts set")
    }
}

func TestStaticSource_Validate_Fails_WhenBothHostAndHosts(t *testing.T) {
    h := "example.com"
    hs := []string{"a.example.com"}
    s := &StaticSource{ Name: "s4", Host: &h, Hosts: &hs, Ip: "127.0.0.1" }
    if err := s.Validate(); err == nil {
        t.Fatalf("expected error when both host and hosts set")
    }
}

func TestStaticSource_Validate_Fails_WhenHostEmpty(t *testing.T) {
    h := ""
    s := &StaticSource{ Name: "s5", Host: &h, Ip: "127.0.0.1" }
    if err := s.Validate(); err == nil {
        t.Fatalf("expected error when host is empty")
    }
}

func TestStaticSource_Validate_Fails_WhenHostsEmpty(t *testing.T) {
    hs := []string{}
    s := &StaticSource{ Name: "s6", Hosts: &hs, Ip: "127.0.0.1" }
    if err := s.Validate(); err == nil {
        t.Fatalf("expected error when hosts is empty")
    }
}

func TestStaticSource_Validate_Fails_WhenIpEmpty(t *testing.T) {
    h := "example.com"
    s := &StaticSource{ Name: "s7", Host: &h, Ip: "" }
    if err := s.Validate(); err == nil {
        t.Fatalf("expected error when ip is empty")
    }
}
