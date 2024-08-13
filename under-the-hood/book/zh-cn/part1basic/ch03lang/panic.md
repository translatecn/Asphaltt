---
weight: 1305
title: "3.5 恐慌与恢复内建函数"
---

# 3.5 恐慌与恢复内建函数

<!-- panic/recover 由 russ cox 完成 -->


<!-- 在 `buildssa` 阶段完成后，`genssa` 将初始化好的 `ssaGenBlock`，
用来最终生成 `defer` 的函数尾声：

```go
// src/cmd/compile/internal/gc/ssa.go
func genssa(f *ssa.Func, pp *Progs) {
	var s SSAGenState
	...
	for i, b := range f.Blocks {
		...
		thearch.SSAGenBlock(&s, b, next) // 在内部调用 ssaGenBlock
		...
	}
	...
}
func ssaGenBlock(s *gc.SSAGenState, b, next *ssa.Block) {
	switch b.Kind {
	case ssa.BlockDefer:
		p := s.Prog(x86.ATESTL) // TESTL
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x86.REG_AX // AX
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x86.REG_AX   // AX
		p = s.Prog(x86.AJNE)    // JNE
		p.To.Type = obj.TYPE_BRANCH
		s.Branches = append(s.Branches, gc.Branch{P: p, B: b.Succs[1].Block()})
		if b.Succs[0].Block() != next {
			p := s.Prog(obj.AJMP) // JMP
			p.To.Type = obj.TYPE_BRANCH
			s.Branches = append(s.Branches, gc.Branch{P: p, B: b.Succs[0].Block()})
		}
	case ...:
	}
}
``` -->

<!-- 

		// 根据 TESTL AX, AX 判断 deferproc 或者 deferprocStack 的返回值：
		// 0 表示我们应该继续执行
		// 1 表示我们应该跳转到 deferreturn 调用

原因在于 `deferproc` 或者 `deferprocStack` 的返回值可能是 0 也可能是 1。
当返回 0 时，说明 defer 语句的记录被成功创建，返回值会记录在 AX 上：

```go
//go:nosplit
func deferproc(siz int32, fn *funcval) {
	...
	return0()
}
```

```asm
TEXT runtime·return0(SB), NOSPLIT, $0
	MOVL	$0, AX
	RET
```

而当发生 panic 时，AX 的值会记录为 1，从而会跳转到函数尾声，执行 `deferreturn`。 -->



<!-- TODO: 当发生 panic 时，如何回到一个开放编码的 defer
对于一个开放编码式 defer 而言，它比普通类型的 defer 多记录了以下字段：

```go
// src/runtime/panic.go
type _defer struct {
	started   bool	// 指定 defer 是否已经执行
	openDefer bool	// 指定是否是一个开放编码式 defer
	fd        unsafe.Pointer
	varp      uintptr
	framepc   uintptr
	...
}
```

这些字段的主要功能是在需要进行调用时候 -->

<!-- ### 函数信息

TODO: 转移到 panic 中，函数信息仅用于 panic 内建发生时

```go
// src/cmd/compile/internal/gc/ssa.go
func buildssa(fn *Node, worker int) *ssa.Func {
	...
	if s.hasOpenDefers {
		s.emitOpenDeferInfo()
	}
	...
}
```

```go
func (s *state) emitOpenDeferInfo() {
	x := Ctxt.Lookup(s.curfn.Func.lsym.Name + ".opendefer")
	s.curfn.Func.lsym.Func.OpenCodedDeferInfo = x
	off := 0

	// Compute maxargsize (max size of arguments for all defers)
	// first, so we can output it first to the funcdata
	var maxargsize int64
	for i := len(s.openDefers) - 1; i >= 0; i-- {
		r := s.openDefers[i]
		argsize := r.n.Left.Type.ArgWidth()
		if argsize > maxargsize {
			maxargsize = argsize
		}
	}
	off = dvarint(x, off, maxargsize)
	off = dvarint(x, off, -s.deferBitsTemp.Xoffset)
	off = dvarint(x, off, int64(len(s.openDefers)))

	// Write in reverse-order, for ease of running in that order at runtime
	for i := len(s.openDefers) - 1; i >= 0; i-- {
		r := s.openDefers[i]
		off = dvarint(x, off, r.n.Left.Type.ArgWidth())
		off = dvarint(x, off, -r.closureNode.Xoffset)
		numArgs := len(r.argNodes)
		if r.rcvrNode != nil {
			// If there's an interface receiver, treat/place it as the first
			// arg. (If there is a method receiver, it's already included as
			// first arg in r.argNodes.)
			numArgs++
		}
		off = dvarint(x, off, int64(numArgs))
		if r.rcvrNode != nil {
			off = dvarint(x, off, -r.rcvrNode.Xoffset)
			off = dvarint(x, off, s.config.PtrSize)
			off = dvarint(x, off, 0)
		}
		for j, arg := range r.argNodes {
			f := getParam(r.n, j)
			off = dvarint(x, off, -arg.Xoffset)
			off = dvarint(x, off, f.Type.Size())
			off = dvarint(x, off, f.Offset)
		}
	}
}
``` -->


