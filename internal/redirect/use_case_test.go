package redirect

import (
	"context"
	"errors"
	"testing"
)

type SelectLinkRepoMock struct {
	Case string
}

func (mock SelectLinkRepoMock) SelectLink(ctx context.Context, alias string) (string, error) {
	switch mock.Case {
	case "normal":
		return "test.com", nil
	case "not_found":
		return "", ErrAliasNotFound
	case "cache_can_not_set":
		return "test.com", ErrCantSetToCache
	}

	return "", nil
}

func TestUseCase_GetLink(t *testing.T) {
	type fields struct {
		repo SelectLinkRepo
	}
	type args struct {
		ctx   context.Context
		alias string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Normal case",
			fields: fields{
				repo: SelectLinkRepoMock{Case: "normal"},
			},
			args: args{
				ctx:   context.Background(),
				alias: "test",
			},
			want:    "test.com",
			wantErr: false,
		},
		{
			name: "Not found case",
			fields: fields{
				repo: SelectLinkRepoMock{Case: "not_found"},
			},
			args: args{
				ctx:   context.Background(),
				alias: "test",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := UseCase{
				repo: tt.fields.repo,
			}
			got, err := uc.GetLink(tt.args.ctx, tt.args.alias)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_GetLink_Cache_Cases(t *testing.T) {
	type fields struct {
		repo SelectLinkRepo
	}
	type args struct {
		ctx   context.Context
		alias string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          string
		wantErr       bool
		wantAsErrType error
	}{
		{
			name: "Cache can not set case",
			fields: fields{
				repo: SelectLinkRepoMock{Case: "cache_can_not_set"},
			},
			args: args{
				ctx:   context.Background(),
				alias: "test",
			},
			want:          "test.com",
			wantErr:       true,
			wantAsErrType: ErrCantSetToCache,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := UseCase{
				repo: tt.fields.repo,
			}
			got, err := uc.GetLink(tt.args.ctx, tt.args.alias)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !errors.Is(err, tt.wantAsErrType) {
				t.Errorf("GetLink() error = %v, wantAsErrType %v", err, tt.wantAsErrType)
				return
			}
			if got != tt.want {
				t.Errorf("GetLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}
