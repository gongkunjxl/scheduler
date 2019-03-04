package main
import (
    "fmt"
    "strconv"
)
func main(){
    fmt.Println("help")
    var i int
    var podName string = "hadoop-slave-"
    for i=1;i<4;i++{
       fmt.Println(podName+strconv.Itoa(i))
    }
}
