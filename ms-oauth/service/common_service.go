package service

// Service 服务接口，需实现健康检查接口
type Service interface {
	HealthCheck() bool
}

// CommentService 实现类
type CommentService struct {
}

// HealthCheck 健康检查，返回true
func (s *CommentService) HealthCheck() bool {
	return true
}

// NewCommentService 构建对象
func NewCommentService() *CommentService {
	return &CommentService{}
}
