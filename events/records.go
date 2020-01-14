package events

import (
	"context"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/pkg/storage"
)

const (
	AppendRecordEvent       EventType = 1
	UpdateRecordUpEvent     EventType = 2
	UpdateRecordDownEvent   EventType = 3
	UpdateRecordStatusEvent EventType = 4
)

func init() {
	// append record
	RegisterEvent(AppendRecordEvent, func(ctx context.Context, v interface{}) error {
		r, ok := v.(*RecordEntity)
		if !ok {
			return errors.Errorf("[%s] is not RecordEntity", reflect.TypeOf(v).Kind().String())
		}
		AppendRecord(ctx, r)
		return nil
	})
	// update record up
	RegisterEvent(UpdateRecordUpEvent, func(ctx context.Context, v interface{}) error {
		r, ok := v.(*RecordEntity)
		if !ok {
			return errors.Errorf("[%s] is not RecordEntity", reflect.TypeOf(v).Kind().String())
		}
		UpdateRecord(ctx, r.ID, func(re *RecordEntity) {
			re.Up += r.Up
		})
		return nil
	})
	// update record down
	RegisterEvent(UpdateRecordDownEvent, func(ctx context.Context, v interface{}) error {
		r, ok := v.(*RecordEntity)
		if !ok {
			return errors.Errorf("[%s] is not RecordEntity", reflect.TypeOf(v).Kind().String())
		}
		UpdateRecord(ctx, r.ID, func(re *RecordEntity) {
			re.Down += r.Down
		})
		return nil
	})
	// update record status
	RegisterEvent(UpdateRecordStatusEvent, func(ctx context.Context, v interface{}) error {
		r, ok := v.(*RecordEntity)
		if !ok {
			return errors.Errorf("[%s] is not RecordEntity", reflect.TypeOf(v).Kind().String())
		}
		UpdateRecord(ctx, r.ID, func(re *RecordEntity) {
			re.Status = r.Status
		})
		return nil
	})
}

type ConnEntity struct {
	ID         int64
	SourceAddr string
	DestAddr   string
}

type RecordStatus struct {
	Text string
}

func (r *RecordStatus) String() string {
	return r.Text
}

var (
	ActiveStatus    = &RecordStatus{"Active"}
	CompletedStatus = &RecordStatus{"Completed"}
)

type RecordEntity struct {
	ID        int64
	DestAddr  string
	Policy    string
	Up        int64
	Down      int64
	Status    *RecordStatus
	Timestamp time.Time
	Protocol  string
	Duration  time.Duration
	Conn      ConnEntity
}

var recordStarge = storage.NewLRUList(1000)

func AppendRecord(ctx context.Context, record *RecordEntity) {
	logrus.WithField("recode", *record).Debug("events: append record")
	recordStarge.PushBack(record)
}
func RangeRecord(ctx context.Context, f func(*RecordEntity) bool) {
	recordStarge.Range(func(v interface{}) bool {
		return f(v.(*RecordEntity))
	})
}
func UpdateRecord(ctx context.Context, id int64, f func(*RecordEntity)) {
	recordStarge.Range(func(v interface{}) bool {
		if r, ok := v.(*RecordEntity); ok && r.ID == id {
			f(r)
			return true
		}
		return false
	})
}
