package layer_fetcher

import (
	"errors"
	"fmt"
	"io"
)

type QuotaedReader struct {
	DelegateReader            io.Reader
	QuotaLeft                 int64
	SkipValidation            bool
	QuotaExceededErrorHandler func() error
}

func NewQuotaedReader(delegateReader io.Reader, quotaLeft int64, skipValidation bool, errorMsg string) *QuotaedReader {
	return &QuotaedReader{
		DelegateReader: delegateReader,
		QuotaLeft:      quotaLeft,
		SkipValidation: skipValidation,
		QuotaExceededErrorHandler: func() error {
			return fmt.Errorf(errorMsg)
		},
	}
}

func (q *QuotaedReader) Read(p []byte) (int, error) {
	if q.QuotaLeft < 0 || q.SkipValidation {
		return q.DelegateReader.Read(p)
	}

	if int64(len(p)) > q.QuotaLeft {
		p = p[0 : q.QuotaLeft+1]
	}

	n, err := q.DelegateReader.Read(p)
	q.QuotaLeft = q.QuotaLeft - int64(n)

	if q.QuotaLeft < 0 {
		return n, q.QuotaExceededErrorHandler()
	}

	return n, err
}

func (q *QuotaedReader) AnyQuotaLeft() bool {
	return q.QuotaLeft > 0 && !q.SkipValidation
}

func (q *QuotaedReader) Close() error {
	return errors.New("should not be called")
}
