package main

import (
    "fmt"
    "strconv"
)

func main() {
    var name string
    fmt.Println("Please enter your name: ")
    fmt.Scanln(&name)

    var subjectNo int
    for {
        fmt.Println("Please enter the number of subjects: ")
        _, err := fmt.Scanln(&subjectNo)
        if err == nil && subjectNo > 0 {
            break
        }
        fmt.Println("Invalid input. Please enter a positive integer for the number of subjects.")
    }

    subjectHolder := make(map[string]float64)
    originalSubjectNo := subjectNo // Store the original number of subjects

    for subjectNo != 0 {
        var subjectName string
        for {
            fmt.Println("Please enter the subject name: ")
            fmt.Scanln(&subjectName)
            if subjectName != "" && !isNumeric(subjectName) {
                break
            }
            fmt.Println("Subject name cannot be empty or a number. Please enter a valid subject name.")
        }

        var subjectMark float64
        for {
            fmt.Println("Please enter the subject mark: ")
            _, err := fmt.Scanln(&subjectMark)
            if err == nil && subjectMark >= 0 && subjectMark <= 100 {
                break
            }
            fmt.Println("Invalid mark. Please enter a mark between 0 and 100.")
        }

        subjectNo = subjectNo - 1
        subjectHolder[subjectName] = subjectMark
    }

    total := 0.0
    for _, value := range subjectHolder {
        total += value
    }

    fmt.Println("Name: ", name)
    for key, value := range subjectHolder {
        fmt.Println(key, " : ", value)
    }
    fmt.Println("Average Grade: ", total/float64(originalSubjectNo))
}


func isNumeric(s string) bool {
    _, err := strconv.ParseFloat(s, 64)
    return err == nil
}