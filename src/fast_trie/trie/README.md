https://github.com/derekparker/trie

## Usage

Create a Trie with:

```Go
t := trie.New()
```

Add Keys with:

```Go
// Add can take in meta information which can be stored with the key.
// i.e. you could store any information you would like to associate with
// this particular key.
t.Add("foobar", 1)
```

Find a key with:

```Go
node, ok := t.Find("foobar")
meta := node.Meta()
// use meta with meta.(type)
```

Remove Keys with:

```Go
t.Remove("foobar")
```

Prefix search with:

```Go
t.PrefixSearch("foo")
```

Fast test for valid prefix:
```Go
t.HasKeysWithPrefix("foo")
```

Fuzzy search with:

```Go
t.FuzzySearch("fb")
```