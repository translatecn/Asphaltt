# 剑指 Offer 22. 链表中倒数第k个节点

链接：https://leetcode-cn.com/problems/lian-biao-zhong-dao-shu-di-kge-jie-dian-lcof/

## 双指针解法

1. 让一个指针先走 k 步
2. 此后两个指针同时往下走

```go
/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func getKthFromEnd(head *ListNode, k int) *ListNode {
    p, i := head, 0
    for ; p != nil && i < k; i++ {
        p = p.Next
    }

    for ; p != nil; p = p.Next {
        head = head.Next
    }
    return head
}
```

## 解法效果

![剑指 Offer 22. lian-biao-zhong-dao-shu-di-kge-jie-dian-lcof](./img/剑指 Offer 22. lian-biao-zhong-dao-shu-di-kge-jie-dian-lcof.png)

## 测试用例

```txt
[1,2,3,4,5]
2
[1,2,3,4,5]
5
[1]
1
```

