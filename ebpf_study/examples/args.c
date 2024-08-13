//go:build ignore

// #define PT_REGS_PARM1(x) ((x)->di)
// #define PT_REGS_PARM2(x) ((x)->si)
// #define PT_REGS_PARM3(x) ((x)->dx)
// #define PT_REGS_PARM4(x) ((x)->cx)
// #define PT_REGS_PARM5(x) ((x)->r8)
// #define PT_REGS_PARM6(x) ((x)->r9)

//内核函数  __get_user_pages

// constexpr int MAX_CALLING_CONV_REGS = 6;
//const char *calling_conv_regs_x86[] = {
//  "di", "si", "dx", "cx", "r8", "r9"
//};
//
//bool BTypeVisitor::VisitFunctionDecl(FunctionDecl *D) {
//    if (D->param_size() > MAX_CALLING_CONV_REGS + 1) {
//      error(GET_BEGINLOC(D->getParamDecl(MAX_CALLING_CONV_REGS + 1)),
//            "too many arguments, bcc only supports in-register parameters");
//      return false;
//    }
//}

//        BCC 代码中明确表明：只支持寄存器参数。
//        那什么是寄存器参数呢？其实就是内核函数调用约定中的前 6 个参数要通过寄存器传递，只支持这前六个寄存器参数。


