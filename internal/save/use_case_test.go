package save

import (
	"context"
	"errors"
	"testing"
)

type MockRepo struct {
	Case string
}

func (m MockRepo) InsertLink(ctx context.Context, link, alias string) error {
	switch m.Case {
	case "normal":
		return nil
	case "alias_exists":
		return ErrAliasAlreadyExists
	}

	return nil
}

func TestUseCase_AddLink(t *testing.T) {
	type fields struct {
		repo InsertLinkRepo
	}
	type args struct {
		ctx   context.Context
		link  string
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
				repo: MockRepo{Case: "normal"},
			},
			args: args{
				ctx:   context.Background(),
				link:  "test.com",
				alias: "test",
			},
			want:    "test",
			wantErr: false,
		},
		{
			name: "Alias already exists",
			fields: fields{
				repo: MockRepo{Case: "alias_exists"},
			},
			args: args{
				ctx:   context.Background(),
				link:  "test.com",
				alias: "test",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := UseCase{
				repo: tt.fields.repo,
			}
			got, err := ls.AddLink(tt.args.ctx, tt.args.link, tt.args.alias)
			if (err != nil) != tt.wantErr {
				if t.Name() == "Auto alias aleady exists" {
					if !errors.Is(err, ErrAliasAlreadyExists) {
						t.Errorf("AddLink() error = %v, wantErr %v", err, ErrAliasAlreadyExists)
						return
					}
				}

				t.Errorf("AddLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				if t.Name() == "Empty alias" {
					return
				}
				t.Errorf("AddLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_AddLink_Empty_Aliases(t *testing.T) {
	type fields struct {
		repo InsertLinkRepo
	}
	type args struct {
		ctx   context.Context
		link  string
		alias string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantErrType error
	}{
		{
			name: "Empty alias",
			fields: fields{
				repo: MockRepo{Case: "normal"},
			},
			args: args{
				ctx:   context.Background(),
				link:  "test.com",
				alias: "",
			},
			wantErr: false,
		},
		{
			name: "Auto alias aleady exists",
			fields: fields{
				repo: MockRepo{Case: "normal"},
			},
			args: args{
				ctx:   context.Background(),
				link:  "test.com",
				alias: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := UseCase{
				repo: tt.fields.repo,
			}
			_, err := ls.AddLink(tt.args.ctx, tt.args.link, tt.args.alias)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && errors.Is(err, tt.wantErrType) {
				t.Errorf("AddLink() error = %v, wantErr %v", err, tt.wantErrType)
				return
			}
		})
	}
}
