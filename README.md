UBC Scheduler
---

A small tool for finding courses that fits in your time schedule.

Usage
---
```bash
go run ubc.go -day Mon -term 2 -dept ASIA -time 16:00 -upper
```

Options:

`-day` Day of the week (Mon, Tue, Wed, Thu, Fri). Default: `Mon`

`-term` Term `1` or `2`. Default: `2`

`-dept` Department code in UPPER case. Default: `CPSC`

`-time` Start time of the class in format `HH:mm`. Default: `13:00`

`-upper` Show only courses in upper level. Default: `false`
