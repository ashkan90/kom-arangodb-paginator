## Next

- [x] Paginate with keeping original data type safe. 
- [ ] Paginate with no data type concern.
- [ ] Paginate result set (example is below).
- [ ] Show query option.
- [x] Export Mock functions which `arangodb` driver uses.
- [ ] Mock test via `arangodb` driver.
- [ ] Integration test via `arangodb` driver.

## Paginate Result Sets

### Type Safe Result Set

```go
type PaginationSafeResult struct {
    CurrentPage  int
    PrevPage     int
    NextPage     int
    TotalPage    int
    TotalRecords int64
}
```

### Classic Pagination Result Set

```go
type PaginationResult struct {
	Result PaginationSafeResult
	Data   []interface{}
}
```