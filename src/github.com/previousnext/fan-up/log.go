package main

import (
    "fmt"
)

func Fatal(m string) {
    fmt.Printf("FATAL: %s\n", m)
}

func Info(m string) {
    fmt.Printf("INFO: %s\n", m)
}