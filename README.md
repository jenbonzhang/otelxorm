# otelxorm

`otelxorm` is a hook for [xorm](https://github.com/go-xorm/xorm) that enables OpenTelemetry tracing for database operations.

## Installation

To install `otelxorm`, use `go get`:

```bash
go get github.com/jenbonzhang/otelxorm
```


## Usage

To use `otelxorm`, first create an instance of `xorm.Engine`:

```go
engine, err := xorm.NewEngine(dbType, dsn)
if err != nil {
    logrus.Errorf("Database connection failed err: %v. Database name: %s", err, name)
    panic(err)
}
```

Then, add the `otelxorm` hook to the engine:

```go
engine.AddHook(otelxorm.Hook(
    otelxorm.WithDBName(name),       // Set the database name
    otelxorm.WithFormatSQLReplace(), // Set the method for replacing parameters in formatted SQL statements
))
```

This will enable tracing for all database operations performed by the engine.


## Configuration

`otelxorm` provides several options for configuration:

- `WithDBName(name string)`: Sets the name of the database being traced.
- `WithFormatSQL(formatSQL func(sql string, args []interface{}) string) Option `: Sets the method for formatted SQL statements. By default, `otelxorm` uses `otelxorm.defaultFormatSQL` to format SQL statements and  parameters.
- `WithFormatSQLReplace()` this is use args to replace the sql parameters with `$d` in the sql statement.

You can define your own implementation of the `formatSQL` method and pass it to the `WithFormatSQL` option when adding the `otelxorm` hook to the engine. This will override the default implementation used by `otelxorm`. Here's an example:

```go
func myFormatSQL(sql string, args []interface{}) string {
    // Custom implementation for formatting SQL statements and parameters
    // ...
    return formattedQuery
}

// ...

engine.AddHook(otelxorm.Hook(
    otelxorm.WithDBName(name),
    otelxorm.WithFormatSQL(myFormatSQL),
))
```

This will enable tracing for all database operations performed by the engine using your custom implementation of the `formatSQL` method.
## Contributing

We welcome contributions! Please see our [contributing guidelines](CONTRIBUTING.md) for more information.

## License

`otelxorm` is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for more information.