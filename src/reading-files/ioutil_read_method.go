// Reading an entire file

bytes, err := ioutil.ReadFile("_config.yml")
if err != nil {
  log.Fatal(err)
}

fmt.Println("Bytes read: ", len(bytes))
fmt.Println("String read: ", string(bytes))

// Reading an entire directory of files
filelist, err := ioutil.ReadDir(".")
if err != nil {
  log.Fatal(err)
}

for _, fileinfo := range filelist {
  if fileinfo.Mode().IsRegular() {
    bytes, err := ioutil.ReadFile(fileinfo.Name())
    if err != nil {
      log.Fatal(err)
    }
    fmt.Println("Bytes read: ", len(bytes))
    fmt.Println("String read: ", string(bytes))
  }
}

ioutil.ReadAll() -> Takes an io-like object and returns the entire data as a byte array
io.ReadFull()
io.ReadAtLeast()
io.MultiReader -> A very useful primitive to combine multiple io-like objects.
So you can have a list of files to be read, and treat them as a single contiguous block
of data rather than managing the complexity of switching the file objects at the end of
each of the previous objects.