project = "test"

sources = [
  "remote"
]

http "remote" {
  endpoint = "https://raw.githubusercontent.com/boljen/hostman/refs/heads/master/examples/http/response.json"
}