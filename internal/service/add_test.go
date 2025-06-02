package service_test

import (
	"errors"
	"testing"

	"upfile/internal/index"
	"upfile/internal/service"
	"upfile/internal/userfile"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errAny = errors.New("ANY")

func TestAdd(t *testing.T) {
	type testCase struct {
		name          string
		path          string
		setupMocks    func(indexStore *index.MockStore, userfileStore *userfile.MockStore)
		expectedError error
	}

	cases := []testCase{
		{
			name: "successfully adds and sets upstream",
			path: "some/dir/file.txt",
			setupMocks: func(idx *index.MockStore, usr *userfile.MockStore) {
				idx.EXPECT().CheckEntry(t.Context(), "file.txt", "some/dir").Return(false, nil)
				idx.EXPECT().CreateEntry(t.Context(), "file.txt", "some/dir").Return(nil)
				idx.EXPECT().CheckUpstream(t.Context(), "file.txt").Return(false, nil)
				usr.EXPECT().ReadFile(t.Context(), "some/dir/file.txt").Return("file content", nil)
				idx.EXPECT().SetUpstream(t.Context(), "file.txt", "file content").Return(nil)
			},
		},
		{
			name: "entry already exists",
			path: "some/dir/file.txt",
			setupMocks: func(idx *index.MockStore, usr *userfile.MockStore) {
				idx.EXPECT().CheckEntry(t.Context(), "file.txt", "some/dir").Return(true, nil)
			},
			expectedError: service.ErrAlreadyTracked,
		},
		{
			name: "check entry fails",
			path: "some/dir/file.txt",
			setupMocks: func(idx *index.MockStore, usr *userfile.MockStore) {
				idx.EXPECT().CheckEntry(t.Context(), "file.txt", "some/dir").Return(false, errors.New("no permissions"))
			},
			expectedError: errAny,
		},
		{
			name: "upstream already exists, no read or write",
			path: "some/dir/file.txt",
			setupMocks: func(idx *index.MockStore, usr *userfile.MockStore) {
				idx.EXPECT().CheckEntry(t.Context(), "file.txt", "some/dir").Return(false, nil)
				idx.EXPECT().CreateEntry(t.Context(), "file.txt", "some/dir").Return(nil)
				idx.EXPECT().CheckUpstream(t.Context(), "file.txt").Return(true, nil)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			idx := index.NewMockStore(ctrl)
			usr := userfile.NewMockStore(ctrl)
			tc.setupMocks(idx, usr)

			s := service.New(idx, usr)
			err := s.Add(t.Context(), tc.path)

			if tc.expectedError != nil {
				if errors.Is(tc.expectedError, errAny) {
					require.Error(t, err)
				} else {
					require.ErrorIs(t, err, tc.expectedError)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
