package main

import (
	"flag"
	"log"
	"strings"

	"github.com/gocolly/colly"

	"github.com/azusa0127/course-watcher-go/ubc"
)

func main() {
	deptP := flag.String("dept", "CPSC", "department of the course, e.g. BIOL")
	termP := flag.String("term", "2", "term 1 or 2")
	dayP := flag.String("day", "Mon", "day of the week, e.g. Mon")
	startTimeP := flag.String("time", "13:00", "lecture starting time in HH:mm, e.g 13:00")
	upperOnlyP := flag.Bool("upper", false, "Filter out all courses in lower level.")
	flag.Parse()
	println("Dept:", *deptP, "Term:", *termP, "Day:", *dayP, "StartTime:", *startTimeP, "UpperLevelOnly:", *upperOnlyP)
	println("------------")
	filteredSects := filterCoursesBySchedule(*deptP, *termP, *dayP, *startTimeP, *upperOnlyP)
	for _, sect := range filteredSects {
		println(sect.Section)
	}
	println("------------")
}

func filterCoursesBySchedule(dept, term, day, startTime string, upperOnly bool) []*ubc.CourseSectMini {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: coursera.org, www.coursera.org
		colly.AllowedDomains("ubc.ca", "courses.students.ubc.ca"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./ubc_cache"),
	)

	// Create another collector to scrape course details
	sects := make([]*ubc.CourseSectMini, 0, 200)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// If attribute class is this long string return from callback
		// As this a is irrelevant
		link := e.Attr("href")
		if !strings.HasPrefix(link, "/cs/courseschedule?pname=subjarea&tname=subj-course&") {
			return
		}
		cq, err := ubc.NewCourseQueryFromHref(link)
		if err != nil {
			log.Fatalln("CourseQuery Error:", err, " with link: ", link)
		}

		// log.Println("Visiting", cq.BuildURL().String())

		// start scaping the page under the link found
		cm, err := cq.GetCourseMini()
		if err != nil {
			log.Fatalln("GetCourseMini Error:", err)
		}

		if csm := cm.GetLectureSectOn(term, day, startTime, upperOnly); csm != nil {
			sects = append(sects, csm)
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		// log.Println("visiting", r.URL.String())
	})

	c.Visit("https://courses.students.ubc.ca/cs/courseschedule?pname=subjarea&tname=subj-department&dept=" + dept)

	return sects
}
