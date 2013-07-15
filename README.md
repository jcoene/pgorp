Provides conversion of Postgres types for [gorp](https://github.com/coopernurse/gorp).

Currently only implements []int64 as ArrayInt64. Patches welcome.

## Usage

During gorp initialization, set the instance TypeConverter to pgorp's TypeConverter:

```go
orm = new(gorp.DbMap)
...
orm.TypeConverter = pgorp.TypeConverter{}
```

Later on, use the pgorp.ArrayInt64 type for integer array column definitions:

```go
type SchoolDistrict {
  Name string
  SchoolIds pgorp.ArrayInt64 `db:"school_ids"`
}
```

## Credits

Based largely on work done by Andrew Harris in [this gist](https://gist.github.com/adharris/4163702), referenced in [gorp issue #5](https://github.com/coopernurse/gorp/issues/5).

## License

MIT License
