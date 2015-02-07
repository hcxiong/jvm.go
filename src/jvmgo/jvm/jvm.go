package jvm

import (
    "fmt"
    . "jvmgo/any"
    "jvmgo/cmdline"
    _ "jvmgo/native"
    "jvmgo/jvm/rtda"
    rtc "jvmgo/jvm/rtda/class"
)

func Startup(cmd *cmdline.Command) {
    classPath := cmd.Options().Classpath()
    classLoader := rtc.NewClassLoader(classPath)
    mainThread := createMainThread(classLoader, cmd.Class(), cmd.Args())

    // todo
    defer func() {
        if r := recover(); r != nil {
            for !mainThread.IsStackEmpty() {
                frame := mainThread.PopFrame()
                fmt.Printf("%v %v\n", frame.Method().Class(), frame.Method())
            }

            err, ok := r.(error)
            if !ok {
                err = fmt.Errorf("%v", r)
                panic(err.Error())
            } else {
                panic(err.Error())
            }
        }
    }()

    loop(mainThread)
}

func createMainThread(classLoader Any, className string, args []string) (*rtda.Thread) {
    fakeMethod := rtc.NewStartupMethod([]byte{0xff, 0xb1}, classLoader)
    mainThread := rtda.NewThread(128, nil)
    mainFrame := mainThread.NewFrame(fakeMethod)
    mainThread.PushFrame(mainFrame)
    
    stack := mainFrame.OperandStack()
    stack.Push(args)
    stack.Push(className)
    stack.Push(classLoader)

    return mainThread
}
