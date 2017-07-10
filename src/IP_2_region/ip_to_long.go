// ipv4的地址本来就是用32位来表示的,分成4个8位来书写, 所以ipv4和地址是可以和32位unsigned int一一对应的,转换的算法就很显然了,把32位的整型4个字节的数分别计算出来； 反之则是ip地址4个节的数乘上对应的权值(256^3, 256^2, 256^1, 256^0)加起来即可

func Iptolong(IpStr string) (int64, error) {
  bits := strings.Split(IpStr, ".")
  if len(bits) != 4 {
    return 0, errors.New("ip format error")
  }

  var sum int64
  for i, n := range bits {
    bit, _ := strconv.ParseInt(n, 10, 64)
    sum += bit << uint(24 - 8 * i)
  }

  return sum, nil
}