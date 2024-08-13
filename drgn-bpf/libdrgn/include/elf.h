// Copyright (c) Meta Platforms, Inc. and affiliates.
// SPDX-License-Identifier: LGPL-2.1-or-later

// Wrapper around elf.h adding definitions that are not available in older
// versions of glibc.

#ifndef DRGN_ELF_H
#define DRGN_ELF_H

#include_next <elf.h>

// Generated by scripts/gen_elf_compat.py.
#ifndef NT_FILE
#define NT_FILE 0x46494c45
#endif
#ifndef EM_RISCV
#define EM_RISCV 243
#endif
#ifndef R_RISCV_NONE
#define R_RISCV_NONE 0
#endif
#ifndef R_RISCV_32
#define R_RISCV_32 1
#endif
#ifndef R_RISCV_64
#define R_RISCV_64 2
#endif
#ifndef R_RISCV_RELATIVE
#define R_RISCV_RELATIVE 3
#endif
#ifndef R_RISCV_COPY
#define R_RISCV_COPY 4
#endif
#ifndef R_RISCV_JUMP_SLOT
#define R_RISCV_JUMP_SLOT 5
#endif
#ifndef R_RISCV_TLS_DTPMOD32
#define R_RISCV_TLS_DTPMOD32 6
#endif
#ifndef R_RISCV_TLS_DTPMOD64
#define R_RISCV_TLS_DTPMOD64 7
#endif
#ifndef R_RISCV_TLS_DTPREL32
#define R_RISCV_TLS_DTPREL32 8
#endif
#ifndef R_RISCV_TLS_DTPREL64
#define R_RISCV_TLS_DTPREL64 9
#endif
#ifndef R_RISCV_TLS_TPREL32
#define R_RISCV_TLS_TPREL32 10
#endif
#ifndef R_RISCV_TLS_TPREL64
#define R_RISCV_TLS_TPREL64 11
#endif
#ifndef R_RISCV_BRANCH
#define R_RISCV_BRANCH 16
#endif
#ifndef R_RISCV_JAL
#define R_RISCV_JAL 17
#endif
#ifndef R_RISCV_CALL
#define R_RISCV_CALL 18
#endif
#ifndef R_RISCV_CALL_PLT
#define R_RISCV_CALL_PLT 19
#endif
#ifndef R_RISCV_GOT_HI20
#define R_RISCV_GOT_HI20 20
#endif
#ifndef R_RISCV_TLS_GOT_HI20
#define R_RISCV_TLS_GOT_HI20 21
#endif
#ifndef R_RISCV_TLS_GD_HI20
#define R_RISCV_TLS_GD_HI20 22
#endif
#ifndef R_RISCV_PCREL_HI20
#define R_RISCV_PCREL_HI20 23
#endif
#ifndef R_RISCV_PCREL_LO12_I
#define R_RISCV_PCREL_LO12_I 24
#endif
#ifndef R_RISCV_PCREL_LO12_S
#define R_RISCV_PCREL_LO12_S 25
#endif
#ifndef R_RISCV_HI20
#define R_RISCV_HI20 26
#endif
#ifndef R_RISCV_LO12_I
#define R_RISCV_LO12_I 27
#endif
#ifndef R_RISCV_LO12_S
#define R_RISCV_LO12_S 28
#endif
#ifndef R_RISCV_TPREL_HI20
#define R_RISCV_TPREL_HI20 29
#endif
#ifndef R_RISCV_TPREL_LO12_I
#define R_RISCV_TPREL_LO12_I 30
#endif
#ifndef R_RISCV_TPREL_LO12_S
#define R_RISCV_TPREL_LO12_S 31
#endif
#ifndef R_RISCV_TPREL_ADD
#define R_RISCV_TPREL_ADD 32
#endif
#ifndef R_RISCV_ADD8
#define R_RISCV_ADD8 33
#endif
#ifndef R_RISCV_ADD16
#define R_RISCV_ADD16 34
#endif
#ifndef R_RISCV_ADD32
#define R_RISCV_ADD32 35
#endif
#ifndef R_RISCV_ADD64
#define R_RISCV_ADD64 36
#endif
#ifndef R_RISCV_SUB8
#define R_RISCV_SUB8 37
#endif
#ifndef R_RISCV_SUB16
#define R_RISCV_SUB16 38
#endif
#ifndef R_RISCV_SUB32
#define R_RISCV_SUB32 39
#endif
#ifndef R_RISCV_SUB64
#define R_RISCV_SUB64 40
#endif
#ifndef R_RISCV_GNU_VTINHERIT
#define R_RISCV_GNU_VTINHERIT 41
#endif
#ifndef R_RISCV_GNU_VTENTRY
#define R_RISCV_GNU_VTENTRY 42
#endif
#ifndef R_RISCV_ALIGN
#define R_RISCV_ALIGN 43
#endif
#ifndef R_RISCV_RVC_BRANCH
#define R_RISCV_RVC_BRANCH 44
#endif
#ifndef R_RISCV_RVC_JUMP
#define R_RISCV_RVC_JUMP 45
#endif
#ifndef R_RISCV_RVC_LUI
#define R_RISCV_RVC_LUI 46
#endif
#ifndef R_RISCV_GPREL_I
#define R_RISCV_GPREL_I 47
#endif
#ifndef R_RISCV_GPREL_S
#define R_RISCV_GPREL_S 48
#endif
#ifndef R_RISCV_TPREL_I
#define R_RISCV_TPREL_I 49
#endif
#ifndef R_RISCV_TPREL_S
#define R_RISCV_TPREL_S 50
#endif
#ifndef R_RISCV_RELAX
#define R_RISCV_RELAX 51
#endif
#ifndef R_RISCV_SUB6
#define R_RISCV_SUB6 52
#endif
#ifndef R_RISCV_SET6
#define R_RISCV_SET6 53
#endif
#ifndef R_RISCV_SET8
#define R_RISCV_SET8 54
#endif
#ifndef R_RISCV_SET16
#define R_RISCV_SET16 55
#endif
#ifndef R_RISCV_SET32
#define R_RISCV_SET32 56
#endif
#ifndef R_RISCV_32_PCREL
#define R_RISCV_32_PCREL 57
#endif
#ifndef NT_ARM_PAC_MASK
#define NT_ARM_PAC_MASK 0x406
#endif

#endif /* DRGN_ELF_H */