package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound              = errors.New("record not found")
	ErrBannerNotFound        = errors.New("banner not found")
	ErrBlogCategoryNotFound  = errors.New("blog category not found")
	ErrBlogNotFound          = errors.New("blog article not found")
	ErrGalleryNotFound       = errors.New("gallery not found")
	ErrGalleryItemNotFound   = errors.New("gallery item not found")
	ErrCustomBlockNotFound   = errors.New("custom content block not found")
	ErrFAQNotFound           = errors.New("faq not found")
	ErrTestimonialNotFound   = errors.New("testimonial not found")
	ErrDuplicateSlug         = errors.New("slug already exists in this application")
	ErrAppIDRequired         = errors.New("app_id is required")
	ErrUnauthorized          = errors.New("unauthorized: authentication required")
	ErrInvalidInput          = errors.New("invalid input data")
)

type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
