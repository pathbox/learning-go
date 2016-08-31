package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// A result is the product of reading and summing a file using MD5.

type result struc {
  path string
  sum [md5.Size]byte
  err error
}

// sumFiles starts goroutines to walk the directory tree at root and digest each
// regular file.  These goroutines send the results of the digests on the result
// channel and send the result of the walk on the error channel.  If done is
// closed, sumFiles abandons its work.

func sumFiles(done <-chan struct{}, root string)(<-chan result, <-chan error){
  // For each regular file, start a goroutine that sums the file and sends
	// the result on c.  Send the result of the walk on errc.
  c := make(chan result)
  errc := make(chan error, 1)
  go func(){
    var wg sync.WaitGroup
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
      if err != nil{
        return err
      }
      if !info.Mode().IsRegular(){
        return nil
      }
      wg.Add(1)
      go func ()  {
        data, err := ioutil.ReadFile(path)
        select {
        case c<- result{path, md5.Sum(data), err}:
        case <-done:
        }
        wg.Done()
      }()
      select{
      case <-done:
        reutrn errors.New("Walk canceled")
      default:
        return nil
      }
    })
    go func(){
      wg.Wait()
      close(c)
    }()
    errc <-err
  }()
  return c, errc
}

// MD5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, MD5All returns an error.  In that case,
// MD5All does not wait for inflight read operations to complete.

func MD5All(root string) (map[string][md5.Size]byte, error) {
	// MD5All closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc.
	done := make(chan struct{}) // HLdone
	defer close(done)           // HLdone

	c, errc := sumFiles(done, root) // HLdone

	m := make(map[string][md5.Size]byte)
	for r := range c { // HLrange
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	if err := <-errc; err != nil {
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
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x  %s\n", m[path], path)
	}
}
