package route_tests

import "pancast-server/config"

func GetServerRoute(route string, conf config.StartParameters) string {
	return "https://" + config.GetServerURL(conf) + route
}