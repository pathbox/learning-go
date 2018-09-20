package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	sql "gopkg.in/src-d/go-mysql-server.v0/sql"
)

type database struct {
	path   string
	tables map[string]sql.Table
}

func NewDatabse(dir string) (sql.Databse, error) {
	fis, err := ioutil.ReadDir(dir) // read files from dir
	if err != nil {
		return nil, fmt.Errorf("could not read directory %s: %v", dir, err)
	}

	tables := make(map[string]sql.Table)
	for _, fi := range fis {
		name := fi.Name()
		if filepath.Ext(name) != ".csv" { // if file is csv
			continue
		}

		t, err := NewTable(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		tables[t.Name()] = t
	}
	return &database{path: dir, tables: tables}, nil
}

func (db *database) Name() string                 { return db.path }
func (db *database) Tables() map[string]sql.Table { return db.tables }

func NewTable(path string) (sql.Table, error) {
	name := strings.TrimSuffix(filepath.Base(path), ".csv")
	t := &table{name: name, path: path}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %v", path, err)
	}

	defer f.Close()

	cr := csv.NewReader(f)
	cols, err := cr.Read()
	if err != nil {
		return nil, err
	}

	for _, col := range cols {
		t.schema = append(t.schema, &sql.Column{
			Name:   strings.ToLower(strings.TrimSpace(col)),
			Type:   sql.Text,
			Source: name,
		})
	}
	return t, nil
}

type table struct {
	name   string
	path   string
	schema []*sql.Column
}

func (t *table) Name() string       { return t.name }
func (t *table) String() string     { return t.path }
func (t *table) Schema() sql.Schema { return t.schema }

func (t *table) Partitions(ctx *sql.Context) (sql.Partitions, error) {
	return &partitionIter{}, nil
}

func (t *table) PartitionRows(ctx *sql.Context, p sql.Partition) (sql.RowIter, error) {
	return newRowIter(t.path)
}

type partitionIter struct{ done bool }

func (p *partitionIter) Close() error { return nil }

func (p *partitionIter) Next() (sql.Partition, error) {
	if p.done {
		return nil, io.EOF
	}
	p.done = true
	return partition{}, nil
}

type partition struct{}

func (p partition) Key() []byte { return []byte("key") }

func newRowIter(path string) (sql.RowIter, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	r.Read() // skip titles
	return &rowIter{f, r}, nil
}

type rowIter struct {
	io.Closer
	*csv.Reader
}

func (r *rowIter) Next() (sql.Row, error) {
	cols, err := r.Read()
	if err != nil {
		return nil, err
	}
	args := make([]interface{}, len(cols))
	for i, col := range cols {
		args[i] = strings.TrimSpace(col)
	}
	return sql.NewRow(args...), err
}
