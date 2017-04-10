type Item struct {
  Id int
  Name string
}

var item Item

res.Body.FromJsonTo(&item)