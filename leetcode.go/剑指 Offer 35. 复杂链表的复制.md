# 剑指 Offer 35. 复杂链表的复制

链接：https://leetcode-cn.com/problems/fu-za-lian-biao-de-fu-zhi-lcof/

## 哨兵解法

1. 将当前节点的 Next 指向复制的当前节点
2. 更新复制的节点的 Random
3. 新增哨兵作为链表头结点，取出复制的节点，并且恢复原来的节点的 Next

```go
/**
 * Definition for a Node.
 * type Node struct {
 *     Val int
 *     Next *Node
 *     Random *Node
 * }
 */

func copyRandomList(head *Node) *Node {
    for p := head; p != nil; {
        next := p.Next
        p.Next = &Node{
            Val: p.Val,
            Next: next,
            Random: p.Random,
        }
        p = next
    }

    for p := head; p != nil; p = p.Next.Next {
        ptr := p.Next
        if ptr.Random != nil {
            ptr.Random = ptr.Random.Next
        }
    }

    root := &Node{}
    ptr := root
    for p := head; p != nil; {
        ptr.Next = p.Next
        ptr = ptr.Next
        p.Next = ptr.Next
        p = p.Next
    }

    return root.Next
}
```

## 解法效果

![剑指 Offer 35. fu-za-lian-biao-de-fu-zhi-lcof](./img/剑指 Offer 35. fu-za-lian-biao-de-fu-zhi-lcof.png)

## 测试用例

```txt
[[7,null],[13,0],[11,4],[10,2],[1,0]]
[]
[[1,1],[2,1]]
[[3,null],[3,0],[3,null]]
[[1,0]]
[[1,0],[2,0]]
```

