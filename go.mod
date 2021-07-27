module github.com/jasonjoo2010/goschedule-console

go 1.14

require (
	github.com/gin-gonic/gin v1.7.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/jasonjoo2010/goschedule v1.1.0
	github.com/jasonjoo2010/goschedule/store/database v1.1.0
	github.com/jasonjoo2010/goschedule/store/etcdv2 v1.1.0
	github.com/jasonjoo2010/goschedule/store/etcdv3 v1.1.0
	github.com/jasonjoo2010/goschedule/store/redis v1.1.0
	github.com/jasonjoo2010/goschedule/store/zookeeper v1.1.0
	github.com/robfig/cron v1.2.0
	github.com/stretchr/testify v1.7.0 // test
	gopkg.in/yaml.v3 v3.0.0-20200506231410-2ff61e1afc86
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