panic 能中断一个程序的执行，同时也能在一定情况下进行恢复。本节我们就来看一看 panic 和 recover 这对关键字
的实现机制。根据我们对 Go 的实践，可以预见的是，他们的实现跟调度器和 defer 关键字也紧密相关。

最好的方式当然是了解编译器究竟做了什么事情：

```go
package main

func main() {
	defer func() {
		recover()
	}()
	panic(nil)
}
```

其汇编的形式为：

```asm
TEXT main.main(SB) /Users/changkun/dev/go-under-the-hood/demo/7-lang/panic/main.go
  (...)
  main.go:7		0x104e05b		0f57c0			XORPS X0, X0				
  main.go:7		0x104e05e		0f110424		MOVUPS X0, 0(SP)			
  main.go:7		0x104e062		e8d935fdff		CALL runtime.gopanic(SB)		
  main.go:7		0x104e067		0f0b			UD2					
  (...)
```

可以看到 panic 这个关键词本质上只是一个 `runtime.gopanic` 调用。

而与之对应的 `recover` 则：

```asm
TEXT main.main.func1(SB) /Users/changkun/dev/go-under-the-hood/demo/7-lang/panic/main.go
  (...)
  main.go:5		0x104e09d		488d442428		LEAQ 0x28(SP), AX			
  main.go:5		0x104e0a2		48890424		MOVQ AX, 0(SP)				
  main.go:5		0x104e0a6		e8153bfdff		CALL runtime.gorecover(SB)		
  (...)
```

其实也只是一个 `runtime.gorecover` 调用。

## 9.3.1 `gopanic` 和 `gorecover`

正如前面所探究得来，panic 关键字不过是一个 `gopanic` 调用，接受一个参数。
在处理 panic 期间，会先判断当前 panic 的类型，确定 panic 是否可恢复。

```go
// 预先声明的函数 panic 的实现
func gopanic(e interface{}) {
	gp := getg()

	// 判断在系统栈上还是在用户栈上
	// 如果执行在系统或信号栈时，getg() 会返回当前 m 的 g0 或 gsignal
	// 因此可以通过 gp.m.curg == gp 来判断所在栈
	// 系统栈上的 panic 无法恢复
	if gp.m.curg != gp {
		print("panic: ") // 打印
		printany(e)      // 打印
		print("\n")      // 继续打印，下同
		throw("panic on system stack")
	}

	// 如果正在进行 malloc 时发生 panic 也无法恢复
	if gp.m.mallocing != 0 {
		print("panic: ")
		printany(e)
		print("\n")
		throw("panic during malloc")
	}

	// 在禁止抢占时发生 panic 也无法恢复
	if gp.m.preemptoff != "" {
		print("panic: ")
		printany(e)
		print("\n")
		print("preempt off reason: ")
		print(gp.m.preemptoff)
		print("\n")
		throw("panic during preemptoff")
	}

	// 在 g 锁在 m 上时发生 panic 也无法恢复
	if gp.m.locks != 0 {
		print("panic: ")
		printany(e)
		print("\n")
		throw("panic holding locks")
	}
	...
}
```

其他情况，panic 可以从运行时进行恢复，这时候会创建一个 `_panic` 实例。`_panic` 类型
定义了一个 `_panic` 链表：

```go
// _panic 保存了一个活跃的 panic
//
// 这个标记了 go:notinheap 因为 _panic 的值必须位于栈上
//
// argp 和 link 字段为栈指针，但在栈增长时不需要特殊处理：因为他们是指针类型且
// _panic 值只位于栈上，正常的栈指针调整会处理他们。
//
//go:notinheap
type _panic struct {
	argp      unsafe.Pointer // panic 期间 defer 调用参数的指针; 无法移动 - liblink 已知
	arg       interface{}    // panic 的参数
	link      *_panic        // link 链接到更早的 panic
	recovered bool           // 表明 panic 是否结束
	aborted   bool           // 表明 panic 是否忽略
}
```

在创建过程中，panic 保存了对应的消息，并指向了保存在 goroutine 链表中先前的 panic
链表：

```go
	var p _panic
	p.arg = e
	p.link = gp._panic
	gp._panic = (*_panic)(noescape(unsafe.Pointer(&p)))

	atomic.Xadd(&runningPanicDefers, 1)
```

