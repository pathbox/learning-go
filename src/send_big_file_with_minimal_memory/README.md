https://medium.com/@owlwalks/sending-big-file-with-minimal-memory-in-golang-8f3fc280d2c

multipart.Writer will automatically enclose the parts (files, fields) in boundary markup before sending them ❤️ We don’t need to get our hands dirty.

Above works great until you do some benchmark and see the memory allocation grows linearly with the file size, so what went wrong? It turns out that buf is fully buffered the whole file content, buf sequentially reads a modest 32kB from the file but won’t stop until it reaches EOF, so to hold the file content buf needs to be at least at the file size, plus some additional boundary markup.

HTTP/1.1 has a method of transfer data in chunks unboundedly, without specifying Content-Lengthof request body, this is an important feature we can utilize.

So buf causes problem, but how can we synchronize file content to the request without it? io.Pipe is born for these tasks, as the name suggests, it pipes writer to reader: