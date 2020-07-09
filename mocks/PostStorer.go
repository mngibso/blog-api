// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/mngibso/blog-api/models"
)

// PostStorer is an autogenerated mock type for the PostStorer type
type PostStorer struct {
	mock.Mock
}

// DeleteMany provides a mock function with given fields: ctx, username
func (_m *PostStorer) DeleteMany(ctx context.Context, username string) error {
	ret := _m.Called(ctx, username)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: ctx, username
func (_m *PostStorer) Find(ctx context.Context, username string) ([]models.Post, error) {
	ret := _m.Called(ctx, username)

	var r0 []models.Post
	if rf, ok := ret.Get(0).(func(context.Context, string) []models.Post); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Post)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOne provides a mock function with given fields: ctx, postID
func (_m *PostStorer) FindOne(ctx context.Context, postID string) (models.Post, error) {
	ret := _m.Called(ctx, postID)

	var r0 models.Post
	if rf, ok := ret.Get(0).(func(context.Context, string) models.Post); ok {
		r0 = rf(ctx, postID)
	} else {
		r0 = ret.Get(0).(models.Post)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, postID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneAndDelete provides a mock function with given fields: ctx, postID
func (_m *PostStorer) FindOneAndDelete(ctx context.Context, postID string) error {
	ret := _m.Called(ctx, postID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, postID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindOneAndReplace provides a mock function with given fields: ctx, post
func (_m *PostStorer) FindOneAndReplace(ctx context.Context, post models.Post) error {
	ret := _m.Called(ctx, post)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Post) error); ok {
		r0 = rf(ctx, post)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertOne provides a mock function with given fields: ctx, user
func (_m *PostStorer) InsertOne(ctx context.Context, user models.CreatePostInput) (interface{}, error) {
	ret := _m.Called(ctx, user)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, models.CreatePostInput) interface{}); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.CreatePostInput) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}