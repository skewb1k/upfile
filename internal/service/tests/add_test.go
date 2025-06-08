package service_test

import (
	"errors"
	"testing"

	"github.com/skewb1k/upfile/internal/index"
	"github.com/skewb1k/upfile/internal/service"
	"github.com/skewb1k/upfile/internal/userfile"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errAny = errors.New("ANY")

func TestAdd(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name        string
		path        string
		setupMocks  func(indexProvider *index.MockIndexProvider, userfileProvider *userfile.MockUserFileProvider)
		expectedErr error
	}

	cases := []testCase{
		{
			name: "successfully adds and sets upstream",
			path: "some/dir/file.txt",
			setupMocks: func(idx *index.MockIndexProvider, usr *userfile.MockUserFileProvider) {
				usr.EXPECT().ReadFile(t.Context(), "some/dir/file.txt").Return("file content", nil)
				idx.EXPECT().CreateEntry(t.Context(), "file.txt", "some/dir").Return(nil)
				idx.EXPECT().CheckUpstream(t.Context(), "file.txt").Return(false, nil)
				idx.EXPECT().SetUpstream(t.Context(), "file.txt", index.NewUpstream("file content")).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "entry already exists",
			path: "some/dir/file.txt",
			setupMocks: func(idx *index.MockIndexProvider, usr *userfile.MockUserFileProvider) {
				usr.EXPECT().ReadFile(t.Context(), "some/dir/file.txt").Return("file content", nil)
				idx.EXPECT().CreateEntry(t.Context(), "file.txt", "some/dir").Return(index.ErrExists)
			},
			expectedErr: service.ErrAlreadyTracked,
		},
		{
			name: "check entry fails",
			path: "some/dir/file.txt",
			setupMocks: func(idx *index.MockIndexProvider, usr *userfile.MockUserFileProvider) {
				usr.EXPECT().ReadFile(t.Context(), "some/dir/file.txt").Return("file content", nil)
				idx.EXPECT().CreateEntry(t.Context(), "file.txt", "some/dir").Return(errors.New("no permissions"))
			},
			expectedErr: errAny,
		},
		{
			name: "input file does not exists",
			path: "some/dir/file.txt",
			setupMocks: func(idx *index.MockIndexProvider, usr *userfile.MockUserFileProvider) {
				usr.EXPECT().ReadFile(t.Context(), "some/dir/file.txt").Return("", userfile.ErrNotFound)
			},
			expectedErr: service.ErrFileNotFound,
		},
		{
			name: "upstream already exists, no read or write",
			path: "some/dir/file.txt",
			setupMocks: func(idx *index.MockIndexProvider, usr *userfile.MockUserFileProvider) {
				usr.EXPECT().ReadFile(t.Context(), "some/dir/file.txt").Return("file content", nil)
				idx.EXPECT().CreateEntry(t.Context(), "file.txt", "some/dir").Return(nil)
				idx.EXPECT().CheckUpstream(t.Context(), "file.txt").Return(true, nil)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			idx := index.NewMockIndexProvider(ctrl)
			usr := userfile.NewMockUserFileProvider(ctrl)
			tc.setupMocks(idx, usr)

			s := service.New(idx, usr)
			err := s.Add(t.Context(), tc.path)

			if tc.expectedErr != nil {
				if errors.Is(tc.expectedErr, errAny) {
					require.Error(t, err)
				} else {
					require.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
