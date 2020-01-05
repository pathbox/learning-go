http://vinllen.com/golangshi-xian-mapreducedan-jin-cheng-ban-ben/

归纳来说主要分为5部分：用户程序、Master、Mapper、Reducer、Combiner（上图未给出）。

用户程序。用户程序主要对输入数据进行分割，制定Mapper、Reducer、Combiner的代码。
Master：中控系统。控制分发Mapper、Reduer的个数，比如生成m个进程处理Mapper，n个进程处理Reducer。其实对Master来说，Mapper和Reduer都属于worker，只不过跑的程序不一样，Mapper跑用户输入的map代码，Reduer跑用户输入的reduce代码。Master还作为管道负责中间路径传递，比如将Mapper生成的中间文件传递给Reduer，将Reduer生成的结果文件返回，或者传递给Combiner（如果有需要的话）。由于Master是单点，性能瓶颈，所以可以做集群：主备模式或者分布式模式。可以用zookeeper进行选主，用一些消息中间件进行数据同步。Master还可以进行一些策略处理：比如某个Worker执行时间特别长，很有可能卡住了，对分配给该Worker的数据重新分配给别的Worker执行，当然需要对多份数据返回去重处理。
Mapper：负责将输入数据切成key-value格式。Mapper处理完后，将中间文件的路径告知Master，Master获悉后传递给Reduer进行后续处理。如果Mapper未处理完，或者已经处理完但是Reduer未读完其中间输出文件，分配给该Mapper的输入将重新被别的Mapper执行。
Reducer: 接受Master发送的Mapper输出文件的消息，RPC读取文件并处理，并输出结果文件。n个Reduer将产生n个输出文件。
Combiner: 做最后的归并处理，通常不需要。
  总的来说，架构不复杂。组件间通信用啥都可以，比如RPC、HTTP或者私有协议等。