package dao_test

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestStatus_FindByID_Delete(t *testing.T) {
	// set up
	dao := NewDao(t)
	defer dao.InitAll()
	ctx := context.Background()
	CreateBaseTable(t, dao)

	if s, err := dao.Status().FindByID(ctx, 1); err != nil {
		t.Fatalf("FindByID %v", err)
	} else if s.Content != testStatusContent1 {
		t.Errorf("FindByID returned: %v", s)
	}

	if err := dao.Status().Delete(ctx, 1, 1); err != nil {
		t.Fatalf("Delete %v", err)
	}

	if s, err := dao.Status().FindByID(ctx, 1); err != nil {
		t.Fatalf("FindByID %v", err)
	} else if s != nil {
		t.Errorf("FindByID returned: %v", s)
	}
}

//func TestStatusAllStatuses(t *testing.T) {
//	// set up
//	dao := NewDao(t)
//	defer dao.InitAll()
//	ctx := context.Background()
//	CreateBaseTable(t, dao)

//	type args struct {
//		ctx      context.Context
//		since_id int64
//		max_id   int64
//		limit    int64
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    []object.Status
//		wantErr bool
//	}{
//		{
//			name: "simple",
//			args: args{
//				ctx:      ctx,
//				since_id: 0,
//				max_id:   10,
//				limit:    10,
//			},
//			want: []object.Status{
//				{
//					Content: "",
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		tt := tt
//		t.Run(tt.name, func(t *testing.T) {
//			t.Parallel()
//			r := &status{
//				db: tt.fields.db,
//			}
//			got, err := r.AllStatuses(tt.args.ctx, tt.args.since_id, tt.args.max_id, tt.args.limit)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("status.AllStatuses() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("status.AllStatuses() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
