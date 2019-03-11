### 爬楼梯问题

假设有一个10台阶的阶梯，每次你可以爬一个阶梯或两个阶梯。请问爬完这10个阶梯可以有多少种爬法？

##### 思考：

爬上第10个阶梯，相当于爬上第9有多少种爬法 + 爬上第8个阶梯有多少中爬法。然后，

爬上第9个阶梯，相当于爬上第8有多少种爬法 + 爬上第7个阶梯有多少中爬法。然后，

爬上第8个阶梯, 相当于爬上第7有多少种爬法 + 爬上第6个阶梯有多少中爬法

...


得到的一个公式就是

```
F(1) = 1;
F(2) = 2;
F(n) = F(n-1)+F(n-2)（n>=3）
```

参考链接：

https://mp.weixin.qq.com/s?__biz=MzI1MTIzMzI2MA==&mid=2650561116&idx=1&sn=a6cb8c7bf52bc94b1a5d9feae21effa2&chksm=f1feecdfc68965c9bf20c1ef9373118dcf79fdc93a29dd172796b0dbef3e9ef6830b44ef072f&mpshare=1&scene=24&srcid=0616KUfpBJBFKKcCh0Z69GtW&key=b10a7c153a57fbb9c9340e20573980776f826657ab7cdb9b076ed3ca96bcd772a26f9f91086f829947b1285734530a1356d6fd2476485bba6da9f8f4f4637fc5ac69a884093e6b5ceb391093e924d291&ascene=0&uin=OTcxOTY1NTU%3D&devicetype=iMac+MacBookPro12%2C1+OSX+OSX+10.12+build(16A323)&version=12020710&nettype=WIFI&fontScale=100&pass_ticket=mmYmEVv3gqNbe2uX0CV7S0tVNyYKDKJ9qiaR9Jf5%2Fno%3D