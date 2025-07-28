package router

import (
	"fmt"
	"strings"
)

// PathBuilder URL 경로 구성을 위한 유틸리티
type PathBuilder struct {
	version string
	service string
}

// NewPathBuilder 새로운 PathBuilder 인스턴스 생성
func NewPathBuilder(version, service string) *PathBuilder {
	return &PathBuilder{
		version: version,
		service: service,
	}
}

// BasePath 기본 경로 반환 (/{version}/{service})
func (pb *PathBuilder) BasePath() string {
	return fmt.Sprintf("/%s/%s", pb.version, pb.service)
}

// Path 가변 인자로 경로 세그먼트들을 받아 완전한 경로 반환
// 예: Path("users", "123", "orders") -> /{version}/{service}/users/123/orders
func (pb *PathBuilder) Path(segments ...string) string {
	if len(segments) == 0 {
		return pb.BasePath()
	}

	path := pb.BasePath()
	for _, segment := range segments {
		if segment != "" {
			path += "/" + strings.TrimPrefix(segment, "/")
		}
	}

	return path
}

// Version 현재 버전 반환
func (pb *PathBuilder) Version() string {
	return pb.version
}

// Service 현재 서비스명 반환
func (pb *PathBuilder) Service() string {
	return pb.service
}
