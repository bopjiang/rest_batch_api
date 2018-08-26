package main

func main() {
	s := NewServer()
	s.Addr = "127.0.0.1:8080"
	s.Start()
}
