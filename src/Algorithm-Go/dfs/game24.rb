
class Game

  class << self

    def dfs24(item, index)
      $r << item
      $book[index] = true

      if $r.size == 7
        rc = $r.join('')
        rs = eval rc rescue 0  # 使用ruby eval 进行执行
        if rs == 24
          puts rc+"=24"
        end
      end

      $a.each_with_index do |ai, ix|  # 循环 item的下一个结点元素,理论上是$a数组中除了自己的其他任意元素
        if $book[ix] != true  # 排除已经取出过的元素
          $mc.each do |mi|
            $r << mi  # 取一个运算符加入到数组
            dfs24(ai, ix)
            $r.pop  # 出一个数字
            $r.pop  # 出数字前面的运算符合
            $book[ix] = false # 回退到上一个结点了,设置book对应索引值为false
          end
        end
      end

    end

    def count24
      $a = [1,0,0,23]
      $mc = ['+', '-', '*', '/']

      $a.each_with_index do |item, index|
        $r = []
        $book = [false, false, false, false]
        dfs24(item, index)
      end
    end
  end
end

Game.count24


# 想必大家都玩过一个游戏，叫做“24点”：给出4个整数，要求用加减乘除4个运算使其运算结果变成24

# 例如给出4个数：1、2、3、4。我可以用以下运算得到结果24：

# 1*2*3*4 = 24；2*3*4/1 = 24

# 上面的算法是类似全排列,会取到实际上是重复的算式,因为上面是按照图的数据结构做的,
# 如果按照二叉树的数据结构做,应该就能避免重复性质的算式