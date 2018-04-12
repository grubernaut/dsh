DSH - Distributed Execution package for Golang
----------------------------------------------

`dsh` allows for CLI tools to easily create a list of target "nodes" to run a `command` against concurrently.

## Example:

```go
config := dsh.ExecOpts{
  RemoteShell: "ssh",
  RemoteCommand: "sudo pkill docker",
  RemoteUser: "admin",
}

targets := []dsh.Endpoint{
  {
    DisplayName: "my-fancy-host",
    Machine: "10.0.0.3",
  }
  {
    DisplayName: "my-not-fancy-host",
    Machine: "10.0.0.4",
  },
}

if err := config.Execute(targets); err != nil {
  log.Errorf("whoops: %v", err)
}
```

More configuration options can be found inside the `ExecOpts{}` struct, inside `structs.go`.

Obviously the `targets` input to `config.Execute()` can be generated and filtered from any upstream API (Ex: Consul, AWS, DO, GCE, etc...), as this is likely the best path forward instead of hard-coding instance IPs.
