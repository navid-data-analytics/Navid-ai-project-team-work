"""
All constants needed for pipeline to work.

The smallest unit used is second, every time period used
within the pipeline should be represented as seconds.
With this one convention there is no confusion while
handling time values.
"""
WEEKLY = 7
SECOND = 1
MINUTE = 60 * SECOND
HOUR = 60 * MINUTE
DAY = 24 * HOUR
WEEK = 7 * DAY
MONTH = 30 * DAY
YEAR = 365 * DAY
CONVERT_KILO = 1e3
CONVERT_PERCENT = 1e2
CONVERT_GIGA = 1e9
CONVERT_NANO = 1e-9
