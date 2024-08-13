// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"runtime/internal/atomic"
	"unsafe"
)

// 这些函数不能拥有 go:noescape 标记，因为虽然 ptr 没有逃逸，但是 new 会逃逸。
// 如果 new 被标记为非逃逸，则编译器将不正确的对该指针变量存储的值进行逃逸分析决策。

// atomicwb 在原子指针写入之前执行 write barrier，调用方应使用 "if writeBarrier.enabled" 对调用
// 进行保护
//
//go:nosplit
func atomicwb(ptr *unsafe.Pointer, new unsafe.Pointer) {
	slot := (*uintptr)(unsafe.Pointer(ptr))
	if !getg().m.p.ptr().wbBuf.putFast(*slot, uintptr(new)) {
		wbBufFlush(slot, uintptr(new))
	}
}

// atomicstorep performs *ptr = new atomically and invokes a write barrier.
//
//go:nosplit
func atomicstorep(ptr unsafe.Pointer, new unsafe.Pointer) {
	if writeBarrier.enabled {
		atomicwb((*unsafe.Pointer)(ptr), new)
	}
	atomic.StorepNoWB(noescape(ptr), new)
}

// Like above, but implement in terms of sync/atomic's uintptr operations.
// We cannot just call the runtime routines, because the race detector expects
// to be able to intercept the sync/atomic forms but not the runtime forms.

//go:linkname sync_atomic_StoreUintptr sync/atomic.StoreUintptr
func sync_atomic_StoreUintptr(ptr *uintptr, new uintptr)

//go:linkname sync_atomic_StorePointer sync/atomic.StorePointer
//go:nosplit
func sync_atomic_StorePointer(ptr *unsafe.Pointer, new unsafe.Pointer) {
	if writeBarrier.enabled {
		atomicwb(ptr, new)
	}
	sync_atomic_StoreUintptr((*uintptr)(unsafe.Pointer(ptr)), uintptr(new))
}

//go:linkname sync_atomic_SwapUintptr sync/atomic.SwapUintptr
func sync_atomic_SwapUintptr(ptr *uintptr, new uintptr) uintptr

//go:linkname sync_atomic_SwapPointer sync/atomic.SwapPointer
//go:nosplit
func sync_atomic_SwapPointer(ptr *unsafe.Pointer, new unsafe.Pointer) unsafe.Pointer {
	if writeBarrier.enabled {
		atomicwb(ptr, new)
	}
	old := unsafe.Pointer(sync_atomic_SwapUintptr((*uintptr)(noescape(unsafe.Pointer(ptr))), uintptr(new)))
	return old
}

//go:linkname sync_atomic_CompareAndSwapUintptr sync/atomic.CompareAndSwapUintptr
func sync_atomic_CompareAndSwapUintptr(ptr *uintptr, old, new uintptr) bool

//go:linkname sync_atomic_CompareAndSwapPointer sync/atomic.CompareAndSwapPointer
//go:nosplit
func sync_atomic_CompareAndSwapPointer(ptr *unsafe.Pointer, old, new unsafe.Pointer) bool {
	if writeBarrier.enabled {
		atomicwb(ptr, new)
	}
	return sync_atomic_CompareAndSwapUintptr((*uintptr)(noescape(unsafe.Pointer(ptr))), uintptr(old), uintptr(new))
}
