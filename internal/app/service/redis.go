package service

func (s *Service) RedisSetKeyValue(key string, value int) error {
	err := s.rdb.SetKeyValue(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RedisGetValue(key string) (string, error) {
	val, err := s.rdb.GetValue(key)
	if err != nil {
		return "", err
	}
	return val, nil
}
