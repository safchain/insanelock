# InsaneLock

Golang dead lock detector. This aims to provide a way to display a full stack back
when a lock is retain for a while.

## Install

```
go get github.com/safchain/insanelock
```

## Usage

```
package main

import (
        "fmt"
        "github.com/safchain/insanelock"
)

func main() {
        var l insanelock.RWMutex
        l.Lock()
}
```

By default an `insanelock` RWMutex is just an alias of a RWMutex. In order to use the dead lock detector
the following tag as to be used during the build process.

```
go build -tags insanelock main.go
```

## License

This software is licensed under the Apache License, Version 2.0 (the
"License"); you may not use this software except in compliance with the
License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
