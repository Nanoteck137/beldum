package utils

import (
	"log"
	"math"
	"strings"

	"github.com/gosimple/slug"
	"github.com/nrednav/cuid2"
)

var CreateId = createIdGenerator(32)
var CreateSmallId = createIdGenerator(8)

var CreateProjectId = createIdGenerator(8)
var CreateBoardId = createIdGenerator(8)
var CreateTaskId = createIdGenerator(16)

var CreateApiTokenId = createIdGenerator(32)

func createIdGenerator(length int) func() string {
	res, err := cuid2.Init(cuid2.WithLength(length))
	if err != nil {
		log.Fatal("Failed to create id generator", "err", err)
	}

	return res
}

func ParseAuthHeader(authHeader string) string {
	splits := strings.Split(authHeader, " ")
	if len(splits) != 2 {
		return ""
	}

	if splits[0] != "Bearer" {
		return ""
	}

	return splits[1]
}

func Slug(s string) string {
	return slug.Make(s)
}

func SplitString(s string) []string {
	tags := []string{}
	if s != "" {
		tags = strings.Split(s, ",")
	}

	return tags
}

func TotalPages(perPage, totalItems int) int {
	return int(math.Ceil(float64(totalItems) / float64(perPage)))
}
