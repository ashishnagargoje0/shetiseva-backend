package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/ulule/limiter/v3"
    ginmiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
    memory "github.com/ulule/limiter/v3/drivers/store/memory"
    "log"
)

func RateLimitMiddleware() gin.HandlerFunc {
    // Limit: 5 requests per second
    rate, err := limiter.NewRateFromFormatted("5-S")
    if err != nil {
        log.Fatal(err)
    }

    store := memory.NewStore()
    middleware := ginmiddleware.NewMiddleware(limiter.New(store, rate))
    return middleware
}
