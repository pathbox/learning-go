package main

type Foo struct {
  Bar string
}

func (f *Foo) MarshalJSON() ([]byte, error){
  return []byte(`{"Bar":`+ f.Bar +`}`), nil
}
