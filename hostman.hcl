project = "demo2"

sources = [
	"single",
	"multiple"
]

static "single" {
	ip = "127.0.0.1"
	host = "example.com"
}

static "multiple" {
	ip = "127.0.0.1"
	hosts = [
		"a.example.com", 
		"b.example.com"
	]
}

http "remote" {
	endpoint = "https://raw.githubusercontent.com/boljen/hostman/refs/heads/master/examples/http/response.json"
	refresh_interval = 30
}
