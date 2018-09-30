https://colobu.com/2018/09/02/generate-random-string-in-Go/#more

##### 需要翻墙吧


BenchmarkRunes                   1000000              1703 ns/op
BenchmarkBytes                   1000000              1328 ns/op
BenchmarkBytesRmndr              1000000              1012 ns/op
BenchmarkBytesMask               1000000              1214 ns/op
BenchmarkBytesMaskImpr           5000000               395 ns/op
BenchmarkBytesMaskImprSrc        5000000               303 ns/op
