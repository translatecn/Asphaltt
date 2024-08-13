---
weight: 1202
title: "2.2 Go 程序编译流程"
---

# 2.2 Go 程序编译流程



## 2.2.1 第一阶段：词法和语法分析

- `cmd/compile/internal/syntax`（词法分析器，解析器，语法树）

在编译的第一阶段，源代码被 token 化（词法分析），解析（语法分析），并为每个源构造语法树文件。每个语法树都是相应源文件的精确表示对应于源的各种元素的节点，如表达式，声明和陈述。语法树还包括位置信息用于错误报告和调试信息的创建。

```
main -> gc.Main -> amd64.Init -> amd64.LinkArch.Init
-> typecheck -> typecheck -> saveerrors -> typecheckslice
-> checkreturn -> checkMapKeys -> capturevars -> 
typecheckinl -> inlcalls -> escapes -> 
newNowritebarrierrecChecker -> transformclosure
```

## 2.2.2 第二阶段：语义分析

- `cmd/compile/internal/gc`（类型检查，AST变换）

对 AST 进行类型检查。第一步是名称解析和类型推断，它们确定哪个对象属于哪个标识符，以及每个表达式具有的类型。类型检查包括某些额外的检查，例如 “声明和未使用” 以及确定函数是否终止。

在 AST 上也进行了某些转换。一些节点基于类型信息被细化，例如从算术加法节点类型分割的字符串添加。其他一些例子是死代码消除，函数调用内联和转义分析。

语义分析的过程中包含几个重要的操作：逃逸分析、变量捕获、函数内联、闭包处理。

## 2.2.3 第三阶段：SSA 生成

- `cmd/compile/internal/gc` (转换为SSA)
- `cmd/compile/internal/ssa` (SSA 传递与规则)

在此阶段，AST将转换为静态单一分配（SSA）形式，这是一种具有特定属性的低级中间表示，可以更轻松地实现优化并最终从中生成机器代码。

在此转换期间，将应用函数内在函数。 这些是特殊功能，编译器已经教导它们根据具体情况用大量优化的代码替换。

在AST到SSA转换期间，某些节点也被降级为更简单的组件，因此编译器的其余部分可以使用它们。 例如，内置复制替换为内存移动，并且范围循环被重写为for循环。 其中一些目前发生在转换为SSA之前，由于历史原因，但长期计划是将所有这些都移到这里。

然后，应用一系列与机器无关的传递和规则。 这些不涉及任何单个计算机体系结构，因此可以在所有 `GOARCH` 变体上运行。

这些通用过程的一些示例包括消除死代码，删除不需要的零检查以及删除未使用的分支。通用重写规则主要涉及表达式，例如用常量值替换某些表达式，以及优化乘法和浮点运算。

```
initssaconfig -> peekitabs -> funccompile ->
finit -> compileFunctions -> compileSSA -> buildssa -> genssa ->
-> typecheck -> checkMapKeys -> dumpdata -> dumpobj
```

## 2.2.4 第四阶段：机器码生成

- `cmd/compile/internal/ssa` (底层SSA和架构特定的传递)
- `cmd/internal/obj` (生成机器码)

编译器的机器相关阶段以“底层”传递开始，该传递将通用值重写为其机器特定的变体。例如，在 amd64 存储器操作数上是可能的，因此可以组合许多加载存储操作。

请注意，较低的通道运行所有特定于机器的重写规则，因此它当前也应用了大量优化。

一旦SSA“降低”并且更加特定于目标体系结构，就会运行最终的代码优化过程。这包括另一个死代码消除传递，移动值更接近它们的使用，删除从未读取的局部变量，以及寄存器分配。

作为此步骤的一部分完成的其他重要工作包括堆栈框架布局，它将堆栈偏移分配给局部变量，以及指针活动分析，它计算每个 GC 安全点上的堆栈指针。

在SSA生成阶段结束时，Go 函数已转换为一系列 `obj.Prog` 指令。它们会被传递给装载器（`cmd/internal/obj`），将它们转换为机器代码并写出最终的目标文件。目标文件还将包含反射数据，导出数据和调试信息。

## 许可

[Go under the hood](https://github.com/golang-design/under-the-hood) | CC-BY-NC-ND 4.0 & MIT &copy; [changkun](https://changkun.de)