# Bhojpur DBM - KeyValue Database

It is an simple implementation of `Key/Value` databases.

## Installation

```go
	go get github.com/bhojpur/dbm/pkg/keyvalue
```

## Requirements

* You need at least `Go` 1.17 or newer.

## Library Usage

Firstly, you need to instantiate the KeyValue database

### Create or open a KeyValue database

```go
// The returned KeyValue DB instance is safe for concurrent use. It means that all
// DB's methods may be called concurrently from multiple goroutines.
db, err := keyvalue.OpenFile("path/to/db", nil)
...
defer db.Close()
...
```

### Read or modify the database content

```go
// Remember that the contents of the returned slice should not be modified.
data, err := db.Get([]byte("key"), nil)
...
err = db.Put([]byte("key"), []byte("value"), nil)
...
err = db.Delete([]byte("key"), nil)
...
```

### Iterate over the database content

```go
iter := db.NewIterator(nil, nil)
for iter.Next() {
	// Remember that the contents of the returned slice should not be modified, and
	// only valid until the next call to Next.
	key := iter.Key()
	value := iter.Value()
	...
}
iter.Release()
err = iter.Error()
...
```

### Seek-then-Iterate

```go
iter := db.NewIterator(nil, nil)
for ok := iter.Seek(key); ok; ok = iter.Next() {
	// Use key/value.
	...
}
iter.Release()
err = iter.Error()
...
```

### Iterate over subset of database content

```go
iter := db.NewIterator(&util.Range{Start: []byte("foo"), Limit: []byte("xoo")}, nil)
for iter.Next() {
	// Use key/value.
	...
}
iter.Release()
err = iter.Error()
...
```

### Iterate over subset of database content with a particular prefix

```go
iter := db.NewIterator(util.BytesPrefix([]byte("foo-")), nil)
for iter.Next() {
	// Use key/value.
	...
}
iter.Release()
err = iter.Error()
...
```

### Batch writes

```go
batch := new(keyvalue.Batch)
batch.Put([]byte("foo"), []byte("value"))
batch.Put([]byte("bar"), []byte("another value"))
batch.Delete([]byte("baz"))
err = db.Write(batch, nil)
...
```

### Use the Bloom filter

```go
o := &opt.Options{
	Filter: filter.NewBloomFilter(10),
}
db, err := keyvalue.OpenFile("path/to/db", o)
...
defer db.Close()
...
```
