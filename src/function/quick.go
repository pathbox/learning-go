func(){
  fmt.Println("Hello")
}

fn := func(){
  fmt.Println("Hello")
}

x :=5
fn := func(){
  fmt.Println("x is", x)
}

fn()
x++
fn()

type op struct{
  name string
  fn func(int, int) int
}

func main() {
  rand.Seed(time.Now().Unix())

  ops := []op{
        {"add",func(x, y int) int {return x+y}},
        {"sub",func(x, y int) int {return x-y}},
        {"mul" func(x, y int) int {return x*y}},
        {"div" func(x, y int) int {return x/y}},
        {"mod" func(x, y int) int {return x%y}},
    }

  o := ops[rand.Intn(len(ops))]
  x, y := 12, 5
  fmt.Println(o.name, x, y)
  fmt.Println(o.fn(x, y))
}

type walkFn func(*int) walkFn

func walkEqual(i *int) walkFn{
  *i += rand.Intn(7) - 3
}

var fnRegistry = map[string]binFunc{
  "add": func(x, y int) int { return x + y},
  "sub": func(x, y int) int { return x - y},
  "mul": func(x, y int) int { return x * y },
  "div": func(x, y int) int { return x / y },
  "mod": func(x, y int) int { return x % y },
}

