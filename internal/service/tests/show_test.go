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

func TestShow(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		fname       string
		setupMocks  func(indexProvider *index.MockIndexProvider, userfileProvider *userfile.MockUserFileProvider)
		expectedRes string
		expectedErr error
	}

	cases := []testCase{
		{
			name:  "successfully get upstream",
			fname: "file.txt",
			setupMocks: func(idx *index.MockIndexProvider, usr *userfile.MockUserFileProvider) {
				idx.EXPECT().GetUpstream(t.Context(), "file.txt").Return("content", nil)
			},
			expectedRes: "content",
			expectedErr: nil,
		},
		{
			name:  "file not tracked",
			fname: "file.txt",
			setupMocks: func(idx *index.MockIndexProvider, usr *userfile.MockUserFileProvider) {
				idx.EXPECT().GetUpstream(t.Context(), "file.txt").Return("", index.ErrNotFound)
			},
			expectedRes: "",
			expectedErr: service.ErrNotTracked,
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
			res, err := s.Show(t.Context(), tc.fname)

			if tc.expectedErr != nil {
				if errors.Is(tc.expectedErr, errAny) {
					require.Error(t, err)
				} else {
					require.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedRes, res)
			}
		})
	}
}
