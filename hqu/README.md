# Hight QUality, for learning data structure and algorithm

> data structure: queue, stack
>
> algorithm: quick sort, merge sort, fibonacci number, topn, topk, max income

The source code of them is simple, please do the reading.

## Queue

Base on bucket. Queue has two array of buckets. One is used to store data. The another is used to reuse buckets as a freelist.

> Bucket is a small array to keep some elements.

## Stack

Like queue, stack also has two array of buckets.

## Quick sort

A standard quick sort algorithm implementation with Go.

## Merge sort

A standard merge sort algorithm implementation with Go.

## Fibonacci number

I implement fibonacci number with three ways, iteration, recursion and polynomial.

1. Iteration is a `for` loop.
2. Recursion is a recursive function.
3. Polynomial is a polynomial formula.

## TopN and TopK

Refering to divide and conquer method, resolve it with reduce and conquer method. Like quick sort, get the index in partition, and do partition again with low~index or index~high based on whether the index is greater or less than the n/k.

## Max income

Resolving it with dynamic programming is complicate. Refering to **Introduction to Algorithm**, find the right index for the max income, and then find the left index. So simple it is.
