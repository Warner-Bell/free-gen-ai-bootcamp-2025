Required changes for handler files:

1. In internal/handlers/words/handler.go:
```go
package words

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"
    "time"  // Add this import
)
```

2. In internal/handlers/groups/handler.go:
- Change all instances of `Word` to `models.Word`
- Ensure models package is imported:
```go
import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"
    "backend_go/internal/models"
)
```

The three locations that need Word -> models.Word changes are:
1. Line 102: function return type
2. Line 116: slice declaration
3. Line 118: variable declaration

These changes will resolve the build errors:
- "undefined: time" in words/handler.go
- "undefined: Word" in groups/handler.go