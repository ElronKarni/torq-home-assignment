package ip2country

// MongoDBService implements Service by reading data from a MongoDB database
type MongoDBService struct {
	uri string
}

func NewMongoDBService(uri string) (*MongoDBService, error) {
	// ... MongoDBService implementation ...
	return &MongoDBService{uri: uri}, nil
}

func (s *MongoDBService) LookupIP(ip string) (*Result, error) {
	// ... MongoDBService lookup implementation ...
	return nil, nil
}
