# conf
Conf's goal is supply a handy conf manager tools.

## usage

You can use this package just like this

```go
	if err := conf.New(&configStruct).
		Load(toml.New("fake.toml"), toml.New("example.toml")).
		Over(json.New("example.json")).
		Err(); err != nil {
		// error handler
	}
```

