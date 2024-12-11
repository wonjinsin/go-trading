package service

import (
	"testing"

	"go.uber.org/mock/gomock"
)

func TestService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
}
