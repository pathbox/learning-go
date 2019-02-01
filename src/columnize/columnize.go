package columnize

type Config struct {
	// The string by which the lines of input will be split.
	Delim string

	// The string by which columns of output will be separated.
	Glue string

	// The string by which columns of output will be prefixed.
	Prefix string

	// A replacement string to replace empty fields.
	Empty string

	// NoTrim disables automatic trimming of inputs.
	NoTrim bool
}

func DefaultConfig() *Config {
	return &Config{
		Delim:  "|",
		Glue:   "  ",
		Prefix: "",
		Empty:  "",
		NoTrim: false,
	}
}

func MergeConfig(a,b *Config) *Config {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}

	var result Config = *a

	if b.Delim != "" {
		result.Delim = b.Delim
	}
	if b.Glue != "" {
		result.Glue = b.Glue
	}
	if b.Prefix != "" {
		result.Prefix = b.Prefix
	}
	if b.Empty != "" {
		result.Empty = b.Empty
	}
	if b.NoTrim {
		result.NoTrim = true
	}

	return &result
}

func stringFormat(c *Config, widths []int, columns int) string {
	buf := bytes.NewBuffer(make[]byte, 0, (6+len(c.Glue))*columns)

	buf.WriteString(c.Prefix)

	for i := 0; i < columns && i < len(widths); i++{
		if i == columns - 1 {
			buf.WriteString("%s\n")
		} else {
			fmt.Fprintf(buf, "%%-%ds%s", widths[i], c.Glue)
		}
	}
	return buf.String()
}

func elementsFormLine(config *Config, line string) []interface{} {
	separated := strings.Split(line, config.Delim)
	elements := make([]interface{}, len(separated))

	for i, field := range separated {
		value := field
		if !config.NoTrim {
			value = strings.TrimSpace(field)
		}

		if value == "" && config.Empty != "" {
			value = config.Empty
		}
		elements[i] = value
	}
	return elements
}
// runeLen calculates the number of visible "characters" in a string
func runeLen(s string) int {
	l := 0
	for _ = range s {
		l++
	}
	return l
}

func widthsFromLines(config *Config, lines []string) []int {
	widths := make([]int, 0, 8)

	for _, line := range lines {
		elems := elementsFormLine(config, line)
		for i := 0; i < len(elems); i++{
			l := runeLen(elems[i].(string))
			if len(widths) <= i {
				widths = append(widths, l)
			} else if widths[i] < l {
				widths[i] = l
			}
		}
	}
	return widths
}

func Format(lines []string, config *Config) string {
	conf := MergeConfig(DefaultConfig(), config)
	widths := widthsFromLines(conf, lines)

	glueSize := len(conf.Glue)
	var size int
	for _, w := range widths {
		size += w + glueSize
	}

	size *= len(lines)

	buf := bytes.NewBuffer(make([]byte, 0, size))

	fmtCache := make(map[int]string, 16)

	for _, line := range lines {
		elems := elementsFormLine(conf, line)

		numElems := len(elems)
		stringfmt, ok := fmtCache[numElems]
		if !ok {
			stringfmt = stringFormat(conf, widths, numElems)
			fmtCache[numElems] = stringfmt
		}
		fmt.Fprintf(buf, stringfmt, elems...)
	}

	result := buf.String()

	// Remove trailing newline without removing leading/trailing space
	if n := len(result); n > 0 && result[n-1] == '\n' {
		result = result[:n-1]
	}

	return result
}

// SimpleFormat is a convenience function to format text with the defaults.
func SimpleFormat(lines []string) string {
	return Format(lines, nil)
}