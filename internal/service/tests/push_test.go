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

func TestPush(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		path        string
		setupMocks  func(indexProvider *index.MockIndexProvider, userfileProvider *userfile.MockUserFileProvider)
		expectedErr error
	}

	cases := []testCase{
		{
			name: "successfully push tracked file",
			path: "/home/user/file.txt",
			setupMocks: func(idx *index.MockIndexProvider, usr *userfile.MockUserFileProvider) {
				idx.EXPECT().CheckEntry(t.Context(), "file.txt", "/home/user").Return(true, nil)
				usr.EXPECT().ReadFile(t.Context(), "/home/user/file.txt").Return("content2", nil)
				idx.EXPECT().GetUpstream(t.Context(), "file.txt").Return("content1", nil)
				idx.EXPECT().SetUpstream(t.Context(), "file.txt", "content2").Return(nil)
			},
			expectedErr: nil,
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
			err := s.Push(t.Context(), tc.path)

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
