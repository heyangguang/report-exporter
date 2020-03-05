package collector

import (
	"bufio"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"time"
)

func fileLineCount(fileName string) int {
	file, err := os.Open(fileName)
	if err != nil {
		return 0
	}
	defer file.Close()
	fd := bufio.NewReader(file)
	count := 0
	for {
		_, err := fd.ReadString('\n')
		if err != nil {
			break
		}
		count++
	}
	return count
}

type checkFileLineCollector struct {
	dReportFileMetric *prometheus.Desc
}

func NewCheckFileCollector() *checkFileLineCollector {
	v := []string{"filename"}
	return &checkFileLineCollector{
		dReportFileMetric: prometheus.NewDesc("dReportFileMetric", "检测DCI报告数据是否准确数据采集", v, nil),
	}
}

func (collect *checkFileLineCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collect.dReportFileMetric
}

func (collect *checkFileLineCollector) Collect(ch chan<- prometheus.Metric) {

	currentTime := time.Now()
	subOneHour, _ := time.ParseDuration("-24h")
	currentTimeFront := currentTime.Add(subOneHour).Format("2006.01.02")
	currentTimeAfter := time.Now().Format("15:04")
	fileName := "BandwidthInternetDCI-" + currentTimeFront + ".csv"

	// 匹配时间每天01:10 凌晨1点10分检测
	// 30S采样必能匹配时间
	if "01:50" == currentTimeAfter {
		// 获取文件条数
		fileCount := fileLineCount("/tmp/" + fileName)
		// 匹配条数不正确告警
		if fileCount != 34561 {
			ch <- prometheus.MustNewConstMetric(collect.dReportFileMetric, prometheus.GaugeValue, 0, fileName)
			// 告警发给Chanel然后goto到END结束
			goto END
		}
	}

	ch <- prometheus.MustNewConstMetric(collect.dReportFileMetric, prometheus.GaugeValue, 1, fileName)
END:
}
