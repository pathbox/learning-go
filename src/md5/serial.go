package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)


// MD5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, MD5All returns an error.

func MD5All(root string) (map[string][md5.Size]byte, error) {
  m := make(map[string][md5.Size]byte)
  err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
    if err != nil{
      return err
    }
    if !info.Mode().IsRegular(){
      return nil
    }
    data, err := ioutil.ReadFile(path)
    if err != nil {
      return nil
    }
    m[path] = md5.Sum(data)
    return nil
  })
  if err != nil {
    return nil, err
  }
  return m, nil
}

func main() {
  // Calculate the MD5 sum of all files under the specified directory,
	// then print the results sorted by path name.
  m, err := MD5All(os.Args[1])
  if err != nil {
		fmt.Println(err)
		return
	}

  var paths []string
  for path := range m{
    paths = append(paths, path)
  }
  sort.String(paths)
  for _, path := range paths{
    f,t.Printf("%x %s\n", m[path], path)
  }
}
