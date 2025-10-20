package timex

import (
	"errors"
	"iter"
	"time"
)

// ErrInvalidTimeRange 表示无效的时间范围错误
var ErrInvalidTimeRange = errors.New("invalid time range")

// InclusiveTimeRange 表示一个包含起始和结束时间的时间范围
type InclusiveTimeRange struct {
	start time.Time
	end   time.Time
}

// MustNewInclusiveTimeRange 创建InclusiveTimeRange, 如果参数无效则 panic
func MustNewInclusiveTimeRange(startTime, endTime time.Time, startTimeInclusive, endTimeInclusive bool) *TimeRange {
	t, err := NewTimeRange(startTime, endTime, startTimeInclusive, endTimeInclusive)
	if err != nil {
		panic(err)
	}
	return t
}

// NewInclusiveTimeRange 创建InclusiveTimeRange
func NewInclusiveTimeRange(startTime, endTime time.Time) (*InclusiveTimeRange, error) {
	if startTime.After(endTime) {
		return nil, ErrInvalidTimeRange
	}
	return &InclusiveTimeRange{
		start: startTime,
		end:   endTime,
	}, nil
}

// StartTime 返回开始时间
func (tr *InclusiveTimeRange) StartTime() time.Time {
	return tr.start
}

// EndTime 返回结束时间
func (tr *InclusiveTimeRange) EndTime() time.Time {
	return tr.end
}

// IterTimeBy 按照指定时间间隔迭代时间范围内的时间点
func (tr *InclusiveTimeRange) IterTimeBy(interval time.Duration) iter.Seq[time.Time] {
	return func(yield func(time.Time) bool) {
		start, end := tr.StartTime(), tr.EndTime()
		t := start
		for t.Before(end) || t.Equal(end) {
			if !yield(t) {
				return
			}
			t = t.Add(interval)
		}
	}
}

// IsBeforeStart 这个方法判断给定时间是否在开始时间之前
func (tr *InclusiveTimeRange) IsBeforeStart(t time.Time) bool {
	return t.Before(tr.start)
}

// IsAfterEnd 这个方法判断给定时间是否在结束时间之后
func (tr *InclusiveTimeRange) IsAfterEnd(t time.Time) bool {
	return t.After(tr.end)
}

// Contains 判断时间是否在时间区间内
func (tr *InclusiveTimeRange) Contains(t time.Time) bool {
	return (t.After(tr.start) || t.Equal(tr.start)) && (t.Before(tr.end) || t.Equal(tr.end))
}

// TimeRange 表示一个更通用时间范围类型,可以指定起始和结束时间是否包含在范围内
type TimeRange struct {
	start          time.Time
	end            time.Time
	startInclusive bool // start 是否落在日期范围中, 即是否左闭区间.
	endInclusive   bool // end 是否落在日期范围中, 即是否右闭区间.
}

// MustNewTimeRange 创建TimeRange, 如果参数无效则 panic
func MustNewTimeRange(startTime, endTime time.Time, startTimeInclusive, endTimeInclusive bool) *TimeRange {
	t, err := NewTimeRange(startTime, endTime, startTimeInclusive, endTimeInclusive)
	if err != nil {
		panic(err)
	}
	return t
}

// NewTimeRange 创建TimeRange
func NewTimeRange(startTime, endTime time.Time, startTimeInclusive, endTimeInclusive bool) (*TimeRange, error) {
	st, et := startTime, endTime
	if startTimeInclusive {
		st = st.Add(time.Nanosecond)
	}
	if endTimeInclusive {
		et = et.Add(-time.Nanosecond)
	}
	if st.After(et) {
		return nil, ErrInvalidTimeRange
	}

	return &TimeRange{
		start:          startTime,
		end:            endTime,
		startInclusive: startTimeInclusive,
		endInclusive:   endTimeInclusive,
	}, nil
}

// StartTime 返回开始时间
func (tr *TimeRange) StartTime() time.Time {
	return tr.start
}

// StartTimeInclusive 返回包含在范围内的开始时间
func (tr *TimeRange) StartTimeInclusive() time.Time {
	if tr.startInclusive {
		return tr.start
	}

	return tr.start.Add(time.Nanosecond)
}

// EndTime 返回结束时间
func (tr *TimeRange) EndTime() time.Time {
	return tr.end
}

// EndTimeInclusive 返回包含在范围内的结束时间
func (tr *TimeRange) EndTimeInclusive() time.Time {
	if tr.startInclusive {
		return tr.end
	}

	return tr.end.Add(-time.Nanosecond)
}

// IsStartTimeInclusive 返回开始时间是否包含在范围内
func (tr *TimeRange) IsStartTimeInclusive() bool {
	return tr.startInclusive
}

// IsEndTimeInclusive 返回结束时间是否包含在范围内
func (tr *TimeRange) IsEndTimeInclusive() bool {
	return tr.endInclusive
}

// IsBeforeStart 这个方法判断给定时间是否在开始时间之前, 如果开始时间是包含的, 则等于开始时间也算在外部
func (tr *TimeRange) IsBeforeStart(t time.Time) bool {
	if tr.startInclusive {
		return t.Before(tr.start)
	}
	return t.Before(tr.start) || t.Equal(tr.start)
}

// IsAfterEnd 这个方法判断给定时间是否在结束时间之后, 如果结束时间是包含的, 则等于结束时间也算在外部
func (tr *TimeRange) IsAfterEnd(t time.Time) bool {
	if tr.endInclusive {
		return t.After(tr.end)
	}
	return t.After(tr.end) || t.Equal(tr.end)
}

// Contains 判断时间是否在时间区间内
func (tr *TimeRange) Contains(t time.Time) bool {
	s := t.After(tr.start)
	if tr.startInclusive {
		s = s || t.Equal(tr.start)
	}

	e := t.Before(tr.end)
	if tr.endInclusive {
		e = e || t.Equal(tr.end)
	}

	return s && e
}

// ToInclusiveTimeRange 转换为 InclusiveTimeRange
func (tr *TimeRange) ToInclusiveTimeRange() (*InclusiveTimeRange, error) {
	st := tr.start
	if !tr.startInclusive {
		st = st.Add(time.Nanosecond)
	}

	et := tr.end
	if !tr.endInclusive {
		et = et.Add(-time.Nanosecond)
	}

	return NewInclusiveTimeRange(st, et)
}
