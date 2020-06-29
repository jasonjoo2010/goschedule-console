module github.com/jasonjoo2010/goschedule-console

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/jasonjoo2010/goschedule v0.1.2
	github.com/jasonjoo2010/goschedule/store/redis v0.1.2
	github.com/jasonjoo2010/goschedule/store/zookeeper v0.1.2
	github.com/jasonjoo2010/goschedule/store/database v0.1.2
	github.com/jasonjoo2010/goschedule/store/etcdv2 v0.1.2
	github.com/jasonjoo2010/goschedule/store/etcdv3 v0.1.2
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.5.1 // test
	gopkg.in/yaml.v3 v3.0.0-20200506231410-2ff61e1afc86
	github.com/go-sql-driver/mysql v1.5.0
)
