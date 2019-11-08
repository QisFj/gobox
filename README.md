# gobox
A humble memory database. 

简陋的内存数据库。

Use this in your production is not a good idea.  

不建议用在生产环境。

## Usage

```go
package main

import (
    "github.com/qisfj/gobox"
)

type T struct {
    ID int // required
    // ...
}

func main() {
	box := gobox.New()
    t := T{}
    _ = box.Set(&t)
    _ = box.Get(&t)
}
```