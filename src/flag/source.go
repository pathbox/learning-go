一、FlagSet是该包的核心类型：

type FlagSet struct {
    // Usage is the function called when an error occurs while parsing flags.
    // The field is a function (not a method) that may be changed to point to
    // a custom error handler.
    Usage func()
    name          string
    parsed        bool
    actual        map[string]*Flag    // 存放命令行中实际输入的标志映射
    formal        map[string]*Flag  // 存放该实例可以处理的命令行标志映射
    args          []string // arguments after flags 存放非标志的参数列表（即标志后面的参数）
    exitOnError  bool    // does the program exit if there's an error?
    errorHandling ErrorHandling
    output        io.Writer // nil means stderr; use out() accessor
}

该类型同时提供了一系列的方法集合【MethodSet】，通过该方法集用于可以实现灵活的命令行标志处理。


二、flag包export的变量:CommandLine

// CommandLine is the default set of command-line flags, parsed from os.Args.
// The top-level functions such as BoolVar, Arg, and on are wrappers for the
// methods of CommandLine.
var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

该包提供了一个默认变量：CommandLine，其为FlatSet的一个变量（用面向对象的术语叫做：FlagSet的一个实例）

该flag包export的所有函数本质上都是调用FlagSet类型变量（实例）：CommandLine的方法实现。如下：

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func Int(name string, value int, usage string) *int {
    return CommandLine.Int(name, value, usage)
}



三、flag包支持的标志格式

Command line flag syntax:
        -flag // 代表bool值，相当于-flag=true
        -flag=x
        -flag x  // non-boolean flags only  不支持bool值标志
    One or two minus signs may be used; they are equivalent.  // -flag=value与--flag=value是等效的
    The last form is not permitted for boolean flags because the
    meaning of the command
        cmd -x *
    will change if there is a file called 0, false, etc.  You must
    use the -flag=false form to turn off a boolean flag.


    Flag parsing stops just before the first non-flag argument
    ("-" is a non-flag argument) or after the terminator "--".  // 碰到连续两个"-"号且参数长度为2时则终止标志解析



四、标志绑定相关方法：以绑定int类型为例。

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
// FlagSet提供的绑定int类型标志的方法，无返回值。通过传入int类型指针变量进行绑定，当调用该方法后，会将绑定信息存入FlagSet.formal映射中
func (f *FlagSet) IntVar(p *int, name string, value int, usage string) {
    f.Var(newIntValue(value, p), name, usage)
}
//与上述方法相对应的flag包export的函数：
// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func IntVar(p *int, name string, value int, usage string) {
    CommandLine.Var(newIntValue(value, p), name, usage)
}
//使用示例
var flagvar int
flag.IntVar(&flagvar, "flagname2", 1234, "help message for flagname2")


// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
// FlagSet提供的绑定int类型标志的方法，有返回值，返回int类型指针，当调用该方法后，会将绑定信息存入FlagSet.formal映射中
func (f *FlagSet) Int(name string, value int, usage string) *int {
    p := new(int)
    f.IntVar(p, name, value, usage)
    return p
}
//与上述方法相对应的flag包export的函数：
// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func Int(name string, value int, usage string) *int {
    return CommandLine.Int(name, value, usage)
}
//使用示例
var flagvar = flag.Int("flagname", 1234, "help message for flagname")



五、解析标志的相关关键源码

// Parse parses the command-line flags from os.Args[1:].  Must be called
// after all flags are defined and before flags are accessed by the program.
// flag包export的函数，调用时机为：在设置好标志与变量的绑定关系后，调用flag.Parse()。
func Parse() {
    // Ignore errors; CommandLine is set for ExitOnError.
    CommandLine.Parse(os.Args[1:])
}


// Parse parses flag definitions from the argument list, which should not
// include the command name.  Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help was set but not defined.
// FlagSet类型提供的实现方法
func (f *FlagSet) Parse(arguments []string) error {
    f.parsed = true
    f.args = arguments
    for {
        seen, err := f.parseOne()
        if seen {
            continue
        }
        if err == nil {
            break
        }
        switch f.errorHandling {
        case ContinueOnError:
            return err
        case ExitOnError:
            os.Exit(2)
        case PanicOnError:
            panic(err)
        }
    }
    return nil
}


// parseOne parses one flag. It reports whether a flag was seen.
// 解析每个标志并返回相关结果，若碰到 '-' 或 '--' 时也会直接终止整个标志解析过程，每解析成功一个标志就会将该标志信息放入FlagSet.actual映射中
func (f *FlagSet) parseOne() (bool, error) {
    if len(f.args) == 0 {
        return false, nil
    }
    s := f.args[0]
    if len(s) == 0 || s[0] != '-' || len(s) == 1 {
        return false, nil
    }
    num_minuses := 1
    if s[1] == '-' {
        num_minuses++
        if len(s) == 2 { // "--" terminates the flags
            f.args = f.args[1:]
            return false, nil
        }
    }
    name := s[num_minuses:]
    if len(name) == 0 || name[0] == '-' || name[0] == '=' {
        return false, f.failf("bad flag syntax: %s", s)
    }


    // it's a flag. does it have an argument?
    f.args = f.args[1:]
    has_value := false
    value := ""
    for i := 1; i < len(name); i++ { // equals cannot be first
        if name[i] == '=' {
            value = name[i+1:]
            has_value = true
            name = name[0:i]
            break
        }
    }
    m := f.formal
    flag, alreadythere := m[name] // BUG
    if !alreadythere {
        if name == "help" || name == "h" { // special case for nice help message.
            f.usage()
            return false, ErrHelp
        }
        return false, f.failf("flag provided but not defined: -%s", name)
    }
    if fv, ok := flag.Value.(boolFlag); ok && fv.IsBoolFlag() { // special case: doesn't need an arg
        if has_value {
            if err := fv.Set(value); err != nil {
                return false, f.failf("invalid boolean value %q for  -%s: %v", value, name, err)
            }
        } else {
            fv.Set("true")
        }
    } else {
        // It must have a value, which might be the next argument.
        if !has_value && len(f.args) > 0 {
            // value is the next arg
            has_value = true
            value, f.args = f.args[0], f.args[1:]
        }
        if !has_value {
            return false, f.failf("flag needs an argument: -%s", name)
        }
        if err := flag.Value.Set(value); err != nil {
            return false, f.failf("invalid value %q for flag -%s: %v", value, name, err)
        }
    }
    if f.actual == nil {
        f.actual = make(map[string]*Flag)
    }
    f.actual[name] = flag
    return true, nil
}