接下来开始逐一调用当前 goroutine 的 defer 方法，
检查用户态代码是否需要对 panic 进行恢复：

```go
	for {
		// 开始逐个取当前 goroutine 的 defer 调用
		d := gp._defer
		// 如果没有 defer 调用，则跳出循环
		if d == nil {
			break
		}

		// 如果 defer 是由早期的 panic 或 Goexit 开始的（并且，因为我们回到这里，这引发了新的 panic），
		// 则将 defer 带离链表。更早的 panic 或 Goexit 将无法继续运行。
		if d.started {
			if d._panic != nil {
				d._panic.aborted = true
			}
			d._panic = nil
			d.fn = nil
			gp._defer = d.link
			freedefer(d)
			continue
		}

		// 如果栈增长或者垃圾回收在 reflectcall 开始执行 d.fn 前发生
		// 标记 defer 已经开始执行，但仍将其保存在列表中，从而 traceback 可以找到并更新这个 defer 的参数帧
		d.started = true

		// 记录正在运行 defer 的 panic。如果在 defer 调用期间出现新的 panic，该 panic 将在列表中
		// 找到 d 并标记 d._panic（该 panic）中止。
		d._panic = (*_panic)(noescape(unsafe.Pointer(&p)))

		p.argp = unsafe.Pointer(getargp(0))
		reflectcall(nil, unsafe.Pointer(d.fn), deferArgs(d), uint32(d.siz), uint32(d.siz))
		p.argp = nil

		// reflectcall 不会 panic. 移出 d.
		if gp._defer != d {
			throw("bad defer entry in panic")
		}
		d._panic = nil
		d.fn = nil
		gp._defer = d.link

		pc := d.pc
		sp := unsafe.Pointer(d.sp) // 必须是指针，以便在栈复制期间进行调整
		freedefer(d)
		if p.recovered {
			atomic.Xadd(&runningPanicDefers, -1)

			gp._panic = p.link
			// 忽略的 panic 会被标记，但仍然保留在 g.panic 列表中
			// 这里将它们移出列表
			for gp._panic != nil && gp._panic.aborted {
				gp._panic = gp._panic.link
			}
			if gp._panic == nil { // 必须由 signal 完成
				gp.sig = 0
			}
			// 传递关于恢复帧的信息
			gp.sigcode0 = uintptr(sp)
			gp.sigcode1 = pc
			// 调用 recover，并重新进入调度循环，不再返回
			mcall(recovery)
			// 如果无法重新进入调度循环，则无法恢复错误
			throw("recovery failed")
		}
	}
```

这个循环说明了很多问题。首先，当 panic 发生时，如果错误是可恢复的错误，那么
会逐一遍历该 goroutine 对应 defer 链表中的 defer 函数链表，直到 defer 遍历完毕、
或者再次进入调度循环（recover 的 mcall 调用）
后才会停止。

defer 并非简单的遍历，每个在 panic 和 recover 之间的 defer 都会在这里通过 `reflectcall` 执行。

```go
// reflectcall 使用 arg 指向的 n 个参数字节的副本调用 fn。
// fn 返回后，reflectcall 在返回之前将 n-retoffset 结果字节复制回 arg+retoffset。
// 如果重新复制结果字节，则调用者应将参数帧类型作为 argtype 传递，以便该调用可以在复制期间执行适当的写障碍。
// reflect 包传递帧类型。在 runtime 包中，只有一个调用将结果复制回来，即 cgocallbackg1，
// 并且它不传递帧类型，这意味着没有调用写障碍。参见该调用的页面了解相关理由。
//
// 包 reflect 通过 linkname 访问此符号
func reflectcall(argtype *_type, fn, arg unsafe.Pointer, argsize uint32, retoffset uint32)
```

如果某个包含了 recover 的调用（即 gorecover 调用）被执行，这时 `_panic` 实例 `p.recovered` 会被标记为 `true`：

```go
// 执行预先声明的函数 recover。
// 不允许分段栈，因为它需要可靠地找到其调用者的栈段。
//
// TODO(rsc): Once we commit to CopyStackAlways,
// this doesn't need to be nosplit.
//go:nosplit
func gorecover(argp uintptr) interface{} {
	// 必须在 panic 期间作为 defer 调用的一部分在函数中运行。
	// 必须从调用的最顶层函数（ defer 语句中使用的函数）调用。
	// p.argp 是最顶层 defer 函数调用的参数指针。
	// 比较调用方报告的 argp，如果匹配，则调用者可以恢复。
	gp := getg()
	p := gp._panic
	if p != nil && !p.recovered && argp == uintptr(p.argp) {
		p.recovered = true
		return p.arg
	}
	return nil
}
```

