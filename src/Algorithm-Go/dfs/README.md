DFS 广度优先搜索 图搜索算法

```
/*
 * DFS核心伪代码
 * 前置条件是visit数组全部设置成false
 * @param n 当前开始搜索的节点
 * @param d 当前到达的深度
 * @return 是否有解
*/
```

bool DFS(Node n, int d) {
  if (isEnd(n, d)) { // 一旦搜索深度到达一个结束状态,就返回true
    return true;
  }

  for (Node nexNode in n) { // 遍历n相邻的结点nextNode
    if(!visit[nextNode]) {
      visit[nextNode] = true; // 在下一步搜索中,nextNode不能再次出现
      if(DFS(nextNode, d+1)) { // n的某一个下一个结点,以该结点开始,递归寻找该节点的某一个下一个结点,深度+1
        // 如果搜索出有解,做其他事情,例如揭露结果深度
        return true;
      }

      // 重新设置成false,因为它有可能出现在下一次搜索的别的路径中
      visit[nextNode] = false;
    }
  }

  return false; // 本地搜索无解
}


 DFS广度优先搜索 图搜索算法, 可以把具体场景条件转为图的结构

 图的结构是怎样的呢?

 结点
 线路
 线路有向无向

 结点的数量, 结点可以拥有各种的属性

 一个结点有多少个连接的线路

 线路是有方向还是无方向

 结点可以看成是某个元素, 线路表示这个元素通过某种方法或函数,到达另一个结点元素
 DFS适合此类题目：给定初始状态跟目标状态，要求判断从初始状态到目标状态是否有解或有几种解

DFS在回溯时要取消原先的标记，而BFS不存在回溯也就不存在取消标记这一问题。DFS可以用递归来写，也可以用栈来写

1. 访问v0顶点
2. 置visited[v0] = 1
3. 以v0出发,访问v0的邻结点(下一个结点),若邻结点存在w, 则DFS(w),以w结点重复第一步,查找w结点的邻结点
4. 回朔一个结点时,设置visited[x] = 0,回到x的上一个结点w,从w开始,查找除了x之外的其他邻结点



参考链接

```
http://blog.csdn.net/saltriver/article/details/54429068
```