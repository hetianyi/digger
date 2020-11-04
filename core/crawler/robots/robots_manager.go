package robots

import (
	"digger/models"
	"digger/utils"
	"fmt"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/timer"
	"github.com/temoto/robotstxt"
	"net/url"
	"sync"
	"time"
)

const (
	// max age of a robots
	robotsAge = time.Minute * 1
)

var (
	robotsCache = make(map[string]*robots)
	lock        = new(sync.Mutex)
)

type robots struct {
	Domain     string
	RobotsTXT  *robotstxt.RobotsData
	ExpireTime time.Time
}

func init() {
	scheduleExpire()
}

func scheduleExpire() {
	timer.Start(0, time.Second*30, 0, func(t *timer.Timer) {
		lock.Lock()
		defer lock.Unlock()
		for k, v := range robotsCache {
			if v.ExpireTime.Before(time.Now()) {
				logger.Warn("expire robots of: ", k)
				delete(robotsCache, k)
			}
		}
	})
}

func TestAgent(link string, project *models.Project) bool {
	u, err := url.Parse(link)
	if err != nil {
		logger.Warn("无法解析链接: ", link, ": ", err.Error())
		return false
	}

	domain := u.Host

	lock.Lock()
	robot := robotsCache[domain]
	lock.Unlock()
	agent := project.Headers["User-Agent"]
	if agent == "" {
		agent = project.Headers["user-agent"]
	}
	if agent == "" {
		agent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"
	}
	if robot != nil {
		return robot.RobotsTXT.TestAgent(u.Path, agent)
	}

	robotsTxtUrl := fmt.Sprintf("%s://%s%s/robots.txt", u.Scheme, u.Host, gox.TValue(u.Port() == "", "", ":"+u.Port()).(string))

	bytes, err := utils.LoadRobotsTxt(robotsTxtUrl, project)
	if err != nil {
		logger.Warn("无法加载Robots文件: ", link, ": ", err.Error())
		return false
	}
	robotsTxt, err := robotstxt.FromBytes(bytes)
	if err != nil {
		logger.Warn("无法解析Robots文件: ", link, ": ", err.Error())
		return false
	}
	robot = &robots{
		Domain:     domain,
		RobotsTXT:  robotsTxt,
		ExpireTime: time.Now().Add(robotsAge),
	}
	robotsCache[domain] = robot
	return robot.RobotsTXT.TestAgent(u.Path, agent)
}
