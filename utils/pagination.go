package utils

import (
    "strconv"

    "github.com/gin-gonic/gin"
)

func GetPaginationParams(c *gin.Context) (limit, offset int) {
    limitStr := c.DefaultQuery("limit", "10")
    offsetStr := c.DefaultQuery("offset", "0")

    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 {
        limit = 10
    }

    offset, err = strconv.Atoi(offsetStr)
    if err != nil || offset < 0 {
        offset = 0
    }

    return limit, offset
}
