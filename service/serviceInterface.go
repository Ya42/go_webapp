package service

type ServiceInterface interface{
	NewService(path string) *ServiceInterface
  Dispose()
}