同时 `recover()` 这个函数还会返回 panic 的保存相关信息 `p.arg`。
恢复的原则取决于 `gorecover` 这个方法调用方报告的 argp 是否与 `p.argp` 相同，仅当相同才可恢复。

当 `reflectcall` 执行完毕后，这时如果一个 panic 是可恢复的，`p.recovered` 已经被标记为 `true`，
从而会通过 `mcall` 的方式来执行 `recovery` 函数来重新进入调度循环：

```go
// 在发生 panic 后 defer 函数调用 recover 后展开栈。然后安排继续运行，
// 就像 defer 函数的调用方正常返回一样。
func recovery(gp *g) {
	// 传递到 G 结构的 defer 信息
	sp := gp.sigcode0
	pc := gp.sigcode1

	...

	// 使 deferproc 为此 d 返回
	// 这时候返回 1。调用函数将跳转到标准的返回尾声
	gp.sched.sp = sp
	gp.sched.pc = pc
	gp.sched.lr = 0
	gp.sched.ret = 1
	gogo(&gp.sched)
}
```

当然如果所有的 defer 都没有指明显式的 recover，那么这时候则直接在运行时抛出 panic 信息：

```go
	// 消耗完所有的 defer 调用，保守地进行 panic
	// 因为在冻结之后调用任意用户代码是不安全的，所以我们调用 preprintpanics 来调用
	// 所有必要的 Error 和 String 方法来在 startpanic 之前准备 panic 字符串。
	preprintpanics(gp._panic)

	fatalpanic(gp._panic) // 不应该返回
	*(*int)(nil) = 0      // 无法触及
}
```

从而完成 `gopanic` 的调用。

至于 `preprintpanics` 和 `fatalpanic` 无非是一些错误输出，不再赘述：

```go
// 在停止前调用所有的 Error 和 String 方法
func preprintpanics(p *_panic) {
	defer func() {
		if recover() != nil {
			throw("panic while printing panic value")
		}
	}()
	for p != nil {
		switch v := p.arg.(type) {
		case error:
			p.arg = v.Error()
		case stringer:
			p.arg = v.String()
		}
		p = p.link
	}
}
// fatalpanic 实现了不可恢复的 panic。类似于 fatalthrow，
// 要求如果 msgs != nil，则 fatalpanic 仍然能够打印 panic 的消息并在 main 在退出时候减少 runningPanicDefers。
//
//go:nosplit
func fatalpanic(msgs *_panic) {
	pc := getcallerpc()
	sp := getcallersp()
	gp := getg()
	var docrash bool
	// 切换到系统栈来避免栈增长，如果运行时状态较差则可能导致更糟糕的事情
	systemstack(func() {
		if startpanic_m() && msgs != nil {
			// 有 panic 消息和 startpanic_m 则可以尝试打印它们

			// startpanic_m 设置 panic 会从阻止 main 的退出，
			// 因此现在可以开始减少 runningPanicDefers 了
			atomic.Xadd(&runningPanicDefers, -1)

			printpanics(msgs)
		}

		docrash = dopanic_m(gp, pc, sp)
	})

	if docrash {=
		// 通过在上述 systemstack 调用之外崩溃，调试器在生成回溯时不会混淆。
		// 函数崩溃标记为 nosplit 以避免堆栈增长。
		crash()
	}

	// 从系统栈退出
	systemstack(func() {
		exit(2)
	})

	*(*int)(nil) = 0 // 不可达
}
```

## 小结

从 panic 和 recover 这对关键字的实现上可以看出，可恢复的 panic 必须要 recover 的配合。
而且，这个 recover 必须位于同一 goroutine 的直接调用链上，否则无法对 panic 进行恢复。

例如，如果 A 依次调用了 B 和 C，而 B 包含了 recover，而 C 发生了 panic，则这时 B 的 recover 无法恢复 C 的 panic；

```go
func A () {
    B()
    C()
}

func B() {
    defer func () {
        recover() // 无法恢复 panic("C")
    }()

    println("B")
}

func C() {
    panic("C")
}
```

又例如 A 调用了 B 而 B 又调用了 C，那么 C 发生 panic 时，如果 A 要求了 recover 则仍然可以恢复。

```go
func A () {
    defer func () {
        recover() // 可以恢复 panic("C")
    }()

    B()
}

func B() {
    C()
}

func C() {
    panic("C")
}
```

当一个 panic 被恢复后，调度并因此中断，会重新进入调度循环，进而继续执行 recover 后面的代码，
包括比 recover 更早的 defer（因为已经执行过得 defer 已经被释放，
而尚未执行的 defer 仍在 goroutine 的 defer 链表中），或者 recover 所在函数的调用方。

## 许可

[Go under the hood](https://github.com/golang-design/under-the-hood) | CC-BY-NC-ND 4.0 & MIT &copy; [changkun](https://changkun.de)