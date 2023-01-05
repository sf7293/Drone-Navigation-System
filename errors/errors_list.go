package errors

var (
	ErrInternal = Cause{
		Severity:        SeverityCritical,
		Exposable:       true,
		TitleLabel:      "error_title.internal",
		MessageLabel:    "error_message.internal",
	}

	ErrUnexpected = Cause{
		Severity:        SeverityCritical,
		Exposable:       false,
		TitleLabel:      "error_title.unexpected",
		MessageLabel:    "error_message.unexpected",
	}

	ErrRequestBinding = Cause{
		Severity:        SeverityWarning,
		Exposable:       true,
		TitleLabel:      "error_title.request_binding",
		MessageLabel:    "error_message.request_binding",
	}
)